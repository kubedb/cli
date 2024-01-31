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
	"fmt"
	"strings"
)

func ConvertedResource(resource string) string {
	// standardizing the resource name
	res := strings.ToLower(resource)
	switch res {
	case "es", "elasticsearch", "elasticsearches":
		res = "elasticsearch"
	case "kf", "kafka", "kafkas":
		res = "kafka"
	case "md", "mariadb", "mariadbs":
		res = "mariadb"
	case "mg", "mongodb", "mongodbs":
		res = "mongodb"
	case "my", "mysql", "mysqls":
		res = "mysqls"
	case "px", "perconaxtradb", "perconaxtradbs":
		res = "perconaxtradb"
	case "pg", "postgres", "postgreses":
		res = "postgres"
	case "prx", "proxysql", "proxysqls":
		res = "proxysql"
	case "rd", "redis", "redises":
		res = "redis"
	default:
		fmt.Printf("%s is not a valid resource type \n", resource)
	}
	return res
}
