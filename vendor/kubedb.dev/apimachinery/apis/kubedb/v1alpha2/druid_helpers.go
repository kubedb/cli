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
	"strconv"
	"strings"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/Masterminds/semver/v3"
	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (d *Druid) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDruid))
}

func (d *Druid) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(d, SchemeGroupVersion.WithKind(d.ResourceKind()))
}

func (d *Druid) ResourceKind() string {
	return ResourceKindDruid
}

func (d *Druid) ResourceSingular() string {
	return ResourceSingularDruid
}

func (d *Druid) ResourcePlural() string {
	return ResourcePluralDruid
}

func (d *Druid) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", d.ResourcePlural(), kubedb.GroupName)
}

func (d *Druid) OffShootName() string {
	return d.Name
}

func (d *Druid) ServiceName() string {
	return d.OffShootName()
}

func (d *Druid) CoordinatorsServiceName() string {
	return meta_util.NameWithSuffix(d.ServiceName(), "coordinators")
}

func (d *Druid) OverlordsServiceName() string {
	return meta_util.NameWithSuffix(d.ServiceName(), "overlords")
}

func (d *Druid) BrokersServiceName() string {
	return meta_util.NameWithSuffix(d.ServiceName(), "brokers")
}

func (d *Druid) RoutersServiceName() string {
	return meta_util.NameWithSuffix(d.ServiceName(), "routers")
}

func (d *Druid) GoverningServiceName() string {
	return meta_util.NameWithSuffix(d.ServiceName(), "pods")
}

func (d *Druid) OffShootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      d.ResourceFQN(),
		meta_util.InstanceLabelKey:  d.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (d *Druid) offShootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, d.Labels, override))
}

func (d *Druid) OffShootLabels() map[string]string {
	return d.offShootLabels(d.OffShootSelectors(), nil)
}

