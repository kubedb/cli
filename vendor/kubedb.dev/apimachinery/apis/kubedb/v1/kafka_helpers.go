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

package v1

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/google/uuid"
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
	ofst_util "kmodules.xyz/offshoot-api/util"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
)

func (*Kafka) Hub() {}
func (k *Kafka) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralKafka))
}

func (k *Kafka) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(ResourceKindKafka))
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

func (k *Kafka) GoverningServiceNameCruiseControl() string {
	return meta_util.NameWithSuffix(k.ServiceName(), kubedb.KafkaNodeRolesCruiseControl)
}

func (k *Kafka) StandbyServiceName() string {
	return meta_util.NameWithPrefix(k.ServiceName(), kubedb.KafkaStandbyServiceSuffix)
}

func (k *Kafka) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
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
		k.NodeRoleSpecificLabelKey(KafkaNodeRoleController): kubedb.KafkaNodeRoleSet,
	})
}

func (k *Kafka) BrokerNodeSelectors() map[string]string {
	return meta_util.OverwriteKeys(k.OffshootSelectors(), map[string]string{
		k.NodeRoleSpecificLabelKey(KafkaNodeRoleBroker): kubedb.KafkaNodeRoleSet,
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
	return kubedb.DefaultStatsPath
}

func (ks kafkaStatsService) Scheme() string {
	return ""
}

func (k *Kafka) StatsService() mona.StatsAccessor {
	return &kafkaStatsService{k}
}

func (k *Kafka) StatsServiceLabels() map[string]string {
	return k.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (k *Kafka) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Controller.Labels)
}

func (k *Kafka) PodLabels(extraLabels ...map[string]string) map[string]string {
	return k.offshootLabels(meta_util.OverwriteKeys(k.OffshootSelectors(), extraLabels...), k.Spec.PodTemplate.Labels)
}

func (k *Kafka) PetSetName() string {
	return k.OffshootName()
}

func (k *Kafka) CombinedPetSetName() string {
	return k.PetSetName()
}

func (k *Kafka) ControllerPetSetName() string {
	if k.Spec.Topology.Controller.Suffix != "" {
		return meta_util.NameWithSuffix(k.PetSetName(), k.Spec.Topology.Controller.Suffix)
	}
	return meta_util.NameWithSuffix(k.PetSetName(), string(KafkaNodeRoleController))
}

func (k *Kafka) BrokerPetSetName() string {
	if k.Spec.Topology.Broker.Suffix != "" {
		return meta_util.NameWithSuffix(k.PetSetName(), k.Spec.Topology.Broker.Suffix)
	}
	return meta_util.NameWithSuffix(k.PetSetName(), string(KafkaNodeRoleBroker))
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

func (k *Kafka) GetAuthSecretName() string {
	if k.Spec.AuthSecret != nil && k.Spec.AuthSecret.Name != "" {
		return k.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(k.OffshootName(), "auth")
}

func (k *Kafka) GetKeystoreSecretName() string {
	if k.Spec.KeystoreCredSecret != nil && k.Spec.KeystoreCredSecret.Name != "" {
		return k.Spec.KeystoreCredSecret.Name
	}
	return meta_util.NameWithSuffix(k.OffshootName(), "keystore-cred")
}

func (k *Kafka) GetPersistentSecrets() []string {
	var secrets []string
	if k.Spec.AuthSecret != nil {
		secrets = append(secrets, k.Spec.AuthSecret.Name)
	}
	if k.Spec.KeystoreCredSecret != nil {
		secrets = append(secrets, k.Spec.KeystoreCredSecret.Name)
	}
	return secrets
}

func (k *Kafka) CruiseControlConfigSecretName() string {
	return meta_util.NameWithSuffix(k.OffshootName(), "cruise-control-config")
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

// CertSecretVolumeName returns the CertSecretVolumeName
// Values will be like: client-certs, server-certs etc.
func (k *Kafka) CertSecretVolumeName(alias KafkaCertificateAlias) string {
	return string(alias) + "-certs"
}

// CertSecretVolumeMountPath returns the CertSecretVolumeMountPath
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
		k.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if k.Spec.HealthChecker.TimeoutSeconds == nil {
		k.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if k.Spec.HealthChecker.FailureThreshold == nil {
		k.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (k *Kafka) SetDefaults() {
	if k.Spec.Halted {
		if k.Spec.DeletionPolicy == DeletionPolicyDoNotTerminate {
			klog.Errorf(`Can't halt, since deletion policy is 'DoNotTerminate'`)
			return
		}
		k.Spec.DeletionPolicy = DeletionPolicyHalt
	}

	if k.Spec.DeletionPolicy == "" {
		k.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if k.Spec.StorageType == "" {
		k.Spec.StorageType = StorageTypeDurable
	}

	var kfVersion catalog.KafkaVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: k.Spec.Version}, &kfVersion)
	if err != nil {
		klog.Errorf("can't get the kafka version object %s for %s \n", err.Error(), k.Spec.Version)
		return
	}

	if k.Spec.CruiseControl != nil {
		k.setDefaultContainerSecurityContext(&kfVersion, &k.Spec.CruiseControl.PodTemplate)
	}

	k.Spec.Monitor.SetDefaults()
	if k.Spec.Monitor != nil && k.Spec.Monitor.Prometheus != nil {
		if k.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			k.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = kfVersion.Spec.SecurityContext.RunAsUser
		}
		if k.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			k.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = kfVersion.Spec.SecurityContext.RunAsUser
		}
	}

	if k.Spec.Topology != nil {
		if k.Spec.Topology.Controller != nil {
			if k.Spec.Topology.Controller.Suffix == "" {
				k.Spec.Topology.Controller.Suffix = string(KafkaNodeRoleController)
			}
			if k.Spec.Topology.Controller.Replicas == nil {
				k.Spec.Topology.Controller.Replicas = pointer.Int32P(1)
			}
			k.setDefaultContainerSecurityContext(&kfVersion, &k.Spec.Topology.Controller.PodTemplate)

			dbContainer := coreutil.GetContainerByName(k.Spec.Topology.Controller.PodTemplate.Spec.Containers, kubedb.KafkaContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
			}
		}

		if k.Spec.Topology.Broker != nil {
			if k.Spec.Topology.Broker.Suffix == "" {
				k.Spec.Topology.Broker.Suffix = string(KafkaNodeRoleBroker)
			}
			if k.Spec.Topology.Broker.Replicas == nil {
				k.Spec.Topology.Broker.Replicas = pointer.Int32P(1)
			}
			k.setDefaultContainerSecurityContext(&kfVersion, &k.Spec.Topology.Broker.PodTemplate)

			dbContainer := coreutil.GetContainerByName(k.Spec.Topology.Broker.PodTemplate.Spec.Containers, kubedb.KafkaContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
			}
		}
	} else {
		if k.Spec.Replicas == nil {
			k.Spec.Replicas = pointer.Int32P(1)
		}
		k.setDefaultContainerSecurityContext(&kfVersion, &k.Spec.PodTemplate)

		dbContainer := coreutil.GetContainerByName(k.Spec.PodTemplate.Spec.Containers, kubedb.KafkaContainerName)
		if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
		}
	}
	k.SetDefaultEnvs()

	if k.Spec.EnableSSL {
		k.SetTLSDefaults()
	}
	k.SetHealthCheckerDefaults()
}

func (k *Kafka) setDefaultContainerSecurityContext(kfVersion *catalog.KafkaVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = kfVersion.Spec.SecurityContext.RunAsUser
	}
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.KafkaContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: kubedb.KafkaContainerName,
		}
	}
	if dbContainer.SecurityContext == nil {
		dbContainer.SecurityContext = &core.SecurityContext{}
	}
	k.assignDefaultContainerSecurityContext(kfVersion, dbContainer.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *dbContainer)
}

