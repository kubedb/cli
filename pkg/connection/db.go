package connection

import (
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"kubedb.dev/cli/pkg/lib"
	"log"
	"strconv"
)

// kubectl dba monitor check-connection mongodb -n demo sample_mg -> all connection check and report

type PromSvc struct {
	Name      string
	Namespace string
	Port      int
}
type metrics struct {
	metric string
	label  string
}

func Run(f cmdutil.Factory, args []string, prom PromSvc) {
	if len(args) < 2 {
		log.Fatal("Enter database and specific database name as argument")
	}

	database := args[0]
	databaseName := args[1]
	namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		klog.Error(err, "failed to get current namespace")
	}

	identicalMetrics := []metrics{
		{metric: cAdvisorMetric},
		{metric: kubeletMetric},
		{metric: ksmMetric},
		{metric: nodeExporterMetric},
	}
	identicalMetrics = append(identicalMetrics, getDBMetrics(database, databaseName, namespace)...)

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

}
