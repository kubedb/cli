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
	storageapi "kubestash.dev/apimachinery/apis/storage/v1alpha1"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindKubeVirtArchiver     = "KubeVirtArchiver"
	ResourceSingularKubeVirtArchiver = "kubevirtarchiver"
	ResourcePluralKubeVirtArchiver   = "kubevirtarchivers"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubevirtarchivers,singular=kubevirtarchiver,shortName=kvarchiver,categories={archiver,kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Paused",type="boolean",JSONPath=".spec.pause"
// +kubebuilder:printcolumn:name="VMs",type="integer",JSONPath=".status.virtualMachinesMatched"
// +kubebuilder:printcolumn:name="Last-Incremental",type="string",JSONPath=".status.virtualMachines[*].lastIncrementalBackupTime"
// +kubebuilder:printcolumn:name="Incremental-Failures",type="string",JSONPath=".status.virtualMachines[*].incrementalBackupFailures"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// KubeVirtArchiver ties together the two halves of incremental VirtualMachine backup:
// a scheduled full that re-anchors the chain once a day, and a resident Job that
// fills every gap between with changed-block checkpoints.
type KubeVirtArchiver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeVirtArchiverSpec   `json:"spec,omitempty"`
	Status KubeVirtArchiverStatus `json:"status,omitempty"`
}

// KubeVirtArchiverSpec defines the desired state of KubeVirtArchiver
type KubeVirtArchiverSpec struct {
	// VirtualMachines selects the VirtualMachines this archiver backs up.
	// +optional
	VirtualMachines *AllowedVirtualMachines `json:"virtualMachines,omitempty"`

	// BackupStorage is the backend the archiver's Repositories are created in.
	// +optional
	BackupStorage *ArchiverBackupStorage `json:"backupStorage,omitempty"`

	// EncryptionSecret refers to the Secret holding the key the backed up data is
	// encrypted with.
	// +optional
	EncryptionSecret *kmapi.ObjectReference `json:"encryptionSecret,omitempty"`

	// RetentionPolicy governs how long the full Snapshots are kept. Each full's
	// incremental checkpoint chain is retired together with it.
	// +optional
	RetentionPolicy *kmapi.ObjectReference `json:"retentionPolicy,omitempty"`

	// FullBackup configures the scheduled full backup session.
	// +optional
	FullBackup *FullBackupOptions `json:"fullBackup,omitempty"`

	// ManifestBackup configures the scheduled backup of the VM's object graph.
	// +optional
	ManifestBackup *ManifestBackupOptions `json:"manifestBackup,omitempty"`

	// CBTBackup configures the resident checkpoint loop. It carries no scheduler.
	// +optional
	CBTBackup *CBTBackupOptions `json:"cbtBackup,omitempty"`

	// Pause stops the resident loop and suspends the scheduled sessions. The
	// incremental Snapshot is finalised to Succeeded, which is the only terminal
	// transition it has.
	// +optional
	Pause bool `json:"pause,omitempty"`

	// DeletionPolicy specifies what happens to the created Repositories when this
	// archiver is deleted.
	// +optional
	DeletionPolicy *storageapi.BackupConfigDeletionPolicy `json:"deletionPolicy,omitempty"`
}

