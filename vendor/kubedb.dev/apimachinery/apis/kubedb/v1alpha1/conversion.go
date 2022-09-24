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
	"reflect"
	"unsafe"

	"kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	"gomodules.xyz/encoding/json/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/util/intstr"
	apiv1 "kmodules.xyz/monitoring-agent-api/api/v1"
	apiv1alpha1 "kmodules.xyz/monitoring-agent-api/api/v1alpha1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

func Convert_v1alpha1_InitSpec_To_v1alpha2_InitSpec(in *InitSpec, out *v1alpha2.InitSpec, s conversion.Scope) error {
	if in.ScriptSource != nil {
		in, out := &in.ScriptSource, &out.Script
		*out = new(v1alpha2.ScriptSourceSpec)
		if err := Convert_v1alpha1_ScriptSourceSpec_To_v1alpha2_ScriptSourceSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Script = nil
	}
	// WARNING: in.SnapshotSource requires manual conversion: does not exist in peer-type
	// WARNING: in.PostgresWAL requires manual conversion: does not exist in peer-type
	// WARNING: in.StashRestoreSession requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_InitSpec_To_v1alpha1_InitSpec(in *v1alpha2.InitSpec, out *InitSpec, s conversion.Scope) error {
	if in.Script != nil {
		in, out := &in.Script, &out.ScriptSource
		*out = new(ScriptSourceSpec)
		if err := Convert_v1alpha2_ScriptSourceSpec_To_v1alpha1_ScriptSourceSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ScriptSource = nil
	}
	// WARNING: in.Initialized requires manual conversion: does not exist in peer-type
	// WARNING: in.WaitForInitialRestore requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_types_IntHash_To_int64(in *types.IntHash, out *int64, s conversion.Scope) error {
	*out = in.Generation()
	return nil
}

func Convert_int64_To_types_IntHash(in *int64, out *types.IntHash, s conversion.Scope) error {
	g := types.IntHashForGeneration(*in)
	*out = *g
	return nil
}

