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
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindMongoDBMigration     = "MongoDBMigration"
	ResourceSingularMongoDBMigration = "mongodbmigration"
	ResourcePluralMongoDBMigrations  = "mongodbmigrations"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mongodbmigrations,singular=mongodbmigration,shortName=momig,categories={kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Stage",type="string",JSONPath=".status.progress.info.Stage"
// +kubebuilder:printcolumn:name="Lag",type="string",JSONPath=".status.progress.info.Lag"
// +kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress.info.Progress"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MongoDBMigration struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of MongoDBMigration
	// +required
	Spec MongoDBMigrationSpec `json:"spec"`

	// status defines the observed state of MongoDBMigration.
	// It reuses the shared MigrationStatus so that the Migration duck type can
	// project it and the operator's status patches replay onto it unchanged.
	// +optional
	Status MigrationStatus `json:"status,omitzero"`
}

// MongoDBMigrationSpec defines the desired state of MongoDBMigration
type MongoDBMigrationSpec struct {
	// Source defines the source MongoDB database configuration
	Source MongoDBSource `json:"source"`

	// Target defines the target MongoDB database configuration
	Target MongoDBTarget `json:"target"`

	// JobDefaults specifies default settings for migration jobs
	// +optional
	JobDefaults *JobDefaults `json:"jobDefaults,omitempty"`

	// JobTemplate specifies runtime configurations for the migration Job
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
}

// MongoDBMigrationList contains a list of MongoDBMigration

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []MongoDBMigration `json:"items"`
}

type MongoDBSource struct {
	ConnectionInfo ConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
	Mongoshake     *Mongoshake    `yaml:"mongoshake" json:"mongoshake,omitempty"`
}

type MongoDBTarget struct {
	ConnectionInfo ConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
}

type Mongoshake struct {
	// SyncMode controls synchronization mode.
	// Supported values: all, full, incr
	// - all  : full synchronization + incremental synchronization
	// - full : full synchronization only
	// - incr : incremental synchronization only
	// Equivalent behavior: defines replication mode of mongoshake
	SyncMode string `yaml:"syncMode" json:"syncMode,omitempty" config:"sync_mode"`
	// FilterOpTypes filters MongoDB oplog operation types to include.
	// Example values: "i" (insert), "u" (update), "d" (delete)
	// Equivalent behavior: acts as oplog.op type filter
	// +optional
	FilterOpTypes []string `yaml:"filterOpTypes" json:"filterOpTypes,omitempty" config:"filter.op_types"`
	// FilterNamespaceBlack excludes specified namespaces (db.collection or db).
	// Format: db.collection or db
	// Multiple entries separated by ';'
	// Example: db1.col1;db2
	// If set, listed namespaces will be filtered out.
	// +optional
	FilterNamespaceBlack []string `yaml:"filterNamespaceBlack" json:"filterNamespaceBlack,omitempty" config:"filter.namespace.black"`
	// FilterNamespaceWhite includes only specified namespaces (db.collection or db).
	// Format: db.collection or db
	// Multiple entries separated by ';'
	// Example: db1.col1;db2
	// If set, only listed namespaces will be allowed.
	// +optional
	FilterNamespaceWhite []string `yaml:"filterNamespaceWhite" json:"filterNamespaceWhite,omitempty" config:"filter.namespace.white"`
	// FilterPassSpecialDb allows special system databases to be included.
	// Example: admin;local;config;mongoshake
	// Note: collection-level filtering is not supported here.
	// +optional
	FilterPassSpecialDb []string `yaml:"filterPassSpecialDb" json:"filterPassSpecialDb,omitempty" config:"filter.pass.special.db"`
	// FilterDDLEnable controls whether DDL operations are filtered or passed.
	// When enabled, only oplog operations (i/u/d) are synced.
	// When disabled, DDL operations like create index or drop database are included.
	// +optional
	FilterDDLEnable *bool `yaml:"filterDdlEnable" json:"filterDdlEnable,omitempty" config:"filter.ddl_enable"`
	// FilterOplogGids enables filtering of oplog by GID.
	// +optional
	FilterOplogGids *bool `yaml:"filterOplogGids" json:"filterOplogGids,omitempty" config:"filter.oplog.gids"`
	// CheckpointStartPosition defines initial oplog position (UTC timestamp).
	// Used only when no checkpoint exists.
	// Note: UTC time is 8 hours ahead of CST.
	// +optional
	CheckpointStartPosition int64 `yaml:"checkpointStartPosition" json:"checkpointStartPosition,omitempty" config:"checkpoint.start_position" type:"date"`
	// TransformNamespace maps source namespaces to destination namespaces.
	// Format:
	//   fromDb.fromCollection:toDb.toCollection
	//   fromDb:toDb
	// Multiple mappings separated by ';'
	// Example: db1.col1:db2.col1;db3:db4
	// +optional
	TransformNamespace []string `yaml:"transformNamespace" json:"transformNamespace,omitempty" config:"transform.namespace"`
	// ExtraConfiguration allows additional raw mongoshake configuration.
	// Key-value pairs passed directly without schema validation.
	// +optional
	ExtraConfiguration map[string]string `yaml:"extraConfiguration" json:"extraConfiguration,omitempty"`
}

func init() {
	SchemeBuilder.Register(&MongoDBMigration{}, &MongoDBMigrationList{})
}
