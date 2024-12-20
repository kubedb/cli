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
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindBackupVerifier   = "BackupVerifier"
	ResourceSingularBackupVerier = "backupverifier"
	ResourcePluralBackupVerifier = "backupverificatiers"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=backupverifier,singular=backupverifier,categories={kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// BackupVerifier represents backup verification configurations
type BackupVerifier struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BackupVerifierSpec `json:"spec,omitempty"`
}

// BackupVerifierSpec specifies the information related to the respective restore target, verification schedule, and verification type.
type BackupVerifierSpec struct {
	// RestoreOption specifies the restore target, and addonInfo for backup verification
	// +optional
	RestoreOption *RestoreOption `json:"restoreOption,omitempty"`

	// Scheduler specifies the configuration for verification triggering CronJob
	Scheduler *SchedulerSpec `json:"scheduler,omitempty"`

	// Function specifies the name of a Function CR that defines a container definition
	// which will execute the verification logic for a particular application.
	Function string `json:"function,omitempty"`

	// Volumes indicates the list of volumes that should be mounted on the verification job.
	Volumes []ofst.Volume `json:"volumes,omitempty"`

	// VolumeMounts specifies the mount for the volumes specified in `Volumes` section
	VolumeMounts []core.VolumeMount `json:"volumeMounts,omitempty"`

	// Type indicates the type of verifier that will verify the backup.
	// Valid values are:
	// - "RestoreOnly": KubeStash will create a RestoreSession with the tasks provided in BackupVerifier.
	// - "Query": KubeStash operator will restore data and then create a job to run the queries.
	// - "Script": KubeStash operator will restore data and then create a job to run the script.
	Type VerificationType `json:"type,omitempty"`

	// Query specifies the queries to be run to verify backup.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Query *runtime.RawExtension `json:"query,omitempty"`

	// Script specifies the script to be run to verify backup.
	// +optional
	Script *ScriptVerifierSpec `json:"script,omitempty"`

	// RetryConfig specifies the behavior of the retry mechanism in case of a verification failure.
	// +optional
	RetryConfig *RetryConfig `json:"retryConfig,omitempty"`

	// SessionHistoryLimit specifies how many BackupVerificationSessions and associate resources KubeStash should keep for debugging purpose.
	// The default value is 1.
	// +kubebuilder:default=1
	// +optional
	SessionHistoryLimit int32 `json:"sessionHistoryLimit,omitempty"`

	// RuntimeSettings allow to specify Resources, NodeSelector, Affinity, Toleration, ReadinessProbe etc.
	// for the verification job.
	// +optional
	RuntimeSettings ofst.RuntimeSettings `json:"runtimeSettings,omitempty"`
}

type RestoreOption struct {
	// Target indicates the target application where the data will be restored
	// +optional
	Target *kmapi.TypedObjectReference `json:"target,omitempty"`

	// AddonInfo specifies addon configuration that will be used to restore this target.
	AddonInfo *AddonInfo `json:"addonInfo,omitempty"`
}

// VerificationType specifies the type of verifier that will verify the backup
// +kubebuilder:validation:Enum=RestoreOnly;Query;Script
type VerificationType string

const (
	RestoreOnlyVerificationType VerificationType = "RestoreOnly"
	QueryVerificationType       VerificationType = "Query"
	ScriptVerificationType      VerificationType = "Script"
)

// ScriptVerifierSpec defines the script location in verifier job and the args to be provided with the script.
type ScriptVerifierSpec struct {
	// Location specifies the absolute path of the script file's location.
	Location string `json:"location,omitempty"`

	// Args specifies the arguments to be provided with the script.
	Args []string `json:"args,omitempty"`
}

//+kubebuilder:object:root=true

// BackupVerifierList contains a list of BackupVerifier
type BackupVerifierList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupVerifier `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupVerifier{}, &BackupVerifierList{})
}
