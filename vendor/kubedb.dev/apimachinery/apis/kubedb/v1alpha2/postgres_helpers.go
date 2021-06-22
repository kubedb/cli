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
	"math"
	"strconv"
	"time"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

func (_ Postgres) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPostgres))
}

var _ apis.ResourceInfo = &Postgres{}

func (p Postgres) OffshootName() string {
	return p.Name
}

func (p Postgres) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      p.ResourceFQN(),
		meta_util.InstanceLabelKey:  p.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (p Postgres) OffshootLabels() map[string]string {
	out := p.OffshootSelectors()
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, out, p.Labels)
}

func (p Postgres) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPostgres, kubedb.GroupName)
}

func (p Postgres) ResourceShortCode() string {
	return ResourceCodePostgres
}

func (p Postgres) ResourceKind() string {
	return ResourceKindPostgres
}

func (p Postgres) ResourceSingular() string {
	return ResourceSingularPostgres
}

func (p Postgres) ResourcePlural() string {
	return ResourcePluralPostgres
}

func (p Postgres) ServiceName() string {
	return p.OffshootName()
}

func (p Postgres) StandbyServiceName() string {
	return meta_util.NameWithPrefix(p.ServiceName(), "standby")
}

func (p Postgres) GoverningServiceName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "pods")
}

type postgresApp struct {
	*Postgres
}

func (r postgresApp) Name() string {
	return r.Postgres.Name
}

func (r postgresApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPostgres))
}

func (p Postgres) AppBindingMeta() appcat.AppBindingMeta {
	return &postgresApp{&p}
}

type postgresStatsService struct {
	*Postgres
}

func (p postgresStatsService) GetNamespace() string {
	return p.Postgres.GetNamespace()
}

func (p postgresStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p postgresStatsService) ServiceMonitorName() string {
	return p.ServiceName()
}

func (p postgresStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return p.OffshootLabels()
}

func (p postgresStatsService) Path() string {
	return DefaultStatsPath
}

func (p postgresStatsService) Scheme() string {
	return ""
}

func (p Postgres) StatsService() mona.StatsAccessor {
	return &postgresStatsService{&p}
}

func (p Postgres) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, p.OffshootSelectors(), p.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (p *Postgres) SetDefaults(postgresVersion *catalog.PostgresVersion, topology *core_util.Topology) {
	if p == nil {
		return
	}

	if p.Spec.StorageType == "" {
		p.Spec.StorageType = StorageTypeDurable
	}
	if p.Spec.TerminationPolicy == "" {
		p.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if p.Spec.LeaderElection == nil {
		p.Spec.LeaderElection = &PostgreLeaderElectionConfig{
			//we have set this default to 33554432. if the difference between primary and replica is more then this,
			//the replica node is going to manually sync itself.
			Period:                   metav1.Duration{Duration: 100 * time.Millisecond},
			MaximumLagBeforeFailover: 32 * 1024 * 1024,
			ElectionTick:             10,
			HeartbeatTick:            1,
		}
	}

	if p.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		p.Spec.PodTemplate.Spec.ServiceAccountName = p.OffshootName()
	}

	if p.Spec.TLS != nil {
		if p.Spec.SSLMode == "" {
			p.Spec.SSLMode = PostgresSSLModeVerifyFull
		}
		if p.Spec.ClientAuthMode == "" {
			p.Spec.ClientAuthMode = ClientAuthModeMD5
		}
	} else {
		if p.Spec.SSLMode == "" {
			p.Spec.SSLMode = PostgresSSLModeDisable
		}
		if p.Spec.ClientAuthMode == "" {
			p.Spec.ClientAuthMode = ClientAuthModeMD5
		}
	}

	if p.Spec.PodTemplate.Spec.ContainerSecurityContext == nil {
		p.Spec.PodTemplate.Spec.ContainerSecurityContext = &core.SecurityContext{
			RunAsUser:  postgresVersion.Spec.SecurityContext.RunAsUser,
			RunAsGroup: postgresVersion.Spec.SecurityContext.RunAsUser,
			Privileged: pointer.BoolP(false),
			Capabilities: &core.Capabilities{
				Add: []core.Capability{"IPC_LOCK", "SYS_RESOURCE"},
			},
		}
	} else {
		if p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsUser == nil {
			p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsUser = postgresVersion.Spec.SecurityContext.RunAsUser
		}
		if p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup == nil {
			p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup = p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsUser
		}
	}

	if p.Spec.PodTemplate.Spec.SecurityContext == nil {
		p.Spec.PodTemplate.Spec.SecurityContext = &core.PodSecurityContext{
			RunAsUser:  p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsUser,
			RunAsGroup: p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup,
		}
	} else {
		if p.Spec.PodTemplate.Spec.SecurityContext.RunAsUser == nil {
			p.Spec.PodTemplate.Spec.SecurityContext.RunAsUser = p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsUser
		}
		if p.Spec.PodTemplate.Spec.SecurityContext.RunAsGroup == nil {
			p.Spec.PodTemplate.Spec.SecurityContext.RunAsGroup = p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup
		}
	}
	// Need to set FSGroup equal to  p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup.
	// So that /var/pv directory have the group permission for the RunAsGroup user GID.
	// Otherwise, We will get write permission denied.
	p.Spec.PodTemplate.Spec.SecurityContext.FSGroup = p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup

	p.Spec.Monitor.SetDefaults()
	p.SetTLSDefaults()
	SetDefaultResourceLimits(&p.Spec.PodTemplate.Spec.Resources, DefaultResources)
	p.setDefaultAffinity(&p.Spec.PodTemplate, p.OffshootSelectors(), topology)
}

