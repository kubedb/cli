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

func (HanaDBVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralHanaDBVersion))
}

var _ apis.ResourceInfo = &HanaDBVersion{}

func (h HanaDBVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralHanaDBVersion, catalog.GroupName)
}

func (h HanaDBVersion) ResourceShortCode() string {
	return ResourceCodeHanaDBVersion
}

func (h HanaDBVersion) ResourceKind() string {
	return ResourceKindHanaDBVersion
}

func (h HanaDBVersion) ResourceSingular() string {
	return ResourceSingularHanaDBVersion
}

func (h HanaDBVersion) ResourcePlural() string {
	return ResourcePluralHanaDBVersion
}

func (h HanaDBVersion) ValidateSpecs() error {
	if h.Spec.Version == "" ||
		h.Spec.DB.Image == "" {
		fields := []string{
			"spec.version",
			"spec.db.image",
		}
		return fmt.Errorf("atleast one of the following specs is not set for HanaDBVersion %q: %s", h.Name, strings.Join(fields, ", "))
	}
	return nil
}
