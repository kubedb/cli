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

	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"

	"gomodules.xyz/x/arrays"
	core "k8s.io/api/core/v1"
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
var ferretdblog = logf.Log.WithName("ferretdb-resource")

func (f *FerretDB) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(f).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kubedb-com-v1alpha2-ferretdb,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=ferretdbs,verbs=create;update,versions=v1alpha2,name=mferretdb.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &FerretDB{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (f *FerretDB) Default() {
	if f == nil {
		return
	}
	ferretdblog.Info("default", "name", f.Name)
	f.SetDefaults()
}

//+kubebuilder:webhook:path=/validate-kubedb-com-v1alpha2-ferretdb,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubedb.com,resources=ferretdbs,verbs=create;update;delete,versions=v1alpha2,name=vferretdb.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &FerretDB{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (f *FerretDB) ValidateCreate() (admission.Warnings, error) {
	ferretdblog.Info("validate create", "name", f.Name)

	allErr := f.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "FerretDB"}, f.Name, allErr)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (f *FerretDB) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	ferretdblog.Info("validate update", "name", f.Name)

	_ = old.(*FerretDB)
	allErr := f.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "FerretDB"}, f.Name, allErr)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (f *FerretDB) ValidateDelete() (admission.Warnings, error) {
	ferretdblog.Info("validate delete", "name", f.Name)

	var allErr field.ErrorList
	if f.Spec.TerminationPolicy == TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			f.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "FerretDB"}, f.Name, allErr)
	}
	return nil, nil
}

func (f *FerretDB) ValidateCreateOrUpdate() field.ErrorList {
	var allErr field.ErrorList

	err := f.validateFerretDBVersion()
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			f.Name,
			err.Error()))
	}
	if f.Spec.Replicas == nil || *f.Spec.Replicas < 1 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
			f.Name,
			fmt.Sprintf(`spec.replicas "%v" invalid. Must be greater than zero`, f.Spec.Replicas)))
	}

	if f.Spec.PodTemplate != nil {
		if err := FerretDBValidateEnvVar(getMainContainerEnvs(f), forbiddenEnvVars, f.ResourceKind()); err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate"),
				f.Name,
				err.Error()))
		}
	}

	// Storage related
	if f.Spec.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			f.Name,
			`'spec.storageType' is missing`))
	}
	if f.Spec.StorageType != StorageTypeDurable && f.Spec.StorageType != StorageTypeEphemeral {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			f.Name,
			fmt.Sprintf(`'spec.storageType' %s is invalid`, f.Spec.StorageType)))
	}
	if f.Spec.StorageType == StorageTypeEphemeral && f.Spec.Storage != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			f.Name,
			`'spec.storageType' is set to Ephemeral, so 'spec.storage' needs to be empty`))
	}
	if !f.Spec.Backend.ExternallyManaged && f.Spec.StorageType == StorageTypeDurable && f.Spec.Storage == nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storage"),
			f.Name,
			`'spec.storage' is missing for durable storage type when postgres is internally managed`))
	}

	// Auth secret related
	if f.Spec.AuthSecret != nil && f.Spec.AuthSecret.ExternallyManaged != f.Spec.Backend.ExternallyManaged {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("authSecret"),
			f.Name,
			`when 'spec.backend' is internally managed, 'spec.authSecret' can't be externally managed and vice versa`))
	}
	if f.Spec.AuthSecret != nil && f.Spec.AuthSecret.ExternallyManaged && f.Spec.AuthSecret.Name == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("authSecret"),
			f.Name,
			`'spec.authSecret.name' needs to specify when auth secret is externally managed`))
	}

	if f.Spec.StorageType == StorageTypeEphemeral && f.Spec.TerminationPolicy == TerminationPolicyHalt {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			f.Name,
			`'spec.terminationPolicy: Halt' can not be used for 'Ephemeral' storage`))
	}
	if f.Spec.TerminationPolicy == TerminationPolicyHalt || f.Spec.TerminationPolicy == TerminationPolicyDelete {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			f.Name,
			`'spec.terminationPolicy' value 'Halt' or 'Delete' is not supported yet for FerretDB`))
	}

	// Backend related
	if f.Spec.Backend.ExternallyManaged {
		if f.Spec.Backend.Postgres == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("backend"),
				f.Name,
				`'spec.postgres' is missing when backend is externally managed`))
		}
		if f.Spec.Backend.Postgres != nil && f.Spec.Backend.Postgres.URL == nil && f.Spec.Backend.Postgres.Service == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("backend"),
				f.Name,
				`Have to provide 'backend.postgres.url' or 'backend.postgres.service' when backend is externally managed`))
		}
	}
	if !f.Spec.Backend.ExternallyManaged && f.Spec.Backend.Postgres != nil && f.Spec.Backend.Postgres.Version != nil {
		err := f.validatePostgresVersion()
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("backend"),
				f.Name,
				err.Error()))
		}
	}

	if f.Spec.SSLMode == SSLModeAllowSSL || f.Spec.SSLMode == SSLModePreferSSL {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("sslMode"),
			f.Name,
			`'spec.sslMode' value 'allowSSL' or 'preferSSL' is not supported yet for FerretDB`))
	}

	return allErr
}

func FerretDBValidateEnvVar(envs []core.EnvVar, forbiddenEnvs []string, resourceType string) error {
	for _, env := range envs {
		present, _ := arrays.Contains(forbiddenEnvs, env.Name)
		if present {
			return fmt.Errorf("environment variable %s is forbidden to use in %s spec", env.Name, resourceType)
		}
	}
	return nil
}

var forbiddenEnvVars = []string{
	EnvFerretDBUser, EnvFerretDBPassword, EnvFerretDBHandler, EnvFerretDBPgURL,
	EnvFerretDBTLSPort, EnvFerretDBCAPath, EnvFerretDBCertPath, EnvFerretDBKeyPath,
}

func getMainContainerEnvs(f *FerretDB) []core.EnvVar {
	for _, container := range f.Spec.PodTemplate.Spec.Containers {
		if container.Name == FerretDBContainerName {
			return container.Env
		}
	}
	return []core.EnvVar{}
}

func (f *FerretDB) validateFerretDBVersion() error {
	frVersion := v1alpha1.FerretDBVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: f.Spec.Version}, &frVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	return nil
}

func (f *FerretDB) validatePostgresVersion() error {
	pgVersion := v1alpha1.PostgresVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: *f.Spec.Backend.Postgres.Version}, &pgVersion)
	if err != nil {
		return errors.New("postgres version not supported in KubeDB")
	}
	return nil
}
