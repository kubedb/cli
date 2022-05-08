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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var mongodbdatabaselog = logf.Log.WithName("mongodbdatabase-resource")

func (in *MongoDBDatabase) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-schema-kubedb-com-v1alpha1-mongodbdatabase,mutating=true,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=mongodbdatabases,verbs=create;update,versions=v1alpha1,name=mmongodbdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &MongoDBDatabase{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *MongoDBDatabase) Default() {
}

// +kubebuilder:webhook:path=/validate-schema-kubedb-com-v1alpha1-mongodbdatabase,mutating=false,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=mongodbdatabases,verbs=create;update;delete,versions=v1alpha1,name=vmongodbdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MongoDBDatabase{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *MongoDBDatabase) ValidateCreate() error {
	mongodbdatabaselog.Info("validate create", "name", in.Name)
	return in.ValidateMongoDBDatabase()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *MongoDBDatabase) ValidateUpdate(old runtime.Object) error {
	mongodbdatabaselog.Info("validate update", "name", in.Name)
	var allErrs field.ErrorList
	path := field.NewPath("spec")
	oldDb := old.(*MongoDBDatabase)

	// if phase is 'Current', do not give permission to change the DatabaseConfig.Name
	if oldDb.Status.Phase == DatabaseSchemaPhaseCurrent && oldDb.Spec.Database.Config.Name != in.Spec.Database.Config.Name {
		allErrs = append(allErrs, field.Invalid(path.Child("database").Child("config"), in.Name, MongoDBValidateDatabaseNameChangeError))
		return apierrors.NewInvalid(in.GroupVersionKind().GroupKind(), in.Name, allErrs)
	}

	// If Initialized==true, Do not give permission to unset it
	if oldDb.Spec.Init != nil && oldDb.Spec.Init.Initialized { // initialized is already set in old object
		// If user updated the Schema-yaml with no Spec.Init
		// Or
		// user updated the Schema-yaml with Spec.Init.Initialized = false
		if in.Spec.Init == nil || (in.Spec.Init != nil && !in.Spec.Init.Initialized) {
			allErrs = append(allErrs, field.Invalid(path.Child("init").Child("initialized"), in.Name, MongoDBValidateInitializedUnsetError))
			return apierrors.NewInvalid(in.GroupVersionKind().GroupKind(), in.Name, allErrs)
		}
	}

	// making VaultRef & DatabaseRef fields immutable
	if oldDb.Spec.Database.ServerRef != in.Spec.Database.ServerRef {
		allErrs = append(allErrs, field.Invalid(path.Child("database").Child("serverRef"), in.Name, MongoDBValidateDBServerRefChangeError))
	}
	if oldDb.Spec.VaultRef != in.Spec.VaultRef {
		allErrs = append(allErrs, field.Invalid(path.Child("vaultRef"), in.Name, MongoDBValidateVaultRefChangeError))
	}
	if len(allErrs) > 0 {
		return apierrors.NewInvalid(in.GroupVersionKind().GroupKind(), in.Name, allErrs)
	}
	return in.ValidateMongoDBDatabase()
}

const (
	// these constants are here for purpose of Testing of MongoDB Operator code

	MongoDBValidateDeletionPolicyError     = "schema can't be deleted if the deletion policy is DoNotDelete"
	MongoDBValidateInitTypeBothError       = "cannot initialize database using both restore and initSpec"
	MongoDBValidateInitializedUnsetError   = "cannot unset the initialized field directly"
	MongoDBValidateDatabaseNameChangeError = "you can't change the Database Config name now"
	MongoDBValidateDBServerRefChangeError  = "cannot change mongodb reference"
	MongoDBValidateVaultRefChangeError     = "cannot change vault reference"
)

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *MongoDBDatabase) ValidateDelete() error {
	mongodbdatabaselog.Info("validate delete", "name", in.Name)
	if in.Spec.DeletionPolicy == DeletionPolicyDoNotDelete {
		var allErrs field.ErrorList
		path := field.NewPath("spec").Child("deletionPolicy")
		allErrs = append(allErrs, field.Invalid(path, in.Name, MongoDBValidateDeletionPolicyError))
		return apierrors.NewInvalid(in.GroupVersionKind().GroupKind(), in.Name, allErrs)
	}
	return nil
}

func (in *MongoDBDatabase) ValidateMongoDBDatabase() error {
	var allErrs field.ErrorList
	if err := in.validateSchemaInitRestore(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := in.validateMongoDBDatabaseSchemaName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := in.CheckIfNameFieldsAreOkOrNot(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(in.GroupVersionKind().GroupKind(), in.Name, allErrs)
}

func (in *MongoDBDatabase) validateSchemaInitRestore() *field.Error {
	path := field.NewPath("spec").Child("init")
	if in.Spec.Init != nil && in.Spec.Init.Script != nil && in.Spec.Init.Snapshot != nil {
		return field.Invalid(path, in.Name, MongoDBValidateInitTypeBothError)
	}
	return nil
}

func (in *MongoDBDatabase) validateMongoDBDatabaseSchemaName() *field.Error {
	path := field.NewPath("spec").Child("database").Child("config").Child("name")
	name := in.Spec.Database.Config.Name

	if name == MongoDatabaseNameForEntry || name == DatabaseNameAdmin || name == DatabaseNameConfig || name == DatabaseNameLocal {
		str := fmt.Sprintf("cannot use \"%v\" as the database name", name)
		return field.Invalid(path, in.Name, str)
	}
	return nil
}

/*
Ensure that the name of database, vault & repository are not empty
*/

func (in *MongoDBDatabase) CheckIfNameFieldsAreOkOrNot() *field.Error {
	if in.Spec.Database.ServerRef.Name == "" {
		str := "Database Ref name cant be empty"
		return field.Invalid(field.NewPath("spec").Child("database").Child("serverRef").Child("name"), in.Name, str)
	}
	if in.Spec.VaultRef.Name == "" {
		str := "Vault Ref name cant be empty"
		return field.Invalid(field.NewPath("spec").Child("vaultRef").Child("name"), in.Name, str)
	}
	if in.Spec.Init != nil && in.Spec.Init.Snapshot != nil && in.Spec.Init.Snapshot.Repository.Name == "" {
		str := "Repository name cant be empty"
		return field.Invalid(field.NewPath("spec").Child("init").Child("snapshot").Child("repository").Child("name"), in.Name, str)
	}
	return nil
}
