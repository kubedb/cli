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
	"kubestash.dev/apimachinery/apis/storage/v1alpha1"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindBackupConfiguration     = "BackupConfiguration"
	ResourceSingularBackupConfiguration = "backupconfiguration"
	ResourcePluralBackupConfiguration   = "backupconfigurations"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=backupconfigurations,singular=backupconfiguration,shortName=bc,categories={kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Paused",type="boolean",JSONPath=".spec.paused"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// BackupConfiguration specifies the configuration for taking backup of a target application.
type BackupConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupConfigurationSpec   `json:"spec,omitempty"`
	Status BackupConfigurationStatus `json:"status,omitempty"`
}

// BackupConfigurationSpec defines the target of backup, the backends where the data will be stored,
// and the sessions that specifies when and how to take backup.
type BackupConfigurationSpec struct {
	// Target refers to the target of backup. The target must be in the same namespace as the BackupConfiguration.
	Target *kmapi.TypedObjectReference `json:"target,omitempty"`

	// Backends specifies a list of storage references where the backed up data will be stored.
	// The respective BackupStorages can be in a different namespace than the BackupConfiguration.
	// However, it must be allowed by the `usagePolicy` of the BackupStorage to refer from this namespace.
	//
	// This field is optional, if you don't provide any backend here, KubeStash will use the default BackupStorage for the namespace.
	// If a default BackupStorage does not exist in the same namespace, then KubeStash will look for a default BackupStorage
	// in other namespaces that allows using it from the BackupConfiguration namespace.
	// +optional
	Backends []BackendReference `json:"backends,omitempty"`

	// Sessions defines a list of session configuration that specifies when and how to take backup.
	Sessions []Session `json:"sessions,omitempty"`

	// Paused indicates that the BackupConfiguration has been paused from taking backup. Default value is 'false'.
	// If you set `paused` field to `true`, KubeStash will suspend the respective backup triggering CronJob and
	// skip processing any further events for this BackupConfiguration.
	// +optional
	Paused bool `json:"paused,omitempty"`
}

// BackendReference specifies reference to a storage where the backed up data will be stored.
type BackendReference struct {
	// Name provides an identifier for this storage.
	Name string `json:"name,omitempty"`

	// StorageRef refers to the CR that holds the information of a storage.
	// You can refer to the BackupStorage CR of a different namespace as long as it is allowed
	// by the `usagePolicy` of the BackupStorage.`
	StorageRef *kmapi.ObjectReference `json:"storageRef,omitempty"`

	// RetentionPolicy refers to a RetentionPolicy CRs which defines how to cleanup the old Snapshots.
	// This field is optional. If you don't provide this field, KubeStash will use the default RetentionPolicy for
	// the namespace. If there is no default RetentionPolicy for the namespace, then KubeStash will find a
	// RetentionPolicy from other namespaces that is allowed to use from the current namespace.
	// +optional
	RetentionPolicy *kmapi.ObjectReference `json:"retentionPolicy,omitempty"`
}

// Session specifies a backup session configuration for the target
type Session struct {
	*SessionConfig `json:",inline"`

	// Addon specifies addon configuration that will be used to backup the target.
	Addon *AddonInfo `json:"addon,omitempty"`

	// Repositories specifies a list of repository information where the backed up data will be stored.
	// KubeStash will create the respective Repository CRs using this information.
	Repositories []RepositoryInfo `json:"repositories,omitempty"`
}

