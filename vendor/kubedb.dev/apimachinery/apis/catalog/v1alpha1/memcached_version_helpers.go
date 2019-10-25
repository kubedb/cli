package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
)

var _ apis.ResourceInfo = &MemcachedVersion{}

func (m MemcachedVersion) ResourceShortCode() string {
	return ResourceCodeMemcachedVersion
}

func (m MemcachedVersion) ResourceKind() string {
	return ResourceKindMemcachedVersion
}

func (m MemcachedVersion) ResourceSingular() string {
	return ResourceSingularMemcachedVersion
}

func (m MemcachedVersion) ResourcePlural() string {
	return ResourcePluralMemcachedVersion
}

func (m MemcachedVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralMemcachedVersion,
		Singular:      ResourceSingularMemcachedVersion,
		Kind:          ResourceKindMemcachedVersion,
		ShortNames:    []string{ResourceCodeMemcachedVersion},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.MemcachedVersion",
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

func (m MemcachedVersion) ValidateSpecs() error {
	if m.Spec.Version == "" ||
		m.Spec.DB.Image == "" ||
		m.Spec.Exporter.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for memcachedVersion "%v":
spec.version,
spec.db.image,
spec.exporter.image,`, m.Name)
	}
	return nil
}
