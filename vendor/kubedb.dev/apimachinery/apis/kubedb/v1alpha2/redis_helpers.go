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

	"gomodules.xyz/pointer"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	core_util "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	RedisShardAffinityTemplateVar = "SHARD_INDEX"
)

func (r Redis) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedis))
}

var _ apis.ResourceInfo = &Redis{}

func (r Redis) OffshootName() string {
	return r.Name
}

func (r Redis) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      r.ResourceFQN(),
		meta_util.InstanceLabelKey:  r.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (r Redis) OffshootLabels() map[string]string {
	out := r.OffshootSelectors()
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, out, r.Labels)
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

func (r Redis) ServiceName() string {
	return r.OffshootName()
}

func (r Redis) GoverningServiceName() string {
	return meta_util.NameWithSuffix(r.ServiceName(), "pods")
}

func (r Redis) ConfigSecretName() string {
	return r.OffshootName()
}

func (r Redis) BaseNameForShard() string {
	return fmt.Sprintf("%s-shard", r.OffshootName())
}

func (r Redis) StatefulSetNameWithShard(i int) string {
	return fmt.Sprintf("%s%d", r.BaseNameForShard(), i)
}

func (r Redis) Address() string {
	return fmt.Sprintf("%v.%v.svc:%d", r.Name, r.Namespace, RedisDatabasePort)
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
	return DefaultStatsPath
}

func (r redisStatsService) Scheme() string {
	return ""
}

func (r Redis) StatsService() mona.StatsAccessor {
	return &redisStatsService{&r}
}

func (r Redis) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, r.OffshootSelectors(), r.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (r *Redis) SetDefaults(topology *core_util.Topology) {
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
		if r.Spec.Cluster.Master == nil {
			r.Spec.Cluster.Master = pointer.Int32P(3)
		}
		if r.Spec.Cluster.Replicas == nil {
			r.Spec.Cluster.Replicas = pointer.Int32P(1)
		}
	}
	if r.Spec.StorageType == "" {
		r.Spec.StorageType = StorageTypeDurable
	}
	if r.Spec.TerminationPolicy == "" {
		r.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if r.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		r.Spec.PodTemplate.Spec.ServiceAccountName = r.OffshootName()
	}

	labels := r.OffshootSelectors()
	if r.Spec.Mode == RedisModeCluster {
		labels[RedisShardKey] = r.ShardNodeTemplate()
	}
	r.setDefaultAffinity(&r.Spec.PodTemplate, labels, topology)

	r.Spec.Monitor.SetDefaults()

	r.SetTLSDefaults()
	SetDefaultResourceLimits(&r.Spec.PodTemplate.Spec.Container.Resources, DefaultResourceLimits)
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
	return nil
}

func (r *Redis) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
	if podTemplate == nil {
		return
	} else if podTemplate.Spec.Affinity != nil {
		topology.ConvertAffinity(podTemplate.Spec.Affinity)
		return
	}

	podTemplate.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				// Prefer to not schedule multiple pods on the same node
				{
					Weight: 100,
					PodAffinityTerm: corev1.PodAffinityTerm{
						Namespaces: []string{r.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: labels,
						},

						TopologyKey: corev1.LabelHostname,
					},
				},
				// Prefer to not schedule multiple pods on the node with same zone
				{
					Weight: 50,
					PodAffinityTerm: corev1.PodAffinityTerm{
						Namespaces: []string{r.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: labels,
						},
						TopologyKey: topology.LabelZone,
					},
				},
			},
		},
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

// MustCertSecretName returns the secret name for a certificate alias
func (r *Redis) MustCertSecretName(alias RedisCertificateAlias) string {
	if r == nil {
		panic("missing Redis database")
	} else if r.Spec.TLS == nil {
		panic(fmt.Errorf("Redis %s/%s is missing tls spec", r.Namespace, r.Name))
	}
	name, ok := kmapi.GetCertificateSecretName(r.Spec.TLS.Certificates, string(alias))
	if !ok {
		panic(fmt.Errorf("Redis %s/%s is missing secret name for %s certificate", r.Namespace, r.Name, alias))
	}
	return name
}

func (r *Redis) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	if r.Spec.Cluster != nil {
		expectedItems = int(pointer.Int32(r.Spec.Cluster.Master))
	}
	return checkReplicas(lister.StatefulSets(r.Namespace), labels.SelectorFromSet(r.OffshootLabels()), expectedItems)
}
