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
)

const (
	ResourceKindSecretMount = "SecretMount"
	ResourceSecretMount     = "secretmount"
	ResourceSecretMounts    = "secretmounts"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretMount holds secret data of a certain type. The total bytes of the values in
// the Data field must be less than MaxSecretMountSize bytes.
type SecretMount struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More meta: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretMountSpec  `json:"spec,omitempty"`
	Status SecretMountStaus `json:"status,omitempty"`
}

type SecretMountSpec struct {
	Name string `json:"name"`
}

type SecretMountStaus struct {
	Name string `json:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretMountList is a list of SecretMount.
type SecretMountList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More meta: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of secret objects.
	// More meta: https://kubernetes.io/docs/concepts/configuration/secret
	Items []SecretMount `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&SecretMount{})
	SchemeBuilder.Register(&SecretMountList{})
}
