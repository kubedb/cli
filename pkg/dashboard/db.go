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

package dashboard

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"kubedb.dev/cli/pkg/lib"

	"github.com/prometheus/common/model"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type queryInformation struct {
	metric     string
	labelNames []string
}
type PromSvc struct {
	Name      string
	Namespace string
	Port      int
}
type unknownLabel struct {
	metric    string
	labelName string
}

func Run(f cmdutil.Factory, args []string, branch string, prom PromSvc) {
	if len(args) < 2 {
		log.Fatal("Enter database and grafana dashboard name as argument")
	}

	database := args[0]
	dashboard := args[1]

	url := getURL(branch, database, dashboard)

	dashboardData := getDashboard(url)

	queries := parseAllExpressions(dashboardData)

	config, err := f.ToRESTConfig()
	if err != nil {
		log.Fatal(err)
	}
	// Port forwarding cluster prometheus service for that grafana dashboard's prom datasource.
	tunnel, err := lib.TunnelToDBService(config, prom.Name, prom.Namespace, prom.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer tunnel.Close()

	promClient := getPromClient(strconv.Itoa(tunnel.Local))

	var unknownMetrics []string
	var unknownLabels []unknownLabel

	for _, query := range queries {
		metricName := query.metric
		for _, labelKey := range query.labelNames {

			endTime := time.Now()

			result, _, err := promClient.Query(context.TODO(), metricName, endTime)
			if err != nil {
				log.Fatal("Error querying Prometheus:", err, " metric: ", metricName)
			}

			matrix := result.(model.Vector)
			if len(matrix) > 0 {
				// Check if the label exists for any result in the matrix
				labelExists := false

				for _, sample := range matrix {
					if sample.Metric != nil {
						if _, ok := sample.Metric[model.LabelName(labelKey)]; ok {
							labelExists = true
							break
						}
					}
				}

				if !labelExists {
					unknownLabels = uniqueAppend(unknownLabels, unknownLabel{
						metric:    metricName,
						labelName: labelKey,
					})
				}
			} else {
				unknownMetrics = uniqueAppend(unknownMetrics, metricName)
			}
		}
	}
	if len(unknownMetrics) > 0 {
		fmt.Printf("List of unknown metrics:\n%s\n", strings.Join(unknownMetrics, "\n"))
	}
	if len(unknownLabels) > 0 {
		fmt.Println("List of unknown labels:")
		for _, unknown := range unknownLabels {
			fmt.Printf(`Metric: "%s" Label: "%s"\n`, unknown.metric, unknown.labelName)
		}
	}
	if len(unknownMetrics) == 0 && len(unknownLabels) == 0 {
		fmt.Println("All metrics found")
	}
}
