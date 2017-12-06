package decoder

import (
	"fmt"

	"github.com/appscode/kutil/meta"
	"github.com/ghodss/yaml"
	tapi "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	_ "github.com/spf13/cobra/doc"
	"github.com/the-redback/go-oneliners"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

func Decode(kind string, data []byte) (runtime.Object, error) {
	switch kind {
	case tapi.ResourceKindElasticsearch:
		var elastic *tapi.Elasticsearch
		if err := yaml.Unmarshal(data, &elastic); err != nil {
			return nil, err
		}
		return elastic, nil
	case tapi.ResourceKindPostgres:
		//var postgres *tapi.Postgres
		dataPost, err := meta.UnmarshalToYAML(data, tapi.SchemeGroupVersion)
		if err != nil {
			return nil, err
		}
		return dataPost, nil
		//if err := yaml.Unmarshal(data, &postgres); err != nil {
		//	return nil, err
		//}
		//return postgres, nil
	case tapi.ResourceKindMySQL:
		var mysql *tapi.MySQL
		if err := yaml.Unmarshal(data, &mysql); err != nil {
			return nil, err
		}
		return mysql, nil
	case tapi.ResourceKindMongoDB:
		var mongodb *tapi.MongoDB
		if err := yaml.Unmarshal(data, &mongodb); err != nil {
			return nil, err
		}
		return mongodb, nil
	case tapi.ResourceKindRedis:
		var redis *tapi.Redis
		if err := yaml.Unmarshal(data, &redis); err != nil {
			return nil, err
		}
		return redis, nil
	case tapi.ResourceKindMemcached:
		var memcached *tapi.Memcached
		if err := yaml.Unmarshal(data, &memcached); err != nil {
			return nil, err
		}
		return memcached, nil
	case tapi.ResourceKindSnapshot:
		var snapshot *tapi.Snapshot
		if err := yaml.Unmarshal(data, &snapshot); err != nil {
			return nil, err
		}
		return snapshot, nil
	case tapi.ResourceKindDormantDatabase:
		var deletedDb *tapi.DormantDatabase
		if err := yaml.Unmarshal(data, &deletedDb); err != nil {
			return nil, err
		}
		return deletedDb, nil
	}

	return nil, fmt.Errorf(`Invalid kind: "%v"`, kind)
}

func unmarshalToYAML(data []byte, gv schema.GroupVersion) (runtime.Object, error) {

	mediaType := "application/yaml"
	info, ok := runtime.SerializerInfoForMediaType(clientsetscheme.Codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return nil, fmt.Errorf("unsupported media type %q", mediaType)
	}
	oneliners.PrettyJson(info, "info")

	decoder := clientsetscheme.Codecs.DecoderToVersion(info.Serializer, gv)
	dataaa, err := runtime.Decode(decoder, data)

	return dataaa, err
}