// setDefaultAffinity
func (p *Postgres) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
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
						Namespaces: []string{p.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: labels,
						},
						TopologyKey: core.LabelHostname,
					},
				},
				// Prefer to not schedule multiple pods on the node with same zone
				{
					Weight: 50,
					PodAffinityTerm: core.PodAffinityTerm{
						Namespaces: []string{p.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: labels,
						},
						TopologyKey: topology.LabelZone,
					},
				},
			},
		},
	}
}

func (p *Postgres) SetTLSDefaults() {
	if p.Spec.TLS == nil || p.Spec.TLS.IssuerRef == nil {
		return
	}
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PostgresServerCert), p.CertificateName(PostgresServerCert))
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PostgresClientCert), p.CertificateName(PostgresClientCert))
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PostgresMetricsExporterCert), p.CertificateName(PostgresMetricsExporterCert))
}

func (p *PostgresSpec) GetPersistentSecrets() []string {
	if p == nil {
		return nil
	}

	var secrets []string
	if p.AuthSecret != nil {
		secrets = append(secrets, p.AuthSecret.Name)
	}
	return secrets
}

func (p *Postgres) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(p.Namespace), labels.SelectorFromSet(p.OffshootLabels()), expectedItems)
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (p *Postgres) CertificateName(alias PostgresCertificateAlias) string {
	return meta_util.NameWithSuffix(p.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any provide,
// otherwise returns default certificate secret name for the given alias.
func (p *Postgres) GetCertSecretName(alias PostgresCertificateAlias) string {
	if p.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(p.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return p.CertificateName(alias)
}

// GetSharedBufferSizeForPostgres this func takes a input type int64 which is in bytes
// return the 25% of the input in Bytes, KiloBytes, MegaBytes, GigaBytes, or TeraBytes
func GetSharedBufferSizeForPostgres(resource *resource.Quantity) string {
	// no more than 25% of main memory (RAM)
	minSharedBuffer := int64(128 * 1024 * 1024)
	ret := minSharedBuffer
	if resource != nil {
		ret = (resource.Value() / 100) * 25
	}
	// the shared buffer value can't be less then this
	//128 MB  is the minimum
	if ret < minSharedBuffer {
		ret = minSharedBuffer
	}

	sharedBuffer := ConvertBytesInMB(ret)
	return sharedBuffer
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	// this func take a float and return the int and fractional part separately
	// math.modf(100.4) will return int part = 100 and fractional part = 0.40000000000000000
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return newVal
}

// ConvertBytesInMB this func takes a input type int64 which is in bytes
// return the input in Bytes, KiloBytes, MegaBytes, GigaBytes, or TeraBytes
func ConvertBytesInMB(value int64) string {
	var suffixes [5]string
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	// here base is the type we are going to represent the value in string
	// if base is 2 then we will represent the value in MB.
	// if base is 0 then represent the value in B.
	if value == 0 {
		return "0B"
	}
	base := math.Log(float64(value)) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]

	valueMB := strconv.FormatFloat(getSize, 'f', -1, 64) + string(getSuffix)
	return valueMB
}
