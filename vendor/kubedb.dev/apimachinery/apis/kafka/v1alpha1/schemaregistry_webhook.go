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

package v1alpha1

import (
	"context"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	coreutil "kmodules.xyz/client-go/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var schemaregistrylog = logf.Log.WithName("schemaregistry-resource")

var _ webhook.Defaulter = &SchemaRegistry{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (k *SchemaRegistry) Default() {
	if k == nil {
		return
	}
	schemaregistrylog.Info("default", "name", k.Name)
	k.SetDefaults()
}

var _ webhook.Validator = &SchemaRegistry{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (k *SchemaRegistry) ValidateCreate() (admission.Warnings, error) {
	schemaregistrylog.Info("validate create", "name", k.Name)
	allErr := k.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "SchemaRegistry"}, k.Name, allErr)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (k *SchemaRegistry) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	schemaregistrylog.Info("validate update", "name", k.Name)

	oldRegistry := old.(*SchemaRegistry)
	allErr := k.ValidateCreateOrUpdate()

	if *oldRegistry.Spec.Replicas == 1 && *k.Spec.Replicas > 1 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			k.Name,
			"Cannot scale up from 1 to more than 1 in standalone mode"))
	}

	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "SchemaRegistry"}, k.Name, allErr)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (k *SchemaRegistry) ValidateDelete() (admission.Warnings, error) {
	schemaregistrylog.Info("validate delete", "name", k.Name)

	var allErr field.ErrorList
	if k.Spec.DeletionPolicy == dbapi.DeletionPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			k.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "SchemaRegistry"}, k.Name, allErr)
	}
	return nil, nil
}

func (k *SchemaRegistry) ValidateCreateOrUpdate() field.ErrorList {
	var allErr field.ErrorList

	if k.Spec.DeletionPolicy == dbapi.DeletionPolicyHalt {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			k.Name,
			"DeletionPolicyHalt is not supported for SchemaRegistry"))
	}

	// number of replicas can not be 0 or less
	if k.Spec.Replicas != nil && *k.Spec.Replicas <= 0 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			k.Name,
			"number of replicas can not be 0 or less"))
	}

	err := k.validateVersion()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			k.Name,
			err.Error()))
	}

	err = k.validateVolumes()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
			k.Name,
			err.Error()))
	}

	err = k.validateContainerVolumeMountPaths()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("containers").Child("volumeMounts"),
			k.Name,
			err.Error()))
	}

	err = k.validateEnvVars()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("containers").Child("envs"),
			k.Name,
			err.Error()))
	}

	if len(allErr) == 0 {
		return nil
	}
	return allErr
}

func (k *SchemaRegistry) validateEnvVars() error {
	return nil
}

func (k *SchemaRegistry) validateVersion() error {
	ksrVersion := &catalog.SchemaRegistryVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: k.Spec.Version}, ksrVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	return nil
}

var schemaRegistryReservedVolumes = []string{
	KafkaClientCertVolumeName,
}

func (k *SchemaRegistry) validateVolumes() error {
	if k.Spec.PodTemplate.Spec.Volumes == nil {
		return nil
	}
	rsv := make([]string, len(schemaRegistryReservedVolumes))
	copy(rsv, schemaRegistryReservedVolumes)
	volumes := k.Spec.PodTemplate.Spec.Volumes
	for _, rv := range rsv {
		for _, ugv := range volumes {
			if ugv.Name == rv {
				return errors.New("Cannot use a reserve volume name: " + rv)
			}
		}
	}
	return nil
}

var schemaRegistryReservedVolumeMountPaths = []string{
	KafkaClientCertDir,
}

func (k *SchemaRegistry) validateContainerVolumeMountPaths() error {
	container := coreutil.GetContainerByName(k.Spec.PodTemplate.Spec.Containers, SchemaRegistryContainerName)
	if container == nil {
		return errors.New("container not found")
	}
	rPaths := schemaRegistryReservedVolumeMountPaths
	volumeMountPaths := container.VolumeMounts
	for _, rvm := range rPaths {
		for _, ugv := range volumeMountPaths {
			if ugv.MountPath == rvm {
				return errors.New("Cannot use a reserve volume mount path: " + rvm)
			}
		}
	}
	return nil
}
