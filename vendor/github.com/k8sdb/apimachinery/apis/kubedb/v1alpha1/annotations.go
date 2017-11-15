package v1alpha1

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"
	LabelJobType      = GenericKey + "/job-type"

	PostgresKey             = ResourceTypePostgres + "." + GenericKey
	PostgresDatabaseVersion = PostgresKey + "/version"

	ElasticsearchKey             = ResourceTypeElasticsearch + "." + GenericKey
	ElasticsearchDatabaseVersion = ElasticsearchKey + "/version"

	MySQLKey             = ResourceTypeMySQL + "." + GenericKey
	MySQLDatabaseVersion = MySQLKey + "/version"

	SnapshotKey         = ResourceTypeSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	PostgresInitSpec      = PostgresKey + "/init"
	ElasticsearchInitSpec = ElasticsearchKey + "/init"
	MySQLInitSpec         = MySQLKey + "/init"

	PostgresIgnore      = PostgresKey + "/ignore"
	ElasticsearchIgnore = ElasticsearchKey + "/ignore"
	MySQLIgnore         = MySQLKey + "/ignore"

	AgentCoreosPrometheus        = "coreos-prometheus-operator"
	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "http"
)
