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

func (_ *KafkaAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralKafkaAutoscaler))
}

var _ apis.ResourceInfo = &KafkaAutoscaler{}

func (k *KafkaAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralKafkaAutoscaler, autoscaling.GroupName)
}

func (k *KafkaAutoscaler) ResourceShortCode() string {
	return ResourceCodeKafkaAutoscaler
}

func (k *KafkaAutoscaler) ResourceKind() string {
	return ResourceKindKafkaAutoscaler
}

func (k *KafkaAutoscaler) ResourceSingular() string {
	return ResourceSingularKafkaAutoscaler
}

func (k *KafkaAutoscaler) ResourcePlural() string {
	return ResourcePluralKafkaAutoscaler
}

func (k *KafkaAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &KafkaAutoscaler{}

func (k *KafkaAutoscaler) GetStatus() AutoscalerStatus {
	return k.Status
}

func (k *KafkaAutoscaler) SetStatus(s AutoscalerStatus) {
	k.Status = s
}
