package compare

import (
	"fmt"
	"sort"

	"github.com/k8sdb/postgres/pkg/audit/summary"
	"github.com/k8sdb/postgres/pkg/audit/type"
)

func CompareSummaryReport(reportData1, reportData2 map[string]*types.DBInfo, dbname string) []string {
	srData1 := toSerializeReport(reportData1, dbname)
	srData2 := toSerializeReport(reportData2, dbname)

	return diffSerializeReport(srData1, srData2)
}

func toSerializeReport(report map[string]*types.DBInfo, dbname string) map[string]int64 {
	d := make(map[string]int64)

	for dbName, dbInfo := range report {
		if dbname == "" || dbname == dbName {
			for schemaName, schemaInfo := range dbInfo.Schema {
				for tableName, tableInfo := range schemaInfo.Table {
					prefix := fmt.Sprintf("%v.%v.%v", dbName, schemaName, tableName)
					totalRowKey := fmt.Sprintf("%v.%v", prefix, summary.TotalRow)
					maxIDKey := fmt.Sprintf("%v.%v", prefix, summary.MaxID)
					nextIDKey := fmt.Sprintf("%v.%v", prefix, summary.NextID)

					d[totalRowKey] = tableInfo.TotalRow
					d[maxIDKey] = tableInfo.MaxID
					d[nextIDKey] = tableInfo.NextID
				}
			}
		}
	}
	return d
}

func diffSerializeReport(srData1, srData2 map[string]int64) []string {
	keyList := make(map[string]bool)

	for key := range srData1 {
		keyList[key] = true
	}

	for key := range srData2 {
		keyList[key] = true
	}

	diff := make([]string, 0)

	var keys []string
	for k := range keyList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		_, foundInOriginal := srData1[key]
		_, foundInDuplicate := srData2[key]

		if foundInOriginal && !foundInDuplicate {
			diff = append(diff, fmt.Sprintf("---    %v: %v", key, srData1[key]))
		} else if !foundInOriginal && foundInDuplicate {
			diff = append(diff, fmt.Sprintf("+++    %v: %v", key, srData2[key]))
		} else if foundInOriginal && foundInDuplicate {
			if srData1[key] != srData2[key] {
				diff = append(diff, fmt.Sprintf("+++--- %v: %v --> %v", key, srData1[key], srData2[key]))
			}
		}
	}

	return diff
}
