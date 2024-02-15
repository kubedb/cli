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
	"path/filepath"
	"strings"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kafka"
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

func (k *ConnectCluster) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralConnectCluster))
}

func (k *ConnectCluster) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(ResourceKindConnectCluster))
}

func (k *ConnectCluster) ResourceShortCode() string {
	return ResourceCodeConnectCluster
}

func (k *ConnectCluster) ResourceKind() string {
	return ResourceKindConnectCluster
}

func (k *ConnectCluster) ResourceSingular() string {
	return ResourceSingularConnectCluster
}

func (k *ConnectCluster) ResourcePlural() string {
	return ResourcePluralConnectCluster
}

func (k *ConnectCluster) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", k.ResourcePlural(), kafka.GroupName)
}

// Owner returns owner reference to resources
func (k *ConnectCluster) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(k.ResourceKind()))
}

func (k *ConnectCluster) OffshootName() string {
	return k.Name
}

func (k *ConnectCluster) ServiceName() string {
	return k.OffshootName()
}

func (k *ConnectCluster) GoverningServiceName() string {
	return meta_util.NameWithSuffix(k.ServiceName(), "pods")
}

func (k *ConnectCluster) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentKafka
	return meta_util.FilterKeys(kafka.GroupName, selector, meta_util.OverwriteKeys(nil, k.Labels, override))
}

func (k *ConnectCluster) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      k.ResourceFQN(),
		meta_util.InstanceLabelKey:  k.Name,
		meta_util.ManagedByLabelKey: kafka.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (k *ConnectCluster) OffshootLabels() map[string]string {
	return k.offshootLabels(k.OffshootSelectors(), nil)
}

// GetServiceTemplate returns a pointer to the desired serviceTemplate referred by "aliaS". Otherwise, it returns nil.
func (k *ConnectCluster) GetServiceTemplate(templates []api.NamedServiceTemplateSpec, alias api.ServiceAlias) ofst.ServiceTemplateSpec {
	for i := range templates {
		c := templates[i]
		if c.Alias == alias {
			return c.ServiceTemplateSpec
		}
	}
	return ofst.ServiceTemplateSpec{}
}

func (k *ConnectCluster) ServiceLabels(alias api.ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := k.GetServiceTemplate(k.Spec.ServiceTemplates, alias)
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (k *ConnectCluster) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Controller.Labels)
}

type connectClusterStatsService struct {
	*ConnectCluster
}

func (ks connectClusterStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (ks connectClusterStatsService) GetNamespace() string {
	return ks.ConnectCluster.GetNamespace()
}

func (ks connectClusterStatsService) ServiceName() string {
	return ks.OffshootName() + "-stats"
}

func (ks connectClusterStatsService) ServiceMonitorName() string {
	return ks.ServiceName()
}

func (ks connectClusterStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return ks.OffshootLabels()
}

func (ks connectClusterStatsService) Path() string {
	return DefaultStatsPath
}

func (ks connectClusterStatsService) Scheme() string {
	return ""
}

func (k *ConnectCluster) StatsService() mona.StatsAccessor {
	return &connectClusterStatsService{k}
}

func (k *ConnectCluster) StatsServiceLabels() map[string]string {
	return k.ServiceLabels(api.StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}

func (k *ConnectCluster) PodLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Labels)
}

func (k *ConnectCluster) StatefulSetName() string {
	return k.OffshootName()
}

func (k *ConnectCluster) ConfigSecretName() string {
	return meta_util.NameWithSuffix(k.OffshootName(), "config")
}

func (k *ConnectCluster) GetPersistentSecrets() []string {
	var secrets []string
	if k.Spec.AuthSecret != nil {
		secrets = append(secrets, k.Spec.AuthSecret.Name)
	}
	if k.Spec.KeystoreCredSecret != nil {
		secrets = append(secrets, k.Spec.KeystoreCredSecret.Name)
	}
	return secrets
}

func (k *ConnectCluster) KafkaClientCredentialsSecretName() string {
	return meta_util.NameWithSuffix(k.Name, "kafka-client-cred")
}

