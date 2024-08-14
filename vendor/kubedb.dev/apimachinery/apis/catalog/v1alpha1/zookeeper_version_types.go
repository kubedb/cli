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
	ResourceCodeZooKeeperVersion     = "zkversion"
	ResourceKindZooKeeperVersion     = "ZooKeeperVersion"
	ResourceSingularZooKeeperVersion = "zookeeperversion"
	ResourcePluralZooKeeperVersion   = "zookeeperversions"
)

// ZooKeeperVersion defines a ZooKeeper database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=zookeeperversions,singular=zookeeperversion,scope=Cluster,shortName=zkversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ZooKeeperVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ZooKeeperVersionSpec `json:"spec,omitempty"`
}

// ZooKeeperVersionSpec is the spec for zookeeper version
type ZooKeeperVersionSpec struct {
	// Version
	Version string `json:"version"`
	// init container image
	// +optional
	InitContainer ZooKeeperVersionInitContainer `json:"initContainer,omitempty"`
	// Database Image
	DB ZooKeeperVersionDatabase `json:"db"`
	// Exporter Image
	// +optional
	Exporter ZooKeeperVersionExporter `json:"exporter"`
	// Coordinator Image
	Coordinator ZooKeeperVersionCoordinator `json:"coordinator,omitempty"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// PSP names
	// +optional
	PodSecurityPolicies ZooKeeperVersionPodSecurityPolicy `json:"podSecurityPolicies"`
	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`
	// update constraints
	// +optional
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
	// +optional
	GitSyncer GitSyncer `json:"gitSyncer,omitempty"`
}

// ZooKeeperVersionInitContainer is the ZooKeeper init container image
type ZooKeeperVersionInitContainer struct {
	Image string `json:"image"`
}

// ZooKeeperVersionDatabase is the ZooKeeper Database image
type ZooKeeperVersionDatabase struct {
	Image string `json:"image"`
}

// ZooKeeperVersionCoordinator is the ZooKeeper coordinator image
type ZooKeeperVersionCoordinator struct {
	Image string `json:"image"`
}

// ZooKeeperVersionExporter is the image for the ZooKeeper exporter
type ZooKeeperVersionExporter struct {
	Image string `json:"image"`
}

// ZooKeeperVersionPodSecurityPolicy is the ZooKeeper pod security policies
type ZooKeeperVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZooKeeperVersionList is a list of ZooKeeperVersions
type ZooKeeperVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ZooKeeperVersion CRD objects
	Items []ZooKeeperVersion `json:"items,omitempty"`
}
