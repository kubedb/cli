package validator

import (
	"fmt"

	api "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/k8sdb/apimachinery/pkg/docker"
	amv "github.com/k8sdb/apimachinery/pkg/validator"
	"k8s.io/client-go/kubernetes"
)

func ValidateRedis(client kubernetes.Interface, redis *api.Redis) error {
	if redis.Spec.Version == "" {
		return fmt.Errorf(`Object 'Version' is missing in '%v'`, redis.Spec)
	}

	// Set Database Image version
	version := string(redis.Spec.Version)
	if err := docker.CheckDockerImageVersion(docker.ImageRedis, version); err != nil {
		return fmt.Errorf(`Image %v:%v not found`, docker.ImageRedis, version)
	}

	if redis.Spec.Storage != nil {
		var err error
		if err = amv.ValidateStorage(client, redis.Spec.Storage); err != nil {
			return err
		}
	}

	monitorSpec := redis.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}

	}
	return nil
}
