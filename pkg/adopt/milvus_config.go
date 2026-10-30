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
	"fmt"
	"sort"
	"strings"

	"sigs.k8s.io/yaml"
)

// Milvus reads its configuration files from /milvus/configs in this order, each
// overriding the one before; a chart mounts some of them from a ConfigMap.
var milvusConfigFiles = []string{"milvus.yaml", "default.yaml", "user.yaml"}

const defaultEtcdRootPath = "by-dev"

// The minio.* keys the KubeDB operator reads from spec.objectStorage.configSecret
// and writes back into the rendered milvus.yaml.
var objectStorageKeys = []string{
	"address", "port", "useSSL", "bucketName", "rootPath", "useIAM", "cloudProvider",
	"region", "iamEndpoint", "useVirtualHost", "accessKeyID", "secretAccessKey",
}

// mergeConfigFiles merges the files of one Milvus in load order. Unknown file
// names are returned so the caller can report them.
func mergeConfigFiles(files map[string]string) (map[string]any, []string, error) {
	out := map[string]any{}
	for _, name := range milvusConfigFiles {
		raw, ok := files[name]
		if !ok {
			continue
		}
		var doc map[string]any
		if err := yaml.Unmarshal([]byte(raw), &doc); err != nil {
			return nil, nil, fmt.Errorf("failed to parse %s: %w", name, err)
		}
		out = deepMerge(out, doc)
	}
	var unknown []string
	for name := range files {
		if !contains(milvusConfigFiles, name) {
			unknown = append(unknown, name)
		}
	}
	sort.Strings(unknown)
	return out, unknown, nil
}

func deepMerge(dst, src map[string]any) map[string]any {
	for k, v := range src {
		if sm, ok := v.(map[string]any); ok {
			if dm, ok := dst[k].(map[string]any); ok {
				dst[k] = deepMerge(dm, sm)
				continue
			}
		}
		dst[k] = v
	}
	return dst
}

// lookup returns the value at a dotted path, e.g. "minio.bucketName".
func lookup(cfg map[string]any, path string) (any, bool) {
	var cur any = cfg
	for _, p := range strings.Split(path, ".") {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		if cur, ok = m[p]; !ok {
			return nil, false
		}
	}
	return cur, true
}

func lookupString(cfg map[string]any, path string) string {
	v, ok := lookup(cfg, path)
	if !ok || v == nil {
		return ""
	}
	return scalarString(v)
}

func lookupBool(cfg map[string]any, path string) bool {
	return strings.EqualFold(lookupString(cfg, path), "true")
}

// scalarString renders a YAML scalar the way it would be written back: numbers
// without exponent, booleans as true/false.
func scalarString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%v", t)
	default:
		return fmt.Sprintf("%v", t)
	}
}

// etcdSettings is what the running Milvus uses to reach its metadata.
type etcdSettings struct {
	Endpoints   []string
	RootPath    string
	AuthEnabled bool
	TLSEnabled  bool
}

func etcdFromConfig(cfg map[string]any) etcdSettings {
	s := etcdSettings{RootPath: lookupString(cfg, "etcd.rootPath")}
	if s.RootPath == "" {
		s.RootPath = defaultEtcdRootPath
	}
	s.AuthEnabled = lookupBool(cfg, "etcd.auth.enabled")
	s.TLSEnabled = lookupBool(cfg, "etcd.ssl.enabled")
	scheme := "http://"
	if s.TLSEnabled {
		scheme = "https://"
	}
	var eps []string
	switch v, _ := lookup(cfg, "etcd.endpoints"); t := v.(type) {
	case []any:
		for _, e := range t {
			eps = append(eps, scalarString(e))
		}
	case string:
		eps = strings.Split(t, ",")
	}
	for _, e := range eps {
		if e = strings.TrimSpace(e); e == "" {
			continue
		}
		if !strings.Contains(e, "://") {
			e = scheme + e
		}
		s.Endpoints = append(s.Endpoints, e)
	}
	return s
}

// walType returns the message queue the config selects ("" when it does not say).
func walType(cfg map[string]any) string {
	if t := lookupString(cfg, "mq.type"); t != "" && t != "default" {
		return t
	}
	return lookupString(cfg, "messageQueue")
}

// objectStorageData returns the Secret data for spec.objectStorage.configSecret:
// every minio.* key the running config sets.
func objectStorageData(cfg map[string]any) map[string]string {
	out := map[string]string{}
	for _, k := range objectStorageKeys {
		if v, ok := lookup(cfg, "minio."+k); ok && v != nil && scalarString(v) != "" {
			out[k] = scalarString(v)
		}
	}
	return out
}

// carriedConfig is the milvus.yaml handed to KubeDB through
// spec.configuration.secretName. The etcd section is dropped (the Milvus object
// sets the endpoints and root path); everything else is kept as it runs today.
// Milvus serves its RESTful API on the gRPC port unless proxy.http.port is set,
// while KubeDB defaults it to 8080, so an unset port is pinned to "".
func carriedConfig(cfg map[string]any) (string, error) {
	out := deepCopy(cfg)
	delete(out, "etcd")
	if v, ok := lookup(out, "proxy.http.port"); !ok || v == nil {
		proxy, _ := out["proxy"].(map[string]any)
		if proxy == nil {
			proxy = map[string]any{}
			out["proxy"] = proxy
		}
		http, _ := proxy["http"].(map[string]any)
		if http == nil {
			http = map[string]any{}
			proxy["http"] = http
		}
		http["port"] = ""
	}
	b, err := yaml.Marshal(out)
	return string(b), err
}

func deepCopy(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		if m, ok := v.(map[string]any); ok {
			out[k] = deepCopy(m)
			continue
		}
		out[k] = v
	}
	return out
}

// criticalConfigPaths are the settings that decide where the data is read from
// and how clients log in. After the Milvus is created, the config the operator
// renders must agree with the running one on each of them.
var criticalConfigPaths = []string{
	"minio.address", "minio.port", "minio.bucketName", "minio.rootPath", "minio.useSSL",
	"minio.useIAM", "minio.cloudProvider", "minio.region", "minio.accessKeyID", "minio.secretAccessKey",
	"mq.type", "woodpecker.storage.type", "woodpecker.storage.rootPath",
	"common.security.authorizationEnabled", "proxy.port", "proxy.http.port",
}

// configMismatches compares the running config with the one the operator
// rendered. A path the running config does not set is compared only for the
// etcd root path, whose Milvus and KubeDB defaults differ.
func configMismatches(running, rendered map[string]any) []string {
	var out []string
	want := etcdFromConfig(running).RootPath
	if got := lookupString(rendered, "etcd.rootPath"); got != want {
		out = append(out, fmt.Sprintf("etcd.rootPath: running %q, rendered %q", want, got))
	}
	for _, p := range criticalConfigPaths {
		rv, ok := lookup(running, p)
		if !ok || rv == nil {
			if p == "proxy.http.port" {
				rv = "" // carriedConfig pins it
			} else {
				continue
			}
		}
		a, b := scalarString(rv), lookupString(rendered, p)
		if a != b {
			if strings.Contains(strings.ToLower(p), "secret") || strings.HasSuffix(p, "accessKeyID") {
				out = append(out, p+": differs")
				continue
			}
			out = append(out, fmt.Sprintf("%s: running %q, rendered %q", p, a, b))
		}
	}
	return out
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
