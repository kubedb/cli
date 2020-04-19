/*
Copyright The KubeDB Authors.

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

	"kubedb.dev/apimachinery/api/crds"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ MySQL) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMySQL))
}

var _ apis.ResourceInfo = &MySQL{}

func (m MySQL) OffshootName() string {
	return m.Name
}

func (m MySQL) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: m.Name,
		LabelDatabaseKind: ResourceKindMySQL,
	}
}

func (m MySQL) OffshootLabels() map[string]string {
	out := m.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularMySQL
	out[meta_util.VersionLabelKey] = string(m.Spec.Version)
	out[meta_util.InstanceLabelKey] = m.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, m.Labels)
}

func (m MySQL) ResourceShortCode() string {
	return ResourceCodeMySQL
}

func (m MySQL) ResourceKind() string {
	return ResourceKindMySQL
}

func (m MySQL) ResourceSingular() string {
	return ResourceSingularMySQL
}

func (m MySQL) ResourcePlural() string {
	return ResourcePluralMySQL
}

func (m MySQL) ServiceName() string {
	return m.OffshootName()
}

func (m MySQL) GoverningServiceName() string {
	return m.OffshootName() + "-gvr"
}

func (m MySQL) PeerName(idx int) string {
	return fmt.Sprintf("%s-%d.%s.%s", m.OffshootName(), idx, m.GoverningServiceName(), m.Namespace)
}

func (m MySQL) GetDatabaseSecretName() string {
	return m.Spec.DatabaseSecret.SecretName
}

type mysqlApp struct {
	*MySQL
}

func (r mysqlApp) Name() string {
	return r.MySQL.Name
}

func (r mysqlApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMySQL))
}

func (m MySQL) AppBindingMeta() appcat.AppBindingMeta {
	return &mysqlApp{&m}
}

type mysqlStatsService struct {
	*MySQL
}

func (m mysqlStatsService) GetNamespace() string {
	return m.MySQL.GetNamespace()
}

func (m mysqlStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m mysqlStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", m.Namespace, m.Name)
}

func (m mysqlStatsService) Path() string {
	return DefaultStatsPath
}

func (m mysqlStatsService) Scheme() string {
	return ""
}

func (m MySQL) StatsService() mona.StatsAccessor {
	return &mysqlStatsService{&m}
}

func (m MySQL) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, m.OffshootSelectors(), m.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (m *MySQL) GetMonitoringVendor() string {
	if m.Spec.Monitor != nil {
		return m.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (m *MySQL) SetDefaults() {
	if m == nil {
		return
	}
	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}
	if m.Spec.UpdateStrategy.Type == "" {
		m.Spec.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if m.Spec.TerminationPolicy == "" {
		m.Spec.TerminationPolicy = TerminationPolicyDelete
	} else if m.Spec.TerminationPolicy == TerminationPolicyPause {
		m.Spec.TerminationPolicy = TerminationPolicyHalt
	}

	if m.Spec.Topology != nil && m.Spec.Topology.Mode != nil && *m.Spec.Topology.Mode == MySQLClusterModeGroup {
		if m.Spec.Replicas == nil {
			m.Spec.Replicas = types.Int32P(MySQLDefaultGroupSize)
		}
		m.Spec.setDefaultProbes()
	} else {
		if m.Spec.Replicas == nil {
			m.Spec.Replicas = types.Int32P(1)
		}
	}

	if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
	}

	m.Spec.Monitor.SetDefaults()
}

// setDefaultProbes sets defaults only when probe fields are nil.
// In operator, check if the value of probe fields is "{}".
// For "{}", ignore readinessprobe or livenessprobe in statefulset.
// Ref: https://github.com/mattlord/Docker-InnoDB-Cluster/blob/master/healthcheck.sh#L10
func (m *MySQLSpec) setDefaultProbes() {
	probe := &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: []string{
					"bash",
					"-c",
					`
export MYSQL_PWD=${MYSQL_ROOT_PASSWORD}
mysql -h localhost -nsLNE -e "select member_state from performance_schema.replication_group_members where member_id=@@server_uuid;" 2>/dev/null | grep -v "*" | egrep -v "ERROR|OFFLINE"
`,
				},
			},
		},
		InitialDelaySeconds: 30,
		PeriodSeconds:       5,
	}

	if m.PodTemplate.Spec.LivenessProbe == nil {
		m.PodTemplate.Spec.LivenessProbe = probe
	}
	if m.PodTemplate.Spec.ReadinessProbe == nil {
		m.PodTemplate.Spec.ReadinessProbe = probe
	}
}

func (e *MySQLSpec) GetSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	if e.DatabaseSecret != nil {
		secrets = append(secrets, e.DatabaseSecret.SecretName)
	}
	return secrets
}
