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
	nodemeta "kmodules.xyz/resource-metadata/apis/node/v1alpha1"
)

type OpsRequestStatus struct {
	// Specifies the current phase of the ops request
	// +optional
	Phase OpsRequestPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
	// PausedBackups represents the list of backups that have been paused.
	// +optional
	PausedBackups []kmapi.TypedObjectReference `json:"pausedBackups,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Progressing;Successful;WaitingForApproval;Failed;Approved;Denied;Skipped
type OpsRequestPhase string

const (
	// used for ops requests that are currently in queue
	OpsRequestPhasePending OpsRequestPhase = "Pending"
	// used for ops requests that are currently Progressing
	OpsRequestPhaseProgressing OpsRequestPhase = "Progressing"
	// used for ops requests that are executed successfully
	OpsRequestPhaseSuccessful OpsRequestPhase = "Successful"
	// used for ops requests that are failed
	OpsRequestPhaseFailed OpsRequestPhase = "Failed"
	// used for ops requests that are skipped
	OpsRequestPhaseSkipped OpsRequestPhase = "Skipped"

	// Approval-related Phases

	// used for ops requests that are waiting for approval
	OpsRequestPhaseWaitingForApproval OpsRequestPhase = "WaitingForApproval"
	// used for ops requests that are approved
	OpsRequestApproved OpsRequestPhase = "Approved"
	// used for ops requests that are denied
	OpsRequestDenied OpsRequestPhase = "Denied"
)

// +kubebuilder:validation:Enum=Offline;Online
type VolumeExpansionMode string

const (
	// used to define a Online volume expansion mode
	VolumeExpansionModeOnline VolumeExpansionMode = "Online"
	// used to define a Offline volume expansion mode
	VolumeExpansionModeOffline VolumeExpansionMode = "Offline"
)

type RestartSpec struct{}

type Reprovision struct{}

type TLSSpec struct {
	// TLSConfig contains updated tls configurations for client and server.
	// +optional
	kmapi.TLSConfig `json:",inline,omitempty"`

	// RotateCertificates tells operator to initiate certificate rotation
	// +optional
	RotateCertificates bool `json:"rotateCertificates,omitempty"`

	// Remove tells operator to remove TLS configuration
	// +optional
	Remove bool `json:"remove,omitempty"`
}

// +kubebuilder:validation:Enum=IfReady;Always
type ApplyOption string

const (
	ApplyOptionIfReady ApplyOption = "IfReady"
	ApplyOptionAlways  ApplyOption = "Always"
)

type Accessor interface {
	GetObjectMeta() metav1.ObjectMeta
	GetDBRefName() string
	GetRequestType() any
	GetStatus() OpsRequestStatus
	SetStatus(_ OpsRequestStatus)
}

// +kubebuilder:validation:Enum=ConfigureArchiver;DisableArchiver
type ArchiverOperation string

const (
	ArchiverOperationConfigure ArchiverOperation = "ConfigureArchiver"
	ArchiverOperationDisable   ArchiverOperation = "DisableArchiver"
)

type ArchiverOptions struct {
	Operation ArchiverOperation     `json:"operation"`
	Ref       kmapi.ObjectReference `json:"ref"`
}

// ContainerResources is the spec for vertical scaling of containers
type ContainerResources struct {
	// Compute Resources required by the sidecar container.
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`
}

// PodResources is the spec for vertical scaling of pods
type PodResources struct {
	// +optional
	NodeSelectionPolicy nodemeta.NodeSelectionPolicy `json:"nodeSelectionPolicy,omitempty"`
	Topology            *Topology                    `json:"topology,omitempty"`
	// Compute Resources required by the sidecar container.
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`
}

// Topology is the spec for placement of pods onto nodes
type Topology struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
