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

	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var zookeeperlog = logf.Log.WithName("zookeeper-resource")

func (z *ZooKeeper) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(z).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-zookeeper-kubedb-com-v1alpha1-zookeeper,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=zookeepers,verbs=create;update,versions=v1alpha1,name=mzookeeper.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &ZooKeeper{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (z *ZooKeeper) Default() {
	if z == nil {
		return
	}
	zookeeperlog.Info("default", "name", z.Name)
	z.SetDefaults()
}

//+kubebuilder:webhook:path=/validate-zookeeper-kubedb-com-v1alpha1-zookeeper,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=zookeepers,verbs=create;update,versions=v1alpha1,name=vzookeeper.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &ZooKeeper{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (z *ZooKeeper) ValidateCreate() (admission.Warnings, error) {
	zookeeperlog.Info("validate create", "name", z.Name)
	return z.ValidateCreateOrUpdate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (z *ZooKeeper) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	zookeeperlog.Info("validate update", "name", z.Name)
	return z.ValidateCreateOrUpdate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (z *ZooKeeper) ValidateDelete() (admission.Warnings, error) {
	zookeeperlog.Info("validate delete", "name", z.Name)

	var allErr field.ErrorList
	if z.Spec.TerminationPolicy == TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("teminationPolicy"),
			z.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "zookeeper.kubedb.com", Kind: "ZooKeeper"}, z.Name, allErr)
	}
	return nil, nil
}

func (z *ZooKeeper) ValidateCreateOrUpdate() (admission.Warnings, error) {
	var allErr field.ErrorList
	if z.Spec.Replicas != nil && *z.Spec.Replicas == 2 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			z.Name,
			"zookeeper ensemble should have 3 or more replicas"))
	}

	err := z.validateVersion()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			z.Name,
			err.Error()))
	}

	err = z.validateVolumes()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
			z.Name,
			err.Error()))
	}

	err = z.validateVolumesMountPaths()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumeMounts"),
			z.Name,
			err.Error()))
	}

	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "zookeeper.kubedb.com", Kind: "ZooKeeper"}, z.Name, allErr)
}

func (z *ZooKeeper) validateVersion() error {
	zkVersion := v1alpha1.ZooKeeperVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: z.Spec.Version}, &zkVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	return nil
}

var zookeeperReservedVolumes = []string{
	ZooKeeperDataVolumeName,
	ZooKeeperScriptVolumeName,
	ZooKeeperConfigVolumeName,
}

func (z *ZooKeeper) validateVolumes() error {
	if z.Spec.PodTemplate.Spec.Volumes == nil {
		return nil
	}
	rsv := make([]string, len(zookeeperReservedVolumes))
	copy(rsv, zookeeperReservedVolumes)

	volumes := z.Spec.PodTemplate.Spec.Volumes
	for _, rv := range rsv {
		for _, ugv := range volumes {
			if ugv.Name == rv {
				return errors.New("Cannot use a reserve volume name: " + rv)
			}
		}
	}
	return nil
}

var zookeeperReservedVolumeMountPaths = []string{
	ZooKeeperScriptVolumePath,
	ZooKeeperConfigVolumePath,
	ZooKeeperDataVolumePath,
}

func (z *ZooKeeper) validateVolumesMountPaths() error {
	if z.Spec.PodTemplate.Spec.VolumeMounts == nil {
		return nil
	}
	rPaths := zookeeperReservedVolumeMountPaths
	volumeMountPaths := z.Spec.PodTemplate.Spec.VolumeMounts
	for _, rvm := range rPaths {
		for _, ugv := range volumeMountPaths {
			if ugv.Name == rvm {
				return errors.New("Cannot use a reserve volume name: " + rvm)
			}
		}
	}
	return nil
}
