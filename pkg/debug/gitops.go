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

package debug

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	gitops "kubedb.dev/apimachinery/apis/gitops/v1alpha1"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"
	kubedbscheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubedbscheme.AddToScheme(scheme))
}

type GitOpsStatus struct {
	GitOps gitops.GitOpsStatus `json:"gitops,omitempty" yaml:"gitops,omitempty"`
}

type GitOps struct {
	Status GitOpsStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type dbInfo struct {
	resource  string
	name      string
	namespace string
}

type gitOpsOpts struct {
	kc         client.Client
	config     *rest.Config
	db         dbInfo
	kubeClient kubernetes.Interface

	operatorNamespace string
	dir               string
	errWriter         *bytes.Buffer
	resMap            map[string]string
}

func GitOpsDebugCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string
	opts := newGitOpsOpts(f)
	gitOpsDebugCmd := &cobra.Command{
		Use: "gitops",
		Aliases: []string{
			"git",
		},
		Short:   "Debug helper for gitops databases",
		Example: `kubectl dba debug gitops --db-type mysql -n demo sample-mysql --operator-namespace kubedb`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mysql object's name as an argument")
			}
			dbName = args[0]

			pwd, _ := os.Getwd()
			dir := path.Join(pwd, dbName)
			err := os.MkdirAll(path.Join(dir, logsDir), dirPerm)
			if err != nil {
				log.Fatalln(fmt.Errorf("failed to create directory %s: %w", dir, err))
			}
			err = os.MkdirAll(path.Join(dir, yamlsDir), dirPerm)
			if err != nil {
				log.Fatalln(fmt.Errorf("failed to create directory %s: %w", dir, err))
			}
			opts.dir = dir
			opts.db.name = dbName

			err = opts.populateResourceMap()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectGitOpsDatabase()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectDatabase()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectOperatorLogs(opts.operatorNamespace)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	gitOpsDebugCmd.Flags().StringVarP(&opts.db.namespace, "namespace", "n", "demo", "Database namespace")
	gitOpsDebugCmd.Flags().StringVarP(&opts.operatorNamespace, "operator-namespace", "o", "kubedb", "the namespace where the kubedb gitops operator is installed")
	gitOpsDebugCmd.Flags().StringVarP(&opts.db.resource, "db-type", "t", "postgres", "database type")

	return gitOpsDebugCmd
}

func newGitOpsOpts(f cmdutil.Factory) *gitOpsOpts {
	config, err := f.ToRESTConfig()
	if err != nil {
		log.Fatalln(err)
	}

	kc, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create kube client: %v", err)
	}

	opts := &gitOpsOpts{
		kc:         kc,
		config:     config,
		errWriter:  &bytes.Buffer{},
		kubeClient: cs,
	}
	return opts
}

func (g *gitOpsOpts) collectGitOpsDatabase() error {
	var uns unstructured.Unstructured
	uns.SetGroupVersionKind(gitops.SchemeGroupVersion.WithKind(g.getKindFromResource(g.db.resource)))
	err := g.kc.Get(context.Background(), types.NamespacedName{
		Namespace: g.db.namespace,
		Name:      g.db.name,
	}, &uns)
	if err != nil {
		log.Fatalf("failed to get gitops database obj: %v", err)
		return err
	}

	var gitOpsObj GitOps
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(uns.Object, &gitOpsObj)
	if err != nil {
		log.Fatalf("failed to convert unstructured to gitops obj: %v", err)
		return err
	}

	if err := g.collectOpsRequests(gitOpsObj.Status); err != nil {
		return err
	}

	return writeYaml(&uns, g.dir)
}

