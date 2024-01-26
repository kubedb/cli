package connection

import (
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"log"
)

const (
	cAdvisorMetric     = "container_cpu_usage_seconds_total"
	kubeletMetric      = "kubelet_active_pods"
	ksmMetric          = "kube_pod_status_phase"
	nodeExporterMetric = "node_cpu_seconds_total"
)

func getIdenticalMetrics(database, databaseName, namespace string) []metrics {
	identicalMetrics := []metrics{
		{metric: cAdvisorMetric},
		{metric: kubeletMetric},
		{metric: ksmMetric},
		{metric: nodeExporterMetric},
	}
	identicalMetrics = append(identicalMetrics, getDBMetrics(database, databaseName, namespace)...)
	return identicalMetrics
}
func getDBMetrics(database, name, namespace string) []metrics {
	var dbMetric []metrics
	label := fmt.Sprintf("%s-stats", name)
	metric := fmt.Sprintf("kubedb_com_%s", database)
	switch database {
	case "mongodb":
		dbMetric = append(dbMetric, metrics{
			metric: "",
			label:  "",
		})
	case "postgres":
	case "mysql":
	case "redis":
	case "mariadb":
	case "proxysql":
	case "elasticsearch":
	case "perconaxtradb":
	case "kafka":
	default:
		log.Fatal("database invalid")
	}
	panopticon := getPanopticonMetric()
}
func getPanopticonMetric() {

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
