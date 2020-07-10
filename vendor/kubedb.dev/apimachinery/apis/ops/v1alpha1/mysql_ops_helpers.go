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
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
)

func (_ MySQLOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMySQLOpsRequest))
}

var _ apis.ResourceInfo = &MySQLOpsRequest{}

func (m MySQLOpsRequest) ResourceShortCode() string {
	return ResourceCodeMySQLOpsRequest
}

func (m MySQLOpsRequest) ResourceKind() string {
	return ResourceKindMySQLOpsRequest
}

func (m MySQLOpsRequest) ResourceSingular() string {
	return ResourceSingularMySQLOpsRequest
}

func (m MySQLOpsRequest) ResourcePlural() string {
	return ResourcePluralMySQLOpsRequest
}

func (m MySQLOpsRequest) ValidateSpecs() error {
	return nil
}

func (m MySQLOpsRequest) GetKey() string {
	return m.Namespace + "/" + m.Name
}

func (m MySQLOpsRequest) OffshootName() string {
	return m.Name
}

func (m MySQLOpsRequest) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelOpsRequestKind: ResourceSingularMySQLOpsRequest,
		LabelOpsRequestName: m.Name,
	}
}

func (m MySQLOpsRequest) OffshootLabels() map[string]string {
	out := m.OffshootSelectors()
	return meta_util.FilterKeys(GenericKey, out, m.Labels)
}
