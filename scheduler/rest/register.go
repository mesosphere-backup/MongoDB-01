package rest

import(
	"github.com/emicklei/go-restful"
	
)

func Register(container *restful.Container, cors bool) {
	router := Router{}
	router.registerAll(container, cors)
}

type Router struct {
}
type StandaloneService struct{
	
}
type ReplicaSetService struct{
	
}
type ShardClusterService struct{
	
}

func (r *Router)registerAll(container *restful.Container, cors bool){
	standalone := &StandaloneService{}
	rs := &ReplicaSetService{}
	shard := &ShardClusterService{}
	
	
	wss := []*restful.WebService{}
	wss = append(wss,
		standalone.StandaloneService(),rs.ReplicaSetService(),shard.ShardClusterService())

	for _, ws := range wss {
		// Cross Origin Resource Sharing filter
		if cors {
			corsRule := restful.CrossOriginResourceSharing{ExposeHeaders: []string{"Content-Type"}, CookiesAllowed: false, Container: container}
			ws.Filter(corsRule.Filter)
		}
		// Add webservice to container
		container.Add(ws)
	}
}