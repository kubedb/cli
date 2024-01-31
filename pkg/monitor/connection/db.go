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
	"context"
	"fmt"
	"log"
	"time"

	"kubedb.dev/cli/pkg/monitor"

	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// kubectl dba monitor check-connection mongodb -n demo sample_mg -> all connection check and report

type metrics struct {
	metric     string
	label      string
	labelValue string
}

func Run(f cmdutil.Factory, args []string, prom monitor.PromSvc) {
	if len(args) < 2 {
		log.Fatal("Enter database and specific database name as argument")
	}

	database := monitor.ConvertedResource(args[0])
	databaseName := args[1]
	namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		klog.Error(err, "failed to get current namespace")
	}

	config, err := f.ToRESTConfig()
	if err != nil {
		log.Fatal(err)
	}

	promClient, tunnel := monitor.GetPromClientAndTunnel(config, prom)
	defer tunnel.Close()

	queries := getIdenticalMetrics(database, databaseName)
	var notFound []string

	for target, query := range queries {
		metricName := query.metric
		endTime := time.Now()

		result, _, err := promClient.Query(context.TODO(), metricName, endTime)
		if err != nil {
			log.Fatal("Error querying Prometheus:", err, " metric: ", metricName)
		}

		matrix := result.(model.Vector)

		if len(matrix) > 0 {
			if query.label != "" {
				// DB and panopticon related metrics.So we need to check label also
				// Check if the label exists for any result in the matrix
				exist := false
				for _, sample := range matrix {
					if sample.Metric != nil {
						if labelVal, ok := sample.Metric[model.LabelName(query.label)]; ok {
							if string(labelVal) == query.labelValue {
								if namespaceVal, ok := sample.Metric[model.LabelName("namespace")]; ok {
									if string(namespaceVal) == namespace {
										exist = true
									}
								}
							}
						}
					}
				}
				if !exist {
					notFound = append(notFound, target)
				}
			}
		} else {
			notFound = append(notFound, target)
		}
	}

	if len(notFound) == 0 {
		fmt.Printf("All monitoring connection established successfully for %s : %s/%s\n", database, namespace, databaseName)
	} else {
		for _, target := range notFound {
			fmt.Printf("%s monitoring connection not found for %s : %s/%s\n", target, database, namespace, databaseName)
		}
	}
}
