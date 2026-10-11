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

// Package gentype provides a generic, create-only client-go typed client
// for API types that are not real stored Kubernetes objects -- request/
// response "action" payloads (e.g. an aggregated-apiserver endpoint that
// only ever accepts POST) which have no metav1.ObjectMeta and thus no real
// identity, labels, owner references, etc.
//
// k8s.io/client-go/gentype's Client[T] (which every generated client built
// since https://github.com/kubernetes/kubernetes/pull/121439 embeds)
// requires its type parameter to implement metav1.Object, even when the
// generated client only ever exposes Create -- Client[T].Create doesn't
// call any metav1.Object method on T, but the constraint is declared on the
// whole generic struct rather than per-method, so every type needs it
// regardless of which verbs its generated client actually has.
//
// This package provides the same shape (Client, NewClient, Option,
// PrefersProtobuf) but constrained to plain runtime.Object, for use by
// kmodules/code-generator's client-gen when it detects a +genclient type
// that is both create-only (+genclient:onlyVerbs=create, no other verb) and
// has no metav1.ObjectMeta member -- see cmd/client-gen/generators/util.
// IsCreateOnly/HasObjectMeta in that fork. Get/List/Update/Delete/Watch/
// Patch/Apply all need real object identity and have no equivalent here;
// a type needing any of those must have a real metav1.ObjectMeta and use
// k8s.io/client-go/gentype directly instead.
package gentype

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// Client represents a create-only client, optionally namespaced, for a
// runtime.Object type with no real persisted identity.
type Client[T runtime.Object] struct {
	resource       string
	client         rest.Interface
	namespace      string // "" for non-namespaced clients
	newObject      func() T
	parameterCodec runtime.ParameterCodec

	prefersProtobuf bool
}

// Option configures a Client.
type Option[T runtime.Object] func(*Client[T])

// PrefersProtobuf marks the client as preferring protobuf, matching
// k8s.io/client-go/gentype.PrefersProtobuf.
func PrefersProtobuf[T runtime.Object]() Option[T] {
	return func(c *Client[T]) { c.prefersProtobuf = true }
}

// NewClient constructs a create-only client, namespaced or not. Non-
// namespaced clients are constructed by passing an empty namespace ("").
func NewClient[T runtime.Object](
	resource string, client rest.Interface, parameterCodec runtime.ParameterCodec, namespace string, emptyObjectCreator func() T,
	options ...Option[T],
) *Client[T] {
	c := &Client[T]{
		resource:       resource,
		client:         client,
		parameterCodec: parameterCodec,
		namespace:      namespace,
		newObject:      emptyObjectCreator,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// GetClient returns the REST interface.
func (c *Client[T]) GetClient() rest.Interface {
	return c.client
}

// GetNamespace returns the client's namespace, if any.
func (c *Client[T]) GetNamespace() string {
	return c.namespace
}

// Create takes the representation of a resource and creates it. Returns the
// server's representation of the resource, and an error, if there is any.
func (c *Client[T]) Create(ctx context.Context, obj T, opts metav1.CreateOptions) (T, error) {
	result := c.newObject()
	err := c.client.Post().
		UseProtobufAsDefaultIfPreferred(c.prefersProtobuf).
		NamespaceIfScoped(c.namespace, c.namespace != "").
		Resource(c.resource).
		VersionedParams(&opts, c.parameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)
	return result, err
}
