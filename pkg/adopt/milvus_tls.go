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
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Secret keys the KubeDB operator reads a certificate from.
const (
	tlsCertKey = "tls.crt"
	tlsKeyKey  = "tls.key"
	tlsCAKey   = "ca.crt"

	certManagerGroup = "cert-manager.io"
	annIssuerName    = "cert-manager.io/issuer-name"
	annIssuerKind    = "cert-manager.io/issuer-kind"
	milvusWorkDir    = "/milvus"
	certExpiryWarn   = 30 * 24 * time.Hour
)

var (
	certIssuerGVR        = schema.GroupVersionResource{Group: certManagerGroup, Version: "v1", Resource: "issuers"}
	certClusterIssuerGVR = schema.GroupVersionResource{Group: certManagerGroup, Version: "v1", Resource: "clusterissuers"}
)

// tlsMaterial is a certificate as the running Milvus uses it. Secret is set when
// a Secret already has the keys KubeDB reads, so it can be referenced as it is
// and its owner's renewal keeps working.
type tlsMaterial struct {
	Cert, Key, CA []byte
	Secret        string
	// Source is the Secret the certificate file comes from.
	Source string
}

func (m *tlsMaterial) secretData() map[string]string {
	return map[string]string{tlsCertKey: string(m.Cert), tlsKeyKey: string(m.Key), tlsCAKey: string(m.CA)}
}

// adoptionTLS is the TLS of the adopted Milvus. With an issuer, KubeDB issues and
// renews every certificate from the CA clients already trust. Without one it
// uses the running server certificate as it is.
type adoptionTLS struct {
	ExternalMode string // Disabled, TLS or mTLS
	Internal     bool
	Server       *tlsMaterial // the running certificate
	ClientSecret string
	IssuerKind   string
	IssuerName   string
	// CASecret is set when the Issuer does not exist yet: run creates it from this Secret.
	CASecret string
	// ExtraDNSNames and ExtraIPs are the names of the running certificate that the
	// operator does not add by itself.
	ExtraDNSNames []string
	ExtraIPs      []string
}

// managed reports whether KubeDB issues the certificates.
func (t *adoptionTLS) managed() bool { return t.IssuerName != "" }

// parseIssuerRef parses Kind/name; a bare name is an Issuer.
func parseIssuerRef(s string) (kind, name string, err error) {
	kind, name = "Issuer", s
	if k, n, ok := strings.Cut(s, "/"); ok {
		kind, name = k, n
	}
	if name == "" || (kind != "Issuer" && kind != "ClusterIssuer") {
		return "", "", fmt.Errorf("invalid issuer %q: want name, Issuer/name or ClusterIssuer/name", s)
	}
	return kind, name, nil
}

// externalTLSMode reads common.security.tlsMode: 0 off, 1 TLS, 2 mTLS.
func externalTLSMode(cfg map[string]any) (string, error) {
	switch v := lookupString(cfg, "common.security.tlsMode"); v {
	case "", "0":
		return "Disabled", nil
	case "1":
		return "TLS", nil
	case "2":
		return "mTLS", nil
	default:
		return "", fmt.Errorf("common.security.tlsMode is %q; Milvus knows 0, 1 and 2", v)
	}
}

// fileSource is one place a mounted file may come from.
type fileSource struct {
	Kind string // Secret or ConfigMap
	Name string
	Key  string
}

// mountedFile lists where a file in the Milvus container comes from. Several
// sources are possible only for a projected volume without items.
func mountedFile(t core.PodTemplateSpec, file string) ([]fileSource, error) {
	if !path.IsAbs(file) {
		file = path.Join(milvusWorkDir, file)
	}
	c := milvusContainer(t)
	if c == nil {
		return nil, errors.New("no Milvus container")
	}
	var mount *core.VolumeMount
	for i := range c.VolumeMounts {
		m := &c.VolumeMounts[i]
		mp := strings.TrimSuffix(m.MountPath, "/")
		if file != mp && !strings.HasPrefix(file, mp+"/") {
			continue
		}
		if mount == nil || len(mp) > len(strings.TrimSuffix(mount.MountPath, "/")) {
			mount = m
		}
	}
	if mount == nil {
		return nil, fmt.Errorf("%s is not on a mounted volume", file)
	}
	rel := strings.TrimPrefix(strings.TrimPrefix(file, strings.TrimSuffix(mount.MountPath, "/")), "/")
	rel = path.Join(mount.SubPath, rel)
	if rel == "" || rel == "." {
		return nil, fmt.Errorf("%s is a directory", file)
	}
	for _, v := range t.Spec.Volumes {
		if v.Name == mount.Name {
			return volumeFile(v, rel, file)
		}
	}
	return nil, fmt.Errorf("volume %s of %s does not exist", mount.Name, file)
}

