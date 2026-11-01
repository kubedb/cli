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
	"fmt"
	"time"

	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/yaml"
)

// Handover methods.
const (
	// methodOverlap starts the KubeDB pods of a role before it stops the chart's,
	// so the role never stops serving. It needs room for one extra set of pods
	// (one extra QueryNode at a time).
	methodOverlap = "overlap"
	// methodStopStart stops the chart's pods of a role before it starts KubeDB's.
	// Each role is briefly unavailable; no extra capacity is needed.
	methodStopStart = "stop-start"
)

// handoverOrder is the order the roles change hands in. With the overlap
// method the proxies go first, so clients reach the KubeDB pods early; a
// stop-start handover moves the coordinators first and the proxies last.
func handoverOrder(method string) []string {
	if method == methodStopStart {
		return []string{roleMixCoord, roleDataNode, roleQueryNode, roleStreamingNode, roleProxy}
	}
	return []string{roleProxy, roleMixCoord, roleDataNode, roleStreamingNode, roleQueryNode}
}

// rollbackOrder returns the roles to the chart, the proxies last so the
// client Service moves once every other role serves from the chart again.
var rollbackOrder = []string{roleQueryNode, roleStreamingNode, roleDataNode, roleMixCoord, roleProxy}

type timing struct {
	Timeout        time.Duration // per wait
	Settle         time.Duration // after each role
	QueryNodeDelay time.Duration // after each new QueryNode, before an old one stops
}

// step is one replica change on the way from one side to the other: grow
// `to` to toN, then shrink `from` to fromN.
type step struct {
	toN, fromN int32
}

// swapSteps returns the steps that move n replicas from one side to the
// other. QueryNodes go one at a time, so the cluster never holds more than one
// extra QueryNode and every old one stops only after a new one serves.
func swapSteps(n int32, oneByOne bool) []step {
	if !oneByOne || n <= 1 {
		return []step{{toN: n, fromN: 0}}
	}
	var out []step
	for i := int32(1); i <= n; i++ {
		out = append(out, step{toN: i, fromN: n - i})
	}
	return out
}

// overlapSwap moves a Unit from one side to the other without a gap: each
// step grows the target first, waits for it to be Ready, then shrinks the
// source. It is safe to repeat: the target never shrinks and the source never
// grows.
func (c *cluster) overlapSwap(ctx context.Context, from, to side, n int32, oneByOne bool, t timing, beforeShrink func() error) error {
	steps := swapSteps(n, oneByOne)
	for i, s := range steps {
		cur, err := to.desired(ctx)
		if err != nil {
			return err
		}
		if cur < s.toN {
			c.logf("%s -> %d replicas", to, s.toN)
			if err := to.scale(ctx, s.toN); err != nil {
				return err
			}
		}
		if err := c.waitReady(ctx, to, s.toN, t.Timeout); err != nil {
			return err
		}
		if oneByOne && t.QueryNodeDelay > 0 && cur < s.toN {
			c.logf("  %s serves %d; waiting %s before stopping one of %s", to, s.toN, t.QueryNodeDelay, from)
			sleep(ctx, t.QueryNodeDelay)
		}
		if i == len(steps)-1 && beforeShrink != nil {
			if err := beforeShrink(); err != nil {
				return err
			}
		}
		cur, err = from.desired(ctx)
		if err != nil {
			return err
		}
		if cur > s.fromN {
			c.logf("%s -> %d replicas", from, s.fromN)
			if err := from.scale(ctx, s.fromN); err != nil {
				return err
			}
		}
		if err := c.waitGone(ctx, from, s.fromN, t.Timeout); err != nil {
			return err
		}
	}
	return nil
}

