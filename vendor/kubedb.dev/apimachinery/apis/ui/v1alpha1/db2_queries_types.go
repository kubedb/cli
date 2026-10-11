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
	ResourceKindDB2Queries = "DB2Queries"
	ResourceDB2Queries     = "db2queries"
	// ResourceDB2Querieses mirrors the MongoDB/Postgres quirk: the Kind already ends in
	// "s", so the plural constant is the doubled-es name but the value is the same string.
	ResourceDB2Querieses = "db2queries"
)

// DB2QueriesSpec defines the desired state of DB2Queries.
//
// Rows come from MON_GET_PKG_CACHE_STMT, DB2's package cache of recently executed
// statements. It is a cache, not a permanent digest: entries age out, so an empty list on a
// quiet database is a legitimate result.
type DB2QueriesSpec struct {
	Queries []DB2QuerySpec `json:"queries"`
}

type DB2QuerySpec struct {
	// StatementText is STMT_TEXT.
	StatementText string `json:"statementText"`
	// +optional
	NumExecutions *int64 `json:"numExecutions,omitempty"`
	// +optional
	TotalCPUTimeMicroSeconds *int64 `json:"totalCPUTimeMicroSeconds,omitempty"`
	// +optional
	TotalActTimeMilliSeconds *int64 `json:"totalActTimeMilliSeconds,omitempty"`
	// +optional
	AvgExecutionTimeMilliSeconds *int64 `json:"avgExecutionTimeMilliSeconds,omitempty"`
	// +optional
	RowsRead *int64 `json:"rowsRead,omitempty"`
	// +optional
	RowsReturned *int64 `json:"rowsReturned,omitempty"`
}

// DB2Queries is the Schema for the DB2Queries API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DB2Queries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DB2QueriesSpec `json:"spec,omitempty"`
}

// DB2QueriesList contains a list of DB2Queries

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DB2QueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DB2Queries `json:"items"`
}
