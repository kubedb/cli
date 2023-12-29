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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kubedb.dev/cli/pkg/lib"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
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

func Run(f cmdutil.Factory, args []string, branch string, prom PromSvc) {
	if len(args) < 2 {
		log.Fatal("Enter database and grafana dashboard name as argument")
	}

	database := args[0]
	dashboard := args[1]

	url := getURL(branch, database, dashboard)

	dashboardData := getDashboard(url)

	queries := getQueryInformation(dashboardData)

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

	var unknownMetrics, unknownLabels []string

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
					unknownLabels = uniqueAppend(unknownLabels, fmt.Sprintf(`label: "%s" metric: "%s"`, labelKey, metricName))
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
		fmt.Printf("List of unknown labels:\n%s", strings.Join(unknownLabels, "\n"))
	}
	if len(unknownMetrics) == 0 && len(unknownLabels) == 0 {
		fmt.Println("All metrics found")
	}
}

func getURL(branch, database, dashboard string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/appscode/grafana-dashboards/%s/%s/%s.json", branch, database, dashboard)
}

func getDashboard(url string) map[string]interface{} {
	var dashboardData map[string]interface{}
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error fetching url. status : %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading JSON file: ", err)
	}

	err = json.Unmarshal(body, &dashboardData)
	if err != nil {
		log.Fatal("Error unmarshalling JSON data:", err)
	}
	return dashboardData
}

// Steps:
// - if current character is '{'
//   - extract metric name by matching metric regex
//   - get label selector substring inside { }
//   - get label name from this substring by matching label regex
//   - move i to its closing bracket position.
func getMetricAndLabels(query string) []queryInformation {
	var queries []queryInformation
	for i := 0; i < len(query); i++ {
		if query[i] == '{' {
			j := i
			for {
				if j-1 < 0 || (!matchMetricRegex(rune(query[j-1]))) {
					break
				}
				j--
			}
			metric := query[j:i]
			labelSelector, closingPosition := substringInsideLabelSelector(query, i)
			labelNames := getLabelNames(labelSelector)
			queries = append(queries, queryInformation{
				metric:     metric,
				labelNames: labelNames,
			})
			i = closingPosition
		}
	}
	return queries
}

func getPromClient(localPort string) v1.API {
	prometheusURL := fmt.Sprintf("http://localhost:%s/", localPort)

	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})
	if err != nil {
		log.Fatal("Error creating Prometheus client:", err)
	}

	// Create a new Prometheus API client
	return v1.NewAPI(client)
}

func getQueryInformation(dashboardData map[string]interface{}) []queryInformation {
	var queries []queryInformation
	if panels, ok := dashboardData["panels"].([]interface{}); ok {
		for _, panel := range panels {
			if targets, ok := panel.(map[string]interface{})["targets"].([]interface{}); ok {
				for _, target := range targets {
					if expr, ok := target.(map[string]interface{})["expr"]; ok {
						if expr != "" {
							query := expr.(string)
							queries = append(queries, getMetricAndLabels(query)...)
						}
					}
				}
			}
		}
	}
	return queries
}
