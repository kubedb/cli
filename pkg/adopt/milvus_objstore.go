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
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultObjectStorageCAKey = "ca.crt"

// objectStorageTLS is the TLS the Milvus uses to reach the object storage after
// the adoption: spec.objectStorage.tls, and the endpoint when it changes. The
// running cluster may not use TLS at all.
type objectStorageTLS struct {
	Kind    string // "Secret" or "ConfigMap"
	Name    string
	Key     string
	Address string // "" keeps the running address
	Port    int    // 0 keeps the running port
}

// parseCARef splits NAME or NAME/KEY.
func parseCARef(s string) (name, key string, err error) {
	name, key, _ = strings.Cut(s, "/")
	if key == "" {
		key = defaultObjectStorageCAKey
	}
	if name == "" || strings.Contains(key, "/") {
		return "", "", fmt.Errorf("%q is not NAME or NAME/KEY", s)
	}
	return name, key, nil
}

// specTLS is spec.objectStorage.tls.
func (t *objectStorageTLS) specTLS() map[string]any {
	field := "caSecret"
	if t.Kind == "ConfigMap" {
		field = "caConfigMap"
	}
	return map[string]any{field: map[string]any{"name": t.Name, "key": t.Key}}
}

// storageData applies the endpoint and TLS switch to the object storage Secret data.
func (t *objectStorageTLS) storageData(data map[string]string) map[string]string {
	if t == nil {
		return data
	}
	data["useSSL"] = "true"
	if t.Address != "" {
		data["address"] = t.Address
	}
	if t.Port > 0 {
		data["port"] = fmt.Sprint(t.Port)
	}
	return data
}

// carry returns the running config without the settings the operator renders
// from spec.objectStorage: spec.configuration wins a merge, so the running
// plain-HTTP settings would keep the new pods off TLS.
func (t *objectStorageTLS) carry(cfg map[string]any) map[string]any {
	if t == nil {
		return cfg
	}
	out := deepCopy(cfg)
	if minio, ok := out["minio"].(map[string]any); ok {
		delete(minio, "useSSL")
		delete(minio, "ssl")
		if t.Address != "" {
			delete(minio, "address")
		}
		if t.Port > 0 {
			delete(minio, "port")
		}
	}
	return out
}

// expected is the running config as the operator should render it after the
// switch, for the check that nothing else changed.
func (t *objectStorageTLS) expected(cfg map[string]any) map[string]any {
	if t == nil {
		return cfg
	}
	out := deepCopy(cfg)
	minio, ok := out["minio"].(map[string]any)
	if !ok {
		minio = map[string]any{}
		out["minio"] = minio
	}
	minio["useSSL"] = true
	if t.Address != "" {
		minio["address"] = t.Address
	}
	if t.Port > 0 {
		minio["port"] = t.Port
	}
	return out
}

