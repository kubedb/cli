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

func (h HazelcastVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralHazelcastVersion))
}

var _ apis.ResourceInfo = &HazelcastVersion{}

func (h HazelcastVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralHazelcastVersion, catalog.GroupName)
}

func (h HazelcastVersion) ResourceShortCode() string {
	return ResourceCodeHazelcastVersion
}

func (h HazelcastVersion) ResourceKind() string {
	return ResourceKindHazelcastVersion
}

func (h HazelcastVersion) ResourceSingular() string {
	return ResourceSingularHazelcastVersion
}

func (h HazelcastVersion) ResourcePlural() string {
	return ResourcePluralHazelcastVersion
}

func (h HazelcastVersion) ValidateSpecs() error {
	if h.Spec.Version == "" ||
		h.Spec.DB.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for HazelcastVersion "%v":
spec.version,
spec.db.image`, h.Name)
	}
	return nil
}
