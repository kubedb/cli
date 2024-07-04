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

package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// AllowedConsumers defines which consumers may refer to a database instance.
type AllowedConsumers struct {
	// Namespaces indicates namespaces from which Consumers may be attached to
	//
	// +optional
	// +kubebuilder:default={from: Same}
	Namespaces *ConsumerNamespaces `json:"namespaces,omitempty"`

	// Selector specifies a selector for consumers that are allowed to bind
	// to this database instance.
	//
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

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

// ConsumerNamespaces indicate which namespaces Consumers should be selected from.
type ConsumerNamespaces struct {
	// From indicates where Consumers will be selected for the database instance. Possible
	// values are:
	// * All: Consumers in all namespaces.
	// * Selector: Consumers in namespaces selected by the selector
	// * Same: Only Consumers in the same namespace
	//
	// +optional
	// +kubebuilder:default=Same
	From *FromNamespaces `json:"from,omitempty"`

	// Selector must be specified when From is set to "Selector". In that case,
	// only Consumers in Namespaces matching this Selector will be selected by the
	// database instance. This field is ignored for other values of "From".
	//
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}
