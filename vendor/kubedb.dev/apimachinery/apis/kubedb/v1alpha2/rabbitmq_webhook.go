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

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"

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
var rabbitmqlog = logf.Log.WithName("rabbitmq-resource")

//+kubebuilder:webhook:path=/mutate-rabbitmq-kubedb-com-v1alpha1-rabbitmq,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=rabbitmqs,verbs=create;update,versions=v1alpha1,name=mrabbitmq.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &RabbitMQ{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *RabbitMQ) Default() {
	if r == nil {
		return
	}
	rabbitmqlog.Info("default", "name", r.Name)
	r.SetDefaults()
}

//+kubebuilder:webhook:path=/validate-rabbitmq-kubedb-com-v1alpha1-rabbitmq,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=rabbitmqs,verbs=create;update,versions=v1alpha1,name=vrabbitmq.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &RabbitMQ{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *RabbitMQ) ValidateCreate() (admission.Warnings, error) {
	rabbitmqlog.Info("validate create", "name", r.Name)
	return nil, r.ValidateCreateOrUpdate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *RabbitMQ) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	rabbitmqlog.Info("validate update", "name", r.Name)
	return nil, r.ValidateCreateOrUpdate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *RabbitMQ) ValidateDelete() (admission.Warnings, error) {
	rabbitmqlog.Info("validate delete", "name", r.Name)

	var allErr field.ErrorList
	if r.Spec.DeletionPolicy == TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("teminationPolicy"),
			r.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "rabbitmq.kubedb.com", Kind: "RabbitMQ"}, r.Name, allErr)
	}
	return nil, nil
}

func (r *RabbitMQ) ValidateCreateOrUpdate() error {
	var allErr field.ErrorList
	if r.Spec.EnableSSL {
		if r.Spec.TLS == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("enableSSL"),
				r.Name,
				".spec.tls can't be nil, if .spec.enableSSL is true"))
		}
	} else {
		if r.Spec.TLS != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("enableSSL"),
				r.Name,
				".spec.tls must be nil, if .spec.enableSSL is disabled"))
		}
	}

	// number of replicas can not be 0 or less
	if r.Spec.Replicas != nil && *r.Spec.Replicas <= 0 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			r.Name,
			"number of replicas can not be 0 or less"))
	}

	if r.Spec.Version == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			r.Name,
			"spec.version' is missing"))
	} else {
		err := r.ValidateVersion(r)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
				r.Name,
				err.Error()))
		}
	}

	err := r.validateVolumes(r)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
			r.Name,
			err.Error()))
	}

	err = r.validateVolumesMountPaths(&r.Spec.PodTemplate)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumeMounts"),
			r.Name,
			err.Error()))
	}

	if r.Spec.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			r.Name,
			"StorageType can not be empty"))
	} else {
		if r.Spec.StorageType != StorageTypeDurable && r.Spec.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				r.Name,
				"StorageType should be either durable or ephemeral"))
		}
	}

	if r.Spec.ConfigSecret != nil && r.Spec.ConfigSecret.Name == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("configSecret").Child("name"),
			r.Name,
			"ConfigSecret Name can not be empty"))
	}

	if len(allErr) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "rabbitmq.kubedb.com", Kind: "RabbitMQ"}, r.Name, allErr)
}

func (r *RabbitMQ) ValidateVersion(db *RabbitMQ) error {
	rmVersion := catalog.RabbitMQVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: db.Spec.Version}, &rmVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	return nil
}

var rabbitmqReservedVolumes = []string{
	kubedb.RabbitMQVolumeData,
	kubedb.RabbitMQVolumeConfig,
	kubedb.RabbitMQVolumeTempConfig,
}

func (r *RabbitMQ) validateVolumes(db *RabbitMQ) error {
	if db.Spec.PodTemplate.Spec.Volumes == nil {
		return nil
	}
	rsv := make([]string, len(rabbitmqReservedVolumes))
	copy(rsv, rabbitmqReservedVolumes)
	if db.Spec.TLS != nil && db.Spec.TLS.Certificates != nil {
		for _, c := range db.Spec.TLS.Certificates {
			rsv = append(rsv, db.CertSecretVolumeName(RabbitMQCertificateAlias(c.Alias)))
		}
	}
	volumes := db.Spec.PodTemplate.Spec.Volumes
	for _, rv := range rsv {
		for _, ugv := range volumes {
			if ugv.Name == rv {
				return errors.New("Cannot use a reserve volume name: " + rv)
			}
		}
	}
	return nil
}

var rabbitmqReservedVolumeMountPaths = []string{
	kubedb.RabbitMQConfigDir,
	kubedb.RabbitMQTempConfigDir,
	kubedb.RabbitMQDataDir,
	kubedb.RabbitMQPluginsDir,
	kubedb.RabbitMQCertDir,
}

func (r *RabbitMQ) validateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Containers == nil {
		return nil
	}

	for _, rvmp := range rabbitmqReservedVolumeMountPaths {
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

	for _, rvmp := range rabbitmqReservedVolumeMountPaths {
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