func itemKey(items []core.KeyToPath, rel string) (string, bool) {
	if len(items) == 0 {
		return rel, true
	}
	for _, it := range items {
		if it.Path == rel {
			return it.Key, true
		}
	}
	return "", false
}

func volumeFile(v core.Volume, rel, file string) ([]fileSource, error) {
	var out []fileSource
	add := func(kind, name string, items []core.KeyToPath) {
		if key, ok := itemKey(items, rel); ok {
			out = append(out, fileSource{Kind: kind, Name: name, Key: key})
		}
	}
	switch {
	case v.Secret != nil:
		add("Secret", v.Secret.SecretName, v.Secret.Items)
	case v.ConfigMap != nil:
		add("ConfigMap", v.ConfigMap.Name, v.ConfigMap.Items)
	case v.Projected != nil:
		for _, s := range v.Projected.Sources {
			if s.Secret != nil {
				add("Secret", s.Secret.Name, s.Secret.Items)
			}
			if s.ConfigMap != nil {
				add("ConfigMap", s.ConfigMap.Name, s.ConfigMap.Items)
			}
		}
	default:
		return nil, fmt.Errorf("%s does not come from a Secret or ConfigMap (volume %s), so it cannot be carried over", file, v.Name)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("volume %s does not provide %s", v.Name, file)
	}
	return out, nil
}

func (c *cluster) readSource(ctx context.Context, src fileSource) ([]byte, bool, error) {
	if src.Kind == "ConfigMap" {
		cm, err := c.kube.CoreV1().ConfigMaps(c.ns).Get(ctx, src.Name, metav1.GetOptions{})
		if err != nil {
			return nil, false, err
		}
		v, ok := cm.Data[src.Key]
		return []byte(v), ok, nil
	}
	s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, src.Name, metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	v, ok := s.Data[src.Key]
	return v, ok, nil
}

func (c *cluster) readMounted(ctx context.Context, t core.PodTemplateSpec, file string) ([]byte, fileSource, error) {
	sources, err := mountedFile(t, file)
	if err != nil {
		return nil, fileSource{}, err
	}
	for _, src := range sources {
		data, ok, err := c.readSource(ctx, src)
		if err != nil {
			return nil, src, fmt.Errorf("failed to read %s %s: %w", src.Kind, src.Name, err)
		}
		if ok {
			return data, src, nil
		}
	}
	return nil, fileSource{}, fmt.Errorf("no mounted Secret or ConfigMap has %s", file)
}

// readTLS reads the certificate, key and CA at the three config paths.
func (c *cluster) readTLS(ctx context.Context, t core.PodTemplateSpec, cfg map[string]any, certPath, keyPath, caPath string) (*tlsMaterial, error) {
	m := &tlsMaterial{}
	var srcs [3]fileSource
	for i, p := range []string{certPath, keyPath, caPath} {
		file := lookupString(cfg, p)
		if file == "" {
			return nil, fmt.Errorf("the config does not set %s", p)
		}
		data, src, err := c.readMounted(ctx, t, file)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		srcs[i] = src
		switch i {
		case 0:
			m.Cert = data
		case 1:
			m.Key = data
		case 2:
			m.CA = data
		}
	}
	name := srcs[0].Name
	if srcs[0].Kind == "Secret" {
		m.Source = name
	}
	if srcs == [3]fileSource{{"Secret", name, tlsCertKey}, {"Secret", name, tlsKeyKey}, {"Secret", name, tlsCAKey}} {
		m.Secret = name
	}
	return m, nil
}

