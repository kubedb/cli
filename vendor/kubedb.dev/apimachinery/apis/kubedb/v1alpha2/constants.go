/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"time"

	"kubedb.dev/apimachinery/apis/kubedb"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	// Deprecated
	DatabaseNamePrefix = "kubedb"

	KubeDBOrganization = "kubedb"

	LabelRole = kubedb.GroupName + "/role"

	ReplicationModeDetectorContainerName = "replication-mode-detector"
	DatabasePodPrimary                   = "primary"
	DatabasePodStandby                   = "standby"

	ComponentDatabase         = "database"
	ComponentConnectionPooler = "connection-pooler"
	RoleStats                 = "stats"
	DefaultStatsPath          = "/metrics"
	DefaultPasswordLength     = 16
	HealthCheckInterval       = 10 * time.Second

	ContainerExporterName = "exporter"
	LocalHost             = "localhost"
	LocalHostIP           = "127.0.0.1"

	DBCustomConfigName             = "custom-config"
	DefaultVolumeClaimTemplateName = "data"

	DBTLSVolume         = "tls-volume"
	DBExporterTLSVolume = "exporter-tls-volume"

	CACert = "ca.crt"

	// =========================== Database key Constants ============================
	PostgresKey      = ResourceSingularPostgres + "." + kubedb.GroupName
	ElasticsearchKey = ResourceSingularElasticsearch + "." + kubedb.GroupName
	MySQLKey         = ResourceSingularMySQL + "." + kubedb.GroupName
	MariaDBKey       = ResourceSingularMariaDB + "." + kubedb.GroupName
	PerconaXtraDBKey = ResourceSingularPerconaXtraDB + "." + kubedb.GroupName
	MongoDBKey       = ResourceSingularMongoDB + "." + kubedb.GroupName
	RedisKey         = ResourceSingularRedis + "." + kubedb.GroupName
	MemcachedKey     = ResourceSingularMemcached + "." + kubedb.GroupName
	EtcdKey          = ResourceSingularEtcd + "." + kubedb.GroupName
	ProxySQLKey      = ResourceSingularProxySQL + "." + kubedb.GroupName

	// =========================== Elasticsearch Constants ============================
	ElasticsearchRestPort                        = 9200
	ElasticsearchRestPortName                    = "http"
	ElasticsearchTransportPort                   = 9300
	ElasticsearchTransportPortName               = "transport"
	ElasticsearchPerformanceAnalyzerPort         = 9600
	ElasticsearchPerformanceAnalyzerPortName     = "analyzer"
	ElasticsearchNodeRoleSet                     = "set"
	ElasticsearchConfigDir                       = "/usr/share/elasticsearch/config"
	ElasticsearchOpenSearchConfigDir             = "/usr/share/opensearch/config"
	ElasticsearchSecureSettingsDir               = "/elasticsearch/secure-settings"
	ElasticsearchTempConfigDir                   = "/elasticsearch/temp-config"
	ElasticsearchCustomConfigDir                 = "/elasticsearch/custom-config"
	ElasticsearchDataDir                         = "/usr/share/elasticsearch/data"
	ElasticsearchOpenSearchDataDir               = "/usr/share/opensearch/data"
	ElasticsearchTempDir                         = "/tmp"
	ElasticsearchOpendistroSecurityConfigDir     = "/usr/share/elasticsearch/plugins/opendistro_security/securityconfig"
	ElasticsearchOpenSearchSecurityConfigDir     = "/usr/share/opensearch/plugins/opensearch-security/securityconfig"
	ElasticsearchOpenSearchSecurityConfigDirV2   = "/usr/share/opensearch/config/opensearch-security"
	ElasticsearchSearchGuardSecurityConfigDir    = "/usr/share/elasticsearch/plugins/search-guard-%v/sgconfig"
	ElasticsearchOpendistroReadallMonitorRole    = "readall_and_monitor"
	ElasticsearchOpenSearchReadallMonitorRole    = "readall_and_monitor"
	ElasticsearchSearchGuardReadallMonitorRoleV7 = "SGS_READALL_AND_MONITOR"
	ElasticsearchSearchGuardReadallMonitorRoleV6 = "sg_readall_and_monitor"
	ElasticsearchStatusGreen                     = "green"
	ElasticsearchStatusYellow                    = "yellow"
	ElasticsearchStatusRed                       = "red"
	ElasticsearchInitSysctlContainerName         = "init-sysctl"
	ElasticsearchInitConfigMergerContainerName   = "config-merger"
	ElasticsearchContainerName                   = "elasticsearch"
	ElasticsearchExporterContainerName           = "exporter"
	ElasticsearchSearchGuardRolesMappingFileName = "sg_roles_mapping.yml"
	ElasticsearchSearchGuardInternalUserFileName = "sg_internal_users.yml"
	ElasticsearchOpendistroRolesMappingFileName  = "roles_mapping.yml"
	ElasticsearchOpendistroInternalUserFileName  = "internal_users.yml"
	ElasticsearchJavaOptsEnv                     = "ES_JAVA_OPTS"
	ElasticsearchOpenSearchJavaOptsEnv           = "OPENSEARCH_JAVA_OPTS"
	ElasticsearchVolumeConfig                    = "esconfig"
	ElasticsearchVolumeTempConfig                = "temp-config"
	ElasticsearchVolumeSecurityConfig            = "security-config"
	ElasticsearchVolumeSecureSettings            = "secure-settings"
	ElasticsearchVolumeCustomConfig              = "custom-config"
	ElasticsearchVolumeData                      = "data"
	ElasticsearchVolumeTemp                      = "temp"

	// Ref:
	//	- https://www.elastic.co/guide/en/elasticsearch/reference/7.6/heap-size.html#heap-size
	//	- no more than 50% of your physical RAM
	//	- no more than 32GB that the JVM uses for compressed object pointers (compressed oops)
	//	- no more than 26GB for zero-based compressed oops;
	// 26 GB is safe on most systems
	ElasticsearchMaxHeapSize = 26 * 1024 * 1024 * 1024
	// 128MB
	ElasticsearchMinHeapSize = 128 * 1024 * 1024

	// =========================== Memcached Constants ============================
	MemcachedDatabasePortName       = "db"
	MemcachedPrimaryServicePortName = "primary"
	MemcachedDatabasePort           = 11211

	// =========================== MongoDB Constants ============================

	MongoDBDatabasePortName       = "db"
	MongoDBPrimaryServicePortName = "primary"
	MongoDBDatabasePort           = 27017
	MongoDBKeyFileSecretSuffix    = "-key"
	MongoDBRootUsername           = "root"
	MongoDBCustomConfigFile       = "mongod.conf"
	MongoDBReplicaSetConfig       = "replicaset.json"
	MongoDBConfigurationJSFile    = "configuration.js"
	NodeTypeMongos                = "mongos"
	NodeTypeShard                 = "shard"
	NodeTypeConfig                = "configsvr"
	NodeTypeArbiter               = "arbiter"
	NodeTypeHidden                = "hidden"
	NodeTypeReplica               = "replica"
	NodeTypeStandalone            = "standalone"

	MongoDBWorkDirectoryName = "workdir"
	MongoDBWorkDirectoryPath = "/work-dir"

	MongoDBCertDirectoryName = "certdir"

	MongoDBDataDirectoryName = "datadir"
	MongoDBDataDirectoryPath = "/data/db"

	MongoDBInitInstallContainerName   = "copy-config"
	MongoDBInitBootstrapContainerName = "bootstrap"

	MongoDBConfigDirectoryName = "config"
	MongoDBConfigDirectoryPath = "/data/configdb"

	MongoDBInitialConfigDirectoryName = "configdir"
	MongoDBInitialConfigDirectoryPath = "/configdb-readonly"

	MongoDBInitScriptDirectoryName = "init-scripts"
	MongoDBInitScriptDirectoryPath = "/init-scripts"

	MongoDBInitialDirectoryName = "initial-script"
	MongoDBInitialDirectoryPath = "/docker-entrypoint-initdb.d"

	MongoDBClientCertDirectoryName = "client-cert"
	MongoDBClientCertDirectoryPath = "/client-cert"

	MongoDBServerCertDirectoryName = "server-cert"
	MongoDBServerCertDirectoryPath = "/server-cert"

	MongoDBInitialKeyDirectoryName = "keydir"
	MongoDBInitialKeyDirectoryPath = "/keydir-readonly"

	MongoDBContainerName = ResourceSingularMongoDB

	MongoDBDefaultVolumeClaimTemplateName = MongoDBDataDirectoryName

	MongodbUser             = "root"
	MongoDBKeyForKeyFile    = "key.txt"
	MongoDBAuthSecretSuffix = "-auth"

	// =========================== MySQL Constants ============================
	MySQLMetricsExporterConfigSecretSuffix = "metrics-exporter-config"
	MySQLDatabasePortName                  = "db"
	MySQLRouterReadWritePortName           = "rw"
	MySQLRouterReadOnlyPortName            = "ro"
	MySQLPrimaryServicePortName            = "primary"
	MySQLStandbyServicePortName            = "standby"
	MySQLDatabasePort                      = 3306
	MySQLRouterReadWritePort               = 6446
	MySQLRouterReadOnlyPort                = 6447

	MySQLCoordinatorClientPort = 2379
	MySQLCoordinatorPort       = 2380
	MySQLCoordinatorStatus     = "Coordinator/Status"

	MySQLGroupComPort    = 33060
	MySQLMaxGroupMembers = 9
	// The recommended MySQL server version for group replication (GR)
	MySQLGRRecommendedVersion = "8.0.23"
	MySQLDefaultGroupSize     = 3
	MySQLRootUserName         = "MYSQL_ROOT_USERNAME"
	MySQLRootPassword         = "MYSQL_ROOT_PASSWORD"
	MySQLName                 = "MYSQL_NAME"
	MySQLRootUser             = "root"

	MySQLTLSConfigCustom     = "custom"
	MySQLTLSConfigSkipVerify = "skip-verify"
	MySQLTLSConfigTrue       = "true"
	MySQLTLSConfigFalse      = "false"
	MySQLTLSConfigPreferred  = "preferred"

	MySQLContainerName            = "mysql"
	MySQLRouterContainerName      = "mysql-router"
	MySQLRouterInitContainerName  = "mysql-router-init"
	MySQLCoordinatorContainerName = "mysql-coordinator"
	MySQLInitContainerName        = "mysql-init"

	MySQLRouterInitScriptDirectoryName = "init-scripts"
	MySQLRouterInitScriptDirectoryPath = "/scripts"
	MySQLRouterConfigDirectoryName     = "router-config-secret"
	MySQLRouterConfigDirectoryPath     = "/etc/mysqlrouter"
	MySQLRouterTLSDirectoryName        = "router-tls-volume"
	MySQLRouterTLSDirectoryPath        = "/etc/mysql/certs"
	MySQLReplicationUser               = "repl"

	MySQLComponentKey    = MySQLKey + "/component"
	MySQLComponentDB     = "database"
	MySQLComponentRouter = "router"

	// mysql volume and volume Mounts

	MySQLVolumeNameTemp      = "tmp"
	MySQLVolumeMountPathTemp = "/tmp"

	MySQLVolumeNameData      = "data"
	MySQLVolumeMountPathData = "/var/lib/mysql"

	MySQLVolumeNameUserInitScript      = "initial-script"
	MySQLVolumeMountPathUserInitScript = "/docker-entrypoint-initdb.d"

	MySQLVolumeNameInitScript      = "init-scripts"
	MySQLVolumeMountPathInitScript = "/scripts"

	MySQLVolumeNameCustomConfig      = "custom-config"
	MySQLVolumeMountPathCustomConfig = "/etc/mysql/conf.d"

	MySQLVolumeNameTLS      = "tls-volume"
	MySQLVolumeMountPathTLS = "/etc/mysql/certs"

	MySQLVolumeNameExporterTLS      = "exporter-tls-volume"
	MySQLVolumeMountPathExporterTLS = "/etc/mysql/certs"

	MySQLVolumeNameSourceCA      = "source-ca"
	MySQLVolumeMountPathSourceCA = "/etc/mysql/server/certs"

	// =========================== PerconaXtraDB Constants ============================
	PerconaXtraDBClusterRecommendedVersion     = "5.7"
	PerconaXtraDBMaxClusterNameLength          = 32
	PerconaXtraDBStandaloneReplicas            = 1
	PerconaXtraDBDefaultClusterSize            = 3
	PerconaXtraDBDataMountPath                 = "/var/lib/mysql"
	PerconaXtraDBDataLostFoundPath             = PerconaXtraDBDataMountPath + "/lost+found"
	PerconaXtraDBInitDBVolumeName              = "initial-script"
	PerconaXtraDBInitDBMountPath               = "/docker-entrypoint-initdb.d"
	PerconaXtraDBCustomConfigMountPath         = "/etc/percona-server.conf.d/"
	PerconaXtraDBClusterCustomConfigMountPath  = "/etc/mysql/custom.conf.d/"
	PerconaXtraDBCustomConfigVolumeName        = "custom-config"
	PerconaXtraDBTLSConfigCustom               = "custom"
	PerconaXtraDBInitContainerName             = "px-init"
	PerconaXtraDBCoordinatorContainerName      = "px-coordinator"
	PerconaXtraDBRunScriptVolumeName           = "run-script"
	PerconaXtraDBRunScriptVolumeMountPath      = "/run-script"
	PerconaXtraDBInitScriptVolumeName          = "init-scripts"
	PerconaXtraDBInitScriptVolumeMountPath     = "/scripts"
	PerconaXtraDBContainerName                 = ResourceSingularPerconaXtraDB
	PerconaXtraDBCertMountPath                 = "/etc/mysql/certs"
	PerconaXtraDBExporterConfigFileName        = "exporter.cnf"
	PerconaXtraDBGaleraClusterPrimaryComponent = "Primary"
	PerconaXtraDBServerTLSVolumeName           = "tls-server-config"
	PerconaXtraDBClientTLSVolumeName           = "tls-client-config"
	PerconaXtraDBExporterTLSVolumeName         = "tls-metrics-exporter-config"
	PerconaXtraDBMetricsExporterTLSVolumeName  = "metrics-exporter-config"
	PerconaXtraDBMetricsExporterConfigPath     = "/etc/mysql/config/exporter"
	PerconaXtraDBDataVolumeName                = "data"
	PerconaXtraDBMySQLUserGroupID              = 1001

	// =========================== MariaDB Constants ============================
	MariaDBMaxClusterNameLength          = 32
	MariaDBStandaloneReplicas            = 1
	MariaDBDefaultClusterSize            = 3
	MariaDBDataMountPath                 = "/var/lib/mysql"
	MariaDBDataLostFoundPath             = MariaDBDataMountPath + "/lost+found"
	MariaDBInitDBVolumeName              = "initial-script"
	MariaDBInitDBMountPath               = "/docker-entrypoint-initdb.d"
	MariaDBCustomConfigMountPath         = "/etc/mysql/conf.d/"
	MariaDBClusterCustomConfigMountPath  = "/etc/mysql/custom.conf.d/"
	MariaDBCustomConfigVolumeName        = "custom-config"
	MariaDBTLSConfigCustom               = "custom"
	MariaDBInitContainerName             = "mariadb-init"
	MariaDBCoordinatorContainerName      = "md-coordinator"
	MariaDBRunScriptVolumeName           = "run-script"
	MariaDBRunScriptVolumeMountPath      = "/run-script"
	MariaDBInitScriptVolumeName          = "init-scripts"
	MariaDBInitScriptVolumeMountPath     = "/scripts"
	MariaDBContainerName                 = ResourceSingularMariaDB
	MariaDBCertMountPath                 = "/etc/mysql/certs"
	MariaDBExporterConfigFileName        = "exporter.cnf"
	MariaDBGaleraClusterPrimaryComponent = "Primary"
	MariaDBServerTLSVolumeName           = "tls-server-config"
	MariaDBClientTLSVolumeName           = "tls-client-config"
	MariaDBExporterTLSVolumeName         = "tls-metrics-exporter-config"
	MariaDBMetricsExporterTLSVolumeName  = "metrics-exporter-config"
	MariaDBMetricsExporterConfigPath     = "/etc/mysql/config/exporter"
	MariaDBDataVolumeName                = "data"

	// =========================== SingleStore Constants ============================
	SinglestoreDatabasePortName       = "db"
	SinglestorePrimaryServicePortName = "primary"
	SinglestoreStudioPortName         = "studio"
	SinglestoreDatabasePort           = 3306
	SinglestoreStudioPort             = 8081
	SinglestoreExporterPort           = 9104
	SinglestoreRootUserName           = "ROOT_USERNAME"
	SinglestoreRootPassword           = "ROOT_PASSWORD"
	SinglestoreRootUser               = "root"
	DatabasePodMaster                 = "Master"
	DatabasePodAggregator             = "Aggregator"
	DatabasePodLeaf                   = "Leaf"
	PetSetTypeAggregator              = "aggregator"
	PetSetTypeLeaf                    = "leaf"
	SinglestoreDatabaseHealth         = "singlestore_health"
	SinglestoreTableHealth            = "singlestore_health_table"

	SinglestoreCoordinatorContainerName = "singlestore-coordinator"
	SinglestoreContainerName            = "singlestore"
	SinglestoreInitContainerName        = "singlestore-init"

	SinglestoreVolumeNameUserInitScript      = "initial-script"
	SinglestoreVolumeMountPathUserInitScript = "/docker-entrypoint-initdb.d"
	SinglestoreVolumeNameCustomConfig        = "custom-config"
	SinglestoreVolumeMountPathCustomConfig   = "/config"
	SinglestoreVolmeNameInitScript           = "init-scripts"
	SinglestoreVolumeMountPathInitScript     = "/scripts"
	SinglestoreVolumeNameData                = "data"
	SinglestoreVolumeMountPathData           = "/var/lib/memsql"

	// =========================== PostgreSQL Constants ============================
	PostgresDatabasePortName          = "db"
	PostgresPrimaryServicePortName    = "primary"
	PostgresStandbyServicePortName    = "standby"
	PostgresDatabasePort              = 5432
	PostgresPodPrimary                = "primary"
	PostgresPodStandby                = "standby"
	EnvPostgresUser                   = "POSTGRES_USER"
	EnvPostgresPassword               = "POSTGRES_PASSWORD"
	PostgresRootUser                  = "postgres"
	PostgresCoordinatorContainerName  = "pg-coordinator"
	PostgresCoordinatorPort           = 2380
	PostgresCoordinatorPortName       = "coordinator"
	PostgresContainerName             = ResourceSingularPostgres
	PostgresInitContainerName         = "postgres-init-container"
	PostgresCoordinatorClientPort     = 2379
	PostgresCoordinatorClientPortName = "coordinatclient"

	RaftMetricsExporterPort     = 23790
	RaftMetricsExporterPortName = "raft-metrics"

	PostgresInitVolumeName           = "initial-script"
	PostgresInitDir                  = "/var/initdb"
	PostgresSharedMemoryVolumeName   = "shared-memory"
	PostgresSharedMemoryDir          = "/dev/shm"
	PostgresDataVolumeName           = "data"
	PostgresDataDir                  = "/var/pv"
	PostgresCustomConfigVolumeName   = "custom-config"
	PostgresCustomConfigDir          = "/etc/config"
	PostgresRunScriptsVolumeName     = "run-scripts"
	PostgresRunScriptsDir            = "/run_scripts"
	PostgresRoleScriptsVolumeName    = "role-scripts"
	PostgresRoleScriptsDir           = "/role_scripts"
	PostgresSharedScriptsVolumeName  = "scripts"
	PostgresSharedScriptsDir         = "/scripts"
	PostgresSharedTlsVolumeName      = "certs"
	PostgresSharedTlsVolumeMountPath = "/tls/certs"

	PostgresKeyFileSecretSuffix = "key"
	PostgresPEMSecretSuffix     = "pem"
	PostgresDefaultUsername     = "postgres"
	PostgresPgCoordinatorStatus = "Coordinator/Status"
	// to pause the failover for postgres. this is helpful for ops request
	PostgresPgCoordinatorStatusPause = "Pause"
	// to resume the failover for postgres. this is helpful for ops request
	PostgresPgCoordinatorStatusResume = "Resume"

	// when we need to resume pg-coordinator as non transferable we are going to set this state.
	// this is useful when we have set a node as primary and you don't want other node rather then this node to become primary.
	PostgresPgCoordinatorStatusResumeNonTransferable = "NonTransferableResume"

	SharedBuffersGbAsByte = 1024 * 1024 * 1024
	SharedBuffersMbAsByte = 1024 * 1024

	SharedBuffersGbAsKiloByte = 1024 * 1024
	SharedBuffersMbAsKiloByte = 1024
	IPS_LOCK                  = "IPC_LOCK"
	SYS_RESOURCE              = "SYS_RESOURCE"
	DropCapabilityALL         = "ALL"

	// =========================== ProxySQL Constants ============================
	LabelProxySQLName                  = ProxySQLKey + "/name"
	LabelProxySQLLoadBalance           = ProxySQLKey + "/load-balance"
	LabelProxySQLLoadBalanceStandalone = "Standalone"

	ProxySQLContainerName          = ResourceSingularProxySQL
	ProxySQLDatabasePort           = 6033
	ProxySQLDatabasePortName       = "db"
	ProxySQLPrimaryServicePortName = "db"
	ProxySQLAdminPort              = 6032
	ProxySQLAdminPortName          = "admin"
	ProxySQLDataMountPath          = "/var/lib/proxysql"
	ProxySQLCustomConfigMountPath  = "/etc/custom-config"

	ProxySQLBackendSSLMountPath  = "/var/lib/certs"
	ProxySQLFrontendSSLMountPath = "/var/lib/frontend"
	ProxySQLClusterAdmin         = "cluster"
	ProxySQLClusterPasswordField = "cluster_password"
	ProxySQLTLSConfigCustom      = "custom"
	ProxySQLTLSConfigSkipVerify  = "skip-verify"

	ProxySQLMonitorUsername = "proxysql"
	ProxySQLAuthUsername    = "cluster"
	ProxySQLConfigSecretKey = "proxysql.cnf"

	// =========================== Redis Constants ============================
	RedisConfigKey = "redis.conf" // RedisConfigKey is going to create for the customize redis configuration
	// DefaultConfigKey is going to create for the default redis configuration
	RedisContainerName             = ResourceSingularRedis
	RedisSentinelContainerName     = "redissentinel"
	DefaultConfigKey               = "default.conf"
	RedisShardKey                  = RedisKey + "/shard"
	RedisDatabasePortName          = "db"
	RedisPrimaryServicePortName    = "primary"
	RedisDatabasePort              = 6379
	RedisSentinelPort              = 26379
	RedisGossipPortName            = "gossip"
	RedisGossipPort                = 16379
	RedisSentinelPortName          = "sentinel"
	RedisInitContainerName         = "redis-init"
	RedisCoordinatorContainerName  = "rd-coordinator"
	RedisSentinelInitContainerName = "sentinel-init"

	RedisScriptVolumeName      = "script-vol"
	RedisScriptVolumePath      = "/scripts"
	RedisDataVolumeName        = "data"
	RedisDataVolumePath        = "/data"
	RedisTLSVolumeName         = "tls-volume"
	RedisExporterTLSVolumeName = "exporter-tls-volume"
	RedisTLSVolumePath         = "/certs"
	RedisSentinelTLSVolumeName = "sentinel-tls-volume"
	RedisSentinelTLSVolumePath = "/sentinel-certs"
	RedisConfigVolumeName      = "redis-config"
	RedisConfigVolumePath      = "/usr/local/etc/redis/"
	RedisInitVolumeName        = "init-volume"
	RedisInitVolumePath        = "/init"

	RedisNodeFlagMaster = "master"
	RedisNodeFlagNoAddr = "noaddr"
	RedisNodeFlagSlave  = "slave"

	RedisKeyFileSecretSuffix = "key"
	RedisPEMSecretSuffix     = "pem"
	RedisRootUsername        = "default"

	EnvRedisUser              = "USERNAME"
	EnvRedisPassword          = "REDISCLI_AUTH"
	EnvRedisMode              = "REDIS_MODE"
	EnvRedisMajorRedisVersion = "MAJOR_REDIS_VERSION"

	// =========================== PgBouncer Constants ============================
	PgBouncerUpstreamServerCA               = "upstream-server-ca.crt"
	PgBouncerUpstreamServerClientCert       = "upstream-server-client.crt"
	PgBouncerUpstreamServerClientKey        = "upstream-server-client.key"
	PgBouncerClientCrt                      = "client.crt"
	PgBouncerClientKey                      = "client.key"
	PgBouncerCACrt                          = "ca.crt"
	PgBouncerTLSCrt                         = "tls.crt"
	PgBouncerTLSKey                         = "tls.key"
	PgBouncerDatabasePortName               = "db"
	PgBouncerPrimaryServicePortName         = "primary"
	PgBouncerDatabasePort                   = 5432
	PgBouncerConfigFile                     = "pgbouncer.ini"
	PgBouncerAdminUsername                  = "pgbouncer"
	PgBouncerDefaultPoolMode                = "session"
	PgBouncerDefaultIgnoreStartupParameters = "empty"

	// =========================== Pgpool Constants ============================
	EnvPostgresUsername                = "POSTGRES_USERNAME"
	EnvPgpoolPcpUser                   = "PGPOOL_PCP_USER"
	EnvPgpoolPcpPassword               = "PGPOOL_PCP_PASSWORD"
	EnvPgpoolPasswordEncryptionMethod  = "PGPOOL_PASSWORD_ENCRYPTION_METHOD"
	EnvEnablePoolPasswd                = "PGPOOL_ENABLE_POOL_PASSWD"
	EnvSkipPasswdEncryption            = "PGPOOL_SKIP_PASSWORD_ENCRYPTION"
	PgpoolConfigSecretMountPath        = "/config"
	PgpoolConfigVolumeName             = "pgpool-config"
	PgpoolContainerName                = "pgpool"
	PgpoolDefaultServicePort           = 9999
	PgpoolMonitoringDefaultServicePort = 9719
	PgpoolExporterDatabase             = "postgres"
	EnvPgpoolExporterDatabase          = "POSTGRES_DATABASE"
	EnvPgpoolService                   = "PGPOOL_SERVICE"
	EnvPgpoolServicePort               = "PGPOOL_SERVICE_PORT"
	EnvPgpoolSSLMode                   = "SSLMODE"
	PgpoolDefaultSSLMode               = "disable"
	PgpoolExporterContainerName        = "exporter"
	PgpoolAuthUsername                 = "pcp"
	SyncPeriod                         = 10
	// ========================================== ZooKeeper Constants =================================================//

	KubeDBZooKeeperRoleName         = "kubedb:zookeeper-version-reader"
	KubeDBZooKeeperRoleBindingName  = "kubedb:zookeeper-version-reader"
	ZooKeeperClientPortName         = "client"
	ZooKeeperClientPort             = 2181
	ZooKeeperQuorumPortName         = "quorum"
	ZooKeeperQuorumPort             = 2888
	ZooKeeperLeaderElectionPortName = "leader-election"
	ZooKeeperLeaderElectionPort     = 3888
	ZooKeeperMetricsPortName        = "metrics"
	ZooKeeperMetricsPort            = 7000
	ZooKeeperAdminServerPortName    = "admin-server"
	ZooKeeperAdminServerPort        = 8080
	ZooKeeperNode                   = "/kubedb_health_checker_node"
	ZooKeeperData                   = "kubedb_health_checker_data"
	ZooKeeperConfigVolumeName       = "zookeeper-config"
	ZooKeeperConfigVolumePath       = "/conf"
	ZooKeeperDataVolumeName         = "data"
	ZooKeeperDataVolumePath         = "/data"
	ZooKeeperScriptVolumeName       = "script-vol"
	ZooKeeperScriptVolumePath       = "/scripts"
	ZooKeeperContainerName          = ResourceSingularZooKeeper
	ZooKeeperInitContainerName      = ResourceSingularZooKeeper + "-init"

	ZooKeeperConfigFileName               = "zoo.cfg"
	ZooKeeperLog4jPropertiesFileName      = "log4j.properties"
	ZooKeeperLog4jQuietPropertiesFileName = "log4j-quiet.properties"

	EnvZooKeeperDomain          = "DOMAIN"
	EnvZooKeeperQuorumPort      = "QUORUM_PORT"
	EnvZooKeeperLeaderPort      = "LEADER_PORT"
	EnvZooKeeperClientHost      = "CLIENT_HOST"
	EnvZooKeeperClientPort      = "CLIENT_PORT"
	EnvZooKeeperAdminServerHost = "ADMIN_SERVER_HOST"
	EnvZooKeeperAdminServerPort = "ADMIN_SERVER_PORT"
	EnvZooKeeperClusterName     = "CLUSTER_NAME"
	EnvZooKeeperClusterSize     = "CLUSTER_SIZE"
	EnvZooKeeperUser            = "ZK_USER"
	EnvZooKeeperPassword        = "ZK_PASSWORD"
	EnvZooKeeperJaasFilePath    = "ZK_JAAS_FILE_PATH"
	EnvZooKeeperJVMFLags        = "JVMFLAGS"

	ZooKeeperSuperUsername       = "super"
	ZooKeeperSASLAuthLoginConfig = "-Djava.security.auth.login.config"
	ZooKeeperJaasFilePath        = "/data/jaas.conf"
)

