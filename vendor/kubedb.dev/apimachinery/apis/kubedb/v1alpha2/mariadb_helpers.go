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

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/appscode/go/types"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ MariaDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMariaDB))
}

var _ apis.ResourceInfo = &MariaDB{}

func (m MariaDB) OffshootName() string {
	return m.Name
}

func (m MariaDB) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: m.Name,
		LabelDatabaseKind: ResourceKindMariaDB,
	}
}

func (m MariaDB) OffshootLabels() map[string]string {
	out := m.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularMariaDB
	out[meta_util.VersionLabelKey] = string(m.Spec.Version)
	out[meta_util.InstanceLabelKey] = m.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = kubedb.GroupName
	return meta_util.FilterKeys(kubedb.GroupName, out, m.Labels)
}

func (m MariaDB) ResourceShortCode() string {
	return ResourceCodeMariaDB
}

func (m MariaDB) ResourceKind() string {
	return ResourceKindMariaDB
}

func (m MariaDB) ResourceSingular() string {
	return ResourceSingularMariaDB
}

func (m MariaDB) ResourcePlural() string {
	return ResourcePluralMariaDB
}

func (m MariaDB) ServiceName() string {
	return m.OffshootName()
}

func (m MariaDB) GoverningServiceName() string {
	return m.OffshootName() + "-gvr"
}

type mariadbApp struct {
	*MariaDB
}

func (m mariadbApp) Name() string {
	return m.MariaDB.Name
}

func (m mariadbApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMariaDB))
}

func (m MariaDB) AppBindingMeta() appcat.AppBindingMeta {
	return &mariadbApp{&m}
}

type mariadbStatsService struct {
	*MariaDB
}

func (m mariadbStatsService) GetNamespace() string {
	return m.MariaDB.GetNamespace()
}

func (m mariadbStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m mariadbStatsService) ServiceMonitorName() string {
	return m.ServiceName()
}

func (m mariadbStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return m.OffshootLabels()
}

func (m mariadbStatsService) Path() string {
	return DefaultStatsPath
}

func (m mariadbStatsService) Scheme() string {
	return ""
}

func (m MariaDB) StatsService() mona.StatsAccessor {
	return &mariadbStatsService{&m}
}

func (m MariaDB) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, m.OffshootSelectors(), m.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (m *MariaDB) SetDefaults() {
	if m == nil {
		return
	}
	if m.Spec.Replicas == nil {
		m.Spec.Replicas = types.Int32P(1)
	}

	// perform defaulting

	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}
	if m.Spec.TerminationPolicy == "" {
		m.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	m.Spec.Monitor.SetDefaults()
}

func (m *MariaDBSpec) GetPersistentSecrets() []string {
	if m == nil {
		return nil
	}

	var secrets []string
	if m.AuthSecret != nil {
		secrets = append(secrets, m.AuthSecret.Name)
	}
	return secrets
}

func (m *MariaDB) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(m.Namespace), labels.SelectorFromSet(m.OffshootLabels()), expectedItems)
}