// SessionConfig specifies common session configurations
type SessionConfig struct {
	// Name specifies the name of the session
	Name string `json:"name,omitempty"`

	// Scheduler specifies the configuration for backup triggering CronJob
	Scheduler *SchedulerSpec `json:"scheduler,omitempty"`

	// Hooks specifies the backup hooks that should be executed before and/or after the backup.
	// +optional
	Hooks *BackupHooks `json:"hooks,omitempty"`

	// FailurePolicy specifies what to do if the backup fail.
	// Valid values are:
	// - "Fail": KubeStash should mark the backup as failed if any component fail to complete its backup. This is the default behavior.
	// - "Retry": KubeStash will retry to backup the failed component according to the `retryConfig`.
	// +optional
	// FailurePolicy FailurePolicy `json:"failurePolicy,omitempty"`

	// RetryConfig specifies the behavior of retry in case of a backup failure.
	// +optional
	RetryConfig *RetryConfig `json:"retryConfig,omitempty"`

	// BackupTimeout specifies the maximum duration of backup. Backup will be considered Failed
	// if backup tasks do not complete within this time limit. By default, KubeStash don't set any timeout for backup.
	// +optional
	BackupTimeout *metav1.Duration `json:"backupTimeout,omitempty"`

	// SessionHistoryLimit specifies how many backup Jobs and associate resources KubeStash should keep for debugging purpose.
	// The default value is 1.
	// +kubebuilder:default=1
	// +optional
	SessionHistoryLimit int32 `json:"sessionHistoryLimit,omitempty"`
}

// SchedulerSpec specifies the configuration for the backup triggering CronJob for a session.
type SchedulerSpec struct {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	Schedule string `json:"schedule"`

	// Optional deadline in seconds for starting the job if it misses scheduled
	// time for any reason.  Missed jobs executions will be counted as failed ones.
	// +optional
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`

	// Specifies how to treat concurrent executions of a Job.
	// Valid values are:
	// - "Allow" (default): allows CronJobs to run concurrently;
	// - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet;
	// - "Replace": cancels currently running job and replaces it with a new one
	// +optional
	ConcurrencyPolicy batchv1.ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`

	// This flag tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// Specifies the job that will be created when executing a CronJob.
	JobTemplate JobTemplate `json:"jobTemplate"`

	// The number of successful finished jobs to retain. Value must be non-negative integer.
	// Defaults to 3.
	// +optional
	SuccessfulJobsHistoryLimit *int32 `json:"successfulJobsHistoryLimit,omitempty"`

	// The number of failed finished jobs to retain. Value must be non-negative integer.
	// Defaults to 1.
	// +optional
	FailedJobsHistoryLimit *int32 `json:"failedJobsHistoryLimit,omitempty"`
}

// JobTemplate specifies the template for the Job created by the scheduler CronJob.
type JobTemplate struct {
	// Specifies the maximum desired number of pods the job should
	// run at any given time. The actual number of pods running in steady state will
	// be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism),
	// i.e. when the work left to do is less than max parallelism.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	// +optional
	Parallelism *int32 `json:"parallelism,omitempty"`

	// Specifies the desired number of successfully finished pods the
	// job should be run with.  Setting to nil means that the success of any
	// pod signals the success of all pods, and allows parallelism to have any positive
	// value.  Setting to 1 means that parallelism is limited to 1 and the success of that
	// pod signals the success of the job.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	// +optional
	Completions *int32 `json:"completions,omitempty"`

	// Specifies the duration in seconds relative to the startTime that the job
	// may be continuously active before the system tries to terminate it; value
	// must be positive integer. If a Job is suspended (at creation or through an
	// update), this timer will effectively be stopped and reset when the Job is
	// resumed again.
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty"`

	// Specifies the number of retries before marking this job failed.
	// Defaults to 6
	// +optional
	BackoffLimit *int32 `json:"backoffLimit,omitempty"`

	// Describes the pod that will be created when executing a job.
	// +optional
	Template ofst.PodTemplateSpec `json:"template"`

	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished
	// execution (either Complete or Failed). If this field is set,
	// ttlSecondsAfterFinished after the Job finishes, it is eligible to be
	// automatically deleted. When the Job is being deleted, its lifecycle
	// guarantees (e.g. finalizers) will be honored. If this field is unset,
	// the Job won't be automatically deleted. If this field is set to zero,
	// the Job becomes eligible to be deleted immediately after it finishes.
	// This field is alpha-level and is only honored by servers that enable the
	// TTLAfterFinished feature.
	// +optional
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty"`

	// CompletionMode specifies how Pod completions are tracked. It can be
	// `NonIndexed` (default) or `Indexed`.
	//
	// `NonIndexed` means that the Job is considered complete when there have
	// been .spec.completions successfully completed Pods. Each Pod completion is
	// homologous to each other.
	//
	// `Indexed` means that the Pods of a
	// Job get an associated completion index from 0 to (.spec.completions - 1),
	// available in the annotation batch.kubernetes.io/job-completion-index.
	// The Job is considered complete when there is one successfully completed Pod
	// for each index.
	// When value is `Indexed`, .spec.completions must be specified and
	// `.spec.parallelism` must be less than or equal to 10^5.
	//
	// This field is alpha-level and is only honored by servers that enable the
	// IndexedJob feature gate. More completion modes can be added in the future.
	// If the Job controller observes a mode that it doesn't recognize, the
	// controller skips updates for the Job.
	// +optional
	CompletionMode *batchv1.CompletionMode `json:"completionMode,omitempty"`

	// Suspend specifies whether the Job controller should create Pods or not. If
	// a Job is created with suspend set to true, no Pods are created by the Job
	// controller. If a Job is suspended after creation (i.e. the flag goes from
	// false to true), the Job controller will delete all active Pods associated
	// with this Job. Users must design their workload to gracefully handle this.
	// Suspending a Job will reset the StartTime field of the Job, effectively
	// resetting the ActiveDeadlineSeconds timer too. This is an alpha field and
	// requires the SuspendJob feature gate to be enabled; otherwise this field
	// may not be set to true. Defaults to false.
	// +optional
	Suspend *bool `json:"suspend,omitempty"`
}

