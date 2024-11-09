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
	ResourceCodeMariaDBVersion     = "mariaversion"
	ResourceKindMariaDBVersion     = "MariaDBVersion"
	ResourceSingularMariaDBVersion = "mariadbversion"
	ResourcePluralMariaDBVersion   = "mariadbversions"
)

// MariaDBVersion defines a MariaDB (percona variation for MariaDB database) version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mariadbversions,singular=mariadbversion,scope=Cluster,shortName=mdversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MariaDBVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MariaDBVersionSpec `json:"spec,omitempty"`
}

// MariaDBVersionSpec is the spec for MariaDB version
type MariaDBVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB MariaDBVersionDatabase `json:"db"`
	// Exporter Image
	Exporter MariaDBVersionExporter `json:"exporter"`
	// Coordinator Image
	Coordinator MariaDBVersionCoordinator `json:"coordinator,omitempty"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Init container Image
	// TODO: remove if not needed
	InitContainer MariaDBVersionInitContainer `json:"initContainer"`
	// PSP names
	PodSecurityPolicies MariaDBVersionPodSecurityPolicy `json:"podSecurityPolicies"`
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
	// Archiver defines the walg & stash-addon related specifications
	Archiver ArchiverSpec `json:"archiver,omitempty"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// MariaDBVersionDatabase is the mariadb image
type MariaDBVersionDatabase struct {
	Image string `json:"image"`
}

// MariaDBVersionExporter is the image for the MariaDB exporter
type MariaDBVersionExporter struct {
	Image string `json:"image"`
}

// MariaDBVersionInitContainer is the MariaDB Container initializer
type MariaDBVersionInitContainer struct {
	Image string `json:"image"`
}

// MariaDBVersionCoordinator is the MariaDB Coordinator image
type MariaDBVersionCoordinator struct {
	Image string `json:"image"`
}

// MariaDBVersionPodSecurityPolicy is the MariaDB pod security policies
type MariaDBVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MariaDBVersionList is a list of MariaDBVersions
type MariaDBVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MariaDBVersion CRD objects
	Items []MariaDBVersion `json:"items,omitempty"`
}
