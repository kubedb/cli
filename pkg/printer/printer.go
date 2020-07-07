/*
Copyright AppsCode Inc. and Contributors

Licensed under the PolyForm Noncommercial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/PolyForm-Noncommercial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package printer

// ref: k8s.io/kubectl/pkg/resource_printer.go

// DescriberSettings holds display configuration for each object
// describer to control what is printed.
type DescriberSettings struct {
	ShowEvents   bool
	ShowWorkload bool
	ShowSecret   bool
}
