package eventer

import (
	"github.com/appscode/log"
	kapi "k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/record"
)

const (
	EventReasonCreating                string = "Creating"
	EventReasonPausing                 string = "Pausing"
	EventReasonWipingOut               string = "WipingOut"
	EventReasonFailedToCreate          string = "Failed"
	EventReasonFailedToPause           string = "Failed"
	EventReasonFailedToDelete          string = "Failed"
	EventReasonFailedToWipeOut         string = "Failed"
	EventReasonFailedToGet             string = "Failed"
	EventReasonFailedToInitialize      string = "Failed"
	EventReasonFailedToList            string = "Failed"
	EventReasonFailedToResume          string = "Failed"
	EventReasonFailedToSchedule        string = "Failed"
	EventReasonFailedToStart           string = "Failed"
	EventReasonFailedToUpdate          string = "Failed"
	EventReasonFailedToAddMonitor      string = "Failed"
	EventReasonFailedToDeleteMonitor   string = "Failed"
	EventReasonFailedToUpdateMonitor   string = "Failed"
	EventReasonIgnoredSnapshot         string = "IgnoredSnapshot"
	EventReasonInitializing            string = "Initializing"
	EventReasonInvalid                 string = "Invalid"
	EventReasonInvalidUpdate           string = "InvalidUpdate"
	EventReasonResuming                string = "Resuming"
	EventReasonSnapshotFailed          string = "SnapshotFailed"
	EventReasonStarting                string = "Starting"
	EventReasonSuccessfulCreate        string = "SuccessfulCreate"
	EventReasonSuccessfulPause         string = "SuccessfulPause"
	EventReasonSuccessfulMonitorAdd    string = "SuccessfulMonitorAdd"
	EventReasonSuccessfulMonitorDelete string = "SuccessfulMonitorDelete"
	EventReasonSuccessfulMonitorUpdate string = "SuccessfulMonitorUpdate"
	EventReasonSuccessfulResume        string = "SuccessfulResume"
	EventReasonSuccessfulWipeOut       string = "SuccessfulWipeOut"
	EventReasonSuccessfulSnapshot      string = "SuccessfulSnapshot"
	EventReasonSuccessfulValidate      string = "SuccessfulValidate"
	EventReasonSuccessfulInitialize    string = "SuccessfulInitialize"
)

func NewEventRecorder(client clientset.Interface, component string) record.EventRecorder {
	// Event Broadcaster
	broadcaster := record.NewBroadcaster()
	broadcaster.StartEventWatcher(
		func(event *kapi.Event) {
			if _, err := client.Core().Events(event.Namespace).Create(event); err != nil {
				log.Errorln(err)
			}
		},
	)

	return broadcaster.NewRecorder(kapi.EventSource{Component: component})
}
