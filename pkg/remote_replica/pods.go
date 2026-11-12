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

package remote_replica

import (
	"context"
	"fmt"

	"kubedb.dev/apimachinery/apis/kubedb"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Containers holding the database server itself, the ones a client query has to
// be exec-ed into.
const (
	PostgresContainerName = "postgres"
	MySQLContainerName    = "mysql"
)

// pickDBPod returns a running database pod to exec into.
//
// primaryOnly is set when the caller intends to run DDL. CREATE ROLE / ALTER
// ROLE / GRANT write to the shared catalog, and a hot standby rejects them with
// "cannot execute ... in a read-only transaction", so only a pod labelled
// kubedb.com/role=primary qualifies. Read-only catalog queries run fine on any
// member, so verification passes primaryOnly=false and merely prefers the
// primary when one happens to be labelled.
func pickDBPod(client kubernetes.Interface, ns string, selector map[string]string, container string, primaryOnly bool) (*core.Pod, error) {
	pods, err := client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set(selector).String(),
	})
	if err != nil {
		return nil, err
	}

	pod, err := selectPod(pods.Items, container, primaryOnly)
	if err != nil {
		return nil, fmt.Errorf("in namespace %s: %v", ns, err)
	}
	return pod, nil
}

// selectPod holds the choice itself, kept free of any client so it can be
// exercised directly. Only pods that actually carry `container` are eligible: a
// KubeDB cluster runs helper pods under the same offshoot labels - the Postgres
// arbiter, for one, runs pg-coordinator alone - and exec-ing a database query
// into those fails in a way that is hard to read.
func selectPod(pods []core.Pod, container string, primaryOnly bool) (*core.Pod, error) {
	var usable []core.Pod
	for i := range pods {
		if pods[i].Status.Phase == core.PodRunning && hasRunningContainer(&pods[i], container) {
			usable = append(usable, pods[i])
		}
	}
	if len(usable) == 0 {
		// Reported explicitly: an empty list used to be swallowed and treated as
		// success, which produced a config for a user that was never created.
		return nil, fmt.Errorf("none of the %d database pod(s) has a running %q container", len(pods), container)
	}

	for i := range usable {
		if usable[i].Labels[kubedb.LabelRole] == kubedb.DatabasePodPrimary {
			return &usable[i], nil
		}
	}
	if primaryOnly {
		return nil, fmt.Errorf("no pod is labelled %s=%s, so there is no writable primary to run DDL against",
			kubedb.LabelRole, kubedb.DatabasePodPrimary)
	}
	return &usable[0], nil
}

func hasRunningContainer(pod *core.Pod, container string) bool {
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Name == container {
			return cs.State.Running != nil
		}
	}
	return false
}
