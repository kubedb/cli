/*
Copyright 2014 The Kubernetes Authors.

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
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/golang/glog"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	"github.com/kubedb/cli/pkg/events"
	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
)

// Each level has 2 spaces for PrefixWriter
const (
	LEVEL_0 = iota
	LEVEL_1
	LEVEL_2
	LEVEL_3
)

// DescriberFn gives a way to easily override the function for unit testing if needed
var DescriberFn cmdutil.DescriberFunc = describer

// Returns a Describer for displaying the specified RESTMapping type or an error.
func describer(restClientGetter genericclioptions.RESTClientGetter, mapping *meta.RESTMapping) (printers.Describer, error) {
	clientConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	// try to get a describer
	if describer, ok := DescriberFor(mapping.GroupVersionKind.GroupKind(), clientConfig); ok {
		return describer, nil
	}
	// if this is a kind we don't have a describer for yet, go generic if possible
	if genericDescriber, genericErr := genericDescriber(restClientGetter, mapping); genericErr == nil {
		return genericDescriber, nil
	}
	// otherwise return an unregistered error
	return nil, fmt.Errorf("no description has been implemented for %s", mapping.GroupVersionKind.String())
}

// helper function to make a generic describer, or return an error
func genericDescriber(restClientGetter genericclioptions.RESTClientGetter, mapping *meta.RESTMapping) (printers.Describer, error) {
	clientConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	// used to fetch the resource
	dynamicClient, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	// used to get events for the resource
	clientSet, err := internalclientset.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	eventsClient := clientSet.Core()
	return printersinternal.GenericDescriberFor(mapping, dynamicClient, eventsClient), nil
}

func describerMap(clientConfig *rest.Config) (map[schema.GroupKind]printers.Describer, error) {
	c, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	k, err := cs.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	m := map[schema.GroupKind]printers.Describer{
		api.Kind(api.ResourceKindEtcd):          &EtcdDescriber{c, k},
		api.Kind(api.ResourceKindElasticsearch): &ElasticsearchDescriber{c, k},
		api.Kind(api.ResourceKindMemcached):     &MemcachedDescriber{c, k},
		api.Kind(api.ResourceKindMongoDB):       &MongoDBDescriber{c, k},
		api.Kind(api.ResourceKindMySQL):         &MySQLDescriber{c, k},
		api.Kind(api.ResourceKindPostgres):      &PostgresDescriber{c, k},
		api.Kind(api.ResourceKindRedis):         &RedisDescriber{c, k},
	}

	return m, nil
}

// DescriberFor returns the default describe functions for each of the standard
// Kubernetes types.
func DescriberFor(kind schema.GroupKind, clientConfig *rest.Config) (printers.Describer, bool) {
	describers, err := describerMap(clientConfig)
	if err != nil {
		glog.V(1).Info(err)
		return nil, false
	}

	f, ok := describers[kind]
	return f, ok
}

func describeStatefulSet(ps *appsv1.StatefulSet, running, waiting, succeeded, failed int, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "StatefulSet:\t\n")
	w.Write(LEVEL_1, "Name:\t%s\n", ps.ObjectMeta.Name)
	w.Write(LEVEL_1, "CreationTimestamp:\t%s\n", timeToString(&ps.CreationTimestamp))
	printLabelsMultiline(LEVEL_1, w, "Labels", ps.Labels)
	printAnnotationsMultiline(LEVEL_1, w, "Annotations", ps.Annotations)
	w.Write(LEVEL_1, "Replicas:\t%d desired | %d total\n", ps.Spec.Replicas, ps.Status.Replicas)
	w.Write(LEVEL_1, "Pods Status:\t%d Running / %d Waiting / %d Succeeded / %d Failed\n", running, waiting, succeeded, failed)
}

func describeDeployment(d *appsv1.Deployment, running, waiting, succeeded, failed int, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "Deployment:\t\n")
	w.Write(LEVEL_1, "Name:\t%s\n", d.ObjectMeta.Name)
	w.Write(LEVEL_1, "CreationTimestamp:\t%s\n", timeToString(&d.CreationTimestamp))
	printLabelsMultiline(LEVEL_1, w, "Labels", d.Labels)
	printAnnotationsMultiline(LEVEL_1, w, "Annotations", d.Annotations)
	w.Write(LEVEL_1, "Replicas:\t%d desired | %d updated | %d total | %d available | %d unavailable\n", *(d.Spec.Replicas), d.Status.UpdatedReplicas, d.Status.Replicas, d.Status.AvailableReplicas, d.Status.UnavailableReplicas)
	w.Write(LEVEL_1, "Pods Status:\t%d Running / %d Waiting / %d Succeeded / %d Failed\n", running, waiting, succeeded, failed)
}

func getPodStatusForController(c coreclient.PodInterface, selector labels.Selector) (running, waiting, succeeded, failed int, err error) {
	options := metav1.ListOptions{LabelSelector: selector.String()}
	rcPods, err := c.List(options)
	if err != nil {
		return
	}
	for _, pod := range rcPods.Items {
		switch pod.Status.Phase {
		case core.PodRunning:
			running++
		case core.PodPending:
			waiting++
		case core.PodSucceeded:
			succeeded++
		case core.PodFailed:
			failed++
		}
	}
	return
}

func buildIngressString(ingress []core.LoadBalancerIngress) string {
	var buffer bytes.Buffer

	for i := range ingress {
		if i != 0 {
			buffer.WriteString(", ")
		}
		if ingress[i].IP != "" {
			buffer.WriteString(ingress[i].IP)
		} else {
			buffer.WriteString(ingress[i].Hostname)
		}
	}
	return buffer.String()
}

func describeService(service *core.Service, endpoints *core.Endpoints, w printersinternal.PrefixWriter) {
	if endpoints == nil {
		endpoints = &core.Endpoints{}
	}
	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "Service:\t\n")
	w.Write(LEVEL_1, "Name:\t%s\n", service.Name)
	printLabelsMultiline(LEVEL_1, w, "Labels", service.Labels)
	printAnnotationsMultiline(LEVEL_1, w, "Annotations", service.Annotations)
	w.Write(LEVEL_1, "Type:\t%s\n", service.Spec.Type)
	w.Write(LEVEL_1, "IP:\t%s\n", service.Spec.ClusterIP)
	if len(service.Spec.ExternalIPs) > 0 {
		w.Write(LEVEL_1, "External IPs:\t%v\n", strings.Join(service.Spec.ExternalIPs, ","))
	}
	if service.Spec.LoadBalancerIP != "" {
		w.Write(LEVEL_1, "IP:\t%s\n", service.Spec.LoadBalancerIP)
	}
	if service.Spec.ExternalName != "" {
		w.Write(LEVEL_1, "External Name:\t%s\n", service.Spec.ExternalName)
	}
	if len(service.Status.LoadBalancer.Ingress) > 0 {
		list := buildIngressString(service.Status.LoadBalancer.Ingress)
		w.Write(LEVEL_1, "LoadBalancer Ingress:\t%s\n", list)
	}
	for i := range service.Spec.Ports {
		sp := &service.Spec.Ports[i]

		name := sp.Name
		if name == "" {
			name = "<unset>"
		}
		w.Write(LEVEL_1, "Port:\t%s\t%d/%s\n", name, sp.Port, sp.Protocol)
		if sp.TargetPort.Type == intstr.Type(intstr.Int) {
			w.Write(LEVEL_1, "TargetPort:\t%d/%s\n", sp.TargetPort.IntVal, sp.Protocol)
		} else {
			w.Write(LEVEL_1, "TargetPort:\t%s/%s\n", sp.TargetPort.StrVal, sp.Protocol)
		}
		if sp.NodePort != 0 {
			w.Write(LEVEL_1, "NodePort:\t%s\t%d/%s\n", name, sp.NodePort, sp.Protocol)
		}
		w.Write(LEVEL_1, "Endpoints:\t%s\n", formatEndpoints(endpoints, sets.NewString(sp.Name)))
	}
}

// describeSecret generates information about a secret
func describeSecret(secret *core.Secret, prefix string, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_0, "\n")
	if prefix == "" {
		w.Write(LEVEL_0, "Secret:\n")
	} else {
		w.Write(LEVEL_0, "%s Secret:\n", prefix)
	}
	w.Write(LEVEL_1, "Name:\t%s\n", secret.Name)
	printLabelsMultiline(LEVEL_1, w, "Labels", secret.Labels)
	skipAnnotations := sets.NewString(meta_util.LastAppliedConfigAnnotation)
	printAnnotationsMultilineWithFilter(LEVEL_1, w, "Annotations", secret.Annotations, skipAnnotations)

	w.Write(LEVEL_1, "\nType:\t%s\n", secret.Type)

	w.Write(LEVEL_1, "\nData\n====\n")
	for k, v := range secret.Data {
		switch {
		case k == core.ServiceAccountTokenKey && secret.Type == core.SecretTypeServiceAccountToken:
			w.Write(LEVEL_1, "%s:\t%s\n", k, string(v))
		default:
			w.Write(LEVEL_1, "%s:\t%d bytes\n", k, len(v))
		}
	}
}

func describeVolume(volume core.VolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_0, "Volume:\n")
	switch {
	case volume.HostPath != nil:
		printHostPathVolumeSource(volume.HostPath, w)
	case volume.EmptyDir != nil:
		printEmptyDirVolumeSource(volume.EmptyDir, w)
	case volume.GCEPersistentDisk != nil:
		printGCEPersistentDiskVolumeSource(volume.GCEPersistentDisk, w)
	case volume.AWSElasticBlockStore != nil:
		printAWSElasticBlockStoreVolumeSource(volume.AWSElasticBlockStore, w)
	case volume.GitRepo != nil:
		printGitRepoVolumeSource(volume.GitRepo, w)
	case volume.Secret != nil:
		printSecretVolumeSource(volume.Secret, w)
	case volume.ConfigMap != nil:
		printConfigMapVolumeSource(volume.ConfigMap, w)
	case volume.NFS != nil:
		printNFSVolumeSource(volume.NFS, w)
	case volume.ISCSI != nil:
		printISCSIVolumeSource(volume.ISCSI, w)
	case volume.Glusterfs != nil:
		printGlusterfsVolumeSource(volume.Glusterfs, w)
	case volume.PersistentVolumeClaim != nil:
		printPersistentVolumeClaimVolumeSource(volume.PersistentVolumeClaim, w)
	case volume.RBD != nil:
		printRBDVolumeSource(volume.RBD, w)
	case volume.Quobyte != nil:
		printQuobyteVolumeSource(volume.Quobyte, w)
	case volume.DownwardAPI != nil:
		printDownwardAPIVolumeSource(volume.DownwardAPI, w)
	case volume.AzureDisk != nil:
		printAzureDiskVolumeSource(volume.AzureDisk, w)
	case volume.VsphereVolume != nil:
		printVsphereVolumeSource(volume.VsphereVolume, w)
	case volume.Cinder != nil:
		printCinderVolumeSource(volume.Cinder, w)
	case volume.PhotonPersistentDisk != nil:
		printPhotonPersistentDiskVolumeSource(volume.PhotonPersistentDisk, w)
	case volume.PortworxVolume != nil:
		printPortworxVolumeSource(volume.PortworxVolume, w)
	case volume.ScaleIO != nil:
		printScaleIOVolumeSource(volume.ScaleIO, w)
	case volume.CephFS != nil:
		printCephFSVolumeSource(volume.CephFS, w)
	case volume.StorageOS != nil:
		printStorageOSVolumeSource(volume.StorageOS, w)
	case volume.FC != nil:
		printFCVolumeSource(volume.FC, w)
	case volume.AzureFile != nil:
		printAzureFileVolumeSource(volume.AzureFile, w)
	case volume.FlexVolume != nil:
		printFlexVolumeSource(volume.FlexVolume, w)
	case volume.Flocker != nil:
		printFlockerVolumeSource(volume.Flocker, w)
	default:
		w.Write(LEVEL_1, "<unknown>\n")
	}
}

func printHostPathVolumeSource(hostPath *core.HostPathVolumeSource, w printersinternal.PrefixWriter) {
	hostPathType := "<none>"
	if hostPath.Type != nil {
		hostPathType = string(*hostPath.Type)
	}
	w.Write(LEVEL_2, "Type:\tHostPath (bare host directory volume)\n"+
		"    Path:\t%v\n"+
		"    HostPathType:\t%v\n",
		hostPath.Path, hostPathType)
}

func printEmptyDirVolumeSource(emptyDir *core.EmptyDirVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tEmptyDir (a temporary directory that shares a pod's lifetime)\n"+
		"    Medium:\t%v\n", emptyDir.Medium)
}

func printGCEPersistentDiskVolumeSource(gce *core.GCEPersistentDiskVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tGCEPersistentDisk (a Persistent Disk resource in Google Compute Engine)\n"+
		"    PDName:\t%v\n"+
		"    FSType:\t%v\n"+
		"    Partition:\t%v\n"+
		"    ReadOnly:\t%v\n",
		gce.PDName, gce.FSType, gce.Partition, gce.ReadOnly)
}

func printAWSElasticBlockStoreVolumeSource(aws *core.AWSElasticBlockStoreVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tAWSElasticBlockStore (a Persistent Disk resource in AWS)\n"+
		"    VolumeID:\t%v\n"+
		"    FSType:\t%v\n"+
		"    Partition:\t%v\n"+
		"    ReadOnly:\t%v\n",
		aws.VolumeID, aws.FSType, aws.Partition, aws.ReadOnly)
}

func printGitRepoVolumeSource(git *core.GitRepoVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tGitRepo (a volume that is pulled from git when the pod is created)\n"+
		"    Repository:\t%v\n"+
		"    Revision:\t%v\n",
		git.Repository, git.Revision)
}

func printSecretVolumeSource(secret *core.SecretVolumeSource, w printersinternal.PrefixWriter) {
	optional := secret.Optional != nil && *secret.Optional
	w.Write(LEVEL_2, "Type:\tSecret (a volume populated by a Secret)\n"+
		"    SecretName:\t%v\n"+
		"    Optional:\t%v\n",
		secret.SecretName, optional)
}

func printConfigMapVolumeSource(configMap *core.ConfigMapVolumeSource, w printersinternal.PrefixWriter) {
	optional := configMap.Optional != nil && *configMap.Optional
	w.Write(LEVEL_2, "Type:\tConfigMap (a volume populated by a ConfigMap)\n"+
		"    Name:\t%v\n"+
		"    Optional:\t%v\n",
		configMap.Name, optional)
}

func printNFSVolumeSource(nfs *core.NFSVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tNFS (an NFS mount that lasts the lifetime of a pod)\n"+
		"    Server:\t%v\n"+
		"    Path:\t%v\n"+
		"    ReadOnly:\t%v\n",
		nfs.Server, nfs.Path, nfs.ReadOnly)
}

func printQuobyteVolumeSource(quobyte *core.QuobyteVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tQuobyte (a Quobyte mount on the host that shares a pod's lifetime)\n"+
		"    Registry:\t%v\n"+
		"    Volume:\t%v\n"+
		"    ReadOnly:\t%v\n",
		quobyte.Registry, quobyte.Volume, quobyte.ReadOnly)
}

func printPortworxVolumeSource(pwxVolume *core.PortworxVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tPortworxVolume (a Portworx Volume resource)\n"+
		"    VolumeID:\t%v\n",
		pwxVolume.VolumeID)
}

func printISCSIVolumeSource(iscsi *core.ISCSIVolumeSource, w printersinternal.PrefixWriter) {
	initiator := "<none>"
	if iscsi.InitiatorName != nil {
		initiator = *iscsi.InitiatorName
	}
	w.Write(LEVEL_2, "Type:\tISCSI (an ISCSI Disk resource that is attached to a kubelet's host machine and then exposed to the pod)\n"+
		"    TargetPortal:\t%v\n"+
		"    IQN:\t%v\n"+
		"    Lun:\t%v\n"+
		"    ISCSIInterface\t%v\n"+
		"    FSType:\t%v\n"+
		"    ReadOnly:\t%v\n"+
		"    Portals:\t%v\n"+
		"    DiscoveryCHAPAuth:\t%v\n"+
		"    SessionCHAPAuth:\t%v\n"+
		"    SecretRef:\t%v\n"+
		"    InitiatorName:\t%v\n",
		iscsi.TargetPortal, iscsi.IQN, iscsi.Lun, iscsi.ISCSIInterface, iscsi.FSType, iscsi.ReadOnly, iscsi.Portals, iscsi.DiscoveryCHAPAuth, iscsi.SessionCHAPAuth, iscsi.SecretRef, initiator)
}

func printGlusterfsVolumeSource(glusterfs *core.GlusterfsVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tGlusterfs (a Glusterfs mount on the host that shares a pod's lifetime)\n"+
		"    EndpointsName:\t%v\n"+
		"    Path:\t%v\n"+
		"    ReadOnly:\t%v\n",
		glusterfs.EndpointsName, glusterfs.Path, glusterfs.ReadOnly)
}

func printPersistentVolumeClaimVolumeSource(claim *core.PersistentVolumeClaimVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tPersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)\n"+
		"    ClaimName:\t%v\n"+
		"    ReadOnly:\t%v\n",
		claim.ClaimName, claim.ReadOnly)
}

func printRBDVolumeSource(rbd *core.RBDVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tRBD (a Rados Block Device mount on the host that shares a pod's lifetime)\n"+
		"    CephMonitors:\t%v\n"+
		"    RBDImage:\t%v\n"+
		"    FSType:\t%v\n"+
		"    RBDPool:\t%v\n"+
		"    RadosUser:\t%v\n"+
		"    Keyring:\t%v\n"+
		"    SecretRef:\t%v\n"+
		"    ReadOnly:\t%v\n",
		rbd.CephMonitors, rbd.RBDImage, rbd.FSType, rbd.RBDPool, rbd.RadosUser, rbd.Keyring, rbd.SecretRef, rbd.ReadOnly)
}

func printDownwardAPIVolumeSource(d *core.DownwardAPIVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tDownwardAPI (a volume populated by information about the pod)\n    Items:\n")
	for _, mapping := range d.Items {
		if mapping.FieldRef != nil {
			w.Write(LEVEL_3, "%v -> %v\n", mapping.FieldRef.FieldPath, mapping.Path)
		}
		if mapping.ResourceFieldRef != nil {
			w.Write(LEVEL_3, "%v -> %v\n", mapping.ResourceFieldRef.Resource, mapping.Path)
		}
	}
}

func printAzureDiskVolumeSource(d *core.AzureDiskVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tAzureDisk (an Azure Data Disk mount on the host and bind mount to the pod)\n"+
		"    DiskName:\t%v\n"+
		"    DiskURI:\t%v\n"+
		"    Kind: \t%v\n"+
		"    FSType:\t%v\n"+
		"    CachingMode:\t%v\n"+
		"    ReadOnly:\t%v\n",
		d.DiskName, d.DataDiskURI, *d.Kind, *d.FSType, *d.CachingMode, *d.ReadOnly)
}

func printVsphereVolumeSource(vsphere *core.VsphereVirtualDiskVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tvSphereVolume (a Persistent Disk resource in vSphere)\n"+
		"    VolumePath:\t%v\n"+
		"    FSType:\t%v\n"+
		"    StoragePolicyName:\t%v\n",
		vsphere.VolumePath, vsphere.FSType, vsphere.StoragePolicyName)
}

func printPhotonPersistentDiskVolumeSource(photon *core.PhotonPersistentDiskVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tPhotonPersistentDisk (a Persistent Disk resource in photon platform)\n"+
		"    PdID:\t%v\n"+
		"    FSType:\t%v\n",
		photon.PdID, photon.FSType)
}

func printCinderVolumeSource(cinder *core.CinderVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tCinder (a Persistent Disk resource in OpenStack)\n"+
		"    VolumeID:\t%v\n"+
		"    FSType:\t%v\n"+
		"    ReadOnly:\t%v\n",
		"    SecretRef:\t%v\n"+
			cinder.VolumeID, cinder.FSType, cinder.ReadOnly, cinder.SecretRef)
}

func printScaleIOVolumeSource(sio *core.ScaleIOVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tScaleIO (a persistent volume backed by a block device in ScaleIO)\n"+
		"    Gateway:\t%v\n"+
		"    System:\t%v\n"+
		"    Protection Domain:\t%v\n"+
		"    Storage Pool:\t%v\n"+
		"    Storage Mode:\t%v\n"+
		"    VolumeName:\t%v\n"+
		"    FSType:\t%v\n"+
		"    ReadOnly:\t%v\n",
		sio.Gateway, sio.System, sio.ProtectionDomain, sio.StoragePool, sio.StorageMode, sio.VolumeName, sio.FSType, sio.ReadOnly)
}

func printCephFSVolumeSource(cephfs *core.CephFSVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tCephFS (a CephFS mount on the host that shares a pod's lifetime)\n"+
		"    Monitors:\t%v\n"+
		"    Path:\t%v\n"+
		"    User:\t%v\n"+
		"    SecretFile:\t%v\n"+
		"    SecretRef:\t%v\n"+
		"    ReadOnly:\t%v\n",
		cephfs.Monitors, cephfs.Path, cephfs.User, cephfs.SecretFile, cephfs.SecretRef, cephfs.ReadOnly)
}

func printStorageOSVolumeSource(storageos *core.StorageOSVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tStorageOS (a StorageOS Persistent Disk resource)\n"+
		"    VolumeName:\t%v\n"+
		"    VolumeNamespace:\t%v\n"+
		"    FSType:\t%v\n"+
		"    ReadOnly:\t%v\n",
		storageos.VolumeName, storageos.VolumeNamespace, storageos.FSType, storageos.ReadOnly)
}

func printFCVolumeSource(fc *core.FCVolumeSource, w printersinternal.PrefixWriter) {
	lun := "<none>"
	if fc.Lun != nil {
		lun = strconv.Itoa(int(*fc.Lun))
	}
	w.Write(LEVEL_2, "Type:\tFC (a Fibre Channel disk)\n"+
		"    TargetWWNs:\t%v\n"+
		"    LUN:\t%v\n"+
		"    FSType:\t%v\n"+
		"    ReadOnly:\t%v\n",
		strings.Join(fc.TargetWWNs, ", "), lun, fc.FSType, fc.ReadOnly)
}

func printAzureFileVolumeSource(azureFile *core.AzureFileVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tAzureFile (an Azure File Service mount on the host and bind mount to the pod)\n"+
		"    SecretName:\t%v\n"+
		"    ShareName:\t%v\n"+
		"    ReadOnly:\t%v\n",
		azureFile.SecretName, azureFile.ShareName, azureFile.ReadOnly)
}

func printFlexVolumeSource(flex *core.FlexVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tFlexVolume (a generic volume resource that is provisioned/attached using an exec based plugin)\n"+
		"    Driver:\t%v\n"+
		"    FSType:\t%v\n"+
		"    SecretRef:\t%v\n"+
		"    ReadOnly:\t%v\n"+
		"    Options:\t%v\n",
		flex.Driver, flex.FSType, flex.SecretRef, flex.ReadOnly, flex.Options)
}

func printFlockerVolumeSource(flocker *core.FlockerVolumeSource, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tFlocker (a Flocker volume mounted by the Flocker agent)\n"+
		"    DatasetName:\t%v\n"+
		"    DatasetUUID:\t%v\n",
		flocker.DatasetName, flocker.DatasetUUID)
}

func DescribeEvents(el *core.EventList, w printersinternal.PrefixWriter) {
	w.Write(LEVEL_0, "\n")
	if len(el.Items) == 0 {
		w.Write(LEVEL_0, "Events:\t<none>\n")
		return
	}
	w.Flush()
	sort.Sort(events.SortableEvents(el.Items))
	w.Write(LEVEL_0, "Events:\n  Type\tReason\tAge\tFrom\tMessage\n")
	w.Write(LEVEL_1, "----\t------\t----\t----\t-------\n")
	for _, e := range el.Items {
		var interval string
		if e.Count > 1 {
			interval = fmt.Sprintf("%s (x%d over %s)", translateTimestamp(e.LastTimestamp), e.Count, translateTimestamp(e.FirstTimestamp))
		} else {
			interval = translateTimestamp(e.FirstTimestamp)
		}
		w.Write(LEVEL_1, "%v\t%v\t%s\t%v\t%v\n",
			e.Type,
			e.Reason,
			interval,
			formatEventSource(e.Source),
			strings.TrimSpace(e.Message),
		)
	}
}

var maxAnnotationLen = 200

// printLabelsMultiline prints multiple labels with a proper alignment.
func printLabelsMultiline(level int, w printersinternal.PrefixWriter, title string, labels map[string]string) {
	printLabelsMultilineWithIndent(level, w, "", title, "\t", labels, sets.NewString())
}

// printLabelsMultiline prints multiple labels with a user-defined alignment.
func printLabelsMultilineWithIndent(level int, w printersinternal.PrefixWriter, initialIndent, title, innerIndent string, labels map[string]string, skip sets.String) {
	w.Write(level, "%s%s:%s", initialIndent, title, innerIndent)

	if labels == nil || len(labels) == 0 {
		w.WriteLine("<none>")
		return
	}

	// to print labels in the sorted order
	keys := make([]string, 0, len(labels))
	for key := range labels {
		if skip.Has(key) {
			continue
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		w.WriteLine("<none>")
		return
	}
	sort.Strings(keys)

	for i, key := range keys {
		if i != 0 {
			w.Write(level, "%s", initialIndent)
			w.Write(level, "%s", innerIndent)
		}
		w.Write(level, "%s=%s\n", key, labels[key])
		i++
	}
}

func tabbedString(f func(io.Writer) error) (string, error) {
	out := new(tabwriter.Writer)
	buf := &bytes.Buffer{}
	out.Init(buf, 0, 8, 2, ' ', 0)

	err := f(out)
	if err != nil {
		return "", err
	}

	out.Flush()
	str := string(buf.String())
	return str, nil
}

// printAnnotationsMultilineWithFilter prints filtered multiple annotations with a proper alignment.
func printAnnotationsMultilineWithFilter(level int, w printersinternal.PrefixWriter, title string, annotations map[string]string, skip sets.String) {
	printAnnotationsMultilineWithIndent(level, w, "", title, "\t", annotations, skip)
}

// printAnnotationsMultiline prints multiple annotations with a proper alignment.
func printAnnotationsMultiline(level int, w printersinternal.PrefixWriter, title string, annotations map[string]string) {
	printAnnotationsMultilineWithIndent(level, w, "", title, "\t", annotations, sets.NewString())
}

// printAnnotationsMultilineWithIndent prints multiple annotations with a user-defined alignment.
// If annotation string is too long, we omit chars more than 200 length.
func printAnnotationsMultilineWithIndent(level int, w printersinternal.PrefixWriter, initialIndent, title, innerIndent string, annotations map[string]string, skip sets.String) {

	w.Write(level, "%s%s:%s", initialIndent, title, innerIndent)

	if len(annotations) == 0 {
		w.WriteLine("<none>")
		return
	}

	// to print labels in the sorted order
	keys := make([]string, 0, len(annotations))
	for key := range annotations {
		if skip.Has(key) {
			continue
		}
		keys = append(keys, key)
	}
	if len(annotations) == 0 {
		w.WriteLine("<none>")
		return
	}
	sort.Strings(keys)

	for i, key := range keys {
		if i != 0 {
			w.Write(level, initialIndent)
			w.Write(level, innerIndent)
		}
		line := fmt.Sprintf("%s=%s", key, annotations[key])
		if len(line) > maxAnnotationLen {
			w.Write(level, "%s...\n", line[:maxAnnotationLen])
		} else {
			w.Write(level, "%s\n", line)
		}
		i++
	}
}
