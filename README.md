# MongoDB Replica Set as Framework

A mesos MongoDB scheduler for MongoDB framework. The scheduler provides a REST API to manipulate your standalone MongoDB and replicaSets with Create,Delete and Query.

For standalone scenarios, The scheduler monitor the MongoDB's status, restart it when the task is killed or failed.

For replicaSet scenarios, The scheduler will autoconfig the cluster, and monitor the MongoDB's status, restart it when the task is killed or failed.

You can get the MongoDB instances detailed information by REST API: http://api.mongodb.com/

# How to start

For example : /scheduler -mesos 172.17.2.91:5050 -zk 172.17.2.91:2181 -name mongodb-mesos -port 37017

Usage: ./scheduler -master $mesos-master-ip:port -zk zk-ip:port -name schedulername -port httpport


# Persistence
All standalone MongoDB and replicaSets' detailed information are saved in the zookeeper. The scheduler will reload these information when its restart. The zk' path is /${your scheduler name},and it is /mongodb-mesos by default.

# Rest API Usage
To create a standalone MongoDB:

curl -X DELETE --header 'Accept: application/json' --header 'Content-Type: application/x-www-form-urlencoded' -d 'cpu=1&mem=128' 'http://172.17.2.254:37017/standalone/1'

To delete a standalone MongoDB:

curl -X DELETE --header 'Accept: application/json' 'http://172.17.2.254:37017/standalone/1'

To list all standalone MongoDBs:

curl -X GET --header 'Accept: application/json' 'http://172.17.2.254:37017/standalone/list'

To get a standalone MongoDB's information:

curl -X GET --header 'Accept: application/json' 'http://172.17.2.254:37017/standalone/1'

To create a replicaSet MongoDB cluster:

curl -X DELETE --header 'Accept: application/json' --header 'Content-Type: application/x-www-form-urlencoded' -d 'cpu=1&mem=128&instances=3' 'http://172.17.2.254:37017/replica/r1'

To delete a replicaSet MongoDB cluster:

curl -X DELETE --header 'Accept: application/json' 'http://172.17.2.254:37017/replica/r1'

To list all replicaSet MongoDB clusters:

curl -X GET --header 'Accept: application/json' 'http://172.17.2.254:37017/replica/list'

To get a replicaSet MongoDB cluster's information:

curl -X GET --header 'Accept: application/json' 'http://172.17.2.254:37017/replica/r1'

# Roadmap

