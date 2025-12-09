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
	"log"
	"regexp"
	"strings"
	"unicode"
)

func parseAllExpressions(dashboardData map[string]any) []queryOpts {
	var queries []queryOpts
	if panels, ok := dashboardData["panels"].([]any); ok {
		for _, panel := range panels {
			if targets, ok := panel.(map[string]any)["targets"].([]any); ok {
				title, ok := panel.(map[string]any)["title"].(string)
				if !ok {
					log.Fatal("panel's title found empty")
				}
				for _, target := range targets {
					if expr, ok := target.(map[string]any)["expr"]; ok {
						if expr != "" {
							query := expr.(string)
							queries = append(queries, parseSingleExpression(query, title)...)
						}
					}
				}
			}
		}
	}
	return queries
}

// Steps:
// - if current character is '{'
//   - extract metric name by matching metric regex
//   - get label selector substring inside { }
//   - get label name from this substring by matching label regex
//   - move i to its closing bracket position.
func parseSingleExpression(query, title string) []queryOpts {
	var queries []queryOpts
	for i := 0; i < len(query); i++ {
		if query[i] == '{' {
			j := i
			for j-1 >= 0 && (!matchMetricRegex(rune(query[j-1]))) {
				j--
			}
			metric := query[j:i]
			fullLabelString, closingPosition := getFullLabelString(query, i)
			labelNames := parseLabelNames(fullLabelString)
			queries = append(queries, queryOpts{
				metric:     metric,
				labelNames: labelNames,
				panelTitle: title,
			})
			i = closingPosition
		}
	}
	return queries
}

func matchMetricRegex(char rune) bool { // Must match the regex [a-zA-Z_:][a-zA-Z0-9_:]*
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' || char == ':'
}

// Finding valid bracket sequence from startPosition
func getFullLabelString(query string, startPosition int) (string, int) {
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

// Labels may contain ASCII letters, numbers, as well as underscores. They must match the regex [a-zA-Z_][a-zA-Z0-9_]*
// So we need to split the selector string by comma. then extract label name with the help of the regex format
// Ref: https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels
func parseLabelNames(fullLabelString string) []string {
	// Define the regular expression pattern to match string inside double quotation
	// Replace all quoted substring with an empty string
	excludeQuotedSubstrings := func(input string) string {
		re := regexp.MustCompile(`"[^"]*"`)
		result := re.ReplaceAllString(input, "")
		return result
	}

	// Define the regular expression pattern to match non-alphanumeric characters except underscore
	// Replace non-alphanumeric or underscore characters with an empty string
	excludeNonAlphanumericUnderscore := func(input string) string {
		pattern := `[^a-zA-Z0-9_]`
		re := regexp.MustCompile(pattern)
		result := re.ReplaceAllString(input, "")
		return result
	}

	var labelNames []string
	unQuoted := excludeQuotedSubstrings(fullLabelString)
	commaSeparated := strings.Split(unQuoted, ",")
	for _, s := range commaSeparated {
		labelName := excludeNonAlphanumericUnderscore(s)
		if labelName != "" {
			labelNames = append(labelNames, labelName)
		}
	}
	return labelNames
}
