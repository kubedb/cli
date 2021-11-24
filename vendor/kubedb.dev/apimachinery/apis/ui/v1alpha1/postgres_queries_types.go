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
	Queries []PostgresQuerySpec `json:"queries" protobuf:"bytes,1,rep,name=queries"`
}

type PostgresQuerySpec struct {
	UserOID                  int64   `json:"userOID" protobuf:"varint,1,opt,name=userOID"`
	DatabaseOID              int64   `json:"databaseOID" protobuf:"varint,2,opt,name=databaseOID"`
	Query                    string  `json:"query" protobuf:"bytes,3,opt,name=query"`
	Calls                    int64   `json:"calls" protobuf:"varint,4,opt,name=calls"`
	Rows                     int64   `json:"rows" protobuf:"varint,5,opt,name=rows"`
	TotalTime                float64 `json:"totalTime" protobuf:"fixed64,6,opt,name=totalTime"`
	MinTime                  float64 `json:"minTime" protobuf:"fixed64,7,opt,name=minTime"`
	MaxTime                  float64 `json:"maxTime" protobuf:"fixed64,8,opt,name=maxTime"`
	SharedBlksHit            int64   `json:"sharedBlksHit" protobuf:"varint,9,opt,name=sharedBlksHit"`
	SharedBlksRead           int64   `json:"sharedBlksRead" protobuf:"varint,10,opt,name=sharedBlksRead"`
	SharedBlksDirtied        int64   `json:"sharedBlksDirtied" protobuf:"varint,11,opt,name=sharedBlksDirtied"`
	SharedBlksWritten        int64   `json:"sharedBlksWritten" protobuf:"varint,12,opt,name=sharedBlksWritten"`
	LocalBlksHit             int64   `json:"localBlksHit" protobuf:"varint,13,opt,name=localBlksHit"`
	LocalBlksRead            int64   `json:"localBlksRead" protobuf:"varint,14,opt,name=localBlksRead"`
	LocalBlksDirtied         int64   `json:"localBlksDirtied" protobuf:"varint,15,opt,name=localBlksDirtied"`
	LocalBlksWritten         int64   `json:"localBlksWritten" protobuf:"varint,16,opt,name=localBlksWritten"`
	TempBlksRead             int64   `json:"tempBlksRead" protobuf:"varint,17,opt,name=tempBlksRead"`
	TempBlksWritten          int64   `json:"tempBlksWritten" protobuf:"varint,18,opt,name=tempBlksWritten"`
	BlkReadTime              float64 `json:"blkReadTime" protobuf:"fixed64,19,opt,name=blkReadTime"`
	BlkWriteTime             float64 `json:"blkWriteTime" protobuf:"fixed64,20,opt,name=blkWriteTime"`
	BufferHitPercentage      string  `json:"bufferHitPercentage" protobuf:"bytes,21,opt,name=bufferHitPercentage"`
	LocalBufferHitPercentage string  `json:"localBufferHitPercentage" protobuf:"bytes,22,opt,name=localBufferHitPercentage"`
}

// PostgresQueries is the Schema for the PostgresQueries API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresQueries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec PostgresQueriesSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// PostgresQueriesList contains a list of PostgresQueries

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresQueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []PostgresQueries `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&PostgresQueries{}, &PostgresQueriesList{})
}
