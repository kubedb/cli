package resumer

import (
	"errors"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Resumer interface {
	Resume(string, string) error
}

func NewResumer(restClientGetter genericclioptions.RESTClientGetter, mapping *meta.RESTMapping) (Resumer, error) {
	clientConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	if mapping == nil {
		return nil, errors.New("mapping is empty")
	}

	switch mapping.GroupVersionKind.Kind {
	case api.ResourceKindElasticsearch:
		return NewElasticsearchResumer(clientConfig)
	case api.ResourceKindMongoDB:
		return NewMongoDBResumer(clientConfig)
	case api.ResourceKindMySQL:
		return NewMySQLResumer(clientConfig)
	case api.ResourceKindMariaDB:
		return NewMariaDBResumer(clientConfig)
	case api.ResourceKindPostgres:
		return NewPostgresResumer(clientConfig)
	case api.ResourceKindRedis:
		return NewRedisResumer(clientConfig)
	default:
		return nil, errors.New("unknown object kind")
	}
}