func parseLeaf(certPEM []byte) (*x509.Certificate, *x509.CertPool, error) {
	var leaf *x509.Certificate
	intermediates := x509.NewCertPool()
	for rest := certPEM; ; {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, nil, err
		}
		if leaf == nil {
			leaf = cert
		} else {
			intermediates.AddCert(cert)
		}
	}
	if leaf == nil {
		return nil, nil, errors.New("no certificate in the PEM data")
	}
	return leaf, intermediates, nil
}

func verifyAgainstCA(certPEM, caPEM []byte, now time.Time) error {
	leaf, intermediates, err := parseLeaf(certPEM)
	if err != nil {
		return err
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caPEM) {
		return errors.New("no certificate in the CA data")
	}
	_, err = leaf.Verify(x509.VerifyOptions{
		Roots: roots, Intermediates: intermediates, CurrentTime: now,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	})
	return err
}

// checkCertificate reports a certificate that cannot serve this cluster.
func checkCertificate(label string, m *tlsMaterial, hosts []string, chain bool, now time.Time, fs *findings) {
	if _, err := tls.X509KeyPair(m.Cert, m.Key); err != nil {
		fs.errorf("%s: the certificate and key do not form a pair: %v", label, err)
		return
	}
	leaf, _, err := parseLeaf(m.Cert)
	if err != nil {
		fs.errorf("%s: %v", label, err)
		return
	}
	if chain {
		if err := verifyAgainstCA(m.Cert, m.CA, now); err != nil {
			fs.errorf("%s: the certificate does not verify against its CA: %v", label, err)
		}
	}
	for _, h := range hosts {
		if err := leaf.VerifyHostname(h); err != nil {
			fs.errorf("%s: the certificate does not cover %s, which KubeDB connects to (it covers %s); reissue it with that name", label, h, strings.Join(leaf.DNSNames, ", "))
		}
	}
	if left := leaf.NotAfter.Sub(now); left > 0 && left < certExpiryWarn {
		fs.warnf("%s: the certificate expires on %s", label, leaf.NotAfter.Format("2006-01-02"))
	}
}

// clientHosts are the names KubeDB and Milvus use to reach the proxies.
func clientHosts(name, ns string) []string {
	return []string{name, name + "." + ns + ".svc", name + "." + ns + ".svc.cluster.local"}
}

// planEtcdTLS carries the client certificate of an etcd that uses TLS.
func (c *cluster) planEtcdTLS(ctx context.Context, t core.PodTemplateSpec, cfg map[string]any, in *renderInput, fs *findings) {
	if !etcdFromConfig(cfg).TLSEnabled {
		return
	}
	m, err := c.readTLS(ctx, t, cfg, "etcd.ssl.tlsCert", "etcd.ssl.tlsKey", "etcd.ssl.tlsCACert")
	if err != nil {
		fs.errorf("etcd TLS: %v", err)
		return
	}
	checkCertificate("etcd TLS", m, nil, false, time.Now(), fs)
	in.EtcdTLS = m
	fs.infof("etcd TLS: %s", materialNote(m))
}

func materialNote(m *tlsMaterial) string {
	if m.Secret != "" {
		return "Secret " + m.Secret + " is used as it is"
	}
	return "the certificate is copied to a new Secret, because its keys are not tls.crt, tls.key and ca.crt; renewing the copy is up to you"
}

