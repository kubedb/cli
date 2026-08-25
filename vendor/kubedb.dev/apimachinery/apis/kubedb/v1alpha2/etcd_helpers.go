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
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"
	apiutils "kubedb.dev/apimachinery/pkg/utils"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	metautil "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
)

func (Etcd) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralEtcd))
}

var _ apis.ResourceInfo = &Etcd{}

func (e *Etcd) ResourceKind() string {
	return ResourceKindEtcd
}

func (e *Etcd) ResourceSingular() string {
	return ResourceSingularEtcd
}

func (e *Etcd) ResourcePlural() string {
	return ResourcePluralEtcd
}

func (e *Etcd) ResourceShortCode() string {
	return ResourceCodeEtcd
}

func (e *Etcd) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralEtcd, kubedb.GroupName)
}

func (e *Etcd) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(e, SchemeGroupVersion.WithKind(ResourceKindEtcd))
}

func (e *Etcd) GetNameSpacedName() string {
	return e.Namespace + "/" + e.Name
}

// ---------------------------------------------------------------------------
// Naming
// ---------------------------------------------------------------------------

func (e *Etcd) OffshootName() string {
	return e.Name
}

func (e *Etcd) PetSetName() string {
	return e.OffshootName()
}

func (e *Etcd) ServiceAccountName() string {
	return e.OffshootName()
}

// ServiceName is the client (load balanced) service used by applications.
func (e *Etcd) ServiceName() string {
	return e.OffshootName()
}

// GoverningServiceName is the headless service backing the PetSet; the per-member
// peer URLs are derived from it.
func (e *Etcd) GoverningServiceName() string {
	return metautil.NameWithSuffix(e.ServiceName(), "pods")
}

func (e *Etcd) ConfigSecretName() string {
	uid := string(e.UID)
	return metautil.NameWithSuffix(e.OffshootName(), uid[len(uid)-6:])
}

// PrimaryServiceDNS is the in-cluster DNS name clients should use. etcd has no
// primary/secondary split for reads or writes -- any member serves the client
// API and proxies to the Raft leader -- so this resolves to the client service.
func (e *Etcd) PrimaryServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", e.ServiceName(), e.Namespace)
}

// ServiceDNS is an alias of PrimaryServiceDNS kept for symmetry with the other
// KubeDB databases.
func (e *Etcd) ServiceDNS() string {
	return e.PrimaryServiceDNS()
}

// GoverningServiceDNS returns the fully qualified DNS name of a single member pod.
func (e *Etcd) GoverningServiceDNS(podName string) string {
	return fmt.Sprintf("%s.%s.%s.svc.%s", podName, e.GoverningServiceName(), e.Namespace, apiutils.FindDomain())
}

// PodName returns the name of the member pod at the given PetSet ordinal.
func (e *Etcd) PodName(ordinal int) string {
	return fmt.Sprintf("%s-%d", e.PetSetName(), ordinal)
}

// ClientURL is the etcd client endpoint of a single member.
func (e *Etcd) ClientURL(podName string) string {
	return fmt.Sprintf("%s://%s:%d", e.Scheme(), e.GoverningServiceDNS(podName), kubedb.EtcdClientPort)
}

// PeerURL is the etcd Raft peer endpoint of a single member.
func (e *Etcd) PeerURL(podName string) string {
	return fmt.Sprintf("%s://%s:%d", e.Scheme(), e.GoverningServiceDNS(podName), kubedb.EtcdPeerPort)
}

// Scheme is http unless TLS is configured.
func (e *Etcd) Scheme() string {
	if e.Spec.TLS != nil {
		return "https"
	}
	return "http"
}

// ---------------------------------------------------------------------------
// Labels & selectors
// ---------------------------------------------------------------------------

func (e *Etcd) offshootLabels(selector, override map[string]string) map[string]string {
	selector[metautil.ComponentLabelKey] = kubedb.ComponentDatabase
	return metautil.FilterKeys(SchemeGroupVersion.Group, selector, metautil.OverwriteKeys(nil, e.Labels, override))
}

func (e *Etcd) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		metautil.NameLabelKey:      e.ResourceFQN(),
		metautil.InstanceLabelKey:  e.Name,
		metautil.ManagedByLabelKey: SchemeGroupVersion.Group,
	}
	return metautil.OverwriteKeys(selector, extraSelectors...)
}

func (e *Etcd) OffshootLabels() map[string]string {
	return e.offshootLabels(e.OffshootSelectors(), nil)
}

