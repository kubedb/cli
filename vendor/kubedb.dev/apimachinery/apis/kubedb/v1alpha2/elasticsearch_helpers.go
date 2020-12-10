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
	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	core_util "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ElasticsearchNodeAffinityTemplateVar = "NODE_ROLE"
)

func (_ Elasticsearch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralElasticsearch))
}

var _ apis.ResourceInfo = &Elasticsearch{}

func (e Elasticsearch) OffshootName() string {
	return e.Name
}

func (e Elasticsearch) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseKind: ResourceKindElasticsearch,
		LabelDatabaseName: e.Name,
	}
}

func (e Elasticsearch) OffshootLabels() map[string]string {
	out := e.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularElasticsearch
	out[meta_util.InstanceLabelKey] = e.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = kubedb.GroupName
	return meta_util.FilterKeys(kubedb.GroupName, out, e.Labels)
}

func (e Elasticsearch) ResourceShortCode() string {
	return ResourceCodeElasticsearch
}

func (e Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (e Elasticsearch) ResourceSingular() string {
	return ResourceSingularElasticsearch
}

func (e Elasticsearch) ResourcePlural() string {
	return ResourcePluralElasticsearch
}

func (e Elasticsearch) ServiceName() string {
	return e.OffshootName()
}

func (e *Elasticsearch) MasterDiscoveryServiceName() string {
	return meta_util.NameWithSuffix(e.ServiceName(), "master")
}

func (e Elasticsearch) GoverningServiceName() string {
	return meta_util.NameWithSuffix(e.ServiceName(), "pods")
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (e *Elasticsearch) CertificateName(alias ElasticsearchCertificateAlias) string {
	return meta_util.NameWithSuffix(e.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// MustCertSecretName returns the secret name for a certificate alias
func (e *Elasticsearch) MustCertSecretName(alias ElasticsearchCertificateAlias) string {
	if e == nil {
		panic("missing Elasticsearch database")
	} else if e.Spec.TLS == nil {
		panic(fmt.Errorf("Elasticsearch %s/%s is missing tls spec", e.Namespace, e.Name))
	}
	name, ok := kmapi.GetCertificateSecretName(e.Spec.TLS.Certificates, string(alias))
	if !ok {
		panic(fmt.Errorf("Elasticsearch %s/%s is missing secret name for %s certificate", e.Namespace, e.Name, alias))
	}
	return name
}

// ClientCertificateCN returns the CN for a client certificate
func (e *Elasticsearch) ClientCertificateCN(alias ElasticsearchCertificateAlias) string {
	return fmt.Sprintf("%s-%s", e.Name, string(alias))
}

// returns the volume name for certificate secret.
// Values will be like: transport-certs, http-certs etc.
func (e *Elasticsearch) CertSecretVolumeName(alias ElasticsearchCertificateAlias) string {
	return string(alias) + "-certs"
}

// returns the mountPath for certificate secrets.
// if configDir is "/usr/share/elasticsearch/config",
// mountPath will be, "/usr/share/elasticsearch/config/certs/<alias>".
func (e *Elasticsearch) CertSecretVolumeMountPath(configDir string, alias ElasticsearchCertificateAlias) string {
	return filepath.Join(configDir, "certs", string(alias))
}

// returns the secret name for the  user credentials (ie. username, password)
// If username contains underscore (_), it will be replaced by hyphen (‐) for
// the Kubernetes naming convention.
func (e *Elasticsearch) UserCredSecretName(userName string) string {
	return meta_util.NameWithSuffix(e.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", userName), "_", "-"))
}

// returns the secret name for the default elasticsearch configuration
func (e *Elasticsearch) ConfigSecretName() string {
	return meta_util.NameWithSuffix(e.Name, "config")
}

func (e *Elasticsearch) GetConnectionScheme() string {
	scheme := "http"
	if e.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}

func (e *Elasticsearch) GetConnectionURL() string {
	return fmt.Sprintf("%v://%s.%s:%d", e.GetConnectionScheme(), e.OffshootName(), e.Namespace, ElasticsearchRestPort)
}

func (e *Elasticsearch) CombinedStatefulSetName() string {
	return e.OffshootName()
}

func (e *Elasticsearch) MasterStatefulSetName() string {
	if e.Spec.Topology.Master.Prefix != "" {
		return fmt.Sprintf("%s-%s", e.Spec.Topology.Master.Prefix, e.OffshootName())
	}
	return fmt.Sprintf("%s-%s", ElasticsearchMasterNodePrefix, e.OffshootName())
}

func (e *Elasticsearch) DataStatefulSetName() string {
	if e.Spec.Topology.Data.Prefix != "" {
		return fmt.Sprintf("%s-%s", e.Spec.Topology.Data.Prefix, e.OffshootName())
	}
	return fmt.Sprintf("%s-%s", ElasticsearchDataNodePrefix, e.OffshootName())
}

func (e *Elasticsearch) IngestStatefulSetName() string {
	if e.Spec.Topology.Ingest.Prefix != "" {
		return fmt.Sprintf("%s-%s", e.Spec.Topology.Ingest.Prefix, e.OffshootName())
	}
	return fmt.Sprintf("%s-%s", ElasticsearchIngestNodePrefix, e.OffshootName())
}

type elasticsearchApp struct {
	*Elasticsearch
}

func (r elasticsearchApp) Name() string {
	return r.Elasticsearch.Name
}

func (r elasticsearchApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularElasticsearch))
}

func (e Elasticsearch) AppBindingMeta() appcat.AppBindingMeta {
	return &elasticsearchApp{&e}
}

type elasticsearchStatsService struct {
	*Elasticsearch
}

func (e elasticsearchStatsService) GetNamespace() string {
	return e.Elasticsearch.GetNamespace()
}

func (e elasticsearchStatsService) ServiceName() string {
	return e.OffshootName() + "-stats"
}

func (e elasticsearchStatsService) ServiceMonitorName() string {
	return e.ServiceName()
}

func (e elasticsearchStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return e.OffshootLabels()
}

func (e elasticsearchStatsService) Path() string {
	return DefaultStatsPath
}

func (e elasticsearchStatsService) Scheme() string {
	return ""
}

func (e Elasticsearch) StatsService() mona.StatsAccessor {
	return &elasticsearchStatsService{&e}
}

func (e Elasticsearch) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, e.OffshootSelectors(), e.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (e *Elasticsearch) SetDefaults(esVersion *v1alpha1.ElasticsearchVersion, topology *core_util.Topology) {
	if e == nil {
		return
	}

	if e.Spec.StorageType == "" {
		e.Spec.StorageType = StorageTypeDurable
	}

	if e.Spec.TerminationPolicy == "" {
		e.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if e.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		e.Spec.PodTemplate.Spec.ServiceAccountName = e.OffshootName()
	}

	// Set default values for internal admin user
	if e.Spec.InternalUsers != nil {
		var userSpec ElasticsearchUserSpec

		// load values
		if value, exist := e.Spec.InternalUsers[string(ElasticsearchInternalUserAdmin)]; exist {
			userSpec = value
		}

		// set defaults
		userSpec.Reserved = true
		userSpec.BackendRoles = []string{"admin"}

		// overwrite values
		e.Spec.InternalUsers[string(ElasticsearchInternalUserAdmin)] = userSpec

	}

	// set default elasticsearch node name prefix
	if e.Spec.Topology != nil {

		// Default to "ingest"
		if e.Spec.Topology.Ingest.Prefix == "" {
			e.Spec.Topology.Ingest.Prefix = ElasticsearchIngestNodePrefix
		}
		setDefaultResourceLimits(&e.Spec.Topology.Ingest.Resources, defaultElasticsearchResourceLimits, defaultElasticsearchResourceLimits)

		// Default to "data"
		if e.Spec.Topology.Data.Prefix == "" {
			e.Spec.Topology.Data.Prefix = ElasticsearchDataNodePrefix
		}
		setDefaultResourceLimits(&e.Spec.Topology.Data.Resources, defaultElasticsearchResourceLimits, defaultElasticsearchResourceLimits)

		// Default to "master"
		if e.Spec.Topology.Master.Prefix == "" {
			e.Spec.Topology.Master.Prefix = ElasticsearchMasterNodePrefix
		}
		setDefaultResourceLimits(&e.Spec.Topology.Master.Resources, defaultElasticsearchResourceLimits, defaultElasticsearchResourceLimits)
	} else {
		setDefaultResourceLimits(&e.Spec.PodTemplate.Spec.Resources, defaultElasticsearchResourceLimits, defaultElasticsearchResourceLimits)
	}

	e.setDefaultAffinity(&e.Spec.PodTemplate, e.OffshootSelectors(), topology)
	e.SetTLSDefaults(esVersion)
	e.Spec.Monitor.SetDefaults()
}

// setDefaultAffinity
func (e *Elasticsearch) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
	if podTemplate == nil {
		return
	} else if podTemplate.Spec.Affinity != nil {
		// Update topologyKey fields according to Kubernetes version
		topology.ConvertAffinity(podTemplate.Spec.Affinity)
		return
	}

	podTemplate.Spec.Affinity = &core.Affinity{
		PodAntiAffinity: &core.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []core.WeightedPodAffinityTerm{
				// Prefer to not schedule multiple pods on the same node
				{
					Weight: 100,
					PodAffinityTerm: core.PodAffinityTerm{
						Namespaces: []string{e.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels:      labels,
							MatchExpressions: e.GetMatchExpressions(),
						},

						TopologyKey: core.LabelHostname,
					},
				},
				// Prefer to not schedule multiple pods on the node with same zone
				{
					Weight: 50,
					PodAffinityTerm: core.PodAffinityTerm{
						Namespaces: []string{e.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels:      labels,
							MatchExpressions: e.GetMatchExpressions(),
						},
						TopologyKey: topology.LabelZone,
					},
				},
			},
		},
	}
}

// set default tls configuration (ie. alias, secretName)
func (e *Elasticsearch) SetTLSDefaults(esVersion *v1alpha1.ElasticsearchVersion) {
	// If security is disabled (ie. DisableSecurity: true), ignore.
	if e.Spec.DisableSecurity {
		return
	}

	tlsConfig := e.Spec.TLS
	if tlsConfig == nil {
		tlsConfig = &kmapi.TLSConfig{}
	}

	// If the issuerRef is nil, the operator will create the CA certificate.
	// It is required even if the spec.EnableSSL is false. Because, the transport
	// layer is always secured with certificates. Unless you turned off all the security
	// by setting spec.DisableSecurity to true.
	if tlsConfig.IssuerRef == nil {
		tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchCACert),
			SecretName: e.CertificateName(ElasticsearchCACert),
			Subject: &kmapi.X509Subject{
				Organizations: []string{KubeDBOrganization},
			},
		})
	}

	// transport layer is always secured with certificate
	tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
		Alias:      string(ElasticsearchTransportCert),
		SecretName: e.CertificateName(ElasticsearchTransportCert),
		Subject: &kmapi.X509Subject{
			Organizations: []string{KubeDBOrganization},
		},
	})

	// If SSL is enabled, set missing certificate spec
	if e.Spec.EnableSSL {
		// http
		tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchHTTPCert),
			SecretName: e.CertificateName(ElasticsearchHTTPCert),
			Subject: &kmapi.X509Subject{
				Organizations: []string{KubeDBOrganization},
			},
		})

		// Set missing admin certificate spec, if authPlugin is either "OpenDistro" or "SearchGuard"
		if esVersion.Spec.AuthPlugin == v1alpha1.ElasticsearchAuthPluginOpenDistro ||
			esVersion.Spec.AuthPlugin == v1alpha1.ElasticsearchAuthPluginSearchGuard {
			tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
				Alias:      string(ElasticsearchAdminCert),
				SecretName: e.CertificateName(ElasticsearchAdminCert),
				Subject: &kmapi.X509Subject{
					Organizations: []string{KubeDBOrganization},
				},
			})
		}

		// Set missing metrics-exporter certificate, if monitoring is enabled.
		if e.Spec.Monitor != nil {
			// matrics-exporter
			tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
				Alias:      string(ElasticsearchMetricsExporterCert),
				SecretName: e.CertificateName(ElasticsearchMetricsExporterCert),
				Subject: &kmapi.X509Subject{
					Organizations: []string{KubeDBOrganization},
				},
			})
		}

		// archiver
		tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchArchiverCert),
			SecretName: e.CertificateName(ElasticsearchArchiverCert),
			Subject: &kmapi.X509Subject{
				Organizations: []string{KubeDBOrganization},
			},
		})
	}

	// Force overwrite the private key encoding type to PKCS#8
	for id := range tlsConfig.Certificates {
		tlsConfig.Certificates[id].PrivateKey = &kmapi.CertificatePrivateKey{
			Encoding: kmapi.PKCS8,
		}
	}

	e.Spec.TLS = tlsConfig
}

func (e *Elasticsearch) GetMatchExpressions() []metav1.LabelSelectorRequirement {
	if e.Spec.Topology == nil {
		return nil
	}

	return []metav1.LabelSelectorRequirement{
		{
			Key:      fmt.Sprintf("${%s}", ElasticsearchNodeAffinityTemplateVar),
			Operator: metav1.LabelSelectorOpExists,
		},
	}
}

func (e *Elasticsearch) GetPersistentSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	// Add Admin/Elastic user secret name
	if e.Spec.AuthSecret != nil {
		secrets = append(secrets, e.Spec.AuthSecret.Name)
	}

	// Skip for Admin/Elastic User.
	// Add other user cred secret names.
	if e.Spec.InternalUsers != nil {
		for user := range e.Spec.InternalUsers {
			if user == string(ElasticsearchInternalUserAdmin) || user == string(ElasticsearchInternalUserElastic) {
				continue
			}
			secrets = append(secrets, e.UserCredSecretName(user))
		}
	}
	return secrets
}

func (e *Elasticsearch) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	if e.Spec.Topology != nil {
		expectedItems = 3
	}
	return checkReplicas(lister.StatefulSets(e.Namespace), labels.SelectorFromSet(e.OffshootLabels()), expectedItems)
}
