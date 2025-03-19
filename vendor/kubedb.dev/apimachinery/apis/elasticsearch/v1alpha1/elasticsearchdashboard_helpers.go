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
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/apimachinery/crds"

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	return meta_util.NameWithSuffix(ed.Name, "elasticsearch")
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
	reqAlias := dbapi.ServiceAlias(alias)
	svcTemplate := dbapi.GetServiceTemplate(ed.Spec.ServiceTemplates, reqAlias)
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

func (ed *ElasticsearchDashboard) SetHealthCheckerDefaults() {
	if ed.Spec.HealthChecker.PeriodSeconds == nil {
		ed.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(20)
	}
	if ed.Spec.HealthChecker.TimeoutSeconds == nil {
		ed.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if ed.Spec.HealthChecker.FailureThreshold == nil {
		ed.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (ed *ElasticsearchDashboard) SetDefaults(kc client.Client) {
	if ed.Spec.DeletionPolicy == "" {
		ed.Spec.DeletionPolicy = dbapi.DeletionPolicyDelete
	}

	db := dbapi.Elasticsearch{}
	esVersion := catalog.ElasticsearchVersion{}
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name:      ed.Spec.DatabaseRef.Name,
		Namespace: ed.Namespace,
	}, &db)
	if err != nil {
		klog.Errorf("can't get the elasticsearch: %v\n", err.Error())
		return
	}

	err = kc.Get(context.TODO(), types.NamespacedName{
		Name: db.Spec.Version,
	}, &esVersion)
	if err != nil {
		klog.Errorf("can't get the elasticsearch version: %v\n", err.Error())
		return
	}

	if ed.Spec.Replicas == nil {
		ed.Spec.Replicas = pointer.Int32P(1)
	}
	if ed.Spec.PodTemplate.Spec.SecurityContext == nil {
		ed.Spec.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{}
	}
	ed.Spec.PodTemplate.Spec.SecurityContext.FSGroup = esVersion.Spec.SecurityContext.RunAsUser
	ed.setDefaultContainerSecurityContext(esVersion, &ed.Spec.PodTemplate)
	ed.setDefaultContainerResourceLimits(&ed.Spec.PodTemplate)

	if ed.Spec.EnableSSL {
		if ed.Spec.TLS == nil {
			ed.Spec.TLS = &kmapi.TLSConfig{}
		}
		if ed.Spec.TLS.IssuerRef == nil {
			ed.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(ed.Spec.TLS.Certificates, kmapi.CertificateSpec{
				Alias:      string(ElasticsearchDashboardCACert),
				SecretName: ed.DefaultCertificateSecretName(ElasticsearchDashboardCACert),
			})
		}
		ed.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(ed.Spec.TLS.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchDashboardServerCert),
			SecretName: ed.DefaultCertificateSecretName(ElasticsearchDashboardServerCert),
		})
	}
}

func (ed *ElasticsearchDashboard) setDefaultContainerSecurityContext(esVersion catalog.ElasticsearchVersion, podTemplate *ofst.PodTemplateSpec) {
	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.ElasticsearchInitConfigMergerContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.ElasticsearchInitConfigMergerContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	ed.assignDefaultContainerSecurityContext(esVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)
	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.ElasticsearchContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	ed.assignDefaultContainerSecurityContext(esVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (ed *ElasticsearchDashboard) setDefaultContainerResourceLimits(podTemplate *ofst.PodTemplateSpec) {
	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
	if container != nil && (container.Resources.Requests == nil && container.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&container.Resources, kubedb.DefaultResources)
	}

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.ElasticsearchInitConfigMergerContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}
}

func (ed *ElasticsearchDashboard) assignDefaultContainerSecurityContext(esVersion catalog.ElasticsearchVersion, sc *core.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(esVersion.Spec.SecurityContext.RunAsAnyNonRoot)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = esVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}
