package describer

import (
	"io"
	"strings"
	"github.com/appscode/go/types"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/client/clientset/versioned/scheme"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
)

type EtcdDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *EtcdDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	item, err := d.kubedb.Etcds(namespace).Get(name, metav1.GetOptions{})
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
		events, err = d.client.Core().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeEtcd(item, selector, snapshots, events)
}

func (d *EtcdDescriber) describeEtcd(item *api.Etcd, selector labels.Selector, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
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

type ElasticsearchDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *ElasticsearchDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	item, err := d.kubedb.Elasticsearches(namespace).Get(name, metav1.GetOptions{})
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
		events, err = d.client.Core().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeElasticsearch(item, selector, snapshots, events)
}

func (d *ElasticsearchDescriber) describeElasticsearch(item *api.Elasticsearch, selector labels.Selector, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
		w.Write(LEVEL_0, "Name:\t%s\n", item.Name)
		w.Write(LEVEL_0, "Namespace:\t%s\n", item.Namespace)
		w.Write(LEVEL_0, "CreationTimestamp:\t%s\n", timeToString(&item.CreationTimestamp))
		printLabelsMultiline(LEVEL_0, w, "Labels", item.Labels)
		printAnnotationsMultiline(LEVEL_0, w, "Annotations", item.Annotations)
		w.Write(LEVEL_0, "Status:\t%s\n", string(item.Status.Phase))
		if len(item.Status.Reason) > 0 {
			w.Write(LEVEL_0, "Reason:\t%s\n", item.Status.Reason)
		}

		if item.Spec.Replicas != nil {
			w.Write(LEVEL_0, "Replicas:\t%d  total\n", types.Int32(item.Spec.Replicas))
		}

		describeInitialization(item.Spec.Init, w)

		describeStorage(item.Spec.StorageType, item.Spec.Storage, w)

		showWorkload(d.client, item.Namespace, selector, w)

		secretVolumes := make(map[string]*core.SecretVolumeSource)
		if item.Spec.DatabaseSecret != nil {
			secretVolumes["Database"] = item.Spec.DatabaseSecret
		}
		if item.Spec.CertificateSecret != nil {
			secretVolumes["Certificate"] = item.Spec.CertificateSecret
		}
		showSecret(d.client, item.Namespace, secretVolumes, w)

		specific := map[string]labels.Selector{
			"master": labels.SelectorFromSet(map[string]string{"node.role.master": "set"}),
			"client": labels.SelectorFromSet(map[string]string{"node.role.client": "set"}),
			"data":   labels.SelectorFromSet(map[string]string{"node.role.data": "set"}),
		}
		showTopology(d.client, item.Namespace, selector, specific, w)

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

type PostgresDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *PostgresDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	item, err := d.kubedb.Postgreses(namespace).Get(name, metav1.GetOptions{})
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

	return d.describePostgres(item, selector, snapshots, events)
}

func (d *PostgresDescriber) describePostgres(item *api.Postgres, selector labels.Selector, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
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

		describeArchiver(item.Spec.Archiver, w)

		describeInitialization(item.Spec.Init, w)

		describeStorage(item.Spec.StorageType, item.Spec.Storage, w)

		showWorkload(d.client, item.Namespace, selector, w)

		secretVolumes := make(map[string]*core.SecretVolumeSource)
		if item.Spec.DatabaseSecret != nil {
			secretVolumes["Database"] = item.Spec.DatabaseSecret
		}
		showSecret(d.client, item.Namespace, secretVolumes, w)

		specific := map[string]labels.Selector{
			"primary": labels.SelectorFromSet(map[string]string{"kubedb.com/role": "primary"}),
			"replica": labels.SelectorFromSet(map[string]string{"kubedb.com/role": "replica"}),
		}
		showTopology(d.client, item.Namespace, selector, specific, w)

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

type MySQLDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *MySQLDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	item, err := d.kubedb.MySQLs(namespace).Get(name, metav1.GetOptions{})
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
		events, err = d.client.Core().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeMySQL(item, selector, snapshots, events)
}

func (d *MySQLDescriber) describeMySQL(item *api.MySQL, selector labels.Selector, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
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

type MongoDBDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *MongoDBDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
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
		events, err = d.client.Core().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeMongoDB(item, selector, snapshots, events)
}

func (d *MongoDBDescriber) describeMongoDB(item *api.MongoDB, selector labels.Selector, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
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

type RedisDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *RedisDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	item, err := d.kubedb.Redises(namespace).Get(name, metav1.GetOptions{})
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
		events, err = d.client.Core().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeRedis(item, selector, snapshots, events)
}

