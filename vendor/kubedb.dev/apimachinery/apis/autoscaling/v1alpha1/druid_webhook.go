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
var druidLog = logf.Log.WithName("druid-autoscaler")

var _ webhook.Defaulter = &DruidAutoscaler{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type
func (d *DruidAutoscaler) Default() {
	druidLog.Info("defaulting", "name", d.Name)
	d.setDefaults()
}

func (d *DruidAutoscaler) setDefaults() {
	var db olddbapi.Druid
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name:      d.Spec.DatabaseRef.Name,
		Namespace: d.Namespace,
	}, &db)
	if err != nil {
		_ = fmt.Errorf("can't get Druid %s/%s \n", d.Namespace, d.Spec.DatabaseRef.Name)
		return
	}

	d.setOpsReqOptsDefaults()

	if d.Spec.Storage != nil {
		if db.Spec.Topology != nil {
			if db.Spec.Topology.MiddleManagers != nil && d.Spec.Storage.MiddleManagers != nil {
				setDefaultStorageValues(d.Spec.Storage.MiddleManagers)
			}
			if db.Spec.Topology.Historicals != nil && d.Spec.Storage.Historicals != nil {
				setDefaultStorageValues(d.Spec.Storage.Historicals)
			}
		}
	}

	if d.Spec.Compute != nil {
		if db.Spec.Topology != nil {
			if db.Spec.Topology.Coordinators != nil && d.Spec.Compute.Coordinators != nil {
				setDefaultComputeValues(d.Spec.Compute.Coordinators)
			}
			if db.Spec.Topology.Overlords != nil && d.Spec.Compute.Overlords != nil {
				setDefaultComputeValues(d.Spec.Compute.Overlords)
			}
			if db.Spec.Topology.MiddleManagers != nil && d.Spec.Compute.MiddleManagers != nil {
				setDefaultComputeValues(d.Spec.Compute.MiddleManagers)
			}
			if db.Spec.Topology.Historicals != nil && d.Spec.Compute.Historicals != nil {
				setDefaultComputeValues(d.Spec.Compute.Historicals)
			}
			if db.Spec.Topology.Brokers != nil && d.Spec.Compute.Brokers != nil {
				setDefaultComputeValues(d.Spec.Compute.Brokers)
			}
			if db.Spec.Topology.Routers != nil && d.Spec.Compute.Routers != nil {
				setDefaultComputeValues(d.Spec.Compute.Routers)
			}

		}
	}
}

func (d *DruidAutoscaler) setOpsReqOptsDefaults() {
	if d.Spec.OpsRequestOptions == nil {
		d.Spec.OpsRequestOptions = &DruidOpsRequestOptions{}
	}
	// Timeout is defaulted to 600s in ops-manager retries.go (to retry 120 times with 5sec pause between each)
	// OplogMaxLagSeconds & ObjectsCountDiffPercentage are defaults to 0
	if d.Spec.OpsRequestOptions.Apply == "" {
		d.Spec.OpsRequestOptions.Apply = opsapi.ApplyOptionIfReady
	}
}

var _ webhook.Validator = &DruidAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (d *DruidAutoscaler) ValidateCreate() (admission.Warnings, error) {
	druidLog.Info("validate create", "name", d.Name)
	return nil, d.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (d *DruidAutoscaler) ValidateUpdate(oldObj runtime.Object) (admission.Warnings, error) {
	druidLog.Info("validate create", "name", d.Name)
	return nil, d.validate()
}

func (_ *DruidAutoscaler) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (d *DruidAutoscaler) validate() error {
	if d.Spec.DatabaseRef == nil {
		return errors.New("databaseRef can't be empty")
	}
	var dr olddbapi.Druid
	err := DefaultClient.Get(context.TODO(), types.NamespacedName{
		Name:      d.Spec.DatabaseRef.Name,
		Namespace: d.Namespace,
	}, &dr)
	if err != nil {
		_ = fmt.Errorf("can't get Druid %s/%s \n", d.Namespace, d.Spec.DatabaseRef.Name)
		return err
	}

	if d.Spec.Compute != nil {
		cm := d.Spec.Compute
		if dr.Spec.Topology != nil {
			if cm.Coordinators == nil && cm.Overlords == nil && cm.MiddleManagers == nil && cm.Historicals == nil && cm.Brokers == nil && cm.Routers == nil {
				return errors.New("Spec.Compute.Coordinators, Spec.Compute.Overlords, Spec.Compute.MiddleManagers, Spec.Compute.Brokers, Spec.Compute.Brokers, Spec.Compute.Routers all are empty")
			}
		}
	}

	if d.Spec.Storage != nil {
		if dr.Spec.Topology != nil {
			if d.Spec.Storage.MiddleManagers == nil && d.Spec.Storage.Historicals == nil {
				return errors.New("Spec.Storage.MiddleManagers and Spec.Storage.Historicals both are empty")
			}
		}
	}
	return nil
}
