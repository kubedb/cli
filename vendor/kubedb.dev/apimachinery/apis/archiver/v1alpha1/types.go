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
	api "kubedb.dev/apimachinery/apis/kubedb/v1"

	batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	"kubestash.dev/apimachinery/apis"
	stashcoreapi "kubestash.dev/apimachinery/apis/core/v1alpha1"
)

type Accessor interface {
	GetObjectMeta() metav1.ObjectMeta
	GetConsumers() *api.AllowedConsumers
}

type ListAccessor interface {
	GetItems() []Accessor
}

type FullBackupOptions struct {
	// +kubebuilder:default:=VolumeSnapshotter
	Driver apis.Driver `json:"driver"`
	// +optional
	Task *Task `json:"task,omitempty"`
	// +optional
	Scheduler *SchedulerOptions `json:"scheduler,omitempty"`
	// +optional
	ContainerRuntimeSettings *ofst.ContainerRuntimeSettings `json:"containerRuntimeSettings,omitempty"`
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
	// +optional
	RetryConfig *stashcoreapi.RetryConfig `json:"retryConfig,omitempty"`
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// +optional
	SessionHistoryLimit int32 `json:"sessionHistoryLimit,omitempty"`
}

type ManifestBackupOptions struct {
	// +optional
	Scheduler *SchedulerOptions `json:"scheduler,omitempty"`
	// +optional
	ContainerRuntimeSettings *ofst.ContainerRuntimeSettings `json:"containerRuntimeSettings,omitempty"`
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
	// +optional
	RetryConfig *stashcoreapi.RetryConfig `json:"retryConfig,omitempty"`
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// +optional
	SessionHistoryLimit int32 `json:"sessionHistoryLimit,omitempty"`
}

type LogBackupOptions struct {
	// +optional
	RuntimeSettings *ofst.RuntimeSettings `json:"runtimeSettings,omitempty"`

	// +optional
	ConfigSecret *GenericSecretReference `json:"configSecret,omitempty"`

	// SuccessfulLogHistoryLimit defines the number of successful Logs backup status that the incremental snapshot will retain
	// The default value is 5.
	// +kubebuilder:default=5
	// +optional
	SuccessfulLogHistoryLimit int32 `json:"successfulLogHistoryLimit,omitempty"`

	// FailedLogHistoryLimit defines the number of failed Logs backup that the incremental snapshot will retain for debugging purposes.
	// The default value is 5.
	// +kubebuilder:default=5
	// +optional
	FailedLogHistoryLimit int32 `json:"failedLogHistoryLimit,omitempty"`
}

type Task struct {
	Params *runtime.RawExtension `json:"params"`
}

type BackupStorage struct {
	Ref *kmapi.ObjectReference `json:"ref,omitempty"`
	// +optional
	SubDir string `json:"subDir,omitempty"`
}

type SchedulerOptions struct {
	Schedule string `json:"schedule"`
	// +optional
	ConcurrencyPolicy batch.ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`
	// +optional
	JobTemplate stashcoreapi.JobTemplate `json:"jobTemplate"`
	// +optional
	SuccessfulJobsHistoryLimit *int32 `json:"successfulJobsHistoryLimit,omitempty"`
	// +optional
	FailedJobsHistoryLimit *int32 `json:"failedJobsHistoryLimit,omitempty"`
}

type ArchiverDatabaseRef struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type GenericSecretReference struct {
	// Name of the provider secret
	Name string `json:"name"`

	EnvToSecretKey map[string]string `json:"envToSecretKey"`
}
