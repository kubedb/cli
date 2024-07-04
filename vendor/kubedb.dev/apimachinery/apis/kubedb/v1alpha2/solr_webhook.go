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
	"strings"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var solrlog = logf.Log.WithName("solr-resource")

var _ webhook.Defaulter = &Solr{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (s *Solr) Default() {
	if s == nil {
		return
	}
	solrlog.Info("default", "name", s.Name)

	slVersion := catalog.SolrVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: s.Spec.Version}, &slVersion)
	if err != nil {
		klog.Errorf("Version does not exist.")
		return
	}

	s.SetDefaults(&slVersion)
}

var _ webhook.Validator = &Solr{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (s *Solr) ValidateCreate() (admission.Warnings, error) {
	solrlog.Info("validate create", "name", s.Name)

	allErr := s.ValidateCreateOrUpdate()
	if len(allErr) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Solr"}, s.Name, allErr)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (s *Solr) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	solrlog.Info("validate update", "name", s.Name)

	_ = old.(*Solr)
	allErr := s.ValidateCreateOrUpdate()

	if len(allErr) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Solr"}, s.Name, allErr)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (s *Solr) ValidateDelete() (admission.Warnings, error) {
	solrlog.Info("validate delete", "name", s.Name)

	var allErr field.ErrorList
	if s.Spec.DeletionPolicy == TerminationPolicyDoNotTerminate {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("terminationPolicy"),
			s.Name,
			"Can not delete as terminationPolicy is set to \"DoNotTerminate\""))
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "kubedb.com", Kind: "Solr"}, s.Name, allErr)
	}
	return nil, nil
}

var solrReservedVolumes = []string{
	kubedb.SolrVolumeConfig,
	kubedb.SolrVolumeDefaultConfig,
	kubedb.SolrVolumeCustomConfig,
	kubedb.SolrVolumeAuthConfig,
}

var solrReservedVolumeMountPaths = []string{
	kubedb.SolrHomeDir,
	kubedb.SolrDataDir,
	kubedb.SolrCustomConfigDir,
	kubedb.SolrSecurityConfigDir,
	kubedb.SolrTempConfigDir,
}

var solrAvailableModules = []string{
	"analysis-extras", "extraction", "hdfs", "langid", "prometheus-exporter", "sql",
	"analytics", "gcs-repository", "jaegertracer-configurator", "ltr", "s3-repository",
	"clustering", "hadoop-auth", "jwt-auth", "opentelemetry", "scripting",
}

