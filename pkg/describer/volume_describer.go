package describer

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	core "k8s.io/api/core/v1"
)

const (
	LEVEL_0 = iota
	LEVEL_1
	LEVEL_2
	LEVEL_3
)

// PrefixWriter can write text at various indentation levels.
type PrefixWriter interface {
	// Write writes text with the specified indentation level.
	Write(level int, format string, a ...interface{})
	// WriteLine writes an entire line with no indentation level.
	WriteLine(a ...interface{})
	// Flush forces indendation to be reset.
	Flush()
}

// prefixWriter implements PrefixWriter
type prefixWriter struct {
	out io.Writer
}

// NewPrefixWriter creates a new PrefixWriter.
func newPrefixWriter(out io.Writer) PrefixWriter {
	return &prefixWriter{out: out}
}

func (pw *prefixWriter) Write(level int, format string, a ...interface{}) {
	levelSpace := "  "
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += levelSpace
	}
	fmt.Fprintf(pw.out, prefix+format, a...)
}

func (pw *prefixWriter) WriteLine(a ...interface{}) {
	fmt.Fprintln(pw.out, a...)
}

func (pw *prefixWriter) Flush() {
	if f, ok := pw.out.(flusher); ok {
		f.Flush()
	}
}

func describeVolumes(volumeSource core.VolumeSource, out io.Writer) {
	w := newPrefixWriter(out)

	switch {
	case volumeSource.HostPath != nil:
		printHostPathVolumeSource(volumeSource.HostPath, w)
	case volumeSource.EmptyDir != nil:
		printEmptyDirVolumeSource(volumeSource.EmptyDir, w)
	case volumeSource.GCEPersistentDisk != nil:
		printGCEPersistentDiskVolumeSource(volumeSource.GCEPersistentDisk, w)
	case volumeSource.AWSElasticBlockStore != nil:
		printAWSElasticBlockStoreVolumeSource(volumeSource.AWSElasticBlockStore, w)
	case volumeSource.GitRepo != nil:
		printGitRepoVolumeSource(volumeSource.GitRepo, w)
	case volumeSource.Secret != nil:
		printSecretVolumeSource(volumeSource.Secret, w)
	case volumeSource.ConfigMap != nil:
		printConfigMapVolumeSource(volumeSource.ConfigMap, w)
	case volumeSource.NFS != nil:
		printNFSVolumeSource(volumeSource.NFS, w)
	case volumeSource.ISCSI != nil:
		printISCSIVolumeSource(volumeSource.ISCSI, w)
	case volumeSource.Glusterfs != nil:
		printGlusterfsVolumeSource(volumeSource.Glusterfs, w)
	case volumeSource.PersistentVolumeClaim != nil:
		printPersistentVolumeClaimVolumeSource(volumeSource.PersistentVolumeClaim, w)
	case volumeSource.RBD != nil:
		printRBDVolumeSource(volumeSource.RBD, w)
	case volumeSource.Quobyte != nil:
		printQuobyteVolumeSource(volumeSource.Quobyte, w)
	case volumeSource.DownwardAPI != nil:
		printDownwardAPIVolumeSource(volumeSource.DownwardAPI, w)
	case volumeSource.AzureDisk != nil:
		printAzureDiskVolumeSource(volumeSource.AzureDisk, w)
	case volumeSource.VsphereVolume != nil:
		printVsphereVolumeSource(volumeSource.VsphereVolume, w)
	case volumeSource.Cinder != nil:
		printCinderVolumeSource(volumeSource.Cinder, w)
	case volumeSource.PhotonPersistentDisk != nil:
		printPhotonPersistentDiskVolumeSource(volumeSource.PhotonPersistentDisk, w)
	case volumeSource.PortworxVolume != nil:
		printPortworxVolumeSource(volumeSource.PortworxVolume, w)
	case volumeSource.ScaleIO != nil:
		printScaleIOVolumeSource(volumeSource.ScaleIO, w)
	case volumeSource.CephFS != nil:
		printCephFSVolumeSource(volumeSource.CephFS, w)
	case volumeSource.StorageOS != nil:
		printStorageOSVolumeSource(volumeSource.StorageOS, w)
	case volumeSource.FC != nil:
		printFCVolumeSource(volumeSource.FC, w)
	case volumeSource.AzureFile != nil:
		printAzureFileVolumeSource(volumeSource.AzureFile, w)
	case volumeSource.FlexVolume != nil:
		printFlexVolumeSource(volumeSource.FlexVolume, w)
	case volumeSource.Flocker != nil:
		printFlockerVolumeSource(volumeSource.Flocker, w)
	default:
		w.Write(LEVEL_1, "<unknown>\n")
	}
	w.Flush()
}

