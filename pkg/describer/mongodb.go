/*
Copyright The KubeDB Authors.

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

package describer

import (
	"io"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	"kubedb.dev/apimachinery/client/clientset/versioned/scheme"
	cs "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"

	"github.com/appscode/go/types"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/describe"
	"k8s.io/kubectl/pkg/describe/versioned"
	stash "stash.appscode.dev/apimachinery/client/clientset/versioned"
)

type MongoDBDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
	stash  stash.Interface
}

func (d *MongoDBDescriber) Describe(namespace, name string, describerSettings describe.DescriberSettings) (string, error) {
	item, err := d.kubedb.MongoDBs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	selector := labels.SelectorFromSet(item.OffshootSelectors())

	snapshots, err := d.kubedb.Snapshots(item.Namespace).List(
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)
	if err != nil {
		return "", err
	}

	var events *core.EventList
	if describerSettings.ShowEvents {
		events, err = d.client.CoreV1().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeMongoDB(item, selector, snapshots, events)
}

func (d *MongoDBDescriber) describeMongoDB(item *api.MongoDB, selector labels.Selector, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := versioned.NewPrefixWriter(out)
		w.Write(LEVEL_0, "Name:\t%s\n", item.Name)
		w.Write(LEVEL_0, "Namespace:\t%s\n", item.Namespace)
		w.Write(LEVEL_0, "CreationTimestamp:\t%s\n", timeToString(&item.CreationTimestamp))
		printLabelsMultiline(LEVEL_0, w, "Labels", item.Labels)
		printAnnotationsMultiline(LEVEL_0, w, "Annotations", item.Annotations)

		if item.Spec.Replicas != nil {
			w.Write(LEVEL_0, "Replicas:\t%d  total\n", types.Int32(item.Spec.Replicas))
		}
		w.Write(LEVEL_0, "Status:\t%s\n", string(item.Status.Phase))
		if len(item.Status.Reason) > 0 {
			w.Write(LEVEL_0, "Reason:\t%s\n", item.Status.Reason)
		}

		describeStorage(item.Spec.StorageType, item.Spec.Storage, w)

		showWorkload(d.client, item.Namespace, selector, w)

		secretVolumes := make(map[string]*core.SecretVolumeSource)
		if item.Spec.DatabaseSecret != nil {
			secretVolumes["Database"] = item.Spec.DatabaseSecret
		}
		showSecret(d.client, item.Namespace, secretVolumes, w)

		if item.Spec.Monitor != nil {
			describeMonitor(item.Spec.Monitor, w)
		}

		if snapshots != nil {
			listSnapshots(snapshots, w)
		}

		if events != nil {
			DescribeEvents(events, w)
		}

		return nil
	})
}
