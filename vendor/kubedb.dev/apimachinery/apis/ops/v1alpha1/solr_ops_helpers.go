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

func (s *SolrOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralSolrOpsRequest))
}

var _ apis.ResourceInfo = &SolrOpsRequest{}

func (s *SolrOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralSolrOpsRequest, ops.GroupName)
}

func (s *SolrOpsRequest) ResourceShortCode() string {
	return ResourceCodeSolrOpsRequest
}

func (s *SolrOpsRequest) ResourceKind() string {
	return ResourceKindSolrOpsRequest
}

func (s *SolrOpsRequest) ResourceSingular() string {
	return ResourceSingularSolrOpsRequest
}

func (s *SolrOpsRequest) ResourcePlural() string {
	return ResourcePluralSolrOpsRequest
}

var _ Accessor = &SolrOpsRequest{}

func (s *SolrOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return s.ObjectMeta
}

func (s *SolrOpsRequest) GetDBRefName() string {
	return s.Spec.DatabaseRef.Name
}

func (s *SolrOpsRequest) GetRequestType() any {
	return s.Spec.Type
}

func (s *SolrOpsRequest) GetStatus() OpsRequestStatus {
	return s.Status
}

func (s *SolrOpsRequest) SetStatus(st OpsRequestStatus) {
	s.Status = st
}
