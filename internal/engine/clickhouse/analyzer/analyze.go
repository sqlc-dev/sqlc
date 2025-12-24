package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/ClickHouse/clickhouse-go/v2" // ClickHouse driver

	core "github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/shfmt"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// Analyzer implements the analyzer.Analyzer interface for ClickHouse.
type Analyzer struct {
	db       config.Database
	conn     *sql.DB
	dbg      opts.Debug
	replacer *shfmt.Replacer
	mu       sync.Mutex
}

// New creates a new ClickHouse analyzer.
func New(db config.Database) *Analyzer {
	return &Analyzer{
		db:       db,
		dbg:      opts.DebugFromEnv(),
		replacer: shfmt.NewReplacer(nil),
	}
}

// Analyze analyzes a query and returns column and parameter information.
func (a *Analyzer) Analyze(ctx context.Context, n ast.Node, query string, migrations []string, ps *named.ParamSet) (*core.Analysis, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		if err := a.connect(ctx, migrations); err != nil {
			return nil, err
		}
	}

	var result core.Analysis

	// For ClickHouse, we use EXPLAIN to get column information
	// First, try to prepare the query to get parameter count
	// ClickHouse uses {name:type} or $1 style placeholders

	// Replace named parameters with positional ones for EXPLAIN
	preparedQuery := query

	// Use DESCRIBE (query) to get column information
	describeQuery := fmt.Sprintf("DESCRIBE (%s)", preparedQuery)
	rows, err := a.conn.QueryContext(ctx, describeQuery)
	if err != nil {
		// If DESCRIBE fails, try executing with LIMIT 0
		limitQuery := addLimit0(preparedQuery)
		rows, err = a.conn.QueryContext(ctx, limitQuery)
		if err != nil {
			return nil, a.extractSqlErr(n, err)
		}
	}
	defer rows.Close()

	// Get column information from result set
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, a.extractSqlErr(n, err)
	}

	for i, colType := range colTypes {
		name := colType.Name()
		dataType := colType.DatabaseTypeName()
		nullable, _ := colType.Nullable()

		col := &core.Column{
			Name:         name,
			OriginalName: name,
			DataType:     normalizeType(dataType),
			NotNull:      !nullable,
		}

		// Try to detect if this is an aggregate function result
		// (ClickHouse doesn't always provide this info)
		_ = i

		result.Columns = append(result.Columns, col)
	}

	// Detect parameters in the query
	// ClickHouse uses {name:Type} syntax or $1, $2, etc.
	params := detectParameters(query)
	for i, param := range params {
		result.Params = append(result.Params, &core.Parameter{
			Number: int32(i + 1),
			Column: &core.Column{
				Name:     param.Name,
				DataType: param.Type,
				NotNull:  false,
			},
		})
	}

	// Override with named params if provided
	if ps != nil {
		for i := range result.Params {
			if name, ok := ps.NameFor(i + 1); ok {
				result.Params[i].Column.Name = name
			}
		}
	}

	return &result, nil
}

func (a *Analyzer) connect(ctx context.Context, migrations []string) error {
	if a.dbg.OnlyManagedDatabases && !a.db.Managed {
		return fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
	}

	uri := a.replacer.Replace(a.db.URI)
	if uri == "" {
		return fmt.Errorf("clickhouse: database URI is required")
	}

	conn, err := sql.Open("clickhouse", uri)
	if err != nil {
		return fmt.Errorf("failed to connect to clickhouse: %w", err)
	}

	// Verify connection
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return fmt.Errorf("failed to ping clickhouse: %w", err)
	}

	a.conn = conn

	// Apply migrations for managed databases
	if a.db.Managed {
		for _, m := range migrations {
			if len(strings.TrimSpace(m)) == 0 {
				continue
			}
			if _, err := a.conn.ExecContext(ctx, m); err != nil {
				a.conn.Close()
				a.conn = nil
				return fmt.Errorf("migration failed: %s: %w", m, err)
			}
		}
	}

	return nil
}

func (a *Analyzer) extractSqlErr(n ast.Node, err error) error {
	if err == nil {
		return nil
	}
	return &sqlerr.Error{
		Message:  err.Error(),
		Location: n.Pos(),
	}
}

// Close closes the database connection.
func (a *Analyzer) Close(_ context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.conn != nil {
		err := a.conn.Close()
		a.conn = nil
		return err
	}
	return nil
}

// EnsureConn initializes the database connection if not already done.
func (a *Analyzer) EnsureConn(ctx context.Context, migrations []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn != nil {
		return nil
	}

	return a.connect(ctx, migrations)
}