// stopStartSwap stops the source, waits until its pods are gone, then starts
// the target.
func (c *cluster) stopStartSwap(ctx context.Context, from, to side, n int32, t timing) error {
	if cur, err := from.desired(ctx); err != nil {
		return err
	} else if cur > 0 {
		c.logf("%s -> 0 replicas", from)
		if err := from.scale(ctx, 0); err != nil {
			return err
		}
	}
	if err := c.waitGone(ctx, from, 0, t.Timeout); err != nil {
		return err
	}
	if cur, err := to.desired(ctx); err != nil {
		return err
	} else if cur < n {
		c.logf("%s -> %d replicas", to, n)
		if err := to.scale(ctx, n); err != nil {
			return err
		}
	}
	return c.waitReady(ctx, to, n, t.Timeout)
}

// handOver moves every Unit of the adoption to KubeDB, recording each one.
func (c *cluster) handOver(ctx context.Context, st *State, method string, t timing) error {
	for _, role := range handoverOrder(method) {
		for _, u := range st.unitsOf(role) {
			if st.isDone(u) {
				continue
			}
			start := time.Now()
			from, to := chartSide{c, u}, kubedbSide{c, st.Name, u}
			c.logf("== %s (%s): %d replicas, %s", u.key(), u.Name, u.Replicas, method)
			var err error
			if method == methodStopStart {
				err = c.stopStartSwap(ctx, from, to, u.Replicas, t)
			} else {
				err = c.overlapSwap(ctx, from, to, u.Replicas, u.Role == roleQueryNode, t, nil)
			}
			if err != nil {
				return fmt.Errorf("%s: %w", u.key(), err)
			}
			st.markDone(u)
			if err := c.saveState(ctx, st); err != nil {
				return err
			}
			c.logf("== %s handed over in %s", u.key(), time.Since(start).Round(time.Second))
			sleep(ctx, t.Settle)
		}
	}
	return nil
}

// rollBack returns every Unit to the chart workloads. The Milvus must be gone
// (deleted with the orphan-on-delete annotation), so its PetSets are scaled
// directly and the operator no longer moves the client Service.
func (c *cluster) rollBack(ctx context.Context, st *State, t timing) error {
	for _, role := range rollbackOrder {
		for _, u := range st.unitsOf(role) {
			from, to := petSetSide{c, st.Name, u}, chartSide{c, u}
			cur, err := from.desired(ctx)
			if err != nil {
				return err
			}
			_, pods, err := c.podCounts(ctx, from.selector())
			if err != nil {
				return err
			}
			var restore func() error
			if u.Role == roleProxy && st.Service != nil {
				restore = func() error { return c.restoreClientService(ctx, st) }
			}
			if cur == 0 && pods == 0 {
				// Nothing of this Unit runs on KubeDB (it was never handed over, or a
				// stop-start handover stopped between the two sides).
				if n, err := to.desired(ctx); err != nil {
					return err
				} else if n < u.Replicas {
					c.logf("%s -> %d replicas", to, u.Replicas)
					if err := to.scale(ctx, u.Replicas); err != nil {
						return err
					}
				}
				if err := c.waitReady(ctx, to, u.Replicas, t.Timeout); err != nil {
					return err
				}
				if restore != nil {
					if err := restore(); err != nil {
						return err
					}
				}
				st.unmarkDone(u)
				continue
			}
			start := time.Now()
			c.logf("== %s (%s): back to the chart, %d replicas", u.key(), u.Name, u.Replicas)
			if err := c.overlapSwap(ctx, from, to, u.Replicas, u.Role == roleQueryNode, t, restore); err != nil {
				return fmt.Errorf("%s: %w", u.key(), err)
			}
			st.unmarkDone(u)
			c.logf("== %s back on the chart in %s", u.key(), time.Since(start).Round(time.Second))
			sleep(ctx, t.Settle)
		}
	}
	return nil
}

