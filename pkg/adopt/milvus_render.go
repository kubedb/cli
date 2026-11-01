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
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	MilvusGVR        = schema.GroupVersionResource{Group: "kubedb.com", Version: "v1alpha2", Resource: "milvuses"}
	MilvusVersionGVR = schema.GroupVersionResource{Group: "catalog.kubedb.com", Version: "v1alpha1", Resource: "milvusversions"}
	PetSetGVR        = schema.GroupVersionResource{Group: "apps.k8s.appscode.com", Version: "v1", Resource: "petsets"}
)

const (
	// LabelRelease marks the objects this command creates for one adoption.
	LabelRelease = "adopt.kubedb.com/release"
	stateKey     = "state.json"

	annOrphanOnDelete = "kubedb.com/milvus-orphan-on-delete"
	labelRole         = "kubedb.com/role"
	labelNodeGroup    = "kubedb.com/node-group"
	labelInstance     = "app.kubernetes.io/instance"
	labelManagedBy    = "app.kubernetes.io/managed-by"

	milvusContainerName = "milvus"
)

func stateConfigMapName(name string) string      { return name + "-adopt" }
func authSecretName(name string) string          { return name + "-auth" }
func objectStorageSecretName(name string) string { return name + "-object-storage" }
func serverTLSSecretName(name string) string     { return name + "-tls-server" }
func etcdTLSSecretName(name string) string       { return name + "-etcd-tls" }
func configSecretName(name string) string        { return name + "-adopt-config" }

// petSetName is the KubeDB PetSet that replaces a Unit.
func petSetName(name string, u Unit) string {
	n := name + "-" + u.Role
	if u.Group != "" {
		n += "-" + u.Group
	}
	return n
}

// kubedbSelector selects the KubeDB pods of a Unit.
func kubedbSelector(name string, u Unit) map[string]string {
	sel := map[string]string{labelInstance: name, labelRole: u.Role}
	if u.Group != "" {
		sel[labelNodeGroup] = u.Group
	}
	return sel
}

// adoptObjects are the objects an adoption creates, in creation order.
type adoptObjects struct {
	Secrets []*core.Secret
	State   *core.ConfigMap
	Milvus  *unstructured.Unstructured
}

// renderInput is everything the objects are built from.
type renderInput struct {
	Namespace      string
	Name           string
	Version        string
	Release        string
	Source         *Source
	AuthSecret     string // existing Secret to use; "" with RootUser set creates one
	RootUser       string
	RootPassword   string
	DisableAuth    bool
	EtcdAuthSecret string
	EtcdTLS        *tlsMaterial
	TLS            *adoptionTLS
	RestrictedSCC  bool
}

func newState(in renderInput) *State {
	st := &State{Release: in.Release, Name: in.Name, Phase: phasePlanned, Units: in.Source.Units}
	if svc := in.Source.Service; svc != nil && svc.Name == in.Name {
		rec := recordService(*svc)
		st.Service = &rec
	}
	return st
}

