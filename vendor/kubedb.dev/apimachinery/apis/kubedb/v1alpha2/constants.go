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
	ElasticsearchOpendistroSecurityConfigDir     = "/usr/share/elasticsearch/plugins/opendistro_security/securityconfig"
	ElasticsearchOpenSearchSecurityConfigDir     = "/usr/share/opensearch/plugins/opensearch-security/securityconfig"
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
	NodeTypeMongos                = "mongos"
	NodeTypeShard                 = "shard"
	NodeTypeConfig                = "configsvr"

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
	MySQLGroupComPort                      = 33060
	MySQLMaxGroupMembers                   = 9
	// The recommended MySQL server version for group replication (GR)
	MySQLGRRecommendedVersion = "8.0.23"
	MySQLDefaultGroupSize     = 3
	MySQLRootUserName         = "MYSQL_ROOT_USERNAME"
	MySQLRootPassword         = "MYSQL_ROOT_PASSWORD"
	MySQLName                 = "MYSQL_NAME"

	MySQLTLSConfigCustom     = "custom"
	MySQLTLSConfigSkipVerify = "skip-verify"
	MySQLTLSConfigTrue       = "true"
	MySQLTLSConfigFalse      = "false"
	MySQLTLSConfigPreferred  = "preferred"

	MySQLRouterContainerName           = "mysql-router"
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

	// =========================== PerconaXtraDB Constants ============================
	PerconaXtraDBClusterRecommendedVersion    = "5.7"
	PerconaXtraDBMaxClusterNameLength         = 32
	PerconaXtraDBStandaloneReplicas           = 1
	PerconaXtraDBDefaultClusterSize           = 3
	PerconaXtraDBDataMountPath                = "/var/lib/mysql"
	PerconaXtraDBDataLostFoundPath            = PerconaXtraDBDataMountPath + "lost+found"
	PerconaXtraDBInitDBMountPath              = "/docker-entrypoint-initdb.d"
	PerconaXtraDBCustomConfigMountPath        = "/etc/percona-server.conf.d/"
	PerconaXtraDBClusterCustomConfigMountPath = "/etc/percona-xtradb-cluster.conf.d/"

	// =========================== MariaDB Constants ============================
	MariaDBMaxClusterNameLength         = 32
	MariaDBStandaloneReplicas           = 1
	MariaDBDefaultClusterSize           = 3
	MariaDBDataMountPath                = "/var/lib/mysql"
	MariaDBDataLostFoundPath            = MariaDBDataMountPath + "lost+found"
	MariaDBInitDBVolumeName             = "initial-script"
	MariaDBInitDBMountPath              = "/docker-entrypoint-initdb.d"
	MariaDBCustomConfigMountPath        = "/etc/mysql/conf.d/"
	MariaDBClusterCustomConfigMountPath = "/etc/mysql/custom.conf.d/"
	MariaDBCustomConfigVolumeName       = "custom-config"
	MariaDBTLSConfigCustom              = "custom"
	MariaDBInitContainerName            = "mariadb-init"
	MariaDBCoordinatorContainerName     = "md-coordinator"
	MariaDBRunScriptVolumeName          = "run-script"
	MariaDBRunScriptVolumeMountPath     = "/run-script"
	MariaDBInitScriptVolumeName         = "init-scripts"
	MariaDBInitScriptVolumeMountPath    = "/scripts"

	// =========================== PostgreSQL Constants ============================
	PostgresDatabasePortName         = "db"
	PostgresPrimaryServicePortName   = "primary"
	PostgresStandbyServicePortName   = "standby"
	PostgresDatabasePort             = 5432
	PostgresPodPrimary               = "primary"
	PostgresPodStandby               = "standby"
	EnvPostgresUser                  = "POSTGRES_USER"
	EnvPostgresPassword              = "POSTGRES_PASSWORD"
	PostgresCoordinatorContainerName = "pg-coordinator"
	PostgresCoordinatorPort          = 2380
	PostgresCoordinatorPortName      = "coordinator"

	PostgresCoordinatorClientPort     = 2379
	PostgresCoordinatorClientPortName = "coordinatclient"

	PostgresRunScriptMountPath  = "/run_scripts"
	PostgresRunScriptVolumeName = "scripts"

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

	// =========================== ProxySQL Constants ============================
	LabelProxySQLName        = ProxySQLKey + "/name"
	LabelProxySQLLoadBalance = ProxySQLKey + "/load-balance"

	ProxySQLDatabasePort           = 6033
	ProxySQLDatabasePortName       = "db"
	ProxySQLPrimaryServicePortName = "db"
	ProxySQLAdminPort              = 6032
	ProxySQLAdminPortName          = "admin"
	ProxySQLDataMountPath          = "/var/lib/proxysql"
	ProxySQLCustomConfigMountPath  = "/etc/custom-config"
	// =========================== Redis Constants ============================
	RedisConfigKey = "redis.conf" // RedisConfigKey is going to create for the customize redis configuration
	//DefaultConfigKey is going to create for the default redis configuration
	DefaultConfigKey            = "default.conf"
	RedisShardKey               = RedisKey + "/shard"
	RedisDatabasePortName       = "db"
	RedisPrimaryServicePortName = "primary"
	RedisDatabasePort           = 6379
	RedisGossipPortName         = "gossip"
	RedisGossipPort             = 16379
	RedisSentinelPortName       = "sentinel"
	RedisScriptVolumeName       = "script-vol"
	RedisScriptVolumePath       = "/scripts"
	RedisSentinelPort           = 26379

	RedisKeyFileSecretSuffix = "key"
	RedisPEMSecretSuffix     = "pem"
	RedisRootUsername        = "root"
	EnvRedisUser             = "USERNAME"
	EnvRedisPassword         = "REDISCLI_AUTH"

	// =========================== PgBouncer Constants ============================
	PgBouncerUpstreamServerCA       = "upstream-server-ca.crt"
	PgBouncerDatabasePortName       = "db"
	PgBouncerPrimaryServicePortName = "primary"
	PgBouncerDatabasePort           = 5432
	PgBouncerConfigFile             = "pgbouncer.ini"
	PgBouncerAdminUsername          = "kubedb"
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

	// Condition reasons
	DataRestoreStartedByExternalInitializer = "DataRestoreStartedByExternalInitializer"
	DatabaseSuccessfullyRestored            = "SuccessfullyDataRestored"
	FailedToRestoreData                     = "FailedToRestoreData"
	AllReplicasAreReady                     = "AllReplicasReady"
	SomeReplicasAreNotReady                 = "SomeReplicasNotReady"
	DatabaseAcceptingConnectionRequest      = "DatabaseAcceptingConnectionRequest"
	DatabaseNotAcceptingConnectionRequest   = "DatabaseNotAcceptingConnectionRequest"
	ReadinessCheckSucceeded                 = "ReadinessCheckSucceeded"
	ReadinessCheckFailed                    = "ReadinessCheckFailed"
	DatabaseProvisioningStartedSuccessfully = "DatabaseProvisioningStartedSuccessfully"
	DatabaseSuccessfullyProvisioned         = "DatabaseSuccessfullyProvisioned"
	DatabaseHaltedSuccessfully              = "DatabaseHaltedSuccessfully"
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
)