// List of possible condition types for a KubeDB object
const (
	// used for Databases that have started provisioning
	DatabaseProvisioningStarted = "ProvisioningStarted"
	// used for Databases which completed provisioning
	DatabaseProvisioned = "Provisioned"
	// used for Databases that are currently being initialized using stash
	DatabaseDataRestoreStarted = "DataRestoreStarted"
	// used for Databases that have been initialized using stash
	DatabaseDataRestored = "DataRestored"
	// used for Databases whose pods are ready
	DatabaseReplicaReady = "ReplicaReady"
	// used for Databases that are currently accepting connection
	DatabaseAcceptingConnection = "AcceptingConnection"
	// used for Databases that report status OK (also implies that we can connect to it)
	DatabaseReady = "Ready"
	// used for database that reports ok when all the instances are available
	ServerReady = "ServerReady"
	// used for Databases that are paused
	DatabasePaused = "Paused"
	// used for Databases that are halted
	DatabaseHalted = "Halted"
	// used for pausing health check of a Database
	DatabaseHealthCheckPaused = "HealthCheckPaused"
	// used for Databases whose internal user credentials are synced
	InternalUsersSynced = "InternalUsersSynced"
	// user for databases that have read access
	DatabaseReadAccess = "DatabaseReadAccess"
	// user for databases that have write access
	DatabaseWriteAccess = "DatabaseWriteAccess"

	// Condition reasons
	DataRestoreStartedByExternalInitializer    = "DataRestoreStartedByExternalInitializer"
	DataRestoreInterrupted                     = "DataRestoreInterrupted"
	DatabaseSuccessfullyRestored               = "SuccessfullyDataRestored"
	FailedToRestoreData                        = "FailedToRestoreData"
	AllReplicasAreReady                        = "AllReplicasReady"
	SomeReplicasAreNotReady                    = "SomeReplicasNotReady"
	DatabaseAcceptingConnectionRequest         = "DatabaseAcceptingConnectionRequest"
	DatabaseNotAcceptingConnectionRequest      = "DatabaseNotAcceptingConnectionRequest"
	ReadinessCheckSucceeded                    = "ReadinessCheckSucceeded"
	ReadinessCheckFailed                       = "ReadinessCheckFailed"
	DatabaseProvisioningStartedSuccessfully    = "DatabaseProvisioningStartedSuccessfully"
	DatabaseSuccessfullyProvisioned            = "DatabaseSuccessfullyProvisioned"
	DatabaseHaltedSuccessfully                 = "DatabaseHaltedSuccessfully"
	DatabaseReadAccessCheckSucceeded           = "DatabaseReadAccessCheckSucceeded"
	DatabaseWriteAccessCheckSucceeded          = "DatabaseWriteAccessCheckSucceeded"
	DatabaseReadAccessCheckFailed              = "DatabaseReadAccessCheckFailed"
	DatabaseWriteAccessCheckFailed             = "DatabaseWriteAccessCheckFailed"
	InternalUsersCredentialSyncFailed          = "InternalUsersCredentialsSyncFailed"
	InternalUsersCredentialsSyncedSuccessfully = "InternalUsersCredentialsSyncedSuccessfully"
)

