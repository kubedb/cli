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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	scheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RabbitMQVersionsGetter has a method to return a RabbitMQVersionInterface.
// A group's client should implement this interface.
type RabbitMQVersionsGetter interface {
	RabbitMQVersions() RabbitMQVersionInterface
}

// RabbitMQVersionInterface has methods to work with RabbitMQVersion resources.
type RabbitMQVersionInterface interface {
	Create(ctx context.Context, rabbitMQVersion *v1alpha1.RabbitMQVersion, opts v1.CreateOptions) (*v1alpha1.RabbitMQVersion, error)
	Update(ctx context.Context, rabbitMQVersion *v1alpha1.RabbitMQVersion, opts v1.UpdateOptions) (*v1alpha1.RabbitMQVersion, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.RabbitMQVersion, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.RabbitMQVersionList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.RabbitMQVersion, err error)
	RabbitMQVersionExpansion
}

// rabbitMQVersions implements RabbitMQVersionInterface
type rabbitMQVersions struct {
	client rest.Interface
}

// newRabbitMQVersions returns a RabbitMQVersions
func newRabbitMQVersions(c *CatalogV1alpha1Client) *rabbitMQVersions {
	return &rabbitMQVersions{
		client: c.RESTClient(),
	}
}

// Get takes name of the rabbitMQVersion, and returns the corresponding rabbitMQVersion object, and an error if there is any.
func (c *rabbitMQVersions) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.RabbitMQVersion, err error) {
	result = &v1alpha1.RabbitMQVersion{}
	err = c.client.Get().
		Resource("rabbitmqversions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RabbitMQVersions that match those selectors.
func (c *rabbitMQVersions) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RabbitMQVersionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.RabbitMQVersionList{}
	err = c.client.Get().
		Resource("rabbitmqversions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested rabbitMQVersions.
func (c *rabbitMQVersions) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("rabbitmqversions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a rabbitMQVersion and creates it.  Returns the server's representation of the rabbitMQVersion, and an error, if there is any.
func (c *rabbitMQVersions) Create(ctx context.Context, rabbitMQVersion *v1alpha1.RabbitMQVersion, opts v1.CreateOptions) (result *v1alpha1.RabbitMQVersion, err error) {
	result = &v1alpha1.RabbitMQVersion{}
	err = c.client.Post().
		Resource("rabbitmqversions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(rabbitMQVersion).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a rabbitMQVersion and updates it. Returns the server's representation of the rabbitMQVersion, and an error, if there is any.
func (c *rabbitMQVersions) Update(ctx context.Context, rabbitMQVersion *v1alpha1.RabbitMQVersion, opts v1.UpdateOptions) (result *v1alpha1.RabbitMQVersion, err error) {
	result = &v1alpha1.RabbitMQVersion{}
	err = c.client.Put().
		Resource("rabbitmqversions").
		Name(rabbitMQVersion.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(rabbitMQVersion).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the rabbitMQVersion and deletes it. Returns an error if one occurs.
func (c *rabbitMQVersions) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("rabbitmqversions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *rabbitMQVersions) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("rabbitmqversions").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched rabbitMQVersion.
func (c *rabbitMQVersions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.RabbitMQVersion, err error) {
	result = &v1alpha1.RabbitMQVersion{}
	err = c.client.Patch(pt).
		Resource("rabbitmqversions").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
