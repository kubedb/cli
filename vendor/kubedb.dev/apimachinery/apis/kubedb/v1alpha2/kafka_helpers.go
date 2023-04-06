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

	"kubedb.dev/apimachinery/apis"
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

func (k *Kafka) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralKafka))
}

func (k *Kafka) ResourceShortCode() string {
	return ResourceCodeKafka
}

func (k *Kafka) ResourceKind() string {
	return ResourceKindKafka
}

func (k *Kafka) ResourceSingular() string {
	return ResourceSingularKafka
}

func (k *Kafka) ResourcePlural() string {
	return ResourcePluralKafka
}

func (k *Kafka) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", k.ResourcePlural(), kubedb.GroupName)
}

// Owner returns owner reference to resources
func (k *Kafka) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(k.ResourceKind()))
}

func (k *Kafka) OffshootName() string {
	return k.Name
}

func (k *Kafka) ServiceName() string {
	return k.OffshootName()
}

func (k *Kafka) GoverningServiceName() string {
	return meta_util.NameWithSuffix(k.ServiceName(), "pods")
}

func (k *Kafka) GoverningServiceNameController() string {
	return meta_util.NameWithSuffix(k.ServiceName(), KafkaNodeRolesController)
}

func (k *Kafka) GoverningServiceNameBroker() string {
	return meta_util.NameWithSuffix(k.ServiceName(), KafkaNodeRolesBrokers)
}

func (k *Kafka) StandbyServiceName() string {
	return meta_util.NameWithPrefix(k.ServiceName(), KafkaStandbyServiceSuffix)
}

func (k *Kafka) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, k.Labels, override))
}

func (k *Kafka) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      k.ResourceFQN(),
		meta_util.InstanceLabelKey:  k.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (k *Kafka) ControllerNodeSelectors() map[string]string {
	return meta_util.OverwriteKeys(k.OffshootSelectors(), map[string]string{
		k.NodeRoleSpecificLabelKey(KafkaNodeRoleController): KafkaNodeRoleSet,
	})
}

func (k *Kafka) BrokerNodeSelectors() map[string]string {
	return meta_util.OverwriteKeys(k.OffshootSelectors(), map[string]string{
		k.NodeRoleSpecificLabelKey(KafkaNodeRoleBroker): KafkaNodeRoleSet,
	})
}

func (k *Kafka) OffshootLabels() map[string]string {
	return k.offshootLabels(k.OffshootSelectors(), nil)
}

func (k *Kafka) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(k.Spec.ServiceTemplates, alias)
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (k *Kafka) ControllerServiceLabels() map[string]string {
	return meta_util.OverwriteKeys(k.offshootLabels(k.OffshootLabels(), k.ControllerNodeSelectors()))
}

func (k *Kafka) BrokerServiceLabels() map[string]string {
	return meta_util.OverwriteKeys(k.offshootLabels(k.OffshootLabels(), k.BrokerNodeSelectors()))
}

type kafkaStatsService struct {
	*Kafka
}

func (ks kafkaStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (ks kafkaStatsService) GetNamespace() string {
	return ks.Kafka.GetNamespace()
}

func (ks kafkaStatsService) ServiceName() string {
	return ks.OffshootName() + "-stats"
}

func (ks kafkaStatsService) ServiceMonitorName() string {
	return ks.ServiceName()
}

func (ks kafkaStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return ks.OffshootLabels()
}

func (ks kafkaStatsService) Path() string {
	return DefaultStatsPath
}

func (ks kafkaStatsService) Scheme() string {
	return ""
}

func (k *Kafka) StatsService() mona.StatsAccessor {
	return &kafkaStatsService{k}
}

func (k *Kafka) StatsServiceLabels() map[string]string {
	return k.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}

func (k *Kafka) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Controller.Labels)
}

func (k *Kafka) PodLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Labels)
}

func (k *Kafka) StatefulSetName() string {
	return k.OffshootName()
}

func (k *Kafka) CombinedStatefulSetName() string {
	return k.StatefulSetName()
}

func (k *Kafka) ControllerStatefulSetName() string {
	if k.Spec.Topology.Controller.Suffix != "" {
		return meta_util.NameWithSuffix(k.StatefulSetName(), k.Spec.Topology.Controller.Suffix)
	}
	return meta_util.NameWithSuffix(k.StatefulSetName(), string(KafkaNodeRoleController))
}

func (k *Kafka) BrokerStatefulSetName() string {
	if k.Spec.Topology.Broker.Suffix != "" {
		return meta_util.NameWithSuffix(k.StatefulSetName(), k.Spec.Topology.Broker.Suffix)
	}
	return meta_util.NameWithSuffix(k.StatefulSetName(), string(KafkaNodeRoleBroker))
}

func (k *Kafka) NodeRoleSpecificLabelKey(role KafkaNodeRoleType) string {
	return kubedb.GroupName + "/role-" + string(role)
}

