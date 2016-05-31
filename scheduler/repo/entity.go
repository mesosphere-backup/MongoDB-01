package repo



//a mongodb node
type DBNode struct {
	Name string  `json:"name"`
	Cpu    float64 `json:"cpu"`
	Memory float64   `json:"memory"`
	State string `json:"state"`
	Cancel bool `json:"cancel"`
	Port uint64 `json:"port"`
	Hostname string `json:"hostname"`
}


//a ReplicaSet ,with name and nodes
type ReplicaSet struct {
	Name string   `json:"name"`
	Nodes  []*DBNode `json:"nodes"`
	State string `json:"state"`
	InitState string `json:"initState"`
	Cancel bool `json:"cancel"`
}

type RouterNode struct {
	Name string  `json:"name"`
	Cpu    float64 `json:"cpu"`
	Memory float64   `json:"memory"`
	State string `json:"state"`
	Cancel bool `json:"cancel"`
	Port uint32 `json:"port"`
	Hostname string `json:"hostname"`
}

//a Shard cluster, with name,configReplicaSet,routers,shards
type ShardCluster struct {
	Name     string       `json:"name"`
	Routers  []*RouterNode `json:"routers"`
	ConfigRS ReplicaSet   `json:"configRS"`
	Shards   []*ReplicaSet `json:"shards"`
	State string `json:"state"`
}

type Meta struct {
	StandaloneMap   map[string]*DBNode
	ReplicaSetMap  map[string]*ReplicaSet
	ShardClusterMap map[string]*ShardCluster
}
