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

package adopt

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	annNetworks    = "k8s.v1.cni.cncf.io/networks"
	annNADResource = "k8s.v1.cni.cncf.io/resourceName"

	// The interface names the operator uses when spec.network.sriov leaves them unset.
	defaultPrimaryInterface = "net1"
	defaultGDSInterface     = "net2"
)

var nadGVR = schema.GroupVersionResource{Group: "k8s.cni.cncf.io", Version: "v1", Resource: "network-attachment-definitions"}

// netAttachment is one entry of the Multus networks annotation.
type netAttachment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Interface string `json:"interface,omitempty"`
}

// parseNetworks reads the annotation: a comma-separated list of [ns/]name[@iface],
// or a JSON array of {name, namespace, interface}.
func parseNetworks(v string) ([]netAttachment, error) {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(v, "[") {
		var out []netAttachment
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			return nil, err
		}
		return out, nil
	}
	var out []netAttachment
	for _, item := range strings.Split(v, ",") {
		if item = strings.TrimSpace(item); item == "" {
			continue
		}
		var a netAttachment
		item, a.Interface, _ = strings.Cut(item, "@")
		if ns, name, ok := strings.Cut(item, "/"); ok {
			a.Namespace, a.Name = ns, name
		} else {
			a.Name = item
		}
		out = append(out, a)
	}
	return out, nil
}

// deviceResources are the extended resources of the Milvus container that the
// Milvus object does not carry as they are: the chart's SR-IOV and RDMA devices.
func deviceResources(t core.PodTemplateSpec) (rdma, other []string) {
	c := milvusContainer(t)
	if c == nil {
		return nil, nil
	}
	seen := map[string]bool{}
	for _, l := range []core.ResourceList{c.Resources.Requests, c.Resources.Limits} {
		for _, n := range sortedResourceNames(l) {
			s := string(n)
			if carriedResource(n) || gpuResource(n) || seen[s] {
				continue
			}
			seen[s] = true
			if strings.HasPrefix(s, "rdma/") {
				rdma = append(rdma, s)
			} else {
				other = append(other, s)
			}
		}
	}
	sort.Strings(rdma)
	sort.Strings(other)
	return rdma, other
}

// mapSRIOV turns the Multus annotation and device resources of a chart pod into
// spec.network.sriov. nadResource returns the device resource a
// NetworkAttachmentDefinition names ("" when unknown). It returns nil and no
// reason when the pod attaches no secondary network, and nil with a reason when
// the settings cannot be mapped without guessing.
func mapSRIOV(ns string, t core.PodTemplateSpec, nadResource func(namespace, name string) string) (map[string]any, string) {
	ann := t.Annotations[annNetworks]
	if strings.TrimSpace(ann) == "" {
		return nil, ""
	}
	atts, err := parseNetworks(ann)
	if err != nil {
		return nil, fmt.Sprintf("the annotation %s does not parse: %v", annNetworks, err)
	}
	if len(atts) == 0 {
		return nil, ""
	}
	if len(atts) > 2 {
		return nil, fmt.Sprintf("it attaches %d networks; KubeDB takes one SR-IOV network and one GDS network", len(atts))
	}
	for _, a := range atts {
		if a.Namespace != "" && a.Namespace != ns {
			return nil, fmt.Sprintf("NetworkAttachmentDefinition %s/%s is in another namespace; KubeDB reads it from the namespace of the Milvus", a.Namespace, a.Name)
		}
	}
	rdma, other := deviceResources(t)
	if len(rdma) > 1 {
		return nil, fmt.Sprintf("it requests several RDMA devices (%s); KubeDB takes one", strings.Join(rdma, ", "))
	}

	left := append([]string(nil), other...)
	take := func(name string) bool {
		for i, r := range left {
			if r == name {
				left = append(left[:i], left[i+1:]...)
				return true
			}
		}
		return false
	}
	// A NetworkAttachmentDefinition names its device resource; use that, and the
	// single remaining resource only when it leaves no doubt.
	resourceFor := func(a netAttachment) string {
		if r := nadResource(ns, a.Name); r != "" && take(r) {
			return r
		}
		return ""
	}
	res := make([]string, len(atts))
	for i, a := range atts {
		res[i] = resourceFor(a)
	}
	for i := range res {
		if res[i] != "" {
			continue
		}
		if len(left) == 1 && countEmpty(res) == 1 {
			res[i], left = left[0], nil
			continue
		}
		return nil, fmt.Sprintf("cannot tell which device resource belongs to network %s (candidates: %s)", atts[i].Name, orNone(strings.Join(left, ", ")))
	}
	if len(left) > 0 {
		return nil, fmt.Sprintf("device resource %s is not tied to a network", strings.Join(left, ", "))
	}

	primary := map[string]any{"attachmentRef": atts[0].Name, "resourceName": res[0]}
	if atts[0].Interface != "" && atts[0].Interface != defaultPrimaryInterface {
		primary["interface"] = atts[0].Interface
	}
	if len(rdma) == 1 {
		primary["rdmaResourceName"] = rdma[0]
	}
	if len(atts) == 2 {
		gds := map[string]any{"attachmentRef": atts[1].Name, "resourceName": res[1]}
		if atts[1].Interface != "" && atts[1].Interface != defaultGDSInterface {
			gds["interface"] = atts[1].Interface
		}
		primary["gds"] = gds
	}
	return map[string]any{"sriov": primary}, ""
}

func countEmpty(s []string) int {
	n := 0
	for _, v := range s {
		if v == "" {
			n++
		}
	}
	return n
}

// planNetworks maps the chart's secondary networks of every unit into
// spec.network.sriov, and says what it could not map.
func (c *cluster) planNetworks(ctx context.Context, src *Source, fs *findings) {
	nad := func(namespace, name string) string {
		obj, err := c.dyn.Resource(nadGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return ""
		}
		return obj.GetAnnotations()[annNADResource]
	}
	src.Networks = map[string]map[string]any{}
	for _, u := range src.Units {
		t := src.Templates[u.key()]
		ann := strings.TrimSpace(t.Annotations[annNetworks])
		if ann == "" {
			continue
		}
		id := u.Kind + " " + u.Name
		network, why := mapSRIOV(c.ns, t, nad)
		if network == nil {
			if why == "" {
				continue
			}
			fs.warnf("%s attaches secondary networks (%s) but %s; set %s by hand before handing it over (see the GPU and SR-IOV guide)",
				id, ann, why, networkPath(u))
			continue
		}
		src.Networks[u.key()] = network
		s := network["sriov"].(map[string]any)
		msg := fmt.Sprintf("%s: network %s with device %s", id, s["attachmentRef"], s["resourceName"])
		if r, ok := s["rdmaResourceName"]; ok {
			msg += ", RDMA " + fmt.Sprint(r)
		}
		if g, ok := s["gds"].(map[string]any); ok {
			msg += fmt.Sprintf(", GDS network %s with device %s", g["attachmentRef"], g["resourceName"])
		}
		fs.infof("%s -> %s.sriov; KubeDB generates its own init container that advertises the net1 address", msg, networkPath(u))
	}
}

func networkPath(u Unit) string {
	if u.Group != "" {
		return "spec.topology.distributed." + u.Role + ".groups[" + u.Group + "].network"
	}
	return "spec.topology.distributed." + u.Role + ".network"
}
