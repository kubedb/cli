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

func (p ProxySQLOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralProxySQLOpsRequest))
}

var _ apis.ResourceInfo = &ProxySQLOpsRequest{}

func (p ProxySQLOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralProxySQLOpsRequest, ops.GroupName)
}

func (p ProxySQLOpsRequest) ResourceShortCode() string {
	return ResourceCodeProxySQLOpsRequest
}

func (p ProxySQLOpsRequest) ResourceKind() string {
	return ResourceKindProxySQLOpsRequest
}

func (p ProxySQLOpsRequest) ResourceSingular() string {
	return ResourceSingularProxySQLOpsRequest
}

func (p ProxySQLOpsRequest) ResourcePlural() string {
	return ResourcePluralProxySQLOpsRequest
}

var _ Accessor = &ProxySQLOpsRequest{}

func (p *ProxySQLOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return p.ObjectMeta
}

func (p ProxySQLOpsRequest) GetRequestType() any {
	return p.Spec.Type
}

func (p ProxySQLOpsRequest) GetUpdateVersionSpec() *ProxySQLUpdateVersionSpec {
	return p.Spec.UpdateVersion
}

func (p *ProxySQLOpsRequest) GetDBRefName() string {
	return p.Spec.ProxyRef.Name
}

func (p *ProxySQLOpsRequest) GetStatus() OpsRequestStatus {
	return p.Status
}

func (p *ProxySQLOpsRequest) SetStatus(s OpsRequestStatus) {
	p.Status = s
}
