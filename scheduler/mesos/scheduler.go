package mesos

import (
	"strconv"
	//"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	//util "github.com/mesos/mesos-go/mesosutil"
	//"time"
	//"encoding/json"
	"github.com/tangmingdong123/mongodb-mesos/scheduler/config"
	"github.com/tangmingdong123/mongodb-mesos/scheduler/repo"
	"strings"
)

type MongodbScheduler struct {
}

func Start() {
	log.Infof("startScheduler master:%v,name:%v", *config.MesosMasterIpAndPort, *config.SchedulerName)

	fwinfo := &mesos.FrameworkInfo{
		User:            proto.String(""),
		Name:            proto.String(*config.SchedulerName),
		Id:              &mesos.FrameworkID{Value: proto.String(*config.SchedulerName)},
		FailoverTimeout: config.FailoverTimeoutSeconds, //in second
		Checkpoint:      proto.Bool(true),
		WebuiUrl:        proto.String("http://" + getDockerbrIP() + ":" + strconv.Itoa(*config.HTTPPort)),
	}

	driverConfig := sched.DriverConfig{
		Scheduler: newMongodbScheduler(),
		Framework: fwinfo,
		Master:    *config.MesosMasterIpAndPort,
	}

	driver, err := sched.NewMesosSchedulerDriver(driverConfig)

	log.Infof("schduler.driver create finish")
	if err != nil {
		log.Errorln("Unable to create a SchedulerDriver ", err.Error())
	}

	stat, err := driver.Run()
	if err != nil {
		log.Infof("Framework stopped with status %s and error: %s", stat.String(), err.Error())
	}

}

func listTasks() []*mesos.TaskStatus {
	var list []*mesos.TaskStatus
	//list standalone
	for _, db := range repo.ListStandalone() {
		if db.State != repo.STATE_DEPLOYING || db.State == repo.STATE_RUNNING {
			state := mesos.TaskState_TASK_RUNNING
			list = append(list, &mesos.TaskStatus{
				TaskId: &mesos.TaskID{Value: proto.String(PREFIX_TASK_STANDALONE + db.Name)},
				State:  &state,
			})
		}
	}

	//list replicaSet
	for _, rs := range repo.ListReplicaSet() {
		for _, db := range rs.Nodes {
			state := mesos.TaskState_TASK_RUNNING
			list = append(list, &mesos.TaskStatus{
				TaskId: &mesos.TaskID{Value: proto.String(PREFIX_TASK_STANDALONE + db.Name)},
				State:  &state,
			})
		}
	}

	//list shards,TODO
	return list
}

func newMongodbScheduler() *MongodbScheduler {
	return &MongodbScheduler{}
}

func (sched *MongodbScheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Infoln("Framework Registered with Master ", masterInfo)

	_, err := driver.ReconcileTasks(listTasks())
	if err != nil {
		log.Infof("ReconcileTasks fail %v", err)
	}

}

func (sched *MongodbScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Infoln("Framework Re-Registered with Master ", masterInfo)

	_, err := driver.ReconcileTasks(listTasks())
	if err != nil {
		log.Infof("ReconcileTasks fail %v", err)
	}
}

func (sched *MongodbScheduler) Disconnected(sched.SchedulerDriver) {
	log.Warningf("disconnected from master")
}

func (sched *MongodbScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	//log.Warningf("Framework resourceOffer")

	/*
		for _, offer := range offers {
			bytes, _ := json.Marshal(offer)
			log.Infof("offer:%s", string(bytes))
		}
	*/

	var idleIDs []*mesos.OfferID
	var usedIDs []*mesos.OfferID
	usedMap := make(map[*mesos.Offer]*Used)

	//handle standalone first
	handleStandalone(driver, offers, idleIDs, usedIDs, usedMap)

	//handle replica second
	handleReplicaSet(driver, offers, idleIDs, usedIDs, usedMap)

	//unused offer
	for _, offer := range offers {
		used := false
		for _, usedid := range usedIDs {
			if offer.GetId() == usedid {
				used = true
				break
			}
		}
		if !used {
			idleIDs = append(idleIDs, offer.GetId())
		}
	}
	//reject offer
	driver.LaunchTasks(idleIDs, make([]*mesos.TaskInfo, 0), &mesos.Filters{RefuseSeconds: proto.Float64(5)})
}

func (sched *MongodbScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infoln("Status update: task", status.TaskId.GetValue(), " is in state ", status.State.Enum().String())
	log.Infof("reason:%v,message:%v,source:%v\n", status.GetReason().Enum(), status.GetMessage(), status.GetSource())

	//bs, _ := json.Marshal(status)
	//log.Infof("Status info %v", string(bs))

	if strings.Contains(status.GetTaskId().GetValue(), PREFIX_TASK_STANDALONE) {
		updateStandaloneStatus(status)
	} else if strings.Contains(status.GetTaskId().GetValue(), PREFIX_TASK_REPLICA) {
		updateReplicaStatus(status)
	} else if strings.Contains(status.GetTaskId().GetValue(), PREFIX_TASK_REPLICA_INIT) {
		updateReplicaInitStatus(status)
	}

}

func (sched *MongodbScheduler) OfferRescinded(_ sched.SchedulerDriver, oid *mesos.OfferID) {
	log.Errorf("offer rescinded: %v", oid)
}
func (sched *MongodbScheduler) FrameworkMessage(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, msg string) {
	log.Errorf("framework message from executor %q slave %q: %q", eid, sid, msg)
}
func (sched *MongodbScheduler) SlaveLost(_ sched.SchedulerDriver, sid *mesos.SlaveID) {
	log.Errorf("slave lost: %v", sid)
}
func (sched *MongodbScheduler) ExecutorLost(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, code int) {
	log.Errorf("executor %q lost on slave %q code %d", eid, sid, code)
}
func (sched *MongodbScheduler) Error(_ sched.SchedulerDriver, err string) {
	log.Errorf("Scheduler received error: %v", err)
}
