package v1alpha1

import (
	"fmt"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kubedb.dev/apimachinery/apis"
)

var _ apis.ResourceInfo = &MongoDBVersion{}

func (m MongoDBVersion) ResourceShortCode() string {
	return ResourceCodeMongoDBVersion
}

func (m MongoDBVersion) ResourceKind() string {
	return ResourceKindMongoDBVersion
}

func (m MongoDBVersion) ResourceSingular() string {
	return ResourceSingularMongoDBVersion
}

func (m MongoDBVersion) ResourcePlural() string {
	return ResourcePluralMongoDBVersion
}

func (m MongoDBVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralMongoDBVersion,
		Singular:      ResourceSingularMongoDBVersion,
		Kind:          ResourceKindMongoDBVersion,
		ShortNames:    []string{ResourceCodeMongoDBVersion},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.MongoDBVersion",
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

func (m MongoDBVersion) ValidateSpecs() error {
	if m.Spec.Version == "" ||
		m.Spec.DB.Image == "" ||
		m.Spec.Tools.Image == "" ||
		m.Spec.Exporter.Image == "" ||
		m.Spec.InitContainer.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for mongodbVersion "%v":
spec.version,
spec.db.image,
spec.tools.image,
spec.exporter.image,
spec.initContainer.image.`, m.Name)
	}
	return nil
}
