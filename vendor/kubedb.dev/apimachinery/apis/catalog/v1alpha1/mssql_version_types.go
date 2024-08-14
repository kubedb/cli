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
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceCodeMSSQLServerVersion     = "msversion"
	ResourceKindMSSQLServerVersion     = "MSSQLServerVersion"
	ResourceSingularMSSQLServerVersion = "mssqlserverversion"
	ResourcePluralMSSQLServerVersion   = "mssqlserverversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mssqlserverversions,singular=mssqlserverversion,scope=Cluster,shortName=msversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MSSQLServerVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MSSQLServerVersionSpec `json:"spec,omitempty"`
}

// MSSQLServerVersionSpec defines the desired state of MSSQLServer Version
type MSSQLServerVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB MSSQLServerDatabase `json:"db"`
	// Coordinator Image
	// +optional
	Coordinator MSSQLServerCoordinator `json:"coordinator,omitempty"`
	// Exporter Image
	Exporter MSSQLServerVersionExporter `json:"exporter"`
	// Init container Image
	InitContainer MSSQLServerInitContainer `json:"initContainer"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
	// Archiver defines the walg & kube-stash-addon related specifications
	Archiver ArchiverSpec `json:"archiver,omitempty"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// MSSQLServerDatabase is the MSSQLServer Database image
type MSSQLServerDatabase struct {
	Image string `json:"image"`
}

// MSSQLServerCoordinator is the MSSQLServer coordinator Container image
type MSSQLServerCoordinator struct {
	Image string `json:"image"`
}

// MSSQLServerVersionExporter is the image for the MSSQL Server exporter
type MSSQLServerVersionExporter struct {
	Image string `json:"image"`
}

// MSSQLServerInitContainer is the MSSQLServer Container initializer
type MSSQLServerInitContainer struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MSSQLServerVersionList contains a list of MSSQLServerVersion
type MSSQLServerVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MSSQLServerVersion `json:"items"`
}
