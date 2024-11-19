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
	"fmt"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"

	"github.com/pkg/errors"
	"gomodules.xyz/x/arrays"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	kmapi "kmodules.xyz/client-go/api/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var pgpoollog = logf.Log.WithName("pgpool-resource")

func (p *Pgpool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(p).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kubedb-com-v1alpha2-pgpool,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=pgpools,verbs=create;update,versions=v1alpha2,name=mpgpool.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Pgpool{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (p *Pgpool) Default() {
	pgpoollog.Info("default", "name", p.Name)
	p.SetDefaults()
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-kubedb-com-v1alpha2-pgpool,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=pgpools,verbs=create;update;delete,versions=v1alpha2,name=vpgpool.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Pgpool{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (p *Pgpool) ValidateCreate() (admission.Warnings, error) {
	pgpoollog.Info("validate create", "name", p.Name)
	errorList := p.ValidateCreateOrUpdate()
	if len(errorList) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Pgpool"}, p.Name, errorList)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (p *Pgpool) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	pgpoollog.Info("validate update", "name", p.Name)

	errorList := p.ValidateCreateOrUpdate()
	if len(errorList) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Pgpool"}, p.Name, errorList)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (p *Pgpool) ValidateDelete() (admission.Warnings, error) {
	pgpoollog.Info("validate delete", "name", p.Name)

	var errorList field.ErrorList
	if p.Spec.DeletionPolicy == DeletionPolicyDoNotTerminate {
		errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			p.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Pgpool"}, p.Name, errorList)
	}
	return nil, nil
}

func (p *Pgpool) ValidateCreateOrUpdate() field.ErrorList {
	var errorList field.ErrorList
	if p.Spec.Version == "" {
		errorList = append(errorList, field.Required(field.NewPath("spec").Child("version"),
			"`spec.version` is missing",
		))
	} else {
		err := PgpoolValidateVersion(p)
		if err != nil {
			errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("version"),
				p.Name,
				err.Error()))
		}
	}

	if p.Spec.PostgresRef == nil {
		errorList = append(errorList, field.Required(field.NewPath("spec").Child("postgresRef"),
			"`spec.postgresRef` is missing",
		))
	}

	if p.Spec.ConfigSecret != nil && (p.Spec.InitConfiguration != nil && p.Spec.InitConfiguration.PgpoolConfig != nil) {
		errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("configSecret"),
			p.Name,
			"use either `spec.configSecret` or `spec.initConfig`"))
		errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("initConfig"),
			p.Name,
			"use either `spec.configSecret` or `spec.initConfig`"))
	}

	if p.ObjectMeta.DeletionTimestamp == nil {
		apb := appcat.AppBinding{}
		err := DefaultClient.Get(context.TODO(), types.NamespacedName{
			Name:      p.Spec.PostgresRef.Name,
			Namespace: p.Spec.PostgresRef.Namespace,
		}, &apb)
		if err != nil {
			errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("postgresRef"),
				p.Name,
				err.Error(),
			))
		}

		backendSSL, err := p.IsBackendTLSEnabled()
		if err != nil {
			errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("postgresRef"),
				p.Name,
				err.Error(),
			))
		}

		if p.Spec.TLS == nil && backendSSL {
			errorList = append(errorList, field.Required(field.NewPath("spec").Child("tls"),
				"`spec.tls` must be set because backend postgres is tls enabled",
			))
		}
	}

	if p.Spec.TLS == nil {
		if p.Spec.SSLMode != "disable" {
			errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("sslMode"),
				p.Name,
				"Tls is not enabled, enable it to use this sslMode",
			))
		}

		if p.Spec.ClientAuthMode == "cert" {
			errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("clientAuthMode"),
				p.Name,
				"Tls is not enabled, enable it to use this clientAuthMode",
			))
		}
	}

	if p.Spec.Replicas != nil {
		if *p.Spec.Replicas <= 0 {
			errorList = append(errorList, field.Required(field.NewPath("spec").Child("replicas"),
				"`spec.replica` must be greater than 0",
			))
		}
		if *p.Spec.Replicas > 9 {
			errorList = append(errorList, field.Required(field.NewPath("spec").Child("replicas"),
				"`spec.replica` must be less than 10",
			))
		}
	}

	if p.Spec.PodTemplate != nil {
		if err := p.ValidateEnvVar(PgpoolGetMainContainerEnvs(p), PgpoolForbiddenEnvVars, p.ResourceKind()); err != nil {
			errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("containers").Child("env"),
				p.Name,
				err.Error(),
			))
		}
		err := PgpoolValidateVolumes(p)
		if err != nil {
			errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
				p.Name,
				err.Error(),
			))
		}

		err = PgpoolValidateVolumesMountPaths(p.Spec.PodTemplate)
		if err != nil {
			errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumeMounts"),
				p.Name,
				err.Error()))
		}
	}

	if err := p.ValidateHealth(&p.Spec.HealthChecker); err != nil {
		errorList = append(errorList, field.Invalid(field.NewPath("spec").Child("healthChecker"),
			p.Name,
			err.Error(),
		))
	}

	if len(errorList) == 0 {
		return nil
	}
	return errorList
}

