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

package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"k8s.io/klog/v2"
)

var (
	domain string
	once   sync.Once
)

func FindDomain() string {
	once.Do(func() {
		var err error
		domain, err = findDomain()
		if err != nil {
			klog.Errorf("failed to find domain: %v", err)
			domain = "cluster.local"
		}
	})
	return domain
}

func findDomain() (string, error) {
	filePath := "/etc/resolv.conf"
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %v", filePath, err)
	}
	defer file.Close() // nolint:errcheck

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "search ") {
			fields := strings.FieldsSeq(line)
			// search demo.svc.cluster.local svc.cluster.local cluster.local
			// search demo.svc.cluster.local svc.cluster.local cluster.local lan
			for field := range fields {
				if strings.HasPrefix(field, "svc.") &&
					!strings.HasPrefix(field, "svc.svc.") {
					return strings.TrimPrefix(field, "svc."), nil
				}
			}
			return "", fmt.Errorf("failed to find domain: %s", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading %s: %v", filePath, err)
	}

	return "", fmt.Errorf("no suitable domain found in %s", filePath)
}