func render(in renderInput, st *State) (*adoptObjects, error) {
	cfg := in.Source.Config
	objs := &adoptObjects{}
	lbl := map[string]string{LabelRelease: in.Release}

	authName := in.AuthSecret
	if !in.DisableAuth && authName == "" {
		authName = authSecretName(in.Name)
		objs.Secrets = append(objs.Secrets, &core.Secret{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
			ObjectMeta: metav1.ObjectMeta{Name: authName, Namespace: in.Namespace, Labels: lbl},
			Type:       core.SecretTypeBasicAuth,
			StringData: map[string]string{core.BasicAuthUsernameKey: in.RootUser, core.BasicAuthPasswordKey: in.RootPassword},
		})
	}
	objs.Secrets = append(objs.Secrets, &core.Secret{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
		ObjectMeta: metav1.ObjectMeta{Name: objectStorageSecretName(in.Name), Namespace: in.Namespace, Labels: lbl},
		Type:       core.SecretTypeOpaque,
		StringData: objectStorageData(cfg),
	})
	carried, err := carriedConfig(cfg)
	if err != nil {
		return nil, err
	}
	objs.Secrets = append(objs.Secrets, &core.Secret{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
		ObjectMeta: metav1.ObjectMeta{Name: configSecretName(in.Name), Namespace: in.Namespace, Labels: lbl},
		Type:       core.SecretTypeOpaque,
		StringData: map[string]string{"milvus.yaml": carried},
	})

	raw, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}
	objs.State = &core.ConfigMap{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{Name: stateConfigMapName(in.Name), Namespace: in.Namespace, Labels: lbl},
		Data:       map[string]string{stateKey: string(raw)},
	}

	// secretFor returns the Secret that holds m: the running one when KubeDB can
	// read it as it is, otherwise a copy that rollback removes again.
	secretFor := func(m *tlsMaterial, copyName string) string {
		if m.Secret != "" {
			return m.Secret
		}
		objs.Secrets = append(objs.Secrets, &core.Secret{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
			ObjectMeta: metav1.ObjectMeta{Name: copyName, Namespace: in.Namespace, Labels: lbl},
			Type:       core.SecretTypeOpaque,
			StringData: m.secretData(),
		})
		return copyName
	}

	etcd := etcdFromConfig(cfg)
	meta := map[string]any{
		"externallyManaged": true,
		"endpoints":         toAnySlice(etcd.Endpoints),
		"rootPath":          etcd.RootPath,
	}
	if in.EtcdAuthSecret != "" {
		meta["authSecret"] = map[string]any{"name": in.EtcdAuthSecret}
	}
	if in.EtcdTLS != nil {
		meta["clientTLSSecret"] = map[string]any{"name": secretFor(in.EtcdTLS, etcdTLSSecretName(in.Name))}
	}
	spec := map[string]any{
		"version":        in.Version,
		"deletionPolicy": "DoNotTerminate",
		"autoOps":        map[string]any{"disabled": true},
		"healthChecker":  map[string]any{"disableWriteCheck": true},
		"metaStorage":    meta,
		"objectStorage":  map[string]any{"configSecret": map[string]any{"name": objectStorageSecretName(in.Name)}},
		"configuration":  map[string]any{"secretName": configSecretName(in.Name)},
		"topology": map[string]any{
			"mode":        "Distributed",
			"distributed": distributedSpec(in.Source),
		},
	}
	if in.DisableAuth {
		spec["disableSecurity"] = true
	} else {
		spec["authSecret"] = map[string]any{"name": authName}
	}
	if in.RestrictedSCC {
		spec["restrictedSCC"] = true
	}
	if in.TLS != nil {
		spec["tls"] = tlsSpec(in.TLS, secretFor(in.TLS.Server, serverTLSSecretName(in.Name)))
	}
	objs.Milvus = &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "kubedb.com/v1alpha2",
		"kind":       "Milvus",
		"metadata": map[string]any{
			"name":        in.Name,
			"namespace":   in.Namespace,
			"labels":      map[string]any{LabelRelease: in.Release},
			"annotations": map[string]any{annOrphanOnDelete: "true"},
		},
		"spec": spec,
	}}
	return objs, nil
}

// distributedSpec stages every role at 0 replicas, with the chart's resources,
// node selector, tolerations and GPU count. A role with several workloads
// becomes node groups.
func distributedSpec(src *Source) map[string]any {
	out := map[string]any{}
	for _, role := range milvusRoles {
		var units []Unit
		for _, u := range src.Units {
			if u.Role == role {
				units = append(units, u)
			}
		}
		if len(units) == 0 {
			continue
		}
		if len(units) == 1 && units[0].Group == "" {
			node := nodeSpec(src.Templates[units[0].key()])
			node["replicas"] = int64(0)
			out[role] = node
			continue
		}
		groups := make([]any, 0, len(units))
		for _, u := range units {
			g := nodeSpec(src.Templates[u.key()])
			g["name"] = u.Group
			g["replicas"] = int64(0)
			groups = append(groups, g)
		}
		out[role] = map[string]any{"groups": groups}
	}
	return out
}

// Resources a KubeDB Milvus container takes from the chart as they are. GPUs
// go to the gpu field; SR-IOV and RDMA devices come from the network field.
func carriedResource(name core.ResourceName) bool {
	n := string(name)
	return name == core.ResourceCPU || name == core.ResourceMemory || name == core.ResourceEphemeralStorage ||
		strings.HasPrefix(n, core.ResourceHugePagesPrefix)
}

