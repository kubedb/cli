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
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

type Interface interface {
	metav1.Object
	GetInit() *InitSpec
	GetStatus() DatabaseStatus
}

type InitSpec struct {
	// Initialized indicates that this database has been initialized.
	// This will be set by the operator to ensure
	// that database is not mistakenly reset when recovered using disaster recovery tools.
	Initialized bool              `json:"initialized"`
	Script      *ScriptSourceSpec `json:"script,omitempty"`

	// Snapshot contains the restore-related details
	Snapshot *SnapshotSourceSpec `json:"snapshot,omitempty"`
}

type ScriptSourceSpec struct {
	ScriptPath        string `json:"scriptPath,omitempty"`
	core.VolumeSource `json:",inline,omitempty"`
	// This will take some database related config from the user
	PodTemplate *core.PodTemplateSpec `json:"podTemplate,omitempty"`
}

type SnapshotSourceSpec struct {
	Repository kmapi.TypedObjectReference `json:"repository,omitempty"`
	// +kubebuilder:default="latest"
	SnapshotID string `json:"snapshotID,omitempty"`
}

// DatabaseStatus defines the observed state of schema api types
type DatabaseStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase DatabaseSchemaPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
	// Database authentication secret
	// +optional
	AuthSecret *core.LocalObjectReference `json:"authSecret,omitempty"`
}

// +kubebuilder:validation:Enum=Delete;DoNotDelete
type DeletionPolicy string

type VaultSecretEngineRole struct {
	Subjects []rbac.Subject `json:"subjects"`
	// +optional
	DefaultTTL string `json:"defaultTTL,omitempty"`
	// +optional
	MaxTTL string `json:"maxTTL,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;InProgress;Terminating;Current;Failed;Expired
type DatabaseSchemaPhase string

type (
	DatabaseSchemaConditionType string
	DatabaseSchemaMessage       string
)
