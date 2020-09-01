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
	"k8s.io/api/autoscaling/v2beta2"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceCodeMemcachedAutoscaler     = "mcscaler"
	ResourceKindMemcachedAutoscaler     = "MemcachedAutoscaler"
	ResourceSingularMemcachedAutoscaler = "memcachedautoscaler"
	ResourcePluralMemcachedAutoscaler   = "memcachedautoscalers"
)

// MemcachedAutoscaler is the configuration for a horizontal pod
// autoscaler, which automatically manages the replica count of any resource
// implementing the scale subresource based on the metrics specified.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=memcachedautoscalers,singular=memcachedautoscaler,shortName=mcscaler,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
type MemcachedAutoscaler struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// spec is the specification for the behaviour of the autoscaler.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status.
	// +optional
	Spec MemcachedAutoscalerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// status is the current information about the autoscaler.
	// +optional
	Status MemcachedAutoscalerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MemcachedAutoscalerSpec describes the desired functionality of the MemcachedAutoscaler.
type MemcachedAutoscalerSpec struct {
	// scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics
	// should be collected, as well as to actually change the replica count.
	ScaleTargetRef core.LocalObjectReference `json:"scaleTargetRef" protobuf:"bytes,1,opt,name=scaleTargetRef"`
	// minReplicas is the lower limit for the number of replicas to which the autoscaler
	// can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the
	// alpha feature gate HPAScaleToZero is enabled and at least one Object or External
	// metric is configured.  Scaling is active as long as at least one metric value is
	// available.
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty" protobuf:"varint,2,opt,name=minReplicas"`
	// maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// It cannot be less that minReplicas.
	MaxReplicas int32 `json:"maxReplicas" protobuf:"varint,3,opt,name=maxReplicas"`
	// metrics contains the specifications for which to use to calculate the
	// desired replica count (the maximum replica count across all metrics will
	// be used).  The desired replica count is calculated multiplying the
	// ratio between the target value and the current value by the current
	// number of pods.  Ergo, metrics used must decrease as the pod count is
	// increased, and vice-versa.  See the individual metric source types for
	// more information about how each type of metric must respond.
	// If not set, the default metric will be set to 80% average CPU utilization.
	// +optional
	Metrics []v2beta2.MetricSpec `json:"metrics,omitempty" protobuf:"bytes,4,rep,name=metrics"`

	// behavior configures the scaling behavior of the target
	// in both Up and Down directions (scaleUp and scaleDown fields respectively).
	// If not set, the default MemcachedScalingRules for scale up and scale down are used.
	// +optional
	Behavior *MemcachedAutoscalerBehavior `json:"behavior,omitempty" protobuf:"bytes,5,opt,name=behavior"`
}

// MemcachedAutoscalerBehavior configures the scaling behavior of the target
// in both Up and Down directions (scaleUp and scaleDown fields respectively).
type MemcachedAutoscalerBehavior struct {
	// scaleUp is scaling policy for scaling Up.
	// If not set, the default value is the higher of:
	//   * increase no more than 4 pods per 60 seconds
	//   * double the number of pods per 60 seconds
	// No stabilization is used.
	// +optional
	ScaleUp *MemcachedScalingRules `json:"scaleUp,omitempty" protobuf:"bytes,1,opt,name=scaleUp"`
	// scaleDown is scaling policy for scaling Down.
	// If not set, the default value is to allow to scale down to minReplicas pods, with a
	// 300 second stabilization window (i.e., the highest recommendation for
	// the last 300sec is used).
	// +optional
	ScaleDown *MemcachedScalingRules `json:"scaleDown,omitempty" protobuf:"bytes,2,opt,name=scaleDown"`
}

// MemcachedScalingRules configures the scaling behavior for one direction.
// These Rules are applied after calculating DesiredReplicas from metrics for the HPA.
// They can limit the scaling velocity by specifying scaling policies.
// They can prevent flapping by specifying the stabilization window, so that the
// number of replicas is not set instantly, instead, the safest value from the stabilization
// window is chosen.
type MemcachedScalingRules struct {
	// StabilizationWindowSeconds is the number of seconds for which past recommendations should be
	// considered while scaling up or scaling down.
	// StabilizationWindowSeconds must be greater than or equal to zero and less than or equal to 3600 (one hour).
	// If not set, use the default values:
	// - For scale up: 0 (i.e. no stabilization is done).
	// - For scale down: 300 (i.e. the stabilization window is 300 seconds long).
	// +optional
	StabilizationWindowSeconds *int32 `json:"stabilizationWindowSeconds" protobuf:"varint,3,opt,name=stabilizationWindowSeconds"`
	// selectPolicy is used to specify which policy should be used.
	// If not set, the default value MaxPolicySelect is used.
	// +optional
	SelectPolicy *v2beta2.ScalingPolicySelect `json:"selectPolicy,omitempty" protobuf:"bytes,1,opt,name=selectPolicy"`
	// policies is a list of potential scaling polices which can be used during scaling.
	// At least one policy must be specified, otherwise the MemcachedScalingRules will be discarded as invalid
	// +optional
	Policies []MemcachedScalingPolicy `json:"policies,omitempty" protobuf:"bytes,2,rep,name=policies"`
}

// MemcachedScalingPolicy is a single policy which must hold true for a specified past interval.
type MemcachedScalingPolicy struct {
	// Type is used to specify the scaling policy.
	Type v2beta2.HPAScalingPolicyType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=HPAScalingPolicyType"`
	// Value contains the amount of change which is permitted by the policy.
	// It must be greater than zero
	Value int32 `json:"value" protobuf:"varint,2,opt,name=value"`
	// PeriodSeconds specifies the window of time for which the policy should hold true.
	// PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
	PeriodSeconds int32 `json:"periodSeconds" protobuf:"varint,3,opt,name=periodSeconds"`
}

// MemcachedAutoscalerStatus describes the current status of a horizontal pod autoscaler.
type MemcachedAutoscalerStatus struct {
	// observedGeneration is the most recent generation observed by this autoscaler.
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// lastScaleTime is the last time the MemcachedAutoscaler scaled the number of pods,
	// used by the autoscaler to control how often the number of pods is changed.
	// +optional
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty" protobuf:"bytes,2,opt,name=lastScaleTime"`

	// currentReplicas is current number of replicas of pods managed by this autoscaler,
	// as last seen by the autoscaler.
	CurrentReplicas int32 `json:"currentReplicas" protobuf:"varint,3,opt,name=currentReplicas"`

	// desiredReplicas is the desired number of replicas of pods managed by this autoscaler,
	// as last calculated by the autoscaler.
	DesiredReplicas int32 `json:"desiredReplicas" protobuf:"varint,4,opt,name=desiredReplicas"`

	// currentMetrics is the last read state of the metrics used by this autoscaler.
	// +optional
	CurrentMetrics []v2beta2.MetricStatus `json:"currentMetrics" protobuf:"bytes,5,rep,name=currentMetrics"`

	// conditions is the set of conditions required for this autoscaler to scale its target,
	// and indicates whether or not those conditions are met.
	Conditions []kmapi.Condition `json:"conditions" protobuf:"bytes,6,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MemcachedAutoscalerList is a list of horizontal pod autoscaler objects.
type MemcachedAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// items is the list of horizontal pod autoscaler objects.
	Items []MemcachedAutoscaler `json:"items" protobuf:"bytes,2,rep,name=items"`
}
