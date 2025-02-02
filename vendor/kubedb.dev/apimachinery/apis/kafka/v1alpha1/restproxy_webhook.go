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
	"fmt"

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
var restproxylog = logf.Log.WithName("restproxy-resource")

var _ webhook.Defaulter = &RestProxy{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (k *RestProxy) Default() {
	if k == nil {
		return
	}
	restproxylog.Info("default", "name", k.Name)
	k.SetDefaults()
}

var _ webhook.Validator = &RestProxy{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (k *RestProxy) ValidateCreate() (admission.Warnings, error) {
	restproxylog.Info("validate create", "name", k.Name)
	allErr := k.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "RestProxy"}, k.Name, allErr)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (k *RestProxy) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	restproxylog.Info("validate update", "name", k.Name)
	allErr := k.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "RestProxy"}, k.Name, allErr)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (k *RestProxy) ValidateDelete() (admission.Warnings, error) {
	restproxylog.Info("validate delete", "name", k.Name)

	var allErr field.ErrorList
	if k.Spec.DeletionPolicy == dbapi.DeletionPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("deletionPolicy"),
			k.Name,
			"Can not delete as deletionPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "RestProxy"}, k.Name, allErr)
	}
	return nil, nil
}

func (k *RestProxy) ValidateCreateOrUpdate() field.ErrorList {
	var allErr field.ErrorList

	err := k.validateVersion()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			k.Name,
			err.Error()))
		return allErr
	}

	if k.Spec.SchemaRegistryRef != nil {
		if k.Spec.SchemaRegistryRef.InternallyManaged && k.Spec.SchemaRegistryRef.ObjectReference != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("schemaRegistryRef").Child("objectReference"),
				k.Name,
				"ObjectReference should be nil when InternallyManaged is true"))
		}
		if !k.Spec.SchemaRegistryRef.InternallyManaged && k.Spec.SchemaRegistryRef.ObjectReference == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("schemaRegistryRef").Child("objectReference"),
				k.Name,
				"ObjectReference should not be nil when InternallyManaged is false"))
		}
	}

	if k.Spec.DeletionPolicy == dbapi.DeletionPolicyHalt {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("deletionPolicy"),
			k.Name,
			"DeletionPolicyHalt is not supported for RestProxy"))
	}

	// number of replicas can not be 0 or less
	if k.Spec.Replicas != nil && *k.Spec.Replicas <= 0 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			k.Name,
			"number of replicas can not be 0 or less"))
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

	if len(allErr) == 0 {
		return nil
	}
	return allErr
}

func (k *RestProxy) validateVersion() error {
	ksrVersion := &catalog.SchemaRegistryVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: k.Spec.Version}, ksrVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	if ksrVersion.Spec.Distribution != catalog.SchemaRegistryDistroAiven {
		return errors.New(fmt.Sprintf("Distribution %s is not supported, only supported distribution is Aiven", ksrVersion.Spec.Distribution))
	}
	return nil
}

var restProxyReservedVolumes = []string{
	KafkaClientCertVolumeName,
	RestProxyOperatorVolumeConfig,
}

func (k *RestProxy) validateVolumes() error {
	if k.Spec.PodTemplate.Spec.Volumes == nil {
		return nil
	}
	rsv := make([]string, len(restProxyReservedVolumes))
	copy(rsv, restProxyReservedVolumes)
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

var restProxyReservedVolumeMountPaths = []string{
	KafkaClientCertDir,
	RestProxyOperatorVolumeConfig,
}

func (k *RestProxy) validateContainerVolumeMountPaths() error {
	container := coreutil.GetContainerByName(k.Spec.PodTemplate.Spec.Containers, RestProxyContainerName)
	if container == nil {
		return errors.New("container not found")
	}
	rPaths := restProxyReservedVolumeMountPaths
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
