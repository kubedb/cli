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

func (p PgBouncerOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgBouncerOpsRequest))
}

var _ apis.ResourceInfo = &PgBouncerOpsRequest{}

func (p PgBouncerOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPgBouncerOpsRequest, ops.GroupName)
}

func (p PgBouncerOpsRequest) ResourceShortCode() string {
	return ResourceCodePgBouncerOpsRequest
}

func (p PgBouncerOpsRequest) ResourceKind() string {
	return ResourceKindPgBouncerOpsRequest
}

func (p PgBouncerOpsRequest) ResourceSingular() string {
	return ResourceSingularPgBouncerOpsRequest
}

func (p PgBouncerOpsRequest) ResourcePlural() string {
	return ResourcePluralPgBouncerOpsRequest
}

func (p PgBouncerOpsRequest) ValidateSpecs() error {
	return nil
}

var _ Accessor = &PgBouncerOpsRequest{}

func (p *PgBouncerOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return p.ObjectMeta
}

func (p PgBouncerOpsRequest) GetRequestType() any {
	return p.Spec.Type
}

func (p PgBouncerOpsRequest) GetUpdateVersionSpec() *PgBouncerUpdateVersionSpec {
	return p.Spec.UpdateVersion
}

func (p *PgBouncerOpsRequest) GetDBRefName() string {
	return p.Spec.ServerRef.Name
}

func (p *PgBouncerOpsRequest) GetStatus() OpsRequestStatus {
	return p.Status
}

func (p *PgBouncerOpsRequest) SetStatus(s OpsRequestStatus) {
	p.Status = s
}
