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
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

// ContainerControlledValues controls which resource value should be autoscaled.
// +kubebuilder:validation:Enum=RequestsAndLimits;RequestsOnly
type ContainerControlledValues string

// AutoscalerTrigger controls if autoscaler is enabled.
type AutoscalerTrigger string

type ComputeAutoscalerSpec struct {
	// Whether compute autoscaler is enabled. The default is Off".
	Trigger AutoscalerTrigger `json:"trigger,omitempty"`
	// Specifies the minimal amount of resources that will be recommended.
	// The default is no minimum.
	// +optional
	MinAllowed core.ResourceList `json:"minAllowed,omitempty"`
	// Specifies the maximum amount of resources that will be recommended.
	// The default is no maximum.
	// +optional
	MaxAllowed core.ResourceList `json:"maxAllowed,omitempty"`

	// Specifies the type of recommendations that will be computed
	// (and possibly applied) by VPA.
	// If not specified, the default of [ResourceCPU, ResourceMemory] will be used.
	// +optional
	// +patchStrategy=merge
	ControlledResources []core.ResourceName `json:"controlledResources,omitempty" patchStrategy:"merge"`

	// Specifies which resource values should be controlled.
	// The default is "RequestsAndLimits".
	// +optional
	ContainerControlledValues *ContainerControlledValues `json:"containerControlledValues,omitempty"`

	// Specifies the minimum resource difference in percentage. The default is 10%.
	// If the difference between current & recommended resource is less than ResourceDiffPercentage,
	// Autoscaler Operator will ignore the updating.
	// +optional
	ResourceDiffPercentage int32 `json:"resourceDiffPercentage,omitempty"`

	// Specifies the minimum pod life time. The default is 12h.
	// If the resource Request is inside the recommended range & there is no quickOOM (out-of-memory),
	// we can still update the pod, if that pod's lifeTime is greater than this threshold.
	// +optional
	PodLifeTimeThreshold metav1.Duration `json:"podLifeTimeThreshold,omitempty"`

	// Specifies the dbStorage scaling when db data is stored in Memory
	InMemoryStorage *ComputeInMemoryStorageSpec `json:"inMemoryStorage,omitempty"`
}

type ComputeInMemoryStorageSpec struct {
	// For InMemory storageType, if db uses more than UsageThresholdPercentage of the total memory() ,
	// memoryStorage should be increased by ScalingThreshold percent
	// Default is 70%
	// +optional
	UsageThresholdPercentage int32 `json:"usageThresholdPercentage,omitempty"`

	// For InMemory storageType, if db uses more than UsageThresholdPercentage
	// of the total memory() memoryStorage should be increased by ScalingFactor percent
	// Default is 50%
	// +optional
	ScalingFactorPercentage int32 `json:"scalingFactorPercentage,omitempty"`
}

type StorageAutoscalerSpec struct {
	// Whether storage autoscaler is enabled. The default is Off".
	Trigger AutoscalerTrigger `json:"trigger,omitempty"`

	// If PVC usage percentage is less than the UsageThreshold,
	// we don't need to scale it. The Default is 80%
	UsageThreshold int32 `json:"usageThreshold,omitempty"`

	// If PVC usage percentage >= UsageThreshold,
	// we need to scale that by ScalingThreshold percentage. The Default is 50%
	ScalingThreshold int32 `json:"scalingThreshold,omitempty"`

	// ExpansionMode can be `Online` or `Offline`
	// Default VolumeExpansionMode is `Online`
	ExpansionMode *opsapi.VolumeExpansionMode `json:"expansionMode,omitempty"`
}

// AutoscalerStatus describes the runtime state of the autoscaler.
type AutoscalerStatus struct {
	// Specifies the current phase of the autoscaler
	// +optional
	Phase AutoscalerPhase `json:"phase,omitempty"`

	// observedGeneration is the most recent generation observed by this autoscaler.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions is the set of conditions required for this autoscaler to scale its target,
	// and indicates whether or not those conditions are met.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []kmapi.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// This field is equivalent to this one:
	// https://github.com/kubernetes/autoscaler/blob/273e35b88cb50c5aac383c5eceb88fb337cb31b6/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1/types.go#L218-L230
	// +optional
	VPAs []VPAStatus `json:"vpas,omitempty"`

	// Checkpoints hold all the Checkpoint those are associated
	// with this Autoscaler object. Equivalent to :
	// https://github.com/kubernetes/autoscaler/blob/273e35b88cb50c5aac383c5eceb88fb337cb31b6/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1/types.go#L354-L378
	// +optional
	Checkpoints []Checkpoint `json:"checkpoints,omitempty"`
}

// +kubebuilder:validation:Enum=InProgress;Current;Terminating;Failed
type AutoscalerPhase string

type StatusAccessor interface {
	GetStatus() AutoscalerStatus
	SetStatus(_ AutoscalerStatus)
}