func (g *gitOpsOpts) collectOpsRequests(gitOpsStatus GitOpsStatus) error {
	opsYamlDir := path.Join(g.dir, yamlsDir, "ops")
	err := os.MkdirAll(opsYamlDir, dirPerm)
	if err != nil {
		return err
	}
	for _, info := range gitOpsStatus.GitOps.GitOpsInfo {
		for _, op := range info.Operations {
			var uns unstructured.Unstructured
			uns.SetGroupVersionKind(opsapi.SchemeGroupVersion.WithKind(g.getKindFromResource(g.db.resource + "opsrequest")))
			err := g.kc.Get(context.Background(), types.NamespacedName{
				Namespace: g.db.namespace,
				Name:      op.Name,
			}, &uns)
			if err != nil {
				log.Fatalf("failed to get opsrequest: %v", err)
				return err
			}
			err = writeYaml(&uns, opsYamlDir)
			if err != nil {
				return err
			}
			var opsStatus opsapi.OpsRequestStatus
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(uns.Object, &opsStatus)
			if err != nil {
				log.Fatalf("failed to convert unstructured to opsrequest obj: %v", err)
				return err
			}
			// if opsStatus.Phase == opsapi.OpsRequestPhaseFailed {
			// 	for _, cond := range opsStatus.Conditions {
			// 		if cond.Type == opsapi.Failed {
			// 			// TODO: ()
			// 		}
			// 	}
			// }
		}
	}

	return nil
}

func (g *gitOpsOpts) collectDatabase() error {
	var uns unstructured.Unstructured
	uns.SetGroupVersionKind(dbapi.SchemeGroupVersion.WithKind(g.getKindFromResource(g.db.resource)))
	err := g.kc.Get(context.Background(), types.NamespacedName{
		Namespace: g.db.namespace,
		Name:      g.db.name,
	}, &uns)
	if err != nil {
		log.Fatalf("failed to get database: %v", err)
	}

	return writeYaml(&uns, path.Join(g.dir, yamlsDir))
}

func (g *gitOpsOpts) collectOperatorLogs(operatorNamespace string) error {
	pods, err := g.kubeClient.CoreV1().Pods(operatorNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		isOperatorPod := false
		for _, container := range pod.Spec.Containers {
			if container.Name == operatorContainerName {
				isOperatorPod = true
			}
		}
		if isOperatorPod {
			err = g.writeLogs(pod.Name, pod.Namespace, operatorContainerName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *gitOpsOpts) populateResourceMap() error {
	dc, err := discovery.NewDiscoveryClientForConfig(g.config)
	if err != nil {
		return err
	}
	g.resMap = make(map[string]string)

	if err := g.populate(dc, "kubedb.com/v1"); err != nil {
		return err
	}
	if err := g.populate(dc, "kubedb.com/v1alpha2"); err != nil {
		return err
	}
	if err := g.populate(dc, "gitops.kubedb.com/v1alpha1"); err != nil {
		return err
	}
	if err := g.populate(dc, "ops.kubedb.com/v1alpha1"); err != nil {
		return err
	}
	return nil
}

func (g *gitOpsOpts) populate(dc *discovery.DiscoveryClient, gv string) error {
	resources, err := dc.ServerResourcesForGroupVersion(gv)
	if err != nil {
		return err
	}
	for _, r := range resources.APIResources {
		if !strings.ContainsAny(r.Name, "/") {
			g.resMap[r.Name] = r.Kind
			g.resMap[r.SingularName] = r.Kind
			for _, s := range r.ShortNames {
				g.resMap[s] = r.Kind
			}
			g.resMap[r.Kind] = r.Kind
		}
	}
	return nil
}

func (g *gitOpsOpts) getKindFromResource(res string) string {
	kind, exists := g.resMap[res]
	if !exists {
		_ = fmt.Errorf("resource %s not supported", res)
	}
	return kind
}

func (g *gitOpsOpts) writeLogs(podName, ns, container string) error {
	req := g.kubeClient.CoreV1().Pods(ns).GetLogs(podName, &corev1.PodLogOptions{
		Container: container,
	})

	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer podLogs.Close()

	logFile, err := os.Create(path.Join(g.dir, logsDir, podName+"_"+container+".log"))
	if err != nil {
		return err
	}
	defer logFile.Close()

	buf := make([]byte, 1024)
	for {
		bytesRead, err := podLogs.Read(buf)
		if err != nil {
			break
		}
		_, _ = logFile.Write(buf[:bytesRead])
	}
	return nil
}
