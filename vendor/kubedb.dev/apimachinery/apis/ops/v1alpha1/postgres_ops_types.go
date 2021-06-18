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
	apis "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceCodePostgresOpsRequest     = "pgops"
	ResourceKindPostgresOpsRequest     = "PostgresOpsRequest"
	ResourceSingularPostgresOpsRequest = "postgresopsrequest"
	ResourcePluralPostgresOpsRequest   = "postgresopsrequests"
)

// PostgresOpsRequest defines a PostgreSQL DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=postgresopsrequests,singular=postgresopsrequest,shortName=pgops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PostgresOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              PostgresOpsRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            PostgresOpsRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
type PostgresTLSSpec struct {
	TLSSpec `json:",inline,omitempty" protobuf:"bytes,1,opt,name=tLSSpec"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	// +optional
	SSLMode apis.PostgresSSLMode `json:"sslMode,omitempty" protobuf:"bytes,2,opt,name=sslMode,casttype=PostgresSSLMode"`

	// ClientAuthMode for sidecar or sharding. (default will be md5. [md5;scram;cert])
	// +optional
	ClientAuthMode apis.PostgresClientAuthMode `json:"clientAuthMode,omitempty" protobuf:"bytes,3,opt,name=clientAuthMode,casttype=PostgresClientAuthMode"`
}

// PostgresOpsRequestSpec is the spec for PostgresOpsRequest
type PostgresOpsRequestSpec struct {
	// Specifies the Postgres reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=OpsRequestType"`
	// Specifies information necessary for upgrading Postgres
	Upgrade *PostgresUpgradeSpec `json:"upgrade,omitempty" protobuf:"bytes,3,opt,name=upgrade"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *PostgresHorizontalScalingSpec `json:"horizontalScaling,omitempty" protobuf:"bytes,4,opt,name=horizontalScaling"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *PostgresVerticalScalingSpec `json:"verticalScaling,omitempty" protobuf:"bytes,5,opt,name=verticalScaling"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *PostgresVolumeExpansionSpec `json:"volumeExpansion,omitempty" protobuf:"bytes,6,opt,name=volumeExpansion"`
	// Specifies information necessary for custom configuration of Postgres
	Configuration *PostgresCustomConfigurationSpec `json:"configuration,omitempty" protobuf:"bytes,7,opt,name=configuration"`

	// Specifies information necessary for configuring TLS
	TLS *PostgresTLSSpec `json:"tls,omitempty" protobuf:"bytes,8,opt,name=tls"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty" protobuf:"bytes,9,opt,name=restart"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty" protobuf:"varint,10,opt,name=timeout"`
}

type PostgresUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty" protobuf:"bytes,1,opt,name=targetVersion"`
}

// HorizontalScaling is the spec for Postgres horizontal scaling
type PostgresHorizontalScalingSpec struct {
	Replicas *int32 `json:"replicas,omitempty" protobuf:"bytes,1,opt,name=replicas"`
}

// PostgresVerticalScalingSpec is the spec for Postgres vertical scaling
type PostgresVerticalScalingSpec struct {
	Postgres *core.ResourceRequirements `json:"postgres,omitempty" protobuf:"bytes,1,opt,name=postgres"`
	Exporter *core.ResourceRequirements `json:"exporter,omitempty" protobuf:"bytes,2,opt,name=exporter"`
}

// PostgresVolumeExpansionSpec is the spec for Postgres volume expansion
type PostgresVolumeExpansionSpec struct {
	// volume specification for Postgres
	Postgres *resource.Quantity `json:"postgres,omitempty" protobuf:"bytes,1,opt,name=postgres"`
}

type PostgresCustomConfigurationSpec struct {
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty" protobuf:"bytes,2,opt,name=configSecret"`
	InlineConfig       string                     `json:"inlineConfig,omitempty" protobuf:"bytes,3,opt,name=inlineConfig"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty" protobuf:"varint,4,opt,name=removeCustomConfig"`
}

type PostgresCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty" protobuf:"bytes,1,opt,name=configMap"`
	Data      map[string]string          `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`
	Remove    bool                       `json:"remove,omitempty" protobuf:"varint,3,opt,name=remove"`
}

// PostgresOpsRequestStatus is the status for PostgresOpsRequest
type PostgresOpsRequestStatus struct {
	Phase OpsRequestPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=OpsRequestPhase"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty" protobuf:"bytes,3,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresOpsRequestList is a list of PostgresOpsRequests
type PostgresOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of PostgresOpsRequest CRD objects
	Items []PostgresOpsRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
