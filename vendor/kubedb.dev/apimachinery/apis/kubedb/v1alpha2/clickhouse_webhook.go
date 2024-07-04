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
func (r *ClickHouse) Default() {
	if r == nil {
		return
	}
	clickhouselog.Info("default", "name", r.Name)
	r.SetDefaults()
}

var _ webhook.Validator = &ClickHouse{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ClickHouse) ValidateCreate() (admission.Warnings, error) {
	clickhouselog.Info("validate create", "name", r.Name)
	return nil, r.ValidateCreateOrUpdate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ClickHouse) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	clickhouselog.Info("validate update", "name", r.Name)
	return nil, r.ValidateCreateOrUpdate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ClickHouse) ValidateDelete() (admission.Warnings, error) {
	clickhouselog.Info("validate delete", "name", r.Name)

	var allErr field.ErrorList
	if r.Spec.DeletionPolicy == TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("teminationPolicy"),
			r.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, r.Name, allErr)
	}
	return nil, nil
}

func (r *ClickHouse) ValidateCreateOrUpdate() error {
	var allErr field.ErrorList

	if r.Spec.Version == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			r.Name,
			"spec.version' is missing"))
		return apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, r.Name, allErr)
	} else {
		err := r.ValidateVersion(r)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
				r.Spec.Version,
				err.Error()))
			return apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, r.Name, allErr)
		}
	}

	if r.Spec.ClusterTopology != nil {
		clusterName := map[string]bool{}
		clusters := r.Spec.ClusterTopology.Cluster
		for _, cluster := range clusters {
			if cluster.Shards != nil && *cluster.Shards <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("shards"),
					r.Name,
					"number of shards can not be 0 or less"))
			}
			if cluster.Replicas != nil && *cluster.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("replicas"),
					r.Name,
					"number of replicas can't be 0 or less"))
			}
			if clusterName[cluster.Name] {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child(cluster.Name),
					r.Name,
					"cluster name is duplicated, use different cluster name"))
			}
			clusterName[cluster.Name] = true

			allErr = r.validateClusterStorageType(cluster, allErr)

			err := r.validateVolumes(cluster.PodTemplate)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("podTemplate").Child("spec").Child("volumes"),
					r.Name,
					err.Error()))
			}
			err = r.validateVolumesMountPaths(cluster.PodTemplate)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("clusterTopology").Child("podTemplate").Child("spec").Child("volumeMounts"),
					r.Name,
					err.Error()))
			}
		}
		if r.Spec.PodTemplate != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate"),
				r.Name,
				"PodTemplate should be nil in clusterTopology"))
		}

		if r.Spec.Replicas != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replica"),
				r.Name,
				"replica should be nil in clusterTopology"))
		}

		if r.Spec.StorageType != "" {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				r.Name,
				"StorageType should be empty in clusterTopology"))
		}

		if r.Spec.Storage != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storage"),
				r.Name,
				"storage should be nil in clusterTopology"))
		}

	} else {
		// number of replicas can not be 0 or less
		if r.Spec.Replicas != nil && *r.Spec.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				r.Name,
				"number of replicas can't be 0 or less"))
		}

		// number of replicas can not be greater than 1
		if r.Spec.Replicas != nil && *r.Spec.Replicas > 1 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				r.Name,
				"number of replicas can't be greater than 1 in standalone mode"))
		}
		err := r.validateVolumes(r.Spec.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
				r.Name,
				err.Error()))
		}
		err = r.validateVolumesMountPaths(r.Spec.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumeMounts"),
				r.Name,
				err.Error()))
		}

		allErr = r.validateStandaloneStorageType(r.Spec.StorageType, r.Spec.Storage, allErr)
	}

	if len(allErr) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "ClickHouse.kubedb.com", Kind: "ClickHouse"}, r.Name, allErr)
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

func (r *ClickHouse) ValidateVersion(db *ClickHouse) error {
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

func (r *ClickHouse) validateVolumes(podTemplate *ofst.PodTemplateSpec) error {
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

func (r *ClickHouse) validateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec) error {
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
