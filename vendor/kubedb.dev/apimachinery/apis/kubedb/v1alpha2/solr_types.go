/*
Copyright 2023.

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
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeSolr     = "sl"
	ResourceKindSolr     = "Solr"
	ResourceSingularSolr = "solr"
	ResourcePluralSolr   = "solrs"
)

// Solr is the schema for the Sole API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=solrs,singular=solr,shortName=sl,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Solr struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SolrSpec   `json:"spec,omitempty"`
	Status SolrStatus `json:"status,omitempty"`
}

// SolrSpec defines the desired state of Solr c
type SolrSpec struct {
	// Version of Solr to be deployed
	Version string `json:"version"`

	// Number of instances to deploy for a Solr database
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Solr topology for node specification
	// +optional
	Topology *SolrClusterTopology `json:"topology,omitempty"`

	// StorageType van be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// 	// ZooKeeper contains information for Solr to store configurations for collections
	ZookeeperRef *kmapi.ObjectReference `json:"zookeeperRef,omitempty"`

	// +optional
	SolrModules []string `json:"solrModules,omitempty"`

	// +optional
	SolrOpts []string `json:"solrOpts,omitempty"`

	// To enable ssl for http layer
	EnableSSL bool `json:"enableSSL,omitempty"`

	// Client auth need or want
	ClientAuthSSL string `json:"clientAuthSSL,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// Disable security. It disables authentication security of users.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// +optional
	KeystoreSecret *core.LocalObjectReference `json:"keystoreSecret,omitempty"`

	// +optional
	AuthSecret *core.LocalObjectReference `json:"authSecret,omitempty"`

	// +optional
	ZookeeperDigestSecret *core.LocalObjectReference `json:"zookeeperDigestSecret,omitempty"`

	// +optional
	ZookeeperDigestReadonlySecret *core.LocalObjectReference `json:"zookeeperDigestReadonlySecret,omitempty"`

	// +optional
	AuthConfigSecret *core.LocalObjectReference `json:"authConfigSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

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

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`
}

type SolrClusterTopology struct {
	Overseer    *SolrNode `json:"overseer,omitempty"`
	Data        *SolrNode `json:"data,omitempty"`
	Coordinator *SolrNode `json:"coordinator,omitempty"`
}

type SolrNode struct {
	// Replica represents number of replica for this specific type of nodes
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// suffix to append with node name
	// +optional
	Suffix string `json:"suffix,omitempty"`

	// Storage to specify how storage shall be used.
	// +optional
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`
}

// SolrStatus defines the observed state of Solr
type SolrStatus struct {
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

// +kubebuilder:validation:Enum=overseer;data;coordinator;combined
type SolrNodeRoleType string

const (
	SolrNodeRoleOverseer    SolrNodeRoleType = "overseer"
	SolrNodeRoleData        SolrNodeRoleType = "data"
	SolrNodeRoleCoordinator SolrNodeRoleType = "coordinator"
	SolrNodeRoleSet                          = "set"
)

// +kubebuilder:validation:Enum=ca;transport;http;client;server
type SolrCertificateAlias string

const (
	SolrCACert        SolrCertificateAlias = "ca"
	SolrTransportCert SolrCertificateAlias = "transport"
	SolrHTTPCert      SolrCertificateAlias = "http"
	SolrClientCert    SolrCertificateAlias = "client"
	SolrServerCert    SolrCertificateAlias = "server"
)

//+kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SolrList contains a list of Solr
type SolrList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Solr `json:"items"`
}
