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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceCodeConnector     = "kc"
	ResourceKindConnector     = "Connector"
	ResourceSingularConnector = "connector"
	ResourcePluralConnector   = "connectors"
)

// Connector defines to run in a Kafka Connect cluster to read data from Kafka topics and write the data into another system.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=connectors,singular=connector,shortName=kc,categories={kfstore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="ConnectCluster",type="string",JSONPath=".spec.connectClusterRef.name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Connector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectorSpec   `json:"spec,omitempty"`
	Status ConnectorStatus `json:"status,omitempty"`
}

// ConnectorSpec defines the desired state of Connector
type ConnectorSpec struct {
	// Kafka connect cluster app-binding reference
	// ConnectClusterRef is a required field, where Connector will add tasks to produce or consume data from kafka topics.
	ConnectClusterRef *kmapi.ObjectReference `json:"connectClusterRef"`

	// ConfigSecret is a required field to provide configuration file for Connector to create connectors for Kafka connect cluster(i.e connector.properties).
	ConfigSecret *core.LocalObjectReference `json:"configSecret"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy dbapi.DeletionPolicy `json:"deletionPolicy,omitempty"`
}

// ConnectorStatus defines the observed state of connectors
type ConnectorStatus struct {
	// Specifies the current phase of the Connector
	// +optional
	Phase ConnectorPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the Connector, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Unassigned;Running;Paused;Failed;Restarting;Stopped;Destroyed;Unknown
type ConnectorPhase string

const (
	ConnectorPhasePending    ConnectorPhase = "Pending"
	ConnectorPhaseUnassigned ConnectorPhase = "Unassigned"
	ConnectorPhaseRunning    ConnectorPhase = "Running"
	ConnectorPhasePaused     ConnectorPhase = "Paused"
	ConnectorPhaseFailed     ConnectorPhase = "Failed"
	ConnectorPhaseRestarting ConnectorPhase = "Restarting"
	ConnectorPhaseStopped    ConnectorPhase = "Stopped"
	ConnectorPhaseDestroyed  ConnectorPhase = "Destroyed"
	ConnectorPhaseUnknown    ConnectorPhase = "Unknown"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConnectorList contains a list of Connector
type ConnectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Connector `json:"items"`
}
