/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package restic

import (
	"fmt"
	"os"
	"path/filepath"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	storage "kmodules.xyz/objectstore-api/api/v1"
)

func (w *ResticWrapper) setupEnv() error {
	// Set progress report frequency.
	// 0.016666 is for one report per minute.
	// ref: https://restic.readthedocs.io/en/stable/manual_rest.html
	w.sh.SetEnv(RESTIC_PROGRESS_FPS, "0.016666")
	if w.Config.EnableCache {
		cacheDir := filepath.Join(w.Config.ScratchDir, resticCacheDir)
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			return err
		}
	}
	var errs []error
	for _, b := range w.Config.Backends {
		err := w.setupEnvsForBackend(b)
		if err != nil {
			b.Error = errors.NewAggregate([]error{b.Error, err})
			errs = append(errs, err)
		}
	}
	return errors.NewAggregate(errs)
}

func (w *ResticWrapper) setupEnvsForBackend(b *Backend) error {
	// Use the injected ConfigResolver to get storage configuration
	if b.ConfigResolver == nil {
		return fmt.Errorf("ConfigResolver is not set for backend %s", b.Repository)
	}
	if err := b.ConfigResolver(b); err != nil {
		return fmt.Errorf("failed to resolve storage config: %w", err)
	}

	if b.Envs == nil {
		b.Envs = make(map[string]string)
	}
	if err := w.setEnvFromSecretIfExists(b.Envs, b.EncryptionSecret, RESTIC_PASSWORD, true); err != nil {
		return fmt.Errorf("failed to set secret for backend %s: %w", b.Repository, err)
	}

	tmpDir, err := os.MkdirTemp(w.Config.ScratchDir, fmt.Sprintf("%s-tmp-", filepath.Base(b.Repository)))
	if err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}
	b.Envs[TMPDIR] = tmpDir
	if w.isSecretKeyExist(b.StorageSecret, CA_CERT_DATA) {
		filePath, err := w.writeSecretKeyToFile(tmpDir, b.StorageSecret, CA_CERT_DATA, "ca.crt")
		if err != nil {
			return fmt.Errorf("failed to write secret for backend %s: %w", b.Repository, err)
		}
		b.CaCertFile = filePath
	}

	switch b.Provider {
	case storage.ProviderLocal:
		b.Envs[RESTIC_REPOSITORY] = fmt.Sprintf("%s/%s", b.Bucket, b.Directory)

	case storage.ProviderS3:
		b.Envs[RESTIC_REPOSITORY] = fmt.Sprintf("s3:%s/%s", b.Endpoint, filepath.Join(b.Bucket, b.Prefix, b.Directory))
		if err := w.setEnvFromSecretIfExists(b.Envs, b.StorageSecret, AWS_ACCESS_KEY_ID, false); err != nil {
			return fmt.Errorf("failed to set secret for backend %s: %w", b.Repository, err)
		}
		if err := w.setEnvFromSecretIfExists(b.Envs, b.StorageSecret, AWS_SECRET_ACCESS_KEY, false); err != nil {
			return fmt.Errorf("failed to set secret for backend %s: %w", b.Repository, err)
		}
		if b.Region != "" {
			b.Envs[AWS_DEFAULT_REGION] = b.Region
		}

	case storage.ProviderAzure:
		b.Envs[RESTIC_REPOSITORY] = fmt.Sprintf("azure:%s:/%s", b.Bucket, filepath.Join(b.Prefix, b.Directory))
		if b.AzureStorageAccount == "" {
			return fmt.Errorf("missing storage account for Azure storage")
		}
		b.Envs[AZURE_ACCOUNT_NAME] = b.AzureStorageAccount
		if err := w.setEnvFromSecretIfExists(b.Envs, b.StorageSecret, AZURE_ACCOUNT_KEY, false); err != nil {
			return fmt.Errorf("failed to set secret for backend %s: %w", b.Repository, err)
		}

	case storage.ProviderGCS:
		b.Envs[RESTIC_REPOSITORY] = fmt.Sprintf("gs:%s:/%s", b.Bucket, filepath.Join(b.Prefix, b.Directory))
		if w.isSecretKeyExist(b.StorageSecret, GOOGLE_SERVICE_ACCOUNT_JSON_KEY) {
			filePath, err := w.writeSecretKeyToFile(tmpDir, b.StorageSecret, GOOGLE_SERVICE_ACCOUNT_JSON_KEY, GOOGLE_SERVICE_ACCOUNT_JSON_KEY)
			if err != nil {
				return err
			}
			b.Envs[GOOGLE_APPLICATION_CREDENTIALS] = filePath
		}
	default:
		return fmt.Errorf("unsupported storage provider: %s", b.Provider)
	}
	return nil
}

// nolint: unused
func (w *ResticWrapper) exportSecretKey(secret *core.Secret, key string, required bool) error {
	if v, ok := secret.Data[key]; !ok {
		if required {
			return fmt.Errorf("storage Secret missing %s key", key)
		}
	} else {
		w.sh.SetEnv(key, string(v))
	}
	return nil
}

func (w *ResticWrapper) setEnvFromSecretIfExists(envs map[string]string, secret *core.Secret, key string, required bool) error {
	if required && secret == nil {
		return fmt.Errorf("storage Secret is Required")
	}
	v, ok := secret.Data[key]
	if !ok {
		if required {
			return fmt.Errorf("%s storage Secret missing %s key", secret.Name, key)
		}
	}
	envs[key] = string(v)
	return nil
}

func (w *ResticWrapper) isSecretKeyExist(secret *core.Secret, key string) bool {
	if secret == nil {
		return false
	}
	_, ok := secret.Data[key]
	return ok
}

func (w *ResticWrapper) writeSecretKeyToFile(tmpDir string, secret *core.Secret, key, name string) (string, error) {
	v, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("storage Secret missing %s key", key)
	}

	filePath := filepath.Join(tmpDir, name)

	if err := os.WriteFile(filePath, v, 0o755); err != nil {
		return "", err
	}
	return filePath, nil
}
