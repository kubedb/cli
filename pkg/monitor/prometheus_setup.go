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
	"log"
	"strconv"

	"kubedb.dev/cli/pkg/lib"

	"github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"k8s.io/client-go/rest"
	"kmodules.xyz/client-go/tools/portforward"
)

type PromSvc struct {
	Name      string
	Namespace string
	Port      int
}

func GetPromClientAndTunnel(config *rest.Config, prom PromSvc) (promv1.API, *portforward.Tunnel) {
	tunnel, err := lib.TunnelToDBService(config, prom.Name, prom.Namespace, prom.Port)
	if err != nil {
		log.Fatal(err)
	}

	promClient := getPromClient(strconv.Itoa(tunnel.Local))
	return promClient, tunnel
}

func getPromClient(localPort string) promv1.API {
	prometheusURL := fmt.Sprintf("http://localhost:%s/", localPort)

	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})
	if err != nil {
		log.Fatal("Error creating Prometheus client:", err)
	}

	return promv1.NewAPI(client)
}
