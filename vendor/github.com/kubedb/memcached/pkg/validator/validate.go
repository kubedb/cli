package validator

import (
	"fmt"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/pkg/docker"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	"k8s.io/client-go/kubernetes"
)

func ValidateMemcached(client kubernetes.Interface, memcached *api.Memcached) error {
	if memcached.Spec.Version == "" {
		return fmt.Errorf(`Object 'Version' is missing in '%v'`, memcached.Spec)
	}

	// Set Database Image version
	version := string(memcached.Spec.Version)
	if err := docker.CheckDockerImageVersion(docker.ImageMemcached, version); err != nil {
		return fmt.Errorf(`Image %v:%v not found`, docker.ImageMemcached, version)
	}

	monitorSpec := memcached.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}
	}
	return nil
}
