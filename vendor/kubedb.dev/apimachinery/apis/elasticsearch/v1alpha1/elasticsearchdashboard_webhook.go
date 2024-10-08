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
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	amv "kubedb.dev/apimachinery/pkg/validator"

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	coreutil "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/policy/secomp"
	ofst "kmodules.xyz/offshoot-api/api/v2"
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

// +kubebuilder:webhook:path=/mutate-elasticsearch-kubedb-com-v1alpha1-elasticsearchelasticsearch,mutating=true,failurePolicy=fail,sideEffects=None,groups=elasticsearch.kubedb.com,resources=elasticsearchelasticsearchs,verbs=create;update,versions=v1alpha1,name=melasticsearchelasticsearch.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &ElasticsearchDashboard{}

func (ed *ElasticsearchDashboard) setDefaultContainerSecurityContext(esVersion catalog.ElasticsearchVersion, podTemplate *ofst.PodTemplateSpec) {
	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.ElasticsearchInitConfigMergerContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.ElasticsearchInitConfigMergerContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	ed.assignDefaultContainerSecurityContext(esVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)
	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.ElasticsearchContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	ed.assignDefaultContainerSecurityContext(esVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (ed *ElasticsearchDashboard) setDefaultContainerResourceLimits(podTemplate *ofst.PodTemplateSpec) {
	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
	if container != nil && (container.Resources.Requests == nil && container.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&container.Resources, kubedb.DefaultResources)
	}

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.ElasticsearchInitConfigMergerContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}
}

func (ed *ElasticsearchDashboard) assignDefaultContainerSecurityContext(esVersion catalog.ElasticsearchVersion, sc *core.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(esVersion.Spec.SecurityContext.RunAsAnyNonRoot)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = esVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (ed *ElasticsearchDashboard) Default() {
	ed.SetDefaults()
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

	if ed.Spec.DeletionPolicy == dbapi.DeletionPolicyDoNotTerminate {
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
	container := coreutil.GetContainerByName(ed.Spec.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
	if err := amv.ValidateEnvVar(container.Env, forbiddenEnvVars, ResourceKindElasticsearchDashboard); err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podtemplate").Child("spec").Child("containers").Child("env"), ed.Name,
			"Invalid spec.podtemplate.spec.containers[i].env , avoid using the forbidden env variables"))
	}

	initContainer := coreutil.GetContainerByName(ed.Spec.PodTemplate.Spec.InitContainers, kubedb.ElasticsearchInitConfigMergerContainerName)
	if err := amv.ValidateEnvVar(initContainer.Env, forbiddenEnvVars, ResourceKindElasticsearchDashboard); err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podtemplate").Child("spec").Child("initContainers").Child("env"), ed.Name,
			"Invalid spec.podtemplate.spec.initContainers[i].env , avoid using the forbidden env variables"))
	}

	if len(allErr) == 0 {
		return nil
	}

	return apierrors.NewInvalid(schema.GroupKind{Group: "elasticsearch.kubedb.com", Kind: "ElasticsearchDashboard"}, ed.Name, allErr)
}
