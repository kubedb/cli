package docker

const (
	ImageOperator          = "aerokite/operator"
	ImagePostgresOperator  = "aerokite/pg-operator"
	ImagePostgres          = "aerokite/postgres"
	ImageMySQLOperator     = "aerokite/mysql-operator"
	ImageMySQL             = "library/mysql"
	ImageElasticOperator   = "aerokite/es-operator"
	ImageElasticsearch     = "aerokite/elasticsearch"
	ImageElasticdump       = "aerokite/elasticdump"
	ImageMongoDBOperator   = "aerokite/mongodb-operator"
	ImageMongoDB           = "library/mongo"
	ImageRedisOperator     = "aerokite/redis-operator"
	ImageRedis             = "library/redis"
	ImageMemcachedOperator = "aerokite/mc-operator"
	ImageMemcached         = "library/memcached"
)

const (
	OperatorName       = "kubedb-operator"
	OperatorContainer  = "operator"
	OperatorPortName   = "web"
	OperatorPortNumber = 8080
)
