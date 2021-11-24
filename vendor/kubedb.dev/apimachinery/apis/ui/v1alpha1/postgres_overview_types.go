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
	ResourceKindPostgresOverview = "PostgresOverview"
	ResourcePostgresOverview     = "postgresoverview"
	ResourcePostgresOverviews    = "postgresoverviews"
)

// PostgresOverviewSpec defines the desired state of PostgresOverview
type PostgresOverviewSpec struct {
	Version           string                      `json:"version" protobuf:"bytes,1,opt,name=version"`
	ConnectionURL     string                      `json:"connectionURL" protobuf:"bytes,2,opt,name=connectionURL"`
	Status            string                      `json:"status" protobuf:"bytes,3,opt,name=status"`
	Mode              string                      `json:"mode" protobuf:"bytes,4,opt,name=mode"`
	ReplicationStatus []PostgresReplicationStatus `json:"replicationStatus,omitempty" protobuf:"bytes,5,rep,name=replicationStatus"`
	ConnectionInfo    PostgresConnectionInfo      `json:"connectionInfo,omitempty" protobuf:"bytes,6,opt,name=connectionInfo"`
	BackupInfo        PostgresBackupInfo          `json:"backupInfo,omitempty" protobuf:"bytes,7,opt,name=backupInfo"`
	VacuumInfo        PostgresVacuumInfo          `json:"vacuumInfo,omitempty" protobuf:"bytes,8,opt,name=vacuumInfo"`
}

type PostgresVacuumInfo struct {
	AutoVacuum          string `json:"autoVacuum" protobuf:"bytes,1,opt,name=autoVacuum"`
	ActiveVacuumProcess int64  `json:"activeVacuumProcess" protobuf:"varint,2,opt,name=activeVacuumProcess"`
}

type PostgresBackupInfo struct {
}

type PostgresConnectionInfo struct {
	MaxConnections    int64 `json:"maxConnections" protobuf:"varint,1,opt,name=maxConnections"`
	ActiveConnections int64 `json:"activeConnections" protobuf:"varint,2,opt,name=activeConnections"`
}

// Ref: https://www.postgresql.org/docs/10/monitoring-stats.html#PG-STAT-REPLICATION-VIEW

type PostgresReplicationStatus struct {
	ApplicationName string `json:"applicationName" protobuf:"bytes,1,opt,name=applicationName"`
	State           string `json:"state" protobuf:"bytes,2,opt,name=state"`
	WriteLag        int64  `json:"writeLag" protobuf:"varint,3,opt,name=writeLag"`
	FlushLag        int64  `json:"flushLag" protobuf:"varint,4,opt,name=flushLag"`
	ReplayLag       int64  `json:"replayLag" protobuf:"varint,5,opt,name=replayLag"`
}

// PostgresOverviewStatus defines the observed state of PostgresOverview
type PostgresOverviewStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// PostgresOverview is the Schema for the postgresoverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   PostgresOverviewSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status api.PostgresStatus   `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// PostgresOverviewList contains a list of PostgresOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []PostgresOverview `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&PostgresOverview{}, &PostgresOverviewList{})
}
