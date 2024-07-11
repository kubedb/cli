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

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	core_util "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
)

const (
	RedisShardAffinityTemplateVar = "SHARD_INDEX"
)

func (*Redis) Hub() {}

func (r Redis) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedis))
}

func (r *Redis) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(r, SchemeGroupVersion.WithKind(ResourceKindRedis))
}

var _ apis.ResourceInfo = &Redis{}

func (r Redis) OffshootName() string {
	return r.Name
}

func (r Redis) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      r.ResourceFQN(),
		meta_util.InstanceLabelKey:  r.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (r Redis) OffshootLabels() map[string]string {
	return r.offshootLabels(r.OffshootSelectors(), nil)
}

func (r Redis) PodLabels(extraLabels ...map[string]string) map[string]string {
	return r.offshootLabels(meta_util.OverwriteKeys(r.OffshootSelectors(), extraLabels...), r.Spec.PodTemplate.Labels)
}

func (r Redis) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return r.offshootLabels(meta_util.OverwriteKeys(r.OffshootSelectors(), extraLabels...), r.Spec.PodTemplate.Controller.Labels)
}

func (r Redis) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(r.Spec.ServiceTemplates, alias)
	return r.offshootLabels(meta_util.OverwriteKeys(r.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (r Redis) offshootLabels(selector, overrides map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, r.Labels, overrides))
}

func (r Redis) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralRedis, kubedb.GroupName)
}

func (r Redis) ResourceShortCode() string {
	return ResourceCodeRedis
}

func (r Redis) ResourceKind() string {
	return ResourceKindRedis
}

func (r Redis) ResourceSingular() string {
	return ResourceSingularRedis
}

func (r Redis) ResourcePlural() string {
	return ResourcePluralRedis
}

func (r Redis) GetAuthSecretName() string {
	if r.Spec.AuthSecret != nil && r.Spec.AuthSecret.Name != "" {
		return r.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(r.OffshootName(), "auth")
}

func (r Redis) ServiceName() string {
	return r.OffshootName()
}

func (r Redis) StandbyServiceName() string {
	return meta_util.NameWithPrefix(r.ServiceName(), "standby")
}

func (r Redis) GoverningServiceName() string {
	return meta_util.NameWithSuffix(r.ServiceName(), "pods")
}

func (r Redis) ConfigSecretName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "config")
}

func (r Redis) CustomConfigSecretName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "custom-config")
}

func (r Redis) BaseNameForShard() string {
	return fmt.Sprintf("%s-shard", r.OffshootName())
}

func (r Redis) PetSetNameWithShard(i int) string {
	return fmt.Sprintf("%s%d", r.BaseNameForShard(), i)
}

func (r Redis) Address() string {
	return fmt.Sprintf("%v.%v.svc:%d", r.Name, r.Namespace, kubedb.RedisDatabasePort)
}

type redisApp struct {
	*Redis
}

func (r redisApp) Name() string {
	return r.Redis.Name
}

func (r redisApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularRedis))
}

func (r Redis) AppBindingMeta() appcat.AppBindingMeta {
	return &redisApp{&r}
}

type redisStatsService struct {
	*Redis
}

func (r redisStatsService) GetNamespace() string {
	return r.Redis.GetNamespace()
}

func (r redisStatsService) ServiceName() string {
	return r.OffshootName() + "-stats"
}

func (r redisStatsService) ServiceMonitorName() string {
	return r.ServiceName()
}

func (p redisStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return p.OffshootLabels()
}

func (r redisStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (r redisStatsService) Scheme() string {
	return ""
}

func (r redisStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (r Redis) StatsService() mona.StatsAccessor {
	return &redisStatsService{&r}
}

func (r Redis) StatsServiceLabels() map[string]string {
	return r.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (r *Redis) SetDefaults(rdVersion *catalog.RedisVersion) {
	if r == nil {
		return
	}

	// perform defaulting
	if r.Spec.Mode == "" {
		r.Spec.Mode = RedisModeStandalone
	} else if r.Spec.Mode == RedisModeCluster {
		if r.Spec.Cluster == nil {
			r.Spec.Cluster = &RedisClusterSpec{}
		}
		if r.Spec.Cluster.Shards == nil {
			r.Spec.Cluster.Shards = pointer.Int32P(3)
		}
		if r.Spec.Cluster.Replicas == nil {
			r.Spec.Cluster.Replicas = pointer.Int32P(2)
		}
	}
	if r.Spec.StorageType == "" {
		r.Spec.StorageType = StorageTypeDurable
	}
	if r.Spec.DeletionPolicy == "" {
		r.Spec.DeletionPolicy = DeletionPolicyDelete
	}
	r.setDefaultContainerSecurityContext(rdVersion, &r.Spec.PodTemplate)
	r.setDefaultContainerResourceLimits(&r.Spec.PodTemplate)

	if r.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		r.Spec.PodTemplate.Spec.ServiceAccountName = r.OffshootName()
	}

	labels := r.OffshootSelectors()
	if r.Spec.Mode == RedisModeCluster {
		labels[kubedb.RedisShardKey] = r.ShardNodeTemplate()
	}

	r.SetTLSDefaults()
	r.SetHealthCheckerDefaults()
	r.Spec.Monitor.SetDefaults()
	if r.Spec.Monitor != nil && r.Spec.Monitor.Prometheus != nil {
		if r.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			r.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = rdVersion.Spec.SecurityContext.RunAsUser
		}
		if r.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			r.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = rdVersion.Spec.SecurityContext.RunAsUser
		}
	}
}

