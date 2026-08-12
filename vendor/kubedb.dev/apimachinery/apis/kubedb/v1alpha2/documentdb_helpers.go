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
	"strconv"
	"time"

	"kubedb.dev/apimachinery/apis"
	catalogv1alpha1 "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/ptr"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	metautil "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
	ofst_util "kmodules.xyz/offshoot-api/util"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (d *DocumentDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDocumentDB))
}

// AsOwner returns owner reference to resources
func (d *DocumentDB) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(d, SchemeGroupVersion.WithKind(ResourceKindDocumentDB))
}

var _ apis.ResourceInfo = &DocumentDB{}

func (d *DocumentDB) OffshootName() string {
	return d.Name
}

func (d *DocumentDB) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		metautil.NameLabelKey:      d.ResourceFQN(),
		metautil.InstanceLabelKey:  d.Name,
		metautil.ManagedByLabelKey: SchemeGroupVersion.Group,
	}
	return metautil.OverwriteKeys(selector, extraSelectors...)
}

func (d *DocumentDB) OffshootLabels() map[string]string {
	return d.offshootLabels(d.OffshootSelectors(), nil)
}

func (d *DocumentDB) PodLabels(podTemplate *ofstv2.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if podTemplate != nil && podTemplate.Labels != nil {
		return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), podTemplate.Labels)
	}
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), nil)
}

func (d *DocumentDB) PodControllerLabels(podTemplate *ofstv2.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if podTemplate != nil && podTemplate.Controller.Labels != nil {
		return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), podTemplate.Controller.Labels)
	}
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), nil)
}

func (d *DocumentDB) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(d.Spec.ServiceTemplates, ServiceAlias(alias))
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (d *DocumentDB) offshootLabels(selector, override map[string]string) map[string]string {
	selector[metautil.ComponentLabelKey] = kubedb.ComponentDatabase
	return metautil.FilterKeys(SchemeGroupVersion.Group, selector, metautil.OverwriteKeys(nil, d.Labels, override))
}

func (d *DocumentDB) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", d.ResourcePlural(), SchemeGroupVersion.Group)
}

func (d *DocumentDB) ResourceShortCode() string {
	return ResourceCodeDocumentDB
}

func (d *DocumentDB) ResourceKind() string {
	return ResourceKindDocumentDB
}

func (d *DocumentDB) ResourceSingular() string {
	return ResourceSingularDocumentDB
}

func (d *DocumentDB) ResourcePlural() string {
	return ResourcePluralDocumentDB
}

func (d *DocumentDB) GetAuthSecretName() string {
	if d.Spec.AuthSecret != nil && d.Spec.AuthSecret.Name != "" {
		return d.Spec.AuthSecret.Name
	}
	return metautil.NameWithSuffix(d.OffshootName(), "auth")
}

func (d *DocumentDB) GetAdminAuthSecretName() string {
	if d.Spec.AdminAuthSecret != nil && d.Spec.AdminAuthSecret.Name != "" {
		return d.Spec.AdminAuthSecret.Name
	}
	return metautil.NameWithSuffix(d.OffshootName(), kubedb.DocumentDBAdminAuthSecretSuffix)
}

// ConfigSecretName returns the name of the operator-generated secret that holds
// tuning (pgtune.conf) and inline (inline.conf) configuration for the database.
func (d *DocumentDB) ConfigSecretName() string {
	uid := string(d.UID)
	if len(uid) >= 6 {
		return metautil.NameWithSuffix(d.OffshootName(), uid[len(uid)-6:])
	}
	return metautil.NameWithSuffix(d.OffshootName(), "config")
}

func (d *DocumentDB) GetStorageClassName() string {
	if d.Spec.Storage == nil || d.Spec.Storage.StorageClassName == nil {
		return ""
	}
	return *d.Spec.Storage.StorageClassName
}

func (d *DocumentDB) ServiceName() string {
	return d.OffshootName()
}

func (d *DocumentDB) StandbyServiceName() string {
	return metautil.NameWithPrefix(d.ServiceName(), "standby")
}

func (d *DocumentDB) GoverningServiceName() string {
	return metautil.NameWithSuffix(d.ServiceName(), "pods")
}

func (d *DocumentDB) PetSetName() string {
	return d.OffshootName()
}

