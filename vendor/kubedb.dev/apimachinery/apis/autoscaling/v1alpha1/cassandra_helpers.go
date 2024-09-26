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
	"kubedb.dev/apimachinery/apis/autoscaling"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (r CassandraAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralCassandraAutoscaler))
}

var _ apis.ResourceInfo = &CassandraAutoscaler{}

func (r CassandraAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralCassandraAutoscaler, autoscaling.GroupName)
}

func (r CassandraAutoscaler) ResourceShortCode() string {
	return ResourceCodeCassandraAutoscaler
}

func (r CassandraAutoscaler) ResourceKind() string {
	return ResourceKindCassandraAutoscaler
}

func (r CassandraAutoscaler) ResourceSingular() string {
	return ResourceSingularCassandraAutoscaler
}

func (r CassandraAutoscaler) ResourcePlural() string {
	return ResourcePluralCassandraAutoscaler
}

func (r CassandraAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &CassandraAutoscaler{}

func (r *CassandraAutoscaler) GetStatus() AutoscalerStatus {
	return r.Status
}

func (r *CassandraAutoscaler) SetStatus(s AutoscalerStatus) {
	r.Status = s
}
