/*
Copyright 2024.

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
	ResourceCodeIgniteVersion     = "igversion"
	ResourceKindIgniteVersion     = "IgniteVersion"
	ResourceSingularIgniteVersion = "igniteversion"
	ResourcePluralIgniteVersion   = "igniteversions"
)

// IgniteVersion defines a Ignite database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=igniteversions,singular=igniteversion,scope=Cluster,shortName=igversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type IgniteVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IgniteVersionSpec `json:"spec,omitempty"`
}

// IgniteVersionSpec is the spec for Ignite version
type IgniteVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB IgniteVersionDatabase `json:"db"`
	// Database Image
	InitContainer IgniteInitContainer `json:"initContainer,omitempty"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext IgniteSecurityContext `json:"securityContext"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`

	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
}

// IgniteSecurityContext is for the additional config for the DB container
type IgniteSecurityContext struct {
	RunAsUser  *int64 `json:"runAsUser,omitempty"`
	RunAsGroup *int64 `json:"runAsGroup,omitempty"`
}

// IgniteVersionDatabase is the Ignite Database image
type IgniteVersionDatabase struct {
	Image string `json:"image"`
}

// IgniteInitContainer is the Ignite init Container image
type IgniteInitContainer struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IgniteVersionList is a list of IgniteVersions
type IgniteVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of IgniteVersion CRD objects
	Items []IgniteVersion `json:"items"`
}
