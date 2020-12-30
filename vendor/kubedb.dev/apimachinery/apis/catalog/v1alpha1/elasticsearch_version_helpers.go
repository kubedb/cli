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

func (_ ElasticsearchVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralElasticsearchVersion))
}

var _ apis.ResourceInfo = &ElasticsearchVersion{}

func (e ElasticsearchVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralElasticsearchVersion, catalog.GroupName)
}

func (e ElasticsearchVersion) ResourceShortCode() string {
	return ResourceCodeElasticsearchVersion
}

func (e ElasticsearchVersion) ResourceKind() string {
	return ResourceKindElasticsearchVersion
}

func (e ElasticsearchVersion) ResourceSingular() string {
	return ResourceSingularElasticsearchVersion
}

func (e ElasticsearchVersion) ResourcePlural() string {
	return ResourcePluralElasticsearchVersion
}

func (e ElasticsearchVersion) ValidateSpecs() error {
	if e.Spec.AuthPlugin == "" ||
		e.Spec.Version == "" ||
		e.Spec.DB.Image == "" ||
		e.Spec.Exporter.Image == "" ||
		e.Spec.InitContainer.YQImage == "" ||
		e.Spec.InitContainer.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for elasticsearchVersion "%v":
spec.authPlugin,
spec.version,
spec.db.image,
spec.exporter.image,
spec.initContainer.yqImage,
spec.initContainer.image.`, e.Name)
	}
	return nil
}
