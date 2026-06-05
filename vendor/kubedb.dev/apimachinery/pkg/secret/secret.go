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

// Package secret provides helpers that operate uniformly over a Kubernetes
// core/v1 Secret and a virtual-secrets.dev/v1alpha1 Secret. Callers pass an
// isVirtual flag (typically the result of
// kubedb.dev/apimachinery/apis/kubedb/v{1,1alpha2}.IsVirtualAuthSecretReferred)
// and the helper dispatches to the matching object type. Both types expose
// `metav1.ObjectMeta` (embedded) and `Data map[string][]byte`, so the mutation
// logic in mutator callbacks is identical across the two paths.
package secret

import (
	"context"

	vsecretapi "go.virtual-secrets.dev/apimachinery/apis/virtual/v1alpha1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	cu "kmodules.xyz/client-go/client"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Mutator updates the annotations and data of an auth secret in-place. It
// receives the current annotations/data maps (either may be nil for a
// not-yet-existing object) and returns the maps that should be persisted.
// Allocate a new map if the input is nil.
type Mutator func(annotations map[string]string, data map[string][]byte) (map[string]string, map[string][]byte)

// GetData fetches the auth secret and returns its Data map. When isVirtual is
// true the secret is read as a vsecretapi.Secret; otherwise as a core.Secret.
// Errors (including NotFound) are returned to the caller unchanged.
func GetData(ctx context.Context, c client.Client, namespace, name string, isVirtual bool) (map[string][]byte, error) {
	key := types.NamespacedName{Namespace: namespace, Name: name}
	if isVirtual {
		s := &vsecretapi.Secret{}
		if err := c.Get(ctx, key, s); err != nil {
			return nil, err
		}
		return s.Data, nil
	}
	s := &core.Secret{}
	if err := c.Get(ctx, key, s); err != nil {
		return nil, err
	}
	return s.Data, nil
}

// GetAnnotations fetches the auth secret and returns its Annotations map.
func GetAnnotations(ctx context.Context, c client.Client, namespace, name string, isVirtual bool) (map[string]string, error) {
	key := types.NamespacedName{Namespace: namespace, Name: name}
	if isVirtual {
		s := &vsecretapi.Secret{}
		if err := c.Get(ctx, key, s); err != nil {
			return nil, err
		}
		return s.Annotations, nil
	}
	s := &core.Secret{}
	if err := c.Get(ctx, key, s); err != nil {
		return nil, err
	}
	return s.Annotations, nil
}

// Get fetches the auth secret and returns Data and Annotations in one call.
func Get(ctx context.Context, c client.Client, namespace, name string, isVirtual bool) (data map[string][]byte, annotations map[string]string, err error) {
	key := types.NamespacedName{Namespace: namespace, Name: name}
	if isVirtual {
		s := &vsecretapi.Secret{}
		if err := c.Get(ctx, key, s); err != nil {
			return nil, nil, err
		}
		return s.Data, s.Annotations, nil
	}
	s := &core.Secret{}
	if err := c.Get(ctx, key, s); err != nil {
		return nil, nil, err
	}
	return s.Data, s.Annotations, nil
}

// CreateOrPatch applies mutate to the auth secret using kmodules
// client-go CreateOrPatch semantics. The same mutator runs against either
// secret type because both expose Annotations (via ObjectMeta) and Data.
func CreateOrPatch(ctx context.Context, c client.Client, namespace, name string, isVirtual bool, mutate Mutator) error {
	if isVirtual {
		_, err := cu.CreateOrPatch(ctx, c, &vsecretapi.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		}, func(obj client.Object, createOp bool) client.Object {
			ret := obj.(*vsecretapi.Secret)
			ret.Annotations, ret.Data = mutate(ret.Annotations, ret.Data)
			return ret
		})
		return err
	}
	_, err := cu.CreateOrPatch(ctx, c, &core.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
	}, func(obj client.Object, createOp bool) client.Object {
		ret := obj.(*core.Secret)
		ret.Annotations, ret.Data = mutate(ret.Annotations, ret.Data)
		return ret
	})
	return err
}

// DropAnnotation removes a single annotation key from the auth secret via
// CreateOrPatch. The Data is preserved.
func DropAnnotation(ctx context.Context, c client.Client, namespace, name string, isVirtual bool, key string) error {
	return CreateOrPatch(ctx, c, namespace, name, isVirtual, func(annotations map[string]string, data map[string][]byte) (map[string]string, map[string][]byte) {
		delete(annotations, key)
		return annotations, data
	})
}
