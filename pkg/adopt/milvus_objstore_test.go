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
	"strings"
	"testing"
	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestParseCARef(t *testing.T) {
	for in, want := range map[string][2]string{"ca": {"ca", "ca.crt"}, "ca/root.pem": {"ca", "root.pem"}} {
		n, k, err := parseCARef(in)
		if err != nil || n != want[0] || k != want[1] {
			t.Fatalf("%q: %s %s %v", in, n, k, err)
		}
	}
	for _, bad := range []string{"", "/key", "a/b/c"} {
		if _, _, err := parseCARef(bad); err == nil {
			t.Fatalf("%q must be rejected", bad)
		}
	}
}

func TestValidateCABundle(t *testing.T) {
	ca := newTestCA(t)
	subject, notAfter, err := validateCABundle(ca.pem, time.Now())
	if err != nil || !strings.Contains(subject, "test-ca") || notAfter.IsZero() {
		t.Fatalf("%q %v %v", subject, notAfter, err)
	}
	if _, _, err := validateCABundle([]byte("not a certificate"), time.Now()); err == nil || !strings.Contains(err.Error(), "no PEM certificate") {
		t.Fatalf("garbage: %v", err)
	}
	if _, _, err := validateCABundle(ca.pem, time.Now().Add(10*365*24*time.Hour)); err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expired: %v", err)
	}
}

func TestObjectStorageTLSConfigHelpers(t *testing.T) {
	cfg := chartConfig(t)
	cfg["minio"].(map[string]any)["ssl"] = map[string]any{"tlsCACert": "/old/ca.pem"}
	var none *objectStorageTLS
	if got := none.storageData(map[string]string{"useSSL": "false"}); got["useSSL"] != "false" {
		t.Fatal("no TLS flags must change nothing")
	}
	if none.carry(cfg) == nil || lookupString(none.carry(cfg), "minio.useSSL") != "false" {
		t.Fatal("no TLS flags must keep the carried config")
	}

	tls := &objectStorageTLS{Kind: "Secret", Name: "s3-ca", Key: "ca.crt", Address: "rgw.example.svc", Port: 443}
	data := tls.storageData(objectStorageData(cfg))
	if data["useSSL"] != "true" || data["address"] != "rgw.example.svc" || data["port"] != "443" || data["bucketName"] != "milvus-data" {
		t.Fatalf("storage data %v", data)
	}
	carried := tls.carry(cfg)
	for _, p := range []string{"minio.useSSL", "minio.ssl", "minio.address", "minio.port"} {
		if _, ok := lookup(carried, p); ok {
			t.Fatalf("%s must not be carried: the operator renders it, and spec.configuration wins a merge", p)
		}
	}
	if lookupString(carried, "minio.bucketName") != "milvus-data" || lookupString(carried, "minio.rootPath") != "prod-override" {
		t.Fatalf("the rest of minio must be carried: %v", carried["minio"])
	}
	if _, ok := lookup(cfg, "minio.ssl"); !ok {
		t.Fatal("carry must not change the running config")
	}

	// Only the CA: the endpoint stays carried.
	caOnly := &objectStorageTLS{Kind: "ConfigMap", Name: "ca", Key: "ca.crt"}
	if lookupString(caOnly.carry(cfg), "minio.address") != "minio.storage.svc" || lookupString(caOnly.carry(cfg), "minio.port") != "9000" {
		t.Fatal("without --object-storage-address/port the running endpoint stays")
	}

	// The check expects what the operator renders after the switch, and still
	// reports a change nobody asked for.
	rendered := deepCopy(tls.expected(cfg))
	exp := tls.expected(cfg)
	rendered["etcd"] = map[string]any{"rootPath": lookupString(cfg, "etcd.rootPath")}
	if got := withoutPrefix(configMismatches(exp, rendered), "etcd.rootPath"); len(got) != 0 {
		t.Fatalf("unexpected differences: %v", got)
	}
	rendered["minio"].(map[string]any)["bucketName"] = "other"
	if got := withoutPrefix(configMismatches(exp, rendered), "etcd.rootPath"); len(got) != 1 || !strings.HasPrefix(got[0], "minio.bucketName") {
		t.Fatalf("a changed bucket must be reported: %v", got)
	}
	rendered["minio"].(map[string]any)["bucketName"] = "milvus-data"
	if got := withoutPrefix(configMismatches(cfg, rendered), "etcd.rootPath"); len(got) == 0 {
		t.Fatal("without the expectation the TLS switch is a difference")
	}
}

