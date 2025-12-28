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
	ResourceCodeNeo4jVersion     = "neoversion"
	ResourceKindNeo4jVersion     = "Neo4jVersion"
	ResourceSingularNeo4jVersion = "neo4jversion"
	ResourcePluralNeo4jVersion   = "neo4jversions"
)

// Neo4jVersion defines a Neo4j database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=neo4jversions,singular=neo4jversion,scope=Cluster,shortName=neoversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Neo4jVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Neo4jVersionSpec `json:"spec,omitempty"`
}

// Neo4jVersionSpec is the spec for Neo4j version
type Neo4jVersionSpec struct {
	// Version
	Version string `json:"version"`

	// EndOfLife refers if this version reached into its end of the life or not, based on https://endoflife.date/
	// +optional
	EndOfLife bool `json:"endOfLife"`

	// Database Image
	DB Neo4jVersionDatabase `json:"db"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
}

// Neo4jVersionDatabase is the Neo4j Database image
type Neo4jVersionDatabase struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Neo4jVersionList is a list of Neo4jVersions
type Neo4jVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Neo4jVersion CRD objects
	Items []Neo4jVersion `json:"items,omitempty"`
}
