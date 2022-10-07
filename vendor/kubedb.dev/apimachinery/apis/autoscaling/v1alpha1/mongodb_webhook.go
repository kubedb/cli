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
var mongoLog = logf.Log.WithName("mongodb-autoscaler")

func (in *MongoDBAutoscaler) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-autoscaling-kubedb-com-v1alpha1-mongodbautoscaler,mutating=true,failurePolicy=fail,sideEffects=None,groups=autoscaling.kubedb.com,resources=mongodbautoscaler,verbs=create;update,versions=v1alpha1,name=mmongodbautoscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &MongoDBAutoscaler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *MongoDBAutoscaler) Default() {
	mongoLog.Info("defaulting", "name", in.Name)
	in.setDefaults()
}

func (in *MongoDBAutoscaler) setDefaults() {
	in.setOpsReqOptsDefaults()

	if in.Spec.Storage != nil {
		setDefaultStorageValues(in.Spec.Storage.Standalone)
		setDefaultStorageValues(in.Spec.Storage.ReplicaSet)
		setDefaultStorageValues(in.Spec.Storage.Shard)
		setDefaultStorageValues(in.Spec.Storage.ConfigServer)
		setDefaultStorageValues(in.Spec.Storage.Hidden)
	}

	if in.Spec.Compute != nil {
		setDefaultComputeValues(in.Spec.Compute.Standalone)
		setDefaultComputeValues(in.Spec.Compute.ReplicaSet)
		setDefaultComputeValues(in.Spec.Compute.Shard)
		setDefaultComputeValues(in.Spec.Compute.ConfigServer)
		setDefaultComputeValues(in.Spec.Compute.Mongos)
		setDefaultComputeValues(in.Spec.Compute.Arbiter)
		setDefaultComputeValues(in.Spec.Compute.Hidden)
	}
}

func (in *MongoDBAutoscaler) setOpsReqOptsDefaults() {
	if in.Spec.OpsRequestOptions == nil {
		in.Spec.OpsRequestOptions = &MongoDBOpsRequestOptions{}
	}
	// Timeout is defaulted to 600s in ops-manager retries.go (to retry 120 times with 5sec pause between each)
	// OplogMaxLagSeconds & ObjectsCountDiffPercentage are defaults to 0
	if in.Spec.OpsRequestOptions.Apply == "" {
		in.Spec.OpsRequestOptions.Apply = opsapi.ApplyOptionIfReady
	}
}

func (in *MongoDBAutoscaler) SetDefaults(db *dbapi.MongoDB) {
	if in.Spec.Compute != nil {
		setInMemoryDefaults(in.Spec.Compute.Standalone, db.Spec.StorageEngine)
		setInMemoryDefaults(in.Spec.Compute.ReplicaSet, db.Spec.StorageEngine)
		setInMemoryDefaults(in.Spec.Compute.Shard, db.Spec.StorageEngine)
		setInMemoryDefaults(in.Spec.Compute.ConfigServer, db.Spec.StorageEngine)
		setInMemoryDefaults(in.Spec.Compute.Mongos, db.Spec.StorageEngine)
		// no need for Defaulting the Arbiter & Hidden Node.
		// As arbiter is not a data-node.  And hidden doesn't have the impact of storageEngine (it can't be InMemory).
	}
}

// +kubebuilder:webhook:path=/validate-schema-kubedb-com-v1alpha1-mongodbautoscaler,mutating=false,failurePolicy=fail,sideEffects=None,groups=schema.kubedb.com,resources=mongodbautoscalers,verbs=create;update;delete,versions=v1alpha1,name=vmongodbautoscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MongoDBAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *MongoDBAutoscaler) ValidateCreate() error {
	mongoLog.Info("validate create", "name", in.Name)
	return in.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *MongoDBAutoscaler) ValidateUpdate(old runtime.Object) error {
	mongoLog.Info("validate create", "name", in.Name)
	return in.validate()
}

func (_ MongoDBAutoscaler) ValidateDelete() error {
	return nil
}

func (in *MongoDBAutoscaler) validate() error {
	if in.Spec.DatabaseRef == nil {
		return errors.New("databaseRef can't be empty")
	}
	return nil
}

func (in *MongoDBAutoscaler) ValidateFields(mg *dbapi.MongoDB) error {
	if in.Spec.Compute != nil {
		cm := in.Spec.Compute
		if mg.Spec.ShardTopology != nil {
			if cm.ReplicaSet != nil {
				return errors.New("Spec.Compute.ReplicaSet is invalid for sharded mongoDB")
			}
			if cm.Standalone != nil {
				return errors.New("Spec.Compute.Standalone is invalid for sharded mongoDB")
			}
		} else if mg.Spec.ReplicaSet != nil {
			if cm.Standalone != nil {
				return errors.New("Spec.Compute.Standalone is invalid for replicaSet mongoDB")
			}
			if cm.Shard != nil {
				return errors.New("Spec.Compute.Shard is invalid for replicaSet mongoDB")
			}
			if cm.ConfigServer != nil {
				return errors.New("Spec.Compute.ConfigServer is invalid for replicaSet mongoDB")
			}
			if cm.Mongos != nil {
				return errors.New("Spec.Compute.Mongos is invalid for replicaSet mongoDB")
			}
		} else {
			if cm.ReplicaSet != nil {
				return errors.New("Spec.Compute.Replicaset is invalid for Standalone mongoDB")
			}
			if cm.Shard != nil {
				return errors.New("Spec.Compute.Shard is invalid for Standalone mongoDB")
			}
			if cm.ConfigServer != nil {
				return errors.New("Spec.Compute.ConfigServer is invalid for Standalone mongoDB")
			}
			if cm.Mongos != nil {
				return errors.New("Spec.Compute.Mongos is invalid for Standalone mongoDB")
			}
			if cm.Arbiter != nil {
				return errors.New("Spec.Compute.Arbiter is invalid for Standalone mongoDB")
			}
			if cm.Hidden != nil {
				return errors.New("Spec.Compute.Hidden is invalid for Standalone mongoDB")
			}
		}
	}

	if in.Spec.Storage != nil {
		st := in.Spec.Storage
		if mg.Spec.ShardTopology != nil {
			if st.ReplicaSet != nil {
				return errors.New("Spec.Storage.ReplicaSet is invalid for sharded mongoDB")
			}
			if st.Standalone != nil {
				return errors.New("Spec.Storage.Standalone is invalid for sharded mongoDB")
			}
		} else if mg.Spec.ReplicaSet != nil {
			if st.Standalone != nil {
				return errors.New("Spec.Storage.Standalone is invalid for replicaSet mongoDB")
			}
			if st.Shard != nil {
				return errors.New("Spec.Storage.Shard is invalid for replicaSet mongoDB")
			}
			if st.ConfigServer != nil {
				return errors.New("Spec.Storage.ConfigServer is invalid for replicaSet mongoDB")
			}
		} else {
			if st.ReplicaSet != nil {
				return errors.New("Spec.Storage.Replicaset is invalid for Standalone mongoDB")
			}
			if st.Shard != nil {
				return errors.New("Spec.Storage.Shard is invalid for Standalone mongoDB")
			}
			if st.ConfigServer != nil {
				return errors.New("Spec.Storage.ConfigServer is invalid for Standalone mongoDB")
			}
			if st.Hidden != nil {
				return errors.New("Spec.Storage.Hidden is invalid for Standalone mongoDB")
			}
		}
	}
	return nil
}
