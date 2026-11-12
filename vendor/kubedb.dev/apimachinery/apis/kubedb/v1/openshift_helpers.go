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

package v1

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OpenShift assigns every namespace a block of UIDs and supplemental groups and records them
// in these annotations. Its restricted-v2 SCC only admits pods whose runAsUser lies in the
// uid-range block and whose fsGroup lies in the supplemental-groups block, so a catalog UID
// such as postgres's 70 is rejected outright.
const (
	OpenShiftSCCUIDRangeAnnotation           = "openshift.io/sa.scc.uid-range"
	OpenShiftSCCSupplementalGroupsAnnotation = "openshift.io/sa.scc.supplemental-groups"
)

// NamespaceSecurityIDs are the IDs a database pod should run with in a namespace that carries
// OpenShift's SCC annotations. They are defaults only: a value the user set on the pod or
// container securityContext always wins.
// +k8s:deepcopy-gen=false
// +k8s:openapi-gen=false
type NamespaceSecurityIDs struct {
	RunAsUser  int64
	RunAsGroup int64
	FSGroup    int64
}

// GetNamespaceSecurityIDs reads the OpenShift SCC annotations from a namespace.
//
// It returns nil, nil when the namespace has no uid-range annotation, which is the case on
// every non-OpenShift cluster; callers then keep the catalog defaults.
//
//   - RunAsUser and RunAsGroup are the first ID of the uid-range block.
//   - FSGroup is the first ID of the supplemental-groups block, or of the uid-range block when
//     the namespace has no supplemental-groups annotation.
func GetNamespaceSecurityIDs(ns *core.Namespace) (*NamespaceSecurityIDs, error) {
	if ns == nil {
		return nil, nil
	}
	uidRange, ok := ns.Annotations[OpenShiftSCCUIDRangeAnnotation]
	if !ok {
		return nil, nil
	}
	uid, err := firstIDOfSCCRange(uidRange)
	if err != nil {
		return nil, fmt.Errorf("namespace %s: invalid %s annotation %q: %w", ns.Name, OpenShiftSCCUIDRangeAnnotation, uidRange, err)
	}
	ids := &NamespaceSecurityIDs{RunAsUser: uid, RunAsGroup: uid, FSGroup: uid}

	if groups, ok := ns.Annotations[OpenShiftSCCSupplementalGroupsAnnotation]; ok {
		gid, err := firstIDOfSCCRange(groups)
		if err != nil {
			return nil, fmt.Errorf("namespace %s: invalid %s annotation %q: %w", ns.Name, OpenShiftSCCSupplementalGroupsAnnotation, groups, err)
		}
		ids.FSGroup = gid
	}
	return ids, nil
}

// firstIDOfSCCRange returns the first ID of the first block in an OpenShift SCC range
// annotation. A block is written either as "start/size" or as "start-end", and the
// supplemental-groups annotation may hold several comma separated blocks.
func firstIDOfSCCRange(value string) (int64, error) {
	block := strings.TrimSpace(strings.Split(value, ",")[0])

	var startStr, restStr string
	var isSize bool
	if s, r, found := strings.Cut(block, "/"); found {
		startStr, restStr, isSize = s, r, true
	} else if s, r, found := strings.Cut(block, "-"); found {
		startStr, restStr = s, r
	} else {
		return 0, fmt.Errorf("expected <start>/<size> or <start>-<end>")
	}

	start, err := strconv.ParseInt(strings.TrimSpace(startStr), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid start: %w", err)
	}
	rest, err := strconv.ParseInt(strings.TrimSpace(restStr), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid range: %w", err)
	}
	if start <= 0 {
		return 0, fmt.Errorf("start must be a non-root ID, got %d", start)
	}
	if isSize && rest <= 0 {
		return 0, fmt.Errorf("size must be positive, got %d", rest)
	}
	if !isSize && rest < start {
		return 0, fmt.Errorf("end %d is before start %d", rest, start)
	}
	return start, nil
}

