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
)

const (
	ResourceKindMariaDBSchemaOverview = "MariaDBSchemaOverview"
	ResourceMariaDBSchemaOverview     = "mariadbschemaoverview"
	ResourceMariaDBSchemaOverviews    = "mariadbschemaoverviews"
)

type MariaDBSchemaOverviewSpec = GenericSchemaOverviewSpec

// MariaDBSchemaOverview is the Schema for the MariaDBSchemaOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MariaDBSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MariaDBSchemaOverviewSpec `json:"spec,omitempty"`
}

// MariaDBSchemaOverviewList contains a list of MariaDBSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MariaDBSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MariaDBSchemaOverview `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MariaDBSchemaOverview{}, &MariaDBSchemaOverviewList{})
}
