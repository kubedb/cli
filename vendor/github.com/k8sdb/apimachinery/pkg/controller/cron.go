package controller

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/appscode/log"
	tapi "github.com/k8sdb/apimachinery/api"
	tcs "github.com/k8sdb/apimachinery/client/clientset"
	"github.com/k8sdb/apimachinery/pkg/eventer"
	cmap "github.com/orcaman/concurrent-map"
	"gopkg.in/robfig/cron.v2"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/record"
)

type CronControllerInterface interface {
	StartCron()
	ScheduleBackup(runtime.Object, metav1.ObjectMeta, *tapi.BackupScheduleSpec) error
	StopBackupScheduling(metav1.ObjectMeta)
	StopCron()
}

type cronController struct {
	// ThirdPartyExtension client
	extClient tcs.ExtensionInterface
	// For Internal Cron Job
	cron *cron.Cron
	// Store Cron Job EntryID for further use
	cronEntryIDs cmap.ConcurrentMap
	// Event Recorder
	eventRecorder record.EventRecorder
	// To perform start operation once
	once sync.Once
}

/*
 NewCronController returns CronControllerInterface.
 Need to call StartCron() method to start Cron.
*/
func NewCronController(client clientset.Interface, extClient tcs.ExtensionInterface) CronControllerInterface {
	return &cronController{
		extClient:     extClient,
		cron:          cron.New(),
		cronEntryIDs:  cmap.New(),
		eventRecorder: eventer.NewEventRecorder(client, "Cron Controller"),
	}
}

func (c *cronController) StartCron() {
	c.once.Do(func() {
		c.cron.Start()
	})
}

func (c *cronController) ScheduleBackup(
	// Runtime Object to push event
	runtimeObj runtime.Object,
	// ObjectMeta of Database TPR object
	om metav1.ObjectMeta,
	// BackupScheduleSpec
	spec *tapi.BackupScheduleSpec,
) error {
	// cronEntry name
	cronEntryName := fmt.Sprintf("%v@%v", om.Name, om.Namespace)

	// Remove previous cron job if exist
	if id, exists := c.cronEntryIDs.Pop(cronEntryName); exists {
		c.cron.Remove(id.(cron.EntryID))
	}

	invoker := &snapshotInvoker{
		extClient:     c.extClient,
		runtimeObject: runtimeObj,
		om:            om,
		spec:          spec,
		eventRecorder: c.eventRecorder,
	}

	if err := invoker.validateScheduler(durationCheckSnapshotJob); err != nil {
		return err
	}

	// Set cron job
	entryID, err := c.cron.AddFunc(spec.CronExpression, invoker.createScheduledSnapshot)
	if err != nil {
		return err
	}

	// Add job entryID
	c.cronEntryIDs.Set(cronEntryName, entryID)
	return nil
}

func (c *cronController) StopBackupScheduling(om metav1.ObjectMeta) {
	// cronEntry name
	cronEntryName := fmt.Sprintf("%v@%v", om.Name, om.Namespace)

	if id, exists := c.cronEntryIDs.Pop(cronEntryName); exists {
		c.cron.Remove(id.(cron.EntryID))
	}
}

func (c *cronController) StopCron() {
	c.cron.Stop()
}

type snapshotInvoker struct {
	extClient     tcs.ExtensionInterface
	runtimeObject runtime.Object
	om            metav1.ObjectMeta
	spec          *tapi.BackupScheduleSpec
	eventRecorder record.EventRecorder
}

func (s *snapshotInvoker) validateScheduler(checkDuration time.Duration) error {
	utc := time.Now().UTC()
	snapshotName := fmt.Sprintf("%v-%v", s.om.Name, utc.Format("20060102-150405"))
	if err := s.createSnapshot(snapshotName); err != nil {
		return err
	}

	var snapshotSuccess bool = false

	then := time.Now()
	now := time.Now()
	for now.Sub(then) < checkDuration {
		snapshot, err := s.extClient.Snapshots(s.om.Namespace).Get(snapshotName)
		if err != nil {
			if kerr.IsNotFound(err) {
				time.Sleep(sleepDuration)
				now = time.Now()
				continue
			} else {
				return err
			}
		}

		if snapshot.Status.Phase == tapi.SnapshotPhaseSuccessed {
			snapshotSuccess = true
			break
		}
		if snapshot.Status.Phase == tapi.SnapshotPhaseFailed {
			break
		}

		time.Sleep(sleepDuration)
		now = time.Now()
	}

	if !snapshotSuccess {
		return errors.New("Failed to complete initial snapshot")
	}

	return nil
}

func (s *snapshotInvoker) createScheduledSnapshot() {
	kind := s.runtimeObject.GetObjectKind().GroupVersionKind().Kind
	name := s.om.Name

	labelMap := map[string]string{
		LabelDatabaseKind:   kind,
		LabelDatabaseName:   name,
		LabelSnapshotStatus: string(tapi.SnapshotPhaseRunning),
	}

	snapshotList, err := s.extClient.Snapshots(s.om.Namespace).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelMap).AsSelector().String(),
	})
	if err != nil {
		s.eventRecorder.Eventf(
			s.runtimeObject,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToList,
			"Failed to list Snapshots. Reason: %v",
			err,
		)
		log.Errorln(err)
		return
	}

	if len(snapshotList.Items) > 0 {
		s.eventRecorder.Event(
			s.runtimeObject,
			apiv1.EventTypeNormal,
			eventer.EventReasonIgnoredSnapshot,
			"Skipping scheduled Backup. One is still active.",
		)
		log.Debugln("Skipping scheduled Backup. One is still active.")
		return
	}

	// Set label. Elastic controller will detect this using label selector
	labelMap = map[string]string{
		LabelDatabaseKind: kind,
		LabelDatabaseName: name,
	}

	now := time.Now().UTC()
	snapshotName := fmt.Sprintf("%v-%v", s.om.Name, now.Format("20060102-150405"))

	if err = s.createSnapshot(snapshotName); err != nil {
		log.Errorln(err)
	}
}

func (s *snapshotInvoker) createSnapshot(snapshotName string) error {
	labelMap := map[string]string{
		LabelDatabaseKind: s.runtimeObject.GetObjectKind().GroupVersionKind().Kind,
		LabelDatabaseName: s.om.Name,
	}

	snapshot := &tapi.Snapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshotName,
			Namespace: s.om.Namespace,
			Labels:    labelMap,
		},
		Spec: tapi.SnapshotSpec{
			DatabaseName:        s.om.Name,
			SnapshotStorageSpec: s.spec.SnapshotStorageSpec,
		},
	}

	if _, err := s.extClient.Snapshots(snapshot.Namespace).Create(snapshot); err != nil {
		s.eventRecorder.Eventf(
			s.runtimeObject,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToCreate,
			"Failed to create Snapshot. Reason: %v",
			err,
		)
		return err
	}

	return nil
}
