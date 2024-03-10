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
	"strings"

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
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	metautil "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

func (s *Singlestore) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralSinglestore))
}

type singlestoreApp struct {
	*Singlestore
}

func (s *singlestoreApp) Name() string {
	return s.Singlestore.Name
}

func (s *singlestoreApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularSinglestore))
}

func (s *Singlestore) ResourceShortCode() string {
	return ResourceCodeSinglestore
}

func (s *Singlestore) ResourceKind() string {
	return ResourceKindSinglestore
}

func (s *Singlestore) ResourceSingular() string {
	return ResourceSingularRabbitmq
}

func (s *Singlestore) ResourcePlural() string {
	return ResourcePluralSinglestore
}

func (s *Singlestore) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", s.ResourcePlural(), kubedb.GroupName)
}

// Owner returns owner reference to resources
func (s *Singlestore) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(s, SchemeGroupVersion.WithKind(s.ResourceKind()))
}

type singlestoreStatsService struct {
	*Singlestore
}

func (s singlestoreStatsService) GetNamespace() string {
	return s.Singlestore.GetNamespace()
}

func (s Singlestore) GetNameSpacedName() string {
	return s.Namespace + "/" + s.Name
}

func (s singlestoreStatsService) ServiceName() string {
	return s.OffshootName() + "-stats"
}

func (s singlestoreStatsService) ServiceMonitorName() string {
	return s.ServiceName()
}

func (s singlestoreStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return s.OffshootLabels()
}

func (s singlestoreStatsService) Path() string {
	return DefaultStatsPath
}

func (s singlestoreStatsService) Scheme() string {
	return ""
}

func (s singlestoreStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (s Singlestore) StatsService() mona.StatsAccessor {
	return &singlestoreStatsService{&s}
}

func (s Singlestore) StatsServiceLabels() map[string]string {
	return s.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}

func (s *Singlestore) OffshootName() string {
	return s.Name
}

func (s *Singlestore) ServiceName() string {
	return s.OffshootName()
}

func (s *Singlestore) AppBindingMeta() appcat.AppBindingMeta {
	return &singlestoreApp{s}
}

func (s *Singlestore) GoverningServiceName() string {
	return metautil.NameWithSuffix(s.ServiceName(), "pods")
}

func (s *Singlestore) PrimaryServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", s.ServiceName(), s.Namespace)
}

func (s *Singlestore) DefaultUserCredSecretName(username string) string {
	return metautil.NameWithSuffix(s.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", username), "_", "-"))
}

func (s *Singlestore) offshootLabels(selector, override map[string]string) map[string]string {
	selector[metautil.ComponentLabelKey] = ComponentDatabase
	return metautil.FilterKeys(kubedb.GroupName, selector, metautil.OverwriteKeys(nil, s.Labels, override))
}

func (s *Singlestore) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(s.Spec.ServiceTemplates, alias)
	return s.offshootLabels(metautil.OverwriteKeys(s.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (s *Singlestore) OffshootLabels() map[string]string {
	return s.offshootLabels(s.OffshootSelectors(), nil)
}

func (s *Singlestore) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		metautil.NameLabelKey:      s.ResourceFQN(),
		metautil.InstanceLabelKey:  s.Name,
		metautil.ManagedByLabelKey: kubedb.GroupName,
	}
	return metautil.OverwriteKeys(selector, extraSelectors...)
}

func (s *Singlestore) IsClustering() bool {
	return s.Spec.Topology != nil
}

func (s *Singlestore) IsStandalone() bool {
	return s.Spec.Topology == nil
}

func (s *Singlestore) PVCName(alias string) string {
	return metautil.NameWithSuffix(s.OffshootName(), alias)
	// return s.OffshootName()
}

func (s *Singlestore) AggregatorStatefulSet() string {
	sts := s.OffshootName()
	if s.Spec.Topology.Aggregator.Suffix != "" {
		sts = metautil.NameWithSuffix(sts, s.Spec.Topology.Aggregator.Suffix)
	}
	return metautil.NameWithSuffix(sts, StatefulSetTypeAggregator)
}

func (s *Singlestore) LeafStatefulSet() string {
	sts := s.OffshootName()
	if s.Spec.Topology.Leaf.Suffix != "" {
		sts = metautil.NameWithSuffix(sts, s.Spec.Topology.Leaf.Suffix)
	}
	return metautil.NameWithSuffix(sts, StatefulSetTypeLeaf)
}

func (s *Singlestore) PodLabels(extraLabels ...map[string]string) map[string]string {
	return s.offshootLabels(metautil.OverwriteKeys(s.OffshootSelectors(), extraLabels...), s.Spec.PodTemplate.Labels)
}

func (s *Singlestore) PodLabel(podTemplate *ofst.PodTemplateSpec) map[string]string {
	if podTemplate != nil && podTemplate.Labels != nil {
		return s.offshootLabels(s.OffshootSelectors(), s.Spec.PodTemplate.Labels)
	}
	return s.offshootLabels(s.OffshootSelectors(), nil)
}

func (s *Singlestore) ConfigSecretName() string {
	return metautil.NameWithSuffix(s.OffshootName(), "config")
}

func (s *Singlestore) StatefulSetName() string {
	return s.OffshootName()
}

func (s *Singlestore) ServiceAccountName() string {
	return s.OffshootName()
}

func (s *Singlestore) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return s.offshootLabels(metautil.OverwriteKeys(s.OffshootSelectors(), extraLabels...), s.Spec.PodTemplate.Controller.Labels)
}