// RepositoryInfo specifies information about the repository where the backed up data will be stored.
// KubeStash will create the respective Repository CR from this information.
type RepositoryInfo struct {
	// Name specifies the name of the Repository
	Name string `json:"name,omitempty"`

	// Backend specifies the name of the backend where this repository will be initialized.
	// This should point to a backend name specified in `.spec.backends` section.
	// For using a default backend, keep this field empty.
	// +optional
	Backend string `json:"backend,omitempty"`

	// BackupVerifier specifies the name of the BackupVerifier which will be used to verify the backed up data in this repository.
	// +optional
	BackupVerifier *kmapi.ObjectReference `json:"backupVerifier,omitempty"`

	// Directory specifies the path inside the backend where the backed up data will be stored.
	Directory string `json:"directory,omitempty"`

	// EncryptionSecret refers to the Secret containing the encryption key which will be used to encode/decode the backed up dta.
	// You can refer to a Secret of a different namespace.
	// If you don't provide the namespace field, KubeStash will look for the Secret in the same namespace as the BackupConfiguration / BackupBatch.
	EncryptionSecret *kmapi.ObjectReference `json:"encryptionSecret,omitempty"`

	// DeletionPolicy specifies what to do when you delete a Repository CR.
	// +optional
	DeletionPolicy v1alpha1.BackupConfigDeletionPolicy `json:"deletionPolicy,omitempty"`
}

// BackupHooks specifies the hooks that will be executed before and/or after backup
type BackupHooks struct {
	// PreBackup specifies a list of hooks that will be executed before backup
	// +optional
	PreBackup []HookInfo `json:"preBackup,omitempty"`

	// PostBackup specifies a list of hooks that will be executed after backup
	// +optional
	PostBackup []HookInfo `json:"postBackup,omitempty"`
}

