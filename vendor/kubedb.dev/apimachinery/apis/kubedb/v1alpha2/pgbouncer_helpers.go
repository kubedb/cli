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

	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (p PgBouncer) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgBouncer))
}

var _ apis.ResourceInfo = &PgBouncer{}

func (p PgBouncer) OffshootName() string {
	return p.Name
}

func (p PgBouncer) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      p.ResourceFQN(),
		meta_util.InstanceLabelKey:  p.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (p PgBouncer) OffshootLabels() map[string]string {
	out := p.OffshootSelectors()
	out[meta_util.ComponentLabelKey] = "connection-pooler"
	return meta_util.FilterKeys(kubedb.GroupName, out, p.Labels)
}

func (p PgBouncer) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPgBouncer, kubedb.GroupName)
}

func (p PgBouncer) ResourceShortCode() string {
	return ResourceCodePgBouncer
}

func (p PgBouncer) ResourceKind() string {
	return ResourceKindPgBouncer
}

func (p PgBouncer) ResourceSingular() string {
	return ResourceSingularPgBouncer
}

func (p PgBouncer) ResourcePlural() string {
	return ResourcePluralPgBouncer
}

func (p PgBouncer) ServiceName() string {
	return p.OffshootName()
}

func (p PgBouncer) GoverningServiceName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "pods")
}

func (p PgBouncer) AuthSecretName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "auth")
}

func (p PgBouncer) ConfigSecretName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "config")
}

type pgbouncerApp struct {
	*PgBouncer
}

func (r pgbouncerApp) Name() string {
	return r.PgBouncer.Name
}

func (r pgbouncerApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPgBouncer))
}

func (p PgBouncer) AppBindingMeta() appcat.AppBindingMeta {
	return &pgbouncerApp{&p}
}

type pgbouncerStatsService struct {
	*PgBouncer
}

func (p pgbouncerStatsService) GetNamespace() string {
	return p.PgBouncer.GetNamespace()
}

func (p pgbouncerStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p pgbouncerStatsService) ServiceMonitorName() string {
	return p.ServiceName()
}

func (p pgbouncerStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return p.OffshootLabels()
}

func (p pgbouncerStatsService) Path() string {
	return DefaultStatsPath
}

func (p pgbouncerStatsService) Scheme() string {
	return ""
}

func (p PgBouncer) StatsService() mona.StatsAccessor {
	return &pgbouncerStatsService{&p}
}

func (p PgBouncer) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, p.OffshootSelectors(), p.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (p PgBouncer) ReplicasServiceName() string {
	return fmt.Sprintf("%v-replicas", p.Name)
}

func (p *PgBouncer) SetDefaults() {
	if p == nil {
		return
	}
	p.Spec.Monitor.SetDefaults()

	p.SetTLSDefaults()
	SetDefaultResourceLimits(&p.Spec.PodTemplate.Spec.Container.Resources, DefaultResourceLimits)
}

func (p *PgBouncer) SetTLSDefaults() {
	if p.Spec.TLS == nil || p.Spec.TLS.IssuerRef == nil {
		return
	}

	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PgBouncerServerCert), p.CertificateName(PgBouncerServerCert))
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PgBouncerClientCert), p.CertificateName(PgBouncerClientCert))
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PgBouncerMetricsExporterCert), p.CertificateName(PgBouncerMetricsExporterCert))
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (p *PgBouncer) CertificateName(alias PgBouncerCertificateAlias) string {
	return meta_util.NameWithSuffix(p.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// MustCertSecretName returns the secret name for a certificate alias
func (p *PgBouncer) MustCertSecretName(alias PgBouncerCertificateAlias) string {
	if p == nil {
		panic("missing PgBouncer database")
	} else if p.Spec.TLS == nil {
		panic(fmt.Errorf("PgBouncer %s/%s is missing tls spec", p.Namespace, p.Name))
	}
	name, ok := kmapi.GetCertificateSecretName(p.Spec.TLS.Certificates, string(alias))
	if !ok {
		panic(fmt.Errorf("PgBouncer %s/%s is missing secret name for %s certificate", p.Namespace, p.Name, alias))
	}
	return name
}

func (p *PgBouncer) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(p.Namespace), labels.SelectorFromSet(p.OffshootLabels()), expectedItems)
}
