/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

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
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Pass ports=nil for all ports.
func formatEndpointSlices(endpointSlices []discoveryv1.EndpointSlice, ports sets.Set[string]) string {
	if len(endpointSlices) == 0 {
		return "<none>"
	}
	var list []string
	max := 3
	more := false
	count := 0
	for i := range endpointSlices {
		if len(endpointSlices[i].Ports) == 0 {
			// It's possible to have headless services with no ports.
			for j := range endpointSlices[i].Endpoints {
				if len(list) == max {
					more = true
				}
				isReady := endpointSlices[i].Endpoints[j].Conditions.Ready == nil || *endpointSlices[i].Endpoints[j].Conditions.Ready
				if !isReady {
					// ready indicates that this endpoint is prepared to receive traffic,
					// according to whatever system is managing the endpoint. A nil value
					// indicates an unknown state. In most cases consumers should interpret this
					// unknown state as ready.
					// More info: vendor/k8s.io/api/discovery/v1/types.go
					continue
				}
				if !more {
					list = append(list, endpointSlices[i].Endpoints[j].Addresses[0])
				}
				count++
			}
		} else {
			// "Normal" services with ports defined.
			for j := range endpointSlices[i].Ports {
				port := endpointSlices[i].Ports[j]
				if ports == nil || ports.Has(*port.Name) {
					for k := range endpointSlices[i].Endpoints {
						if len(list) == max {
							more = true
						}
						addr := endpointSlices[i].Endpoints[k].Addresses[0]
						isReady := endpointSlices[i].Endpoints[k].Conditions.Ready == nil || *endpointSlices[i].Endpoints[k].Conditions.Ready
						if !isReady {
							// ready indicates that this endpoint is prepared to receive traffic,
							// according to whatever system is managing the endpoint. A nil value
							// indicates an unknown state. In most cases consumers should interpret this
							// unknown state as ready.
							// More info: vendor/k8s.io/api/discovery/v1/types.go
							continue
						}
						if !more {
							hostPort := net.JoinHostPort(addr, strconv.Itoa(int(*port.Port)))
							list = append(list, hostPort)
						}
						count++
					}
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
