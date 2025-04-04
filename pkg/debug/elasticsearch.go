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
	"log"
	"os"
	"path"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	ps "kubeops.dev/petset/client/clientset/versioned"
)

type elasticsearchOpts struct {
	db        *dbapi.Elasticsearch
	dbClient  *cs.Clientset
	psClient  *ps.Clientset
	podClient *kubernetes.Clientset

	operatorNamespace string
	dir               string
	errWriter         *bytes.Buffer
}

func ElasticsearchDebugCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName            string
		operatorNamespace string
	)

	esDebugCmd := &cobra.Command{
		Use: "elasticsearch",
		Aliases: []string{
			"es",
		},
		Short:   "Debug helper for elasticsearch database",
		Example: `kubectl dba debug elasticsearch -n demo sample-elasticsearch --operator-namespace kubedb`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter elasticsearch object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newElasticsearchOpts(f, dbName, namespace, operatorNamespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.collectOperatorLogs()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectForAllDBPetSets()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectForAllDBPods()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectOtherYamls()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	esDebugCmd.Flags().StringVarP(&operatorNamespace, "operator-namespace", "o", "kubedb", "the namespace where the kubedb operator is installed")

	return esDebugCmd
}

func newElasticsearchOpts(f cmdutil.Factory, dbName, namespace, operatorNS string) (*elasticsearchOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	dbClient, err := cs.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	psClient, err := ps.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	podClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := dbClient.KubedbV1().Elasticsearches(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	pwd, _ := os.Getwd()
	dir := path.Join(pwd, db.Name)
	err = os.MkdirAll(path.Join(dir, logsDir), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(dir, yamlsDir), dirPerm)
	if err != nil {
		return nil, err
	}

	opts := &elasticsearchOpts{
		db:                db,
		dbClient:          dbClient,
		psClient:          psClient,
		podClient:         podClient,
		operatorNamespace: operatorNS,
		dir:               dir,
		errWriter:         &bytes.Buffer{},
	}
	return opts, writeYaml(db, path.Join(opts.dir, yamlsDir))
}

func (opts *elasticsearchOpts) collectOperatorLogs() error {
	pods, err := opts.podClient.CoreV1().Pods(opts.operatorNamespace).List(context.TODO(), metav1.ListOptions{})
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
			err = opts.writeLogs(pod.Name, pod.Namespace, operatorContainerName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (opts *elasticsearchOpts) collectForAllDBPetSets() error {
	psLabels := labels.SelectorFromSet(opts.db.OffshootLabels()).String()
	petsets, err := opts.psClient.AppsV1().PetSets(opts.db.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: psLabels,
	})
	if err != nil {
		return err
	}

	psYamlDir := path.Join(opts.dir, yamlsDir, "petsets")
	err = os.MkdirAll(psYamlDir, dirPerm)
	if err != nil {
		return err
	}

	for _, p := range petsets.Items {
		err = writeYaml(&p, psYamlDir)
		if err != nil {
			return err
		}

	}
	return nil
}

func (opts *elasticsearchOpts) collectForAllDBPods() error {
	dbLabels := labels.SelectorFromSet(opts.db.OffshootLabels()).String()
	pods, err := opts.podClient.CoreV1().Pods(opts.db.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: dbLabels,
	})
	if err != nil {
		return err
	}

	podYamlDir := path.Join(opts.dir, yamlsDir)
	for _, pod := range pods.Items {
		err = opts.writeLogsForSinglePod(pod)
		if err != nil {
			return err
		}

		err = writeYaml(&pod, podYamlDir)
		if err != nil {
			return err
		}

	}
	return nil
}

func (opts *elasticsearchOpts) writeLogsForSinglePod(pod corev1.Pod) error {
	for _, c := range pod.Spec.Containers {
		err := opts.writeLogs(pod.Name, pod.Namespace, c.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (opts *elasticsearchOpts) writeLogs(podName, ns, container string) error {
	req := opts.podClient.CoreV1().Pods(ns).GetLogs(podName, &corev1.PodLogOptions{
		Container: container,
	})

	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer podLogs.Close()

	logFile, err := os.Create(path.Join(opts.dir, logsDir, podName+"_"+container+".log"))
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

func (opts *elasticsearchOpts) collectOtherYamls() error {
	opsReqs, err := opts.dbClient.OpsV1alpha1().ElasticsearchOpsRequests(opts.db.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	opsYamlDir := path.Join(opts.dir, yamlsDir, "ops")
	err = os.MkdirAll(opsYamlDir, dirPerm)
	if err != nil {
		return err
	}
	for _, ops := range opsReqs.Items {
		if ops.Spec.DatabaseRef.Name == opts.db.Name {
			err = writeYaml(&ops, opsYamlDir)
			if err != nil {
				return err
			}
		}
	}

	scalers, err := opts.dbClient.AutoscalingV1alpha1().ElasticsearchAutoscalers(opts.db.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	scalerYamlDir := path.Join(opts.dir, yamlsDir, "scaler")
	err = os.MkdirAll(scalerYamlDir, dirPerm)
	if err != nil {
		return err
	}
	for _, sc := range scalers.Items {
		if sc.Spec.DatabaseRef.Name == opts.db.Name {
			err = writeYaml(&sc, scalerYamlDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
