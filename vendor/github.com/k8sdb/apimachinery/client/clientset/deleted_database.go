package clientset

import (
	aci "github.com/k8sdb/apimachinery/api"
	"k8s.io/kubernetes/pkg/api"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/watch"
)

type DeletedDatabaseNamespacer interface {
	DeletedDatabases(namespace string) DeletedDatabaseInterface
}

type DeletedDatabaseInterface interface {
	List(opts api.ListOptions) (*aci.DeletedDatabaseList, error)
	Get(name string) (*aci.DeletedDatabase, error)
	Create(deleteddatabase *aci.DeletedDatabase) (*aci.DeletedDatabase, error)
	Update(deleteddatabase *aci.DeletedDatabase) (*aci.DeletedDatabase, error)
	Delete(name string) error
	Watch(opts api.ListOptions) (watch.Interface, error)
	UpdateStatus(deleteddatabase *aci.DeletedDatabase) (*aci.DeletedDatabase, error)
}

type DeletedDatabaseImpl struct {
	r  rest.Interface
	ns string
}

func newDeletedDatabase(c *ExtensionsClient, namespace string) *DeletedDatabaseImpl {
	return &DeletedDatabaseImpl{c.restClient, namespace}
}

func (c *DeletedDatabaseImpl) List(opts api.ListOptions) (result *aci.DeletedDatabaseList, err error) {
	result = &aci.DeletedDatabaseList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDeletedDatabase).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *DeletedDatabaseImpl) Get(name string) (result *aci.DeletedDatabase, err error) {
	result = &aci.DeletedDatabase{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDeletedDatabase).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *DeletedDatabaseImpl) Create(deleteddatabase *aci.DeletedDatabase) (result *aci.DeletedDatabase, err error) {
	result = &aci.DeletedDatabase{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDeletedDatabase).
		Body(deleteddatabase).
		Do().
		Into(result)
	return
}

func (c *DeletedDatabaseImpl) Update(deleteddatabase *aci.DeletedDatabase) (result *aci.DeletedDatabase, err error) {
	result = &aci.DeletedDatabase{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDeletedDatabase).
		Name(deleteddatabase.Name).
		Body(deleteddatabase).
		Do().
		Into(result)
	return
}

func (c *DeletedDatabaseImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDeletedDatabase).
		Name(name).
		Do().
		Error()
}

func (c *DeletedDatabaseImpl) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeDeletedDatabase).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *DeletedDatabaseImpl) UpdateStatus(deleteddatabase *aci.DeletedDatabase) (result *aci.DeletedDatabase, err error) {
	result = &aci.DeletedDatabase{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDeletedDatabase).
		Name(deleteddatabase.Name).
		SubResource("status").
		Body(deleteddatabase).
		Do().
		Into(result)
	return
}
