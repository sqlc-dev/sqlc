package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"sync"

	_ "github.com/ClickHouse/clickhouse-go/v2" // ClickHouse driver
	dcast "github.com/sqlc-dev/doubleclick/ast"
	"github.com/sqlc-dev/doubleclick/parser"

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

	// Check if this is a SELECT query that returns columns
	// INSERT, UPDATE, DELETE don't return columns
	isSelectQuery := isSelectStatement(query)

	if isSelectQuery {
		// For ClickHouse, we use DESCRIBE or LIMIT 0 to get column information

		// Replace all parameter placeholders with NULL for introspection
		// This handles both ? placeholders and {name:Type} named parameters
		preparedQuery := replaceParamsWithNull(query)

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
				NotNull:  true, // Parameters are typically not nullable
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
			// For CREATE TABLE statements, drop the table first if it exists
			upper := strings.ToUpper(strings.TrimSpace(m))
			if strings.HasPrefix(upper, "CREATE TABLE") {
				// Extract table name and drop it first
				parts := strings.Fields(m)
				if len(parts) >= 3 {
					tableName := parts[2]
					// Remove any trailing characters like "("
					tableName = strings.TrimSuffix(tableName, "(")
					a.conn.ExecContext(ctx, "DROP TABLE IF EXISTS "+tableName)
				}
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

	// Replace ? placeholders with NULL for introspection
	preparedQuery := strings.ReplaceAll(query, "?", "NULL")

	// Use DESCRIBE (query) to get column information
	describeQuery := fmt.Sprintf("DESCRIBE (%s)", preparedQuery)
	rows, err := a.conn.QueryContext(ctx, describeQuery)
	if err != nil {
		// Fallback to LIMIT 0 if DESCRIBE fails
		limitQuery := addLimit0(preparedQuery)
		rows, err = a.conn.QueryContext(ctx, limitQuery)
		if err != nil {
			return nil, err
		}
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

// detectParameters finds parameters in a ClickHouse query using the doubleclick parser.
// ClickHouse supports {name:Type} and ? style parameters.
func detectParameters(query string) []paramInfo {
	var params []paramInfo

	ctx := context.Background()

	// First, try to parse and walk the query AST for named parameters
	stmts, err := parser.Parse(ctx, strings.NewReader(query))
	if err == nil {
		for _, stmt := range stmts {
			walkStatement(stmt, func(expr dcast.Expression) {
				if param, ok := expr.(*dcast.Parameter); ok {
					if param.Name != "" {
						dataType := "any"
						if param.Type != nil {
							dataType = normalizeType(param.Type.Name)
						}
						params = append(params, paramInfo{
							Name: param.Name,
							Type: dataType,
						})
					}
				}
			})
		}
	}

	// If no named parameters found from AST, try to extract VALUES clause for INSERT statements
	// The doubleclick parser doesn't fully parse VALUES, so we parse it as a SELECT
	if len(params) == 0 {
		params = extractValuesParameters(ctx, query)
	}

	// Count ? placeholders and add them after any named parameters
	count := strings.Count(query, "?")
	for i := 0; i < count; i++ {
		params = append(params, paramInfo{
			Name: fmt.Sprintf("p%d", len(params)+1),
			Type: "any",
		})
	}

	return params
}

// extractValuesParameters extracts parameters from INSERT VALUES clause by parsing it as a SELECT.
// This works around the limitation that doubleclick doesn't parse VALUES clause expressions.
func extractValuesParameters(ctx context.Context, query string) []paramInfo {
	var params []paramInfo

	// Find VALUES clause (case insensitive)
	upperQuery := strings.ToUpper(query)
	valuesIdx := strings.Index(upperQuery, "VALUES")
	if valuesIdx == -1 {
		return params
	}

	// Extract everything after VALUES
	valuesClause := query[valuesIdx+6:]

	// Find the parentheses containing the values
	start := strings.Index(valuesClause, "(")
	if start == -1 {
		return params
	}

	// Find matching closing parenthesis
	depth := 0
	end := -1
	for i := start; i < len(valuesClause); i++ {
		switch valuesClause[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				end = i
				break
			}
		}
		if end != -1 {
			break
		}
	}

	if end == -1 {
		return params
	}

	// Extract the values list and convert to SELECT query
	valuesList := valuesClause[start+1 : end]
	selectQuery := "SELECT " + valuesList

	// Parse the synthetic SELECT query
	stmts, err := parser.Parse(ctx, strings.NewReader(selectQuery))
	if err != nil {
		return params
	}

	// Walk the AST to find Parameter nodes
	for _, stmt := range stmts {
		walkStatement(stmt, func(expr dcast.Expression) {
			if param, ok := expr.(*dcast.Parameter); ok {
				if param.Name != "" {
					dataType := "any"
					if param.Type != nil {
						dataType = normalizeType(param.Type.Name)
					}
					params = append(params, paramInfo{
						Name: param.Name,
						Type: dataType,
					})
				}
			}
		})
	}

	return params
}

// walkStatement walks a statement and calls fn for each expression.
func walkStatement(stmt dcast.Statement, fn func(dcast.Expression)) {
	switch s := stmt.(type) {
	case *dcast.SelectQuery:
		walkSelectQuery(s, fn)
	case *dcast.SelectWithUnionQuery:
		for _, sel := range s.Selects {
			walkStatement(sel, fn)
		}
	case *dcast.InsertQuery:
		if s.Select != nil {
			walkStatement(s.Select, fn)
		}
	}
}

// walkSelectQuery walks a SELECT query and calls fn for each expression.
func walkSelectQuery(s *dcast.SelectQuery, fn func(dcast.Expression)) {
	// Walk columns
	for _, col := range s.Columns {
		walkExpression(col, fn)
	}
	// Walk WHERE clause
	if s.Where != nil {
		walkExpression(s.Where, fn)
	}
	// Walk GROUP BY
	for _, g := range s.GroupBy {
		walkExpression(g, fn)
	}
	// Walk HAVING
	if s.Having != nil {
		walkExpression(s.Having, fn)
	}
	// Walk ORDER BY
	for _, o := range s.OrderBy {
		walkExpression(o.Expression, fn)
	}
	// Walk LIMIT
	if s.Limit != nil {
		walkExpression(s.Limit, fn)
	}
	// Walk OFFSET
	if s.Offset != nil {
		walkExpression(s.Offset, fn)
	}
}

// walkExpression walks an expression and calls fn for each sub-expression.
func walkExpression(expr dcast.Expression, fn func(dcast.Expression)) {
	if expr == nil {
		return
	}
	fn(expr)

	switch e := expr.(type) {
	case *dcast.BinaryExpr:
		walkExpression(e.Left, fn)
		walkExpression(e.Right, fn)
	case *dcast.UnaryExpr:
		walkExpression(e.Operand, fn)
	case *dcast.FunctionCall:
		for _, arg := range e.Arguments {
			walkExpression(arg, fn)
		}
	case *dcast.Subquery:
		walkStatement(e.Query, fn)
	case *dcast.CaseExpr:
		if e.Operand != nil {
			walkExpression(e.Operand, fn)
		}
		for _, when := range e.Whens {
			walkExpression(when.Condition, fn)
			walkExpression(when.Result, fn)
		}
		if e.Else != nil {
			walkExpression(e.Else, fn)
		}
	case *dcast.InExpr:
		walkExpression(e.Expr, fn)
		for _, v := range e.List {
			walkExpression(v, fn)
		}
		if e.Query != nil {
			walkStatement(e.Query, fn)
		}
	case *dcast.BetweenExpr:
		walkExpression(e.Expr, fn)
		walkExpression(e.Low, fn)
		walkExpression(e.High, fn)
	}
}

// namedParamRegex matches ClickHouse named parameters like {name:Type}
var namedParamRegex = regexp.MustCompile(`\{(\w+):(\w+)\}`)

// replaceParamsWithNull replaces all parameter placeholders with NULL for query introspection.
// It handles both ? placeholders and {name:Type} named parameters.
func replaceParamsWithNull(query string) string {
	// Replace {name:Type} named parameters with NULL
	result := namedParamRegex.ReplaceAllString(query, "NULL")
	// Also replace ? placeholders with NULL
	result = strings.ReplaceAll(result, "?", "NULL")
	return result
}

// addLimit0 adds LIMIT 0 to a query for schema introspection.
func addLimit0(query string) string {
	// Simple approach: append LIMIT 0 if not already present
	lower := strings.ToLower(query)
	if strings.Contains(lower, "limit") {
		return query
	}

	// Remove trailing semicolon and whitespace
	trimmed := strings.TrimRight(query, " \t\n\r;")

	return trimmed + " LIMIT 0"
}

// isSelectStatement checks if a query is a SELECT statement that returns columns.
// It skips past comment lines to find the actual statement.
func isSelectStatement(query string) bool {
	lines := strings.Split(query, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines
		if trimmed == "" {
			continue
		}
		// Skip comment lines
		if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		// Check if this is a SELECT or WITH statement
		lower := strings.ToLower(trimmed)
		return strings.HasPrefix(lower, "select") || strings.HasPrefix(lower, "with")
	}
	return false
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
