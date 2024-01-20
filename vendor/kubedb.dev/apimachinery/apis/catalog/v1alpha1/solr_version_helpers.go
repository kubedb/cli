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

func (s SolrVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralSolrVersion))
}

var _ apis.ResourceInfo = &SolrVersion{}

func (s SolrVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralSolrVersion, catalog.GroupName)
}

func (s SolrVersion) ResourceShortCode() string {
	return ResourceCodeSolrVersion
}

func (s SolrVersion) ResourceKind() string {
	return ResourceKindSolrVersion
}

func (s SolrVersion) ResourceSingular() string {
	return ResourceSingularSolrVersion
}

func (s SolrVersion) ResourcePlural() string {
	return ResourcePluralSolrVersion
}

func (s SolrVersion) ValidateSpecs() error {
	if s.Spec.Version == "" ||
		s.Spec.DB.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for solrVersion "%v":
spec.version,
spec.db.image`, s.Name)
	}
	return nil
}
