package types

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
