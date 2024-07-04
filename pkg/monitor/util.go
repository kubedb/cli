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

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
)

func ConvertedResourceToPlural(resource string) string {
	// standardizing the resource name
	res := strings.ToLower(resource)
	switch res {
	case dbapi.ResourceCodeElasticsearch, dbapi.ResourcePluralElasticsearch, dbapi.ResourceSingularElasticsearch:
		res = dbapi.ResourcePluralElasticsearch
	case dbapi.ResourceCodeKafka, dbapi.ResourcePluralKafka, dbapi.ResourceSingularKafka:
		res = dbapi.ResourcePluralKafka
	case dbapi.ResourceCodeMariaDB, dbapi.ResourcePluralMariaDB, dbapi.ResourceSingularMariaDB:
		res = dbapi.ResourcePluralMariaDB
	case dbapi.ResourceCodeMongoDB, dbapi.ResourcePluralMongoDB, dbapi.ResourceSingularMongoDB:
		res = dbapi.ResourcePluralMongoDB
	case dbapi.ResourceCodeMySQL, dbapi.ResourcePluralMySQL, dbapi.ResourceSingularMySQL:
		res = dbapi.ResourcePluralMySQL
	case dbapi.ResourceCodePerconaXtraDB, dbapi.ResourcePluralPerconaXtraDB, dbapi.ResourceSingularPerconaXtraDB:
		res = dbapi.ResourcePluralPerconaXtraDB
	case dbapi.ResourceCodePostgres, dbapi.ResourcePluralPostgres, dbapi.ResourceSingularPostgres:
		res = dbapi.ResourcePluralPostgres
	case dbapi.ResourceCodeProxySQL, dbapi.ResourcePluralProxySQL, dbapi.ResourceSingularProxySQL:
		res = dbapi.ResourcePluralProxySQL
	case dbapi.ResourceCodeRedis, dbapi.ResourcePluralRedis, dbapi.ResourceSingularRedis:
		res = dbapi.ResourcePluralRedis
	default:
		log.Fatalf("%s is not a valid resource type \n", resource)
	}
	return res
}

func ConvertedResourceToSingular(resource string) string {
	// standardizing the resource name
	res := strings.ToLower(resource)
	switch res {
	case dbapi.ResourceCodeElasticsearch, dbapi.ResourcePluralElasticsearch, dbapi.ResourceSingularElasticsearch:
		res = dbapi.ResourceSingularElasticsearch
	case dbapi.ResourceCodeKafka, dbapi.ResourcePluralKafka, dbapi.ResourceSingularKafka:
		res = dbapi.ResourceSingularKafka
	case dbapi.ResourceCodeMariaDB, dbapi.ResourcePluralMariaDB, dbapi.ResourceSingularMariaDB:
		res = dbapi.ResourceSingularMariaDB
	case dbapi.ResourceCodeMongoDB, dbapi.ResourcePluralMongoDB, dbapi.ResourceSingularMongoDB:
		res = dbapi.ResourceSingularMongoDB
	case dbapi.ResourceCodeMySQL, dbapi.ResourcePluralMySQL, dbapi.ResourceSingularMySQL:
		res = dbapi.ResourceSingularMySQL
	case dbapi.ResourceCodePerconaXtraDB, dbapi.ResourcePluralPerconaXtraDB, dbapi.ResourceSingularPerconaXtraDB:
		res = dbapi.ResourceSingularPerconaXtraDB
	case dbapi.ResourceCodePostgres, dbapi.ResourcePluralPostgres, dbapi.ResourceSingularPostgres:
		res = dbapi.ResourceSingularPostgres
	case dbapi.ResourceCodeProxySQL, dbapi.ResourcePluralProxySQL, dbapi.ResourceSingularProxySQL:
		res = dbapi.ResourceSingularProxySQL
	case dbapi.ResourceCodeRedis, dbapi.ResourcePluralRedis, dbapi.ResourceSingularRedis:
		res = dbapi.ResourceSingularRedis
	default:
		log.Fatalf("%s is not a valid resource type \n", resource)
	}
	return res
}
