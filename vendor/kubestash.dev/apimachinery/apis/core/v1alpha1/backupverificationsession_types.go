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
)

const (
	ResourceKindBackupVerificationSession     = "BackupVerificationSession"
	ResourceSingularBackupVerificationSession = "backupverificationsession"
	ResourcePluralBackupVerificationSession   = "backupverificationsessions"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=backupverificationsession,singular=backupverificationsession,categories={kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Duration",type="string",JSONPath=".status.duration"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// BackupVerificationSession represent one backup verification run for the target(s) pointed by the
// respective BackupConfiguration or BackupBatch
type BackupVerificationSession struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupVerificationSessionSpec   `json:"spec,omitempty"`
	Status BackupVerificationSessionStatus `json:"status,omitempty"`
}

// BackupVerificationSessionSpec specifies the information related to the respective backup verifier, session, repository and snapshot.
type BackupVerificationSessionSpec struct {
	// Invoker points to the respective BackupConfiguration or BackupBatch
	// which is responsible for triggering this backup verification.
	Invoker *core.TypedLocalObjectReference `json:"invoker,omitempty"`

	// Session specifies the name of the session that triggered this backup verification
	Session string `json:"session,omitempty"`

	// Repository specifies the name of the repository whose backed-up data will be verified
	Repository string `json:"repository,omitempty"`

	// Snapshot specifies the name of the snapshot that will be verified
	Snapshot string `json:"snapshot,omitempty"`

	// RetryLeft specifies number of retry attempts left for the backup verification session.
	// If this set to non-zero, KubeStash will create a new BackupVerificationSession if the current one fails.
	// +optional
	RetryLeft int32 `json:"retryLeft,omitempty"`
}

// BackupVerificationSessionStatus defines the observed state of BackupVerificationSession
type BackupVerificationSessionStatus struct {
	// Phase represents the current state of the backup verification process.
	// +optional
	Phase BackupVerificationSessionPhase `json:"phase,omitempty"`

	// Duration specifies the time required to complete the backup verification process
	// +optional
	Duration string `json:"duration,omitempty"`

	// Retried specifies whether this session was retried or not.
	// This field will exist only if the `retryConfig` has been set in the respective backup verification strategy.
	// +optional
	Retried *bool `json:"retried,omitempty"`

	// Conditions represents list of conditions regarding this BackupSession
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// BackupVerificationSessionPhase specifies the current state of the backup verification process
// +kubebuilder:validation:Enum=Running;Succeeded;Failed;Skipped
type BackupVerificationSessionPhase string

const (
	BackupVerificationSessionRunning   BackupVerificationSessionPhase = "Running"
	BackupVerificationSessionSucceeded BackupVerificationSessionPhase = "Succeeded"
	BackupVerificationSessionFailed    BackupVerificationSessionPhase = "Failed"
	BackupVerificationSessionSkipped   BackupVerificationSessionPhase = "Skipped"
)

// ============================ Conditions ========================

const (
	// TypeBackupVerificationSkipped indicates that the current session was skipped
	TypeBackupVerificationSkipped = "BackupVerificationSkipped"
	// ReasonSkippedVerifyingNewBackup indicates that the backup verification was skipped because the snapshot has already been verified
	ReasonSkippedVerifyingNewBackup = "SnapshotAlreadyVerified"

	// TypeVerificationSessionHistoryCleaned indicates whether the backup history was cleaned or not according to sessionHistoryLimit
	TypeVerificationSessionHistoryCleaned               = "VerificationSessionHistoryCleaned"
	ReasonSuccessfullyCleanedVerificationSessionHistory = "SuccessfullyCleanedVerificationSessionHistory"
	ReasonFailedToCleanVerificationSessionHistory       = "FailedToCleanVerificationSessionHistory"

	// TypeVerificationExecutorEnsured indicates whether the backup verification executor is ensured or not.
	TypeVerificationExecutorEnsured               = "VerificationExecutorEnsured"
	ReasonSuccessfullyEnsuredVerificationExecutor = "SuccessfullyEnsuredVerificationExecutor"
	ReasonFailedToEnsureVerificationExecutor      = "FailedToEnsureVerificationExecutor"

	// TypeRestoreSucceeded indicates whether the restore is succeeded or not.
	TypeRestoreSucceeded             = "RestoreSucceeded"
	ReasonSuccessfullyRestoredBackup = "SuccessfullyRestoredBackup"
	ReasonFailedToRestoreBackup      = "FailedToRestoreBackup"

	// TypeBackupVerified indicates whether backup is verified or not
	TypeBackupVerified               = "BackupVerified"
	ReasonSuccessfullyVerifiedBackup = "SuccessfullyVerifiedBackup"
	ReasonFailedToVerifyBackup       = "FailedToVerifyBackup"
)

//+kubebuilder:object:root=true

// BackupVerificationSessionList contains a list of BackupVerificationSession
type BackupVerificationSessionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupVerificationSession `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupVerificationSession{}, &BackupVerificationSessionList{})
}
