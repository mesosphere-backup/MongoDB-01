package repo

import(
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/samuel/go-zookeeper/zk"
)


//-----------------for standalone
func IsReplicaSetExist(name string)bool{
	_,ok := meta.ReplicaSetMap[name]
	return ok
}

func AddReplicaSet(db *ReplicaSet){
	meta.ReplicaSetMap[db.Name] = db
	
	SaveReplicaSet(db)
}
func ReplicaSetListJson()([]byte,error){
	return json.Marshal(&meta.ReplicaSetMap)
}
func ReplicaSetJson(db *ReplicaSet)([]byte,error){
	return json.Marshal(db)
}
func FindReplicaSet(name string)(*ReplicaSet){
	return meta.ReplicaSetMap[name]
}
func ListReplicaSet()[]*ReplicaSet{
	arr := make([]*ReplicaSet,len(meta.ReplicaSetMap))
	
	i := 0
	for name := range meta.ReplicaSetMap {
		arr[i] = meta.ReplicaSetMap[name]
		i = i + 1
	}
	return arr
}


func SaveReplicaSet(node *ReplicaSet){
	path := rootPath+"/replica/"+node.Name
	log.Infof("saveRelicaSet %s",path)
	
	if node.Cancel && node.State == STATE_CANCEL {
		log.Infof("replicaSet is Canceled and all instances is closed,so delete it")
		conn.Delete(path,-1)
		delete(meta.ReplicaSetMap,node.Name)	
	}else{
		bytes,_ := json.Marshal(&node)
		
		ex, _, err := conn.Exists(path)
		if err != nil {
			log.Infof("exist %s err:%s", path, err)
			return
		}
		
		if(ex){
			_,err := conn.Set(path,bytes,-1)
			if(err!=nil){
				log.Infof("saveRelicaSet fail %s",err)
			}
		}else{
			_,err := conn.Create(path,bytes,0,zk.WorldACL(zk.PermAll))
			if(err!=nil){
				log.Infof("saveRelicaSet fail %s",err)
			}
		}
	}
}