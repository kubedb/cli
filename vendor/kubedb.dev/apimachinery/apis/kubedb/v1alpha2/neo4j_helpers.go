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
	"slices"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
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

func (*Neo4j) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralNeo4j))
}

func (r *Neo4j) ResourceSingular() string {
	return ResourceSingularNeo4j
}

func (r *Neo4j) GetAuthSecretName() string {
	if r.Spec.AuthSecret != nil && r.Spec.AuthSecret.Name != "" {
		return r.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(r.OffshootName(), "auth")
}

func (r *Neo4j) GetKeystoreSecretName() string {
	if r.Spec.TLS.KeystoreCredSecret != nil && r.Spec.TLS.KeystoreCredSecret.Name != "" {
		return r.Spec.TLS.KeystoreCredSecret.Name
	}
	return meta_util.NameWithSuffix(r.OffshootName(), "keystore-cred")
}

func (r *Neo4j) OffshootName() string {
	return r.Name
}

func (r *Neo4j) ServiceAccountName() string {
	return r.OffshootName()
}

func (r *Neo4j) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      r.ResourceFQN(),
		meta_util.InstanceLabelKey:  r.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (r *Neo4j) ConfigSecretName() string {
	uid := string(r.UID)
	return meta_util.NameWithSuffix(r.OffshootName(), uid[len(uid)-6:])
}

func (r *Neo4j) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", r.ResourcePlural(), kubedb.GroupName)
}

func (r *Neo4j) ResourcePlural() string {
	return ResourcePluralNeo4j
}

func (r *Neo4j) GetPersistentSecrets() []string {
	var secrets []string
	if !r.Spec.DisableSecurity {
		secrets = append(secrets, r.GetAuthSecretName())
	}
	return secrets
}

// Owner returns owner reference to resources
func (r *Neo4j) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(r, SchemeGroupVersion.WithKind(r.ResourceKind()))
}

func (r *Neo4j) ResourceKind() string {
	return ResourceKindNeo4j
}

func (r *Neo4j) SetDefaults(kc client.Client) {
	if r.Spec.Replicas == nil {
		r.Spec.Replicas = ptr.To(int32(1))
	}
	if r.Spec.DeletionPolicy == "" {
		r.Spec.DeletionPolicy = DeletionPolicyDelete
	}
	if r.Spec.StorageType == "" {
		r.Spec.StorageType = StorageTypeDurable
	}

	if !r.Spec.DisableSecurity {
		if r.Spec.AuthSecret == nil {
			r.Spec.AuthSecret = &SecretReference{}
		}
		if r.Spec.AuthSecret.Kind == "" {
			r.Spec.AuthSecret.Kind = kubedb.ResourceKindSecret
		}
	}

	var neoVersion catalog.Neo4jVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name: r.Spec.Version,
	}, &neoVersion)
	if err != nil {
		klog.Errorf("can't get the Neo4j version object %s for %s \n", err.Error(), r.Spec.Version)
		return
	}

	r.setDefaultContainerSecurityContext(&neoVersion, &r.Spec.PodTemplate)
	r.SetHealthCheckerDefaults()

	if r.Spec.Monitor != nil {
		if r.Spec.Monitor.Prometheus == nil {
			r.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if r.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			r.Spec.Monitor.Prometheus.Exporter.Port = kubedb.Neo4jExporterPort
		}
		r.Spec.Monitor.SetDefaults()
		if r.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			r.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = neoVersion.Spec.SecurityContext.RunAsUser
		}
		if r.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			r.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = neoVersion.Spec.SecurityContext.RunAsUser
		}
	}

	r.SetTLSDefaults()

	dbContainer := coreutil.GetContainerByName(r.Spec.PodTemplate.Spec.Containers, kubedb.Neo4jContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesNeo4j)
	}
}

func (r *Neo4j) SetTLSDefaults() {
	if r.Spec.TLS == nil || r.Spec.TLS.IssuerRef == nil {
		return
	}
	if r.Spec.TLS.HTTP == nil {
		r.Spec.TLS.HTTP = &ProtocolTLSConfig{
			Mode: TLSModeTLS,
		}
	}
	if r.Spec.TLS.Bolt == nil {
		r.Spec.TLS.Bolt = &ProtocolTLSConfig{
			Mode: TLSModeTLS,
		}
	}
	if r.Spec.TLS.Cluster == nil {
		r.Spec.TLS.Cluster = &ProtocolTLSConfig{
			Mode: TLSModeMTLS,
		}
	}
	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(Neo4jCertificateTypeServer), r.CertificateName(Neo4jCertificateTypeServer))
	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(Neo4jCertificateTypeClient), r.CertificateName(Neo4jCertificateTypeClient))
}

