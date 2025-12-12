//go:build ignore

// This is a mock PGLite WASM module for testing.
// Build with: GOOS=wasip1 GOARCH=wasm go build -o mock_pglite.wasm mock_pglite.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Request struct {
	Type       string   `json:"type"`
	Migrations []string `json:"migrations"`
	Query      string   `json:"query"`
}

type Response struct {
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Prepare *PrepareResult `json:"prepare,omitempty"`
}

type ErrorResponse struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Position int    `json:"position"`
}

type PrepareResult struct {
	Columns []ColumnInfo    `json:"columns"`
	Params  []ParameterInfo `json:"params"`
}

type ColumnInfo struct {
	Name        string `json:"name"`
	DataType    string `json:"data_type"`
	DataTypeOID uint32 `json:"data_type_oid"`
	NotNull     bool   `json:"not_null"`
	IsArray     bool   `json:"is_array"`
	ArrayDims   int    `json:"array_dims"`
	TableOID    uint32 `json:"table_oid,omitempty"`
	TableName   string `json:"table_name,omitempty"`
	TableSchema string `json:"table_schema,omitempty"`
}

type ParameterInfo struct {
	Number      int    `json:"number"`
	DataType    string `json:"data_type"`
	DataTypeOID uint32 `json:"data_type_oid"`
	IsArray     bool   `json:"is_array"`
	ArrayDims   int    `json:"array_dims"`
}

// Simple schema tracking
type Column struct {
	Name     string
	Type     string
	NotNull  bool
	IsArray  bool
}

type Table struct {
	Schema  string
	Name    string
	Columns []Column
}

var tables = make(map[string]*Table)

func main() {
	// Read all input from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		writeError("READ", fmt.Sprintf("failed to read stdin: %v", err), 0)
		return
	}

	var req Request
	if err := json.Unmarshal(input, &req); err != nil {
		writeError("PARSE", fmt.Sprintf("failed to parse request: %v (input length: %d)", err, len(input)), 0)
		return
	}

	switch req.Type {
	case "init":
		handleInit(req)
	case "prepare":
		handlePrepare(req)
	default:
		writeError("UNKNOWN", fmt.Sprintf("unknown request type: %s", req.Type), 0)
	}
}

func handleInit(req Request) {
	// Parse migrations to build schema
	for _, migration := range req.Migrations {
		parseMigration(migration)
	}

	resp := Response{Success: true}
	writeResponse(resp)
}

func handlePrepare(req Request) {
	// First apply any migrations
	for _, migration := range req.Migrations {
		parseMigration(migration)
	}

	// Parse the query and infer types
	result, err := analyzeQuery(req.Query)
	if err != nil {
		writeError("42601", err.Error(), 0)
		return
	}

	resp := Response{
		Success: true,
		Prepare: result,
	}
	writeResponse(resp)
}

func parseMigration(sql string) {
	// Very simple CREATE TABLE parser
	sql = strings.ToUpper(sql)

	createTableRe := regexp.MustCompile(`CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(?:(\w+)\.)?(\w+)\s*\(([^)]+)\)`)
	matches := createTableRe.FindStringSubmatch(sql)
	if matches == nil {
		return
	}

	schema := "public"
	if matches[1] != "" {
		schema = strings.ToLower(matches[1])
	}
	tableName := strings.ToLower(matches[2])
	columnsDef := matches[3]

	table := &Table{
		Schema: schema,
		Name:   tableName,
	}

	// Parse columns
	colDefs := strings.Split(columnsDef, ",")
	for _, colDef := range colDefs {
		colDef = strings.TrimSpace(colDef)
		if colDef == "" {
			continue
		}
		// Skip constraints
		if strings.HasPrefix(colDef, "PRIMARY") || strings.HasPrefix(colDef, "FOREIGN") ||
			strings.HasPrefix(colDef, "UNIQUE") || strings.HasPrefix(colDef, "CHECK") ||
			strings.HasPrefix(colDef, "CONSTRAINT") {
			continue
		}

		parts := strings.Fields(colDef)
		if len(parts) < 2 {
			continue
		}

		col := Column{
			Name:    strings.ToLower(parts[0]),
			Type:    strings.ToLower(parts[1]),
			NotNull: strings.Contains(colDef, "NOT NULL"),
			IsArray: strings.Contains(parts[1], "[]"),
		}
		table.Columns = append(table.Columns, col)
	}

	tables[schema+"."+tableName] = table
}

