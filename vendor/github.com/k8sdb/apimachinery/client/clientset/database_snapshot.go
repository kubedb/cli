package clientset

import (
	aci "github.com/k8sdb/apimachinery/api"
	"k8s.io/kubernetes/pkg/api"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/watch"
)

type DatabaseSnapshotNamespacer interface {
	DatabaseSnapshots(namespace string) DatabaseSnapshotInterface
}

type DatabaseSnapshotInterface interface {
	List(opts api.ListOptions) (*aci.DatabaseSnapshotList, error)
	Get(name string) (*aci.DatabaseSnapshot, error)
	Create(databasesnapshot *aci.DatabaseSnapshot) (*aci.DatabaseSnapshot, error)
	Update(databasesnapshot *aci.DatabaseSnapshot) (*aci.DatabaseSnapshot, error)
	Delete(name string) error
	Watch(opts api.ListOptions) (watch.Interface, error)
	UpdateStatus(databasesnapshot *aci.DatabaseSnapshot) (*aci.DatabaseSnapshot, error)
}

type DatabaseSnapshotImpl struct {
	r  rest.Interface
	ns string
}

func newDatabaseSnapshot(c *ExtensionsClient, namespace string) *DatabaseSnapshotImpl {
	return &DatabaseSnapshotImpl{c.restClient, namespace}
}

func (c *DatabaseSnapshotImpl) List(opts api.ListOptions) (result *aci.DatabaseSnapshotList, err error) {
	result = &aci.DatabaseSnapshotList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDatabaseSnapshot).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *DatabaseSnapshotImpl) Get(name string) (result *aci.DatabaseSnapshot, err error) {
	result = &aci.DatabaseSnapshot{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDatabaseSnapshot).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *DatabaseSnapshotImpl) Create(databasesnapshot *aci.DatabaseSnapshot) (result *aci.DatabaseSnapshot, err error) {
	result = &aci.DatabaseSnapshot{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDatabaseSnapshot).
		Body(databasesnapshot).
		Do().
		Into(result)
	return
}

func (c *DatabaseSnapshotImpl) Update(databasesnapshot *aci.DatabaseSnapshot) (result *aci.DatabaseSnapshot, err error) {
	result = &aci.DatabaseSnapshot{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDatabaseSnapshot).
		Name(databasesnapshot.Name).
		Body(databasesnapshot).
		Do().
		Into(result)
	return
}

func (c *DatabaseSnapshotImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDatabaseSnapshot).
		Name(name).
		Do().
		Error()
}

func (c *DatabaseSnapshotImpl) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeDatabaseSnapshot).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *DatabaseSnapshotImpl) UpdateStatus(databasesnapshot *aci.DatabaseSnapshot) (result *aci.DatabaseSnapshot, err error) {
	result = &aci.DatabaseSnapshot{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDatabaseSnapshot).
		Name(databasesnapshot.Name).
		SubResource("status").
		Body(databasesnapshot).
		Do().
		Into(result)
	return
}
