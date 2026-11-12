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
	"bytes"
	"strings"
	"testing"

	core "k8s.io/api/core/v1"
)

func envTemplate(env []core.EnvVar, from []core.EnvFromSource) core.PodTemplateSpec {
	return core.PodTemplateSpec{Spec: core.PodSpec{Containers: []core.Container{{Name: "milvus", Image: "milvusdb/milvus:v2.6.15", Env: env, EnvFrom: from}}}}
}

func TestCarriedEnv(t *testing.T) {
	tpl := envTemplate([]core.EnvVar{
		{Name: "TZ", Value: "UTC"},
		{Name: "MILVUS_USER", Value: "x"},
		{Name: "MILVUS_CONF_PROXY_IP", Value: "1.2.3.4"},
		{Name: "TOKEN", ValueFrom: &core.EnvVarSource{SecretKeyRef: &core.SecretKeySelector{LocalObjectReference: core.LocalObjectReference{Name: "s"}, Key: "k"}}},
	}, []core.EnvFromSource{{ConfigMapRef: &core.ConfigMapEnvSource{LocalObjectReference: core.LocalObjectReference{Name: "cm"}}}})

	env, from, skipped := carriedEnv(tpl, false)
	if len(env) != 3 || env[0].Name != "TZ" || env[1].Name != "MILVUS_CONF_PROXY_IP" || env[2].Name != "TOKEN" || len(from) != 1 {
		t.Fatalf("no SR-IOV: %v %v", env, from)
	}
	if len(skipped) != 1 || skipped[0] != "MILVUS_USER" {
		t.Fatalf("skipped %v", skipped)
	}
	// With an SR-IOV network the init container writes MILVUS_CONF_*_IP.
	env, _, skipped = carriedEnv(tpl, true)
	if len(env) != 2 || len(skipped) != 2 || skipped[1] != "MILVUS_CONF_PROXY_IP" {
		t.Fatalf("SR-IOV: %v %v", env, skipped)
	}
	if e, f, s := carriedEnv(core.PodTemplateSpec{}, false); e != nil || f != nil || s != nil {
		t.Fatal("no container")
	}
}

func TestNodeSpecEnv(t *testing.T) {
	tpl := envTemplate([]core.EnvVar{{Name: "TZ", Value: "UTC"}}, nil)
	node := nodeSpec(tpl, nil)
	cs := node["podTemplate"].(map[string]any)["spec"].(map[string]any)["containers"].([]any)
	main := cs[0].(map[string]any)
	if main["name"] != "milvus" {
		t.Fatalf("%v", main)
	}
	env := main["env"].([]any)[0].(map[string]any)
	if env["name"] != "TZ" || env["value"] != "UTC" {
		t.Fatalf("%v", env)
	}
	if _, ok := main["envFrom"]; ok {
		t.Fatal("no envFrom expected")
	}
	// No env, no resources: no container entry.
	if n := nodeSpec(envTemplate(nil, nil), nil); n["podTemplate"] != nil {
		t.Fatalf("%v", n)
	}
}

func TestEnvFindings(t *testing.T) {
	check := func(kind, name, key string) string {
		switch {
		case kind == "Secret" && name == "s" && key == "k", name == "s" && key == "":
			return ""
		case kind == "Secret" && name == "s":
			return "Secret s has no key " + key
		}
		return kind + " " + name + " does not exist in namespace ns"
	}
	secretRef := func(name, key string) core.EnvVar {
		return core.EnvVar{Name: "A_" + name + key, ValueFrom: &core.EnvVarSource{SecretKeyRef: &core.SecretKeySelector{LocalObjectReference: core.LocalObjectReference{Name: name}, Key: key}}}
	}
	yes := true
	optional := secretRef("gone2", "k")
	optional.ValueFrom.SecretKeyRef.Optional = &yes
	u := Unit{Kind: "Deployment", Name: "qn", Role: roleQueryNode}
	src := &Source{Units: []Unit{u}, Templates: map[string]core.PodTemplateSpec{}, Networks: map[string]map[string]any{}}
	src.Templates[u.key()] = envTemplate([]core.EnvVar{
		{Name: "TZ", Value: "UTC"}, {Name: "MILVUS_PASSWORD", Value: "x"},
		secretRef("s", "k"), secretRef("s", "missing"), secretRef("gone", "k"), optional,
	}, []core.EnvFromSource{{ConfigMapRef: &core.ConfigMapEnvSource{LocalObjectReference: core.LocalObjectReference{Name: "nocm"}}}})
	fs := &findings{}
	envFindings(src, fs, check)
	var buf bytes.Buffer
	fs.print(&buf)
	out := buf.String()
	for _, want := range []string{
		"5 environment variables", "MILVUS_PASSWORD, which the KubeDB operator sets itself",
		"Secret s has no key missing", "Secret gone does not exist", "ConfigMap nocm does not exist",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "gone2") && strings.Contains(out, "ERROR") && strings.Count(out, "ERROR") != 3 {
		t.Errorf("an optional reference must not be an error:\n%s", out)
	}
	if strings.Count(out, "ERROR") != 3 {
		t.Errorf("want 3 errors:\n%s", out)
	}
}
