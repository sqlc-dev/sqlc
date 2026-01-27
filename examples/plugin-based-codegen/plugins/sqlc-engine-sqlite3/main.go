// sqlc-engine-sqlite3 demonstrates a custom database engine plugin.
//
// This plugin provides SQLite3 SQL parsing for sqlc. It shows how external
// repositories can implement database support without modifying sqlc core.
//
// Build: go build -o sqlc-engine-sqlite3 .
package main

import (
	"regexp"
	"strings"

	"github.com/sqlc-dev/sqlc/pkg/engine"
)

func main() {
	engine.Run(engine.Handler{
		PluginName:    "sqlite3",
		PluginVersion: "1.0.0",
		Parse:         handleParse,
	})
}

func handleParse(req *engine.ParseRequest) (*engine.ParseResponse, error) {
	sql := req.GetSql()
	
	// Parse schema if provided
	var schema *SchemaInfo
	if schemaSQL := req.GetSchemaSql(); schemaSQL != "" {
		schema = parseSchema(schemaSQL)
	}
	// Note: connection_params support can be added here if needed
	
	// Extract parameters from SQL
	parameters := extractParameters(sql)
	
	// Extract columns from SQL (expand wildcards if schema is available)
	columns := extractColumns(sql, schema)
	
	// Process SQL (expand wildcards if needed)
	processedSQL := processSQL(sql, schema)
	
	return &engine.ParseResponse{
		Sql:        processedSQL,
		Parameters: parameters,
		Columns:    columns,
	}, nil
}

// SchemaInfo represents parsed schema information
type SchemaInfo struct {
	Tables map[string]*TableInfo
}

type TableInfo struct {
	Name    string
	Columns []*ColumnInfo
}

type ColumnInfo struct {
	Name     string
	DataType string
	NotNull  bool
}

// parseSchema parses CREATE TABLE statements from schema SQL
func parseSchema(schemaSQL string) *SchemaInfo {
	info := &SchemaInfo{
		Tables: make(map[string]*TableInfo),
	}
	// (?s) makes . match newlines so we capture multiline table body
	createTableRegex := regexp.MustCompile(`(?is)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(?:["]?(\w+)["]?\.)?["]?(\w+)["]?\s*\((.*?)\)\s*(?:;|$)`)
	// Column: optional quote, name, spaces, type (INTEGER/TEXT/REAL/BLOB/...)
	columnDefRegex := regexp.MustCompile(`["]?(\w+)["]?\s+(INTEGER|INT|TEXT|REAL|BLOB|VARCHAR|CHAR|FLOAT|DOUBLE|NUMERIC|DECIMAL|DATE|DATETIME|BOOLEAN)\b`)

	matches := createTableRegex.FindAllStringSubmatch(schemaSQL, -1)
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}
		tableName := match[2]
		if match[1] != "" {
			tableName = match[1] + "." + match[2]
		}
		table := &TableInfo{
			Name:    tableName,
			Columns: []*ColumnInfo{},
		}
		columnDefs := match[3]
		colMatches := columnDefRegex.FindAllStringSubmatch(columnDefs, -1)
		for _, colMatch := range colMatches {
			if len(colMatch) >= 3 {
				// NOT NULL if that phrase appears after this column name in its segment
				segStart := strings.Index(columnDefs, colMatch[0])
				segEnd := len(columnDefs)
				if idx := strings.Index(columnDefs[segStart:], ","); idx >= 0 {
					segEnd = segStart + idx
				}
				segment := columnDefs[segStart:segEnd]
				notNull := strings.Contains(strings.ToUpper(segment), "NOT NULL")
				table.Columns = append(table.Columns, &ColumnInfo{
					Name:     colMatch[1],
					DataType: normalizeDataType(colMatch[2]),
					NotNull:  notNull,
				})
			}
		}
		info.Tables[tableName] = table
		if match[1] != "" {
			info.Tables[match[2]] = table
		}
	}
	return info
}

// normalizeDataType normalizes SQLite data types
func normalizeDataType(dt string) string {
	dt = strings.ToUpper(dt)
	switch {
	case strings.Contains(dt, "INT"):
		return "INTEGER"
	case strings.Contains(dt, "TEXT") || strings.Contains(dt, "CHAR") || strings.Contains(dt, "VARCHAR"):
		return "TEXT"
	case strings.Contains(dt, "REAL") || strings.Contains(dt, "FLOAT") || strings.Contains(dt, "DOUBLE"):
		return "REAL"
	case strings.Contains(dt, "BLOB"):
		return "BLOB"
	default:
		return dt
	}
}

