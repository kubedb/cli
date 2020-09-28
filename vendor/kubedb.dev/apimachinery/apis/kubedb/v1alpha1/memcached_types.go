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
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeMemcached     = "mc"
	ResourceKindMemcached     = "Memcached"
	ResourceSingularMemcached = "memcached"
	ResourcePluralMemcached   = "memcacheds"
)

// Memcached defines a Memcached database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=memcacheds,singular=memcached,shortName=mc,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Memcached struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MemcachedSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            MemcachedStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type MemcachedSpec struct {
	// Version of Memcached to be deployed.
	Version string `json:"version" protobuf:"bytes,5,opt,name=version"`

	// Number of instances to deploy for a Memcached database.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,6,opt,name=replicas"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,7,opt,name=monitor"`

	// ConfigSource is an optional field to provide custom configuration file for database.
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty" protobuf:"bytes,8,opt,name=configSource"`

	// DataVolume is an optional field to add one volume to each
	// memcached pod.  The volume will be made available under
	// /data and owned by the memcached user.
	//
	// While not mandated by the API and not configured
	// automatically, the intended purpose is to use that volume
	// for memcached's persistent memory support
	// (https://memcached.org/blog/persistent-memory/) by adding
	// the memory-file and memory-limit options to the config
	// (https://github.com/memcached/memcached/wiki/WarmRestart).
	//
	// For that purpose, a CSI inline volume provided by PMEM-CSI
	// can be used, in which case each pod will get its own, empty
	// volume. Warm restarts are not supported.
	//
	// For testing, an empty dir can be used instead.
	DataVolume *core.VolumeSource `json:"dataVolume,omitempty" protobuf:"bytes,9,opt,name=dataVolume"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,10,opt,name=podTemplate"`

	// ServiceTemplate is an optional configuration for service used to expose database
	// +optional
	ServiceTemplate ofst.ServiceTemplateSpec `json:"serviceTemplate,omitempty" protobuf:"bytes,11,opt,name=serviceTemplate"`

	// TLS contains tls configurations
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty" protobuf:"bytes,12,opt,name=tls"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty" protobuf:"varint,13,opt,name=halted"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty" protobuf:"bytes,14,opt,name=terminationPolicy,casttype=TerminationPolicy"`
}

// +kubebuilder:validation:Enum=server;metrics-exporter
type MemcachedCertificateAlias string

const (
	MemcachedServerCert          MemcachedCertificateAlias = "server"
	MemcachedMetricsExporterCert MemcachedCertificateAlias = "metrics-exporter"
)

type MemcachedStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase DatabasePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=DatabasePhase"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty" protobuf:"bytes,3,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MemcachedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Memcached TPR objects
	Items []Memcached `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
