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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"reflect"
	"testing"
	"time"
)

// makeCA returns a self-signed CA certificate and private key as PEM. keyKind
// selects the key algorithm and PEM encoding, to cover the parse variants.
func makeCA(t *testing.T, keyKind string, notAfter time.Time) (certPEM, keyPEM []byte) {
	t.Helper()
	var pub any
	var signer any
	var keyBlock *pem.Block
	switch keyKind {
	case "rsa-pkcs1":
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		pub, signer = key.Public(), key
		keyBlock = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	case "ec":
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		der, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			t.Fatal(err)
		}
		pub, signer = key.Public(), key
		keyBlock = &pem.Block{Type: "EC PRIVATE KEY", Bytes: der}
	case "rsa-pkcs8":
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		der, err := x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			t.Fatal(err)
		}
		pub, signer = key.Public(), key
		keyBlock = &pem.Block{Type: "PRIVATE KEY", Bytes: der}
	default:
		t.Fatalf("unknown keyKind %q", keyKind)
	}

	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              notAfter,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, signer)
	if err != nil {
		t.Fatal(err)
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM = pem.EncodeToMemory(keyBlock)
	return certPEM, keyPEM
}

func parseCertPEM(t *testing.T, certPEM []byte) *x509.Certificate {
	t.Helper()
	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("no PEM block in certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}

func TestIssueClientCertFromCA(t *testing.T) {
	for _, keyKind := range []string{"rsa-pkcs1", "ec", "rsa-pkcs8"} {
		t.Run(keyKind, func(t *testing.T) {
			caPEM, caKeyPEM := makeCA(t, keyKind, time.Now().Add(10*365*24*time.Hour))
			sans := []string{"replica.example.com", "replica-dr.example.com"}

			certPEM, keyPEM, err := issueClientCertFromCA("repluser", sans, caPEM, caKeyPEM)
			if err != nil {
				t.Fatalf("issueClientCertFromCA: %v", err)
			}

			cert := parseCertPEM(t, certPEM)
			caCert := parseCertPEM(t, caPEM)

			if cert.Subject.CommonName != "repluser" {
				t.Errorf("CommonName = %q, want repluser", cert.Subject.CommonName)
			}
			if !reflect.DeepEqual(cert.DNSNames, sans) {
				t.Errorf("DNSNames = %v, want %v", cert.DNSNames, sans)
			}
			if err := cert.CheckSignatureFrom(caCert); err != nil {
				t.Errorf("client cert is not signed by the CA: %v", err)
			}
			hasClientAuth := false
			for _, u := range cert.ExtKeyUsage {
				if u == x509.ExtKeyUsageClientAuth {
					hasClientAuth = true
				}
			}
			if !hasClientAuth {
				t.Error("client cert lacks ExtKeyUsage clientAuth")
			}
			if cert.IsCA {
				t.Error("client cert must not be a CA")
			}

			// The private key must match the certificate.
			keyBlock, _ := pem.Decode(keyPEM)
			key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
			if err != nil {
				t.Fatalf("parse client key: %v", err)
			}
			rsaKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				t.Fatalf("client key is %T, want *rsa.PrivateKey", key)
			}
			if !rsaKey.PublicKey.Equal(cert.PublicKey) {
				t.Error("client key does not match client cert")
			}
		})
	}
}

func TestIssueClientCertValidityClampedToCA(t *testing.T) {
	caExpiry := time.Now().Add(30 * 24 * time.Hour) // CA dies before the 1y default
	caPEM, caKeyPEM := makeCA(t, "rsa-pkcs1", caExpiry)

	certPEM, _, err := issueClientCertFromCA("repluser", nil, caPEM, caKeyPEM)
	if err != nil {
		t.Fatal(err)
	}
	cert := parseCertPEM(t, certPEM)
	if cert.NotAfter.After(caExpiry.Add(time.Minute)) {
		t.Errorf("client cert NotAfter %v outlives CA NotAfter %v", cert.NotAfter, caExpiry)
	}
}

func TestIssueClientCertRejectsMismatchedKey(t *testing.T) {
	caPEM, _ := makeCA(t, "rsa-pkcs1", time.Now().Add(24*time.Hour))
	_, otherKeyPEM := makeCA(t, "rsa-pkcs1", time.Now().Add(24*time.Hour))

	if _, _, err := issueClientCertFromCA("repluser", nil, caPEM, otherKeyPEM); err == nil {
		t.Fatal("expected error for CA cert/key mismatch, got nil")
	}
}

func TestIssueClientCertRejectsNonCA(t *testing.T) {
	// Build a leaf (IsCA=false) and try to use it as the CA.
	caPEM, caKeyPEM := makeCA(t, "rsa-pkcs1", time.Now().Add(24*time.Hour))
	leafPEM, leafKeyPEM, err := issueClientCertFromCA("leaf", nil, caPEM, caKeyPEM)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := issueClientCertFromCA("repluser", nil, leafPEM, leafKeyPEM); err == nil {
		t.Fatal("expected error when --ca-cert is not a CA, got nil")
	}
}

func TestIssueClientCertRejectsGarbage(t *testing.T) {
	if _, _, err := issueClientCertFromCA("u", nil, []byte("not pem"), []byte("not pem")); err == nil {
		t.Fatal("expected error for garbage PEM, got nil")
	}
}