// restoreClientService points the client Service back at the chart's proxies,
// with the labels and ports it had before the adoption.
func (c *cluster) restoreClientService(ctx context.Context, st *State) error {
	rec := st.Service
	svc, err := c.kube.CoreV1().Services(c.ns).Get(ctx, rec.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	svc.Spec.Selector = rec.Selector
	if rec.Labels != nil {
		svc.Labels = rec.Labels
	}
	if len(rec.Ports) > 0 {
		svc.Spec.Ports = rec.Ports
	}
	if _, err := c.kube.CoreV1().Services(c.ns).Update(ctx, svc, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("failed to point Service %s back at the chart proxies: %w", rec.Name, err)
	}
	c.logf("Service %s points at the chart proxies again", rec.Name)
	return nil
}

// releaseMilvus deletes the Milvus with the orphan-on-delete annotation, so
// its PetSets, Services and pods stay as they are.
func (c *cluster) releaseMilvus(ctx context.Context, name string, timeout time.Duration) error {
	ri := c.dyn.Resource(MilvusGVR).Namespace(c.ns)
	if _, err := ri.Get(ctx, name, metav1.GetOptions{}); kerr.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}
	patch := fmt.Sprintf(`{"metadata":{"annotations":{%q:"true"}},"spec":{"deletionPolicy":"Halt"}}`, annOrphanOnDelete)
	if _, err := ri.Patch(ctx, name, types.MergePatchType, []byte(patch), metav1.PatchOptions{}); err != nil {
		return fmt.Errorf("failed to prepare Milvus %s for release: %w", name, err)
	}
	c.logf("deleting Milvus %s with %s; its pods keep running", name, annOrphanOnDelete)
	if err := ri.Delete(ctx, name, metav1.DeleteOptions{}); err != nil && !kerr.IsNotFound(err) {
		return err
	}
	return c.waitFor(ctx, timeout, "Milvus "+name+" to be released", func() (bool, string, error) {
		_, err := ri.Get(ctx, name, metav1.GetOptions{})
		if kerr.IsNotFound(err) {
			return true, "", nil
		}
		return false, "still present", err
	})
}

// createObjects creates what a plan holds, keeping any object that already
// exists (for example one edited from "plan --output-dir").
func (c *cluster) createObjects(ctx context.Context, p *adoptPlan) error {
	for _, s := range p.Objects.Secrets {
		_, err := c.kube.CoreV1().Secrets(c.ns).Create(ctx, s, metav1.CreateOptions{})
		switch {
		case kerr.IsAlreadyExists(err):
			c.logf("Secret %s exists; using it", s.Name)
		case err != nil:
			return fmt.Errorf("failed to create Secret %s: %w", s.Name, err)
		default:
			c.logf("created Secret %s", s.Name)
		}
	}
	_, err := c.kube.CoreV1().ConfigMaps(c.ns).Create(ctx, p.Objects.State, metav1.CreateOptions{})
	switch {
	case kerr.IsAlreadyExists(err):
	case err != nil:
		return fmt.Errorf("failed to create ConfigMap %s: %w", p.Objects.State.Name, err)
	default:
		c.logf("created ConfigMap %s (adoption state)", p.Objects.State.Name)
	}
	return nil
}

// removeGoverningServices deletes the chart's non-headless Services on names
// KubeDB needs headless, recording them so a rollback can recreate them.
func (c *cluster) removeGoverningServices(ctx context.Context, st *State, svcs []core.Service) error {
	for _, svc := range svcs {
		st.RemovedServices = append(st.RemovedServices, recordService(svc))
	}
	if len(svcs) > 0 {
		if err := c.saveState(ctx, st); err != nil {
			return err
		}
	}
	for _, svc := range svcs {
		err := c.kube.CoreV1().Services(c.ns).Delete(ctx, svc.Name, metav1.DeleteOptions{})
		if err != nil && !kerr.IsNotFound(err) {
			return fmt.Errorf("failed to delete Service %s: %w", svc.Name, err)
		}
		c.logf("deleted the chart's Service %s (KubeDB recreates it headless)", svc.Name)
	}
	return nil
}

func recordService(svc core.Service) ServiceRecord {
	return ServiceRecord{
		Name: svc.Name, Selector: svc.Spec.Selector, Labels: svc.Labels,
		Annotations: svc.Annotations, Ports: svc.Spec.Ports, Type: svc.Spec.Type,
	}
}

