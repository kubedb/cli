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

package adopt

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"time"

	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// cluster bundles the clients and the namespace one command works in.
type cluster struct {
	kube kubernetes.Interface
	dyn  dynamic.Interface
	cfg  *rest.Config
	ns   string
	out  io.Writer
	// poll is the interval of every wait.
	poll time.Duration
	// restTLS is set when the running proxies serve TLS.
	restTLS *tls.Config
}

// newCluster connects to the cluster. fileNamespace, from an adoption file, is
// used when -n is not given.
func newCluster(f cmdutil.Factory, out io.Writer, fileNamespace string) (*cluster, error) {
	ns, explicit, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}
	if !explicit && fileNamespace != "" {
		ns = fileNamespace
	}
	cfg, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	kc, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &cluster{kube: kc, dyn: dyn, cfg: cfg, ns: ns, out: out, poll: 3 * time.Second}, nil
}

func (c *cluster) logf(format string, a ...any) {
	_, _ = fmt.Fprintf(c.out, "%s  %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
}

// side is one place a Unit's pods run: the chart workload, the KubeDB Milvus
// role, or (during a rollback, after the Milvus is gone) the orphaned PetSet.
type side interface {
	String() string
	desired(ctx context.Context) (int32, error)
	scale(ctx context.Context, n int32) error
	selector() map[string]string
}

// podCounts returns the Ready pods (not terminating) and all pods a selector matches.
func (c *cluster) podCounts(ctx context.Context, sel map[string]string) (ready, total int32, err error) {
	pods, err := c.kube.CoreV1().Pods(c.ns).List(ctx, metav1.ListOptions{LabelSelector: labels.SelectorFromSet(sel).String()})
	if err != nil {
		return 0, 0, err
	}
	for _, p := range pods.Items {
		total++
		if p.DeletionTimestamp == nil && podReady(p) {
			ready++
		}
	}
	return ready, total, nil
}

func podReady(p core.Pod) bool {
	for _, c := range p.Status.Conditions {
		if c.Type == core.PodReady {
			return c.Status == core.ConditionTrue
		}
	}
	return false
}

// waitReady waits until a side has at least n Ready pods.
func (c *cluster) waitReady(ctx context.Context, s side, n int32, timeout time.Duration) error {
	return c.waitFor(ctx, timeout, fmt.Sprintf("%s: %d Ready", s, n), func() (bool, string, error) {
		r, t, err := c.podCounts(ctx, s.selector())
		return r >= n, fmt.Sprintf("%d/%d Ready", r, t), err
	})
}

// waitGone waits until a side has at most n pods, terminating ones included.
func (c *cluster) waitGone(ctx context.Context, s side, n int32, timeout time.Duration) error {
	return c.waitFor(ctx, timeout, fmt.Sprintf("%s: at most %d pods", s, n), func() (bool, string, error) {
		_, t, err := c.podCounts(ctx, s.selector())
		return t <= n, fmt.Sprintf("%d pods", t), err
	})
}

func (c *cluster) waitFor(ctx context.Context, timeout time.Duration, what string, cond func() (bool, string, error)) error {
	deadline := time.Now().Add(timeout)
	lastReport := time.Now()
	for {
		ok, now, err := cond()
		if err != nil {
			return fmt.Errorf("waiting for %s: %w", what, err)
		}
		if ok {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out after %s waiting for %s (now %s)", timeout, what, now)
		}
		if time.Since(lastReport) >= 30*time.Second {
			c.logf("  waiting for %s (now %s)", what, now)
			lastReport = time.Now()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.poll):
		}
	}
}

// chartSide is the workload the chart created.
type chartSide struct {
	c *cluster
	u Unit
}

func (s chartSide) String() string { return fmt.Sprintf("%s %s", s.u.Kind, s.u.Name) }

func (s chartSide) selector() map[string]string { return s.u.Selector }

