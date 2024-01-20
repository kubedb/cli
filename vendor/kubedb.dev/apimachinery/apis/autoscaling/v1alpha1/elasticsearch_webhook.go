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
	"context"
	"errors"
	"fmt"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var esLog = logf.Log.WithName("elasticsearch-autoscaler")

func (in *ElasticsearchAutoscaler) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-autoscaling-kubedb-com-v1alpha1-elasticsearchautoscaler,mutating=true,failurePolicy=fail,sideEffects=None,groups=autoscaling.kubedb.com,resources=elasticsearchautoscaler,verbs=create;update,versions=v1alpha1,name=melasticsearchautoscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &ElasticsearchAutoscaler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *ElasticsearchAutoscaler) Default() {
	esLog.Info("defaulting", "name", in.Name)
	in.setDefaults()
}

func (in *ElasticsearchAutoscaler) setDefaults() {
	in.setOpsReqOptsDefaults()

	if in.Spec.Storage != nil {
		setDefaultStorageValues(in.Spec.Storage.Node)
		setDefaultStorageValues(in.Spec.Storage.Master)
		setDefaultStorageValues(in.Spec.Storage.Data)
		setDefaultStorageValues(in.Spec.Storage.Ingest)
		setDefaultStorageValues(in.Spec.Storage.DataContent)
		setDefaultStorageValues(in.Spec.Storage.DataCold)
		setDefaultStorageValues(in.Spec.Storage.DataWarm)
		setDefaultStorageValues(in.Spec.Storage.DataFrozen)
		setDefaultStorageValues(in.Spec.Storage.DataHot)
		setDefaultStorageValues(in.Spec.Storage.ML)
		setDefaultStorageValues(in.Spec.Storage.Transform)
		setDefaultStorageValues(in.Spec.Storage.Coordinating)
	}
	if in.Spec.Compute != nil {
		setDefaultComputeValues(in.Spec.Compute.Node)
		setDefaultComputeValues(in.Spec.Compute.Master)
		setDefaultComputeValues(in.Spec.Compute.Data)
		setDefaultComputeValues(in.Spec.Compute.Ingest)
		setDefaultComputeValues(in.Spec.Compute.DataContent)
		setDefaultComputeValues(in.Spec.Compute.DataCold)
		setDefaultComputeValues(in.Spec.Compute.DataWarm)
		setDefaultComputeValues(in.Spec.Compute.DataFrozen)
		setDefaultComputeValues(in.Spec.Compute.DataHot)
		setDefaultComputeValues(in.Spec.Compute.ML)
		setDefaultComputeValues(in.Spec.Compute.Transform)
		setDefaultComputeValues(in.Spec.Compute.Coordinating)
	}
}

func (in *ElasticsearchAutoscaler) setOpsReqOptsDefaults() {
	if in.Spec.OpsRequestOptions == nil {
		in.Spec.OpsRequestOptions = &ElasticsearchOpsRequestOptions{}
	}
	// Timeout is defaulted to 600s in ops-manager retries.go (to retry 120 times with 5sec pause between each)
	// OplogMaxLagSeconds & ObjectsCountDiffPercentage are defaults to 0
	if in.Spec.OpsRequestOptions.Apply == "" {
		in.Spec.OpsRequestOptions.Apply = opsapi.ApplyOptionIfReady
	}
}

// +kubebuilder:webhook:path=/validate-schema-kubedb-com-v1alpha1-elasticsearchautoscaler,mutating=false,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=elasticsearchautoscalers,verbs=create;update;delete,versions=v1alpha1,name=velasticsearchautoscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &ElasticsearchAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *ElasticsearchAutoscaler) ValidateCreate() (admission.Warnings, error) {
	esLog.Info("validate create", "name", in.Name)
	return nil, in.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *ElasticsearchAutoscaler) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return nil, in.validate()
}

func (_ ElasticsearchAutoscaler) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (in *ElasticsearchAutoscaler) validate() error {
	if in.Spec.DatabaseRef == nil {
		return errors.New("databaseRef can't be empty")
	}

	var es dbapi.Elasticsearch
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name:      in.Spec.DatabaseRef.Name,
		Namespace: in.Namespace,
	}, &es)
	if err != nil {
		_ = fmt.Errorf("can't get Elasticsearch %s/%s \n", in.Namespace, in.Spec.DatabaseRef.Name)
		return err
	}

	if in.Spec.Compute != nil {
		cm := in.Spec.Compute
		if es.Spec.Topology != nil {
			if cm.Node != nil {
				return errors.New("Spec.Compute.PodResources is invalid for elastic-search topology")
			}
		} else {
			if cm.Master != nil || cm.Data != nil || cm.Ingest != nil || cm.DataContent != nil || cm.DataCold != nil || cm.DataFrozen != nil ||
				cm.DataWarm != nil || cm.DataHot != nil || cm.ML != nil || cm.Transform != nil || cm.Coordinating != nil {
				return errors.New("only Spec.Compute.Node is valid for basic elastic search structure")
			}
		}
	}

	if in.Spec.Storage != nil {
		st := in.Spec.Storage
		if es.Spec.Topology != nil {
			if st.Node != nil {
				return errors.New("Spec.Storage.PodResources is invalid for elastic-search topology")
			}
		} else {
			if st.Master != nil || st.Data != nil || st.Ingest != nil || st.DataContent != nil || st.DataCold != nil || st.DataFrozen != nil ||
				st.DataWarm != nil || st.DataHot != nil || st.ML != nil || st.Transform != nil || st.Coordinating != nil {
				return errors.New("only Spec.Storage.Node is valid for basic elastic search structure")
			}
		}
	}
	return nil
}
