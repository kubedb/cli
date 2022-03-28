/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Free Trial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Free-Trial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"

	gocmp "github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var postgresdatabaselog = logf.Log.WithName("postgresdatabase-resource")

func (r *PostgresDatabase) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-schema-kubedb-com-v1alpha1-postgresdatabase,mutating=true,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=postgresdatabases,verbs=create;update,versions=v1alpha1,name=mpostgresdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &PostgresDatabase{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PostgresDatabase) Default() {
	postgresdatabaselog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-schema-kubedb-com-v1alpha1-postgresdatabase,mutating=false,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=postgresdatabases,verbs=create;update;delete,versions=v1alpha1,name=vpostgresdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &PostgresDatabase{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PostgresDatabase) ValidateCreate() error {
	postgresdatabaselog.Info("validate create", "name", r.Name)
	if r.Spec.Init != nil && r.Spec.Init.Initialized {
		return field.Invalid(field.NewPath("spec").Child("init").Child("initialized"), r.Spec.Init.Initialized, fmt.Sprintf(`can't set spec.init.initialized true while creating postgresSchema %s/%s`, r.Namespace, r.Name))
	}
	return r.ValidatePostgresDatabase()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PostgresDatabase) ValidateUpdate(old runtime.Object) error {
	postgresdatabaselog.Info("validate update", "name", r.Name)
	oldobj := old.(*PostgresDatabase)
	return r.ValidatePostgresDatabaseUpdate(oldobj, r)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PostgresDatabase) ValidateDelete() error {
	postgresdatabaselog.Info("validate delete", "name", r.Name)
	if r.Spec.DeletionPolicy == DeletionPolicyDoNotDelete {
		return field.Invalid(field.NewPath("spec").Child("deletionPolicy"), r.Spec.DeletionPolicy, fmt.Sprintf(`can't delete postgresSchema %s/%s when deletionPolicy is "DoNotDelete"`, r.Namespace, r.Name))
	}
	return nil
}

func (r *PostgresDatabase) ValidateReadOnly() *field.Error {
	if r.Spec.Database.Config.Params == nil {
		return nil
	}
	for _, param := range r.Spec.Database.Config.Params {
		if param.ConfigParameter == "default_transaction_read_only" && *param.Value == "on" && r.Spec.Init != nil && !r.Spec.Init.Initialized {
			return field.Invalid(field.NewPath("spec").Child("database").Child("config").Child("params"), r.Spec.Init, fmt.Sprintf("can't initialize a read-only database in postgresSchema %s/%s", r.Namespace, r.Name))
		}
	}
	return nil
}

func (r *PostgresDatabase) ValidatePostgresDatabaseUpdate(oldobj *PostgresDatabase, newobj *PostgresDatabase) error {
	if newobj.Finalizers == nil {
		return nil
	}
	path := field.NewPath("spec")
	if !gocmp.Equal(oldobj.Spec.Database.Config.Name, newobj.Spec.Database.Config.Name) {
		return field.Invalid(path.Child("database").Child("config").Child("name"), newobj.Spec.Database.Config.Name, fmt.Sprintf("can't change the database name in postgresSchema %s/%s", r.Namespace, r.Name))
	}
	if !gocmp.Equal(oldobj.Spec.Database.ServerRef, newobj.Spec.Database.ServerRef) {
		return field.Invalid(path.Child("database").Child("serverRef"), newobj.Spec.Database.ServerRef, fmt.Sprintf("can't change the kubedb server reference in postgresSchema %s/%s", r.Namespace, r.Name))
	}
	if !gocmp.Equal(oldobj.Spec.VaultRef, newobj.Spec.VaultRef) {
		return field.Invalid(path.Child("vaultRef"), newobj.Spec.VaultRef, fmt.Sprintf("can't change the vault reference in postgresSchema %s/%s", r.Namespace, r.Name))
	}
	if err := newobj.ValidatePostgresDatabase(); err != nil {
		return err
	}
	if oldobj.Spec.Init != nil && oldobj.Spec.Init.Initialized && !gocmp.Equal(oldobj.Spec.Init, newobj.Spec.Init) {
		return field.Invalid(path.Child("init"), newobj.Spec.Init, fmt.Sprintf("can't change spec.init in postgresSchema %s/%s, is already initialized", r.Namespace, r.Name))
	}
	return nil
}

func (r *PostgresDatabase) ValidatePostgresDBName() *field.Error {
	path := field.NewPath("spec").Child("database").Child("config").Child("name")
	name := r.Spec.Database.Config.Name
	if name == PostgresSchemaKubeSystem || name == DatabaseNameAdmin || name == DatabaseNameConfig || name == DatabaseNameLocal || name == "postgres" || name == "sys" || name == "template0" || name == "template1" {
		str := fmt.Sprintf("can't set spec.database.config.name \"%v\" in postgresSchema %s/%s", name, r.Namespace, r.Name)
		return field.Invalid(path, name, str)
	}
	return nil
}

func (r *PostgresDatabase) ValidateSchemaInitRestore() *field.Error {
	path := field.NewPath("spec").Child("init")
	if r.Spec.Init != nil && r.Spec.Init.Snapshot != nil && r.Spec.Init.Script != nil {
		return field.Invalid(path, r.Name, fmt.Sprintf("can't set both spec.init.snapshot and spec.init.script in postgresSchema %s/%s", r.Namespace, r.Name))
	}
	return nil
}

func (r *PostgresDatabase) ValidateParams() *field.Error {
	if r.Spec.Database.Config.Params == nil {
		return nil
	}
	for _, param := range r.Spec.Database.Config.Params {
		if param.ConfigParameter == "" || param.Value == nil {
			msg := fmt.Sprintf("can't set empty spec.database.config.params.configParameter or spec.database.config.params.value in postgresSchema %s/%s", r.Namespace, r.Name)
			return field.Invalid(field.NewPath("spec").Child("database").Child("config").Child("params"), r.Spec.Database.Config.Params, msg)
		}
	}
	return nil
}

func (r *PostgresDatabase) ValidateFields() *field.Error {
	if r.Spec.Database.ServerRef.Name == "" {
		str := fmt.Sprintf("spec.database.serverRef.Name can't set empty in postgresSchema %s/%s", r.Namespace, r.Name)
		return field.Invalid(field.NewPath("spec").Child("database").Child("serverRef").Child("name"), r.Spec.Database.ServerRef, str)
	}
	if r.Spec.VaultRef.Name == "" {
		str := fmt.Sprintf("spec.database.vaultRef.Name can't set empty in postgresSchema %s/%s", r.Namespace, r.Name)
		return field.Invalid(field.NewPath("spec").Child("vaultRef").Child("name"), r.Spec.VaultRef, str)
	}
	if r.Spec.Init != nil && r.Spec.Init.Snapshot != nil {
		if r.Spec.Init.Snapshot.Repository.Name == "" {
			str := fmt.Sprintf("spec.init.snapshot.repository.name can't set empty in postgresSchema %s/%s", r.Namespace, r.Name)
			return field.Invalid(field.NewPath("spec").Child("init").Child("snapshot").Child("repository").Child("name"), r.Spec.Init.Snapshot.Repository.Name, str)
		}
	}
	if r.Spec.AccessPolicy.Subjects == nil {
		str := fmt.Sprintf("spec.accessPolicy.subjects can't set empty in postgresSchema %s/%s", r.Namespace, r.Name)
		return field.Invalid(field.NewPath("spec").Child("accessPolicy").Child("subjects"), r.Spec.AccessPolicy.Subjects, str)
	}
	return nil
}

func (r *PostgresDatabase) ValidatePostgresDatabase() error {
	// check if Init and Restore both are present
	if err := r.ValidateSchemaInitRestore(); err != nil {
		return err
	}
	// check the database name is conflicted with some constant name
	if err := r.ValidatePostgresDBName(); err != nil {
		return err
	}
	// check the spec fields
	if err := r.ValidateFields(); err != nil {
		return err
	}
	// check configuration params
	if err := r.ValidateParams(); err != nil {
		return err
	}
	// check read-only
	if err := r.ValidateReadOnly(); err != nil {
		return err
	}
	return nil
}
