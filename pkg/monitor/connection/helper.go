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

package connection

import (
	"fmt"
	"log"

	api "kubedb.dev/apimachinery/apis/kubedb/v1"
)

const (
	cAdvisorMetric     = "container_cpu_usage_seconds_total"
	kubeletMetric      = "kubelet_running_pods"
	ksmMetric          = "kube_pod_status_phase"
	nodeExporterMetric = "node_cpu_seconds_total"
)

func getIdenticalMetrics(database, databaseName string) map[string]*metrics {
	queries := make(map[string]*metrics)
	queries["cAdvisor"] = &metrics{metric: cAdvisorMetric}
	queries["kubelet"] = &metrics{metric: kubeletMetric}
	queries["kube-state-metric"] = &metrics{metric: ksmMetric}
	queries["node-exporter"] = &metrics{metric: nodeExporterMetric}

	queries = getDBMetrics(database, databaseName, queries)

	return queries
}

func getDBMetrics(database, name string, queries map[string]*metrics) map[string]*metrics {
	label := "service"
	labelValue := fmt.Sprintf("%s-stats", name)
	switch database {
	case api.ResourceSingularElasticsearch:
		queries[database] = &metrics{
			metric:     "elasticsearch_clusterinfo_up",
			label:      label,
			labelValue: labelValue,
		}
	case api.ResourceSingularKafka:
		queries[database] = &metrics{
			metric:     "kafka_controller_kafkacontroller_activebrokercount",
			label:      label,
			labelValue: labelValue,
		}
	case api.ResourceSingularMariaDB:
		queries[database] = &metrics{
			metric:     "mysql_up",
			label:      label,
			labelValue: labelValue,
		}
	case api.ResourceSingularMongoDB:
		queries[database] = &metrics{
			metric:     "mongodb_up",
			label:      label,
			labelValue: labelValue,
		}
	case api.ResourceSingularMySQL:
		queries[database] = &metrics{
			metric:     "mysql_up",
			label:      label,
			labelValue: labelValue,
		}
	case api.ResourceSingularPerconaXtraDB:
		queries[database] = &metrics{
			metric:     "mysql_up",
			label:      label,
			labelValue: labelValue,
		}
	case api.ResourceSingularPostgres:
		queries[database] = &metrics{
			metric:     "pg_up",
			label:      label,
			labelValue: labelValue,
		}
	case api.ResourceSingularProxySQL:
		queries[database] = &metrics{
			metric:     "proxysql_uptime_seconds_total",
			label:      label,
			labelValue: labelValue,
		}
	case api.ResourceSingularRedis:
		queries[database] = &metrics{
			metric:     "redis_up",
			label:      label,
			labelValue: labelValue,
		}
	default:
		log.Fatal("database invalid!")
	}

	// Panopticon
	queries["panopticon"] = &metrics{
		metric:     fmt.Sprintf("kubedb_com_%s_info", database),
		label:      database,
		labelValue: name,
	}
	return queries
}
