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
	"kubedb.dev/apimachinery/apis/catalog"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (_ MySQLAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMySQLAutoscaler))
}

var _ apis.ResourceInfo = &MySQLAutoscaler{}

func (m MySQLAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMySQLAutoscaler, catalog.GroupName)
}

func (m MySQLAutoscaler) ResourceShortCode() string {
	return ResourceCodeMySQLAutoscaler
}

func (m MySQLAutoscaler) ResourceKind() string {
	return ResourceKindMySQLAutoscaler
}

func (m MySQLAutoscaler) ResourceSingular() string {
	return ResourceSingularMySQLAutoscaler
}

func (m MySQLAutoscaler) ResourcePlural() string {
	return ResourcePluralMySQLAutoscaler
}

func (m MySQLAutoscaler) ValidateSpecs() error {
	return nil
}