const (
	KafkaPortNameREST                  = "http"
	KafkaPortNameController            = "controller"
	KafkaPortNameCruiseControlListener = "cc-listener"
	KafkaPortNameCruiseControlREST     = "cc-rest"
	KafkaBrokerClientPortName          = "broker"
	KafkaControllerClientPortName      = "controller"
	KafkaPortNameLocal                 = "local"
	KafkaTopicNameHealth               = "kafka-health"
	KafkaTopicDeletionThresholdOffset  = 1000
	KafkaBrokerMaxID                   = 1000
	KafkaRESTPort                      = 9092
	KafkaControllerRESTPort            = 9093
	KafkaLocalRESTPort                 = 29092
	KafkaCruiseControlRESTPort         = 9090
	KafkaCruiseControlListenerPort     = 9094
	KafkaCCDefaultInNetwork            = 500000
	KafkaCCDefaultOutNetwork           = 500000

	KafkaContainerName          = "kafka"
	KafkaUserAdmin              = "admin"
	KafkaNodeRoleSet            = "set"
	KafkaNodeRolesCombined      = "controller,broker"
	KafkaNodeRolesController    = "controller"
	KafkaNodeRolesBrokers       = "broker"
	KafkaNodeRolesCruiseControl = "cruise-control"
	KafkaStandbyServiceSuffix   = "standby"

	KafkaBrokerListener     = "KafkaBrokerListener"
	KafkaControllerListener = "KafkaControllerListener"

	KafkaDataDir                              = "/var/log/kafka"
	KafkaMetaDataDir                          = "/var/log/kafka/metadata"
	KafkaCertDir                              = "/var/private/ssl"
	KafkaConfigDir                            = "/opt/kafka/config/kafkaconfig"
	KafkaTempConfigDir                        = "/opt/kafka/config/temp-config"
	KafkaCustomConfigDir                      = "/opt/kafka/config/custom-config"
	KafkaCCTempConfigDir                      = "/opt/cruise-control/temp-config"
	KafkaCCCustomConfigDir                    = "/opt/cruise-control/custom-config"
	KafkaCapacityConfigPath                   = "config/capacity.json"
	KafkaConfigFileName                       = "config.properties"
	KafkaServerCustomConfigFileName           = "server.properties"
	KafkaBrokerCustomConfigFileName           = "broker.properties"
	KafkaControllerCustomConfigFileName       = "controller.properties"
	KafkaSSLPropertiesFileName                = "ssl.properties"
	KafkaClientAuthConfigFileName             = "clientauth.properties"
	KafkaCruiseControlConfigFileName          = "cruisecontrol.properties"
	KafkaCruiseControlCapacityConfigFileName  = "capacity.json"
	KafkaCruiseControlBrokerSetConfigFileName = "brokerSets.json"
	KafkaCruiseControlClusterConfigFileName   = "clusterConfigs.json"
	KafkaCruiseControlLog4jConfigFileName     = "log4j.properties"
	KafkaCruiseControlUIConfigFileName        = "config.csv"

	KafkaListeners                         = "listeners"
	KafkaAdvertisedListeners               = "advertised.listeners"
	KafkaBootstrapServers                  = "bootstrap.servers"
	KafkaListenerSecurityProtocolMap       = "listener.security.protocol.map"
	KafkaControllerNodeCount               = "controller.count"
	KafkaControllerQuorumVoters            = "controller.quorum.voters"
	KafkaControllerListenersName           = "controller.listener.names"
	KafkaInterBrokerListener               = "inter.broker.listener.name"
	KafkaNodeRole                          = "process.roles"
	KafkaClusterID                         = "cluster.id"
	KafkaClientID                          = "client.id"
	KafkaDataDirName                       = "log.dirs"
	KafkaMetadataDirName                   = "metadata.log.dir"
	KafkaKeystorePasswordKey               = "keystore_password"
	KafkaTruststorePasswordKey             = "truststore_password"
	KafkaServerKeystoreKey                 = "server.keystore.jks"
	KafkaServerTruststoreKey               = "server.truststore.jks"
	KafkaSecurityProtocol                  = "security.protocol"
	KafkaGracefulShutdownTimeout           = "task.shutdown.graceful.timeout.ms"
	KafkaTopicConfigProviderClass          = "topic.config.provider.class"
	KafkaCapacityConfigFile                = "capacity.config.file"
	KafkaTwoStepVerification               = "two.step.verification.enabled"
	KafkaBrokerFailureDetection            = "kafka.broker.failure.detection.enable"
	KafkaMetricSamplingInterval            = "metric.sampling.interval.ms"
	KafkaPartitionMetricsWindow            = "partition.metrics.window.ms"
	KafkaPartitionMetricsWindowNum         = "num.partition.metrics.windows"
	KafkaSampleStoreTopicReplicationFactor = "sample.store.topic.replication.factor"

	KafkaEndpointVerifyAlgo  = "ssl.endpoint.identification.algorithm"
	KafkaKeystoreLocation    = "ssl.keystore.location"
	KafkaTruststoreLocation  = "ssl.truststore.location"
	KafkaKeystorePassword    = "ssl.keystore.password"
	KafkaTruststorePassword  = "ssl.truststore.password"
	KafkaKeyPassword         = "ssl.key.password"
	KafkaKeystoreDefaultPass = "changeit"

	KafkaMetricReporters       = "metric.reporters"
	KafkaAutoCreateTopicEnable = "auto.create.topics.enable"

	KafkaEnabledSASLMechanisms       = "sasl.enabled.mechanisms"
	KafkaSASLMechanism               = "sasl.mechanism"
	KafkaMechanismControllerProtocol = "sasl.mechanism.controller.protocol"
	KafkaSASLInterBrokerProtocol     = "sasl.mechanism.inter.broker.protocol"
	KafkaSASLPLAINConfigKey          = "listener.name.SASL_PLAINTEXT.plain.sasl.jaas.config"
	KafkaSASLSSLConfigKey            = "listener.name.SASL_SSL.plain.sasl.jaas.config"
	KafkaSASLJAASConfig              = "sasl.jaas.config"
	KafkaServiceName                 = "serviceName"
	KafkaSASLPlainMechanism          = "PLAIN"

	KafkaCCMetricSamplerClass            = "metric.sampler.class"
	KafkaCCCapacityConfig                = "capacity.config.file"
	KafkaCCTwoStepVerificationEnabled    = "two.step.verification.enabled"
	KafkaCCBrokerFailureDetectionEnabled = "kafka.broker.failure.detection.enable"
	KafkaOffSetTopicReplica              = "offsets.topic.replication.factor"
	KafkaTransactionStateLogReplica      = "transaction.state.log.replication.factor"
	KafkaTransactionSateLogMinISR        = "transaction.state.log.min.isr"
	KafkaLogCleanerMinLagSec             = "log.cleaner.min.compaction.lag.ms"
	KafkaLogCleanerBackoffMS             = "log.cleaner.backoff.ms"

	KafkaCCKubernetesMode                 = "cruise.control.metrics.reporter.kubernetes.mode"
	KafkaCCBootstrapServers               = "cruise.control.metrics.reporter.bootstrap.servers"
	KafkaCCMetricTopicAutoCreate          = "cruise.control.metrics.topic.auto.create"
	KafkaCCMetricTopicNumPartition        = "cruise.control.metrics.topic.num.partitions"
	KafkaCCMetricTopicReplica             = "cruise.control.metrics.topic.replication.factor"
	KafkaCCMetricReporterSecurityProtocol = "cruise.control.metrics.reporter.security.protocol"
	KafkaCCMetricReporterSaslMechanism    = "cruise.control.metrics.reporter.sasl.mechanism"
	KafkaCCSampleLoadingThreadsNum        = "num.sample.loading.threads"
	KafkaCCMinSamplesPerBrokerWindow      = "min.samples.per.broker.metrics.window"

	KafkaVolumeData         = "data"
	KafkaVolumeConfig       = "kafkaconfig"
	KafkaVolumeTempConfig   = "temp-config"
	KafkaVolumeCustomConfig = "custom-config"

	EnvKafkaUser     = "KAFKA_USER"
	EnvKafkaPassword = "KAFKA_PASSWORD"

	KafkaListenerPLAINTEXTProtocol = "PLAINTEXT"
	KafkaListenerSASLProtocol      = "SASL_PLAINTEXT"
	KafkaListenerSASLSSLProtocol   = "SASL_SSL"

	KafkaCCMetricsSampler         = "com.linkedin.kafka.cruisecontrol.monitor.sampling.CruiseControlMetricsReporterSampler"
	KafkaAdminTopicConfigProvider = "com.linkedin.kafka.cruisecontrol.config.KafkaAdminTopicConfigProvider"
	KafkaCCMetricReporter         = "com.linkedin.kafka.cruisecontrol.metricsreporter.CruiseControlMetricsReporter"
	KafkaJMXMetricReporter        = "org.apache.kafka.common.metrics.JmxReporter"

	// =========================== Solr Constants ============================
	ResourceCodeSolr      = "sl"
	ResourceKindSolr      = "Solr"
	ResourceSingularSolr  = "solr"
	ResourcePluralSolr    = "solrs"
	SolrPortName          = "http"
	SolrRestPort          = 8983
	SolrSecretKey         = "solr.xml"
	SolrContainerName     = "solr"
	SolrInitContainerName = "init-solr"
	SolrAdmin             = "admin"
	SecurityJSON          = "security.json"
	SolrZkDigest          = "zk-digest"
	SolrZkReadonlyDigest  = "zk-digest-readonly"

	SolrVolumeDefaultConfig = "default-config"
	SolrVolumeCustomConfig  = "custom-config"
	SolrVolumeAuthConfig    = "auth-config"
	SolrVolumeData          = "data"
	SolrVolumeConfig        = "slconfig"

	DistLibs              = "/opt/solr/dist"
	ContribLibs           = "/opt/solr/contrib/%s/lib"
	SysPropLibPlaceholder = "${solr.sharedLib:}"
	SolrHomeDir           = "/var/solr"
	SolrDataDir           = "/var/solr/data"
	SolrTempConfigDir     = "/temp-config"
	SolrCustomConfigDir   = "/custom-config"
	SolrSecurityConfigDir = "/var/security"

	SolrCloudHostKey                       = "host"
	SolrCloudHostValue                     = ""
	SolrCloudHostPortKey                   = "hostPort"
	SolrCloudHostPortValue                 = 80
	SolrCloudHostContextKey                = "hostContext"
	SolrCloudHostContextValue              = "solr"
	SolrCloudGenericCoreNodeNamesKey       = "genericCoreNodeNames"
	SolrCloudGenericCoreNodeNamesValue     = true
	SolrCloudZKClientTimeoutKey            = "zkClientTimeout"
	SolrCloudZKClientTimeoutValue          = 30000
	SolrCloudDistribUpdateSoTimeoutKey     = "distribUpdateSoTimeout"
	SolrCloudDistribUpdateSoTimeoutValue   = 600000
	SolrCloudDistribUpdateConnTimeoutKey   = "distribUpdateConnTimeout"
	SolrCloudDistribUpdateConnTimeoutValue = 60000
	SolrCloudZKCredentialProviderKey       = "zkCredentialsProvider"
	SolrCloudZKCredentialProviderValue     = "org.apache.solr.common.cloud.DigestZkCredentialsProvider"
	SolrCloudZKAclProviderKey              = "zkACLProvider"
	SolrCloudZKAclProviderValue            = "org.apache.solr.common.cloud.DigestZkACLProvider"
	SolrCloudZKCredentialsInjectorKey      = "zkCredentialsInjector"
	SolrCloudZKCredentialsInjectorValue    = "org.apache.solr.common.cloud.VMParamsZkCredentialsInjector"

	ShardHandlerFactorySocketTimeoutKey   = "socketTimeout"
	ShardHandlerFactorySocketTimeoutValue = 600000
	ShardHandlerFactoryConnTimeoutKey     = "connTimeout"
	ShardHandlerFactoryConnTimeoutValue   = 60000

	SolrKeysMaxBooleanClausesKey   = "maxBooleanClauses"
	SolrKeysMaxBooleanClausesValue = "solr.max.booleanClauses"
	SolrKeysSharedLibKey           = "sharedLib"
	SolrKeysShardLibValue          = "solr.sharedLib"
	SolrKeysHostPortKey            = "hostPort"
	SolrKeysHostPortValue          = "solr.port.advertise"
	SolrKeysAllowPathsKey          = "allowPaths"
	SolrKeysAllowPathsValue        = "solr.allowPaths"

	SolrConfMaxBooleanClausesKey   = "maxBooleanClauses"
	SolrConfMaxBooleanClausesValue = 1024
	SolrConfAllowPathsKey          = "allowPaths"
	SolrConfAllowPathsValue        = ""
	SolrConfSolrCloudKey           = "solrcloud"
	SolrConfShardHandlerFactoryKey = "shardHandlerFactory"
)

