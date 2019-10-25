package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
)

var _ apis.ResourceInfo = &ElasticsearchVersion{}

func (e ElasticsearchVersion) ResourceShortCode() string {
	return ResourceCodeElasticsearchVersion
}

func (e ElasticsearchVersion) ResourceKind() string {
	return ResourceKindElasticsearchVersion
}

func (e ElasticsearchVersion) ResourceSingular() string {
	return ResourceSingularElasticsearchVersion
}

func (e ElasticsearchVersion) ResourcePlural() string {
	return ResourcePluralElasticsearchVersion
}

func (e ElasticsearchVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralElasticsearchVersion,
		Singular:      ResourceSingularElasticsearchVersion,
		Kind:          ResourceKindElasticsearchVersion,
		ShortNames:    []string{ResourceCodeElasticsearchVersion},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.ElasticsearchVersion",
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

func (e ElasticsearchVersion) ValidateSpecs() error {
	if e.Spec.AuthPlugin == "" ||
		e.Spec.Version == "" ||
		e.Spec.DB.Image == "" ||
		e.Spec.Tools.Image == "" ||
		e.Spec.Exporter.Image == "" ||
		e.Spec.InitContainer.YQImage == "" ||
		e.Spec.InitContainer.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for elasticsearchVersion "%v":
spec.authPlugin,
spec.version,
spec.db.image,
spec.tools.image,
spec.exporter.image,
spec.initContainer.yqImage,
spec.initContainer.image.`, e.Name)
	}
	return nil
}
