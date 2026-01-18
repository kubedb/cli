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
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (*Weaviate) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralWeaviate))
}

type WeaviateApp struct {
	*Weaviate
}

func (w WeaviateApp) Name() string {
	return w.Weaviate.Name
}

func (w WeaviateApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularWeaviate))
}

func (w *Weaviate) AppBindingMeta() appcat.AppBindingMeta {
	return &WeaviateApp{w}
}

func (w *Weaviate) GetPersistentSecrets() []string {
	var secrets []string
	secrets = append(secrets, w.GetAuthSecretName())
	return secrets
}

func (w *Weaviate) ResourceShortCode() string {
	return ResourceCodeWeaviate
}

func (w *Weaviate) ResourceKind() string {
	return ResourceKindWeaviate
}

func (w *Weaviate) ResourceSingular() string {
	return ResourceSingularWeaviate
}

func (w *Weaviate) ResourcePlural() string {
	return ResourcePluralWeaviate
}

func (w *Weaviate) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(w, SchemeGroupVersion.WithKind(ResourceKindRabbitmq))
}

func (w *Weaviate) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", w.ResourcePlural(), kubedb.GroupName)
}

// Owner returns owner reference to resources
func (w *Weaviate) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(w, SchemeGroupVersion.WithKind(w.ResourceKind()))
}

func (w *Weaviate) OffshootName() string {
	return w.Name
}

func (w *Weaviate) ServiceName() string {
	return w.OffshootName()
}

func (w *Weaviate) GoverningServiceName() string {
	return meta_util.NameWithSuffix(w.ServiceName(), "pods")
}

func (w *Weaviate) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, w.Labels, override))
}

func (w *Weaviate) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      w.ResourceFQN(),
		meta_util.InstanceLabelKey:  w.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (w *Weaviate) OffshootLabels() map[string]string {
	return w.offshootLabels(w.OffshootSelectors(), nil)
}

func (w *Weaviate) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(w.Spec.ServiceTemplates, alias)
	return w.offshootLabels(meta_util.OverwriteKeys(w.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (w *Weaviate) SetHealthCheckerDefaults() {
	if w.Spec.HealthChecker.PeriodSeconds == nil {
		w.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if w.Spec.HealthChecker.TimeoutSeconds == nil {
		w.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if w.Spec.HealthChecker.FailureThreshold == nil {
		w.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (w *Weaviate) PodLabels(extraLabels ...map[string]string) map[string]string {
	return w.offshootLabels(meta_util.OverwriteKeys(w.OffshootSelectors(), extraLabels...), w.Spec.PodTemplate.Labels)
}

func (w *Weaviate) PetSetName() string {
	return w.OffshootName()
}

func (w *Weaviate) GetAuthSecretName() string {
	if w.Spec.AuthSecret != nil && w.Spec.AuthSecret.Name != "" {
		return w.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(w.OffshootName(), "auth")
}

func (w *Weaviate) Finalizer() string {
	return fmt.Sprintf("%s/%s", apis.Finalizer, w.ResourceSingular())
}

func (w *Weaviate) SetDefaults(kc client.Client) {
	if w.Spec.Replicas == nil {
		w.Spec.Replicas = pointer.Int32P(1)
	}

	if w.Spec.DeletionPolicy == "" {
		w.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if w.Spec.StorageType == "" {
		w.Spec.StorageType = StorageTypeDurable
	}

	var wvVersion catalog.WeaviateVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name:      w.Spec.Version,
		Namespace: "",
	}, &wvVersion)
	if err != nil {
		klog.Errorf("can't get the weaviate version object %s for %s \n", err.Error(), w.Spec.Version)
		return
	}

	w.setDefaultContainerSecurityContext(&wvVersion, &w.Spec.PodTemplate)

	dbContainer := coreutil.GetContainerByName(w.Spec.PodTemplate.Spec.Containers, kubedb.WeaviateContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	w.SetHealthCheckerDefaults()
}

func (w *Weaviate) setDefaultContainerSecurityContext(wvVersion *catalog.WeaviateVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = wvVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.WeaviateContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.WeaviateContainerName,
		}
		podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	w.assignDefaultContainerSecurityContext(wvVersion, container.SecurityContext)
}

func (w *Weaviate) assignDefaultContainerSecurityContext(wvVersion *catalog.WeaviateVersion, rc *core.SecurityContext) {
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
		rc.RunAsUser = wvVersion.Spec.SecurityContext.RunAsUser
	}
	if rc.SeccompProfile == nil {
		rc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (w *Weaviate) GetAPIKey(ctx context.Context, kc client.Client) string {
	secretName := w.GetAuthSecretName()
	var secret core.Secret
	err := kc.Get(ctx, client.ObjectKey{Namespace: w.Namespace, Name: secretName}, &secret)
	if err != nil {
		return ""
	}
	apiKey, ok := secret.Data[kubedb.WeaviateAPIKey]
	if !ok {
		return ""
	}
	return string(apiKey)
}

func (w *Weaviate) GetConnectionScheme() string {
	scheme := "http"
	return scheme
}

func (w *Weaviate) ConfigSecretName() string {
	uid := string(w.UID)
	return meta_util.NameWithSuffix(w.OffshootName(), uid[len(uid)-6:])
}
