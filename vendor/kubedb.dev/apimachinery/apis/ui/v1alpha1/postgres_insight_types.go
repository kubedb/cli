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
	ResourceKindPostgresInsight = "PostgresInsight"
	ResourcePostgresInsight     = "postgresinsight"
	ResourcePostgresInsights    = "postgresinsights"
)

// PostgresInsightSpec defines the desired state of PostgresInsight
type PostgresInsightSpec struct {
	Version           string                      `json:"version"`
	Status            string                      `json:"status"`
	Mode              string                      `json:"mode"`
	ReplicationStatus []PostgresReplicationStatus `json:"replicationStatus"`
	ConnectionInfo    PostgresConnectionInfo      `json:"connectionInfo,omitempty"`
	VacuumInfo        PostgresVacuumInfo          `json:"vacuumInfo,omitempty"`
}

type PostgresVacuumInfo struct {
	AutoVacuum          string `json:"autoVacuum"`
	ActiveVacuumProcess *int64 `json:"activeVacuumProcess,omitempty"`
}

type PostgresConnectionInfo struct {
	MaxConnections    *int64 `json:"maxConnections,omitempty"`
	ActiveConnections *int64 `json:"activeConnections,omitempty"`
}

// Ref: https://www.postgresql.org/docs/10/monitoring-stats.html#PG-STAT-REPLICATION-VIEW

type PostgresReplicationStatus struct {
	ApplicationName string `json:"applicationName"`
	State           string `json:"state"`
	WriteLag        *int64 `json:"writeLag,omitempty"`
	FlushLag        *int64 `json:"flushLag,omitempty"`
	ReplayLag       *int64 `json:"replayLag,omitempty"`
}

// PostgresInsight is the Schema for the postgresinsights API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgresInsightSpec `json:"spec,omitempty"`
	Status api.PostgresStatus  `json:"status,omitempty"`
}

// PostgresInsightList contains a list of PostgresInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresInsight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PostgresInsight{}, &PostgresInsightList{})
}
