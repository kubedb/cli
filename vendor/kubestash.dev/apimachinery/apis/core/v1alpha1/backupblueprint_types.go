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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindBackupBlueprint     = "BackupBlueprint"
	ResourceSingularBackupBlueprint = "backupblueprint"
	ResourcePluralBackupBlueprint   = "backupblueprints"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=backupblueprints,singular=backupblueprint,categories={kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// BackupBlueprint lets you define a common template for taking backup for all the similar applications.
// Then, you can just apply some annotations in the targeted application to enable backup.
// KubeStash will automatically resolve the template and create a BackupConfiguration for the targeted application.
type BackupBlueprint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BackupBlueprintSpec `json:"spec,omitempty"`
}

// BackupBlueprintSpec defines the desired state of BackupBlueprint
type BackupBlueprintSpec struct {
	// BackupConfigurationTemplate Specifies the BackupConfiguration that will be created by BackupBlueprint.
	BackupConfigurationTemplate *BackupConfigurationTemplate `json:"backupConfigurationTemplate,omitempty"`

	// Subjects specify a list of subject to which this BackupBlueprint is applicable. KubeStash will start watcher for these resources.
	// Multiple BackupBlueprints can have common subject. The watcher will find the appropriate blueprint from its annotations.
	Subjects []metav1.TypeMeta `json:"subjects,omitempty"`

	// UsagePolicy specifies a policy of how this BackupBlueprint will be used. For example,
	// you can use `allowedNamespaces` policy to restrict the usage of this BackupBlueprint to particular namespaces.
	// This field is optional. If you don't provide the usagePolicy, then it can be used only from the current namespace.
	// +optional
	UsagePolicy *apis.UsagePolicy `json:"usagePolicy,omitempty"`
}

// BackupConfigurationTemplate specifies the template for the BackupConfiguration created by the BackupBlueprint.
type BackupConfigurationTemplate struct {
	// Namespace specifies the namespace of the BackupConfiguration.
	// The field is optional. If you don't provide the namespace, then BackupConfiguration will be created in the BackupBlueprint namespace.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Backends specifies a list of storage references where the backed up data will be stored.
	// The respective BackupStorages can be in a different namespace than the BackupConfiguration.
	// However, it must be allowed by the `usagePolicy` of the BackupStorage to refer from this namespace.
	//
	// This field is optional, if you don't provide any backend here, KubeStash will use the default BackupStorage for the namespace.
	// If a default BackupStorage does not exist in the same namespace, then KubeStash will look for a default BackupStorage
	// in other namespaces that allows using it from the BackupConfiguration namespace.
	// +optional
	Backends []BackendReference `json:"backends,omitempty"`

	// Sessions specifies a list of session template for backup. You can use custom variables
	// in your template then provide the variable value through annotations.
	Sessions []Session `json:"sessions,omitempty"`

	// DeletionPolicy specifies whether the BackupConfiguration will be deleted on BackupBlueprint deletion
	// This field is optional, if you don't provide deletionPolicy, then BackupConfiguration will not be deleted on BackupBlueprint deletion
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`
}

// DeletionPolicy specifies whether the BackupConfiguration will be deleted on BackupBlueprint deletion
// +kubebuilder:validation:Enum=OnDelete
type DeletionPolicy string

const (
	DeletionPolicyOnDelete DeletionPolicy = "OnDelete"
)

//+kubebuilder:object:root=true

// BackupBlueprintList contains a list of BackupBlueprint
type BackupBlueprintList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupBlueprint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupBlueprint{}, &BackupBlueprintList{})
}