func (d *Druid) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(d.Spec.ServiceTemplates, alias)
	return d.offShootLabels(meta_util.OverwriteKeys(d.OffShootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (r *Druid) Finalizer() string {
	return fmt.Sprintf("%s/%s", apis.Finalizer, r.ResourceSingular())
}

func (d *Druid) DefaultUserCredSecretName(username string) string {
	return meta_util.NameWithSuffix(d.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", username), "_", "-"))
}

type DruidStatsService struct {
	*Druid
}

func (ks DruidStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (ks DruidStatsService) GetNamespace() string {
	return ks.Druid.GetNamespace()
}

func (ks DruidStatsService) ServiceName() string {
	return ks.OffShootName() + "-stats"
}

func (ks DruidStatsService) ServiceMonitorName() string {
	return ks.ServiceName()
}

func (ks DruidStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return ks.OffshootLabels()
}

func (ks DruidStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (ks DruidStatsService) Scheme() string {
	return ""
}

func (d *Druid) StatsService() mona.StatsAccessor {
	return &DruidStatsService{d}
}

func (d *Druid) StatsServiceLabels() map[string]string {
	return d.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (d *Druid) ConfigSecretName() string {
	return meta_util.NameWithSuffix(d.OffShootName(), "config")
}

func (d *Druid) PetSetName(nodeRole DruidNodeRoleType) string {
	return meta_util.NameWithSuffix(d.OffShootName(), d.DruidNodeRoleString(nodeRole))
}

func (d *Druid) PodLabels(nodeType DruidNodeRoleType, extraLabels ...map[string]string) map[string]string {
	nodeSpec, dataNodeSpec := d.GetNodeSpec(nodeType)
	var labels map[string]string
	if nodeSpec != nil {
		labels = nodeSpec.PodTemplate.Labels
	} else {
		labels = dataNodeSpec.PodTemplate.Labels
	}
	return d.offShootLabels(meta_util.OverwriteKeys(d.OffShootSelectors(), extraLabels...), labels)
}

func (d *Druid) PodControllerLabels(nodeType DruidNodeRoleType, extraLabels ...map[string]string) map[string]string {
	nodeSpec, dataNodeSpec := d.GetNodeSpec(nodeType)
	var labels map[string]string
	if nodeSpec != nil {
		labels = nodeSpec.PodTemplate.Controller.Labels
	} else {
		labels = dataNodeSpec.PodTemplate.Controller.Labels
	}
	return d.offShootLabels(meta_util.OverwriteKeys(d.OffShootSelectors(), extraLabels...), labels)
}

func (d *Druid) GetNodeSpec(nodeType DruidNodeRoleType) (*DruidNode, *DruidDataNode) {
	if nodeType == DruidNodeRoleCoordinators {
		return d.Spec.Topology.Coordinators, nil
	} else if nodeType == DruidNodeRoleOverlords {
		return d.Spec.Topology.Overlords, nil
	} else if nodeType == DruidNodeRoleMiddleManagers {
		return nil, d.Spec.Topology.MiddleManagers
	} else if nodeType == DruidNodeRoleHistoricals {
		return nil, d.Spec.Topology.Historicals
	} else if nodeType == DruidNodeRoleBrokers {
		return d.Spec.Topology.Brokers, nil
	} else if nodeType == DruidNodeRoleRouters {
		return d.Spec.Topology.Routers, nil
	}

	panic("Node role name does not match any known types")
}

func (d *Druid) ServiceAccountName() string {
	return d.OffShootName()
}

func (d *Druid) DruidNodeRoleString(nodeRole DruidNodeRoleType) string {
	return strings.ToLower(string(nodeRole))
}

func (d *Druid) DruidNodeRoleStringSingular(nodeRole DruidNodeRoleType) string {
	singularNodeRole := string(nodeRole)[:len(nodeRole)-1]
	return singularNodeRole
}

func (d *Druid) DruidNodeContainerPort(nodeRole DruidNodeRoleType) int32 {
	if nodeRole == DruidNodeRoleCoordinators {
		return kubedb.DruidPortCoordinators
	} else if nodeRole == DruidNodeRoleOverlords {
		return kubedb.DruidPortOverlords
	} else if nodeRole == DruidNodeRoleMiddleManagers {
		return kubedb.DruidPortMiddleManagers
	} else if nodeRole == DruidNodeRoleHistoricals {
		return kubedb.DruidPortHistoricals
	} else if nodeRole == DruidNodeRoleBrokers {
		return kubedb.DruidPortBrokers
	}
	// Routers
	return kubedb.DruidPortRouters
}

func (d *Druid) SetHealthCheckerDefaults() {
	if d.Spec.HealthChecker.PeriodSeconds == nil {
		d.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(30)
	}
	if d.Spec.HealthChecker.TimeoutSeconds == nil {
		d.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if d.Spec.HealthChecker.FailureThreshold == nil {
		d.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

type DruidApp struct {
	*Druid
}

func (d DruidApp) Name() string {
	return d.Druid.Name
}

func (d DruidApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularDruid))
}

func (d *Druid) AppBindingMeta() appcat.AppBindingMeta {
	return &DruidApp{d}
}

func (d *Druid) GetConnectionScheme() string {
	scheme := "http"
	//if d.Spec.EnableSSL {
	//	scheme = "https"
	//}
	return scheme
}

func (d *Druid) GetMetadataStorageConnectURI(appbinding *appcat.AppBinding, metadataStorageType DruidMetadataStorageType) string {
	var url string
	if metadataStorageType == DruidMetadataStorageMySQL {
		url = *appbinding.Spec.ClientConfig.URL
		url = kubedb.DruidMetadataStorageConnectURIPrefixMySQL + url[4:len(url)-2] + "/" + ResourceSingularDruid
	} else if metadataStorageType == DruidMetadataStoragePostgreSQL {
		url = appbinding.Spec.ClientConfig.Service.Name + ":" + strconv.Itoa(int(appbinding.Spec.ClientConfig.Service.Port))
		url = kubedb.DruidMetadataStorageConnectURIPrefixPostgreSQL + url + "/" + ResourceSingularDruid
	}
	return url
}

func (d *Druid) GetZKServiceHost(appbinding *appcat.AppBinding) string {
	return fmt.Sprintf("%s.%s.svc:%d", appbinding.Spec.ClientConfig.Service.Name, appbinding.Namespace, int(appbinding.Spec.ClientConfig.Service.Port))
}

func (d *Druid) AddDruidExtensionLoadList(druidExtensionLoadList string, extension string) string {
	if len(druidExtensionLoadList) == 0 {
		druidExtensionLoadList += "["
	} else {
		druidExtensionLoadList = strings.TrimSuffix(druidExtensionLoadList, "]")
		druidExtensionLoadList += ", "
	}
	druidExtensionLoadList += "\"" + extension + "\"]"
	return druidExtensionLoadList
}

func (d *Druid) GetMetadataStorageType(metadataStorage string) DruidMetadataStorageType {
	if metadataStorage == string(DruidMetadataStorageMySQL) || metadataStorage == strings.ToLower(string(DruidMetadataStorageMySQL)) {
		return DruidMetadataStorageMySQL
	} else if metadataStorage == string(DruidMetadataStoragePostgreSQL) || metadataStorage == strings.ToLower(string(DruidMetadataStoragePostgreSQL)) ||
		metadataStorage == kubedb.DruidMetadataStorageTypePostgres || metadataStorage == strings.ToLower(string(kubedb.DruidMetadataStorageTypePostgres)) {
		return DruidMetadataStoragePostgreSQL
	} else {
		panic(fmt.Sprintf("Unknown metadata storage type: %s", metadataStorage))
	}
}

func (d *Druid) PVCName(alias string) string {
	return meta_util.NameWithSuffix(d.Name, alias)
}

func (d *Druid) GetDruidSegmentCacheConfig() string {
	// Update the storage size according to the druid segment cache configuration
	var storageSize string

	if d.Spec.Topology.MiddleManagers.Storage != nil {
		storageSize = d.Spec.Topology.Historicals.Storage.Resources.Requests.Storage().String()
		storageSize = d.GetDruidStorageSize(storageSize)
	} else {
		storageSize = "1g"
	}

	segmentCache := fmt.Sprintf("[{\"path\":\"%s\",\"maxSize\":\"%s\"}]", kubedb.DruidHistoricalsSegmentCacheDir, storageSize)
	return segmentCache
}

func (d *Druid) GetDruidStorageSize(storageSize string) string {
	lastTwoCharacters := storageSize[len(storageSize)-2:]
	storageSize = storageSize[:len(storageSize)-2]
	intSorageSize, _ := strconv.Atoi(storageSize)

	if lastTwoCharacters == "Gi" {
		intSorageSize *= 1000000000
	} else {
		intSorageSize *= 1000000
	}
	return strconv.Itoa(intSorageSize)
}

func (d *Druid) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      d.ResourceFQN(),
		meta_util.InstanceLabelKey:  d.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (d Druid) OffshootLabels() map[string]string {
	return d.offshootLabels(d.OffshootSelectors(), nil)
}

func (e Druid) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, e.Labels, override))
}

func (d *Druid) SetDefaults() {
	if d.Spec.DeletionPolicy == "" {
		d.Spec.DeletionPolicy = TerminationPolicyDelete
	}

	if d.Spec.DisableSecurity == nil {
		d.Spec.DisableSecurity = pointer.BoolP(false)
	}

	if !*d.Spec.DisableSecurity {
		if d.Spec.AuthSecret == nil {
			d.Spec.AuthSecret = &v1.LocalObjectReference{
				Name: d.DefaultUserCredSecretName(kubedb.DruidUserAdmin),
			}
		}
	}

	var druidVersion catalog.DruidVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name: d.Spec.Version,
	}, &druidVersion)
	if err != nil {
		klog.Errorf("failed to get the druid version object %s: %s\n", d.Spec.Version, err.Error())
		return
	}

	version, err := semver.NewVersion(druidVersion.Spec.Version)
	if err != nil {
		klog.Errorf("failed to parse druid version :%s\n", err.Error())
		return
	}

	if d.Spec.Topology != nil {
		if d.Spec.Topology.Coordinators == nil {
			d.Spec.Topology.Coordinators = &DruidNode{}
		}
		if d.Spec.Topology.Coordinators != nil {
			if d.Spec.Topology.Coordinators.Replicas == nil {
				d.Spec.Topology.Coordinators.Replicas = pointer.Int32P(1)
			}
			if version.Major() > 25 {
				if d.Spec.Topology.Coordinators.PodTemplate.Spec.SecurityContext == nil {
					d.Spec.Topology.Coordinators.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{FSGroup: druidVersion.Spec.SecurityContext.RunAsUser}
				}
				d.setDefaultContainerSecurityContext(&druidVersion, &d.Spec.Topology.Coordinators.PodTemplate)
				d.setDefaultContainerResourceLimits(&d.Spec.Topology.Coordinators.PodTemplate, DruidNodeRoleCoordinators)
			}
		}

		if d.Spec.Topology.Overlords != nil {
			if d.Spec.Topology.Overlords.Replicas == nil {
				d.Spec.Topology.Overlords.Replicas = pointer.Int32P(1)
			}
			if version.Major() > 25 {
				if d.Spec.Topology.Overlords.PodTemplate.Spec.SecurityContext == nil {
					d.Spec.Topology.Overlords.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{FSGroup: druidVersion.Spec.SecurityContext.RunAsUser}
				}
				d.setDefaultContainerSecurityContext(&druidVersion, &d.Spec.Topology.Overlords.PodTemplate)
				d.setDefaultContainerResourceLimits(&d.Spec.Topology.Overlords.PodTemplate, DruidNodeRoleOverlords)
			}
		}

		if d.Spec.Topology.MiddleManagers == nil {
			d.Spec.Topology.MiddleManagers = &DruidDataNode{}
		}
		if d.Spec.Topology.MiddleManagers != nil {
			if d.Spec.Topology.MiddleManagers.Replicas == nil {
				d.Spec.Topology.MiddleManagers.Replicas = pointer.Int32P(1)
			}
			if d.Spec.Topology.MiddleManagers.StorageType == "" {
				d.Spec.Topology.MiddleManagers.StorageType = StorageTypeDurable
			}
			if d.Spec.Topology.MiddleManagers.Storage == nil && d.Spec.Topology.MiddleManagers.StorageType == StorageTypeDurable {
				d.Spec.Topology.MiddleManagers.Storage = d.GetDefaultPVC()
			}
			if version.Major() > 25 {
				if d.Spec.Topology.MiddleManagers.PodTemplate.Spec.SecurityContext == nil {
					d.Spec.Topology.MiddleManagers.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{FSGroup: druidVersion.Spec.SecurityContext.RunAsUser}
				}
				d.setDefaultContainerSecurityContext(&druidVersion, &d.Spec.Topology.MiddleManagers.PodTemplate)
				d.setDefaultContainerResourceLimits(&d.Spec.Topology.MiddleManagers.PodTemplate, DruidNodeRoleMiddleManagers)
			}
		}

		if d.Spec.Topology.Historicals == nil {
			d.Spec.Topology.Historicals = &DruidDataNode{}
		}
		if d.Spec.Topology.Historicals != nil {
			if d.Spec.Topology.Historicals.Replicas == nil {
				d.Spec.Topology.Historicals.Replicas = pointer.Int32P(1)
			}
			if d.Spec.Topology.Historicals.StorageType == "" {
				d.Spec.Topology.Historicals.StorageType = StorageTypeDurable
			}
			if d.Spec.Topology.Historicals.Storage == nil && d.Spec.Topology.Historicals.StorageType == StorageTypeDurable {
				d.Spec.Topology.Historicals.Storage = d.GetDefaultPVC()
			}
			if version.Major() > 25 {
				if d.Spec.Topology.Historicals.PodTemplate.Spec.SecurityContext == nil {
					d.Spec.Topology.Historicals.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{FSGroup: druidVersion.Spec.SecurityContext.RunAsUser}
				}
				d.setDefaultContainerSecurityContext(&druidVersion, &d.Spec.Topology.Historicals.PodTemplate)
				d.setDefaultContainerResourceLimits(&d.Spec.Topology.Historicals.PodTemplate, DruidNodeRoleHistoricals)
			}
		}

		if d.Spec.Topology.Brokers == nil {
			d.Spec.Topology.Brokers = &DruidNode{}
		}
		if d.Spec.Topology.Brokers != nil {
			if d.Spec.Topology.Brokers.Replicas == nil {
				d.Spec.Topology.Brokers.Replicas = pointer.Int32P(1)
			}
			if version.Major() > 25 {
				if d.Spec.Topology.Brokers.PodTemplate.Spec.SecurityContext == nil {
					d.Spec.Topology.Brokers.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{FSGroup: druidVersion.Spec.SecurityContext.RunAsUser}
				}
				d.setDefaultContainerSecurityContext(&druidVersion, &d.Spec.Topology.Brokers.PodTemplate)
				d.setDefaultContainerResourceLimits(&d.Spec.Topology.Brokers.PodTemplate, DruidNodeRoleBrokers)

			}
		}

		if d.Spec.Topology.Routers != nil {
			if d.Spec.Topology.Routers.Replicas == nil {
				d.Spec.Topology.Routers.Replicas = pointer.Int32P(1)
			}
			if version.Major() > 25 {
				if d.Spec.Topology.Routers.PodTemplate.Spec.SecurityContext == nil {
					d.Spec.Topology.Routers.PodTemplate.Spec.SecurityContext = &v1.PodSecurityContext{FSGroup: druidVersion.Spec.SecurityContext.RunAsUser}
				}
				d.setDefaultContainerSecurityContext(&druidVersion, &d.Spec.Topology.Routers.PodTemplate)
				d.setDefaultContainerResourceLimits(&d.Spec.Topology.Routers.PodTemplate, DruidNodeRoleRouters)
			}
		}
	}

	d.SetDefaultsToMetadataStorage()
	d.SetDefaultsToZooKeeperRef()

	if d.Spec.Monitor != nil {
		if d.Spec.Monitor.Prometheus == nil {
			d.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if d.Spec.Monitor.Prometheus != nil && d.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			d.Spec.Monitor.Prometheus.Exporter.Port = kubedb.DruidExporterPort
		}
		d.Spec.Monitor.SetDefaults()
	}
}

func (d *Druid) SetDefaultsToMetadataStorage() {
	if d.Spec.MetadataStorage == nil {
		d.Spec.MetadataStorage = &MetadataStorage{}
	}
	d.SetMetadataStorageObjectRef()
	if d.Spec.MetadataStorage.LinkedDB == "" {
		d.Spec.MetadataStorage.LinkedDB = "druid"
	}
	if d.Spec.MetadataStorage.CreateTables == nil {
		d.Spec.MetadataStorage.CreateTables = ptr.To(true)
	}

	if d.Spec.MetadataStorage.Type == "" {
		if d.Spec.MetadataStorage.ExternallyManaged {
			appBinding, err := d.GetAppBinding(d.Spec.MetadataStorage.Name, d.Spec.MetadataStorage.Namespace)
			if err != nil {
				return
			}
			d.Spec.MetadataStorage.Type = d.GetMetadataStorageType(appBinding.Spec.AppRef.Kind)
		} else {
			d.Spec.MetadataStorage.Type = DruidMetadataStorageMySQL
		}
	}

	if d.Spec.MetadataStorage.Version == nil {
		var defaultVersion string
		if d.Spec.MetadataStorage.Type == DruidMetadataStorageMySQL {
			defaultVersion = "8.0.35"
		} else {
			defaultVersion = "13.13"
		}
		d.Spec.MetadataStorage.Version = &defaultVersion
	}
}

func (d *Druid) SetDefaultsToZooKeeperRef() {
	if d.Spec.ZookeeperRef == nil {
		d.Spec.ZookeeperRef = &ZookeeperRef{}
	}
	d.SetZooKeeperObjectRef()
	if d.Spec.ZookeeperRef.Version == nil {
		defaultVersion := "3.7.2"
		d.Spec.ZookeeperRef.Version = &defaultVersion
	}
}

func (d *Druid) GetDefaultPVC() *core.PersistentVolumeClaimSpec {
	return &core.PersistentVolumeClaimSpec{
		AccessModes: []core.PersistentVolumeAccessMode{
			core.ReadWriteOnce,
		},
		Resources: core.VolumeResourceRequirements{
			Requests: core.ResourceList{
				core.ResourceStorage: resource.MustParse("1Gi"),
			},
		},
	}
}

func (d *Druid) setDefaultContainerSecurityContext(druidVersion *catalog.DruidVersion, podTemplate *ofst.PodTemplateSpec) {
	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.DruidContainerName)
	if container == nil {
		container = &v1.Container{
			Name: kubedb.DruidContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &v1.SecurityContext{}
	}
	d.assignDefaultContainerSecurityContext(druidVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.DruidInitContainerName)
	if initContainer == nil {
		initContainer = &v1.Container{
			Name: kubedb.DruidInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &v1.SecurityContext{}
	}
	d.assignDefaultContainerSecurityContext(druidVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)
}

func (d *Druid) assignDefaultContainerSecurityContext(druidVersion *catalog.DruidVersion, sc *v1.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &v1.Capabilities{
			Drop: []v1.Capability{"ALL"},
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(true)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = druidVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (d *Druid) setDefaultContainerResourceLimits(podTemplate *ofst.PodTemplateSpec, nodeRole DruidNodeRoleType) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.DruidContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		if nodeRole == DruidNodeRoleMiddleManagers {
			apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensiveDruid)
		} else {
			apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
		}
	}

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.DruidInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}
}

func (d *Druid) GetPersistentSecrets() []string {
	if d == nil {
		return nil
	}

	var secrets []string
	if d.Spec.AuthSecret != nil {
		secrets = append(secrets, d.Spec.AuthSecret.Name)
	}
	return secrets
}

func (d *Druid) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	if d.Spec.Topology != nil {
		expectedItems = 4
	}
	if d.Spec.Topology.Routers != nil {
		expectedItems++
	}
	if d.Spec.Topology.Overlords != nil {
		expectedItems++
	}
	return checkReplicasOfPetSet(lister.PetSets(d.Namespace), labels.SelectorFromSet(d.OffshootLabels()), expectedItems)
}

