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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"testing"
	"time"

	core "k8s.io/api/core/v1"
)

type testCA struct {
	cert *x509.Certificate
	key  *ecdsa.PrivateKey
	pem  []byte
}

func newTestCA(t *testing.T) *testCA {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "test-ca"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(24 * 365 * time.Hour),
		IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, _ := x509.ParseCertificate(der)
	return &testCA{cert: cert, key: key, pem: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})}
}

func (ca *testCA) issue(t *testing.T, ttl time.Duration, dns ...string) *tlsMaterial {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2), Subject: pkix.Name{CommonName: "milvus"}, DNSNames: dns,
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(ttl),
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca.cert, &key.PublicKey, ca.key)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return &tlsMaterial{
		Cert: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
		Key:  pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}),
		CA:   ca.pem,
	}
}

func TestCheckCertificate(t *testing.T) {
	ca := newTestCA(t)
	hosts := clientHosts("mv", "demo")
	year := 365 * 24 * time.Hour
	now := time.Now()

	good := ca.issue(t, year, hosts...)
	fs := &findings{}
	checkCertificate("server", good, hosts, true, now, fs)
	if fs.hasErrors() || len(*fs) != 0 {
		t.Fatalf("a good certificate has findings: %s", findingText(*fs))
	}

	fs = &findings{}
	checkCertificate("server", ca.issue(t, year, "other.example.com"), hosts, true, now, fs)
	if !fs.hasErrors() || !strings.Contains(findingText(*fs), "does not cover mv") {
		t.Fatalf("a wrong name must stop the plan: %s", findingText(*fs))
	}

	fs = &findings{}
	foreign := ca.issue(t, year, hosts...)
	foreign.CA = newTestCA(t).pem
	checkCertificate("server", foreign, hosts, true, now, fs)
	if !fs.hasErrors() {
		t.Fatal("a certificate from another CA must stop the plan")
	}

	fs = &findings{}
	checkCertificate("server", ca.issue(t, 5*24*time.Hour, hosts...), hosts, true, now, fs)
	if fs.hasErrors() || !strings.Contains(findingText(*fs), "expires") {
		t.Fatalf("a certificate near expiry is a warning: %s", findingText(*fs))
	}

	fs = &findings{}
	mismatched := ca.issue(t, year, hosts...)
	mismatched.Key = ca.issue(t, year, hosts...).Key
	checkCertificate("server", mismatched, hosts, true, now, fs)
	if !fs.hasErrors() {
		t.Fatal("a key of another certificate must stop the plan")
	}
}

func TestExternalTLSMode(t *testing.T) {
	for v, want := range map[any]string{nil: "Disabled", 0: "Disabled", 1: "TLS", float64(2): "mTLS"} {
		cfg := map[string]any{}
		if v != nil {
			cfg = map[string]any{"common": map[string]any{"security": map[string]any{"tlsMode": v}}}
		}
		if got, err := externalTLSMode(cfg); err != nil || got != want {
			t.Fatalf("tlsMode %v: got %q %v, want %q", v, got, err, want)
		}
	}
	if _, err := externalTLSMode(map[string]any{"common": map[string]any{"security": map[string]any{"tlsMode": 3}}}); err == nil {
		t.Fatal("an unknown mode must fail")
	}
}

func TestMountedFile(t *testing.T) {
	tmpl := core.PodTemplateSpec{Spec: core.PodSpec{
		Containers: []core.Container{{Name: "milvus", Image: "milvusdb/milvus", VolumeMounts: []core.VolumeMount{
			{Name: "certs", MountPath: "/certs"},
			{Name: "items", MountPath: "/items"},
			{Name: "proj", MountPath: "/proj"},
			{Name: "sub", MountPath: "/single/server.pem", SubPath: "tls.crt"},
			{Name: "disk", MountPath: "/disk"},
			{Name: "certs", MountPath: "/milvus/configs/cert"},
		}}},
		Volumes: []core.Volume{
			{Name: "certs", VolumeSource: core.VolumeSource{Secret: &core.SecretVolumeSource{SecretName: "s1"}}},
			{Name: "items", VolumeSource: core.VolumeSource{Secret: &core.SecretVolumeSource{
				SecretName: "s2", Items: []core.KeyToPath{{Key: "tls.crt", Path: "server.pem"}},
			}}},
			{Name: "proj", VolumeSource: core.VolumeSource{Projected: &core.ProjectedVolumeSource{Sources: []core.VolumeProjection{
				{Secret: &core.SecretProjection{LocalObjectReference: core.LocalObjectReference{Name: "s3"}, Items: []core.KeyToPath{{Key: "tls.key", Path: "server.key"}}}},
				{ConfigMap: &core.ConfigMapProjection{LocalObjectReference: core.LocalObjectReference{Name: "cm"}, Items: []core.KeyToPath{{Key: "ca", Path: "ca.pem"}}}},
			}}}},
			{Name: "sub", VolumeSource: core.VolumeSource{Secret: &core.SecretVolumeSource{SecretName: "s4"}}},
			{Name: "disk", VolumeSource: core.VolumeSource{EmptyDir: &core.EmptyDirVolumeSource{}}},
		},
	}}
	for file, want := range map[string]fileSource{
		"/certs/server.pem":   {"Secret", "s1", "server.pem"},
		"/items/server.pem":   {"Secret", "s2", "tls.crt"},
		"/proj/server.key":    {"Secret", "s3", "tls.key"},
		"/proj/ca.pem":        {"ConfigMap", "cm", "ca"},
		"/single/server.pem":  {"Secret", "s4", "tls.crt"},
		"configs/cert/ca.pem": {"Secret", "s1", "ca.pem"}, // relative to the Milvus work dir
	} {
		got, err := mountedFile(tmpl, file)
		if err != nil || len(got) != 1 || got[0] != want {
			t.Fatalf("%s: got %+v %v, want %+v", file, got, err, want)
		}
	}
	for _, file := range []string{"/nowhere/x.pem", "/disk/x.pem", "/items/other.pem"} {
		if _, err := mountedFile(tmpl, file); err == nil {
			t.Fatalf("%s must fail", file)
		}
	}
}