func withoutPrefix(in []string, prefix string) []string {
	var out []string
	for _, d := range in {
		if !strings.HasPrefix(d, prefix) {
			out = append(out, d)
		}
	}
	return out
}

func TestRenderObjectStorageTLS(t *testing.T) {
	refs := allRoles()
	units, tmpls, _, _ := assignUnits(refs, sourceOptions{Release: "rel", RoleLabel: "component"})
	src := &Source{Units: units, Templates: tmpls, Config: chartConfig(t)}
	in := renderInput{Namespace: "ns", Name: "rel-milvus", Version: "2.6.15-gpu", Release: "rel", Source: src, RootUser: "root", RootPassword: "x"}

	objs, err := render(in, newState(in))
	if err != nil {
		t.Fatal(err)
	}
	if _, found, _ := unstructured.NestedMap(objs.Milvus.Object, "spec", "objectStorage", "tls"); found {
		t.Fatal("no TLS flags: spec.objectStorage.tls must not be set")
	}

	in.ObjStore = &objectStorageTLS{Kind: "Secret", Name: "s3-ca", Key: "root.pem", Address: "rgw.example.svc", Port: 443}
	objs, err = render(in, newState(in))
	if err != nil {
		t.Fatal(err)
	}
	ref, _, _ := unstructured.NestedStringMap(objs.Milvus.Object, "spec", "objectStorage", "tls", "caSecret")
	if ref["name"] != "s3-ca" || ref["key"] != "root.pem" {
		t.Fatalf("spec.objectStorage.tls %v", objs.Milvus.Object["spec"])
	}
	var storage, config *core.Secret
	for _, s := range objs.Secrets {
		switch s.Name {
		case "rel-milvus-object-storage":
			storage = s
		case "rel-milvus-adopt-config":
			config = s
		}
	}
	if storage.StringData["useSSL"] != "true" || storage.StringData["address"] != "rgw.example.svc" || storage.StringData["port"] != "443" {
		t.Fatalf("object storage Secret %v", storage.StringData)
	}
	if strings.Contains(config.StringData["milvus.yaml"], "useSSL") || strings.Contains(config.StringData["milvus.yaml"], "minio.storage.svc") {
		t.Fatalf("the carried config must not hold the plain-HTTP storage settings:\n%s", config.StringData["milvus.yaml"])
	}

	in.ObjStore = &objectStorageTLS{Kind: "ConfigMap", Name: "ca", Key: "ca.crt"}
	objs, _ = render(in, newState(in))
	if m, _, _ := unstructured.NestedStringMap(objs.Milvus.Object, "spec", "objectStorage", "tls", "caConfigMap"); m["name"] != "ca" {
		t.Fatalf("caConfigMap %v", objs.Milvus.Object["spec"])
	}
}

func TestAdoptionFileObjectStorage(t *testing.T) {
	af, err := loadAdoptionFile(writeTemp(t, `apiVersion: cli.kubedb.com/v1alpha1
kind: MilvusAdoption
release: rel
objectStorage:
  caSecret: s3-ca/root.pem
  address: rgw.example.svc
  port: 443
`))
	if err != nil {
		t.Fatal(err)
	}
	var o runOptions
	var fo fileOptions
	cmd := newRunCommandForTest(&o, &fo)
	if err := af.applyFlags(cmd); err != nil {
		t.Fatal(err)
	}
	if o.ObjectStorageCASecret != "s3-ca/root.pem" || o.ObjectStorageAddress != "rgw.example.svc" || o.ObjectStoragePort != 443 {
		t.Fatalf("flags from the file: %+v", o.planOptions)
	}
}

