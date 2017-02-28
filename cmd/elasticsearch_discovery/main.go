package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/appscode/log"
	logs "github.com/appscode/log/golog"
	"k8s.io/kubernetes/pkg/api"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	rest "k8s.io/kubernetes/pkg/client/restclient"
)

func flattenSubsets(subsets []api.EndpointSubset) []string {
	ips := []string{}
	for _, ss := range subsets {
		for _, addr := range ss.Addresses {
			ips = append(ips, fmt.Sprintf(`"%s"`, addr.IP))
		}
	}
	return ips
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	log.Info("Kubernetes Elasticsearch Cluster discovery")

	////// Collect service name and namespace //////
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalln("Provide Kubernetes Service Name and Namespace Name")
	}
	var serviceName, namespace string
	serviceName = args[0]
	namespace = args[1]
	log.Infof("Searching for %s.%s", serviceName, namespace)
	////////////////////////////////////////////////
	cnf, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to make client: %v", err)
	}

	c, err := kubernetes.NewForConfig(cnf)
	if err != nil {
		log.Fatalf("Failed to make client: %v", err)
	}

	var elasticsearch *api.Service
	// Look for endpoints associated with the Elasticsearch loggging service.
	// First wait for the service to become available.
	for t := time.Now(); time.Since(t) < 5*time.Minute; time.Sleep(10 * time.Second) {
		elasticsearch, err = c.Core().Services(namespace).Get(serviceName)
		if err == nil {
			break
		}
	}
	// If we did not find an elasticsearch logging service then log a warning
	// and return without adding any unicast hosts.
	if elasticsearch == nil {
		log.Warningf("Failed to find the Kubernetes service: %v", err)
		return
	}

	var endpoints *api.Endpoints
	addrs := []string{}
	// Wait for some endpoints.
	count := 0
	for t := time.Now(); time.Since(t) < 5*time.Minute; time.Sleep(10 * time.Second) {
		endpoints, err = c.Core().Endpoints(namespace).Get(serviceName)
		if err != nil {
			continue
		}
		addrs = flattenSubsets(endpoints.Subsets)
		log.Infof("Found %s", addrs)
		if len(addrs) > 0 && len(addrs) == count {
			break
		}
		count = len(addrs)
	}
	// If there was an error finding endpoints then log a warning and quit.
	if err != nil {
		log.Warningf("Error finding endpoints: %v", err)
		return
	}

	log.Infof("Endpoints = %s", addrs)
	fmt.Printf("discovery.zen.ping.unicast.hosts: [%s]\n", strings.Join(addrs, ", "))
}
