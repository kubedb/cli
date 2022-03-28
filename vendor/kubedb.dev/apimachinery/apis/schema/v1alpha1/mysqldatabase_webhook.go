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
var mysqldatabaselog = logf.Log.WithName("mysqldatabase-resource")

func (in *MySQLDatabase) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-schema-kubedb-com-v1alpha1-mysqldatabase,mutating=true,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=mysqldatabases,verbs=create;update,versions=v1alpha1,name=mmysqldatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &MySQLDatabase{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *MySQLDatabase) Default() {
	mysqldatabaselog.Info("default", "name", in.Name)

	if in.Spec.Init != nil {
		if in.Spec.Init.Snapshot != nil {
			if in.Spec.Init.Snapshot.SnapshotID == "" {
				in.Spec.Init.Snapshot.SnapshotID = "latest"
			}
		}
	}
	val := in.Spec.Database.Config.Encryption
	if val == "enable" || val == MySQLEncryptionEnabled {
		in.Spec.Database.Config.Encryption = MySQLEncryptionEnabled
	} else {
		in.Spec.Database.Config.Encryption = MySQLEncryptionDisabled
	}
	if in.Spec.Database.Config.ReadOnly != 1 {
		in.Spec.Database.Config.ReadOnly = 0
	}
	if in.Spec.Database.Config.CharacterSet == "" {
		in.Spec.Database.Config.CharacterSet = "utf8"
	}
}

// +kubebuilder:webhook:path=/validate-schema-kubedb-com-v1alpha1-mysqldatabase,mutating=false,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=mysqldatabases,verbs=create;update;delete,versions=v1alpha1,name=vmysqldatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MySQLDatabase{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *MySQLDatabase) ValidateCreate() error {
	mysqldatabaselog.Info("validate create", "name", in.Name)
	var allErrs field.ErrorList
	//if in.Spec.Database.Config.ReadOnly == 1 { //todo handle this case if possible
	//	allErrs = append(allErrs, field.Invalid(field.NewPath("spec.database.config"), in.Name, "Cannot create readOnly database"))
	//}
	if err := in.ValidateMySQLDatabase(); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath(""), in.Name, err.Error()))
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "schema.kubedb.com", Kind: "MySQLDatabase"}, in.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *MySQLDatabase) ValidateUpdate(old runtime.Object) error {
	mysqldatabaselog.Info("validate update", "name", in.Name)
	oldobj := old.(*MySQLDatabase)
	return ValidateMySQLDatabaseUpdate(in, oldobj)
}

func ValidateMySQLDatabaseUpdate(newobj *MySQLDatabase, oldobj *MySQLDatabase) error {
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
	er := newobj.ValidateMySQLDatabase()
	if er != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), newobj.Name, er.Error()))
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "schema.kubedb.com", Kind: "MySQLDatabase"}, newobj.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *MySQLDatabase) ValidateDelete() error {
	mysqldatabaselog.Info("validate delete", "name", in.Name)
	if in.Spec.DeletionPolicy == DeletionPolicyDoNotDelete {
		return field.Invalid(field.NewPath("spec").Child("terminationPolicy"), in.Name, `cannot delete object when terminationPolicy is set to "DoNotDelete"`)
	}
	if in.Spec.Database.Config.ReadOnly == 1 {
		return field.Invalid(field.NewPath("spec").Child("databaseConfig.readOnly"), in.Name, `schema manger cannot be deleted : database is read only enabled`)
	}
	return nil
}

func (in *MySQLDatabase) ValidateMySQLDatabase() error {
	var allErrs field.ErrorList
	if err := in.validateInitailizationSchema(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := in.validateMySQLDatabaseConfig(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := in.validateMySQLDatabaseNamespace(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := in.validateMySQLDatabaseName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "schema.kubedb.com", Kind: "MySQLDatabase"}, in.Name, allErrs)
}

func (in *MySQLDatabase) validateInitailizationSchema() *field.Error {
	path := field.NewPath("spec.init")
	if in.Spec.Init != nil {
		if in.Spec.Init.Script != nil && in.Spec.Init.Snapshot != nil {
			return field.Invalid(path, in.Name, `cannot initialize database using both restore and initSpec`)
		}
	}
	return nil
}

func (in *MySQLDatabase) validateMySQLDatabaseConfig() *field.Error {
	path := field.NewPath("spec").Child("database.config").Child("name")
	name := in.Spec.Database.Config.Name
	if name == "sys" {
		return field.Invalid(path, in.Name, `cannot use "sys" as the database name`)
	}
	if name == "performance_schema" {
		return field.Invalid(path, in.Name, `cannot use "performance_schema" as the database name`)
	}
	if name == "mysql" {
		return field.Invalid(path, in.Name, `cannot use "mysql" as the database name`)
	}
	if name == "kubedb_system" {
		return field.Invalid(path, in.Name, `cannot use "kubedb_system" as the database name`)
	}
	if name == "information_schema" {
		return field.Invalid(path, in.Name, `cannot use "information_schema" as the database name`)
	}
	if name == "admin" {
		return field.Invalid(path, in.Name, `cannot use "admin" as the database name`)
	}
	if name == "config" {
		return field.Invalid(path, in.Name, `cannot use "config" as the database name`)
	}
	path = field.NewPath("spec").Child("database.config")
	val := in.Spec.Database.Config.ReadOnly
	if val == 1 {
		if in.Spec.Init != nil {
			if (in.Spec.Init.Script != nil || in.Spec.Init.Snapshot != nil) && in.Status.Phase != DatabaseSchemaPhaseCurrent {
				return field.Invalid(path.Child("readOnly"), in.Name, `cannot make the database readonly , init/restore yet to be applied`)
			}
		}
	} else if in.Spec.Database.Config.Encryption == MySQLEncryptionEnabled {
		if in.Spec.Init != nil {
			if (in.Spec.Init.Script != nil || in.Spec.Init.Snapshot != nil) && in.Status.Phase != DatabaseSchemaPhaseCurrent {
				return field.Invalid(path.Child("encryption"), in.Name, `cannot make the database encryption enables , init/restore yet to be applied`)
			}
		}
	}
	return nil
}

func (in *MySQLDatabase) validateMySQLDatabaseNamespace() *field.Error {
	path := field.NewPath("metadata").Child("namespace")
	ns := in.ObjectMeta.Namespace
	if ns == "cert-manager" {
		return field.Invalid(path, in.Name, `cannot use namespace "cert-manager" to create schema manager`)
	}
	if ns == "kube-system" {
		return field.Invalid(path, in.Name, `cannot use namespace "kube-system" to create schema manager`)
	}
	if ns == "kubedb-system" {
		return field.Invalid(path, in.Name, `cannot use namespace "kubedb-system" to create schema manager`)
	}
	if ns == "kubedb" {
		return field.Invalid(path, in.Name, `cannot use namespace "kubedb" to create schema manager`)
	}
	if ns == "kubevault" {
		return field.Invalid(path, in.Name, `cannot use namespace "kubevault" to create schema manager`)
	}
	if ns == "local-path-storage" {
		return field.Invalid(path, in.Name, `cannot use namespace "local-path-storage" to create schema manager`)
	}
	return nil
}

func (in *MySQLDatabase) validateMySQLDatabaseName() *field.Error {
	if len(in.ObjectMeta.Name) > 45 {
		return field.Invalid(field.NewPath("metadata").Child("name"), in.Name, "must be no more than 30 characters")
	}
	return nil
}
