package validator

import (
	"fmt"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	adr "github.com/kubedb/apimachinery/pkg/docker"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	dr "github.com/kubedb/redis/pkg/docker"
	"k8s.io/client-go/kubernetes"
)

func ValidateRedis(client kubernetes.Interface, redis *api.Redis, docker dr.Docker) error {
	if redis.Spec.Version == "" {
		return fmt.Errorf(`Object 'Version' is missing in '%v'`, redis.Spec)
	}

	// Set Database Image version
	version := string(redis.Spec.Version)
	if err := adr.CheckDockerImageVersion(docker.GetImage(redis), version); err != nil {
		return fmt.Errorf(`Image %v not found`, docker.GetImageWithTag(redis))
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
