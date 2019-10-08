package v1alpha1

import (
	"fmt"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kubedb.dev/apimachinery/apis"
)

var _ apis.ResourceInfo = &MySQLVersion{}

func (m MySQLVersion) ResourceShortCode() string {
	return ResourceCodeMySQLVersion
}

func (m MySQLVersion) ResourceKind() string {
	return ResourceKindMySQLVersion
}

func (m MySQLVersion) ResourceSingular() string {
	return ResourceSingularMySQLVersion
}

func (m MySQLVersion) ResourcePlural() string {
	return ResourcePluralMySQLVersion
}

func (m MySQLVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralMySQLVersion,
		Singular:      ResourceSingularMySQLVersion,
		Kind:          ResourceKindMySQLVersion,
		ShortNames:    []string{ResourceCodeMySQLVersion},
		Categories:    []string{"datastore", "kubedb", "appscode"},
		ResourceScope: string(apiextensions.ClusterScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.MySQLVersion",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: false,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "Version",
				Type:     "string",
				JSONPath: ".spec.version",
			},
			{
				Name:     "DB_IMAGE",
				Type:     "string",
				JSONPath: ".spec.db.image",
			},
			{
				Name:     "Deprecated",
				Type:     "boolean",
				JSONPath: ".spec.deprecated",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	})
}

func (m MySQLVersion) ValidateSpecs() error {
	if m.Spec.Version == "" ||
		m.Spec.DB.Image == "" ||
		m.Spec.Tools.Image == "" ||
		m.Spec.Exporter.Image == "" ||
		m.Spec.InitContainer.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for mysqlVersion "%v":
spec.version,
spec.db.image,
spec.tools.image,
spec.exporter.image,
spec.initContainer.image.`, m.Name)
	}
	return nil
}
