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
	api "kubedb.dev/apimachinery/apis/kubedb/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ Accessor = &MariaDBArchiver{}

func (m *MariaDBArchiver) GetObjectMeta() metav1.ObjectMeta {
	return m.ObjectMeta
}

func (m *MariaDBArchiver) GetConsumers() *api.AllowedConsumers {
	return m.Spec.Databases
}

var _ ListAccessor = &MariaDBArchiverList{}

func (l *MariaDBArchiverList) GetItems() []Accessor {
	var accessors []Accessor
	for _, item := range l.Items {
		accessors = append(accessors, &item)
	}
	return accessors
}