// =========================== Druid Constants ============================
const (
	DruidConfigDirCommon              = "/opt/druid/conf/druid/cluster/_common"
	DruidConfigDirCoordinatorOverlord = "/opt/druid/conf/druid/cluster/master/coordinator-overlord"
	DruidConfigDirHistoricals         = "/opt/druid/conf/druid/cluster/data/historical"
	DruidConfigDirMiddleManagers      = "/opt/druid/conf/druid/cluster/data/middleManager"
	DruidConfigDirBrokers             = "/opt/druid/conf/druid/cluster/query/broker"
	DruidConfigDirRouters             = "/opt/druid/conf/druid/cluster/query/router"
	DruidCConfigDirMySQLMetadata      = "/opt/druid/extensions/mysql-metadata-storage"

	DruidVolumeOperatorConfig = "operator-config-volume"
	DruidVolumeMainConfig     = "main-config-volume"
	DruidVolumeCustomConfig   = "custom-config"

	DruidOperatorConfigDir = "/tmp/config/operator-config"
	DruidMainConfigDir     = "/opt/druid/conf"
	DruidCustomConfigDir   = "/tmp/config/custom-config"

	DruidVolumeCommonConfig          = "common-config-volume"
	DruidCommonConfigFile            = "common.runtime.properties"
	DruidCoordinatorsJVMConfigFile   = "coordinators.jvm.config"
	DruidHistoricalsJVMConfigFile    = "historicals.jvm.config"
	DruidBrokersJVMConfigFile        = "brokers.jvm.config"
	DruidMiddleManagersJVMConfigFile = "middleManagers.jvm.config"
	DruidRoutersJVMConfigFile        = "routers.jvm.config"
	DruidCoordinatorsConfigFile      = "coordinators.properties"
	DruidHistoricalsConfigFile       = "historicals.properties"
	DruidMiddleManagersConfigFile    = "middleManagers.properties"
	DruidBrokersConfigFile           = "brokers.properties"
	DruidRoutersConfigFile           = "routers.properties"
	DruidVolumeMySQLMetadataStorage  = "mysql-metadata-storage"

	DruidContainerName     = "druid"
	DruidInitContainerName = "init-druid"
	DruidUserAdmin         = "admin"

	EnvDruidAdminPassword          = "DRUID_ADMIN_PASSWORD"
	EnvDruidMetdataStoragePassword = "DRUID_METADATA_STORAGE_PASSWORD"
	EnvDruidZKServicePassword      = "DRUID_ZK_SERVICE_PASSWORD"
	EnvDruidCoordinatorAsOverlord  = "DRUID_COORDINATOR_AS_OVERLORD"

	DruidPortCoordinators   = 8081
	DruidPortOverlords      = 8090
	DruidPortHistoricals    = 8083
	DruidPortMiddleManagers = 8091
	DruidPortBrokers        = 8082
	DruidPortRouters        = 8888

	// Common Runtime Configurations Properties
	// ZooKeeperSpec
	DruidZKServiceHost              = "druid.zk.service.host"
	DruidZKPathsBase                = "druid.zk.paths.base"
	DruidZKServiceCompress          = "druid.zk.service.compress"
	DruidZKServiceUserKey           = "druid.zk.service.user"
	DruidZKServicePasswordKey       = "druid.zk.service.pwd"
	DruidZKServicePasswordEnvConfig = "{\"type\": \"environment\", \"variable\": \"DRUID_ZK_SERVICE_PASSWORD\"}"

	// Metadata Storage
	DruidMetadataStorageTypeKey                    = "druid.metadata.storage.type"
	DruidMetadataStorageConnectorConnectURI        = "druid.metadata.storage.connector.connectURI"
	DruidMetadataStorageConnectURIPrefixMySQL      = "jdbc:mysql://"
	DruidMetadataStorageConnectURIPrefixPostgreSQL = "jdbc:postgresql://"
	DruidMetadataStorageConnectorUser              = "druid.metadata.storage.connector.user"
	DruidMetadataStorageConnectorPassword          = "druid.metadata.storage.connector.password"
	DruidMetadataStorageConnectorPasswordEnvConfig = "{\"type\": \"environment\", \"variable\": \"DRUID_METADATA_STORAGE_PASSWORD\"}"
	DruidMetadataStorageCreateTables               = "druid.metadata.storage.connector.createTables"

	// Deep Storage
	DruidDeepStorageTypeKey      = "druid.storage.type"
	DruidDeepStorageTypeS3       = "s3"
	DruidDeepStorageBaseKey      = "druid.storage.baseKey"
	DruidDeepStorageBucket       = "druid.storage.bucket"
	DruidS3AccessKey             = "druid.s3.accessKey"
	DruidS3SecretKey             = "druid.s3.secretKey"
	DruidS3EndpointSigningRegion = "druid.s3.endpoint.signingRegion"
	DruidS3EnablePathStyleAccess = "druid.s3.enablePathStyleAccess"
	DruidS3EndpointURL           = "druid.s3.endpoint.url"

	// Indexing service logs
	DruidIndexerLogsType           = "druid.indexer.logs.type"
	DruidIndexerLogsS3Bucket       = "druid.indexer.logs.s3Bucket"
	DruidIndexerLogsS3Prefix       = "druid.indexer.logs.s3Prefix"
	DruidEnableLookupSyncOnStartup = "druid.lookup.enableLookupSyncOnStartup"

	// Authentication
	DruidAuthAuthenticationChain                             = "druid.auth.authenticatorChain"
	DruidAuthAuthenticationChainValueBasic                   = "[\"basic\"]"
	DruidAuthAuthenticatorBasicType                          = "druid.auth.authenticator.basic.type"
	DruidAuthAuthenticatorBasicTypeValue                     = "basic"
	DruidAuthAuthenticatorBasicInitialAdminPassword          = "druid.auth.authenticator.basic.initialAdminPassword"
	DruidAuthAuthenticatorBasicInitialAdminPasswordEnvConfig = "{\"type\": \"environment\", \"variable\": \"DRUID_ADMIN_PASSWORD\"}"
	DruidAuthAuthenticatorBasicInitialInternalClientPassword = "druid.auth.authenticator.basic.initialInternalClientPassword"
	DruidAuthAuthenticatorBasicCredentialsValidatorType      = "druid.auth.authenticator.basic.credentialsValidator.type"
	DruidAuthAuthenticatorBasicSkipOnFailure                 = "druid.auth.authenticator.basic.skipOnFailure"
	DruidAuthAuthenticatorBasicAuthorizerName                = "druid.auth.authenticator.basic.authorizerName"

	// Escalator
	DruidAuthEscalatorType                   = "druid.escalator.type"
	DruidAuthEscalatorInternalClientUsername = "druid.escalator.internalClientUsername"
	DruidAuthEscalatorInternalClientPassword = "druid.escalator.internalClientPassword"
	DruidAuthEscalatorAuthorizerName         = "druid.escalator.authorizerName"
	DruidAuthAuthorizers                     = "druid.auth.authorizers"
	DruidAuthAuthorizerBasicType             = "druid.auth.authorizer.basic.type"

	// Extension Load List
	DruidExtensionLoadListKey               = "druid.extensions.loadList"
	DruidExtensionLoadList                  = "[\"druid-avro-extensions\", \"druid-s3-extensions\", \"druid-hdfs-storage\", \"druid-kafka-indexing-service\", \"druid-datasketches\", \"mysql-metadata-storage\", \"druid-basic-security\", \"druid-multi-stage-query\"]"
	DruidExtensionAvro                      = "druid-avro-extensions"
	DruidExtensionS3                        = "druid-s3-extensions"
	DruidExtensionHDFS                      = "druid-hdfs-storage"
	DruidExtensionGoogle                    = "druid-google-extensions"
	DruidExtensionAzure                     = "druid-azure-extensions"
	DruidExtensionKafkaIndexingService      = "druid-kafka-indexing-service"
	DruidExtensionDataSketches              = "druid-datasketches"
	DruidExtensionKubernetes                = "druid-kubernetes-extensions"
	DruidExtensionMySQLMetadataStorage      = "mysql-metadata-storage"
	DruidExtensionPostgreSQLMetadataStorage = "postgresql-metadata-storage"
	DruidExtensionBasicSecurity             = "druid-basic-security"
	DruidExtensionMultiStageQuery           = "druid-multi-stage-query"
	DruidExtensionPrometheusEmitter         = "prometheus-emitter"
	DruidService                            = "druid.service"

	// Monitoring Configurations
	DruidEmitter                                = "druid.emitter"
	DruidEmitterPrometheus                      = "prometheus"
	DruidEmitterPrometheusPortKey               = "druid.emitter.prometheus.port"
	DruidEmitterPrometheusPortVal               = 8080
	DruidMonitoringMonitorsKey                  = "druid.monitoring.monitors"
	DruidEmitterPrometheusStrategy              = "druid.emitter.prometheus.strategy"
	DruidMetricsJVMMonitor                      = "org.apache.druid.java.util.metrics.JvmMonitor"
	DruidMetricsServiceStatusMonitor            = "org.apache.druid.server.metrics.ServiceStatusMonitor"
	DruidMetricsQueryCountStatsMonitor          = "org.apache.druid.server.metrics.QueryCountStatsMonitor"
	DruidMonitoringHistoricalMetricsMonitor     = "org.apache.druid.server.metrics.HistoricalMetricsMonitor"
	DruidMonitoringSegmentsStatsMonitor         = "org.apache.druid.server.metrics.SegmentStatsMonitor"
	DruidMonitoringWorkerTaskCountsStatsMonitor = "org.apache.druid.server.metrics.WorkerTaskCountStatsMonitor"
	DruidMonitoringQueryCountStatsMonitor       = "org.apache.druid.server.metrics.QueryCountStatsMonitor"
	DruidMonitoringTaskCountStatsMonitor        = "org.apache.druid.server.metrics.TaskCountStatsMonitor"
	DruidMonitoringSysMonitor                   = "org.apache.druid.java.util.metrics.SysMonitor"

	/// Coordinators Configurations
	DruidCoordinatorStartDelay                = "druid.coordinator.startDelay"
	DruidCoordinatorPeriod                    = "druid.coordinator.period"
	DruidIndexerQueueStartDelay               = "druid.indexer.queue.startDelay"
	DruidManagerSegmentsPollDuration          = "druid.manager.segments.pollDuration"
	DruidCoordinatorKillAuditLogOn            = "druid.coordinator.kill.audit.on"
	DruidMillisToWaitBeforeDeleting           = "millisToWaitBeforeDeleting"
	DruidCoordinatorAsOverlord                = "druid.coordinator.asOverlord.enabled"
	DruidCoordinatorAsOverlordOverlordService = "druid.coordinator.asOverlord.overlordService"

	/// Overlords Configurations
	DruidServiceNameOverlords            = "druid/overlord"
	DruidIndexerStorageType              = "druid.indexer.storage.type"
	DruidIndexerAuditLogEnabled          = "druid.indexer.auditLog.enabled"
	DruidIndexerLogsKillEnables          = "druid.indexer.logs.kill.enabled"
	DruidIndexerLogsKillDurationToRetain = "druid.indexer.logs.kill.durationToRetain"
	DruidIndexerLogsKillInitialDelay     = "druid.indexer.logs.kill.initialDelay"
	DruidIndexerLogsKillDelay            = "druid.indexer.logs.kill.delay"

	DruidEmitterLoggingLogLevel = "druid.emitter.logging.logLevel"

	/// Historicals Configurations
	// Properties
	DruidProcessingNumOfThreads = "druid.processing.numThreads"

	// Segment Cache
	DruidHistoricalsSegmentCacheLocations              = "druid.segmentCache.locations"
	DruidHistoricalsSegmentCacheDropSegmentDelayMillis = "druid.segmentCache.dropSegmentDelayMillis"
	DruidHistoricalsSegmentCacheDir                    = "/druid/data/segments"
	DruidVolumeHistoricalsSegmentCache                 = "segment-cache"

	// Query Cache
	DruidHistoricalCacheUseCache      = "druid.historical.cache.useCache"
	DruidHistoricalCachePopulateCache = "druid.historical.cache.populateCache"
	DruidCacheSizeInBytes             = "druid.cache.sizeInBytes"

	// Values
	DruidSegmentCacheLocationsDefaultValue = "[{\"path\":\"/druid/data/segments\",\"maxSize\":10737418240}]"

	/// MiddleManagers Configurations
	// Properties
	DruidWorkerCapacity                                    = "druid.worker.capacity"
	DruidIndexerTaskBaseTaskDir                            = "druid.indexer.task.baseTaskDir"
	DruidWorkerTaskBaseTaskDirKey                          = "druid.worker.task.baseTaskDir"
	DruidWorkerTaskBaseTaskDir                             = "/var/druid/task"
	DruidWorkerBaseTaskDirSize                             = "druid.worker.baseTaskDirSize"
	DruidIndexerForkPropertyDruidProcessingBufferSizeBytes = "druid.indexer.fork.property.druid.processing.buffer.sizeBytes"
	DruidMiddleManagersVolumeBaseTaskDir                   = "base-task-dir"
	DruidVolumeMiddleManagersBaseTaskDir                   = "base-task-dir"

	// Values
	DruidIndexerTaskBaseTaskDirValue = "/druid/data/baseTaskDir"

	/// Brokers Configurations
	DruidBrokerHTTPNumOfConnections = "druid.broker.http.numConnections"
	DruidSQLEnable                  = "druid.sql.enable"

	/// Routers Configurations
	DruidRouterHTTPNumOfConnections = "druid.router.http.numConnections"
	DruidRouterHTTPNumOfMaxThreads  = "druid.router.http.numMaxThreads"

	// Common Nodes Configurations
	// Properties
	DruidPlaintextPort               = "druid.plaintextPort"
	DruidProcessingBufferSizeBytes   = "druid.processing.buffer.sizeBytes"
	DruidProcessingNumOfMergeBuffers = "druid.processing.numMergeBuffers"
	DruidServerHTTPNumOfThreads      = "druid.server.http.numThreads"

	// Health Check
	DruidHealthDataZero = "0"
	DruidHealthDataOne  = "1"
)

