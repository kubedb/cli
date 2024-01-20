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

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"

	"github.com/pkg/errors"
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
var singlestorelog = logf.Log.WithName("singlestore-resource")

var _ webhook.Defaulter = &Singlestore{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (s *Singlestore) Default() {
	if s == nil {
		return
	}
	singlestorelog.Info("default", "name", s.Name)

	s.SetDefaults()
}

var _ webhook.Validator = &Singlestore{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (s *Singlestore) ValidateCreate() (admission.Warnings, error) {
	singlestorelog.Info("validate create", "name", s.Name)
	allErr := s.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Singlestore"}, s.Name, allErr)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (s *Singlestore) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	singlestorelog.Info("validate update", "name", s.Name)

	oldConnect := old.(*Singlestore)
	allErr := s.ValidateCreateOrUpdate()

	if s.Spec.Topology == nil && *oldConnect.Spec.Replicas == 1 && *s.Spec.Replicas > 1 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			s.Name,
			"Cannot scale up from 1 to more than 1 in standalone mode"))
	}

	if len(allErr) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Singlestore"}, s.Name, allErr)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (s *Singlestore) ValidateDelete() (admission.Warnings, error) {
	singlestorelog.Info("validate delete", "name", s.Name)

	var allErr field.ErrorList
	if s.Spec.TerminationPolicy == TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			s.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Singlestore"}, s.Name, allErr)
	}
	return nil, nil
}

func (s *Singlestore) ValidateCreateOrUpdate() field.ErrorList {
	var allErr field.ErrorList

	if s.Spec.EnableSSL {
		if s.Spec.TLS == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("enableSSL"),
				s.Name,
				".spec.tls can't be nil, if .spec.enableSSL is true"))
		}
	} else {
		if s.Spec.TLS != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("enableSSL"),
				s.Name,
				".spec.tls must be nil, if .spec.enableSSL is disabled"))
		}
	}

	if s.Spec.Version == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			s.Name,
			"spec.version' is missing"))
	} else {
		err := sdbValidateVersion(s)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
				s.Name,
				err.Error()))
		}
	}

	if s.Spec.Topology == nil {
		if *s.Spec.Replicas != 1 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				s.Name,
				"number of replicas for standalone must be one "))
		}
		err := sdbValidateVolumes(s.Spec.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
				s.Name,
				err.Error()))
		}
		err = sdbValidateVolumesMountPaths(s.Spec.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("containers"),
				s.Name,
				err.Error()))
		}

	} else {
		if s.Spec.Topology.Aggregator == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("aggregator"),
				s.Name,
				".spec.topology.aggregator can't be empty in cluster mode"))
		}
		if s.Spec.Topology.Leaf == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("leaf"),
				s.Name,
				".spec.topology.leaf can't be empty in cluster mode"))
		}

		if s.Spec.Topology.Aggregator.Replicas == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("aggregator").Child("replicas"),
				s.Name,
				"doesn't support spec.topology.aggregator.replicas is set"))
		}
		if s.Spec.Topology.Leaf.Replicas == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("leaf").Child("replicas"),
				s.Name,
				"doesn't support spec.topology.leaf.replicas is set"))
		}

		if *s.Spec.Topology.Aggregator.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("aggregator").Child("replicas"),
				s.Name,
				"number of replicas can not be less be 0 or less"))
		}

		if *s.Spec.Topology.Leaf.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("leaf").Child("replicas"),
				s.Name,
				"number of replicas can not be 0 or less"))
		}

		err := sdbValidateVolumes(s.Spec.Topology.Aggregator.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("aggregator").Child("podTemplate").Child("spec").Child("volumes"),
				s.Name,
				err.Error()))
		}
		err = sdbValidateVolumes(s.Spec.Topology.Leaf.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("leaf").Child("podTemplate").Child("spec").Child("volumes"),
				s.Name,
				err.Error()))
		}

		err = sdbValidateVolumesMountPaths(s.Spec.Topology.Aggregator.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("aggregator").Child("podTemplate").Child("spec").Child("containers"),
				s.Name,
				err.Error()))
		}
		err = sdbValidateVolumesMountPaths(s.Spec.Topology.Leaf.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("leaf").Child("podTemplate").Child("spec").Child("containers"),
				s.Name,
				err.Error()))
		}
	}

	if s.Spec.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			s.Name,
			"StorageType can not be empty"))
	} else {
		if s.Spec.StorageType != StorageTypeDurable && s.Spec.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				s.Name,
				"StorageType should be either durable or ephemeral"))
		}
	}

	if len(allErr) == 0 {
		return nil
	}
	return allErr
}

// reserved volume and volumes mounts for singlestore
var sdbReservedVolumes = []string{
	SinglestoreVolumeNameUserInitScript,
	SinglestoreVolumeNameCustomConfig,
	SinglestoreVolmeNameInitScript,
	SinglestoreVolumeNameData,
}

var sdbReservedVolumesMountPaths = []string{
	SinglestoreVolumeMountPathData,
	SinglestoreVolumeMountPathInitScript,
	SinglestoreVolumeMountPathCustomConfig,
	SinglestoreVolumeMountPathUserInitScript,
}

func sdbValidateVersion(s *Singlestore) error {
	var sdbVersion catalog.SinglestoreVersion
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name: s.Spec.Version,
	}, &sdbVersion)
	if err != nil {
		return errors.New("version not supported")
	}

	return nil
}

func sdbValidateVolumes(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Volumes == nil {
		return nil
	}

	for _, rv := range sdbReservedVolumes {
		for _, ugv := range podTemplate.Spec.Volumes {
			if ugv.Name == rv {
				return errors.New("Can't use a reserve volume name: " + rv)
			}
		}
	}

	return nil
}

func sdbValidateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Containers == nil {
		return nil
	}

	for _, rvmp := range sdbReservedVolumesMountPaths {
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

	for _, rvmp := range sdbReservedVolumesMountPaths {
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
