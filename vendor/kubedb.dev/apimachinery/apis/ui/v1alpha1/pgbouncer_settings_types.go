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
	ResourceKindPgBouncerSettings = "PgBouncerSettings"
	ResourcePgBouncerSettings     = "pgbouncersettings"
	ResourcePgBouncerSettingss    = "pgbouncersettings"
)

// PgBouncerSettingsSpec defines the desired state of PgBouncerSettings
type PgBouncerSettingsSpec struct {
	Settings []PBSetting `json:"settings"`
}

type PBSetting struct {
	Name         string `json:"name"`
	CurrentValue string `json:"currentValue"`
	DefaultValue string `json:"defaultValue"`
	Changeable   bool   `json:"changeable"`
}

// PgBouncerSettings is the Schema for the PgBouncerSettingss API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerSettings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PgBouncerSettingsSpec `json:"spec,omitempty"`
}

// PgBouncerSettingsList contains a list of PgBouncerSettings
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerSettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PgBouncerSettings `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PgBouncerSettings{}, &PgBouncerSettingsList{})
}