func getSpace(spaceLevel int) string {
	levelSpace := "  "
	prefix := ""
	for i := 0; i < spaceLevel; i++ {
		prefix += levelSpace
	}

	return prefix
}

func describeSnapshotStorage(snapshot api.SnapshotStorageSpec, out io.Writer, spaceLevel int) {
	w := newPrefixWriter(out)

	switch {
	case snapshot.Local != nil:
		describeVolumes(snapshot.Local.VolumeSource, out)
		w.Write(spaceLevel, "Type:\tLocal\n"+
			getSpace(spaceLevel)+"path:\t%v\n", snapshot.Local.Path)
	case snapshot.S3 != nil:
		w.Write(spaceLevel, "Type:\tS3\n"+
			getSpace(spaceLevel)+"endpoint:\t%v\n"+
			getSpace(spaceLevel)+"bucket:\t%v\n"+
			getSpace(spaceLevel)+"prefix:\t%v\n",
			snapshot.S3.Endpoint,
			snapshot.S3.Bucket,
			snapshot.S3.Prefix)
	case snapshot.GCS != nil:
		w.Write(spaceLevel, "Type:\tGCS\n"+
			getSpace(spaceLevel)+"bucket:\t%v\n"+
			getSpace(spaceLevel)+"prefix:\t%v\n",
			snapshot.GCS.Bucket,
			snapshot.GCS.Prefix)
	case snapshot.Azure != nil:
		w.Write(spaceLevel, "Type:\tAzure\n"+
			getSpace(spaceLevel)+"container:\t%v\n"+
			getSpace(spaceLevel)+"prefix:\t%v\n",
			snapshot.Azure.Container,
			snapshot.Azure.Prefix)
	case snapshot.Swift != nil:
		w.Write(spaceLevel, "Type:\tSwift\n"+
			getSpace(spaceLevel)+"container:\t%v\n"+
			getSpace(spaceLevel)+"prefix:\t%v\n",
			snapshot.Swift.Container,
			snapshot.Swift.Prefix)
	}
	w.Flush()
}

func printHostPathVolumeSource(hostPath *core.HostPathVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tHostPath (bare host directory volume)\n"+
		"    Path:\t%v\n", hostPath.Path)
}

func printEmptyDirVolumeSource(emptyDir *core.EmptyDirVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tEmptyDir (a temporary directory that shares a pod's lifetime)\n"+
		"    Medium:\t%v\n", emptyDir.Medium)
}

func printGCEPersistentDiskVolumeSource(gce *core.GCEPersistentDiskVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tGCEPersistentDisk (a Persistent Disk resource in Google Compute Engine)\n"+
		"    PDName:\t%v\n"+
		"    FSType:\t%v\n"+
		"    Partition:\t%v\n"+
		"    ReadOnly:\t%v\n",
		gce.PDName, gce.FSType, gce.Partition, gce.ReadOnly)
}

func printAWSElasticBlockStoreVolumeSource(aws *core.AWSElasticBlockStoreVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tAWSElasticBlockStore (a Persistent Disk resource in AWS)\n"+
		"    VolumeID:\t%v\n"+
		"    FSType:\t%v\n"+
		"    Partition:\t%v\n"+
		"    ReadOnly:\t%v\n",
		aws.VolumeID, aws.FSType, aws.Partition, aws.ReadOnly)
}

func printGitRepoVolumeSource(git *core.GitRepoVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tGitRepo (a volume that is pulled from git when the pod is created)\n"+
		"    Repository:\t%v\n"+
		"    Directory:\t%v\n"+
		"    Revision:\t%v\n",
		git.Repository, git.Directory, git.Revision)
}