func (k *Kafka) ConfigSecretName(role KafkaNodeRoleType) string {
	if role == KafkaNodeRoleController {
		return meta_util.NameWithSuffix(k.OffshootName(), "controller-config")
	} else if role == KafkaNodeRoleBroker {
		return meta_util.NameWithSuffix(k.OffshootName(), "broker-config")
	}
	return meta_util.NameWithSuffix(k.OffshootName(), "config")
}

func (k *Kafka) DefaultUserCredSecretName(username string) string {
	return meta_util.NameWithSuffix(k.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", username), "_", "-"))
}

func (k *Kafka) DefaultKeystoreCredSecretName() string {
	return meta_util.NameWithSuffix(k.Name, strings.ReplaceAll("keystore-cred", "_", "-"))
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (k *Kafka) CertificateName(alias KafkaCertificateAlias) string {
	return meta_util.NameWithSuffix(k.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// ClientCertificateCN returns the CN for a client certificate
func (k *Kafka) ClientCertificateCN(alias KafkaCertificateAlias) string {
	return fmt.Sprintf("%s-%s", k.Name, string(alias))
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (k *Kafka) GetCertSecretName(alias KafkaCertificateAlias) string {
	if k.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(k.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return k.CertificateName(alias)
}

// returns the CertSecretVolumeName
// Values will be like: client-certs, server-certs etc.
func (k *Kafka) CertSecretVolumeName(alias KafkaCertificateAlias) string {
	return string(alias) + "-certs"
}

// returns CertSecretVolumeMountPath
// if configDir is "/opt/kafka/config",
// mountPath will be, "/opt/kafka/config/<alias>".
func (k *Kafka) CertSecretVolumeMountPath(configDir string, cert string) string {
	return filepath.Join(configDir, cert)
}

func (k *Kafka) PVCName(alias string) string {
	return meta_util.NameWithSuffix(k.Name, alias)
}

func (k *Kafka) SetHealthCheckerDefaults() {
	if k.Spec.HealthChecker.PeriodSeconds == nil {
		k.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(20)
	}
	if k.Spec.HealthChecker.TimeoutSeconds == nil {
		k.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if k.Spec.HealthChecker.FailureThreshold == nil {
		k.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (k *Kafka) SetDefaults() {
	if k.Spec.TerminationPolicy == "" {
		k.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if k.Spec.StorageType == "" {
		k.Spec.StorageType = StorageTypeDurable
	}

	if k.Spec.Topology != nil {
		if k.Spec.Topology.Controller != nil {
			if k.Spec.Topology.Controller.Suffix == "" {
				k.Spec.Topology.Controller.Suffix = string(KafkaNodeRoleController)
			}
			if k.Spec.Topology.Controller.Replicas == nil {
				k.Spec.Topology.Controller.Replicas = pointer.Int32P(1)
			}
			apis.SetDefaultResourceLimits(&k.Spec.Topology.Controller.Resources, DefaultResources)
		}

		if k.Spec.Topology.Broker != nil {
			if k.Spec.Topology.Broker.Suffix == "" {
				k.Spec.Topology.Broker.Suffix = string(KafkaNodeRoleBroker)
			}
			if k.Spec.Topology.Broker.Replicas == nil {
				k.Spec.Topology.Broker.Replicas = pointer.Int32P(1)
			}
			apis.SetDefaultResourceLimits(&k.Spec.Topology.Broker.Resources, DefaultResources)
		}
	} else {
		apis.SetDefaultResourceLimits(&k.Spec.PodTemplate.Spec.Resources, DefaultResources)
		if k.Spec.Replicas == nil {
			k.Spec.Replicas = pointer.Int32P(1)
		}
	}
	if k.Spec.EnableSSL {
		k.SetTLSDefaults()
	}
	k.SetHealthCheckerDefaults()
}

func (k *Kafka) SetTLSDefaults() {
	if k.Spec.TLS == nil || k.Spec.TLS.IssuerRef == nil {
		return
	}
	k.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(k.Spec.TLS.Certificates, string(KafkaServerCert), k.CertificateName(KafkaServerCert))
	k.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(k.Spec.TLS.Certificates, string(KafkaClientCert), k.CertificateName(KafkaClientCert))
}

// ToMap returns ClusterTopology in a Map
func (kfTopology *KafkaClusterTopology) ToMap() map[KafkaNodeRoleType]KafkaNode {
	topology := make(map[KafkaNodeRoleType]KafkaNode)

	if kfTopology.Controller != nil {
		topology[KafkaNodeRoleController] = *kfTopology.Controller
	}
	if kfTopology.Broker != nil {
		topology[KafkaNodeRoleBroker] = *kfTopology.Broker
	}
	return topology
}

type KafkaApp struct {
	*Kafka
}

func (r KafkaApp) Name() string {
	return r.Kafka.Name
}

func (r KafkaApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularKafka))
}

func (k *Kafka) AppBindingMeta() appcat.AppBindingMeta {
	return &KafkaApp{k}
}

func (k *Kafka) GetConnectionScheme() string {
	scheme := "http"
	if k.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}
