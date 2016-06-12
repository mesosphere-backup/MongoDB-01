# task nameing

# for standalone
  the task name is _standalone_{name}
  
# for replicaSet
  the mongod task name is _replica_{name}_{seq}
  the initial task name is _replicainit_{name}
  
# monitor and failover
  The scheduler receive the mesos master message when task status is changed, if the task is lost , fail ,kill , finish, the scheduler will
  restart the task
  
