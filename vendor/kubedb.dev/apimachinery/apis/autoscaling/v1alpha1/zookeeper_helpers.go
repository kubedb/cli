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

func (r ZooKeeperAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralZooKeeperAutoscaler))
}

var _ apis.ResourceInfo = &ZooKeeperAutoscaler{}

func (r ZooKeeperAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralZooKeeperAutoscaler, autoscaling.GroupName)
}

func (r ZooKeeperAutoscaler) ResourceShortCode() string {
	return ResourceCodeZooKeeperAutoscaler
}

func (r ZooKeeperAutoscaler) ResourceKind() string {
	return ResourceKindZooKeeperAutoscaler
}

func (r ZooKeeperAutoscaler) ResourceSingular() string {
	return ResourceSingularZooKeeperAutoscaler
}

func (r ZooKeeperAutoscaler) ResourcePlural() string {
	return ResourcePluralZooKeeperAutoscaler
}

func (r ZooKeeperAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &ZooKeeperAutoscaler{}

func (r *ZooKeeperAutoscaler) GetStatus() AutoscalerStatus {
	return r.Status
}

func (r *ZooKeeperAutoscaler) SetStatus(s AutoscalerStatus) {
	r.Status = s
}
