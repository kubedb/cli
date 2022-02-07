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
	ResourceCodeMongoDBVersion     = "mgversion"
	ResourceKindMongoDBVersion     = "MongoDBVersion"
	ResourceSingularMongoDBVersion = "mongodbversion"
	ResourcePluralMongoDBVersion   = "mongodbversions"
)

// MongoDBVersion defines a MongoDB database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mongodbversions,singular=mongodbversion,scope=Cluster,shortName=mgversion,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Distribution",type="string",JSONPath=".spec.distribution"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MongoDBVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MongoDBVersionSpec `json:"spec,omitempty"`
}

// MongoDBVersionSpec is the spec for mongodb version
type MongoDBVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Distribution
	Distribution MongoDBDistro `json:"distribution,omitempty"`
	// Database Image
	DB MongoDBVersionDatabase `json:"db"`
	// Exporter Image
	Exporter MongoDBVersionExporter `json:"exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Init container Image
	InitContainer MongoDBVersionInitContainer `json:"initContainer"`
	// PSP names
	PodSecurityPolicies MongoDBVersionPodSecurityPolicy `json:"podSecurityPolicies"`
	// ReplicationModeDetector Image
	ReplicationModeDetector ReplicationModeDetector `json:"replicationModeDetector"`
	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`
}

// MongoDBVersionDatabase is the MongoDB Database image
type MongoDBVersionDatabase struct {
	Image string `json:"image"`
}

// MongoDBVersionExporter is the image for the MongoDB exporter
type MongoDBVersionExporter struct {
	Image string `json:"image"`
}

// MongoDBVersionInitContainer is the Elasticsearch Container initializer
type MongoDBVersionInitContainer struct {
	Image string `json:"image"`
}

// MongoDBVersionPodSecurityPolicy is the MongoDB pod security policies
type MongoDBVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBVersionList is a list of MongoDBVersions
type MongoDBVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MongoDBVersion CRD objects
	Items []MongoDBVersion `json:"items,omitempty"`
}

// +kubebuilder:validation:Enum=Official;Percona;KubeDB;MongoDB
type MongoDBDistro string

const (
	MongoDBDistroOfficaial MongoDBDistro = "Official"
	MongoDBDistroPercona   MongoDBDistro = "Percona"
	MongoDBDistroKubeDB    MongoDBDistro = "KubeDB"
)
