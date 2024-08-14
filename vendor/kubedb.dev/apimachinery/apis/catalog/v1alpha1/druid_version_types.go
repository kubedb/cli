/*
Copyright 2023.

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
	ResourceCodeDruidVersion     = "drversion"
	ResourceKindDruidVersion     = "DruidVersion"
	ResourceSingularDruidVersion = "druidversion"
	ResourcePluralDruidVersion   = "druidversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=druidversions,singular=druidversion,scope=Cluster,shortName=drversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DruidVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DruidVersionSpec `json:"spec,omitempty"`
}

// DruidVersionSpec defines the desired state of DruidVersion
type DruidVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB DruidVersionDatabase `json:"db"`
	// Init Container Image
	InitContainer DruidInitContainer `json:"initContainer"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// SecurityContext is for the additional security information for the Druid container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// DruidVersionDatabase is the Druid Database image
type DruidVersionDatabase struct {
	Image string `json:"image"`
}

// Druid is the Druid Init Container image
type DruidInitContainer struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DruidVersionList contains a list of DruidVersion
type DruidVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DruidVersion `json:"items,omitempty"`
}
