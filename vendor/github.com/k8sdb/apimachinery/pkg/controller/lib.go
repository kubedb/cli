package controller

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/appscode/log"
	"github.com/ghodss/yaml"
	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/azure"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/pkg/eventer"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	batch "k8s.io/client-go/pkg/apis/batch/v1"
	"k8s.io/client-go/tools/record"
)

func (c *Controller) ValidateStorageSpec(spec *tapi.StorageSpec) (*tapi.StorageSpec, error) {
	if spec == nil {
		return nil, nil
	}

	if spec.Class == "" {
		return nil, fmt.Errorf(`Object 'Class' is missing in '%v'`, *spec)
	}

	if _, err := c.Client.StorageV1().StorageClasses().Get(spec.Class, metav1.GetOptions{}); err != nil {
		if kerr.IsNotFound(err) {
			return nil, fmt.Errorf(`Spec.Storage.Class "%v" not found`, spec.Class)
		}
		return nil, err
	}

	if len(spec.AccessModes) == 0 {
		spec.AccessModes = []apiv1.PersistentVolumeAccessMode{
			apiv1.ReadWriteOnce,
		}
		log.Infof(`Using "%v" as AccessModes in "%v"`, apiv1.ReadWriteOnce, *spec)
	}

	if val, found := spec.Resources.Requests[apiv1.ResourceStorage]; found {
		if val.Value() <= 0 {
			return nil, errors.New("Invalid ResourceStorage request")
		}
	} else {
		return nil, errors.New("Missing ResourceStorage request")
	}

	return spec, nil
}

func (c *Controller) ValidateBackupSchedule(spec *tapi.BackupScheduleSpec) error {
	if spec == nil {
		return nil
	}
	// CronExpression can't be empty
	if spec.CronExpression == "" {
		return errors.New("Invalid cron expression")
	}

	return c.ValidateSnapshotSpec(spec.SnapshotStorageSpec)
}

func (c *Controller) ValidateSnapshotSpec(spec tapi.SnapshotStorageSpec) error {
	// BucketName can't be empty
	bucketName := spec.BucketName
	if bucketName == "" {
		return fmt.Errorf(`Object 'BucketName' is missing in '%v'`, spec)
	}

	// Need to provide Storage credential secret
	storageSecret := spec.StorageSecret
	if storageSecret == nil {
		return fmt.Errorf(`Object 'StorageSecret' is missing in '%v'`, spec)
	}

	// Credential SecretName  can't be empty
	storageSecretName := storageSecret.SecretName
	if storageSecretName == "" {
		return fmt.Errorf(`Object 'SecretName' is missing in '%v'`, *spec.StorageSecret)
	}
	return nil
}

const (
	keyProvider = "provider"
	keyConfig   = "config"
)

func (c *Controller) CheckBucketAccess(snapshotSpec tapi.SnapshotStorageSpec, namespace string) error {
	secret, err := c.Client.CoreV1().Secrets(namespace).Get(snapshotSpec.StorageSecret.SecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	provider := secret.Data[keyProvider]
	if provider == nil {
		return errors.New("Missing provider key")
	}
	configData := secret.Data[keyConfig]
	if configData == nil {
		return errors.New("Missing config key")
	}

	var config stow.ConfigMap
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return err
	}

	loc, err := stow.Dial(string(provider), config)
	if err != nil {
		return err
	}

	container, err := loc.Container(snapshotSpec.BucketName)
	if err != nil {
		return err
	}

	r := bytes.NewReader([]byte("CheckBucketAccess"))
	item, err := container.Put(".kubedb", r, r.Size(), nil)
	if err != nil {
		return err
	}

	if err := container.RemoveItem(item.ID()); err != nil {
		return err
	}
	return nil
}

func (c *Controller) CheckStatefulSetPodStatus(statefulSet *apps.StatefulSet, checkDuration time.Duration) error {
	podName := fmt.Sprintf("%v-%v", statefulSet.Name, 0)

	podReady := false
	then := time.Now()
	now := time.Now()
	for now.Sub(then) < checkDuration {
		pod, err := c.Client.CoreV1().Pods(statefulSet.Namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			if kerr.IsNotFound(err) {
				_, err := c.Client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Get(statefulSet.Name, metav1.GetOptions{})
				if kerr.IsNotFound(err) {
					break
				}

				time.Sleep(sleepDuration)
				now = time.Now()
				continue
			} else {
				return err
			}
		}
		log.Debugf("Pod Phase: %v", pod.Status.Phase)

		// If job is success
		if pod.Status.Phase == apiv1.PodRunning {
			podReady = true
			break
		}

		time.Sleep(sleepDuration)
		now = time.Now()
	}
	if !podReady {
		return errors.New("Database fails to be Ready")
	}
	return nil
}

