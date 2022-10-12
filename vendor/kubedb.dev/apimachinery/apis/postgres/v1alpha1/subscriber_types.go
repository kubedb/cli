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
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindSubscriber = "Subscriber"
	ResourceSubscriber     = "subscriber"
	ResourceSubscribers    = "subscribers"
)

// SubscriberSpec defines the desired state of Subscriber
type SubscriberSpec struct {
	// Name of the publisher
	Name string `json:"name"`
	// ServerRef specifies the database appbinding reference in any namespace.
	ServerRef core.LocalObjectReference `json:"serverRef"`
	// DatabaseName is the name of the target database inside a Postgres instance.
	DatabaseName string `json:"databaseName"`
	// Parameters to set while creating subscriber
	// +optional
	Parameters *SubscriberParameters `json:"parameters,omitempty"`

	Publisher PublisherInfo `json:"publisher"`

	// +optional
	Disable bool `json:"disable,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	// +kubebuilder:default=Delete
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`
}

type PublisherInfo struct {
	Managed  *ManagedPublisherInfo  `json:"managed,omitempty"`
	External *ExternalPublisherInfo `json:"external,omitempty"`
}

type ManagedPublisherInfo struct {
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// publication crd ref
	Refs []core.LocalObjectReference `json:"refs"`
}

type ExternalPublisherInfo struct {
	// ServerRef specifies the database appbinding reference in any namespace.
	ServerRef kmapi.ObjectReference `json:"serverRef"`
	// DatabaseName is the name of the target database inside a Postgres instance.
	DatabaseName string `json:"databaseName"`

	Publications []string `json:"publications"`
}

type SubscriberParameters struct {
	// +optional
	TableCreationPolicy TableCreationPolicy `json:"tableCreationPolicy,omitempty"`
	// +optional
	CopyData *bool `json:"copyData,omitempty"`
	// +optional
	CreateSlot *bool `json:"createSlot,omitempty"`
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// +optional
	SlotName *string `json:"slotName,omitempty"`
	// +optional
	SynchronousCommit *string `json:"synchronousCommit,omitempty"`
	// +optional
	Connect *bool `json:"connect,omitempty"`
	// +optional
	Streaming *bool `json:"streaming,omitempty"`
	// +optional
	Binary *bool `json:"binary,omitempty"`
}

type TableCreationPolicy string

const (
	TableCreationPolicyDefault            = ""
	TableCreationPolicyCreateIfNotPresent = "IfNotPresent"
)

// SubscriberStatus defines the observed state of Subscriber
type SubscriberStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase ReplicationConfigurationPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

type (
	SubscriberConditionType string
	SubscriberMessage       string
)

const (
	SubscriberConditionTypeDBServerReady  SubscriberConditionType = "DatabaseServerReady"
	SubscriberMessageDBServerNotCreated   SubscriberMessage       = "Database Server is not created yet"
	SubscriberMessageDBServerProvisioning SubscriberMessage       = "Database Server is provisioning"
	SubscriberMessageDBServerNotReady     SubscriberMessage       = "Database Server is not ready"
	SubscriberMessageDBServerCritical     SubscriberMessage       = "Database Server is critical"
	SubscriberMessageDBServerReady        SubscriberMessage       = "Database Server is Ready"

	SubscriberConditionTypeDatabaseIsFound SubscriberConditionType = "DatabaseIsFound"
	SubscriberMessageDatabaseIsFound       SubscriberMessage       = "Database is found"
	SubscriberMessageDatabaseIsNotFound    SubscriberMessage       = "Database is not found"

	SubscriberConditionTypeAllPublisherReady SubscriberConditionType = "AllPublisherReady"
	SubscriberMessageAllPublisherAreReady    SubscriberMessage       = "All Publisher are ready"
	SubscriberMessageAllPublisherAreNotReady SubscriberMessage       = "All Publisher are not ready"

	SubscriberConditionTypeSubscriberIsAllowed SubscriberConditionType = "SubscriberIsAllowed"
	SubscriberMessageSubscriberIsAllowed       SubscriberMessage       = "Subscriber is allowed"
	SubscriberMessageSubscriberIsNotAllowed    SubscriberMessage       = "Subscriber is not allowed"

	SubscriberConditionTypeAllTablesFound SubscriberConditionType = "AllTablesFound"
	SubscriberMessageAllTablesNotFound    SubscriberMessage       = "All tables are not found"
	SubscriberMessageAllTablesFound       SubscriberMessage       = "All tables are found"

	SubscriberConditionTypeSubscriptionIsSuccessful SubscriberConditionType = "SubscriptionIsSuccessful"
	SubscriberMessageSubscriptionIsSuccessful       SubscriberMessage       = "Subscription is successful"
	SubscriberMessageSubscriptionIsNotSuccessful    SubscriberMessage       = "Subscription is not successful"

	SubscriberConditionTypeSubscriptionIsDisabled SubscriberConditionType = "SubscriptionIsDisabled"
	SubscriberMessageSubscriptionIsDisabled       SubscriberMessage       = "Subsription is disabled"
	SubscriberMessageSubscriptionIsNotDisabled    SubscriberMessage       = "Subsription is not disabled"

	SubscriberConditionTypeSubscriptionIsEnabled SubscriberConditionType = "SubscriptionIsEnabled"
	SubscriberMessageSubscriptionIsEnabled       SubscriberMessage       = "Subsription is enabled"
	SubscriberMessageSubscriptionIsNotEnabled    SubscriberMessage       = "Subsription is not enabled"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=subscribers,singular=subscriber,shortName=sub,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Subscriber is the Schema for the subscribers API
type Subscriber struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubscriberSpec   `json:"spec,omitempty"`
	Status SubscriberStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// SubscriberList contains a list of Subscriber
type SubscriberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Subscriber `json:"items"`
}
