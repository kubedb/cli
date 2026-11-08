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
	"os"
	"strings"

	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/yaml"
)

// planOptions are the inputs of an adoption that discovery cannot find.
type planOptions struct {
	sourceOptions
	Name                     string
	Version                  string
	AuthSecret               string
	RootPasswordFile         string
	EtcdAuthSecret           string
	TLSIssuer                string
	TLSCASecret              string
	TLSClientSecret          string
	ObjectStorageCASecret    string
	ObjectStorageCAConfigMap string
	ObjectStorageAddress     string
	ObjectStoragePort        int
	ConfigFile               string
	RestrictedSCC            bool
	SkipCredentialCheck      bool
	// allowExisting accepts a Milvus of that name: run resumes an adoption
	// whose Milvus it already created.
	allowExisting bool
}

// adoptPlan is a checked adoption: the findings and, when nothing blocks it,
// the objects to create.
type adoptPlan struct {
	Findings findings
	Source   *Source
	State    *State
	Objects  *adoptObjects
	TLS      *adoptionTLS
	ObjStore *objectStorageTLS
	// GoverningServices are the chart's Services on names KubeDB needs headless;
	// run deletes them before it creates the Milvus.
	GoverningServices []core.Service
}

// governingRoles have a headless governing Service named <name>-<role>.
var governingRoles = []string{roleMixCoord, roleDataNode, roleQueryNode, roleStreamingNode}

