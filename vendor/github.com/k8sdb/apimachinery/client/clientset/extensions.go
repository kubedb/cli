package clientset

import (
	"fmt"

	schema "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apimachinery/registered"
	rest "k8s.io/kubernetes/pkg/client/restclient"
)

const (
	defaultAPIPath = "/apis"
)

type ExtensionInterface interface {
	RESTClient() rest.Interface
	DatabaseSnapshotNamespacer
	DeletedDatabaseNamespacer
	ElasticNamespacer
	PostgresNamespacer
}

// ExtensionsClient is used to interact with experimental Kubernetes features.
// Features of Extensions group are not supported and may be changed or removed in
// incompatible ways at any time.
type ExtensionsClient struct {
	restClient rest.Interface
}

var _ ExtensionInterface = &ExtensionsClient{}

func (a *ExtensionsClient) DatabaseSnapshots(namespace string) DatabaseSnapshotInterface {
	return newDatabaseSnapshot(a, namespace)
}

func (a *ExtensionsClient) DeletedDatabases(namespace string) DeletedDatabaseInterface {
	return newDeletedDatabase(a, namespace)
}

func (a *ExtensionsClient) Elastics(namespace string) ElasticInterface {
	return newElastic(a, namespace)
}

func (a *ExtensionsClient) Postgreses(namespace string) PostgresInterface {
	return newPostgres(a, namespace)
}

// NewAppsCodeExtensions creates a new ExtensionsClient for the given config. This client
// provides access to experimental Kubernetes features.
// Features of Extensions group are not supported and may be changed or removed in
// incompatible ways at any time.
func NewExtensionsForConfig(c *rest.Config) (*ExtensionsClient, error) {
	config := *c
	if err := setExtensionsDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &ExtensionsClient{client}, nil
}

// NewAppsCodeExtensionsOrDie creates a new ExtensionsClient for the given config and
// panics if there is an error in the config.
// Features of Extensions group are not supported and may be changed or removed in
// incompatible ways at any time.
func NewExtensionsForConfigOrDie(c *rest.Config) *ExtensionsClient {
	client, err := NewExtensionsForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ExtensionsV1beta1Client for the given RESTClient.
func NewExtensions(c rest.Interface) *ExtensionsClient {
	return &ExtensionsClient{c}
}

func setExtensionsDefaults(config *rest.Config) error {
	gv, err := schema.ParseGroupVersion("k8sdb.com/v1beta1")
	if err != nil {
		return err
	}
	// if k8sdb.com/v1beta1 is not enabled, return an error
	if !registered.IsEnabledVersion(gv) {
		return fmt.Errorf("k8sdb.com/v1beta1 is not enabled")
	}
	config.APIPath = defaultAPIPath
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	if config.GroupVersion == nil || config.GroupVersion.Group != "k8sdb.com" {
		g, err := registered.Group("k8sdb.com")
		if err != nil {
			return err
		}
		copyGroupVersion := g.GroupVersion
		config.GroupVersion = &copyGroupVersion
	}

	config.NegotiatedSerializer = DirectCodecFactory{extendedCodec: ExtendedCodec}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ExtensionsClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
