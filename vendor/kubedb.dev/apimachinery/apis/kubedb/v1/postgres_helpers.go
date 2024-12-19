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

package v1

import (
	"fmt"
	"strconv"
	"time"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/ptr"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
	ofst_util "kmodules.xyz/offshoot-api/util"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
)

func (*Postgres) Hub() {}

func (_ Postgres) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPostgres))
}

func (p *Postgres) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(p, SchemeGroupVersion.WithKind(ResourceKindPostgres))
}

var _ apis.ResourceInfo = &Postgres{}

func (p Postgres) OffshootName() string {
	return p.Name
}

func (p Postgres) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      p.ResourceFQN(),
		meta_util.InstanceLabelKey:  p.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (p Postgres) OffshootLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), nil)
}

func (p Postgres) PodLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), p.Spec.PodTemplate.Labels)
}

func (p Postgres) PodControllerLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), p.Spec.PodTemplate.Controller.Labels)
}

func (p Postgres) SidekickLabels(skName string) map[string]string {
	return meta_util.OverwriteKeys(nil, kubedb.CommonSidekickLabels(), map[string]string{
		meta_util.InstanceLabelKey: skName,
		kubedb.SidekickOwnerName:   p.Name,
		kubedb.SidekickOwnerKind:   p.ResourceFQN(),
	})
}

func (p Postgres) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(p.Spec.ServiceTemplates, alias)
	return p.offshootLabels(meta_util.OverwriteKeys(p.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (p Postgres) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, p.Labels, override))
}

func (p Postgres) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPostgres, kubedb.GroupName)
}

func (p Postgres) ResourceShortCode() string {
	return ResourceCodePostgres
}

func (p Postgres) ResourceKind() string {
	return ResourceKindPostgres
}

func (p Postgres) ResourceSingular() string {
	return ResourceSingularPostgres
}

func (p Postgres) ResourcePlural() string {
	return ResourcePluralPostgres
}

func (p Postgres) GetAuthSecretName() string {
	if p.Spec.AuthSecret != nil && p.Spec.AuthSecret.Name != "" {
		return p.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(p.OffshootName(), "auth")
}

func (p Postgres) ServiceName() string {
	return p.OffshootName()
}

func (p Postgres) StandbyServiceName() string {
	return meta_util.NameWithPrefix(p.ServiceName(), "standby")
}

func (p Postgres) GoverningServiceName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "pods")
}

type postgresApp struct {
	*Postgres
}

func (r postgresApp) Name() string {
	return r.Postgres.Name
}

func (r postgresApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPostgres))
}

func (p Postgres) AppBindingMeta() appcat.AppBindingMeta {
	return &postgresApp{&p}
}

type postgresStatsService struct {
	*Postgres
}

func (p postgresStatsService) GetNamespace() string {
	return p.Postgres.GetNamespace()
}

func (p postgresStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p postgresStatsService) ServiceMonitorName() string {
	return p.ServiceName()
}

func (p postgresStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return p.OffshootLabels()
}

