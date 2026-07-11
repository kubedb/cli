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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindMySQLMigration     = "MySQLMigration"
	ResourceSingularMySQLMigration = "mysqlmigration"
	ResourcePluralMySQLMigrations  = "mysqlmigrations"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mysqlmigrations,singular=mysqlmigration,shortName=mymig,categories={kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Stage",type="string",JSONPath=".status.progress.info.Stage"
// +kubebuilder:printcolumn:name="Lag",type="string",JSONPath=".status.progress.info.Lag"
// +kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress.info.Progress"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MySQLMigration struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of MySQLMigration
	// +required
	Spec MySQLMigrationSpec `json:"spec"`

	// status defines the observed state of MySQLMigration.
	// It reuses the shared MigrationStatus so that the Migration duck type can
	// project it and the operator's status patches replay onto it unchanged.
	// +optional
	Status MigrationStatus `json:"status,omitzero"`
}

// MySQLMigrationSpec defines the desired state of MySQLMigration
type MySQLMigrationSpec struct {
	// Source defines the source MySQL database configuration
	Source MySQLSource `json:"source"`

	// Target defines the target MySQL database configuration
	Target MySQLTarget `json:"target"`

	// JobDefaults specifies default settings for migration jobs
	// +optional
	JobDefaults *JobDefaults `json:"jobDefaults,omitempty"`

	// JobTemplate specifies runtime configurations for the migration Job
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
}

// MySQLMigrationList contains a list of MySQLMigration

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MySQLMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []MySQLMigration `json:"items"`
}

type MySQLSource struct {
	// ConnectionInfo refers to the source MySQL database connection information.
	ConnectionInfo *ConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
	Schema         *MySQLSchema    `yaml:"schema" json:"schema,omitempty"`
	Snapshot       *MySQLSnapshot  `yaml:"snapshot" json:"snapshot,omitempty"`
	Streaming      *MySQLStreaming `yaml:"streaming" json:"streaming,omitempty"`
}

type MySQLTarget struct {
	// ConnectionInfo refers to the target MySQL database connection information.
	ConnectionInfo *ConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
}

type MySQLSchema struct {
	// Enabled controls whether the Schema Phase should be executed.
	Enabled bool `yaml:"enabled" json:"enabled"`
	// Database is the list of databases to migrate.
	// +optional
	Database []string `yaml:"database" json:"database,omitempty"`
	// ExcludeDatabase is the list of databases to exclude from migration.
	// +optional
	ExcludeDatabase []string `yaml:"excludeDatabase" json:"excludeDatabase,omitempty"`
}

type MySQLSnapshot struct {
	// Enabled controls whether the Snapshot Phase should be executed.
	// +optional
	Enabled  bool                   `yaml:"enabled" json:"enabled"`
	Pipeline *MySQLSnapshotPipeline `yaml:"pipeline" json:"pipeline,omitempty"`
}

type MySQLStreaming struct {
	// Enabled controls whether the Logical Replication Phase should be executed.
	// +optional
	Enabled bool `yaml:"enabled" json:"enabled"`
}

type MySQLConnectionInfo struct {
	// AppBinding refers to the source database AppBinding name, which contains the connection information.
	// +optional
	AppBinding     *kmapi.ObjectReference `yaml:"appBinding,omitempty" json:"appBinding,omitempty"`
	Address        string                 `yaml:"address" json:"address"`
	User           string                 `yaml:"user" json:"user"`
	Password       string                 `yaml:"password" json:"password"`
	DBName         string                 `yaml:"dbName" json:"dbName"`
	MaxConnections int                    `yaml:"maxConnections" json:"maxConnections,omitempty"`
}

type MySQLSnapshotPipeline struct {
	Workers        *int `yaml:"workers" json:"workers"`
	Sinkers        *int `yaml:"sinkers" json:"sinkers"`
	Buffer         *int `yaml:"buffer" json:"buffer"`
	ReadBatchSize  *int `yaml:"readBatchSize" json:"read_batch_size"`
	WriteBatchSize *int `yaml:"writeBatchSize" json:"write_batch_size"`
}

func init() {
	SchemeBuilder.Register(&MySQLMigration{}, &MySQLMigrationList{})
}
