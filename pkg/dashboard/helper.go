/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dashboard

import (
	"regexp"
	"strings"
	"unicode"
)

func excludeQuotedSubstrings(input string) string {
	// Define the regular expression pattern to match string inside double quotation
	re := regexp.MustCompile(`"[^"]*"`)

	// Replace all quoted substring with an empty string
	result := re.ReplaceAllString(input, "")

	return result
}

func excludeNonAlphanumericUnderscore(input string) string {
	// Define the regular expression pattern to match non-alphanumeric characters except underscore
	pattern := `[^a-zA-Z0-9_]`
	re := regexp.MustCompile(pattern)

	// Replace non-alphanumeric or underscore characters with an empty string
	result := re.ReplaceAllString(input, "")

	return result
}

// Labels may contain ASCII letters, numbers, as well as underscores. They must match the regex [a-zA-Z_][a-zA-Z0-9_]*
// So we need to split the selector string by comma. then extract label name with the help of the regex format
// Ref: https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels
func getLabelNames(labelSelector string) []string {
	var labelNames []string
	unQuoted := excludeQuotedSubstrings(labelSelector)
	commaSeparated := strings.Split(unQuoted, ",")
	for _, s := range commaSeparated {
		labelName := excludeNonAlphanumericUnderscore(s)
		labelNames = append(labelNames, labelName)
	}
	return labelNames
}

// Finding valid bracket sequence from startPosition
func substringInsideLabelSelector(query string, startPosition int) (string, int) {
	balance := 0
	closingPosition := startPosition
	for i := startPosition; i < len(query); i++ {
		if query[i] == '{' {
			balance++
		}
		if query[i] == '}' {
			balance--
		}
		if balance == 0 {
			closingPosition = i
			break
		}
	}
	return query[startPosition+1 : closingPosition], closingPosition
}

// Metric names may contain ASCII letters, digits, underscores, and colons. It must match the regex [a-zA-Z_:][a-zA-Z0-9_:]*
// So we can use this if the character is in a metric name
// Ref: https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels
func matchMetricRegex(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' || char == ':'
}

func uniqueAppend(slice []string, valueToAdd string) []string {
	for _, existingValue := range slice {
		if existingValue == valueToAdd {
			return slice
		}
	}
	return append(slice, valueToAdd)
}
