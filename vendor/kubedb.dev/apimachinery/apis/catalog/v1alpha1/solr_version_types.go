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
	ResourceCodeSolrVersion     = "slversion"
	ResourceKindSolrVersion     = "SolrVersion"
	ResourceSingularSolrVersion = "Solrversion"
	ResourcePluralSolrVersion   = "solrversions"
)

// SolrVersion defines a Solr database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=solrversions,singular=solrversion,scope=Cluster,shortName=slversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SolrVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SolrVersionSpec `json:"spec,omitempty"`
}

// SolrVersionSpec defines the desired state of SolrVersion
type SolrVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB SolrVersionDatabase `json:"db"`
	// Database Image
	InitContainer SolrInitContainer `json:"initContainer"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// SecurityContext is for the additional security information for the Solr container
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
}

// SolrVersionDatabase is the Solr Database image
type SolrVersionDatabase struct {
	Image string `json:"image"`
}

// SolrInitContainer is the Solr init Container image
type SolrInitContainer struct {
	Image string `json:"image"`
}

// SolrVersionList contains a list of SolrVersion

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SolrVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SolrVersion `json:"items"`
}
