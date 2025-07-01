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
	"path/filepath"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (i *Ignite) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralIgnite))
}

func (i *Ignite) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(i, SchemeGroupVersion.WithKind(ResourceKindIgnite))
}

func (i *Ignite) ResourceKind() string {
	return ResourceKindIgnite
}

func (i *Ignite) ResourceSingular() string {
	return ResourceSingularIgnite
}

func (i *Ignite) ResourcePlural() string {
	return ResourcePluralIgnite
}

func (i *Ignite) Finalizer() string {
	return fmt.Sprintf("%s/%s", apis.Finalizer, i.ResourceSingular())
}

func (i *Ignite) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", i.ResourcePlural(), kubedb.GroupName)
}

func (i *Ignite) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      i.ResourceFQN(),
		meta_util.InstanceLabelKey:  i.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (i *Ignite) OffshootName() string {
	return i.Name
}

func (i *Ignite) GetAuthSecretName() string {
	if i.Spec.AuthSecret != nil && i.Spec.AuthSecret.Name != "" {
		return i.Spec.AuthSecret.Name
	}
	return i.DefaultAuthSecretName()
}

func (i *Ignite) GetPersistentSecrets() []string {
	var secrets []string
	if i.Spec.AuthSecret != nil {
		secrets = append(secrets, i.GetAuthSecretName())
	}
	return secrets
}

// Owner returns owner reference to resources
func (i *Ignite) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(i, SchemeGroupVersion.WithKind(i.ResourceKind()))
}

func (i *Ignite) SetDefaults(kc client.Client) {
	if i.Spec.Replicas == nil {
		i.Spec.Replicas = pointer.Int32P(1)
	}

	if i.Spec.DeletionPolicy == "" {
		i.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if i.Spec.StorageType == "" {
		i.Spec.StorageType = StorageTypeDurable
	}

	var igVersion catalog.IgniteVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name: i.Spec.Version,
	}, &igVersion)
	if err != nil {
		return
	}

	i.setDefaultContainerSecurityContext(&igVersion, &i.Spec.PodTemplate)

	dbContainer := coreutil.GetContainerByName(i.Spec.PodTemplate.Spec.Containers, kubedb.IgniteContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil || dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	i.SetHealthCheckerDefaults()

	i.Spec.Monitor.SetDefaults()
}

func (i *Ignite) SetHealthCheckerDefaults() {
	if i.Spec.HealthChecker.PeriodSeconds == nil {
		i.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if i.Spec.HealthChecker.TimeoutSeconds == nil {
		i.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if i.Spec.HealthChecker.FailureThreshold == nil {
		i.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (i *Ignite) setDefaultContainerSecurityContext(igVersion *catalog.IgniteVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = igVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.IgniteContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.IgniteContainerName,
		}
		podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
	}

	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	i.assignDefaultContainerSecurityContext(igVersion, container.SecurityContext)

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.IgniteInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.IgniteInitContainerName,
		}
		podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	i.assignDefaultContainerSecurityContext(igVersion, initContainer.SecurityContext)
}

func (i *Ignite) assignDefaultContainerSecurityContext(igVersion *catalog.IgniteVersion, rc *core.SecurityContext) {
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
		rc.RunAsUser = igVersion.Spec.SecurityContext.RunAsUser
	}
	if rc.SeccompProfile == nil {
		rc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (i *Ignite) PetSetName() string {
	return i.OffshootName()
}

func (i *Ignite) ServiceName() string { return i.OffshootName() }

func (i *Ignite) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, i.Labels, override))
}

func (i *Ignite) OffshootLabels() map[string]string {
	return i.offshootLabels(i.OffshootSelectors(), nil)
}

func (i *Ignite) GoverningServiceName() string {
	return meta_util.NameWithSuffix(i.ServiceName(), "pods")
}

func (i *Ignite) DefaultAuthSecretName() string {
	return meta_util.NameWithSuffix(i.OffshootName(), "auth")
}

func (i *Ignite) ServiceAccountName() string {
	return i.OffshootName()
}

