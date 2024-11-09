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

package v1

import (
	"sync"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofstv1 "kmodules.xyz/offshoot-api/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	once          sync.Once
	DefaultClient client.Client
)

func SetDefaultClient(kc client.Client) {
	once.Do(func() {
		DefaultClient = kc
	})
}

type InitSpec struct {
	// Initialized indicates that this database has been initialized.
	// This will be set by the operator when status.conditions["Provisioned"] is set to ensure
	// that database is not mistakenly reset when recovered using disaster recovery tools.
	Initialized bool `json:"initialized,omitempty"`
	// Wait for initial DataRestore condition
	WaitForInitialRestore bool              `json:"waitForInitialRestore,omitempty"`
	Script                *ScriptSourceSpec `json:"script,omitempty"`

	Archiver *ArchiverRecovery `json:"archiver,omitempty"`
}

type ScriptSourceSpec struct {
	ScriptPath        string `json:"scriptPath,omitempty"`
	core.VolumeSource `json:",inline,omitempty"`
	Git               *GitRepo `json:"git,omitempty"`
}

type GitRepo struct {
	// https://github.com/kubernetes/git-sync/tree/master
	Args []string `json:"args"`
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	Env []core.EnvVar `json:"env,omitempty"`
	// Security options the pod should run with.
	// More info: https://kubernetes.io/docs/concepts/policy/security-context/
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +optional
	SecurityContext *core.SecurityContext `json:"securityContext,omitempty"`
	// Compute Resources required by the sidecar container.
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`
	// Authentication secret for git repository
	// +optional
	AuthSecret *core.LocalObjectReference `json:"authSecret,omitempty"`
}

type RemoteReplicaSpec struct {
	// SourceRef specifies the  source object
	SourceRef core.ObjectReference `json:"sourceRef" protobuf:"bytes,1,opt,name=sourceRef"`
}

// +kubebuilder:validation:Enum=Provisioning;DataRestoring;Ready;Critical;NotReady;Halted;Unknown
type DatabasePhase string

const (
	// used for Databases that are currently provisioning
	DatabasePhaseProvisioning DatabasePhase = "Provisioning"
	// used for Databases for which data is currently restoring
	DatabasePhaseDataRestoring DatabasePhase = "DataRestoring"
	// used for Databases that are currently ReplicaReady, AcceptingConnection and Ready
	DatabasePhaseReady DatabasePhase = "Ready"
	// used for Databases that can connect, ReplicaReady == false || Ready == false (eg, ES yellow)
	DatabasePhaseCritical DatabasePhase = "Critical"
	// used for Databases that can't connect
	DatabasePhaseNotReady DatabasePhase = "NotReady"
	// used for Databases that are halted
	DatabasePhaseHalted DatabasePhase = "Halted"
	// used for Databases for which Phase can't be calculated
	DatabasePhaseUnknown DatabasePhase = "Unknown"
)

// +kubebuilder:validation:Enum=Durable;Ephemeral
type StorageType string

const (
	// default storage type and requires spec.storage to be configured
	StorageTypeDurable StorageType = "Durable"
	// Uses emptyDir as storage
	StorageTypeEphemeral StorageType = "Ephemeral"
)

// +kubebuilder:validation:Enum=Halt;Delete;WipeOut;DoNotTerminate
type DeletionPolicy string

const (
	// Deletes database pods, service but leave the PVCs and stash backup data intact.
	DeletionPolicyHalt DeletionPolicy = "Halt"
	// Deletes database pods, service, pvcs but leave the stash backup data intact.
	DeletionPolicyDelete DeletionPolicy = "Delete"
	// Deletes database pods, service, pvcs and stash backup data.
	DeletionPolicyWipeOut DeletionPolicy = "WipeOut"
	// Rejects attempt to delete database using ValidationWebhook.
	DeletionPolicyDoNotTerminate DeletionPolicy = "DoNotTerminate"
)

// +kubebuilder:validation:Enum=primary;standby;stats
type ServiceAlias string

const (
	PrimaryServiceAlias ServiceAlias = "primary"
	StandbyServiceAlias ServiceAlias = "standby"
	StatsServiceAlias   ServiceAlias = "stats"
)

// +kubebuilder:validation:Enum=fscopy;clone;sync;none
type PITRReplicationStrategy string

const (
	// ReplicationStrategySync means data will be synced from primary to secondary
	ReplicationStrategySync PITRReplicationStrategy = "sync"
	// ReplicationStrategyFSCopy means data will be copied from filesystem
	ReplicationStrategyFSCopy PITRReplicationStrategy = "fscopy"
	// ReplicationStrategyClone means volumeSnapshot will be used to create pvc's
	ReplicationStrategyClone PITRReplicationStrategy = "clone"
	// ReplicationStrategyNone means no replication will be used
	// data will be fully restored in every replicas instead of replication
	ReplicationStrategyNone PITRReplicationStrategy = "none"
)

// +kubebuilder:validation:Enum=DNS;IP;IPv4;IPv6
type AddressType string

const (
	AddressTypeDNS AddressType = "DNS"
	// Uses spec.podIP as address for db pods.
	AddressTypeIP AddressType = "IP"
	// Uses first IPv4 address from spec.podIP, spec.podIPs fields as address for db pods.
	AddressTypeIPv4 AddressType = "IPv4"
	// Uses first IPv6 address from spec.podIP, spec.podIPs fields as address for db pods.
	AddressTypeIPv6 AddressType = "IPv6"
)

func (a AddressType) IsIP() bool {
	return a == AddressTypeIP || a == AddressTypeIPv4 || a == AddressTypeIPv6
}

type NamedServiceTemplateSpec struct {
	// Alias represents the identifier of the service.
	Alias ServiceAlias `json:"alias"`

	// ServiceTemplate is an optional configuration for a service used to expose database
	// +optional
	ofstv1.ServiceTemplateSpec `json:",inline,omitempty"`
}

type KernelSettings struct {
	// DisableDefaults can be set to false to avoid defaulting via mutator
	DisableDefaults bool `json:"disableDefaults,omitempty"`
	// Privileged specifies the status whether the init container
	// requires privileged access to perform the following commands.
	// +optional
	Privileged bool `json:"privileged,omitempty"`
	// Sysctls hold a list of sysctls commands needs to apply to kernel.
	// +optional
	Sysctls []core.Sysctl `json:"sysctls,omitempty"`
}

// AutoOpsSpec defines the specifications of automatic ops-request recommendation generation
type AutoOpsSpec struct {
	// Disabled specifies whether the ops-request recommendation generation will be disabled or not.
	// +optional
	Disabled bool `json:"disabled,omitempty"`
}

type SystemUserSecretsSpec struct {
	// ReplicationUserSecret contains replication system user credentials
	// +optional
	ReplicationUserSecret *SecretReference `json:"replicationUserSecret,omitempty"`

	// MonitorUserSecret contains monitor system user credentials
	// +optional
	MonitorUserSecret *SecretReference `json:"monitorUserSecret,omitempty"`
}

type SecretReference struct {
	core.LocalObjectReference `json:",inline,omitempty"`
	// Recommendation engine will generate RotateAuth opsReq using this field
	// +optional
	RotateAfter *metav1.Duration `json:"rotateAfter,omitempty"`
	// ActiveFrom holds the RFC3339 time. The referred authSecret is in-use from this timestamp.
	// +optional
	ActiveFrom        *metav1.Time `json:"activeFrom,omitempty"`
	ExternallyManaged bool         `json:"externallyManaged,omitempty"`
}

type Age struct {
	// Populated by Provisioner when authSecret is created or Ops Manager when authSecret is updated.
	LastUpdateTimestamp metav1.Time `json:"lastUpdateTimestamp,omitempty"`
}

type Archiver struct {
	// Pause is used to stop the archiver backup for the database
	// +optional
	Pause bool `json:"pause,omitempty"`
	// Ref is the name and namespace reference to the Archiver CR
	Ref kmapi.ObjectReference `json:"ref"`
}

type ArchiverRecovery struct {
	RecoveryTimestamp metav1.Time `json:"recoveryTimestamp"`
	// +optional
	EncryptionSecret *kmapi.ObjectReference `json:"encryptionSecret,omitempty"`
	// +optional
	ManifestRepository *kmapi.ObjectReference `json:"manifestRepository,omitempty"`

	// FullDBRepository means db restore + manifest restore
	FullDBRepository    *kmapi.ObjectReference   `json:"fullDBRepository,omitempty"`
	ReplicationStrategy *PITRReplicationStrategy `json:"replicationStrategy,omitempty"`
}