func (s chartSide) desired(ctx context.Context) (int32, error) {
	switch s.u.Kind {
	case "Deployment":
		d, err := s.c.kube.AppsV1().Deployments(s.c.ns).Get(ctx, s.u.Name, metav1.GetOptions{})
		if err != nil {
			return 0, err
		}
		if d.Spec.Replicas == nil {
			return 1, nil
		}
		return *d.Spec.Replicas, nil
	case "StatefulSet":
		st, err := s.c.kube.AppsV1().StatefulSets(s.c.ns).Get(ctx, s.u.Name, metav1.GetOptions{})
		if err != nil {
			return 0, err
		}
		if st.Spec.Replicas == nil {
			return 1, nil
		}
		return *st.Spec.Replicas, nil
	}
	return 0, fmt.Errorf("unsupported workload kind %s", s.u.Kind)
}

func (s chartSide) scale(ctx context.Context, n int32) error {
	patch := []byte(fmt.Sprintf(`{"spec":{"replicas":%d}}`, n))
	var err error
	switch s.u.Kind {
	case "Deployment":
		_, err = s.c.kube.AppsV1().Deployments(s.c.ns).Patch(ctx, s.u.Name, types.MergePatchType, patch, metav1.PatchOptions{})
	case "StatefulSet":
		_, err = s.c.kube.AppsV1().StatefulSets(s.c.ns).Patch(ctx, s.u.Name, types.MergePatchType, patch, metav1.PatchOptions{})
	default:
		err = fmt.Errorf("unsupported workload kind %s", s.u.Kind)
	}
	if err != nil {
		return fmt.Errorf("failed to scale %s to %d: %w", s, n, err)
	}
	return nil
}

// kubedbSide is the role (or node group) of the KubeDB Milvus.
type kubedbSide struct {
	c    *cluster
	name string
	u    Unit
}

func (s kubedbSide) String() string {
	if s.u.Group != "" {
		return fmt.Sprintf("Milvus %s %s group %s", s.name, s.u.Role, s.u.Group)
	}
	return fmt.Sprintf("Milvus %s %s", s.name, s.u.Role)
}

func (s kubedbSide) selector() map[string]string { return kubedbSelector(s.name, s.u) }

// replicasPath returns the field path of the Unit's replicas in the Milvus.
func (s kubedbSide) replicasPath(db *unstructured.Unstructured) ([]string, int, error) {
	base := []string{"spec", "topology", "distributed", s.u.Role}
	if s.u.Group == "" {
		return append(base, "replicas"), -1, nil
	}
	groups, _, _ := unstructured.NestedSlice(db.Object, append(base, "groups")...)
	for i, g := range groups {
		if m, ok := g.(map[string]any); ok && m["name"] == s.u.Group {
			return nil, i, nil
		}
	}
	return nil, -1, fmt.Errorf("no %s group %q in Milvus %s", s.u.Role, s.u.Group, s.name)
}

func (s kubedbSide) desired(ctx context.Context) (int32, error) {
	db, err := s.c.dyn.Resource(MilvusGVR).Namespace(s.c.ns).Get(ctx, s.name, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}
	path, idx, err := s.replicasPath(db)
	if err != nil {
		return 0, err
	}
	if idx < 0 {
		n, _, _ := unstructured.NestedInt64(db.Object, path...)
		return int32(n), nil
	}
	groups, _, _ := unstructured.NestedSlice(db.Object, "spec", "topology", "distributed", s.u.Role, "groups")
	n, _, _ := unstructured.NestedInt64(groups[idx].(map[string]any), "replicas")
	return int32(n), nil
}

