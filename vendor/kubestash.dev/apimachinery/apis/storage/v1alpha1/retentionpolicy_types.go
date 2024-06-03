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
	"kubestash.dev/apimachinery/apis"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindRetentionPolicy     = "RetentionPolicy"
	ResourceSingularRetentionPolicy = "retentionpolicy"
	ResourcePluralRetentionPolicy   = "retentionpolicies"
)

// +k8s:openapi-gen=true
//+kubebuilder:object:root=true
// +kubebuilder:resource:path=retentionpolicies,singular=retentionpolicy,categories={kubestash,appscode}
// +kubebuilder:printcolumn:name="Max-Retention-Period",type="string",JSONPath=".spec.maxRetentionPeriod"
// +kubebuilder:printcolumn:name="Default",type="boolean",JSONPath=".spec.default"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// RetentionPolicy specifies how the old Snapshots should be cleaned up.
// This is a namespaced CRD. However, you can refer it from other namespaces
// as long as it is permitted via `.spec.usagePolicy`.
type RetentionPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RetentionPolicySpec `json:"spec,omitempty"`
}

// RetentionPolicySpec defines the policy of cleaning old Snapshots
type RetentionPolicySpec struct {
	// MaxRetentionPeriod specifies a duration up to which the old Snapshots should be kept.
	// KubeStash will remove all the Snapshots that are older than the MaxRetentionPeriod.
	// For example, MaxRetentionPeriod of `30d` will keep only the Snapshots of last 30 days.
	// Sample duration format:
	// - years: 	2y
	// - months: 	6mo
	// - days: 		30d
	// - hours: 	12h
	// - minutes: 	30m
	// You can also combine the above durations. For example: 30d12h30m
	// +optional
	MaxRetentionPeriod RetentionPeriod `json:"maxRetentionPeriod,omitempty"`

	// UsagePolicy specifies a policy of how this RetentionPolicy will be used. For example, you can use `allowedNamespaces`
	// policy to restrict the usage of this RetentionPolicy to particular namespaces.
	// This field is optional. If you don't provide the usagePolicy, then it can be used only from the current namespace.
	// +optional
	UsagePolicy *apis.UsagePolicy `json:"usagePolicy,omitempty"`

	// SuccessfulSnapshots specifies how many successful Snapshots should be kept.
	// +optional
	SuccessfulSnapshots *SuccessfulSnapshotsKeepPolicy `json:"successfulSnapshots,omitempty"`

	// FailedSnapshots specifies how many failed Snapshots should be kept.
	// +optional
	FailedSnapshots *FailedSnapshotsKeepPolicy `json:"failedSnapshots,omitempty"`

	// Default specifies whether to use this RetentionPolicy as a default RetentionPolicy for
	// the current namespace as well as the permitted namespaces.
	// One namespace can have at most one default RetentionPolicy configured.
	// +optional
	Default bool `json:"default,omitempty"`
}

// RetentionPeriod represents a duration in the format "1y2mo3w4d5h6m", where
// y=year, mo=month, w=week, d=day, h=hour, m=minute.
type RetentionPeriod string

// SuccessfulSnapshotsKeepPolicy specifies the policy for keeping successful Snapshots
type SuccessfulSnapshotsKeepPolicy struct {
	// Last specifies how many last Snapshots should be kept.
	// +optional
	Last *int32 `json:"last,omitempty"`

	// Hourly specifies how many hourly Snapshots should be kept.
	// +optional
	Hourly *int32 `json:"hourly,omitempty"`

	// Daily specifies how many daily Snapshots should be kept.
	// +optional
	Daily *int32 `json:"daily,omitempty"`

	// Weekly specifies how many weekly Snapshots should be kept.
	// +optional
	Weekly *int32 `json:"weekly,omitempty"`

	// Monthly specifies how many monthly Snapshots should be kept.
	// +optional
	Monthly *int32 `json:"monthly,omitempty"`

	// Yearly specifies how many yearly Snapshots should be kept.
	// +optional
	Yearly *int32 `json:"yearly,omitempty"`
}

// FailedSnapshotsKeepPolicy specifies the policy for keeping failed Snapshots
type FailedSnapshotsKeepPolicy struct {
	// Last specifies how many last failed Snapshots should be kept.
	// By default, KubeStash will keep only the last 1 failed Snapshot.
	// +kubebuilder:default=1
	// +optional
	Last *int32 `json:"last,omitempty"`
}

//+kubebuilder:object:root=true

// RetentionPolicyList contains a list of RetentionPolicy
type RetentionPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RetentionPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RetentionPolicy{}, &RetentionPolicyList{})
}
