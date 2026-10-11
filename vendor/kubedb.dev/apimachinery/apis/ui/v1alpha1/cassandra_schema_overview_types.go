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
	ResourceKindCassandraSchemaOverview = "CassandraSchemaOverview"
	ResourceCassandraSchemaOverview     = "cassandraschemaoverview"
	ResourceCassandraSchemaOverviews    = "cassandraschemaoverviews"
)

// CassandraSchemaOverviewSpec defines the desired state of CassandraSchemaOverview.
//
// Cassandra's unit is a keyspace containing tables, and a keyspace carries replication
// settings that have no equivalent in GenericSchemaOverviewSpec, so this does not alias it.
type CassandraSchemaOverviewSpec struct {
	Keyspaces []CassandraKeyspaceSpec `json:"keyspaces"`
}

// CassandraKeyspaceSpec describes one keyspace and the tables it holds.
type CassandraKeyspaceSpec struct {
	Name string `json:"name"`

	// ReplicationClass is the replication strategy, e.g. SimpleStrategy or
	// NetworkTopologyStrategy, from system_schema.keyspaces.replication['class'].
	// +optional
	ReplicationClass string `json:"replicationClass,omitempty"`

	// ReplicationFactors maps each remaining replication option to its value: for
	// SimpleStrategy that is replication_factor, for NetworkTopologyStrategy one entry
	// per datacenter. Cassandra stores these as strings.
	// +optional
	ReplicationFactors map[string]string `json:"replicationFactors,omitempty"`

	// +optional
	DurableWrites bool `json:"durableWrites,omitempty"`

	// Tables lists the table names in the keyspace, sorted.
	//
	// No per-table size is reported: Cassandra's catalog carries none, and a real size
	// needs nodetool tablestats or JMX, which is a different transport than CQL.
	// +optional
	Tables []string `json:"tables,omitempty"`
}

// CassandraSchemaOverview is the Schema for the CassandraSchemaOverviews API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CassandraSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CassandraSchemaOverviewSpec `json:"spec,omitempty"`
}

// CassandraSchemaOverviewList contains a list of CassandraSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CassandraSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CassandraSchemaOverview `json:"items"`
}
