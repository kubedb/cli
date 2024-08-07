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

package resumer

import (
	"errors"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Resumer interface {
	Resume(string, string) (bool, error) // returns true if backupconfiguration is resumed
}

func NewResumer(restClientGetter genericclioptions.RESTClientGetter, mapping *meta.RESTMapping, onlyDb, onlyBackup, onlyArchiver bool) (Resumer, error) {
	clientConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	if mapping == nil {
		return nil, errors.New("mapping is empty")
	}

	switch mapping.GroupVersionKind.Kind {
	case dbapi.ResourceKindElasticsearch:
		return NewElasticsearchResumer(clientConfig, onlyDb, onlyBackup)
	case dbapi.ResourceKindMongoDB:
		return NewMongoDBResumer(clientConfig, onlyDb, onlyBackup, onlyArchiver)
	case dbapi.ResourceKindMySQL:
		return NewMySQLResumer(clientConfig, onlyDb, onlyBackup, onlyArchiver)
	case dbapi.ResourceKindMariaDB:
		return NewMariaDBResumer(clientConfig, onlyDb, onlyBackup, onlyArchiver)
	case dbapi.ResourceKindPostgres:
		return NewPostgresResumer(clientConfig, onlyDb, onlyBackup, onlyArchiver)
	case dbapi.ResourceKindRedis:
		return NewRedisResumer(clientConfig, onlyDb, onlyBackup)
	default:
		return nil, errors.New("unknown object kind")
	}
}
