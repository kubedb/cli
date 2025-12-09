/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Free Trial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Free-Trial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

import (
	"context"
	"fmt"
	"time"

	"gomodules.xyz/wait"
	batchv1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/ptr"
	core_util "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
)

const (
	migratorPVCTemplate = "data-migrate"

	pvcMounterPodTemplate = "pvcmounter"
	pvcMounterImage       = "alpine:latest"
	pvcMounterVolumeName  = "test-pvc"
	pvcMounterMountPath   = "test-pvc"

	migratorJobTemplate              = "migrator"
	migratorJobSourceVolumeName      = "source-volume"
	migratorJobSourceMountPath       = "/mnt/source"
	migratorJobDestinationVolumeName = "destination-volume"
	migratorJobDestinationMountPath  = "/mnt/destination"

	podMigrationCompleted = "PodMigrationCompleted"
	rsyncImageName        = "ghcr.io/kubedb/rsync:v2025.8.31"
)

func GetDatabasePVCName(pvcTemplate string, podName string) string {
	return meta_util.NameWithSuffix(pvcTemplate, podName)
}

func GetMigratorPVCName(podName string) string {
	return meta_util.NameWithSuffix(migratorPVCTemplate, podName)
}

func GetPVCMounterPodName(podName string) string {
	return meta_util.NameWithSuffix(pvcMounterPodTemplate, podName)
}

func GetStorageMigratorJobName(podName string) string {
	return meta_util.NameWithSuffix(migratorJobTemplate, podName)
}

func GetPodMigrationCompleteCondition(podName string) string {
	return meta_util.NameWithSuffix(podMigrationCompleted, podName)
}

// -----------------------

func CreateMigratorPVC(client kubernetes.Interface, pvcMeta metav1.ObjectMeta, pvc *core.PersistentVolumeClaim, storageClass *string, owner *metav1.OwnerReference) error {
	_, _, err := core_util.CreateOrPatchPVC(context.TODO(), client, pvcMeta, func(claim *core.PersistentVolumeClaim) *core.PersistentVolumeClaim {
		core_util.EnsureOwnerReference(&claim.ObjectMeta, owner)
		claim.Spec.StorageClassName = storageClass
		claim.Spec.AccessModes = pvc.Spec.AccessModes
		claim.Spec.Resources = pvc.Spec.Resources
		claim.Spec.VolumeMode = pvc.Spec.VolumeMode
		return claim
	}, metav1.PatchOptions{})
	return err
}

func CreateDatabasePVC(client kubernetes.Interface, pvcMeta metav1.ObjectMeta, pvc *core.PersistentVolumeClaim, storageClass *string, labels map[string]string) error {
	_, _, err := core_util.CreateOrPatchPVC(context.TODO(), client, pvcMeta, func(claim *core.PersistentVolumeClaim) *core.PersistentVolumeClaim {
		claim.Labels = labels
		claim.Spec.StorageClassName = storageClass
		claim.Spec.AccessModes = pvc.Spec.AccessModes
		claim.Spec.Resources = pvc.Spec.Resources
		claim.Spec.VolumeMode = pvc.Spec.VolumeMode
		claim.Spec.VolumeName = pvc.Spec.VolumeName
		return claim
	}, metav1.PatchOptions{})
	return err
}

func CreatePVCMounterPod(client kubernetes.Interface, dbPod *core.Pod, pvcMounter string) error {
	pvcName := GetMigratorPVCName(dbPod.Name)

	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcMounter,
			Namespace: dbPod.Namespace,
		},
		Spec: core.PodSpec{
			RestartPolicy: core.RestartPolicyNever,
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": dbPod.Spec.NodeName,
			},
			Containers: []core.Container{
				{
					Name:            pvcMounterPodTemplate,
					Image:           pvcMounterImage,
					ImagePullPolicy: core.PullIfNotPresent,
					Command: []string{
						"/bin/sh",
						"-c",
						fmt.Sprintf(
							`echo "Checking PVC mount...";
						if mountpoint -q %v; then
							echo "PVC is mounted!";
							exit 0;
						else
							echo "PVC is NOT mounted!";
							exit 1;
						fi`, pvcMounterMountPath),
					},
					VolumeMounts: []core.VolumeMount{
						{
							Name:      pvcMounterVolumeName,
							MountPath: pvcMounterMountPath,
						},
					},
				},
			},
			Volumes: []core.Volume{
				{
					Name: pvcMounterVolumeName,
					VolumeSource: core.VolumeSource{
						PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
						},
					},
				},
			},
		},
	}

	_, _, err := core_util.CreateOrPatchPod(context.TODO(), client, pod.ObjectMeta, func(p *core.Pod) *core.Pod {
		return pod
	}, metav1.PatchOptions{})
	return err
}

func CreateDataMigratorJob(client kubernetes.Interface, jobMeta metav1.ObjectMeta, pvcTemplate string, podName string) error {
	job := batchv1.Job{
		ObjectMeta: jobMeta,
		Spec: batchv1.JobSpec{
			BackoffLimit: ptr.To(int32(6)),
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					RestartPolicy: core.RestartPolicyNever,
					Containers: []core.Container{
						{
							Name:            migratorJobTemplate,
							Image:           rsyncImageName,
							ImagePullPolicy: core.PullIfNotPresent,
							Command: []string{
								"/bin/sh",
								"-c",
								fmt.Sprintf(
									`echo "Starting rsync from source to destination..."
								 rsync -avh %v/ %v/
            					 STATUS=$?
								 if [ $STATUS -eq 0 ]; then
								   echo "Rsync completed successfully."
								   exit 0
								 else
								   echo "Rsync failed!"
								   exit 1
								 fi`, migratorJobSourceMountPath, migratorJobDestinationMountPath),
							},
							VolumeMounts: []core.VolumeMount{
								{
									Name:      migratorJobSourceVolumeName,
									MountPath: migratorJobSourceMountPath,
								},
								{
									Name:      migratorJobDestinationVolumeName,
									MountPath: migratorJobDestinationMountPath,
								},
							},
						},
					},
					Volumes: []core.Volume{
						{
							Name: migratorJobSourceVolumeName,
							VolumeSource: core.VolumeSource{
								PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
									ClaimName: GetDatabasePVCName(pvcTemplate, podName),
								},
							},
						},
						{
							Name: migratorJobDestinationVolumeName,
							VolumeSource: core.VolumeSource{
								PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
									ClaimName: GetMigratorPVCName(podName),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := client.BatchV1().Jobs(jobMeta.Namespace).Create(context.TODO(), &job, metav1.CreateOptions{})
	return err
}

func WaitForJobToBeCompleted(ctx context.Context, timeOut time.Duration, c kubernetes.Interface, meta metav1.ObjectMeta) error {
	return wait.PollImmediate(time.Second, timeOut, func() (bool, error) {
		job, err := c.BatchV1().Jobs(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		if IsJobConditionTrue(job.Status.Conditions, batchv1.JobComplete) {
			return true, nil
		}
		if IsJobConditionTrue(job.Status.Conditions, batchv1.JobFailed) {
			return false, fmt.Errorf("job is failed")
		}
		return false, nil
	})
}

func IsJobConditionTrue(conditions []batchv1.JobCondition, condType batchv1.JobConditionType) bool {
	for i := range conditions {
		if conditions[i].Type == condType && conditions[i].Status == core.ConditionTrue {
			return true
		}
	}
	return false
}