func gpuResource(name core.ResourceName) bool {
	return strings.HasSuffix(string(name), "/gpu")
}

func nodeSpec(t core.PodTemplateSpec) map[string]any {
	node := map[string]any{}
	podSpec := map[string]any{}
	c := milvusContainer(t)
	if c != nil {
		res := map[string]any{}
		var gpuName core.ResourceName
		var gpuCount int64
		for _, part := range []struct {
			key  string
			list core.ResourceList
		}{{"requests", c.Resources.Requests}, {"limits", c.Resources.Limits}} {
			m := map[string]any{}
			for _, n := range sortedResourceNames(part.list) {
				q := part.list[n]
				switch {
				case carriedResource(n):
					m[string(n)] = q.String()
				case gpuResource(n):
					gpuName, gpuCount = n, q.Value()
				}
			}
			if len(m) > 0 {
				res[part.key] = m
			}
		}
		if len(res) > 0 {
			podSpec["containers"] = []any{map[string]any{"name": milvusContainerName, "resources": res}}
		}
		if gpuCount > 0 {
			gpu := map[string]any{"count": gpuCount}
			if gpuName != "nvidia.com/gpu" {
				gpu["resourceName"] = string(gpuName)
			}
			node["gpu"] = gpu
		}
	}
	if len(t.Spec.NodeSelector) > 0 {
		ns := map[string]any{}
		for k, v := range t.Spec.NodeSelector {
			ns[k] = v
		}
		podSpec["nodeSelector"] = ns
	}
	if len(t.Spec.Tolerations) > 0 {
		tols := make([]any, 0, len(t.Spec.Tolerations))
		for _, tol := range t.Spec.Tolerations {
			b, _ := json.Marshal(tol)
			var m map[string]any
			_ = json.Unmarshal(b, &m)
			tols = append(tols, m)
		}
		podSpec["tolerations"] = tols
	}
	if len(podSpec) > 0 {
		node["podTemplate"] = map[string]any{"spec": podSpec}
	}
	return node
}

func sortedResourceNames(l core.ResourceList) []core.ResourceName {
	out := make([]core.ResourceName, 0, len(l))
	for n := range l {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func toAnySlice(in []string) []any {
	out := make([]any, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}

// imageVersion splits a Milvus image into its version and whether it is a GPU
// build: "milvusdb/milvus:v2.6.15-gpu" -> ("2.6.15", true).
func imageVersion(image string) (string, bool, error) {
	if strings.Contains(image, "@") {
		return "", false, fmt.Errorf("image %s is pinned by digest; pass the version with --version", image)
	}
	i := strings.LastIndex(image, ":")
	if i < 0 || strings.Contains(image[i:], "/") {
		return "", false, fmt.Errorf("image %s has no tag; pass the version with --version", image)
	}
	tag := strings.TrimPrefix(image[i+1:], "v")
	gpu := strings.Contains(tag, "-gpu")
	ver, _, _ := strings.Cut(tag, "-")
	return ver, gpu, nil
}

// pickVersion returns the MilvusVersion for an image: same version, same GPU
// support, not deprecated.
func pickVersion(image string, versions []unstructured.Unstructured) (string, error) {
	ver, gpu, err := imageVersion(image)
	if err != nil {
		return "", err
	}
	var names []string
	for _, v := range versions {
		sv, _, _ := unstructured.NestedString(v.Object, "spec", "version")
		g, _, _ := unstructured.NestedBool(v.Object, "spec", "db", "gpu", "supported")
		dep, _, _ := unstructured.NestedBool(v.Object, "spec", "deprecated")
		if sv == ver && g == gpu && !dep {
			names = append(names, v.GetName())
		}
	}
	if len(names) == 0 {
		build := "CPU"
		if gpu {
			build = "GPU"
		}
		return "", fmt.Errorf("no MilvusVersion for %s (%s build of %s); install the KubeDB catalog that has it, or pass --version", image, build, ver)
	}
	sort.Strings(names)
	for _, n := range names { // prefer the canonical name
		if n == ver || n == ver+"-gpu" {
			return n, nil
		}
	}
	return names[0], nil
}
