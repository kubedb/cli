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
	storage "kubestash.dev/apimachinery/apis/storage/v1alpha1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindBackupSession     = "BackupSession"
	ResourceSingularBackupSession = "backupsession"
	ResourcePluralBackupSession   = "backupsessions"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=backupsessions,singular=backupsession,categories={kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Invoker-Type",type="string",JSONPath=".spec.invoker.kind"
// +kubebuilder:printcolumn:name="Invoker-Name",type="string",JSONPath=".spec.invoker.name"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Duration",type="string",JSONPath=".status.duration"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// BackupSession represent one backup run for the target(s) pointed by the
// respective BackupConfiguration or BackupBatch
type BackupSession struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSessionSpec   `json:"spec,omitempty"`
	Status BackupSessionStatus `json:"status,omitempty"`
}

// BackupSessionSpec specifies the information related to the respective backup invoker and session.
type BackupSessionSpec struct {
	// Invoker points to the respective BackupConfiguration or BackupBatch
	// which is responsible for triggering this backup.
	Invoker *core.TypedLocalObjectReference `json:"invoker,omitempty"`

	// Session specifies the name of the session that triggered this backup
	Session string `json:"session,omitempty"`

	// RetryLeft specifies number of retry attempts left for the session.
	// If this set to non-zero, KubeStash will create a new BackupSession if the current one fails.
	// +optional
	RetryLeft int32 `json:"retryLeft,omitempty"`

	// BackupTimeout specifies the maximum duration of backup. Backup will be considered Failed
	// if backup tasks do not complete within this time limit. By default, KubeStash don't set any timeout for backup.
	// +optional
	BackupTimeout *metav1.Duration `json:"backupTimeout,omitempty"`
}

