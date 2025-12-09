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
	"kubedb.dev/apimachinery/apis/ops"
	"kubedb.dev/apimachinery/crds"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/apiextensions"
)

func (*KafkaOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralKafkaOpsRequest))
}

var _ apis.ResourceInfo = &KafkaOpsRequest{}

func (k *KafkaOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralKafkaOpsRequest, ops.GroupName)
}

func (k *KafkaOpsRequest) ResourceShortCode() string {
	return ResourceCodeKafkaOpsRequest
}

func (k *KafkaOpsRequest) ResourceKind() string {
	return ResourceKindKafkaOpsRequest
}

func (k *KafkaOpsRequest) ResourceSingular() string {
	return ResourceSingularKafkaOpsRequest
}

func (k *KafkaOpsRequest) ResourcePlural() string {
	return ResourcePluralKafkaOpsRequest
}

var _ Accessor = &KafkaOpsRequest{}

func (k *KafkaOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return k.ObjectMeta
}

func (k *KafkaOpsRequest) GetDBRefName() string {
	return k.Spec.DatabaseRef.Name
}

func (k *KafkaOpsRequest) GetRequestType() any {
	return k.Spec.Type
}

func (k *KafkaOpsRequest) GetStatus() OpsRequestStatus {
	return k.Status
}

func (k *KafkaOpsRequest) SetStatus(s OpsRequestStatus) {
	k.Status = s
}
