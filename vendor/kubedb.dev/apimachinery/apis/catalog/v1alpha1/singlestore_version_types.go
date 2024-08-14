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
	ResourceCodeSinglestoreVersion     = "sdbv"
	ResourceKindSinlestoreVersion      = "SinglestoreVersion"
	ResourceSingularSinglestoreVersion = "singlestoreversion"
	ResourcePluralSinglestoreVersion   = "singlestoreversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=singlestoreversions,singular=singlestoreversion,scope=Cluster,shortName=sdbv,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SinglestoreVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SinglestoreVersionSpec `json:"spec,omitempty"`
}

// SinglestoreVersionSpec defines the desired state of SinglestoreVersion
type SinglestoreVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB SinglestoreVersionDatabase `json:"db"`
	// +optional
	Coordinator SinglestoreCoordinator `json:"coordinator,omitempty"`
	// +optional
	Standalone SinglestoreStandaloneVersionDatabase `json:"standalone,omitempty"`
	// +optional
	InitContainer SinglestoreInitContainer `json:"initContainer,omitempty"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`
	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SinglestoreSecurityContext `json:"securityContext"`
}

// SinglestoreSecurityContext is for the additional config for the DB container
type SinglestoreSecurityContext struct {
	RunAsUser  *int64 `json:"runAsUser,omitempty"`
	RunAsGroup *int64 `json:"runAsGroup,omitempty"`
}

// SinglestoreVersionDatabase is the Singlestore Cluster Database image
type SinglestoreVersionDatabase struct {
	Image string `json:"image"`
}

// SinglestoreVersionDatabase is the Singlestore Standalone Database image
type SinglestoreStandaloneVersionDatabase struct {
	Image string `json:"image"`
}

// SinglestoreCoordinator is the Singlestore coordinator Container image
type SinglestoreCoordinator struct {
	Image string `json:"image"`
}

// SinglestoreInitContainer is the Singlestore init Container image
type SinglestoreInitContainer struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SinglestoreVersionList contains a list of SinglestoreVersions
type SinglestoreVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SinglestoreVersion `json:"items"`
}