// CertificateName returns the default certificate name / certificate secret name for a certificate alias.
func (d *DocumentDB) CertificateName(alias DocumentDBCertificateAlias) string {
	return metautil.NameWithSuffix(d.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetGRPCIssuerName is the cert-manager CA Issuer that signs the coordinator's gRPC leaf certs.
func (d *DocumentDB) GetGRPCIssuerName() string {
	return metautil.NameWithSuffix(d.OffshootName(), "grpc-issuer")
}

// GetGRPCSelfSignedIssuerName is the self-signed bootstrap Issuer that signs the grpc-ca certificate.
func (d *DocumentDB) GetGRPCSelfSignedIssuerName() string {
	return metautil.NameWithSuffix(d.OffshootName(), "grpc-selfsigned")
}

// tlsConfigUsable reports whether a config can actually issue certificates. A config with no
// issuerRef cannot, so it reads as "not configured" rather than as an empty configuration.
func tlsConfigUsable(c *kmapi.TLSConfig) bool {
	return c != nil && c.IssuerRef != nil
}

// DBTLSConfig returns the TLS config that signs the database-plane certificates: the Postgres
// server certificate and the streaming-replication client certificate.
func (d *DocumentDB) DBTLSConfig() *kmapi.TLSConfig {
	if d.Spec.TLS == nil {
		return nil
	}
	return d.Spec.TLS.DBTLS
}

// GatewayTLSConfig returns the TLS config that signs the MongoDB-wire gateway certificate.
// It falls back to DBTLS so that setting only dbTLS still covers the gateway with one issuer.
func (d *DocumentDB) GatewayTLSConfig() *kmapi.TLSConfig {
	if d.Spec.TLS == nil {
		return nil
	}
	if d.Spec.TLS.GatewayTLS != nil {
		return d.Spec.TLS.GatewayTLS
	}
	return d.Spec.TLS.DBTLS
}

// TLSConfigForAlias routes a certificate alias to the config that owns it. The grpc-* aliases
// are signed by an operator-managed self-signed CA and belong to neither, so they return nil.
func (d *DocumentDB) TLSConfigForAlias(alias DocumentDBCertificateAlias) *kmapi.TLSConfig {
	switch alias {
	case DocumentDBServerCert, DocumentDBClientCert:
		return d.DBTLSConfig()
	case DocumentDBGatewayCert:
		return d.GatewayTLSConfig()
	}
	return nil
}

// AllCertificates returns the union of both configs' certificate specs, deduplicated by alias
// with the gateway config winning for the gateway alias. Callers that enumerate "every managed
// certificate" — rotation, sync, and cleanup — must use this: iterating a single config's list
// makes the other config's certificates look unrecognized, and cleanup then deletes them.
func (d *DocumentDB) AllCertificates() []kmapi.CertificateSpec {
	if d.Spec.TLS == nil {
		return nil
	}
	var out []kmapi.CertificateSpec
	seen := map[string]bool{}
	// Gateway first so its entry wins if the same alias appears in both lists.
	for _, c := range []*kmapi.TLSConfig{d.Spec.TLS.GatewayTLS, d.Spec.TLS.DBTLS} {
		if c == nil {
			continue
		}
		for _, cert := range c.Certificates {
			if seen[cert.Alias] {
				continue
			}
			seen[cert.Alias] = true
			out = append(out, cert)
		}
	}
	return out
}

// IsTLSEnabled reports whether cert-manager TLS is actually configured. A DocumentDBTLSConfig
// carrying only GatewayMutualTLSEnabled (left behind by a ReconfigureTLS removal, so the
// preference survives a later re-add) is not TLS enabled.
func (d *DocumentDB) IsTLSEnabled() bool {
	return tlsConfigUsable(d.DBTLSConfig()) || tlsConfigUsable(d.GatewayTLSConfig())
}

// GetCertSecretName returns the secret name for a certificate alias if provided,
// otherwise returns the default certificate secret name for the given alias.
func (d *DocumentDB) GetCertSecretName(alias DocumentDBCertificateAlias) string {
	if cfg := d.TLSConfigForAlias(alias); cfg != nil {
		name, ok := kmapi.GetCertificateSecretName(cfg.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return d.CertificateName(alias)
}

// SetTLSDefaults fills missing secret names for each certificate alias, mirroring postgres.
// Each config is defaulted independently, so a gatewayTLS with its own issuer still gets its
// secret name filled even when it is configured alongside a separate dbTLS.
func (d *DocumentDB) SetTLSDefaults() {
	if d.Spec.TLS == nil {
		return
	}
	if db := d.DBTLSConfig(); tlsConfigUsable(db) {
		db.Certificates = kmapi.SetMissingSecretNameForCertificate(db.Certificates, string(DocumentDBServerCert), d.CertificateName(DocumentDBServerCert))
		db.Certificates = kmapi.SetMissingSecretNameForCertificate(db.Certificates, string(DocumentDBClientCert), d.CertificateName(DocumentDBClientCert))
	}
	if gw := d.GatewayTLSConfig(); tlsConfigUsable(gw) {
		gw.Certificates = kmapi.SetMissingSecretNameForCertificate(gw.Certificates, string(DocumentDBGatewayCert), d.CertificateName(DocumentDBGatewayCert))
	}
}

// GatewayMutualTLSEnabled reports whether the MongoDB-wire gateway listener requires mutual TLS.
// It stays gated on TLS actually being configured: the toggle is meaningless without certs.
func (d *DocumentDB) GatewayMutualTLSEnabled() bool {
	if !d.IsTLSEnabled() {
		return false
	}
	return d.Spec.TLS.GatewayMutualTLSEnabled == nil || *d.Spec.TLS.GatewayMutualTLSEnabled
}

func (d *DocumentDB) ServiceAccountName() string {
	return d.OffshootName()
}

type documentDBApp struct {
	*DocumentDB
}

func (r documentDBApp) Name() string {
	return r.DocumentDB.Name
}

func (r documentDBApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularDocumentDB))
}

func (d *DocumentDB) AppBindingMeta() appcat.AppBindingMeta {
	return &documentDBApp{d}
}

type documentDBStatsService struct {
	*DocumentDB
}

func (d documentDBStatsService) GetNamespace() string {
	return d.DocumentDB.GetNamespace()
}

func (d documentDBStatsService) ServiceName() string {
	return d.OffshootName() + "-stats"
}

func (d documentDBStatsService) ServiceMonitorName() string {
	return d.ServiceName()
}

func (d documentDBStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return d.OffshootLabels()
}

func (d documentDBStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (d documentDBStatsService) Scheme() string {
	sc := promapi.SchemeHTTP
	return sc.String()
}

func (d documentDBStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (d *DocumentDB) StatsService() mona.StatsAccessor {
	return &documentDBStatsService{d}
}

func (d *DocumentDB) StatsServiceLabels() map[string]string {
	return d.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (d *DocumentDB) SetDefaults(_ client.Client, documentDBVersion catalogv1alpha1.DocumentDBVersion) {
	if d == nil {
		return
	}
	if d.Spec.StandbyMode == nil {
		d.Spec.StandbyMode = ptr.To(HotDocDBStandbyMode)
	}
	if d.Spec.StreamingMode == nil {
		d.Spec.StreamingMode = ptr.To(AsynchronousDocDBStreamingMode)
	}
	if d.Spec.ClientAuthMode == "" {
		d.Spec.ClientAuthMode = DocDBClientAuthModeScram
	}
	if d.Spec.SSLMode == "" {
		if d.IsTLSEnabled() {
			d.Spec.SSLMode = DocumentDBSSLModeVerifyFull
		} else {
			d.Spec.SSLMode = DocumentDBSSLModeDisable
		}
	}
	d.SetTLSDefaults()
	if d.Spec.StorageType == "" {
		d.Spec.StorageType = StorageTypeDurable
	}
	if d.Spec.DeletionPolicy == "" {
		d.Spec.DeletionPolicy = DeletionPolicyDelete
	}
	if d.Spec.Replicas == nil {
		d.Spec.Replicas = ptr.To(int32(1))
	}

	if d.Spec.AuthSecret == nil {
		d.Spec.AuthSecret = &SecretReference{}
	}
	if d.Spec.AuthSecret.Kind == "" {
		d.Spec.AuthSecret.Kind = kubedb.ResourceKindSecret
	}
	if d.Spec.AdminAuthSecret == nil {
		d.Spec.AdminAuthSecret = &SecretReference{}
	}
	if d.Spec.AdminAuthSecret.Kind == "" {
		d.Spec.AdminAuthSecret.Kind = kubedb.ResourceKindSecret
	}

	if d.Spec.LeaderElection == nil {
		d.Spec.LeaderElection = &DocumentDBLeaderElectionConfig{
			// The upper limit of election timeout is 50000ms (50s), which should only be used when deploying a
			// globally-distributed etcd cluster. A reasonable round-trip time for the continental United States is around 130-150ms,
			// and the time between US and Japan is around 350-400ms. If the network has uneven performance or regular packet
			// delays/loss then it is possible that a couple of retries may be necessary to successfully send a packet.
			// So 5s is a safe upper limit of global round-trip time. As the election timeout should be an order of magnitude
			// bigger than broadcast time, in the case of ~5s for a globally distributed cluster, then 50 seconds becomes
			// a reasonable maximum.
			Period: metav1.Duration{Duration: 1 * time.Second},
			// the amount of HeartbeatTick can be missed before the failOver
			ElectionTick: 15,
			// this value should be one.
			HeartbeatTick: 1,
			// we have set this default to 67108864. if the difference between primary and replica is more then this,
			// the replica node is going to manually sync itself.
			MaximumLagBeforeFailover: 64 * 1024 * 1024,
		}
	}
	if d.Spec.LeaderElection.TransferLeadershipInterval == nil {
		d.Spec.LeaderElection.TransferLeadershipInterval = &metav1.Duration{Duration: 1 * time.Second}
	}
	if d.Spec.LeaderElection.TransferLeadershipTimeout == nil {
		d.Spec.LeaderElection.TransferLeadershipTimeout = &metav1.Duration{Duration: 60 * time.Second}
	}

	d.initializePodTemplates()

	if d.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		d.Spec.PodTemplate.Spec.ServiceAccountName = d.OffshootName()
	}

	d.SetDefaultPodSecurityContext(d.Spec.PodTemplate, &documentDBVersion)
	d.SetInitContainerDefaults(d.Spec.PodTemplate, &documentDBVersion)
	d.SetDocumentDBContainerDefaults(d.Spec.PodTemplate, &documentDBVersion)
	d.SetCoordinatorContainerDefaults(d.Spec.PodTemplate, &documentDBVersion)
	apis.SetDefaultResizePolicy(d.Spec.PodTemplate.Spec.Containers, d.Spec.PodTemplate.Spec.InitContainers)
	d.SetDefaultReplicationMode()
	d.SetHealthCheckerDefaults()
}

// SetDefaultReplicationMode sets the default replication mode.
// WALKeepSize will be the default policy for DocumentDB.
func (d *DocumentDB) SetDefaultReplicationMode() {
	if d.Spec.Replication == nil {
		d.Spec.Replication = &DocumentDBReplication{}
	}
	if d.Spec.Replication.WALLimitPolicy == "" {
		d.Spec.Replication.WALLimitPolicy = WALKeepSize
	}
	if d.Spec.Replication.WALLimitPolicy == WALKeepSegment && d.Spec.Replication.WalKeepSegment == nil {
		d.Spec.Replication.WalKeepSegment = pointer.Int32P(96)
	}
	if d.Spec.Replication.WALLimitPolicy == WALKeepSize && d.Spec.Replication.WalKeepSizeInMegaBytes == nil {
		d.Spec.Replication.WalKeepSizeInMegaBytes = pointer.Int32P(1536)
	}
	if d.Spec.Replication.WALLimitPolicy == ReplicationSlot && d.Spec.Replication.MaxSlotWALKeepSizeInMegaBytes == nil {
		d.Spec.Replication.MaxSlotWALKeepSizeInMegaBytes = pointer.Int32P(-1)
	}
}

func (d *DocumentDB) SetDefaultPodSecurityContext(podTemplate *ofstv2.PodTemplateSpec, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
	if podTemplate == nil {
		return
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
	if podTemplate.Spec.SecurityContext.RunAsUser == nil {
		podTemplate.Spec.SecurityContext.RunAsUser = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
	if podTemplate.Spec.SecurityContext.RunAsGroup == nil {
		podTemplate.Spec.SecurityContext.RunAsGroup = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
}

func (d *DocumentDB) SetInitContainerDefaults(podTemplate *ofstv2.PodTemplateSpec, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
	if podTemplate == nil {
		return
	}
	container := ofst_util.EnsureInitContainerExists(podTemplate, kubedb.DocumentDBInitContainerName)
	d.setContainerDefaultSecurityContext(container, documentDBVersion)
	d.setContainerDefaultResources(container, *kubedb.DefaultInitContainerResource.DeepCopy())
}

func (d *DocumentDB) SetDocumentDBContainerDefaults(podTemplate *ofstv2.PodTemplateSpec, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
	if podTemplate == nil {
		return
	}
	container := ofst_util.EnsureContainerExists(podTemplate, kubedb.DocumentDBContainerName)
	d.setContainerDefaultSecurityContext(container, documentDBVersion)
	d.setContainerDefaultResources(container, *kubedb.DefaultResources.DeepCopy())
}

func (d *DocumentDB) SetCoordinatorContainerDefaults(podTemplate *ofstv2.PodTemplateSpec, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
	if podTemplate == nil {
		return
	}
	container := ofst_util.EnsureContainerExists(podTemplate, kubedb.DocumentDBCoordinatorContainerName)
	d.setContainerDefaultSecurityContext(container, documentDBVersion)
	d.setContainerDefaultResources(container, *kubedb.CoordinatorDefaultResources.DeepCopy())
}

func (d *DocumentDB) setContainerDefaultSecurityContext(container *core.Container, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	d.assignDefaultContainerSecurityContext(container.SecurityContext, documentDBVersion)
}

func (d *DocumentDB) setContainerDefaultResources(container *core.Container, defaultResources core.ResourceRequirements) {
	if container != nil {
		apis.SetDefaultResourceLimits(&container.Resources, defaultResources)
	}
}

func (d *DocumentDB) assignDefaultContainerSecurityContext(sc *core.SecurityContext, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
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
		sc.RunAsUser = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (d *DocumentDB) initializePodTemplates() {
	if d.Spec.PodTemplate == nil {
		d.Spec.PodTemplate = new(ofstv2.PodTemplateSpec)
	}
}

func (d *DocumentDB) GetPersistentSecrets() []string {
	secrets := make([]string, 0, 2)
	if !IsVirtualAuthSecretReferred(d.Spec.AuthSecret) && d.Spec.AuthSecret != nil && d.Spec.AuthSecret.Name != "" {
		secrets = append(secrets, d.GetAuthSecretName())
	}
	secrets = append(secrets, d.GetAdminAuthSecretName())
	return secrets
}

func (d *DocumentDB) Finalizer() string {
	return fmt.Sprintf("%s/%s", apis.Finalizer, d.ResourceSingular())
}

func (d *DocumentDB) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicasOfPetSet(lister.PetSets(d.Namespace), labels.SelectorFromSet(d.OffshootLabels()), expectedItems)
}

func (d *DocumentDB) SetHealthCheckerDefaults() {
	if d.Spec.HealthChecker.PeriodSeconds == nil {
		d.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if d.Spec.HealthChecker.TimeoutSeconds == nil {
		d.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if d.Spec.HealthChecker.FailureThreshold == nil {
		d.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

// GetSharedBufferSizeForDocumentdb this func takes a input type int64 which is in bytes
// return the 25% of the input in Bytes
func GetSharedBufferSizeForDocumentdb(resource *resource.Quantity) string {
	// no more than 25% of main memory (RAM)
	minSharedBuffer := int64(128)
	ret := minSharedBuffer
	if resource != nil {
		ret = resource.Value() / (4 * 1024)
	}
	// the shared buffer value can't be less then this
	// 128 KB  is the minimum
	if ret < minSharedBuffer {
		ret = minSharedBuffer
	}

	// check If the ret value need to convert into MB
	// why need this? -> PostgreSQL officially stores shared_buffers as an int32 that's why if the value is greater than 2147483648B.
	// It's going to through and error that the value is going to cross the limit.

	sharedBuffer := fmt.Sprintf("%skB", strconv.FormatInt(ret, 10))
	if ret > kubedb.SharedBuffersGbAsKiloByte {
		// convert the ret as MB devide by SharedBuffersMbAsByte
		ret /= kubedb.SharedBuffersMbAsKiloByte
		sharedBuffer = fmt.Sprintf("%sMB", strconv.FormatInt(ret, 10))
	}

	return sharedBuffer
}

func (d *DocumentDB) GetDeletionPolicy() string {
	return string(d.Spec.DeletionPolicy)
}
