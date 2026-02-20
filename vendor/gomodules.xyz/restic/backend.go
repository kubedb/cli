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

	core "k8s.io/api/core/v1"
	storage "kmodules.xyz/objectstore-api/api/v1"
)

// StorageConfig contains provider-agnostic storage configuration.
// This struct decouples the restic package from specific storage CRD types.
type StorageConfig struct {
	Provider            string
	Bucket              string
	Endpoint            string
	Region              string
	Prefix              string
	InsecureTLS         bool
	MaxConnections      int64
	AzureStorageAccount string
}

// StorageConfigResolver is a function type that resolves storage configuration.
// This allows callers to inject their own logic for fetching storage info
// from any storage type (e.g., BackupStorage, or custom types in other projects).
type StorageConfigResolver func(b *Backend) error

// Backend represents a backup storage backend with its configuration and runtime state
type Backend struct {
	*StorageConfig
	// ConfigResolver is called during setup to populate storage configuration.
	// Callers must provide this function to resolve storage info from their storage type.
	ConfigResolver StorageConfigResolver

	Repository string
	Directory  string
	MountPath  string

	// Secrets for accessing the Backend Storage
	StorageSecret    *core.Secret
	EncryptionSecret *core.Secret

	CaCertFile string
	Envs       map[string]string
	Error      error
}

func (b *Backend) createLocalDir() error {
	if b.Provider == storage.ProviderLocal {
		return os.MkdirAll(b.Bucket, 0o755)
	}
	return nil
}

func (b *Backend) appendInsecureTLSFlag(args []any) []any {
	if b.InsecureTLS {
		return append(args, "--insecure-tls")
	}
	return args
}

func (b *Backend) appendCaCertFlag(args []any) []any {
	if b.CaCertFile != "" {
		return append(args, "--cacert", b.CaCertFile)
	}
	return args
}

func (b *Backend) appendMaxConnectionsFlag(args []any) []any {
	var maxConOption string
	if b.MaxConnections > 0 {
		switch b.Provider {
		case storage.ProviderGCS:
			maxConOption = fmt.Sprintf("gs.connections=%d", b.MaxConnections)
		case storage.ProviderAzure:
			maxConOption = fmt.Sprintf("azure.connections=%d", b.MaxConnections)
		case storage.ProviderB2:
			maxConOption = fmt.Sprintf("b2.connections=%d", b.MaxConnections)
		}
	}
	if maxConOption != "" {
		return append(args, "--option", maxConOption)
	}
	return args
}

func (b *Backend) GetCaPath() string {
	return b.CaCertFile
}
