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

package monitor

import (
	"log"
	"strings"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
)

func ConvertedResource(resource string) string {
	// standardizing the resource name
	res := strings.ToLower(resource)
	switch res {
	case api.ResourceCodeElasticsearch, api.ResourcePluralElasticsearch, api.ResourceSingularElasticsearch:
		res = api.ResourcePluralElasticsearch
	case api.ResourceCodeKafka, api.ResourcePluralKafka, api.ResourceSingularKafka:
		res = api.ResourcePluralKafka
	case api.ResourceCodeMariaDB, api.ResourcePluralMariaDB, api.ResourceSingularMariaDB:
		res = api.ResourcePluralMariaDB
	case api.ResourceCodeMongoDB, api.ResourcePluralMongoDB, api.ResourceSingularMongoDB:
		res = api.ResourcePluralMongoDB
	case api.ResourceCodeMySQL, api.ResourcePluralMySQL, api.ResourceSingularMySQL:
		res = api.ResourcePluralMySQL
	case api.ResourceCodePerconaXtraDB, api.ResourcePluralPerconaXtraDB, api.ResourceSingularPerconaXtraDB:
		res = api.ResourcePluralPerconaXtraDB
	case api.ResourceCodePostgres, api.ResourcePluralPostgres, api.ResourceSingularPostgres:
		res = api.ResourcePluralPostgres
	case api.ResourceCodeProxySQL, api.ResourcePluralProxySQL, api.ResourceSingularProxySQL:
		res = api.ResourcePluralProxySQL
	case api.ResourceCodeRedis, api.ResourcePluralRedis, api.ResourceSingularRedis:
		res = api.ResourcePluralRedis
	default:
		log.Fatalf("%s is not a valid resource type \n", resource)
	}
	return res
}
