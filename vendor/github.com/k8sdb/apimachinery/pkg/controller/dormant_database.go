package controller

import (
	"errors"
	"reflect"
	"time"

	"github.com/appscode/go/wait"
	"github.com/appscode/log"
	tapi "github.com/k8sdb/apimachinery/api"
	tcs "github.com/k8sdb/apimachinery/client/clientset"
	"github.com/k8sdb/apimachinery/pkg/analytics"
	"github.com/k8sdb/apimachinery/pkg/eventer"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

type Deleter interface {
	// Check Database TPR
	Exists(*metav1.ObjectMeta) (bool, error)
	// Pause operation
	PauseDatabase(*tapi.DormantDatabase) error
	// Wipe out operation
	WipeOutDatabase(*tapi.DormantDatabase) error
	// Resume operation
	ResumeDatabase(*tapi.DormantDatabase) error
}

type DormantDbController struct {
	// Kubernetes client
	client clientset.Interface
	// ThirdPartyExtension client
	extClient tcs.ExtensionInterface
	// Deleter interface
	deleter Deleter
	// ListerWatcher
	lw *cache.ListWatch
	// Event Recorder
	eventRecorder record.EventRecorder
	// sync time to sync the list.
	syncPeriod time.Duration
}

// NewDormantDbController creates a new DormantDatabase Controller
func NewDormantDbController(
	client clientset.Interface,
	extClient tcs.ExtensionInterface,
	deleter Deleter,
	lw *cache.ListWatch,
	syncPeriod time.Duration,
) *DormantDbController {
	// return new DormantDatabase Controller
	return &DormantDbController{
		client:        client,
		extClient:     extClient,
		deleter:       deleter,
		lw:            lw,
		eventRecorder: eventer.NewEventRecorder(client, "DormantDatabase Controller"),
		syncPeriod:    syncPeriod,
	}
}

func (c *DormantDbController) Run() {
	// Ensure DormantDatabase TPR
	c.ensureThirdPartyResource()
	// Watch DormantDatabase with provided ListerWatcher
	c.watch()
}

// Ensure DormantDatabase ThirdPartyResource
func (c *DormantDbController) ensureThirdPartyResource() {
	log.Infoln("Ensuring DormantDatabase ThirdPartyResource")

	resourceName := tapi.ResourceNameDormantDatabase + "." + tapi.V1alpha1SchemeGroupVersion.Group
	var err error
	if _, err = c.client.ExtensionsV1beta1().ThirdPartyResources().Get(resourceName, metav1.GetOptions{}); err == nil {
		return
	}
	if !kerr.IsNotFound(err) {
		log.Fatalln(err)
	}

	thirdPartyResource := &extensions.ThirdPartyResource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "ThirdPartyResource",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
			Labels: map[string]string{
				"app": "kubedb",
			},
		},
		Description: "Dormant KubeDB databases",
		Versions: []extensions.APIVersion{
			{
				Name: tapi.V1alpha1SchemeGroupVersion.Version,
			},
		},
	}
	if _, err := c.client.ExtensionsV1beta1().ThirdPartyResources().Create(thirdPartyResource); err != nil {
		log.Fatalln(err)
	}
}

