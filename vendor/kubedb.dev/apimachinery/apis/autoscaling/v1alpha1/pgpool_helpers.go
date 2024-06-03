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

func (r PgpoolAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgpoolAutoscaler))
}

var _ apis.ResourceInfo = &PgpoolAutoscaler{}

func (r PgpoolAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPgpoolAutoscaler, autoscaling.GroupName)
}

func (r PgpoolAutoscaler) ResourceShortCode() string {
	return ResourceCodePgpoolAutoscaler
}

func (r PgpoolAutoscaler) ResourceKind() string {
	return ResourceKindPgpoolAutoscaler
}

func (r PgpoolAutoscaler) ResourceSingular() string {
	return ResourceSingularPgpoolAutoscaler
}

func (r PgpoolAutoscaler) ResourcePlural() string {
	return ResourcePluralPgpoolAutoscaler
}

func (r PgpoolAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &PgpoolAutoscaler{}

func (r *PgpoolAutoscaler) GetStatus() AutoscalerStatus {
	return r.Status
}

func (r *PgpoolAutoscaler) SetStatus(s AutoscalerStatus) {
	r.Status = s
}
