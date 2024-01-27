/*
Copyright 2023.

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

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"

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
var druidlog = logf.Log.WithName("druid-resource")

//+kubebuilder:webhook:path=/mutate-kubedb-com-v1alpha2-druid,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=druids,verbs=create;update,versions=v1alpha2,name=mdruid.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Druid{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (d *Druid) Default() {
	if d == nil {
		return
	}
	druidlog.Info("default", "name", d.Name)

	d.SetDefaults()
}

//+kubebuilder:webhook:path=/validate-kubedb-com-v1alpha2-druid,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=druids,verbs=create;update,versions=v1alpha2,name=vdruid.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Druid{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (d *Druid) ValidateCreate() (admission.Warnings, error) {
	druidlog.Info("validate create", "name", d.Name)

	allErr := d.validateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Druid"}, d.Name, allErr)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (d *Druid) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	druidlog.Info("validate update", "name", d.Name)
	_ = old.(*Druid)
	allErr := d.validateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Druid"}, d.Name, allErr)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (d *Druid) ValidateDelete() (admission.Warnings, error) {
	druidlog.Info("validate delete", "name", d.Name)
	return nil, nil
}

var druidReservedVolumes = []string{
	DruidVolumeOperatorConfig,
	DruidVolumeMainConfig,
	DruidVolumeCustomConfig,
	DruidVolumeMySQLMetadataStorage,
}

var druidReservedVolumeMountPaths = []string{
	DruidCConfigDirMySQLMetadata,
	DruidOperatorConfigDir,
	DruidMainConfigDir,
	DruidCustomConfigDir,
}

func (d *Druid) validateCreateOrUpdate() field.ErrorList {
	var allErr field.ErrorList

	if d.Spec.Version == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			d.Name,
			"spec.version is missing"))
	} else {
		err := druidValidateVersion(d)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
				d.Name,
				err.Error()))
		}
	}

	if d.Spec.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			d.Name,
			"StorageType can not be empty"))
	} else {
		if d.Spec.StorageType != StorageTypeDurable && d.Spec.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				d.Name,
				"StorageType should be either durable or ephemeral"))
		}
	}

	if d.Spec.DeepStorage == nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("deepStorage"),
			d.Name,
			"spec.deepStorage is missing"))
	} else {
		if d.Spec.DeepStorage.Type == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("deepStorage").Child("type"),
				d.Name,
				"spec.deepStorage.type is missing"))
		}
	}

	if d.Spec.MetadataStorage == nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("metadataStorage"),
			d.Name,
			"spec.metadataStorage is missing"))
	} else {
		if d.Spec.MetadataStorage.Name == nil && d.Spec.MetadataStorage.Type == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("metadataStorage").Child("name"),
				d.Name,
				"spec.metadataStorage.type and spec.metadataStorage.name both can not be empty simultaneously"))
		}
	}

	if d.Spec.Topology == nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology"),
			d.Name,
			"spec.topology can not be empty"))
	} else {
		if d.Spec.Topology.Coordinators == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("coordinators"),
				d.Name,
				"spec.topology.coordinators can not be empty"))
		} else {

			if *d.Spec.Topology.Coordinators.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("coordinators").Child("replicas"),
					d.Name,
					"number of replicas can not be 0 or less"))
			}

			err := druidValidateVolumes(&d.Spec.Topology.Coordinators.PodTemplate, DruidNodeRoleCoordinators)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("coordinators").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
			err = druidValidateVolumesMountPaths(&d.Spec.Topology.Coordinators.PodTemplate, DruidNodeRoleCoordinators)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("coordinators").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
		}

		if d.Spec.Topology.MiddleManagers == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("middleManagers"),
				d.Name,
				"spec.topology.middleManagers can not be empty"))
		} else {
			if *d.Spec.Topology.MiddleManagers.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("middleManagers").Child("replicas"),
					d.Name,
					"number of replicas can not be 0 or less"))
			}

			err := druidValidateVolumes(&d.Spec.Topology.MiddleManagers.PodTemplate, DruidNodeRoleMiddleManagers)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("middleManagers").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
			err = druidValidateVolumesMountPaths(&d.Spec.Topology.MiddleManagers.PodTemplate, DruidNodeRoleMiddleManagers)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("middleManagers").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
		}

		if d.Spec.Topology.Historicals == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("historicals"),
				d.Name,
				"spec.topology.historicals can not be empty"))
		} else {
			if *d.Spec.Topology.Historicals.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("historicals").Child("replicas"),
					d.Name,
					"number of replicas can not be 0 or less"))
			}

			err := druidValidateVolumes(&d.Spec.Topology.Historicals.PodTemplate, DruidNodeRoleHistoricals)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("historicals").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
			err = druidValidateVolumesMountPaths(&d.Spec.Topology.Historicals.PodTemplate, DruidNodeRoleHistoricals)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("historicals").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
		}

		if d.Spec.Topology.Brokers == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("brokers").Child("replicas"),
				d.Name,
				"spec.topology.brokers.replicas can not be empty"))
		} else {
			if *d.Spec.Topology.Brokers.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("brokers").Child("replicas"),
					d.Name,
					"number of replicas can not be 0 or less"))
			}

			err := druidValidateVolumes(&d.Spec.Topology.Brokers.PodTemplate, DruidNodeRoleBrokers)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("brokers").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
			err = druidValidateVolumesMountPaths(&d.Spec.Topology.Brokers.PodTemplate, DruidNodeRoleBrokers)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("brokers").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
		}

		if d.Spec.Topology.Overlords != nil {
			if *d.Spec.Topology.Overlords.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("overlords").Child("replicas"),
					d.Name,
					"number of replicas can not be 0 or less"))
			}

			err := druidValidateVolumes(&d.Spec.Topology.Overlords.PodTemplate, DruidNodeRoleOverlords)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("overlords").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
			err = druidValidateVolumesMountPaths(&d.Spec.Topology.Overlords.PodTemplate, DruidNodeRoleOverlords)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("overlords").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
		}
		if d.Spec.Topology.Routers != nil {
			if *d.Spec.Topology.Routers.Replicas <= 0 {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("routers").Child("replicas"),
					d.Name,
					"number of replicas can not be 0 or less"))
			}

			err := druidValidateVolumes(&d.Spec.Topology.Routers.PodTemplate, DruidNodeRoleRouters)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("routers").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
			err = druidValidateVolumesMountPaths(&d.Spec.Topology.Routers.PodTemplate, DruidNodeRoleRouters)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("routers").Child("podTemplate").Child("spec").Child("volumes"),
					d.Name,
					err.Error()))
			}
		}
	}
	if len(allErr) == 0 {
		return nil
	}
	return allErr
}

func druidValidateVersion(d *Druid) error {
	var druidVersion catalog.DruidVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name: d.Spec.Version,
	}, &druidVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	return nil
}

func druidValidateVolumes(podTemplate *ofst.PodTemplateSpec, nodeType DruidNodeRoleType) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Volumes == nil {
		return nil
	}

	if nodeType == DruidNodeRoleHistoricals {
		druidReservedVolumes = append(druidReservedVolumes, DruidVolumeHistoricalsSegmentCache)
	} else if nodeType == DruidNodeRoleMiddleManagers {
		druidReservedVolumes = append(druidReservedVolumes, DruidVolumeMiddleManagersBaseTaskDir)
	}

	for _, rv := range druidReservedVolumes {
		for _, ugv := range podTemplate.Spec.Volumes {
			if ugv.Name == rv {
				return errors.New("Can't use a reserve volume name: " + rv)
			}
		}
	}

	return nil
}

func druidValidateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec, nodeType DruidNodeRoleType) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Containers == nil {
		return nil
	}

	if nodeType == DruidNodeRoleHistoricals {
		druidReservedVolumeMountPaths = append(druidReservedVolumeMountPaths, DruidHistoricalsSegmentCacheDir)
	}
	if nodeType == DruidNodeRoleMiddleManagers {
		druidReservedVolumeMountPaths = append(druidReservedVolumeMountPaths, DruidWorkerTaskBaseTaskDir)
	}

	for _, rvmp := range druidReservedVolumeMountPaths {
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

	for _, rvmp := range druidReservedVolumeMountPaths {
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
