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

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var redisLog = logf.Log.WithName("redis-autoscaler")

func (in *RedisAutoscaler) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-autoscaling-kubedb-com-v1alpha1-redisautoscaler,mutating=true,failurePolicy=fail,sideEffects=None,groups=autoscaling.kubedb.com,resources=redisautoscaler,verbs=create;update,versions=v1alpha1,name=mredisautoscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &RedisAutoscaler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *RedisAutoscaler) Default() {
	redisLog.Info("defaulting", "name", in.Name)
	in.setDefaults()
}

func (in *RedisAutoscaler) setDefaults() {
	in.setOpsReqOptsDefaults()

	if in.Spec.Storage != nil {
		setDefaultStorageValues(in.Spec.Storage.Standalone)
		setDefaultStorageValues(in.Spec.Storage.Cluster)
		setDefaultStorageValues(in.Spec.Storage.Sentinel)
	}

	if in.Spec.Compute != nil {
		setDefaultComputeValues(in.Spec.Compute.Standalone)
		setDefaultComputeValues(in.Spec.Compute.Cluster)
		setDefaultComputeValues(in.Spec.Compute.Sentinel)
	}
}

func (in *RedisAutoscaler) setOpsReqOptsDefaults() {
	if in.Spec.OpsRequestOptions == nil {
		in.Spec.OpsRequestOptions = &RedisOpsRequestOptions{}
	}
	// Timeout is defaulted to 600s in ops-manager retries.go (to retry 120 times with 5sec pause between each)
	// OplogMaxLagSeconds & ObjectsCountDiffPercentage are defaults to 0
	if in.Spec.OpsRequestOptions.Apply == "" {
		in.Spec.OpsRequestOptions.Apply = opsapi.ApplyOptionIfReady
	}
}

func (in *RedisAutoscaler) SetDefaults() {
}

// +kubebuilder:webhook:path=/validate-schema-kubedb-com-v1alpha1-redisautoscaler,mutating=false,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=redisautoscalers,verbs=create;update;delete,versions=v1alpha1,name=vredisautoscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &RedisAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *RedisAutoscaler) ValidateCreate() error {
	redisLog.Info("validate create", "name", in.Name)
	return in.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *RedisAutoscaler) ValidateUpdate(old runtime.Object) error {
	redisLog.Info("validate create", "name", in.Name)
	return in.validate()
}

func (_ RedisAutoscaler) ValidateDelete() error {
	return nil
}

func (in *RedisAutoscaler) validate() error {
	if in.Spec.DatabaseRef == nil {
		return errors.New("databaseRef can't be empty")
	}
	return nil
}

func (in *RedisAutoscaler) ValidateFields(rd *dbapi.Redis) error {
	if in.Spec.Compute != nil {
		cm := in.Spec.Compute
		if rd.Spec.Mode == dbapi.RedisModeCluster {
			if cm.Standalone != nil {
				return errors.New("Spec.Compute.Standalone is invalid for clustered redis")
			}
			if cm.Sentinel != nil {
				return errors.New("Spec.Compute.Sentinel is invalid for clustered redis")
			}
		} else if rd.Spec.Mode == dbapi.RedisModeSentinel {
			if cm.Standalone != nil {
				return errors.New("Spec.Compute.Standalone is invalid for redis sentinel")
			}
			if cm.Cluster != nil {
				return errors.New("Spec.Compute.Cluster is invalid for redis sentinel")
			}
		} else if rd.Spec.Mode == dbapi.RedisModeStandalone {
			if cm.Cluster != nil {
				return errors.New("Spec.Compute.Cluster is invalid for standalone redis")
			}
			if cm.Cluster != nil {
				return errors.New("Spec.Compute.Sentinel is invalid for standalone redis")
			}

		}
	}

	if in.Spec.Storage != nil {
		st := in.Spec.Storage
		if rd.Spec.Mode == dbapi.RedisModeCluster {
			if st.Standalone != nil {
				return errors.New("Spec.Storage.Standalone is invalid for clustered redis")
			}
			if st.Sentinel != nil {
				return errors.New("Spec.Storage.Sentinel is invalid for clustered redis")
			}
		} else if rd.Spec.Mode == dbapi.RedisModeSentinel {
			if st.Standalone != nil {
				return errors.New("Spec.Storage.Standalone is invalid for redis sentinel")
			}
			if st.Cluster != nil {
				return errors.New("Spec.Storage.Cluster is invalid for redis sentinel")
			}
		} else if rd.Spec.Mode == dbapi.RedisModeStandalone {
			if st.Cluster != nil {
				return errors.New("Spec.Storage.Cluster is invalid for standalone redis")
			}
			if st.Sentinel != nil {
				return errors.New("Spec.Storage.Sentinel is invalid for standalone redis")
			}
		}
	}
	return nil
}