func (e *Etcd) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(e.Spec.ServiceTemplates, alias)
	return e.offshootLabels(metautil.OverwriteKeys(e.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (e *Etcd) PodLabels(extraLabels ...map[string]string) map[string]string {
	return e.offshootLabels(metautil.OverwriteKeys(e.OffshootSelectors(), extraLabels...), e.Spec.PodTemplate.Labels)
}

func (e *Etcd) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return e.offshootLabels(metautil.OverwriteKeys(e.OffshootSelectors(), extraLabels...), e.Spec.PodTemplate.Controller.Labels)
}

func (e *Etcd) StatsServiceLabels() map[string]string {
	return e.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

// ---------------------------------------------------------------------------
// Secrets
// ---------------------------------------------------------------------------

func (e *Etcd) GetAuthSecretName() string {
	if e.Spec.AuthSecret != nil && e.Spec.AuthSecret.Name != "" {
		return e.Spec.AuthSecret.Name
	}
	return metautil.NameWithSuffix(e.OffshootName(), "auth")
}

func (e *Etcd) GetPersistentSecrets() []string {
	var secrets []string
	if !IsVirtualAuthSecretReferred(e.Spec.AuthSecret) && e.Spec.AuthSecret != nil && e.Spec.AuthSecret.Name != "" {
		secrets = append(secrets, e.GetAuthSecretName())
	}
	return secrets
}

func (e *Etcd) GetDeletionPolicy() string {
	return string(e.Spec.DeletionPolicy)
}

// ---------------------------------------------------------------------------
// Certificates
// ---------------------------------------------------------------------------

// CertificateName returns the default secret name holding the certificate of the
// given alias.
func (e *Etcd) CertificateName(alias EtcdCertificateAlias) string {
	return metautil.NameWithSuffix(e.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the referenced secret name for the certificate alias.
// If the alias exists in spec.tls.certificates with an empty secretName, the
// conventional "<db-name>-<alias>-cert" name is returned.
func (e *Etcd) GetCertSecretName(alias EtcdCertificateAlias) string {
	if e.Spec.TLS == nil {
		return ""
	}
	if _, cert := kmapi.GetCertificate(e.Spec.TLS.Certificates, string(alias)); cert != nil {
		if cert.SecretName != "" {
			return cert.SecretName
		}
		return e.CertificateName(alias)
	}
	return ""
}

// GetCertSecretVolumeName returns the pod volume name that carries the given
// certificate alias.
func (e *Etcd) GetCertSecretVolumeName(alias EtcdCertificateAlias) string {
	switch alias {
	case EtcdServerCert:
		return kubedb.EtcdServerTLSVolumeName
	case EtcdClientCert:
		return kubedb.EtcdClientTLSVolumeName
	case EtcdPeerCert:
		return kubedb.EtcdPeerTLSVolumeName
	case EtcdMetricsExporterCert:
		return kubedb.EtcdExporterTLSVolumeName
	}
	return ""
}

// GetCertSecretVolumeMountPath returns the mount path of the given certificate alias.
func (e *Etcd) GetCertSecretVolumeMountPath(alias EtcdCertificateAlias) string {
	switch alias {
	case EtcdServerCert:
		return kubedb.EtcdServerTLSMountPath
	case EtcdClientCert:
		return kubedb.EtcdClientTLSMountPath
	case EtcdPeerCert:
		return kubedb.EtcdPeerTLSMountPath
	case EtcdMetricsExporterCert:
		return kubedb.EtcdExporterTLSMountPath
	}
	return ""
}

// CertificateDNSNames are the SANs shared by the server and peer certificates.
// Both the client service and every member pod (via the governing service) must
// be covered, because etcd validates the URL it dialed.
func (e *Etcd) CertificateDNSNames() []string {
	names := []string{
		e.ServiceName(),
		fmt.Sprintf("%s.%s", e.ServiceName(), e.Namespace),
		fmt.Sprintf("%s.%s.svc", e.ServiceName(), e.Namespace),
		fmt.Sprintf("%s.%s.svc.%s", e.ServiceName(), e.Namespace, apiutils.FindDomain()),
		fmt.Sprintf("%s.%s", e.GoverningServiceName(), e.Namespace),
		fmt.Sprintf("%s.%s.svc", e.GoverningServiceName(), e.Namespace),
		fmt.Sprintf("%s.%s.svc.%s", e.GoverningServiceName(), e.Namespace, apiutils.FindDomain()),
		fmt.Sprintf("*.%s.%s.svc", e.GoverningServiceName(), e.Namespace),
		fmt.Sprintf("*.%s.%s.svc.%s", e.GoverningServiceName(), e.Namespace, apiutils.FindDomain()),
	}

	seen := map[string]struct{}{}
	out := make([]string, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}

// ---------------------------------------------------------------------------
// AppBinding & monitoring
// ---------------------------------------------------------------------------

type etcdApp struct {
	*Etcd
}

func (r etcdApp) Name() string {
	return r.Etcd.Name
}

func (r etcdApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", SchemeGroupVersion.Group, ResourceSingularEtcd))
}

func (e Etcd) AppBindingMeta() appcat.AppBindingMeta {
	return &etcdApp{&e}
}

type etcdStatsService struct {
	*Etcd
}

func (s etcdStatsService) GetNamespace() string {
	return s.Etcd.GetNamespace()
}

func (s etcdStatsService) ServiceName() string {
	return s.OffshootName() + "-stats"
}

func (s etcdStatsService) ServiceMonitorName() string {
	return s.ServiceName()
}

func (s etcdStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return s.OffshootLabels()
}

func (s etcdStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (s etcdStatsService) Scheme() string {
	sc := promapi.SchemeHTTP
	return sc.String()
}

func (s etcdStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (e *Etcd) StatsService() mona.StatsAccessor {
	return &etcdStatsService{e}
}

func (e *Etcd) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicasOfPetSet(lister.PetSets(e.Namespace), labels.SelectorFromSet(e.OffshootLabels()), expectedItems)
}

// ---------------------------------------------------------------------------
// Defaulting
// ---------------------------------------------------------------------------

func (e *Etcd) SetHealthCheckerDefaults() {
	if e.Spec.HealthChecker.PeriodSeconds == nil {
		e.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if e.Spec.HealthChecker.TimeoutSeconds == nil {
		e.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if e.Spec.HealthChecker.FailureThreshold == nil {
		e.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (e *Etcd) SetDefaults(etcdVersion *catalog.EtcdVersion) {
	if e == nil {
		return
	}

	// An odd member count keeps the Raft quorum optimal; three is the smallest
	// size that tolerates a single failure.
	if e.Spec.Replicas == nil {
		e.Spec.Replicas = pointer.Int32P(3)
	}
	if e.Spec.StorageType == "" {
		e.Spec.StorageType = StorageTypeDurable
	}
	if e.Spec.DeletionPolicy == "" {
		e.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if e.Spec.AuthSecret == nil {
		e.Spec.AuthSecret = &SecretReference{}
	}
	if e.Spec.AuthSecret.Kind == "" {
		e.Spec.AuthSecret.Kind = kubedb.ResourceKindSecret
	}

	if e.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		e.Spec.PodTemplate.Spec.ServiceAccountName = e.OffshootName()
	}

	e.SetTLSDefaults()
	e.SetHealthCheckerDefaults()
	e.setDefaultContainerSecurityContext(etcdVersion, &e.Spec.PodTemplate)
	e.setDefaultContainerResourceLimits(&e.Spec.PodTemplate)

	if e.Spec.Monitor != nil {
		if e.Spec.Monitor.Prometheus == nil {
			e.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if e.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			e.Spec.Monitor.Prometheus.Exporter.Port = kubedb.EtcdExporterPort
		}
		e.Spec.Monitor.SetDefaults()
		if etcdVersion != nil {
			if e.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
				e.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = etcdVersion.Spec.SecurityContext.RunAsUser
			}
			if e.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
				e.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = etcdVersion.Spec.SecurityContext.RunAsUser
			}
		}
	}
}

func (e *Etcd) SetTLSDefaults() {
	if e.Spec.TLS == nil || e.Spec.TLS.IssuerRef == nil {
		return
	}

	for _, alias := range []EtcdCertificateAlias{
		EtcdServerCert,
		EtcdClientCert,
		EtcdPeerCert,
		EtcdMetricsExporterCert,
	} {
		defaultOrg := []string{kubedb.KubeDBOrganization}
		defaultOrgUnit := []string{string(alias)}
		if _, cert := kmapi.GetCertificate(e.Spec.TLS.Certificates, string(alias)); cert != nil && cert.Subject != nil {
			if cert.Subject.Organizations != nil {
				defaultOrg = cert.Subject.Organizations
			}
			if cert.Subject.OrganizationalUnits != nil {
				defaultOrgUnit = cert.Subject.OrganizationalUnits
			}
		}

		spec := kmapi.CertificateSpec{
			Alias:      string(alias),
			SecretName: e.CertificateName(alias),
			Subject: &kmapi.X509Subject{
				Organizations:       defaultOrg,
				OrganizationalUnits: defaultOrgUnit,
			},
		}
		// Only the certificates presented on a listening socket need SANs; the
		// client certificate is used to authenticate outbound connections.
		if alias == EtcdServerCert || alias == EtcdPeerCert {
			spec.DNSNames = e.CertificateDNSNames()
		}

		e.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(e.Spec.TLS.Certificates, spec)
	}
}

func (e *Etcd) setDefaultContainerSecurityContext(etcdVersion *catalog.EtcdVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = etcdRunAsUser(etcdVersion)
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.EtcdContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.EtcdContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	e.assignDefaultContainerSecurityContext(etcdVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (e *Etcd) assignDefaultContainerSecurityContext(etcdVersion *catalog.EtcdVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = etcdRunAsUser(etcdVersion)
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = etcdRunAsUser(etcdVersion)
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

// etcdRunAsUser falls back to the uid baked into the upstream etcd image when the
// catalog object does not pin one.
func etcdRunAsUser(etcdVersion *catalog.EtcdVersion) *int64 {
	if etcdVersion != nil && etcdVersion.Spec.SecurityContext.RunAsUser != nil {
		return etcdVersion.Spec.SecurityContext.RunAsUser
	}
	return pointer.Int64P(kubedb.EtcdUserID)
}

func (e *Etcd) setDefaultContainerResourceLimits(podTemplate *ofst.PodTemplateSpec) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.EtcdContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: kubedb.EtcdContainerName,
		}
	}
	apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.EtcdDefaultResources)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *dbContainer)

	apis.SetDefaultResizePolicy(podTemplate.Spec.Containers, podTemplate.Spec.InitContainers)
}
