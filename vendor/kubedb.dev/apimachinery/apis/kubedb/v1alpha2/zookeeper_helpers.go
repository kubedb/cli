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
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (z *ZooKeeper) PetSetName() string {
	return z.OffshootName()
}

func (z *ZooKeeper) ServiceAccountName() string {
	return z.OffshootName()
}

func (z *ZooKeeper) ConfigSecretName() string {
	return meta_util.NameWithSuffix(z.OffshootName(), "config")
}

func (z *ZooKeeper) PVCName(alias string) string {
	return alias
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
	return fmt.Sprintf("%v.%v.svc:%d", z.ServiceName(), z.Namespace, kubedb.ZooKeeperClientPort)
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
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
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

func (z *ZooKeeper) GetKeystoreSecretName() string {
	if z.Spec.KeystoreCredSecret != nil && z.Spec.KeystoreCredSecret.Name != "" {
		return z.Spec.KeystoreCredSecret.Name
	}
	return meta_util.NameWithSuffix(z.OffshootName(), "keystore-cred")
}

func (k *ZooKeeper) DefaultUserCredSecretName(username string) string {
	return meta_util.NameWithSuffix(k.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", username), "_", "-"))
}

func (z *ZooKeeper) DefaultKeystoreCredSecretName() string {
	return meta_util.NameWithSuffix(z.Name, strings.ReplaceAll("keystore-cred", "_", "-"))
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

func (z *ZooKeeper) SetDefaults(kc client.Client) {
	if z.Spec.DeletionPolicy == "" {
		z.Spec.DeletionPolicy = DeletionPolicyDelete
	}
	if z.Spec.Replicas == nil {
		z.Spec.Replicas = pointer.Int32P(1)
	}

	if z.Spec.Halted {
		if z.Spec.DeletionPolicy == DeletionPolicyDoNotTerminate {
			klog.Errorf(`Can't halt, since termination policy is 'DoNotTerminate'`)
			return
		}
		z.Spec.DeletionPolicy = DeletionPolicyHalt
	}

	if !z.Spec.DisableAuth {
		if z.Spec.AuthSecret == nil {
			z.Spec.AuthSecret = &SecretReference{}
		}
		if z.Spec.AuthSecret.Kind == "" {
			z.Spec.AuthSecret.Kind = kubedb.ResourceKindSecret
		}
	}

	var zkVersion catalog.ZooKeeperVersion
	err := kc.Get(context.TODO(), types.NamespacedName{Name: z.Spec.Version}, &zkVersion)
	if err != nil {
		klog.Errorf("can't get the zookeeper version object %s for %s \n", err.Error(), z.Spec.Version)
		return
	}

	z.setDefaultContainerSecurityContext(&zkVersion, &z.Spec.PodTemplate)

	dbContainer := coreutil.GetContainerByName(z.Spec.PodTemplate.Spec.Containers, kubedb.ZooKeeperContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	initContainer := coreutil.GetContainerByName(z.Spec.PodTemplate.Spec.InitContainers, kubedb.ZooKeeperInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}

	if z.Spec.EnableSSL {
		z.SetTLSDefaults()
	}

	z.SetHealthCheckerDefaults()
	if z.Spec.Monitor != nil {
		if z.Spec.Monitor.Prometheus == nil {
			z.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if z.Spec.Monitor.Prometheus != nil && z.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			z.Spec.Monitor.Prometheus.Exporter.Port = kubedb.ZooKeeperMetricsPort
		}
		z.Spec.Monitor.SetDefaults()
		if z.Spec.Monitor.Prometheus != nil {
			if z.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
				z.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = zkVersion.Spec.SecurityContext.RunAsUser
			}
			if z.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
				z.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = zkVersion.Spec.SecurityContext.RunAsUser
			}
		}
	}
}

func (z *ZooKeeper) SetTLSDefaults() {
	if z.Spec.TLS == nil || z.Spec.TLS.IssuerRef == nil {
		return
	}
	z.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(z.Spec.TLS.Certificates, string(ZooKeeperServerCert), z.CertificateName(ZooKeeperServerCert))
	z.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(z.Spec.TLS.Certificates, string(ZooKeeperClientCert), z.CertificateName(ZooKeeperClientCert))
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

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.ZooKeeperContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.ZooKeeperContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}

	z.assignDefaultContainerSecurityContext(zkVersion, container.SecurityContext)

	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.ZooKeeperInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.ZooKeeperInitContainerName,
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
	return kubedb.DefaultStatsPath
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
	return z.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
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

func (z *ZooKeeper) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicasOfPetSet(lister.PetSets(z.Namespace), labels.SelectorFromSet(z.OffshootLabels()), expectedItems)
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (z *ZooKeeper) CertificateName(alias ZooKeeperCertificateAlias) string {
	return meta_util.NameWithSuffix(z.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (z *ZooKeeper) GetCertSecretName(alias ZooKeeperCertificateAlias) string {
	if z.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(z.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return z.CertificateName(alias)
}

// CertSecretVolumeName returns the CertSecretVolumeName
// Values will be like: client-certs, server-certs etc.
func (k *ZooKeeper) CertSecretVolumeName(alias ZooKeeperCertificateAlias) string {
	return string(alias) + "-certs"
}