func buildPlan(ctx context.Context, c *cluster, o planOptions) (*adoptPlan, error) {
	src, err := discover(ctx, c.kube, c.ns, o.sourceOptions)
	if err != nil {
		return nil, err
	}
	p := &adoptPlan{Source: src, Findings: src.Findings}
	fs := &p.Findings
	if o.ConfigFile != "" {
		raw, err := os.ReadFile(o.ConfigFile)
		if err != nil {
			return nil, err
		}
		cfg, _, err := mergeConfigFiles(map[string]string{"user.yaml": string(raw)})
		if err != nil {
			return nil, err
		}
		src.Config = cfg
		src.ConfigFrom = []string{o.ConfigFile}
		p.Findings = withoutConfigError(p.Findings)
	}
	if fs.hasErrors() {
		return p, nil
	}
	for _, u := range src.Units {
		fs.infof("%-13s %s %s, %d replicas%s", u.Role, u.Kind, u.Name, u.Replicas, groupNote(u))
	}
	fs.infof("config read from %s", strings.Join(src.ConfigFrom, "; "))
	cfg := src.Config

	// The Milvus is named after the client Service, so KubeDB takes it over.
	name := o.Name
	switch {
	case name == "" && src.Service != nil:
		name = src.Service.Name
	case name == "":
		fs.errorf("no Service selects the proxy pods; name the Milvus with --name")
		return p, nil
	case src.Service != nil && name != src.Service.Name:
		fs.warnf("clients use Service %s; the Milvus %s gets a new Service, and clients must move to it", src.Service.Name, name)
	}
	if src.Service != nil && src.Service.Name == name {
		fs.infof("client Service %s (%s) is taken over: same name, IP and ports", name, src.Service.Spec.ClusterIP)
	}
	if _, err := c.dyn.Resource(MilvusGVR).Namespace(c.ns).Get(ctx, name, metav1.GetOptions{}); err == nil {
		if !o.allowExisting {
			fs.errorf("a Milvus %s already exists", name)
		}
	} else if !kerr.IsNotFound(err) {
		return nil, err
	}

	// Write-ahead log: KubeDB runs Woodpecker on object storage only.
	switch wt := walType(cfg); wt {
	case "woodpecker":
		fs.infof("write-ahead log: woodpecker")
	case "":
		if c.streamingNodeLogsMention(ctx, src, "woodpecker") {
			fs.infof("write-ahead log: woodpecker (from the StreamingNode logs; the config does not set mq.type)")
		} else {
			fs.errorf("cannot confirm the write-ahead log is Woodpecker: the config sets no mq.type and the StreamingNode logs do not mention it")
		}
	default:
		fs.errorf("write-ahead log is %s; KubeDB runs Woodpecker only, and changing the WAL of a running cluster is a migration", wt)
	}
	if t := lookupString(cfg, "woodpecker.storage.type"); t != "" && t != "minio" {
		fs.errorf("woodpecker.storage.type is %s; KubeDB keeps the write-ahead log on the object storage (minio)", t)
	}

	// Metadata.
	etcd := etcdFromConfig(cfg)
	if len(etcd.Endpoints) == 0 {
		fs.errorf("the config sets no etcd.endpoints")
	} else {
		fs.infof("etcd %s, root path %q (kept)", strings.Join(etcd.Endpoints, ","), etcd.RootPath)
	}
	if etcd.AuthEnabled {
		if o.EtcdAuthSecret == "" {
			fs.errorf("the etcd uses authentication; give its root user in a basic-auth Secret with --etcd-auth-secret")
		} else if err := c.basicAuthSecretExists(ctx, o.EtcdAuthSecret); err != nil {
			fs.errorf("--etcd-auth-secret: %v", err)
		}
	}

	// Object storage.
	osData := objectStorageData(cfg)
	if osData["bucketName"] == "" || osData["address"] == "" {
		fs.errorf("the config does not set minio.address and minio.bucketName")
	} else {
		fs.infof("object storage %s, bucket %s, root path %q (kept)", osData["address"], osData["bucketName"], osData["rootPath"])
	}
	if osData["accessKeyID"] == "" && !strings.EqualFold(osData["useIAM"], "true") {
		fs.warnf("the config has no minio.accessKeyID and minio.useIAM is not true; if the credentials come from environment variables, add them to Secret %s", objectStorageSecretName(name))
	}

	// Version.
	version := o.Version
	if version == "" {
		list, err := c.dyn.Resource(MilvusVersionGVR).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list MilvusVersions: %w", err)
		}
		if version, err = pickVersion(src.Image, list.Items); err != nil {
			fs.errorf("%v", err)
		}
	}
	if version != "" {
		fs.infof("image %s -> MilvusVersion %s", src.Image, version)
	}

	// Credentials.
	in := renderInput{
		Namespace: c.ns, Name: name, Version: version, Release: o.Release, Source: src,
		EtcdAuthSecret: o.EtcdAuthSecret, RestrictedSCC: o.RestrictedSCC,
	}
	if err := c.planTLS(ctx, o, src, &in, fs); err != nil {
		return nil, err
	}
	if err := c.planObjectStorageTLS(ctx, o, cfg, &in, fs); err != nil {
		return nil, err
	}
	c.planNetworks(ctx, src, fs)
	if !lookupBool(cfg, "common.security.authorizationEnabled") {
		in.DisableAuth = true
		fs.infof("authentication is off; the Milvus gets spec.disableSecurity: true")
	} else if err := c.resolveRootCredentials(ctx, o, src, &in, fs); err != nil {
		return nil, err
	}

	// Services KubeDB needs headless.
	for _, role := range governingRoles {
		svc, err := c.kube.CoreV1().Services(c.ns).Get(ctx, name+"-"+role, metav1.GetOptions{})
		if kerr.IsNotFound(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if svc.Spec.ClusterIP == core.ClusterIPNone {
			continue
		}
		if belongsToRelease(*svc, src) {
			p.GoverningServices = append(p.GoverningServices, *svc)
			fs.infof("Service %s is not headless; run deletes it (KubeDB recreates it headless; Milvus does not use it)", svc.Name)
		} else {
			fs.errorf("Service %s exists, is not headless and does not belong to the release; KubeDB needs that name for a headless Service", svc.Name)
		}
	}

	if fs.hasErrors() {
		return p, nil
	}
	p.TLS = in.TLS
	p.ObjStore = in.ObjStore
	p.State = newState(in)
	if p.Objects, err = render(in, p.State); err != nil {
		return nil, err
	}
	return p, nil
}

func groupNote(u Unit) string {
	if u.Group == "" {
		return ""
	}
	return ", node group " + u.Group
}

func withoutConfigError(fs findings) findings {
	var out findings
	for _, f := range fs {
		if f.Level == levelError && strings.Contains(f.Msg, "--config-file") {
			continue
		}
		out = append(out, f)
	}
	return out
}

// belongsToRelease reports whether a Service selects pods of the release, or
// carries the release's instance label.
func belongsToRelease(svc core.Service, src *Source) bool {
	for _, u := range src.Units {
		if len(svc.Spec.Selector) > 0 && labels.SelectorFromSet(svc.Spec.Selector).Matches(labels.Set(src.Templates[u.key()].Labels)) {
			return true
		}
	}
	for _, u := range src.Units {
		if v := u.Selector[labelInstance]; v != "" && svc.Labels[labelInstance] == v {
			return true
		}
	}
	return false
}

