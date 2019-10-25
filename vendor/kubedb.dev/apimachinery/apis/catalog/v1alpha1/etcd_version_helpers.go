package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
)

var _ apis.ResourceInfo = &EtcdVersion{}

func (e EtcdVersion) ResourceShortCode() string {
	return ResourceCodeEtcdVersion
}

func (e EtcdVersion) ResourceKind() string {
	return ResourceKindEtcdVersion
}

func (e EtcdVersion) ResourceSingular() string {
	return ResourceSingularEtcdVersion
}

func (e EtcdVersion) ResourcePlural() string {
	return ResourcePluralEtcdVersion
}

func (e EtcdVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralEtcdVersion,
		Singular:      ResourceSingularEtcdVersion,
		Kind:          ResourceKindEtcdVersion,
		ShortNames:    []string{ResourceCodeEtcdVersion},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.EtcdVersion",
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

func (e EtcdVersion) ValidateSpecs() error {
	if e.Spec.Version == "" ||
		e.Spec.DB.Image == "" ||
		e.Spec.Tools.Image == "" ||
		e.Spec.Exporter.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for etcdVersion "%v":
spec.version,
spec.db.image,
spec.tools.image,
spec.exporter.image.`, e.Name)
	}
	return nil
}