func (p *Pgpool) ValidateEnvVar(envs []core.EnvVar, forbiddenEnvs []string, resourceType string) error {
	for _, env := range envs {
		present, _ := arrays.Contains(forbiddenEnvs, env.Name)
		if present {
			return fmt.Errorf("environment variable %s is forbidden to use in %s spec", env.Name, resourceType)
		}
	}
	return nil
}

func (p *Pgpool) ValidateHealth(health *kmapi.HealthCheckSpec) error {
	if health.PeriodSeconds != nil && *health.PeriodSeconds <= 0 {
		return fmt.Errorf(`spec.healthCheck.periodSeconds: can not be less than 1`)
	}

	if health.TimeoutSeconds != nil && *health.TimeoutSeconds <= 0 {
		return fmt.Errorf(`spec.healthCheck.timeoutSeconds: can not be less than 1`)
	}

	if health.FailureThreshold != nil && *health.FailureThreshold <= 0 {
		return fmt.Errorf(`spec.healthCheck.failureThreshold: can not be less than 1`)
	}
	return nil
}

func PgpoolValidateVersion(p *Pgpool) error {
	ppVersion := catalog.PgpoolVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name: p.Spec.Version,
	}, &ppVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	return nil
}

var PgpoolReservedVolumes = []string{
	kubedb.PgpoolConfigVolumeName,
	kubedb.PgpoolTlsVolumeName,
}

func PgpoolValidateVolumes(p *Pgpool) error {
	if p.Spec.PodTemplate.Spec.Volumes == nil {
		return nil
	}

	for _, rv := range PgpoolReservedVolumes {
		for _, ugv := range p.Spec.PodTemplate.Spec.Volumes {
			if ugv.Name == rv {
				return errors.New("Cannot use a reserve volume name: " + rv)
			}
		}
	}
	return nil
}

var PgpoolForbiddenEnvVars = []string{
	kubedb.EnvPostgresUsername, kubedb.EnvPostgresPassword, kubedb.EnvPgpoolPcpUser, kubedb.EnvPgpoolPcpPassword,
	kubedb.EnvPgpoolPasswordEncryptionMethod, kubedb.EnvEnablePoolPasswd, kubedb.EnvSkipPasswdEncryption,
}

func PgpoolGetMainContainerEnvs(p *Pgpool) []core.EnvVar {
	for _, container := range p.Spec.PodTemplate.Spec.Containers {
		if container.Name == kubedb.PgpoolContainerName {
			return container.Env
		}
	}
	return []core.EnvVar{}
}

func PgpoolValidateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Containers == nil {
		return nil
	}

	for _, rvmp := range PgpoolReservedVolumesMountPaths {
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

	for _, rvmp := range PgpoolReservedVolumesMountPaths {
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

var PgpoolReservedVolumesMountPaths = []string{
	kubedb.PgpoolConfigSecretMountPath,
	kubedb.PgpoolTlsVolumeMountPath,
}
