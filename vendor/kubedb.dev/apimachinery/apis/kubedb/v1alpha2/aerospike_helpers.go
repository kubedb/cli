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
	"kubedb.dev/apimachinery/apis/kubedb"

	meta_util "kmodules.xyz/client-go/meta"
)

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
