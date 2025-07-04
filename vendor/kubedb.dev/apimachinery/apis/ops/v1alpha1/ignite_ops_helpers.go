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

func (r *IgniteOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralIgniteOpsRequest))
}

var _ apis.ResourceInfo = &IgniteOpsRequest{}

func (r *IgniteOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralIgniteOpsRequest, ops.GroupName)
}

func (r *IgniteOpsRequest) ResourceShortCode() string {
	return ResourceCodeIgniteOpsRequest
}

func (r *IgniteOpsRequest) ResourceKind() string {
	return ResourceKindIgniteOpsRequest
}

func (r *IgniteOpsRequest) ResourceSingular() string {
	return ResourceSingularIgniteOpsRequest
}

func (r *IgniteOpsRequest) ResourcePlural() string {
	return ResourcePluralIgniteOpsRequest
}

var _ Accessor = &IgniteOpsRequest{}

func (r *IgniteOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r *IgniteOpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *IgniteOpsRequest) GetRequestType() any {
	return r.Spec.Type
}

func (r *IgniteOpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *IgniteOpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
