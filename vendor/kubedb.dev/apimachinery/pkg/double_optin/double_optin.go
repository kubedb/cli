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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
)

// CheckIfDoubleOptInPossible is the intended function to be called from operator
// It checks if the namespace, where SchemaDatabase or Archiver is applied, is allowed.
// It also checks the labels of schemaDatabase OR archiver, to decide if that is allowed or not.
//
// Here, clientMeta is the ObjectMeta of SchemaDatabase or Archiver
// & clientNSMeta is the ObjectMeta of the namespace where they belong.
func CheckIfDoubleOptInPossible(clientMeta metav1.ObjectMeta, clientNSMeta metav1.ObjectMeta, dbNSMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers == nil {
		return false, nil
	}
	matchNamespace, err := IsInAllowedNamespaces(clientNSMeta, dbNSMeta, consumers)
	if err != nil {
		return false, err
	}
	matchLabels, err := IsMatchByLabels(clientMeta, consumers)
	if err != nil {
		return false, err
	}
	return matchNamespace && matchLabels, nil
}

func IsInAllowedNamespaces(clientNSMeta metav1.ObjectMeta, dbNSMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers.Namespaces == nil || consumers.Namespaces.From == nil {
		return false, nil
	}

	if *consumers.Namespaces.From == dbapi.NamespacesFromAll {
		return true, nil
	}
	if *consumers.Namespaces.From == dbapi.NamespacesFromSame {
		return clientNSMeta.GetName() == dbNSMeta.GetName(), nil
	}
	if *consumers.Namespaces.From == dbapi.NamespacesFromSelector {
		if consumers.Namespaces.Selector == nil {
			// this says, Select namespace from the Selector, but the Namespace.Selector field is nil. So, no way to select namespace here.
			return false, nil
		}
		ret, err := selectorMatches(consumers.Namespaces.Selector, clientNSMeta.GetLabels())
		if err != nil {
			return false, err
		}
		return ret, nil
	}
	return false, nil
}

func IsMatchByLabels(clientMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers.Selector != nil {
		ret, err := selectorMatches(consumers.Selector, clientMeta.Labels)
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