// AllowedVirtualMachines selects VirtualMachines by namespace and by label.
//
// Opt-in is selector-only. KubeStash does not own the VirtualMachine CRD, so there
// is no field on the VM to set — a VM joins an archiver by carrying the label.
type AllowedVirtualMachines struct {
	// Namespaces indicates which namespaces the VirtualMachines are selected from.
	// +optional
	Namespaces *apis.AllowedNamespaces `json:"namespaces,omitempty"`

	// Selector matches the VirtualMachines within those namespaces. An empty
	// selector matches every VirtualMachine in them.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// ArchiverBackupStorage points at the BackupStorage the archiver's Repositories are created in.
type ArchiverBackupStorage struct {
	// Ref refers to the BackupStorage CR. It may live in another namespace as long
	// as that namespace is allowed by the BackupStorage's usage policy.
	// +optional
	Ref *kmapi.ObjectReference `json:"ref,omitempty"`

	// SubDir is the directory inside the BackupStorage the archiver writes under.
	// +optional
	SubDir string `json:"subDir,omitempty"`
}

// FullBackupOptions configures the scheduled half that produces full backups.
// A full re-reads every allocated extent, so it depends on no earlier data and
// becomes the point the checkpoint chain can be pruned back to.
type FullBackupOptions struct {
	// Driver specifies the tool that moves the data. KubeVirt native CBT
	// archivers store map-and-extent artifacts in Restic; the resident loop and
	// scheduled full must use the same driver so a chain is recoverable.
	// +kubebuilder:default:=Restic
	// +optional
	Driver apis.Driver `json:"driver,omitempty"`

	// Scheduler specifies the configuration for the backup triggering CronJob.
	// +optional
	Scheduler *SchedulerSpec `json:"scheduler,omitempty"`

	// JobTemplate specifies the pod template for the backup Job.
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`

	// RetryConfig specifies the behavior of retry in case of a backup failure.
	// +optional
	RetryConfig *RetryConfig `json:"retryConfig,omitempty"`

	// Timeout specifies the maximum duration of the backup.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// SessionHistoryLimit specifies how many backup Jobs and associated resources
	// KubeStash should keep for debugging purpose.
	// +optional
	SessionHistoryLimit int32 `json:"sessionHistoryLimit,omitempty"`
}

// ManifestBackupOptions configures the scheduled backup of the VirtualMachine's
// object graph — the VM, its VMI and its DataVolumes — which Restic stores in a
// repository of its own, separate from the disk blocks.
type ManifestBackupOptions struct {
	// Scheduler specifies the configuration for the backup triggering CronJob.
	// +optional
	Scheduler *SchedulerSpec `json:"scheduler,omitempty"`

	// JobTemplate specifies the pod template for the backup Job.
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`

	// RetryConfig specifies the behavior of retry in case of a backup failure.
	// +optional
	RetryConfig *RetryConfig `json:"retryConfig,omitempty"`

	// Timeout specifies the maximum duration of the backup.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// SessionHistoryLimit specifies how many backup Jobs and associated resources
	// KubeStash should keep for debugging purpose.
	// +optional
	SessionHistoryLimit int32 `json:"sessionHistoryLimit,omitempty"`
}

// CBTBackupOptions configures the resident half: a Job that loops on Interval,
// reading the blocks KubeVirt reports dirty since the previous checkpoint.
//
// There is deliberately no Scheduler field here. Its absence is what marks this
// half resident — the same convention KubeDB uses for LogBackupOptions. Adding a
// scheduler would turn the loop back into a CronJob, whose ~70-90s of Job and PVC
// churn per run is more than the interval it is meant to fill.
type CBTBackupOptions struct {
	// Interval is how long the loop sleeps between checkpoints. The floor is set
	// by the KubeVirt export handshake, which costs ~15-20s per increment.
	// +kubebuilder:default="10m"
	// +optional
	Interval *metav1.Duration `json:"interval,omitempty"`

	// ScratchStorageClass is the StorageClass for the scratch PVC. KubeVirt's
	// webhook requires one even in Pull mode, where nothing reads it back.
	// +optional
	ScratchStorageClass string `json:"scratchStorageClass,omitempty"`

	// ScratchSize is the size of the scratch PVC, 2Gi when unset. Size it for
	// the data the guest overwrites while one export is open.
	// +optional
	ScratchSize *resource.Quantity `json:"scratchSize,omitempty"`

	// SkipQuiesce skips guest filesystem quiescing for each native Pull export.
	// Set this for guests without a QEMU guest agent. Keeping the default false
	// preserves application-consistent backups when a guest agent is available.
	// +optional
	SkipQuiesce bool `json:"skipQuiesce,omitempty"`

	// TTLDuration optionally overrides the ceiling on how long one resident
	// checkpoint's export and CBT Lease hold may last (default 24h). Runs are not
	// bounded by it in practice: a stalled stream is aborted after ten minutes
	// without progress, a crashed run's leftovers are removed by the next run,
	// and deleting the VirtualMachineBackup is what completes a run. Set it only
	// to cap a checkpoint below the default; it must outlast the longest
	// checkpoint, because KubeVirt may treat an expired pull backup as complete.
	// +optional
	TTLDuration *metav1.Duration `json:"ttlDuration,omitempty"`

	// RuntimeSettings allows tuning the resident Job's pod and container.
	// +optional
	RuntimeSettings *ofst.RuntimeSettings `json:"runtimeSettings,omitempty"`
}

// KubeVirtArchiverStatus defines the observed state of KubeVirtArchiver
type KubeVirtArchiverStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions represents the current state of the archiver as a whole.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	// VirtualMachinesMatched is how many VirtualMachines the selector matched.
	// +optional
	VirtualMachinesMatched int32 `json:"virtualMachinesMatched,omitempty"`

	// VirtualMachines reports, per matched VirtualMachine, whether it is actually
	// backupable. Enabling changed-block tracking requires restarting the VM, so
	// the controller validates and reports rather than mutating — silently
	// restarting a production VM because a label changed would be hostile. This
	// list is where the user reads what to do about it.
	// +optional
	// +listType=map
	// +listMapKey=namespace
	// +listMapKey=name
	VirtualMachines []VirtualMachineArchiverStatus `json:"virtualMachines,omitempty"`
}