// planTLS keeps the TLS the running Milvus uses. What exists is reused, and
// KubeDB issues only a certificate that is missing. TLS that is off stays off:
// switching it on during a handover would split the old pods from the new ones.
func (c *cluster) planTLS(ctx context.Context, o planOptions, src *Source, in *renderInput, fs *findings) error {
	cfg := src.Config
	proxy := src.Templates[src.firstUnit(roleProxy).key()]
	c.planEtcdTLS(ctx, proxy, cfg, in, fs)
	mode, err := externalTLSMode(cfg)
	if err != nil {
		fs.errorf("%v", err)
		return nil
	}
	internal := lookupBool(cfg, "common.security.internaltlsEnabled")
	if mode == "Disabled" && !internal {
		fs.infof("client and internal TLS are off and stay off")
		return nil
	}
	out := &adoptionTLS{ExternalMode: mode, Internal: internal}

	var external, inner *tlsMaterial
	if mode != "Disabled" {
		if external, err = c.readTLS(ctx, proxy, cfg, "tls.serverPemPath", "tls.serverKeyPath", "tls.caPemPath"); err != nil {
			fs.errorf("client TLS: %v", err)
			return nil
		}
		out.Server = external
	}
	if internal {
		if inner, err = c.readTLS(ctx, proxy, cfg, "internaltls.serverPemPath", "internaltls.serverKeyPath", "internaltls.caPemPath"); err != nil {
			fs.errorf("internal TLS: %v", err)
			return nil
		}
		if out.Server == nil {
			out.Server = inner
		}
	}

	if err := c.resolveIssuer(ctx, o, in.Name, out, fs); err != nil {
		return err
	}
	now := time.Now()
	if out.managed() {
		// The new certificate gets the Service names from the operator, so only
		// the names clients may use on top of those are carried over.
		checkCertificate("server certificate", out.Server, nil, true, now, fs)
		if external != nil && inner != nil && !bytes.Equal(external.Cert, inner.Cert) {
			if err := verifyAgainstCA(inner.Cert, external.CA, now); err != nil {
				fs.errorf("client and internal TLS use different CAs: %v", err)
				return nil
			}
		}
		if err := out.carryNames(in.Name, c.ns); err != nil {
			fs.errorf("server certificate: %v", err)
			return nil
		}
		fs.infof("TLS: client %s, internal %v; KubeDB issues and renews the certificates from %s %s, the CA clients already trust", mode, internal, out.IssuerKind, out.IssuerName)
		if len(out.ExtraDNSNames)+len(out.ExtraIPs) > 0 {
			fs.infof("the new server certificate also covers %s", strings.Join(append(append([]string{}, out.ExtraDNSNames...), out.ExtraIPs...), ", "))
		}
	} else {
		if external != nil && inner != nil && !bytes.Equal(external.Cert, inner.Cert) {
			fs.errorf("client and internal TLS use different certificates and KubeDB mounts one for both; give --tls-issuer or --tls-ca-secret so KubeDB issues one")
			return nil
		}
		checkCertificate("server certificate", out.Server, clientHosts(in.Name, c.ns), true, now, fs)
		fs.infof("TLS: client %s, internal %v; %s", mode, internal, materialNote(out.Server))
	}

	if mode == "mTLS" {
		if err := c.planClientCert(ctx, o, out, fs); err != nil {
			return err
		}
	}
	in.TLS = out
	c.setRESTTLS(out, in.Name)
	return nil
}

// issuerFromAnnotations reads the issuer cert-manager records on a Secret.
func issuerFromAnnotations(ann map[string]string) (kind, name string) {
	name = ann[annIssuerName]
	if name == "" {
		return "", ""
	}
	if kind = ann[annIssuerKind]; kind == "" {
		kind = "Issuer"
	}
	return kind, name
}

// resolveIssuer finds what lets KubeDB manage the certificates: an issuer the
// user names, the CA Secret the user gives, or the issuer that already signs the
// running certificate. Finding none leaves the running certificate in use.
func (c *cluster) resolveIssuer(ctx context.Context, o planOptions, name string, out *adoptionTLS, fs *findings) error {
	switch {
	case o.TLSCASecret != "":
		s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, o.TLSCASecret, metav1.GetOptions{})
		if err != nil {
			fs.errorf("--tls-ca-secret: %v", err)
			return nil
		}
		for _, k := range []string{tlsCertKey, tlsKeyKey} {
			if _, ok := s.Data[k]; !ok {
				fs.errorf("--tls-ca-secret: Secret %s has no key %s", s.Name, k)
				return nil
			}
		}
		if err := verifyAgainstCA(out.Server.Cert, s.Data[tlsCertKey], time.Now()); err != nil {
			fs.errorf("--tls-ca-secret: Secret %s is not the CA that signed the running certificate, so clients would reject what it issues: %v", s.Name, err)
			return nil
		}
		out.IssuerKind, out.IssuerName, out.CASecret = "Issuer", caIssuerName(name), s.Name
	case o.TLSIssuer != "":
		kind, issuer, err := parseIssuerRef(o.TLSIssuer)
		if err != nil {
			fs.errorf("--tls-issuer: %v", err)
			return nil
		}
		if c.checkIssuer(ctx, kind, issuer, out.Server.Cert, fs) {
			out.IssuerKind, out.IssuerName = kind, issuer
		}
	case out.Server.Source != "":
		s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, out.Server.Source, metav1.GetOptions{})
		if err != nil {
			return err
		}
		kind, issuer := issuerFromAnnotations(s.Annotations)
		if issuer == "" {
			return nil
		}
		if c.checkIssuer(ctx, kind, issuer, out.Server.Cert, fs) {
			out.IssuerKind, out.IssuerName = kind, issuer
			fs.infof("%s %s signs the running certificate", kind, issuer)
		}
	}
	return nil
}

