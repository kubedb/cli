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
	ResourceCodePgBouncerVersion     = "pbversion"
	ResourceKindPgBouncerVersion     = "PgBouncerVersion"
	ResourceSingularPgBouncerVersion = "pgbouncerversion"
	ResourcePluralPgBouncerVersion   = "pgbouncerversions"
)

// PgBouncerVersion defines a PgBouncer database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pgbouncerversions,singular=pgbouncerversion,scope=Cluster,shortName=pbversion,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="PGBOUNCER_IMAGE",type="string",JSONPath=".spec.pgBouncer.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PgBouncerVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PgBouncerVersionSpec `json:"spec,omitempty"`
}

// PgBouncerVersionSpec is the spec for pgbouncer version
type PgBouncerVersionSpec struct {
	// Version
	Version string `json:"version"`
	// init container image
	InitContainer PgBouncerVersionInitContainer `json:"initContainer,omitempty"`
	// Database Image
	PgBouncer PgBouncerVersionDatabase `json:"pgBouncer"`
	// Exporter Image
	Exporter PgBouncerVersionExporter `json:"exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// upgrade constraints
	UpgradeConstraints UpgradeConstraints `json:"upgradeConstraints,omitempty"`
}

// PgBouncerVersionInitContainer is the PgBouncer init container image
type PgBouncerVersionInitContainer struct {
	Image string `json:"image"`
}

// PgBouncerVersionDatabase is the PgBouncer Database image
type PgBouncerVersionDatabase struct {
	Image string `json:"image"`
}

// PgBouncerVersionExporter is the image for the PgBouncer exporter
type PgBouncerVersionExporter struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgBouncerVersionList is a list of PgBouncerVersions
type PgBouncerVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PgBouncerVersion CRD objects
	Items []PgBouncerVersion `json:"items,omitempty"`
}
