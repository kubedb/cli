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

// +k8s:openapi-gen=true
// +kubebuilder:object:generate=true
package apis

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

// Driver specifies the name of underlying tool that is being used to upload the backed up data.
// +kubebuilder:validation:Enum=Restic;WalG;VolumeSnapshotter;Solr;Medusa
type Driver string

const (
	DriverRestic            Driver = "Restic"
	DriverWalG              Driver = "WalG"
	DriverMedusa            Driver = "Medusa"
	DriverVolumeSnapshotter Driver = "VolumeSnapshotter"
	DriverSolr              Driver = "Solr"
)

// VolumeSource specifies the source of volume to mount in the backup/restore executor
// +k8s:openapi-gen=true
type VolumeSource struct {
	ofst.VolumeSource `json:",inline"`

	// VolumeClaimTemplate specifies a template for volume to use by the backup/restore executor
	// +optional
	VolumeClaimTemplate *ofst.PersistentVolumeClaimTemplate `json:"volumeClaimTemplate,omitempty"`
}

// ParameterDefinition defines the parameter names, their usage, their requirements etc.
// +k8s:openapi-gen=true
type ParameterDefinition struct {
	// Name specifies the name of the parameter
	Name string `json:"name,omitempty"`

	// Usage specifies the usage of this parameter
	Usage string `json:"usage,omitempty"`

	// Required specify whether this parameter is required or not
	// +optional
	Required bool `json:"required,omitempty"`

	// Default specifies a default value for the parameter
	// +optional
	Default string `json:"default,omitempty"`
}

// UsagePolicy specifies a policy that restrict the usage of a resource across namespaces.
// +k8s:openapi-gen=true
type UsagePolicy struct {
	// AllowedNamespaces specifies which namespaces are allowed to use the resource
	// +optional
	AllowedNamespaces AllowedNamespaces `json:"allowedNamespaces,omitempty"`
}

// AllowedNamespaces indicate which namespaces the resource should be selected from.
// +k8s:openapi-gen=true
type AllowedNamespaces struct {
	// From indicates how to select the namespaces that are allowed to use this resource.
	// Possible values are:
	// * All: All namespaces can use this resource.
	// * Selector: Namespaces that matches the selector can use this resource.
	// * Same: Only current namespace can use the resource.
	//
	// +optional
	// +kubebuilder:default=Same
	From *FromNamespaces `json:"from,omitempty"`

	// Selector must be specified when From is set to "Selector". In that case,
	// only the selected namespaces are allowed to use this resource.
	// This field is ignored for other values of "From".
	//
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// FromNamespaces specifies namespace from which namespaces are allowed to use the resource.
// +kubebuilder:validation:Enum=All;Selector;Same
type FromNamespaces string

const (
	// NamespacesFromAll specifies that all namespaces can use the resource.
	NamespacesFromAll FromNamespaces = "All"

	// NamespacesFromSelector specifies that only the namespace that matches the selector can use the resource.
	NamespacesFromSelector FromNamespaces = "Selector"

	// NamespacesFromSame specifies that only the current namespace can use the resource.
	NamespacesFromSame FromNamespaces = "Same"
)