// carryNames keeps the names of the running certificate that the operator does
// not add, such as an ingress host.
func (t *adoptionTLS) carryNames(name, ns string) error {
	leaf, _, err := parseLeaf(t.Server.Cert)
	if err != nil {
		return err
	}
	known := append(clientHosts(name, ns), "localhost")
	for _, d := range leaf.DNSNames {
		if !contains(known, d) {
			t.ExtraDNSNames = append(t.ExtraDNSNames, d)
		}
	}
	for _, ip := range leaf.IPAddresses {
		if !ip.IsLoopback() {
			t.ExtraIPs = append(t.ExtraIPs, ip.String())
		}
	}
	return nil
}

// planClientCert finds the certificate KubeDB presents when Milvus asks clients
// for one: the user's, or else one from the issuer.
func (c *cluster) planClientCert(ctx context.Context, o planOptions, out *adoptionTLS, fs *findings) error {
	if o.TLSClientSecret == "" {
		if !out.managed() {
			fs.errorf("client mTLS is on, so KubeDB needs a client certificate that Milvus trusts; give --tls-client-secret, or give --tls-issuer or --tls-ca-secret so KubeDB issues one")
		}
		return nil
	}
	s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, o.TLSClientSecret, metav1.GetOptions{})
	if err != nil {
		fs.errorf("--tls-client-secret: %v", err)
		return nil
	}
	for _, k := range []string{tlsCertKey, tlsKeyKey, tlsCAKey} {
		if _, ok := s.Data[k]; !ok {
			fs.errorf("--tls-client-secret: Secret %s has no key %s", s.Name, k)
			return nil
		}
	}
	checkCertificate("client certificate", &tlsMaterial{Cert: s.Data[tlsCertKey], Key: s.Data[tlsKeyKey], CA: out.Server.CA}, nil, true, time.Now(), fs)
	if err := verifyAgainstCA(out.Server.Cert, s.Data[tlsCAKey], time.Now()); err != nil {
		fs.errorf("--tls-client-secret: its ca.crt does not verify the server certificate: %v", err)
	}
	out.ClientSecret = s.Name
	fs.infof("client certificate: Secret %s is used as it is", s.Name)
	return nil
}

// checkIssuer reports whether KubeDB can use the issuer. Clients trust only the
// CA that signed the running certificate, so an issuer on another CA is refused.
// What plan cannot read here is checked by run, on the issued certificate.
func (c *cluster) checkIssuer(ctx context.Context, kind, name string, serverCert []byte, fs *findings) bool {
	if kind == "ClusterIssuer" {
		if _, err := c.dyn.Resource(certClusterIssuerGVR).Get(ctx, name, metav1.GetOptions{}); err != nil {
			fs.errorf("issuer: %v", err)
			return false
		}
		fs.warnf("cannot read the CA of ClusterIssuer %s here; run checks the certificate it issues before any pod moves", name)
		return true
	}
	obj, err := c.dyn.Resource(certIssuerGVR).Namespace(c.ns).Get(ctx, name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		fs.errorf("Issuer %s does not exist", name)
		return false
	}
	if err != nil {
		fs.errorf("issuer: %v", err)
		return false
	}
	spec, _ := obj.Object["spec"].(map[string]any)
	ca, _ := spec["ca"].(map[string]any)
	secret, _ := ca["secretName"].(string)
	if secret == "" {
		fs.warnf("Issuer %s is not a CA issuer; run checks the certificate it issues before any pod moves", name)
		return true
	}
	s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, secret, metav1.GetOptions{})
	if err != nil {
		fs.errorf("issuer: %v", err)
		return false
	}
	if err := verifyAgainstCA(serverCert, s.Data[tlsCertKey], time.Now()); err != nil {
		fs.errorf("Issuer %s signs with a different CA than the one clients trust, so they would reject what it issues", name)
		return false
	}
	return true
}

