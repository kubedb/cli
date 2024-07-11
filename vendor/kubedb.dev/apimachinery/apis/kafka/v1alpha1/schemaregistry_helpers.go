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

func (k *SchemaRegistry) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralSchemaRegistry))
}

func (k *SchemaRegistry) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(ResourceKindSchemaRegistry))
}

func (k *SchemaRegistry) ResourceShortCode() string {
	return ResourceCodeSchemaRegistry
}

func (k *SchemaRegistry) ResourceKind() string {
	return ResourceKindSchemaRegistry
}

func (k *SchemaRegistry) ResourceSingular() string {
	return ResourceSingularSchemaRegistry
}

func (k *SchemaRegistry) ResourcePlural() string {
	return ResourcePluralSchemaRegistry
}

func (k *SchemaRegistry) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", k.ResourcePlural(), kafka.GroupName)
}

// Owner returns owner reference to resources
func (k *SchemaRegistry) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(k.ResourceKind()))
}

func (k *SchemaRegistry) OffshootName() string {
	return k.Name
}

func (k *SchemaRegistry) GoverningServiceName() string {
	return meta_util.NameWithSuffix(k.ServiceName(), "pods")
}

func (k *SchemaRegistry) ServiceName() string {
	return k.OffshootName()
}

func (k *SchemaRegistry) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentKafka
	return meta_util.FilterKeys(kafka.GroupName, selector, meta_util.OverwriteKeys(nil, k.Labels, override))
}

func (k *SchemaRegistry) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      k.ResourceFQN(),
		meta_util.InstanceLabelKey:  k.Name,
		meta_util.ManagedByLabelKey: kafka.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (k *SchemaRegistry) OffshootLabels() map[string]string {
	return k.offshootLabels(k.OffshootSelectors(), nil)
}

// GetServiceTemplate returns a pointer to the desired serviceTemplate referred by "aliaS". Otherwise, it returns nil.
func (k *SchemaRegistry) GetServiceTemplate(templates []dbapi.NamedServiceTemplateSpec, alias dbapi.ServiceAlias) ofst.ServiceTemplateSpec {
	for i := range templates {
		c := templates[i]
		if c.Alias == alias {
			return c.ServiceTemplateSpec
		}
	}
	return ofst.ServiceTemplateSpec{}
}

func (k *SchemaRegistry) ServiceLabels(alias dbapi.ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := k.GetServiceTemplate(k.Spec.ServiceTemplates, alias)
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (k *SchemaRegistry) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Controller.Labels)
}

//type schemaRegistryStatsService struct {
//	*SchemaRegistry
//}
//
//func (ks schemaRegistryStatsService) TLSConfig() *promapi.TLSConfig {
//	return nil
//}
//
//func (ks schemaRegistryStatsService) GetNamespace() string {
//	return ks.SchemaRegistry.GetNamespace()
//}
//
//func (ks schemaRegistryStatsService) ServiceName() string {
//	return ks.OffshootName() + "-stats"
//}
//
//func (ks schemaRegistryStatsService) ServiceMonitorName() string {
//	return ks.ServiceName()
//}
//
//func (ks schemaRegistryStatsService) ServiceMonitorAdditionalLabels() map[string]string {
//	return ks.OffshootLabels()
//}
//
//func (ks schemaRegistryStatsService) Path() string {
//	return DefaultStatsPath
//}
//
//func (ks schemaRegistryStatsService) Scheme() string {
//	return ""
//}
//
//func (k *SchemaRegistry) StatsService() mona.StatsAccessor {
//	return &schemaRegistryStatsService{k}
//}
//
//func (k *SchemaRegistry) StatsServiceLabels() map[string]string {
//	return k.ServiceLabels(api.StatsServiceAlias, map[string]string{LabelRole: RoleStats})
//}

func (k *SchemaRegistry) PodLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Labels)
}

func (k *SchemaRegistry) PetSetName() string {
	return k.OffshootName()
}

func (k *SchemaRegistry) ConfigSecretName() string {
	return meta_util.NameWithSuffix(k.OffshootName(), "config")
}

func (k *SchemaRegistry) GetPersistentSecrets() []string {
	var secrets []string
	return secrets
}

func (k *SchemaRegistry) KafkaClientCredentialsSecretName() string {
	return meta_util.NameWithSuffix(k.Name, "kafka-client-cred")
}

func (k *SchemaRegistry) SetHealthCheckerDefaults() {
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

func (k *SchemaRegistry) SetDefaults() {
	if k.Spec.DeletionPolicy == "" {
		k.Spec.DeletionPolicy = dbapi.DeletionPolicyDelete
	}

	if k.Spec.Replicas == nil {
		k.Spec.Replicas = pointer.Int32P(1)
	}

	var ksrVersion catalog.SchemaRegistryVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: k.Spec.Version}, &ksrVersion)
	if err != nil {
		klog.Errorf("can't get the schema-registry version object %s for %s \n", err.Error(), k.Spec.Version)
		return
	}

	k.setDefaultContainerSecurityContext(&ksrVersion, &k.Spec.PodTemplate)

	dbContainer := coreutil.GetContainerByName(k.Spec.PodTemplate.Spec.Containers, SchemaRegistryContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	k.SetHealthCheckerDefaults()
}

func (k *SchemaRegistry) setDefaultContainerSecurityContext(ksrVersion *catalog.SchemaRegistryVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = ksrVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, SchemaRegistryContainerName)
	if container == nil {
		container = &core.Container{
			Name: SchemaRegistryContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	k.assignDefaultContainerSecurityContext(ksrVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (k *SchemaRegistry) assignDefaultContainerSecurityContext(ksrVersion *catalog.SchemaRegistryVersion, sc *core.SecurityContext) {
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

type SchemaRegistryApp struct {
	*SchemaRegistry
}

func (r SchemaRegistryApp) Name() string {
	return r.SchemaRegistry.Name
}

func (r SchemaRegistryApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kafka.GroupName, ResourceSingularSchemaRegistry))
}

func (k *SchemaRegistry) AppBindingMeta() appcat.AppBindingMeta {
	return &SchemaRegistryApp{k}
}

func (k *SchemaRegistry) GetConnectionScheme() string {
	scheme := "http"
	return scheme
}
