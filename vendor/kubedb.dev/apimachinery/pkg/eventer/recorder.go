/*
Copyright AppsCode Inc. and Contributors

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

package eventer

import (
	"context"
	"fmt"
	"time"

	"gomodules.xyz/x/log"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
)

const (
	EventReasonPausing                          string = "Pausing"
	EventReasonWipingOut                        string = "WipingOut"
	EventReasonFailedToCreate                   string = "Failed"
	EventReasonFailedToPause                    string = "Failed"
	EventReasonFailedToDelete                   string = "Failed"
	EventReasonFailedToWipeOut                  string = "Failed"
	EventReasonFailedToGet                      string = "Failed"
	EventReasonFailedToInitialize               string = "Failed"
	EventReasonFailedToList                     string = "Failed"
	EventReasonFailedToResume                   string = "Failed"
	EventReasonFailedToSchedule                 string = "Failed"
	EventReasonFailedToStart                    string = "Failed"
	EventReasonFailedToUpdate                   string = "Failed"
	EventReasonIgnoredSnapshot                  string = "IgnoredSnapshot"
	EventReasonInitializing                     string = "Initializing"
	EventReasonInvalid                          string = "Invalid"
	EventReasonResuming                         string = "Resuming"
	EventReasonSnapshotFailed                   string = "SnapshotFailed"
	EventReasonSnapshotError                    string = "SnapshotError"
	EventReasonStarting                         string = "Starting"
	EventReasonSuccessful                       string = "Successful"
	EventReasonSuccessfulCreate                 string = "SuccessfulCreate"
	EventReasonSuccessfulPause                  string = "SuccessfulPause"
	EventReasonSuccessfulWipeOut                string = "SuccessfulWipeOut"
	EventReasonSuccessfulSnapshot               string = "SuccessfulSnapshot"
	EventReasonSuccessfulInitialize             string = "SuccessfulInitialize"
	EventReasonAdmissionWebhookNotActivated     string = "AdmissionWebhookNotActivated"
	EventReasonSuccessfulDatabaseInitialization string = "DatabaseInitializedSuccessfully"
	EventReasonFailureInDatabaseInitialization  string = "FailedToInitializeDatabase"
)

func NewEventRecorder(client kubernetes.Interface, component string) record.EventRecorder {
	// Event Broadcaster
	broadcaster := record.NewBroadcaster()
	broadcaster.StartEventWatcher(
		func(event *core.Event) {
			if _, err := client.CoreV1().Events(event.Namespace).Create(context.TODO(), event, metav1.CreateOptions{}); err != nil {
				log.Errorln(err)
			}
		},
	)

	return broadcaster.NewRecorder(scheme.Scheme, core.EventSource{Component: component})
}

func CreateEvent(client kubernetes.Interface, component string, obj runtime.Object, eventType, reason, message string) (*core.Event, error) {
	ref, err := reference.GetReference(scheme.Scheme, obj)
	if err != nil {
		return nil, err
	}

	t := metav1.Time{Time: time.Now()}

	return client.CoreV1().Events(ref.Namespace).Create(context.TODO(), &core.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v.%x", ref.Name, t.UnixNano()),
			Namespace: ref.Namespace,
		},
		InvolvedObject: *ref,
		Reason:         reason,
		Message:        message,
		FirstTimestamp: t,
		LastTimestamp:  t,
		Count:          1,
		Type:           eventType,
		Source:         core.EventSource{Component: component},
	}, metav1.CreateOptions{})
}

func CreateEventWithLog(client kubernetes.Interface, component string, obj runtime.Object, eventType, reason, message string) {
	event, err := CreateEvent(client, component, obj, eventType, reason, message)
	if err != nil {
		log.Errorln("Failed to write event, reason: ", err)
	} else {
		log.Infoln("Event created: ", event.Name)
	}
}
