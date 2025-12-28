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
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindDatabaseInfo = "DatabaseInfo"
	ResourceDatabaseInfo     = "databaseinfo"
	ResourceDatabaseInfos    = "databaseinfos"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=get,list,update,delete,watch
// +genclient:onlyVerbs=create
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=databaseinfos,singular=databaseinfo,scope=Cluster
type DatabaseInfo struct {
	metav1.TypeMeta `json:",inline"`
	// Request describes the attributes for the graph request.
	// +optional
	Request *DatabaseInfoRequest `json:"request,omitempty"`
	// Response describes the attributes for the graph response.
	// +optional
	Response *DatabaseInfoResponse `json:"response,omitempty"`
}

type DatabaseInfoRequest struct {
	Source kmapi.ObjectInfo `json:"source"`
	Keys   []string         `json:"keys,omitempty"`
}

type DatabaseInfoResponse struct {
	Configurations   []SingleComponentConfiguration `json:"configurations,omitempty"`
	AvailableSecrets []string                       `json:"availableSecrets,omitempty"`
}

type SingleComponentConfiguration struct {
	// +optional
	ComponentName string            `json:"componentName,omitempty"`
	Data          map[string][]byte `json:"data,omitempty"`
	// +optional
	SecretName string `json:"secretName,omitempty"`
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
}
