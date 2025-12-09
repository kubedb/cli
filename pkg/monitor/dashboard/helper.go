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
	"io"
	"log"
	"net/http"
	"os"

	kapi "kubedb.dev/apimachinery/apis/kafka/v1alpha1"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	olddbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"kubedb.dev/cli/pkg/monitor"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func getDB(f cmdutil.Factory, resource, ns, name string) (*unstructured.Unstructured, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	if resource == kapi.ResourcePluralConnectCluster {
		kGvk := kapi.SchemeGroupVersion
		kRes := schema.GroupVersionResource{Group: kGvk.Group, Version: kGvk.Version, Resource: resource}
		return dc.Resource(kRes).Namespace(ns).Get(context.TODO(), name, metav1.GetOptions{})
	}

	gvk := dbapi.SchemeGroupVersion
	if monitor.IsOldAPI(resource) {
		gvk = olddbapi.SchemeGroupVersion
	}

	dbRes := schema.GroupVersionResource{Group: gvk.Group, Version: gvk.Version, Resource: resource}
	return dc.Resource(dbRes).Namespace(ns).Get(context.TODO(), name, metav1.GetOptions{})
}

func getURL(branch, database, dashboard string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/appscode/grafana-dashboards/%s/%s/%s.json", branch, database, dashboard)
}

func getDashboardFromURL(url string) map[string]any {
	var dashboardData map[string]any
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
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

func getDashboardFromFile(file string) map[string]any {
	body, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("Error on ReadFile:", err)
	}
	var dashboardData map[string]any
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

func ignoreModeSpecificExpressions(unknown map[string]*missingOpts, database string, db *unstructured.Unstructured) map[string]*missingOpts {
	if database == dbapi.ResourceSingularMongoDB {
		return ignoreMongoDBModeSpecificExpressions(unknown, db)
	}
	if database == dbapi.ResourceSingularElasticsearch {
		return ignoreElasticsearchModeSpecificExpressions(unknown, db)
	}
	return unknown
}
