package docker

import (
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
)

const (
	ImageKubedbOperator  = "operator"
	ImageMongoDBOperator = "mg-operator"
	ImageMongoDB         = "mongo"
	ImageMongoDBTools    = "mongo-tools"
)

type Docker struct {
	// Docker Registry
	Registry string
	// Exporter tag
	ExporterTag string
}

func (d Docker) GetImage(mongodb *api.MongoDB) string {
	return d.Registry + "/" + ImageMongoDB
}

func (d Docker) GetImageWithTag(mongodb *api.MongoDB) string {
	return d.GetImage(mongodb) + ":" + string(mongodb.Spec.Version)
}

func (d Docker) GetOperatorImage(mongodb *api.MongoDB) string {
	return d.Registry + "/" + ImageKubedbOperator
}

func (d Docker) GetOperatorImageWithTag(mongodb *api.MongoDB) string {
	return d.GetOperatorImage(mongodb) + ":" + d.ExporterTag
}

func (d Docker) GetToolsImage(mongodb *api.MongoDB) string {
	return d.Registry + "/" + ImageMongoDBTools
}

func (d Docker) GetToolsImageWithTag(mongodb *api.MongoDB) string {
	return d.GetToolsImage(mongodb) + ":" + string(mongodb.Spec.Version)
}
