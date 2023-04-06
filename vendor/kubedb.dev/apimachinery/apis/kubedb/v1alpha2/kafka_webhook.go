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
	"errors"

	errors2 "github.com/pkg/errors"
	"gomodules.xyz/pointer"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var kafkalog = logf.Log.WithName("kafka-resource")

func (k *Kafka) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(k).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kafka-kubedb-com-v1alpha1-kafka,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=kafkas,verbs=create,versions=v1alpha1,name=mkafka.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &Kafka{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (k *Kafka) Default() {
	if k == nil {
		return
	}
	kafkalog.Info("default", "name", k.Name)
	k.SetDefaults()
}

//+kubebuilder:webhook:path=/validate-kafka-kubedb-com-v1alpha1-kafka,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=kafkas,verbs=create;update,versions=v1alpha1,name=vkafka.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Kafka{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (k *Kafka) ValidateCreate() error {
	kafkalog.Info("validate create", "name", k.Name)
	return k.ValidateCreateOrUpdate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (k *Kafka) ValidateUpdate(old runtime.Object) error {
	kafkalog.Info("validate update", "name", k.Name)
	return k.ValidateCreateOrUpdate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (k *Kafka) ValidateDelete() error {
	kafkalog.Info("validate delete", "name", k.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	var allErr field.ErrorList
	if k.Spec.TerminationPolicy == TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("teminationPolicy"),
			k.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "Kafka"}, k.Name, allErr)
	}
	return nil
}

func (k *Kafka) ValidateCreateOrUpdate() error {
	var allErr field.ErrorList
	// TODO(user): fill in your validation logic upon object creation.
	if k.Spec.Topology != nil {
		if k.Spec.Topology.Controller == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("controller"),
				k.Name,
				".spec.topology.controller can't be empty in topology cluster"))
		}
		if k.Spec.Topology.Broker == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("broker"),
				k.Name,
				".spec.topology.broker can't be empty in topology cluster"))
		}

		if k.Spec.Replicas != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				k.Name,
				"doesn't support spec.replicas when spec.topology is set"))
		}
		if k.Spec.Storage != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("broker"),
				k.Name,
				"doesn't support spec.storage when spec.topology is set"))
		}
		if k.Spec.PodTemplate.Spec.Resources.Size() != 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("resources"),
				k.Name,
				"doesn't support spec.podTemplate.spec.resources when spec.topology is set"))
		}

		if *k.Spec.Topology.Controller.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("controller").Child("replicas"),
				k.Name,
				"number of replicas can not be less be 0 or less"))
		}

		if *k.Spec.Topology.Broker.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("broker").Child("replicas"),
				k.Name,
				"number of replicas can not be 0 or less"))
		}

		// validate that multiple nodes don't have same suffixes
		err := validateNodeSuffix(k.Spec.Topology)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology"),
				k.Name,
				err.Error()))
		}

		// validate that node replicas are not 0 or negative
		err = validateNodeReplicas(k.Spec.Topology)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology"),
				k.Name,
				err.Error()))
		}
	} else {
		// number of replicas can not be 0 or less
		if k.Spec.Replicas != nil && *k.Spec.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				k.Name,
				"number of replicas can not be 0 or less"))
		}
	}

	err := validateVersion(k)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			k.Name,
			err.Error()))
	}

	err = validateVolumes(k)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
			k.Name,
			err.Error()))
	}

	err = validateVolumesMountPaths(k)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumeMounts"),
			k.Name,
			err.Error()))
	}

	if k.Spec.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			k.Name,
			"StorageType can not be empty"))
	} else {
		if k.Spec.StorageType != StorageTypeDurable && k.Spec.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				k.Name,
				"StorageType should be either durable or ephemeral"))
		}
	}

	if len(allErr) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "Kafka"}, k.Name, allErr)
}

var availableVersions = []string{
	"3.3.0",
	"3.3.2",
	"3.4.0",
}

func validateVersion(db *Kafka) error {
	version := db.Spec.Version
	for _, v := range availableVersions {
		if v == version {
			return nil
		}
	}
	return errors.New("version not supported")
}

func validateNodeSuffix(topology *KafkaClusterTopology) error {
	tMap := topology.ToMap()
	names := make(map[string]bool)
	for _, value := range tMap {
		names[value.Suffix] = true
	}
	if len(tMap) != len(names) {
		return errors.New("two or more node cannot have same suffix")
	}
	return nil
}

func validateNodeReplicas(topology *KafkaClusterTopology) error {
	tMap := topology.ToMap()
	for key, node := range tMap {
		if pointer.Int32(node.Replicas) <= 0 {
			return errors2.Errorf("replicas for node role %s must be alteast 1", string(key))
		}
	}
	return nil
}

var reservedVolumes = []string{
	KafkaVolumeData,
	KafkaVolumeConfig,
	KafkaVolumeTempConfig,
}

func validateVolumes(db *Kafka) error {
	if db.Spec.PodTemplate.Spec.Volumes == nil {
		return nil
	}
	rsv := make([]string, len(reservedVolumes))
	copy(rsv, reservedVolumes)
	if db.Spec.TLS != nil && db.Spec.TLS.Certificates != nil {
		for _, c := range db.Spec.TLS.Certificates {
			rsv = append(rsv, db.CertSecretVolumeName(KafkaCertificateAlias(c.Alias)))
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

var reservedVolumeMountPaths = []string{
	KafkaConfigDir,
	KafkaTempConfigDir,
	KafkaDataDir,
	KafkaMetaDataDir,
	KafkaCertDir,
}

func validateVolumesMountPaths(db *Kafka) error {
	if db.Spec.PodTemplate.Spec.VolumeMounts == nil {
		return nil
	}
	rPaths := reservedVolumeMountPaths
	volumeMountPaths := db.Spec.PodTemplate.Spec.VolumeMounts
	for _, rvm := range rPaths {
		for _, ugv := range volumeMountPaths {
			if ugv.Name == rvm {
				return errors.New("Cannot use a reserve volume name: " + rvm)
			}
		}
	}
	return nil
}
