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
	"kubestash.dev/apimachinery/apis"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

// AddonInfo specifies addon configuration that will be used to backup/restore the respective target.
type AddonInfo struct {
	// Name specifies the name of the addon that will be used for the backup/restore purpose
	Name string `json:"name,omitempty"`

	// Tasks specifies a list of backup/restore tasks and their configuration parameters
	Tasks []TaskReference `json:"tasks,omitempty"`

	// ContainerRuntimeSettings specifies runtime settings for the backup/restore executor container
	// +optional
	ContainerRuntimeSettings *ofst.ContainerRuntimeSettings `json:"containerRuntimeSettings,omitempty"`

	// JobTemplate specifies runtime configurations for the backup/restore Job
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
}

// TaskReference specifies a task and its configuration parameters
type TaskReference struct {
	// Name indicates to the name of the task
	Name string `json:"name,omitempty"`

	// Variables specifies a list of variables and their sources that will be used to resolve the task.
	// +optional
	Variables []core.EnvVar `json:"variables,omitempty"`

	// Params specifies parameters for the task. You must provide the parameter in the Addon desired structure.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Params *runtime.RawExtension `json:"params,omitempty"`

	// TargetVolumes specifies which volumes from the target should be mounted in the backup/restore job/container.
	// +optional
	TargetVolumes *TargetVolumeInfo `json:"targetVolumes,omitempty"`

	// AddonVolumes lets you overwrite the volume sources used in the VolumeTemplate section of Addon.
	// Make sure that name of your volume matches with the name of the volume you want to overwrite.
	// +optional
	AddonVolumes []AddonVolumeInfo `json:"addonVolumes,omitempty"`
}

// TargetVolumeInfo specifies the volumes and their mounts of the targeted application that should
// be mounted in backup/restore Job/container.
type TargetVolumeInfo struct {
	// Volumes indicates the list of volumes of targeted application that should be mounted on the backup/restore job.
	Volumes []ofst.Volume `json:"volumes,omitempty"`

	// VolumeMounts specifies the mount for the volumes specified in `Volumes` section
	VolumeMounts []core.VolumeMount `json:"volumeMounts,omitempty"`

	// VolumeClaimTemplates specifies a template for the PersistentVolumeClaims that will be created for each Pod in a StatefulSet.
	VolumeClaimTemplates []ofst.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
}

// AddonVolumeInfo specifies the name and the source of volume
type AddonVolumeInfo struct {
	// Name specifies the name of the volume
	Name string `json:"name,omitempty"`

	// Source specifies the source of this volume.
	Source *apis.VolumeSource `json:"source,omitempty"`
}

// HookInfo specifies the information about the backup/restore hooks
type HookInfo struct {
	// Name specifies a name for the hook
	Name string `json:"name,omitempty"`

	// HookTemplate points to a HookTemplate CR that will be used to execute the hook.
	// You can refer to a HookTemplate from other namespaces as long as your current
	// namespace is allowed by the `usagePolicy` in the respective HookTemplate.
	HookTemplate *kmapi.ObjectReference `json:"hookTemplate,omitempty"`

	// Params specifies parameters for the hook. You must provide the parameter in the HookTemplates desired structure.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Params *runtime.RawExtension `json:"params,omitempty"`

	// MaxRetry specifies how many times KubeStash should retry the hook execution in case of failure.
	// The default value of this field is 0 which means no retry.
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxRetry int32 `json:"maxRetry,omitempty"`

	// Timeout specifies a duration in seconds that KubeStash should wait for the hook execution to be completed.
	// If the hook execution does not finish within this time period, KubeStash will consider this hook execution as failure.
	// Then, it will be re-tried according to MaxRetry policy.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// ExecutionPolicy specifies when to execute the hook.
	// Valid values are:
	// - "Always": KubeStash will execute this hook no matter the backup/restore failed. This is the default execution policy.
	// - "OnSuccess": KubeStash will execute this hook only if the backup/restore has succeeded.
	// - "OnFailure": KubeStash will execute this hook only if the backup/restore has failed.
	// +kubebuilder:default=Always
	// +optional
	ExecutionPolicy HookExecutionPolicy `json:"executionPolicy,omitempty"`

	// Variables specifies a list of variables and their sources that will be used to resolve the HookTemplate.
	// +optional
	Variables []core.EnvVar `json:"variables,omitempty"`

	// Volumes indicates the list of volumes of targeted application that should be mounted on the hook executor.
	// Use this field only for `Function` type hook executor.
	// +optional
	Volumes []ofst.Volume `json:"volumes,omitempty"`

	// VolumeMounts specifies the mount for the volumes specified in `Volumes` section
	// Use this field only for `Function` type hook executor.
	// +optional
	VolumeMounts []core.VolumeMount `json:"volumeMounts,omitempty"`

	// RuntimeSettings specifies runtime configurations for the hook executor Job.
	// Use this field only for `Function` type hook executor.
	// +optional
	RuntimeSettings *ofst.RuntimeSettings `json:"runtimeSettings,omitempty"`
}

