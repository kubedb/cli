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

package alerts

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"kmodules.xyz/client-go/tools/portforward"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PromSvc struct {
	Name      string
	Namespace string
	Port      int
}

type dbOpts struct {
	db         client.Object
	config     *rest.Config
	kubeClient *kubernetes.Clientset
	resource   string
}

func Run(f cmdutil.Factory, args []string, prom PromSvc) {
	if len(args) < 2 {
		log.Fatal("Enter db object's name as an argument")
	}
	resource := args[0]
	dbName := args[1]

	namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		_ = fmt.Errorf("failed to get current namespace")
	}

	opts, err := newDBOpts(f, dbName, namespace, convertedResource(resource))
	if err != nil {
		log.Fatalln(err)
	}

	p, err := opts.ForwardPort("services", prom)
	if err != nil {
		log.Fatalln(err)
	}
	opts.work(p)
}

func convertedResource(resource string) string {
	// standardizing the resource name
	res := strings.ToLower(resource)
	switch res {
	case "es", "elasticsearch", "elasticsearches":
		res = "elasticsearches"
	case "md", "mariadb", "mariadbs":
		res = "mariadbs"
	case "mg", "mongodb", "mongodbs":
		res = "mongodbs"
	case "my", "mysql", "mysqls":
		res = "mysqls"
	case "pg", "postgres", "postgreses":
		res = "postgreses"
	case "rd", "redis", "redises":
		res = "redises"
	default:
		fmt.Printf("%s is not a valid resource type \n", resource)
	}
	return res
}

func newDBOpts(f cmdutil.Factory, dbName, namespace, resource string) (*dbOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	dbRes := schema.GroupVersionResource{Group: "kubedb.com", Version: "v1alpha2", Resource: resource}
	db, err := dc.Resource(dbRes).Namespace(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	opts := &dbOpts{
		db:         db,
		config:     config,
		kubeClient: kubeClient,
		resource:   resource,
	}
	return opts, nil
}

func (opts *dbOpts) ForwardPort(resource string, prom PromSvc) (*portforward.Tunnel, error) {
	tunnel := portforward.NewTunnel(portforward.TunnelOptions{
		Client:    opts.kubeClient.CoreV1().RESTClient(),
		Config:    opts.config,
		Resource:  resource,
		Namespace: prom.Namespace,
		Name:      prom.Name,
		Remote:    prom.Port,
	})

	if err := tunnel.ForwardPort(); err != nil {
		return nil, err
	}
	return tunnel, nil
}

func (opts *dbOpts) work(p *portforward.Tunnel) {
	pc, err := promapi.NewClient(promapi.Config{
		Address: fmt.Sprintf("http://localhost:%d", p.Local),
	})
	if err != nil {
		panic(err)
	}

	promAPI := promv1.NewAPI(pc)
	alertQuery := fmt.Sprintf("ALERTS{alertstate=\"firing\",k8s_group=\"kubedb.com\",k8s_resource=\"%s\",app=\"%s\",app_namespace=\"%s\"}",
		opts.resource, opts.db.GetName(), opts.db.GetNamespace())
	result, warnings, err := promAPI.QueryRange(context.TODO(), alertQuery, promv1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute * 2,
	})
	if err != nil {
		panic(err)
	}
	if len(warnings) > 0 {
		fmt.Println("Warnings:", warnings)
	}

	// Access the elements of the matrix using indexing.
	matrix := result.(model.Matrix)
	fmt.Printf("The number of firing alerts for %s %s/%s is %d \n", opts.resource, opts.db.GetNamespace(), opts.db.GetName(), matrix.Len())
	for _, sample := range matrix {
		alertName := sample.Metric["alertname"]
		alertValue := sample.Values[0].Value
		labels := sample.Metric
		fmt.Printf("Alert Name: %s, Value: %s, Labels: %s\n", alertName, alertValue, labels)
	}
}