func TestMapSRIOVAndNetworks(t *testing.T) {
	podWith := func(ann string, limits ...string) core.PodTemplateSpec {
		l := core.ResourceList{core.ResourceMemory: resource.MustParse("8Gi"), "nvidia.com/gpu": resource.MustParse("1")}
		for _, n := range limits {
			l[core.ResourceName(n)] = resource.MustParse("1")
		}
		t := core.PodTemplateSpec{Spec: core.PodSpec{Containers: []core.Container{{Name: "milvus", Resources: core.ResourceRequirements{Limits: l}}}}}
		if ann != "" {
			t.Annotations = map[string]string{annNetworks: ann}
		}
		return t
	}
	noNAD := func(string, string) string { return "" }

	if m, why := mapSRIOV("ns", podWith(""), noNAD); m != nil || why != "" {
		t.Fatalf("no network: %v %q", m, why)
	}

	// One network, one device: no NetworkAttachmentDefinition lookup needed.
	m, why := mapSRIOV("ns", podWith("sriov-net", "intel.com/sriov_a", "rdma/rdma_shared"), noNAD)
	s, _ := m["sriov"].(map[string]any)
	if why != "" || s["attachmentRef"] != "sriov-net" || s["resourceName"] != "intel.com/sriov_a" || s["rdmaResourceName"] != "rdma/rdma_shared" || s["gds"] != nil || s["interface"] != nil {
		t.Fatalf("%v %q", m, why)
	}

	// Two networks (RDMA and GDS) are told apart by the resource their definitions name.
	nad := func(ns, name string) string {
		return map[string]string{"rdma-net": "intel.com/sriov_rdma", "gds-net": "intel.com/sriov_gds"}[name]
	}
	m, why = mapSRIOV("ns", podWith(`[{"name":"rdma-net","interface":"net1"},{"name":"gds-net","interface":"storage0"}]`, "intel.com/sriov_gds", "intel.com/sriov_rdma"), nad)
	s, _ = m["sriov"].(map[string]any)
	g, _ := s["gds"].(map[string]any)
	if why != "" || s["resourceName"] != "intel.com/sriov_rdma" || g["attachmentRef"] != "gds-net" || g["resourceName"] != "intel.com/sriov_gds" || g["interface"] != "storage0" || s["interface"] != nil {
		t.Fatalf("%v %q", m, why)
	}

	// Never guess.
	for name, tc := range map[string]struct {
		ann string
		res []string
	}{
		"two devices, one network":    {"net", []string{"a/x", "a/y"}},
		"two networks, no definition": {"n1,n2", []string{"a/x", "a/y"}},
		"another namespace":           {"other/net", []string{"a/x"}},
		"three networks":              {"a,b,c", []string{"a/x", "a/y", "a/z"}},
		"network without a device":    {"net", nil},
		"bad annotation":              {"[oops", []string{"a/x"}},
	} {
		if m, why := mapSRIOV("ns", podWith(tc.ann, tc.res...), noNAD); m != nil || why == "" {
			t.Fatalf("%s: must not map, got %v %q", name, m, why)
		}
	}

	// The mapping reaches the node spec; no mapping adds nothing.
	node := nodeSpec(podWith("sriov-net", "intel.com/sriov_a"), m)
	if node["network"] == nil {
		t.Fatalf("network missing: %v", node)
	}
	if nodeSpec(podWith(""), nil)["network"] != nil {
		t.Fatal("no network, no spec.network")
	}
}

func TestParseNetworks(t *testing.T) {
	got, err := parseNetworks("ns/a@net3, b")
	if err != nil || len(got) != 2 || got[0] != (netAttachment{Name: "a", Namespace: "ns", Interface: "net3"}) || got[1].Name != "b" {
		t.Fatalf("%v %v", got, err)
	}
}

func TestSecretKeysDiffering(t *testing.T) {
	want := &core.Secret{StringData: map[string]string{"useSSL": "true", "address": "a", "port": "8443"}}
	cur := &core.Secret{Data: map[string][]byte{"useSSL": []byte("false"), "address": []byte("a")}}
	got := secretKeysDiffering(cur, want)
	if len(got) != 2 || got[0] != "port" || got[1] != "useSSL" {
		t.Fatalf("want the keys that differ or are missing, sorted: %v", got)
	}
	cur.Data["useSSL"], cur.Data["port"] = []byte("true"), []byte("8443")
	if got := secretKeysDiffering(cur, want); len(got) != 0 {
		t.Fatalf("equal Secrets: %v", got)
	}
}
