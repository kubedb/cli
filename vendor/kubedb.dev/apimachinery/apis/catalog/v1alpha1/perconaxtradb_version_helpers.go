package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
)

var _ apis.ResourceInfo = &PerconaXtraDBVersion{}

func (p PerconaXtraDBVersion) ResourceShortCode() string {
	return ResourceCodePerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourceKind() string {
	return ResourceKindPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourceSingular() string {
	return ResourceSingularPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourcePlural() string {
	return ResourcePluralPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralPerconaXtraDBVersion,
		Singular:      ResourceSingularPerconaXtraDBVersion,
		Kind:          ResourceKindPerconaXtraDBVersion,
		ShortNames:    []string{ResourceCodePerconaXtraDBVersion},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.PerconaXtraDBVersion",
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

func (p PerconaXtraDBVersion) ValidateSpecs() error {
	if p.Spec.Version == "" ||
		p.Spec.DB.Image == "" ||
		p.Spec.Exporter.Image == "" ||
		p.Spec.InitContainer.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for perconaxtradbversion "%v":
spec.version,
spec.db.image,
spec.exporter.image,
spec.initContainer.image.`, p.Name)
	}
	return nil
}