func (p postgresStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (p postgresStatsService) Scheme() string {
	return ""
}

func (p postgresStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (p Postgres) StatsService() mona.StatsAccessor {
	return &postgresStatsService{&p}
}

func (p Postgres) StatsServiceLabels() map[string]string {
	return p.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (p *Postgres) SetDefaults(postgresVersion *catalog.PostgresVersion) {
	if p == nil {
		return
	}
	if p.Spec.StandbyMode == nil {
		p.Spec.StandbyMode = ptr.To(HotPostgresStandbyMode)
	}
	if p.Spec.StorageType == "" {
		p.Spec.StorageType = StorageTypeDurable
	}
	if p.Spec.DeletionPolicy == "" {
		p.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if p.Spec.LeaderElection == nil {
		p.Spec.LeaderElection = &PostgreLeaderElectionConfig{
			// The upper limit of election timeout is 50000ms (50s), which should only be used when deploying a
			// globally-distributed etcd cluster. A reasonable round-trip time for the continental United States is around 130-150ms,
			// and the time between US and Japan is around 350-400ms. If the network has uneven performance or regular packet
			// delays/loss then it is possible that a couple of retries may be necessary to successfully send a packet.
			// So 5s is a safe upper limit of global round-trip time. As the election timeout should be an order of magnitude
			// bigger than broadcast time, in the case of ~5s for a globally distributed cluster, then 50 seconds becomes
			// a reasonable maximum.
			Period: metav1.Duration{Duration: 300 * time.Millisecond},
			// the amount of HeartbeatTick can be missed before the failOver
			ElectionTick: 10,
			// this value should be one.
			HeartbeatTick: 1,
			// we have set this default to 67108864. if the difference between primary and replica is more then this,
			// the replica node is going to manually sync itself.
			MaximumLagBeforeFailover: 64 * 1024 * 1024,
		}
	}
	if p.Spec.LeaderElection.TransferLeadershipInterval == nil {
		p.Spec.LeaderElection.TransferLeadershipInterval = &metav1.Duration{Duration: 1 * time.Second}
	}
	if p.Spec.LeaderElection.TransferLeadershipTimeout == nil {
		p.Spec.LeaderElection.TransferLeadershipTimeout = &metav1.Duration{Duration: 60 * time.Second}
	}

	if p.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		p.Spec.PodTemplate.Spec.ServiceAccountName = p.OffshootName()
	}

	if p.Spec.TLS != nil {
		if p.Spec.SSLMode == "" {
			p.Spec.SSLMode = PostgresSSLModeVerifyFull
		}
		if p.Spec.ClientAuthMode == "" {
			p.Spec.ClientAuthMode = ClientAuthModeMD5
		}
	} else {
		if p.Spec.SSLMode == "" {
			p.Spec.SSLMode = PostgresSSLModeDisable
		}
		if p.Spec.ClientAuthMode == "" {
			p.Spec.ClientAuthMode = ClientAuthModeMD5
		}
	}

	p.setDefaultPodSecurityContext(&p.Spec.PodTemplate, postgresVersion)
	p.setPostgresContainerDefaults(&p.Spec.PodTemplate, postgresVersion)
	p.setCoordinatorContainerDefaults(&p.Spec.PodTemplate, postgresVersion)
	p.setInitContainerDefaults(&p.Spec.PodTemplate, postgresVersion)

	// Need to set FSGroup equal to  p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup.
	// So that /var/pv directory have the group permission for the RunAsGroup user GID.
	// Otherwise, We will get write permission denied.
	p.SetDefaultReplicationMode(postgresVersion)
	p.SetArbiterDefault()
	p.SetTLSDefaults()
	p.SetHealthCheckerDefaults()

	p.Spec.Monitor.SetDefaults()
	if p.Spec.Monitor != nil && p.Spec.Monitor.Prometheus != nil {
		if p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = postgresVersion.Spec.SecurityContext.RunAsUser
		}
		if p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = postgresVersion.Spec.SecurityContext.RunAsUser
		}
	}
	if p.Spec.Init != nil && p.Spec.Init.Archiver != nil && p.Spec.Init.Archiver.ReplicationStrategy == nil {
		p.Spec.Init.Archiver.ReplicationStrategy = ptr.To(ReplicationStrategyNone)
	}
}

func getMajorPgVersion(postgresVersion *catalog.PostgresVersion) (uint64, error) {
	ver, err := semver.NewVersion(postgresVersion.Spec.Version)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to get postgres major.")
	}
	return ver.Major(), nil
}

// SetDefaultReplicationMode set the default replication mode.
// Replication slot will be prioritized if no WalLimitPolicy is mentioned
func (p *Postgres) SetDefaultReplicationMode(postgresVersion *catalog.PostgresVersion) {
	majorVersion, _ := getMajorPgVersion(postgresVersion)
	if p.Spec.Replication == nil {
		p.Spec.Replication = &PostgresReplication{}
	}
	if p.Spec.Replication.WALLimitPolicy == "" {
		if majorVersion <= uint64(12) {
			p.Spec.Replication.WALLimitPolicy = WALKeepSegment
		} else {
			p.Spec.Replication.WALLimitPolicy = WALKeepSize
		}
	}
	if p.Spec.Replication.WALLimitPolicy == WALKeepSegment && p.Spec.Replication.WalKeepSegment == nil {
		p.Spec.Replication.WalKeepSegment = pointer.Int32P(96)
	}
	if p.Spec.Replication.WALLimitPolicy == WALKeepSize && p.Spec.Replication.WalKeepSizeInMegaBytes == nil {
		p.Spec.Replication.WalKeepSizeInMegaBytes = pointer.Int32P(1536)
	}
	if p.Spec.Replication.WALLimitPolicy == ReplicationSlot && p.Spec.Replication.MaxSlotWALKeepSizeInMegaBytes == nil {
		p.Spec.Replication.MaxSlotWALKeepSizeInMegaBytes = pointer.Int32P(-1)
	}
}

func (p *Postgres) SetArbiterDefault() {
	if ptr.Deref(p.Spec.Replicas, 0)%2 == 0 && p.Spec.Arbiter == nil {
		p.Spec.Arbiter = &ArbiterSpec{
			Resources: core.ResourceRequirements{},
		}
		apis.SetDefaultResourceLimits(&p.Spec.Arbiter.Resources, kubedb.DefaultArbiter(false))
	}
}

