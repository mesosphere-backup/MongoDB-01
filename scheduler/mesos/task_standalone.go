package mesos

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"strings"
	//"encoding/json"
	//"github.com/tangmingdong123/mongodb-mesos/scheduler/config"
	"github.com/tangmingdong123/mongodb-mesos/scheduler/repo"
	//"strconv"
)

func handleStandalone(driver sched.SchedulerDriver,
	offers []*mesos.Offer,
	idleIDs []*mesos.OfferID,
	usedIDs []*mesos.OfferID,
	usedMap map[*mesos.Offer]*Used) {
	for _, db := range repo.ListStandalone() {
		if !db.Cancel { //TO BE DEPLOYING
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
					u.Cpu = u.Cpu + float64(db.Cpu)
					u.Mem = u.Mem + float64(db.Memory)
					hostPort := selectPort(offer, u)

					if hostPort != 0 {
						log.Infof("toLaunchTask,%v", *db)
						db.State = repo.STATE_DEPLOYING
						db.Port = hostPort
						db.Hostname = *offer.Hostname
						repo.SaveStandalone(db)

						driver.LaunchTasks([]*mesos.OfferID{offer.GetId()},
							[]*mesos.TaskInfo{genStandaloneTask(db, offer, hostPort)},
							&mesos.Filters{RefuseSeconds: proto.Float64(5)})
					} else {
						log.Errorf("no useful port")
					}
				}
			}
		} else if db.Cancel { //to be cancel
			if db.State == repo.STATE_DEPLOYING || db.State == repo.STATE_RUNNING {
				taskID := &mesos.TaskID{
					Value: proto.String(PREFIX_TASK_STANDALONE + db.Name),
				}
				driver.KillTask(taskID)
			} else {
				db.State = repo.STATE_CANCEL
			}

			if db.State == repo.STATE_CANCEL && db.Cancel {
				repo.SaveStandalone(db)
			}
		}
	}
}

func genStandaloneTask(db *repo.DBNode, offer *mesos.Offer, hostPort uint64) *mesos.TaskInfo {
	taskID := &mesos.TaskID{
		Value: proto.String(PREFIX_TASK_STANDALONE + db.Name),
	}
	taskType := mesos.ContainerInfo_DOCKER

	containerPort := uint32(27017)
	protocol := "tcp"
	network := mesos.ContainerInfo_DockerInfo_BRIDGE
	hostPort32 := uint32(hostPort)
	//visibility := mesos.DiscoveryInfo_EXTERNAL
	//domainname := PREFIX_TASK_STANDALONE + db.Name + "." + *config.SchedulerName

	task := &mesos.TaskInfo{
		Name:    proto.String(PREFIX_TASK_STANDALONE + db.Name),
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
			Arguments: []string{},
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
		//Discovery: &mesos.DiscoveryInfo{
		//	Visibility: &visibility,
		//	Name:       &domainname,
		//},
	}
	return task
}

func isMatchStandalone(db *repo.DBNode, offers []*mesos.Offer, usedMap map[*mesos.Offer]*Used) *mesos.Offer {
	for _, offer := range offers {
		summary := sum(offer)
		merge(summary, usedMap[offer])

		if float64(db.Cpu) <= summary.Cpu && float64(db.Memory) <= summary.Mem {
			return offer
		}
	}

	return nil
}

func updateStandaloneStatus(status *mesos.TaskStatus) {
	name := strings.Replace(status.GetTaskId().GetValue(), PREFIX_TASK_STANDALONE, "", -1)
	db := repo.FindStandalone(name)

	if db != nil {
		//bs, _ := repo.DBNodeJson(db)
		//log.Infof("db status update before,%v\n", string(bs))

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

		repo.SaveStandalone(db)
		//log.Infof("db status update after,%v\n", string(bs))
	}
}