// planObjectStorageTLS checks --object-storage-* and sets in.ObjStore.
func (c *cluster) planObjectStorageTLS(ctx context.Context, o planOptions, cfg map[string]any, in *renderInput, fs *findings) error {
	if o.ObjectStorageCASecret == "" && o.ObjectStorageCAConfigMap == "" {
		if o.ObjectStorageAddress != "" || o.ObjectStoragePort != 0 {
			fs.errorf("--object-storage-address and --object-storage-port need --object-storage-ca-secret or --object-storage-ca-configmap: a new endpoint without its CA cannot be verified")
		}
		return nil
	}
	if o.ObjectStorageCASecret != "" && o.ObjectStorageCAConfigMap != "" {
		fs.errorf("give one of --object-storage-ca-secret and --object-storage-ca-configmap, not both")
		return nil
	}
	if o.ObjectStoragePort < 0 || o.ObjectStoragePort > 65535 {
		fs.errorf("--object-storage-port %d is not a port", o.ObjectStoragePort)
		return nil
	}
	t := &objectStorageTLS{Kind: "Secret", Address: o.ObjectStorageAddress, Port: o.ObjectStoragePort}
	flag, ref := "--object-storage-ca-secret", o.ObjectStorageCASecret
	if o.ObjectStorageCAConfigMap != "" {
		t.Kind, flag, ref = "ConfigMap", "--object-storage-ca-configmap", o.ObjectStorageCAConfigMap
	}
	var err error
	if t.Name, t.Key, err = parseCARef(ref); err != nil {
		fs.errorf("%s: %v", flag, err)
		return nil
	}
	var pemBytes []byte
	if t.Kind == "Secret" {
		s, err := c.kube.CoreV1().Secrets(c.ns).Get(ctx, t.Name, metav1.GetOptions{})
		if kerr.IsNotFound(err) {
			fs.errorf("%s: Secret %s not found in namespace %s (the Milvus pods mount it, so it must be there)", flag, t.Name, c.ns)
			return nil
		} else if err != nil {
			return err
		}
		pemBytes = s.Data[t.Key]
	} else {
		cm, err := c.kube.CoreV1().ConfigMaps(c.ns).Get(ctx, t.Name, metav1.GetOptions{})
		if kerr.IsNotFound(err) {
			fs.errorf("%s: ConfigMap %s not found in namespace %s (the Milvus pods mount it, so it must be there)", flag, t.Name, c.ns)
			return nil
		} else if err != nil {
			return err
		}
		pemBytes = []byte(cm.Data[t.Key])
	}
	if len(pemBytes) == 0 {
		fs.errorf("%s: %s %s has no key %s", flag, t.Kind, t.Name, t.Key)
		return nil
	}
	subject, notAfter, err := validateCABundle(pemBytes, time.Now())
	if err != nil {
		fs.errorf("%s: %s %s key %s: %v", flag, t.Kind, t.Name, t.Key, err)
		return nil
	}
	in.ObjStore = t

	osData := t.storageData(objectStorageData(cfg))
	fs.infof("object storage over TLS: %s:%s, bucket %s; CA %q from %s %s key %s (valid until %s)",
		osData["address"], osData["port"], osData["bucketName"], subject, t.Kind, t.Name, t.Key, notAfter.Format("2006-01-02"))
	if lookupBool(cfg, "minio.useSSL") {
		fs.warnf("the running cluster already uses TLS to the object storage; the CA you give replaces its minio.ssl settings")
	}
	fs.infof("run tests TLS to %s:%s from a short-lived Pod before it creates anything (skip with --skip-object-storage-check)", osData["address"], osData["port"])
	if t.Address == "" && t.Port == 0 {
		fs.warnf("the object storage endpoint stays %s:%s; it must serve TLS there, and its certificate must name that host", osData["address"], osData["port"])
	}
	if strings.EqualFold(osData["useIAM"], "true") {
		fs.warnf("minio.useIAM is true; confirm the TLS endpoint accepts the same IAM credentials")
	}
	return nil
}

// validateCABundle checks that data holds PEM certificates and that at least one
// is valid at now. It returns the subject of the first valid one and its expiry.
func validateCABundle(data []byte, now time.Time) (subject string, notAfter time.Time, err error) {
	var n int
	for rest := data; ; {
		var blk *pem.Block
		blk, rest = pem.Decode(rest)
		if blk == nil {
			break
		}
		if blk.Type != "CERTIFICATE" {
			continue
		}
		cert, perr := x509.ParseCertificate(blk.Bytes)
		if perr != nil {
			return "", time.Time{}, fmt.Errorf("a certificate in the bundle does not parse: %v", perr)
		}
		n++
		if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
			continue
		}
		if subject == "" {
			subject, notAfter = cert.Subject.String(), cert.NotAfter
		}
	}
	switch {
	case n == 0:
		return "", time.Time{}, fmt.Errorf("holds no PEM certificate")
	case subject == "":
		return "", time.Time{}, fmt.Errorf("every certificate in it is expired or not yet valid")
	}
	return subject, notAfter, nil
}
