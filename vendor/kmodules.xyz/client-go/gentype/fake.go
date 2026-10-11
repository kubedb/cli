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

package gentype

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/testing"
)

// FakeClient represents a fake create-only client for a runtime.Object type
// with no real persisted identity, matching Client's shape.
type FakeClient[T runtime.Object] struct {
	*testing.Fake
	ns        string
	resource  schema.GroupVersionResource
	kind      schema.GroupVersionKind
	newObject func() T
}

// NewFakeClient constructs a fake create-only client, namespaced or not.
// Non-namespaced clients are constructed by passing an empty namespace ("").
func NewFakeClient[T runtime.Object](
	fake *testing.Fake, namespace string, resource schema.GroupVersionResource, kind schema.GroupVersionKind, emptyObjectCreator func() T,
) *FakeClient[T] {
	return &FakeClient[T]{fake, namespace, resource, kind, emptyObjectCreator}
}

// Create takes the representation of a resource and creates it. Returns the
// server's representation of the resource, and an error, if there is any.
func (c *FakeClient[T]) Create(ctx context.Context, resource T, opts metav1.CreateOptions) (result T, err error) {
	emptyResult := c.newObject()
	obj, err := c.Invokes(testing.NewCreateActionWithOptions(c.resource, c.ns, resource, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(T), err
}
