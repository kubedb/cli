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
	"path/filepath"
	"strings"

	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (r *RabbitMQ) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRabbitmq))
}

type RabbitmqApp struct {
	*RabbitMQ
}

func (r RabbitmqApp) Name() string {
	return r.RabbitMQ.Name
}

func (r RabbitmqApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularRabbitmq))
}

func (r *RabbitMQ) AppBindingMeta() appcat.AppBindingMeta {
	return &RabbitmqApp{r}
}

func (r *RabbitMQ) GetConnectionScheme() string {
	scheme := "http"
	if r.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}

func (r *RabbitMQ) ResourceShortCode() string {
	return ResourceCodeRabbitmq
}

func (r *RabbitMQ) ResourceKind() string {
	return ResourceKindRabbitmq
}

func (r *RabbitMQ) ResourceSingular() string {
	return ResourceSingularRabbitmq
}

func (r *RabbitMQ) ResourcePlural() string {
	return ResourcePluralRabbitmq
}

func (r *RabbitMQ) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(r, SchemeGroupVersion.WithKind(ResourceKindRabbitmq))
}

func (r *RabbitMQ) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", r.ResourcePlural(), kubedb.GroupName)
}

// Owner returns owner reference to resources
func (r *RabbitMQ) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(r, SchemeGroupVersion.WithKind(r.ResourceKind()))
}

func (r *RabbitMQ) OffshootName() string {
	return r.Name
}

func (r *RabbitMQ) ServiceName() string {
	return r.OffshootName()
}

func (r *RabbitMQ) GoverningServiceName() string {
	return meta_util.NameWithSuffix(r.ServiceName(), "pods")
}

func (r *RabbitMQ) StandbyServiceName() string {
	return meta_util.NameWithPrefix(r.ServiceName(), KafkaStandbyServiceSuffix)
}

func (r *RabbitMQ) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, r.Labels, override))
}

func (r *RabbitMQ) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      r.ResourceFQN(),
		meta_util.InstanceLabelKey:  r.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (r *RabbitMQ) OffshootLabels() map[string]string {
	return r.offshootLabels(r.OffshootSelectors(), nil)
}

func (r *RabbitMQ) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(r.Spec.ServiceTemplates, alias)
	return r.offshootLabels(meta_util.OverwriteKeys(r.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

type RabbitmqStatsService struct {
	*RabbitMQ
}

func (ks RabbitmqStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (ks RabbitmqStatsService) GetNamespace() string {
	return ks.RabbitMQ.GetNamespace()
}

func (ks RabbitmqStatsService) ServiceName() string {
	return ks.OffshootName() + "-stats"
}

func (ks RabbitmqStatsService) ServiceMonitorName() string {
	return ks.ServiceName()
}

func (ks RabbitmqStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return ks.OffshootLabels()
}

func (ks RabbitmqStatsService) Path() string {
	return DefaultStatsPath
}

func (ks RabbitmqStatsService) Scheme() string {
	return ""
}

func (r *RabbitMQ) StatsService() mona.StatsAccessor {
	return &RabbitmqStatsService{r}
}

func (r *RabbitMQ) StatsServiceLabels() map[string]string {
	return r.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}

func (r *RabbitMQ) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return r.offshootLabels(meta_util.OverwriteKeys(r.OffshootSelectors(), extraLabels...), r.Spec.PodTemplate.Controller.Labels)
}

func (r *RabbitMQ) PodLabels(extraLabels ...map[string]string) map[string]string {
	return r.offshootLabels(meta_util.OverwriteKeys(r.OffshootSelectors(), extraLabels...), r.Spec.PodTemplate.Labels)
}

func (r *RabbitMQ) StatefulSetName() string {
	return r.OffshootName()
}

func (r *RabbitMQ) ServiceAccountName() string {
	return r.OffshootName()
}

func (r *RabbitMQ) DefaultPodRoleName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "role")
}

func (r *RabbitMQ) DefaultPodRoleBindingName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "rolebinding")
}

func (r *RabbitMQ) ConfigSecretName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "config")
}

func (r *RabbitMQ) DefaultUserCredSecretName(username string) string {
	return meta_util.NameWithSuffix(r.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", username), "_", "-"))
}

func (r *RabbitMQ) DefaultErlangCookieSecretName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "erlang-cookie")
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (r *RabbitMQ) CertificateName(alias RabbitMQCertificateAlias) string {
	return meta_util.NameWithSuffix(r.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// ClientCertificateCN returns the CN for a client certificate
func (r *RabbitMQ) ClientCertificateCN(alias RabbitMQCertificateAlias) string {
	return fmt.Sprintf("%s-%s", r.Name, string(alias))
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (r *RabbitMQ) GetCertSecretName(alias RabbitMQCertificateAlias) string {
	if r.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(r.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return r.CertificateName(alias)
}

// returns the CertSecretVolumeName
// Values will be like: client-certs, server-certs etc.
func (r *RabbitMQ) CertSecretVolumeName(alias RabbitMQCertificateAlias) string {
	return string(alias) + "-certs"
}

// returns CertSecretVolumeMountPath
// if configDir is "/opt/kafka/config",
// mountPath will be, "/opt/kafka/config/<alias>".
func (r *RabbitMQ) CertSecretVolumeMountPath(configDir string, cert string) string {
	return filepath.Join(configDir, cert)
}

func (r *RabbitMQ) PVCName(alias string) string {
	return meta_util.NameWithSuffix(r.Name, alias)
}

func (r *RabbitMQ) SetHealthCheckerDefaults() {
	if r.Spec.HealthChecker.PeriodSeconds == nil {
		r.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if r.Spec.HealthChecker.TimeoutSeconds == nil {
		r.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if r.Spec.HealthChecker.FailureThreshold == nil {
		r.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (r *RabbitMQ) SetDefaults() {
	if r.Spec.TerminationPolicy == "" {
		r.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if r.Spec.StorageType == "" {
		r.Spec.StorageType = StorageTypeDurable
	}

	r.SetHealthCheckerDefaults()
}

func (r *RabbitMQ) SetTLSDefaults() {
	if r.Spec.TLS == nil || r.Spec.TLS.IssuerRef == nil {
		return
	}
	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(RabbitmqServerCert), r.CertificateName(RabbitmqServerCert))
	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(RabbitmqClientCert), r.CertificateName(RabbitmqClientCert))
}
