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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
)

// CheckIfDoubleOptInPossible is the intended function to be called from operators.
// In Archiver - DB connection, DB is the requester, Archiver is the allower
// In Schema - DB connection, Schema is the requester, DB is the allower
func CheckIfDoubleOptInPossible(requesterMeta, requesterNSMeta, allowerNSMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers == nil {
		return false, nil
	}
	matchNamespace, err := IsInAllowedNamespaces(requesterNSMeta, allowerNSMeta, consumers)
	if err != nil {
		return false, err
	}
	matchLabels, err := IsMatchByLabels(requesterMeta, consumers)
	if err != nil {
		return false, err
	}
	return matchNamespace && matchLabels, nil
}

func IsInAllowedNamespaces(requesterNSMeta, allowerNSMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers.Namespaces == nil || consumers.Namespaces.From == nil {
		return false, nil
	}

	if *consumers.Namespaces.From == dbapi.NamespacesFromAll {
		return true, nil
	}
	if *consumers.Namespaces.From == dbapi.NamespacesFromSame {
		return requesterNSMeta.GetName() == allowerNSMeta.GetName(), nil
	}
	if *consumers.Namespaces.From == dbapi.NamespacesFromSelector {
		if consumers.Namespaces.Selector == nil {
			// this says, Select namespace from the Selector, but the Namespace.Selector field is nil. So, no way to select namespace here.
			return false, nil
		}
		ret, err := selectorMatches(consumers.Namespaces.Selector, requesterNSMeta.GetLabels())
		if err != nil {
			return false, err
		}
		return ret, nil
	}
	return false, nil
}

func IsMatchByLabels(requesterMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers.Selector != nil {
		ret, err := selectorMatches(consumers.Selector, requesterMeta.Labels)
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
