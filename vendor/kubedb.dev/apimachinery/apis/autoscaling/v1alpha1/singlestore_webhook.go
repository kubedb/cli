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

	olddbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var singlestoreLog = logf.Log.WithName("singlestore-autoscaler")

var _ webhook.Defaulter = &SinglestoreAutoscaler{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type
func (s *SinglestoreAutoscaler) Default() {
	singlestoreLog.Info("defaulting", "name", s.Name)
	s.setDefaults()
}

func (s *SinglestoreAutoscaler) setDefaults() {
	var db olddbapi.Singlestore
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name:      s.Spec.DatabaseRef.Name,
		Namespace: s.Namespace,
	}, &db)
	if err != nil {
		_ = fmt.Errorf("can't get Singlestore %s/%s \n", s.Namespace, s.Spec.DatabaseRef.Name)
		return
	}

	s.setOpsReqOptsDefaults()

	if s.Spec.Storage != nil {
		if db.Spec.Topology != nil {
			setDefaultStorageValues(s.Spec.Storage.Aggregator)
			setDefaultStorageValues(s.Spec.Storage.Leaf)
		} else {
			setDefaultStorageValues(s.Spec.Storage.Node)
		}
	}

	if s.Spec.Compute != nil {
		if db.Spec.Topology != nil {
			setDefaultComputeValues(s.Spec.Compute.Aggregator)
			setDefaultComputeValues(s.Spec.Compute.Leaf)
		} else {
			setDefaultComputeValues(s.Spec.Compute.Node)
		}
	}
}

func (s *SinglestoreAutoscaler) setOpsReqOptsDefaults() {
	if s.Spec.OpsRequestOptions == nil {
		s.Spec.OpsRequestOptions = &SinglestoreOpsRequestOptions{}
	}
	// Timeout is defaulted to 600s in ops-manager retries.go (to retry 120 times with 5sec pause between each)
	// OplogMaxLagSeconds & ObjectsCountDiffPercentage are defaults to 0
	if s.Spec.OpsRequestOptions.Apply == "" {
		s.Spec.OpsRequestOptions.Apply = opsapi.ApplyOptionIfReady
	}
}

var _ webhook.Validator = &SinglestoreAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (s *SinglestoreAutoscaler) ValidateCreate() (admission.Warnings, error) {
	kafkaLog.Info("validate create", "name", s.Name)
	return nil, s.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (s *SinglestoreAutoscaler) ValidateUpdate(oldObj runtime.Object) (admission.Warnings, error) {
	kafkaLog.Info("validate update", "name", s.Name)
	return nil, s.validate()
}

func (_ *SinglestoreAutoscaler) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (s *SinglestoreAutoscaler) validate() error {
	if s.Spec.DatabaseRef == nil {
		return errors.New("databaseRef can't be empty")
	}
	var sdb olddbapi.Singlestore
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name:      s.Spec.DatabaseRef.Name,
		Namespace: s.Namespace,
	}, &sdb)
	if err != nil {
		_ = fmt.Errorf("can't get Singlestore %s/%s \n", s.Namespace, s.Spec.DatabaseRef.Name)
		return err
	}

	if s.Spec.Compute != nil {
		cm := s.Spec.Compute
		if sdb.Spec.Topology != nil {
			if cm.Node != nil {
				return errors.New("Spec.Compute.Node is invalid for singlestore with cluster")
			}
		} else {
			if cm.Aggregator != nil {
				return errors.New("Spec.Compute.Aggregator is invalid for standalone")
			}
			if cm.Leaf != nil {
				return errors.New("Spec.Compute.Leaf is invalid for combined standalone")
			}
		}
	}

	if s.Spec.Storage != nil {
		st := s.Spec.Storage
		if sdb.Spec.Topology != nil {
			if st.Node != nil {
				return errors.New("Spec.Storage.Node is invalid for Singlestore with cluster")
			}
		} else {
			if st.Aggregator != nil {
				return errors.New("Spec.Storage.Aggregator is invalid for standalone")
			}
			if st.Leaf != nil {
				return errors.New("Spec.Storage.Leaf is invalid for standalone")
			}
		}
	}

	return nil
}
