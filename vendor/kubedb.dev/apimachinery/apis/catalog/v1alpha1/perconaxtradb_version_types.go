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
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceCodePerconaXtraDBVersion     = "pxversion"
	ResourceKindPerconaXtraDBVersion     = "PerconaXtraDBVersion"
	ResourceSingularPerconaXtraDBVersion = "perconaxtradbversion"
	ResourcePluralPerconaXtraDBVersion   = "perconaxtradbversions"
)

// PerconaXtraDBVersion defines a PerconaXtraDB (percona variation for MySQL database) version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=perconaxtradbversions,singular=perconaxtradbversion,scope=Cluster,shortName=pxversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PerconaXtraDBVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PerconaXtraDBVersionSpec `json:"spec,omitempty"`
}

// PerconaXtraDBVersionSpec is the spec for PerconaXtraDB version
type PerconaXtraDBVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB PerconaXtraDBVersionDatabase `json:"db"`
	// Exporter Image
	Exporter PerconaXtraDBVersionExporter `json:"exporter"`
	// Coordinator Image
	Coordinator PerconaXtraDBVersionCoordinator `json:"coordinator,omitempty"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Init container Image
	// TODO: remove if not needed
	InitContainer PerconaXtraDBVersionInitContainer `json:"initContainer"`
	// PSP names
	PodSecurityPolicies PerconaXtraDBVersionPodSecurityPolicy `json:"podSecurityPolicies"`
	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
	// +optional
	GitSyncer GitSyncer `json:"gitSyncer,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// PerconaXtraDBVersionDatabase is the perconaxtradb image
type PerconaXtraDBVersionDatabase struct {
	Image string `json:"image"`
}

// PerconaXtraDBVersionExporter is the image for the PerconaXtraDB exporter
type PerconaXtraDBVersionExporter struct {
	Image string `json:"image"`
}

// PerconaXtraDBVersionInitContainer is the PerconaXtraDB Container initializer
type PerconaXtraDBVersionInitContainer struct {
	Image string `json:"image"`
}

// PerconaXtraDBVersionCoordinator is the PerconaXtraDB Coordinator image
type PerconaXtraDBVersionCoordinator struct {
	Image string `json:"image"`
}

// PerconaXtraDBVersionPodSecurityPolicy is the PerconaXtraDB pod security policies
type PerconaXtraDBVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PerconaXtraDBVersionList is a list of PerconaXtraDBVersions
type PerconaXtraDBVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PerconaXtraDBVersion CRD objects
	Items []PerconaXtraDBVersion `json:"items,omitempty"`
}
