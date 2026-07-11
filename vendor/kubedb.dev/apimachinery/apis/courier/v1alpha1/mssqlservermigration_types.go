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

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindMSSQLServerMigration     = "MSSQLServerMigration"
	ResourceSingularMSSQLServerMigration = "mssqlservermigration"
	ResourcePluralMSSQLServerMigrations  = "mssqlservermigrations"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mssqlservermigrations,singular=mssqlservermigration,shortName=msmig,categories={kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Stage",type="string",JSONPath=".status.progress.info.Stage"
// +kubebuilder:printcolumn:name="Lag",type="string",JSONPath=".status.progress.info.Lag"
// +kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress.info.Progress"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MSSQLServerMigration struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of MSSQLServerMigration
	// +required
	Spec MSSQLServerMigrationSpec `json:"spec"`

	// status defines the observed state of MSSQLServerMigration.
	// It reuses the shared MigrationStatus so that the Migration duck type can
	// project it and the operator's status patches replay onto it unchanged.
	// +optional
	Status MigrationStatus `json:"status,omitzero"`
}

// MSSQLServerMigrationSpec defines the desired state of MSSQLServerMigration
type MSSQLServerMigrationSpec struct {
	// Source defines the source MSSQL Server database configuration
	Source MSSQLServerSource `json:"source"`

	// Target defines the target MSSQL Server database configuration
	Target MSSQLServerTarget `json:"target"`

	// JobDefaults specifies default settings for migration jobs
	// +optional
	JobDefaults *JobDefaults `json:"jobDefaults,omitempty"`

	// JobTemplate specifies runtime configurations for the migration Job
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
}

// MSSQLServerMigrationList contains a list of MSSQLServerMigration

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MSSQLServerMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []MSSQLServerMigration `json:"items"`
}

type MSSQLServerSource struct {
	// ConnectionInfo refers to the source MSSQL Server database connection information.
	ConnectionInfo *MSSQLServerConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
	Schema         *MSSQLServerSchema         `yaml:"schema" json:"schema,omitempty"`
	Snapshot       *MSSQLServerSnapshot       `yaml:"snapshot" json:"snapshot,omitempty"`
	Streaming      *MSSQLServerStreaming      `yaml:"streaming" json:"streaming,omitempty"`
}

type MSSQLServerTarget struct {
	// ConnectionInfo refers to the target MSSQL Server database connection information.
	ConnectionInfo *MSSQLServerConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
}

type MSSQLServerSchema struct {
	// Enabled controls whether the Schema Phase should be executed.
	Enabled bool `yaml:"enabled" json:"enabled"`
	// Database is the list of databases to migrate.
	// +optional
	Database []string `yaml:"database" json:"database,omitempty"`
	// Schema is the list of SQL Server schemas (e.g. "dbo") to include.
	// +optional
	Schema []string `yaml:"schema" json:"schema,omitempty"`
	// ExcludeSchema is the list of SQL Server schemas to exclude.
	// +optional
	ExcludeSchema []string `yaml:"excludeSchema" json:"excludeSchema,omitempty"`
	// Table is the list of schema-qualified tables (e.g. "dbo.Users") to include.
	// +optional
	Table []string `yaml:"table" json:"table,omitempty"`
	// ExcludeTable is the list of schema-qualified tables to exclude.
	// +optional
	ExcludeTable []string `yaml:"excludeTable" json:"excludeTable,omitempty"`
}

type MSSQLServerSnapshot struct {
	// Enabled controls whether the Snapshot Phase should be executed.
	// +optional
	Enabled  bool                         `yaml:"enabled" json:"enabled"`
	Pipeline *MSSQLServerSnapshotPipeline `yaml:"pipeline" json:"pipeline,omitempty"`
}

type MSSQLServerStreaming struct {
	// Enabled controls whether the CDC Streaming Phase should be executed.
	// +optional
	Enabled bool `yaml:"enabled" json:"enabled"`
	// PollInterval controls how often CDC changes are polled from the source.
	// +optional
	PollInterval time.Duration `yaml:"pollInterval" json:"pollInterval,omitempty"`
	// AutoEnableCDC enables Change Data Capture on the source database/tables
	// automatically when set to true. If false, CDC must be pre-enabled.
	// +optional
	AutoEnableCDC bool `yaml:"autoEnableCDC" json:"autoEnableCDC,omitempty"`
	// BatchSize is the maximum number of CDC changes to apply in a single
	// target transaction. Larger batches improve throughput but extend the
	// window of lost work if the migration is cancelled mid-batch.
	// +optional
	BatchSize *int `yaml:"batchSize" json:"batchSize,omitempty"`
}

type MSSQLServerConnectionInfo struct {
	AppBinding             *kmapi.ObjectReference `json:"appBinding,omitempty"`
	Database               string                 `json:"database"`
	MaxConnections         int                    `json:"maxConnections,omitempty"`
	Encrypt                bool                   `json:"encrypt,omitempty"`
	TrustServerCertificate bool                   `json:"trustServerCertificate,omitempty"`

	Address  string `json:"-"`
	User     string `json:"-"`
	Password string `json:"-"`
}

type MSSQLServerSnapshotPipeline struct {
	Workers        *int `yaml:"workers" json:"workers"`
	Sinkers        *int `yaml:"sinkers" json:"sinkers"`
	Buffer         *int `yaml:"buffer" json:"buffer"`
	ReadBatchSize  *int `yaml:"readBatchSize" json:"read_batch_size"`
	WriteBatchSize *int `yaml:"writeBatchSize" json:"write_batch_size"`
}

func init() {
	SchemeBuilder.Register(&MSSQLServerMigration{}, &MSSQLServerMigrationList{})
}
