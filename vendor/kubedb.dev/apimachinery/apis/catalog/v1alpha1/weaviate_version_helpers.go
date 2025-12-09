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

func (WeaviateVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralWeaviateVersion))
}

var _ apis.ResourceInfo = &WeaviateVersion{}

func (w WeaviateVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralWeaviateVersion, catalog.GroupName)
}

func (w WeaviateVersion) ResourceShortCode() string {
	return ResourceCodeWeaviateVersion
}

func (w WeaviateVersion) ResourceKind() string {
	return ResourceKindWeaviateVersion
}

func (w WeaviateVersion) ResourceSingular() string {
	return ResourceSingularWeaviateVersion
}

func (w WeaviateVersion) ResourcePlural() string {
	return ResourcePluralWeaviateVersion
}

func (w WeaviateVersion) ValidateSpecs() error {
	if w.Spec.Version == "" ||
		w.Spec.DB.Image == "" {
		fields := []string{
			"spec.version",
			"spec.db.image",
		}
		return fmt.Errorf("atleast one of the following specs is not set for QdrantVersion %q: %s", w.Name, strings.Join(fields, ", "))
	}
	return nil
}