// BackupConfigurationStatus defines the observed state of BackupConfiguration
type BackupConfigurationStatus struct {
	// +optional
	OffshootStatus `json:",inline"`

	// Phase represents the current state of the Backup Invoker.
	// +optional
	Phase BackupInvokerPhase `json:"phase,omitempty"`

	// TargetFound specifies whether the backup target exist or not
	// +optional
	TargetFound *bool `json:"targetFound,omitempty"`

	// Conditions represents list of conditions regarding this BackupConfiguration
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// BackupInvokerPhase specifies the current state of the backup setup process
// +kubebuilder:validation:Enum=NotReady;Ready;Invalid
type BackupInvokerPhase string

const (
	BackupInvokerNotReady BackupInvokerPhase = "NotReady"
	BackupInvokerReady    BackupInvokerPhase = "Ready"
	BackupInvokerInvalid  BackupInvokerPhase = "Invalid"
)

// OffshootStatus specifies the status that are common between BackupConfiguration and BackupBatch
type OffshootStatus struct {
	// Backends specifies whether the backends exist or not
	// +optional
	Backends []BackendStatus `json:"backends,omitempty"`

	// Repositories specifies whether the repositories have been successfully initialized or not
	// +optional
	Repositories []RepoStatus `json:"repositories,omitempty"`

	// Dependencies specifies whether the objects required by this BackupConfiguration exist or not
	// +optional
	Dependencies []ResourceFoundStatus `json:"dependencies,omitempty"`

	// Sessions specifies status of the session specific resources
	// +optional
	Sessions []SessionStatus `json:"sessions,omitempty"`
}

// BackendStatus specifies the status of the backends
type BackendStatus struct {
	// Name indicates the backend name
	Name string `json:"name,omitempty"`

	// Ready indicates whether the respective BackupStorage is ready or not
	// +optional
	Ready *bool `json:"ready,omitempty"`

	// Storage indicates the status of the respective BackupStorage
	// +optional
	Storage *StorageStatus `json:"storage,omitempty"`

	// RetentionPolicy indicates the status of the respective RetentionPolicy
	// +optional
	RetentionPolicy *RetentionPolicyStatus `json:"retentionPolicy,omitempty"`
}

type StorageStatus struct {
	// Ref indicates to the BackupStorage object.
	Ref kmapi.ObjectReference `json:"ref,omitempty"`

	// Phase indicates the current phase of the respective BackupStorage.
	// +optional
	Phase v1alpha1.BackupStoragePhase `json:"phase,omitempty"`

	// Reason specifies the error messages found while checking the BackupStorage phase
	// +optional
	Reason string `json:"reason,omitempty"`
}

type RetentionPolicyStatus struct {
	// Ref indicates the RetentionPolicy object reference.
	Ref kmapi.ObjectReference `json:"ref,omitempty"`

	// Found indicates whether the RetentionPolicy is Found or not
	// +optional
	Found *bool `json:"found,omitempty"`

	// Reason specifies the error messages found while checking the RetentionPolicy
	// +optional
	Reason string `json:"reason,omitempty"`
}

// RepoStatus specifies the status of a Repository
type RepoStatus struct {
	// Name indicate the name of the Repository
	Name string `json:"name,omitempty"`

	// Ready indicates whether the respective Repository is ready or not
	// +optional
	Phase v1alpha1.RepositoryPhase `json:"phase,omitempty"`

	// Reason specifies the error messages found while ensuring the respective Repository
	// +optional
	Reason string `json:"reason,omitempty"`

	// VerificationConfigured indicates whether the verification for this repository is configured or not
	// +optional
	VerificationConfigured bool `json:"verificationConfigured,omitempty"`
}

// SessionStatus specifies the status of a session specific fields.
type SessionStatus struct {
	// Name indicates the name of the session
	Name string `json:"name,omitempty"`

	// NextSchedule specifies when the next backup will execute for this session
	// +optional
	NextSchedule string `json:"nextSchedule,omitempty"`

	// Conditions specifies a list of conditions related to this session
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

const (
	// TypeValidationPassed indicates the validation conditions of the CRD are passed or not.
	TypeValidationPassed           = "ValidationPassed"
	ReasonResourceValidationPassed = "ResourceValidationPassed"
	ReasonResourceValidationFailed = "ResourceValidationFailed"

	// TypeSchedulerEnsured indicates whether the Scheduler is ensured or not.
	TypeSchedulerEnsured      = "SchedulerEnsured"
	ReasonSchedulerNotEnsured = "SchedulerNotEnsured"
	ReasonSchedulerEnsured    = "SchedulerEnsured"

	// TypeInitialBackupTriggered indicates whether the initial backup is triggered or not.
	TypeInitialBackupTriggered               = "InitialBackupTriggered"
	ReasonFailedToTriggerInitialBackup       = "FailedToTriggerInitialBackup"
	ReasonSuccessfullyTriggeredInitialBackup = "SuccessfullyTriggeredInitialBackup"
)

//+kubebuilder:object:root=true

// BackupConfigurationList contains a list of BackupConfiguration
type BackupConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupConfiguration{}, &BackupConfigurationList{})
}
