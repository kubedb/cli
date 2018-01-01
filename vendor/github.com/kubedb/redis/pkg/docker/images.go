package docker

import (
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
)

const (
	ImageKubedbOperator = "operator"
	ImageRedisOperator  = "rd-operator"
	ImageRedis          = "redis"
)

type Docker struct {
	// Docker Registry
	Registry string
	// Exporter tag
	ExporterTag string
}

func (d Docker) GetImage(redis *api.Redis) string {
	return d.Registry + "/" + ImageRedis
}

func (d Docker) GetImageWithTag(redis *api.Redis) string {
	return d.GetImage(redis) + ":" + string(redis.Spec.Version)
}

func (d Docker) GetOperatorImage(redis *api.Redis) string {
	return d.Registry + "/" + ImageKubedbOperator
}

func (d Docker) GetOperatorImageWithTag(redis *api.Redis) string {
	return d.GetOperatorImage(redis) + ":" + d.ExporterTag
}
