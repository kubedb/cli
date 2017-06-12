package controller

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/appscode/log"
	"github.com/ghodss/yaml"
	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/pkg/eventer"
	kapi "k8s.io/kubernetes/pkg/api"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	kapps "k8s.io/kubernetes/pkg/apis/apps"
	kbatch "k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/client/record"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
)

func (c *Controller) ValidateStorageSpec(spec *tapi.StorageSpec) (*tapi.StorageSpec, error) {
	if spec == nil {
		return nil, nil
	}

	if spec.Class == "" {
		return nil, fmt.Errorf(`Object 'Class' is missing in '%v'`, *spec)
	}

	if _, err := c.Client.Storage().StorageClasses().Get(spec.Class); err != nil {
		if k8serr.IsNotFound(err) {
			return nil, fmt.Errorf(`Spec.Storage.Class "%v" not found`, spec.Class)
		}
		return nil, err
	}

	if len(spec.AccessModes) == 0 {
		spec.AccessModes = []kapi.PersistentVolumeAccessMode{
			kapi.ReadWriteOnce,
		}
		log.Infof(`Using "%v" as AccessModes in "%v"`, kapi.ReadWriteOnce, *spec)
	}

	if val, found := spec.Resources.Requests[kapi.ResourceStorage]; found {
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
	secret, err := c.Client.Core().Secrets(namespace).Get(snapshotSpec.StorageSecret.SecretName)
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

func (c *Controller) CheckStatefulSetPodStatus(statefulSet *kapps.StatefulSet, checkDuration time.Duration) error {
	podName := fmt.Sprintf("%v-%v", statefulSet.Name, 0)

	podReady := false
	then := time.Now()
	now := time.Now()
	for now.Sub(then) < checkDuration {
		pod, err := c.Client.Core().Pods(statefulSet.Namespace).Get(podName)
		if err != nil {
			if k8serr.IsNotFound(err) {
				_, err := c.Client.Apps().StatefulSets(statefulSet.Namespace).Get(statefulSet.Name)
				if k8serr.IsNotFound(err) {
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
		if pod.Status.Phase == kapi.PodRunning {
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
	pvcList, err := c.Client.Core().PersistentVolumeClaims(namespace).List(
		kapi.ListOptions{
			LabelSelector: selector,
		},
	)
	if err != nil {
		return err
	}

	for _, pvc := range pvcList.Items {
		if err := c.Client.Core().PersistentVolumeClaims(pvc.Namespace).Delete(pvc.Name, nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) DeleteSnapshotData(snapshot *tapi.Snapshot) error {
	secret, err := c.Client.Core().Secrets(snapshot.Namespace).Get(snapshot.Spec.StorageSecret.SecretName)
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
		kapi.ListOptions{
			LabelSelector: selector,
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
	job *kbatch.Job,
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
		job, err = c.Client.Batch().Jobs(job.Namespace).Get(job.Name)
		if err != nil {
			if k8serr.IsNotFound(err) {
				time.Sleep(sleepDuration)
				now = time.Now()
				continue
			}
			recorder.Eventf(
				runtimeObj,
				kapi.EventTypeWarning,
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

	podList, err := c.Client.Core().Pods(job.Namespace).List(
		kapi.ListOptions{
			LabelSelector: labels.SelectorFromSet(job.Spec.Selector.MatchLabels),
		},
	)
	if err != nil {
		recorder.Eventf(
			runtimeObj,
			kapi.EventTypeWarning,
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
				kapi.EventTypeWarning,
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
			err := c.Client.Core().PersistentVolumeClaims(job.Namespace).Delete(claim.ClaimName, nil)
			if err != nil {
				recorder.Eventf(
					runtimeObj,
					kapi.EventTypeWarning,
					eventer.EventReasonFailedToDelete,
					"Failed to delete PersistentVolumeClaim. Reason: %v",
					err,
				)
				log.Errorln(err)
			}
		}
	}

	if err := c.Client.Batch().Jobs(job.Namespace).Delete(job.Name, nil); err != nil {
		recorder.Eventf(
			runtimeObj,
			kapi.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			"Failed to delete Job. Reason: %v",
			err,
		)
		log.Errorln(err)
	}

	return jobSuccess
}

func (c *Controller) checkGoverningService(name, namespace string) (bool, error) {
	_, err := c.Client.Core().Services(namespace).Get(name)
	if err != nil {
		if k8serr.IsNotFound(err) {
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

	service := &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name: name,
		},
		Spec: kapi.ServiceSpec{
			Type:      kapi.ServiceTypeClusterIP,
			ClusterIP: kapi.ClusterIPNone,
		},
	}
	_, err = c.Client.Core().Services(namespace).Create(service)
	return err
}

func (c *Controller) DeleteService(name, namespace string) error {
	service, err := c.Client.Core().Services(namespace).Get(name)
	if err != nil {
		if k8serr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	if service.Spec.Selector[LabelDatabaseName] != name {
		return nil
	}

	return c.Client.Core().Services(namespace).Delete(name, nil)
}

func (c *Controller) DeleteStatefulSet(name, namespace string) error {
	statefulSet, err := c.Client.Apps().StatefulSets(namespace).Get(name)
	if err != nil {
		if k8serr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	// Update StatefulSet
	statefulSet.Spec.Replicas = 0
	if _, err := c.Client.Apps().StatefulSets(statefulSet.Namespace).Update(statefulSet); err != nil {
		return err
	}

	labelSelector := labels.SelectorFromSet(statefulSet.Spec.Selector.MatchLabels)

	var checkSuccess bool = false
	then := time.Now()
	now := time.Now()
	for now.Sub(then) < time.Minute*10 {
		podList, err := c.Client.Core().Pods(kapi.NamespaceAll).List(kapi.ListOptions{
			LabelSelector: labelSelector,
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
	return c.Client.Apps().StatefulSets(statefulSet.Namespace).Delete(statefulSet.Name, nil)
}

func (c *Controller) DeleteSecret(name, namespace string) error {
	if _, err := c.Client.Core().Secrets(namespace).Get(name); err != nil {
		if k8serr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	return c.Client.Core().Secrets(namespace).Delete(name, nil)
}
