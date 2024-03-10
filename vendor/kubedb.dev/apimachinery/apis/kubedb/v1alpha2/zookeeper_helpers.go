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
	appslister "k8s.io/client-go/listers/apps/v1"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

func (z *ZooKeeper) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralZooKeeper))
}

// Owner returns owner reference to resources
func (z *ZooKeeper) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(z, SchemeGroupVersion.WithKind(z.ResourceKind()))
}

func (z *ZooKeeper) OffshootName() string {
	return z.Name
}

func (z *ZooKeeper) ResourceKind() string {
	return ResourceKindZooKeeper
}

func (z *ZooKeeper) ResourceShortCode() string {
	return ResourceCodeZooKeeper
}

func (z *ZooKeeper) ResourceSingular() string {
	return ResourceSingularZooKeeper
}

func (z *ZooKeeper) ResourcePlural() string {
	return ResourcePluralZooKeeper
}

func (z *ZooKeeper) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", z.ResourcePlural(), kubedb.GroupName)
}

func (z *ZooKeeper) StatefulSetName() string {
	return z.OffshootName()
}

func (z *ZooKeeper) ServiceAccountName() string {
	return z.OffshootName()
}

func (z *ZooKeeper) ConfigSecretName() string {
	return meta_util.NameWithSuffix(z.OffshootName(), "config")
}

func (z *ZooKeeper) PVCName(alias string) string {
	return meta_util.NameWithSuffix(z.Name, alias)
}

func (z *ZooKeeper) ServiceName() string {
	return z.OffshootName()
}

func (z *ZooKeeper) AdminServerServiceName() string {
	return fmt.Sprintf("%s-admin-server", z.ServiceName())
}

func (z *ZooKeeper) GoverningServiceName() string {
	return meta_util.NameWithSuffix(z.ServiceName(), "pods")
}

func (z *ZooKeeper) Address() string {
	return fmt.Sprintf("%v.%v.svc:%d", z.ServiceName(), z.Namespace, ZooKeeperClientPort)
}

func (z *ZooKeeper) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      z.ResourceFQN(),
		meta_util.InstanceLabelKey:  z.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (z *ZooKeeper) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, z.Labels, override))
}

func (z *ZooKeeper) OffshootLabels() map[string]string {
	return z.offshootLabels(z.OffshootSelectors(), nil)
}

func (z *ZooKeeper) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return z.offshootLabels(meta_util.OverwriteKeys(z.OffshootSelectors(), extraLabels...), z.Spec.PodTemplate.Controller.Labels)
}

func (z *ZooKeeper) PodLabels(extraLabels ...map[string]string) map[string]string {
	return z.offshootLabels(meta_util.OverwriteKeys(z.OffshootSelectors(), extraLabels...), z.Spec.PodTemplate.Labels)
}

func (z *ZooKeeper) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(z.Spec.ServiceTemplates, alias)
	return z.offshootLabels(meta_util.OverwriteKeys(z.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (z *ZooKeeper) GetAuthSecretName() string {
	if z.Spec.AuthSecret != nil && z.Spec.AuthSecret.Name != "" {
		return z.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(z.OffshootName(), "auth")
}

func (z *ZooKeeper) GetPersistentSecrets() []string {
	if z == nil {
		return nil
	}

	var secrets []string
	if z.Spec.AuthSecret != nil {
		secrets = append(secrets, z.Spec.AuthSecret.Name)
	}
	return secrets
}

func (z *ZooKeeper) SetHealthCheckerDefaults() {
	if z.Spec.HealthChecker.PeriodSeconds == nil {
		z.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if z.Spec.HealthChecker.TimeoutSeconds == nil {
		z.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if z.Spec.HealthChecker.FailureThreshold == nil {
		z.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (z *ZooKeeper) SetDefaults() {
	if z.Spec.TerminationPolicy == "" {
		z.Spec.TerminationPolicy = TerminationPolicyDelete
	}
	if z.Spec.Replicas == nil {
		z.Spec.Replicas = pointer.Int32P(1)
	}

	if z.Spec.Halted {
		if z.Spec.TerminationPolicy == TerminationPolicyDoNotTerminate {
			klog.Errorf(`Can't halt, since termination policy is 'DoNotTerminate'`)
			return
		}
		z.Spec.TerminationPolicy = TerminationPolicyHalt
	}

	var zkVersion catalog.ZooKeeperVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: z.Spec.Version}, &zkVersion)
	if err != nil {
		klog.Errorf("can't get the zookeeper version object %s for %s \n", err.Error(), z.Spec.Version)
		return
	}

	z.setDefaultContainerSecurityContext(&zkVersion, &z.Spec.PodTemplate)

	dbContainer := coreutil.GetContainerByName(z.Spec.PodTemplate.Spec.Containers, ZooKeeperContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, DefaultResources)
	}

	initContainer := coreutil.GetContainerByName(z.Spec.PodTemplate.Spec.InitContainers, ZooKeeperInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, DefaultInitContainerResource)
	}

	z.SetHealthCheckerDefaults()
	if z.Spec.Monitor != nil {
		if z.Spec.Monitor.Prometheus == nil {
			z.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if z.Spec.Monitor.Prometheus != nil && z.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			z.Spec.Monitor.Prometheus.Exporter.Port = ZooKeeperMetricsPort
		}
		z.Spec.Monitor.SetDefaults()
	}
}

func (z *ZooKeeper) setDefaultContainerSecurityContext(zkVersion *catalog.ZooKeeperVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = zkVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, ZooKeeperContainerName)
	if container == nil {
		container = &core.Container{
			Name: ZooKeeperContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}

	z.assignDefaultContainerSecurityContext(zkVersion, container.SecurityContext)

	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, ZooKeeperInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: ZooKeeperInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	z.assignDefaultContainerSecurityContext(zkVersion, initContainer.SecurityContext)

	podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)
}

func (z *ZooKeeper) assignDefaultContainerSecurityContext(zkVersion *catalog.ZooKeeperVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = zkVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = zkVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

type zookeeperStatsService struct {
	*ZooKeeper
}

func (z zookeeperStatsService) GetNamespace() string {
	return z.ZooKeeper.GetNamespace()
}

func (z zookeeperStatsService) ServiceName() string {
	return z.OffshootName() + "-stats"
}

func (z zookeeperStatsService) ServiceMonitorName() string {
	return z.ServiceName()
}

func (z zookeeperStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return z.OffshootLabels()
}

func (z zookeeperStatsService) Path() string {
	return DefaultStatsPath
}

func (z zookeeperStatsService) Scheme() string {
	return ""
}

func (z zookeeperStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (z *ZooKeeper) StatsService() mona.StatsAccessor {
	return &zookeeperStatsService{z}
}

func (z *ZooKeeper) StatsServiceLabels() map[string]string {
	return z.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}

type ZooKeeperApp struct {
	*ZooKeeper
}

func (z ZooKeeperApp) Name() string {
	return z.ZooKeeper.Name
}

func (z ZooKeeperApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularZooKeeper))
}

func (z *ZooKeeper) AppBindingMeta() appcat.AppBindingMeta {
	return &ZooKeeperApp{z}
}

func (z *ZooKeeper) GetConnectionScheme() string {
	scheme := "http"
	//if z.Spec.EnableSSL {
	//	scheme = "https"
	//}
	return scheme
}

func (z *ZooKeeper) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(z.Namespace), labels.SelectorFromSet(z.OffshootLabels()), expectedItems)
}
