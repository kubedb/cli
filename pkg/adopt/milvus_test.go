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
	"reflect"
	"strings"
	"testing"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/yaml"
)

// The ConfigMap the upstream Milvus Helm chart (5.0.x) renders, trimmed.
const chartDefaultYAML = `
etcd:
  endpoints:
  - mv-etcd-0.mv-etcd-headless.ns.svc.cluster.local:2379
  - mv-etcd-1.mv-etcd-headless.ns.svc.cluster.local:2379
minio:
  address: minio.storage.svc
  port: 9000
  accessKeyID: ak
  secretAccessKey: sk
  useSSL: false
  bucketName: milvus-data
  rootPath: prod
  useIAM: false
  cloudProvider: aws
  iamEndpoint:
  region:
mq:
  type: woodpecker
messageQueue: woodpecker
proxy:
  port: 19530
woodpecker:
  storage:
    type: minio
    rootPath: /woodpecker/data
`

const chartUserYAML = `
common:
  security:
    authorizationEnabled: true
    defaultRootPassword: S3cret
minio:
  rootPath: prod-override
`

func chartConfig(t *testing.T) map[string]any {
	t.Helper()
	cfg, unknown, err := mergeConfigFiles(map[string]string{"default.yaml": chartDefaultYAML, "user.yaml": chartUserYAML})
	if err != nil {
		t.Fatal(err)
	}
	if len(unknown) != 0 {
		t.Fatalf("unexpected unknown files %v", unknown)
	}
	return cfg
}

func TestMergeConfigFilesUserOverridesDefault(t *testing.T) {
	cfg := chartConfig(t)
	if got := lookupString(cfg, "minio.rootPath"); got != "prod-override" {
		t.Fatalf("user.yaml must override default.yaml: rootPath = %q", got)
	}
	if got := lookupString(cfg, "minio.bucketName"); got != "milvus-data" {
		t.Fatalf("keys only in default.yaml are kept: bucketName = %q", got)
	}
	_, unknown, err := mergeConfigFiles(map[string]string{"user.yaml": "a: 1", "extra.yaml": "b: 2"})
	if err != nil || !reflect.DeepEqual(unknown, []string{"extra.yaml"}) {
		t.Fatalf("unknown files = %v, %v", unknown, err)
	}
}

func TestEtcdFromConfig(t *testing.T) {
	s := etcdFromConfig(chartConfig(t))
	want := []string{
		"http://mv-etcd-0.mv-etcd-headless.ns.svc.cluster.local:2379",
		"http://mv-etcd-1.mv-etcd-headless.ns.svc.cluster.local:2379",
	}
	if !reflect.DeepEqual(s.Endpoints, want) {
		t.Fatalf("endpoints = %v", s.Endpoints)
	}
	if s.RootPath != "by-dev" {
		t.Fatalf("an unset etcd.rootPath is the Milvus default by-dev, got %q", s.RootPath)
	}
	s = etcdFromConfig(map[string]any{"etcd": map[string]any{
		"endpoints": "a:2379, b:2379", "rootPath": "custom",
		"ssl": map[string]any{"enabled": true}, "auth": map[string]any{"enabled": "true"},
	}})
	if !reflect.DeepEqual(s.Endpoints, []string{"https://a:2379", "https://b:2379"}) || s.RootPath != "custom" || !s.AuthEnabled || !s.TLSEnabled {
		t.Fatalf("got %+v", s)
	}
}

