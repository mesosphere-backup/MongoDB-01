package mesos

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"strconv"
	"strings"
	//"encoding/json"
	//"github.com/tangmingdong123/mongodb-mesos/scheduler/mongocmd"
	"github.com/tangmingdong123/mongodb-mesos/scheduler/repo"
	//"strconv"
)

func handleReplicaSet(driver sched.SchedulerDriver,
	offers []*mesos.Offer,
	idleIDs []*mesos.OfferID,
	usedIDs []*mesos.OfferID,
	usedMap map[*mesos.Offer]*Used) {
	for _, rs := range repo.ListReplicaSet() {
		if !rs.Cancel { //TO BE DEPLOYING
			for _, db := range rs.Nodes {
				if db.State == repo.STATE_FAIL {
					db.State = repo.STATE_INIT
				}

				if db.State == repo.STATE_INIT {
					offer := isMatchStandalone(db, offers, usedMap)
					if offer != nil {
						usedIDs = append(usedIDs, offer.GetId())

						u := usedMap[offer]
						if u == nil {
							u = &Used{Cpu: 0, Mem: 0, Ports: []uint64{}}
							usedMap[offer] = u
						}
						u.Cpu = u.Cpu + db.Cpu
						u.Mem = u.Mem + db.Memory
						hostPort := selectPort(offer, u)

						log.Infof("toLaunchTask,%v", *rs)
						db.State = repo.STATE_DEPLOYING
						db.Port = hostPort
						db.Hostname = *offer.Hostname
						rs.State = repo.STATE_DEPLOYING
						repo.SaveReplicaSet(rs)

						driver.LaunchTasks([]*mesos.OfferID{offer.GetId()},
							[]*mesos.TaskInfo{genReplicaTask(rs, db, offer, hostPort)},
							&mesos.Filters{RefuseSeconds: proto.Float64(5)})
					}
				}
			}
			//if init fail before, then retry
			if rs.State == repo.STATE_RUNNING && (rs.InitState == repo.INIT_STATE_INIT || rs.InitState == repo.INIT_STATE_FAIL) {
				db := &repo.DBNode{Cpu: PREFIX_TASK_REPLICA_INIT_CPU, Memory: PREFIX_TASK_REPLICA_INIT_MEM}

				offer := isMatchStandalone(db, offers, usedMap)
				if offer != nil {
					usedIDs = append(usedIDs, offer.GetId())

					u := usedMap[offer]
					if u == nil {
						u = &Used{Cpu: 0, Mem: 0, Ports: []uint64{}}
						usedMap[offer] = u
					}
					u.Cpu = u.Cpu + db.Cpu
					u.Mem = u.Mem + db.Memory

					log.Infof("toLaunchReplicaInitTask,%v", *rs)
					rs.InitState = repo.INIT_STATE_DEPLOYING

					driver.LaunchTasks([]*mesos.OfferID{offer.GetId()},
						[]*mesos.TaskInfo{genReplicaInitTask(rs, offer)},
						&mesos.Filters{RefuseSeconds: proto.Float64(5)})
				}
			}
		} else if rs.Cancel { //to be cancel
			for _, db := range rs.Nodes {
				if db.State == repo.STATE_DEPLOYING || db.State == repo.STATE_RUNNING {

					taskID := &mesos.TaskID{
						Value: proto.String(replicaTaskID(rs, db)),
					}
					driver.KillTask(taskID)
				} else {
					db.State = repo.STATE_CANCEL
				}
			}

			checkReplicaState(rs)
			if rs.State == repo.STATE_CANCEL && rs.Cancel {
				repo.SaveReplicaSet(rs)
			}
		}
	}
}

//taskID:  _replica_rsname_0
func genReplicaTask(rs *repo.ReplicaSet, db *repo.DBNode, offer *mesos.Offer, hostPort uint64) *mesos.TaskInfo {
	taskID := &mesos.TaskID{
		Value: proto.String(replicaTaskID(rs, db)),
	}
	taskType := mesos.ContainerInfo_DOCKER

	containerPort := uint32(27017)
	protocol := "tcp"
	network := mesos.ContainerInfo_DockerInfo_BRIDGE
	hostPort32 := uint32(hostPort)

	task := &mesos.TaskInfo{
		Name:    proto.String(replicaTaskID(rs, db)),
		TaskId:  taskID,
		SlaveId: offer.SlaveId,
		Container: &mesos.ContainerInfo{
			Type: &taskType,
			Docker: &mesos.ContainerInfo_DockerInfo{
				Image: proto.String("mongo:3.2.6"),
				PortMappings: []*mesos.ContainerInfo_DockerInfo_PortMapping{
					&mesos.ContainerInfo_DockerInfo_PortMapping{
						HostPort:      &hostPort32,
						ContainerPort: &containerPort,
						Protocol:      &protocol}},
				Network: &network,
			},
		},
		Command: &mesos.CommandInfo{
			Shell:     proto.Bool(false),
			Arguments: []string{"--replSet", rs.Name},
		},
		Resources: []*mesos.Resource{
			util.NewScalarResource("cpus", float64(db.Cpu)),
			util.NewScalarResource("mem", float64(db.Memory)),
			util.NewRangesResource("ports", []*mesos.Value_Range{
				&mesos.Value_Range{
					Begin: &hostPort,
					End:   &hostPort,
				}}),
		},
	}
	return task
}

