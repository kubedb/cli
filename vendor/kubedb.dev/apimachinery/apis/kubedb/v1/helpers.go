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

package v1

import (
	"context"
	"fmt"
	"time"

	"kubedb.dev/apimachinery/apis/kubedb"

	cm_api "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	meta_util "kmodules.xyz/client-go/meta"
	ofstv1 "kmodules.xyz/offshoot-api/api/v1"
	petsetutil "kubeops.dev/petset/client/clientset/versioned/typed/apps/v1"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func checkReplicas(lister pslister.PetSetNamespaceLister, selector labels.Selector, expectedItems int) (bool, string, error) {
	items, err := lister.List(selector)
	if err != nil {
		return false, "", err
	}
	if len(items) < expectedItems {
		return false, fmt.Sprintf("All PetSets are not available. Desire number of PetSet: %d, Available: %d", expectedItems, len(items)), nil
	}

	// return isReplicasReady, message, error
	ready, msg := petsetutil.PetSetsAreReady(items)
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
func GetServiceTemplate(templates []NamedServiceTemplateSpec, alias ServiceAlias) ofstv1.ServiceTemplateSpec {
	for i := range templates {
		c := templates[i]
		if c.Alias == alias {
			return c.ServiceTemplateSpec
		}
	}
	return ofstv1.ServiceTemplateSpec{}
}

func GetDatabasePods(db metav1.Object, psLister pslister.PetSetLister, pods []core.Pod) ([]core.Pod, error) {
	var dbPods []core.Pod

	for i := range pods {
		owner := metav1.GetControllerOf(&pods[i])
		if owner == nil {
			continue
		}

		// If the Pod is not control by a PetSet, then it is not a KubeDB database Pod
		if owner.Kind == kubedb.ResourceKindPetSet {
			// Find the controlling PetSet
			sts, err := psLister.PetSets(db.GetNamespace()).Get(owner.Name)
			if err != nil {
				return nil, err
			}

			// Check if the PetSet is controlled by the database
			if metav1.IsControlledBy(sts, db) {
				dbPods = append(dbPods, pods[i])
			}
		}
	}

	return dbPods, nil
}

func GetDatabasePodsByPetSetLister(db metav1.Object, psLister pslister.PetSetLister, pods []core.Pod) ([]core.Pod, error) {
	var dbPods []core.Pod

	for i := range pods {
		owner := metav1.GetControllerOf(&pods[i])
		if owner == nil {
			continue
		}

		// If the Pod is not control by a PetSet, then it is not a KubeDB database Pod
		if owner.Kind == kubedb.ResourceKindPetSet {
			// Find the controlling PetSet
			ps, err := psLister.PetSets(db.GetNamespace()).Get(owner.Name)
			if err != nil {
				return nil, err
			}

			// Check if the PetSet is controlled by the database
			if metav1.IsControlledBy(ps, db) {
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

func GetSelectorForNetworkPolicy() map[string]string {
	return map[string]string{
		meta_util.ComponentLabelKey: kubedb.ComponentDatabase,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func GetActivationTimeFromSecret(secretName *core.Secret) (*metav1.Time, error) {
	if val, exists := secretName.Annotations[kubedb.AuthActiveFromAnnotation]; exists {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		return &metav1.Time{Time: t}, nil
	}
	return nil, nil
}
