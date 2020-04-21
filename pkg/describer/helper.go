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
	"strings"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/describe/versioned"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
)

func describeStorage(st api.StorageType, pvcSpec *core.PersistentVolumeClaimSpec, w versioned.PrefixWriter) {
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

func describeArchiver(archiver *api.PostgresArchiverSpec, w versioned.PrefixWriter) {
	if archiver == nil {
		return
	}
	w.WriteLine("Archiver:")
	if archiver.Storage != nil {
		describeSnapshotStorage(*archiver.Storage, w)
	}
}

func describeInitialization(init *api.InitSpec, w versioned.PrefixWriter) {
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

func describeSnapshotStorage(snapshot store.Backend, w versioned.PrefixWriter) {
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

func describeMonitor(monitor *mona.AgentSpec, w versioned.PrefixWriter) {
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

func listSnapshots(snapshotList *api.SnapshotList, w versioned.PrefixWriter) {
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

func describeOrigin(origin api.Origin, w versioned.PrefixWriter) {
	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "Origin:\n")
	w.Write(LEVEL_0, "  Name:\t%s\n", origin.Name)
	w.Write(LEVEL_0, "  Namespace:\t%s\n", origin.Namespace)
	printLabelsMultiline(LEVEL_0, w, "Labels", origin.Labels)
	printAnnotationsMultiline(LEVEL_0, w, "Annotations", origin.Annotations)
}

func showWorkload(client kubernetes.Interface, namespace string, selector labels.Selector, w versioned.PrefixWriter) {
	pc := client.CoreV1().Pods(namespace)
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

	if services, err := client.CoreV1().Services(namespace).List(opts); err == nil {
		for _, s := range services.Items {
			endpoints, _ := client.CoreV1().Endpoints(namespace).Get(s.Name, metav1.GetOptions{})
			describeService(&s, endpoints, w)
		}
	}
}

func showSecret(client kubernetes.Interface, namespace string, secretVolumes map[string]*core.SecretVolumeSource, w versioned.PrefixWriter) {
	sc := client.CoreV1().Secrets(namespace)

	for key, sv := range secretVolumes {
		secret, err := sc.Get(sv.SecretName, metav1.GetOptions{})
		if err != nil {
			continue
		}
		describeSecret(secret, key, w)
	}
}

func showTopology(client kubernetes.Interface, namespace string, selector labels.Selector, specific map[string]labels.Selector, w versioned.PrefixWriter) {
	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "Topology:\n")
	w.Write(LEVEL_0, "  Type\tPod\tStartTime\tPhase\n")
	w.Write(LEVEL_0, "  ----\t---\t---------\t-----\n")

	pods, _ := client.CoreV1().Pods(namespace).List(metav1.ListOptions{
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
	var modesStr []string
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