func genReplicaInitTask(rs *repo.ReplicaSet, offer *mesos.Offer) *mesos.TaskInfo {
	taskID := &mesos.TaskID{
		Value: proto.String(replicaInitTaskID(rs)),
	}
	taskType := mesos.ContainerInfo_DOCKER

	var db0 *repo.DBNode //node to be connect

	initial := "rs.initiate({_id:\"" + rs.Name + "\",members:["
	for i, db := range rs.Nodes {
		if db0 == nil {
			db0 = db
		}

		initial = initial + "{_id:" + strconv.Itoa(i) + ",host:\"" + db.Hostname + ":" + strconv.Itoa(int(db.Port)) + "\"}"
		if i != len(rs.Nodes)-1 {
			initial = initial + ","
		}
	}
	initial = initial + "]})"
	//initial = "'"+initial+"'"

	log.Infof("initial command %s", initial)
	node := db0.Hostname + ":" + strconv.Itoa(int(db0.Port))

	task := &mesos.TaskInfo{
		Name:    proto.String(replicaInitTaskID(rs)),
		TaskId:  taskID,
		SlaveId: offer.SlaveId,
		Container: &mesos.ContainerInfo{
			Type: &taskType,
			Docker: &mesos.ContainerInfo_DockerInfo{
				Image: proto.String("mongo:3.2.6"),
			},
		},
		Command: &mesos.CommandInfo{
			Shell:     proto.Bool(false),
			Arguments: []string{"/usr/bin/mongo", node, "--eval=" + initial},
		},
		Resources: []*mesos.Resource{
			util.NewScalarResource("cpus", float64(PREFIX_TASK_REPLICA_INIT_CPU)),
			util.NewScalarResource("mem", float64(PREFIX_TASK_REPLICA_INIT_MEM)),
		},
	}
	return task
}

func replicaTaskID(rs *repo.ReplicaSet, db *repo.DBNode) string {
	return PREFIX_TASK_REPLICA + rs.Name + "_" + db.Name
}

func replicaInitTaskID(rs *repo.ReplicaSet) string {
	return PREFIX_TASK_REPLICA_INIT + rs.Name
}

//_replica_rsname_0
func updateReplicaStatus(status *mesos.TaskStatus) {
	fname := strings.Replace(status.GetTaskId().GetValue(), PREFIX_TASK_REPLICA, "", -1)
	pos := strings.LastIndex(fname, "_")
	rsname := fname[0:pos]
	dbname := fname[pos+1 : len(fname)]

	rs := repo.FindReplicaSet(rsname)

	if rs != nil {
		//bs, _ := repo.ReplicaSetJson(rs)
		//log.Infof("rs status update before,%v\n", string(bs))

		var db *repo.DBNode
		for _, element := range rs.Nodes {
			if element.Name == dbname {
				db = element
				break
			}
		}

		if db != nil {
			if db.Cancel {
				if IsFail(status) {
					db.State = repo.STATE_CANCEL
				}
			} else {
				if IsFail(status) {
					db.State = repo.STATE_FAIL
				} else if IsRunning(status) {
					db.State = repo.STATE_RUNNING
				} else {
					db.State = repo.STATE_DEPLOYING
				}
			}

			checkReplicaState(rs)
			repo.SaveReplicaSet(rs)

			//log.Infof("rs status update after,%v\n", string(bs))
		} else {
			log.Errorf("db %v of replicaSet %s not exist", dbname, rsname)
		}
	} else {
		log.Errorf("replicaSet %s not exist", rsname)
	}
}

func updateReplicaInitStatus(status *mesos.TaskStatus) {
	rsname := strings.Replace(status.GetTaskId().GetValue(), PREFIX_TASK_REPLICA_INIT, "", -1)

	rs := repo.FindReplicaSet(rsname)

	if rs != nil {
		//bs, _ := repo.ReplicaSetJson(rs)
		//log.Infof("rs status update before,%v\n", string(bs))

		if status.GetState() == mesos.TaskState_TASK_FINISHED {
			rs.InitState = repo.INIT_STATE_FINISH
		} else if IsFail(status) {
			rs.InitState = repo.INIT_STATE_FAIL
		} else {
			rs.InitState = repo.INIT_STATE_DEPLOYING
		}

		repo.SaveReplicaSet(rs)
		//log.Infof("rs status update after,%v\n", string(bs))
	} else {
		log.Errorf("replicaSet %s not exist", rsname)
	}
}

//check if all dbs of the rs are canceled
//check if one db of the rs are running
func checkReplicaState(rs *repo.ReplicaSet) {
	if rs.Cancel {
		var allCanceled = true
		for _, db := range rs.Nodes {
			if db.State != repo.STATE_DEPLOYING && db.State != repo.STATE_RUNNING {
				db.State = repo.STATE_CANCEL
			}
			db.Cancel = true
			if db.State != repo.STATE_CANCEL {
				allCanceled = false
				break
			}
		}

		if allCanceled { //if all canceled ,then the rs is canceld
			rs.State = repo.STATE_CANCEL
		}
	} else {
		var allRunning = true
		for _, db := range rs.Nodes {
			if db.State != repo.STATE_RUNNING {
				allRunning = false
				break
			}
		}
		if allRunning { //if one is running ,then the rs is running
			rs.State = repo.STATE_RUNNING
		}
	}
}

func initReplicaSet(rs *repo.ReplicaSet) {

	var db0 *repo.DBNode
	var ipports []string
	for _, db := range rs.Nodes {
		if db0 == nil {
			db0 = db
		}

		ipports = append(ipports, db.Hostname+""+strconv.Itoa(int(db.Port)))
	}

}