// HookStatus represents the status of the hooks
type HookStatus struct {
	// PreHooks represents the pre-restore hook execution status
	// +optional
	PreHooks []HookExecutionStatus `json:"preHooks,omitempty"`

	// PostHooks represents the post-restore hook execution status
	// +optional
	PostHooks []HookExecutionStatus `json:"postHooks,omitempty"`
}

// HookExecutionPolicy specifies when to execute the hook.
// +kubebuilder:validation:Enum=Always;OnSuccess;OnFailure
type HookExecutionPolicy string

const (
	ExecuteAlways    HookExecutionPolicy = "Always"
	ExecuteOnSuccess HookExecutionPolicy = "OnSuccess"
	ExecuteOnFailure HookExecutionPolicy = "OnFailure"
)

// HookExecutionStatus represents the state of the hook execution
type HookExecutionStatus struct {
	// Name indicates the name of the hook whose status is being shown here.
	Name string `json:"name,omitempty"`

	// Phase represents the hook execution phase
	// +optional
	Phase HookExecutionPhase `json:"phase,omitempty"`
}

// HookExecutionPhase specifies the state of the hook execution
// +kubebuilder:validation:Enum=Succeeded;Failed;Pending
type HookExecutionPhase string

const (
	HookExecutionSucceeded HookExecutionPhase = "Succeeded"
	HookExecutionFailed    HookExecutionPhase = "Failed"
	HookExecutionPending   HookExecutionPhase = "Pending"
)

// ResourceFoundStatus specifies whether a resource was found or not
type ResourceFoundStatus struct {
	kmapi.TypedObjectReference `json:",inline"`
	// Found indicates whether the resource was found or not
	Found *bool `json:"found,omitempty"`
}

// FailurePolicy specifies what to do if a backup/restore fails
// +kubebuilder:validation:Enum=Fail;Retry
type FailurePolicy string

const (
	FailurePolicyFail  FailurePolicy = "Fail"
	FailurePolicyRetry FailurePolicy = "Retry"
)

// RetryConfig specifies the behavior of retry
type RetryConfig struct {
	// MaxRetry specifies the maximum number of times KubeStash should retry the backup/restore process.
	// By default, KubeStash will retry only 1 time.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	MaxRetry int32 `json:"maxRetry,omitempty"`

	// The amount of time to wait before next retry. If you don't specify this field, KubeStash will retry immediately.
	// Format: 30s, 2m, 1h etc.
	// +optional
	Delay metav1.Duration `json:"delay,omitempty"`
}

const (
	// TypeMetricsPushed indicates whether Metrics are pushed or not
	TypeMetricsPushed               = "MetricsPushed"
	ReasonSuccessfullyPushedMetrics = "SuccessfullyPushedMetrics"
	ReasonFailedToPushMetrics       = "FailedToPushMetrics"
)
