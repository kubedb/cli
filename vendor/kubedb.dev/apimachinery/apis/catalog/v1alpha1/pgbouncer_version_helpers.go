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

func (_ PgBouncerVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgBouncerVersion))
}

var _ apis.ResourceInfo = &PgBouncerVersion{}

func (p PgBouncerVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPgBouncerVersion, catalog.GroupName)
}

func (p PgBouncerVersion) ResourceShortCode() string {
	return ResourceCodePgBouncerVersion
}

func (p PgBouncerVersion) ResourceKind() string {
	return ResourceKindPgBouncerVersion
}

func (p PgBouncerVersion) ResourceSingular() string {
	return ResourceSingularPgBouncerVersion
}

func (p PgBouncerVersion) ResourcePlural() string {
	return ResourcePluralPgBouncerVersion
}

func (p PgBouncerVersion) ValidateSpecs() error {
	if p.Spec.Version == "" ||
		p.Spec.Exporter.Image == "" ||
		p.Spec.PgBouncer.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for pgbouncerversion "%v":
spec.version,
spec.pgBouncer.image,
spec.exporter.image.`, p.Name)
	}
	return nil
}

func (p PgBouncerVersion) IsDeprecated() bool {
	return p.Spec.Deprecated
}
