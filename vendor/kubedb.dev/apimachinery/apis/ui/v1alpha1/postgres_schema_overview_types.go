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
	ResourceKindPostgresSchemaOverview = "PostgresSchemaOverview"
	ResourcePostgresSchemaOverview     = "postgresschemaoverview"
	ResourcePostgresSchemaOverviews    = "postgresschemaoverviews"
)

type PostgresSchemaOverviewSpec = GenericSchemaOverviewSpec

// PostgresSchemaOverview is the Schema for the PostgresSchemaOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PostgresSchemaOverviewSpec `json:"spec,omitempty"`
}

// PostgresSchemaOverviewList contains a list of PostgresSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresSchemaOverview `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PostgresSchemaOverview{}, &PostgresSchemaOverviewList{})
}