// GetColumnNames returns the column names for a query.
func (a *Analyzer) GetColumnNames(ctx context.Context, query string) ([]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	// Add LIMIT 0 to avoid fetching data
	limitQuery := addLimit0(query)

	rows, err := a.conn.QueryContext(ctx, limitQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	return columns, nil
}

// IntrospectSchema queries the database to build a catalog containing tables and columns.
func (a *Analyzer) IntrospectSchema(ctx context.Context, schemas []string) (*catalog.Catalog, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	cat := &catalog.Catalog{
		DefaultSchema: "default",
	}

	// Get current database
	var currentDB string
	if err := a.conn.QueryRowContext(ctx, "SELECT currentDatabase()").Scan(&currentDB); err != nil {
		currentDB = "default"
	}

	// Create default schema
	mainSchema := &catalog.Schema{Name: currentDB}
	cat.Schemas = append(cat.Schemas, mainSchema)

	// Query tables from system.tables
	tableQuery := `SELECT name FROM system.tables WHERE database = currentDatabase() AND engine != 'View'`
	rows, err := a.conn.QueryContext(ctx, tableQuery)
	if err != nil {
		return nil, fmt.Errorf("introspect tables: %w", err)
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, name)
	}
	rows.Close()

	// For each table, get column information from system.columns
	for _, tableName := range tableNames {
		tbl := &catalog.Table{
			Rel: &ast.TableName{
				Name: tableName,
			},
		}

		colQuery := `SELECT name, type, default_kind FROM system.columns WHERE database = currentDatabase() AND table = ?`
		colRows, err := a.conn.QueryContext(ctx, colQuery, tableName)
		if err != nil {
			return nil, fmt.Errorf("introspect columns for %s: %w", tableName, err)
		}

		for colRows.Next() {
			var name, colType, defaultKind string
			if err := colRows.Scan(&name, &colType, &defaultKind); err != nil {
				colRows.Close()
				return nil, err
			}

			// Determine if column is NOT NULL
			notNull := !isNullable(colType)

			tbl.Columns = append(tbl.Columns, &catalog.Column{
				Name:      name,
				Type:      ast.TypeName{Name: normalizeType(colType)},
				IsNotNull: notNull,
			})
		}
		colRows.Close()

		mainSchema.Tables = append(mainSchema.Tables, tbl)
	}

	return cat, nil
}

// paramInfo holds information about a detected parameter.
type paramInfo struct {
	Name string
	Type string
}

// detectParameters finds parameters in a ClickHouse query.
// ClickHouse supports {name:Type} and $1, $2 style parameters.
func detectParameters(query string) []paramInfo {
	var params []paramInfo

	// Find {name:Type} style parameters
	i := 0
	for i < len(query) {
		if query[i] == '{' {
			j := i + 1
			for j < len(query) && query[j] != '}' {
				j++
			}
			if j < len(query) {
				paramStr := query[i+1 : j]
				parts := strings.SplitN(paramStr, ":", 2)
				if len(parts) == 2 {
					params = append(params, paramInfo{
						Name: parts[0],
						Type: normalizeType(parts[1]),
					})
				} else if len(parts) == 1 {
					params = append(params, paramInfo{
						Name: parts[0],
						Type: "any",
					})
				}
			}
			i = j + 1
		} else {
			i++
		}
	}

	// Find $1, $2 style parameters (simpler approach)
	for i := 1; i <= 100; i++ {
		placeholder := fmt.Sprintf("$%d", i)
		if strings.Contains(query, placeholder) {
			params = append(params, paramInfo{
				Name: fmt.Sprintf("p%d", i),
				Type: "any",
			})
		} else {
			break
		}
	}

	// Find ? placeholders
	count := strings.Count(query, "?")
	for i := len(params); i < count; i++ {
		params = append(params, paramInfo{
			Name: fmt.Sprintf("p%d", i+1),
			Type: "any",
		})
	}

	return params
}

// addLimit0 adds LIMIT 0 to a query for schema introspection.
func addLimit0(query string) string {
	// Simple approach: append LIMIT 0 if not already present
	lower := strings.ToLower(query)
	if strings.Contains(lower, "limit") {
		return query
	}
	return query + " LIMIT 0"
}

// isNullable checks if a ClickHouse type is nullable.
func isNullable(dataType string) bool {
	return strings.HasPrefix(dataType, "Nullable(") ||
		strings.HasPrefix(strings.ToLower(dataType), "nullable(")
}

// normalizeType normalizes a ClickHouse type name to a standard form.
func normalizeType(dataType string) string {
	if dataType == "" {
		return "any"
	}

	// Strip Nullable wrapper
	dt := dataType
	if strings.HasPrefix(dt, "Nullable(") && strings.HasSuffix(dt, ")") {
		dt = dt[9 : len(dt)-1]
	}

	// Normalize common types
	lower := strings.ToLower(dt)
	switch {
	case strings.HasPrefix(lower, "int8"):
		return "Int8"
	case strings.HasPrefix(lower, "int16"):
		return "Int16"
	case strings.HasPrefix(lower, "int32"):
		return "Int32"
	case strings.HasPrefix(lower, "int64"):
		return "Int64"
	case strings.HasPrefix(lower, "uint8"):
		return "UInt8"
	case strings.HasPrefix(lower, "uint16"):
		return "UInt16"
	case strings.HasPrefix(lower, "uint32"):
		return "UInt32"
	case strings.HasPrefix(lower, "uint64"):
		return "UInt64"
	case strings.HasPrefix(lower, "float32"):
		return "Float32"
	case strings.HasPrefix(lower, "float64"):
		return "Float64"
	case lower == "string" || strings.HasPrefix(lower, "fixedstring"):
		return "String"
	case strings.HasPrefix(lower, "date32"):
		return "Date32"
	case lower == "date":
		return "Date"
	case strings.HasPrefix(lower, "datetime64"):
		return "DateTime64"
	case strings.HasPrefix(lower, "datetime"):
		return "DateTime"
	case lower == "bool" || lower == "boolean":
		return "Bool"
	case lower == "uuid":
		return "UUID"
	case strings.HasPrefix(lower, "decimal"):
		return dt // Keep original precision/scale
	case strings.HasPrefix(lower, "array"):
		return dt // Keep original array type
	case strings.HasPrefix(lower, "map"):
		return dt // Keep original map type
	case strings.HasPrefix(lower, "tuple"):
		return dt // Keep original tuple type
	case strings.HasPrefix(lower, "enum"):
		return dt // Keep original enum type
	case strings.HasPrefix(lower, "lowcardinality"):
		// Extract inner type
		if strings.HasPrefix(dt, "LowCardinality(") && strings.HasSuffix(dt, ")") {
			inner := dt[15 : len(dt)-1]
			return normalizeType(inner)
		}
		return dt
	default:
		return dt
	}
}
