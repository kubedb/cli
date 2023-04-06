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

//go:generate go-enum --mustparse --names --values
package v1alpha1

import (
	"kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ResourceCodeProxySQLOpsRequest     = "prxops"
	ResourceKindProxySQLOpsRequest     = "ProxySQLOpsRequest"
	ResourceSingularProxySQLOpsRequest = "proxysqlopsrequest"
	ResourcePluralProxySQLOpsRequest   = "proxysqlopsrequests"
)

// ProxySQLOpsRequest defines a ProxySQL load-balancer DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=proxysqlopsrequests,singular=proxysqlopsrequest,shortName=prxops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ProxySQLOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProxySQLOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus       `json:"status,omitempty"`
}

// ProxySQLOpsRequestSpec is the spec for ProxySQLOpsRequest
type ProxySQLOpsRequestSpec struct {
	// Specifies the ProxySQL reference
	ProxyRef core.LocalObjectReference `json:"proxyRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type ProxySQLOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading ProxySQL
	UpdateVersion *ProxySQLUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *ProxySQLHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *ProxySQLVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for custom configuration of ProxySQL
	Configuration *ProxySQLCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;Restart;Reconfigure;ReconfigureTLS
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, Restart, Reconfigure, ReconfigureTLS)
type ProxySQLOpsRequestType string

// ProxySQLReplicaReadinessCriteria is the criteria for checking readiness of a ProxySQL pod
// after updating, horizontal scaling etc.
type ProxySQLReplicaReadinessCriteria struct{}

type ProxySQLUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                            `json:"targetVersion,omitempty"`
	ReadinessCriteria *ProxySQLReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

// HorizontalScaling is the spec for ProxySQL horizontal scaling
type ProxySQLHorizontalScalingSpec struct {
	// Number of nodes/members of the group
	Member *int32 `json:"member,omitempty"`
}

// ProxySQLVerticalScalingSpec is the spec for ProxySQL vertical scaling
type ProxySQLVerticalScalingSpec struct {
	ProxySQL *core.ResourceRequirements `json:"proxysql,omitempty"`
}

type ProxySQLCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty"`
	Data      map[string]string          `json:"data,omitempty"`
	Remove    bool                       `json:"remove,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProxySQLOpsRequestList is a list of ProxySQLOpsRequests
type ProxySQLOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	//+optional
	// Items is a list of ProxySQLOpsRequest CRD objects
	Items []ProxySQLOpsRequest `json:"items,omitempty"`
}

type ProxySQLCustomConfigurationSpec struct {
	//+optional
	MySQLUsers *MySQLUsers `json:"mysqlUsers,omitempty"`

	//+optional
	MySQLQueryRules *MySQLQueryRules `json:"mysqlQueryRules,omitempty"`

	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	AdminVariables *runtime.RawExtension `json:"adminVariables,omitempty"`

	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	MySQLVariables *runtime.RawExtension `json:"mysqlVariables,omitempty"`
}

type MySQLUsers struct {
	Users       []v1alpha2.MySQLUser `json:"users"`
	RequestType OperationType        `json:"reqType"`
}

type MySQLQueryRules struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	Rules       []*runtime.RawExtension `json:"rules"`
	RequestType OperationType           `json:"reqType"`
}

type OperationType string

const (
	ProxySQLConfigurationAdd    OperationType = "add"
	ProxySQLConfigurationDelete OperationType = "delete"
	ProxySQLConfigurationUpdate OperationType = "update"
)
