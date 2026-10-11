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
	ResourceKindDocumentDBInsight = "DocumentDBInsight"
	ResourceDocumentDBInsight     = "documentdbinsight"
	ResourceDocumentDBInsights    = "documentdbinsights"
)

// DocumentDBMode is the deployment shape of a DocumentDB object.
//
// DocumentDB has no spec.topology discriminator: the operator decides HA purely on
// spec.replicas > 1 (pkg/controllers/health.go), and this mirrors that.
// +kubebuilder:validation:Enum=Standalone;Cluster
type DocumentDBMode string

const (
	DocumentDBModeStandalone DocumentDBMode = "Standalone"
	DocumentDBModeCluster    DocumentDBMode = "Cluster"
)

// DocumentDBInsightSpec defines the desired state of DocumentDBInsight.
//
// DocumentDB is PostgreSQL with the pg_documentdb extensions behind a gateway that speaks
// the MongoDB wire protocol. Everything here is read over the PostgreSQL side (port 9712)
// as the admin user, which is the path the operator itself uses and therefore the proven one.
type DocumentDBInsightSpec struct {
	Version string                 `json:"version"`
	Type    DocumentDBMode         `json:"type"`
	Status  dboldapi.DatabasePhase `json:"status"`

	// ServerVersion is the PostgreSQL backend version reported by version().
	// spec.version is the DocumentDB catalog version, which is a different thing.
	// +optional
	ServerVersion string `json:"serverVersion,omitempty"`

	// InRecovery is pg_is_in_recovery(). The primary Service pins to the primary pod, so
	// this is expected to be false; true means the connection landed on a standby.
	// +optional
	InRecovery bool `json:"inRecovery,omitempty"`

	// +optional
	ConnectionInfo DocumentDBConnectionInfo `json:"connectionInfo,omitempty"`

	// ReplicationStatus is one entry per connected standby, from pg_stat_replication.
	// Empty on a standalone object.
	// +optional
	ReplicationStatus []DocumentDBReplicationStatus `json:"replicationStatus,omitempty"`

	// Replicas compares spec.replicas against Ready pods.
	// +optional
	Replicas DocumentDBReplicaStat `json:"replicas,omitempty"`

	// +optional
	TotalDatabases *int32 `json:"totalDatabases,omitempty"`
	// +optional
	TotalCollections *int32 `json:"totalCollections,omitempty"`
}

type DocumentDBConnectionInfo struct {
	// +optional
	MaxConnections *int64 `json:"maxConnections,omitempty"`
	// +optional
	ActiveConnections *int64 `json:"activeConnections,omitempty"`
	// +optional
	IdleConnections *int64 `json:"idleConnections,omitempty"`
}

type DocumentDBReplicationStatus struct {
	// ApplicationName is the standby's pod name: the operator sets
	// application_name=$HOSTNAME in primary_conninfo.
	ApplicationName string `json:"applicationName"`
	// +optional
	State string `json:"state,omitempty"`
	// +optional
	SyncState string `json:"syncState,omitempty"`
	// +optional
	ClientAddr string `json:"clientAddr,omitempty"`
	// WriteLagBytes is the primary's WAL position minus the standby's write position.
	// +optional
	WriteLagBytes *int64 `json:"writeLagBytes,omitempty"`
	// +optional
	ReplayLagBytes *int64 `json:"replayLagBytes,omitempty"`
}

type DocumentDBReplicaStat struct {
	Ready int32 `json:"ready"`
	// +optional
	Desired *int32 `json:"desired,omitempty"`
}

// DocumentDBInsight is the Schema for the DocumentDBInsights API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DocumentDBInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DocumentDBInsightSpec     `json:"spec,omitempty"`
	Status dboldapi.DocumentDBStatus `json:"status,omitempty"`
}

// DocumentDBInsightList contains a list of DocumentDBInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DocumentDBInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DocumentDBInsight `json:"items"`
}
