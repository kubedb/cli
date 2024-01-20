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
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeDruid     = "dr"
	ResourceKindDruid     = "Druid"
	ResourceSingularDruid = "druid"
	ResourcePluralDruid   = "druids"
)

// Druid is the Schema for the druids API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=dr,scope=Namespaced
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Druid struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DruidSpec   `json:"spec,omitempty"`
	Status DruidStatus `json:"status,omitempty"`
}

// DruidSpec defines the desired state of Druid
type DruidSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Version of Druid to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Druid database
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Druid topology for node specification
	// +optional
	Topology *DruidClusterTopology `json:"topology,omitempty"`

	// StorageType can be durable (default) or ephemeral.
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	// Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// To enable ssl for http layer.
	// +optional
	// EnableSSL bool `json:"enableSSL,omitempty"`

	// disable security. It disables authentication security of user.
	// If unset, default is false
	// +optional
	DisableSecurity *bool `json:"disableSecurity,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *core.LocalObjectReference `json:"authSecret,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e. config.properties).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	//// TLS contains tls configurations
	//// +optional
	//TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// MetadataStorage contains information for Druid to connect to external dependency metadata storage
	// +optional
	MetadataStorage *MetadataStorage `json:"metadataStorage,omitempty"`

	// DeepStorage contains specification for druid to connect to the deep storage
	DeepStorage *DeepStorageSpec `json:"deepStorage"`

	// ZooKeeper contains information for Druid to connect to external dependency metadata storage
	ZooKeeper *ZooKeeperRef `json:"zooKeeper"`

	// PodTemplate is an optional configuration
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 30, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`
}

type DruidClusterTopology struct {
	Coordinators *DruidNode `json:"coordinators"`
	// +optional
	Overlords *DruidNode `json:"overlords,omitempty"`

	MiddleManagers *DruidNode `json:"middleManagers"`

	Historicals *DruidNode `json:"historicals"`

	Brokers *DruidNode `json:"brokers"`
	// +optional
	Routers *DruidNode `json:"routers,omitempty"`
}

type DruidNode struct {
	// Replicas represents number of replica for the specific type of node
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Suffix to append with node name
	// +optional
	Suffix string `json:"suffix,omitempty"`

	// Storage to specify how storage shall be used.
	// +optional
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	// +mapType=atomic
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
}

type MetadataStorage struct {
	// Name of the appbinding of metadata storage
	Name string `json:"name"`

	// Namespace of the appbinding of metadata storage
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// If Druid has the permission to create new tables
	// +optional
	CreateTables *bool `json:"createTables,omitempty"`
}

type DeepStorageSpec struct {
	// Specifies the storage type to be used by druid
	// Possible values: s3, google, azure, hdfs
	Type *string `json:"type"`

	// deepStorage.configSecret should contain the necessary data
	// to connect to the deep storage
	ConfigSecret *core.LocalObjectReference `json:"configSecret"`
}

type ZooKeeperRef struct {
	// Name of the appbinding of zookeeper
	Name *string `json:"name"`

	// Namespace of the appbinding of zookeeper
	// +optional
	Namespace string `json:"namespace"`

	// Base ZooKeeperSpec path
	// +optional
	PathsBase string `json:"pathsBase"`
}

// DruidStatus defines the observed state of Druid
type DruidStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Specifies the current phase of the database
	// +optional
	Phase DruidPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DruidList contains a list of Druid
type DruidList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Druid `json:"items"`
}

// +kubebuilder:validation:Enum=Provisioning;Ready;NotReady;Critical
type DruidPhase string

const (
	DruidPhaseProvisioning DruidPhase = "Provisioning"
	DruidPhaseReady        DruidPhase = "Ready"
	DruidPhaseNotReady     DruidPhase = "NotReady"
	DruidPhaseCritical     DruidPhase = "Critical"
)

type DruidNodeRoleType string

const (
	DruidNodeRoleCoordinators   DruidNodeRoleType = "coordinators"
	DruidNodeRoleOverlords      DruidNodeRoleType = "overlords"
	DruidNodeRoleBrokers        DruidNodeRoleType = "brokers"
	DruidNodeRoleRouters        DruidNodeRoleType = "routers"
	DruidNodeRoleMiddleManagers DruidNodeRoleType = "middleManagers"
	DruidNodeRoleHistoricals    DruidNodeRoleType = "historicals"
)
