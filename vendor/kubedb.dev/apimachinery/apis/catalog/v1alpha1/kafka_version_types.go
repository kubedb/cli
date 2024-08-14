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
	ResourceCodeKafkaVersion     = "kfversion"
	ResourceKindKafkaVersion     = "KafkaVersion"
	ResourceSingularKafkaVersion = "kafkaversion"
	ResourcePluralKafkaVersion   = "kafkaversions"
)

// KafkaVersion defines a Kafka database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=kafkaversions,singular=kafkaversion,scope=Cluster,shortName=kfversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type KafkaVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KafkaVersionSpec `json:"spec,omitempty"`
}

// KafkaVersionSpec is the spec for kafka version
type KafkaVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB KafkaVersionDatabase `json:"db"`
	// Connect Image
	ConnectCluster ConnectClusterVersion `json:"connectCluster"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Database Image
	CruiseControl CruiseControlVersionDatabase `json:"cruiseControl"`
	// PSP names
	// +optional
	PodSecurityPolicies KafkaVersionPodSecurityPolicy `json:"podSecurityPolicies"`
	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// KafkaVersionDatabase is the Kafka Database image
type KafkaVersionDatabase struct {
	Image string `json:"image"`
}

// ConnectClusterVersion is the Kafka Connect Cluster image
type ConnectClusterVersion struct {
	Image string `json:"image"`
}

// CruiseControlVersionDatabase is the Kafka Database image
type CruiseControlVersionDatabase struct {
	Image string `json:"image"`
}

// KafkaVersionPodSecurityPolicy is the Kafka pod security policies
type KafkaVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KafkaVersionList is a list of KafkaVersions
type KafkaVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of RedisVersion CRD objects
	Items []KafkaVersion `json:"items,omitempty"`
}
