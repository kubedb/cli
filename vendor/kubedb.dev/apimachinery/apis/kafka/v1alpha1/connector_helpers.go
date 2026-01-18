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
	"fmt"

	"kubedb.dev/apimachinery/apis/kafka"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/apimachinery/crds"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/apiextensions"
)

func (*Connector) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralConnector))
}

func (k *Connector) AsOwner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(ResourceKindConnector))
}

func (k *Connector) Default() {
	if k.Spec.DeletionPolicy == "" {
		k.Spec.DeletionPolicy = dbapi.DeletionPolicyDelete
	}

	k.Spec.Configuration = copyConfigurationField(k.Spec.Configuration, &k.Spec.ConfigSecret)
}

func (k *Connector) ResourceShortCode() string {
	return ResourceCodeConnector
}

func (k *Connector) ResourceKind() string {
	return ResourceKindConnector
}

func (k *Connector) ResourceSingular() string {
	return ResourceSingularConnector
}

func (k *Connector) ResourcePlural() string {
	return ResourcePluralConnector
}

func (k *Connector) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", k.ResourcePlural(), kafka.GroupName)
}

// Owner returns owner reference to resources
func (k *Connector) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(k, SchemeGroupVersion.WithKind(k.ResourceKind()))
}

func (k *Connector) OffshootName() string {
	return k.Name
}
