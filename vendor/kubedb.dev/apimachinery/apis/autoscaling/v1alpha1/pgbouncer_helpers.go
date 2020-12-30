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

func (_ PgBouncerAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgBouncerAutoscaler))
}

var _ apis.ResourceInfo = &PgBouncerAutoscaler{}

func (p PgBouncerAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPgBouncerAutoscaler, catalog.GroupName)
}

func (p PgBouncerAutoscaler) ResourceShortCode() string {
	return ResourceCodePgBouncerAutoscaler
}

func (p PgBouncerAutoscaler) ResourceKind() string {
	return ResourceKindPgBouncerAutoscaler
}

func (p PgBouncerAutoscaler) ResourceSingular() string {
	return ResourceSingularPgBouncerAutoscaler
}

func (p PgBouncerAutoscaler) ResourcePlural() string {
	return ResourcePluralPgBouncerAutoscaler
}

func (p PgBouncerAutoscaler) ValidateSpecs() error {
	return nil
}
