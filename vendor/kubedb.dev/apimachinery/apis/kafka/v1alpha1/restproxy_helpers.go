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
	"context"
	"fmt"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kafka"
	"kubedb.dev/apimachinery/apis/kubedb"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/apimachinery/crds"

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

func (k *RestProxy) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRestProxy))
}

func (k *RestProxy) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(ResourceKindRestProxy))
}

func (k *RestProxy) ResourceShortCode() string {
	return ResourceCodeRestProxy
}

func (k *RestProxy) ResourceKind() string {
	return ResourceKindRestProxy
}

func (k *RestProxy) ResourceSingular() string {
	return ResourceSingularRestProxy
}

func (k *RestProxy) ResourcePlural() string {
	return ResourcePluralRestProxy
}

func (k *RestProxy) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", k.ResourcePlural(), kafka.GroupName)
}

// Owner returns owner reference to resources
func (k *RestProxy) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(k.ResourceKind()))
}

func (k *RestProxy) OffshootName() string {
	return k.Name
}

func (k *RestProxy) GoverningServiceName() string {
	return meta_util.NameWithSuffix(k.ServiceName(), "pods")
}

func (k *RestProxy) ServiceName() string {
	return k.OffshootName()
}

func (k *RestProxy) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentKafka
	return meta_util.FilterKeys(kafka.GroupName, selector, meta_util.OverwriteKeys(nil, k.Labels, override))
}

func (k *RestProxy) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      k.ResourceFQN(),
		meta_util.InstanceLabelKey:  k.Name,
		meta_util.ManagedByLabelKey: kafka.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (k *RestProxy) OffshootLabels() map[string]string {
	return k.offshootLabels(k.OffshootSelectors(), nil)
}

// GetServiceTemplate returns a pointer to the desired serviceTemplate referred by "aliaS". Otherwise, it returns nil.
func (k *RestProxy) GetServiceTemplate(templates []dbapi.NamedServiceTemplateSpec, alias dbapi.ServiceAlias) ofst.ServiceTemplateSpec {
	for i := range templates {
		c := templates[i]
		if c.Alias == alias {
			return c.ServiceTemplateSpec
		}
	}
	return ofst.ServiceTemplateSpec{}
}

func (k *RestProxy) ServiceLabels(alias dbapi.ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := k.GetServiceTemplate(k.Spec.ServiceTemplates, alias)
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (k *RestProxy) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Controller.Labels)
}

func (k *RestProxy) PodLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Labels)
}

func (k *RestProxy) PetSetName() string {
	return k.OffshootName()
}

func (k *RestProxy) ConfigSecretName() string {
	return meta_util.NameWithSuffix(k.OffshootName(), "config")
}

func (k *RestProxy) GetPersistentSecrets() []string {
	var secrets []string
	return secrets
}

func (k *RestProxy) KafkaClientCredentialsSecretName() string {
	return meta_util.NameWithSuffix(k.Name, "kafka-client-cred")
}

func (k *RestProxy) SetHealthCheckerDefaults() {
	if k.Spec.HealthChecker.PeriodSeconds == nil {
		k.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if k.Spec.HealthChecker.TimeoutSeconds == nil {
		k.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if k.Spec.HealthChecker.FailureThreshold == nil {
		k.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (k *RestProxy) SetDefaults() {
	if k.Spec.DeletionPolicy == "" {
		k.Spec.DeletionPolicy = dbapi.DeletionPolicyDelete
	}

	if k.Spec.Replicas == nil {
		k.Spec.Replicas = pointer.Int32P(1)
	}

	var ksrVersion catalog.SchemaRegistryVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: k.Spec.Version}, &ksrVersion)
	if err != nil {
		klog.Errorf("can't get the version object %s for %s \n", err.Error(), k.Spec.Version)
		return
	}

	k.setDefaultContainerSecurityContext(&ksrVersion, &k.Spec.PodTemplate)

	dbContainer := coreutil.GetContainerByName(k.Spec.PodTemplate.Spec.Containers, RestProxyContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	k.SetHealthCheckerDefaults()
}

func (k *RestProxy) setDefaultContainerSecurityContext(ksrVersion *catalog.SchemaRegistryVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = ksrVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, RestProxyContainerName)
	if container == nil {
		container = &core.Container{
			Name: RestProxyContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	k.assignDefaultContainerSecurityContext(ksrVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (k *RestProxy) assignDefaultContainerSecurityContext(ksrVersion *catalog.SchemaRegistryVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = ksrVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

type RestProxyApp struct {
	*RestProxy
}

func (r RestProxyApp) Name() string {
	return r.RestProxy.Name
}

func (r RestProxyApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kafka.GroupName, ResourceSingularRestProxy))
}

func (k *RestProxy) AppBindingMeta() appcat.AppBindingMeta {
	return &RestProxyApp{k}
}

func (k *RestProxy) GetConnectionScheme() string {
	scheme := "http"
	return scheme
}
