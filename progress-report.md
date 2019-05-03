# KubeDB Work Progress Report

## May, 2019

### Completed Tasks

- [[MongoDB] MongoDB sharding](https://github.com/kubedb/project/issues/234)
  
  Finally, MongoDB Sharding is here in KubeDB! it will be a part of kubedb 0.12.0 release. Documentation is added. Try it out and feel free to give your feedback. SSL support is not added yet. Hopefully, it will be a part of 0.12.0. 
  
- [[MySQL] MySQL clustering - Group replication without InnoDB](https://github.com/kubedb/project/issues/18)
    
  MySQL Group replication is added in 0.12.0. Documentation is also added. More MySQL clustering feature will be added in the future. Feel free to try and give your feedback. 
  
- [[Postgres] WAL-archiving to Minio broken](https://github.com/kubedb/project/issues/492)
  
  Postgres WAL-archiving to Minio is supported now. Both SSL or non-SSL will work. Documentation is also added for minio.
  
- [[Postgres] wal-g: Postgres Continuous Archiving to Swift failing](https://github.com/kubedb/project/issues/486)
  
  The issue is fixed in master. It will be a part of 0.12.0 release.
  
- [[Postgres] Support Local volume for PostgreSQL WAL archiving](https://github.com/kubedb/project/issues/475)
  
  Postgres WAL-archiving to Local volume is supported now. Documentation is also added.
  
- [[Postgres] PostgreSQL 11.2](https://github.com/kubedb/project/issues/483)
  
  Postgres 11.2 support is added.

- [[MongoDB] Use official Mongo GO Driver for testing MongoDB](https://github.com/kubedb/project/issues/491)
  
  Mongo has released their official [GO driver](https://github.com/mongodb/mongo-go-driver) recently. We have updated our driver to mongodb official go driver in mongodb E2E testing.

- [[CLI] kubedb - MethodNotAllowed](https://github.com/kubedb/project/issues/467)
  
  Users sometimes got this error while running `kubedb get dormantdatabases`, because there were two resources whose plural name was `dormantdatabases`. These two resources were `dormantdatabases.kubedb.com` & `dormantdatabases.validators.kubedb.com`. The workaround is given in [issue comment](https://github.com/kubedb/project/issues/467#issuecomment-475384569). It is fixed in the master branch, so, it will be a part of 0.12.0 release.


- [Assign resources for init containers](https://github.com/kubedb/project/issues/503) 
  
  Resources for init containers was not applied in the past. This issue is fixed in the master branch.

- [Snapshot to minio with SSL](https://github.com/kubedb/project/issues/457)
  
  Snapshot to minio with SSL was failing. It is fixed in the master branch. This will be a part of 0.12.0 release.

  
- [Kubedb 0.11 missing rbac permissions error](https://github.com/kubedb/project/issues/481) 
  
  There was an issue while upgrading from 0.10.0 to 0.11.0. The workaround is given in [issue comment](https://github.com/kubedb/project/issues/481#issuecomment-481171356). The issue is solved in the master branch, so, it will not be an issue while upgrading to 0.12.0.
  
### Ongoing Tasks

- [Support restic as snapshot uploader](https://github.com/kubedb/project/issues/168)
  
  This has been a long term goal to integrate [stash](https://github.com/stashed/stash) with KubeDB. Stash has introduced new API `v1beta1` which supports [`AppBinding`](https://blog.byte.builders/post/the-case-for-appbinding/) and `Functions-Tasks` (Don't bother. The terms are new but interesting. You will know once the stash new release 0.9.0-stable/alpha/beta happens). 
  
  Hopefully, we will be able to introduce stash integration with the kubedb-operator in 0.13.0 release. We are now working on it.
  
- [[MongoDB] SSL support in MongoDB](https://github.com/kubedb/project/issues/352)

  SSL support is not added to MongoDB. We are working on it. Hopefully, we will be able to add some parts (if not all) of these tasks mentioned in the issue.
  
- [[MySQL] MySQL Clustering - Percona XtraDB]()
  
  We are trying to add MySQL clustering for Percona XtraDB. Currently, we are researching how to do it. Hopefully, it will be a part of kubedb 0.13.0 release.
