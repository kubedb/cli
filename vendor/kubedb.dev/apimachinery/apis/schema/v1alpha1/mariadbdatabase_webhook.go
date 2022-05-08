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
	gocmp "github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var mariadbdatabaselog = logf.Log.WithName("mariadbdatabase-resource")

func (r *MariaDBDatabase) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-schema-kubedb-com-v1alpha1-mariadbdatabase,mutating=true,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=mariadbdatabases,verbs=create;update,versions=v1alpha1,name=mmariadbdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &MariaDBDatabase{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MariaDBDatabase) Default() {
	mariadbdatabaselog.Info("default", "name", r.Name)

	if r.Spec.Init != nil {
		if r.Spec.Init.Snapshot != nil {
			if r.Spec.Init.Snapshot.SnapshotID == "" {
				r.Spec.Init.Snapshot.SnapshotID = "latest"
			}
		}
	}
	if r.Spec.Database.Config.CharacterSet == "" {
		r.Spec.Database.Config.CharacterSet = "utf8mb4"
	}
}

//+kubebuilder:webhook:path=/validate-schema-kubedb-com-v1alpha1-mariadbdatabase,mutating=false,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=mariadbdatabases,verbs=create;update;delete,versions=v1alpha1,name=vmariadbdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MariaDBDatabase{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MariaDBDatabase) ValidateCreate() error {
	mariadbdatabaselog.Info("validate create", "name", r.Name)
	var allErrs field.ErrorList
	if err := r.ValidateMariaDBDatabase(); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath(""), r.Name, err.Error()))
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "schema.kubedb.com", Kind: "MariaDBDatabase"}, r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MariaDBDatabase) ValidateUpdate(old runtime.Object) error {
	mariadbdatabaselog.Info("validate update", "name", r.Name)
	oldobj := old.(*MariaDBDatabase)
	return ValidateMariaDBDatabaseUpdate(r, oldobj)
}

func ValidateMariaDBDatabaseUpdate(newobj *MariaDBDatabase, oldobj *MariaDBDatabase) error {
	if newobj.Finalizers == nil {
		return nil
	}
	var allErrs field.ErrorList
	if !gocmp.Equal(oldobj.Spec.Database.ServerRef, newobj.Spec.Database.ServerRef) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.database.serverRef"), newobj.Name, "cannot change database serverRef"))
	}
	if !gocmp.Equal(oldobj.Spec.Database.Config.Name, newobj.Spec.Database.Config.Name) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.database.config.name"), newobj.Name, "cannot change database name configuration"))
	}
	if !gocmp.Equal(oldobj.Spec.VaultRef, newobj.Spec.VaultRef) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.vaultRef"), newobj.Name, "cannot change vaultRef"))
	}
	if !gocmp.Equal(oldobj.Spec.AccessPolicy, newobj.Spec.AccessPolicy) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.accessPolicy"), newobj.Name, "cannot change accessPolicy"))
	}
	if newobj.Spec.Init != nil {
		if oldobj.Spec.Init == nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.init"), newobj.Name, "cannot change init"))
		}
	}
	if oldobj.Spec.Init != nil {
		if newobj.Spec.Init == nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.init"), newobj.Name, "cannot change init"))
		}
	}
	if newobj.Spec.Init != nil && oldobj.Spec.Init != nil {
		if !gocmp.Equal(newobj.Spec.Init.Script, oldobj.Spec.Init.Script) || !gocmp.Equal(newobj.Spec.Init.Snapshot, oldobj.Spec.Init.Snapshot) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.init"), newobj.Name, "cannot change init"))
		}
	}
	er := newobj.ValidateMariaDBDatabase()
	if er != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), newobj.Name, er.Error()))
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "schema.kubedb.com", Kind: "MariaDBDatabase"}, newobj.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MariaDBDatabase) ValidateDelete() error {
	mariadbdatabaselog.Info("validate delete", "name", r.Name)
	if r.Spec.DeletionPolicy == DeletionPolicyDoNotDelete {
		return field.Invalid(field.NewPath("spec").Child("terminationPolicy"), r.Name, `cannot delete object when terminationPolicy is set to "DoNotDelete"`)
	}
	return nil
}

func (in *MariaDBDatabase) ValidateMariaDBDatabase() error {
	var allErrs field.ErrorList
	if err := in.validateInitailizationSchema(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := in.validateMariaDBDatabaseConfig(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "schema.kubedb.com", Kind: "MariaDBDatabase"}, in.Name, allErrs)
}

func (in *MariaDBDatabase) validateInitailizationSchema() *field.Error {
	path := field.NewPath("spec.init")
	if in.Spec.Init != nil {
		if in.Spec.Init.Script != nil && in.Spec.Init.Snapshot != nil {
			return field.Invalid(path, in.Name, `cannot initialize database using both restore and initSpec`)
		}
	}
	return nil
}

func (in *MariaDBDatabase) validateMariaDBDatabaseConfig() *field.Error {
	path := field.NewPath("spec").Child("database.config").Child("name")
	name := in.Spec.Database.Config.Name
	if name == SYSDatabase {
		return field.Invalid(path, in.Name, `cannot use "sys" as the database name`)
	}
	if name == "performance_schema" {
		return field.Invalid(path, in.Name, `cannot use "performance_schema" as the database name`)
	}
	if name == "mysql" {
		return field.Invalid(path, in.Name, `cannot use "mysql" as the database name`)
	}
	if name == DatabaseForEntry {
		return field.Invalid(path, in.Name, `cannot use "kubedb_system" as the database name`)
	}
	if name == "information_schema" {
		return field.Invalid(path, in.Name, `cannot use "information_schema" as the database name`)
	}
	if name == DatabaseNameAdmin {
		return field.Invalid(path, in.Name, `cannot use "admin" as the database name`)
	}
	if name == DatabaseNameConfig {
		return field.Invalid(path, in.Name, `cannot use "config" as the database name`)
	}
	return nil
}
