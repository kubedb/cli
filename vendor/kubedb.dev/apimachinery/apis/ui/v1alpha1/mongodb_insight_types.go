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
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindMongoDBInsight = "MongoDBInsight"
	ResourceMongoDBInsight     = "mongodbinsight"
	ResourceMongoDBInsights    = "mongodbinsights"
)

// MongoDBInsightSpec defines the desired state of MongoDBInsight
type MongoDBInsightSpec struct {
	Version        string                  `json:"version"`
	Type           MongoDBMode             `json:"type"`
	Status         api.DatabasePhase       `json:"status"`
	Connections    *MongoDBConnectionsInfo `json:"connections,omitempty"`
	DBStats        *MongoDBDatabaseStats   `json:"dbStats,omitempty"`
	ShardsInfo     *MongoDBShardsInfo      `json:"shardsInfo,omitempty"`
	ReplicaSetInfo *MongoDBReplicaSetInfo  `json:"replicaSetInfo,omitempty"`
}

type MongoDBDatabaseStats struct {
	TotalCollections int32 `json:"totalCollections"`
	DataSize         int64 `json:"dataSize"`
	TotalIndexes     int32 `json:"totalIndexes"`
	IndexSize        int64 `json:"indexSize"`
}

type MongoDBConnectionsInfo struct {
	CurrentConnections   int32 `json:"currentConnections"`
	TotalConnections     int32 `json:"totalConnections"`
	AvailableConnections int32 `json:"availableConnections"`
	ActiveConnections    int32 `json:"activeConnections"`
}

type MongoDBReplicaSetInfo struct {
	NumberOfReplicas int32 `json:"numberOfReplicas"`
}

type MongoDBShardsInfo struct {
	NumberOfShards    int32 `json:"numberOfShards"`
	ReplicasPerShards int32 `json:"replicasPerShards"`
	NumberOfChunks    int32 `json:"numberOfChunks"`
	BalancerEnabled   bool  `json:"balancerEnabled,omitempty"`
	ChunksBalanced    bool  `json:"chunksBalanced,omitempty"`
}

// MongoDBInsight is the Schema for the MongoDBInsights API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoDBInsightSpec `json:"spec,omitempty"`
	Status api.MongoDBStatus  `json:"status,omitempty"`
}

// MongoDBInsightList contains a list of MongoDBInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDBInsight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBInsight{}, &MongoDBInsightList{})
}