func (c *Controller) DeletePersistentVolumeClaims(namespace string, selector labels.Selector) error {
	pvcList, err := c.Client.CoreV1().PersistentVolumeClaims(namespace).List(
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)
	if err != nil {
		return err
	}

	for _, pvc := range pvcList.Items {
		if err := c.Client.CoreV1().PersistentVolumeClaims(pvc.Namespace).Delete(pvc.Name, nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) DeleteSnapshotData(snapshot *tapi.Snapshot) error {
	secret, err := c.Client.CoreV1().Secrets(snapshot.Namespace).Get(snapshot.Spec.StorageSecret.SecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	provider := secret.Data[keyProvider]
	if provider == nil {
		return errors.New("Missing provider key")
	}
	configData := secret.Data[keyConfig]
	if configData == nil {
		return errors.New("Missing config key")
	}

	var config stow.ConfigMap
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return err
	}

	loc, err := stow.Dial(string(provider), config)
	if err != nil {
		return err
	}

	container, err := loc.Container(snapshot.Spec.BucketName)
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("%v/%v/%v/%v", DatabaseNamePrefix, snapshot.Namespace, snapshot.Spec.DatabaseName, snapshot.Name)
	cursor := stow.CursorStart
	for {
		items, next, err := container.Items(prefix, cursor, 50)
		if err != nil {
			return err
		}
		for _, item := range items {
			if err := container.RemoveItem(item.ID()); err != nil {
				return err
			}
		}
		cursor = next
		if stow.IsCursorEnd(cursor) {
			break
		}
	}

	return nil
}

func (c *Controller) DeleteSnapshots(namespace string, selector labels.Selector) error {
	snapshotList, err := c.ExtClient.Snapshots(namespace).List(
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)
	if err != nil {
		return err
	}

	for _, snapshot := range snapshotList.Items {
		if err := c.ExtClient.Snapshots(snapshot.Namespace).Delete(snapshot.Name); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) CheckDatabaseRestoreJob(
	job *batch.Job,
	runtimeObj runtime.Object,
	recorder record.EventRecorder,
	checkDuration time.Duration,
) bool {
	var jobSuccess bool = false
	var err error

	then := time.Now()
	now := time.Now()
	for now.Sub(then) < checkDuration {
		log.Debugln("Checking for Job ", job.Name)
		job, err = c.Client.BatchV1().Jobs(job.Namespace).Get(job.Name, metav1.GetOptions{})
		if err != nil {
			if kerr.IsNotFound(err) {
				time.Sleep(sleepDuration)
				now = time.Now()
				continue
			}
			recorder.Eventf(
				runtimeObj,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToList,
				"Failed to get Job. Reason: %v",
				err,
			)
			log.Errorln(err)
			return jobSuccess
		}
		log.Debugf("Pods Statuses:	%d Running / %d Succeeded / %d Failed",
			job.Status.Active, job.Status.Succeeded, job.Status.Failed)
		// If job is success
		if job.Status.Succeeded > 0 {
			jobSuccess = true
			break
		} else if job.Status.Failed > 0 {
			break
		}

		time.Sleep(sleepDuration)
		now = time.Now()
	}

	if err != nil {
		return false
	}

	podList, err := c.Client.CoreV1().Pods(job.Namespace).List(
		metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(job.Spec.Selector.MatchLabels).String(),
		},
	)
	if err != nil {
		recorder.Eventf(
			runtimeObj,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToList,
			"Failed to list Pods. Reason: %v",
			err,
		)
		log.Errorln(err)
		return jobSuccess
	}

	for _, pod := range podList.Items {
		if err := c.Client.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
			recorder.Eventf(
				runtimeObj,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				"Failed to delete Pod. Reason: %v",
				err,
			)
			log.Errorln(err)
		}
	}

	for _, volume := range job.Spec.Template.Spec.Volumes {
		claim := volume.PersistentVolumeClaim
		if claim != nil {
			err := c.Client.CoreV1().PersistentVolumeClaims(job.Namespace).Delete(claim.ClaimName, nil)
			if err != nil {
				recorder.Eventf(
					runtimeObj,
					apiv1.EventTypeWarning,
					eventer.EventReasonFailedToDelete,
					"Failed to delete PersistentVolumeClaim. Reason: %v",
					err,
				)
				log.Errorln(err)
			}
		}
	}

	if err := c.Client.BatchV1().Jobs(job.Namespace).Delete(job.Name, nil); err != nil {
		recorder.Eventf(
			runtimeObj,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			"Failed to delete Job. Reason: %v",
			err,
		)
		log.Errorln(err)
	}

	return jobSuccess
}

func (c *Controller) checkGoverningService(name, namespace string) (bool, error) {
	_, err := c.Client.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (c *Controller) CreateGoverningService(name, namespace string) error {
	// Check if service name exists
	found, err := c.checkGoverningService(name, namespace)
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: apiv1.ServiceSpec{
			Type:      apiv1.ServiceTypeClusterIP,
			ClusterIP: apiv1.ClusterIPNone,
		},
	}
	_, err = c.Client.CoreV1().Services(namespace).Create(service)
	return err
}

func (c *Controller) DeleteService(name, namespace string) error {
	service, err := c.Client.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	if service.Spec.Selector[LabelDatabaseName] != name {
		return nil
	}

	return c.Client.CoreV1().Services(namespace).Delete(name, nil)
}

func (c *Controller) DeleteStatefulSet(name, namespace string) error {
	statefulSet, err := c.Client.AppsV1beta1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	// Update StatefulSet
	replicas := int32(0)
	statefulSet.Spec.Replicas = &replicas
	if _, err := c.Client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Update(statefulSet); err != nil {
		return err
	}

	var checkSuccess bool = false
	then := time.Now()
	now := time.Now()
	for now.Sub(then) < time.Minute*10 {
		podList, err := c.Client.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
			LabelSelector: labels.Set(statefulSet.Spec.Selector.MatchLabels).AsSelector().String(),
		})
		if err != nil {
			return err
		}
		if len(podList.Items) == 0 {
			checkSuccess = true
			break
		}

		time.Sleep(sleepDuration)
		now = time.Now()
	}

	if !checkSuccess {
		return errors.New("Fail to delete StatefulSet Pods")
	}

	// Delete StatefulSet
	return c.Client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Delete(statefulSet.Name, nil)
}

func (c *Controller) DeleteSecret(name, namespace string) error {
	if _, err := c.Client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{}); err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	return c.Client.CoreV1().Secrets(namespace).Delete(name, nil)
}
