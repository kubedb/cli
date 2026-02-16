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

package v1alpha1

import kmapi "kmodules.xyz/client-go/api/v1"

// MigrationConfig defines the desired state of Migrator
type MigrationConfig struct {
	Source Source `yaml:"source" json:"source"`
	Target Target `yaml:"target" json:"target"`
}

// Source defines the source database configuration
type Source struct {
	// Postgres refers to the source Postgres database configuration
	Postgres *PostgresSource `yaml:"postgres" json:"postgres,omitempty"`
}

// Target defines the target database configuration
type Target struct {
	// Postgres refers to the target Postgres database configuration
	Postgres *PostgresTarget `yaml:"postgres" json:"postgres,omitempty"`
}

type ConnectionInfo struct {
	// AppBinding refers to the source database AppBinding name, Who contains the connection information.
	// +optional
	AppBinding kmapi.ObjectReference `yaml:"appBinding,omitempty" json:"appBinding,omitempty"`

	// DBName refers to the database name.
	// +optional
	DBName string `yaml:"dbName" json:"dbName"`

	// URL refers to the database connection string.e.g postgres://postgres:password@localhost:5432/postgres
	// +optional
	URL string `yaml:"url" json:"url,omitempty"`

	// MaxConnections refers to the `MaxConns`,which means the maximum size of the pool.
	// The default is the greater of 4 or runtime.NumCPU().
	// +optional
	MaxConnections *int32 `yaml:"maxConnections" json:"maxConnections,omitempty"`
}

type DBMigratorImages struct {
	// CLI specifies the migrator CLI image
	// +optional
	CLI string `json:"cli,omitempty"`
	// StatusReporter is the sidecar image used to report migration progress
	// +optional
	StatusReporter string `json:"statusReporter,omitempty"`
}
