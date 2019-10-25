package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
)

var _ apis.ResourceInfo = &PgBouncerVersion{}

func (p PgBouncerVersion) ResourceShortCode() string {
	return ResourceCodePgBouncerVersion
}

func (p PgBouncerVersion) ResourceKind() string {
	return ResourceKindPgBouncerVersion
}

func (p PgBouncerVersion) ResourceSingular() string {
	return ResourceSingularPgBouncerVersion
}

func (p PgBouncerVersion) ResourcePlural() string {
	return ResourcePluralPgBouncerVersion
}

func (p PgBouncerVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralPgBouncerVersion,
		Singular:      ResourceSingularPgBouncerVersion,
		Kind:          ResourceKindPgBouncerVersion,
		ShortNames:    []string{ResourceCodePgBouncerVersion},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.PgBouncerVersion",
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

func (p PgBouncerVersion) ValidateSpecs() error {
	if p.Spec.Version == "" ||
		p.Spec.Exporter.Image == "" ||
		p.Spec.Server.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for pgbouncerversion "%v":
spec.version,
spec.server.image,
spec.exporter.image.`, p.Name)
	}
	return nil
}
