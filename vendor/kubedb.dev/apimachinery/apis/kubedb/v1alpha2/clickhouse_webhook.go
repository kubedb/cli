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
	"errors"
	"fmt"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"

	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var clickhouselog = logf.Log.WithName("clickhouse-resource")

var _ webhook.Defaulter = &ClickHouse{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (c *ClickHouse) Default() {
	if c == nil {
		return
	}
	clickhouselog.Info("default", "name", c.Name)
	c.SetDefaults()
}

var _ webhook.Validator = &ClickHouse{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (c *ClickHouse) ValidateCreate() (admission.Warnings, error) {
	clickhouselog.Info("validate create", "name", c.Name)
	return nil, c.ValidateCreateOrUpdate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (c *ClickHouse) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	clickhouselog.Info("validate update", "name", c.Name)
	return nil, c.ValidateCreateOrUpdate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (c *ClickHouse) ValidateDelete() (admission.Warnings, error) {
	clickhouselog.Info("validate delete", "name", c.Name)

	var allErr field.ErrorList
	if c.Spec.DeletionPolicy == DeletionPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("teminationPolicy"),
			c.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, c.Name, allErr)
	}
	return nil, nil
}

func (c *ClickHouse) ValidateCreateOrUpdate() error {
	var allErr field.ErrorList

	if c.Spec.Version == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			c.Name,
			"spec.version' is missing"))
		return apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, c.Name, allErr)
	} else {
		err := c.ValidateVersion(c)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
				c.Spec.Version,
				err.Error()))
			return apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, c.Name, allErr)
		}
	}

	if c.Spec.DisableSecurity {
		if c.Spec.AuthSecret != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("authSecret"),
				c.Name,
				"authSecret should be nil when security is disabled"))
			return apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, c.Name, allErr)
		}
	}

	if c.Spec.ClusterTopology != nil {
		clusterName := map[string]bool{}
		clusters := c.Spec.ClusterTopology.Cluster
		if c.Spec.ClusterTopology.ClickHouseKeeper != nil {
			if !c.Spec.ClusterTopology.ClickHouseKeeper.ExternallyManaged {
				if c.Spec.ClusterTopology.ClickHouseKeeper.Spec == nil {
					allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("spec"),
						c.Name,
						"spec can't be nil when externally managed is false"))
				} else {
					if *c.Spec.ClusterTopology.ClickHouseKeeper.Spec.Replicas < 1 {
						allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("spec").Child("replica"),
							c.Name,
							"number of replica can not be 0 or less"))
					}
					allErr = c.validateClickHouseKeeperStorageType(c.Spec.ClusterTopology.ClickHouseKeeper.Spec.StorageType, c.Spec.ClusterTopology.ClickHouseKeeper.Spec.Storage, allErr)
				}
				if c.Spec.ClusterTopology.ClickHouseKeeper.Node != nil {
					allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("node"),
						c.Name,
						"ClickHouse Keeper node should be empty when externally managed is false"))
				}
			} else {
				if c.Spec.ClusterTopology.ClickHouseKeeper.Node == nil {
					allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("node"),
						c.Name,
						"ClickHouse Keeper node can't be empty when externally managed is true"))
				} else {
					if c.Spec.ClusterTopology.ClickHouseKeeper.Node.Host == "" {
						allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("node").Child("host"),
							c.Name,
							"ClickHouse Keeper host can't be empty"))
					}
					if c.Spec.ClusterTopology.ClickHouseKeeper.Node.Port == nil {
						allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("node").Child("port"),
							c.Name,
							"ClickHouse Keeper port can't be empty"))
					}
				}
				if c.Spec.ClusterTopology.ClickHouseKeeper.Spec != nil {
					allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("spec"),
						c.Name,
						"ClickHouse Keeper spec should be empty when externally managed is true"))
				}
			}
		}
		for _, cluster := range clusters {
			if cluster.Shards != nil && *cluster.Shards <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("shards"),
					c.Name,
					"number of shards can not be 0 or less"))
			}
			if cluster.Replicas != nil && *cluster.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("replicas"),
					c.Name,
					"number of replicas can't be 0 or less"))
			}
			if clusterName[cluster.Name] {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child(cluster.Name),
					c.Name,
					"cluster name is already exists, use different cluster name"))
			}
			clusterName[cluster.Name] = true

			allErr = c.validateClusterStorageType(cluster, allErr)

			err := c.validateVolumes(cluster.PodTemplate)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("podTemplate").Child("spec").Child("volumes"),
					c.Name,
					err.Error()))
			}
			err = c.validateVolumesMountPaths(cluster.PodTemplate)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("podTemplate").Child("spec").Child("volumeMounts"),
					c.Name,
					err.Error()))
			}
		}
		if c.Spec.PodTemplate != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate"),
				c.Name,
				"PodTemplate should be nil in clusterTopology"))
		}

		if c.Spec.Replicas != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replica"),
				c.Name,
				"replica should be nil in clusterTopology"))
		}

		if c.Spec.StorageType != "" {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				c.Name,
				"StorageType should be empty in clusterTopology"))
		}

		if c.Spec.Storage != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storage"),
				c.Name,
				"storage should be nil in clusterTopology"))
		}

	} else {
		// number of replicas can not be 0 or less
		if c.Spec.Replicas != nil && *c.Spec.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				c.Name,
				"number of replicas can't be 0 or less"))
		}

		// number of replicas can not be greater than 1
		if c.Spec.Replicas != nil && *c.Spec.Replicas > 1 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				c.Name,
				"number of replicas can't be greater than 1 in standalone mode"))
		}
		err := c.validateVolumes(c.Spec.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
				c.Name,
				err.Error()))
		}
		err = c.validateVolumesMountPaths(c.Spec.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumeMounts"),
				c.Name,
				err.Error()))
		}

		allErr = c.validateStandaloneStorageType(c.Spec.StorageType, c.Spec.Storage, allErr)
	}

	if len(allErr) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, c.Name, allErr)
}

