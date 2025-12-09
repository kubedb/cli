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

func (Neo4jVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralNeo4jVersion))
}

var _ apis.ResourceInfo = &Neo4jVersion{}

func (r Neo4jVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralNeo4jVersion, catalog.GroupName)
}

func (r Neo4jVersion) ResourceShortCode() string {
	return ResourceCodeNeo4jVersion
}

func (r Neo4jVersion) ResourceKind() string {
	return ResourceKindNeo4jVersion
}

func (r Neo4jVersion) ResourceSingular() string {
	return ResourceSingularNeo4jVersion
}

func (r Neo4jVersion) ResourcePlural() string {
	return ResourcePluralNeo4jVersion
}

func (r Neo4jVersion) ValidateSpecs() error {
	if r.Spec.Version == "" ||
		r.Spec.DB.Image == "" {
		fields := []string{
			"spec.version",
			"spec.db.image",
		}
		return fmt.Errorf("atleast one of the following specs is not set for Neo4jVersion %q: %s", r.Name, strings.Join(fields, ", "))
	}
	return nil
}
