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
	"fmt"
	"regexp"
	"strconv"
	"time"

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

	// RetentionPeriod is the retention policy to be used for Logs (i.e. '60d') means how long logs will be retained before being pruned.
	// The retention policy is expressed in the form of `XXu` where `XX` is a positive integer and `u` is in `[dwm]` - days, weeks, months, years.
	// time.RFC3339 We need to parse the time to RFC3339 format
	// +kubebuilder:validation:Pattern=^[1-9][0-9]*[dwmy]$
	// +kubebuilder:default="1y"
	// +optional
	RetentionPeriod string `json:"retentionPeriod,omitempty"`

	// RetentionSchedule defines the cron expression when the log retention (pruning) task will run.
	// Cron format, e.g. "0 0 1 * *" (monthly on the 1st at 12).
	// +kubebuilder:default="0 0 1 * *"
	// +optional
	RetentionSchedule string `json:"retentionSchedule,omitempty"`

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

	// LogRetentionHistoryLimit defines the number of retention status the incremental snapshot will retain for debugging purposes.
	// The default value is 5.
	// +kubebuilder:default=5
	// +optional
	LogRetentionHistoryLimit int32 `json:"logRetentionHistoryLimit,omitempty"`
}

func ParseCutoffTimeFromPeriod(period string, now time.Time) (time.Time, error) {
	regexPolicy := regexp.MustCompile(`^([1-9][0-9]*)([dwmy])$`)
	unitFunc := map[string]func(int) time.Time{
		"d": func(v int) time.Time { return now.AddDate(0, 0, -v) },
		"w": func(v int) time.Time { return now.AddDate(0, 0, -v*7) },
		"m": func(v int) time.Time { return now.AddDate(0, -v, 0) },
		"y": func(v int) time.Time { return now.AddDate(-v, 0, 0) },
	}
	matches := regexPolicy.FindStringSubmatch(period)
	if len(matches) < 3 {
		return time.Time{}, fmt.Errorf("not a valid period")
	}
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid numeric value: %w", err)
	}
	return unitFunc[matches[2]](value), nil
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
