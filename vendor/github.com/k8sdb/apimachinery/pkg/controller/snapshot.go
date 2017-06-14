package controller

import (
	"fmt"
	"time"

	"github.com/appscode/go/wait"
	"github.com/appscode/log"
	tapi "github.com/k8sdb/apimachinery/api"
	tcs "github.com/k8sdb/apimachinery/client/clientset"
	"github.com/k8sdb/apimachinery/pkg/analytics"
	"github.com/k8sdb/apimachinery/pkg/eventer"
	kapi "k8s.io/kubernetes/pkg/api"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	kbatch "k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/cache"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/record"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
)

type Snapshotter interface {
	ValidateSnapshot(*tapi.Snapshot) error
	GetDatabase(*tapi.Snapshot) (runtime.Object, error)
	GetSnapshotter(*tapi.Snapshot) (*kbatch.Job, error)
	WipeOutSnapshot(*tapi.Snapshot) error
}

type SnapshotController struct {
	// Kubernetes client
	client clientset.Interface
	// ThirdPartyExtension client
	extClient tcs.ExtensionInterface
	// Snapshotter interface
	snapshoter Snapshotter
	// ListerWatcher
	lw *cache.ListWatch
	// Event Recorder
	eventRecorder record.EventRecorder
	// sync time to sync the list.
	syncPeriod time.Duration
}

const (
	LabelJobType        = "job.kubedb.com/type"
	LabelSnapshotStatus = "snapshot.kubedb.com/status"
)

// NewSnapshotController creates a new SnapshotController
func NewSnapshotController(
	client clientset.Interface,
	extClient tcs.ExtensionInterface,
	snapshoter Snapshotter,
	lw *cache.ListWatch,
	syncPeriod time.Duration,
) *SnapshotController {

	// return new DormantDatabase Controller
	return &SnapshotController{
		client:        client,
		extClient:     extClient,
		snapshoter:    snapshoter,
		lw:            lw,
		eventRecorder: eventer.NewEventRecorder(client, "Snapshot Controller"),
		syncPeriod:    syncPeriod,
	}
}

func (c *SnapshotController) Run() {
	// Ensure DormantDatabase TPR
	c.ensureThirdPartyResource()
	// Watch DormantDatabase with provided ListerWatcher
	c.watch()
}

// Ensure Snapshot ThirdPartyResource
func (c *SnapshotController) ensureThirdPartyResource() {
	log.Infoln("Ensuring Snapshot ThirdPartyResource")

	resourceName := tapi.ResourceNameSnapshot + "." + tapi.V1alpha1SchemeGroupVersion.Group
	var err error
	if _, err = c.client.Extensions().ThirdPartyResources().Get(resourceName); err == nil {
		return
	}
	if !k8serr.IsNotFound(err) {
		log.Fatalln(err)
	}

	thirdPartyResource := &extensions.ThirdPartyResource{
		TypeMeta: unversioned.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "ThirdPartyResource",
		},
		ObjectMeta: kapi.ObjectMeta{
			Name: resourceName,
		},
		Description: "Snapshot of kubedb databases",
		Versions: []extensions.APIVersion{
			{
				Name: tapi.V1alpha1SchemeGroupVersion.Version,
			},
		},
	}
	if _, err := c.client.Extensions().ThirdPartyResources().Create(thirdPartyResource); err != nil {
		log.Fatalln(err)
	}
}

func (c *SnapshotController) watch() {
	_, cacheController := cache.NewInformer(c.lw,
		&tapi.Snapshot{},
		c.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				snapshot := obj.(*tapi.Snapshot)
				if snapshot.Status.StartTime == nil {
					if err := c.create(snapshot); err != nil {
						snapshotFailedToCreate()
						log.Errorln(err)
					} else {
						snapshotSuccessfullyCreated()
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				snapshot := obj.(*tapi.Snapshot)
				if err := c.delete(snapshot); err != nil {
					snapshotFailedToDelete()
					log.Errorln(err)
				} else {
					snapshotSuccessfullyDeleted()
				}
			},
		},
	)
	cacheController.Run(wait.NeverStop)
}

const (
	durationCheckSnapshotJob = time.Minute * 30
)

