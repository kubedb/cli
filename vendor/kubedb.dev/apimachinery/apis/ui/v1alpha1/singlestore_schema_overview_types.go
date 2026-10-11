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
	ResourceKindSinglestoreSchemaOverview = "SinglestoreSchemaOverview"
	ResourceSinglestoreSchemaOverview     = "singlestoreschemaoverview"
	ResourceSinglestoreSchemaOverviews    = "singlestoreschemaoverviews"
)

// SinglestoreSchemaOverviewSpec defines the desired state of SinglestoreSchemaOverview.
type SinglestoreSchemaOverviewSpec struct {
	Databases []SinglestoreDatabaseSpec `json:"databases"`
}

// SinglestoreDatabaseSpec contains logical and physical storage statistics for a table.
type SinglestoreDatabaseSpec struct {
	DatabaseName string `json:"databaseName"`
	TableName    string `json:"tableName"`
	// StorageType is the deterministic, comma-separated set reported by TABLE_STATISTICS.
	StorageType string `json:"storageType,omitempty"`
	// Rows sums TABLE_STATISTICS master partitions so replicas are not double-counted.
	Rows *int64 `json:"rows,omitempty"`
	// MemoryUseBytes sums TABLE_STATISTICS master partitions so replicas are not double-counted.
	MemoryUseBytes *int64 `json:"memoryUseBytes,omitempty"`
	// DiskUseBytes is the cluster-wide physical columnstore file size, which can include redundant copies.
	DiskUseBytes *int64 `json:"diskUseBytes,omitempty"`
}

// SinglestoreSchemaOverview is the Schema for the singlestoreschemaoverviews API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SinglestoreSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SinglestoreSchemaOverviewSpec `json:"spec,omitempty"`
}

// SinglestoreSchemaOverviewList contains a list of SinglestoreSchemaOverview.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SinglestoreSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SinglestoreSchemaOverview `json:"items"`
}
