package validator

import (
	"fmt"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	adr "github.com/kubedb/apimachinery/pkg/docker"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	dr "github.com/kubedb/memcached/pkg/docker"
	"k8s.io/client-go/kubernetes"
)

func ValidateMemcached(client kubernetes.Interface, memcached *api.Memcached, docker dr.Docker) error {
	if memcached.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, memcached.Spec)
	}

	version := string(memcached.Spec.Version)
	if err := adr.CheckDockerImageVersion(docker.GetImage(memcached), version); err != nil {
		return fmt.Errorf(`image %s not found`, docker.GetImageWithTag(memcached))
	}

	monitorSpec := memcached.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}
	}
	return nil
}
