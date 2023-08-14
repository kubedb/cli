/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package describer

import (
	"context"
	"io"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"kubedb.dev/apimachinery/client/clientset/versioned/scheme"
	cs "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha2"

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/describe"
	condutil "kmodules.xyz/client-go/conditions"
	"kmodules.xyz/client-go/discovery"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned"
	stashV1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	stash "stash.appscode.dev/apimachinery/client/clientset/versioned"
)

type ElasticsearchDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha2Interface
	stash  stash.Interface
	appcat appcat_cs.Interface
}

func (d *ElasticsearchDescriber) Describe(namespace, name string, describerSettings describe.DescriberSettings) (string, error) {
	item, err := d.kubedb.Elasticsearches(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	selector := labels.SelectorFromSet(item.OffshootSelectors())

	var events *core.EventList
	if describerSettings.ShowEvents {
		events, err = d.client.CoreV1().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeElasticsearch(item, selector, events)
}

func (d *ElasticsearchDescriber) describeElasticsearch(item *api.Elasticsearch, selector labels.Selector, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := describe.NewPrefixWriter(out)
		w.Write(LEVEL_0, "Name:\t%s\n", item.Name)
		w.Write(LEVEL_0, "Namespace:\t%s\n", item.Namespace)
		w.Write(LEVEL_0, "CreationTimestamp:\t%s\n", timeToString(&item.CreationTimestamp))
		printLabelsMultiline(LEVEL_0, w, "Labels", item.Labels)
		printAnnotationsMultiline(LEVEL_0, w, "Annotations", item.Annotations)
		w.Write(LEVEL_0, "Status:\t%s\n", string(item.Status.Phase))

		if item.Spec.Replicas != nil {
			w.Write(LEVEL_0, "Replicas:\t%d  total\n", pointer.Int32(item.Spec.Replicas))
		}

		describeStorage(item.Spec.StorageType, item.Spec.Storage, w)

		w.Write(LEVEL_0, "Paused:\t%v\n", condutil.IsConditionTrue(item.Status.Conditions, api.DatabasePaused))
		w.Write(LEVEL_0, "Halted:\t%v\n", item.Spec.Halted)
		w.Write(LEVEL_0, "Termination Policy:\t%v\n", item.Spec.TerminationPolicy)

		showWorkload(d.client, item.Namespace, selector, w)

		secrets := make(map[string]*api.SecretReference)
		if item.Spec.AuthSecret != nil {
			secrets["Auth"] = item.Spec.AuthSecret
		}
		showSecret(d.client, item.Namespace, secrets, w)

		specific := map[string]labels.Selector{
			"master": labels.SelectorFromSet(map[string]string{"node.role.master": "set"}),
			"client": labels.SelectorFromSet(map[string]string{"node.role.client": "set"}),
			"data":   labels.SelectorFromSet(map[string]string{"node.role.data": "set"}),
		}
		showTopology(d.client, item.Namespace, selector, specific, w)

		if item.Spec.Monitor != nil {
			describeMonitor(item.Spec.Monitor, w)
		}

		// Show Initialization information
		describeInitialization(item.Spec.Init, w)

		ab, err := d.appcat.AppcatalogV1alpha1().AppBindings(item.Namespace).Get(context.TODO(), item.Name, metav1.GetOptions{})
		if err != nil && !kerr.IsNotFound(err) {
			return err
		}

		// Show Backup information
		if discovery.ExistsGroupKind(d.client.Discovery(), stashV1beta1.SchemeGroupVersion.Group, stashV1beta1.ResourceKindBackupBlueprint) {
			err = showBackups(d.stash, ab, w)
			if err != nil {
				return err
			}
		}

		// Show AppBinding
		if ab != nil {
			err = showAppBinding(ab, w)
			if err != nil {
				return err
			}
		}

		if events != nil {
			DescribeEvents(events, w)
		}

		return nil
	})
}
