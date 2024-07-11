/*
Copyright 2024.

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

package v1alpha2

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceKindClickHouse     = "ClickHouse"
	ResourceSingularClickHouse = "clickhouse"
	ResourcePluralClickHouse   = "clickhouses"
	ResourceCodeClickHouse     = "ch"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ch,scope=Namespaced
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ClickHouse struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClickHouseSpec   `json:"spec,omitempty"`
	Status ClickHouseStatus `json:"status,omitempty"`
}

// ClickHouseSpec defines the desired state of ClickHouse
type ClickHouseSpec struct {
	// Version of ClickHouse to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a ClickHouse database.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Cluster
	// +optional
	ClusterTopology *ClusterTopology `json:"clusterTopology,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// disable security. It disables authentication security of user.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy TerminationPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 20, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`
}

type ClusterTopology struct {
	// Clickhouse Cluster Structure
	Cluster []ClusterSpec `json:"cluster,omitempty"`

	// ClickHouse Keeper server name
	ClickHouseKeeper *ClickHouseKeeperConfig `json:"clickHouseKeeper,omitempty"`
}

type ClusterSpec struct {
	// Cluster Name
	Name string `json:"name,omitempty"`
	// Number of replica for each shard to deploy for a cluster.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Number of shard to deploy for a cluster.
	// +optional
	Shards *int32 `json:"shards,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`
}

type ClickHouseKeeperConfig struct {
	Node ClickHouseKeeperNode `json:"node,omitempty"`
}

// ClickHouseKeeperNode defines item of nodes section of .spec.clusterTopology.
type ClickHouseKeeperNode struct {
	Host string `json:"host,omitempty"`

	// +optional
	Port int32 `json:"port,omitempty"`
}

// ClickHouseStatus defines the observed state of ClickHouse
type ClickHouseStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase DatabasePhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClickHouseList contains a list of ClickHouse
type ClickHouseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClickHouse `json:"items"`
}