func (d *RedisDescriber) describeRedis(item *api.Redis, selector labels.Selector, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
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

type MemcachedDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *MemcachedDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	item, err := d.kubedb.Memcacheds(namespace).Get(name, metav1.GetOptions{})
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
		events, err = d.client.Core().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeMemcached(item, selector, snapshots, events)
}

func (d *MemcachedDescriber) describeMemcached(item *api.Memcached, selector labels.Selector, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
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

		showWorkload(d.client, item.Namespace, selector, w)

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

type SnapshotDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *SnapshotDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	item, err := d.kubedb.Snapshots(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	var events *core.EventList
	if describerSettings.ShowEvents {
		events, err = d.client.Core().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeSnapshot(item, events)
}

func (d *SnapshotDescriber) describeSnapshot(item *api.Snapshot, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
		w.Write(LEVEL_0, "Name:\t%s\n", item.Name)
		w.Write(LEVEL_0, "Namespace:\t%s\n", item.Namespace)
		w.Write(LEVEL_0, "CreationTimestamp:\t%s\n", timeToString(&item.CreationTimestamp))
		if item.Status.CompletionTime != nil {
			w.Write(LEVEL_0, "CompletionTimestamp:\t%s\n", timeToString(item.Status.CompletionTime))
		}
		printLabelsMultiline(LEVEL_0, w, "Labels", item.Labels)
		printAnnotationsMultiline(LEVEL_0, w, "Annotations", item.Annotations)

		w.Write(LEVEL_0, "Status:\t%s\n", string(item.Status.Phase))
		if len(item.Status.Reason) > 0 {
			w.Write(LEVEL_0, "Reason:\t%s\n", item.Status.Reason)
		}

		w.Write(LEVEL_0, "Storage:\n")
		describeSnapshotStorage(item.Spec.Backend, w)

		secretVolumes := make(map[string]*core.SecretVolumeSource)
		if item.Spec.StorageSecretName != "" {
			secretVolumes["Database"] = &core.SecretVolumeSource{SecretName: item.Spec.StorageSecretName}
		}
		showSecret(d.client, item.Namespace, secretVolumes, w)

		if events != nil {
			DescribeEvents(events, w)
		}

		return nil
	})
}

type DormantDatabaseDescriber struct {
	client kubernetes.Interface
	kubedb cs.KubedbV1alpha1Interface
}

func (d *DormantDatabaseDescriber) Describe(namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	item, err := d.kubedb.DormantDatabases(namespace).Get(name, metav1.GetOptions{})
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
		events, err = d.client.Core().Events(item.Namespace).Search(scheme.Scheme, item)
		if err != nil {
			return "", err
		}
	}

	return d.describeDormantDatabase(item, snapshots, events)
}

func (d *DormantDatabaseDescriber) describeDormantDatabase(item *api.DormantDatabase, snapshots *api.SnapshotList, events *core.EventList) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := printersinternal.NewPrefixWriter(out)
		w.Write(LEVEL_0, "Name:\t%s\n", item.Name)
		w.Write(LEVEL_0, "Namespace:\t%s\n", item.Namespace)
		w.Write(LEVEL_0, "CreationTimestamp:\t%s\n", timeToString(&item.CreationTimestamp))
		if item.Status.PausingTime != nil {
			w.Write(LEVEL_0, "PausedTimestamp:\t%s\n", timeToString(item.Status.PausingTime))
		}
		if item.Status.WipeOutTime != nil {
			w.Write(LEVEL_0, "WipeOutTimestamp:\t%s\n", timeToString(item.Status.WipeOutTime))
		}
		printLabelsMultiline(LEVEL_0, w, "Labels", item.Labels)
		printAnnotationsMultiline(LEVEL_0, w, "Annotations", item.Annotations)

		w.Write(LEVEL_0, "Status:\t%s\n", string(item.Status.Phase))
		if len(item.Status.Reason) > 0 {
			w.Write(LEVEL_0, "Reason:\t%s\n", item.Status.Reason)
		}

		describeOrigin(item.Spec.Origin, w)

		if item.Status.Phase != api.DormantDatabasePhaseWipedOut && snapshots != nil {
			listSnapshots(snapshots, w)
		}

		if events != nil {
			DescribeEvents(events, w)
		}

		return nil
	})
}

