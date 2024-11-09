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
var cassandralog = logf.Log.WithName("cassandra-resource")

var _ webhook.Defaulter = &Cassandra{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Cassandra) Default() {
	if r == nil {
		return
	}
	cassandralog.Info("default", "name", r.Name)
	r.SetDefaults()
}

var _ webhook.Validator = &Cassandra{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Cassandra) ValidateCreate() (admission.Warnings, error) {
	cassandralog.Info("validate create", "name", r.Name)
	return nil, r.ValidateCreateOrUpdate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Cassandra) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	cassandralog.Info("validate update", "name", r.Name)
	return nil, r.ValidateCreateOrUpdate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Cassandra) ValidateDelete() (admission.Warnings, error) {
	cassandralog.Info("validate delete", "name", r.Name)

	var allErr field.ErrorList
	if r.Spec.DeletionPolicy == TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("deletionPolicy"),
			r.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "Cassandra.kubedb.com", Kind: "Cassandra"}, r.Name, allErr)
	}
	return nil, nil
}

func (r *Cassandra) ValidateCreateOrUpdate() error {
	var allErr field.ErrorList

	if r.Spec.Version == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			r.Name,
			"spec.version' is missing"))
		return apierrors.NewInvalid(schema.GroupKind{Group: "Cassandra.kubedb.com", Kind: "Cassandra"}, r.Name, allErr)
	} else {
		err := r.ValidateVersion(r)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
				r.Spec.Version,
				err.Error()))
			return apierrors.NewInvalid(schema.GroupKind{Group: "cassandra.kubedb.com", Kind: "cassandra"}, r.Name, allErr)
		}
	}

	if r.Spec.Topology != nil {
		rackName := map[string]bool{}
		racks := r.Spec.Topology.Rack
		for _, rack := range racks {
			if rack.Replicas != nil && *rack.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("Topology").Child("replicas"),
					r.Name,
					"number of replicas can't be 0 or less"))
			}
			if rackName[rack.Name] {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("Topology").Child(rack.Name),
					r.Name,
					"rack name is duplicated, use different rack name"))
			}
			rackName[rack.Name] = true

			allErr = r.validateClusterStorageType(rack, allErr)

			err := r.validateVolumes(rack.PodTemplate)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("Topology").Child("podTemplate").Child("spec").Child("volumes"),
					r.Name,
					err.Error()))
			}
			err = r.validateVolumesMountPaths(rack.PodTemplate)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("Topology").Child("podTemplate").Child("spec").Child("volumeMounts"),
					r.Name,
					err.Error()))
			}
		}
		if r.Spec.PodTemplate != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate"),
				r.Name,
				"PodTemplate should be nil in Topology"))
		}

		if r.Spec.Replicas != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replica"),
				r.Name,
				"replica should be nil in Topology"))
		}

		if r.Spec.StorageType != "" {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				r.Name,
				"StorageType should be empty in Topology"))
		}

		if r.Spec.Storage != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storage"),
				r.Name,
				"storage should be nil in Topology"))
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
	return apierrors.NewInvalid(schema.GroupKind{Group: "Cassandra.kubedb.com", Kind: "Cassandra"}, r.Name, allErr)
}

func (c *Cassandra) validateStandaloneStorageType(storageType StorageType, storage *core.PersistentVolumeClaimSpec, allErr field.ErrorList) field.ErrorList {
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

func (c *Cassandra) validateClusterStorageType(rack RackSpec, allErr field.ErrorList) field.ErrorList {
	if rack.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("Topology").Child(rack.Name).Child("storageType"),
			c.Name,
			"StorageType can not be empty"))
	} else {
		if rack.StorageType != StorageTypeDurable && rack.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("Topology").Child(rack.Name).Child("storageType"),
				rack.StorageType,
				"StorageType should be either durable or ephemeral"))
		}
	}
	if rack.Storage == nil && rack.StorageType == StorageTypeDurable {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("Topology").Child(rack.Name).Child("storage"),
			c.Name,
			"Storage can't be empty when StorageType is durable"))
	}
	return allErr
}

func (r *Cassandra) ValidateVersion(db *Cassandra) error {
	casVersion := catalog.CassandraVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: db.Spec.Version}, &casVersion)
	if err != nil {
		return errors.New(fmt.Sprint("version ", db.Spec.Version, " not supported"))
	}
	return nil
}

var cassandraReservedVolumes = []string{
	kubedb.CassandraVolumeData,
}

func (r *Cassandra) validateVolumes(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate.Spec.Volumes == nil {
		return nil
	}
	rsv := make([]string, len(cassandraReservedVolumes))
	copy(rsv, cassandraReservedVolumes)
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

var cassandraReservedVolumeMountPaths = []string{
	kubedb.CassandraDataDir,
}

func (r *Cassandra) validateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Containers == nil {
		return nil
	}

	for _, rvmp := range cassandraReservedVolumeMountPaths {
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

	for _, rvmp := range cassandraReservedVolumeMountPaths {
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
