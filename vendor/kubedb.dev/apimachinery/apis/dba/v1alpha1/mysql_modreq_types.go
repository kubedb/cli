/*
Copyright The KubeDB Authors.

Licensed under the Apache License, ModificationRequest 2.0 (the "License");
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
	ResourceCodeMySQLModificationRequest     = "mymodreq"
	ResourceKindMySQLModificationRequest     = "MySQLModificationRequest"
	ResourceSingularMySQLModificationRequest = "mysqlmodificationrequest"
	ResourcePluralMySQLModificationRequest   = "mysqlmodificationrequests"
)

// MySQLModificationRequest defines a MySQL database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mysqlmodificationrequests,singular=mysqlmodificationrequest,shortName=mymodreq,categories={datastore,kubedb,appscode}
type MySQLModificationRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MySQLModificationRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            MySQLModificationRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MySQLModificationRequestSpec is the spec for elasticsearch version
type MySQLModificationRequestSpec struct {
}

// MySQLModificationRequestStatus is the status for elasticsearch version
type MySQLModificationRequestStatus struct {
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []MySQLModificationRequestCondition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`
}

type MySQLModificationRequestCondition struct {
	// request approval state, currently Approved or Denied.
	Type RequestConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=RequestConditionType"`

	// brief reason for the request state
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`

	// human readable message with details about the request state
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`

	// timestamp for the last update to this condition
	// +optional
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty" protobuf:"bytes,4,opt,name=lastUpdateTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLModificationRequestList is a list of MySQLModificationRequests
type MySQLModificationRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of MySQLModificationRequest CRD objects
	Items []MySQLModificationRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
