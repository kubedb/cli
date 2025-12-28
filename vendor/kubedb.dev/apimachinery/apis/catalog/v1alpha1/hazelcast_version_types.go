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

// HazelcastVersionSpec defines the desired state of HazelcastVersion.

const (
	ResourceCodeHazelcastVersion     = "hzversion"
	ResourceKindHazelcastVersion     = "HazelcastVersion"
	ResourceSingularHazelcastVersion = "hazelcastversion"
	ResourcePluralHazelcastVersion   = "hazelcastversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=hazelcastversions,singular=hazelcastversion,scope=Cluster,shortName=hzversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type HazelcastVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec HazelcastVersionSpec `json:"spec,omitempty"`
}
type HazelcastVersionSpec struct {
	// Version
	Version string `json:"version"`

	// EndOfLife refers if this version reached into its end of the life or not, based on https://endoflife.date/
	// +optional
	EndOfLife bool `json:"endOfLife"`

	// Database Image
	DB HazelcastVersionDatabase `json:"db"`
	// Database Image
	InitContainer HazelcastInitContainer `json:"initContainer"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// SecurityContext is for the additional security information for the Hazelcast container
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
}

// HazelcastVersionDatabase is the Hazelcast Database image
type HazelcastVersionDatabase struct {
	Image string `json:"image"`
}

// HazelcastInitContainer is the Hazelcast init Container image
type HazelcastInitContainer struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HazelcastVersionList contains a list of HazelcastVersion.
type HazelcastVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HazelcastVersion `json:"items"`
}