func (r *Neo4j) setDefaultContainerSecurityContext(neoVersion *catalog.Neo4jVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = neoVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.Neo4jContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.Neo4jContainerName,
		}
		podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
	}

	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}

	r.assignDefaultContainerSecurityContext(neoVersion, container.SecurityContext)
}

func (r *Neo4j) assignDefaultContainerSecurityContext(n4Version *catalog.Neo4jVersion, rc *core.SecurityContext) {
	if rc.AllowPrivilegeEscalation == nil {
		rc.AllowPrivilegeEscalation = ptr.To(false)
	}
	if rc.Capabilities == nil {
		rc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if rc.RunAsNonRoot == nil {
		rc.RunAsNonRoot = ptr.To(true)
	}
	if rc.RunAsUser == nil {
		rc.RunAsUser = n4Version.Spec.SecurityContext.RunAsUser
	}
	if rc.SeccompProfile == nil {
		rc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (r *Neo4j) PetSetName() string {
	return r.OffshootName()
}

func (r *Neo4j) ServiceName() string {
	return r.OffshootName()
}

func (r *Neo4j) GoverningServiceName() string {
	return meta_util.NameWithSuffix(r.ServiceName(), "pods")
}

func (r *Neo4j) InternalServiceName(id int32) string {
	return meta_util.NameWithSuffix(r.ServiceName(), fmt.Sprintf("%d", id))
}

func (r *Neo4j) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, r.Labels, override))
}

func (r *Neo4j) OffshootLabels() map[string]string {
	return r.offshootLabels(r.OffshootSelectors(), nil)
}

func (r *Neo4j) IsProtocolDisabled(protocol Neo4jProtocol) bool {
	return slices.Contains(r.Spec.DisabledProtocols, protocol)
}

func (r *Neo4j) DefaultPodRoleName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "role")
}

func (r *Neo4j) DefaultPodRoleBindingName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "rolebinding")
}

func (r *Neo4j) PodLabels(extraLabels ...map[string]string) map[string]string {
	return r.offshootLabels(meta_util.OverwriteKeys(r.OffshootSelectors(), extraLabels...), r.Spec.PodTemplate.Labels)
}

type Neo4jApp struct {
	*Neo4j
}

func (r *Neo4jApp) Name() string {
	return r.Neo4j.Name
}

func (r *Neo4jApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceKindNeo4j))
}

func (r *Neo4j) AppBindingMeta() appcat.AppBindingMeta {
	return &Neo4jApp{r}
}

func (r *Neo4j) GetConnectionScheme() string {
	scheme := "http"
	if r.Spec.TLS != nil && r.Spec.TLS.HTTP.Mode != TLSModeDisabled {
		scheme = "https"
	}
	return scheme
}

func (r *Neo4j) SetHealthCheckerDefaults() {
	if r.Spec.HealthChecker.PeriodSeconds == nil {
		r.Spec.HealthChecker.PeriodSeconds = ptr.To(int32(10))
	}
	if r.Spec.HealthChecker.TimeoutSeconds == nil {
		r.Spec.HealthChecker.TimeoutSeconds = ptr.To(int32(10))
	}
	if r.Spec.HealthChecker.FailureThreshold == nil {
		r.Spec.HealthChecker.FailureThreshold = ptr.To(int32(3))
	}
}

func (r *Neo4j) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(r.Spec.ServiceTemplates, alias)
	return r.offshootLabels(meta_util.OverwriteKeys(r.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

type neo4jStatsService struct {
	*Neo4j
}

func (r neo4jStatsService) GetNamespace() string {
	return r.Neo4j.GetNamespace()
}

func (r neo4jStatsService) ServiceName() string {
	return r.OffshootName() + "-stats"
}

func (r neo4jStatsService) ServiceMonitorName() string {
	return r.ServiceName()
}

func (r neo4jStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return r.OffshootLabels()
}

func (r neo4jStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (r neo4jStatsService) Scheme() string {
	scheme := "http"
	if r.Spec.TLS != nil && r.Spec.TLS.HTTP.Mode != TLSModeDisabled {
		scheme = "https"
	}
	return scheme
}

func (r neo4jStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (r *Neo4j) StatsService() mona.StatsAccessor {
	return &neo4jStatsService{r}
}

func (r *Neo4j) StatsServiceLabels() map[string]string {
	return r.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (r *Neo4j) GetCertSecretName(alias Neo4jCertificateType) string {
	if r.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(r.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return r.CertificateName(alias)
}

func (r *Neo4j) CertificateName(alias Neo4jCertificateType) string {
	return meta_util.NameWithSuffix(r.Name, fmt.Sprintf("%s-cert", string(alias)))
}
