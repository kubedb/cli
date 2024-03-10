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
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var kafkaLog = logf.Log.WithName("kafka-autoscaler")

var _ webhook.Defaulter = &KafkaAutoscaler{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type
func (k *KafkaAutoscaler) Default() {
	kafkaLog.Info("defaulting", "name", k.Name)
	k.setDefaults()
}

func (k *KafkaAutoscaler) setDefaults() {
	var db dbapi.Kafka
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name:      k.Spec.DatabaseRef.Name,
		Namespace: k.Namespace,
	}, &db)
	if err != nil {
		_ = fmt.Errorf("can't get Kafka %s/%s \n", k.Namespace, k.Spec.DatabaseRef.Name)
		return
	}

	k.setOpsReqOptsDefaults()

	if k.Spec.Storage != nil {
		if db.Spec.Topology != nil {
			setDefaultStorageValues(k.Spec.Storage.Broker)
			setDefaultStorageValues(k.Spec.Storage.Controller)
		} else {
			setDefaultStorageValues(k.Spec.Storage.Node)
		}
	}

	if k.Spec.Compute != nil {
		if db.Spec.Topology != nil {
			setDefaultComputeValues(k.Spec.Compute.Broker)
			setDefaultComputeValues(k.Spec.Compute.Controller)
		} else {
			setDefaultComputeValues(k.Spec.Compute.Node)
		}
	}
}

func (k *KafkaAutoscaler) setOpsReqOptsDefaults() {
	if k.Spec.OpsRequestOptions == nil {
		k.Spec.OpsRequestOptions = &KafkaOpsRequestOptions{}
	}
	// Timeout is defaulted to 600s in ops-manager retries.go (to retry 120 times with 5sec pause between each)
	// OplogMaxLagSeconds & ObjectsCountDiffPercentage are defaults to 0
	if k.Spec.OpsRequestOptions.Apply == "" {
		k.Spec.OpsRequestOptions.Apply = opsapi.ApplyOptionIfReady
	}
}

var _ webhook.Validator = &KafkaAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (k *KafkaAutoscaler) ValidateCreate() (admission.Warnings, error) {
	kafkaLog.Info("validate create", "name", k.Name)
	return nil, k.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (k *KafkaAutoscaler) ValidateUpdate(oldObj runtime.Object) (admission.Warnings, error) {
	kafkaLog.Info("validate create", "name", k.Name)
	return nil, k.validate()
}

func (_ *KafkaAutoscaler) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (k *KafkaAutoscaler) validate() error {
	if k.Spec.DatabaseRef == nil {
		return errors.New("databaseRef can't be empty")
	}
	var kf dbapi.Kafka
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name:      k.Spec.DatabaseRef.Name,
		Namespace: k.Namespace,
	}, &kf)
	if err != nil {
		_ = fmt.Errorf("can't get Kafka %s/%s \n", k.Namespace, k.Spec.DatabaseRef.Name)
		return err
	}

	if k.Spec.Compute != nil {
		cm := k.Spec.Compute
		if kf.Spec.Topology != nil {
			if cm.Node != nil {
				return errors.New("Spec.Compute.Node is invalid for kafka with topology")
			}
		} else {
			if cm.Broker != nil {
				return errors.New("Spec.Compute.Broker is invalid for combined kafka")
			}
			if cm.Controller != nil {
				return errors.New("Spec.Compute.Controller is invalid for combined kafka")
			}
		}
	}

	if k.Spec.Storage != nil {
		st := k.Spec.Storage
		if kf.Spec.Topology != nil {
			if st.Node != nil {
				return errors.New("Spec.Storage.Node is invalid for kafka with topology")
			}
		} else {
			if st.Broker != nil {
				return errors.New("Spec.Storage.Broker is invalid for combined kafka")
			}
			if st.Controller != nil {
				return errors.New("Spec.Storage.Controller is invalid for combined kafka")
			}
		}
	}

	return nil
}
