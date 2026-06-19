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

func (r *AerospikeOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralAerospikeOpsRequest))
}

var _ apis.ResourceInfo = &AerospikeOpsRequest{}

func (r *AerospikeOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralAerospikeOpsRequest, ops.GroupName)
}

func (r *AerospikeOpsRequest) ResourceShortCode() string {
	return ResourceCodeAerospikeOpsRequest
}

func (r *AerospikeOpsRequest) ResourceKind() string {
	return ResourceKindAerospikeOpsRequest
}

func (r *AerospikeOpsRequest) ResourceSingular() string {
	return ResourceSingularAerospikeOpsRequest
}

func (r *AerospikeOpsRequest) ResourcePlural() string {
	return ResourcePluralAerospikeOpsRequest
}

var _ Accessor = &AerospikeOpsRequest{}

func (r *AerospikeOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r *AerospikeOpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *AerospikeOpsRequest) GetRequestType() string {
	return string(r.Spec.Type)
}

func (r *AerospikeOpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *AerospikeOpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
