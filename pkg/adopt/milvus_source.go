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
	"path"
	"regexp"
	"sort"
	"strings"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Distributed Milvus 2.6 roles, in the order the CR lists them.
const (
	roleProxy         = "proxy"
	roleMixCoord      = "mixcoord"
	roleDataNode      = "datanode"
	roleQueryNode     = "querynode"
	roleStreamingNode = "streamingnode"
)

var milvusRoles = []string{roleProxy, roleMixCoord, roleDataNode, roleQueryNode, roleStreamingNode}

// Components of Milvus before 2.6 (separate coordinators, IndexNode) and the
// Standalone mode. A cluster with any of them is not a 2.6 Distributed cluster.
var unsupportedComponents = []string{"rootcoord", "querycoord", "datacoord", "indexcoord", "indexnode", "standalone"}

// Unit is one workload of the running cluster and the KubeDB PetSet that
// replaces it: a role, or one node group of a role when the role runs as
// several workloads (for example QueryNodes on two GPU pools).
type Unit struct {
	Role     string            `json:"role"`
	Group    string            `json:"group,omitempty"`
	Kind     string            `json:"kind"`
	Name     string            `json:"name"`
	Replicas int32             `json:"replicas"`
	Selector map[string]string `json:"selector"`
}

func (u Unit) key() string {
	if u.Group == "" {
		return u.Role
	}
	return u.Role + "/" + u.Group
}