func (p *Postgres) setDefaultPodSecurityContext(podTemplate *ofstv2.PodTemplateSpec, pgVersion *catalog.PostgresVersion) {
	if podTemplate == nil {
		return
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = pgVersion.Spec.SecurityContext.RunAsUser
	}
	if podTemplate.Spec.SecurityContext.RunAsUser == nil {
		podTemplate.Spec.SecurityContext.RunAsUser = pgVersion.Spec.SecurityContext.RunAsUser
	}
	if podTemplate.Spec.SecurityContext.RunAsGroup == nil {
		podTemplate.Spec.SecurityContext.RunAsGroup = pgVersion.Spec.SecurityContext.RunAsUser
	}
}

func (p *Postgres) setInitContainerDefaults(podTemplate *ofstv2.PodTemplateSpec, pgVersion *catalog.PostgresVersion) {
	if podTemplate == nil {
		return
	}
	container := ofst_util.EnsureInitContainerExists(podTemplate, kubedb.PostgresInitContainerName)
	p.setContainerDefaultSecurityContext(container, pgVersion)
	p.setContainerDefaultResources(container, *kubedb.DefaultInitContainerResource.DeepCopy())
}

func (p *Postgres) setPostgresContainerDefaults(podTemplate *ofstv2.PodTemplateSpec, pgVersion *catalog.PostgresVersion) {
	if podTemplate == nil {
		return
	}
	container := ofst_util.EnsureContainerExists(podTemplate, kubedb.PostgresContainerName)
	p.setContainerDefaultSecurityContext(container, pgVersion)
	p.setContainerDefaultResources(container, *kubedb.DefaultResources.DeepCopy())
}

func (p *Postgres) setCoordinatorContainerDefaults(podTemplate *ofstv2.PodTemplateSpec, pgVersion *catalog.PostgresVersion) {
	if podTemplate == nil {
		return
	}
	container := ofst_util.EnsureContainerExists(podTemplate, kubedb.PostgresCoordinatorContainerName)
	p.setContainerDefaultSecurityContext(container, pgVersion)
	p.setContainerDefaultResources(container, *kubedb.CoordinatorDefaultResources.DeepCopy())
}

func (p *Postgres) setContainerDefaultSecurityContext(container *core.Container, pgVersion *catalog.PostgresVersion) {
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	p.assignDefaultContainerSecurityContext(container.SecurityContext, pgVersion)
}

func (p *Postgres) setContainerDefaultResources(container *core.Container, defaultResources core.ResourceRequirements) {
	if container.Resources.Requests == nil && container.Resources.Limits == nil {
		apis.SetDefaultResourceLimits(&container.Resources, defaultResources)
	}
}

func (p *Postgres) assignDefaultContainerSecurityContext(sc *core.SecurityContext, pgVersion *catalog.PostgresVersion) {
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
		sc.RunAsUser = pgVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = pgVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (p *Postgres) SetTLSDefaults() {
	if p.Spec.TLS == nil || p.Spec.TLS.IssuerRef == nil {
		return
	}
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PostgresServerCert), p.CertificateName(PostgresServerCert))
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PostgresClientCert), p.CertificateName(PostgresClientCert))
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PostgresMetricsExporterCert), p.CertificateName(PostgresMetricsExporterCert))
}

func (p *PostgresSpec) GetPersistentSecrets() []string {
	if p == nil {
		return nil
	}

	var secrets []string
	if p.AuthSecret != nil {
		secrets = append(secrets, p.AuthSecret.Name)
	}
	return secrets
}

func (p *Postgres) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicas(lister.PetSets(p.Namespace), labels.SelectorFromSet(p.OffshootLabels()), expectedItems)
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (p *Postgres) CertificateName(alias PostgresCertificateAlias) string {
	return meta_util.NameWithSuffix(p.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any provide,
// otherwise returns default certificate secret name for the given alias.
func (p *Postgres) GetCertSecretName(alias PostgresCertificateAlias) string {
	if p.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(p.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return p.CertificateName(alias)
}

// GetSharedBufferSizeForPostgres this func takes a input type int64 which is in bytes
// return the 25% of the input in Bytes
func GetSharedBufferSizeForPostgres(resource *resource.Quantity) string {
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

func (m *Postgres) SetHealthCheckerDefaults() {
	if m.Spec.HealthChecker.PeriodSeconds == nil {
		m.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if m.Spec.HealthChecker.TimeoutSeconds == nil {
		m.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if m.Spec.HealthChecker.FailureThreshold == nil {
		m.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (m *Postgres) IsRemoteReplica() bool {
	return m.Spec.RemoteReplica != nil
}