// BackupSessionStatus defines the observed state of BackupSession
type BackupSessionStatus struct {
	// Phase represents the current state of the backup process.
	// +optional
	Phase BackupSessionPhase `json:"phase,omitempty"`

	// Duration specifies the time required to complete the backup process
	// +optional
	Duration string `json:"duration,omitempty"`

	// BackupDeadline specifies the deadline of backup. Backup will be
	// considered Failed if it does not complete within this deadline
	// +optional
	BackupDeadline *metav1.Time `json:"backupDeadline,omitempty"`

	// TotalSnapshots specifies the total number of snapshots created for this backupSession.
	// +optional
	TotalSnapshots *int32 `json:"totalSnapshots,omitempty"`

	// Snapshots specifies the Snapshots status
	// +optional
	Snapshots []SnapshotStatus `json:"snapshots,omitempty"`

	// Hooks represents the hook execution status
	// +optional
	Hooks HookStatus `json:"hooks,omitempty"`

	// RetentionPolices specifies whether the retention policies were properly applied on the repositories or not
	// +optional
	RetentionPolicies []RetentionPolicyApplyStatus `json:"retentionPolicy,omitempty"`

	// Retried specifies whether this session was retried or not.
	// This field will exist only if the `retryConfig` has been set in the respective backup invoker.
	// +optional
	Retried *bool `json:"retried,omitempty"`

	// NextRetry specifies the time when KubeStash should retry the current failed backup.
	// This field will exist only if the `retryConfig` has been set in the respective backup invoker.
	// +optional
	NextRetry *metav1.Time `json:"nextRetry,omitempty"`

	// Conditions represents list of conditions regarding this BackupSession
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// BackupSessionPhase specifies the current state of the backup process
// +kubebuilder:validation:Enum=Pending;Running;Succeeded;Failed;Skipped
type BackupSessionPhase string

const (
	BackupSessionPending   BackupSessionPhase = "Pending"
	BackupSessionRunning   BackupSessionPhase = "Running"
	BackupSessionSucceeded BackupSessionPhase = "Succeeded"
	BackupSessionFailed    BackupSessionPhase = "Failed"
	BackupSessionSkipped   BackupSessionPhase = "Skipped"
)

// SnapshotStatus represents the current state of respective the Snapshot
type SnapshotStatus struct {
	// Name indicates to the name of the Snapshot
	Name string `json:"name,omitempty"`

	// Phase indicate the phase of the Snapshot
	// +optional
	Phase storage.SnapshotPhase `json:"phase,omitempty"`

	// AppRef points to the application that is being backed up in this Snapshot
	AppRef *kmapi.TypedObjectReference `json:"appRef,omitempty"`

	// Repository indicates the name of the Repository where the Snapshot is being stored.
	Repository string `json:"repository,omitempty"`
}

// RetentionPolicyApplyStatus represents the state of the applying retention policy
type RetentionPolicyApplyStatus struct {
	// Ref points to the RetentionPolicy CR that is being used to cleanup the old Snapshots for this session.
	Ref kmapi.ObjectReference `json:"ref,omitempty"`

	// Repository specifies the name of the Repository on which the RetentionPolicy has been applied.
	Repository string `json:"repository,omitempty"`
	// Phase specifies the state of retention policy apply process
	// +optional
	Phase RetentionPolicyApplyPhase `json:"phase,omitempty"`

	// Error represents the reason if the retention policy applying fail
	// +optional
	Error string `json:"error,omitempty"`
}

// RetentionPolicyApplyPhase represents the state of the retention policy apply process
// +kubebuilder:validation:Enum=Pending;Applied;FailedToApply
type RetentionPolicyApplyPhase string

const (
	RetentionPolicyPending       RetentionPolicyApplyPhase = "Pending"
	RetentionPolicyApplied       RetentionPolicyApplyPhase = "Applied"
	RetentionPolicyFailedToApply RetentionPolicyApplyPhase = "FailedToApply"
)

// ============================ Conditions ========================

const (
	// TypeBackupSkipped indicates that the current session was skipped
	TypeBackupSkipped = "BackupSkipped"
	// ReasonSkippedTakingNewBackup indicates that the backup was skipped because another backup was running or backup invoker is not ready state.
	ReasonSkippedTakingNewBackup = "PreRequisitesNotSatisfied"

	// TypeSessionHistoryCleaned indicates whether the backup history was cleaned or not according to backupHistoryLimit
	TypeSessionHistoryCleaned               = "SessionHistoryCleaned"
	ReasonSuccessfullyCleanedSessionHistory = "SuccessfullyCleanedSessionHistory"
	ReasonFailedToCleanSessionHistory       = "FailedToCleanSessionHistory"

	// TypePreBackupHooksExecutionSucceeded indicates whether the pre-backup hooks were executed successfully or not
	TypePreBackupHooksExecutionSucceeded     = "PreBackupHooksExecutionSucceeded"
	ReasonSuccessfullyExecutedPreBackupHooks = "SuccessfullyExecutedPreBackupHooks"
	ReasonFailedToExecutePreBackupHooks      = "FailedToExecutePreBackupHooks"

	// TypePostBackupHooksExecutionSucceeded indicates whether the pre-backup hooks were executed successfully or not
	TypePostBackupHooksExecutionSucceeded     = "PostBackupHooksExecutionSucceeded"
	ReasonSuccessfullyExecutedPostBackupHooks = "SuccessfullyExecutedPostBackupHooks"
	ReasonFailedToExecutePostBackupHooks      = "FailedToExecutePostBackupHooks"

	// TypeBackupExecutorEnsured indicates whether the Backup Executor is ensured or not.
	TypeBackupExecutorEnsured               = "BackupExecutorEnsured"
	ReasonSuccessfullyEnsuredBackupExecutor = "SuccessfullyEnsuredBackupExecutor"
	ReasonFailedToEnsureBackupExecutor      = "FailedToEnsureBackupExecutor"

	// TypeRetentionPolicyExecutorEnsured indicates whether the Backup Executor is ensured or not.
	TypeRetentionPolicyExecutorEnsured               = "RetentionPolicyExecutorEnsured"
	ReasonSuccessfullyEnsuredRetentionPolicyExecutor = "SuccessfullyEnsuredRetentionPolicyExecutor"
	ReasonFailedToEnsureRetentionPolicyExecutor      = "FailedToEnsureRetentionPolicyExecutor"

	// TypeSnapshotsEnsured indicates whether Snapshots are ensured for each Repository or not
	TypeSnapshotsEnsured               = "SnapshotsEnsured"
	ReasonSuccessfullyEnsuredSnapshots = "SuccessfullyEnsuredSnapshots"
	ReasonFailedToEnsureSnapshots      = "FailedToEnsureSnapshots"

	// TypeSnapshotCleanupIncomplete indicates whether Snapshot cleanup incomplete or not
	TypeSnapshotCleanupIncomplete                   = "SnapshotCleanupIncomplete"
	ReasonSnapshotCleanupTerminatedBeforeCompletion = "SnapshotCleanupTerminatedBeforeCompletion"
)

//+kubebuilder:object:root=true

// BackupSessionList contains a list of BackupSession
type BackupSessionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupSession `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupSession{}, &BackupSessionList{})
}
