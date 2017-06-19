package postgres

import (
	"fmt"
	"github.com/k8sdb/cli/pkg/audit/type"
	"sort"
)

func ToSerializePostgresReport(report *types.Summary) map[string]int64 {
	d := make(map[string]int64)

	for dbName, dbInfo := range report.SummaryReport.Postgres {
		for schemaName, schemaInfo := range dbInfo.Schema {
			for tableName, tableInfo := range schemaInfo.Table {
				totalRowKey := fmt.Sprintf("%v.%v.%v.total_row", dbName, schemaName, tableName)
				maxIDKey := fmt.Sprintf("%v.%v.%v.max_id", dbName, schemaName, tableName)
				nextIDKey := fmt.Sprintf("%v.%v.%v.next_id", dbName, schemaName, tableName)

				d[totalRowKey] = tableInfo.TotalRow
				d[maxIDKey] = tableInfo.MaxID
				d[nextIDKey] = tableInfo.NextID
			}
		}
	}
	return d
}

func DiffSerializePostgresReport(original, duplicate map[string]int64) []string {
	keyList := make(map[string]bool)

	for key, _ := range original {
		keyList[key] = true
	}

	for key, _ := range duplicate {
		keyList[key] = true
	}

	result := make([]string, 0)
	diff := make([]string, 0)

	var keys []string
	for k, _ := range keyList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		_, foundInOriginal := original[key]
		_, foundInDuplicate := duplicate[key]

		if foundInOriginal && !foundInDuplicate {
			diff = append(diff, fmt.Sprintf("---    %v: %v", key, original[key]))
		} else if !foundInOriginal && foundInDuplicate {
			diff = append(diff, fmt.Sprintf("+++    %v: %v", key, duplicate[key]))
		} else if foundInOriginal && foundInDuplicate {
			if original[key] != duplicate[key] {
				diff = append(diff, fmt.Sprintf("+++--- %v: %v --> %v", key, original[key], duplicate[key]))
			}
		}
	}

	if len(diff) == 0 {
		result = append(result, "No change\n")
	} else {
		result = append(result, fmt.Sprintf("Total mismatch: %d\n", len(diff)))
		diff = append(diff, "")
	}

	return append(result, diff...)
}
