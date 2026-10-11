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
	ResourceKindMilvusSchemaOverview = "MilvusSchemaOverview"
	ResourceMilvusSchemaOverview     = "milvusschemaoverview"
	ResourceMilvusSchemaOverviews    = "milvusschemaoverviews"
)

// MilvusSchemaOverviewSpec defines the desired state of MilvusSchemaOverview.
//
// Milvus organises data as databases -> collections, so this does not reuse
// GenericSchemaOverviewSpec: a collection carries load state, shard count, partitions and
// indexes, none of which fit the generic {databaseName, tableName, tableSizeBytes} triple.
type MilvusSchemaOverviewSpec struct {
	Collections []MilvusCollectionSpec `json:"collections"`
}

// MilvusCollectionSpec describes one Milvus collection.
type MilvusCollectionSpec struct {
	DatabaseName string `json:"databaseName"`
	Name         string `json:"name"`

	// RowCount is the number of entities in the collection, from GetCollectionStats.
	// +optional
	RowCount *int64 `json:"rowCount,omitempty"`

	// StorageSizeBytes is the size of the collection's flushed binlogs in object storage.
	// It excludes growing segments that have not been flushed yet, so it lags recent
	// writes. Unset when the server does not report it.
	// +optional
	StorageSizeBytes *int64 `json:"storageSizeBytes,omitempty"`

	// Loaded reports whether the collection is loaded into memory. Only loaded collections
	// can be searched.
	Loaded bool `json:"loaded"`

	// LoadProgress is the percentage of the collection loaded into memory. It is reported
	// while a load is in progress or has not started, and omitted once the collection is
	// fully loaded.
	// +optional
	LoadProgress *int64 `json:"loadProgress,omitempty"`

	// +optional
	ShardNum int32 `json:"shardNum,omitempty"`

	// +optional
	ConsistencyLevel string `json:"consistencyLevel,omitempty"`

	// +optional
	Partitions []string `json:"partitions,omitempty"`

	// +optional
	Indexes []string `json:"indexes,omitempty"`
}

// MilvusSchemaOverview is the Schema for the MilvusSchemaOverviews API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MilvusSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MilvusSchemaOverviewSpec `json:"spec,omitempty"`
}

// MilvusSchemaOverviewList contains a list of MilvusSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MilvusSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MilvusSchemaOverview `json:"items"`
}