func (i *Ignite) DefaultPodRoleName() string {
	return meta_util.NameWithSuffix(i.OffshootName(), "role")
}

func (i *Ignite) DefaultPodRoleBindingName() string {
	return meta_util.NameWithSuffix(i.OffshootName(), "rolebinding")
}

type IgniteApp struct {
	*Ignite
}

func (i *IgniteApp) Name() string {
	return i.Ignite.Name
}

func (i IgniteApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularIgnite))
}

func (i *Ignite) AppBindingMeta() appcat.AppBindingMeta {
	return &IgniteApp{i}
}

func (i *Ignite) GetConnectionScheme() string {
	scheme := "http"
	return scheme
}

func (i *Ignite) PodLabels(extraLabels ...map[string]string) map[string]string {
	return i.offshootLabels(meta_util.OverwriteKeys(i.OffshootSelectors(), extraLabels...), i.Spec.PodTemplate.Labels)
}

func (i *Ignite) ConfigSecretName() string {
	return meta_util.NameWithSuffix(i.OffshootName(), "config")
}

func (i *Ignite) PVCName(alias string) string {
	return meta_util.NameWithSuffix(i.Name, alias)
}

func (i *Ignite) Address() string {
	return fmt.Sprintf("%v.%v.svc.cluster.local", i.Name, i.Namespace)
}

type igniteStatsService struct {
	*Ignite
}

func (i igniteStatsService) GetNamespace() string {
	return i.Ignite.GetNamespace()
}

func (i igniteStatsService) ServiceName() string {
	return i.OffshootName() + "-stats"
}

func (i igniteStatsService) ServiceMonitorName() string {
	return i.ServiceName()
}

func (i igniteStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return i.OffshootLabels()
}

func (i igniteStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (i igniteStatsService) Scheme() string {
	return ""
}

func (i igniteStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (i Ignite) StatsService() mona.StatsAccessor {
	return &igniteStatsService{&i}
}

func (i Ignite) StatsServiceLabels() map[string]string {
	return i.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (i Ignite) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(i.Spec.ServiceTemplates, alias)
	return i.offshootLabels(meta_util.OverwriteKeys(i.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (i Ignite) GetIgniteCertSecretName(alias IgniteCertificateAlias) string {
	if i.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(i.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return i.IgniteCertificateName(alias)
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (i Ignite) IgniteCertificateName(alias IgniteCertificateAlias) string {
	return meta_util.NameWithSuffix(i.Name, fmt.Sprintf("%s-cert", string(alias)))
}

func (i Ignite) GetIgniteConnectionScheme() string {
	scheme := "http"
	if i.Spec.TLS != nil {
		scheme = "https"
	}
	return scheme
}

func (i Ignite) GetIgniteKeystoreSecretName() string {
	if i.Spec.KeystoreCredSecret != nil && i.Spec.KeystoreCredSecret.Name != "" {
		return i.Spec.KeystoreCredSecret.Name
	}
	return meta_util.NameWithSuffix(i.OffshootName(), "keystore-cred")
}

// CertSecretVolumeName returns the CertSecretVolumeName
// Values will be like: client-certs, server-certs etc.
func (i Ignite) IgniteCertSecretVolumeName(alias IgniteCertificateAlias) string {
	return string(alias) + "-certs"
}

// CertSecretVolumeMountPath returns the CertSecretVolumeMountPath
func (i Ignite) IgniteCertSecretVolumeMountPath(configDir string, cert string) string {
	return filepath.Join(configDir, cert)
}

func (i Ignite) SetTLSDefaults() {
	if i.Spec.TLS == nil || i.Spec.TLS.IssuerRef == nil {
		return
	}
	i.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(i.Spec.TLS.Certificates, string(IgniteServerCert), i.IgniteCertificateName(IgniteServerCert))
	i.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(i.Spec.TLS.Certificates, string(IgniteClientCert), i.IgniteCertSecretVolumeName(IgniteClientCert))
}