func printSecretVolumeSource(secret *core.SecretVolumeSource, w PrefixWriter) {
	optional := secret.Optional != nil && *secret.Optional
	w.Write(LEVEL_2, "Type:\tSecret (a volume populated by a Secret)\n"+
		"    SecretName:\t%v\n"+
		"    Optional:\t%v\n",
		secret.SecretName, optional)
}

func printConfigMapVolumeSource(configMap *core.ConfigMapVolumeSource, w PrefixWriter) {
	optional := configMap.Optional != nil && *configMap.Optional
	w.Write(LEVEL_2, "Type:\tConfigMap (a volume populated by a ConfigMap)\n"+
		"    Name:\t%v\n"+
		"    Optional:\t%v\n",
		configMap.Name, optional)
}

func printNFSVolumeSource(nfs *core.NFSVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tNFS (an NFS mount that lasts the lifetime of a pod)\n"+
		"    Server:\t%v\n"+
		"    Path:\t%v\n"+
		"    ReadOnly:\t%v\n",
		nfs.Server, nfs.Path, nfs.ReadOnly)
}

func printQuobyteVolumeSource(quobyte *core.QuobyteVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tQuobyte (a Quobyte mount on the host that shares a pod's lifetime)\n"+
		"    Registry:\t%v\n"+
		"    Volume:\t%v\n"+
		"    ReadOnly:\t%v\n",
		quobyte.Registry, quobyte.Volume, quobyte.ReadOnly)
}

func printPortworxVolumeSource(pwxVolume *core.PortworxVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tPortworxVolume (a Portworx Volume resource)\n"+
		"    VolumeID:\t%v\n",
		pwxVolume.VolumeID)
}

func printISCSIVolumeSource(iscsi *core.ISCSIVolumeSource, w PrefixWriter) {
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

func printGlusterfsVolumeSource(glusterfs *core.GlusterfsVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tGlusterfs (a Glusterfs mount on the host that shares a pod's lifetime)\n"+
		"    EndpointsName:\t%v\n"+
		"    Path:\t%v\n"+
		"    ReadOnly:\t%v\n",
		glusterfs.EndpointsName, glusterfs.Path, glusterfs.ReadOnly)
}

func printPersistentVolumeClaimVolumeSource(claim *core.PersistentVolumeClaimVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tPersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)\n"+
		"    ClaimName:\t%v\n"+
		"    ReadOnly:\t%v\n",
		claim.ClaimName, claim.ReadOnly)
}

func printRBDVolumeSource(rbd *core.RBDVolumeSource, w PrefixWriter) {
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

func printDownwardAPIVolumeSource(d *core.DownwardAPIVolumeSource, w PrefixWriter) {
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

func printAzureDiskVolumeSource(d *core.AzureDiskVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tAzureDisk (an Azure Data Disk mount on the host and bind mount to the pod)\n"+
		"    DiskName:\t%v\n"+
		"    DiskURI:\t%v\n"+
		"    Kind: \t%v\n"+
		"    FSType:\t%v\n"+
		"    CachingMode:\t%v\n"+
		"    ReadOnly:\t%v\n",
		d.DiskName, d.DataDiskURI, *d.Kind, *d.FSType, *d.CachingMode, *d.ReadOnly)
}

func printVsphereVolumeSource(vsphere *core.VsphereVirtualDiskVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tvSphereVolume (a Persistent Disk resource in vSphere)\n"+
		"    VolumePath:\t%v\n"+
		"    FSType:\t%v\n",
		"    StoragePolicyName:\t%v\n",
		vsphere.VolumePath, vsphere.FSType, vsphere.StoragePolicyName)
}

func printPhotonPersistentDiskVolumeSource(photon *core.PhotonPersistentDiskVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tPhotonPersistentDisk (a Persistent Disk resource in photon platform)\n"+
		"    PdID:\t%v\n"+
		"    FSType:\t%v\n",
		photon.PdID, photon.FSType)
}

func printCinderVolumeSource(cinder *core.CinderVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tCinder (a Persistent Disk resource in OpenStack)\n"+
		"    VolumeID:\t%v\n"+
		"    FSType:\t%v\n"+
		"    ReadOnly:\t%v\n",
		cinder.VolumeID, cinder.FSType, cinder.ReadOnly)
}

