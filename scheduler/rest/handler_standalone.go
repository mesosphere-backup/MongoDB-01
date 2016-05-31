package rest

import (
	"github.com/emicklei/go-restful"
	"github.com/tangmingdong123/mongodb-mesos/scheduler/repo"
	"strconv"
	log "github.com/Sirupsen/logrus"
)

func (r *StandaloneService) CreateStandalone(req *restful.Request, resp *restful.Response){
	resp.AddHeader("Content-Type","application/json")
	
	
	id := req.PathParameter("id")
	cpu,_ := strconv.ParseFloat(req.QueryParameter("cpu"),32)
	mem,_ := strconv.Atoi(req.QueryParameter("mem"))
	
	
	if(repo.IsStandaloneExist(id)){
		log.Errorf("standalone mongodb '%s' already exist.",id)
		r := &Response{Code:CODE_ALREADY_EXIST,Desc:"already exist"}
		bs,_ := r.Byte()
		resp.Write(bs)
		return
	}else{
		db := &repo.DBNode{Name:id,
			Cpu:cpu,
			Memory:float64(mem),
			State:repo.STATE_INIT}
		
		repo.AddStandalone(db)
		
		bs,_ := repo.DBNodeJson(db)
		resp.Write(bs)
	}
}

func (r *StandaloneService) ListStandalone(req *restful.Request, resp *restful.Response){
	resp.AddHeader("Content-Type","application/json")
	
	bs,_ := repo.StandaloneListJson()
	resp.Write(bs)
}

func (r *StandaloneService) KillStandalone(req *restful.Request, resp *restful.Response){
	resp.AddHeader("Content-Type","application/json")
	name := req.PathParameter("id")
	
	db := repo.FindStandalone(name)
	if(db != nil){
		db.Cancel = true
		
		repo.SaveStandalone(db)
		bs,_ := repo.DBNodeJson(db)

		resp.Write(bs)
	}else{
		r := &Response{Code:CODE_NOT_EXIST,Desc:"not exist"}
		bs,_ := r.Byte()
		
		resp.Write(bs)
	}
}

func (r *StandaloneService) GetStandalone(req *restful.Request, resp *restful.Response){
	resp.AddHeader("Content-Type","application/json")
	name := req.PathParameter("id")
	
	db := repo.FindStandalone(name)
	if(db != nil){
		bs,_ := repo.DBNodeJson(db)

		resp.Write(bs)
	}else{
		r := &Response{Code:CODE_NOT_EXIST,Desc:"not exist"}
		bs,_ := r.Byte()
		
		resp.Write(bs)
	}
}