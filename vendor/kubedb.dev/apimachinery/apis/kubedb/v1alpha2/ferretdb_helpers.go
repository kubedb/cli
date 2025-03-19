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

package v1alpha2

import (
	"context"
	"fmt"
	"net/url"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/Masterminds/semver/v3"
	"github.com/fatih/structs"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (f *FerretDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralFerretDB))
}

func (f *FerretDB) ResourcePlural() string {
	return ResourcePluralFerretDB
}

func (f *FerretDB) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", f.ResourcePlural(), "kubedb.com")
}

func (f *FerretDB) ServiceName() string {
	return f.Name
}

type FerretDBApp struct {
	*FerretDB
}

func (f FerretDBApp) Name() string {
	return f.FerretDB.Name
}

func (f FerretDBApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularFerretDB))
}

func (f *FerretDB) AppBindingMeta() appcat.AppBindingMeta {
	return &FerretDBApp{f}
}

func (f *FerretDB) OffshootName() string {
	return f.Name
}

func (f *FerretDB) OffshootSelectors() map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      f.ResourceFQN(),
		meta_util.InstanceLabelKey:  f.Name,
		meta_util.ManagedByLabelKey: "kubedb.com",
	}
	return selector
}

func (f *FerretDB) PodControllerLabels(podControllerLabels map[string]string, extraLabels ...map[string]string) map[string]string {
	return f.offshootLabels(meta_util.OverwriteKeys(f.OffshootSelectors(), extraLabels...), podControllerLabels)
}

func (f *FerretDB) OffshootLabels() map[string]string {
	return f.offshootLabels(f.OffshootSelectors(), nil)
}

func (f *FerretDB) SecondaryServerName() string {
	if f.Spec.Server.Secondary == nil {
		return ""
	}
	return fmt.Sprintf("%s-%s", f.OffshootName(), kubedb.FerretDBServerTypeSecondary)
}

func (f *FerretDB) PrimaryServerSelectors() map[string]string {
	return meta_util.OverwriteKeys(f.OffshootSelectors(), map[string]string{
		kubedb.FerretDBPrimaryLabelKey: f.OffshootName(),
	})
}

func (f *FerretDB) PrimaryServerLabels() map[string]string {
	return meta_util.OverwriteKeys(f.OffshootLabels(), f.PrimaryServerSelectors())
}

func (f *FerretDB) SecondaryServerSelectors() map[string]string {
	return meta_util.OverwriteKeys(f.OffshootSelectors(), map[string]string{
		kubedb.FerretDBSecondaryLabelKey: f.SecondaryServerName(),
	})
}

func (f *FerretDB) SecondaryServerLabels() map[string]string {
	return meta_util.OverwriteKeys(f.OffshootLabels(), f.SecondaryServerSelectors())
}

func (f *FerretDB) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys("kubedb.com", selector, meta_util.OverwriteKeys(nil, f.Labels, override))
}

func (f *FerretDB) GetAuthSecretName() string {
	if f.Spec.AuthSecret != nil && f.Spec.AuthSecret.Name != "" {
		return f.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(f.OffshootName(), "auth")
}

func (f *FerretDB) GetPersistentSecrets() []string {
	var secrets []string
	if f.Spec.AuthSecret != nil {
		secrets = append(secrets, f.GetAuthSecretName())
	}
	return secrets
}

// AsOwner returns owner reference to resources
func (f *FerretDB) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(f, SchemeGroupVersion.WithKind(f.ResourceKind()))
}

func (f *FerretDB) ResourceKind() string {
	return ResourceKindFerretDB
}

func (f *FerretDB) PgBackendName() string {
	return f.OffshootName() + "-pg-backend"
}

func (f *FerretDB) PodLabels(podTemplateLabels map[string]string, extraLabels ...map[string]string) map[string]string {
	return f.offshootLabels(meta_util.OverwriteKeys(f.OffshootSelectors(), extraLabels...), podTemplateLabels)
}

func (f *FerretDB) CertificateName(alias FerretDBCertificateAlias) string {
	return meta_util.NameWithSuffix(f.Name, fmt.Sprintf("%s-cert", string(alias)))
}

func (f *FerretDB) GetCertSecretName(alias FerretDBCertificateAlias) string {
	name, ok := kmapi.GetCertificateSecretName(f.Spec.TLS.Certificates, string(alias))
	if ok {
		return name
	}

	return f.CertificateName(alias)
}

func (f *FerretDB) GetExternalBackendClientSecretName() string {
	return f.Name + "-ext-pg-client-cert"
}

func (f *FerretDB) GetBackendConnectionSecretName() string {
	return f.OffshootName() + "-backend-connection"
}

func (f *FerretDB) GetSecretVolumeName(secretName string) string {
	return secretName + "-vol"
}