// checkIssuedCertificates stops an adoption whose new certificates clients would
// reject. It runs after the Milvus is created and before any pod moves.
func (c *cluster) checkIssuedCertificates(ctx context.Context, name string, t *adoptionTLS, timeout time.Duration) error {
	if t == nil || !t.managed() {
		return nil
	}
	names := []string{name + "-server-cert"}
	if t.ClientSecret == "" {
		names = append(names, name+"-client-cert")
	}
	deadline := time.Now().Add(timeout)
	for _, n := range names {
		var s *core.Secret
		for {
			var err error
			s, err = c.kube.CoreV1().Secrets(c.ns).Get(ctx, n, metav1.GetOptions{})
			if err == nil && len(s.Data[tlsCertKey]) > 0 {
				break
			}
			if err != nil && !kerr.IsNotFound(err) {
				return err
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("no pod was moved; certificate Secret %s was not issued in %s, check the Certificate and %s %s", n, timeout, t.IssuerKind, t.IssuerName)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.poll):
			}
		}
		if len(s.Data[tlsCAKey]) == 0 {
			return fmt.Errorf("no pod was moved; %s %s does not fill ca.crt in Secret %s, and KubeDB needs it to reach Milvus", t.IssuerKind, t.IssuerName, n)
		}
		if err := verifyAgainstCA(s.Data[tlsCertKey], t.Server.CA, time.Now()); err != nil {
			return fmt.Errorf("no pod was moved; the certificate in Secret %s does not verify against the CA clients trust: %w", n, err)
		}
	}
	c.logf("the issued certificates verify against the CA clients trust")
	return nil
}

// setRESTTLS makes the credential check speak TLS to the running proxies.
func (c *cluster) setRESTTLS(t *adoptionTLS, name string) {
	if t.ExternalMode == "Disabled" {
		return
	}
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(t.Server.CA)
	c.restTLS = &tls.Config{RootCAs: roots, ServerName: name, MinVersion: tls.VersionTLS12}
}

// tlsSpec is spec.tls of the adopted Milvus. serverSecret is empty when KubeDB
// issues the server certificate.
func tlsSpec(t *adoptionTLS, serverSecret string) map[string]any {
	internal := "Disabled"
	if t.Internal {
		internal = "TLS"
	}
	existing := map[string]any{}
	if serverSecret != "" {
		existing["server"] = map[string]any{"name": serverSecret}
	}
	if t.ClientSecret != "" {
		existing["client"] = map[string]any{"name": t.ClientSecret}
	}
	spec := map[string]any{
		"external": map[string]any{"mode": t.ExternalMode},
		"internal": map[string]any{"mode": internal},
	}
	if len(existing) > 0 {
		spec["existingSecrets"] = existing
	}
	if t.managed() {
		spec["issuerRef"] = map[string]any{"apiGroup": certManagerGroup, "kind": t.IssuerKind, "name": t.IssuerName}
		if len(t.ExtraDNSNames)+len(t.ExtraIPs) > 0 {
			server := map[string]any{"alias": "server"}
			if len(t.ExtraDNSNames) > 0 {
				server["dnsNames"] = toAnySlice(t.ExtraDNSNames)
			}
			if len(t.ExtraIPs) > 0 {
				server["ipAddresses"] = toAnySlice(t.ExtraIPs)
			}
			spec["certificates"] = []any{server}
		}
	}
	return spec
}
