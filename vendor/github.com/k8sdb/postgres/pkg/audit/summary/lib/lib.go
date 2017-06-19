package lib

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-xorm/xorm"
	pg "github.com/lib/pq"
)

type TableInfo struct {
	TotalRow int64 `json:"total_row"`
	MaxID    int64 `json:"max_id"`
	NextID   int64 `json:"next_id"`
}

type SchemaInfo struct {
	Table map[string]*TableInfo `json:"table"`
}

type DBInfo struct {
	Schema map[string]*SchemaInfo `json:"schema"`
}

func DumpDBInfo(engine *xorm.Engine) (*DBInfo, error) {
	defer engine.Close()
	engine.ShowSQL(true)
	session := engine.NewSession()
	session.Close()
	schemaRowSlice, err := session.Query("select schema_name from information_schema.schemata")
	if err != nil {
		return nil, err
	}

	schemaList := make(map[string]*SchemaInfo, 0)
	for _, row := range schemaRowSlice {
		schemaName := string(row["schema_name"])
		schemaInfo, err := getDataFromSchema(session, schemaName)
		if err != nil {
			return nil, err
		}
		schemaList[schemaName] = schemaInfo
	}

	return &DBInfo{
		Schema: schemaList,
	}, nil
}

func getDataFromSchema(session *xorm.Session, schemaName string) (*SchemaInfo, error) {
	tableRowSlice, err := session.Query("SELECT tablename FROM pg_tables where schemaname=$1", schemaName)
	if err != nil {
		log.Fatalln(err)
	}

	schemaInfo := &SchemaInfo{
		Table: make(map[string]*TableInfo),
	}

	for _, row := range tableRowSlice {
		for _, val := range row {
			tableName := string(val)
			tableInfo, err := getDataFromTable(session, schemaName, tableName)
			if err != nil {
				return nil, err
			}
			schemaInfo.Table[tableName] = tableInfo
		}
	}

	return schemaInfo, nil
}

const (
	errorUndefinedColumn  = "undefined_column"
	errorDatatypeMismatch = "datatype_mismatch"
	invalidData           = -1
)

func getDataFromTable(session *xorm.Session, schemaName, tableName string) (*TableInfo, error) {
	table := fmt.Sprintf(`"%v".%v`, schemaName, tableName)
	dataRows, err := session.Query(fmt.Sprintf(`SELECT count(*) as total_row, coalesce(max(id),0) as max_id FROM %v`, table))

	var totalRow, maxID, nextID int64
	var errorName string

	if driverErr, ok := err.(*pg.Error); ok {
		errorName = driverErr.Code.Name()
		if errorName == errorUndefinedColumn || errorName == errorDatatypeMismatch {
			dataRows, err = session.Query(fmt.Sprintf("SELECT count(*) as total_row FROM %v", table))
			if err != nil {
				return &TableInfo{}, err
			}

			if totalRow, err = strconv.ParseInt(string(dataRows[0]["total_row"]), 10, 64); err != nil {
				return &TableInfo{}, err
			}
			maxID = invalidData
			nextID = invalidData

		} else {
			return &TableInfo{}, err
		}
	} else {
		if len(dataRows) == 0 {
			log.Println("Data missing: ", err)
			totalRow = invalidData
			maxID = invalidData
			nextID = invalidData
		} else {
			if totalRow, err = strconv.ParseInt(string(dataRows[0]["total_row"]), 10, 64); err != nil {
				return &TableInfo{}, err
			}

			if maxID, err = strconv.ParseInt(string(dataRows[0]["max_id"]), 10, 64); err != nil {
				return &TableInfo{}, err
			}

			dataRows, err = session.Query(fmt.Sprintf(`select (last_value+1) as next_id from %v_id_seq`, table))
			if err != nil {
				return &TableInfo{}, err
			}
			if len(dataRows) == 0 {
				nextID = invalidData
			} else {
				if nextID, err = strconv.ParseInt(string(dataRows[0]["next_id"]), 10, 64); err != nil {
					return &TableInfo{}, err
				}
			}
		}
	}

	return &TableInfo{
		TotalRow: totalRow,
		MaxID:    maxID,
		NextID:   nextID,
	}, nil
}

func GetAllDatabase(engine *xorm.Engine) ([]string, error) {
	defer engine.Close()
	engine.ShowSQL(true)
	session := engine.NewSession()
	session.Close()
	rows, err := session.Query("SELECT datname FROM pg_database where datistemplate = false")
	if err != nil {
		return nil, err
	}

	databases := make([]string, 0)

	for _, row := range rows {
		databases = append(databases, string(row["datname"]))
	}
	return databases, nil
}
