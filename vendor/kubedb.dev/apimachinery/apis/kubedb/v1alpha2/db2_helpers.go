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
	"strings"

	"kubedb.dev/apimachinery/apis"
	catalogv1alpha1 "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"kmodules.xyz/client-go/apiextensions"
	metautil "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (d *DB2) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDB2))
}

func (d *DB2) ResourcePlural() string {
	return ResourcePluralDB2
}

func (d *DB2) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", d.ResourcePlural(), SchemeGroupVersion.Group)
}

func (d *DB2) OffshootName() string {
	return d.Name
}

func (d *DB2) ServiceName() string {
	return d.OffshootName()
}

func (d *DB2) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(d.Spec.ServiceTemplates, ServiceAlias(alias))
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (d *DB2) PetSetName() string {
	return d.OffshootName()
}

func (d *DB2) GoverningServiceName() string {
	return metautil.NameWithSuffix(d.ServiceName(), "pods")
}

func (d *DB2) ResourceKind() string {
	return ResourceKindDB2
}

func (d *DB2) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		metautil.NameLabelKey:      d.ResourceFQN(),
		metautil.InstanceLabelKey:  d.Name,
		metautil.ManagedByLabelKey: SchemeGroupVersion.Group,
	}
	return metautil.OverwriteKeys(selector, extraSelectors...)
}

func (d *DB2) OffshootLabels() map[string]string {
	return d.offshootLabels(d.OffshootSelectors(), nil)
}

func (d *DB2) offshootLabels(selector, override map[string]string) map[string]string {
	selector[metautil.ComponentLabelKey] = kubedb.ComponentDatabase
	return metautil.FilterKeys(SchemeGroupVersion.Group, selector, metautil.OverwriteKeys(nil, d.Labels, override))
}

type db2App struct {
	*DB2
}

func (d db2App) Name() string {
	return d.DB2.Name
}

func (d db2App) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularDB2))
}

func (d *DB2) AppBindingMeta() appcat.AppBindingMeta {
	return &db2App{d}
}

type db2StatsService struct {
	*DB2
}

func (d db2StatsService) GetNamespace() string {
	return d.DB2.GetNamespace()
}

func (d db2StatsService) ServiceName() string {
	return d.OffshootName() + "-stats"
}

func (d db2StatsService) ServiceMonitorName() string {
	return d.ServiceName()
}

func (d db2StatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return d.OffshootLabels()
}

func (d db2StatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (d db2StatsService) Scheme() string {
	sc := promapi.SchemeHTTP
	return sc.String()
}

func (d db2StatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (d *DB2) StatsService() mona.StatsAccessor {
	return &db2StatsService{d}
}

func (d *DB2) StatsServiceLabels() map[string]string {
	return d.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (d *DB2) PodLabels(podTemplate *ofst.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if podTemplate != nil && podTemplate.Labels != nil {
		return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), podTemplate.Labels)
	}
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), nil)
}

func (d *DB2) PodControllerLabels(podTemplate *ofst.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if podTemplate != nil && podTemplate.Controller.Labels != nil {
		return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), podTemplate.Controller.Labels)
	}
	return d.offshootLabels(metautil.OverwriteKeys(d.OffshootSelectors(), extraLabels...), nil)
}

func (d *DB2) GetAuthSecretName() string {
	if d.Spec.AuthSecret != nil && d.Spec.AuthSecret.Name != "" {
		return d.Spec.AuthSecret.Name
	}
	return metautil.NameWithSuffix(d.OffshootName(), "auth")
}

func (d *DB2) GetPersistentSecrets() []string {
	var secrets []string
	if !IsVirtualAuthSecretReferred(d.Spec.AuthSecret) && d.Spec.AuthSecret != nil && d.Spec.AuthSecret.Name != "" {
		secrets = append(secrets, d.GetAuthSecretName())
	}
	return secrets
}

// StandbyServiceName is the Service selecting standby pods. Only created when
// HADR is enabled.
func (d *DB2) StandbyServiceName() string {
	return fmt.Sprintf("%s-%s", d.OffshootName(), kubedb.DB2StandbyServiceSuffix)
}

// PodName returns the PetSet pod name for an ordinal.
func (d *DB2) PodName(ordinal int) string {
	return fmt.Sprintf("%s-%d", d.OffshootName(), ordinal)
}

// PodFQDN is the stable per-pod DNS name via the governing headless Service.
// HADR_LOCAL_HOST and HADR_REMOTE_HOST must both use this form: DB2 binds
// HADR_LOCAL_HOST's resolved address specifically, and a ClusterIP Service is
// rejected by the HADR handshake even though it is TCP-reachable.
func (d *DB2) PodFQDN(ordinal int) string {
	return fmt.Sprintf("%s.%s.%s.svc", d.PodName(ordinal), d.GoverningServiceName(), d.Namespace)
}

// HADRLeaseName is the coordination.k8s.io Lease the primary pod holds. See
// kubedb.DB2HADRLeaseSuffix.
func (d *DB2) HADRLeaseName() string {
	return fmt.Sprintf("%s-%s", d.OffshootName(), kubedb.DB2HADRLeaseSuffix)
}