// ServiceRecord is the client Service as the chart left it, so a rollback can
// give it back unchanged.
type ServiceRecord struct {
	Name        string             `json:"name"`
	Type        core.ServiceType   `json:"type,omitempty"`
	Selector    map[string]string  `json:"selector"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Ports       []core.ServicePort `json:"ports"`
}

// Phases of an adoption, recorded in the state ConfigMap.
const (
	phasePlanned     = "Planned"
	phaseHandingOver = "HandingOver"
	phaseAdopted     = "Adopted"
	phaseRolledBack  = "RolledBack"
)

// State is the record of one adoption, kept in the ConfigMap <name>-adopt so
// that run can resume and rollback knows what to restore.
type State struct {
	Release string         `json:"release"`
	Name    string         `json:"name"`
	Method  string         `json:"method,omitempty"`
	Phase   string         `json:"phase"`
	Units   []Unit         `json:"units"`
	Service *ServiceRecord `json:"service,omitempty"`
	// RemovedServices are chart Services deleted because KubeDB needs their
	// names for headless Services; a rollback recreates them.
	RemovedServices []ServiceRecord `json:"removedServices,omitempty"`
	Done            []string        `json:"done,omitempty"`
}

func (s *State) unitsOf(role string) []Unit {
	var out []Unit
	for _, u := range s.Units {
		if u.Role == role {
			out = append(out, u)
		}
	}
	return out
}

func (s *State) isDone(u Unit) bool { return contains(s.Done, u.key()) }

func (s *State) markDone(u Unit) {
	if !s.isDone(u) {
		s.Done = append(s.Done, u.key())
	}
}

func (s *State) unmarkDone(u Unit) {
	out := s.Done[:0]
	for _, k := range s.Done {
		if k != u.key() {
			out = append(out, k)
		}
	}
	s.Done = out
}

// Source is what discovery found about the running cluster.
type Source struct {
	Units     []Unit
	Templates map[string]core.PodTemplateSpec // by Unit.key()
	// Networks is spec.network of each unit, mapped from the chart pod's Multus annotation (by Unit.key()).
	Networks   map[string]map[string]any
	Config     map[string]any
	ConfigFrom []string
	Image      string
	Service    *core.Service
	Findings   findings
}

type workloadRef struct {
	Kind     string
	Name     string
	Labels   map[string]string
	Replicas int32
	Selector map[string]string
	Template core.PodTemplateSpec
}

func refFromDeployment(d apps.Deployment) workloadRef {
	r := workloadRef{Kind: "Deployment", Name: d.Name, Labels: d.Labels, Replicas: 1, Template: d.Spec.Template}
	if d.Spec.Replicas != nil {
		r.Replicas = *d.Spec.Replicas
	}
	if d.Spec.Selector != nil {
		r.Selector = d.Spec.Selector.MatchLabels
	}
	return r
}

func refFromStatefulSet(s apps.StatefulSet) workloadRef {
	r := workloadRef{Kind: "StatefulSet", Name: s.Name, Labels: s.Labels, Replicas: 1, Template: s.Spec.Template}
	if s.Spec.Replicas != nil {
		r.Replicas = *s.Spec.Replicas
	}
	if s.Spec.Selector != nil {
		r.Selector = s.Spec.Selector.MatchLabels
	}
	return r
}

// sourceOptions select the running cluster.
type sourceOptions struct {
	Release   string
	Selector  string   // label selector of the release's workloads
	RoleLabel string   // label that names a workload's role
	Workloads []string // explicit role[/group]=Kind/name, instead of the selector
}

func (o sourceOptions) selector() string {
	if o.Selector != "" {
		return o.Selector
	}
	return "app.kubernetes.io/instance=" + o.Release
}

// assignUnits turns the release's workloads into Units. Workloads of one role
// become node groups named after the workload.
func assignUnits(refs []workloadRef, o sourceOptions) ([]Unit, map[string]core.PodTemplateSpec, findings, error) {
	var fs findings
	explicit := map[string]string{} // Kind/name -> role[/group]
	for _, w := range o.Workloads {
		k, v, ok := strings.Cut(w, "=")
		if !ok {
			return nil, nil, nil, fmt.Errorf("--workload %q: want role[/group]=Kind/name", w)
		}
		explicit[strings.ToLower(v)] = k
	}
	byRole := map[string][]workloadRef{}
	groupOf := map[string]string{}
	for _, r := range refs {
		var role, group string
		if len(explicit) > 0 {
			rg, ok := explicit[strings.ToLower(r.Kind+"/"+r.Name)]
			if !ok {
				continue
			}
			role, group, _ = strings.Cut(rg, "/")
		} else {
			role = r.Template.Labels[o.RoleLabel]
			if role == "" {
				role = r.Labels[o.RoleLabel]
			}
		}
		switch {
		case role == "":
			continue // etcd, MinIO and other workloads of the release
		case contains(unsupportedComponents, role):
			fs.errorf("%s %s is a %q component: only a Distributed Milvus 2.6 (MixCoord, Proxy, DataNode, QueryNode, StreamingNode) can be adopted", r.Kind, r.Name, role)
			continue
		case !contains(milvusRoles, role):
			fs.infof("%s %s has %s=%s, not a Milvus role; ignored", r.Kind, r.Name, o.RoleLabel, role)
			continue
		}
		byRole[role] = append(byRole[role], r)
		if group != "" {
			groupOf[r.Kind+"/"+r.Name] = group
		}
	}
	if len(explicit) > 0 {
		found := 0
		for _, rs := range byRole {
			found += len(rs)
		}
		if found != len(explicit) {
			return nil, nil, nil, fmt.Errorf("found %d of the %d workloads given with --workload", found, len(explicit))
		}
	}

	var units []Unit
	templates := map[string]core.PodTemplateSpec{}
	for _, role := range milvusRoles {
		rs := byRole[role]
		if len(rs) == 0 {
			fs.errorf("no workload found for the %s role (selector %q, role label %q)", role, o.selector(), o.RoleLabel)
			continue
		}
		sort.Slice(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
		used := map[string]bool{}
		for _, r := range rs {
			u := Unit{Role: role, Kind: r.Kind, Name: r.Name, Replicas: r.Replicas, Selector: r.Selector}
			if len(rs) > 1 {
				g := groupOf[r.Kind+"/"+r.Name]
				if g == "" {
					g = groupName(r.Name, role)
				}
				if used[g] {
					return nil, nil, nil, fmt.Errorf("workloads of role %s map to the same group %q; name them with --workload", role, g)
				}
				used[g] = true
				u.Group = g
			}
			if len(u.Selector) == 0 {
				return nil, nil, nil, fmt.Errorf("%s %s has no matchLabels selector", r.Kind, r.Name)
			}
			if u.Replicas == 0 {
				fs.errorf("%s %s runs 0 replicas: nothing to hand over. If an adoption already moved it, run \"status\"", r.Kind, r.Name)
			}
			units = append(units, u)
			templates[u.key()] = r.Template
		}
	}
	return units, templates, fs, nil
}

var dnsUnsafe = regexp.MustCompile(`[^a-z0-9-]+`)

// groupName derives a node group name from a workload name, e.g.
// "milvus-querynode-h200" -> "h200".
func groupName(workload, role string) string {
	n := strings.ToLower(workload)
	if i := strings.LastIndex(n, role+"-"); i >= 0 && i+len(role)+1 < len(n) {
		n = n[i+len(role)+1:]
	}
	n = strings.Trim(dnsUnsafe.ReplaceAllString(n, "-"), "-")
	if len(n) > 20 {
		n = strings.Trim(n[len(n)-20:], "-")
	}
	return n
}

// milvusContainer returns the Milvus container of a pod template.
func milvusContainer(t core.PodTemplateSpec) *core.Container {
	for i := range t.Spec.Containers {
		c := &t.Spec.Containers[i]
		if strings.Contains(c.Image, "milvus") {
			return c
		}
	}
	if len(t.Spec.Containers) > 0 {
		return &t.Spec.Containers[0]
	}
	return nil
}

// configMountDir is where Milvus reads its configuration files.
const configMountDir = "/milvus/configs"

// configSource is one configuration file mounted into the Milvus container.
type configSource struct {
	File   string // milvus.yaml, default.yaml, user.yaml
	Kind   string // ConfigMap or Secret
	Name   string
	Key    string // "" for every key of a directory mount
	IsDir  bool
	Volume string
}

// configSources lists the configuration files a pod template mounts.
func configSources(t core.PodTemplateSpec) []configSource {
	c := milvusContainer(t)
	if c == nil {
		return nil
	}
	vols := map[string]core.Volume{}
	for _, v := range t.Spec.Volumes {
		vols[v.Name] = v
	}
	var out []configSource
	for _, m := range c.VolumeMounts {
		v, ok := vols[m.Name]
		if !ok {
			continue
		}
		src := configSource{Volume: m.Name}
		switch {
		case v.ConfigMap != nil:
			src.Kind, src.Name = "ConfigMap", v.ConfigMap.Name
		case v.Secret != nil:
			src.Kind, src.Name = "Secret", v.Secret.SecretName
		default:
			continue
		}
		mp := strings.TrimSuffix(m.MountPath, "/")
		switch {
		case mp == configMountDir:
			src.IsDir = true
		case path.Dir(mp) == configMountDir:
			src.File = path.Base(mp)
			src.Key = m.SubPath
			if src.Key == "" {
				continue // a directory mounted over a file name; not a config file
			}
		default:
			continue
		}
		out = append(out, src)
	}
	return out
}

// discover reads the release's workloads, configuration, image and client
// Service. Problems that block an adoption are returned as error findings.
func discover(ctx context.Context, kc kubernetes.Interface, ns string, o sourceOptions) (*Source, error) {
	sel := o.selector()
	if _, err := labels.Parse(sel); err != nil {
		return nil, fmt.Errorf("invalid selector %q: %w", sel, err)
	}
	lo := metav1.ListOptions{LabelSelector: sel}
	if len(o.Workloads) > 0 {
		lo = metav1.ListOptions{}
	}
	var refs []workloadRef
	deps, err := kc.AppsV1().Deployments(ns).List(ctx, lo)
	if err != nil {
		return nil, fmt.Errorf("failed to list Deployments: %w", err)
	}
	for _, d := range deps.Items {
		refs = append(refs, refFromDeployment(d))
	}
	sts, err := kc.AppsV1().StatefulSets(ns).List(ctx, lo)
	if err != nil {
		return nil, fmt.Errorf("failed to list StatefulSets: %w", err)
	}
	for _, s := range sts.Items {
		refs = append(refs, refFromStatefulSet(s))
	}
	units, templates, fs, err := assignUnits(refs, o)
	if err != nil {
		return nil, err
	}
	src := &Source{Units: units, Templates: templates, Findings: fs}
	if src.Findings.hasErrors() {
		return src, nil
	}

	proxy := templates[src.firstUnit(roleProxy).key()]
	if c := milvusContainer(proxy); c != nil {
		src.Image = c.Image
	}
	if err := src.readConfig(ctx, kc, ns, proxy); err != nil {
		return nil, err
	}
	if err := src.findClientService(ctx, kc, ns, proxy); err != nil {
		return nil, err
	}
	src.checkTemplates()
	return src, nil
}

func (s *Source) firstUnit(role string) Unit {
	for _, u := range s.Units {
		if u.Role == role {
			return u
		}
	}
	return Unit{}
}

func (s *Source) readConfig(ctx context.Context, kc kubernetes.Interface, ns string, t core.PodTemplateSpec) error {
	files := map[string]string{}
	for _, src := range configSources(t) {
		var data map[string]string
		switch src.Kind {
		case "ConfigMap":
			cm, err := kc.CoreV1().ConfigMaps(ns).Get(ctx, src.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to read the Milvus config ConfigMap %s: %w", src.Name, err)
			}
			data = cm.Data
		case "Secret":
			sec, err := kc.CoreV1().Secrets(ns).Get(ctx, src.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to read the Milvus config Secret %s: %w", src.Name, err)
			}
			data = map[string]string{}
			for k, v := range sec.Data {
				data[k] = string(v)
			}
		}
		if src.IsDir {
			for k, v := range data {
				if strings.HasSuffix(k, ".yaml") {
					files[k] = v
					s.ConfigFrom = append(s.ConfigFrom, fmt.Sprintf("%s %s key %s", src.Kind, src.Name, k))
				}
			}
			continue
		}
		v, ok := data[src.Key]
		if !ok {
			return fmt.Errorf("%s %s has no key %q mounted as %s", src.Kind, src.Name, src.Key, src.File)
		}
		files[src.File] = v
		s.ConfigFrom = append(s.ConfigFrom, fmt.Sprintf("%s %s key %s as %s", src.Kind, src.Name, src.Key, src.File))
	}
	if len(files) == 0 {
		s.Findings.errorf("found no configuration file mounted in %s; pass the effective milvus.yaml with --config-file", configMountDir)
		s.Config = map[string]any{}
		return nil
	}
	cfg, unknown, err := mergeConfigFiles(files)
	if err != nil {
		return err
	}
	for _, f := range unknown {
		s.Findings.warnf("%s/%s is mounted but Milvus does not read it; not carried over", configMountDir, f)
	}
	s.Config = cfg
	return nil
}

// findClientService finds the Service whose selector picks the proxy pods.
func (s *Source) findClientService(ctx context.Context, kc kubernetes.Interface, ns string, proxy core.PodTemplateSpec) error {
	svcs, err := kc.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list Services: %w", err)
	}
	var matches []core.Service
	for _, svc := range svcs.Items {
		if len(svc.Spec.Selector) == 0 || svc.Spec.ClusterIP == core.ClusterIPNone {
			continue
		}
		if labels.SelectorFromSet(svc.Spec.Selector).Matches(labels.Set(proxy.Labels)) {
			matches = append(matches, svc)
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		pi, pj := hasPort(matches[i], 19530), hasPort(matches[j], 19530)
		if pi != pj {
			return pi
		}
		return matches[i].Name < matches[j].Name
	})
	switch len(matches) {
	case 0:
		s.Findings.warnf("no Service selects the proxy pods; clients will need the new Service <name>")
	case 1:
	default:
		var names []string
		for _, m := range matches[1:] {
			names = append(names, m.Name)
		}
		s.Findings.warnf("Services %s also select the proxy pods; they keep pointing at the chart proxies and stop working once those are gone",
			strings.Join(names, ", "))
	}
	if len(matches) > 0 {
		s.Service = &matches[0]
	}
	return nil
}

func hasPort(svc core.Service, port int32) bool {
	for _, p := range svc.Spec.Ports {
		if p.Port == port {
			return true
		}
	}
	return false
}

// checkTemplates reports pod settings that the Milvus object does not carry over
// by itself.
func (s *Source) checkTemplates() {
	for _, u := range s.Units {
		t := s.Templates[u.key()]
		id := u.Kind + " " + u.Name
		if t.Spec.Affinity != nil {
			s.Findings.warnf("%s sets affinity, which is not carried over; use a PlacementPolicy (spec.topology.distributed.%s.podTemplate.spec.podPlacementPolicy)", id, u.Role)
		}
		for _, ic := range t.Spec.InitContainers {
			s.Findings.warnf("%s has init container %q, which is not carried over (on a role with an SR-IOV network, KubeDB generates its own init container)", id, ic.Name)
		}
		if c := milvusContainer(t); c != nil {
			for _, e := range c.Env {
				s.Findings.warnf("%s sets environment variable %s; Milvus may read it as configuration, and it is not carried over", id, e.Name)
			}
			if c.Image != s.Image {
				s.Findings.errorf("%s runs image %s, the proxy runs %s; adopt a cluster that runs one version", id, c.Image, s.Image)
			}
		}
	}
}