func TestParseIssuerRef(t *testing.T) {
	for in, want := range map[string][2]string{
		"ca": {"Issuer", "ca"}, "Issuer/ca": {"Issuer", "ca"}, "ClusterIssuer/ca": {"ClusterIssuer", "ca"},
	} {
		if k, n, err := parseIssuerRef(in); err != nil || k != want[0] || n != want[1] {
			t.Fatalf("%s: got %s %s %v", in, k, n, err)
		}
	}
	for _, in := range []string{"Certificate/ca", "Issuer/", ""} {
		if _, _, err := parseIssuerRef(in); err == nil {
			t.Fatalf("%q must fail", in)
		}
	}
}

func TestCarriedConfigDropsTLS(t *testing.T) {
	raw, err := carriedConfig(map[string]any{
		"tls":         map[string]any{"serverPemPath": "/old/server.pem"},
		"internaltls": map[string]any{"serverPemPath": "/old/server.pem"},
		"common": map[string]any{"security": map[string]any{
			"tlsMode": 1, "internaltlsEnabled": true, "authorizationEnabled": true,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"/old/", "tlsMode", "internaltlsEnabled"} {
		if strings.Contains(raw, s) {
			t.Fatalf("%s must not be carried, the operator renders it: %s", s, raw)
		}
	}
	if !strings.Contains(raw, "authorizationEnabled") {
		t.Fatalf("other security settings stay: %s", raw)
	}
}

func TestConfigMismatchesTLS(t *testing.T) {
	sec := func(v map[string]any) map[string]any {
		return map[string]any{"common": map[string]any{"security": v}, "etcd": map[string]any{"rootPath": "by-dev"}}
	}
	if d := configMismatches(sec(map[string]any{"tlsMode": 0, "internaltlsEnabled": false}), sec(map[string]any{})); len(d) != 0 {
		t.Fatalf("off and not rendered are the same: %v", d)
	}
	if d := configMismatches(sec(map[string]any{"tlsMode": 1}), sec(map[string]any{})); len(d) != 1 {
		t.Fatalf("TLS on and not rendered must be reported: %v", d)
	}
}

func TestTLSSpec(t *testing.T) {
	got := tlsSpec(&adoptionTLS{ExternalMode: "mTLS", Internal: true, IssuerKind: "Issuer", IssuerName: "ca"}, "theirs")
	if got["external"].(map[string]any)["mode"] != "mTLS" || got["internal"].(map[string]any)["mode"] != "TLS" {
		t.Fatalf("modes: %v", got)
	}
	if got["existingSecrets"].(map[string]any)["server"].(map[string]any)["name"] != "theirs" {
		t.Fatalf("server Secret: %v", got)
	}
	if got["issuerRef"].(map[string]any)["apiGroup"] != certManagerGroup {
		t.Fatalf("issuer: %v", got)
	}
	got = tlsSpec(&adoptionTLS{ExternalMode: "TLS"}, "theirs")
	if _, ok := got["issuerRef"]; ok {
		t.Fatalf("nothing to issue, no issuerRef: %v", got)
	}

	got = tlsSpec(&adoptionTLS{
		ExternalMode: "TLS", IssuerKind: "Issuer", IssuerName: "ca",
		ExtraDNSNames: []string{"milvus.example.com"}, ExtraIPs: []string{"10.0.0.1"},
	}, "")
	if _, ok := got["existingSecrets"]; ok {
		t.Fatalf("KubeDB issues everything, no existing Secret: %v", got)
	}
	server := got["certificates"].([]any)[0].(map[string]any)
	if server["alias"] != "server" || server["dnsNames"].([]any)[0] != "milvus.example.com" || server["ipAddresses"].([]any)[0] != "10.0.0.1" {
		t.Fatalf("extra names: %v", got)
	}
}

func TestIssuerFromAnnotations(t *testing.T) {
	if k, n := issuerFromAnnotations(map[string]string{annIssuerName: "ca"}); k != "Issuer" || n != "ca" {
		t.Fatalf("got %s %s", k, n)
	}
	if k, n := issuerFromAnnotations(map[string]string{annIssuerName: "ca", annIssuerKind: "ClusterIssuer"}); k != "ClusterIssuer" || n != "ca" {
		t.Fatalf("got %s %s", k, n)
	}
	if _, n := issuerFromAnnotations(map[string]string{"other": "x"}); n != "" {
		t.Fatalf("no issuer expected, got %s", n)
	}
}

func TestCarryNames(t *testing.T) {
	ca := newTestCA(t)
	hosts := clientHosts("mv", "demo")
	m := ca.issue(t, 365*24*time.Hour, append(hosts, "localhost", "milvus.example.com")...)
	out := &adoptionTLS{Server: m}
	if err := out.carryNames("mv", "demo"); err != nil {
		t.Fatal(err)
	}
	if len(out.ExtraDNSNames) != 1 || out.ExtraDNSNames[0] != "milvus.example.com" {
		t.Fatalf("only names the operator does not add are carried: %v", out.ExtraDNSNames)
	}
}
