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
	ResourceCodeMySQLVersion     = "myversion"
	ResourceKindMySQLVersion     = "MySQLVersion"
	ResourceSingularMySQLVersion = "mysqlversion"
	ResourcePluralMySQLVersion   = "mysqlversions"
)

// MySQLVersion defines a MySQL database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mysqlversions,singular=mysqlversion,scope=Cluster,shortName=myversion,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Distribution",type="string",JSONPath=".spec.distribution"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MySQLVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MySQLVersionSpec `json:"spec,omitempty"`
}

// MySQLVersionSpec is the spec for MySQL version
type MySQLVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Distribution
	Distribution MySQLDistro `json:"distribution,omitempty"`
	// Database Image
	DB MySQLVersionDatabase `json:"db"`
	// Exporter Image
	Exporter MySQLVersionExporter `json:"exporter"`
	// Coordinator Image
	// +optional
	Coordinator MySQLVersionCoordinator `json:"coordinator,omitempty"`
	// ReplicationModeDetector Image
	// +optional
	ReplicationModeDetector ReplicationModeDetector `json:"replicationModeDetector,omitempty"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Init container Image
	InitContainer MySQLVersionInitContainer `json:"initContainer"`
	// PSP names
	PodSecurityPolicies MySQLVersionPodSecurityPolicy `json:"podSecurityPolicies"`
	// upgrade constraints
	UpgradeConstraints MySQLUpgradeConstraints `json:"upgradeConstraints,omitempty"`
	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`
	// Router image
	// +optional
	Router MySQLVersionRouter `json:"router,omitempty"`
	// +optional
	RouterInitContainer MySQLVersionRouterInitContainer `json:"routerInitContainer,omitempty"`
}

// MySQLVersionDatabase is the MySQL Database image
type MySQLVersionDatabase struct {
	Image string `json:"image"`
}

// MySQLVersionExporter is the image for the MySQL exporter
type MySQLVersionExporter struct {
	Image string `json:"image"`
}

// MySQLVersionCoordinator is the image for coordinator
type MySQLVersionCoordinator struct {
	Image string `json:"image"`
}

// MySQLVersionInitContainer is the MySQL Container initializer
type MySQLVersionInitContainer struct {
	Image string `json:"image"`
}

// MySQLVersionRouter is the MySQL Router lightweight middleware
// that provides transparent routing between your application and back-end MySQL Servers
type MySQLVersionRouter struct {
	Image string `json:"image"`
}

// MySQLVersionRouterInitContainer is mysql router init container
type MySQLVersionRouterInitContainer struct {
	Image string `json:"image"`
}

// MySQLVersionPodSecurityPolicy is the MySQL pod security policies
type MySQLVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

type MySQLUpgradeConstraints struct {
	// List of all accepted versions for upgrade request
	Allowlist MySQLVersionAllowlist `json:"allowlist,omitempty"`
	// List of all rejected versions for upgrade request
	Denylist MySQLVersionDenylist `json:"denylist,omitempty"`
}

type MySQLVersionAllowlist struct {
	// List of all accepted versions for upgrade request of a Standalone server. empty indicates all accepted
	Standalone []string `json:"standalone,omitempty"`
	// List of all accepted versions for upgrade request of a GroupReplication cluster. empty indicates all accepted
	GroupReplication []string `json:"groupReplication,omitempty"`
}

type MySQLVersionDenylist struct {
	// List of all rejected versions for upgrade request of a Standalone server
	Standalone []string `json:"standalone,omitempty"`
	// List of all rejected versions for upgrade request of a GroupReplication cluster
	GroupReplication []string `json:"groupReplication,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLVersionList is a list of MySQLVersions
type MySQLVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MySQLVersion CRD objects
	Items []MySQLVersion `json:"items,omitempty"`
}

// +kubebuilder:validation:Enum=Official;Oracle;Percona;KubeDB;MySQL
type MySQLDistro string

const (
	MySQLDistroOfficial MySQLDistro = "Official"
	MySQLDistroMySQL    MySQLDistro = "MySQL"
	MySQLDistroPercona  MySQLDistro = "Percona"
	MySQLDistroKubeDB   MySQLDistro = "KubeDB"
)
