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
	ResourceKindRestoreSession     = "RestoreSession"
	ResourceSingularRestoreSession = "restoresession"
	ResourcePluralRestoreSession   = "restoresessions"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=restoresessions,singular=restoresession,shortName=restore,categories={kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Repository",type="string",JSONPath=".spec.dataSource.repository"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Duration",type="string",JSONPath=".status.duration"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// RestoreSession represents one restore run for the targeted application
type RestoreSession struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RestoreSessionSpec   `json:"spec,omitempty"`
	Status RestoreSessionStatus `json:"status,omitempty"`
}

// RestoreSessionSpec specifies the necessary configurations for restoring data into a target
type RestoreSessionSpec struct {
	// Target indicates the target application where the data will be restored.
	// +optional
	Target *kmapi.TypedObjectReference `json:"target,omitempty"`

	// DataSource specifies the information about the data that will be restored
	DataSource *RestoreDataSource `json:"dataSource,omitempty"`

	// Addon specifies addon configuration that will be used to restore the target.
	Addon *AddonInfo `json:"addon,omitempty"`

	// Hooks specifies the restore hooks that should be executed before and/or after the restore.
	// +optional
	Hooks *RestoreHooks `json:"hooks,omitempty"`

	// RestoreTimeout specifies a duration that KubeStash should wait for the restore to be completed.
	// If the restore tasks do not finish within this time period, KubeStash will consider this restore as a failure.
	// +optional
	RestoreTimeout *metav1.Duration `json:"restoreTimeout,omitempty"`

	// ManifestOptions provide options to select particular manifest object to restore
	// +optional
	ManifestOptions *ManifestRestoreOptions `json:"manifestOptions,omitempty"`
}

type ManifestRestoreOptions struct {
	// Workload specifies the options for Workload components to restore in manifest restore
	// +optional
	Workload *WorkloadManifestOptions `json:"workload,omitempty"`

	// MongoDB specifies the options for selecting particular MongoDB components to restore in manifest restore
	// +optional
	MongoDB *KubeDBManifestOptions `json:"mongoDB,omitempty"`

	// Postgres specifies the options for selecting particular Postgres components to restore in manifest restore
	// +optional
	Postgres *KubeDBManifestOptions `json:"postgres,omitempty"`

	// MySQL specifies the options for selecting particular MySQL components to restore in manifest restore
	// +optional
	MySQL *KubeDBManifestOptions `json:"mySQL,omitempty"`

	// MariaDB specifies the options for selecting particular MariaDB components to restore in manifest restore
	// +optional
	MariaDB *KubeDBManifestOptions `json:"mariaDB,omitempty"`

	// MSSQLServer specifies the options for selecting particular MSSQLServer components to restore in manifest restore
	// +optional
	MSSQLServer *MSSQLServerManifestOptions `json:"msSQLServer,omitempty"`

	// Druid specifies the options for selecting particular Druid components to restore in manifest restore
	// +optional
	Druid *DruidManifestOptions `json:"druid,omitempty"`

	// ZooKeeper specifies the options for selecting particular ZooKeeper components to restore in manifest restore
	// +optional
	ZooKeeper *KubeDBManifestOptions `json:"zooKeeper,omitempty"`

	// Singlestore specifies the options for selecting particular Singlestore components to restore in manifest restore
	// +optional
	Singlestore *KubeDBManifestOptions `json:"singlestore,omitempty"`

	// Redis specifies the options for selecting particular Redis components to restore in manifest restore
	// +optional
	Redis *KubeDBManifestOptions `json:"redis,omitempty"`

	// RedisSentinel specifies the options for selecting particular RedisSentinel components to restore in manifest restore
	// +optional
	RedisSentinel *RedisSentinelManifestOptions `json:"redisSentinel,omitempty"`
}

