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
	"encoding/json"
	"fmt"
	"strings"

	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// operatorEnv are the variables the KubeDB operator sets on the Milvus container
// itself. A copy from the chart would override its value.
var operatorEnv = map[string]bool{
	"MINIO_PORT":             true,
	"MILVUS_USER":            true,
	"MILVUS_PASSWORD":        true,
	"MILVUS_ETCD_USER":       true,
	"MILVUS_ETCD_PASSWORD":   true,
	"MILVUS_CONFIG_SECTIONS": true,
}

// operatorOwned reports whether the operator sets the variable. With an SR-IOV
// network the init container writes MILVUS_CONF_<SECTION>_IP.
func operatorOwned(name string, sriov bool) bool {
	if operatorEnv[name] {
		return true
	}
	return sriov && strings.HasPrefix(name, "MILVUS_CONF_") && strings.HasSuffix(name, "_IP")
}

// carriedEnv splits the Milvus container's environment into what is copied to the
// Milvus object and the names left out because the operator sets them.
func carriedEnv(t core.PodTemplateSpec, sriov bool) (env []core.EnvVar, envFrom []core.EnvFromSource, skipped []string) {
	c := milvusContainer(t)
	if c == nil {
		return nil, nil, nil
	}
	for _, e := range c.Env {
		if operatorOwned(e.Name, sriov) {
			skipped = append(skipped, e.Name)
			continue
		}
		env = append(env, e)
	}
	return env, c.EnvFrom, skipped
}

// toMaps turns typed values into the maps the unstructured Milvus object holds.
func toMaps[T any](in []T) []any {
	out := make([]any, 0, len(in))
	for _, v := range in {
		b, _ := json.Marshal(v)
		var m map[string]any
		_ = json.Unmarshal(b, &m)
		out = append(out, m)
	}
	return out
}

// planEnv reports the variables of every unit that are carried over, those left
// out, and the Secrets and ConfigMaps they point to that are missing.
func (c *cluster) planEnv(ctx context.Context, src *Source, fs *findings) {
	cache := map[string]string{}
	check := func(kind, name, key string) string {
		id := kind + "/" + name + "/" + key
		if why, ok := cache[id]; ok {
			return why
		}
		why := ""
		switch kind {
		case "Secret":
			s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, name, metav1.GetOptions{})
			if kerr.IsNotFound(err) {
				why = fmt.Sprintf("Secret %s does not exist in namespace %s", name, c.ns)
			} else if err == nil && key != "" {
				if _, ok := s.Data[key]; !ok {
					why = fmt.Sprintf("Secret %s has no key %s", name, key)
				}
			}
		case "ConfigMap":
			cm, err := c.kube.CoreV1().ConfigMaps(c.ns).Get(ctx, name, metav1.GetOptions{})
			if kerr.IsNotFound(err) {
				why = fmt.Sprintf("ConfigMap %s does not exist in namespace %s", name, c.ns)
			} else if err == nil && key != "" {
				_, inData := cm.Data[key]
				_, inBin := cm.BinaryData[key]
				if !inData && !inBin {
					why = fmt.Sprintf("ConfigMap %s has no key %s", name, key)
				}
			}
		}
		cache[id] = why
		return why
	}
	envFindings(src, fs, check)
}

// envFindings does the reporting; check returns why a Secret or ConfigMap (and
// key, when given) is not usable, or "".
func envFindings(src *Source, fs *findings, check func(kind, name, key string) string) {
	skip := func(optional *bool) bool { return optional != nil && *optional }
	for _, u := range src.Units {
		t := src.Templates[u.key()]
		env, envFrom, skipped := carriedEnv(t, src.Networks[u.key()] != nil)
		id := u.Kind + " " + u.Name
		if len(env) > 0 {
			names := make([]string, len(env))
			for i, e := range env {
				names[i] = e.Name
			}
			fs.infof("%s: %d environment variables carried over: %s", id, len(env), strings.Join(names, ", "))
		}
		if len(envFrom) > 0 {
			fs.infof("%s: %d envFrom sources carried over", id, len(envFrom))
		}
		for _, n := range skipped {
			fs.warnf("%s sets %s, which the KubeDB operator sets itself; it is not carried over", id, n)
		}
		for _, e := range env {
			if e.ValueFrom == nil {
				continue
			}
			var why string
			if r := e.ValueFrom.SecretKeyRef; r != nil && !skip(r.Optional) {
				why = check("Secret", r.Name, r.Key)
			} else if r := e.ValueFrom.ConfigMapKeyRef; r != nil && !skip(r.Optional) {
				why = check("ConfigMap", r.Name, r.Key)
			}
			if why != "" {
				fs.errorf("%s: environment variable %s: %s, so the new pods would not start", id, e.Name, why)
			}
		}
		for _, f := range envFrom {
			var why string
			if r := f.SecretRef; r != nil && !skip(r.Optional) {
				why = check("Secret", r.Name, "")
			} else if r := f.ConfigMapRef; r != nil && !skip(r.Optional) {
				why = check("ConfigMap", r.Name, "")
			}
			if why != "" {
				fs.errorf("%s: envFrom: %s, so the new pods would not start", id, why)
			}
		}
	}
}