func (c *ClickHouse) validateStandaloneStorageType(storageType StorageType, storage *core.PersistentVolumeClaimSpec, allErr field.ErrorList) field.ErrorList {
	if storageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			c.Name,
			"StorageType can not be empty"))
	} else {
		if storageType != StorageTypeDurable && c.Spec.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				c.Name,
				"StorageType should be either durable or ephemeral"))
		}
	}

	if storage == nil && c.Spec.StorageType == StorageTypeDurable {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storage"),
			c.Name,
			"Storage can't be empty when StorageType is durable"))
	}

	return allErr
}

func (c *ClickHouse) validateClusterStorageType(cluster ClusterSpec, allErr field.ErrorList) field.ErrorList {
	if cluster.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child(cluster.Name).Child("storageType"),
			c.Name,
			"StorageType can not be empty"))
	} else {
		if cluster.StorageType != StorageTypeDurable && cluster.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child(cluster.Name).Child("storageType"),
				cluster.StorageType,
				"StorageType should be either durable or ephemeral"))
		}
	}
	if cluster.Storage == nil && cluster.StorageType == StorageTypeDurable {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child(cluster.Name).Child("storage"),
			c.Name,
			"Storage can't be empty when StorageType is durable"))
	}
	return allErr
}

func (c *ClickHouse) validateClickHouseKeeperStorageType(storageType StorageType, storage *core.PersistentVolumeClaimSpec, allErr field.ErrorList) field.ErrorList {
	if storageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("spec").Child("storageType"),
			c.Name,
			"StorageType can not be empty"))
	} else {
		if storageType != StorageTypeDurable && c.Spec.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("spec").Child("storageType"),
				c.Name,
				"StorageType should be either durable or ephemeral"))
		}
	}
	if storage == nil && c.Spec.StorageType == StorageTypeDurable {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("clickHouseKeeper").Child("spec").Child("storage"),
			c.Name,
			"Storage can't be empty when StorageType is durable"))
	}

	return allErr
}

func (c *ClickHouse) ValidateVersion(db *ClickHouse) error {
	chVersion := catalog.ClickHouseVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: db.Spec.Version}, &chVersion)
	if err != nil {
		// fmt.Sprint(db.Spec.Version, "version not supported")
		return errors.New(fmt.Sprint("version ", db.Spec.Version, " not supported"))
	}
	return nil
}

var clickhouseReservedVolumes = []string{
	kubedb.ClickHouseVolumeData,
}

func (c *ClickHouse) validateVolumes(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate.Spec.Volumes == nil {
		return nil
	}
	rsv := make([]string, len(clickhouseReservedVolumes))
	copy(rsv, clickhouseReservedVolumes)
	volumes := podTemplate.Spec.Volumes
	for _, rv := range rsv {
		for _, ugv := range volumes {
			if ugv.Name == rv {
				return errors.New("Cannot use a reserve volume name: " + rv)
			}
		}
	}
	return nil
}

var clickhouseReservedVolumeMountPaths = []string{
	kubedb.ClickHouseDataDir,
}

func (c *ClickHouse) validateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Containers == nil {
		return nil
	}

	for _, rvmp := range clickhouseReservedVolumeMountPaths {
		containerList := podTemplate.Spec.Containers
		for i := range containerList {
			mountPathList := containerList[i].VolumeMounts
			for j := range mountPathList {
				if mountPathList[j].MountPath == rvmp {
					return errors.New("Can't use a reserve volume mount path name: " + rvmp)
				}
			}
		}
	}

	if podTemplate.Spec.InitContainers == nil {
		return nil
	}

	for _, rvmp := range clickhouseReservedVolumeMountPaths {
		containerList := podTemplate.Spec.InitContainers
		for i := range containerList {
			mountPathList := containerList[i].VolumeMounts
			for j := range mountPathList {
				if mountPathList[j].MountPath == rvmp {
					return errors.New("Can't use a reserve volume mount path name: " + rvmp)
				}
			}
		}
	}
	return nil
}
