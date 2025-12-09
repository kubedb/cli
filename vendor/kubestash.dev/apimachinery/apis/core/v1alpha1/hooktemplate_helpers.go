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

	corev1 "k8s.io/api/core/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
)

func (HookTemplate) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralHookTemplate))
}

func (h *HookTemplate) UsageAllowed(srcNamespace *corev1.Namespace) bool {
	allowedNamespace := h.Spec.UsagePolicy.AllowedNamespaces
	if *allowedNamespace.From == apis.NamespacesFromAll {
		return true
	}

	if *allowedNamespace.From == apis.NamespacesFromSame {
		return h.Namespace == srcNamespace.Name
	}

	return selectorMatches(allowedNamespace.Selector, srcNamespace.Labels)
}

func (h *HookTemplate) OffshootLabels() map[string]string {
	newLabels := make(map[string]string)
	newLabels[meta_util.ManagedByLabelKey] = apis.KubeStashKey
	newLabels[apis.KubeStashInvokerName] = h.Name
	newLabels[apis.KubeStashInvokerNamespace] = h.Namespace

	return apis.UpsertLabels(h.Labels, newLabels)
}