// openShiftByMapper caches IsOpenShift per RESTMapper. A cluster does not switch between
// OpenShift and plain Kubernetes while an operator runs.
var openShiftByMapper sync.Map // meta.RESTMapper -> bool

// IsOpenShift reports whether the cluster serves OpenShift's project.openshift.io Project kind,
// the same check as kmodules' cluster.IsOpenShiftManaged, but remembered per RESTMapper.
//
// The cache matters on clusters that are not OpenShift: controller-runtime's RESTMapper does not
// remember a group the server lacks, so every uncached lookup re-fetches /api and /apis while
// holding the mapper's write lock, and operators call this on every reconcile.
//
// Only a definite answer is cached: the kind was found (true) or the server has no such kind
// (a NoMatch error, false). Any other error, such as a failed discovery request, returns false
// for this call and is retried on the next one, so a transient failure on OpenShift is not
// remembered as "not OpenShift".
func IsOpenShift(mapper meta.RESTMapper) bool {
	if v, ok := openShiftByMapper.Load(mapper); ok {
		return v.(bool)
	}
	_, err := mapper.RESTMappings(schema.GroupKind{Group: "project.openshift.io", Kind: "Project"})
	switch {
	case err == nil:
		openShiftByMapper.Store(mapper, true)
		return true
	case meta.IsNoMatchError(err):
		openShiftByMapper.Store(mapper, false)
		return false
	default:
		return false
	}
}

// GetNamespaceSecurityIDsFor reads a namespace and returns its OpenShift SCC IDs, or nil when it
// carries no SCC annotations.
func GetNamespaceSecurityIDsFor(ctx context.Context, kc client.Reader, namespace string) (*NamespaceSecurityIDs, error) {
	var ns core.Namespace
	if err := kc.Get(ctx, types.NamespacedName{Name: namespace}, &ns); err != nil {
		return nil, fmt.Errorf("failed to get namespace %s: %w", namespace, err)
	}
	return GetNamespaceSecurityIDs(&ns)
}

// SetDefaultsInNamespace reads the namespace's OpenShift SCC IDs and runs
// SetDefaultsWithNamespaceSecurityIDs with them. It is what the mutating webhooks call in place
// of a database's plain SetDefaults; setDefaults's error, if any, is returned.
//
// On a cluster that is not OpenShift it is exactly setDefaults(): the namespace is not read.
func SetDefaultsInNamespace(ctx context.Context, kc client.Client, namespace string, spec any, setDefaults func() error) error {
	if !IsOpenShift(kc.RESTMapper()) {
		return setDefaults()
	}
	nsIDs, err := GetNamespaceSecurityIDsFor(ctx, kc, namespace)
	if err != nil {
		return err
	}
	SetDefaultsWithNamespaceSecurityIDs(spec, nsIDs, func() { err = setDefaults() })
	return err
}

// SetDefaultsInNamespaceForVersion is SetDefaultsInNamespace for databases whose catalog
// version can pin the UID (the RunAsAnyNonRoot field of Postgres, Elasticsearch, PgBouncer,
// Pgpool and Etcd versions). When runAsAnyNonRoot is false the image only works as the catalog
// UID and the validating webhooks reject any other, so the catalog defaults are kept exactly as
// before and the namespace is not read. Such a version still cannot run under restricted-v2.
func SetDefaultsInNamespaceForVersion(ctx context.Context, kc client.Client, namespace string, runAsAnyNonRoot bool, spec any, setDefaults func() error) error {
	if !runAsAnyNonRoot {
		return setDefaults()
	}
	return SetDefaultsInNamespace(ctx, kc, namespace, spec, setDefaults)
}

