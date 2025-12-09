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
	"kubestash.dev/apimachinery/apis"
	"kubestash.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	cutil "kmodules.xyz/client-go/conditions"
	"kmodules.xyz/client-go/meta"
)

func (Repository) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralRepository))
}

func (r *Repository) CalculatePhase() RepositoryPhase {
	if cutil.IsConditionTrue(r.Status.Conditions, TypeRepositoryInitialized) {
		return RepositoryReady
	}
	return RepositoryNotReady
}

func (r *Repository) OffshootLabels() map[string]string {
	newLabels := make(map[string]string)
	newLabels[meta.ManagedByLabelKey] = apis.KubeStashKey
	newLabels[apis.KubeStashInvokerKind] = ResourceKindRepository
	newLabels[apis.KubeStashInvokerName] = r.Name
	newLabels[apis.KubeStashInvokerNamespace] = r.Namespace
	return apis.UpsertLabels(r.Labels, newLabels)
}
