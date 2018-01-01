package docker

import (
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
)

const (
	ImagePostgresOperator = "pg-operator"
	ImagePostgres         = "postgres"
	ImagePostgresTools    = "postgres-tools"
)

type Docker struct {
	// Docker Registry
	Registry string
	// Exporter tag
	ExporterTag string
}

func (d Docker) GetImage(postgres *api.Postgres) string {
	return d.Registry + "/" + ImagePostgres
}

func (d Docker) GetImageWithTag(postgres *api.Postgres) string {
	return d.GetImage(postgres) + ":" + string(postgres.Spec.Version)
}

func (d Docker) GetOperatorImage(postgres *api.Postgres) string {
	return d.Registry + "/" + ImagePostgresOperator
}

func (d Docker) GetOperatorImageWithTag(postgres *api.Postgres) string {
	return d.GetOperatorImage(postgres) + ":" + d.ExporterTag
}

func (d Docker) GetToolsImage(postgres *api.Postgres) string {
	return d.Registry + "/" + ImagePostgresTools
}

func (d Docker) GetToolsImageWithTag(postgres *api.Postgres) string {
	return d.GetToolsImage(postgres) + ":" + string(postgres.Spec.Version)
}
