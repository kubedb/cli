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
	"context"
	"fmt"

	cm_api "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	apps_util "kmodules.xyz/client-go/apps/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func checkReplicas(lister appslister.StatefulSetNamespaceLister, selector labels.Selector, expectedItems int) (bool, string, error) {
	items, err := lister.List(selector)
	if err != nil {
		return false, "", err
	}
	if len(items) < expectedItems {
		return false, fmt.Sprintf("All StatefulSets are not available. Desire number of StatefulSet: %d, Available: %d", expectedItems, len(items)), nil
	}

	// return isReplicasReady, message, error
	ready, msg := apps_util.StatefulSetsAreReady(items)
	return ready, msg, nil
}

// HasServiceTemplate returns "true" if the desired serviceTemplate provided in "aliaS" is present in the serviceTemplate list.
// Otherwise, it returns "false".
func HasServiceTemplate(templates []NamedServiceTemplateSpec, alias ServiceAlias) bool {
	for i := range templates {
		if templates[i].Alias == alias {
			return true
		}
	}
	return false
}

// GetServiceTemplate returns a pointer to the desired serviceTemplate referred by "aliaS". Otherwise, it returns nil.
func GetServiceTemplate(templates []NamedServiceTemplateSpec, alias ServiceAlias) ofst.ServiceTemplateSpec {
	for i := range templates {
		c := templates[i]
		if c.Alias == alias {
			return c.ServiceTemplateSpec
		}
	}
	return ofst.ServiceTemplateSpec{}
}

func GetDatabasePods(db metav1.Object, stsLister appslister.StatefulSetLister, pods []core.Pod) ([]core.Pod, error) {
	var dbPods []core.Pod

	for i := range pods {
		owner := metav1.GetControllerOf(&pods[i])
		if owner == nil {
			continue
		}

		// If the Pod is not control by a StatefulSet, then it is not a KubeDB database Pod
		if owner.Kind == ResourceKindStatefulSet {
			// Find the controlling StatefulSet
			sts, err := stsLister.StatefulSets(db.GetNamespace()).Get(owner.Name)
			if err != nil {
				return nil, err
			}

			// Check if the StatefulSet is controlled by the database
			if metav1.IsControlledBy(sts, db) {
				dbPods = append(dbPods, pods[i])
			}
		}
	}

	return dbPods, nil
}

// Upsert elements to string slice
func upsertStringSlice(inSlice []string, values ...string) []string {
	upsert := func(m string) {
		for _, v := range inSlice {
			if v == m {
				return
			}
		}
		inSlice = append(inSlice, m)
	}

	for _, value := range values {
		upsert(value)
	}
	return inSlice
}

func UsesAcmeIssuer(kc client.Client, ns string, issuerRef core.TypedLocalObjectReference) (bool, error) {
	switch issuerRef.Kind {
	case cm_api.IssuerKind:
		var issuer cm_api.Issuer
		err := kc.Get(context.TODO(), client.ObjectKey{Namespace: ns, Name: issuerRef.Name}, &issuer)
		if err != nil {
			return false, err
		}
		return issuer.Spec.ACME != nil, nil
	case cm_api.ClusterIssuerKind:
		var issuer cm_api.ClusterIssuer
		err := kc.Get(context.TODO(), client.ObjectKey{Name: issuerRef.Name}, &issuer)
		if err != nil {
			return false, err
		}
		return issuer.Spec.ACME != nil, nil
	default:
		return false, fmt.Errorf("invalid issuer %+v", issuerRef)
	}
}
