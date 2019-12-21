/*
Copyright The KubeDB Authors.

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

package v1alpha1

import (
	"fmt"
	"path/filepath"

	"kubedb.dev/apimachinery/api/crds"
	"kubedb.dev/apimachinery/apis"

	"github.com/pkg/errors"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (_ Snapshot) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralSnapshot))
}

var _ apis.ResourceInfo = &Snapshot{}

func (s Snapshot) OffshootName() string {
	return s.Name
}

func (s Snapshot) Location() (string, error) {
	spec := s.Spec.Backend
	if spec.S3 != nil {
		return filepath.Join(spec.S3.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.GCS != nil {
		return filepath.Join(spec.GCS.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Azure != nil {
		return filepath.Join(spec.Azure.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Local != nil {
		return filepath.Join(DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Swift != nil {
		return filepath.Join(spec.Swift.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	}
	return "", errors.New("no storage provider is configured")
}

func (s Snapshot) ResourceShortCode() string {
	return ResourceCodeSnapshot
}

func (s Snapshot) ResourceKind() string {
	return ResourceKindSnapshot
}

func (s Snapshot) ResourceSingular() string {
	return ResourceSingularSnapshot
}

func (s Snapshot) ResourcePlural() string {
	return ResourcePluralSnapshot
}

func (s Snapshot) OSMSecretName() string {
	return fmt.Sprintf("osm-%v", s.Name)
}

func (s *Snapshot) SetDefaults() {
	if s == nil {
		return
	}
	// Add snapshot defaulting here
}
