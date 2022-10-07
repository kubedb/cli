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

func (_ PerconaXtraDBOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPerconaXtraDBOpsRequest))
}

var _ apis.ResourceInfo = &PerconaXtraDBOpsRequest{}

func (p PerconaXtraDBOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPerconaXtraDBOpsRequest, ops.GroupName)
}

func (p PerconaXtraDBOpsRequest) ResourceShortCode() string {
	return ResourceCodePerconaXtraDBOpsRequest
}

func (p PerconaXtraDBOpsRequest) ResourceKind() string {
	return ResourceKindPerconaXtraDBOpsRequest
}

func (p PerconaXtraDBOpsRequest) ResourceSingular() string {
	return ResourceSingularPerconaXtraDBOpsRequest
}

func (p PerconaXtraDBOpsRequest) ResourcePlural() string {
	return ResourcePluralPerconaXtraDBOpsRequest
}

func (p PerconaXtraDBOpsRequest) ValidateSpecs() error {
	return nil
}

var _ Accessor = &PerconaXtraDBOpsRequest{}

func (e *PerconaXtraDBOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return e.ObjectMeta
}

func (e *PerconaXtraDBOpsRequest) GetRequestType() string {
	return string(e.Spec.Type)
}

func (e *PerconaXtraDBOpsRequest) GetDBRefName() string {
	return e.Spec.DatabaseRef.Name
}

func (e *PerconaXtraDBOpsRequest) GetStatus() OpsRequestStatus {
	return e.Status
}

func (e *PerconaXtraDBOpsRequest) SetStatus(s OpsRequestStatus) {
	e.Status = s
}
