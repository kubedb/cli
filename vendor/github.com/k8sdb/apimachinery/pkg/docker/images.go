package docker

const (
	ImageOperator          = "kubedb/operator"
	ImagePostgresOperator  = "kubedb/pg-operator"
	ImagePostgres          = "kubedb/postgres"
	ImageMySQLOperator     = "kubedb/mysql-operator"
	ImageMySQL             = "library/mysql"
	ImageElasticOperator   = "kubedb/es-operator"
	ImageElasticsearch     = "kubedb/elasticsearch"
	ImageElasticdump       = "kubedb/elasticdump"
	ImageMongoDBOperator   = "kubedb/mongodb-operator"
	ImageMongoDB           = "library/mongo"
	ImageRedisOperator     = "kubedb/redis-operator"
	ImageRedis             = "library/redis"
	ImageMemcachedOperator = "kubedb/mc-operator"
	ImageMemcached         = "library/memcached"
)

const (
	OperatorName       = "kubedb-operator"
	OperatorContainer  = "operator"
	OperatorPortName   = "web"
	OperatorPortNumber = 8080
)
