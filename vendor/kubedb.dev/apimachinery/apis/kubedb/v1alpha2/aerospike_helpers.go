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

package v1alpha2

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"gomodules.xyz/pointer"
	corev1 "k8s.io/api/core/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

func (a *Aerospike) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralAerospike))
}

func (a *Aerospike) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", a.ResourcePlural(), kubedb.GroupName)
}

func (a *Aerospike) ResourcePlural() string {
	return ResourcePluralAerospike
}

func (a *Aerospike) ConfigSecretName() string {
	uid := string(a.UID)
	return meta_util.NameWithSuffix(a.OffshootName(), uid[len(uid)-6:])
}

func (a *Aerospike) OffshootName() string {
	return a.Name
}

func (a *Aerospike) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.InstanceLabelKey:  a.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (a *Aerospike) GoverningServiceName() string {
	return meta_util.NameWithSuffix(a.ServiceName(), "pods")
}

func (a *Aerospike) ServiceName() string {
	return a.OffshootName()
}

func (a Aerospike) OffshootLabels() map[string]string {
	return a.offshootLabels(a.OffshootSelectors(), nil)
}

func (a Aerospike) offshootLabels(selector, overrides map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, a.Labels, overrides))
}

func (a Aerospike) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(a.Spec.ServiceTemplates, alias)
	return a.offshootLabels(meta_util.OverwriteKeys(a.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (a *Aerospike) SetDefaults(arVersion *catalog.AerospikeVersion) {
	if a == nil {
		return
	}

	// perform defaulting
	switch a.Spec.Mode {
	case "":
		a.Spec.Mode = AerospikeModeStandalone
	case AerospikeModeCluster:
		if a.Spec.Cluster == nil {
			a.Spec.Cluster = &AerospikeClusterSpec{}
		}
		if a.Spec.Cluster.Replicas == nil {
			a.Spec.Cluster.Replicas = pointer.Int32P(3)
		}
		if a.Spec.Cluster.ReplicationFactor == nil {
			a.Spec.Cluster.ReplicationFactor = pointer.Int32P(2)
		}
	}
	if a.Spec.DeletionPolicy == "" {
		a.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if !a.Spec.DisableAuth {
		if a.Spec.AuthSecret == nil {
			a.Spec.AuthSecret = &SecretReference{}
		}
		if a.Spec.AuthSecret.Kind == "" {
			a.Spec.AuthSecret.Kind = kubedb.ResourceKindSecret
		}
	}

	a.setDefaultContainerSecurityContext(arVersion, &a.Spec.PodTemplate)
	if a.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		a.Spec.PodTemplate.Spec.ServiceAccountName = a.OffshootName()
	}

	container := coreutil.GetContainerByName(a.Spec.PodTemplate.Spec.Containers, kubedb.AerospikeContainerName)
	if container == nil {
		container = &corev1.Container{
			Name: kubedb.AerospikeContainerName,
		}
	}
	apis.SetDefaultResourceLimits(&container.Resources, kubedb.DefaultResources)
	a.Spec.PodTemplate.Spec.Containers = coreutil.UpsertContainer(a.Spec.PodTemplate.Spec.Containers, *container)
}

func (a Aerospike) setDefaultContainerSecurityContext(arVersion *catalog.AerospikeVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &corev1.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = arVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.AerospikeContainerName)
	if container == nil {
		container = &corev1.Container{
			Name: kubedb.AerospikeContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &corev1.SecurityContext{}
	}

	a.assignDefaultContainerSecurityContext(arVersion, container.SecurityContext)

	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (a *Aerospike) assignDefaultContainerSecurityContext(arVersion *catalog.AerospikeVersion, sc *corev1.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(true)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = arVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = arVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}