// extractParameters extracts parameters from SQL (?, $1, :name, sqlc.arg(), etc.)
func extractParameters(sql string) []*engine.Parameter {
	var params []*engine.Parameter
	position := 1
	
	// Extract ? positional parameters
	questionMarkRegex := regexp.MustCompile(`\?`)
	matches := questionMarkRegex.FindAllStringIndex(sql, -1)
	for range matches {
		params = append(params, &engine.Parameter{
			Position: int32(position),
			DataType: "TEXT", // Default, could be improved with better parsing
			Nullable: true,
		})
		position++
	}
	
	// Extract sqlc.arg() calls
	sqlcArgRegex := regexp.MustCompile(`sqlc\.arg\(["']?(\w+)["']?\)`)
	argMatches := sqlcArgRegex.FindAllStringSubmatch(sql, -1)
	for _, argMatch := range argMatches {
		if len(argMatch) >= 2 {
			params = append(params, &engine.Parameter{
				Name:     argMatch[1],
				Position: int32(position),
				DataType: "TEXT",
				Nullable: true,
			})
			position++
		}
	}
	
	return params
}

// extractColumns extracts result columns from SELECT statements
func extractColumns(sql string, schema *SchemaInfo) []*engine.Column {
	var columns []*engine.Column
	
	// Check if it's a SELECT statement
	upperSQL := strings.ToUpper(strings.TrimSpace(sql))
	if !strings.HasPrefix(upperSQL, "SELECT") {
		return columns
	}
	
	// Try to extract columns from SELECT clause
	selectRegex := regexp.MustCompile(`(?i)SELECT\s+(.*?)\s+FROM`)
	match := selectRegex.FindStringSubmatch(sql)
	if len(match) < 2 {
		return columns
	}
	
	selectClause := match[1]
	
	// Check for wildcard
	if strings.Contains(selectClause, "*") {
			// Extract table name from FROM clause
		fromRegex := regexp.MustCompile(`(?i)FROM\s+(?:["]?(\w+)["]?\.)?["]?(\w+)["]?`)
		fromMatch := fromRegex.FindStringSubmatch(sql)
		if len(fromMatch) >= 3 && schema != nil {
			tableName := fromMatch[2]
			if table, ok := schema.Tables[tableName]; ok {
				for _, col := range table.Columns {
					columns = append(columns, &engine.Column{
						Name:     col.Name,
						DataType: col.DataType,
						Nullable: !col.NotNull,
						TableName: tableName,
					})
				}
			}
		}
	} else {
		// Extract explicit column names
		columnRegex := regexp.MustCompile(`["]?(\w+)["]?(?:\s+AS\s+["]?(\w+)["]?)?`)
		colMatches := columnRegex.FindAllStringSubmatch(selectClause, -1)
		for _, colMatch := range colMatches {
			colName := colMatch[1]
			if len(colMatch) >= 3 && colMatch[2] != "" {
				colName = colMatch[2] // Use alias if present
			}
			columns = append(columns, &engine.Column{
				Name:     colName,
				DataType: "TEXT", // Default, could be improved
				Nullable: true,
			})
		}
	}
	
	return columns
}

// processSQL processes SQL and expands wildcards if schema is available
func processSQL(sql string, schema *SchemaInfo) string {
	if schema == nil {
		return sql
	}
	
	// Check for SELECT * and expand if we have schema
	if strings.Contains(strings.ToUpper(sql), "SELECT *") {
		fromRegex := regexp.MustCompile(`(?i)(SELECT\s+)\*(\s+FROM\s+(?:["]?(\w+)["]?\.)?["]?(\w+)["]?)`)
		processed := fromRegex.ReplaceAllStringFunc(sql, func(match string) string {
			parts := fromRegex.FindStringSubmatch(match)
			if len(parts) >= 5 {
				tableName := parts[4]
				if table, ok := schema.Tables[tableName]; ok {
					var colNames []string
					for _, col := range table.Columns {
						colNames = append(colNames, `"`+col.Name+`"`)
					}
					return parts[1] + strings.Join(colNames, ", ") + parts[2]
				}
			}
			return match
		})
		return processed
	}
	
	return sql
}

func splitStatements(sql string) []string {
	var statements []string
	var current strings.Builder

	for _, line := range strings.Split(sql, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		// Include sqlc metadata comments (-- name: ...) with the statement
		if strings.HasPrefix(trimmedLine, "--") {
			// Check if it's a sqlc query annotation
			if strings.Contains(trimmedLine, "name:") {
				current.WriteString(trimmedLine)
				current.WriteString("\n")
			}
			// Skip other comments
			continue
		}
		current.WriteString(trimmedLine)
		current.WriteString(" ")
		if strings.HasSuffix(trimmedLine, ";") {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" && stmt != ";" {
				statements = append(statements, stmt)
			}
			current.Reset()
		}
	}
	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}
	return statements
}

func detectStatementType(sql string) string {
	sql = strings.ToUpper(strings.TrimSpace(sql))
	switch {
	case strings.HasPrefix(sql, "SELECT"):
		return "SelectStmt"
	case strings.HasPrefix(sql, "INSERT"):
		return "InsertStmt"
	case strings.HasPrefix(sql, "UPDATE"):
		return "UpdateStmt"
	case strings.HasPrefix(sql, "DELETE"):
		return "DeleteStmt"
	case strings.HasPrefix(sql, "CREATE TABLE"):
		return "CreateTableStmt"
	default:
		return "Unknown"
	}
}
