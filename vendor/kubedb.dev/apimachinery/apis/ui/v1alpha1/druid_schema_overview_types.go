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

const (
	ResourceKindDruidSchemaOverview = "DruidSchemaOverview"
	ResourceDruidSchemaOverview     = "druidschemaoverview"
	ResourceDruidSchemaOverviews    = "druidschemaoverviews"
)

type DruidSchemaOverviewSpec struct {
	CollectedAt *metav1.Time            `json:"collectedAt,omitempty"`
	Datasources []DruidDatasourceSchema `json:"datasources,omitempty"`
}

type DruidDatasourceSchema struct {
	Name      string        `json:"name"`
	Joinable  bool          `json:"joinable"`
	Broadcast bool          `json:"broadcast"`
	Columns   []DruidColumn `json:"columns,omitempty"`
}

type DruidColumn struct {
	Name            string `json:"name"`
	OrdinalPosition *int64 `json:"ordinalPosition,omitempty"`
	DataType        string `json:"dataType"`
	Nullable        bool   `json:"nullable"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DruidSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DruidSchemaOverviewSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DruidSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DruidSchemaOverview `json:"items"`
}
