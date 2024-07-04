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
	"time"

	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/cli/pkg/monitor"

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

type dbOpts struct {
	db         client.Object
	config     *rest.Config
	kubeClient *kubernetes.Clientset
	resource   string
}

func Run(f cmdutil.Factory, args []string, prom monitor.PromSvc) {
	if len(args) < 2 {
		log.Fatal("Enter db object's name as an argument")
	}
	resource := args[0]
	dbName := args[1]

	namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		_ = fmt.Errorf("failed to get current namespace")
	}

	opts, err := newDBOpts(f, dbName, namespace, monitor.ConvertedResourceToPlural(resource))
	if err != nil {
		log.Fatalln(err)
	}

	promClient, tunnel := monitor.GetPromClientAndTunnel(opts.config, prom)
	defer tunnel.Close()

	err = opts.work(promClient)
	if err != nil {
		log.Fatalln(err)
	}
}

func newDBOpts(f cmdutil.Factory, dbName, namespace, resource string) (*dbOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	gvk := api.SchemeGroupVersion
	dbRes := schema.GroupVersionResource{Group: gvk.Group, Version: gvk.Version, Resource: resource}
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

func (opts *dbOpts) ForwardPort(resource string, prom monitor.PromSvc) (*portforward.Tunnel, error) {
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

func (opts *dbOpts) work(promAPI promv1.API) error {
	alertQuery := fmt.Sprintf("ALERTS{alertstate=\"firing\",k8s_group=\"kubedb.com\",k8s_resource=\"%s\",app=\"%s\",app_namespace=\"%s\"}",
		opts.resource, opts.db.GetName(), opts.db.GetNamespace())
	result, warnings, err := promAPI.QueryRange(context.TODO(), alertQuery, promv1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute * 2,
	})
	if err != nil {
		return err
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
	return nil
}
