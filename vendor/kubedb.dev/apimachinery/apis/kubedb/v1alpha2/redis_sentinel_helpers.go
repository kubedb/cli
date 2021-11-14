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

func (rs RedisSentinel) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedisSentinel))
}

var _ apis.ResourceInfo = &RedisSentinel{}

func (rs RedisSentinel) OffshootName() string {
	return rs.Name
}

func (rs RedisSentinel) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      rs.ResourceFQN(),
		meta_util.InstanceLabelKey:  rs.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (rs RedisSentinel) OffshootLabels() map[string]string {
	return rs.offshootLabels(rs.OffshootSelectors(), nil)
}

func (rs RedisSentinel) PodLabels() map[string]string {
	return rs.offshootLabels(rs.OffshootSelectors(), rs.Spec.PodTemplate.Labels)
}

func (rs RedisSentinel) PodControllerLabels() map[string]string {
	return rs.offshootLabels(rs.OffshootSelectors(), rs.Spec.PodTemplate.Controller.Labels)
}

func (rs RedisSentinel) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(rs.Spec.ServiceTemplates, alias)
	return rs.offshootLabels(meta_util.OverwriteKeys(rs.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (rs RedisSentinel) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(rs.Labels, override))
}

func (rs RedisSentinel) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralRedisSentinel, kubedb.GroupName)
}

func (rs RedisSentinel) ResourceShortCode() string {
	return ResourceCodeRedisSentinel
}

func (rs RedisSentinel) ResourceKind() string {
	return ResourceKindRedisSentinel
}

func (rs RedisSentinel) ResourceSingular() string {
	return ResourceSingularRedisSentinel
}

func (rs RedisSentinel) ResourcePlural() string {
	return ResourcePluralRedisSentinel
}

func (rs RedisSentinel) GoverningServiceName() string {
	return meta_util.NameWithSuffix(rs.OffshootName(), "pods")
}

func (rs RedisSentinel) ConfigSecretName() string {
	return rs.OffshootName()
}

type redisSentinelApp struct {
	*RedisSentinel
}

func (rs redisSentinelApp) Name() string {
	return rs.RedisSentinel.Name
}

func (rs redisSentinelApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularRedisSentinel))
}

func (rs RedisSentinel) AppBindingMeta() appcat.AppBindingMeta {
	return &redisSentinelApp{&rs}
}

type redisSentinelStatsService struct {
	*RedisSentinel
}

func (rs redisSentinelStatsService) GetNamespace() string {
	return rs.RedisSentinel.GetNamespace()
}

func (rs redisSentinelStatsService) ServiceName() string {
	return rs.OffshootName() + "-stats"
}

func (rs redisSentinelStatsService) ServiceMonitorName() string {
	return rs.ServiceName()
}

func (rs redisSentinelStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return rs.OffshootLabels()
}

func (rs redisSentinelStatsService) Path() string {
	return DefaultStatsPath
}

func (r redisSentinelStatsService) Scheme() string {
	return ""
}

func (rs RedisSentinel) StatsService() mona.StatsAccessor {
	return &redisSentinelStatsService{&rs}
}

func (rs RedisSentinel) StatsServiceLabels() map[string]string {
	return rs.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}

func (rs *RedisSentinel) SetDefaults(topology *core_util.Topology) {
	if rs == nil {
		return
	}

	if rs.Spec.StorageType == "" {
		rs.Spec.StorageType = StorageTypeDurable
	}
	if rs.Spec.TerminationPolicy == "" {
		rs.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if rs.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		rs.Spec.PodTemplate.Spec.ServiceAccountName = rs.OffshootName()
	}

	rs.setDefaultAffinity(&rs.Spec.PodTemplate, rs.OffshootSelectors(), topology)

	rs.Spec.Monitor.SetDefaults()
	rs.SetTLSDefaults()
	SetDefaultResourceLimits(&rs.Spec.PodTemplate.Spec.Resources, DefaultResources)
}

func (rs *RedisSentinel) SetTLSDefaults() {
	if rs.Spec.TLS == nil || rs.Spec.TLS.IssuerRef == nil {
		return
	}
	rs.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(rs.Spec.TLS.Certificates, string(RedisServerCert), rs.CertificateName(RedisServerCert))
	rs.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(rs.Spec.TLS.Certificates, string(RedisClientCert), rs.CertificateName(RedisClientCert))
	rs.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(rs.Spec.TLS.Certificates, string(RedisMetricsExporterCert), rs.CertificateName(RedisMetricsExporterCert))
}

func (r *RedisSentinel) GetPersistentSecrets() []string {
	if r == nil {
		return nil
	}

	var secrets []string
	if r.Spec.AuthSecret != nil {
		secrets = append(secrets, r.Spec.AuthSecret.Name)
	}
	return secrets
}

func (rs *RedisSentinel) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
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
						Namespaces: []string{rs.Namespace},
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
						Namespaces: []string{rs.Namespace},
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

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (rs *RedisSentinel) CertificateName(alias RedisCertificateAlias) string {
	return meta_util.NameWithSuffix(rs.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any provide,
// otherwise returns default certificate secret name for the given alias.
func (rs *RedisSentinel) GetCertSecretName(alias RedisCertificateAlias) string {
	if rs.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(rs.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return rs.CertificateName(alias)
}

func (rs *RedisSentinel) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(rs.Namespace), labels.SelectorFromSet(rs.OffshootLabels()), expectedItems)
}
