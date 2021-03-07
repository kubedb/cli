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
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MySQLVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MySQLVersionSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// MySQLVersionSpec is the spec for MySQL version
type MySQLVersionSpec struct {
	// Version
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`
	// Distribution
	Distribution MySQLDistro `json:"distribution,omitempty" protobuf:"bytes,2,opt,name=distribution,casttype=MySQLDistro"`
	// Database Image
	DB MySQLVersionDatabase `json:"db" protobuf:"bytes,3,opt,name=db"`
	// Exporter Image
	Exporter MySQLVersionExporter `json:"exporter" protobuf:"bytes,4,opt,name=exporter"`
	// ReplicationModeDetector Image
	ReplicationModeDetector ReplicationModeDetector `json:"replicationModeDetector" protobuf:"bytes,5,opt,name=replicationModeDetector"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,6,opt,name=deprecated"`
	// Init container Image
	InitContainer MySQLVersionInitContainer `json:"initContainer" protobuf:"bytes,7,opt,name=initContainer"`
	// PSP names
	PodSecurityPolicies MySQLVersionPodSecurityPolicy `json:"podSecurityPolicies" protobuf:"bytes,8,opt,name=podSecurityPolicies"`
	//upgrade constraints
	UpgradeConstraints MySQLUpgradeConstraints `json:"upgradeConstraints" protobuf:"bytes,9,opt,name=upgradeConstraints"`
	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty" protobuf:"bytes,10,opt,name=stash"`
}

// MySQLVersionDatabase is the MySQL Database image
type MySQLVersionDatabase struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// MySQLVersionExporter is the image for the MySQL exporter
type MySQLVersionExporter struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// MySQLVersionInitContainer is the MySQL Container initializer
type MySQLVersionInitContainer struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// MySQLVersionPodSecurityPolicy is the MySQL pod security policies
type MySQLVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName" protobuf:"bytes,1,opt,name=databasePolicyName"`
}

type MySQLUpgradeConstraints struct {
	// List of all accepted versions for uprade request
	Allowlist Allowlist `json:"allowlist,omitempty" protobuf:"bytes,1,opt,name=allowlist"`
	// List of all rejected versions for uprade request
	Denylist Denylist `json:"denylist,omitempty" protobuf:"bytes,2,opt,name=denylist"`
}

type Allowlist struct {
	// List of all accepted versions for upgrade request of a Standalone server. empty indicates all accepted
	Standalone []string `json:"standalone,omitempty" protobuf:"bytes,1,opt,name=standalone"`
	// List of all accepted versions for upgrade request of a GroupReplication cluster. empty indicates all accepted
	GroupReplication []string `json:"groupReplication,omitempty" protobuf:"bytes,2,opt,name=groupReplication"`
}

type Denylist struct {
	// List of all rejected versions for upgrade request of a Standalone server
	Standalone []string `json:"standalone,omitempty" protobuf:"bytes,1,opt,name=standalone"`
	// List of all rejected versions for upgrade request of a GroupReplication cluster
	GroupReplication []string `json:"groupReplication,omitempty" protobuf:"bytes,2,opt,name=groupReplication"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLVersionList is a list of MySQLVersions
type MySQLVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of MySQLVersion CRD objects
	Items []MySQLVersion `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

// +kubebuilder:validation:Enum=Oracle;Percona
type MySQLDistro string

const (
	MySQLDistroOracle  MySQLDistro = "Oracle"
	MySQLDistroPercona MySQLDistro = "Percona"
)
