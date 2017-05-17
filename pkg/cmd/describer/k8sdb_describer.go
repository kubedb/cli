package describer

import (
	"fmt"
	"io"
	"time"

	"github.com/golang/glog"
	tapi "github.com/k8sdb/apimachinery/api"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/labels"
)

const (
	LabelDatabaseKind = "k8sdb.com/kind"
	LabelDatabaseName = "k8sdb.com/name"
)

func (d *humanReadableDescriber) describeElastic(item *tapi.Elastic, describerSettings *kubectl.DescriberSettings) (string, error) {
	clientSet, err := d.ClientSet()
	if err != nil {
		return "", err
	}

	snapshots, err := d.extensionsClient.Snapshots(item.Namespace).List(
		kapi.ListOptions{
			LabelSelector: labels.SelectorFromSet(
				map[string]string{
					LabelDatabaseKind: tapi.ResourceKindElastic,
					LabelDatabaseName: item.Name,
				},
			),
		},
	)
	if err != nil {
		return "", err
	}

	var events *kapi.EventList
	if describerSettings.ShowEvents {
		if ref, err := kapi.GetReference(item); err != nil {
			glog.Errorf("Unable to construct reference to '%#v': %v", item, err)
		} else {
			ref.Kind = ""
			events, err = clientSet.Core().Events(item.Namespace).Search(ref)
			if err != nil {
				return "", err
			}
		}
	}

	return tabbedString(func(out io.Writer) error {
		fmt.Fprintf(out, "Name:\t%s\n", item.Name)
		fmt.Fprintf(out, "Namespace:\t%s\n", item.Namespace)
		fmt.Fprintf(out, "CreationTimestamp:\t%s\n", item.CreationTimestamp.Time.Format(time.RFC1123Z))
		if item.Labels != nil {
			printLabelsMultiline(out, "Labels:", item.Labels)
		}
		fmt.Fprintf(out, "Status:\t%s\n", string(item.Status.Phase))
		if len(item.Status.Reason) > 0 {
			fmt.Fprintf(out, "Reason:\t%s\n", item.Status.Reason)
		}
		fmt.Fprintf(out, "Replicas:\t%d  total\n", item.Spec.Replicas)
		if item.Annotations != nil {
			printLabelsMultiline(out, "Annotations", item.Annotations)
		}

		describeStorage(item.Spec.Storage, out)

		statefulSetName := fmt.Sprintf("%v-%v", item.Name, tapi.ResourceCodeElastic)

		d.describeStatefulSet(item.Namespace, statefulSetName, out)
		d.describeService(item.Namespace, item.Name, out)

		listSnapshots(snapshots, out)

		if events != nil {
			describeEvents(events, out)
		}

		return nil
	})
}

func (d *humanReadableDescriber) describePostgres(item *tapi.Postgres, describerSettings *kubectl.DescriberSettings) (string, error) {
	clientSet, err := d.ClientSet()
	if err != nil {
		return "", err
	}

	snapshots, err := d.extensionsClient.Snapshots(item.Namespace).List(
		kapi.ListOptions{
			LabelSelector: labels.SelectorFromSet(
				map[string]string{
					LabelDatabaseKind: tapi.ResourceKindPostgres,
					LabelDatabaseName: item.Name,
				},
			),
		},
	)
	if err != nil {
		return "", err
	}

	var events *kapi.EventList
	if describerSettings.ShowEvents {
		if ref, err := kapi.GetReference(item); err != nil {
			glog.Errorf("Unable to construct reference to '%#v': %v", item, err)
		} else {
			ref.Kind = ""
			events, err = clientSet.Core().Events(item.Namespace).Search(ref)
			if err != nil {
				return "", err
			}
		}
	}

	return tabbedString(func(out io.Writer) error {
		fmt.Fprintf(out, "Name:\t%s\n", item.Name)
		fmt.Fprintf(out, "Namespace:\t%s\n", item.Namespace)
		fmt.Fprintf(out, "StartTimestamp:\t%s\n", item.CreationTimestamp.Time.Format(time.RFC1123Z))
		if item.Labels != nil {
			printLabelsMultiline(out, "Labels:", item.Labels)
		}
		fmt.Fprintf(out, "Status:\t%s\n", string(item.Status.Phase))
		if len(item.Status.Reason) > 0 {
			fmt.Fprintf(out, "Reason:\t%s\n", item.Status.Reason)
		}
		if item.Annotations != nil {
			printLabelsMultiline(out, "Annotations", item.Annotations)
		}

		describeStorage(item.Spec.Storage, out)

		statefulSetName := fmt.Sprintf("%v-%v", item.Name, tapi.ResourceCodePostgres)

		d.describeStatefulSet(item.Namespace, statefulSetName, out)
		d.describeService(item.Namespace, item.Name, out)
		if item.Spec.DatabaseSecret != nil {
			d.describeSecret(item.Namespace, item.Spec.DatabaseSecret.SecretName, "Database", out)
		}

		listSnapshots(snapshots, out)

		if events != nil {
			describeEvents(events, out)
		}

		return nil
	})
}

