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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var connectorlog = logf.Log.WithName("connector-resource")

var _ webhook.Defaulter = &Connector{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (k *Connector) Default() {
	if k == nil {
		return
	}
	connectClusterLog.Info("default", "name", k.Name)
	if k.Spec.DeletionPolicy == "" {
		k.Spec.DeletionPolicy = dbapi.DeletionPolicyDelete
	}
}

var _ webhook.Validator = &Connector{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (k *Connector) ValidateCreate() (admission.Warnings, error) {
	connectClusterLog.Info("validate create", "name", k.Name)
	return nil, k.ValidateCreateOrUpdate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (k *Connector) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	connectClusterLog.Info("validate update", "name", k.Name)
	return nil, k.ValidateCreateOrUpdate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (k *Connector) ValidateDelete() (admission.Warnings, error) {
	connectorlog.Info("validate delete", "name", k.Name)

	var allErr field.ErrorList
	if k.Spec.DeletionPolicy == dbapi.DeletionPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("deletionPolicy"),
			k.Name,
			"Can not delete as deletionPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "Connector"}, k.Name, allErr)
	}
	return nil, nil
}

func (k *Connector) ValidateCreateOrUpdate() error {
	var allErr field.ErrorList
	if k.Spec.DeletionPolicy == dbapi.DeletionPolicyHalt {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("deletionPolicy"),
			k.Name,
			"DeletionPolicyHalt isn't supported for Connector"))
	}
	if len(allErr) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "kafka.kubedb.com", Kind: "ConnectCluster"}, k.Name, allErr)
}
