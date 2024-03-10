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

func (p *PgpoolVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgpoolVersion))
}

var _ apis.ResourceInfo = &PgpoolVersion{}

func (p *PgpoolVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPgpoolVersion, catalog.GroupName)
}

func (p *PgpoolVersion) ResourceShortCode() string {
	return ResourceCodePgpoolVersion
}

func (p *PgpoolVersion) ResourceKind() string {
	return ResourceKindPgpoolVersion
}

func (p *PgpoolVersion) ResourceSingular() string {
	return ResourceSingularPgpoolVersion
}

func (p *PgpoolVersion) ResourcePlural() string {
	return ResourcePluralPgpoolVersion
}

func (p *PgpoolVersion) ValidateSpecs() error {
	if p.Spec.Version == "" ||
		p.Spec.Pgpool.Image == "" ||
		p.Spec.Exporter.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for pgpoolVersion "%v":
spec.version,
spec.pgpool.image,
spec.exporter.image`, p.Name)
	}
	return nil
}
