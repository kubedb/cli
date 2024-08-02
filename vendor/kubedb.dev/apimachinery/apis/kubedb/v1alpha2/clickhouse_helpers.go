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

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
)

type ClickhouseApp struct {
	*ClickHouse
}

func (r *ClickHouse) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralClickHouse))
}

func (c *ClickHouse) AppBindingMeta() appcat.AppBindingMeta {
	return &ClickhouseApp{c}
}

func (r ClickhouseApp) Name() string {
	return r.ClickHouse.Name
}

func (r ClickhouseApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularClickHouse))
}

// Owner returns owner reference to resources
func (c *ClickHouse) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(c, SchemeGroupVersion.WithKind(c.ResourceKind()))
}

func (c *ClickHouse) ResourceKind() string {
	return ResourceKindClickHouse
}

func (c *ClickHouse) OffshootName() string {
	return c.Name
}

func (c *ClickHouse) OffshootClusterName(value string) string {
	return meta_util.NameWithSuffix(c.OffshootName(), value)
}

func (c *ClickHouse) OffshootClusterPetSetName(clusterName string, shardNo int) string {
	shard := meta_util.NameWithSuffix("shard", strconv.Itoa(shardNo))
	cluster := meta_util.NameWithSuffix(clusterName, shard)
	return meta_util.NameWithSuffix(c.OffshootName(), cluster)
}

func (c *ClickHouse) OffshootLabels() map[string]string {
	return c.offshootLabels(c.OffshootSelectors(), nil)
}

func (c *ClickHouse) OffshootClusterLabels(petSetName string) map[string]string {
	return c.offshootLabels(c.OffshootClusterSelectors(petSetName), nil)
}

func (c *ClickHouse) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, c.Labels, override))
}

func (c *ClickHouse) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      c.ResourceFQN(),
		meta_util.InstanceLabelKey:  c.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (c *ClickHouse) OffshootClusterSelectors(petSetName string, extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      c.ResourceFQN(),
		meta_util.InstanceLabelKey:  c.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
		meta_util.PartOfLabelKey:    petSetName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (c *ClickHouse) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", c.ResourcePlural(), kubedb.GroupName)
}

func (c *ClickHouse) ResourcePlural() string {
	return ResourcePluralClickHouse
}

func (c *ClickHouse) ServiceName() string {
	return c.OffshootName()
}

func (c *ClickHouse) PrimaryServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", c.ServiceName(), c.Namespace)
}

func (c *ClickHouse) GoverningServiceName() string {
	return meta_util.NameWithSuffix(c.ServiceName(), "pods")
}

func (c *ClickHouse) ClusterGoverningServiceName(name string) string {
	return meta_util.NameWithSuffix(name, "pods")
}

func (c *ClickHouse) ClusterGoverningServiceDNS(petSetName string, replicaNo int) string {
	return fmt.Sprintf("%s-%d.%s.%s.svc", petSetName, replicaNo, c.ClusterGoverningServiceName(petSetName), c.GetNamespace())
}

func (c *ClickHouse) GetAuthSecretName() string {
	if c.Spec.AuthSecret != nil && c.Spec.AuthSecret.Name != "" {
		return c.Spec.AuthSecret.Name
	}
	return c.DefaultUserCredSecretName("admin")
}

func (r *ClickHouse) ConfigSecretName() string {
	return meta_util.NameWithSuffix(r.OffshootName(), "config")
}

