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
