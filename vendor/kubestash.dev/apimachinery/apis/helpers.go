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

package apis

import (
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	once sync.Once
	kc   client.Client
)

func GetRuntimeClient() client.Client {
	if kc == nil {
		panic("runtime client is not initialized!")
	}
	return kc
}

func SetRuntimeClient(client client.Client) {
	once.Do(func() {
		kc = client
	})
}

func UpsertLabels(oldLabels, newLabels map[string]string) map[string]string {
	if oldLabels == nil {
		oldLabels = make(map[string]string, len(newLabels))
	}
	for k, v := range newLabels {
		oldLabels[k] = v
	}
	return oldLabels
}
