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
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

type InitSpec struct {
	// Initialized indicates that this database has been initialized.
	// This will be set by the operator when status.conditions["Provisioned"] is set to ensure
	// that database is not mistakenly reset when recovered using disaster recovery tools.
	Initialized bool `json:"initialized,omitempty"`
	// Wait for initial DataRestore condition
	WaitForInitialRestore bool              `json:"waitForInitialRestore,omitempty"`
	Script                *ScriptSourceSpec `json:"script,omitempty"`
}

type ScriptSourceSpec struct {
	ScriptPath        string `json:"scriptPath,omitempty"`
	core.VolumeSource `json:",inline,omitempty"`
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
type TerminationPolicy string

const (
	// Deletes database pods, service but leave the PVCs and stash backup data intact.
	TerminationPolicyHalt TerminationPolicy = "Halt"
	// Deletes database pods, service, pvcs but leave the stash backup data intact.
	TerminationPolicyDelete TerminationPolicy = "Delete"
	// Deletes database pods, service, pvcs and stash backup data.
	TerminationPolicyWipeOut TerminationPolicy = "WipeOut"
	// Rejects attempt to delete database using ValidationWebhook.
	TerminationPolicyDoNotTerminate TerminationPolicy = "DoNotTerminate"
)

// +kubebuilder:validation:Enum=primary;standby;stats
type ServiceAlias string

const (
	PrimaryServiceAlias ServiceAlias = "primary"
	StandbyServiceAlias ServiceAlias = "standby"
	StatsServiceAlias   ServiceAlias = "stats"
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
	ofst.ServiceTemplateSpec `json:",inline,omitempty"`
}

type KernelSettings struct {
	// Privileged specifies the status whether the init container
	// requires privileged access to perform the following commands.
	// +optional
	Privileged bool `json:"privileged,omitempty"`
	// Sysctls hold a list of sysctls commands needs to apply to kernel.
	// +optional
	Sysctls []core.Sysctl `json:"sysctls,omitempty"`
}

// CoordinatorSpec defines attributes of the coordinator container
type CoordinatorSpec struct {
	// Compute Resources required by coordinator container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`

	// Security options the coordinator container should run with.
	// More info: https://kubernetes.io/docs/concepts/policy/security-context/
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +optional
	SecurityContext *core.SecurityContext `json:"securityContext,omitempty"`
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
	ExternallyManaged         bool `json:"externallyManaged,omitempty"`
}

type Age struct {
	// Populated by Provisioner when authSecret is created or Ops Manager when authSecret is updated.
	LastUpdateTimestamp metav1.Time `json:"lastUpdateTimestamp,omitempty"`
}
