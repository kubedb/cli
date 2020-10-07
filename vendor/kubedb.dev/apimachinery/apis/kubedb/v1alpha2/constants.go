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

import "kubedb.dev/apimachinery/apis/kubedb"

const (
	// Deprecated
	DatabaseNamePrefix = "kubedb"

	KubeDBOrganization = "kubedb"

	LabelDatabaseKind = kubedb.GroupName + "/kind"
	LabelDatabaseName = kubedb.GroupName + "/name"
	LabelRole         = kubedb.GroupName + "/role"

	ComponentDatabase     = "database"
	RoleStats             = "stats"
	DefaultStatsPath      = "/metrics"
	DefaultPasswordLength = 16

	PostgresKey      = ResourceSingularPostgres + "." + kubedb.GroupName
	ElasticsearchKey = ResourceSingularElasticsearch + "." + kubedb.GroupName
	MySQLKey         = ResourceSingularMySQL + "." + kubedb.GroupName
	PerconaXtraDBKey = ResourceSingularPerconaXtraDB + "." + kubedb.GroupName
	MongoDBKey       = ResourceSingularMongoDB + "." + kubedb.GroupName
	RedisKey         = ResourceSingularRedis + "." + kubedb.GroupName
	MemcachedKey     = ResourceSingularMemcached + "." + kubedb.GroupName
	EtcdKey          = ResourceSingularEtcd + "." + kubedb.GroupName
	ProxySQLKey      = ResourceSingularProxySQL + "." + kubedb.GroupName

	ElasticsearchRestPort                        = 9200
	ElasticsearchRestPortName                    = "http"
	ElasticsearchTransportPort                   = 9300
	ElasticsearchTransportPortName               = "transport"
	ElasticsearchMetricsPort                     = 9600
	ElasticsearchMetricsPortName                 = "metrics"
	ElasticsearchIngestNodePrefix                = "ingest"
	ElasticsearchDataNodePrefix                  = "data"
	ElasticsearchMasterNodePrefix                = "master"
	ElasticsearchNodeRoleMaster                  = "node.role.master"
	ElasticsearchNodeRoleIngest                  = "node.role.ingest"
	ElasticsearchNodeRoleData                    = "node.role.data"
	ElasticsearchNodeRoleSet                     = "set"
	ElasticsearchConfigDir                       = "/usr/share/elasticsearch/config"
	ElasticsearchTempConfigDir                   = "/elasticsearch/temp-config"
	ElasticsearchCustomConfigDir                 = "/elasticsearch/custom-config"
	ElasticsearchDataDir                         = "/usr/share/elasticsearch/data"
	ElasticsearchOpendistroSecurityConfigDir     = "/usr/share/elasticsearch/plugins/opendistro_security/securityconfig"
	ElasticsearchSearchGuardSecurityConfigDir    = "/usr/share/elasticsearch/plugins/search-guard-%v/sgconfig"
	ElasticsearchOpendistroReadallMonitorRole    = "readall_and_monitor"
	ElasticsearchSearchGuardReadallMonitorRoleV7 = "SGS_READALL_AND_MONITOR"
	ElasticsearchSearchGuardReadallMonitorRoleV6 = "sg_readall_and_monitor"

	// Ref:
	//	- https://www.elastic.co/guide/en/elasticsearch/reference/7.6/heap-size.html#heap-size
	//	- no more than 50% of your physical RAM
	//	- no more than 32GB that the JVM uses for compressed object pointers (compressed oops)
	//	- no more than 26GB for zero-based compressed oops;
	// 26 GB is safe on most systems
	ElasticsearchMaxHeapSize = 26 * 1024 * 1024 * 1024
	// 128MB
	ElasticsearchMinHeapSize = 128 * 1024 * 1024

	MongoDBShardPort           = 27017
	MongoDBConfigdbPort        = 27017
	MongoDBMongosPort          = 27017
	MongoDBKeyFileSecretSuffix = "key"
	MongoDBRootUsername        = "root"

	MySQLMetricsExporterConfigSecretSuffix = "metrics-exporter-config"
	MySQLNodePort                          = 3306
	MySQLGroupComPort                      = 33060
	MySQLMaxGroupMembers                   = 9
	// The recommended MySQL server version for group replication (GR)
	MySQLGRRecommendedVersion       = "5.7.25"
	MySQLDefaultGroupSize           = 3
	MySQLDefaultBaseServerID  int64 = 1
	// The server id for each group member must be unique and in the range [1, 2^32 - 1]
	// And the maximum group size is 9. So MySQLMaxBaseServerID is the maximum safe value
	// for BaseServerID calculated as max MySQL server_id value - max Replication Group size.
	// xref: https://dev.mysql.com/doc/refman/5.7/en/replication-options.html
	MySQLMaxBaseServerID int64 = ((1 << 32) - 1) - 9
	MySQLRootUserName          = "MYSQL_ROOT_USERNAME"
	MySQLRootPassword          = "MYSQL_ROOT_PASSWORD"
	MySQLName                  = "MYSQL_NAME"

	MySQLContainerReplicationModeDetectorName = "replication-mode-detector"
	MySQLPodPrimary                           = "primary"
	MySQLPodSecondary                         = "secondary"
	MySQLLabelRole                            = MySQLKey + "/role"

	PerconaXtraDBClusterRecommendedVersion    = "5.7"
	PerconaXtraDBMaxClusterNameLength         = 32
	PerconaXtraDBStandaloneReplicas           = 1
	PerconaXtraDBDefaultClusterSize           = 3
	PerconaXtraDBDataMountPath                = "/var/lib/mysql"
	PerconaXtraDBDataLostFoundPath            = PerconaXtraDBDataMountPath + "lost+found"
	PerconaXtraDBInitDBMountPath              = "/docker-entrypoint-initdb.d"
	PerconaXtraDBCustomConfigMountPath        = "/etc/percona-server.conf.d/"
	PerconaXtraDBClusterCustomConfigMountPath = "/etc/percona-xtradb-cluster.conf.d/"

	LabelProxySQLName        = ProxySQLKey + "/name"
	LabelProxySQLLoadBalance = ProxySQLKey + "/load-balance"

	ProxySQLMySQLNodePort         = 6033
	ProxySQLAdminPort             = 6032
	ProxySQLAdminPortName         = "admin"
	ProxySQLDataMountPath         = "/var/lib/proxysql"
	ProxySQLCustomConfigMountPath = "/etc/custom-config"

	RedisShardKey   = RedisKey + "/shard"
	RedisNodePort   = 6379
	RedisGossipPort = 16379

	RedisKeyFileSecretSuffix = "key"
	RedisPEMSecretSuffix     = "pem"
	RedisRootUsername        = "root"

	PgBouncerUpstreamServerCA = "upstream-server-ca.crt"

	ContainerExporterName = "exporter"
	LocalHost             = "localhost"
	LocalHostIP           = "127.0.0.1"

	DBCustomConfigName = "custom-config"
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
	// used for Databases that are paused
	DatabasePaused = "Paused"

	// Condition reasons
	DataRestoreStartedByExternalInitializer = "DataRestoreStartedByExternalInitializer"
	DatabaseSuccessfullyRestored            = "SuccessfullyDataRestored"
	FailedToRestoreData                     = "FailedToRestoreData"
	AllReplicasAreReady                     = "AllReplicasReady"
	SomeReplicasAreNotReady                 = "SomeReplicasNotReady"
)

// Resource kind related constants
const (
	ResourceKindStatefulSet = "StatefulSet"
)
