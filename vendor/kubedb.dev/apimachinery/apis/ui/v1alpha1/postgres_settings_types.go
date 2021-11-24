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
	ResourceKindPostgresSettings = "PostgresSettings"
	ResourcePostgresSettings     = "postgressettings"
	ResourcePostgresSettingss    = "postgressettings"
)

// PostgresSettingsSpec defines the desired state of PostgresSettings
type PostgresSettingsSpec struct {
	Settings []PGSetting `json:"settings" protobuf:"bytes,1,rep,name=settings"`
}

type PGSetting struct {
	Name         string `json:"name" protobuf:"bytes,1,opt,name=name"`
	CurrentValue string `json:"currentValue" protobuf:"bytes,2,opt,name=currentValue"`
	DefaultValue string `json:"defaultValue" protobuf:"bytes,3,opt,name=defaultValue"`
	Unit         string `json:"unit" protobuf:"bytes,4,opt,name=unit"`
	Source       string `json:"source" protobuf:"bytes,5,opt,name=source"`
}

// PostgresSettings is the Schema for the PostgresSettingss API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresSettings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec PostgresSettingsSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// PostgresSettingsList contains a list of PostgresSettings

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresSettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []PostgresSettings `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&PostgresSettings{}, &PostgresSettingsList{})
}
