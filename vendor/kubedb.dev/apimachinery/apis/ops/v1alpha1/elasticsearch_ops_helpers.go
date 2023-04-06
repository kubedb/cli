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

func (_ ElasticsearchOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralElasticsearchOpsRequest))
}

var _ apis.ResourceInfo = &ElasticsearchOpsRequest{}

func (e ElasticsearchOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralElasticsearchOpsRequest, ops.GroupName)
}

func (e ElasticsearchOpsRequest) ResourceShortCode() string {
	return ResourceCodeElasticsearchOpsRequest
}

func (e ElasticsearchOpsRequest) ResourceKind() string {
	return ResourceKindElasticsearchOpsRequest
}

func (e ElasticsearchOpsRequest) ResourceSingular() string {
	return ResourceSingularElasticsearchOpsRequest
}

func (e ElasticsearchOpsRequest) ResourcePlural() string {
	return ResourcePluralElasticsearchOpsRequest
}

func (e ElasticsearchOpsRequest) ValidateSpecs() error {
	return nil
}

var _ Accessor = &ElasticsearchOpsRequest{}

func (e *ElasticsearchOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return e.ObjectMeta
}

func (e ElasticsearchOpsRequest) GetRequestType() any {
	switch e.Spec.Type {
	case ElasticsearchOpsRequestTypeUpgrade:
		return ElasticsearchOpsRequestTypeUpdateVersion
	}
	return e.Spec.Type
}

func (e ElasticsearchOpsRequest) GetUpdateVersionSpec() *ElasticsearchUpdateVersionSpec {
	if e.Spec.UpdateVersion != nil {
		return e.Spec.UpdateVersion
	}
	return e.Spec.Upgrade
}

func (e *ElasticsearchOpsRequest) GetDBRefName() string {
	return e.Spec.DatabaseRef.Name
}

func (e *ElasticsearchOpsRequest) GetStatus() OpsRequestStatus {
	return e.Status
}

func (e *ElasticsearchOpsRequest) SetStatus(s OpsRequestStatus) {
	e.Status = s
}