func (d *humanReadableDescriber) describeSnapshot(item *tapi.Snapshot, describerSettings *kubectl.DescriberSettings) (string, error) {
	clientSet, err := d.ClientSet()
	if err != nil {
		return "", err
	}

	var events *kapi.EventList
	if describerSettings.ShowEvents {
		if ref, err := kapi.GetReference(item); err != nil {
			glog.Errorf("Unable to construct reference to '%#v': %v", item, err)
		} else {
			ref.Kind = ""
			events, err = clientSet.Core().Events(item.Namespace).Search(ref)
			if err != nil {
				return "", err
			}
		}
	}

	return tabbedString(func(out io.Writer) error {
		fmt.Fprintf(out, "Name:\t%s\n", item.Name)
		fmt.Fprintf(out, "Namespace:\t%s\n", item.Namespace)
		fmt.Fprintf(out, "CreationTimestamp:\t%s\n", item.CreationTimestamp.Format(time.RFC1123Z))
		if item.Status.CompletionTime != nil {
			fmt.Fprintf(out, "CompletionTimestamp:\t%s\n", item.Status.CompletionTime.Format(time.RFC1123Z))
		}
		if item.Labels != nil {
			printLabelsMultiline(out, "Labels:", item.Labels)
		}
		fmt.Fprintf(out, "Status:\t%s\n", string(item.Status.Phase))
		if len(item.Status.Reason) > 0 {
			fmt.Fprintf(out, "Reason:\t%s\n", item.Status.Reason)
		}
		if item.Annotations != nil {
			printLabelsMultiline(out, "Annotations", item.Annotations)
		}

		d.describeSecret(item.Namespace, item.Spec.StorageSecret.SecretName, "Storage", out)

		if events != nil {
			describeEvents(events, out)
		}

		return nil
	})
}

func (d *humanReadableDescriber) describeDeletedDatabase(item *tapi.DeletedDatabase, describerSettings *kubectl.DescriberSettings) (string, error) {
	clientSet, err := d.ClientSet()
	if err != nil {
		return "", err
	}

	var events *kapi.EventList
	if describerSettings.ShowEvents {
		if ref, err := kapi.GetReference(item); err != nil {
			glog.Errorf("Unable to construct reference to '%#v': %v", item, err)
		} else {
			ref.Kind = ""
			events, err = clientSet.Core().Events(item.Namespace).Search(ref)
			if err != nil {
				return "", err
			}
		}
	}

	return tabbedString(func(out io.Writer) error {
		fmt.Fprintf(out, "Name:\t%s\n", item.Name)
		fmt.Fprintf(out, "Namespace:\t%s\n", item.Namespace)
		fmt.Fprintf(out, "CreationTimestamp:\t%s\n", item.CreationTimestamp.Format(time.RFC1123Z))
		if item.Status.DeletionTime != nil {
			fmt.Fprintf(out, "DeletionTimestamp:\t%s\n", item.Status.DeletionTime.Format(time.RFC1123Z))
		}
		if item.Status.WipeOutTime != nil {
			fmt.Fprintf(out, "WipeOutTimestamp:\t%s\n", item.Status.WipeOutTime.Format(time.RFC1123Z))
		}
		if item.Labels != nil {
			printLabelsMultiline(out, "Labels:", item.Labels)
		}
		fmt.Fprintf(out, "Status:\t%s\n", string(item.Status.Phase))
		if len(item.Status.Reason) > 0 {
			fmt.Fprintf(out, "Reason:\t%s\n", item.Status.Reason)
		}
		if item.Annotations != nil {
			printLabelsMultiline(out, "Annotations", item.Annotations)
		}

		describeOrigin(item.Spec.Origin, out)

		if events != nil {
			describeEvents(events, out)
		}

		return nil
	})
}

func describeStorage(storage *tapi.StorageSpec, out io.Writer) {
	if storage == nil {
		fmt.Fprint(out, "No volumes.\n")
		return
	}

	accessModes := kapi.GetAccessModesAsString(storage.AccessModes)
	val, _ := storage.Resources.Requests[kapi.ResourceStorage]
	capacity := val.String()
	fmt.Fprint(out, "Volume:\n")
	fmt.Fprintf(out, "  StorageClass:\t%s\n", storage.Class)
	fmt.Fprintf(out, "  Capacity:\t%s\n", capacity)
	fmt.Fprintf(out, "  Access Modes:\t%s\n", accessModes)
}

func listSnapshots(snapshotList *tapi.SnapshotList, out io.Writer) {
	fmt.Fprint(out, "\n")

	if len(snapshotList.Items) == 0 {
		fmt.Fprint(out, "No Snapshots.\n")
		return
	}

	fmt.Fprint(out, "Snapshots:\n")
	w := kubectl.GetNewTabWriter(out)

	fmt.Fprint(w, "  Name\tBucket\tStartTime\tCompletionTime\tPhase\n")
	fmt.Fprint(w, "  ----\t------\t---------\t--------------\t-----\n")
	for _, e := range snapshotList.Items {
		fmt.Fprintf(w, "  %s\t%s\t%s\t%s\t%s\n",
			e.Name,
			e.Spec.BucketName,
			timeToString(e.Status.StartTime),
			timeToString(e.Status.CompletionTime),
			e.Status.Phase,
		)
	}
	w.Flush()
}

func describeOrigin(origin tapi.Origin, out io.Writer) {
	fmt.Fprint(out, "\n")
	fmt.Fprint(out, "Origin:\n")
	fmt.Fprintf(out, "  Name:\t%s\n", origin.Name)
	fmt.Fprintf(out, "  Namespace:\t%s\n", origin.Namespace)
	if origin.Labels != nil {
		printLabelsMultiline(out, "  Labels:", origin.Labels)
	}
	if origin.Annotations != nil {
		printLabelsMultiline(out, "  Annotations", origin.Annotations)
	}
}
