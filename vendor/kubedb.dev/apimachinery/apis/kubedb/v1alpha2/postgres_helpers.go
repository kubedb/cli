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

	"gomodules.xyz/pointer"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
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

func (p *Postgres) SetDefaults() {
	if p == nil {
		return
	}

	if p.Spec.StorageType == "" {
		p.Spec.StorageType = StorageTypeDurable
	}
	if p.Spec.TerminationPolicy == "" {
		p.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if p.Spec.Init != nil && p.Spec.Init.PostgresWAL != nil && p.Spec.Init.PostgresWAL.PITR != nil {
		pitr := p.Spec.Init.PostgresWAL.PITR

		if pitr.TargetInclusive == nil {
			pitr.TargetInclusive = pointer.BoolP(true)
		}

		p.Spec.Init.PostgresWAL.PITR = pitr
	}

	if p.Spec.LeaderElection == nil {
		// Default values: https://github.com/kubernetes/apiserver/blob/e85ad7b666fef0476185731329f4cff1536efff8/pkg/apis/config/v1alpha1/defaults.go#L26-L52
		p.Spec.LeaderElection = &LeaderElectionConfig{
			LeaseDurationSeconds: 15,
			RenewDeadlineSeconds: 10,
			RetryPeriodSeconds:   2,
		}
	}

	if p.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		p.Spec.PodTemplate.Spec.ServiceAccountName = p.OffshootName()
	}

	p.Spec.Monitor.SetDefaults()
	setDefaultResourceLimits(&p.Spec.PodTemplate.Spec.Resources, defaultResourceLimits, defaultResourceLimits)
}

func (e *PostgresSpec) GetPersistentSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	if e.AuthSecret != nil {
		secrets = append(secrets, e.AuthSecret.Name)
	}
	return secrets
}

func (p *Postgres) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(p.Namespace), labels.SelectorFromSet(p.OffshootLabels()), expectedItems)
}
