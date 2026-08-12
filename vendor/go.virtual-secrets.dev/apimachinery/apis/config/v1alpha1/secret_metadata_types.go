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
)

const (
	ResourceKindSecretMetadata = "SecretMetadata"
	ResourceSecretMetadata     = "secretmetadata"
	ResourceSecretMetadatas    = "secretmetadatas"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=secretmetadatas,singular=secretmetadata,shortName=scmetadata,categories={meta,virtual-secrets,appscode}
// +kubebuilder:printcolumn:name="TYPE",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SecretMetadata struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SecretMetadataSpec `json:"spec,omitempty"`
}

// SecretMetadataSpec defines the desired state of SecretMetadata
type SecretMetadataSpec struct {
	// Name of the SecretStoreName object
	SecretStoreName string `json:"secretStoreName"`

	// Used to facilitate programmatic handling of secret data.
	// More meta: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
	// +optional
	Type core.SecretType `json:"type,omitempty"`

	// Immutable, if set to true, ensures that data stored in the Secret cannot
	// be updated (only object metadata can be modified).
	// If not set to true, the field can be modified at any time.
	// Defaulted to nil.
	// +optional
	Immutable *bool `json:"immutable,omitempty"`

	// DataLength specifies the count of data stored in the Virtual Secret
	// +optional
	DataLength int `json:"dataLength,omitempty"`

	// DataHash specifies the hash value of the data stored in the Virtual Secret
	// +optional
	DataHash string `json:"dataHash,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretMetadataList contains a list of SecretMetadata
type SecretMetadataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretMetadata `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretMetadata{})
	SchemeBuilder.Register(&SecretMetadataList{})
}
