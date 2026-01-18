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

func (r *ClickHouseOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralClickHouseOpsRequest))
}

var _ apis.ResourceInfo = &ClickHouseOpsRequest{}

func (r *ClickHouseOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralClickHouseOpsRequest, ops.GroupName)
}

func (r *ClickHouseOpsRequest) ResourceShortCode() string {
	return ResourceCodeClickHouseOpsRequest
}

func (r *ClickHouseOpsRequest) ResourceKind() string {
	return ResourceKindClickHouseOpsRequest
}

func (r *ClickHouseOpsRequest) ResourceSingular() string {
	return ResourceSingularClickHouseOpsRequest
}

func (r *ClickHouseOpsRequest) ResourcePlural() string {
	return ResourcePluralClickHouseOpsRequest
}

var _ Accessor = &ClickHouseOpsRequest{}

func (r *ClickHouseOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r *ClickHouseOpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *ClickHouseOpsRequest) GetRequestType() string {
	return string(r.Spec.Type)
}

func (r *ClickHouseOpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *ClickHouseOpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
