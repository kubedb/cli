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
	ResourceKindPostgresMigration     = "PostgresMigration"
	ResourceSingularPostgresMigration = "postgresmigration"
	ResourcePluralPostgresMigrations  = "postgresmigrations"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=postgresmigrations,singular=postgresmigration,shortName=pgmig,categories={kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Stage",type="string",JSONPath=".status.progress.info.Stage"
// +kubebuilder:printcolumn:name="Lag",type="string",JSONPath=".status.progress.info.Lag"
// +kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress.info.Progress"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PostgresMigration struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of PostgresMigration
	// +required
	Spec PostgresMigrationSpec `json:"spec"`

	// status defines the observed state of PostgresMigration.
	// It reuses the shared MigrationStatus so that the Migration duck type can
	// project it and the operator's status patches replay onto it unchanged.
	// +optional
	Status MigrationStatus `json:"status,omitzero"`
}

// PostgresMigrationSpec defines the desired state of PostgresMigration
type PostgresMigrationSpec struct {
	// Source defines the source Postgres database configuration
	Source PostgresSource `json:"source"`

	// Target defines the target Postgres database configuration
	Target PostgresTarget `json:"target"`

	// JobDefaults specifies default settings for migration jobs
	// +optional
	JobDefaults *JobDefaults `json:"jobDefaults,omitempty"`

	// JobTemplate specifies runtime configurations for the migration Job
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
}

// PostgresMigrationList contains a list of PostgresMigration

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []PostgresMigration `json:"items"`
}

type PostgresSource struct {
	// ConnectionInfo refers to the source Postgres database connection information.
	ConnectionInfo ConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`

	// PgDump refers to the CLI name which will be used to dump the schema or data from the source Postgres database.
	PgDump *PgDump `yaml:"pgDump" json:"pgDump,omitempty"`

	// LogicalReplication refers to the logical replication configuration. URL: https://www.postgresql.org/docs/current/logical-replication.html
	LogicalReplication *LogicalReplication `yaml:"logicalReplication" json:"logicalReplication,omitempty"`
}

type PostgresTarget struct {
	// ConnectionInfo refers to the target Postgres database connection information.
	ConnectionInfo ConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
}

type LogicalReplication struct {
	// CopyData refers to whether to copy data the initial snapshot when creating the subscription.
	// +kubebuilder:default=true
	// +optional
	CopyData bool `yaml:"copyData" json:"copyData,omitempty"`

	// Publication refers to the publication configuration.
	Publication *Publication `yaml:"publication" json:"publication,omitempty"`

	// Subscription refers to the subscription configuration.
	Subscription *Subscription `yaml:"subscription" json:"subscription,omitempty"`
}

type PgDump struct {
	// SchemaOnly indicates dump only the schema, no data
	// Equivalent to: pg_dump --schema-only
	// +optional
	SchemaOnly bool `yaml:"schemaOnly" json:"schemaOnly,omitempty"`

	// Schema specifies dump the specified schema(s) only
	// Equivalent to: pg_dump --schema=<schema>
	// +optional
	Schema []string `yaml:"schema" json:"schema,omitempty"`

	// ExcludeSchema specifies PATTERN do NOT dump the specified schema(s)
	// Equivalent to: pg_dump --exclude-schema=<schema>
	// +optional
	ExcludeSchema []string `yaml:"excludeSchema" json:"excludeSchema,omitempty"`

	// Table specifies dump only the specified table(s)
	// Equivalent to: pg_dump --table=<table>
	// +optional
	Table []string `yaml:"table" json:"table,omitempty"`

	// ExcludeTable specifies do NOT dump the specified table(s)
	// Equivalent to: pg_dump --exclude-table=<table>
	// +optional
	ExcludeTable []string `yaml:"excludeTable" json:"excludeTable,omitempty"`

	// ExtraOptions contains additional raw pg_dump command-line flags
	// that are not explicitly modeled by the CRD fields.
	// +optional
	ExtraOptions []string `yaml:"extraOptions" json:"extraOptions,omitempty"`
}

type Publication struct {
	// Name is the identifier of the PostgreSQL publication.
	// This name will be used when creating or referencing the publication in logical replication.
	Name string `yaml:"name" json:"name,omitempty"`

	// Mode defines how tables are selected for the publication.
	//
	// Supported values:
	//   - default: Applies filtering behavior similar to pg_dump (manual selection).
	//   - table: Publishes only the specified tables (FOR TABLE ...).
	//   - allTable: Publishes all tables in the database (FOR ALL TABLES).
	//   - tableInSchema: Publishes all tables within specified schemas (FOR TABLES IN SCHEMA ...).
	// +kubebuilder:validation:Enum=default;table;allTable;tableInSchema
	// +kubebuilder:default=default
	// +optional
	Mode string `yaml:"mode" json:"mode,omitempty"`

	// Args contains additional publication parameters,
	// such as table names or schema names depending on the selected Mode.
	//
	// For example:
	//   - Mode=table -> Args may include table names
	//   - Mode=tableInSchema -> Args may include schema names
	//
	// +optional
	Args []string `yaml:"args" json:"args,omitempty"`
}

type Subscription struct {
	// Name is the identifier of the PostgreSQL subscription.
	// This name will be used when creating or referencing the subscription in logical replication.
	Name string `yaml:"name" json:"name,omitempty"`
}

func init() {
	SchemeBuilder.Register(&PostgresMigration{}, &PostgresMigrationList{})
}