const (
	RabbitMQAMQPPort          = 5672
	RabbitMQPeerDiscoveryPort = 4369
	RabbitMQManagementUIPort  = 15672
	RabbitMQExporterPort      = 15692
	RabbitMQInterNodePort     = 25672

	RabbitMQVolumeData         = "data"
	RabbitMQVolumeConfig       = "rabbitmqconfig"
	RabbitMQVolumeTempConfig   = "temp-config"
	RabbitMQVolumeCustomConfig = "custom-config"

	RabbitMQDataDir         = "/var/lib/rabbitmq/mnesia"
	RabbitMQConfigDir       = "/config/"
	RabbitMQPluginsDir      = "/etc/rabbitmq/"
	RabbitMQCertDir         = "/var/private/ssl"
	RabbitMQTempConfigDir   = "/tmp/config/"
	RabbitMQCustomConfigDir = "/tmp/config/custom_config/"

	RabbitMQConfigVolName     = "rabbitmq-config"
	RabbitMQPluginsVolName    = "rabbitmq-plugins"
	RabbitMQTempConfigVolName = "temp-config"

	RabbitMQContainerName          = "rabbitmq"
	RabbitMQInitContainerName      = "rabbitmq-init"
	RabbitMQManagementPlugin       = "rabbitmq_management"
	RabbitMQPeerdiscoveryPlugin    = "rabbitmq_peer_discovery_k8s"
	RabbitMQLoopBackUserKey        = "loopback_users"
	RabbitMQLoopBackUserVal        = "none"
	RabbitMQDefaultTCPListenerKey  = "listeners.tcp.default"
	RabbitMQDefaultTCPListenerVal  = "5672"
	RabbitMQQueueMasterLocatorKey  = "queue_master_locator"
	RabbitMQQueueMasterLocatorVal  = "min-masters"
	RabbitMQDiskFreeLimitKey       = "disk_free_limit.absolute"
	RabbitMQDiskFreeLimitVal       = "2GB"
	RabbitMQPartitionHandingKey    = "cluster_partition_handling"
	RabbitMQPartitionHandingVal    = "pause_minority"
	RabbitMQPeerDiscoveryKey       = "cluster_formation.peer_discovery_backend"
	RabbitMQPeerDiscoveryVal       = "rabbit_peer_discovery_k8s"
	RabbitMQK8sHostKey             = "cluster_formation.k8s.host"
	RabbitMQK8sHostVal             = "kubernetes.default.svc.cluster.local"
	RabbitMQK8sAddressTypeKey      = "cluster_formation.k8s.address_type"
	RabbitMQK8sAddressTypeVal      = "hostname"
	RabbitMQNodeCleanupWarningKey  = "cluster_formation.node_cleanup.only_log_warning"
	RabbitMQNodeCleanupWarningVal  = "true"
	RabbitMQLogFileLevelKey        = "log.file.level"
	RabbitMQLogFileLevelVal        = "info"
	RabbitMQLogConsoleKey          = "log.console"
	RabbitMQLogConsoleVal          = "true"
	RabbitMQLogConsoleLevelKey     = "log.console.level"
	RabbitMQLogConsoleLevelVal     = "info"
	RabbitMQDefaultUserKey         = "default_user"
	RabbitMQDefaultUserVal         = "$(RABBITMQ_DEFAULT_USER)"
	RabbitMQDefaultPasswordKey     = "default_pass"
	RabbitMQDefaultPasswordVal     = "$(RABBITMQ_DEFAULT_PASS)"
	RabbitMQClusterNameKey         = "cluster_name"
	RabbitMQK8sSvcNameKey          = "cluster_formation.k8s.service_name"
	RabbitMQConfigFileName         = "rabbitmq.conf"
	RabbitMQEnabledPluginsFileName = "enabled_plugins"
	RabbitMQHealthCheckerQueueName = "kubedb-system"
)

