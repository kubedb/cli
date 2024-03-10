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
	"strings"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"gomodules.xyz/pointer"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

type SolrApp struct {
	*Solr
}

func (s *Solr) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralSolr))
}

func (s *Solr) StatefulsetName(suffix string) string {
	sts := []string{s.Name}
	if suffix != "" {
		sts = append(sts, suffix)
	}
	return strings.Join(sts, "-")
}

// Owner returns owner reference to resources
func (s *Solr) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(s, SchemeGroupVersion.WithKind(s.ResourceKind()))
}

func (s *Solr) ResourceKind() string {
	return ResourceKindSolr
}

func (s *Solr) GoverningServiceName() string {
	return meta_util.NameWithSuffix(s.ServiceName(), "pods")
}

func (s *Solr) OverseerDiscoveryServiceName() string {
	return meta_util.NameWithSuffix(s.ServiceName(), "overseer")
}

func (s *Solr) ServiceAccountName() string { return s.OffshootName() }

func (s *Solr) DefaultPodRoleName() string {
	return meta_util.NameWithSuffix(s.OffshootName(), "role")
}

func (s *Solr) DefaultPodRoleBindingName() string {
	return meta_util.NameWithSuffix(s.OffshootName(), "rolebinding")
}

func (s *Solr) ServiceName() string {
	return s.OffshootName()
}

func (s *Solr) SolrSecretName(suffix string) string {
	return strings.Join([]string{s.Name, suffix}, "-")
}

func (s *Solr) SolrSecretKey() string {
	return SolrSecretKey
}

func (s *Solr) Merge(opt map[string]string) map[string]string {
	if len(s.Spec.SolrOpts) == 0 {
		return opt
	}
	for _, y := range s.Spec.SolrOpts {
		sr := strings.Split(y, "=")
		_, ok := opt[sr[0]]
		if !ok || sr[0] != "-Dsolr.node.roles" {
			opt[sr[0]] = sr[1]
		}
	}
	return opt
}

func (s *Solr) Append(opt map[string]string) string {
	fl := 0
	as := ""
	for x, y := range opt {
		if fl == 1 {
			as += " "
		}
		as += fmt.Sprintf("%s=%s", x, y)
		fl = 1

	}
	return as
}

func (s *Solr) OffshootName() string {
	return s.Name
}

func (s *Solr) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return s.offshootLabels(meta_util.OverwriteKeys(s.OffshootSelectors(), extraLabels...), s.Spec.PodTemplate.Controller.Labels)
}

func (s *Solr) PodLabels(extraLabels ...map[string]string) map[string]string {
	return s.offshootLabels(meta_util.OverwriteKeys(s.OffshootSelectors(), extraLabels...), s.Spec.PodTemplate.Labels)
}

func (s *Solr) OffshootLabels() map[string]string {
	return s.offshootLabels(s.OffshootSelectors(), nil)
}

func (s *Solr) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, s.Labels, override))
}

func (s *Solr) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      s.ResourceFQN(),
		meta_util.InstanceLabelKey:  s.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (s *Solr) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", s.ResourcePlural(), kubedb.GroupName)
}

func (s *Solr) ResourcePlural() string {
	return ResourcePluralSolr
}

func (s SolrApp) Name() string {
	return s.Solr.Name
}

func (s SolrApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularSolr))
}

func (s *Solr) AppBindingMeta() appcat.AppBindingMeta {
	return &SolrApp{s}
}

func (s *Solr) GetConnectionScheme() string {
	scheme := "http"
	if s.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}

func (s *Solr) PVCName(alias string) string {
	return meta_util.NameWithSuffix(s.Name, alias)
}