// VirtualMachineArchiverStatus is the archiver's view of one matched VirtualMachine.
type VirtualMachineArchiverStatus struct {
	// Name of the VirtualMachine.
	Name string `json:"name"`

	// Namespace of the VirtualMachine.
	Namespace string `json:"namespace"`

	// ChangedBlockTrackingReady is true only when vmi.status.changedBlockTracking
	// reports Enabled. Until then no backup for this VM can run at all: KubeVirt
	// refuses the VirtualMachineBackup outright.
	// +optional
	ChangedBlockTrackingReady bool `json:"changedBlockTrackingReady,omitempty"`

	// ChangedBlockTrackingState mirrors vmi.status.changedBlockTracking.state
	// verbatim. Each value maps to a different user action, so the raw value is
	// more useful than a boolean alone: Enabled, PendingRestart, Initializing,
	// Disabled, IncrementalBackupFeatureGateDisabled, or empty when the VM does
	// not match the cluster-level CBT label selector.
	// +optional
	ChangedBlockTrackingState string `json:"changedBlockTrackingState,omitempty"`

	// Reason is a short, machine-readable cause when ChangedBlockTrackingReady is
	// false.
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message tells the user what to do about Reason — restart the VM, label it,
	// or enable the IncrementalBackup feature gate.
	// +optional
	Message string `json:"message,omitempty"`

	// LastObservedTime is when the controller last looked at this VM.
	// +optional
	LastObservedTime *metav1.Time `json:"lastObservedTime,omitempty"`

	// LastIncrementalBackupTime is when the resident loop last stored a complete
	// checkpoint for this VM.
	// +optional
	LastIncrementalBackupTime *metav1.Time `json:"lastIncrementalBackupTime,omitempty"`

	// LastIncrementalBackupError is why the loop's latest attempt failed. It is
	// empty once an attempt succeeds.
	// +optional
	LastIncrementalBackupError string `json:"lastIncrementalBackupError,omitempty"`

	// IncrementalBackupFailures counts the loop's failed attempts since its last
	// success. A failed attempt is retried at the next interval.
	// +optional
	IncrementalBackupFailures int32 `json:"incrementalBackupFailures,omitempty"`
}

// AnnotationIncrementalBackupReport is where the resident loop reports its
// latest attempt on the Snapshot it appends to. The loop owns no status of its
// own; the controller mirrors the report into the KubeVirtArchiver's status.
const AnnotationIncrementalBackupReport = "kubevirtarchiver.kubestash.com/incremental-backup-report"

// IncrementalBackupReport is the value of AnnotationIncrementalBackupReport.
type IncrementalBackupReport struct {
	// ReportedAt is when the attempt ended.
	ReportedAt metav1.Time `json:"reportedAt"`

	// LastSuccessTime is when a checkpoint was last stored on this Snapshot.
	// +optional
	LastSuccessTime *metav1.Time `json:"lastSuccessTime,omitempty"`

	// LastError is the latest attempt's failure, empty when it succeeded.
	// +optional
	LastError string `json:"lastError,omitempty"`

	// Failures counts failed attempts since the last success.
	// +optional
	Failures int32 `json:"failures,omitempty"`
}

const (
	// TypeVirtualMachinesReady is True only when every matched VirtualMachine has
	// changed-block tracking Enabled.
	TypeVirtualMachinesReady           = "VirtualMachinesReady"
	ReasonAllVirtualMachinesReady      = "AllVirtualMachinesReady"
	ReasonChangedBlockTrackingNotReady = "ChangedBlockTrackingNotReady"

	// TypeBackupConfigurationsEnsured covers the scheduled half.
	TypeBackupConfigurationsEnsured          = "BackupConfigurationsEnsured"
	ReasonBackupConfigurationEnsureSucceeded = "BackupConfigurationEnsureSucceeded"
	ReasonBackupConfigurationEnsureFailed    = "BackupConfigurationEnsureFailed"

	// TypeArchiverJobsEnsured covers the resident half.
	TypeArchiverJobsEnsured          = "ArchiverJobsEnsured"
	ReasonArchiverJobEnsureSucceeded = "ArchiverJobEnsureSucceeded"
	ReasonArchiverJobEnsureFailed    = "ArchiverJobEnsureFailed"
)

// Reasons a VirtualMachine is not backupable, mapped from the CBT readiness state.
const (
	ReasonVMNotSelectedForCBT      = "NotSelectedForChangedBlockTracking"
	ReasonVMCBTPendingRestart      = "ChangedBlockTrackingPendingRestart"
	ReasonVMCBTInitializing        = "ChangedBlockTrackingInitializing"
	ReasonVMCBTFeatureGateDisabled = "IncrementalBackupFeatureGateDisabled"
	ReasonVMNotRunning             = "VirtualMachineNotRunning"
)

//+kubebuilder:object:root=true

// KubeVirtArchiverList contains a list of KubeVirtArchiver
type KubeVirtArchiverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeVirtArchiver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeVirtArchiver{}, &KubeVirtArchiverList{})
}
