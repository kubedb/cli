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
	ResourceKindPostgresQueries = "PostgresQueries"
	ResourcePostgresQueries     = "postgresqueries"
	ResourcePostgresQuerieses   = "postgresqueries"
)

// PostgresQueriesSpec defines the desired state of PostgresQueries
type PostgresQueriesSpec struct {
	Queries []PostgresQuerySpec `json:"queries"`
}

type PostgresQuerySpec struct {
	UserOID                  int64    `json:"userOID"`
	DatabaseOID              int64    `json:"databaseOID"`
	Query                    string   `json:"query"`
	Calls                    *int64   `json:"calls,omitempty"`
	Rows                     *int64   `json:"rows,omitempty"`
	TotalTimeMilliSeconds    *float64 `json:"totalTimeMilliSeconds,omitempty"`
	MinTimeMilliSeconds      *float64 `json:"minTimeMilliSeconds,omitempty"`
	MaxTimeMilliSeconds      *float64 `json:"maxTimeMilliSeconds,omitempty"`
	SharedBlksHit            *int64   `json:"sharedBlksHit,omitempty"`
	SharedBlksRead           *int64   `json:"sharedBlksRead,omitempty"`
	SharedBlksDirtied        *int64   `json:"sharedBlksDirtied,omitempty"`
	SharedBlksWritten        *int64   `json:"sharedBlksWritten,omitempty"`
	LocalBlksHit             *int64   `json:"localBlksHit,omitempty"`
	LocalBlksRead            *int64   `json:"localBlksRead,omitempty"`
	LocalBlksDirtied         *int64   `json:"localBlksDirtied,omitempty"`
	LocalBlksWritten         *int64   `json:"localBlksWritten,omitempty"`
	TempBlksRead             *int64   `json:"tempBlksRead,omitempty"`
	TempBlksWritten          *int64   `json:"tempBlksWritten,omitempty"`
	BlkReadTimeMilliSeconds  *float64 `json:"blkReadTimeMilliSeconds,omitempty"`
	BlkWriteTime             *float64 `json:"blkWriteTime,omitempty"`
	BufferHitPercentage      *float64 `json:"bufferHitPercentage,omitempty"`
	LocalBufferHitPercentage *float64 `json:"localBufferHitPercentage,omitempty"`
}

// PostgresQueries is the Schema for the PostgresQueries API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresQueries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PostgresQueriesSpec `json:"spec,omitempty"`
}

// PostgresQueriesList contains a list of PostgresQueries

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresQueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresQueries `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PostgresQueries{}, &PostgresQueriesList{})
}
