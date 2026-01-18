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

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MilvusApp struct {
	*Milvus
}

func (Milvus) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMilvus))
}

func (m *Milvus) ResourceKind() string {
	return ResourceKindMilvus
}

func (m *Milvus) ResourceSingular() string {
	return ResourceSingularMilvus
}

func (m *Milvus) ResourcePlural() string {
	return ResourcePluralMilvus
}

func (m *Milvus) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", m.ResourcePlural(), kubedb.GroupName)
}

func (m *Milvus) AppBindingMeta() appcat.AppBindingMeta {
	return &MilvusApp{m}
}

func (r MilvusApp) Name() string {
	return r.Milvus.Name
}

func (m Milvus) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, m.ResourceSingular()))
}

func (m *Milvus) GetConnectionScheme() string {
	scheme := "http"
	return scheme
}

func (m *Milvus) Owner() *metav1.OwnerReference {
	return metav1.NewControllerRef(m, SchemeGroupVersion.WithKind(m.ResourceKind()))
}

func (m *Milvus) OffshootName() string {
	return m.Name
}

func (m *Milvus) ServiceName() string {
	return m.OffshootName()
}

func (m *Milvus) PetSetName() string {
	return m.OffshootName()
}

func (m *Milvus) ServiceAccountName() string {
	return m.OffshootName()
}

func (m *Milvus) GetAuthSecretName() string {
	if m.Spec.AuthSecret != nil && m.Spec.AuthSecret.Name != "" {
		return m.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(m.OffshootName(), "auth")
}

func (m *Milvus) ConfigSecretName() string {
	uid := string(m.UID)
	return meta_util.NameWithSuffix(m.OffshootName(), uid[len(uid)-6:])
}

func (m *Milvus) GetPersistentSecrets() []string {
	var secrets []string
	if m.Spec.AuthSecret != nil {
		secrets = append(secrets, m.GetAuthSecretName())
	}
	secrets = append(secrets, m.ConfigSecretName())
	return secrets
}

func (m *Milvus) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, m.Labels, override))
}

func (m *Milvus) OffshootLabels() map[string]string {
	return m.offshootLabels(m.OffshootSelectors(), nil)
}

func (m *Milvus) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      m.ResourceFQN(),
		meta_util.InstanceLabelKey:  m.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (m *Milvus) PodLabel(podTemplate *ofstv2.PodTemplateSpec) map[string]string {
	if podTemplate != nil && podTemplate.Labels != nil {
		return m.offshootLabels(m.OffshootSelectors(), m.Spec.PodTemplate.Labels)
	}
	return m.offshootLabels(m.OffshootSelectors(), nil)
}

func (m *Milvus) ServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", m.ServiceName(), m.Namespace, kubedb.MilvusGrpcPort)
}

func (m *Milvus) SetHealthCheckerDefaults() {
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

func (m *Milvus) EtcdServiceName() string {
	return fmt.Sprintf("%s-%s", m.Namespace, kubedb.EtcdName)
}

func (m *Milvus) MetaStorageEndpoints() []string {
	if m.Spec.MetaStorage.ExternallyManaged {
		if len(m.Spec.MetaStorage.Endpoints) == 0 {
			klog.Errorf("metadata storage is externally managed but no endpoints were provided")
			return []string{}
		}
		return m.Spec.MetaStorage.Endpoints
	}

	size := m.Spec.MetaStorage.Size

	endpoints := make([]string, size)
	for i := 0; i < size; i++ {
		// Use pod DNS names for the etcd cluster
		endpoints[i] = fmt.Sprintf(
			"http://%s-%d.%s.%s.svc.cluster.local:%d",
			m.EtcdServiceName(), i,
			m.EtcdServiceName(), m.Namespace,
			2379,
		)
	}

	return endpoints
}

func (m *Milvus) SetDefaults(kc client.Client) {
	if m.Spec.Topology.Mode == nil {
		mode := MilvusMode("Standalone")
		m.Spec.Topology.Mode = &mode
	}

	if m.Spec.DeletionPolicy == "" {
		m.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}

	if m.Spec.AuthSecret == nil {
		m.Spec.AuthSecret = &SecretReference{}
	}

	if m.Spec.AuthSecret.Kind == "" {
		m.Spec.AuthSecret.Kind = kubedb.ResourceKindSecret
	}

	if m.Spec.PodTemplate == nil {
		m.Spec.PodTemplate = &ofstv2.PodTemplateSpec{}
	}

	var mvVersion catalog.MilvusVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name: m.Spec.Version,
	}, &mvVersion)
	if err != nil {
		return
	}

	m.setMetaStorageDefaults()

	m.setDefaultContainerSecurityContext(&mvVersion, m.Spec.PodTemplate)

	m.SetHealthCheckerDefaults()

	m.setDefaultContainerResourceLimits(m.Spec.PodTemplate)
}

func (m *Milvus) setMetaStorageDefaults() {
	if m.Spec.MetaStorage == nil {
		m.Spec.MetaStorage = &MetaStorageSpec{}
	}

	if m.Spec.MetaStorage.StorageType == "" {
		m.Spec.MetaStorage.StorageType = StorageTypeDurable
	}

	if !m.Spec.MetaStorage.ExternallyManaged {
		if m.Spec.MetaStorage.Size == 0 {
			m.Spec.MetaStorage.Size = 3
		}

		if m.Spec.MetaStorage.Storage == nil {
			m.Spec.MetaStorage.Storage = &core.PersistentVolumeClaimSpec{}
		}

		if len(m.Spec.MetaStorage.Storage.AccessModes) == 0 {
			m.Spec.MetaStorage.Storage.AccessModes = []core.PersistentVolumeAccessMode{
				core.ReadWriteOnce,
			}
		}

		if m.Spec.MetaStorage.Storage.Resources.Requests == nil {
			m.Spec.MetaStorage.Storage.Resources.Requests = core.ResourceList{
				core.ResourceStorage: resource.MustParse("1Gi"),
			}
		}
	}
}

func (m *Milvus) setDefaultContainerSecurityContext(mvVersion *catalog.MilvusVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = mvVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.MilvusContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.MilvusContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	m.AssignDefaultContainerSecurityContext(mvVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (m *Milvus) AssignDefaultContainerSecurityContext(mvVersion *catalog.MilvusVersion, rc *core.SecurityContext) {
	if rc.AllowPrivilegeEscalation == nil {
		rc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if rc.Capabilities == nil {
		rc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if rc.RunAsNonRoot == nil {
		rc.RunAsNonRoot = pointer.BoolP(true)
	}
	if rc.RunAsUser == nil {
		rc.RunAsUser = mvVersion.Spec.SecurityContext.RunAsUser
	}
	if rc.RunAsGroup == nil {
		rc.RunAsGroup = mvVersion.Spec.SecurityContext.RunAsUser
	}
	if rc.SeccompProfile == nil {
		rc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (m *Milvus) setDefaultContainerResourceLimits(podTemplate *ofstv2.PodTemplateSpec) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.MilvusContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}
}
