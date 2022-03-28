/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"errors"
	"fmt"
	"path/filepath"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"kubedb.dev/apimachinery/crds"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
)

func (_ ElasticsearchDashboard) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceElasticsearchDashboards))
}

var _ apis.ResourceInfo = &ElasticsearchDashboard{}

func (ed ElasticsearchDashboard) OffshootName() string {
	return ed.Name
}

func (ed ElasticsearchDashboard) ServiceName() string {
	return ed.OffshootName()
}

func (ed ElasticsearchDashboard) DeploymentName() string {
	return ed.OffshootName()
}

func (ed ElasticsearchDashboard) DashboardContainerName() string {
	return meta_util.NameWithSuffix(ed.Name, "dashboard")
}

func (ed ElasticsearchDashboard) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourceElasticsearchDashboards, kubedb.GroupName)
}

func (ed ElasticsearchDashboard) ResourceShortCode() string {
	return ResourceCodeElasticsearchDashboard
}

func (ed ElasticsearchDashboard) ResourceKind() string {
	return ResourceKindElasticsearchDashboard
}

func (ed ElasticsearchDashboard) ResourceSingular() string {
	return ResourceElasticsearchDashboard
}

func (ed ElasticsearchDashboard) ResourcePlural() string {
	return ResourceElasticsearchDashboards
}

// DefaultCertificateSecretName returns the default certificate name and/or certificate secret name for a certificate alias
func (ed *ElasticsearchDashboard) DefaultCertificateSecretName(alias ElasticsearchDashboardCertificateAlias) string {
	return meta_util.NameWithSuffix(ed.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// Owner returns owner reference to resources
func (ed *ElasticsearchDashboard) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(ed, SchemeGroupVersion.WithKind(ResourceKindElasticsearchDashboard))
}

// returns the volume name for config Secret.

func (ed ElasticsearchDashboard) GetCertVolumeName(alias ElasticsearchDashboardCertificateAlias) string {
	return meta_util.NameWithSuffix(string(alias), "volume")
}

func (ed ElasticsearchDashboard) AuthSecretName() string {
	if ed.Spec.AuthSecret != nil {
		return ed.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(ed.Name, "database-cred")
}

func (ed ElasticsearchDashboard) GetSecretName(alias ElasticsearchDashboardCertificateAlias) string {
	return meta_util.NameWithSuffix(ed.Name, string(alias))
}

func (ed ElasticsearchDashboard) DatabaseClientSecretName() string {
	return meta_util.NameWithSuffix(ed.Name, "database-client")
}

func (ed ElasticsearchDashboard) ClientCertificateCN(alias ElasticsearchDashboardCertificateAlias) string {
	return fmt.Sprintf("%s-%s", ed.Name, string(alias))
}

func (ed *ElasticsearchDashboard) GetDatabaseClientCertName(databaseName string) string {
	return fmt.Sprintf("%s-%s", databaseName, ed.CertificateSecretName(DefaultElasticsearchClientCertAlias))
}

func (ed *ElasticsearchDashboard) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      ed.ResourceFQN(),
		meta_util.InstanceLabelKey:  ed.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

// Selectors returns a labels elector by offshoring extra selector (if any)
func (ed *ElasticsearchDashboard) Selectors() *meta.LabelSelector {
	extraLabels := map[string]string{
		meta_util.InstanceLabelKey: ed.Name,
	}
	return &meta.LabelSelector{
		MatchLabels: ed.OffshootSelectors(extraLabels),
	}
}

func (ed *ElasticsearchDashboard) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentDashboard
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, ed.getLabels(), override))
}

func (ed *ElasticsearchDashboard) OffshootLabels() map[string]string {
	return ed.offshootLabels(ed.OffshootSelectors(), nil)
}

func (ed *ElasticsearchDashboard) getLabels(extraLabels ...map[string]string) map[string]string {
	return meta_util.OverwriteKeys(ed.OffshootSelectors(), extraLabels...)
}

func (ed *ElasticsearchDashboard) PodLabels(extraLabels ...map[string]string) map[string]string {
	return meta_util.OverwriteKeys(ed.OffshootSelectors(), extraLabels...)
}

func (ed *ElasticsearchDashboard) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return ed.offshootLabels(meta_util.OverwriteKeys(ed.OffshootSelectors(), extraLabels...), ed.Spec.PodTemplate.Controller.Labels)
}

func (ed *ElasticsearchDashboard) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	return meta_util.OverwriteKeys(ed.OffshootSelectors(), extraLabels...)
}

func (ed *ElasticsearchDashboard) GetServiceSelectors() map[string]string {
	extraSelectors := map[string]string{
		"app.kubernetes.io/instance": ed.Name,
	}
	return ed.OffshootSelectors(extraSelectors)
}

// returns the mountPath for certificate secrets.
// if configDir is "/usr/share/kibana/config",
// mountPath will be, "/usr/share/kibana/config/certs/<alias>/filename".

func (ed *ElasticsearchDashboard) CertSecretVolumeMountPath(configDir string, alias ElasticsearchDashboardCertificateAlias) string {
	return filepath.Join(configDir, "certs", string(alias))
}

// returns a certificate file path  for a specific file using the certificate alias

func (ed *ElasticsearchDashboard) CertificateFilePath(configDir string, alias ElasticsearchDashboardCertificateAlias, filename string) string {
	return filepath.Join(ed.CertSecretVolumeMountPath(configDir, alias), filename)
}

func (ed *ElasticsearchDashboard) GetServicePort(alias ServiceAlias) int32 {
	reqAlias := v1alpha2.ServiceAlias(alias)
	svcTemplate := v1alpha2.GetServiceTemplate(ed.Spec.ServiceTemplates, reqAlias)
	return svcTemplate.Spec.Ports[0].Port
}

func (ed *ElasticsearchDashboard) DatabaseConnectionURL(servicePort int32) (string, error) {
	if ed.Spec.DatabaseRef != nil {
		if ed.Spec.DatabaseRef.Name == "" {
			return "", errors.New("required database fields not found")
		}
		return fmt.Sprintf("%s://%s.%s.svc:%d", ed.GetConnectionScheme(), ed.Spec.DatabaseRef.Name, ed.Namespace, servicePort), nil
	}
	return fmt.Sprintf("%s://%s.%s.svc:%d", ed.GetConnectionScheme(), ed.Spec.DatabaseRef.Name, ed.Namespace, servicePort), nil
}

func (ed *ElasticsearchDashboard) GetConnectionScheme() string {
	scheme := "http"
	if ed.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}

// CertificateSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.

func (ed *ElasticsearchDashboard) CertificateSecretName(alias ElasticsearchDashboardCertificateAlias) string {
	if ed.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(ed.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return ed.DefaultCertificateSecretName(alias)
}

func (ed ElasticsearchDashboard) DefaultConfigSecretName() string {
	return meta_util.NameWithSuffix(ed.Name, "config")
}

func (ed ElasticsearchDashboard) CustomConfigSecretName() string {
	return ed.Spec.ConfigSecret.Name
}

func (ed *ElasticsearchDashboard) CertSecretExists(alias ElasticsearchDashboardCertificateAlias) bool {
	if ed.Spec.TLS != nil {
		_, ok := kmapi.GetCertificateSecretName(ed.Spec.TLS.Certificates, string(alias))
		if ok {
			return true
		}
	}
	return false
}
