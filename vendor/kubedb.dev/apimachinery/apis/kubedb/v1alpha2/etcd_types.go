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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeEtcd     = "etcd"
	ResourceKindEtcd     = "Etcd"
	ResourceSingularEtcd = "etcd"
	ResourcePluralEtcd   = "etcds"
)

// +kubebuilder:validation:Enum=server;client;peer;metrics-exporter
type EtcdCertificateAlias string

const (
	// EtcdServerCert is used by an etcd member to serve the client (gRPC/HTTP) API.
	EtcdServerCert EtcdCertificateAlias = "server"
	// EtcdClientCert is used by KubeDB (health checker, exporter, ops) to talk to etcd.
	EtcdClientCert EtcdCertificateAlias = "client"
	// EtcdPeerCert is used for the member-to-member (Raft) traffic.
	EtcdPeerCert EtcdCertificateAlias = "peer"
	// EtcdMetricsExporterCert is used to serve the exporter metrics endpoint over TLS.
	EtcdMetricsExporterCert EtcdCertificateAlias = "metrics-exporter"
)

// +kubebuilder:validation:Enum=periodic;revision
type EtcdAutoCompactionMode string

const (
	EtcdAutoCompactionModePeriodic EtcdAutoCompactionMode = "periodic"
	EtcdAutoCompactionModeRevision EtcdAutoCompactionMode = "revision"
)

// Etcd is the Schema for the etcds API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:path=etcds,singular=etcd,shortName=etcd,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Etcd struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec defines the desired state of Etcd
	// +optional
	Spec EtcdSpec `json:"spec,omitempty"`

	// status defines the observed state of Etcd
	// +optional
	Status EtcdStatus `json:"status,omitempty"`
}

// EtcdSpec defines the desired state of Etcd.
//
// Etcd is a single-mode, natively clustered (Raft) datastore. There is no
// standalone-vs-cluster topology switch: a one member cluster is simply
// spec.replicas=1. Members are added/removed through the etcd membership API,
// so no external leader-election sidecar is involved.
type EtcdSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of Etcd to be deployed. It refers to the name of an EtcdVersion
	// catalog object.
	Version string `json:"version"`

	// Number of etcd members to deploy. An odd number is strongly recommended
	// so that the Raft quorum tolerates the maximum number of failures.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// StorageType can be durable (default) or ephemeral
	// +optional
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used
	// +optional
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for
	// the etcd members (i.e. etcd.conf.yaml). Deprecated in favour of
	// spec.configuration, kept for parity with the other KubeDB databases.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// Configuration holds the custom config for etcd
	// +optional
	Configuration *EtcdConfiguration `json:"configuration,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// TLS configures certificates issued from spec.tls.issuerRef for the etcd
	// client API, the peer (Raft) traffic and the metrics exporter.
	// When TLS is specified, issuerRef must be set.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// AllowedSchemas defines the types of database schemas that may refer to
	// a database instance and the trusted namespaces where those schema resources may be
	// present.
	// +kubebuilder:default={namespaces:{from: Same}}
	// +optional
	AllowedSchemas *AllowedConsumers `json:"allowedSchemas,omitempty"`

	// Archiver controls database backup using Archiver CR
	// +optional
	Archiver *Archiver `json:"archiver,omitempty"`
}

// EtcdConfiguration holds the user supplied etcd configuration.
type EtcdConfiguration struct {
	ConfigurationSpec `json:",inline"`

	// Tuning exposes the handful of etcd knobs that KubeDB manages explicitly
	// (they end up on the etcd command line rather than being merged blindly
	// from the config secret).
	// +optional
	Tuning *EtcdTuningConfig `json:"tuning,omitempty"`
}

// EtcdTuningConfig holds the KubeDB-managed etcd tuning knobs.
type EtcdTuningConfig struct {
	// QuotaBackendBytes sets the maximum size of the backend database
	// (--quota-backend-bytes). etcd goes into a read-only alarm state once the
	// backend grows past this value.
	// +optional
	QuotaBackendBytes *int64 `json:"quotaBackendBytes,omitempty"`

	// AutoCompactionMode selects the interpretation of autoCompactionRetention
	// (--auto-compaction-mode). Either "periodic" or "revision".
	// +optional
	AutoCompactionMode *EtcdAutoCompactionMode `json:"autoCompactionMode,omitempty"`

	// AutoCompactionRetention configures the auto compaction retention
	// (--auto-compaction-retention). It is a duration (e.g. "1h") when mode is
	// periodic and a revision count (e.g. "1000") when mode is revision.
	// +optional
	AutoCompactionRetention *string `json:"autoCompactionRetention,omitempty"`

	// SnapshotCount is the number of committed transactions that trigger a
	// snapshot to disk (--snapshot-count).
	// +optional
	SnapshotCount *uint64 `json:"snapshotCount,omitempty"`
}

// EtcdStatus defines the observed state of Etcd.
type EtcdStatus struct {
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
	AuthSecret *Age `json:"authSecret,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdList contains a list of Etcd
type EtcdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Etcd CRD objects
	Items []Etcd `json:"items"`
}

var _ Accessor = &Etcd{}

func (e *Etcd) GetObjectMeta() metav1.ObjectMeta {
	return e.ObjectMeta
}

func (e *Etcd) GetConditions() []kmapi.Condition {
	return e.Status.Conditions
}

func (e *Etcd) SetCondition(cond kmapi.Condition) {
	e.Status.Conditions = setCondition(e.Status.Conditions, cond)
}

func (e *Etcd) RemoveCondition(typ string) {
	e.Status.Conditions = removeCondition(e.Status.Conditions, typ)
}

func (e *Etcd) HasCondition(typ string) bool {
	for _, c := range e.Status.Conditions {
		if string(c.Type) == typ {
			return true
		}
	}
	return false
}
