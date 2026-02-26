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

package v1alpha1

import (
	"encoding/json"
	"fmt"

	"kubedb.dev/apimachinery/crds"

	"k8s.io/apimachinery/pkg/runtime"
	"kmodules.xyz/client-go/apiextensions"
)

func GetFinalizer() string {
	return SchemeGroupVersion.Group
}

func (MongoDBArchiver) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMongoDBArchiver))
}

func (PostgresArchiver) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPostgresArchiver))
}

func (MySQLArchiver) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMySQLArchiver))
}

func (MariaDBArchiver) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMariaDBArchiver))
}

func (MSSQLServerArchiver) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMSSQLServerArchiver))
}

func SetDefaultLogBackupOptions(log *LogBackupOptions) *LogBackupOptions {
	if log == nil {
		log = &LogBackupOptions{
			SuccessfulLogHistoryLimit: 5,
			FailedLogHistoryLimit:     5,
		}
	}
	return log
}

func GetValueFromExtraArgs(args map[string]runtime.RawExtension, key string, valType any) (any, error) {
	var err error
	if val, ok := args[key]; ok {
		err = json.Unmarshal(val.Raw, &valType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value for key %s from extra arg maps. Reason: %w", key, err)
		}
		return valType, nil
	}

	return nil, fmt.Errorf("key %s not found in extra arg maps", key)
}

func SetKeyValueToExtraArgs(args map[string]runtime.RawExtension, key string, value any) error {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s to json. Reason: %w", key, err)
	}
	args[key] = runtime.RawExtension{
		Raw: jsonVal,
	}
	return nil
}
