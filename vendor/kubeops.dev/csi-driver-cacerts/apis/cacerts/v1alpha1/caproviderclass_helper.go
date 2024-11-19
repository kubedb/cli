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
	"kubeops.dev/csi-driver-cacerts/crds"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"kmodules.xyz/client-go/apiextensions"
)

func (_ CAProviderClass) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceCAProviderClasses))
}

type ObjectRef struct {
	APIGroup  string `json:"apiGroup"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
	Key       string `json:"key,omitempty"`
}

func RefFrom(pc CAProviderClass, ref TypedObjectReference) ObjectRef {
	result := ObjectRef{
		APIGroup:  "",
		Kind:      ref.Kind,
		Namespace: ref.Namespace,
		Name:      ref.Name,
		Key:       ref.Key,
	}
	if ref.APIGroup != nil {
		result.APIGroup = *ref.APIGroup
	} else {
		result.APIGroup = "v1"
	}
	if result.Namespace == "" {
		result.Namespace = pc.Namespace
	}
	return result
}

func (ref ObjectRef) GroupKind() schema.GroupKind {
	return schema.GroupKind{Group: ref.APIGroup, Kind: ref.Kind}
}

func (ref ObjectRef) ObjKey() types.NamespacedName {
	return types.NamespacedName{
		Namespace: ref.Namespace,
		Name:      ref.Name,
	}
}