func describeStorage(st api.StorageType, pvcSpec *core.PersistentVolumeClaimSpec, w printersinternal.PrefixWriter) {
	if st == api.StorageTypeEphemeral {
		w.Write(LEVEL_0, "  StorageType:\t%s\n", api.StorageTypeEphemeral)
	} else {
		w.Write(LEVEL_0, "  StorageType:\t%s\n", api.StorageTypeDurable)
	}
	if pvcSpec == nil {
		w.Write(LEVEL_0, "No volumes.\n")
		return
	}

	accessModes := getAccessModesAsString(pvcSpec.AccessModes)
	val, _ := pvcSpec.Resources.Requests[core.ResourceStorage]
	capacity := val.String()
	w.Write(LEVEL_0, "Volume:\n")
	if pvcSpec.StorageClassName != nil {
		w.Write(LEVEL_0, "  StorageClass:\t%s\n", *pvcSpec.StorageClassName)
	}
	w.Write(LEVEL_0, "  Capacity:\t%s\n", capacity)
	if accessModes != "" {
		w.Write(LEVEL_0, "  Access Modes:\t%s\n", accessModes)
	}
}

func describeArchiver(archiver *api.PostgresArchiverSpec, w printersinternal.PrefixWriter) {
	if archiver == nil {
		return
	}
	w.WriteLine("Archiver:")
	if archiver.Storage != nil {
		describeSnapshotStorage(*archiver.Storage, w)
	}
}

func describeInitialization(init *api.InitSpec, w printersinternal.PrefixWriter) {
	if init == nil {
		return
	}

	w.WriteLine("Init:")
	if init.ScriptSource != nil {
		w.WriteLine("  scriptSource:")
		describeVolume(init.ScriptSource.VolumeSource, w)
	}
	if init.SnapshotSource != nil {
		w.WriteLine("  snapshotSource:")
		w.Write(LEVEL_0, "    namespace:\t%s\n", init.SnapshotSource.Namespace)
		w.Write(LEVEL_0, "    name:\t%s\n", init.SnapshotSource.Name)
	}
	if init.PostgresWAL != nil {
		w.WriteLine("  postgresWAL:")
		describeSnapshotStorage(init.PostgresWAL.Backend, w)
	}
}

func describeSnapshotStorage(snapshot store.Backend, w printersinternal.PrefixWriter) {
	switch {
	case snapshot.Local != nil:
		describeVolume(snapshot.Local.VolumeSource, w)
		w.Write(LEVEL_0, "Type:\tLocal\n")
		w.Write(LEVEL_0, "path:\t%v\n", snapshot.Local.MountPath)
	case snapshot.S3 != nil:
		w.Write(LEVEL_0, "Type:\tS3\n")
		w.Write(LEVEL_0, "endpoint:\t%v\n", snapshot.S3.Endpoint)
		w.Write(LEVEL_0, "bucket:\t%v\n", snapshot.S3.Bucket)
		w.Write(LEVEL_0, "prefix:\t%v\n", snapshot.S3.Prefix)
	case snapshot.GCS != nil:
		w.Write(LEVEL_0, "Type:\tGCS\n")
		w.Write(LEVEL_0, "bucket:\t%v\n", snapshot.GCS.Bucket)
		w.Write(LEVEL_0, "prefix:\t%v\n", snapshot.GCS.Prefix)
	case snapshot.Azure != nil:
		w.Write(LEVEL_0, "Type:\tAzure\n")
		w.Write(LEVEL_0, "container:\t%v\n", snapshot.Azure.Container)
		w.Write(LEVEL_0, "prefix:\t%v\n", snapshot.Azure.Prefix)
	case snapshot.Swift != nil:
		w.Write(LEVEL_0, "Type:\tSwift\n")
		w.Write(LEVEL_0, "container:\t%v\n", snapshot.Swift.Container)
		w.Write(LEVEL_0, "prefix:\t%v\n", snapshot.Swift.Prefix)
	}
}

func describeMonitor(monitor *mona.AgentSpec, w printersinternal.PrefixWriter) {
	if monitor == nil {
		return
	}

	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "Monitoring System:\n")
	w.Write(LEVEL_0, "  Agent:\t%s\n", monitor.Agent)
	if monitor.Prometheus != nil {
		prom := monitor.Prometheus
		w.Write(LEVEL_0, "  Prometheus:\n")
		if prom.Port != 0 {
			w.Write(LEVEL_0, "    Port:\t%v\n", prom.Port)
		}
		if prom.Namespace != "" {
			w.Write(LEVEL_0, "    Namespace:\t%s\n", prom.Namespace)
		}
		if prom.Labels != nil {
			printLabelsMultiline(LEVEL_0, w, "    Labels", prom.Labels)
		}
		if prom.Interval != "" {
			w.Write(LEVEL_0, "    Interval:\t%s\n", prom.Interval)
		}

	}
}