func (s kubedbSide) scale(ctx context.Context, n int32) error {
	db, err := s.c.dyn.Resource(MilvusGVR).Namespace(s.c.ns).Get(ctx, s.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	_, idx, err := s.replicasPath(db)
	if err != nil {
		return err
	}
	var patch []byte
	pt := types.MergePatchType
	if idx < 0 {
		patch = []byte(fmt.Sprintf(`{"spec":{"topology":{"distributed":{%q:{"replicas":%d}}}}}`, s.u.Role, n))
	} else {
		pt = types.JSONPatchType
		patch = []byte(fmt.Sprintf(`[{"op":"add","path":"/spec/topology/distributed/%s/groups/%d/replicas","value":%d}]`, s.u.Role, idx, n))
	}
	if _, err := s.c.dyn.Resource(MilvusGVR).Namespace(s.c.ns).Patch(ctx, s.name, pt, patch, metav1.PatchOptions{}); err != nil {
		return fmt.Errorf("failed to set %s to %d replicas: %w", s, n, err)
	}
	return nil
}

// petSetSide is a KubeDB PetSet after its Milvus was deleted with the
// orphan-on-delete annotation: the pods keep running under the PetSet controller.
type petSetSide struct {
	c    *cluster
	name string
	u    Unit
}

func (s petSetSide) String() string { return "PetSet " + petSetName(s.name, s.u) }

func (s petSetSide) selector() map[string]string { return kubedbSelector(s.name, s.u) }

func (s petSetSide) desired(ctx context.Context) (int32, error) {
	ps, err := s.c.dyn.Resource(PetSetGVR).Namespace(s.c.ns).Get(ctx, petSetName(s.name, s.u), metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	n, found, _ := unstructured.NestedInt64(ps.Object, "spec", "replicas")
	if !found {
		return 1, nil
	}
	return int32(n), nil
}

// remove deletes the PetSet and, with it, its pods.
func (s petSetSide) remove(ctx context.Context) error {
	bg := metav1.DeletePropagationBackground
	err := s.c.dyn.Resource(PetSetGVR).Namespace(s.c.ns).Delete(ctx, petSetName(s.name, s.u), metav1.DeleteOptions{PropagationPolicy: &bg})
	if err != nil && !kerr.IsNotFound(err) {
		return fmt.Errorf("failed to delete %s: %w", s, err)
	}
	return nil
}

func (s petSetSide) scale(ctx context.Context, n int32) error {
	patch := []byte(fmt.Sprintf(`{"spec":{"replicas":%d}}`, n))
	_, err := s.c.dyn.Resource(PetSetGVR).Namespace(s.c.ns).Patch(ctx, petSetName(s.name, s.u), types.MergePatchType, patch, metav1.PatchOptions{})
	if kerr.IsNotFound(err) && n == 0 {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to scale %s to %d: %w", s, n, err)
	}
	return nil
}

// loadState finds the adoption of a release in the namespace.
func (c *cluster) loadState(ctx context.Context, release string) (*State, *core.ConfigMap, error) {
	cms, err := c.kube.CoreV1().ConfigMaps(c.ns).List(ctx, metav1.ListOptions{LabelSelector: LabelRelease + "=" + release})
	if err != nil {
		return nil, nil, err
	}
	for i := range cms.Items {
		cm := &cms.Items[i]
		raw, ok := cm.Data[stateKey]
		if !ok {
			continue
		}
		st := &State{}
		if err := json.Unmarshal([]byte(raw), st); err != nil {
			return nil, nil, fmt.Errorf("ConfigMap %s holds an unreadable adoption state: %w", cm.Name, err)
		}
		return st, cm, nil
	}
	return nil, nil, nil
}

func (c *cluster) saveState(ctx context.Context, st *State) error {
	raw, err := json.Marshal(st)
	if err != nil {
		return err
	}
	cm, err := c.kube.CoreV1().ConfigMaps(c.ns).Get(ctx, stateConfigMapName(st.Name), metav1.GetOptions{})
	if err != nil {
		return err
	}
	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	cm.Data[stateKey] = string(raw)
	_, err = c.kube.CoreV1().ConfigMaps(c.ns).Update(ctx, cm, metav1.UpdateOptions{})
	return err
}

// renderedConfig returns the milvus.yaml the operator rendered for a Milvus:
// the Secret it owns that carries that key.
func (c *cluster) renderedConfig(ctx context.Context, db *unstructured.Unstructured) (string, error) {
	secs, err := c.kube.CoreV1().Secrets(c.ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, s := range secs.Items {
		for _, o := range s.OwnerReferences {
			if o.UID == db.GetUID() {
				if v, ok := s.Data["milvus.yaml"]; ok {
					return string(v), nil
				}
			}
		}
	}
	return "", nil
}
