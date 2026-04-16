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

func (f DocumentDBVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDocumentDBVersion))
}

var _ apis.ResourceInfo = &DocumentDBVersion{}

func (f DocumentDBVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralDocumentDBVersion, catalog.GroupName)
}

func (f DocumentDBVersion) ResourceShortCode() string {
	return ResourceCodeDocumentDBVersion
}

func (f DocumentDBVersion) ResourceKind() string {
	return ResourceKindDocumentDBVersion
}

func (f DocumentDBVersion) ResourceSingular() string {
	return ResourceSingularDocumentDBVersion
}

func (f DocumentDBVersion) ResourcePlural() string {
	return ResourcePluralDocumentDBVersion
}

func (f DocumentDBVersion) ValidateSpecs() error {
	if f.Spec.Version == "" ||
		f.Spec.DB.Image == "" {
		fields := []string{
			"spec.version",
			"spec.db.image",
		}
		return fmt.Errorf("atleast one of the following specs is not set for documentdbVersion %q: %s", f.Name, strings.Join(fields, ", "))
	}
	return nil
}