func (c *SnapshotController) create(snapshot *tapi.Snapshot) error {
	var err error
	if snapshot, err = c.extClient.Snapshots(snapshot.Namespace).Get(snapshot.Name); err != nil {
		return err
	}

	t := unversioned.Now()
	snapshot.Status.StartTime = &t
	if _, err = c.extClient.Snapshots(snapshot.Namespace).Update(snapshot); err != nil {
		c.eventRecorder.Eventf(
			snapshot,
			kapi.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			`Fail to update Elastic: "%v". Reason: %v`,
			snapshot.Name,
			err,
		)
		log.Errorln(err)
	}

	// Validate DatabaseSnapshot spec
	if err := c.snapshoter.ValidateSnapshot(snapshot); err != nil {
		c.eventRecorder.Event(snapshot, kapi.EventTypeWarning, eventer.EventReasonInvalid, err.Error())
		return err
	}

	runtimeObj, err := c.snapshoter.GetDatabase(snapshot)
	if err != nil {
		c.eventRecorder.Event(snapshot, kapi.EventTypeWarning, eventer.EventReasonFailedToGet, err.Error())
		return err
	}

	if snapshot, err = c.extClient.Snapshots(snapshot.Namespace).Get(snapshot.Name); err != nil {
		return err
	}

	snapshot.Labels[LabelDatabaseName] = snapshot.Spec.DatabaseName
	snapshot.Labels[LabelSnapshotStatus] = string(tapi.SnapshotPhaseRunning)
	snapshot.Status.Phase = tapi.SnapshotPhaseRunning
	if _, err = c.extClient.Snapshots(snapshot.Namespace).Update(snapshot); err != nil {
		c.eventRecorder.Eventf(
			snapshot,
			kapi.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			"Failed to update Snapshot. Reason: %v",
			err,
		)
		log.Errorln(err)
	}

	c.eventRecorder.Event(runtimeObj, kapi.EventTypeNormal, eventer.EventReasonStarting, "Backup running")
	c.eventRecorder.Event(snapshot, kapi.EventTypeNormal, eventer.EventReasonStarting, "Backup running")

	job, err := c.snapshoter.GetSnapshotter(snapshot)
	if err != nil {
		message := fmt.Sprintf("Failed to take snapshot. Reason: %v", err)
		c.eventRecorder.Event(runtimeObj, kapi.EventTypeWarning, eventer.EventReasonSnapshotFailed, message)
		c.eventRecorder.Event(snapshot, kapi.EventTypeWarning, eventer.EventReasonSnapshotFailed, message)
		return err
	}

	if _, err := c.client.Batch().Jobs(snapshot.Namespace).Create(job); err != nil {
		message := fmt.Sprintf("Failed to take snapshot. Reason: %v", err)
		c.eventRecorder.Event(runtimeObj, kapi.EventTypeWarning, eventer.EventReasonSnapshotFailed, message)
		c.eventRecorder.Event(snapshot, kapi.EventTypeWarning, eventer.EventReasonSnapshotFailed, message)
		return err
	}

	go func() {
		if err := c.checkSnapshotJob(snapshot, job.Name, durationCheckSnapshotJob); err != nil {
			log.Errorln(err)
		}
	}()

	return nil
}

func (c *SnapshotController) delete(snapshot *tapi.Snapshot) error {
	runtimeObj, err := c.snapshoter.GetDatabase(snapshot)
	if err != nil {
		if !k8serr.IsNotFound(err) {
			c.eventRecorder.Event(
				snapshot,
				kapi.EventTypeWarning,
				eventer.EventReasonFailedToGet,
				err.Error(),
			)
			return err
		}
	}

	if runtimeObj != nil {
		c.eventRecorder.Eventf(
			runtimeObj,
			kapi.EventTypeNormal,
			eventer.EventReasonWipingOut,
			"Wiping out Snapshot: %v",
			snapshot.Name,
		)
	}

	if err := c.snapshoter.WipeOutSnapshot(snapshot); err != nil {
		if runtimeObj != nil {
			c.eventRecorder.Eventf(
				runtimeObj,
				kapi.EventTypeWarning,
				eventer.EventReasonFailedToWipeOut,
				"Failed to  wipeOut. Reason: %v",
				err,
			)
		}
		return err
	}

	if runtimeObj != nil {
		c.eventRecorder.Eventf(
			runtimeObj,
			kapi.EventTypeNormal,
			eventer.EventReasonSuccessfulWipeOut,
			"Successfully wiped out Snapshot: %v",
			snapshot.Name,
		)
	}
	return nil
}

