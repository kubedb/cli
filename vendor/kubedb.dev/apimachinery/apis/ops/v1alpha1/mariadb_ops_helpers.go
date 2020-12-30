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

	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
)

func (_ MariaDBOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMariaDBOpsRequest))
}

var _ apis.ResourceInfo = &MariaDBOpsRequest{}

func (m MariaDBOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMariaDBOpsRequest, ops.GroupName)
}

func (m MariaDBOpsRequest) ResourceShortCode() string {
	return ResourceCodeMariaDBOpsRequest
}

func (m MariaDBOpsRequest) ResourceKind() string {
	return ResourceKindMariaDBOpsRequest
}

func (m MariaDBOpsRequest) ResourceSingular() string {
	return ResourceSingularMariaDBOpsRequest
}

func (m MariaDBOpsRequest) ResourcePlural() string {
	return ResourcePluralMariaDBOpsRequest
}

func (m MariaDBOpsRequest) ValidateSpecs() error {
	return nil
}

func (m MariaDBOpsRequest) GetKey() string {
	return m.Namespace + "/" + m.Name
}

func (m MariaDBOpsRequest) OffshootName() string {
	return m.Name
}

func (m MariaDBOpsRequest) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelOpsRequestKind: ResourceSingularMariaDBOpsRequest,
		LabelOpsRequestName: m.Name,
	}
}

func (m MariaDBOpsRequest) OffshootLabels() map[string]string {
	out := m.OffshootSelectors()
	return meta_util.FilterKeys(GenericKey, out, m.Labels)
}
