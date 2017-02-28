 # Pgpool-II
 - https://github.com/paunin/postgres-docker-cluster
 - https://github.com/RasterBurn/docker-pgpool2
 - https://project.altservice.com/issues/697
 - http://www.pgpool.net/docs/latest/tutorial-en.html

# PgpoolAdmin
 - http://pgpool.projects.pgfoundry.org/pgpoolAdmin/doc/en/install.html
 - http://pgpool.projects.pgfoundry.org/pgpoolAdmin/doc/en/errorCode.html
 - superuser: unknown issue: http://www.sraoss.jp/pipermail/pgpool-general/2014-January/002495.html
```
    The solution is that "pcp_user" needs to be a real user for the backend.
	$ psql -p {pgpool port} -U {login_user} -W template1 -c "SELECT usesuper
	FROM pg_user WHERE usename = 'postgres'" Password for user {login user}:
	 usesuper
	----------
	 t
	(1 row)

	If failed, pg_hba.conf might have something incorrect.
```

Mount a file to /srv/pgpool2/secrets/.admin with secrets

```
name: "PCP_USER"
value: "pcp_user"

name: "PCP_PASSWORD"
value: "pcp_pass"

name: "PGPOOL_START_DELAY"
value: "120"

name: "REPLICATION_USER"
value: "replication_user"

name: "REPLICATION_PASSWORD"
value: "replication_pass"

name: "DB_USERS"
value: "monkey_user:monkey_pass"

name: "BACKENDS"
value: "0:database-node1-service:5432:1:/var/lib/postgresql/data:ALLOW_TO_FAILOVER,1:database-node2-service::::,2:database-node3-service::::,3:database-node4-service::::,4:database-node5-service::::"

```
to configure pgpool.
