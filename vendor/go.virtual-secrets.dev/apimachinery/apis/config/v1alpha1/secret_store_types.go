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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindSecretStore = "SecretStore"
	ResourceSecretStore     = "secretstore"
	ResourceSecretStores    = "secretstores"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=secretstores,singular=secretstore,scope=Cluster,shortName=scsource,categories={meta,virtual-secrets,appscode}
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SecretStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SecretStoreSpec `json:"spec,omitempty"`
}

// SecretStoreSpec defines the desired state of SecretStore
type SecretStoreSpec struct {
	Vault *Vault `json:"vault,omitempty"`

	// +optional
	AWS *AWS `json:"aws,omitempty"`

	// +optional
	Azure *Azure `json:"azure,omitempty"`

	// +optional
	GCP *GCP `json:"gcp,omitempty"`

	// **For Dev Mode Only**
	// We can use a secret as the Secret Store for testing
	// +optional
	Secret *Secret `json:"secret,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretStoreList contains a list of SecretStore
type SecretStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretStore `json:"items"`
}

type Vault struct {
	// Connection url to the secret manager
	URL string `json:"url"`

	// Name of the vault role to use for the operator
	// +optional
	RoleName string `json:"roleName,omitempty"`
}

type AWS struct {
	// SecretRef defines a secret that contains the
	// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
	// +optional
	SecretRef *kmapi.ObjectReference `json:"secretRef,omitempty"`

	// Region specifies the AWS region where the Secret will be stored
	Region string `json:"region,omitempty"`
}

type Azure struct {
	// SecretRef defines a secret that contains the
	// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
	SecretRef *kmapi.ObjectReference `json:"secretRef,omitempty"`

	// The Azure Key Vault Provider offers four modes for accessing a Key Vault instance
	// Workload Identity, Pod Identity, Managed Identities, Service Principal
	// Pod Identity, Managed Identities are not supported yet
	// +kubebuilder:validation:Enum=WorkloadIdentity;ServicePrincipal
	AccessMode string `json:"accessMode,omitempty"`

	KeyVaultName string `json:"keyVaultName,omitempty"`
}

type GCP struct {
	// SecretRef defines a secret that contains the json file with all data
	// client_id, client_secret etc
	SecretRef *kmapi.ObjectReference `json:"secretRef,omitempty"`

	// Region specifies the GCP region where the Secret will be stored
	Region string `json:"region,omitempty"`
}

type Secret struct {
	// Name and namespace of the Secret which will work as the Secret Manager
	// +optional
	*kmapi.ObjectReference `json:",omitempty"`
}

func init() {
	SchemeBuilder.Register(&SecretStore{})
	SchemeBuilder.Register(&SecretStoreList{})
}