func (k *Kafka) assignDefaultContainerSecurityContext(kfVersion *catalog.KafkaVersion, sc *core.SecurityContext) {
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

func (k *Kafka) SetDefaultEnvs() {
	clusterID := k.GenerateClusterID()
	if k.Spec.Topology != nil {
		if k.Spec.Topology.Controller != nil {
			k.setClusterIDEnv(&k.Spec.Topology.Controller.PodTemplate, clusterID)
		}
		if k.Spec.Topology.Broker != nil {
			k.setClusterIDEnv(&k.Spec.Topology.Broker.PodTemplate, clusterID)
		}
	} else {
		k.setClusterIDEnv(&k.Spec.PodTemplate, clusterID)
	}
}

func (k *Kafka) setClusterIDEnv(podTemplate *ofst.PodTemplateSpec, clusterID string) {
	container := ofst_util.EnsureContainerExists(podTemplate, kubedb.KafkaContainerName)
	env := coreutil.GetEnvByName(container.Env, kubedb.EnvKafkaClusterID)
	if env == nil {
		container.Env = coreutil.UpsertEnvVars(container.Env, core.EnvVar{
			Name:  kubedb.EnvKafkaClusterID,
			Value: clusterID,
		})
	}
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

func (k *Kafka) GetCruiseControlClientID() string {
	return meta_util.NameWithSuffix(k.Name, "cruise-control")
}

func (k *Kafka) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	if k.Spec.Topology != nil {
		expectedItems = 2
	}
	return checkReplicas(lister.PetSets(k.Namespace), labels.SelectorFromSet(k.OffshootLabels()), expectedItems)
}

// GenerateClusterID Kafka uses Leach-Salz UUIDs for cluster ID. It requires 16 bytes of base64 encoded RFC 4122 version 1 UUID.
// Here, the generated uuid is 32 bytes hexadecimal string and have 5 hyphen separated parts: 8-4-4-4-12
// part 3 contains version number, part 4 is a randomly generated clock sequence and
// part 5 is node field that contains MAC address of the host machine
// These 3 parts will be used as cluster ID
// ref: https://kafka.apache.org/31/javadoc/org/apache/kafka/common/Uuid.html
// ref: https://go-recipes.dev/how-to-generate-uuids-with-go-be3988e771a6
func (k *Kafka) GenerateClusterID() string {
	clusterUUID, _ := uuid.NewUUID()
	slicedUUID := strings.Split(clusterUUID.String(), "-")
	trimmedUUID := slicedUUID[2:]
	generatedUUID := strings.Join(trimmedUUID, "-")
	return generatedUUID[:len(generatedUUID)-1] + "w"
}

func (k *Kafka) KafkaSaslListenerProtocolConfigKey(protocol string, mechanism string) string {
	return fmt.Sprintf("listener.name.%s.%s.sasl.jaas.config", strings.ToLower(protocol), strings.ToLower(mechanism))
}

func (k *Kafka) KafkaEnabledSASLMechanismsKey(protocol string) string {
	return fmt.Sprintf("listener.name.%s.sasl.enabled.mechanisms", strings.ToLower(protocol))
}