func (k *ConnectCluster) DefaultUserCredSecretName(username string) string {
	return meta_util.NameWithSuffix(k.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", username), "_", "-"))
}

func (k *ConnectCluster) DefaultKeystoreCredSecretName() string {
	return meta_util.NameWithSuffix(k.Name, strings.ReplaceAll("connect-keystore-cred", "_", "-"))
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (k *ConnectCluster) CertificateName(alias ConnectClusterCertificateAlias) string {
	return meta_util.NameWithSuffix(k.Name, fmt.Sprintf("%s-connect-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (k *ConnectCluster) GetCertSecretName(alias ConnectClusterCertificateAlias) string {
	if k.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(k.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return k.CertificateName(alias)
}

// returns CertSecretVolumeMountPath
// if configDir is "/opt/kafka/config",
// mountPath will be, "/opt/kafka/config/<alias>".
func (k *ConnectCluster) CertSecretVolumeMountPath(configDir string, cert string) string {
	return filepath.Join(configDir, cert)
}

func (k *ConnectCluster) SetHealthCheckerDefaults() {
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

func (k *ConnectCluster) SetDefaults() {
	if k.Spec.TerminationPolicy == "" {
		k.Spec.TerminationPolicy = api.TerminationPolicyDelete
	}

	if k.Spec.Replicas == nil {
		k.Spec.Replicas = pointer.Int32P(1)
	}

	var kfVersion catalog.KafkaVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: k.Spec.Version}, &kfVersion)
	if err != nil {
		klog.Errorf("can't get the kafka version object %s for %s \n", err.Error(), k.Spec.Version)
		return
	}

	k.setDefaultContainerSecurityContext(&kfVersion, &k.Spec.PodTemplate)
	k.setDefaultInitContainerSecurityContext(&k.Spec.PodTemplate)

	dbContainer := coreutil.GetContainerByName(k.Spec.PodTemplate.Spec.Containers, ConnectClusterContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, api.DefaultResources)
	}

	k.Spec.Monitor.SetDefaults()
	if k.Spec.Monitor != nil && k.Spec.Monitor.Prometheus != nil && k.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
		k.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = kfVersion.Spec.SecurityContext.RunAsUser
	}

	if k.Spec.EnableSSL {
		k.SetTLSDefaults()
	}
	k.SetHealthCheckerDefaults()

	k.SetDefaultEnvs()
}

func (k *ConnectCluster) SetDefaultEnvs() {
	container := coreutil.GetContainerByName(k.Spec.PodTemplate.Spec.Containers, ConnectClusterContainerName)
	if container != nil {
		env := coreutil.GetEnvByName(container.Env, ConnectClusterModeEnv)
		if env == nil {
			if *k.Spec.Replicas == 1 {
				container.Env = coreutil.UpsertEnvVars(container.Env, core.EnvVar{
					Name:  ConnectClusterModeEnv,
					Value: string(ConnectClusterNodeRoleStandalone),
				})
			} else {
				container.Env = coreutil.UpsertEnvVars(container.Env, core.EnvVar{
					Name:  ConnectClusterModeEnv,
					Value: string(ConnectClusterNodeRoleDistributed),
				})
			}
		}
	}
}

func (k *ConnectCluster) setDefaultInitContainerSecurityContext(podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	for _, name := range k.Spec.ConnectorPlugins {
		connectorVersion := &catalog.KafkaConnectorVersion{}
		err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: name}, connectorVersion)
		if err != nil {
			klog.Errorf("can't get the kafka connector version object %s for %s \n", err.Error(), name)
			return
		}

		initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, strings.ToLower(connectorVersion.Spec.Type))

		if initContainer == nil {
			initContainer = &core.Container{
				Name: strings.ToLower(connectorVersion.Spec.Type),
			}
			podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)
		}
		if initContainer.SecurityContext == nil {
			initContainer.SecurityContext = &core.SecurityContext{}
		}
		k.assignDefaultInitContainerSecurityContext(connectorVersion, initContainer.SecurityContext)
	}
}

func (k *ConnectCluster) setDefaultContainerSecurityContext(kfVersion *catalog.KafkaVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = kfVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, ConnectClusterContainerName)
	if container == nil {
		container = &core.Container{
			Name: ConnectClusterContainerName,
		}
		podTemplate.Spec.Containers = append(podTemplate.Spec.Containers, *container)
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	k.assignDefaultContainerSecurityContext(kfVersion, container.SecurityContext)
}

func (k *ConnectCluster) assignDefaultInitContainerSecurityContext(connectorVersion *catalog.KafkaConnectorVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = connectorVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = connectorVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (k *ConnectCluster) assignDefaultContainerSecurityContext(kfVersion *catalog.KafkaVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = kfVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = kfVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (k *ConnectCluster) SetTLSDefaults() {
	if k.Spec.TLS == nil || k.Spec.TLS.IssuerRef == nil {
		return
	}
	k.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(k.Spec.TLS.Certificates, string(ConnectClusterServerCert), k.CertificateName(ConnectClusterServerCert))
	k.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(k.Spec.TLS.Certificates, string(ConnectClusterClientCert), k.CertificateName(ConnectClusterClientCert))
}

type ConnectClusterApp struct {
	*ConnectCluster
}

func (r ConnectClusterApp) Name() string {
	return r.ConnectCluster.Name
}

func (r ConnectClusterApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kafka.GroupName, ResourceSingularConnectCluster))
}

func (k *ConnectCluster) AppBindingMeta() appcat.AppBindingMeta {
	return &ConnectClusterApp{k}
}

func (k *ConnectCluster) GetConnectionScheme() string {
	scheme := "http"
	if k.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}