func printScaleIOVolumeSource(sio *core.ScaleIOVolumeSource, w PrefixWriter) {
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

func printLocalVolumeSource(ls *core.LocalVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tLocalVolume (a persistent volume backed by local storage on a node)\n"+
		"    Path:\t%v\n",
		ls.Path)
}

func printCephFSVolumeSource(cephfs *core.CephFSVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tCephFS (a CephFS mount on the host that shares a pod's lifetime)\n"+
		"    Monitors:\t%v\n"+
		"    Path:\t%v\n"+
		"    User:\t%v\n"+
		"    SecretFile:\t%v\n"+
		"    SecretRef:\t%v\n"+
		"    ReadOnly:\t%v\n",
		cephfs.Monitors, cephfs.Path, cephfs.User, cephfs.SecretFile, cephfs.SecretRef, cephfs.ReadOnly)
}

func printCephFSPersistentVolumeSource(cephfs *core.CephFSPersistentVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tCephFS (a CephFS mount on the host that shares a pod's lifetime)\n"+
		"    Monitors:\t%v\n"+
		"    Path:\t%v\n"+
		"    User:\t%v\n"+
		"    SecretFile:\t%v\n"+
		"    SecretRef:\t%v\n"+
		"    ReadOnly:\t%v\n",
		cephfs.Monitors, cephfs.Path, cephfs.User, cephfs.SecretFile, cephfs.SecretRef, cephfs.ReadOnly)
}

func printStorageOSVolumeSource(storageos *core.StorageOSVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tStorageOS (a StorageOS Persistent Disk resource)\n"+
		"    VolumeName:\t%v\n"+
		"    VolumeNamespace:\t%v\n"+
		"    FSType:\t%v\n"+
		"    ReadOnly:\t%v\n",
		storageos.VolumeName, storageos.VolumeNamespace, storageos.FSType, storageos.ReadOnly)
}

func printStorageOSPersistentVolumeSource(storageos *core.StorageOSPersistentVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tStorageOS (a StorageOS Persistent Disk resource)\n"+
		"    VolumeName:\t%v\n"+
		"    VolumeNamespace:\t%v\n"+
		"    FSType:\t%v\n"+
		"    ReadOnly:\t%v\n",
		storageos.VolumeName, storageos.VolumeNamespace, storageos.FSType, storageos.ReadOnly)
}

func printFCVolumeSource(fc *core.FCVolumeSource, w PrefixWriter) {
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

func printAzureFileVolumeSource(azureFile *core.AzureFileVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tAzureFile (an Azure File Service mount on the host and bind mount to the pod)\n"+
		"    SecretName:\t%v\n"+
		"    ShareName:\t%v\n"+
		"    ReadOnly:\t%v\n",
		azureFile.SecretName, azureFile.ShareName, azureFile.ReadOnly)
}

func printAzureFilePersistentVolumeSource(azureFile *core.AzureFilePersistentVolumeSource, w PrefixWriter) {
	ns := ""
	if azureFile.SecretNamespace != nil {
		ns = *azureFile.SecretNamespace
	}
	w.Write(LEVEL_2, "Type:\tAzureFile (an Azure File Service mount on the host and bind mount to the pod)\n"+
		"    SecretName:\t%v\n"+
		"    SecretNamespace:\t%v\n"+
		"    ShareName:\t%v\n"+
		"    ReadOnly:\t%v\n",
		azureFile.SecretName, ns, azureFile.ShareName, azureFile.ReadOnly)
}

func printFlexVolumeSource(flex *core.FlexVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tFlexVolume (a generic volume resource that is provisioned/attached using an exec based plugin)\n"+
		"    Driver:\t%v\n"+
		"    FSType:\t%v\n"+
		"    SecretRef:\t%v\n"+
		"    ReadOnly:\t%v\n",
		"    Options:\t%v\n",
		flex.Driver, flex.FSType, flex.SecretRef, flex.ReadOnly, flex.Options)
}

func printFlockerVolumeSource(flocker *core.FlockerVolumeSource, w PrefixWriter) {
	w.Write(LEVEL_2, "Type:\tFlocker (a Flocker volume mounted by the Flocker agent)\n"+
		"    DatasetName:\t%v\n"+
		"    DatasetUUID:\t%v\n",
		flocker.DatasetName, flocker.DatasetUUID)
}

type flusher interface {
	Flush()
}
