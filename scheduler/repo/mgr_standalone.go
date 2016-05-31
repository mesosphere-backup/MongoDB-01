package repo

import(
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/samuel/go-zookeeper/zk"
)


//-----------------for standalone
func IsStandaloneExist(name string)bool{
	_,ok := meta.StandaloneMap[name]
	return ok
}

func AddStandalone(db *DBNode){
	meta.StandaloneMap[db.Name] = db
	
	SaveStandalone(db)
}
func StandaloneListJson()([]byte,error){
	return json.Marshal(&meta.StandaloneMap)
}
func DBNodeJson(db *DBNode)([]byte,error){
	return json.Marshal(db)
}
func FindStandalone(name string)(*DBNode){
	return meta.StandaloneMap[name]
}
func ListStandalone()[]*DBNode{
	arr := make([]*DBNode,len(meta.StandaloneMap))
	
	i := 0
	for name := range meta.StandaloneMap {
		arr[i] = meta.StandaloneMap[name]
		i = i + 1
	}
	return arr
}
func SaveStandalone(node *DBNode){
	path := rootPath+"/standalone/"+node.Name
	log.Infof("saveStandalone %s",path)
		
	if node.Cancel && node.State == STATE_CANCEL{
		//delete from zk and memory
		log.Infof("standalone is Canceled and closed,so delete it")
		conn.Delete(path,-1)
		delete(meta.StandaloneMap,node.Name)
	}else {
		
		
		bytes,_ := json.Marshal(&node)
		
		ex, _, err := conn.Exists(path)
		if err != nil {
			log.Infof("exist %s err:%s", path, err)
			return
		}
		
		if(ex){
			_,err := conn.Set(path,bytes,-1)
			if(err!=nil){
				log.Infof("saveStandalone fail %s",err)
			}
		}else{
			_,err := conn.Create(path,bytes,0,zk.WorldACL(zk.PermAll))
			if(err!=nil){
				log.Infof("saveStandalone fail %s",err)
			}
		}
	}
}