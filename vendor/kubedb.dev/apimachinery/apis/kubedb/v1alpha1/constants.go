package v1alpha1

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"
	LabelRole         = GenericKey + "/role"

	ComponentDatabase = "database"
	RoleStats         = "stats"
	DefaultStatsPath  = "/metrics"

	PostgresKey         = ResourceSingularPostgres + "." + GenericKey
	ElasticsearchKey    = ResourceSingularElasticsearch + "." + GenericKey
	MySQLKey            = ResourceSingularMySQL + "." + GenericKey
	PerconaXtraDBKey    = ResourceSingularPerconaXtraDB + "." + GenericKey
	MongoDBKey          = ResourceSingularMongoDB + "." + GenericKey
	RedisKey            = ResourceSingularRedis + "." + GenericKey
	MemcachedKey        = ResourceSingularMemcached + "." + GenericKey
	EtcdKey             = ResourceSingularEtcd + "." + GenericKey
	ProxySQLKey         = ResourceSingularProxySQL + "." + GenericKey
	SnapshotKey         = ResourceSingularSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	AnnotationInitialized = GenericKey + "/initialized"
	AnnotationJobType     = GenericKey + "/job-type"

	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "prom-http"

	JobTypeBackup  = "backup"
	JobTypeRestore = "restore"

	ElasticsearchRestPort     = 9200
	ElasticsearchRestPortName = "http"
	ElasticsearchNodePort     = 9300
	ElasticsearchNodePortName = "transport"

	MongoDBShardPort    = 27017
	MongoDBConfigdbPort = 27017
	MongoDBMongosPort   = 27017

	MySQLUserKey         = "username"
	MySQLPasswordKey     = "password"
	MySQLNodePort        = 3306
	MySQLGroupComPort    = 33060
	MySQLMaxGroupMembers = 9
	// The recommended MySQL server version for group replication (GR)
	MySQLGRRecommendedVersion = "5.7.25"
	MySQLDefaultGroupSize     = 3
	MySQLDefaultBaseServerID  = uint(1)
	// The server id for each group member must be unique and in the range [1, 2^32 - 1]
	// And the maximum group size is 9. So MySQLMaxBaseServerID is the maximum safe value
	// for BaseServerID calculated as max MySQL server_id value - max Replication Group size.
	MySQLMaxBaseServerID = uint(4294967295 - 9)

	PerconaXtraDBClusterRecommendedVersion    = "5.7"
	PerconaXtraDBMaxClusterNameLength         = 32
	PerconaXtraDBStandaloneReplicas           = 1
	PerconaXtraDBDefaultClusterSize           = 3
	PerconaXtraDBDataMountPath                = "/var/lib/mysql"
	PerconaXtraDBDataLostFoundPath            = PerconaXtraDBDataMountPath + "lost+found"
	PerconaXtraDBInitDBMountPath              = "/docker-entrypoint-initdb.d"
	PerconaXtraDBCustomConfigMountPath        = "/etc/percona-server.conf.d/"
	PerconaXtraDBClusterCustomConfigMountPath = "/etc/mysql/percona-xtradb-cluster.conf.d/"

	LabelProxySQLName        = ProxySQLKey + "/name"
	LabelProxySQLLoadBalance = ProxySQLKey + "/load-balance"

	ProxySQLUserKey               = "proxysqluser"
	ProxySQLPasswordKey           = "proxysqlpass"
	ProxySQLMySQLNodePort         = 6033
	ProxySQLAdminPort             = 6032
	ProxySQLAdminPortName         = "proxyadm"
	ProxySQLDataMountPath         = "/var/lib/proxysql"
	ProxySQLCustomConfigMountPath = "/etc/custom-proxysql.cnf"

	RedisShardKey   = RedisKey + "/shard"
	RedisNodePort   = 6379
	RedisGossipPort = 16379
)