func (s *Solr) SetDefaults(slVersion *catalog.SolrVersion) {
	if s.Spec.TerminationPolicy == "" {
		s.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if s.Spec.StorageType == "" {
		s.Spec.StorageType = StorageTypeDurable
	}

	if s.Spec.AuthSecret == nil {
		s.Spec.AuthSecret = &v1.LocalObjectReference{
			Name: s.SolrSecretName("admin-cred"),
		}
	}

	if s.Spec.ZookeeperDigestSecret == nil {
		s.Spec.ZookeeperDigestSecret = &v1.LocalObjectReference{
			Name: s.SolrSecretName("zk-digest"),
		}
	}

	if s.Spec.ZookeeperDigestReadonlySecret == nil {
		s.Spec.ZookeeperDigestReadonlySecret = &v1.LocalObjectReference{
			Name: s.SolrSecretName("zk-digest-readonly"),
		}
	}

	if s.Spec.AuthConfigSecret == nil {
		s.Spec.AuthConfigSecret = &v1.LocalObjectReference{
			Name: s.SolrSecretName("auth-config"),
		}
	}

	if s.Spec.Topology != nil {
		if s.Spec.Topology.Data != nil {
			if s.Spec.Topology.Data.Suffix == "" {
				s.Spec.Topology.Data.Suffix = string(SolrNodeRoleData)
			}
			if s.Spec.Topology.Data.Replicas == nil {
				s.Spec.Topology.Data.Replicas = pointer.Int32P(1)
			}

			if s.Spec.Topology.Data.PodTemplate.Spec.SecurityContext == nil {
				s.Spec.Topology.Data.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{}
			}
			s.Spec.Topology.Data.PodTemplate.Spec.SecurityContext.FSGroup = slVersion.Spec.SecurityContext.RunAsUser
			s.setDefaultContainerSecurityContext(slVersion, &s.Spec.Topology.Data.PodTemplate)
			s.setDefaultContainerResourceLimits(&s.Spec.Topology.Data.PodTemplate)
		}

		if s.Spec.Topology.Overseer != nil {
			if s.Spec.Topology.Overseer.Suffix == "" {
				s.Spec.Topology.Overseer.Suffix = string(SolrNodeRoleOverseer)
			}
			if s.Spec.Topology.Overseer.Replicas == nil {
				s.Spec.Topology.Overseer.Replicas = pointer.Int32P(1)
			}

			if s.Spec.Topology.Overseer.PodTemplate.Spec.SecurityContext == nil {
				s.Spec.Topology.Overseer.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{}
			}
			s.Spec.Topology.Overseer.PodTemplate.Spec.SecurityContext.FSGroup = slVersion.Spec.SecurityContext.RunAsUser
			s.setDefaultContainerSecurityContext(slVersion, &s.Spec.Topology.Overseer.PodTemplate)
			s.setDefaultContainerResourceLimits(&s.Spec.Topology.Overseer.PodTemplate)
		}

		if s.Spec.Topology.Coordinator != nil {
			if s.Spec.Topology.Coordinator.Suffix == "" {
				s.Spec.Topology.Coordinator.Suffix = string(SolrNodeRoleCoordinator)
			}
			if s.Spec.Topology.Coordinator.Replicas == nil {
				s.Spec.Topology.Coordinator.Replicas = pointer.Int32P(1)
			}

			if s.Spec.Topology.Coordinator.PodTemplate.Spec.SecurityContext == nil {
				s.Spec.Topology.Coordinator.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{}
			}
			s.Spec.Topology.Coordinator.PodTemplate.Spec.SecurityContext.FSGroup = slVersion.Spec.SecurityContext.RunAsUser
			s.setDefaultContainerSecurityContext(slVersion, &s.Spec.Topology.Coordinator.PodTemplate)
			s.setDefaultContainerResourceLimits(&s.Spec.Topology.Coordinator.PodTemplate)
		}
	} else {

		if s.Spec.Replicas == nil {
			s.Spec.Replicas = pointer.Int32P(1)
		}

		if s.Spec.PodTemplate.Spec.SecurityContext == nil {
			s.Spec.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{}
		}

		s.Spec.PodTemplate.Spec.SecurityContext.FSGroup = slVersion.Spec.SecurityContext.RunAsUser
		s.setDefaultContainerSecurityContext(slVersion, &s.Spec.PodTemplate)
		s.setDefaultContainerResourceLimits(&s.Spec.PodTemplate)
	}
}

func (s *Solr) setDefaultContainerSecurityContext(slVersion *catalog.SolrVersion, podTemplate *ofst.PodTemplateSpec) {
	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, SolrInitContainerName)
	if initContainer == nil {
		initContainer = &v1.Container{
			Name: SolrInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &v1.SecurityContext{}
	}
	s.assignDefaultContainerSecurityContext(slVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, SolrContainerName)
	if container == nil {
		container = &v1.Container{
			Name: SolrContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &v1.SecurityContext{}
	}
	s.assignDefaultContainerSecurityContext(slVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (s *Solr) assignDefaultContainerSecurityContext(slVersion *catalog.SolrVersion, sc *v1.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &v1.Capabilities{
			Drop: []v1.Capability{"ALL"},
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(true)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = slVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (s *Solr) setDefaultContainerResourceLimits(podTemplate *ofst.PodTemplateSpec) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, SolrContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, DefaultResourcesMemoryIntensive)
	}

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, SolrInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, DefaultInitContainerResource)
	}
}

func (s *Solr) SetHealthCheckerDefaults() {
	if s.Spec.HealthChecker.PeriodSeconds == nil {
		s.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(20)
	}
	if s.Spec.HealthChecker.TimeoutSeconds == nil {
		s.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if s.Spec.HealthChecker.FailureThreshold == nil {
		s.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (s *Solr) GetPersistentSecrets() []string {
	if s == nil {
		return nil
	}

	var secrets []string
	// Add Admin/Elastic user secret name
	if s.Spec.AuthSecret != nil {
		secrets = append(secrets, s.Spec.AuthSecret.Name)
	}

	if s.Spec.AuthConfigSecret != nil {
		secrets = append(secrets, s.Spec.AuthConfigSecret.Name)
	}

	if s.Spec.ZookeeperDigestSecret != nil {
		secrets = append(secrets, s.Spec.ZookeeperDigestSecret.Name)
	}

	if s.Spec.ZookeeperDigestReadonlySecret != nil {
		secrets = append(secrets, s.Spec.ZookeeperDigestReadonlySecret.Name)
	}

	return secrets
}

func (s *Solr) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	if s.Spec.Topology != nil {
		expectedItems = 3
	}
	return checkReplicas(lister.StatefulSets(s.Namespace), labels.SelectorFromSet(s.OffshootLabels()), expectedItems)
}
