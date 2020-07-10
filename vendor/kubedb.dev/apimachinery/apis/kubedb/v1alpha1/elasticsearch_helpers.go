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
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/apiextensions"
	core_util "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ElasticsearchNodeAffinityTemplateVar = "NODE_ROLE"
)

func (_ Elasticsearch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralElasticsearch))
}

var _ apis.ResourceInfo = &Elasticsearch{}

func (e Elasticsearch) OffshootName() string {
	return e.Name
}

func (e Elasticsearch) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseKind: ResourceKindElasticsearch,
		LabelDatabaseName: e.Name,
	}
}

func (e Elasticsearch) OffshootLabels() map[string]string {
	out := e.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularElasticsearch
	out[meta_util.VersionLabelKey] = string(e.Spec.Version)
	out[meta_util.InstanceLabelKey] = e.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, e.Labels)
}

func (e Elasticsearch) ResourceShortCode() string {
	return ResourceCodeElasticsearch
}

func (e Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (e Elasticsearch) ResourceSingular() string {
	return ResourceSingularElasticsearch
}

func (e Elasticsearch) ResourcePlural() string {
	return ResourcePluralElasticsearch
}

func (e Elasticsearch) ServiceName() string {
	return e.OffshootName()
}

func (e *Elasticsearch) MasterServiceName() string {
	return fmt.Sprintf("%v-master", e.ServiceName())
}

// Governing Service Name
func (e Elasticsearch) GvrSvcName() string {
	return e.OffshootName() + "-gvr"
}

func (e *Elasticsearch) GetConnectionScheme() string {
	scheme := "http"
	if e.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}

func (e *Elasticsearch) GetConnectionURL() string {
	return fmt.Sprintf("%v://%s.%s:%d", e.GetConnectionScheme(), e.OffshootName(), e.Namespace, ElasticsearchRestPort)
}

type elasticsearchApp struct {
	*Elasticsearch
}

func (r elasticsearchApp) Name() string {
	return r.Elasticsearch.Name
}

func (r elasticsearchApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularElasticsearch))
}

func (e Elasticsearch) AppBindingMeta() appcat.AppBindingMeta {
	return &elasticsearchApp{&e}
}

type elasticsearchStatsService struct {
	*Elasticsearch
}

func (e elasticsearchStatsService) GetNamespace() string {
	return e.Elasticsearch.GetNamespace()
}

func (e elasticsearchStatsService) ServiceName() string {
	return e.OffshootName() + "-stats"
}

func (e elasticsearchStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", e.Namespace, e.Name)
}

func (e elasticsearchStatsService) Path() string {
	return DefaultStatsPath
}

func (e elasticsearchStatsService) Scheme() string {
	return ""
}

func (e Elasticsearch) StatsService() mona.StatsAccessor {
	return &elasticsearchStatsService{&e}
}

func (e Elasticsearch) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, e.OffshootSelectors(), e.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (e *Elasticsearch) GetMonitoringVendor() string {
	if e.Spec.Monitor != nil {
		return e.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (e *Elasticsearch) SetDefaults(topology *core_util.Topology) {
	if e == nil {
		return
	}
	if !e.Spec.DisableSecurity && e.Spec.AuthPlugin == v1alpha1.ElasticsearchAuthPluginNone {
		e.Spec.DisableSecurity = true
	}
	e.Spec.AuthPlugin = ""
	if e.Spec.StorageType == "" {
		e.Spec.StorageType = StorageTypeDurable
	}
	if e.Spec.UpdateStrategy.Type == "" {
		e.Spec.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if e.Spec.TerminationPolicy == "" {
		e.Spec.TerminationPolicy = TerminationPolicyDelete
	} else if e.Spec.TerminationPolicy == TerminationPolicyPause {
		e.Spec.TerminationPolicy = TerminationPolicyHalt
	}

	if e.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		e.Spec.PodTemplate.Spec.ServiceAccountName = e.OffshootName()
	}

	e.setDefaultAffinity(&e.Spec.PodTemplate, e.OffshootSelectors(), topology)

	e.Spec.Monitor.SetDefaults()
}

// setDefaultAffinity
func (e *Elasticsearch) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
	if podTemplate == nil {
		return
	} else if podTemplate.Spec.Affinity != nil {
		// Update topologyKey fields according to Kubernetes version
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
						Namespaces: []string{e.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels:      labels,
							MatchExpressions: e.GetMatchExpressions(),
						},

						TopologyKey: corev1.LabelHostname,
					},
				},
				// Prefer to not schedule multiple pods on the node with same zone
				{
					Weight: 50,
					PodAffinityTerm: corev1.PodAffinityTerm{
						Namespaces: []string{e.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels:      labels,
							MatchExpressions: e.GetMatchExpressions(),
						},
						TopologyKey: topology.LabelZone,
					},
				},
			},
		},
	}
}

func (e *Elasticsearch) GetMatchExpressions() []metav1.LabelSelectorRequirement {
	if e.Spec.Topology == nil {
		return nil
	}

	return []metav1.LabelSelectorRequirement{
		{
			Key:      fmt.Sprintf("${%s}", ElasticsearchNodeAffinityTemplateVar),
			Operator: metav1.LabelSelectorOpExists,
		},
	}
}

func (e *ElasticsearchSpec) GetSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	if e.DatabaseSecret != nil {
		secrets = append(secrets, e.DatabaseSecret.SecretName)
	}
	if e.CertificateSecret != nil {
		secrets = append(secrets, e.CertificateSecret.SecretName)
	}
	return secrets
}
