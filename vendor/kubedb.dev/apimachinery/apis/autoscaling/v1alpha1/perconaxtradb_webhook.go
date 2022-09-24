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
	"errors"

	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var pxLog = logf.Log.WithName("perconaxtradb-autoscaler")

func (in *PerconaXtraDBAutoscaler) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-autoscaling-kubedb-com-v1alpha1-perconaxtradbautoscaler,mutating=true,failurePolicy=fail,sideEffects=None,groups=autoscaling.kubedb.com,resources=perconaxtradbautoscaler,verbs=create;update,versions=v1alpha1,name=mperconaxtradbautoscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &PerconaXtraDBAutoscaler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *PerconaXtraDBAutoscaler) Default() {
	pxLog.Info("defaulting", "name", in.Name)
	in.setDefaults()
}

func (in *PerconaXtraDBAutoscaler) setDefaults() {
	in.setOpsReqOptsDefaults()

	if in.Spec.Storage != nil {
		setDefaultStorageValues(in.Spec.Storage.PerconaXtraDB)
	}
	if in.Spec.Compute != nil {
		setDefaultComputeValues(in.Spec.Compute.PerconaXtraDB)
	}
}

func (in *PerconaXtraDBAutoscaler) setOpsReqOptsDefaults() {
	if in.Spec.OpsRequestOptions == nil {
		in.Spec.OpsRequestOptions = &PerconaXtraDBOpsRequestOptions{}
	}
	// Timeout is defaulted to 600s in ops-manager retries.go (to retry 120 times with 5sec pause between each)
	// OplogMaxLagSeconds & ObjectsCountDiffPercentage are defaults to 0
	if in.Spec.OpsRequestOptions.Apply == "" {
		in.Spec.OpsRequestOptions.Apply = opsapi.ApplyOptionIfReady
	}
}

// +kubebuilder:webhook:path=/validate-schema-kubedb-com-v1alpha1-perconaxtradbautoscaler,mutating=false,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=perconaxtradbautoscalers,verbs=create;update;delete,versions=v1alpha1,name=vperconaxtradbautoscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &PerconaXtraDBAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *PerconaXtraDBAutoscaler) ValidateCreate() error {
	pxLog.Info("validate create", "name", in.Name)
	return in.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *PerconaXtraDBAutoscaler) ValidateUpdate(old runtime.Object) error {
	return in.validate()
}

func (_ PerconaXtraDBAutoscaler) ValidateDelete() error {
	return nil
}

func (in *PerconaXtraDBAutoscaler) validate() error {
	if in.Spec.DatabaseRef == nil {
		return errors.New("databaseRef can't be empty")
	}
	return nil
}
