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
	// If this set to non-zero, Stash will create a new BackupSession if the current one fails.
	// +optional
	RetryLeft int32 `json:"retryLeft,omitempty"`
}

// BackupSessionStatus defines the observed state of BackupSession
type BackupSessionStatus struct {
	// Phase represents the current state of the backup process.
	// +optional
	Phase BackupSessionPhase `json:"phase,omitempty"`

	// Duration specifies the time required to complete the backup process
	// +optional
	Duration string `json:"duration,omitempty"`

	// Deadline specifies the deadline of backup. BackupSession will be
	// considered Failed if backup does not complete within this deadline
	// +optional
	Deadline *metav1.Time `json:"sessionDeadline,omitempty"`

	// Snapshots specifies the Snapshots status
	// +optional
	Snapshots []SnapshotStatus `json:"snapshots,omitempty"`

	// Hooks represents the hook execution status
	// +optional
	Hooks HookStatus `json:"hooks,omitempty"`

	// Verifications specifies the backup verification status
	// +optional
	Verifications []VerificationStatus `json:"verifications,omitempty"`

	// RetentionPolices specifies whether the retention policies were properly applied on the repositories or not
	// +optional
	RetentionPolicies []RetentionPolicyApplyStatus `json:"retentionPolicy,omitempty"`

	// Retried specifies whether this session was retried or not.
	// This field will exist only if the `retryConfig` has been set in the respective backup invoker.
	// +optional
	Retried *bool `json:"retried,omitempty"`

	// NextRetry specifies the time when Stash should retry the current failed backup.
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

// VerificationStatus specifies the status of a backup verification
type VerificationStatus struct {
	// Name indicates the name of the respective verification strategy
	Name string `json:"name,omitempty"`

	// Phase represents the state of the verification process
	// +optional
	Phase BackupVerificationPhase `json:"phase,omitempty"`
}

// BackupVerificationPhase represents the state of the backup verification process
// +kubebuilder:validation:Enum=Verified;NotVerified;VerificationFailed
type BackupVerificationPhase string

const (
	Verified           BackupVerificationPhase = "Verified"
	NotVerified        BackupVerificationPhase = "NotVerified"
	VerificationFailed BackupVerificationPhase = "VerificationFailed"
)

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

	// TypeSnapshotsEnsured indicates whether Snapshots are ensured for each Repository or not
	TypeSnapshotsEnsured               = "SnapshotsEnsured"
	ReasonSuccessfullyEnsuredSnapshots = "SuccessfullyEnsuredSnapshots"
	ReasonFailedToEnsureSnapshots      = "FailedToEnsureSnapshots"
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
