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

func (z *ZooKeeperOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralZooKeeperOpsRequest))
}

var _ apis.ResourceInfo = &ZooKeeperOpsRequest{}

func (z *ZooKeeperOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralZooKeeperOpsRequest, ops.GroupName)
}

func (z *ZooKeeperOpsRequest) ResourceShortCode() string {
	return ResourceCodeZooKeeperOpsRequest
}

func (z *ZooKeeperOpsRequest) ResourceKind() string {
	return ResourceKindZooKeeperOpsRequest
}

func (z *ZooKeeperOpsRequest) ResourceSingular() string {
	return ResourceSingularZooKeeperOpsRequest
}

func (z *ZooKeeperOpsRequest) ResourcePlural() string {
	return ResourcePluralZooKeeperOpsRequest
}

var _ Accessor = &ZooKeeperOpsRequest{}

func (z *ZooKeeperOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return z.ObjectMeta
}

func (z *ZooKeeperOpsRequest) GetDBRefName() string {
	return z.Spec.DatabaseRef.Name
}

func (z *ZooKeeperOpsRequest) GetRequestType() string {
	return string(z.Spec.Type)
}

func (z *ZooKeeperOpsRequest) GetStatus() OpsRequestStatus {
	return z.Status
}

func (z *ZooKeeperOpsRequest) SetStatus(st OpsRequestStatus) {
	z.Status = st
}
