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

// MySQLQuery specifies query for MySQL database
type MySQLQuery struct {
	// Database refers to the database name being checked for existence
	Database string `json:"database,omitempty"`

	// Table refers to the table name being checked for existence in specified Database
	// +optional
	Table string `json:"table,omitempty"`

	// RowCount represents the number of row to be checked in the specified Table
	// +optional
	RowCount *MatchExpression `json:"rowCount,omitempty"`
}

// MariaDBQuery specifies query for MariaDB database
type MariaDBQuery struct {
	// Database refers to the database name being checked for existence
	Database string `json:"database,omitempty"`

	// Table refers to the table name being checked for existence in specified Database
	// +optional
	Table string `json:"table,omitempty"`

	// RowCount represents the number of row to be checked in the specified Table
	// +optional
	RowCount *MatchExpression `json:"rowCount,omitempty"`
}

// PostgresQuery specifies query for Postgres database
type PostgresQuery struct {
	// Database refers to the database name being checked for existence
	Database string `json:"database,omitempty"`

	// Schema refers to the schema name being checked for existence in specified Database
	// +optional
	Schema string `json:"schema,omitempty"`

	// Table refers to the table name being checked for existence in specified Database
	// +optional
	Table string `json:"table,omitempty"`

	// RowCount represents the number of row to be checked in the specified Table
	// +optional
	RowCount *MatchExpression `json:"rowCount,omitempty"`
}

// MongoDBQuery specifies query for MongoDB database
type MongoDBQuery struct {
	// Database refers to the database name being checked for existence
	Database string `json:"database,omitempty"`

	// Collection refers to the collection name being checked for existence in specified Database
	// +optional
	Collection string `json:"collection,omitempty"`

	// RowCount represents the number of document to be checked in the specified Collection
	// +optional
	DocumentCount *MatchExpression `json:"documentCount,omitempty"`
}

// ElasticsearchQuery specifies query for Elasticsearch database
type ElasticsearchQuery struct {
	// Index refers to the index name being checked for existence
	Index string `json:"index,omitempty"`
}

// RedisQuery specifies query for Redis database
type RedisQuery struct {
	// Index refers to the database index being checked for existence
	Index int `json:"index,omitempty"`

	// DbSize specifies the number of keys in the specified Database
	// +optional
	DbSize *MatchExpression `json:"dbSize,omitempty"`
}

// SinglestoreQuery specifies query for Singlestore database
type SinglestoreQuery struct {
	// Database refers to the database name being checked for existence
	Database string `json:"database,omitempty"`

	// Table refers to the table name being checked for existence in specified Database
	// +optional
	Table string `json:"table,omitempty"`

	// RowCount represents the number of row to be checked in the specified Table
	// +optional
	RowCount *MatchExpression `json:"rowCount,omitempty"`
}

// MSSQLServerQuery specifies query for MSSQLServer database
type MSSQLServerQuery struct {
	// Database refers to the database name being checked for existence
	Database string `json:"database,omitempty"`

	// Schema refers to the schema name being checked for existence in specified Database
	// +optional
	Schema string `json:"schema,omitempty"`

	// Table refers to the table name being checked for existence in specified Database
	// +optional
	Table string `json:"table,omitempty"`

	// RowCount represents the number of row to be checked in the specified Table
	// +optional
	RowCount *MatchExpression `json:"rowCount,omitempty"`
}

type MatchExpression struct {
	// Operator represents the operation that will be done on the given Value
	Operator Operator `json:"operator,omitempty"`

	// Value represents the numerical value of the desired output
	Value *int64 `json:"value,omitempty"`
}

// Operator represents the operation that will be done
// +kubebuilder:validation:Enum=Equal;NotEqual;LessThan;LessThanOrEqual;GreaterThan;GreaterThanOrEqual
type Operator string

const (
	EqualOperator              Operator = "Equal"
	NotEqualOperator           Operator = "NotEqual"
	LessThanOperator           Operator = "LessThan"
	LessThanOrEqualOperator    Operator = "LessThanOrEqual"
	GreaterThanOperator        Operator = "GreaterThan"
	GreaterThanOrEqualOperator Operator = "GreaterThanOrEqual"
)
