package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	//"os"
	//"path/filepath"

	//"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	//"github.com/magiconair/properties"
	"github.com/tangmingdong123/mongodb-mesos/scheduler/rest"
	"github.com/tangmingdong123/mongodb-mesos/scheduler/repo"
	mongodbschd "github.com/tangmingdong123/mongodb-mesos/scheduler/mesos"
	config "github.com/tangmingdong123/mongodb-mesos/scheduler/config"
)

/**
mongodb-mesos -mesos=172.17.2.91:5050 -zk 192.168.3.223:2181 -name mongodb-mesos -port 37017 -failoverTimeoutSeconds 300
*/
func main() {
	//parse args
	config.MesosMasterIpAndPort = flag.String("mesos","172.17.2.91:5050","mesos master ip and port")
	config.ZK = flag.String("zk","192.168.3.223:2181","repo of mongodb-scheduler")
	config.SchedulerName = flag.String("name","mongodb-mesos","framework's name")	
	config.HTTPPort = flag.Int("port",37017,"framework's http port")
	config.FailoverTimeoutSeconds = flag.Float64("failoverTimeoutSeconds",300,"Failover Timeout in Second")
	
	flag.Parse()
	
	fmt.Println("mongodb-mesos scheduler start...")
	fmt.Printf("mongodb-mesos scheduler mesos:%s,zk:%s,name:%s,port:%d\n",*config.MesosMasterIpAndPort,*config.ZK,*config.SchedulerName,*config.HTTPPort)
	
	//launch HTTP REST service
	go launchHTTP(*config.HTTPPort)
	
	//int zk
	repo.InitZK([]string{*config.ZK},"/"+*config.SchedulerName)
	
	//launch framework
	mongodbschd.Start()
	select{}
}

func launchHTTP(port int){
	fmt.Printf("mongodb-mesos framework listen on %d\n",port)
	
	// accept and respond in JSON unless told otherwise
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultContainer.EnableContentEncoding(true)
	restful.DefaultContainer.Router(restful.CurlyRouter{})
	restful.SetCacheReadEntity(false)
	rest.Register(restful.DefaultContainer, false)
	
	config := swagger.Config{
			WebServices:     restful.DefaultContainer.RegisteredWebServices(),
			WebServicesUrl:  "",
			ApiPath:         "/apidocs.json",
			SwaggerPath:     "/apidocs/",
			SwaggerFilePath: "swagger-ui/dist",
		}
	
	swagger.RegisterSwaggerService(config,restful.DefaultContainer)

	// If swagger is not on `/` redirect to it
	http.HandleFunc("/", index)
		
	server := &http.Server{Addr: ":"+strconv.Itoa(port), Handler: restful.DefaultContainer}
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/apidocs/", http.StatusMovedPermanently)
}