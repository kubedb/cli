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
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/appscode/go/types"
	core "k8s.io/api/core/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ PerconaXtraDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPerconaXtraDB))
}

var _ apis.ResourceInfo = &PerconaXtraDB{}

func (p PerconaXtraDB) OffshootName() string {
	return p.Name
}

func (p PerconaXtraDB) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindPerconaXtraDB,
	}
}

func (p PerconaXtraDB) OffshootLabels() map[string]string {
	out := p.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularPerconaXtraDB
	out[meta_util.VersionLabelKey] = string(p.Spec.Version)
	out[meta_util.InstanceLabelKey] = p.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = kubedb.GroupName
	return meta_util.FilterKeys(kubedb.GroupName, out, p.Labels)
}

func (p PerconaXtraDB) ResourceShortCode() string {
	return ResourceCodePerconaXtraDB
}

func (p PerconaXtraDB) ResourceKind() string {
	return ResourceKindPerconaXtraDB
}

func (p PerconaXtraDB) ResourceSingular() string {
	return ResourceSingularPerconaXtraDB
}

func (p PerconaXtraDB) ResourcePlural() string {
	return ResourcePluralPerconaXtraDB
}

func (p PerconaXtraDB) ServiceName() string {
	return p.OffshootName()
}

func (p PerconaXtraDB) IsCluster() bool {
	return types.Int32(p.Spec.Replicas) > 1
}

func (p PerconaXtraDB) GoverningServiceName() string {
	return p.OffshootName() + "-gvr"
}

func (p PerconaXtraDB) PeerName(idx int) string {
	return fmt.Sprintf("%s-%d.%s.%s", p.OffshootName(), idx, p.GoverningServiceName(), p.Namespace)
}

func (p PerconaXtraDB) GetDatabaseSecretName() string {
	return p.Spec.DatabaseSecret.SecretName
}

func (p PerconaXtraDB) ClusterName() string {
	return p.OffshootName()
}

type perconaXtraDBApp struct {
	*PerconaXtraDB
}

func (p perconaXtraDBApp) Name() string {
	return p.PerconaXtraDB.Name
}

func (p perconaXtraDBApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPerconaXtraDB))
}

func (p PerconaXtraDB) AppBindingMeta() appcat.AppBindingMeta {
	return &perconaXtraDBApp{&p}
}

type perconaXtraDBStatsService struct {
	*PerconaXtraDB
}

func (p perconaXtraDBStatsService) GetNamespace() string {
	return p.PerconaXtraDB.GetNamespace()
}

func (p perconaXtraDBStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p perconaXtraDBStatsService) ServiceMonitorName() string {
	return p.ServiceName()
}

func (p perconaXtraDBStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return p.OffshootLabels()
}

func (p perconaXtraDBStatsService) Path() string {
	return DefaultStatsPath
}

func (p perconaXtraDBStatsService) Scheme() string {
	return ""
}

func (p PerconaXtraDB) StatsService() mona.StatsAccessor {
	return &perconaXtraDBStatsService{&p}
}

func (p PerconaXtraDB) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, p.OffshootSelectors(), p.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (p *PerconaXtraDB) GetMonitoringVendor() string {
	if p.Spec.Monitor != nil {
		return p.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (p *PerconaXtraDB) SetDefaults() {
	if p == nil {
		return
	}

	if p.Spec.Replicas == nil {
		p.Spec.Replicas = types.Int32P(1)
	}

	if p.Spec.StorageType == "" {
		p.Spec.StorageType = StorageTypeDurable
	}
	if p.Spec.TerminationPolicy == "" {
		p.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	p.Spec.setDefaultProbes()
	p.Spec.Monitor.SetDefaults()
}

// setDefaultProbes sets defaults only when probe fields are nil.
// In operator, check if the value of probe fields is "{}".
// For "{}", ignore readinessprobe or livenessprobe in statefulset.
// Ref: https://github.com/mattlord/Docker-InnoDB-Cluster/blob/master/healthcheck.sh#L10
func (p *PerconaXtraDBSpec) setDefaultProbes() {
	if p == nil {
		return
	}

	var readynessProbeCmd []string
	if types.Int32(p.Replicas) > 1 {
		readynessProbeCmd = []string{
			"/cluster-check.sh",
		}
	} else {
		readynessProbeCmd = []string{
			"bash",
			"-c",
			`export MYSQL_PWD="${MYSQL_ROOT_PASSWORD}"
ping_resp=$(mysqladmin -uroot ping)
if [[ "$ping_resp" != "mysqld is alive" ]]; then
    echo "[ERROR] server is not ready. PING_RESPONSE: $ping_resp"
    exit 1
fi
`,
		}
	}

	readinessProbe := &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: readynessProbeCmd,
			},
		},
		InitialDelaySeconds: 30,
		PeriodSeconds:       10,
	}
	if p.PodTemplate.Spec.ReadinessProbe == nil {
		p.PodTemplate.Spec.ReadinessProbe = readinessProbe
	}
}

func (p *PerconaXtraDBSpec) GetSecrets() []string {
	if p == nil {
		return nil
	}

	var secrets []string
	if p.DatabaseSecret != nil {
		secrets = append(secrets, p.DatabaseSecret.SecretName)
	}
	return secrets
}
