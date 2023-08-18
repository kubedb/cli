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

package factory

import (
	kubedbscheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"

	cmscheme "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned/scheme"
	promscheme "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/scheme"
	crdscheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	aggscheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme"
	metricsscheme "k8s.io/metrics/pkg/client/clientset/versioned/scheme"
	crscheme "kmodules.xyz/custom-resources/client/clientset/versioned/scheme"
	sidekickapi "kubeops.dev/sidekick/apis/apps/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	stashscheme "stash.appscode.dev/apimachinery/client/clientset/versioned/scheme"
)

func NewUncachedClient(cfg *rest.Config) (client.Client, error) {
	mapper, err := apiutil.NewDynamicRESTMapper(cfg)
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := kubedbscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := sidekickapi.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// crd
	if err := crdscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// apiservices
	if err := aggscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// metrics
	if err := metricsscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// appbinding
	if err := crscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// cert-manager
	if err := cmscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// stash
	if err := stashscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// prometheus
	if err := promscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		//Opts: client.WarningHandlerOptions{
		//	SuppressWarnings:   false,
		//	AllowDuplicateLogs: false,
		//},
	})
}

func MustUncachedClient(cfg *rest.Config) client.Client {
	c, err := NewUncachedClient(cfg)
	if err != nil {
		panic(err)
	}
	return c
}
