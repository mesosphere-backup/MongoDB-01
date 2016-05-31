package rest

import (
	//"encoding/json"
	//"errors"

	"github.com/emicklei/go-restful"
	//"github.com/Sirupsen/logrus"
	//"github.com/compose/mejson"
	//"gopkg.in/mgo.v2/bson"
)




func (r *StandaloneService)StandaloneService()*restful.WebService{
	ws := new(restful.WebService)
	ws.Path("/standalone")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	
	ws.Route(ws.POST("/{id}").
		To(r.CreateStandalone).
		Doc("createStandalone").
		Operation("createStandalone").
		Param(ws.PathParameter("id","Standalone Mongodb id")).
		Param(ws.FormParameter("cpu","mongodb's cpu").DataType("float")).
		Param(ws.FormParameter("mem","mongodb's mem in MB").DataType("int32")).
		Produces("application/json"))
		
	ws.Route(ws.GET("/{id}").To(r.GetStandalone).
		// docs
		Doc("getStandalone Info").
		Operation("getStandalone").
		Param(ws.PathParameter("id", "identifier of the standalone").DataType("string")).
		Consumes("text/plain").Produces("application/json"))
		
	ws.Route(ws.DELETE("/{id}").To(r.KillStandalone).
		// docs
		Doc("killStandalone Info").
		Operation("killStandalone").
		Param(ws.PathParameter("id", "identifier of the standalone").DataType("string")).
		Consumes("text/plain").Produces("application/json"))
	
	ws.Route(ws.GET("/list").To(r.ListStandalone).
		// docs
		Doc("getStandalone list").
		Operation("listStandalone").
		Consumes("application/json").Produces("application/json"))
	
	return ws
}


func (r *ReplicaSetService)ReplicaSetService()*restful.WebService{
	ws := new(restful.WebService)
	ws.Path("/rs")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	
	ws.Route(ws.POST("/{id}").
		To(r.CreateReplicaSet).
		Doc("createReplicaSet").
		Operation("createReplicaSet").
		Param(ws.PathParameter("id","ReplicaSet Mongodb id")).
		Param(ws.FormParameter("cpu","mongodb's cpu").DataType("float")).
		Param(ws.FormParameter("mem","mongodb's mem in MB").DataType("int32")).
		Param(ws.FormParameter("instances","mongodb's num of instances").DataType("int32")).
		Produces("application/json"))
	
	ws.Route(ws.GET("/{id}").To(r.GetReplicaSet).
		// docs
		Doc("getReplicaSet Info").
		Operation("getReplicaSet").
		Param(ws.PathParameter("id", "identifier of the ReplicaSet").DataType("string")).
		Consumes("text/plain").Produces("application/json"))
		
	ws.Route(ws.DELETE("/{id}").To(r.KillReplicaSet).
		// docs
		Doc("killReplicaSet Info").
		Operation("killReplicaSet").
		Param(ws.PathParameter("id", "identifier of the ReplicaSet").DataType("string")).
		Consumes("text/plain").Produces("application/json"))
		
	ws.Route(ws.GET("/list").To(r.ListReplicaSets).
		// docs
		Doc("getReplicaSet list").
		Operation("listReplicaSet").
		Consumes("application/json").Produces("application/json"))
	
	return ws
}


func (r *ShardClusterService)ShardClusterService()*restful.WebService{
	ws := new(restful.WebService)
	ws.Path("/shard")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	
	return ws
}