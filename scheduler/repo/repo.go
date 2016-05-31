package repo

import (
	"encoding/json"
	"github.com/samuel/go-zookeeper/zk"
	"time"
	log "github.com/Sirupsen/logrus"
)

var meta *Meta = &Meta{StandaloneMap: make(map[string]*DBNode),
	ReplicaSetMap:  make(map[string]*ReplicaSet),
	ShardClusterMap: make(map[string]*ShardCluster)}
var conn *zk.Conn
var rootPath string

func InitZK(zkHosts []string, root string) (zkconn *zk.Conn, err error) {
	log.Infof("connect zk %v, root:%v", zkHosts,root)
	zkconn, _, err = zk.Connect(zkHosts, time.Second)

	if err != nil {
		log.Infof("connect zk err:%s", err)
	} else {
		conn = zkconn
		rootPath = root

		load(root)
	}
	return
}

func load(root string) {
	//make sure the rootpath exist
	createIfNotExist(root, []byte("root"))
	createIfNotExist(root+"/standalone", []byte("standalone"))
	createIfNotExist(root+"/replica", []byte("replicaSet"))
	createIfNotExist(root+"/shard", []byte("shardCluster"))

	loadStandalone()
	loadReplicatSet()
	loadShardCluster()
	
	log.Infof("load meta finish :%s",meta)
}

func createIfNotExist(path string, data []byte) {
	log.Infof("createIfNotExist %s", path)
	ex, _, err := conn.Exists(path)
	if err != nil {
		log.Infof("exist %s err:%s", path, err)
		return
	}
	if !ex {
		path, err := conn.Create(path,
			data,
			0,
			zk.WorldACL(zk.PermAll))

		if err != nil {
			log.Infof("create %s,err:%s", path, err)
			return
		}
	}
}

func loadStandalone() {
	standalonePath := rootPath + "/standalone"

	childs, _, err := conn.Children(standalonePath)
	if err != nil {
		log.Infof("fetch standalonePath's children fail,%s", err)
		return
	}

	for i, child := range childs {
		log.Infof("standalonePath child %d = %s", i, child)

		bytes, _, err := conn.Get(standalonePath + "/" + child)
		if err != nil {
			log.Infof("fetch standalone fail %s", err)
		} else {
			var dbNode DBNode
			err := json.Unmarshal(bytes, &dbNode)
			if err != nil {

			} else {
				meta.StandaloneMap[dbNode.Name] = &dbNode
			}
		}
	}
}
func loadReplicatSet() {
	replicaPath := rootPath + "/replica"

	childs, _, err := conn.Children(replicaPath)
	if err != nil {
		log.Infof("fetch replicaPath's children fail,%s", err)
		return
	}

	for i, child := range childs {
		log.Infof("replicaPath child %d = %s", i, child)

		bytes, _, err := conn.Get(replicaPath + "/" + child)
		if err != nil {
			log.Infof("fetch replica fail %s", err)
		} else {
			var rs ReplicaSet
			err := json.Unmarshal(bytes, &rs)
			if err != nil {

			} else {
				meta.ReplicaSetMap[rs.Name] = &rs
			}
		}
	}
}
func loadShardCluster() {

}


