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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading JSON file: ", err)
	}

	err = json.Unmarshal(body, &dashboardData)
	if err != nil {
		log.Fatal("Error unmarshalling JSON data:", err)
	}
	return dashboardData
}

func uniqueAppend(slice []string, valueToAdd string) []string {
	for _, existingValue := range slice {
		if existingValue == valueToAdd {
			return slice
		}
	}
	return append(slice, valueToAdd)
}