func listSnapshots(snapshotList *api.SnapshotList, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_0, "\n")

	if len(snapshotList.Items) == 0 {
		w.Write(LEVEL_0, "No Snapshots.\n")
		return
	}

	w.Write(LEVEL_0, "Snapshots:\n")

	w.Write(LEVEL_0, "  Name\tBucket\tStartTime\tCompletionTime\tPhase\n")
	w.Write(LEVEL_0, "  ----\t------\t---------\t--------------\t-----\n")
	for _, e := range snapshotList.Items {
		location, err := e.Spec.Backend.Location()
		if err != nil {
			location = "<invalid>"
		}
		w.Write(LEVEL_0, "  %s\t%s\t%s\t%s\t%s\n",
			e.Name,
			location,
			timeToString(e.Status.StartTime),
			timeToString(e.Status.CompletionTime),
			e.Status.Phase,
		)
	}
	w.Flush()
}

func describeOrigin(origin api.Origin, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "Origin:\n")
	w.Write(LEVEL_0, "  Name:\t%s\n", origin.Name)
	w.Write(LEVEL_0, "  Namespace:\t%s\n", origin.Namespace)
	printLabelsMultiline(LEVEL_0, w, "Labels", origin.Labels)
	printAnnotationsMultiline(LEVEL_0, w, "Annotations", origin.Annotations)
}

func showWorkload(client kubernetes.Interface, namespace string, selector labels.Selector, w printersinternal.PrefixWriter) {
	pc := client.Core().Pods(namespace)
	opts := metav1.ListOptions{LabelSelector: selector.String()}

	if statefulSets, err := client.AppsV1().StatefulSets(namespace).List(opts); err == nil {
		for _, s := range statefulSets.Items {
			selector, err := metav1.LabelSelectorAsSelector(s.Spec.Selector)
			if err != nil {
				continue
			}

			running, waiting, succeeded, failed, err := getPodStatusForController(pc, selector)
			if err != nil {
				continue
			}

			describeStatefulSet(&s, running, waiting, succeeded, failed, w)
		}
	}

	if deployments, err := client.AppsV1().Deployments(namespace).List(opts); err == nil {
		for _, d := range deployments.Items {
			selector, err := metav1.LabelSelectorAsSelector(d.Spec.Selector)
			if err != nil {
				continue
			}

			running, waiting, succeeded, failed, err := getPodStatusForController(pc, selector)
			if err != nil {
				continue
			}

			describeDeployment(&d, running, waiting, succeeded, failed, w)
		}
	}

	if services, err := client.Core().Services(namespace).List(opts); err == nil {
		for _, s := range services.Items {
			endpoints, _ := client.Core().Endpoints(namespace).Get(s.Name, metav1.GetOptions{})
			describeService(&s, endpoints, w)
		}
	}
}

func showSecret(client kubernetes.Interface, namespace string, secretVolumes map[string]*core.SecretVolumeSource, w printersinternal.PrefixWriter) {
	sc := client.Core().Secrets(namespace)

	for key, sv := range secretVolumes {
		secret, err := sc.Get(sv.SecretName, metav1.GetOptions{})
		if err != nil {
			continue
		}
		describeSecret(secret, key, w)
	}
}

func showTopology(client kubernetes.Interface, namespace string, selector labels.Selector, specific map[string]labels.Selector, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "Topology:\n")
	w.Write(LEVEL_0, "  Type\tPod\tStartTime\tPhase\n")
	w.Write(LEVEL_0, "  ----\t---\t---------\t-----\n")

	pods, _ := client.Core().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: selector.String(),
	})

	for _, pod := range pods.Items {
		types := make([]string, 0)
		for key, val := range specific {
			if val.Matches(labels.Set(pod.Labels)) {
				types = append(types, key)
			}
		}
		w.Write(LEVEL_0, "  %s\t%s\t%s\t%s\n",
			strings.Join(types, "|"),
			pod.Name,
			pod.Status.StartTime,
			pod.Status.Phase,
		)
	}

	w.Flush()
}

func getAccessModesAsString(modes []core.PersistentVolumeAccessMode) string {
	modes = removeDuplicateAccessModes(modes)
	modesStr := []string{}
	if containsAccessMode(modes, core.ReadWriteOnce) {
		modesStr = append(modesStr, "RWO")
	}
	if containsAccessMode(modes, core.ReadOnlyMany) {
		modesStr = append(modesStr, "ROX")
	}
	if containsAccessMode(modes, core.ReadWriteMany) {
		modesStr = append(modesStr, "RWX")
	}
	return strings.Join(modesStr, ",")
}

func removeDuplicateAccessModes(modes []core.PersistentVolumeAccessMode) []core.PersistentVolumeAccessMode {
	var accessModes []core.PersistentVolumeAccessMode
	for _, m := range modes {
		if !containsAccessMode(accessModes, m) {
			accessModes = append(accessModes, m)
		}
	}
	return accessModes
}

func containsAccessMode(modes []core.PersistentVolumeAccessMode, mode core.PersistentVolumeAccessMode) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}
