package docker

import (
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
)

const (
	ImageKubedbOperator        = "operator"
	ImageElasticsearchOperator = "es-operator"
	ImageElasticsearch         = "elasticsearch"
	ImageElasticsearchTools    = "elasticsearch-tools"
)

type Docker struct {
	// Docker Registry
	Registry string
	// Exporter tag
	ExporterTag string
}

func (d Docker) GetImage(elasticsearch *api.Elasticsearch) string {
	return d.Registry + "/" + ImageElasticsearch
}

func (d Docker) GetImageWithTag(elasticsearch *api.Elasticsearch) string {
	return d.GetImage(elasticsearch) + ":" + string(elasticsearch.Spec.Version)
}

func (d Docker) GetOperatorImage(elasticsearch *api.Elasticsearch) string {
	return d.Registry + "/" + ImageKubedbOperator
}

func (d Docker) GetOperatorImageWithTag(elasticsearch *api.Elasticsearch) string {
	return d.GetOperatorImage(elasticsearch) + ":" + d.ExporterTag
}

func (d Docker) GetToolsImage(elasticsearch *api.Elasticsearch) string {
	return d.Registry + "/" + ImageElasticsearchTools
}

func (d Docker) GetToolsImageWithTag(elasticsearch *api.Elasticsearch) string {
	return d.GetToolsImage(elasticsearch) + ":" + string(elasticsearch.Spec.Version)
}