func (f *FerretDB) SetHealthCheckerDefaults() {
	if f.Spec.HealthChecker.PeriodSeconds == nil {
		f.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if f.Spec.HealthChecker.TimeoutSeconds == nil {
		f.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if f.Spec.HealthChecker.FailureThreshold == nil {
		f.Spec.HealthChecker.FailureThreshold = pointer.Int32P(2)
	}
}

func (f *FerretDB) SetDefaults(kc client.Client) {
	if f == nil {
		return
	}
	if f.Spec.StorageType == "" {
		f.Spec.StorageType = StorageTypeDurable
	}

	if f.Spec.DeletionPolicy == "" {
		f.Spec.DeletionPolicy = DeletionPolicyWipeOut
	}

	if f.Spec.SSLMode == "" {
		f.Spec.SSLMode = SSLModeDisabled
	}

	var frVersion catalog.FerretDBVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name: f.Spec.Version,
	}, &frVersion)
	if err != nil {
		klog.Errorf("can't get the FerretDB version object %s for %s \n", err.Error(), f.Spec.Version)
		return
	}

	if f.Spec.Server == nil {
		f.Spec.Server = &FerretDBServer{
			Primary: &FerretDBServerSpec{
				Replicas:    pointer.Int32P(1),
				PodTemplate: &ofst.PodTemplateSpec{},
			},
		}
	}

	if f.Spec.Server.Primary != nil {
		if f.Spec.Server.Primary.Replicas == nil {
			f.Spec.Server.Primary.Replicas = pointer.Int32P(1)
		}
		if f.Spec.Server.Primary.PodTemplate == nil {
			f.Spec.Server.Primary.PodTemplate = &ofst.PodTemplateSpec{}
		}
		f.setDefaultPodTemplateValues(f.Spec.Server.Primary.PodTemplate, &frVersion)
	}

	if f.Spec.Server.Secondary != nil {
		if f.Spec.Server.Secondary.Replicas == nil {
			f.Spec.Server.Secondary.Replicas = pointer.Int32P(1)
		}
		if f.Spec.Server.Secondary.PodTemplate == nil {
			f.Spec.Server.Secondary.PodTemplate = &ofst.PodTemplateSpec{}
		}
		f.setDefaultPodTemplateValues(f.Spec.Server.Secondary.PodTemplate, &frVersion)
	}

	if f.Spec.Backend.LinkedDB == "" {
		if f.Spec.Backend.ExternallyManaged {
			f.Spec.Backend.LinkedDB = "postgres"
		} else {
			f.Spec.Backend.LinkedDB = "ferretdb"
		}
	}

	if f.Spec.AuthSecret == nil {
		f.Spec.AuthSecret = &SecretReference{
			ExternallyManaged: false,
		}
	}

	f.Spec.Monitor.SetDefaults()
	if f.Spec.Monitor != nil && f.Spec.Monitor.Prometheus != nil {
		if f.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			f.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = frVersion.Spec.SecurityContext.RunAsUser
		}
		if f.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			f.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = frVersion.Spec.SecurityContext.RunAsUser
		}
	}
	defaultVersion := "16.4-bookworm"
	if f.IsLaterVersion(&frVersion, 2) {
		defaultVersion = "16.7-doc"
	}
	if !f.Spec.Backend.ExternallyManaged {
		if f.Spec.Backend.Version == nil {
			f.Spec.Backend.Version = &defaultVersion
		}
	}

	if f.Spec.Backend.PostgresRef != nil && f.Spec.Backend.PostgresRef.Name != "" && f.Spec.Backend.PostgresRef.Namespace == "" {
		f.Spec.Backend.PostgresRef.Namespace = f.Namespace
	}

	f.SetTLSDefaults()
	f.SetHealthCheckerDefaults()
}

func (f *FerretDB) setDefaultPodTemplateValues(podTemplate *ofst.PodTemplateSpec, frVersion *catalog.FerretDBVersion) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.FerretDBContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: kubedb.FerretDBContainerName,
		}
		podTemplate.Spec.Containers = append(podTemplate.Spec.Containers, *dbContainer)
	}
	if structs.IsZero(dbContainer.Resources) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}
	if dbContainer.SecurityContext == nil {
		dbContainer.SecurityContext = &core.SecurityContext{}
	}
	f.setDefaultContainerSecurityContext(frVersion, dbContainer.SecurityContext)
	f.setDefaultPodTemplateSecurityContext(frVersion, podTemplate)
}

func (f *FerretDB) IsLaterVersion(frVersion *catalog.FerretDBVersion, version uint64) bool {
	v, _ := semver.NewVersion(frVersion.Spec.Version)
	return v.Major() >= version
}

func (f *FerretDB) setDefaultPodTemplateSecurityContext(frVersion *catalog.FerretDBVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = frVersion.Spec.SecurityContext.RunAsUser
	}
}