type RedisSentinelManifestOptions struct {
	// RestoreNamespace specifies the Namespace where the restored files will be applied
	// +optional
	RestoreNamespace string `json:"restoreNamespace,omitempty"`

	// Sentinel specifies whether to restore the Sentinel manifest or not
	// +kubebuilder:default=true
	// +optional
	Sentinel *bool `json:"sentinel,omitempty"`

	// SentinelName specifies the new name of the Sentinel yaml after restore
	// +optional
	SentinelName string `json:"SentinelName,omitempty"`

	// AuthSecret specifies whether to restore the AuthSecret manifest or not
	// +kubebuilder:default=true
	// +optional
	AuthSecret *bool `json:"authSecret,omitempty"`

	// AuthSecretName specifies new name of the AuthSecret yaml after restore
	// +optional
	AuthSecretName string `json:"authSecretName,omitempty"`

	// TLSIssuerRef specifies the name of the IssuerRef used for TLS configurations for both client and server
	// +optional
	TLSIssuerRef *core.TypedLocalObjectReference `json:"tlsIssuerRef,omitempty"`
}

type WorkloadManifestOptions struct {
	// RestoreNamespace specifies the Namespace where the restored files will be applied
	// +optional
	RestoreNamespace string `json:"restoreNamespace,omitempty"`
}

type MSSQLServerManifestOptions struct {
	// RestoreNamespace specifies the Namespace where the restored files will be applied
	// +optional
	RestoreNamespace string `json:"restoreNamespace,omitempty"`

	// DB specifies whether to restore the DB manifest or not
	// +kubebuilder:default=true
	// +optional
	DB *bool `json:"db,omitempty"`

	// DBName specifies the new name of the DB yaml after restore
	// +optional
	DBName string `json:"dbName,omitempty"`

	// AuthSecret specifies whether to restore the AuthSecret manifest or not
	// +kubebuilder:default=true
	// +optional
	AuthSecret *bool `json:"authSecret,omitempty"`

	// AuthSecretName specifies new name of the AuthSecret yaml after restore
	// +optional
	AuthSecretName string `json:"authSecretName,omitempty"`

	// InitScript specifies whether to restore the InitScript manifest or not
	// +kubebuilder:default=true
	// +optional
	InitScript *bool `json:"initScript,omitempty"`

	// Archiver specifies whether to restore the Archiver manifest or not
	// +kubebuilder:default=false
	// +optional
	Archiver *bool `json:"archiver,omitempty"`

	// ArchiverRef specifies the new name and namespace of the Archiver yaml after restore
	// +optional
	ArchiverRef *kmapi.ObjectReference `json:"archiverRef,omitempty"`

	// TLSIssuerRef specifies the name of the IssuerRef used for TLS configurations for both client and server.
	// +optional
	TLSIssuerRef *core.TypedLocalObjectReference `json:"tlsIssuerRef,omitempty"`
}

type DruidManifestOptions struct {
	// RestoreNamespace specifies the Namespace where the restored files will be applied
	// +optional
	RestoreNamespace string `json:"restoreNamespace,omitempty"`

	// DB specifies whether to restore the DB manifest or not
	// +kubebuilder:default=true
	// +optional
	DB *bool `json:"db,omitempty"`

	// DBName specifies the new name of the DB yaml after restore
	// +optional
	DBName string `json:"dbName,omitempty"`

	// AuthSecret specifies whether to restore the AuthSecret manifest or not
	// +kubebuilder:default=true
	// +optional
	AuthSecret *bool `json:"authSecret,omitempty"`

	// AuthSecretName specifies new name of the AuthSecret yaml after restore
	// +optional
	AuthSecretName string `json:"authSecretName,omitempty"`

	// ConfigSecret specifies whether to restore the ConfigSecret manifest or not
	// +kubebuilder:default=true
	// +optional
	ConfigSecret *bool `json:"configSecret,omitempty"`

	// ConfigSecretName specifies new name of the ConfigSecret yaml after restore
	// +optional
	ConfigSecretName string `json:"configSecretName,omitempty"`

	// DeepStorageSecret specifies whether to restore the DeepStorageSecret manifest or not
	// +kubebuilder:default=true
	// +optional
	DeepStorageSecret *bool `json:"deepStorageSecret,omitempty"`
}

