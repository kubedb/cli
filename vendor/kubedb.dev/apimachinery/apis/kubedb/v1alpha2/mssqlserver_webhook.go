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
	"fmt"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"

	"gomodules.xyz/x/arrays"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var mssqllog = logf.Log.WithName("mssql-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *MSSQLServer) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kubedb-com-v1alpha2-mssqlserver,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=mssqlservers,verbs=create;update,versions=v1alpha2,name=mmssqlserver.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &MSSQLServer{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (m *MSSQLServer) Default() {
	if m == nil {
		return
	}
	mssqllog.Info("default", "name", m.Name)

	m.SetDefaults()
}

//+kubebuilder:webhook:path=/validate-kubedb-com-v1alpha2-mssqlserver,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=mssqlservers,verbs=create;update,versions=v1alpha2,name=vmssqlserver.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &MSSQLServer{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (m *MSSQLServer) ValidateCreate() (admission.Warnings, error) {
	mssqllog.Info("validate create", "name", m.Name)

	allErr := m.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: kubedb.GroupName, Kind: ResourceKindMSSQLServer}, m.Name, allErr)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (m *MSSQLServer) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	mssqllog.Info("validate update", "name", m.Name)

	allErr := m.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(schema.GroupKind{Group: kubedb.GroupName, Kind: ResourceKindMSSQLServer}, m.Name, allErr)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (m *MSSQLServer) ValidateDelete() (admission.Warnings, error) {
	mssqllog.Info("validate delete", "name", m.Name)

	var allErr field.ErrorList
	if m.Spec.DeletionPolicy == DeletionPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			m.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: kubedb.GroupName, Kind: ResourceKindMSSQLServer}, m.Name, allErr)
	}
	return nil, nil
}

func (m *MSSQLServer) ValidateCreateOrUpdate() field.ErrorList {
	var allErr field.ErrorList

	err := mssqlValidateVersion(m)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			m.Name,
			err.Error()))
	}

	if m.IsStandalone() {
		if ptr.Deref(m.Spec.Replicas, 0) != 1 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				m.Name,
				"number of replicas for standalone must be one "))
		}
	} else {
		if m.Spec.Topology.Mode == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("mode"),
				m.Name,
				".spec.topology.mode can't be empty in cluster mode"))
		}

		if ptr.Deref(m.Spec.Replicas, 0) <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				m.Name,
				"number of replicas can not be nil and can not be less than or equal to 0"))
		}
	}

	if m.Spec.TLS == nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("tls"),
			m.Name, "spec.tls is missing"))
	} else {
		if m.Spec.TLS.IssuerRef == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("tls").Child("issuerRef"),
				m.Name, "spec.tls.issuerRef' is missing"))
		}
	}

	if m.Spec.PodTemplate != nil {
		if err = ValidateMSSQLServerEnvVar(getMSSQLServerContainerEnvs(m), forbiddenMSSQLServerEnvVars, m.ResourceKind()); err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate"),
				m.Name,
				err.Error()))
		}
	}

	err = mssqlValidateVolumes(m.Spec.PodTemplate)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
			m.Name,
			err.Error()))
	}

	err = mssqlValidateVolumesMountPaths(m.Spec.PodTemplate)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("containers"),
			m.Name,
			err.Error()))
	}

	if m.Spec.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			m.Name,
			"StorageType can not be empty"))
	} else {
		if m.Spec.StorageType != StorageTypeDurable && m.Spec.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				m.Name,
				"StorageType should be either durable or ephemeral"))
		}
	}

	if len(allErr) == 0 {
		return nil
	}

	return allErr
}

// reserved volume and volumes mounts for mssql
var mssqlReservedVolumes = []string{
	kubedb.MSSQLVolumeNameData,
	kubedb.MSSQLVolumeNameConfig,
	kubedb.MSSQLVolumeNameInitScript,
	kubedb.MSSQLVolumeNameEndpointCert,
	kubedb.MSSQLVolumeNameCerts,
	kubedb.MSSQLVolumeNameTLS,
	kubedb.MSSQLVolumeNameSecurityCACertificates,
	kubedb.MSSQLVolumeNameCACerts,
}

var mssqlReservedVolumesMountPaths = []string{
	kubedb.MSSQLVolumeMountPathData,
	kubedb.MSSQLVolumeMountPathConfig,
	kubedb.MSSQLVolumeMountPathInitScript,
	kubedb.MSSQLVolumeMountPathEndpointCert,
	kubedb.MSSQLVolumeMountPathCerts,
	kubedb.MSSQLVolumeMountPathTLS,
	kubedb.MSSQLVolumeMountPathSecurityCACertificates,
	kubedb.MSSQLVolumeMountPathCACerts,
}

func mssqlValidateVersion(m *MSSQLServer) error {
	var mssqlVersion catalog.MSSQLServerVersion

	return DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name: m.Spec.Version,
	}, &mssqlVersion)
}

func mssqlValidateVolumes(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Volumes == nil {
		return nil
	}

	for _, rv := range mssqlReservedVolumes {
		for _, ugv := range podTemplate.Spec.Volumes {
			if ugv.Name == rv {
				return errors.New("Can't use a reserved volume name: " + rv)
			}
		}
	}

	return nil
}

func mssqlValidateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}

	if podTemplate.Spec.Containers != nil {
		// Check container volume mounts
		for _, rvmp := range mssqlReservedVolumesMountPaths {
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
	}

	if podTemplate.Spec.InitContainers != nil {
		// Check init container volume mounts
		for _, rvmp := range mssqlReservedVolumesMountPaths {
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
	}

	return nil
}

var forbiddenMSSQLServerEnvVars = []string{
	kubedb.EnvMSSQLSAUsername,
	kubedb.EnvMSSQLSAPassword,
	kubedb.EnvMSSQLEnableHADR,
	kubedb.EnvMSSQLAgentEnabled,
	kubedb.EnvMSSQLVersion,
}

func getMSSQLServerContainerEnvs(m *MSSQLServer) []core.EnvVar {
	for _, container := range m.Spec.PodTemplate.Spec.Containers {
		if container.Name == kubedb.MSSQLContainerName {
			return container.Env
		}
	}
	return []core.EnvVar{}
}

func ValidateMSSQLServerEnvVar(envs []core.EnvVar, forbiddenEnvs []string, resourceType string) error {
	presentMSSQL_PID := false
	for _, env := range envs {
		present, _ := arrays.Contains(forbiddenEnvs, env.Name)
		if present {
			return fmt.Errorf("environment variable %s is forbidden to use in %s spec", env.Name, resourceType)
		}
		if env.Name == "MSSQL_PID" {
			presentMSSQL_PID = true
		}
	}
	if !presentMSSQL_PID {
		return fmt.Errorf("environment variable %s must be provided in %s spec", kubedb.EnvMSSQLPid, resourceType)
	}
	return nil
}