func (c *ClickHouse) DefaultUserCredSecretName(username string) string {
	return meta_util.NameWithSuffix(c.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", username), "_", "-"))
}

func (c *ClickHouse) PVCName(alias string) string {
	return meta_util.NameWithSuffix(c.Name, alias)
}

func (c *ClickHouse) PetSetName() string {
	return c.OffshootName()
}

func (c *ClickHouse) PodLabels(extraLabels ...map[string]string) map[string]string {
	return c.offshootLabels(meta_util.OverwriteKeys(c.OffshootSelectors(), extraLabels...), c.Spec.PodTemplate.Labels)
}

func (c *ClickHouse) ClusterPodLabels(petSetName string, labels map[string]string, extraLabels ...map[string]string) map[string]string {
	return c.offshootLabels(meta_util.OverwriteKeys(c.OffshootClusterSelectors(petSetName), extraLabels...), labels)
}

func (c *ClickHouse) GetConnectionScheme() string {
	scheme := "http"
	return scheme
}

func (c *ClickHouse) SetHealthCheckerDefaults() {
	if c.Spec.HealthChecker.PeriodSeconds == nil {
		c.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if c.Spec.HealthChecker.TimeoutSeconds == nil {
		c.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if c.Spec.HealthChecker.FailureThreshold == nil {
		c.Spec.HealthChecker.FailureThreshold = pointer.Int32P(3)
	}
}

func (c *ClickHouse) Finalizer() string {
	return fmt.Sprintf("%s/%s", apis.Finalizer, c.ResourceSingular())
}

func (c *ClickHouse) ResourceSingular() string {
	return ResourceSingularClickHouse
}

func (c *ClickHouse) SetDefaults() {
	var chVersion catalog.ClickHouseVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name: c.Spec.Version,
	}, &chVersion)
	if err != nil {
		klog.Errorf("can't get the clickhouse version object %s for %s \n", err.Error(), c.Spec.Version)
		return
	}

	if c.Spec.ClusterTopology != nil {
		clusterName := map[string]bool{}
		clusters := c.Spec.ClusterTopology.Cluster
		for index, cluster := range clusters {
			if cluster.Shards == nil {
				cluster.Shards = pointer.Int32P(1)
			}
			if cluster.Replicas == nil {
				cluster.Replicas = pointer.Int32P(1)
			}
			if cluster.Name == "" {
				for i := 1; ; i += 1 {
					cluster.Name = c.OffshootClusterName(strconv.Itoa(i))
					if !clusterName[cluster.Name] {
						clusterName[cluster.Name] = true
						break
					}
				}
			} else {
				clusterName[cluster.Name] = true
			}
			if cluster.StorageType == "" {
				cluster.StorageType = StorageTypeDurable
			}

			if cluster.PodTemplate == nil {
				cluster.PodTemplate = &ofst.PodTemplateSpec{}
			}

			dbContainer := coreutil.GetContainerByName(cluster.PodTemplate.Spec.Containers, kubedb.ClickHouseContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
			}
			c.setDefaultContainerSecurityContext(&chVersion, cluster.PodTemplate)
			clusters[index] = cluster
		}
		c.Spec.ClusterTopology.Cluster = clusters
	} else {
		if c.Spec.Replicas == nil {
			c.Spec.Replicas = pointer.Int32P(1)
		}
		if c.Spec.DeletionPolicy == "" {
			c.Spec.DeletionPolicy = TerminationPolicyDelete
		}
		if c.Spec.StorageType == "" {
			c.Spec.StorageType = StorageTypeDurable
		}

		if c.Spec.PodTemplate == nil {
			c.Spec.PodTemplate = &ofst.PodTemplateSpec{}
		}
		c.setDefaultContainerSecurityContext(&chVersion, c.Spec.PodTemplate)
		dbContainer := coreutil.GetContainerByName(c.Spec.PodTemplate.Spec.Containers, kubedb.ClickHouseContainerName)
		if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
		}
		c.SetHealthCheckerDefaults()
	}
}

func (c *ClickHouse) setDefaultContainerSecurityContext(chVersion *catalog.ClickHouseVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = chVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.ClickHouseContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.ClickHouseContainerName,
		}
		podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	c.assignDefaultContainerSecurityContext(chVersion, container.SecurityContext)

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.ClickHouseInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.ClickHouseInitContainerName,
		}
		podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	c.assignDefaultContainerSecurityContext(chVersion, initContainer.SecurityContext)
}

func (c *ClickHouse) assignDefaultContainerSecurityContext(chVersion *catalog.ClickHouseVersion, rc *core.SecurityContext) {
	if rc.AllowPrivilegeEscalation == nil {
		rc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if rc.Capabilities == nil {
		rc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if rc.RunAsNonRoot == nil {
		rc.RunAsNonRoot = pointer.BoolP(true)
	}
	if rc.RunAsUser == nil {
		rc.RunAsUser = chVersion.Spec.SecurityContext.RunAsUser
	}
	if rc.SeccompProfile == nil {
		rc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (c *ClickHouse) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 0
	if c.Spec.ClusterTopology != nil {
		for _, cluster := range c.Spec.ClusterTopology.Cluster {
			expectedItems += int(*cluster.Shards)
		}
	} else {
		expectedItems += 1
	}
	return checkReplicasOfPetSet(lister.PetSets(c.Namespace), labels.SelectorFromSet(c.OffshootLabels()), expectedItems)
}