// createMilvus creates the staged Milvus and waits until the operator has
// rendered its config and created a PetSet for every Unit.
func (c *cluster) createMilvus(ctx context.Context, p *adoptPlan, timeout time.Duration) (*unstructured.Unstructured, error) {
	ri := c.dyn.Resource(MilvusGVR).Namespace(c.ns)
	db, err := ri.Create(ctx, p.Objects.Milvus, metav1.CreateOptions{})
	if kerr.IsAlreadyExists(err) {
		db, err = ri.Get(ctx, p.State.Name, metav1.GetOptions{})
	} else if err == nil {
		c.logf("created Milvus %s, every role at 0 replicas", p.State.Name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create Milvus %s: %w", p.State.Name, err)
	}
	return db, c.waitOperator(ctx, p.State, timeout)
}

func (c *cluster) waitOperator(ctx context.Context, st *State, timeout time.Duration) error {
	return c.waitFor(ctx, timeout, "the operator to create the PetSets and render the config of Milvus "+st.Name, func() (bool, string, error) {
		db, err := c.dyn.Resource(MilvusGVR).Namespace(c.ns).Get(ctx, st.Name, metav1.GetOptions{})
		if err != nil {
			return false, "", err
		}
		cfg, err := c.renderedConfig(ctx, db)
		if err != nil {
			return false, "", err
		}
		if cfg == "" {
			return false, "no rendered config yet; see kubectl describe milvus " + st.Name, nil
		}
		for _, u := range st.Units {
			if _, err := c.dyn.Resource(PetSetGVR).Namespace(c.ns).Get(ctx, petSetName(st.Name, u), metav1.GetOptions{}); kerr.IsNotFound(err) {
				return false, "PetSet " + petSetName(st.Name, u) + " not created yet", nil
			} else if err != nil {
				return false, "", err
			}
		}
		return true, "", nil
	})
}

// checkRenderedConfig compares the config the operator rendered with the one
// the cluster runs, on every setting that decides where the data is read from.
func (c *cluster) checkRenderedConfig(ctx context.Context, name string, running map[string]any) ([]string, error) {
	db, err := c.dyn.Resource(MilvusGVR).Namespace(c.ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	raw, err := c.renderedConfig(ctx, db)
	if err != nil {
		return nil, err
	}
	var rendered map[string]any
	if err := yaml.Unmarshal([]byte(raw), &rendered); err != nil {
		return nil, fmt.Errorf("failed to parse the rendered config: %w", err)
	}
	return configMismatches(running, rendered), nil
}

// waitMilvusReady waits for the Milvus to report Ready.
func (c *cluster) waitMilvusReady(ctx context.Context, name string, timeout time.Duration) error {
	return c.waitFor(ctx, timeout, "Milvus "+name+" to be Ready", func() (bool, string, error) {
		db, err := c.dyn.Resource(MilvusGVR).Namespace(c.ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, "", err
		}
		phase, _, _ := unstructured.NestedString(db.Object, "status", "phase")
		return phase == "Ready", "phase " + phase, nil
	})
}

// cleanupKubeDBObjects deletes what the released Milvus left behind (its
// PetSets, Services other than the client Service, rendered config and RBAC),
// then recreates the chart Services the adoption removed. PVCs are kept.
func (c *cluster) cleanupKubeDBObjects(ctx context.Context, st *State) error {
	sel := metav1.ListOptions{LabelSelector: labelInstance + "=" + st.Name + "," + labelManagedBy + "=kubedb.com"}
	pss, err := c.dyn.Resource(PetSetGVR).Namespace(c.ns).List(ctx, sel)
	if err != nil {
		return err
	}
	for _, ps := range pss.Items {
		if n, _, _ := unstructured.NestedInt64(ps.Object, "spec", "replicas"); n > 0 {
			return fmt.Errorf("PetSet %s still runs %d replicas; not deleting", ps.GetName(), n)
		}
		if err := c.dyn.Resource(PetSetGVR).Namespace(c.ns).Delete(ctx, ps.GetName(), metav1.DeleteOptions{}); err != nil && !kerr.IsNotFound(err) {
			return err
		}
		c.logf("deleted PetSet %s", ps.GetName())
	}
	svcs, err := c.kube.CoreV1().Services(c.ns).List(ctx, sel)
	if err != nil {
		return err
	}
	for _, svc := range svcs.Items {
		if st.Service != nil && svc.Name == st.Service.Name {
			continue
		}
		if err := c.kube.CoreV1().Services(c.ns).Delete(ctx, svc.Name, metav1.DeleteOptions{}); err != nil && !kerr.IsNotFound(err) {
			return err
		}
		c.logf("deleted Service %s", svc.Name)
	}
	secs, err := c.kube.CoreV1().Secrets(c.ns).List(ctx, sel)
	if err != nil {
		return err
	}
	for _, s := range secs.Items {
		if _, ok := s.Data["milvus.yaml"]; !ok {
			continue // credentials stay
		}
		if err := c.kube.CoreV1().Secrets(c.ns).Delete(ctx, s.Name, metav1.DeleteOptions{}); err != nil && !kerr.IsNotFound(err) {
			return err
		}
		c.logf("deleted Secret %s (rendered config)", s.Name)
	}
	// The carried config and copies of certificates: a later adoption must read
	// them afresh, because run keeps a Secret that exists. The label keeps a Secret
	// of the user's with the same name out of this.
	for _, name := range []string{configSecretName(st.Name), serverTLSSecretName(st.Name), etcdTLSSecretName(st.Name)} {
		s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, name, metav1.GetOptions{})
		if kerr.IsNotFound(err) {
			continue
		}
		if err != nil {
			return err
		}
		if s.Labels[LabelRelease] == "" {
			continue
		}
		if err := c.kube.CoreV1().Secrets(c.ns).Delete(ctx, name, metav1.DeleteOptions{}); err != nil && !kerr.IsNotFound(err) {
			return err
		}
		c.logf("deleted Secret %s (created by the adoption)", name)
	}
	for _, rec := range st.RemovedServices {
		svc := &core.Service{
			ObjectMeta: metav1.ObjectMeta{Name: rec.Name, Namespace: c.ns, Labels: rec.Labels, Annotations: rec.Annotations},
			Spec:       core.ServiceSpec{Type: rec.Type, Selector: rec.Selector, Ports: withoutNodePorts(rec.Ports)},
		}
		if _, err := c.kube.CoreV1().Services(c.ns).Create(ctx, svc, metav1.CreateOptions{}); err != nil && !kerr.IsAlreadyExists(err) {
			return fmt.Errorf("failed to recreate the chart's Service %s: %w", rec.Name, err)
		}
		c.logf("recreated the chart's Service %s", rec.Name)
	}
	return nil
}

func withoutNodePorts(in []core.ServicePort) []core.ServicePort {
	out := make([]core.ServicePort, len(in))
	for i, p := range in {
		p.NodePort = 0
		out[i] = p
	}
	return out
}

// fallbackServiceRecord rebuilds the client Service record for a Milvus that
// was adopted by hand: the chart proxies' selector, and the Service's current
// ports with the port names only KubeDB pods define replaced by numbers.
func fallbackServiceRecord(svc core.Service, proxy Unit, t core.PodTemplateSpec) *ServiceRecord {
	rec := &ServiceRecord{Name: svc.Name, Selector: proxy.Selector}
	named := map[string]bool{}
	numbers := map[int32]bool{}
	if c := milvusContainer(t); c != nil {
		for _, p := range c.Ports {
			named[p.Name] = true
			numbers[p.ContainerPort] = true
		}
	}
	for _, p := range svc.Spec.Ports {
		if p.TargetPort.Type == intstr.String && !named[p.TargetPort.StrVal] {
			if !numbers[p.Port] {
				continue // a port only the KubeDB pods serve
			}
			p.TargetPort = intstr.FromInt32(p.Port)
		}
		rec.Ports = append(rec.Ports, p)
	}
	return rec
}

func sleep(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
	}
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
