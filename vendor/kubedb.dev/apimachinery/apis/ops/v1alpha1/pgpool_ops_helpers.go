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

func (r *PgpoolOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgpoolOpsRequest))
}

var _ apis.ResourceInfo = &PgpoolOpsRequest{}

func (r *PgpoolOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPgpoolOpsRequest, ops.GroupName)
}

func (r *PgpoolOpsRequest) ResourceShortCode() string {
	return ResourceCodePgpoolOpsRequest
}

func (r *PgpoolOpsRequest) ResourceKind() string {
	return ResourceKindPgpoolOpsRequest
}

func (r *PgpoolOpsRequest) ResourceSingular() string {
	return ResourceSingularPgpoolOpsRequest
}

func (r *PgpoolOpsRequest) ResourcePlural() string {
	return ResourcePluralPgpoolOpsRequest
}

var _ Accessor = &PgpoolOpsRequest{}

func (r *PgpoolOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r *PgpoolOpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *PgpoolOpsRequest) GetRequestType() any {
	return r.Spec.Type
}

func (r *PgpoolOpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *PgpoolOpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