func (s *Solr) ValidateCreateOrUpdate() field.ErrorList {
	var allErr field.ErrorList

	if s.Spec.EnableSSL {
		if s.Spec.TLS == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("enableSSL"),
				s.Name,
				".spec.tls can't be nil, if .spec.enableSSL is true"))
		}
	} else {
		if s.Spec.TLS != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("enableSSL"),
				s.Name,
				".spec.tls must be nil, if .spec.enableSSL is disabled"))
		}
	}

	if s.Spec.Version == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
			s.Name,
			"spec.version' is missing"))
	} else {
		err := solrValidateVersion(s)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("version"),
				s.Name,
				err.Error()))
		}
	}

	err := solrValidateModules(s)
	if err != nil {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("solrmodules"),
			s.Name,
			err.Error()))
	}

	if s.Spec.Topology == nil {
		if s.Spec.Replicas != nil && *s.Spec.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("replicas"),
				s.Name,
				"number of replicas can not be less be 0 or less"))
		}
		err := solrValidateVolumes(&s.Spec.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("volumes"),
				s.Name,
				err.Error()))
		}
		err = solrValidateVolumesMountPaths(&s.Spec.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("podTemplate").Child("spec").Child("containers"),
				s.Name,
				err.Error()))
		}

	} else {
		if s.Spec.Topology.Data == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("data"),
				s.Name,
				".spec.topology.data can't be empty in cluster mode"))
		}
		if s.Spec.Topology.Data.Replicas != nil && *s.Spec.Topology.Data.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("data").Child("replicas"),
				s.Name,
				"number of replicas can not be less be 0 or less"))
		}
		err := solrValidateVolumes(&s.Spec.Topology.Data.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("data").Child("podTemplate").Child("spec").Child("volumes"),
				s.Name,
				err.Error()))
		}
		err = solrValidateVolumesMountPaths(&s.Spec.Topology.Data.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("data").Child("podTemplate").Child("spec").Child("containers"),
				s.Name,
				err.Error()))
		}

		if s.Spec.Topology.Overseer == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("overseer"),
				s.Name,
				".spec.topology.overseer can't be empty in cluster mode"))
		}
		if s.Spec.Topology.Overseer.Replicas != nil && *s.Spec.Topology.Overseer.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("overseer").Child("replicas"),
				s.Name,
				"number of replicas can not be less be 0 or less"))
		}
		err = solrValidateVolumes(&s.Spec.Topology.Overseer.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("overseer").Child("podTemplate").Child("spec").Child("volumes"),
				s.Name,
				err.Error()))
		}
		err = solrValidateVolumesMountPaths(&s.Spec.Topology.Overseer.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("overseer").Child("podTemplate").Child("spec").Child("containers"),
				s.Name,
				err.Error()))
		}

		if s.Spec.Topology.Coordinator == nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("coordinator"),
				s.Name,
				".spec.topology.coordinator can't be empty in cluster mode"))
		}
		if s.Spec.Topology.Coordinator.Replicas != nil && *s.Spec.Topology.Coordinator.Replicas <= 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("coordinator").Child("replicas"),
				s.Name,
				"number of replicas can not be less be 0 or less"))
		}
		err = solrValidateVolumes(&s.Spec.Topology.Coordinator.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("coordinator").Child("podTemplate").Child("spec").Child("volumes"),
				s.Name,
				err.Error()))
		}
		err = solrValidateVolumesMountPaths(&s.Spec.Topology.Coordinator.PodTemplate)
		if err != nil {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("topology").Child("coordinator").Child("podTemplate").Child("spec").Child("containers"),
				s.Name,
				err.Error()))
		}
	}

	if s.Spec.StorageType == "" {
		allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
			s.Name,
			"StorageType can not be empty"))
	} else {
		if s.Spec.StorageType != StorageTypeDurable && s.Spec.StorageType != StorageTypeEphemeral {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("storageType"),
				s.Name,
				"StorageType should be either durable or ephemeral"))
		}
	}

	for _, x := range s.Spec.SolrOpts {
		if strings.Count(x, " ") > 0 {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("solropts"),
				s.Name,
				"solropt jvm env variables must not contain space"))
		}
		if x[0] != '-' || x[1] != 'D' {
			allErr = append(allErr, field.Invalid(field.NewPath("spec").Child("solropts"),
				s.Name,
				"solropt jvm env variables must start with -D"))
		}
	}

	if len(allErr) == 0 {
		return nil
	}
	return allErr
}

func solrValidateVersion(s *Solr) error {
	slVersion := catalog.SolrVersion{}
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{Name: s.Spec.Version}, &slVersion)
	if err != nil {
		return errors.New("version not supported")
	}
	return nil
}

func solrValidateModules(s *Solr) error {
	modules := s.Spec.SolrModules
	for _, mod := range modules {
		fl := false
		for _, av_mod := range solrAvailableModules {
			if mod == av_mod {
				fl = true
				break
			}
		}
		if !fl {
			return fmt.Errorf("%s does not exist in available modules", mod)
		}
	}
	return nil
}

func solrValidateVolumes(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Volumes == nil {
		return nil
	}

	for _, rv := range solrReservedVolumes {
		for _, ugv := range podTemplate.Spec.Volumes {
			if ugv.Name == rv {
				return errors.New("Can't use a reserve volume name: " + rv)
			}
		}
	}

	return nil
}

func solrValidateVolumesMountPaths(podTemplate *ofst.PodTemplateSpec) error {
	if podTemplate == nil {
		return nil
	}
	if podTemplate.Spec.Containers == nil {
		return nil
	}

	for _, rvmp := range solrReservedVolumeMountPaths {
		containerList := podTemplate.Spec.Containers
		for i := range containerList {
			mountPathList := containerList[i].VolumeMounts
			for j := range mountPathList {
				if mountPathList[j].MountPath == rvmp {
					return errors.New("Can't use a reserve volume mount path name: " + rvmp)
				}
			}
		}
	}

	if podTemplate.Spec.InitContainers == nil {
		return nil
	}

	for _, rvmp := range solrReservedVolumeMountPaths {
		containerList := podTemplate.Spec.InitContainers
		for i := range containerList {
			mountPathList := containerList[i].VolumeMounts
			for j := range mountPathList {
				if mountPathList[j].MountPath == rvmp {
					return errors.New("Can't use a reserve volume mount path name: " + rvmp)
				}
			}
		}
	}

	return nil
}
