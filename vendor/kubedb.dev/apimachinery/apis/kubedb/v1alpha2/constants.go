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

	// =========================== PostgreSQL Constants ============================
	PostgresDatabasePortName         = "db"
	PostgresPrimaryServicePortName   = "primary"
	PostgresStandbyServicePortName   = "standby"
	PostgresDatabasePort             = 5432
	PostgresPodPrimary               = "primary"
	PostgresPodStandby               = "standby"
	EnvPostgresUser                  = "POSTGRES_USER"
	EnvPostgresPassword              = "POSTGRES_PASSWORD"
	PostgresRootUser                 = "postgres"
	PostgresCoordinatorContainerName = "pg-coordinator"
	PostgresCoordinatorPort          = 2380
	PostgresCoordinatorPortName      = "coordinator"
	PostgresContainerName            = ResourceSingularPostgres

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

	// =========================== ProxySQL Constants ============================
	LabelProxySQLName        = ProxySQLKey + "/name"
	LabelProxySQLLoadBalance = ProxySQLKey + "/load-balance"

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
	RedisContainerName          = ResourceSingularRedis
	RedisSentinelContainerName  = "redissentinel"
	DefaultConfigKey            = "default.conf"
	RedisShardKey               = RedisKey + "/shard"
	RedisDatabasePortName       = "db"
	RedisPrimaryServicePortName = "primary"
	RedisDatabasePort           = 6379
	RedisSentinelPort           = 26379
	RedisGossipPortName         = "gossip"
	RedisGossipPort             = 16379
	RedisSentinelPortName       = "sentinel"

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

	RedisNodeFlagMaster = "master"
	RedisNodeFlagNoAddr = "noaddr"
	RedisNodeFlagSlave  = "slave"

	RedisKeyFileSecretSuffix = "key"
	RedisPEMSecretSuffix     = "pem"
	RedisRootUsername        = "default"
	EnvRedisUser             = "USERNAME"
	EnvRedisPassword         = "REDISCLI_AUTH"

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
	DatabaseWriteAccessCheckFailed             = "DatabaseReadAccessCheckFailed"
	InternalUsersCredentialSyncFailed          = "InternalUsersCredentialsSyncFailed"
	InternalUsersCredentialsSyncedSuccessfully = "InternalUsersCredentialsSyncedSuccessfully"
)

const (
	KafkaPortNameREST                 = "http"
	KafkaPortNameController           = "controller"
	KafkaBrokerClientPortName         = "broker"
	KafkaControllerClientPortName     = "controller"
	KafkaPortNameInternal             = "internal"
	KafkaTopicNameHealth              = "kafka-health"
	KafkaTopicDeletionThresholdOffset = 1000
	KafkaRESTPort                     = 9092
	KafkaControllerRESTPort           = 9093
	KafkaInternalRESTPort             = 29092

	KafkaContainerName        = "kafka"
	KafkaUserAdmin            = "admin"
	KafkaNodeRoleSet          = "set"
	KafkaNodeRolesCombined    = "controller,broker"
	KafkaNodeRolesController  = "controller"
	KafkaNodeRolesBrokers     = "broker"
	KafkaStandbyServiceSuffix = "standby"

	KafkaBrokerListener     = "KafkaBrokerListener"
	KafkaControllerListener = "KafkaControllerListener"

	KafkaDataDir                  = "/var/log/kafka"
	KafkaMetaDataDir              = "/var/log/kafka/metadata"
	KafkaCertDir                  = "/var/private/ssl"
	KafkaConfigDir                = "/opt/kafka/config/kafkaconfig"
	KafkaTempConfigDir            = "/opt/kafka/config/temp-config"
	KafkaConfigFileName           = "config.properties"
	KafkaSSLPropertiesFileName    = "ssl.properties"
	KafkaClientAuthConfigFileName = "clientauth.properties"

	KafkaListeners                   = "listeners"
	KafkaAdvertisedListeners         = "advertised.listeners"
	KafkaListenerSecurityProtocolMap = "listener.security.protocol.map"
	KafkaControllerNodeCount         = "controller.count"
	KafkaControllerQuorumVoters      = "controller.quorum.voters"
	KafkaControllerListenersName     = "controller.listener.names"
	KafkaInterBrokerListener         = "inter.broker.listener.name"
	KafkaNodeRole                    = "process.roles"
	KafkaClusterID                   = "cluster.id"
	KafkaDataDirName                 = "log.dirs"
	KafkaMetadataDirName             = "metadata.log.dir"
	KafkaKeystorePasswordKey         = "keystore_password"
	KafkaTruststorePasswordKey       = "truststore_password"
	KafkaServerKeystoreKey           = "server.keystore.jks"
	KafkaServerTruststoreKey         = "server.truststore.jks"
	KafkaSecurityProtocol            = "security.protocol"
	KafkaGracefulShutdownTimeout     = "task.shutdown.graceful.timeout.ms"

	KafkaEndpointVerifyAlgo  = "ssl.endpoint.identification.algorithm"
	KafkaKeystoreLocation    = "ssl.keystore.location"
	KafkaTruststoreLocation  = "ssl.truststore.location"
	KafkaKeystorePassword    = "ssl.keystore.password"
	KafkaTruststorePassword  = "ssl.truststore.password"
	KafkaKeyPassword         = "ssl.key.password"
	KafkaKeystoreDefaultPass = "changeit"

	KafkaEnabledSASLMechanisms       = "sasl.enabled.mechanisms"
	KafkaSASLMechanism               = "sasl.mechanism"
	KafkaMechanismControllerProtocol = "sasl.mechanism.controller.protocol"
	KafkaSASLInterBrokerProtocol     = "sasl.mechanism.inter.broker.protocol"
	KafkaSASLPLAINConfigKey          = "listener.name.SASL_PLAINTEXT.plain.sasl.jaas.config"
	KafkaSASLSSLConfigKey            = "listener.name.SASL_SSL.plain.sasl.jaas.config"
	KafkaSASLJAASConfig              = "sasl.jaas.config"
	KafkaServiceName                 = "serviceName"
	KafkaSASLPlainMechanism          = "PLAIN"

	KafkaVolumeData       = "data"
	KafkaVolumeConfig     = "kafkaconfig"
	KafkaVolumeTempConfig = "temp-config"

	EnvKafkaUser     = "KAFKA_USER"
	EnvKafkaPassword = "KAFKA_PASSWORD"

	KafkaListenerPLAINTEXTProtocol = "PLAINTEXT"
	KafkaListenerSASLProtocol      = "SASL_PLAINTEXT"
	KafkaListenerSASLSSLProtocol   = "SASL_SSL"
)

// Resource kind related constants
const (
	ResourceKindStatefulSet = "StatefulSet"
)

var (
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

	// DefaultResourcesElasticSearch must be used for elasticsearch
	// to avoid OOMKILLED while deploying ES V8
	DefaultResourcesElasticSearch = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".500"),
			core.ResourceMemory: resource.MustParse("1.5Gi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("1.5Gi"),
		},
	}
)
