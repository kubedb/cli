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
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/fatih/structs"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
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

func (f *FerretDB) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys("kubedb.com", selector, meta_util.OverwriteKeys(nil, f.Labels, override))
}

func (f *FerretDB) GetAuthSecretName() string {
	if f.Spec.AuthSecret != nil && f.Spec.AuthSecret.Name != "" {
		return f.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(f.PgBackendName(), "auth")
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

func (f *FerretDB) SetDefaults() {
	if f == nil {
		return
	}
	if f.Spec.StorageType == "" {
		f.Spec.StorageType = StorageTypeDurable
	}

	if f.Spec.TerminationPolicy == "" {
		f.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if f.Spec.SSLMode == "" {
		f.Spec.SSLMode = SSLModeDisabled
	}

	if f.Spec.Replicas == nil {
		f.Spec.Replicas = pointer.Int32P(1)
	}

	if f.Spec.PodTemplate == nil {
		f.Spec.PodTemplate = &ofst.PodTemplateSpec{}
	}

	dbContainer := coreutil.GetContainerByName(f.Spec.PodTemplate.Spec.Containers, FerretDBContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: FerretDBContainerName,
		}
		f.Spec.PodTemplate.Spec.Containers = append(f.Spec.PodTemplate.Spec.Containers, *dbContainer)
	}
	if structs.IsZero(dbContainer.Resources) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, DefaultResources)
	}

	if f.Spec.Backend.LinkedDB == "" {
		if f.Spec.Backend.ExternallyManaged {
			f.Spec.Backend.LinkedDB = "postgres"
		} else {
			f.Spec.Backend.LinkedDB = "ferretdb"
		}
	}
	if f.Spec.Monitor != nil && f.Spec.Monitor.Prometheus.Exporter.Port == 0 {
		// 56790 is default port for Prometheus operator.
		f.Spec.Monitor.Prometheus.Exporter.Port = 56790
	}
	defaultVersion := "13.13"
	if !f.Spec.Backend.ExternallyManaged && f.Spec.Backend.Postgres == nil {
		f.Spec.Backend.Postgres = &PostgresRef{
			Version: &defaultVersion,
		}
	}
	f.SetTLSDefaults()
	f.SetHealthCheckerDefaults()
}

func (f *FerretDB) SetTLSDefaults() {
	if f.Spec.TLS == nil || f.Spec.TLS.IssuerRef == nil {
		return
	}

	defaultServerOrg := []string{KubeDBOrganization}
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
	defaultClientOrg := []string{KubeDBOrganization}
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
	return FerretDBMetricsPath
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

func (f *FerretDB) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(f.Spec.ServiceTemplates, alias)
	return f.offshootLabels(meta_util.OverwriteKeys(f.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (f *FerretDB) StatsServiceLabels() map[string]string {
	return f.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}