func (f *FerretDB) setDefaultContainerSecurityContext(frVersion *catalog.FerretDBVersion, sc *core.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(true)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = frVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = frVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (f *FerretDB) SetTLSDefaults() {
	if f.Spec.TLS == nil || f.Spec.TLS.IssuerRef == nil {
		return
	}

	defaultServerOrg := []string{kubedb.KubeDBOrganization}
	defaultServerOrgUnit := []string{string(FerretDBServerCert)}

	_, cert := kmapi.GetCertificate(f.Spec.TLS.Certificates, string(FerretDBServerCert))
	if cert != nil && cert.Subject != nil {
		if cert.Subject.Organizations != nil {
			defaultServerOrg = cert.Subject.Organizations
		}
		if cert.Subject.OrganizationalUnits != nil {
			defaultServerOrgUnit = cert.Subject.OrganizationalUnits
		}
	}
	f.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(f.Spec.TLS.Certificates, kmapi.CertificateSpec{
		Alias:      string(FerretDBServerCert),
		SecretName: f.GetCertSecretName(FerretDBServerCert),
		Subject: &kmapi.X509Subject{
			Organizations:       defaultServerOrg,
			OrganizationalUnits: defaultServerOrgUnit,
		},
	})

	// Client-cert
	defaultClientOrg := []string{kubedb.KubeDBOrganization}
	defaultClientOrgUnit := []string{string(FerretDBClientCert)}
	_, cert = kmapi.GetCertificate(f.Spec.TLS.Certificates, string(FerretDBClientCert))
	if cert != nil && cert.Subject != nil {
		if cert.Subject.Organizations != nil {
			defaultClientOrg = cert.Subject.Organizations
		}
		if cert.Subject.OrganizationalUnits != nil {
			defaultClientOrgUnit = cert.Subject.OrganizationalUnits
		}
	}
	f.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(f.Spec.TLS.Certificates, kmapi.CertificateSpec{
		Alias:      string(FerretDBClientCert),
		SecretName: f.GetCertSecretName(FerretDBClientCert),
		Subject: &kmapi.X509Subject{
			Organizations:       defaultClientOrg,
			OrganizationalUnits: defaultClientOrgUnit,
		},
	})
}

type FerretDBStatsService struct {
	*FerretDB
}

func (fs FerretDBStatsService) ServiceMonitorName() string {
	return fs.OffshootName() + "-stats"
}

func (fs FerretDBStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return fs.OffshootLabels()
}

func (fs FerretDBStatsService) Path() string {
	return kubedb.FerretDBMetricsPath
}

func (fs FerretDBStatsService) Scheme() string {
	return ""
}

func (fs FerretDBStatsService) TLSConfig() *v1.TLSConfig {
	return nil
}

func (fs FerretDBStatsService) ServiceName() string {
	return fs.OffshootName() + "-stats"
}

func (f *FerretDB) StatsService() mona.StatsAccessor {
	return &FerretDBStatsService{f}
}

func (f *FerretDB) GoverningServiceName() string {
	return f.OffshootName() + "-pods"
}

func (f *FerretDB) SecondaryGoverningServiceName() string {
	return f.SecondaryServerName() + "-pods"
}

func (f *FerretDB) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(f.Spec.ServiceTemplates, alias)
	return f.offshootLabels(meta_util.OverwriteKeys(f.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (f *FerretDB) StatsServiceLabels() map[string]string {
	return f.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (f *FerretDB) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicasOfPetSet(lister.PetSets(f.Namespace), labels.SelectorFromSet(f.OffshootLabels()), expectedItems)
}

func (f *FerretDB) GetSSLModeFromAppBinding(apb *appcat.AppBinding) (PostgresSSLMode, error) {
	var sslMode string
	if apb.Spec.ClientConfig.URL != nil {
		parsedURL, err := url.Parse(*apb.Spec.ClientConfig.URL)
		if err != nil {
			return "", fmt.Errorf("parse error in appbinding %s/%s 'spec.clientConfig.url'. error: %v", apb.Namespace, apb.Name, err)
		}
		if parsedURL.Scheme != "postgres" && parsedURL.Scheme != "postgresql" {
			return "", fmt.Errorf("invalid scheme provided in URL in provided appbinding %s/%s", apb.Namespace, apb.Name)
		}
		sslMode = parsedURL.Query().Get("sslmode")
	}
	if apb.Spec.ClientConfig.Service != nil {
		values, err := url.ParseQuery(apb.Spec.ClientConfig.Service.Query)
		if err != nil {
			return "", fmt.Errorf("parse error in appbinding %s/%s 'spec.clientConfig.service.query'. error: %v", apb.Namespace, apb.Name, err)
		}
		if sslMode != "" && sslMode != values.Get("sslmode") {
			return "", fmt.Errorf("sslMode is not same in 'spec.clientConfig.service.query' and 'spec.clientConfig.url' of appbinding %s/%s", apb.Namespace, apb.Name)
		}
		sslMode = values.Get("sslmode")
	}
	// If sslMode is not specified anywhere, it will be disabled
	if sslMode == "" {
		sslMode = "disable"
	}
	return PostgresSSLMode(sslMode), nil
}