func (r *Redis) SetHealthCheckerDefaults() {
	if r.Spec.HealthChecker.PeriodSeconds == nil {
		r.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if r.Spec.HealthChecker.TimeoutSeconds == nil {
		r.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if r.Spec.HealthChecker.FailureThreshold == nil {
		r.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (r *Redis) SetTLSDefaults() {
	if r.Spec.TLS == nil || r.Spec.TLS.IssuerRef == nil {
		return
	}
	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(RedisServerCert), r.CertificateName(RedisServerCert))
	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(RedisClientCert), r.CertificateName(RedisClientCert))
	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(RedisMetricsExporterCert), r.CertificateName(RedisMetricsExporterCert))
}

func (r *RedisSpec) GetPersistentSecrets() []string {
	if r == nil {
		return nil
	}

	var secrets []string
	if r.AuthSecret != nil {
		secrets = append(secrets, r.AuthSecret.Name)
	}
	return secrets
}

func (r *Redis) setDefaultContainerSecurityContext(rdVersion *catalog.RedisVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		podTemplate = &ofstv2.PodTemplateSpec{}
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = rdVersion.Spec.SecurityContext.RunAsUser
	}

	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.RedisContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: kubedb.RedisContainerName,
		}
	}
	if dbContainer.SecurityContext == nil {
		dbContainer.SecurityContext = &core.SecurityContext{}
	}
	r.assignDefaultContainerSecurityContext(rdVersion, dbContainer.SecurityContext)
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *dbContainer)

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.RedisInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.RedisInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	r.assignDefaultContainerSecurityContext(rdVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = core_util.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)

	if r.Spec.Mode == RedisModeSentinel {
		coordinatorContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.RedisCoordinatorContainerName)
		if coordinatorContainer == nil {
			coordinatorContainer = &core.Container{
				Name: kubedb.RedisCoordinatorContainerName,
			}
		}
		if coordinatorContainer.SecurityContext == nil {
			coordinatorContainer.SecurityContext = &core.SecurityContext{}
		}
		r.assignDefaultContainerSecurityContext(rdVersion, coordinatorContainer.SecurityContext)
		podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *coordinatorContainer)
	}
}

func (r *Redis) assignDefaultContainerSecurityContext(rdVersion *catalog.RedisVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = rdVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = rdVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (r *Redis) setDefaultContainerResourceLimits(podTemplate *ofstv2.PodTemplateSpec) {
	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.RedisContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.RedisInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}

	if r.Spec.Mode == RedisModeSentinel {
		coordinatorContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.RedisCoordinatorContainerName)
		if coordinatorContainer != nil && (coordinatorContainer.Resources.Requests == nil && coordinatorContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&coordinatorContainer.Resources, kubedb.CoordinatorDefaultResources)
		}
	}
}

func (r Redis) ShardNodeTemplate() string {
	if r.Spec.Mode == RedisModeStandalone {
		panic("shard template is not applicable to a standalone redis server")
	}
	return fmt.Sprintf("${%s}", RedisShardAffinityTemplateVar)
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (r *Redis) CertificateName(alias RedisCertificateAlias) string {
	return meta_util.NameWithSuffix(r.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any provide,
// otherwise returns default certificate secret name for the given alias.
func (r *Redis) GetCertSecretName(alias RedisCertificateAlias) string {
	if r.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(r.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return r.CertificateName(alias)
}

func (r *Redis) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	if r.Spec.Cluster != nil {
		expectedItems = int(pointer.Int32(r.Spec.Cluster.Shards))
	}
	return checkReplicas(lister.PetSets(r.Namespace), labels.SelectorFromSet(r.OffshootLabels()), expectedItems)
}
