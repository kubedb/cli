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

func (_ EtcdOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralEtcdOpsRequest))
}

var _ apis.ResourceInfo = &EtcdOpsRequest{}

func (e EtcdOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralEtcdOpsRequest, ops.GroupName)
}

func (e EtcdOpsRequest) ResourceShortCode() string {
	return ResourceCodeEtcdOpsRequest
}

func (e EtcdOpsRequest) ResourceKind() string {
	return ResourceKindEtcdOpsRequest
}

func (e EtcdOpsRequest) ResourceSingular() string {
	return ResourceSingularEtcdOpsRequest
}

func (e EtcdOpsRequest) ResourcePlural() string {
	return ResourcePluralEtcdOpsRequest
}

func (e EtcdOpsRequest) ValidateSpecs() error {
	return nil
}

var _ Accessor = &EtcdOpsRequest{}

func (e *EtcdOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return e.ObjectMeta
}

func (e EtcdOpsRequest) GetRequestType() any {
	return e.Spec.Type
}

func (e EtcdOpsRequest) GetUpdateVersionSpec() *EtcdUpdateVersionSpec {
	return e.Spec.UpdateVersion
}

func (e *EtcdOpsRequest) GetDBRefName() string {
	return e.Spec.DatabaseRef.Name
}

func (e *EtcdOpsRequest) GetStatus() OpsRequestStatus {
	return e.Status
}

func (e *EtcdOpsRequest) SetStatus(s OpsRequestStatus) {
	e.Status = s
}
