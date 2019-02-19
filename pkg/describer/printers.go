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
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Pass ports=nil for all ports.
func formatEndpoints(endpoints *core.Endpoints, ports sets.String) string {
	if len(endpoints.Subsets) == 0 {
		return "<none>"
	}
	list := []string{}
	max := 3
	more := false
	count := 0
	for i := range endpoints.Subsets {
		ss := &endpoints.Subsets[i]
		for i := range ss.Ports {
			port := &ss.Ports[i]
			if ports == nil || ports.Has(port.Name) {
				for i := range ss.Addresses {
					if len(list) == max {
						more = true
					}
					addr := &ss.Addresses[i]
					if !more {
						hostPort := net.JoinHostPort(addr.IP, strconv.Itoa(int(port.Port)))
						list = append(list, hostPort)
					}
					count++
				}
			}
		}
	}
	ret := strings.Join(list, ",")
	if more {
		return fmt.Sprintf("%s + %d more...", ret, count-max)
	}
	return ret
}

// formatEventSource formats EventSource as a comma separated string excluding Host when empty
func formatEventSource(es core.EventSource) string {
	EventSourceString := []string{es.Component}
	if len(es.Host) > 0 {
		EventSourceString = append(EventSourceString, es.Host)
	}
	return strings.Join(EventSourceString, ", ")
}

// translateTimestamp returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestamp(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.ShortHumanDuration(time.Since(timestamp.Time))
}

func timeToString(t *metav1.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}

	return t.Format(time.RFC1123Z)
}