// =========================== FerretDB Constants ============================
const (

	// envs
	EnvFerretDBUser     = "FERRETDB_PG_USER"
	EnvFerretDBPassword = "FERRETDB_PG_PASSWORD"
	EnvFerretDBHandler  = "FERRETDB_HANDLER"
	EnvFerretDBPgURL    = "FERRETDB_POSTGRESQL_URL"
	EnvFerretDBTLSPort  = "FERRETDB_LISTEN_TLS"
	EnvFerretDBCAPath   = "FERRETDB_LISTEN_TLS_CA_FILE"
	EnvFerretDBCertPath = "FERRETDB_LISTEN_TLS_CERT_FILE"
	EnvFerretDBKeyPath  = "FERRETDB_LISTEN_TLS_KEY_FILE"

	FerretDBContainerName = "ferretdb"
	FerretDBMainImage     = "ghcr.io/ferretdb/ferretdb"
	FerretDBUser          = "postgres"

	FerretDBServerPath = "/etc/certs/server"

	FerretDBDefaultPort = 27017
	FerretDBMetricsPort = 8080
	FerretDBTLSPort     = 27018

	FerretDBMetricsPath = "/debug/metrics"
)

// Resource kind related constants
const (
	ResourceKindStatefulSet = "StatefulSet"
	ResourceKindPetSet      = "PetSet"
)

