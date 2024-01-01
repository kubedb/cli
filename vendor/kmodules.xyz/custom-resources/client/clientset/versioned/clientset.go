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

package versioned

import (
	"fmt"
	"net/http"

	appcatalogv1alpha1 "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
	auditorv1alpha1 "kmodules.xyz/custom-resources/client/clientset/versioned/typed/auditor/v1alpha1"
	metricsv1alpha1 "kmodules.xyz/custom-resources/client/clientset/versioned/typed/metrics/v1alpha1"

	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	AppcatalogV1alpha1() appcatalogv1alpha1.AppcatalogV1alpha1Interface
	AuditorV1alpha1() auditorv1alpha1.AuditorV1alpha1Interface
	MetricsV1alpha1() metricsv1alpha1.MetricsV1alpha1Interface
}

// Clientset contains the clients for groups.
type Clientset struct {
	*discovery.DiscoveryClient
	appcatalogV1alpha1 *appcatalogv1alpha1.AppcatalogV1alpha1Client
	auditorV1alpha1    *auditorv1alpha1.AuditorV1alpha1Client
	metricsV1alpha1    *metricsv1alpha1.MetricsV1alpha1Client
}

// AppcatalogV1alpha1 retrieves the AppcatalogV1alpha1Client
func (c *Clientset) AppcatalogV1alpha1() appcatalogv1alpha1.AppcatalogV1alpha1Interface {
	return c.appcatalogV1alpha1
}

// AuditorV1alpha1 retrieves the AuditorV1alpha1Client
func (c *Clientset) AuditorV1alpha1() auditorv1alpha1.AuditorV1alpha1Interface {
	return c.auditorV1alpha1
}

// MetricsV1alpha1 retrieves the MetricsV1alpha1Client
func (c *Clientset) MetricsV1alpha1() metricsv1alpha1.MetricsV1alpha1Interface {
	return c.metricsV1alpha1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c

	if configShallowCopy.UserAgent == "" {
		configShallowCopy.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	// share the transport between all clients
	httpClient, err := rest.HTTPClientFor(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	return NewForConfigAndClient(&configShallowCopy, httpClient)
}

// NewForConfigAndClient creates a new Clientset for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfigAndClient will generate a rate-limiter in configShallowCopy.
func NewForConfigAndClient(c *rest.Config, httpClient *http.Client) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}

	var cs Clientset
	var err error
	cs.appcatalogV1alpha1, err = appcatalogv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.auditorV1alpha1, err = auditorv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.metricsV1alpha1, err = metricsv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	cs, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.appcatalogV1alpha1 = appcatalogv1alpha1.New(c)
	cs.auditorV1alpha1 = auditorv1alpha1.New(c)
	cs.metricsV1alpha1 = metricsv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
