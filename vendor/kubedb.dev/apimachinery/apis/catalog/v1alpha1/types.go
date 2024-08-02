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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// ReplicationModeDetector is the image for the MySQL replication mode detector
type ReplicationModeDetector struct {
	Image string `json:"image"`
}

// UpdateConstraints specifies the constraints that need to be considered during version upgrade
type UpdateConstraints struct {
	// List of all accepted versions for upgrade request.
	// An empty list indicates all versions are accepted except the denylist.
	Allowlist []string `json:"allowlist,omitempty"`
	// List of all rejected versions for upgrade request.
	// An empty list indicates no version is rejected.
	Denylist []string `json:"denylist,omitempty"`
}

type ArchiverSpec struct {
	Walg  WalgSpec  `json:"walg,omitempty"`
	Addon AddonSpec `json:"addon,omitempty"`
}

type WalgSpec struct {
	Image string `json:"image"`
}

type AddonSpec struct {
	Name  AddonType  `json:"name,omitempty"`
	Tasks AddonTasks `json:"tasks,omitempty"`
}

// +kubebuilder:validation:Enum=mongodb-addon;postgres-addon;mysql-addon;mariadb-addon;mssqlserver-addon
type AddonType string

type AddonTasks struct {
	VolumeSnapshot    VolumeSnapshot    `json:"volumeSnapshot,omitempty"`
	ManifestBackup    ManifestBackup    `json:"manifestBackup,omitempty"`
	ManifestRestore   ManifestRestore   `json:"manifestRestore,omitempty"`
	FullBackup        FullBackup        `json:"fullBackup,omitempty"`
	FullBackupRestore FullBackupRestore `json:"fullBackupRestore,omitempty"`
}

type FullBackup struct {
	Name string `json:"name"`
}

type VolumeSnapshot struct {
	Name string `json:"name"`
}

type ManifestBackup struct {
	Name string `json:"name"`
}

type ManifestRestore struct {
	Name string `json:"name"`
}

type FullBackupRestore struct {
	Name string `json:"name"`
}

// GitSyncer is the image for the kubernetes/git-sync
// https://github.com/kubernetes/git-sync
type GitSyncer struct {
	Image string `json:"image"`
}

// SecurityContext is for the additional config for the DB container
type SecurityContext struct {
	RunAsUser *int64 `json:"runAsUser,omitempty"`
}

type ChartInfo struct {
	// Name specifies the name of the chart
	Name string `json:"name"`
	// Version specifies the version of the chart.
	Version string `json:"version,omitempty"`
	// Disable installing this chart
	// +optional
	Disable bool `json:"disable,omitempty"`
	// Values holds the values for this Helm release.
	// +optional
	Values *apiextensionsv1.JSON `json:"values,omitempty"`
}