func TestObjectStorageData(t *testing.T) {
	got := objectStorageData(chartConfig(t))
	want := map[string]string{
		"address": "minio.storage.svc", "port": "9000", "accessKeyID": "ak", "secretAccessKey": "sk",
		"useSSL": "false", "bucketName": "milvus-data", "rootPath": "prod-override", "useIAM": "false", "cloudProvider": "aws",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestCarriedConfig(t *testing.T) {
	raw, err := carriedConfig(chartConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	var cfg map[string]any
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatal(err)
	}
	if _, ok := cfg["etcd"]; ok {
		t.Fatal("the etcd section must be dropped: the Milvus object sets it")
	}
	if v, ok := lookup(cfg, "proxy.http.port"); !ok || v != "" {
		t.Fatalf("an unset proxy.http.port must be pinned to \"\" to keep REST on the gRPC port, got %v %v", v, ok)
	}
	if lookupString(cfg, "woodpecker.storage.rootPath") != "/woodpecker/data" {
		t.Fatal("the woodpecker section must be kept as it runs")
	}
	raw, _ = carriedConfig(map[string]any{"proxy": map[string]any{"http": map[string]any{"port": nil}}})
	_ = yaml.Unmarshal([]byte(raw), &cfg)
	if v, ok := lookup(cfg, "proxy.http.port"); !ok || v != "" {
		t.Fatalf("a null proxy.http.port must be pinned to \"\" as well, got %v %v", v, ok)
	}
	raw, _ = carriedConfig(map[string]any{"proxy": map[string]any{"http": map[string]any{"port": float64(8080)}}})
	if !strings.Contains(raw, "port: 8080") {
		t.Fatalf("a set proxy.http.port is kept: %s", raw)
	}
}

func TestWalType(t *testing.T) {
	for _, tc := range []struct {
		cfg  map[string]any
		want string
	}{
		{map[string]any{"mq": map[string]any{"type": "woodpecker"}}, "woodpecker"},
		{map[string]any{"mq": map[string]any{"type": "default"}, "messageQueue": "pulsar"}, "pulsar"},
		{map[string]any{"messageQueue": "kafka"}, "kafka"},
		{map[string]any{}, ""},
	} {
		if got := walType(tc.cfg); got != tc.want {
			t.Errorf("walType(%v) = %q, want %q", tc.cfg, got, tc.want)
		}
	}
}

func TestConfigMismatches(t *testing.T) {
	running := chartConfig(t)
	carried, _ := carriedConfig(running)
	var rendered map[string]any
	_ = yaml.Unmarshal([]byte(carried), &rendered)
	rendered["etcd"] = map[string]any{"rootPath": "by-dev"}
	rendered["minio"].(map[string]any)["port"] = float64(9000) // the operator writes the port as a number
	if d := configMismatches(running, rendered); len(d) != 0 {
		t.Fatalf("identical settings reported as different: %v", d)
	}

	rendered["etcd"] = map[string]any{"rootPath": "mv-etcd"}
	rendered["minio"].(map[string]any)["rootPath"] = "files"
	rendered["minio"].(map[string]any)["secretAccessKey"] = "other"
	d := configMismatches(running, rendered)
	if len(d) != 3 {
		t.Fatalf("want 3 differences, got %v", d)
	}
	for _, s := range d {
		if strings.Contains(s, "other") || strings.Contains(s, `"sk"`) {
			t.Fatalf("a credential leaked into the report: %s", s)
		}
	}
}

func ref(name, role string, replicas int32) workloadRef {
	sel := map[string]string{"app.kubernetes.io/instance": "rel", "component": role}
	t := core.PodTemplateSpec{Spec: core.PodSpec{Containers: []core.Container{{Name: role, Image: "milvusdb/milvus:v2.6.15"}}}}
	t.Labels = sel
	return workloadRef{Kind: "Deployment", Name: name, Replicas: replicas, Selector: sel, Template: t}
}

func allRoles() []workloadRef {
	return []workloadRef{
		ref("rel-milvus-proxy", "proxy", 3),
		ref("rel-milvus-mixcoord", "mixcoord", 2),
		ref("rel-milvus-datanode", "datanode", 2),
		ref("rel-milvus-querynode", "querynode", 6),
		ref("rel-milvus-streamingnode", "streamingnode", 2),
		{Kind: "StatefulSet", Name: "rel-etcd", Replicas: 3, Selector: map[string]string{"app": "etcd"}},
	}
}

func TestAssignUnits(t *testing.T) {
	o := sourceOptions{Release: "rel", RoleLabel: "component"}
	units, tmpls, fs, err := assignUnits(allRoles(), o)
	if err != nil || fs.hasErrors() {
		t.Fatalf("err %v, findings %v", err, fs)
	}
	if len(units) != 5 || len(tmpls) != 5 {
		t.Fatalf("want the 5 roles and not etcd, got %v", units)
	}
	for _, u := range units {
		if u.Group != "" {
			t.Fatalf("one workload per role needs no group: %+v", u)
		}
	}

	// Two QueryNode pools become node groups named after the workloads.
	refs := allRoles()
	refs[3] = ref("rel-milvus-querynode-l40s", "querynode", 4)
	refs = append(refs, ref("rel-milvus-querynode-h200", "querynode", 2))
	units, _, fs, err = assignUnits(refs, o)
	if err != nil || fs.hasErrors() {
		t.Fatalf("err %v, findings %v", err, fs)
	}
	var groups []string
	for _, u := range units {
		if u.Role == roleQueryNode {
			groups = append(groups, u.Group)
		}
	}
	if !reflect.DeepEqual(groups, []string{"h200", "l40s"}) {
		t.Fatalf("groups = %v", groups)
	}

	// A pre-2.6 component blocks the adoption; a missing role does too.
	refs = append(allRoles()[1:], ref("rel-milvus-rootcoord", "rootcoord", 1))
	_, _, fs, _ = assignUnits(refs, o)
	if !fs.hasErrors() || !strings.Contains(findingText(fs), "rootcoord") || !strings.Contains(findingText(fs), "proxy role") {
		t.Fatalf("findings %v", fs)
	}

	// A role at 0 replicas has nothing to hand over.
	refs = allRoles()
	refs[0].Replicas = 0
	_, _, fs, _ = assignUnits(refs, o)
	if !fs.hasErrors() {
		t.Fatal("a workload at 0 replicas must be reported")
	}

	// Explicit mapping, for charts without a role label.
	refs = allRoles()
	for i := range refs {
		delete(refs[i].Template.Labels, "component")
	}
	o.Workloads = []string{
		"proxy=Deployment/rel-milvus-proxy", "mixcoord=Deployment/rel-milvus-mixcoord", "datanode=Deployment/rel-milvus-datanode",
		"querynode/gpu=Deployment/rel-milvus-querynode", "streamingnode=Deployment/rel-milvus-streamingnode",
	}
	units, _, fs, err = assignUnits(refs, o)
	if err != nil || fs.hasErrors() || len(units) != 5 {
		t.Fatalf("explicit mapping: units %v, findings %v, err %v", units, fs, err)
	}
}

func findingText(fs findings) string {
	var b strings.Builder
	for _, f := range fs {
		b.WriteString(f.Msg + "\n")
	}
	return b.String()
}

func TestGroupName(t *testing.T) {
	for in, want := range map[string]string{
		"milvus-querynode-h200":     "h200",
		"milvus-querynode-L40S_gpu": "l40s-gpu",
		"qn-pool-a":                 "qn-pool-a",
	} {
		if got := groupName(in, roleQueryNode); got != want {
			t.Errorf("groupName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestConfigSources(t *testing.T) {
	tpl := core.PodTemplateSpec{Spec: core.PodSpec{
		Containers: []core.Container{{Name: "proxy", Image: "milvusdb/milvus:v2.6.15", VolumeMounts: []core.VolumeMount{
			{Name: "milvus-config", MountPath: "/milvus/configs/default.yaml", SubPath: "default.yaml"},
			{Name: "milvus-config", MountPath: "/milvus/configs/user.yaml", SubPath: "user.yaml"},
			{Name: "tools", MountPath: "/milvus/tools"},
		}}},
		Volumes: []core.Volume{
			{Name: "milvus-config", VolumeSource: core.VolumeSource{ConfigMap: &core.ConfigMapVolumeSource{LocalObjectReference: core.LocalObjectReference{Name: "rel-milvus"}}}},
			{Name: "tools", VolumeSource: core.VolumeSource{EmptyDir: &core.EmptyDirVolumeSource{}}},
		},
	}}
	got := configSources(tpl)
	if len(got) != 2 || got[0].File != "default.yaml" || got[1].Key != "user.yaml" || got[0].Name != "rel-milvus" {
		t.Fatalf("got %+v", got)
	}
}

func TestImageVersionAndPickVersion(t *testing.T) {
	v, gpu, err := imageVersion("milvusdb/milvus:v2.6.15-gpu")
	if err != nil || v != "2.6.15" || !gpu {
		t.Fatalf("got %s %v %v", v, gpu, err)
	}
	if _, _, err := imageVersion("registry:5000/milvus"); err == nil {
		t.Fatal("an image without a tag must fail")
	}
	mk := func(name, ver string, gpu, deprecated bool) unstructured.Unstructured {
		u := unstructured.Unstructured{Object: map[string]any{"spec": map[string]any{
			"version": ver, "deprecated": deprecated, "db": map[string]any{"gpu": map[string]any{"supported": gpu}},
		}}}
		u.SetName(name)
		return u
	}
	versions := []unstructured.Unstructured{
		mk("2.6.15", "2.6.15", false, false), mk("2.6.15-gpu", "2.6.15", true, false),
		mk("2.6.11", "2.6.11", false, false), mk("2.6.9-old", "2.6.9", false, true),
	}
	for image, want := range map[string]string{
		"milvusdb/milvus:v2.6.15":     "2.6.15",
		"milvusdb/milvus:v2.6.15-gpu": "2.6.15-gpu",
	} {
		if got, err := pickVersion(image, versions); err != nil || got != want {
			t.Errorf("pickVersion(%s) = %s, %v; want %s", image, got, err, want)
		}
	}
	for _, image := range []string{"milvusdb/milvus:v2.6.9", "milvusdb/milvus:v2.6.11-gpu"} {
		if _, err := pickVersion(image, versions); err == nil {
			t.Errorf("pickVersion(%s): a deprecated or missing build must fail", image)
		}
	}
}

func TestRenderStagedMilvus(t *testing.T) {
	refs := allRoles()
	qn := &refs[3].Template.Spec.Containers[0]
	qn.Resources = core.ResourceRequirements{
		Requests: core.ResourceList{core.ResourceCPU: resource.MustParse("8"), core.ResourceMemory: resource.MustParse("144Gi")},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("144Gi"), "nvidia.com/gpu": resource.MustParse("1"),
			"intel.com/sriov_rdma": resource.MustParse("1"),
		},
	}
	refs[3].Template.Spec.NodeSelector = map[string]string{"pool": "gpu"}
	refs[3].Template.Spec.Tolerations = []core.Toleration{{Key: "nvidia.com/gpu", Operator: core.TolerationOpExists, Effect: core.TaintEffectNoSchedule}}
	units, tmpls, _, _ := assignUnits(refs, sourceOptions{Release: "rel", RoleLabel: "component"})
	svc := &core.Service{Spec: core.ServiceSpec{Selector: map[string]string{"component": "proxy"}, Ports: []core.ServicePort{{Name: "milvus", Port: 19530}}}}
	svc.Name = "rel-milvus"
	src := &Source{Units: units, Templates: tmpls, Config: chartConfig(t), Service: svc}
	in := renderInput{
		Namespace: "ns", Name: "rel-milvus", Version: "2.6.15-gpu", Release: "rel", Source: src,
		RootUser: "root", RootPassword: "S3cret",
	}
	st := newState(in)
	if st.Service == nil || st.Service.Name != "rel-milvus" || len(st.Service.Ports) != 1 {
		t.Fatalf("the client Service must be recorded for a rollback: %+v", st.Service)
	}
	objs, err := render(in, st)
	if err != nil {
		t.Fatal(err)
	}
	if len(objs.Secrets) != 3 || objs.Secrets[0].Name != "rel-milvus-auth" || objs.Secrets[0].StringData["password"] != "S3cret" {
		t.Fatalf("secrets %v", objs.Secrets)
	}
	m := objs.Milvus.Object
	str := func(path ...string) string { s, _, _ := unstructured.NestedString(m, path...); return s }
	if str("spec", "metaStorage", "rootPath") != "by-dev" || str("spec", "authSecret", "name") != "rel-milvus-auth" ||
		str("spec", "deletionPolicy") != "DoNotTerminate" || str("metadata", "annotations", annOrphanOnDelete) != "true" {
		t.Fatalf("Milvus %v", m)
	}
	if b, _, _ := unstructured.NestedBool(m, "spec", "metaStorage", "externallyManaged"); !b {
		t.Fatal("the etcd must be externally managed")
	}
	for _, role := range milvusRoles {
		n, found, _ := unstructured.NestedInt64(m, "spec", "topology", "distributed", role, "replicas")
		if !found || n != 0 {
			t.Fatalf("%s must be staged at 0 replicas, got %d %v", role, n, found)
		}
	}
	q, _, _ := unstructured.NestedMap(m, "spec", "topology", "distributed", "querynode")
	if c, _, _ := unstructured.NestedInt64(q, "gpu", "count"); c != 1 {
		t.Fatalf("GPU count not carried: %v", q)
	}
	ctrs, _, _ := unstructured.NestedSlice(q, "podTemplate", "spec", "containers")
	lim := ctrs[0].(map[string]any)["resources"].(map[string]any)["limits"].(map[string]any)
	if lim["memory"] != "144Gi" || lim["nvidia.com/gpu"] != nil || lim["intel.com/sriov_rdma"] != nil {
		t.Fatalf("limits %v: memory is carried, GPU and SR-IOV devices are not", lim)
	}
	if ns, _, _ := unstructured.NestedStringMap(q, "podTemplate", "spec", "nodeSelector"); ns["pool"] != "gpu" {
		t.Fatalf("nodeSelector not carried: %v", q)
	}

	// Without authentication there is no auth Secret.
	in.DisableAuth = true
	objs, _ = render(in, st)
	if len(objs.Secrets) != 2 {
		t.Fatalf("want no auth Secret, got %d secrets", len(objs.Secrets))
	}
	if b, _, _ := unstructured.NestedBool(objs.Milvus.Object, "spec", "disableSecurity"); !b {
		t.Fatal("disableSecurity must be set")
	}
}

func TestRenderNodeGroups(t *testing.T) {
	refs := allRoles()
	refs[3] = ref("rel-milvus-querynode-l40s", "querynode", 4)
	refs = append(refs, ref("rel-milvus-querynode-h200", "querynode", 2))
	units, tmpls, _, _ := assignUnits(refs, sourceOptions{Release: "rel", RoleLabel: "component"})
	spec := distributedSpec(&Source{Units: units, Templates: tmpls})
	groups, _, _ := unstructured.NestedSlice(spec, "querynode", "groups")
	if len(groups) != 2 || groups[0].(map[string]any)["name"] != "h200" || groups[0].(map[string]any)["replicas"] != int64(0) {
		t.Fatalf("groups %v", groups)
	}
	u := units[3]
	if petSetName("db", u) != "db-querynode-h200" || kubedbSelector("db", u)[labelNodeGroup] != "h200" {
		t.Fatalf("group naming: %s %v", petSetName("db", u), kubedbSelector("db", u))
	}
}

func TestSwapSteps(t *testing.T) {
	if got := swapSteps(3, false); !reflect.DeepEqual(got, []step{{3, 0}}) {
		t.Fatalf("whole role: %v", got)
	}
	if got := swapSteps(3, true); !reflect.DeepEqual(got, []step{{1, 2}, {2, 1}, {3, 0}}) {
		t.Fatalf("one by one: %v", got)
	}
}

func TestHandoverOrder(t *testing.T) {
	if o := handoverOrder(methodOverlap); o[0] != roleProxy || o[len(o)-1] != roleQueryNode {
		t.Fatalf("overlap moves the proxies first and the QueryNodes last: %v", o)
	}
	if o := handoverOrder(methodStopStart); o[0] != roleMixCoord || o[len(o)-1] != roleProxy {
		t.Fatalf("stop-start moves the coordinators first and the proxies last: %v", o)
	}
	if len(handoverOrder(methodOverlap)) != len(milvusRoles) || len(rollbackOrder) != len(milvusRoles) {
		t.Fatal("every role is handed over and rolled back")
	}
}

func TestStateDone(t *testing.T) {
	st := &State{}
	u := Unit{Role: roleQueryNode, Group: "h200"}
	st.markDone(u)
	st.markDone(u)
	if !st.isDone(u) || len(st.Done) != 1 || st.Done[0] != "querynode/h200" {
		t.Fatalf("done %v", st.Done)
	}
	st.unmarkDone(u)
	if st.isDone(u) {
		t.Fatal("unmarkDone")
	}
}

func TestFallbackServiceRecord(t *testing.T) {
	// A client Service KubeDB took over: its ports target the KubeDB container's
	// port names. The chart's proxy names its gRPC port "milvus".
	svc := core.Service{Spec: core.ServiceSpec{Ports: []core.ServicePort{
		{Name: "proxy", Port: 19530, TargetPort: intstr.FromString("proxy")},
		{Name: "metrics", Port: 9091, TargetPort: intstr.FromString("metrics")},
		{Name: "http", Port: 8080, TargetPort: intstr.FromString("http")},
		{Name: "admin", Port: 9999, TargetPort: intstr.FromInt32(9999)},
	}}}
	svc.Name = "rel-milvus"
	proxy := Unit{Role: roleProxy, Selector: map[string]string{"component": "proxy"}}
	tpl := core.PodTemplateSpec{Spec: core.PodSpec{Containers: []core.Container{{
		Name: "proxy", Image: "milvusdb/milvus:v2.6.15",
		Ports: []core.ContainerPort{{Name: "milvus", ContainerPort: 19530}, {Name: "metrics", ContainerPort: 9091}},
	}}}}
	rec := fallbackServiceRecord(svc, proxy, tpl)
	if !reflect.DeepEqual(rec.Selector, proxy.Selector) {
		t.Fatalf("selector %v", rec.Selector)
	}
	got := make([]string, 0, len(rec.Ports))
	for _, p := range rec.Ports {
		got = append(got, p.Name+"->"+p.TargetPort.String())
	}
	want := []string{"proxy->19530", "metrics->metrics", "admin->9999"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ports %v, want %v (8080 only exists on KubeDB pods)", got, want)
	}
}

func TestRestPort(t *testing.T) {
	if restPort(map[string]any{}) != 19530 {
		t.Fatal("default")
	}
	if restPort(map[string]any{"proxy": map[string]any{"http": map[string]any{"port": ""}, "port": float64(19531)}}) != 19531 {
		t.Fatal("an empty http port serves REST on the gRPC port")
	}
	if restPort(map[string]any{"proxy": map[string]any{"http": map[string]any{"port": float64(8080)}}}) != 8080 {
		t.Fatal("a set http port")
	}
}
