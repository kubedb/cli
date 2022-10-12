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
	ResourceKindPublisher = "Publisher"
	ResourcePublisher     = "publisher"
	ResourcePublishers    = "publishers"
)

// PublisherSpec defines the desired state of Publisher
type PublisherSpec struct {
	// Name of the publisher
	Name string `json:"name"`
	// ServerRef specifies the database appbinding reference in any namespace.
	ServerRef core.LocalObjectReference `json:"serverRef"`
	// DatabaseName is the name of the target database inside a Postgres instance.
	DatabaseName string `json:"databaseName"`
	// PublishAllTables is the option to publish all tables in the database
	// +optional
	PublishAllTables bool `json:"publishAllTables"`
	// Tables is the list of table to publish
	// +optional
	Tables []string `json:"tables,omitempty"`
	// Parameters to set while creating publisher
	// +optional
	Parameters *PublisherParameters `json:"parameters,omitempty"`

	// AllowedSubscribers defines the types of database schemas that MAY refer to
	// a database instance and the trusted namespaces where those schema resources MAY be
	// present.
	//
	// +kubebuilder:default={namespaces:{from: Same}}
	// +optional
	AllowedSubscribers *AllowedSubscribers `json:"allowedSubscribers,omitempty"`

	// +optional
	Disable bool `json:"disable,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	// +kubebuilder:default=Delete
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`
}

// +kubebuilder:validation:Enum=insert;update;delete;truncate
type DMLOperation string

const (
	DMLOpInsert   = "insert"
	DMLOpUpdate   = "update"
	DMLOpDelete   = "delete"
	DMLOpTruncate = "truncate"
)

type PublisherParameters struct {
	// Publish parameters determines which DML operations will be published by the new publication to the subscribers.
	// The allowed operations are insert, update, delete, and truncate.
	// The default is to publish all actions, and so the default value for this option is 'insert, update, delete, truncate'.
	// +optional
	Operations []DMLOperation `json:"operations,omitempty"`
	// PublishViaPartitionRoot parameter determines whether changes in a partitioned table (or on its partitions)
	// contained in the publication will be published using the identity
	// and schema of the partitioned table rather than that of the individual partitions that are actually changed; the latter is the default.
	// Enabling this allows the changes to be replicated into a non-partitioned table
	// or a partitioned table consisting of a different set of partitions.
	// +optional
	PublishViaPartitionRoot *bool `json:"publishViaPartitionRoot,omitempty"`
}

// PublisherStatus defines the observed state of Publisher
type PublisherStatus struct {
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
	// Database authentication secret
	// +optional
	Subscribers []kmapi.ObjectReference `json:"subscribers,omitempty"`
}

type (
	PublisherConditionType string
	PublisherMessage       string
)

const (
	PublisherConditionTypeDBServerReady  PublisherConditionType = "DatabaseServerReady"
	PublisherMessageDBServerNotCreated   PublisherMessage       = "Database Server is not created yet"
	PublisherMessageDBServerProvisioning PublisherMessage       = "Database Server is provisioning"
	PublisherMessageDBServerNotReady     PublisherMessage       = "Database Server is not ready"
	PublisherMessageDBServerCritical     PublisherMessage       = "Database Server is critical"
	PublisherMessageDBServerReady        PublisherMessage       = "Database Server is Ready"

	PublisherConditionTypeDatabaseIsFound PublisherConditionType = "DatabaseIsFound"
	PublisherMessageDatabaseIsFound       PublisherMessage       = "Database is found"
	PublisherMessageDatabaseIsNotFound    PublisherMessage       = "Database is not found"

	PublisherConditionTypeAllTablesFound PublisherConditionType = "AllTablesFound"
	PublisherMessageAllTablesNotFound    PublisherMessage       = "All tables are not found"
	PublisherMessageAllTablesFound       PublisherMessage       = "All tables are found"

	PublisherConditionTypePublicationSuccessful PublisherConditionType = "PublicationSuccessful"
	PublisherMessagePublicationIsSuccessful     PublisherMessage       = "Publication is successful"
	PublisherMessagePublicationIsNotSuccessful  PublisherMessage       = "Publication is not successful"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=publishers,singular=publisher,shortName=pub,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Publisher is the Schema for the publishers API
type Publisher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PublisherSpec   `json:"spec,omitempty"`
	Status PublisherStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// PublisherList contains a list of Publisher
type PublisherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Publisher `json:"items"`
}
