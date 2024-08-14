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
	ResourceCodeSchemaRegistryVersion     = "ksrversion"
	ResourceKindSchemaRegistryVersion     = "SchemaRegistryVersion"
	ResourceSingularSchemaRegistryVersion = "schemaregistryversion"
	ResourcePluralSchemaRegistryVersion   = "schemaregistryversions"
)

// SchemaRegistryVersion defines a SchemaRegistry version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=schemaregistryversions,singular=schemaregistryversion,scope=Cluster,shortName=ksrversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Distribution",type="string",JSONPath=".spec.distribution"
// +kubebuilder:printcolumn:name="REGISTRY_IMAGE",type="string",JSONPath=".spec.registry.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SchemaRegistryVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SchemaRegistryVersionSpec `json:"spec,omitempty"`
}

// SchemaRegistryVersionSpec is the spec for SchemaRegistry version
type SchemaRegistryVersionSpec struct {
	Distribution SchemaRegistryDistro `json:"distribution"`
	// Version
	Version string `json:"version"`
	// Registry Image
	Registry RegistryImage `json:"registry"`
	// Schema Registry In Memory Image
	InMemory ApicurioInMemory `json:"inMemory,omitempty"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
}

// RegistryImage is the SchemaRegistry image
type RegistryImage struct {
	Image string `json:"image"`
}

// ApicurioInMemory is the Apicurio Registry In-Memory image
type ApicurioInMemory struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SchemaRegistryVersionList is a list of SchemaRegistryVersion
type SchemaRegistryVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of SchemaRegistryVersion CRD objects
	Items []SchemaRegistryVersion `json:"items,omitempty"`
}

// +kubebuilder:validation:Enum=Apicurio;Aiven
type SchemaRegistryDistro string

const (
	SchemaRegistryDistroApicurio SchemaRegistryDistro = "Apicurio"
	SchemaRegistryDistroAiven    SchemaRegistryDistro = "Aiven"
)
