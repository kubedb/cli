package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
)

var _ apis.ResourceInfo = &RedisVersion{}

func (r RedisVersion) ResourceShortCode() string {
	return ResourceCodeRedisVersion
}

func (r RedisVersion) ResourceKind() string {
	return ResourceKindRedisVersion
}

func (r RedisVersion) ResourceSingular() string {
	return ResourceSingularRedisVersion
}

func (r RedisVersion) ResourcePlural() string {
	return ResourcePluralRedisVersion
}

func (r RedisVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralRedisVersion,
		Singular:      ResourceSingularRedisVersion,
		Kind:          ResourceKindRedisVersion,
		ShortNames:    []string{ResourceCodeRedisVersion},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.RedisVersion",
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

func (r RedisVersion) ValidateSpecs() error {
	if r.Spec.Version == "" ||
		r.Spec.DB.Image == "" ||
		r.Spec.Exporter.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for redisVersion "%v":
spec.version,
spec.db.image,
spec.exporter.image.`, r.Name)
	}
	return nil
}
