/*
Copyright 2025.

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
	ResourceCodeDB2Version     = "db2v"
	ResourceKindDB2Version     = "DB2Version"
	ResourceSingularDB2Version = "db2version"
	ResourcePluralDB2Version   = "db2versions"
)

// DB2Version defines a DB2 database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=db2versions,singular=db2version,scope=Cluster,shortName=db2v,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DB2Version struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DB2VersionSpec `json:"spec,omitempty"`
}

// DB2VersionSpec is the spec for oracle version
type DB2VersionSpec struct {
	// Version
	Version string `json:"version"`

	// EndOfLife refers if this version reached into its end of the life or not, based on https://endoflife.date/
	// +optional
	EndOfLife bool `json:"endOfLife"`

	// Database Image
	DB DB2VersionDatabase `json:"db"`
	// Coordinator Image
	// +optional
	Coordinator DB2Coordinator `json:"coordinator,omitempty"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`

	// SecurityContext is for the additional config for oracle DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`

	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`

	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// DB2VersionDatabase is the DB2 Database image
type DB2VersionDatabase struct {
	Image string `json:"image"`
}
type DB2Coordinator struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DB2VersionList is a list of DB2Versions
type DB2VersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of DB2Version CRD objects
	Items []DB2Version `json:"items,omitempty"`
}
