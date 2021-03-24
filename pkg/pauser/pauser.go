package pauser

import (
	"errors"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Pauser interface {
	Pause(string, string) error
}

func NewPauser(restClientGetter genericclioptions.RESTClientGetter, mapping *meta.RESTMapping) (Pauser, error) {
	clientConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	if mapping == nil {
		return nil, errors.New("mapping is empty")
	}

	switch mapping.GroupVersionKind.Kind {
	case api.ResourceKindElasticsearch:
		return NewElasticsearchPauser(clientConfig)
	default:
		return nil, errors.New("unknown kind")
	}

	return nil, nil
}
