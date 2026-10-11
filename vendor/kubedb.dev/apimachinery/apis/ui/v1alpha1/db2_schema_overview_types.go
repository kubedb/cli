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
)

const (
	ResourceKindDB2SchemaOverview = "DB2SchemaOverview"
	ResourceDB2SchemaOverview     = "db2schemaoverview"
	ResourceDB2SchemaOverviews    = "db2schemaoverviews"
)

// DB2SchemaOverviewSpec defines the desired state of DB2SchemaOverview.
//
// DB2's unit is schema.table, which maps onto GenericSchemaOverviewSpec's
// databaseName/tableName pair, but DB2 also reports separate data and index sizes that the
// generic triple cannot hold, so this defines its own shape.
type DB2SchemaOverviewSpec struct {
	Tables []DB2TableSpec `json:"tables"`
}

// DB2TableSpec describes one table. System schemas (SYS*) are excluded.
type DB2TableSpec struct {
	// SchemaName is TABSCHEMA. DB2 upper-cases identifiers by default, so expect
	// upper-case values here.
	SchemaName string `json:"schemaName"`
	// TableName is TABNAME.
	TableName string `json:"tableName"`
	// Type is the TYPE column: T for table, V for view, and so on.
	// +optional
	Type string `json:"type,omitempty"`
	// +optional
	RowCount *int64 `json:"rowCount,omitempty"`
	// DataSizeBytes is DATA_OBJECT_L_SIZE converted from KB.
	// +optional
	DataSizeBytes *int64 `json:"dataSizeBytes,omitempty"`
	// IndexSizeBytes is INDEX_OBJECT_L_SIZE converted from KB.
	// +optional
	IndexSizeBytes *int64 `json:"indexSizeBytes,omitempty"`
	// +optional
	ColumnCount *int32 `json:"columnCount,omitempty"`
}

// DB2SchemaOverview is the Schema for the DB2SchemaOverviews API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DB2SchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DB2SchemaOverviewSpec `json:"spec,omitempty"`
}

// DB2SchemaOverviewList contains a list of DB2SchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DB2SchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DB2SchemaOverview `json:"items"`
}
