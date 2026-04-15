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
	catalogv1alpha1 "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"kmodules.xyz/client-go/apiextensions"
	metautil "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
	ofst_util "kmodules.xyz/offshoot-api/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (d *DocumentDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDocumentDB))
}

func (d *DocumentDB) ResourcePlural() string {
	return ResourcePluralDocumentDB
}

func (d *DocumentDB) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", d.ResourcePlural(), SchemeGroupVersion.Group)
}

func (d *DocumentDB) OffshootName() string {
	return d.Name
}

func (d *DocumentDB) ServiceName() string {
	return d.OffshootName()
}

func (d *DocumentDB) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(d.Spec.ServiceTemplates, ServiceAlias(alias))
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (d *DocumentDB) PetSetName() string {
	return d.OffshootName()
}

func (d *DocumentDB) GoverningServiceName() string {
	return metautil.NameWithSuffix(d.ServiceName(), "pods")
}

// Owner returns owner reference to resources
func (d *DocumentDB) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(d, SchemeGroupVersion.WithKind(d.ResourceKind()))
}

func (d *DocumentDB) ResourceKind() string {
	return ResourceKindDocumentDB
}

func (d *DocumentDB) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		metautil.NameLabelKey:      d.ResourceFQN(),
		metautil.InstanceLabelKey:  d.Name,
		metautil.ManagedByLabelKey: SchemeGroupVersion.Group,
	}
	return metautil.OverwriteKeys(selector, extraSelectors...)
}

func (d *DocumentDB) OffshootLabels() map[string]string {
	return d.offshootLabels(d.OffshootSelectors(), nil)
}

func (d *DocumentDB) offshootLabels(selector, override map[string]string) map[string]string {
	selector[metautil.ComponentLabelKey] = kubedb.ComponentDatabase
	return metautil.FilterKeys(SchemeGroupVersion.Group, selector, metautil.OverwriteKeys(nil, d.Labels, override))
}

func (d *DocumentDB) PodLabels(podTemplate *ofstv2.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if podTemplate != nil && podTemplate.Labels != nil {
		return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), podTemplate.Labels)
	}
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), nil)
}

func (d *DocumentDB) PodControllerLabels(podTemplate *ofstv2.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if podTemplate != nil && podTemplate.Controller.Labels != nil {
		return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), podTemplate.Controller.Labels)
	}
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), nil)
}

func (d *DocumentDB) GetAuthSecretName() string {
	if d.Spec.AuthSecret != nil && d.Spec.AuthSecret.Name != "" {
		return d.Spec.AuthSecret.Name
	}
	return metautil.NameWithSuffix(d.OffshootName(), "auth")
}

func (d *DocumentDB) GetPersistentSecrets() []string {
	var secrets []string
	secrets = append(secrets, d.GetAuthSecretName())
	return secrets
}

func (d *DocumentDB) Finalizer() string {
	return fmt.Sprintf("%s/%s", apis.Finalizer, d.ResourceSingular())
}

func (d *DocumentDB) ResourceSingular() string {
	return ResourceSingularDocumentDB
}

func (d *DocumentDB) SetDefaults(_ client.Client, documentDBVersion catalogv1alpha1.DocumentDBVersion) {
	if d.Spec.DeletionPolicy == "" {
		d.Spec.DeletionPolicy = DeletionPolicyDelete
	}
	if d.Spec.StorageType == "" {
		d.Spec.StorageType = StorageTypeDurable
	}
	if d.Spec.Replicas == nil {
		d.Spec.Replicas = ptr.To(int32(1))
	}
	d.initializePodTemplates()

	d.SetDefaultPodSecurityContext(d.Spec.PodTemplate, &documentDBVersion)
	d.SetDocumentDBContainerDefaults(d.Spec.PodTemplate, &documentDBVersion)
}

func (d *DocumentDB) SetDefaultPodSecurityContext(podTemplate *ofstv2.PodTemplateSpec, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
	if podTemplate == nil {
		return
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
	if podTemplate.Spec.SecurityContext.RunAsUser == nil {
		podTemplate.Spec.SecurityContext.RunAsUser = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
	if podTemplate.Spec.SecurityContext.RunAsGroup == nil {
		podTemplate.Spec.SecurityContext.RunAsGroup = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
}

func (d *DocumentDB) SetDocumentDBContainerDefaults(podTemplate *ofstv2.PodTemplateSpec, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
	if podTemplate == nil {
		return
	}
	container := ofst_util.EnsureContainerExists(podTemplate, kubedb.DocumentDBContainerName)
	d.setContainerDefaultSecurityContext(container, documentDBVersion)
	d.setContainerDefaultResources(container, *kubedb.DefaultResources.DeepCopy())
}

func (d *DocumentDB) setContainerDefaultSecurityContext(container *core.Container, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	d.assignDefaultContainerSecurityContext(container.SecurityContext, documentDBVersion)
}

func (d *DocumentDB) setContainerDefaultResources(container *core.Container, defaultResources core.ResourceRequirements) {
	if container != nil {
		apis.SetDefaultResourceLimits(&container.Resources, defaultResources)
	}
}

func (d *DocumentDB) assignDefaultContainerSecurityContext(sc *core.SecurityContext, documentDBVersion *catalogv1alpha1.DocumentDBVersion) {
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
		sc.RunAsUser = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = documentDBVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (d *DocumentDB) initializePodTemplates() {
	if d.Spec.PodTemplate == nil {
		d.Spec.PodTemplate = new(ofstv2.PodTemplateSpec)
	}
}

func (d *DocumentDB) SetHealthCheckerDefaults() {
	if d.Spec.HealthChecker.PeriodSeconds == nil {
		d.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if d.Spec.HealthChecker.TimeoutSeconds == nil {
		d.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if d.Spec.HealthChecker.FailureThreshold == nil {
		d.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}
