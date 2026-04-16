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

package double_optin

import (
	"context"
	"fmt"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CheckIfDoubleOptInPossible is the intended function to be called from operators.
// In Archiver - DB connection, DB is the requester, Archiver is the allower
// In Schema - DB connection, Schema is the requester, DB is the allower
func CheckIfDoubleOptInPossible(requesterMeta, requesterNSMeta, allowerNSMeta metav1.ObjectMeta, consumerSpec *dbapi.AllowedConsumers) (bool, error) {
	if consumerSpec == nil {
		return false, nil
	}
	matchNamespace, err := IsInAllowedNamespaces(requesterNSMeta, allowerNSMeta, consumerSpec)
	if err != nil {
		return false, err
	}
	matchLabels, err := IsMatchByLabels(requesterMeta, consumerSpec)
	if err != nil {
		return false, err
	}
	return matchNamespace && matchLabels, nil
}

func IsInAllowedNamespaces(requesterNSMeta, allowerNSMeta metav1.ObjectMeta, consumerSpec *dbapi.AllowedConsumers) (bool, error) {
	if consumerSpec.Namespaces == nil || consumerSpec.Namespaces.From == nil {
		return false, nil
	}

	if *consumerSpec.Namespaces.From == dbapi.NamespacesFromAll {
		return true, nil
	}
	if *consumerSpec.Namespaces.From == dbapi.NamespacesFromSame {
		return requesterNSMeta.GetName() == allowerNSMeta.GetName(), nil
	}
	if *consumerSpec.Namespaces.From == dbapi.NamespacesFromSelector {
		if consumerSpec.Namespaces.Selector == nil {
			// this says, Select namespace from the Selector, but the Namespace.Selector field is nil. So, no way to select namespace here.
			return false, nil
		}
		ret, err := selectorMatches(consumerSpec.Namespaces.Selector, requesterNSMeta.GetLabels())
		if err != nil {
			return false, err
		}
		return ret, nil
	}
	return false, nil
}

func IsMatchByLabels(requesterMeta metav1.ObjectMeta, consumerSpec *dbapi.AllowedConsumers) (bool, error) {
	if consumerSpec.Selector != nil {
		ret, err := selectorMatches(consumerSpec.Selector, requesterMeta.Labels)
		if err != nil {
			return false, err
		}
		return ret, nil
	}
	// if Selector is not given, all the Schemas are allowed of the selected namespace
	return true, nil
}

func selectorMatches(ls *metav1.LabelSelector, srcLabels map[string]string) (bool, error) {
	selector, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		klog.Infoln("invalid selector: ", ls)
		return false, err
	}
	return selector.Matches(labels.Set(srcLabels)), nil
}

// ListConsumers lists all the consumers of a corresponding gvk
// In Archiver - DB connection, DB is the requester, Archiver is the allower
// In Schema - DB connection, Schema is the requester, DB is the allower
func ListConsumers(kc client.Client, allowerMeta metav1.ObjectMeta, requesterGVK schema.GroupVersionKind, consumerSpec *dbapi.AllowedConsumers) (*unstructured.UnstructuredList, error) {
	nsList, err := getAllowedNamespaceList(kc, consumerSpec)
	if err != nil {
		klog.Warningf("failed to get allowed namespace list for archiver: %s/%s. Reason: %v", allowerMeta.Namespace, allowerMeta.Name, err)
		return nil, err
	}

	consumerList, err := listAllConsumers(kc, requesterGVK, consumerSpec)
	if err != nil {
		klog.Warningf("failed to list dbs for archiver: %s/%s. Reason: %v", allowerMeta.Namespace, allowerMeta.Name, err)
		return nil, err
	}

	isAllowed := func(requesterNS, allowerNs string, from dbapi.FromNamespaces, namespaceAllowlist map[string]struct{}) bool {
		if namespaceAllowlist != nil {
			_, ok := namespaceAllowlist[requesterNS]
			return ok
		}
		return from != dbapi.NamespacesFromSame || requesterNS == allowerNs
	}

	for _, c := range consumerList.Items {
		if !isAllowed(c.GetNamespace(), allowerMeta.Namespace, *consumerSpec.Namespaces.From, nsList) {
			continue
		}
	}
	return consumerList, nil
}

func getAllowedNamespaceList(kc client.Client, consumerSpec *dbapi.AllowedConsumers) (map[string]struct{}, error) {
	if *consumerSpec.Namespaces.From != dbapi.NamespacesFromSelector {
		return nil, nil
	}
	nsSelector, err := metav1.LabelSelectorAsSelector(consumerSpec.Namespaces.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to converting namespace selector. Reason: %v", err)
	}

	nsList := &core.NamespaceList{}
	err = kc.List(context.TODO(), nsList, client.MatchingLabelsSelector{Selector: nsSelector})
	if err != nil {
		return nil, fmt.Errorf("failed to listing namespaces. Reason: %v", err)
	}

	allowlist := make(map[string]struct{}, len(nsList.Items))
	for _, ns := range nsList.Items {
		allowlist[ns.Name] = struct{}{}
	}

	return allowlist, nil
}

func listAllConsumers(kc client.Client, requesterGVK schema.GroupVersionKind, consumerSpec *dbapi.AllowedConsumers) (*unstructured.UnstructuredList, error) {
	sel, err := metav1.LabelSelectorAsSelector(consumerSpec.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to converting namespace selector. Reason: %v", err)
	}

	consumerList := &unstructured.UnstructuredList{}
	consumerList.SetGroupVersionKind(requesterGVK)
	if err := kc.List(context.TODO(), consumerList, []client.ListOption{
		client.MatchingLabelsSelector{Selector: sel},
	}...); err != nil {
		return nil, fmt.Errorf("failed to list databases. Reason: %v", err)
	}

	return consumerList, nil
}
