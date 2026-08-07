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

package remote_replica

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"kubedb.dev/cli/pkg/common"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// clientCertValidity is how long a locally issued client certificate stays valid.
// It is clamped to the CA's own NotAfter, since a leaf outliving its CA is useless.
const clientCertValidity = 365 * 24 * time.Hour

// parseCAPair decodes a PEM CA certificate and its private key and verifies they
// belong together. Signing a client certificate requires the CA's PRIVATE key;
// the public ca.crt alone cannot issue anything.
func parseCAPair(caCertPEM, caKeyPEM []byte) (*x509.Certificate, crypto.Signer, error) {
	block, _ := pem.Decode(caCertPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("--ca-cert does not contain a PEM CERTIFICATE block")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse --ca-cert: %v", err)
	}
	if !caCert.IsCA {
		return nil, nil, fmt.Errorf("--ca-cert is not a CA certificate (BasicConstraints CA=false); it cannot sign client certificates")
	}

	keyBlock, _ := pem.Decode(caKeyPEM)
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("--ca-key does not contain a PEM block")
	}
	var key any
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(keyBlock.Bytes)
	case "PRIVATE KEY":
		key, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	default:
		return nil, nil, fmt.Errorf("--ca-key: unsupported PEM block type %q", keyBlock.Type)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse --ca-key: %v", err)
	}
	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, nil, fmt.Errorf("--ca-key is not a usable signing key")
	}

	type pubEqualer interface{ Equal(crypto.PublicKey) bool }
	pub, ok := caCert.PublicKey.(pubEqualer)
	if !ok || !pub.Equal(signer.Public()) {
		return nil, nil, fmt.Errorf("--ca-key does not match --ca-cert (public keys differ)")
	}
	return caCert, signer, nil
}

// issueClientCertFromCA generates a fresh RSA key pair and a client-auth
// certificate for userName signed by the given CA. The certificate's CommonName
// is the username — that is what PostgreSQL cert authentication maps to the
// database role — and dnsSANs land in the SAN extension.
func issueClientCertFromCA(userName string, dnsSANs []string, caCertPEM, caKeyPEM []byte) (certPEM, keyPEM []byte, err error) {
	caCert, caKey, err := parseCAPair(caCertPEM, caKeyPEM)
	if err != nil {
		return nil, nil, err
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate client key: %v", err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %v", err)
	}

	now := time.Now()
	notAfter := now.Add(clientCertValidity)
	if notAfter.After(caCert.NotAfter) {
		notAfter = caCert.NotAfter
	}
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: userName},
		DNSNames:              dnsSANs,
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign client certificate: %v", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal client key: %v", err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	return certPEM, keyPEM, nil
}

// generateTlsSecretFromLocalCA issues the client certificate locally from the CA
// files given on the command line — no cert-manager involved — and packages it as
// the kubernetes.io/tls Secret the remote replica will mount (ca.crt, tls.crt,
// tls.key; the operator remaps tls.* to client.* at mount time).
func generateTlsSecretFromLocalCA(userName, ns, caCertPath, caKeyPath string, dnsSANs []string, opts *common.PostgresOpts) ([]byte, string, error) {
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read --ca-cert: %v", err)
	}
	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read --ca-key: %v", err)
	}
	certPEM, keyPEM, err := issueClientCertFromCA(userName, dnsSANs, caCertPEM, caKeyPEM)
	if err != nil {
		return nil, "", err
	}

	tlsSecret := core.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindSecret,
			APIVersion: ApiversionV1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-remote-replica-client-cert", opts.DB.Name),
			Namespace: ns,
		},
		Type: core.SecretTypeTLS,
		Data: map[string][]byte{
			"ca.crt":  caCertPEM,
			"tls.crt": certPEM,
			"tls.key": keyPEM,
		},
	}
	tlsSecretYaml, err := yaml.Marshal(tlsSecret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal tls secret yaml %v", err)
	}
	buffer := make([]byte, 0, len(tlsSecretYaml)+4)
	buffer = append(buffer, tlsSecretYaml...)
	buffer = append(buffer, []byte("---\n")...)
	return buffer, tlsSecret.Name, nil
}
