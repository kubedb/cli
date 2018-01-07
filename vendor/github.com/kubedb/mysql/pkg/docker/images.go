package docker

import (
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
)

const (
	ImageKubedbOperator = "operator"
	ImageMySQLOperator  = "my-operator"
	ImageMySQL          = "mysql"
	ImageMySQLTools     = "mysql-tools"
)

type Docker struct {
	// Docker Registry
	Registry string
	// Exporter tag
	ExporterTag string
}

func (d Docker) GetImage(mysql *api.MySQL) string {
	return d.Registry + "/" + ImageMySQL
}

func (d Docker) GetImageWithTag(mysql *api.MySQL) string {
	return d.GetImage(mysql) + ":" + string(mysql.Spec.Version)
}

func (d Docker) GetOperatorImage(mysql *api.MySQL) string {
	return d.Registry + "/" + ImageKubedbOperator
}

func (d Docker) GetOperatorImageWithTag(mysql *api.MySQL) string {
	return d.GetOperatorImage(mysql) + ":" + d.ExporterTag
}

func (d Docker) GetToolsImage(mysql *api.MySQL) string {
	return d.Registry + "/" + ImageMySQLTools
}

func (d Docker) GetToolsImageWithTag(mysql *api.MySQL) string {
	return d.GetToolsImage(mysql) + ":" + string(mysql.Spec.Version)
}
