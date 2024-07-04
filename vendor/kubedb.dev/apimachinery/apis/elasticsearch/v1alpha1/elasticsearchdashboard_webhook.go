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

	"kubedb.dev/apimachinery/apis"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	amv "kubedb.dev/apimachinery/pkg/validator"

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/policy/secomp"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var forbiddenEnvVars = []string{
	ES_USER_ENV,
	ES_PASSWORD_ENV,
	ES_USER_KEY,
	ES_PASSWORD_KEY,
	OS_USER_KEY,
	OS_PASSWORD_KEY,
	DashboardServerHostKey,
	DashboardServerNameKey,
	DashboardServerPortKey,
	DashboardServerSSLCaKey,
	DashboardServerSSLCertKey,
	DashboardServerSSLKey,
	DashboardServerSSLEnabledKey,
	ElasticsearchSSLCaKey,
	ElasticsearchHostsKey,
	OpensearchHostsKey,
	OpensearchSSLCaKey,
}

// log is for logging in this package.
var edLog = logf.Log.WithName("elasticsearchelasticsearch-validation")

func (ed *ElasticsearchDashboard) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(ed).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-elasticsearch-kubedb-com-v1alpha1-elasticsearchelasticsearch,mutating=true,failurePolicy=fail,sideEffects=None,groups=elasticsearch.kubedb.com,resources=elasticsearchelasticsearchs,verbs=create;update,versions=v1alpha1,name=melasticsearchelasticsearch.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &ElasticsearchDashboard{}

func (ed *ElasticsearchDashboard) setDefaultContainerSecurityContext(podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.ContainerSecurityContext == nil {
		podTemplate.Spec.ContainerSecurityContext = &core.SecurityContext{}
	}
	ed.assignDefaultContainerSecurityContext(podTemplate.Spec.ContainerSecurityContext)
}

func (ed *ElasticsearchDashboard) assignDefaultContainerSecurityContext(sc *core.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(true)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = pointer.Int64P(1000)
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (ed *ElasticsearchDashboard) Default() {
	if ed.Spec.Replicas == nil {
		ed.Spec.Replicas = pointer.Int32P(1)
		edLog.Info(".Spec.Replicas have been set to default")
	}

	apis.SetDefaultResourceLimits(&ed.Spec.PodTemplate.Spec.Resources, DashboardsDefaultResources)
	edLog.Info(".PodTemplate.Spec.Resources have been set to default")

	if len(ed.Spec.TerminationPolicy) == 0 {
		ed.Spec.TerminationPolicy = dbapi.DeletionPolicyWipeOut
		edLog.Info(".Spec.DeletionPolicy have been set to DeletionPolicyWipeOut")
	}

	ed.setDefaultContainerSecurityContext(&ed.Spec.PodTemplate)

	if ed.Spec.EnableSSL {
		if ed.Spec.TLS == nil {
			ed.Spec.TLS = &kmapi.TLSConfig{}
		}
		if ed.Spec.TLS.IssuerRef == nil {
			ed.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(ed.Spec.TLS.Certificates, kmapi.CertificateSpec{
				Alias:      string(ElasticsearchDashboardCACert),
				SecretName: ed.DefaultCertificateSecretName(ElasticsearchDashboardCACert),
			})
		}
		ed.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(ed.Spec.TLS.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchDashboardServerCert),
			SecretName: ed.DefaultCertificateSecretName(ElasticsearchDashboardServerCert),
		})
	}
}

// +kubebuilder:webhook:path=/validate-elasticsearch-kubedb-com-v1alpha1-elasticsearchelasticsearch,mutating=false,failurePolicy=fail,sideEffects=None,groups=elasticsearch.kubedb.com,resources=elasticsearchelasticsearchs,verbs=create;update;delete,versions=v1alpha1,name=velasticsearchelasticsearch.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &ElasticsearchDashboard{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (ed *ElasticsearchDashboard) ValidateCreate() (admission.Warnings, error) {
	edLog.Info("validate create", "name", ed.Name)
	return nil, ed.Validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (ed *ElasticsearchDashboard) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	// Skip validation, if UPDATE operation is called after deletion.
	// Case: Removing Finalizer
	if ed.DeletionTimestamp != nil {
		return nil, nil
	}
	return nil, ed.Validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (ed *ElasticsearchDashboard) ValidateDelete() (admission.Warnings, error) {
	edLog.Info("validate delete", "name", ed.Name)

	var allErr field.ErrorList

	if ed.Spec.TerminationPolicy == dbapi.DeletionPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationpolicy"), ed.Name,
			fmt.Sprintf("ElasticsearchDashboard %s/%s can't be deleted. Change .spec.terminationpolicy", ed.Namespace, ed.Name)))
	}

	if len(allErr) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "elasticsearch.kubedb.com", Kind: "ElasticsearchDashboard"},
		ed.Name, allErr)
}

func (ed *ElasticsearchDashboard) Validate() error {
	var allErr field.ErrorList

	// database ref is required
	if ed.Spec.DatabaseRef == nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("databaseref"), ed.Name,
			"spec.databaseref can't be empty"))
	}

	// validate if user provided replicas are non-negative
	// user may provide 0 replicas
	if *ed.Spec.Replicas < 0 {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"), ed.Name,
			fmt.Sprintf("spec.replicas %v invalid. Must be greater than zero", ed.Spec.Replicas)))
	}

	// env variables needs to be validated
	// so that variables provided in config secret
	// and credential env may not be overwritten
	if err := amv.ValidateEnvVar(ed.Spec.PodTemplate.Spec.Env, forbiddenEnvVars, ResourceKindElasticsearchDashboard); err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podtemplate").Child("spec").Child("env"), ed.Name,
			"Invalid spec.podtemplate.spec.env , avoid using the forbidden env variables"))
	}

	if len(allErr) == 0 {
		return nil
	}

	return apierrors.NewInvalid(schema.GroupKind{Group: "elasticsearch.kubedb.com", Kind: "ElasticsearchDashboard"}, ed.Name, allErr)
}
