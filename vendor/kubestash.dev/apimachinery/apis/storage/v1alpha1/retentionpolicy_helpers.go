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
	"errors"
	"strconv"
	"strings"
	"unicode"

	"kubestash.dev/apimachinery/apis"
	"kubestash.dev/apimachinery/crds"

	core "k8s.io/api/core/v1"
	"kmodules.xyz/client-go/apiextensions"
)

func (RetentionPolicy) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralRetentionPolicy))
}

func (r RetentionPolicy) UsageAllowed(srcNamespace *core.Namespace) bool {
	allowedNamespaces := r.Spec.UsagePolicy.AllowedNamespaces

	if allowedNamespaces.From == nil {
		return false
	}

	if *allowedNamespaces.From == apis.NamespacesFromAll {
		return true
	}

	if *allowedNamespaces.From == apis.NamespacesFromSame {
		return r.Namespace == srcNamespace.Name
	}

	return selectorMatches(allowedNamespaces.Selector, srcNamespace.Labels)
}

func (r RetentionPeriod) ToMinutes() (int, error) {
	d, err := ParseDuration(string(r))
	if err != nil {
		return 0, err
	}
	minutes := d.Minutes
	minutes += d.Hours * 60
	minutes += d.Days * 24 * 60
	minutes += d.Weeks * 7 * 24 * 60
	minutes += d.Months * 30 * 24 * 60
	minutes += d.Years * 365 * 24 * 60
	return minutes, nil
}

type Duration struct {
	Minutes int
	Hours   int
	Days    int
	Weeks   int
	Months  int
	Years   int
}

var errInvalidDuration = errors.New("invalid duration provided")

// ParseDuration parses a duration from a string. The format is `6y5m234d37h`
func ParseDuration(s string) (Duration, error) {
	var (
		d   Duration
		num int
		err error
	)

	s = strings.TrimSpace(s)

	for s != "" {
		num, s, err = nextNumber(s)
		if err != nil {
			return Duration{}, err
		}

		if len(s) == 0 {
			return Duration{}, errInvalidDuration
		}

		if len(s) > 1 && s[0] == 'm' && s[1] == 'o' {
			d.Months = num
			s = s[2:]
			continue
		}

		switch s[0] {
		case 'y':
			d.Years = num
		case 'w':
			d.Weeks = num
		case 'd':
			d.Days = num
		case 'h':
			d.Hours = num
		case 'm':
			d.Minutes = num
		default:
			return Duration{}, errInvalidDuration
		}

		s = s[1:]
	}

	return d, nil
}

func nextNumber(input string) (num int, rest string, err error) {
	if len(input) == 0 {
		return 0, "", nil
	}

	var (
		n        string
		negative bool
	)

	if input[0] == '-' {
		negative = true
		input = input[1:]
	}

	for i, s := range input {
		if !unicode.IsNumber(s) {
			rest = input[i:]
			break
		}

		n += string(s)
	}

	if len(n) == 0 {
		return 0, input, errInvalidDuration
	}

	num, err = strconv.Atoi(n)
	if err != nil {
		return 0, input, err
	}

	if negative {
		num = -num
	}

	return num, rest, nil
}