type KubeDBManifestOptions struct {
	// RestoreNamespace specifies the Namespace where the restored files will be applied
	// +optional
	RestoreNamespace string `json:"restoreNamespace,omitempty"`

	// DB specifies whether to restore the DB manifest or not
	// +kubebuilder:default=true
	// +optional
	DB *bool `json:"db,omitempty"`

	// DBName specifies the new name of the DB yaml after restore
	// +optional
	DBName string `json:"dbName,omitempty"`

	// AuthSecret specifies whether to restore the AuthSecret manifest or not
	// +kubebuilder:default=true
	// +optional
	AuthSecret *bool `json:"authSecret,omitempty"`

	// AuthSecretName specifies new name of the AuthSecret yaml after restore
	// +optional
	AuthSecretName string `json:"authSecretName,omitempty"`

	// Archiver specifies whether to restore the Archiver manifest or not
	// +kubebuilder:default=false
	// +optional
	Archiver *bool `json:"archiver,omitempty"`

	// ArchiverRef specifies the new name and namespace of the Archiver yaml after restore
	// +optional
	ArchiverRef *kmapi.ObjectReference `json:"archiverRef,omitempty"`

	// ConfigSecret specifies whether to restore the ConfigSecret manifest or not
	// +kubebuilder:default=true
	// +optional
	ConfigSecret *bool `json:"configSecret,omitempty"`

	// ConfigSecretName specifies new name of the ConfigSecret yaml after restore
	// +optional
	ConfigSecretName string `json:"configSecretName,omitempty"`

	// InitScript specifies whether to restore the InitScript manifest or not
	// +kubebuilder:default=true
	// +optional
	InitScript *bool `json:"initScript,omitempty"`

	// TLSIssuerRef specifies the name of the IssuerRef used for TLS configurations for both client and server
	// +optional
	TLSIssuerRef *core.TypedLocalObjectReference `json:"tlsIssuerRef,omitempty"`
}

// RestoreDataSource specifies the information about the data that will be restored
type RestoreDataSource struct {
	// Namespace specifies the namespace of the DataSource (i.e. Repository, Snapshot).
	Namespace string `json:"namespace,omitempty"`

	// Repository points to the Repository name from which the data will be restored
	Repository string `json:"repository,omitempty"`

	// Snapshot specifies the Snapshot name that will be restored.
	// If you want to use Point-In-Time recovery option, don't specify this field. Specify `pitr` field instead.
	// +optional
	Snapshot string `json:"snapshot,omitempty"`

	// PITR stands for Point-In-Time Recovery. You can provide a target time instead of specifying a particular Snapshot.
	// KubeStash will automatically find the latest Snapshot that satisfies the targeted time and restore it.
	// +optional
	PITR *PITR `json:"pitr,omitempty"`

	// Components specifies the components that will be restored. If you keep this field empty, then all
	// the components that were backed up in the desired Snapshot will be restored.
	// +optional
	Components []string `json:"components,omitempty"`

	// EncryptionSecret refers to the Secret containing the encryption key which will be used to encode/decode the backed up data.
	// You can refer to a Secret of a different namespace.
	// If you don't provide the namespace field, KubeStash will look for the Secret in the same namespace as the RestoreSession.
	// +optional
	EncryptionSecret *kmapi.ObjectReference `json:"encryptionSecret,omitempty"`
}

// PITR specifies the target time and behavior of Point-In-Time Recovery
type PITR struct {
	// TargetTime specifies the desired date and time at which you want to roll back your application data
	TargetTime *metav1.Time `json:"targetTime,omitempty"`

	// Exclusive specifies whether to exclude the Snapshot that falls in the exact time specified
	// in the `targetTime` field. By default, KubeStash will select the Snapshot that fall in the exact time.
	// +optional
	Exclusive bool `json:"exclusive,omitempty"`
}

// RestoreHooks specifies the hooks that will be executed before and/or after restore
type RestoreHooks struct {
	// PreRestore specifies a list of hooks that will be executed before restore
	// +optional
	PreRestore []HookInfo `json:"preRestore,omitempty"`

	// PostRestore specifies a list of hooks that will be executed after restore
	// +optional
	PostRestore []HookInfo `json:"postRestore,omitempty"`
}