func (c *cluster) basicAuthSecretExists(ctx context.Context, name string) error {
	s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if len(s.Data[core.BasicAuthUsernameKey]) == 0 || len(s.Data[core.BasicAuthPasswordKey]) == 0 {
		return fmt.Errorf("the Secret %s needs the keys %s and %s", name, core.BasicAuthUsernameKey, core.BasicAuthPasswordKey)
	}
	return nil
}

// resolveRootCredentials picks the root credentials (an existing Secret, a
// password file, or the chart's defaultRootPassword) and logs in with them.
// Milvus applies a root password only on its first start, so a password that
// does not log in today must not be handed to KubeDB.
func (c *cluster) resolveRootCredentials(ctx context.Context, o planOptions, src *Source, in *renderInput, fs *findings) error {
	user, password, from := "root", "", ""
	switch {
	case o.AuthSecret != "":
		s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, o.AuthSecret, metav1.GetOptions{})
		if err != nil {
			fs.errorf("--auth-secret: %v", err)
			return nil
		}
		user, password = string(s.Data[core.BasicAuthUsernameKey]), string(s.Data[core.BasicAuthPasswordKey])
		if user == "" || password == "" {
			fs.errorf("Secret %s needs the keys %s and %s", o.AuthSecret, core.BasicAuthUsernameKey, core.BasicAuthPasswordKey)
			return nil
		}
		in.AuthSecret, from = o.AuthSecret, "Secret "+o.AuthSecret
	case o.RootPasswordFile != "":
		raw, err := os.ReadFile(o.RootPasswordFile)
		if err != nil {
			return err
		}
		password, from = strings.TrimRight(string(raw), "\r\n"), o.RootPasswordFile
	default:
		password = lookupString(src.Config, "common.security.defaultRootPassword")
		if password == "" {
			fs.errorf("authentication is on and the config has no defaultRootPassword; give the current root password with --auth-secret or --root-password-file")
			return nil
		}
		from = "common.security.defaultRootPassword"
	}
	in.RootUser, in.RootPassword = user, password
	if o.SkipCredentialCheck {
		fs.warnf("root credentials from %s are not checked (--skip-credential-check)", from)
		return nil
	}
	res, detail := c.checkCredentials(ctx, src.firstUnit(roleProxy), restPort(src.Config), user, password)
	switch res {
	case credentialsAccepted:
		fs.infof("root credentials from %s log in to %s", from, detail)
	case credentialsRejected:
		fs.errorf("root credentials from %s are rejected by the running cluster (%s); Milvus applies a root password only on its first start, so give the current one with --auth-secret or --root-password-file", from, detail)
	default:
		fs.warnf("root credentials from %s could not be checked: %s", from, detail)
	}
	return nil
}

// streamingNodeLogsMention reports whether a StreamingNode's recent log mentions s.
func (c *cluster) streamingNodeLogsMention(ctx context.Context, src *Source, s string) bool {
	u := src.firstUnit(roleStreamingNode)
	pods, err := c.kube.CoreV1().Pods(c.ns).List(ctx, metav1.ListOptions{LabelSelector: labels.SelectorFromSet(u.Selector).String()})
	if err != nil || len(pods.Items) == 0 {
		return false
	}
	limit := int64(4 << 20)
	raw, err := c.kube.CoreV1().Pods(c.ns).GetLogs(pods.Items[0].Name, &core.PodLogOptions{LimitBytes: &limit}).DoRaw(ctx)
	return err == nil && strings.Contains(strings.ToLower(string(raw)), s)
}

// yamlDocs renders the plan's objects as YAML, one document set per file.
func (p *adoptPlan) yamlDocs(redact bool) (map[string]string, error) {
	out := map[string]string{}
	var secrets []string
	for _, s := range p.Objects.Secrets {
		cp := s.DeepCopy()
		if redact {
			for k := range cp.StringData {
				cp.StringData[k] = "<redacted: use --output-dir>"
			}
		}
		b, err := yaml.Marshal(cp)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, string(b))
	}
	out["00-secrets.yaml"] = strings.Join(secrets, "---\n")
	if p.Objects.Issuer != nil {
		b, err := yaml.Marshal(p.Objects.Issuer.Object)
		if err != nil {
			return nil, err
		}
		out["00-issuer.yaml"] = string(b)
	}
	b, err := yaml.Marshal(p.Objects.State)
	if err != nil {
		return nil, err
	}
	out["01-adopt-state.yaml"] = string(b)
	b, err = yaml.Marshal(p.Objects.Milvus.Object)
	if err != nil {
		return nil, err
	}
	out["02-milvus.yaml"] = string(b)
	return out, nil
}
