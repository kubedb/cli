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

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	core_util "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
)

func (p *Pgpool) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgpool))
}

func (p *Pgpool) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", p.ResourcePlural(), kubedb.GroupName)
}

func (p *Pgpool) ResourceShortCode() string {
	return ResourceCodePgpool
}

func (p *Pgpool) ResourceKind() string {
	return ResourceKindPgpool
}

func (p *Pgpool) ResourceSingular() string {
	return ResourceSingularPgpool
}

func (p *Pgpool) ResourcePlural() string {
	return ResourcePluralPgpool
}

func (p *Pgpool) ConfigSecretName() string {
	return meta_util.NameWithSuffix(p.OffshootName(), "config")
}

func (p *Pgpool) ServiceAccountName() string {
	return p.OffshootName()
}

func (p *Pgpool) GoverningServiceName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "pods")
}

func (p *Pgpool) ServiceName() string {
	return p.OffshootName()
}

// Owner returns owner reference to resources
func (p *Pgpool) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(p, SchemeGroupVersion.WithKind(p.ResourceKind()))
}

func (p *Pgpool) PodLabels(extraLabels ...map[string]string) map[string]string {
	var labels map[string]string
	if p.Spec.PodTemplate != nil {
		return p.offshootLabels(meta_util.OverwriteKeys(p.OffshootSelectors(), extraLabels...), p.Spec.PodTemplate.Labels)
	} else {
		return p.offshootLabels(meta_util.OverwriteKeys(p.OffshootSelectors(), extraLabels...), labels)
	}
}

func (p *Pgpool) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	var labels map[string]string
	if p.Spec.PodTemplate != nil {
		return p.offshootLabels(meta_util.OverwriteKeys(p.OffshootSelectors(), extraLabels...), p.Spec.PodTemplate.Controller.Labels)
	} else {
		return p.offshootLabels(meta_util.OverwriteKeys(p.OffshootSelectors(), extraLabels...), labels)
	}
}

func (p *Pgpool) OffshootLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), nil)
}

func (p *Pgpool) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentConnectionPooler
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, p.Labels, override))
}

func (p *Pgpool) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      p.ResourceFQN(),
		meta_util.InstanceLabelKey:  p.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (p *Pgpool) PetSetName() string {
	return p.OffshootName()
}

func (p *Pgpool) OffshootName() string {
	return p.Name
}

func (p *Pgpool) GetAuthSecretName() string {
	if p.Spec.AuthSecret != nil && p.Spec.AuthSecret.Name != "" {
		return p.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(p.OffshootName(), "auth")
}

func (p *Pgpool) SetHealthCheckerDefaults() {
	if p.Spec.HealthChecker.PeriodSeconds == nil {
		p.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if p.Spec.HealthChecker.TimeoutSeconds == nil {
		p.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if p.Spec.HealthChecker.FailureThreshold == nil {
		p.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

// PrimaryServiceDNS make primary host dns with require template
func (p *Pgpool) PrimaryServiceDNS() string {
	return fmt.Sprintf("%v.%v.svc", p.ServiceName(), p.Namespace)
}

func (p *Pgpool) GetNameSpacedName() string {
	return p.Namespace + "/" + p.Name
}

type PgpoolStatsService struct {
	*Pgpool
}

func (p PgpoolStatsService) GetNamespace() string {
	return p.Pgpool.GetNamespace()
}

func (p PgpoolStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p PgpoolStatsService) ServiceMonitorName() string {
	return p.ServiceName()
}

func (p PgpoolStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return p.OffshootLabels()
}

func (p PgpoolStatsService) Path() string {
	return DefaultStatsPath
}

func (p PgpoolStatsService) Scheme() string {
	return ""
}

func (p PgpoolStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (p Pgpool) StatsService() mona.StatsAccessor {
	return &PgpoolStatsService{&p}
}

func (p Pgpool) StatsServiceLabels() map[string]string {
	return p.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}

func (p *Pgpool) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(p.Spec.ServiceTemplates, alias)
	return p.offshootLabels(meta_util.OverwriteKeys(p.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (p *Pgpool) SetSecurityContext(ppVersion *catalog.PgpoolVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = ppVersion.Spec.SecurityContext.RunAsUser
	}

	container := core_util.GetContainerByName(podTemplate.Spec.Containers, PgpoolContainerName)
	if container == nil {
		container = &core.Container{
			Name: PgpoolContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	p.assignContainerSecurityContext(ppVersion, container.SecurityContext)
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (p *Pgpool) assignContainerSecurityContext(ppVersion *catalog.PgpoolVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = ppVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = ppVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (p *Pgpool) setContainerResourceLimits(podTemplate *ofst.PodTemplateSpec) {
	ppContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, PgpoolContainerName)
	if ppContainer != nil && (ppContainer.Resources.Requests == nil && ppContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&ppContainer.Resources, DefaultResources)
	}
}

func (p *Pgpool) SetDefaults() {
	if p == nil {
		return
	}
	if p.Spec.Replicas == nil {
		p.Spec.Replicas = pointer.Int32P(1)
	}
	if p.Spec.TerminationPolicy == "" {
		p.Spec.TerminationPolicy = TerminationPolicyDelete
	}
	if p.Spec.PodTemplate == nil {
		p.Spec.PodTemplate = &ofst.PodTemplateSpec{}
		p.Spec.PodTemplate.Spec.Containers = []core.Container{}
	}

	ppVersion := catalog.PgpoolVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name: p.Spec.Version,
	}, &ppVersion)
	if err != nil {
		klog.Errorf("can't get the pgpool version object %s for %s \n", err.Error(), p.Spec.Version)
		return
	}

	if p.Spec.Monitor != nil {
		if p.Spec.Monitor.Prometheus == nil {
			p.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if p.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			p.Spec.Monitor.Prometheus.Exporter.Port = PgpoolMonitoringDefaultServicePort
		}
		p.Spec.Monitor.SetDefaults()
		if p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = ppVersion.Spec.SecurityContext.RunAsUser
		}
		if p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser
		}
	}

	p.SetHealthCheckerDefaults()
	p.SetSecurityContext(&ppVersion, p.Spec.PodTemplate)
	p.setContainerResourceLimits(p.Spec.PodTemplate)
}

func (p *Pgpool) GetPersistentSecrets() []string {
	var secrets []string
	if p.Spec.AuthSecret != nil {
		secrets = append(secrets, p.Spec.AuthSecret.Name)
		secrets = append(secrets, p.ConfigSecretName())
	}
	return secrets
}

func (p *Pgpool) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicasOfPetSet(lister.PetSets(p.Namespace), labels.SelectorFromSet(p.OffshootLabels()), expectedItems)
}