func (s *Singlestore) PodControllerLabel(podTemplate *ofst.PodTemplateSpec) map[string]string {
	if podTemplate != nil && podTemplate.Controller.Labels != nil {
		return s.offshootLabels(s.OffshootSelectors(), podTemplate.Controller.Labels)
	}
	return s.offshootLabels(s.OffshootSelectors(), nil)
}

func (s *Singlestore) SetHealthCheckerDefaults() {
	if s.Spec.HealthChecker.PeriodSeconds == nil {
		s.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if s.Spec.HealthChecker.TimeoutSeconds == nil {
		s.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if s.Spec.HealthChecker.FailureThreshold == nil {
		s.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (s *Singlestore) GetAuthSecretName() string {
	if s.Spec.AuthSecret != nil && s.Spec.AuthSecret.Name != "" {
		return s.Spec.AuthSecret.Name
	}
	return metautil.NameWithSuffix(s.OffshootName(), "auth")
}

func (s *Singlestore) GetPersistentSecrets() []string {
	var secrets []string
	if s.Spec.AuthSecret != nil {
		secrets = append(secrets, s.Spec.AuthSecret.Name)
	}
	return secrets
}

func (s *Singlestore) SetDefaults() {
	if s == nil {
		return
	}
	if s.Spec.StorageType == "" {
		s.Spec.StorageType = StorageTypeDurable
	}
	if s.Spec.TerminationPolicy == "" {
		s.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if s.Spec.Topology == nil {
		if s.Spec.Replicas == nil {
			s.Spec.Replicas = pointer.Int32P(1)
		}
		if s.Spec.PodTemplate == nil {
			s.Spec.PodTemplate = &ofst.PodTemplateSpec{}
		}
	} else {
		if s.Spec.Topology.Aggregator.Replicas == nil {
			s.Spec.Topology.Aggregator.Replicas = pointer.Int32P(3)
		}

		if s.Spec.Topology.Leaf.Replicas == nil {
			s.Spec.Topology.Leaf.Replicas = pointer.Int32P(2)
		}
		if s.Spec.Topology.Aggregator.PodTemplate == nil {
			s.Spec.Topology.Aggregator.PodTemplate = &ofst.PodTemplateSpec{}
		}
		if s.Spec.Topology.Leaf.PodTemplate == nil {
			s.Spec.Topology.Leaf.PodTemplate = &ofst.PodTemplateSpec{}
		}
	}

	var sdbVersion catalog.SinglestoreVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name: s.Spec.Version,
	}, &sdbVersion)
	if err != nil {
		klog.Errorf("can't get the singlestore version object %s for %s \n", err.Error(), s.Spec.Version)
		return
	}

	if s.IsStandalone() {
		s.setDefaultContainerSecurityContext(&sdbVersion, s.Spec.PodTemplate)
	} else {
		s.setDefaultContainerSecurityContext(&sdbVersion, s.Spec.Topology.Aggregator.PodTemplate)
		s.setDefaultContainerSecurityContext(&sdbVersion, s.Spec.Topology.Leaf.PodTemplate)
	}

	s.SetTLSDefaults()

	s.SetHealthCheckerDefaults()
	if s.Spec.Monitor != nil {
		if s.Spec.Monitor.Prometheus == nil {
			s.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if s.Spec.Monitor.Prometheus != nil && s.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			s.Spec.Monitor.Prometheus.Exporter.Port = SinglestoreExporterPort
		}
		s.Spec.Monitor.SetDefaults()
	}

	if s.IsClustering() {
		s.setDefaultContainerResourceLimits(s.Spec.Topology.Aggregator.PodTemplate)
		s.setDefaultContainerResourceLimits(s.Spec.Topology.Leaf.PodTemplate)
	} else {
		s.setDefaultContainerResourceLimits(s.Spec.PodTemplate)
	}
}

func (s *Singlestore) setDefaultContainerSecurityContext(sdbVersion *catalog.SinglestoreVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = sdbVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, SinglestoreContainerName)
	if container == nil {
		container = &core.Container{
			Name: SinglestoreContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	s.assignDefaultContainerSecurityContext(sdbVersion, container.SecurityContext)

	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, SinglestoreInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: SinglestoreInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	s.assignDefaultInitContainerSecurityContext(sdbVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)

	if s.IsClustering() {
		coordinatorContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, SinglestoreCoordinatorContainerName)
		if coordinatorContainer == nil {
			coordinatorContainer = &core.Container{
				Name: SinglestoreCoordinatorContainerName,
			}
		}
		if coordinatorContainer.SecurityContext == nil {
			coordinatorContainer.SecurityContext = &core.SecurityContext{}
		}
		s.assignDefaultContainerSecurityContext(sdbVersion, coordinatorContainer.SecurityContext)
		podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *coordinatorContainer)
	}
}

func (s *Singlestore) assignDefaultInitContainerSecurityContext(sdbVersion *catalog.SinglestoreVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = sdbVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = sdbVersion.Spec.SecurityContext.RunAsGroup
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (s *Singlestore) assignDefaultContainerSecurityContext(sdbVersion *catalog.SinglestoreVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = sdbVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = sdbVersion.Spec.SecurityContext.RunAsGroup
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (s *Singlestore) setDefaultContainerResourceLimits(podTemplate *ofst.PodTemplateSpec) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, SinglestoreContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, DefaultResourcesMemoryIntensive)
	}

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, SinglestoreInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, DefaultInitContainerResource)
	}

	if s.IsClustering() {
		coordinatorContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, SinglestoreCoordinatorContainerName)
		if coordinatorContainer != nil && (coordinatorContainer.Resources.Requests == nil && coordinatorContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&coordinatorContainer.Resources, CoordinatorDefaultResources)
		}
	}
}

func (s *Singlestore) SetTLSDefaults() {
	if s.Spec.TLS == nil || s.Spec.TLS.IssuerRef == nil {
		return
	}
	s.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(s.Spec.TLS.Certificates, string(SinglestoreServerCert), s.CertificateName(SinglestoreServerCert))
	s.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(s.Spec.TLS.Certificates, string(SinglestoreClientCert), s.CertificateName(SinglestoreClientCert))
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (s *Singlestore) CertificateName(alias SinglestoreCertificateAlias) string {
	return metautil.NameWithSuffix(s.Name, fmt.Sprintf("%s-cert", string(alias)))
}

func (s *Singlestore) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	if s.Spec.Topology != nil {
		expectedItems = 2
	}
	return checkReplicas(lister.StatefulSets(s.Namespace), labels.SelectorFromSet(s.OffshootLabels()), expectedItems)
}