func (d *Druid) GetAppBinding(name string, namespace string) (*appcat.AppBinding, error) {
	appbinding := &appcat.AppBinding{}
	appbinding.Namespace = namespace
	appbinding.Name = name

	if err := DefaultClient.Get(context.TODO(), client.ObjectKeyFromObject(appbinding), appbinding); err != nil {
		klog.Error(err, fmt.Sprintf("failed to get appbinding for metadata storage %s/%s", name, namespace))
		return nil, err
	}
	return appbinding, nil
}

func (d *Druid) SetMetadataStorageObjectRef() {
	if d.Spec.MetadataStorage.ObjectReference == nil {
		d.Spec.MetadataStorage.ObjectReference = &kmapi.ObjectReference{}
	}
	if d.Spec.MetadataStorage.Name == "" {
		d.Spec.MetadataStorage.ExternallyManaged = false
		d.Spec.MetadataStorage.Name = d.GetMetadataStorageName()
	}
	if d.Spec.MetadataStorage.Namespace == "" {
		d.Spec.MetadataStorage.Namespace = d.Namespace
	}
}

func (d *Druid) GetMetadataStorageName() string {
	if d.Spec.MetadataStorage.Type == DruidMetadataStoragePostgreSQL {
		return d.OffShootName() + "-pg-metadata"
	}
	return d.OffShootName() + "-mysql-metadata"
}

func (d *Druid) SetZooKeeperObjectRef() {
	if d.Spec.ZookeeperRef.ObjectReference == nil {
		d.Spec.ZookeeperRef.ObjectReference = &kmapi.ObjectReference{}
	}
	if d.Spec.ZookeeperRef.Name == "" {
		d.Spec.ZookeeperRef.ExternallyManaged = false
		d.Spec.ZookeeperRef.Name = d.GetZooKeeperName()
	}
	if d.Spec.ZookeeperRef.Namespace == "" {
		d.Spec.ZookeeperRef.Namespace = d.Namespace
	}
}

func (d *Druid) GetZooKeeperName() string {
	return d.OffShootName() + "-zk"
}

func (d *Druid) GetInitConfigMapName() string {
	return d.OffShootName() + "-init-script"
}
