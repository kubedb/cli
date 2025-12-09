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

func (QdrantVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralQdrantVersion))
}

var _ apis.ResourceInfo = &QdrantVersion{}

func (q QdrantVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralQdrantVersion, catalog.GroupName)
}

func (q QdrantVersion) ResourceShortCode() string {
	return ResourceCodeQdrantVersion
}

func (q QdrantVersion) ResourceKind() string {
	return ResourceKindQdrantVersion
}

func (q QdrantVersion) ResourceSingular() string {
	return ResourceSingularQdrantVersion
}

func (q QdrantVersion) ResourcePlural() string {
	return ResourcePluralQdrantVersion
}

func (q QdrantVersion) ValidateSpecs() error {
	if q.Spec.Version == "" ||
		q.Spec.DB.Image == "" {
		fields := []string{
			"spec.version",
			"spec.db.image",
		}
		return fmt.Errorf("atleast one of the following specs is not set for QdrantVersion %q: %s", q.Name, strings.Join(fields, ", "))
	}
	return nil
}
