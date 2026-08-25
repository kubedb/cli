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
	ResourceCodeEtcdVersion     = "etcdversion"
	ResourceKindEtcdVersion     = "EtcdVersion"
	ResourceSingularEtcdVersion = "etcdversion"
	ResourcePluralEtcdVersion   = "etcdversions"
)

// EtcdVersion defines an Etcd database version.
//
// Unlike most other databases, etcd has no vendor distributions -- there is a
// single upstream implementation -- so this catalog type deliberately has no
// spec.distribution field.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=etcdversions,singular=etcdversion,scope=Cluster,shortName=etcdversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type EtcdVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              EtcdVersionSpec `json:"spec,omitempty"`
}

// EtcdVersionSpec is the spec for etcd version
type EtcdVersionSpec struct {
	// Version
	Version string `json:"version"`

	// EndOfLife refers if this version reached into its end of the life or not, based on https://endoflife.date/
	// +optional
	EndOfLife bool `json:"endOfLife"`

	// init container image
	// +optional
	InitContainer EtcdVersionInitContainer `json:"initContainer,omitempty"`

	// Database Image
	DB EtcdVersionDatabase `json:"db"`

	// Exporter Image
	Exporter EtcdVersionExporter `json:"exporter"`

	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`

	// PSP names
	// +optional
	PodSecurityPolicies EtcdVersionPodSecurityPolicy `json:"podSecurityPolicies"`

	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`

	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext EtcdSecurityContext `json:"securityContext"`

	// update constraints
	// +optional
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`

	// +optional
	GitSyncer GitSyncer `json:"gitSyncer,omitempty"`

	// Archiver defines the walg & kube-stash-addon related specifications
	// +optional
	Archiver ArchiverSpec `json:"archiver,omitempty"`

	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// EtcdVersionInitContainer is the Etcd init container image
type EtcdVersionInitContainer struct {
	Image string `json:"image"`
}

// EtcdVersionDatabase is the Etcd Database image
type EtcdVersionDatabase struct {
	Image string `json:"image"`
	// +optional
	BaseOS string `json:"baseOS,omitempty"`
}

// EtcdVersionExporter is the image for the Etcd exporter
type EtcdVersionExporter struct {
	Image string `json:"image"`
}

// EtcdVersionPodSecurityPolicy is the Etcd pod security policies
type EtcdVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// EtcdSecurityContext is the additional features for the Etcd DB container
type EtcdSecurityContext struct {
	// RunAsUser is the default UID for the DB container. The upstream etcd image
	// runs as uid 1000.
	// +optional
	RunAsUser *int64 `json:"runAsUser,omitempty"`

	// RunAsAnyNonRoot will be true if the user can change the default db container
	// user to something other than the image's default user.
	// +optional
	RunAsAnyNonRoot bool `json:"runAsAnyNonRoot,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdVersionList is a list of EtcdVersions
type EtcdVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of EtcdVersion CRD objects
	Items []EtcdVersion `json:"items,omitempty"`
}