// IsClustered reports whether this DB2 runs HADR. Replica count alone decides it:
// 1 is standalone, more is a cluster.
func (d *DB2) IsClustered() bool {
	return d.Spec.Replicas != nil && *d.Spec.Replicas > 1
}

// HADRDatabases returns the databases HADR protects, in order. The first entry
// is the anchor: the operator elects the primary pod from it, and every other
// database is kept primary on the same pod.
//
// It never returns an empty list. Without spec.hadr it returns the image's
// default database, which is what archive logging is kept on for a standalone
// DB2. Entries whose port has not been defaulted yet get the port they would be
// defaulted to, so callers never see port 0.
func (d *DB2) HADRDatabases() []DB2HADRDatabase {
	var dbs []DB2HADRDatabase
	switch {
	case d.Spec.HADR != nil && len(d.Spec.HADR.Databases) > 0:
		dbs = make([]DB2HADRDatabase, len(d.Spec.HADR.Databases))
		copy(dbs, d.Spec.HADR.Databases)
	case d.Spec.HADR != nil && d.Spec.HADR.DatabaseName != "":
		dbs = []DB2HADRDatabase{{Name: d.Spec.HADR.DatabaseName}}
	default:
		dbs = []DB2HADRDatabase{{Name: kubedb.DB2DefaultDatabase}}
	}
	for i := range dbs {
		dbs[i].Name = strings.ToUpper(dbs[i].Name)
	}
	AssignDB2HADRPorts(dbs)
	return dbs
}

// HADRDatabaseName is the anchor database: the first entry of HADRDatabases.
func (d *DB2) HADRDatabaseName() string {
	return d.HADRDatabases()[0].Name
}

// AssignDB2HADRPorts gives every entry without a port the lowest free port from
// kubedb.DB2HadrPort upwards. Ports that are already set are never changed, so
// adding or removing a database does not move any other database's port.
func AssignDB2HADRPorts(dbs []DB2HADRDatabase) {
	used := map[int32]bool{}
	for _, db := range dbs {
		if db.Port != 0 {
			used[db.Port] = true
		}
	}
	next := int32(kubedb.DB2HadrPort)
	for i := range dbs {
		if dbs[i].Port != 0 {
			continue
		}
		for used[next] {
			next++
		}
		dbs[i].Port = next
		used[next] = true
	}
}

func (d *DB2) Finalizer() string {
	return fmt.Sprintf("%s/%s", apis.Finalizer, d.ResourceSingular())
}

func (d *DB2) ResourceSingular() string {
	return ResourceSingularDB2
}

func (d *DB2) SetDefaults(kc client.Client) {
	if d.Spec.DeletionPolicy == "" {
		d.Spec.DeletionPolicy = DeletionPolicyDelete
	}
	if d.Spec.StorageType == "" {
		d.Spec.StorageType = StorageTypeDurable
	}
	if d.Spec.Replicas == nil {
		d.Spec.Replicas = ptr.To(int32(1))
	}
	d.setHADRDefaults()
	d.initializePodTemplates()
	db2Version := &catalogv1alpha1.DB2Version{}
	err := kc.Get(context.Background(), types.NamespacedName{Name: d.Spec.Version}, db2Version)
	if err != nil {
		klog.Errorf("Failed to get database version %s: %s", err.Error(), d.Spec.Version)
		return
	}
	apis.SetDefaultResizePolicy(d.Spec.PodTemplate.Spec.Containers, d.Spec.PodTemplate.Spec.InitContainers)
}

// setHADRDefaults turns the deprecated single databaseName into a one-entry
// list, upper-cases names the way Db2 reports them, and assigns ports. It runs
// before anything that can fail, so a missing DB2Version never leaves the list
// half-defaulted.
func (d *DB2) setHADRDefaults() {
	if d.Spec.HADR == nil {
		return
	}
	if len(d.Spec.HADR.Databases) == 0 {
		name := d.Spec.HADR.DatabaseName
		if name == "" {
			name = kubedb.DB2DefaultDatabase
		}
		d.Spec.HADR.Databases = []DB2HADRDatabase{{Name: name}}
	}
	for i := range d.Spec.HADR.Databases {
		d.Spec.HADR.Databases[i].Name = strings.ToUpper(d.Spec.HADR.Databases[i].Name)
	}
	AssignDB2HADRPorts(d.Spec.HADR.Databases)
}

func (d *DB2) initializePodTemplates() {
	if d.Spec.PodTemplate == nil {
		d.Spec.PodTemplate = new(ofst.PodTemplateSpec)
	}
}

func (d *DB2) SetHealthCheckerDefaults() {
	if d.Spec.HealthChecker.PeriodSeconds == nil {
		d.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if d.Spec.HealthChecker.TimeoutSeconds == nil {
		d.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if d.Spec.HealthChecker.FailureThreshold == nil {
		d.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (d *DB2) GetDeletionPolicy() string {
	return string(d.Spec.DeletionPolicy)
}

func (d *DB2) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(d, SchemeGroupVersion.WithKind(d.ResourceKind()))
}
