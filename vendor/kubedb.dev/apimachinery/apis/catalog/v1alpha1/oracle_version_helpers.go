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
	"strings"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/catalog"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (OracleVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralOracleVersion))
}

var _ apis.ResourceInfo = &OracleVersion{}

func (p OracleVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralOracleVersion, catalog.GroupName)
}

func (p OracleVersion) ResourceShortCode() string {
	return ResourceCodeOracleVersion
}

func (p OracleVersion) ResourceKind() string {
	return ResourceKindOracleVersion
}

func (p OracleVersion) ResourceSingular() string {
	return ResourceSingularOracleVersion
}

func (p OracleVersion) ResourcePlural() string {
	return ResourcePluralOracleVersion
}

func (p OracleVersion) ValidateSpecs() error {
	if p.Spec.Version == "" ||
		p.Spec.DB.Image == "" ||
		p.Spec.Exporter.Image == "" {
		fields := []string{
			"spec.version",
			"spec.db.image",
			"spec.exporter.image",
		}
		return fmt.Errorf("atleast one of the following specs is not set for postgresVersion %q: %s", p.Name, strings.Join(fields, ", "))
	}
	return nil
}
