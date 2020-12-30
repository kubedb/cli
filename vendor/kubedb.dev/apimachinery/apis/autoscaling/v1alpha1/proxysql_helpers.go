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

func (_ ProxySQLAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralProxySQLAutoscaler))
}

var _ apis.ResourceInfo = &ProxySQLAutoscaler{}

func (p ProxySQLAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralProxySQLAutoscaler, catalog.GroupName)
}

func (p ProxySQLAutoscaler) ResourceShortCode() string {
	return ""
}

func (p ProxySQLAutoscaler) ResourceKind() string {
	return ResourceKindProxySQLAutoscaler
}

func (p ProxySQLAutoscaler) ResourceSingular() string {
	return ResourceSingularProxySQLAutoscaler
}

func (p ProxySQLAutoscaler) ResourcePlural() string {
	return ResourcePluralProxySQLAutoscaler
}
