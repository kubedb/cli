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
	"kubedb.dev/apimachinery/apis/kubedb"

	"gomodules.xyz/x/arrays"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
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
	if f.Spec.DeletionPolicy == DeletionPolicyDoNotTerminate {
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
	if f.Spec.AuthSecret != nil && f.Spec.AuthSecret.ExternallyManaged && f.Spec.AuthSecret.Name == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("authSecret"),
			f.Name,
			`'spec.authSecret.name' need to specify when auth secret is externally managed`))
	}

	// Termination policy related
	if f.Spec.DeletionPolicy == DeletionPolicyHalt {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			f.Name,
			`'spec.terminationPolicy' value 'Halt' is not supported yet for FerretDB`))
	}

	// FerretDBBackend related
	if f.Spec.Backend.ExternallyManaged {
		if f.Spec.Backend.PostgresRef == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("backend"),
				f.Name,
				`'backend.postgresRef' is missing when backend is externally managed`))
		} else {
			if f.Spec.Backend.PostgresRef.Namespace == "" {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("backend"),
					f.Name,
					`'backend.postgresRef.namespace' is needed when backend is externally managed`))
			}
			apb := appcat.AppBinding{}
			err := DefaultClient.Get(context.TODO(), types.NamespacedName{
				Name:      f.Spec.Backend.PostgresRef.Name,
				Namespace: f.Spec.Backend.PostgresRef.Namespace,
			}, &apb)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("postgresRef"),
					f.Name,
					err.Error(),
				))
			}

			if apb.Spec.Secret == nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("backend"),
					f.Name,
					`spec.secret needed in external pg appbinding`))
			}

			if apb.Spec.ClientConfig.Service == nil && apb.Spec.ClientConfig.URL == nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("postgresRef"),
					f.Name,
					`'clientConfig.url' or 'clientConfig.service' needed in the external pg appbinding`,
				))
			}
			sslMode, err := f.GetSSLModeFromAppBinding(&apb)
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("postgresRef"),
					f.Name,
					err.Error(),
				))
			}

			if sslMode == PostgresSSLModeRequire || sslMode == PostgresSSLModeVerifyCA || sslMode == PostgresSSLModeVerifyFull {
				if apb.Spec.ClientConfig.CABundle == nil && apb.Spec.TLSSecret == nil {
					allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("postgresRef"),
						f.Name,
						"backend postgres connection is ssl encrypted but 'spec.clientConfig.caBundle' or 'spec.tlsSecret' is not provided in appbinding",
					))
				}
			}
			if (apb.Spec.ClientConfig.CABundle != nil || apb.Spec.TLSSecret != nil) && sslMode == PostgresSSLModeDisable {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("postgresRef"),
					f.Name,
					"no client certificate or ca bundle possible when sslMode set to disable in backend postgres",
				))
			}
		}
	} else {
		if f.Spec.Backend.Version != nil {
			err := f.validatePostgresVersion()
			if err != nil {
				allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("backend"),
					f.Name,
					err.Error()))
			}
		}
	}

	// TLS related
	if f.Spec.SSLMode == SSLModeAllowSSL || f.Spec.SSLMode == SSLModePreferSSL {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("sslMode"),
			f.Name,
			`'spec.sslMode' value 'allowSSL' or 'preferSSL' is not supported yet for FerretDB`))
	}
	if f.Spec.SSLMode == SSLModeRequireSSL && f.Spec.TLS == nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("sslMode"),
			f.Name,
			`'spec.sslMode' is requireSSL but 'spec.tls' is not set`))
	}
	if f.Spec.SSLMode == SSLModeDisabled && f.Spec.TLS != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("sslMode"),
			f.Name,
			`'spec.tls' is can't set when 'spec.sslMode' is disabled`))
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
	kubedb.EnvFerretDBUser, kubedb.EnvFerretDBPassword, kubedb.EnvFerretDBHandler, kubedb.EnvFerretDBPgURL,
	kubedb.EnvFerretDBTLSPort, kubedb.EnvFerretDBCAPath, kubedb.EnvFerretDBCertPath, kubedb.EnvFerretDBKeyPath,
}

func getMainContainerEnvs(f *FerretDB) []core.EnvVar {
	for _, container := range f.Spec.PodTemplate.Spec.Containers {
		if container.Name == kubedb.FerretDBContainerName {
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
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: *f.Spec.Backend.Version}, &pgVersion)
	if err != nil {
		return errors.New("postgres version not supported in KubeDB")
	}
	return nil
}
