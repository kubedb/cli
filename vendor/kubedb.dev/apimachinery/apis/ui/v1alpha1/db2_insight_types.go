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
	ResourceKindDB2Insight = "DB2Insight"
	ResourceDB2Insight     = "db2insight"
	ResourceDB2Insights    = "db2insights"
)

// DB2InsightSpec defines the desired state of DB2Insight.
//
// The IBM DB2 Go driver is cgo and cannot be linked into the ui-server image, so every
// live value here is read over HTTP from the db2-coordinator sidecar rather than from a
// direct SQL connection.
//
// DB2Spec has no TLS, topology, monitor or disableSecurity field, so this reports no
// TLS/mode information: there is none to report.
type DB2InsightSpec struct {
	Version string                 `json:"version"`
	Status  dboldapi.DatabasePhase `json:"status"`

	// ServiceLevel is the DB2 service level string, e.g. "DB2 v11.5.8.0", from
	// SYSIBMADM.ENV_INST_INFO.
	// +optional
	ServiceLevel string `json:"serviceLevel,omitempty"`
	// InstanceName is the DB2 instance owner, normally db2inst1.
	// +optional
	InstanceName string `json:"instanceName,omitempty"`
	// DatabaseName is the database the coordinator is connected to.
	// +optional
	DatabaseName string `json:"databaseName,omitempty"`

	// Reachable is the outcome of the coordinator's own readiness probe against DB2.
	// +optional
	Reachable bool `json:"reachable,omitempty"`

	// Replicas compares spec.replicas against Ready pods.
	// +optional
	Replicas DB2ReplicaStat `json:"replicas,omitempty"`

	// +optional
	Connections DB2ConnectionInfo `json:"connections,omitempty"`

	// BufferPoolHitRatioPercent is the overall buffer pool hit ratio, derived from
	// MON_GET_BUFFERPOOL logical and physical reads.
	// +optional
	BufferPoolHitRatioPercent *float64 `json:"bufferPoolHitRatioPercent,omitempty"`

	// +optional
	TotalTablespaces *int32 `json:"totalTablespaces,omitempty"`
	// +optional
	TotalSchemas *int32 `json:"totalSchemas,omitempty"`
	// +optional
	TotalTables *int32 `json:"totalTables,omitempty"`
}

type DB2ReplicaStat struct {
	Ready int32 `json:"ready"`
	// +optional
	Desired *int32 `json:"desired,omitempty"`
}

type DB2ConnectionInfo struct {
	// +optional
	CurrentConnections *int64 `json:"currentConnections,omitempty"`
	// +optional
	TotalConnections *int64 `json:"totalConnections,omitempty"`
}

// DB2Insight is the Schema for the DB2Insights API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DB2Insight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DB2InsightSpec     `json:"spec,omitempty"`
	Status dboldapi.DB2Status `json:"status,omitempty"`
}

// DB2InsightList contains a list of DB2Insight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DB2InsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DB2Insight `json:"items"`
}
