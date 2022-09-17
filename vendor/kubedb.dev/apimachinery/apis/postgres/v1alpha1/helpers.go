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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:validation:Enum=Pending;InProgress;Current;Failed
type ReplicationConfigurationPhase string

const (
	ReplicationConfigurationPhasePending    ReplicationConfigurationPhase = "Pending"
	ReplicationConfigurationPhaseInProgress ReplicationConfigurationPhase = "InProgress"
	ReplicationConfigurationPhaseCurrent    ReplicationConfigurationPhase = "Current"
	ReplicationConfigurationPhaseFailed     ReplicationConfigurationPhase = "Failed"
)

// +kubebuilder:validation:Enum=Delete;Retain
type DeletionPolicy string

const (
	// Deletes database pods, service, pvcs but leave the snapshot data intact. This will not create a DormantDatabase.
	DeletionPolicyDelete DeletionPolicy = "Delete"
	// Pauses database into a DormantDatabase
	DeletionPolicyRetain DeletionPolicy = "Retain"
)

// FromNamespaces specifies namespace from which Consumers may be attached to a
// database instance.
//
// +kubebuilder:validation:Enum=All;Selector;Same
type FromNamespaces string

const (
	// Consumers in all namespaces may be attached to the database instance.
	NamespacesFromAll FromNamespaces = "All"
	// Only Consumers in namespaces selected by the selector may be attached to the database instance.
	NamespacesFromSelector FromNamespaces = "Selector"
	// Only Consumers in the same namespace as the database instance may be attached to it.
	NamespacesFromSame FromNamespaces = "Same"
)

// AllowedSubscribers defines which consumers may refer to a database instance.
type AllowedSubscribers struct {
	// Namespaces indicates namespaces from which Subscribers may be attached to
	//
	// +optional
	// +kubebuilder:default={from: Same}
	Namespaces *SubscriberNamespaces `json:"namespaces,omitempty"`

	// Selector specifies a selector for consumers that are allowed to bind
	// to this database instance.
	//
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// SubscriberNamespaces indicate which namespaces Subscribers should be selected from.
type SubscriberNamespaces struct {
	// From indicates where Subscribers will be selected for the database instance. Possible
	// values are:
	// * All: Subscribers in all namespaces.
	// * Selector: Subscribers in namespaces selected by the selector
	// * Same: Only Subscribers in the same namespace
	//
	// +optional
	// +kubebuilder:default=Same
	From *FromNamespaces `json:"from,omitempty"`

	// Selector must be specified when From is set to "Selector". In that case,
	// only Subscribers in Namespaces matching this Selector will be selected by the
	// database instance. This field is ignored for other values of "From".
	//
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}
