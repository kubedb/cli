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
	"os"
	"path"

	autoscalerapi "kubedb.dev/apimachinery/apis/autoscaling/v1alpha1"
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"
	kubedbscheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	psapi "kubeops.dev/petset/apis/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubedbscheme.AddToScheme(scheme))
	utilruntime.Must(psapi.AddToScheme(scheme))
}

type dbOpts struct {
	kc         client.Client
	kubeClient *kubernetes.Clientset
	kind       string
	db         metav1.ObjectMeta
	selectors  map[string]string

	operatorNamespace string
	dir               string
	errWriter         *bytes.Buffer
}

func newDBOpts(f cmdutil.Factory, gvk schema.GroupVersionKind, dbName, namespace, operatorNS string) (*dbOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	kc, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	pwd, _ := os.Getwd()
	dir := path.Join(pwd, dbName)
	err = os.MkdirAll(path.Join(dir, logsDir), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(dir, yamlsDir), dirPerm)
	if err != nil {
		return nil, err
	}

	opts := &dbOpts{
		kc:                kc,
		kubeClient:        kubeClient,
		kind:              gvk.Kind,
		db:                metav1.ObjectMeta{Namespace: namespace, Name: dbName},
		operatorNamespace: operatorNS,
		dir:               dir,
		errWriter:         &bytes.Buffer{},
	}
	return opts, nil
}

func (opts *dbOpts) collectALl() error {
	klog.Infof("11111")
	err := opts.collectOperatorLogs()
	if err != nil {
		return err
	}
	err = opts.collectForAllDBPetSets()
	if err != nil {
		return err
	}
	err = opts.collectForAllDBPods()
	if err != nil {
		return err
	}
	err = opts.collectOtherYamls()
	if err != nil {
		return err
	}
	return nil
}

func (opts *dbOpts) collectOperatorLogs() error {
	var pods corev1.PodList
	err := opts.kc.List(context.TODO(), &pods)
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

func (opts *dbOpts) collectForAllDBPetSets() error {
	var petsets psapi.PetSetList
	err := opts.kc.List(context.TODO(), &petsets, client.MatchingLabels(opts.selectors), client.InNamespace(opts.db.GetNamespace()))
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

func (opts *dbOpts) collectForAllDBPods() error {
	var pods corev1.PodList
	err := opts.kc.List(context.TODO(), &pods, client.MatchingLabels(opts.selectors), client.InNamespace(opts.db.GetNamespace()))
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

func (opts *dbOpts) writeLogsForSinglePod(pod corev1.Pod) error {
	for _, c := range pod.Spec.Containers {
		err := opts.writeLogs(pod.Name, pod.Namespace, c.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (opts *dbOpts) writeLogs(podName, ns, container string) error {
	req := opts.kubeClient.CoreV1().Pods(ns).GetLogs(podName, &corev1.PodLogOptions{
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

func (opts *dbOpts) collectOtherYamls() error {
	var opsReqs unstructured.UnstructuredList
	opsReqs.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   opsapi.SchemeGroupVersion.Group,
		Version: opsapi.SchemeGroupVersion.Version,
		Kind:    opts.kind + "OpsRequest",
	})
	if err := opts.kc.List(context.Background(), &opsReqs, client.InNamespace(opts.db.GetNamespace())); err != nil {
		return err
	}

	opsYamlDir := path.Join(opts.dir, yamlsDir, "ops")
	err := os.MkdirAll(opsYamlDir, dirPerm)
	if err != nil {
		return err
	}
	for _, o := range opsReqs.Items {
		var ops OpsRequest
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(o.Object, &ops); err != nil {
			return fmt.Errorf("failed to unmarshal binding %s: %w", o.GetName(), err)
		}
		if ops.Spec.DatabaseRef.Name == opts.db.GetName() {
			err = writeYaml(&o, opsYamlDir)
			if err != nil {
				return err
			}
		}
	}

	var autoscalers unstructured.UnstructuredList
	autoscalers.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   autoscalerapi.SchemeGroupVersion.Group,
		Version: autoscalerapi.SchemeGroupVersion.Version,
		Kind:    opts.kind + "Autoscaler",
	})
	if err := opts.kc.List(context.Background(), &autoscalers, client.InNamespace(opts.db.GetNamespace())); err != nil {
		return err
	}

	scalerYamlDir := path.Join(opts.dir, yamlsDir, "scaler")
	err = os.MkdirAll(scalerYamlDir, dirPerm)
	if err != nil {
		return err
	}
	for _, sc := range autoscalers.Items {
		var autoscaler Autoscaler
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(sc.Object, &autoscaler); err != nil {
			return fmt.Errorf("failed to unmarshal binding %s: %w", sc.GetName(), err)
		}
		if autoscaler.Spec.DatabaseRef.Name == opts.db.GetName() {
			err = writeYaml(&sc, scalerYamlDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
