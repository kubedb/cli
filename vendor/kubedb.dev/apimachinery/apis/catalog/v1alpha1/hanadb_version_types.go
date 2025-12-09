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
	ResourceCodeHanaDBVersion     = "hdbversion"
	ResourceKindHanaDBVersion     = "HanaDBVersion"
	ResourceSingularHanaDBVersion = "hanadbversion"
	ResourcePluralHanaDBVersion   = "hanadbversions"
)

// HanaDBVersion defines a HanaDB database version
// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=hanadbversions,singular=hanadbversion,scope=Cluster,shortName=hdbversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type HanaDBVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HanaDBVersionSpec `json:"spec,omitempty"`
}

// HanaDBVersionSpec is the spec for  HanaDB version
type HanaDBVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB HanaDatabase `json:"db"`
	// Deprecated versions usable but considered as obsolete and best avoided typically superseded
	Deprecated bool `json:"deprecated,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext HanaDBSecurityContext `json:"securityContext"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
}

// HanaDBSecurityContext is for the additional config for the DB container
type HanaDBSecurityContext struct {
	RunAsUser  *int64 `json:"runAsUser,omitempty"`
	RunAsGroup *int64 `json:"runAsGroup,omitempty"`
}

// HanaDBVersionDatabase is the HanaDB Database image
type HanaDatabase struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HanaDBVersionList is a list of HanaDBVersions
type HanaDBVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HanaDBVersion `json:"items,omitempty"`
}
