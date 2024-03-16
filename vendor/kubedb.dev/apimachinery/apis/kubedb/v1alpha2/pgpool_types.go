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

package v1alpha2

import (
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodePgpool     = "pp"
	ResourceKindPgpool     = "Pgpool"
	ResourceSingularPgpool = "pgpool"
	ResourcePluralPgpool   = "pgpools"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Pgpool is the Schema for the pgpools API
// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=pgpools,singular=pgpool,shortName=pp,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Pgpool struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            PgpoolSpec   `json:"spec,omitempty"`
	Status          PgpoolStatus `json:"status,omitempty"`
}

// PgpoolSpec defines the desired state of Pgpool
type PgpoolSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// SyncUsers is a boolean type and when enabled, operator fetches all users created in the backend server to the
	// Pgpool server . Password changes are also synced in pgpool when it is enabled.
	// +optional
	SyncUsers bool `json:"syncUsers,omitempty"`

	// Version of Pgpool to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Pgpool instance.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// PostgresRef refers to the AppBinding of the backend PostgreSQL server
	PostgresRef *kmapi.ObjectReference `json:"postgresRef"`

	// Pgpool secret containing username and password for pgpool pcp user
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// ConfigSecret is a configuration secret which will be created with default and InitConfiguration
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose Pgpool
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// InitConfiguration contains information with which the Pgpool will bootstrap
	// +optional
	InitConfiguration *PgpoolConfiguration `json:"initConfig,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose Pgpool
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Monitor is used to monitor Pgpool instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// TerminationPolicy controls the delete operation for Pgpool
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty"`

	// PodPlacementPolicy is the reference of the podPlacementPolicy
	// +kubebuilder:default={name: "default"}
	// +optional
	PodPlacementPolicy *core.LocalObjectReference `json:"podPlacementPolicy,omitempty"`
}

// PgpoolStatus defines the observed state of Pgpool
type PgpoolStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

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
	// +optional
	Gateway *Gateway `json:"gateway,omitempty"`
}

type PgpoolConfiguration struct {
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	PgpoolConfig *runtime.RawExtension `json:"pgpoolConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgpoolList contains a list of Pgpool
type PgpoolList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []Pgpool `json:"items"`
}
