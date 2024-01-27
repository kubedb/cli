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
	"strings"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	coreutil "kmodules.xyz/client-go/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var connectClusterLog = logf.Log.WithName("connectCluster-resource")

var _ webhook.Defaulter = &ConnectCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (k *ConnectCluster) Default() {
	if k == nil {
		return
	}
	connectClusterLog.Info("default", "name", k.Name)
	k.SetDefaults()
}

var _ webhook.Validator = &ConnectCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (k *ConnectCluster) ValidateCreate() (admission.Warnings, error) {
	connectClusterLog.Info("validate create", "name", k.Name)
	allErr := k.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "Kafka"}, k.Name, allErr)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (k *ConnectCluster) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	connectClusterLog.Info("validate update", "name", k.Name)

	oldConnect := old.(*ConnectCluster)
	allErr := k.ValidateCreateOrUpdate()

	if *oldConnect.Spec.Replicas == 1 && *k.Spec.Replicas > 1 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			k.Name,
			"Cannot scale up from 1 to more than 1 in standalone mode"))
	}

	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "ConnectCluster"}, k.Name, allErr)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (k *ConnectCluster) ValidateDelete() (admission.Warnings, error) {
	connectClusterLog.Info("validate delete", "name", k.Name)

	var allErr field.ErrorList
	if k.Spec.TerminationPolicy == api.TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			k.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "ConnectCluster"}, k.Name, allErr)
	}
	return nil, nil
}

func (k *ConnectCluster) ValidateCreateOrUpdate() field.ErrorList {
	var allErr field.ErrorList
	if k.Spec.EnableSSL {
		if k.Spec.TLS == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("enableSSL"),
				k.Name,
				".spec.tls can't be nil, if .spec.enableSSL is true"))
		}
	} else {
		if k.Spec.TLS != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("enableSSL"),
				k.Name,
				".spec.tls must be nil, if .spec.enableSSL is disabled"))
		}
	}

	if k.Spec.TerminationPolicy == api.TerminationPolicyHalt {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			k.Name,
			"TerminationPolicyHalt is not supported for ConnectCluster"))
	}

	// number of replicas can not be 0 or less
	if k.Spec.Replicas != nil && *k.Spec.Replicas <= 0 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			k.Name,
			"number of replicas can not be 0 or less"))
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

	err = validateContainerVolumeMountPaths(k)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("containers").Child("volumeMounts"),
			k.Name,
			err.Error()))
	}

	err = validateInitContainerVolumeMountPaths(k)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("initContainers").Child("volumeMounts"),
			k.Name,
			err.Error()))
	}

	err = validateEnvVars(k)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("envs"),
			k.Name,
			err.Error()))
	}

	if len(allErr) == 0 {
		return nil
	}
	return allErr
}

func validateEnvVars(connect *ConnectCluster) error {
	container := coreutil.GetContainerByName(connect.Spec.PodTemplate.Spec.Containers, ConnectClusterContainerName)
	if container == nil {
		return errors.New("container not found")
	}
	env := coreutil.GetEnvByName(container.Env, ConnectClusterModeEnv)
	if env != nil {
		if *connect.Spec.Replicas > 1 && env.Value == string(ConnectClusterNodeRoleStandalone) {
			return errors.New("can't use standalone mode as env, if replicas is more than 1")
		}
	}
	return nil
}

func validateVersion(connect *ConnectCluster) error {
	kccVersion := &catalog.KafkaVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: connect.Spec.Version}, kccVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	return nil
}

var reservedVolumes = []string{
	ConnectClusterOperatorVolumeConfig,
	ConnectClusterCustomVolumeConfig,
	ConnectorPluginsVolumeName,
	ConnectClusterAuthSecretVolumeName,
	ConnectClusterOffsetFileDirName,
	KafkaClientCertVolumeName,
	ConnectClusterServerCertsVolumeName,
}

func validateVolumes(connect *ConnectCluster) error {
	if connect.Spec.PodTemplate.Spec.Volumes == nil {
		return nil
	}
	rsv := make([]string, len(reservedVolumes))
	copy(rsv, reservedVolumes)
	volumes := connect.Spec.PodTemplate.Spec.Volumes
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
	ConnectClusterOperatorConfigPath,
	ConnectorPluginsVolumeDir,
	ConnectClusterAuthSecretVolumePath,
	ConnectClusterOffsetFileDir,
	ConnectClusterCustomConfigPath,
	KafkaClientCertDir,
	ConnectClusterServerCertVolumeDir,
}

func validateContainerVolumeMountPaths(connect *ConnectCluster) error {
	container := coreutil.GetContainerByName(connect.Spec.PodTemplate.Spec.Containers, ConnectClusterContainerName)
	if container == nil {
		return errors.New("container not found")
	}
	rPaths := reservedVolumeMountPaths
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

func validateInitContainerVolumeMountPaths(connect *ConnectCluster) error {
	for _, name := range connect.Spec.ConnectorPlugins {
		connectorVersion := &catalog.KafkaConnectorVersion{}
		err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: name}, connectorVersion)
		if err != nil {
			klog.Errorf("can't get the kafka connector version object %s for %s \n", err.Error(), name)
			return errors.New("no connector version found for " + name)
		}
		initContainer := coreutil.GetContainerByName(connect.Spec.PodTemplate.Spec.InitContainers, strings.ToLower(connectorVersion.Spec.Type))
		if initContainer == nil {
			return errors.New("init container not found for " + strings.ToLower(connectorVersion.Spec.Type))
		}
		volumeMount := coreutil.GetVolumeMountByName(initContainer.VolumeMounts, ConnectorPluginsVolumeName)
		if volumeMount != nil && volumeMount.MountPath == ConnectorPluginsVolumeDir {
			return errors.New("Cannot use a reserve volume mount path: " + ConnectorPluginsVolumeDir)
		}
	}
	return nil
}
