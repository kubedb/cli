/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package credentials

import (
	"errors"
	"fmt"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type ShowCred interface {
	GetCred(string, string) (map[string][]byte, error)
}

func NewShowCredentials(restClientGetter genericclioptions.RESTClientGetter, mapping *meta.RESTMapping) (ShowCred, error) {
	clientConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	if mapping == nil {
		return nil, errors.New("mapping is empty")
	}

	switch mapping.GroupVersionKind.Kind {
	case api.ResourceKindElasticsearch:
		return NewElasticsearchShowCred(clientConfig)
	case api.ResourceKindMongoDB:
		return NewMongoDBShowCred(clientConfig)
	case api.ResourceKindMySQL:
		return NewMySQLShowCred(clientConfig)
	case api.ResourceKindRedis:
		return NewRedisShowCred(clientConfig)
	case api.ResourceKindMariaDB:
		return NewMariaDBShowCred(clientConfig)
	case api.ResourceKindPostgres:
		return NewPostgresShowCred(clientConfig)
	default:
		return nil, fmt.Errorf("unsupporterd kind %s", mapping.GroupVersionKind.Kind)
	}
}