func Convert_v1alpha1_ElasticsearchSpec_To_v1alpha2_ElasticsearchSpec(in *ElasticsearchSpec, out *v1alpha2.ElasticsearchSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if in.Topology != nil {
		in, out := &in.Topology, &out.Topology
		*out = new(v1alpha2.ElasticsearchClusterTopology)
		if err := Convert_v1alpha1_ElasticsearchClusterTopology_To_v1alpha2_ElasticsearchClusterTopology(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Topology = nil
	}
	out.EnableSSL = in.EnableSSL
	// WARNING: in.CertificateSecret requires manual conversion: does not exist in peer-type
	// WARNING: in.AuthPlugin requires manual conversion: does not exist in peer-type
	if in.DatabaseSecret != nil {
		out.AuthSecret = &v1alpha2.SecretReference{
			LocalObjectReference: v1.LocalObjectReference{
				Name: in.DatabaseSecret.SecretName,
			},
		}
	}
	out.StorageType = v1alpha2.StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(v1alpha2.InitSpec)
		if err := Convert_v1alpha1_InitSpec_To_v1alpha2_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	// WARNING: in.BackupSchedule requires manual conversion: does not exist in peer-type
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}
	out.PodTemplate = in.PodTemplate
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	out.MaxUnavailable = (*intstr.IntOrString)(unsafe.Pointer(in.MaxUnavailable))
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_ElasticsearchSpec_To_v1alpha1_ElasticsearchSpec(in *v1alpha2.ElasticsearchSpec, out *ElasticsearchSpec, s conversion.Scope) error {
	out.Version = types.StrYo(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if in.Topology != nil {
		in, out := &in.Topology, &out.Topology
		*out = new(ElasticsearchClusterTopology)
		if err := Convert_v1alpha2_ElasticsearchClusterTopology_To_v1alpha1_ElasticsearchClusterTopology(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Topology = nil
	}
	out.EnableSSL = in.EnableSSL
	// WARNING: in.DisableSecurity requires manual conversion: does not exist in peer-type
	// WARNING: in.AuthSecret requires manual conversion: does not exist in peer-type
	out.StorageType = StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(InitSpec)
		if err := Convert_v1alpha2_InitSpec_To_v1alpha1_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1alpha1.AgentSpec)
		if err := apiv1alpha1.Convert_v1_AgentSpec_To_v1alpha1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	// WARNING: in.ConfigSecret requires manual conversion: does not exist in peer-type
	// WARNING: in.SecureConfigSecret requires manual conversion: does not exist in peer-type
	out.PodTemplate = in.PodTemplate
	// WARNING: in.ServiceTemplates requires manual conversion: does not exist in peer-type
	out.MaxUnavailable = (*intstr.IntOrString)(unsafe.Pointer(in.MaxUnavailable))
	// WARNING: in.TLS requires manual conversion: does not exist in peer-type
	// WARNING: in.InternalUsers requires manual conversion: does not exist in peer-type
	// WARNING: in.RolesMapping requires manual conversion: does not exist in peer-type
	// WARNING: in.Halted requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = TerminationPolicy(in.TerminationPolicy)
	// WARNING: in.KernelSettings requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_ElasticsearchClusterTopology_To_v1alpha2_ElasticsearchClusterTopology(in *ElasticsearchClusterTopology, out *v1alpha2.ElasticsearchClusterTopology, s conversion.Scope) error {
	if err := Convert_v1alpha1_ElasticsearchNode_To_v1alpha2_ElasticsearchNode(&in.Master, &out.Master, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ElasticsearchNode_To_v1alpha2_ElasticsearchNode(&in.Client, &out.Ingest, s); err != nil {
		return err
	}
	{
		out := &out.Data
		*out = new(v1alpha2.ElasticsearchNode)
		if err := Convert_v1alpha1_ElasticsearchNode_To_v1alpha2_ElasticsearchNode(&in.Data, *out, s); err != nil {
			return err
		}
	}

	return nil
}

func Convert_v1alpha2_ElasticsearchClusterTopology_To_v1alpha1_ElasticsearchClusterTopology(in *v1alpha2.ElasticsearchClusterTopology, out *ElasticsearchClusterTopology, s conversion.Scope) error {
	if err := Convert_v1alpha2_ElasticsearchNode_To_v1alpha1_ElasticsearchNode(&in.Master, &out.Master, s); err != nil {
		return err
	}
	// WARNING: in.Ingest requires manual conversion: does not exist in peer-type
	// WARNING: in.Data requires manual conversion: inconvertible types (*kubedb.dev/apimachinery/apis/kubedb/v1alpha2.ElasticsearchNode vs kubedb.dev/apimachinery/apis/kubedb/v1alpha1.ElasticsearchNode)
	// WARNING: in.DataContent requires manual conversion: does not exist in peer-type
	// WARNING: in.DataHot requires manual conversion: does not exist in peer-type
	// WARNING: in.DataWarm requires manual conversion: does not exist in peer-type
	// WARNING: in.DataCold requires manual conversion: does not exist in peer-type
	// WARNING: in.DataFrozen requires manual conversion: does not exist in peer-type
	// WARNING: in.ML requires manual conversion: does not exist in peer-type
	// WARNING: in.Transform requires manual conversion: does not exist in peer-type
	// WARNING: in.Coordinating requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_ElasticsearchNode_To_v1alpha2_ElasticsearchNode(in *ElasticsearchNode, out *v1alpha2.ElasticsearchNode, s conversion.Scope) error {
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.Suffix = in.Prefix
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	out.Resources = in.Resources
	out.MaxUnavailable = (*intstr.IntOrString)(unsafe.Pointer(in.MaxUnavailable))
	return nil
}

func Convert_v1alpha2_ElasticsearchNode_To_v1alpha1_ElasticsearchNode(in *v1alpha2.ElasticsearchNode, out *ElasticsearchNode, s conversion.Scope) error {
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.Prefix = in.Suffix
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	out.Resources = in.Resources
	out.MaxUnavailable = (*intstr.IntOrString)(unsafe.Pointer(in.MaxUnavailable))
	return nil
}

func Convert_v1alpha1_ElasticsearchStatus_To_v1alpha2_ElasticsearchStatus(in *ElasticsearchStatus, out *v1alpha2.ElasticsearchStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_ElasticsearchStatus_To_v1alpha1_ElasticsearchStatus(in *v1alpha2.ElasticsearchStatus, out *ElasticsearchStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_EtcdSpec_To_v1alpha2_EtcdSpec(in *EtcdSpec, out *v1alpha2.EtcdSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.StorageType = v1alpha2.StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.DatabaseSecret != nil {
		out.AuthSecret = &v1alpha2.SecretReference{
			LocalObjectReference: v1.LocalObjectReference{
				Name: in.DatabaseSecret.SecretName,
			},
		}
	}
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(v1alpha2.InitSpec)
		if err := Convert_v1alpha1_InitSpec_To_v1alpha2_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	// WARNING: in.BackupSchedule requires manual conversion: does not exist in peer-type
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	out.TLS = (*v1alpha2.TLSPolicy)(unsafe.Pointer(in.TLS))
	out.PodTemplate = in.PodTemplate
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_EtcdSpec_To_v1alpha1_EtcdSpec(in *v1alpha2.EtcdSpec, out *EtcdSpec, s conversion.Scope) error {
	out.Version = types.StrYo(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.StorageType = StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	// WARNING: in.AuthSecret requires manual conversion: does not exist in peer-type
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(InitSpec)
		if err := Convert_v1alpha2_InitSpec_To_v1alpha1_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1alpha1.AgentSpec)
		if err := apiv1alpha1.Convert_v1_AgentSpec_To_v1alpha1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	out.TLS = (*TLSPolicy)(unsafe.Pointer(in.TLS))
	out.PodTemplate = in.PodTemplate
	// WARNING: in.ServiceTemplates requires manual conversion: does not exist in peer-type
	// WARNING: in.Halted requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha1_EtcdStatus_To_v1alpha2_EtcdStatus(in *EtcdStatus, out *v1alpha2.EtcdStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_EtcdStatus_To_v1alpha1_EtcdStatus(in *v1alpha2.EtcdStatus, out *EtcdStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_MariaDBSpec_To_v1alpha2_MariaDBSpec(in *MariaDBSpec, out *v1alpha2.MariaDBSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.StorageType = v1alpha2.StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.DatabaseSecret != nil {
		out.AuthSecret = &v1alpha2.SecretReference{
			LocalObjectReference: v1.LocalObjectReference{
				Name: in.DatabaseSecret.SecretName,
			},
		}
	}
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(v1alpha2.InitSpec)
		if err := Convert_v1alpha1_InitSpec_To_v1alpha2_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}
	out.PodTemplate = in.PodTemplate
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_MariaDBSpec_To_v1alpha1_MariaDBSpec(in *v1alpha2.MariaDBSpec, out *MariaDBSpec, s conversion.Scope) error {
	out.Version = types.StrYo(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.StorageType = StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	// WARNING: in.AuthSecret requires manual conversion: does not exist in peer-type
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(InitSpec)
		if err := Convert_v1alpha2_InitSpec_To_v1alpha1_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1alpha1.AgentSpec)
		if err := apiv1alpha1.Convert_v1_AgentSpec_To_v1alpha1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	// WARNING: in.ConfigSecret requires manual conversion: does not exist in peer-type
	out.PodTemplate = in.PodTemplate
	// WARNING: in.ServiceTemplates requires manual conversion: does not exist in peer-type
	// WARNING: in.RequireSSL requires manual conversion: does not exist in peer-type
	// WARNING: in.TLS requires manual conversion: does not exist in peer-type
	// WARNING: in.Halted requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha1_MariaDBStatus_To_v1alpha2_MariaDBStatus(in *MariaDBStatus, out *v1alpha2.MariaDBStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_MariaDBStatus_To_v1alpha1_MariaDBStatus(in *v1alpha2.MariaDBStatus, out *MariaDBStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_MemcachedSpec_To_v1alpha2_MemcachedSpec(in *MemcachedSpec, out *v1alpha2.MemcachedSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}
	out.PodTemplate = in.PodTemplate
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_MemcachedSpec_To_v1alpha1_MemcachedSpec(in *v1alpha2.MemcachedSpec, out *MemcachedSpec, s conversion.Scope) error {
	out.Version = types.StrYo(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1alpha1.AgentSpec)
		if err := apiv1alpha1.Convert_v1_AgentSpec_To_v1alpha1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	// WARNING: in.ConfigSecret requires manual conversion: does not exist in peer-type
	// WARNING: in.DataVolume requires manual conversion: does not exist in peer-type
	out.PodTemplate = in.PodTemplate
	// WARNING: in.ServiceTemplates requires manual conversion: does not exist in peer-type
	// WARNING: in.TLS requires manual conversion: does not exist in peer-type
	// WARNING: in.Halted requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha1_MemcachedStatus_To_v1alpha2_MemcachedStatus(in *MemcachedStatus, out *v1alpha2.MemcachedStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_MemcachedStatus_To_v1alpha1_MemcachedStatus(in *v1alpha2.MemcachedStatus, out *MemcachedStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_MongoDBSpec_To_v1alpha2_MongoDBSpec(in *MongoDBSpec, out *v1alpha2.MongoDBSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if in.ReplicaSet != nil {
		in, out := &in.ReplicaSet, &out.ReplicaSet
		*out = new(v1alpha2.MongoDBReplicaSet)
		if err := Convert_v1alpha1_MongoDBReplicaSet_To_v1alpha2_MongoDBReplicaSet(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ReplicaSet = nil
	}
	if in.ShardTopology != nil {
		in, out := &in.ShardTopology, &out.ShardTopology
		*out = new(v1alpha2.MongoDBShardingTopology)
		if err := Convert_v1alpha1_MongoDBShardingTopology_To_v1alpha2_MongoDBShardingTopology(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ShardTopology = nil
	}
	out.StorageType = v1alpha2.StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.DatabaseSecret != nil {
		out.AuthSecret = &v1alpha2.SecretReference{
			LocalObjectReference: v1.LocalObjectReference{
				Name: in.DatabaseSecret.SecretName,
			},
		}
	}
	// FIXIT
	// WARNING: in.CertificateSecret requires manual conversion: does not exist in peer-type
	if in.ReplicaSet != nil && in.ReplicaSet.KeyFile != nil {
		out.KeyFileSecret = &v1.LocalObjectReference{
			Name: in.ReplicaSet.KeyFile.SecretName,
		}
	}
	out.ClusterAuthMode = v1alpha2.ClusterAuthMode(in.ClusterAuthMode)
	out.SSLMode = v1alpha2.SSLMode(in.SSLMode)
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(v1alpha2.InitSpec)
		if err := Convert_v1alpha1_InitSpec_To_v1alpha2_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	// WARNING: in.BackupSchedule requires manual conversion: does not exist in peer-type
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}
	out.PodTemplate = (*ofst.PodTemplateSpec)(unsafe.Pointer(in.PodTemplate))
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_MongoDBSpec_To_v1alpha1_MongoDBSpec(in *v1alpha2.MongoDBSpec, out *MongoDBSpec, s conversion.Scope) error {
	out.Version = types.StrYo(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if in.ReplicaSet != nil {
		in, out := &in.ReplicaSet, &out.ReplicaSet
		*out = new(MongoDBReplicaSet)
		if err := Convert_v1alpha2_MongoDBReplicaSet_To_v1alpha1_MongoDBReplicaSet(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ReplicaSet = nil
	}
	if in.ShardTopology != nil {
		in, out := &in.ShardTopology, &out.ShardTopology
		*out = new(MongoDBShardingTopology)
		if err := Convert_v1alpha2_MongoDBShardingTopology_To_v1alpha1_MongoDBShardingTopology(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ShardTopology = nil
	}
	out.StorageType = StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	// WARNING: in.AuthSecret requires manual conversion: does not exist in peer-type
	out.ClusterAuthMode = ClusterAuthMode(in.ClusterAuthMode)
	out.SSLMode = SSLMode(in.SSLMode)
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(InitSpec)
		if err := Convert_v1alpha2_InitSpec_To_v1alpha1_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1alpha1.AgentSpec)
		if err := apiv1alpha1.Convert_v1_AgentSpec_To_v1alpha1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	// WARNING: in.ConfigSecret requires manual conversion: does not exist in peer-type
	out.PodTemplate = (*ofst.PodTemplateSpec)(unsafe.Pointer(in.PodTemplate))
	// WARNING: in.ServiceTemplates requires manual conversion: does not exist in peer-type
	// WARNING: in.TLS requires manual conversion: does not exist in peer-type
	// WARNING: in.KeyFileSecret requires manual conversion: does not exist in peer-type
	// WARNING: in.Halted requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = TerminationPolicy(in.TerminationPolicy)
	// WARNING: in.StorageEngine requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_MongoDBMongosNode_To_v1alpha2_MongoDBMongosNode(in *MongoDBMongosNode, out *v1alpha2.MongoDBMongosNode, s conversion.Scope) error {
	if err := Convert_v1alpha1_MongoDBNode_To_v1alpha2_MongoDBNode(&in.MongoDBNode, &out.MongoDBNode, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_MongoDBNode_To_v1alpha2_MongoDBNode(in *MongoDBNode, out *v1alpha2.MongoDBNode, s conversion.Scope) error {
	out.Replicas = in.Replicas
	out.Prefix = in.Prefix
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}
	out.PodTemplate = in.PodTemplate
	return nil
}

func Convert_v1alpha2_MongoDBNode_To_v1alpha1_MongoDBNode(in *v1alpha2.MongoDBNode, out *MongoDBNode, s conversion.Scope) error {
	out.Replicas = in.Replicas
	out.Prefix = in.Prefix
	if in.ConfigSecret != nil {
		out.ConfigSource.Secret = &v1.SecretVolumeSource{
			SecretName: in.ConfigSecret.Name,
		}
	}
	out.PodTemplate = in.PodTemplate
	return nil
}

func Convert_v1alpha1_MongoDBReplicaSet_To_v1alpha2_MongoDBReplicaSet(in *MongoDBReplicaSet, out *v1alpha2.MongoDBReplicaSet, s conversion.Scope) error {
	out.Name = in.Name
	// FIXIT:
	// WARNING: in.KeyFile requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_MongoDBStatus_To_v1alpha2_MongoDBStatus(in *MongoDBStatus, out *v1alpha2.MongoDBStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_MongoDBStatus_To_v1alpha1_MongoDBStatus(in *v1alpha2.MongoDBStatus, out *MongoDBStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_MySQLSpec_To_v1alpha2_MySQLSpec(in *MySQLSpec, out *v1alpha2.MySQLSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if in.Topology != nil {
		in, out := &in.Topology, &out.Topology
		*out = new(v1alpha2.MySQLTopology)
		if err := Convert_v1alpha1_MySQLTopology_To_v1alpha2_MySQLTopology(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Topology = nil
	}
	out.StorageType = v1alpha2.StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.DatabaseSecret != nil {
		out.AuthSecret = &v1alpha2.SecretReference{
			LocalObjectReference: v1.LocalObjectReference{
				Name: in.DatabaseSecret.SecretName,
			},
		}
	}
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(v1alpha2.InitSpec)
		if err := Convert_v1alpha1_InitSpec_To_v1alpha2_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	// WARNING: in.BackupSchedule requires manual conversion: does not exist in peer-type
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}
	out.PodTemplate = in.PodTemplate
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_MySQLSpec_To_v1alpha1_MySQLSpec(in *v1alpha2.MySQLSpec, out *MySQLSpec, s conversion.Scope) error {
	out.Version = types.StrYo(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if in.Topology != nil {
		in, out := &in.Topology, &out.Topology
		*out = new(MySQLTopology)
		if err := Convert_v1alpha2_MySQLTopology_To_v1alpha1_MySQLTopology(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Topology = nil
	}
	out.StorageType = StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	// WARNING: in.AuthSecret requires manual conversion: does not exist in peer-type
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(InitSpec)
		if err := Convert_v1alpha2_InitSpec_To_v1alpha1_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1alpha1.AgentSpec)
		if err := apiv1alpha1.Convert_v1_AgentSpec_To_v1alpha1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	// WARNING: in.ConfigSecret requires manual conversion: does not exist in peer-type
	out.PodTemplate = in.PodTemplate
	// WARNING: in.ServiceTemplates requires manual conversion: does not exist in peer-type
	// WARNING: in.RequireSSL requires manual conversion: does not exist in peer-type
	// WARNING: in.TLS requires manual conversion: does not exist in peer-type
	// WARNING: in.Halted requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = TerminationPolicy(in.TerminationPolicy)
	// WARNING: in.UseAddressType requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_MySQLGroupSpec_To_v1alpha2_MySQLGroupSpec(in *MySQLGroupSpec, out *v1alpha2.MySQLGroupSpec, s conversion.Scope) error {
	out.Mode = (*v1alpha2.MySQLGroupMode)(unsafe.Pointer(in.Mode))
	out.Name = in.Name
	// WARNING: in.BaseServerID requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_MySQLStatus_To_v1alpha2_MySQLStatus(in *MySQLStatus, out *v1alpha2.MySQLStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_MySQLStatus_To_v1alpha1_MySQLStatus(in *v1alpha2.MySQLStatus, out *MySQLStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_PerconaXtraDBSpec_To_v1alpha2_PerconaXtraDBSpec(in *PerconaXtraDBSpec, out *v1alpha2.PerconaXtraDBSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	// WARNING: in.PXC requires manual conversion: does not exist in peer-type
	out.StorageType = v1alpha2.StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.DatabaseSecret != nil {
		out.AuthSecret = &v1alpha2.SecretReference{
			LocalObjectReference: v1.LocalObjectReference{
				Name: in.DatabaseSecret.SecretName,
			},
		}
	}
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(v1alpha2.InitSpec)
		if err := Convert_v1alpha1_InitSpec_To_v1alpha2_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}
	out.PodTemplate = in.PodTemplate
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_PerconaXtraDBSpec_To_v1alpha1_PerconaXtraDBSpec(in *v1alpha2.PerconaXtraDBSpec, out *PerconaXtraDBSpec, s conversion.Scope) error {
	out.Version = types.StrYo(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.StorageType = StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	// WARNING: in.AuthSecret requires manual conversion: does not exist in peer-type
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(InitSpec)
		if err := Convert_v1alpha2_InitSpec_To_v1alpha1_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1alpha1.AgentSpec)
		if err := apiv1alpha1.Convert_v1_AgentSpec_To_v1alpha1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	// WARNING: in.ConfigSecret requires manual conversion: does not exist in peer-type
	out.PodTemplate = in.PodTemplate
	// WARNING: in.ServiceTemplates requires manual conversion: does not exist in peer-type
	// WARNING: in.TLS requires manual conversion: does not exist in peer-type
	// WARNING: in.Halted requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha1_PerconaXtraDBStatus_To_v1alpha2_PerconaXtraDBStatus(in *PerconaXtraDBStatus, out *v1alpha2.PerconaXtraDBStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_PostgresSpec_To_v1alpha2_PostgresSpec(in *PostgresSpec, out *v1alpha2.PostgresSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.StandbyMode = (*v1alpha2.PostgresStandbyMode)(unsafe.Pointer(in.StandbyMode))
	out.StreamingMode = (*v1alpha2.PostgresStreamingMode)(unsafe.Pointer(in.StreamingMode))
	// WARNING: in.Archiver requires manual conversion: does not exist in peer-type
	if in.LeaderElection != nil {
		in, out := &in.LeaderElection, &out.LeaderElection
		*out = new(v1alpha2.PostgreLeaderElectionConfig)
		if err := Convert_v1alpha1_LeaderElectionConfig_To_v1alpha2_PostgreLeaderElectionConfig(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.LeaderElection = nil
	}
	if in.DatabaseSecret != nil {
		out.AuthSecret = &v1alpha2.SecretReference{
			LocalObjectReference: v1.LocalObjectReference{
				Name: in.DatabaseSecret.SecretName,
			},
		}
	}
	out.StorageType = v1alpha2.StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.Init != nil {
		in, out := &in.Init, &out.Init
		*out = new(v1alpha2.InitSpec)
		if err := Convert_v1alpha1_InitSpec_To_v1alpha2_InitSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Init = nil
	}
	// WARNING: in.BackupSchedule requires manual conversion: does not exist in peer-type
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}
	out.PodTemplate = in.PodTemplate
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	if !reflect.DeepEqual(in.ReplicaServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.StandbyServiceAlias,
			ServiceTemplateSpec: in.ReplicaServiceTemplate,
		})
	}
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_PostgresSpec_To_v1alpha1_PostgresSpec(in *v1alpha2.PostgresSpec, out *PostgresSpec, s conversion.Scope) error {
	return autoConvert_v1alpha2_PostgresSpec_To_v1alpha1_PostgresSpec(in, out, s)
}

func Convert_v1alpha2_PerconaXtraDBStatus_To_v1alpha1_PerconaXtraDBStatus(in *v1alpha2.PerconaXtraDBStatus, out *PerconaXtraDBStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_LeaderElectionConfig_To_v1alpha2_PostgreLeaderElectionConfig(in *LeaderElectionConfig, out *v1alpha2.PostgreLeaderElectionConfig, s conversion.Scope) error {
	out.LeaseDurationSeconds = in.LeaseDurationSeconds
	out.RenewDeadlineSeconds = in.RenewDeadlineSeconds
	out.RetryPeriodSeconds = in.RetryPeriodSeconds
	return nil
}

func Convert_v1alpha2_PostgreLeaderElectionConfig_To_v1alpha1_LeaderElectionConfig(in *v1alpha2.PostgreLeaderElectionConfig, out *LeaderElectionConfig, s conversion.Scope) error {
	out.LeaseDurationSeconds = in.LeaseDurationSeconds
	out.RenewDeadlineSeconds = in.RenewDeadlineSeconds
	out.RetryPeriodSeconds = in.RetryPeriodSeconds
	return nil
}

func Convert_v1alpha1_PostgresStatus_To_v1alpha2_PostgresStatus(in *PostgresStatus, out *v1alpha2.PostgresStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_PostgresStatus_To_v1alpha1_PostgresStatus(in *v1alpha2.PostgresStatus, out *PostgresStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha1_RedisSpec_To_v1alpha2_RedisSpec(in *RedisSpec, out *v1alpha2.RedisSpec, s conversion.Scope) error {
	out.Version = string(in.Version)
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	out.Mode = v1alpha2.RedisMode(in.Mode)
	out.Cluster = (*v1alpha2.RedisClusterSpec)(unsafe.Pointer(in.Cluster))
	out.StorageType = v1alpha2.StorageType(in.StorageType)
	out.Storage = (*v1.PersistentVolumeClaimSpec)(unsafe.Pointer(in.Storage))
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		if err := apiv1alpha1.Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Monitor = nil
	}
	if in.ConfigSource != nil {
		if in.ConfigSource.Secret != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: in.ConfigSource.Secret.SecretName,
			}
		} else if in.ConfigSource.ConfigMap != nil {
			out.ConfigSecret = &v1.LocalObjectReference{
				Name: "FIX_CONVERT_TO_SECRET_" + in.ConfigSource.ConfigMap.Name,
			}
		}
	}

	out.PodTemplate = in.PodTemplate
	if !reflect.DeepEqual(in.ServiceTemplate, ofst.ServiceTemplateSpec{}) {
		out.ServiceTemplates = append(out.ServiceTemplates, v1alpha2.NamedServiceTemplateSpec{
			Alias:               v1alpha2.PrimaryServiceAlias,
			ServiceTemplateSpec: in.ServiceTemplate,
		})
	}
	// WARNING: in.UpdateStrategy requires manual conversion: does not exist in peer-type
	out.TerminationPolicy = v1alpha2.TerminationPolicy(in.TerminationPolicy)
	return nil
}

func Convert_v1alpha2_RedisSpec_To_v1alpha1_RedisSpec(in *v1alpha2.RedisSpec, out *RedisSpec, s conversion.Scope) error {
	return autoConvert_v1alpha2_RedisSpec_To_v1alpha1_RedisSpec(in, out, s)
}

func Convert_v1alpha1_RedisStatus_To_v1alpha2_RedisStatus(in *RedisStatus, out *v1alpha2.RedisStatus, s conversion.Scope) error {
	out.Phase = v1alpha2.DatabasePhase(in.Phase)
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = in.ObservedGeneration.Generation()
	}
	// WARNING: in.Reason requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v1alpha2_RedisStatus_To_v1alpha1_RedisStatus(in *v1alpha2.RedisStatus, out *RedisStatus, s conversion.Scope) error {
	out.Phase = DatabasePhase(in.Phase)
	out.ObservedGeneration = types.IntHashForGeneration(in.ObservedGeneration)
	// WARNING: in.Conditions requires manual conversion: does not exist in peer-type
	return nil
}
