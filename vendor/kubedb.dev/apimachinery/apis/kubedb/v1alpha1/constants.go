/*
Copyright The KubeDB Authors.

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

	PostgresKey      = ResourceSingularPostgres + "." + GenericKey
	ElasticsearchKey = ResourceSingularElasticsearch + "." + GenericKey
	MySQLKey         = ResourceSingularMySQL + "." + GenericKey
	PerconaXtraDBKey = ResourceSingularPerconaXtraDB + "." + GenericKey
	MongoDBKey       = ResourceSingularMongoDB + "." + GenericKey
	RedisKey         = ResourceSingularRedis + "." + GenericKey
	MemcachedKey     = ResourceSingularMemcached + "." + GenericKey
	EtcdKey          = ResourceSingularEtcd + "." + GenericKey
	ProxySQLKey      = ResourceSingularProxySQL + "." + GenericKey

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

	MongoDBShardPort                  = 27017
	MongoDBConfigdbPort               = 27017
	MongoDBMongosPort                 = 27017
	MongoDBKeyFileSecretSuffix        = "-key"
	MongoDBExternalClientSecretSuffix = "-client-cert"
	MongoDBExporterClientSecretSuffix = "-exporter-cert"
	MongoDBServerSecretSuffix         = "-server-cert"
	MongoDBPEMSecretSuffix            = "-pem"
	MongoDBClientCertOrganization     = DatabaseNamePrefix + ":client"
	MongoDBCertificateCN              = "root"

	MySQLUserKey         = "username"
	MySQLPasswordKey     = "password"
	MySQLNodePort        = 3306
	MySQLGroupComPort    = 33060
	MySQLMaxGroupMembers = 9
	// The recommended MySQL server version for group replication (GR)
	MySQLGRRecommendedVersion       = "5.7.25"
	MySQLDefaultGroupSize           = 3
	MySQLDefaultBaseServerID  int64 = 1
	// The server id for each group member must be unique and in the range [1, 2^32 - 1]
	// And the maximum group size is 9. So MySQLMaxBaseServerID is the maximum safe value
	// for BaseServerID calculated as max MySQL server_id value - max Replication Group size.
	// xref: https://dev.mysql.com/doc/refman/5.7/en/replication-options.html
	MySQLMaxBaseServerID int64 = ((1 << 32) - 1) - 9

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

	ProxySQLUserKey               = "proxysqluser"
	ProxySQLPasswordKey           = "proxysqlpass"
	ProxySQLMySQLNodePort         = 6033
	ProxySQLAdminPort             = 6032
	ProxySQLAdminPortName         = "proxyadm"
	ProxySQLDataMountPath         = "/var/lib/proxysql"
	ProxySQLCustomConfigMountPath = "/etc/custom-config"

	RedisShardKey   = RedisKey + "/shard"
	RedisNodePort   = 6379
	RedisGossipPort = 16379

	PgBouncerServingClientSuffix      = "-serving-client-cert"
	PgBouncerExporterClientCertSuffix = "-exporter-cert"
	PgBouncerServingServerSuffix      = "-serving-server-cert"
	PgBouncerUpstreamServerCA         = "upstream-server-ca.crt"

	MySQLClientCertSuffix         = "client-cert"
	MySQLExporterClientCertSuffix = "exporter-cert"
	MySQLServerCertSuffix         = "server-cert"

	LocalHost   = "localhost"
	LocalHostIP = "127.0.0.1"
)
