# default port 37017

# swagger support
visit http://172.17.2.254:37017/apidocs ,and input http://172.17.2.254:37017/apidocs.json ,click 'Explorer' button. 

# REST API
# create a standalone mongodb
curl -X DELETE --header 'Accept: application/json' --header 'Content-Type: application/x-www-form-urlencoded' -d 'cpu=1&mem=128' 'http://172.17.2.254:37017/standalone/1'

# delete a standalone mongodb
curl -X DELETE --header 'Accept: application/json' 'http://172.17.2.254:37017/standalone/1'

# list all standalone mongodbs
curl -X GET --header 'Accept: application/json' 'http://172.17.2.254:37017/standalone/list'

# get  a standalone mongodb's info
curl -X GET --header 'Accept: application/json' 'http://172.17.2.254:37017/standalone/1'

# create a replicaSet mongodb cluster
curl -X DELETE --header 'Accept: application/json' --header 'Content-Type: application/x-www-form-urlencoded' -d 'cpu=1&mem=128&instances=3' 'http://172.17.2.254:37017/replica/r1'

# delete a replicaSet mongodb cluster
curl -X DELETE --header 'Accept: application/json' 'http://172.17.2.254:37017/replica/r1'

# list all replicaSet mongodb clusters
curl -X GET --header 'Accept: application/json' 'http://172.17.2.254:37017/replica/list'

# get a replicaSet mongodb cluster's info
curl -X GET --header 'Accept: application/json' 'http://172.17.2.254:37017/replica/r1'
