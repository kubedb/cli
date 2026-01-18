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
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/ops"
	"kubedb.dev/apimachinery/crds"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/apiextensions"
)

func (r *MSSQLServerOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMSSQLServerOpsRequest))
}

var _ apis.ResourceInfo = &MSSQLServerOpsRequest{}

func (r *MSSQLServerOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMSSQLServerOpsRequest, ops.GroupName)
}

func (r *MSSQLServerOpsRequest) ResourceShortCode() string {
	return ResourceCodeMSSQLServerOpsRequest
}

func (r *MSSQLServerOpsRequest) ResourceKind() string {
	return ResourceKindMSSQLServerOpsRequest
}

func (r *MSSQLServerOpsRequest) ResourceSingular() string {
	return ResourceSingularMSSQLServerOpsRequest
}

func (r *MSSQLServerOpsRequest) ResourcePlural() string {
	return ResourcePluralMSSQLServerOpsRequest
}

var _ Accessor = &MSSQLServerOpsRequest{}

func (r *MSSQLServerOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r *MSSQLServerOpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *MSSQLServerOpsRequest) GetRequestType() string {
	return string(r.Spec.Type)
}

func (r *MSSQLServerOpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *MSSQLServerOpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