func (c *DormantDbController) watch() {
	_, cacheController := cache.NewInformer(c.lw,
		&tapi.DormantDatabase{},
		c.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				dormantDb := obj.(*tapi.DormantDatabase)
				if dormantDb.Status.CreationTime == nil {
					if err := c.create(dormantDb); err != nil {
						dormantDbFailedToCreate()
						log.Errorln(err)
					} else {
						dormantDbSuccessfullyCreated()
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if err := c.delete(obj.(*tapi.DormantDatabase)); err != nil {
					dormantDbFailedToDelete()
					log.Errorln(err)
				} else {
					dormantDbSuccessfullyDeleted()
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldDormantDb, ok := old.(*tapi.DormantDatabase)
				if !ok {
					return
				}
				newDormantDb, ok := new.(*tapi.DormantDatabase)
				if !ok {
					return
				}
				// TODO: Find appropriate checking
				// Only allow if Spec varies
				if !reflect.DeepEqual(oldDormantDb.Spec, newDormantDb.Spec) {
					if err := c.update(oldDormantDb, newDormantDb); err != nil {
						log.Errorln(err)
					}
				}
			},
		},
	)
	cacheController.Run(wait.NeverStop)
}

func (c *DormantDbController) create(dormantDb *tapi.DormantDatabase) error {

	var err error
	if dormantDb, err = c.extClient.DormantDatabases(dormantDb.Namespace).Get(dormantDb.Name); err != nil {
		return err
	}

	// Set DormantDatabase Phase: Deleting
	t := metav1.Now()
	dormantDb.Status.CreationTime = &t
	if _, err := c.extClient.DormantDatabases(dormantDb.Namespace).Update(dormantDb); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			"Failed to update DormantDatabase. Reason: %v",
			err,
		)
		return err
	}

	// Check if DB TPR object exists
	found, err := c.deleter.Exists(&dormantDb.ObjectMeta)
	if err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToPause,
			"Failed to pause Database. Reason: %v",
			err,
		)
		return err
	}

	if found {
		message := "Failed to pause Database. Delete Database TPR object first"
		c.eventRecorder.Event(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToPause,
			message,
		)

		// Delete DormantDatabase object
		if err := c.extClient.DormantDatabases(dormantDb.Namespace).Delete(dormantDb.Name); err != nil {
			c.eventRecorder.Eventf(
				dormantDb,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				"Failed to delete DormantDatabase. Reason: %v",
				err,
			)
			log.Errorln(err)
		}
		return errors.New(message)
	}

	if dormantDb, err = c.extClient.DormantDatabases(dormantDb.Namespace).Get(dormantDb.Name); err != nil {
		return err
	}

	// Set DormantDatabase Phase: Deleting
	t = metav1.Now()
	dormantDb.Status.Phase = tapi.DormantDatabasePhasePausing
	if _, err = c.extClient.DormantDatabases(dormantDb.Namespace).Update(dormantDb); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			"Failed to update DormantDatabase. Reason: %v",
			err,
		)
		return err
	}

	c.eventRecorder.Event(dormantDb, apiv1.EventTypeNormal, eventer.EventReasonPausing, "Pausing Database")

	// Pause Database workload
	if err := c.deleter.PauseDatabase(dormantDb); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			"Failed to pause. Reason: %v",
			err,
		)
		return err
	}

	c.eventRecorder.Event(
		dormantDb,
		apiv1.EventTypeNormal,
		eventer.EventReasonSuccessfulPause,
		"Successfully paused Database workload",
	)

	if dormantDb, err = c.extClient.DormantDatabases(dormantDb.Namespace).Get(dormantDb.Name); err != nil {
		return err
	}

	// Set DormantDatabase Phase: Paused
	t = metav1.Now()
	dormantDb.Status.PausingTime = &t
	dormantDb.Status.Phase = tapi.DormantDatabasePhasePaused
	if _, err = c.extClient.DormantDatabases(dormantDb.Namespace).Update(dormantDb); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			"Failed to update DormantDatabase. Reason: %v",
			err,
		)
		return err
	}

	return nil
}

func (c *DormantDbController) delete(dormantDb *tapi.DormantDatabase) error {
	phase := dormantDb.Status.Phase
	if phase != tapi.DormantDatabasePhaseResuming && phase != tapi.DormantDatabasePhaseWipedOut {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			`DormantDatabase "%v" is not %v.`,
			dormantDb.Name,
			tapi.DormantDatabasePhaseWipedOut,
		)

		if err := c.reCreateDormantDatabase(dormantDb); err != nil {
			c.eventRecorder.Eventf(
				dormantDb,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToCreate,
				`Failed to recreate DormantDatabase: "%v". Reason: %v`,
				dormantDb.Name,
				err,
			)
			return err
		}
	}
	return nil
}

func (c *DormantDbController) update(oldDormantDb, updatedDormantDb *tapi.DormantDatabase) error {
	if oldDormantDb.Spec.WipeOut != updatedDormantDb.Spec.WipeOut && updatedDormantDb.Spec.WipeOut {
		return c.wipeOut(updatedDormantDb)
	}

	if oldDormantDb.Spec.Resume != updatedDormantDb.Spec.Resume && updatedDormantDb.Spec.Resume {
		if oldDormantDb.Status.Phase == tapi.DormantDatabasePhasePaused {
			return c.resume(updatedDormantDb)
		} else {
			message := "Failed to resume Database. " +
				"Only DormantDatabase of \"Paused\" Phase can be resumed"
			c.eventRecorder.Event(
				updatedDormantDb,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToUpdate,
				message,
			)
		}
	}
	return nil
}

