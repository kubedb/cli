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

func (_ MongoDBAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMongoDBAutoscaler))
}

var _ apis.ResourceInfo = &MongoDBAutoscaler{}

func (m MongoDBAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMongoDBAutoscaler, catalog.GroupName)
}

func (m MongoDBAutoscaler) ResourceShortCode() string {
	return ResourceCodeMongoDBAutoscaler
}

func (m MongoDBAutoscaler) ResourceKind() string {
	return ResourceKindMongoDBAutoscaler
}

func (m MongoDBAutoscaler) ResourceSingular() string {
	return ResourceSingularMongoDBAutoscaler
}

func (m MongoDBAutoscaler) ResourcePlural() string {
	return ResourcePluralMongoDBAutoscaler
}

func (m MongoDBAutoscaler) ValidateSpecs() error {
	return nil
}
