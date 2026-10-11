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
	dboldapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindMilvusInsight = "MilvusInsight"
	ResourceMilvusInsight     = "milvusinsight"
	ResourceMilvusInsights    = "milvusinsights"
)

// MilvusInsightSpec defines the desired state of MilvusInsight
type MilvusInsightSpec struct {
	Version string                 `json:"version"`
	Type    dboldapi.MilvusMode    `json:"type"`
	Status  dboldapi.DatabasePhase `json:"status"`

	// Healthy reports the outcome of the CheckHealth RPC against the Milvus proxy.
	Healthy bool `json:"healthy"`

	// UnhealthyReasons carries the reasons returned by CheckHealth when Healthy is false.
	// +optional
	UnhealthyReasons []string `json:"unhealthyReasons,omitempty"`

	// DeployMode is reported by the running server, e.g. STANDALONE or DISTRIBUTED.
	// It is read from the system_info metrics and may disagree with spec.topology.mode
	// while a mode change is in flight.
	// +optional
	DeployMode string `json:"deployMode,omitempty"`

	// Nodes reports, per Milvus component role, how many members the server currently
	// sees against how many the Milvus object asks for.
	// +optional
	Nodes []MilvusNodeStat `json:"nodes,omitempty"`

	// +optional
	TotalDatabases *int32 `json:"totalDatabases,omitempty"`
	// +optional
	TotalCollections *int32 `json:"totalCollections,omitempty"`
	// LoadedCollections counts the collections currently loaded into memory. Only loaded
	// collections are searchable.
	// +optional
	LoadedCollections *int32 `json:"loadedCollections,omitempty"`
	// +optional
	TotalRows *int64 `json:"totalRows,omitempty"`
	// +optional
	TotalResourceGroups *int32 `json:"totalResourceGroups,omitempty"`
}

// MilvusNodeStat reports the observed and desired member count of one Milvus component role.
type MilvusNodeStat struct {
	// Type is the Milvus node role: proxy, mixcoord, querynode, datanode or streamingnode.
	Type string `json:"type"`

	// Ready is the number of members of this role the server reports as alive.
	Ready int32 `json:"ready"`

	// Desired is read from the Milvus object. It is unset in Standalone mode, where the
	// object declares no per-role replica counts.
	// +optional
	Desired *int32 `json:"desired,omitempty"`

	// HasError is true when at least one member of this role reported an error while
	// collecting its metrics.
	// +optional
	HasError bool `json:"hasError,omitempty"`
}

// MilvusInsight is the Schema for the MilvusInsights API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MilvusInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MilvusInsightSpec     `json:"spec,omitempty"`
	Status dboldapi.MilvusStatus `json:"status,omitempty"`
}

// MilvusInsightList contains a list of MilvusInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MilvusInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MilvusInsight `json:"items"`
}