// SetDefaultsWithNamespaceSecurityIDs runs a database's own SetDefaults and then makes the
// namespace IDs replace the catalog IDs it filled in. spec must be a pointer to the database's
// spec; setDefaults is the database's usual defaulting call, e.g.
//
//	SetDefaultsWithNamespaceSecurityIDs(&db.Spec, nsIDs, func() { db.SetDefaults(version) })
//
// Every pod and container securityContext reachable from spec is covered: pod templates,
// topology nodes, sidecars, the exporter and so on. A runAsUser, runAsGroup or fsGroup is
// replaced only when all of these hold:
//
//   - it was unset before setDefaults ran, so values the user set (or that an earlier
//     defaulting already stored on an existing database) are never changed;
//   - setDefaults filled it in, so a field left empty stays empty and the SCC assigns it;
//   - it is not 0 and the container is not privileged, since root is a deliberate choice that
//     only a privileged SCC can admit anyway.
//
// With nil nsIDs (any namespace without the annotations) this is exactly setDefaults().
func SetDefaultsWithNamespaceSecurityIDs(spec any, nsIDs *NamespaceSecurityIDs, setDefaults func()) {
	if nsIDs == nil {
		setDefaults()
		return
	}

	userSet := map[string]bool{}
	for path, f := range securityIDFields(spec) {
		if *f.ptr != nil {
			userSet[path] = true
		}
	}

	setDefaults()

	for path, f := range securityIDFields(spec) {
		if userSet[path] || *f.ptr == nil || **f.ptr == 0 {
			continue
		}
		id := nsIDs.RunAsUser
		switch f.kind {
		case securityIDGroup:
			id = nsIDs.RunAsGroup
		case securityIDFSGroup:
			id = nsIDs.FSGroup
		}
		// Assign a new pointer: defaulting often shares the catalog's own pointer, which must
		// not be modified.
		*f.ptr = ptr.To(id)
	}
}

type securityIDKind int

const (
	securityIDUser securityIDKind = iota
	securityIDGroup
	securityIDFSGroup
)

// securityIDField points at one runAsUser, runAsGroup or fsGroup field inside a spec.
// +k8s:deepcopy-gen=false
// +k8s:openapi-gen=false
type securityIDField struct {
	ptr  **int64
	kind securityIDKind
}

var (
	podSecurityContextType = reflect.TypeOf(core.PodSecurityContext{})
	securityContextType    = reflect.TypeOf(core.SecurityContext{})
	containerType          = reflect.TypeOf(core.Container{})
)

// securityIDFields returns the security ID fields reachable from spec, keyed by a path that
// stays the same across defaulting. Containers are keyed by name, because defaulting may add
// or reorder them.
func securityIDFields(spec any) map[string]securityIDField {
	out := map[string]securityIDField{}
	walkSecurityIDFields(reflect.ValueOf(spec), "", out)
	return out
}

func walkSecurityIDFields(v reflect.Value, path string, out map[string]securityIDField) {
	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			walkSecurityIDFields(v.Elem(), path, out)
		}
	case reflect.Struct:
		if !v.CanAddr() {
			return
		}
		switch v.Type() {
		case podSecurityContextType:
			sc := v.Addr().Interface().(*core.PodSecurityContext)
			out[path+".runAsUser"] = securityIDField{&sc.RunAsUser, securityIDUser}
			out[path+".runAsGroup"] = securityIDField{&sc.RunAsGroup, securityIDGroup}
			out[path+".fsGroup"] = securityIDField{&sc.FSGroup, securityIDFSGroup}
			return
		case securityContextType:
			sc := v.Addr().Interface().(*core.SecurityContext)
			if ptr.Deref(sc.Privileged, false) {
				return
			}
			out[path+".runAsUser"] = securityIDField{&sc.RunAsUser, securityIDUser}
			out[path+".runAsGroup"] = securityIDField{&sc.RunAsGroup, securityIDGroup}
			return
		}
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			if f := t.Field(i); f.IsExported() {
				walkSecurityIDFields(v.Field(i), path+"."+f.Name, out)
			}
		}
	case reflect.Slice, reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Struct, reflect.Pointer, reflect.Slice, reflect.Array:
		default:
			return // []byte, []string and the like hold no security contexts
		}
		isContainer := v.Type().Elem() == containerType
		for i := 0; i < v.Len(); i++ {
			key := strconv.Itoa(i)
			if isContainer {
				key = v.Index(i).Interface().(core.Container).Name
			}
			walkSecurityIDFields(v.Index(i), path+"["+key+"]", out)
		}
	}
}