var (
	DefaultInitContainerResource = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".200"),
			core.ResourceMemory: resource.MustParse("256Mi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("512Mi"),
		},
	}
	DefaultResources = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".500"),
			core.ResourceMemory: resource.MustParse("1024Mi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("1024Mi"),
		},
	}
	// CoordinatorDefaultResources must be used for raft backed coordinators to avoid unintended leader switches
	CoordinatorDefaultResources = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".200"),
			core.ResourceMemory: resource.MustParse("256Mi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("256Mi"),
		},
	}
	defaultArbiter = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceStorage: resource.MustParse("2Gi"),
			// these are the default cpu & memory for a coordinator container
			core.ResourceCPU:    resource.MustParse(".200"),
			core.ResourceMemory: resource.MustParse("256Mi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("256Mi"),
		},
	}

	// DefaultResourcesCPUIntensive is for MongoDB versions >= 6
	DefaultResourcesCPUIntensive = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".800"),
			core.ResourceMemory: resource.MustParse("1024Mi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("1024Mi"),
		},
	}

	// DefaultResourcesMemoryIntensive must be used for elasticsearch
	// to avoid OOMKILLED while deploying ES V8
	DefaultResourcesMemoryIntensive = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".500"),
			core.ResourceMemory: resource.MustParse("1.5Gi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("1.5Gi"),
		},
	}

	// DefaultResourcesCoreAndMemoryIntensive must be used for Solr
	DefaultResourcesCoreAndMemoryIntensiveSolr = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".900"),
			core.ResourceMemory: resource.MustParse("2Gi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("2Gi"),
		},
	}

	// DefaultResourcesMemoryIntensiveSDB must be used for Singlestore when enabled monitoring or version >= 8.5.x
	DefaultResourcesMemoryIntensiveSDB = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".500"),
			core.ResourceMemory: resource.MustParse("2Gi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("2Gi"),
		},
	}
)

func DefaultArbiter(computeOnly bool) core.ResourceRequirements {
	cp := defaultArbiter.DeepCopy()
	if computeOnly {
		delete(cp.Requests, core.ResourceStorage)
	}
	return *cp
}

const (
	InitFromGit          = "init-from-git"
	InitFromGitMountPath = "/git"
	GitSecretVolume      = "git-secret"
	GitSecretMountPath   = "/etc/git-secret"
	GitSyncContainerName = "git-sync"
)
