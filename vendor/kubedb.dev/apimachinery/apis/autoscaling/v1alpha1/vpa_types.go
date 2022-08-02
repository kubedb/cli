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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VerticalPodAutopilotRecommenderSelector points to a specific Vertical Pod autopilot recommender.
// In the future it might pass parameters to the recommender.
type VerticalPodAutopilotRecommenderSelector struct {
	// Name of the recommender responsible for generating recommendation for this object.
	Name string `json:"name"`
}

// PodUpdatePolicy describes the rules on how changes are applied to the pods.
type PodUpdatePolicy struct {
	// Controls when autopilot applies changes to the pod resources.
	// The default is 'Auto'.
	// +optional
	UpdateMode *UpdateMode `json:"updateMode,omitempty"`

	// Minimal number of replicas which need to be alive for Updater to attempt
	// pod eviction (pending other checks like PDB). Only positive values are
	// allowed. Overrides global '--min-replicas' flag.
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`
}

// UpdateMode controls when autopilot applies changes to the pod resoures.
// +kubebuilder:validation:Enum=Off;Initial;Recreate;Auto
type UpdateMode string

const (
	// UpdateModeOff means that autopilot never changes Pod resources.
	// The recommender still sets the recommended resources in the
	// VerticalPodAutopilot object. This can be used for a "dry run".
	UpdateModeOff UpdateMode = "Off"
	// UpdateModeInitial means that autopilot only assigns resources on pod
	// creation and does not change them during the lifetime of the pod.
	UpdateModeInitial UpdateMode = "Initial"
	// UpdateModeRecreate means that autopilot assigns resources on pod
	// creation and additionally can update them during the lifetime of the
	// pod by deleting and recreating the pod.
	UpdateModeRecreate UpdateMode = "Recreate"
	// UpdateModeAuto means that autopilot assigns resources on pod creation
	// and additionally can update them during the lifetime of the pod,
	// using any available update method. Currently this is equivalent to
	// Recreate, which is the only available update method.
	UpdateModeAuto UpdateMode = "Auto"
)

// PodResourcePolicy controls how autopilot computes the recommended resources
// for containers belonging to the pod. There can be at most one entry for every
// named container and optionally a single wildcard entry with `containerName` = '*',
// which handles all containers that don't have individual policies.
type PodResourcePolicy struct {
	// Per-container resource policies.
	// +optional
	// +patchMergeKey=containerName
	// +patchStrategy=merge
	ContainerPolicies []ContainerResourcePolicy `json:"containerPolicies,omitempty" patchStrategy:"merge" patchMergeKey:"containerName"`
}

// ContainerResourcePolicy controls how autopilot computes the recommended
// resources for a specific container.
type ContainerResourcePolicy struct {
	// Name of the container or DefaultContainerResourcePolicy, in which
	// case the policy is used by the containers that don't have their own
	// policy specified.
	ContainerName string `json:"containerName,omitempty"`
	// Whether autopilot is enabled for the container. The default is "Auto".
	// +optional
	Mode *ContainerScalingMode `json:"mode,omitempty"`
}

const (
	// DefaultContainerResourcePolicy can be passed as
	// ContainerResourcePolicy.ContainerName to specify the default policy.
	DefaultContainerResourcePolicy = "*"
)

// ContainerScalingMode controls whether autopilot is enabled for a specific
// container.
// +kubebuilder:validation:Enum=Auto;Off
type ContainerScalingMode string

const (
	// ContainerScalingModeAuto means autopilot is enabled for a container.
	ContainerScalingModeAuto ContainerScalingMode = "Auto"
	// ContainerScalingModeOff means autopilot is disabled for a container.
	ContainerScalingModeOff ContainerScalingMode = "Off"
)

// VPAStatus describes the runtime state of the autopilot.
type VPAStatus struct {
	// The name of the VerticalPodAutoscaler.
	// This field will help us to get the corresponding vpa.
	VPAName string `json:"vpaName,omitempty"`

	// The most recently computed amount of resources recommended by the
	// autopilot for the controlled pods.
	// +optional
	Recommendation *RecommendedPodResources `json:"recommendation,omitempty"`

	// Conditions is the set of conditions required for this autopilot to scale its target,
	// and indicates whether or not those conditions are met.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []VerticalPodAutopilotCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// RecommendedPodResources is the recommendation of resources computed by
// autopilot. It contains a recommendation for each container in the pod
// (except for those with `ContainerScalingMode` set to 'Off').
type RecommendedPodResources struct {
	// Resources recommended by the autopilot for each container.
	// +optional
	ContainerRecommendations []RecommendedContainerResources `json:"containerRecommendations,omitempty"`
}

// RecommendedContainerResources is the recommendation of resources computed by
// autopilot for a specific container. Respects the container resource policy
// if present in the spec. In particular the recommendation is not produced for
// containers with `ContainerScalingMode` set to 'Off'.
type RecommendedContainerResources struct {
	// Name of the container.
	ContainerName string `json:"containerName,omitempty" `
	// Recommended amount of resources. Observes ContainerResourcePolicy.
	Target v1.ResourceList `json:"target" `
	// Minimum recommended amount of resources. Observes ContainerResourcePolicy.
	// This amount is not guaranteed to be sufficient for the application to operate in a stable way, however
	// running with less resources is likely to have significant impact on performance/availability.
	// +optional
	LowerBound v1.ResourceList `json:"lowerBound,omitempty"`
	// Maximum recommended amount of resources. Observes ContainerResourcePolicy.
	// Any resources allocated beyond this value are likely wasted. This value may be larger than the maximum
	// amount of application is actually capable of consuming.
	// +optional
	UpperBound v1.ResourceList `json:"upperBound,omitempty"`
	// The most recent recommended resources target computed by the autopilot
	// for the controlled pods, based only on actual resource usage, not taking
	// into account the ContainerResourcePolicy.
	// May differ from the Recommendation if the actual resource usage causes
	// the target to violate the ContainerResourcePolicy (lower than MinAllowed
	// or higher that MaxAllowed).
	// Used only as status indication, will not affect actual resource assignment.
	// +optional
	UncappedTarget v1.ResourceList `json:"uncappedTarget,omitempty"`
}

// VerticalPodAutopilotConditionType are the valid conditions of
// a VerticalPodAutopilot.
type VerticalPodAutopilotConditionType string

var (
	// RecommendationProvided indicates whether the VPA recommender was able to calculate a recommendation.
	RecommendationProvided VerticalPodAutopilotConditionType = "RecommendationProvided"
	// LowConfidence indicates whether the VPA recommender has low confidence in the recommendation for
	// some of containers.
	LowConfidence VerticalPodAutopilotConditionType = "LowConfidence"
	// NoPodsMatched indicates that label selector used with VPA object didn't match any pods.
	NoPodsMatched VerticalPodAutopilotConditionType = "NoPodsMatched"
	// FetchingHistory indicates that VPA recommender is in the process of loading additional history samples.
	FetchingHistory VerticalPodAutopilotConditionType = "FetchingHistory"
	// ConfigDeprecated indicates that this VPA configuration is deprecated
	// and will stop being supported soon.
	ConfigDeprecated VerticalPodAutopilotConditionType = "ConfigDeprecated"
	// ConfigUnsupported indicates that this VPA configuration is unsupported
	// and recommendations will not be provided for it.
	ConfigUnsupported VerticalPodAutopilotConditionType = "ConfigUnsupported"
)

// VerticalPodAutopilotCondition describes the state of
// a VerticalPodAutopilot at a certain point.
type VerticalPodAutopilotCondition struct {
	// type describes the current condition
	Type VerticalPodAutopilotConditionType `json:"type" `
	// status is the status of the condition (True, False, Unknown)
	Status v1.ConditionStatus `json:"status"`
	// lastTransitionTime is the last time the condition transitioned from
	// one status to another
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// reason is the reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// message is a human-readable explanation containing details about
	// the transition
	// +optional
	Message string `json:"message,omitempty"`
}
