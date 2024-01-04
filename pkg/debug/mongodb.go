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

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type mongodbOpts struct {
	db        *api.MongoDB
	dbClient  *cs.Clientset
	podClient *kubernetes.Clientset

	operatorNamespace string
	dir               string
	errWriter         *bytes.Buffer
}

func MongoDBDebugCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName            string
		operatorNamespace string
	)

	mgDebugCmd := &cobra.Command{
		Use: "mongodb",
		Aliases: []string{
			"mg",
		},
		Short:   "Debug helper for mongodb database",
		Example: `kubectl dba debug mongodb -n demo sample-mongodb --operator-namespace kubedb`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mongodb object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMongodbOpts(f, dbName, namespace, operatorNamespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.collectOperatorLogs()
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
	mgDebugCmd.Flags().StringVarP(&operatorNamespace, "operator-namespace", "o", "kubedb", "the namespace where the kubedb operator is installed")

	return mgDebugCmd
}

func newMongodbOpts(f cmdutil.Factory, dbName, namespace, operatorNS string) (*mongodbOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	dbClient, err := cs.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	podClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := dbClient.KubedbV1alpha2().MongoDBs(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
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

	opts := &mongodbOpts{
		db:                db,
		dbClient:          dbClient,
		podClient:         podClient,
		operatorNamespace: operatorNS,
		dir:               dir,
		errWriter:         &bytes.Buffer{},
	}
	return opts, writeYaml(db, path.Join(opts.dir, yamlsDir))
}

func (opts *mongodbOpts) collectOperatorLogs() error {
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

func (opts *mongodbOpts) collectForAllDBPods() error {
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

func (opts *mongodbOpts) writeLogsForSinglePod(pod corev1.Pod) error {
	for _, c := range pod.Spec.Containers {
		err := opts.writeLogs(pod.Name, pod.Namespace, c.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (opts *mongodbOpts) writeLogs(podName, ns, container string) error {
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

func (opts *mongodbOpts) collectOtherYamls() error {
	opsReqs, err := opts.dbClient.OpsV1alpha1().MongoDBOpsRequests(opts.db.Namespace).List(context.TODO(), metav1.ListOptions{})
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

	scalers, err := opts.dbClient.AutoscalingV1alpha1().MongoDBAutoscalers(opts.db.Namespace).List(context.TODO(), metav1.ListOptions{})
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
