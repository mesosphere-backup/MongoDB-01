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
)

/**
mongodb-mesos -mesos=172.17.2.91:5050 -zk 192.168.3.223:2181 -name mongodb-mesos 
*/
func main() {
	//parse args
	mesos := flag.String("mesos","172.17.2.91:5050","zk of mesos")
	zk := flag.String("zk","192.168.3.223:2181","repo of mongodb-scheduler")
	name := flag.String("name","mongodb-mesos","framework's name")	
	port := flag.Int("port",37017,"framework's http port")
	
	flag.Parse()
	
	fmt.Println("mongodb-mesos scheduler start...")
	fmt.Printf("mongodb-mesos scheduler mesos:%s,zk:%s,name:%s,port:%d\n",*mesos,*zk,*name,*port)
	
	//launch HTTP REST service
	go launchHTTP(*port)
	
	//int zk
	repo.InitZK([]string{*zk},"/"+*name)
	
	//launch framework
	mongodbschd.Start(mesos)
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
			SwaggerFilePath: "d:/swagger-ui/dist",
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