// RestoreSessionStatus defines the observed state of RestoreSession
type RestoreSessionStatus struct {
	// Phase represents the current state of the restore process
	// +optional
	Phase RestorePhase `json:"phase,omitempty"`

	// TargetFound specifies whether the restore target exist or not
	// +optional
	TargetFound *bool `json:"targetFound,omitempty"`

	// Duration specifies the total time taken to complete the restore process
	// +optional
	Duration string `json:"duration,omitempty"`

	// RestoreDeadline specifies the deadline of restore. Restore will be
	// considered Failed if it does not complete within this deadline
	// +optional
	RestoreDeadline *metav1.Time `json:"restoreDeadline,omitempty"`

	// TotalComponents represents the number of total components for this RestoreSession
	// +optional
	TotalComponents int32 `json:"totalComponents,omitempty"`

	// Components represents the individual component restore status
	// +optional
	// +mapType=granular
	Components map[string]ComponentRestoreStatus `json:"components,omitempty"`

	// Hooks represents the hook execution status
	// +optional
	Hooks HookStatus `json:"hooks,omitempty"`

	// Dependencies specifies whether the objects required by this RestoreSession exist or not
	// +optional
	Dependencies []ResourceFoundStatus `json:"dependencies,omitempty"`

	// PausedBackups represents the list of backups that have been paused before restore.
	// +optional
	PausedBackups []kmapi.TypedObjectReference `json:"pausedBackups,omitempty"`

	// Conditions specifies a list of conditions related to this restore session
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// RestorePhase represents the current state of the restore process
// +kubebuilder:validation:Enum=Pending;Running;Failed;Succeeded;Invalid;Unknown
type RestorePhase string

const (
	RestorePending      RestorePhase = "Pending"
	RestoreRunning      RestorePhase = "Running"
	RestoreFailed       RestorePhase = "Failed"
	RestoreSucceeded    RestorePhase = "Succeeded"
	RestoreInvalid      RestorePhase = "Invalid"
	RestorePhaseUnknown RestorePhase = "Unknown"
)

// ComponentRestoreStatus represents the restore status of individual components
type ComponentRestoreStatus struct {
	// Phase represents the restore phase of the component
	// +optional
	Phase RestorePhase `json:"phase,omitempty"`

	// Duration specifies the total time taken to complete the restore process for this component
	// +optional
	Duration string `json:"duration,omitempty"`

	// Error specifies the reason in case of restore failure for the component
	// +optional
	Error string `json:"error,omitempty"`
}

const (
	TypeRestoreExecutorEnsured               = "RestoreExecutorEnsured"
	ReasonSuccessfullyEnsuredRestoreExecutor = "SuccessfullyEnsuredRestoreExecutor"
	ReasonFailedToEnsureRestoreExecutor      = "FailedToEnsureRestoreExecutor"

	TypePreRestoreHooksExecutionSucceeded     = "PreRestoreHooksExecutionSucceeded"
	ReasonSuccessfullyExecutedPreRestoreHooks = "SuccessfullyExecutedPreRestoreHooks"
	ReasonFailedToExecutePreRestoreHooks      = "FailedToExecutePreRestoreHooks"

	TypePostRestoreHooksExecutionSucceeded     = "PostRestoreHooksExecutionSucceeded"
	ReasonSuccessfullyExecutedPostRestoreHooks = "SuccessfullyExecutedPostRestoreHooks"
	ReasonFailedToExecutePostRestoreHooks      = "FailedToExecutePostRestoreHooks"

	TypeRestoreTargetFound                = "RestoreTargetFound"
	ReasonUnableToCheckTargetAvailability = "UnableToCheckTargetAvailability"

	TypeRestoreIncomplete                           = "RestoreIncomplete"
	ReasonRestoreExecutorTerminatedBeforeCompletion = "RestoreExecutorTerminatedBeforeCompletion"
)

//+kubebuilder:object:root=true

// RestoreSessionList contains a list of RestoreSession
type RestoreSessionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RestoreSession `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RestoreSession{}, &RestoreSessionList{})
}
