package rest

import (
	"github.com/emicklei/go-restful"
	"github.com/tangmingdong123/mongodb-mesos/scheduler/repo"
	"strconv"
	log "github.com/Sirupsen/logrus"
)

func (r *ReplicaSetService) CreateReplicaSet(req *restful.Request, resp *restful.Response){
	resp.AddHeader("Content-Type","application/json")
	
	
	id := req.PathParameter("id")
	cpu,_ := strconv.ParseFloat(req.QueryParameter("cpu"),32)
	mem,_ := strconv.Atoi(req.QueryParameter("mem"))
	instances,_ := strconv.Atoi(req.QueryParameter("instances"))
	
	
	if(repo.IsReplicaSetExist(id)){
		log.Errorf("ReplicaSet '%s' already exist.",id)
		r := &Response{Code:CODE_ALREADY_EXIST,Desc:"already exist"}
		bs,_ := r.Byte()
		resp.Write(bs)
		return
	}else{
		rs := &repo.ReplicaSet{}
		rs.Name = id
		rs.Cancel = false
		rs.State = repo.STATE_INIT
		rs.InitState = repo.INIT_STATE_INIT
		
		nodes := make([]*repo.DBNode,instances)
		for i:=0;i<instances;i++ {
			nodes[i] = &repo.DBNode{
				Name:strconv.Itoa(i),
				Cpu:cpu,
				Memory:float64(mem),
				State:repo.STATE_INIT,
			}
		}
		rs.Nodes = nodes
		
		repo.AddReplicaSet(rs)
		
		bs,_ := repo.ReplicaSetJson(rs)
		resp.Write(bs)
	}
}

func (r *ReplicaSetService) ListReplicaSets(req *restful.Request, resp *restful.Response){
	resp.AddHeader("Content-Type","application/json")
	
	bs,_ := repo.ReplicaSetListJson()
	resp.Write(bs)
}

func (r *ReplicaSetService) KillReplicaSet(req *restful.Request, resp *restful.Response){
	resp.AddHeader("Content-Type","application/json")
	name := req.PathParameter("id")
	
	rs := repo.FindReplicaSet(name)
	if(rs != nil){
		rs.Cancel = true
		
		for _, db := range rs.Nodes {
			db.Cancel = true
			
			if db.State !=repo.STATE_DEPLOYING && db.State != repo.STATE_RUNNING {
				db.State = repo.STATE_CANCEL
			}
		}
		
		repo.SaveReplicaSet(rs)
		bs,_ := repo.ReplicaSetJson(rs)

		resp.Write(bs)
	}else{
		r := &Response{Code:CODE_NOT_EXIST,Desc:"not exist"}
		bs,_ := r.Byte()
		
		resp.Write(bs)
	}
}

func (r *ReplicaSetService) GetReplicaSet(req *restful.Request, resp *restful.Response){
	resp.AddHeader("Content-Type","application/json")
	name := req.PathParameter("id")
	
	rs := repo.FindReplicaSet(name)
	if(rs != nil){
		bs,_ := repo.ReplicaSetJson(rs)

		resp.Write(bs)
	}else{
		r := &Response{Code:CODE_NOT_EXIST,Desc:"not exist"}
		bs,_ := r.Byte()
		
		resp.Write(bs)
	}
}