func (c *SnapshotController) checkSnapshotJob(snapshot *tapi.Snapshot, jobName string, checkDuration time.Duration) error {

	var jobSuccess bool = false
	var job *kbatch.Job
	var err error
	then := time.Now()
	now := time.Now()
	for now.Sub(then) < checkDuration {
		log.Debugln("Checking for Job ", jobName)
		job, err = c.client.Batch().Jobs(snapshot.Namespace).Get(jobName)
		if err != nil {
			if k8serr.IsNotFound(err) {
				time.Sleep(sleepDuration)
				now = time.Now()
				continue
			}
			c.eventRecorder.Eventf(
				snapshot,
				kapi.EventTypeWarning,
				eventer.EventReasonFailedToList,
				"Failed to get Job. Reason: %v",
				err,
			)
			return err
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
		return err
	}

	podList, err := c.client.Core().Pods(job.Namespace).List(
		kapi.ListOptions{
			LabelSelector: labels.SelectorFromSet(job.Spec.Selector.MatchLabels),
		},
	)
	if err != nil {
		c.eventRecorder.Eventf(
			snapshot,
			kapi.EventTypeWarning,
			eventer.EventReasonFailedToList,
			"Failed to list Pods. Reason: %v",
			err,
		)
		return err
	}

	for _, pod := range podList.Items {
		if err := c.client.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
			c.eventRecorder.Eventf(
				snapshot,
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
			err := c.client.Core().PersistentVolumeClaims(job.Namespace).Delete(claim.ClaimName, nil)
			if err != nil {
				c.eventRecorder.Eventf(
					snapshot,
					kapi.EventTypeWarning,
					eventer.EventReasonFailedToDelete,
					"Failed to delete PersistentVolumeClaim. Reason: %v",
					err,
				)
				log.Errorln(err)
			}
		}
	}

	if err := c.client.Batch().Jobs(job.Namespace).Delete(job.Name, nil); err != nil {
		c.eventRecorder.Eventf(
			snapshot,
			kapi.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			"Failed to delete Job. Reason: %v",
			err,
		)
		log.Errorln(err)
	}

	if snapshot, err = c.extClient.Snapshots(snapshot.Namespace).Get(snapshot.Name); err != nil {
		return err
	}

	runtimeObj, err := c.snapshoter.GetDatabase(snapshot)
	if err != nil {
		c.eventRecorder.Event(snapshot, kapi.EventTypeWarning, eventer.EventReasonFailedToGet, err.Error())
		return err
	}

	t := unversioned.Now()
	snapshot.Status.CompletionTime = &t
	if jobSuccess {
		snapshot.Status.Phase = tapi.SnapshotPhaseSuccessed
		c.eventRecorder.Event(
			runtimeObj,
			kapi.EventTypeNormal,
			eventer.EventReasonSuccessfulSnapshot,
			"Successfully completed snapshot",
		)
		c.eventRecorder.Event(
			snapshot,
			kapi.EventTypeNormal,
			eventer.EventReasonSuccessfulSnapshot,
			"Successfully completed snapshot",
		)
	} else {
		snapshot.Status.Phase = tapi.SnapshotPhaseFailed
		c.eventRecorder.Event(
			runtimeObj,
			kapi.EventTypeWarning,
			eventer.EventReasonSnapshotFailed,
			"Failed to complete snapshot",
		)
		c.eventRecorder.Event(
			snapshot,
			kapi.EventTypeWarning,
			eventer.EventReasonSnapshotFailed,
			"Failed to complete snapshot",
		)
	}

	delete(snapshot.Labels, LabelSnapshotStatus)
	if _, err := c.extClient.Snapshots(snapshot.Namespace).Update(snapshot); err != nil {
		c.eventRecorder.Eventf(
			snapshot,
			kapi.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			"Failed to update Snapshot. Reason: %v",
			err,
		)
		log.Errorln(err)
	}
	return nil
}

func snapshotSuccessfullyCreated() {
	analytics.SendEvent(tapi.ResourceNameSnapshot, "created", "success")
}

func snapshotFailedToCreate() {
	analytics.SendEvent(tapi.ResourceNameSnapshot, "created", "failure")
}

func snapshotSuccessfullyDeleted() {
	analytics.SendEvent(tapi.ResourceNameSnapshot, "deleted", "success")
}

func snapshotFailedToDelete() {
	analytics.SendEvent(tapi.ResourceNameSnapshot, "deleted", "failure")
}