func analyzeQuery(query string) (*PrepareResult, error) {
	query = strings.ToUpper(strings.TrimSpace(query))
	result := &PrepareResult{}

	// Count parameters
	paramCount := strings.Count(query, "$")
	for i := 1; i <= paramCount; i++ {
		result.Params = append(result.Params, ParameterInfo{
			Number:      i,
			DataType:    "text", // Default to text
			DataTypeOID: 25,
		})
	}

	// Very simple SELECT parser
	if strings.HasPrefix(query, "SELECT") {
		// Find FROM clause
		fromIdx := strings.Index(query, "FROM")
		if fromIdx == -1 {
			// SELECT without FROM (e.g., SELECT 1)
			selectPart := query[6:]
			if whereIdx := strings.Index(selectPart, "WHERE"); whereIdx != -1 {
				selectPart = selectPart[:whereIdx]
			}

			cols := strings.Split(selectPart, ",")
			for _, col := range cols {
				col = strings.TrimSpace(col)
				result.Columns = append(result.Columns, ColumnInfo{
					Name:        strings.ToLower(col),
					DataType:    "integer",
					DataTypeOID: 23,
				})
			}
			return result, nil
		}

		selectPart := strings.TrimSpace(query[6:fromIdx])
		fromPart := query[fromIdx+4:]

		// Get table name
		tableName := ""
		parts := strings.Fields(fromPart)
		if len(parts) > 0 {
			tableName = strings.ToLower(parts[0])
		}

		// Look up table
		table := findTable(tableName)

		// Handle SELECT *
		if strings.TrimSpace(selectPart) == "*" {
			if table != nil {
				for _, col := range table.Columns {
					result.Columns = append(result.Columns, ColumnInfo{
						Name:        col.Name,
						DataType:    mapType(col.Type),
						DataTypeOID: mapTypeOID(col.Type),
						NotNull:     col.NotNull,
						IsArray:     col.IsArray,
						TableName:   table.Name,
						TableSchema: table.Schema,
						TableOID:    16384, // Fake OID
					})
				}
			}
			return result, nil
		}

		// Parse individual columns
		cols := strings.Split(selectPart, ",")
		for _, col := range cols {
			col = strings.TrimSpace(col)
			// Handle aliases (col AS alias)
			alias := col
			if asIdx := strings.Index(col, " AS "); asIdx != -1 {
				alias = strings.TrimSpace(col[asIdx+4:])
				col = strings.TrimSpace(col[:asIdx])
			}

			colInfo := ColumnInfo{
				Name:        strings.ToLower(alias),
				DataType:    "text",
				DataTypeOID: 25,
			}

			// Try to find column in table
			if table != nil {
				for _, tc := range table.Columns {
					if strings.ToUpper(tc.Name) == col || strings.ToUpper(table.Name+"."+tc.Name) == col {
						colInfo.DataType = mapType(tc.Type)
						colInfo.DataTypeOID = mapTypeOID(tc.Type)
						colInfo.NotNull = tc.NotNull
						colInfo.IsArray = tc.IsArray
						colInfo.TableName = table.Name
						colInfo.TableSchema = table.Schema
						colInfo.TableOID = 16384
						break
					}
				}
			}

			result.Columns = append(result.Columns, colInfo)
		}
	}

	return result, nil
}

func findTable(name string) *Table {
	// Try with public schema
	if t, ok := tables["public."+name]; ok {
		return t
	}
	// Try as-is
	if t, ok := tables[name]; ok {
		return t
	}
	// Search all tables
	for _, t := range tables {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func mapType(t string) string {
	t = strings.ToLower(t)
	t = strings.TrimSuffix(t, "[]")
	switch t {
	case "int", "integer", "int4":
		return "integer"
	case "bigint", "int8":
		return "bigint"
	case "smallint", "int2":
		return "smallint"
	case "text", "varchar", "character varying":
		return "text"
	case "boolean", "bool":
		return "boolean"
	case "timestamp", "timestamptz", "timestamp with time zone", "timestamp without time zone":
		return "pg_catalog.timestamp"
	case "uuid":
		return "uuid"
	case "jsonb":
		return "jsonb"
	case "json":
		return "json"
	default:
		return t
	}
}

func mapTypeOID(t string) uint32 {
	t = strings.ToLower(t)
	t = strings.TrimSuffix(t, "[]")
	switch t {
	case "int", "integer", "int4":
		return 23
	case "bigint", "int8":
		return 20
	case "smallint", "int2":
		return 21
	case "text", "varchar", "character varying":
		return 25
	case "boolean", "bool":
		return 16
	case "timestamp", "timestamptz":
		return 1114
	case "uuid":
		return 2950
	case "jsonb":
		return 3802
	case "json":
		return 114
	default:
		return 25 // default to text
	}
}

func writeError(code, message string, position int) {
	resp := Response{
		Success: false,
		Error: &ErrorResponse{
			Code:     code,
			Message:  message,
			Position: position,
		},
	}
	writeResponse(resp)
}

func writeResponse(resp Response) {
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}
