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
	ResourceKindCassandraInsight = "CassandraInsight"
	ResourceCassandraInsight     = "cassandrainsight"
	ResourceCassandraInsights    = "cassandrainsights"
)

// CassandraMode is the deployment shape of a Cassandra object.
// +kubebuilder:validation:Enum=Standalone;Topology
type CassandraMode string

const (
	CassandraModeStandalone CassandraMode = "Standalone"
	CassandraModeTopology   CassandraMode = "Topology"
)

// CassandraInsightSpec defines the desired state of CassandraInsight
type CassandraInsightSpec struct {
	Version string                 `json:"version"`
	Type    CassandraMode          `json:"type"`
	Status  dboldapi.DatabasePhase `json:"status"`

	// ClusterName is the cluster name the nodes agree on, from system.local.
	// +optional
	ClusterName string `json:"clusterName,omitempty"`
	// Partitioner in use, e.g. org.apache.cassandra.dht.Murmur3Partitioner.
	// +optional
	Partitioner string `json:"partitioner,omitempty"`
	// CQLVersion reported by the contacted node.
	// +optional
	CQLVersion string `json:"cqlVersion,omitempty"`
	// NativeProtocolVersion reported by the contacted node.
	// +optional
	NativeProtocolVersion string `json:"nativeProtocolVersion,omitempty"`

	// SchemaAgreement is true when every reachable node reports the same
	// schema_version. A false value means a schema change is still propagating.
	// +optional
	SchemaAgreement bool `json:"schemaAgreement,omitempty"`

	// Nodes is one entry per declared replica, each read from that node's own
	// system.local. An unreachable node is absent rather than failing the resource.
	// +optional
	Nodes []CassandraNodeStat `json:"nodes,omitempty"`

	// Racks reflects spec.topology.rack: how many replicas each rack asks for
	// against how many pods are Ready.
	// +optional
	Racks []CassandraRackStat `json:"racks,omitempty"`

	// +optional
	TotalNodes *int32 `json:"totalNodes,omitempty"`
	// +optional
	TotalKeyspaces *int32 `json:"totalKeyspaces,omitempty"`
	// +optional
	TotalTables *int32 `json:"totalTables,omitempty"`
}

// CassandraNodeStat describes one node of the ring.
type CassandraNodeStat struct {
	// Address is the node's rpc_address, read from that node's own system.local.
	Address string `json:"address"`
	// +optional
	DataCenter string `json:"dataCenter,omitempty"`
	// +optional
	Rack string `json:"rack,omitempty"`
	// +optional
	ReleaseVersion string `json:"releaseVersion,omitempty"`
	// +optional
	HostID string `json:"hostID,omitempty"`
	// +optional
	SchemaVersion string `json:"schemaVersion,omitempty"`
}

// CassandraRackStat compares a rack's requested replicas against Ready pods.
type CassandraRackStat struct {
	Name string `json:"name"`
	// Ready is the number of Ready pods in the rack's PetSet.
	Ready int32 `json:"ready"`
	// Desired is spec.topology.rack[].replicas. Unset for a Standalone object,
	// which declares no racks.
	// +optional
	Desired *int32 `json:"desired,omitempty"`
}

// CassandraInsight is the Schema for the CassandraInsights API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CassandraInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CassandraInsightSpec     `json:"spec,omitempty"`
	Status dboldapi.CassandraStatus `json:"status,omitempty"`
}

// CassandraInsightList contains a list of CassandraInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CassandraInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CassandraInsight `json:"items"`
}