func (c *DormantDbController) wipeOut(dormantDb *tapi.DormantDatabase) error {
	// Check if DB TPR object exists
	found, err := c.deleter.Exists(&dormantDb.ObjectMeta)
	if err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			"Failed to wipeOut Database. Reason: %v",
			err,
		)
		return err
	}

	if found {
		message := "Failed to wipeOut Database. Delete Database TPR object first"
		c.eventRecorder.Event(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToWipeOut,
			message,
		)

		// Delete DormantDatabase object
		if err := c.extClient.DormantDatabases(dormantDb.Namespace).Delete(dormantDb.Name); err != nil {
			c.eventRecorder.Eventf(
				dormantDb,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				"Failed to delete DormantDatabase. Reason: %v",
				err,
			)
			log.Errorln(err)
		}
		return errors.New(message)
	}

	if dormantDb, err = c.extClient.DormantDatabases(dormantDb.Namespace).Get(dormantDb.Name); err != nil {
		return err
	}

	// Set DormantDatabase Phase: Wiping out
	t := metav1.Now()
	dormantDb.Status.Phase = tapi.DormantDatabasePhaseWipingOut

	if _, err := c.extClient.DormantDatabases(dormantDb.Namespace).Update(dormantDb); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			"Failed to update DormantDatabase. Reason: %v",
			err,
		)
		return err
	}

	// Wipe out Database workload
	c.eventRecorder.Event(dormantDb, apiv1.EventTypeNormal, eventer.EventReasonWipingOut, "Wiping out Database")
	if err := c.deleter.WipeOutDatabase(dormantDb); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToWipeOut,
			"Failed to wipeOut. Reason: %v",
			err,
		)
		return err
	}

	c.eventRecorder.Event(
		dormantDb,
		apiv1.EventTypeNormal,
		eventer.EventReasonSuccessfulWipeOut,
		"Successfully wiped out Database workload",
	)

	if dormantDb, err = c.extClient.DormantDatabases(dormantDb.Namespace).Get(dormantDb.Name); err != nil {
		return err
	}

	// Set DormantDatabase Phase: Deleted
	t = metav1.Now()
	dormantDb.Status.WipeOutTime = &t
	dormantDb.Status.Phase = tapi.DormantDatabasePhaseWipedOut
	if _, err = c.extClient.DormantDatabases(dormantDb.Namespace).Update(dormantDb); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			"Failed to update DormantDatabase. Reason: %v",
			err,
		)
		return err
	}

	return nil
}

func (c *DormantDbController) resume(dormantDb *tapi.DormantDatabase) error {
	c.eventRecorder.Event(
		dormantDb,
		apiv1.EventTypeNormal,
		eventer.EventReasonResuming,
		"Resuming DormantDatabase",
	)

	// Check if DB TPR object exists
	found, err := c.deleter.Exists(&dormantDb.ObjectMeta)
	if err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToResume,
			"Failed to resume DormantDatabase. Reason: %v",
			err,
		)
		return err
	}

	if found {
		message := "Failed to resume DormantDatabase. One Database TPR object exists with same name"
		c.eventRecorder.Event(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToResume,
			message,
		)
		return errors.New(message)
	}

	if dormantDb, err = c.extClient.DormantDatabases(dormantDb.Namespace).Get(dormantDb.Name); err != nil {
		return err
	}

	_dormantDb := dormantDb
	_dormantDb.Status.Phase = tapi.DormantDatabasePhaseResuming
	if _, err = c.extClient.DormantDatabases(_dormantDb.Namespace).Update(_dormantDb); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			"Failed to update DormantDatabase. Reason: %v",
			err,
		)
		return err
	}

	if err = c.extClient.DormantDatabases(dormantDb.Namespace).Delete(dormantDb.Name); err != nil {
		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			"Failed to delete DormantDatabase. Reason: %v",
			err,
		)
		return err
	}

	if err = c.deleter.ResumeDatabase(dormantDb); err != nil {
		if err := c.reCreateDormantDatabase(dormantDb); err != nil {
			c.eventRecorder.Eventf(
				dormantDb,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToCreate,
				`Failed to recreate DormantDatabase: "%v". Reason: %v`,
				dormantDb.Name,
				err,
			)
			return err
		}

		c.eventRecorder.Eventf(
			dormantDb,
			apiv1.EventTypeWarning,
			eventer.EventReasonFailedToResume,
			"Failed to resume Database. Reason: %v",
			err,
		)
		return err
	}
	return nil
}

func (c *DormantDbController) reCreateDormantDatabase(dormantDb *tapi.DormantDatabase) error {
	_dormantDb := &tapi.DormantDatabase{
		ObjectMeta: metav1.ObjectMeta{
			Name:        dormantDb.Name,
			Namespace:   dormantDb.Namespace,
			Labels:      dormantDb.Labels,
			Annotations: dormantDb.Annotations,
		},
		Spec:   dormantDb.Spec,
		Status: dormantDb.Status,
	}

	if _, err := c.extClient.DormantDatabases(_dormantDb.Namespace).Create(_dormantDb); err != nil {
		return err
	}

	return nil
}

func dormantDbSuccessfullyCreated() {
	analytics.SendEvent(tapi.ResourceNameDormantDatabase, "created", "success")
}

func dormantDbFailedToCreate() {
	analytics.SendEvent(tapi.ResourceNameDormantDatabase, "created", "failure")
}

func dormantDbSuccessfullyDeleted() {
	analytics.SendEvent(tapi.ResourceNameDormantDatabase, "deleted", "success")
}

func dormantDbFailedToDelete() {
	analytics.SendEvent(tapi.ResourceNameDormantDatabase, "deleted", "failure")
}
