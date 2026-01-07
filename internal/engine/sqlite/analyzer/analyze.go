package analyzer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"

	core "github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/shfmt"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

type Analyzer struct {
	db       config.Database
	conn     *sqlite3.Conn
	dbg      opts.Debug
	replacer *shfmt.Replacer
	mu       sync.Mutex
}

func New(db config.Database) *Analyzer {
	return &Analyzer{
		db:       db,
		dbg:      opts.DebugFromEnv(),
		replacer: shfmt.NewReplacer(nil),
	}
}

func (a *Analyzer) Analyze(ctx context.Context, n ast.Node, query string, migrations []string, ps *named.ParamSet) (*core.Analysis, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		var uri string
		applyMigrations := a.db.Managed
		if a.db.Managed {
			// For managed databases, create an in-memory database
			uri = ":memory:"
		} else if a.dbg.OnlyManagedDatabases {
			return nil, fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
		} else {
			uri = a.replacer.Replace(a.db.URI)
			// For in-memory databases, we need to apply migrations since the database starts empty
			if isInMemoryDatabase(uri) {
				applyMigrations = true
			}
		}

		conn, err := sqlite3.Open(uri)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite database: %w", err)
		}
		a.conn = conn

		// Apply migrations for managed or in-memory databases
		if applyMigrations {
			for _, m := range migrations {
				if len(strings.TrimSpace(m)) == 0 {
					continue
				}
				if err := a.conn.Exec(m); err != nil {
					a.conn.Close()
					a.conn = nil
					return nil, fmt.Errorf("migration failed: %s: %w", m, err)
				}
			}
		}
	}

	// Prepare the statement to get column and parameter information
	stmt, _, err := a.conn.Prepare(query)
	if err != nil {
		return nil, a.extractSqlErr(n, err)
	}
	defer stmt.Close()

	var result core.Analysis

	// Get column information
	colCount := stmt.ColumnCount()
	for i := 0; i < colCount; i++ {
		name := stmt.ColumnName(i)
		declType := stmt.ColumnDeclType(i)
		tableName := stmt.ColumnTableName(i)
		originName := stmt.ColumnOriginName(i)
		dbName := stmt.ColumnDatabaseName(i)

		// Normalize the data type
		dataType := normalizeType(declType)

		// Determine if column is NOT NULL
		// SQLite doesn't provide this info directly from prepared statements,
		// so we default to nullable (false)
		notNull := false

		col := &core.Column{
			Name:         name,
			OriginalName: originName,
			DataType:     dataType,
			NotNull:      notNull,
		}

		if tableName != "" {
			col.Table = &core.Identifier{
				Schema: dbName,
				Name:   tableName,
			}
		}

		result.Columns = append(result.Columns, col)
	}

	// Get parameter information
	bindCount := stmt.BindCount()
	for i := 1; i <= bindCount; i++ {
		paramName := stmt.BindName(i)

		// SQLite doesn't provide parameter types from prepared statements
		// We use "any" as the default type
		name := ""
		if paramName != "" {
			// Remove the prefix (?, :, @, $) from parameter names
			name = strings.TrimLeft(paramName, "?:@$")
		}
		if ps != nil {
			if n, ok := ps.NameFor(i); ok {
				name = n
			}
		}

		result.Params = append(result.Params, &core.Parameter{
			Number: int32(i),
			Column: &core.Column{
				Name:     name,
				DataType: "any",
				NotNull:  false,
			},
		})
	}

	return &result, nil
}

func (a *Analyzer) extractSqlErr(n ast.Node, err error) error {
	if err == nil {
		return nil
	}
	// Try to extract SQLite error details
	var sqliteErr *sqlite3.Error
	if e, ok := err.(*sqlite3.Error); ok {
		sqliteErr = e
	}
	if sqliteErr != nil {
		return &sqlerr.Error{
			Code:     fmt.Sprintf("%d", sqliteErr.Code()),
			Message:  sqliteErr.Error(),
			Location: n.Pos(),
		}
	}
	return &sqlerr.Error{
		Message:  err.Error(),
		Location: n.Pos(),
	}
}

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
// This is useful for database-only mode where we need to connect before analyzing queries.
func (a *Analyzer) EnsureConn(ctx context.Context, migrations []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn != nil {
		return nil
	}

	var uri string
	applyMigrations := a.db.Managed
	if a.db.Managed {
		// For managed databases, create an in-memory database
		uri = ":memory:"
	} else if a.dbg.OnlyManagedDatabases {
		return fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
	} else {
		uri = a.replacer.Replace(a.db.URI)
		// For in-memory databases, we need to apply migrations since the database starts empty
		if isInMemoryDatabase(uri) {
			applyMigrations = true
		}
	}

	conn, err := sqlite3.Open(uri)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}
	a.conn = conn

	// Apply migrations for managed or in-memory databases
	if applyMigrations {
		for _, m := range migrations {
			if len(strings.TrimSpace(m)) == 0 {
				continue
			}
			if err := a.conn.Exec(m); err != nil {
				a.conn.Close()
				a.conn = nil
				return fmt.Errorf("migration failed: %s: %w", m, err)
			}
		}
	}

	return nil
}

// GetColumnNames implements the expander.ColumnGetter interface.
// It prepares a query and returns the column names from the result set description.
func (a *Analyzer) GetColumnNames(ctx context.Context, query string) ([]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	stmt, _, err := a.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	colCount := stmt.ColumnCount()
	columns := make([]string, colCount)
	for i := 0; i < colCount; i++ {
		columns[i] = stmt.ColumnName(i)
	}

	return columns, nil
}

// IntrospectSchema queries the database to build a catalog containing
// tables and columns for the database.
func (a *Analyzer) IntrospectSchema(ctx context.Context, schemas []string) (*catalog.Catalog, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	// Build catalog
	cat := &catalog.Catalog{
		DefaultSchema: "main",
	}

	// Create default schema
	mainSchema := &catalog.Schema{Name: "main"}
	cat.Schemas = append(cat.Schemas, mainSchema)

	// Query tables from sqlite_master
	stmt, _, err := a.conn.Prepare("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return nil, fmt.Errorf("introspect tables: %w", err)
	}

	tableNames := []string{}
	for stmt.Step() {
		tableName := stmt.ColumnText(0)
		tableNames = append(tableNames, tableName)
	}
	stmt.Close()

	// For each table, get column information using PRAGMA table_info
	for _, tableName := range tableNames {
		tbl := &catalog.Table{
			Rel: &ast.TableName{
				Name: tableName,
			},
		}

		pragmaStmt, _, err := a.conn.Prepare(fmt.Sprintf("PRAGMA table_info('%s')", tableName))
		if err != nil {
			return nil, fmt.Errorf("pragma table_info for %s: %w", tableName, err)
		}

		for pragmaStmt.Step() {
			// PRAGMA table_info returns: cid, name, type, notnull, dflt_value, pk
			colName := pragmaStmt.ColumnText(1)
			colType := pragmaStmt.ColumnText(2)
			notNull := pragmaStmt.ColumnInt(3) != 0

			tbl.Columns = append(tbl.Columns, &catalog.Column{
				Name:      colName,
				Type:      ast.TypeName{Name: normalizeType(colType)},
				IsNotNull: notNull,
			})
		}
		pragmaStmt.Close()

		mainSchema.Tables = append(mainSchema.Tables, tbl)
	}

	return cat, nil
}

// isInMemoryDatabase checks if a SQLite URI refers to an in-memory database
func isInMemoryDatabase(uri string) bool {
	if uri == ":memory:" || uri == "" {
		return true
	}
	// Check for file URI with mode=memory parameter
	// e.g., "file:test?mode=memory&cache=shared"
	if strings.Contains(uri, "mode=memory") {
		return true
	}
	return false
}

// normalizeType converts SQLite type declarations to standard type names
func normalizeType(declType string) string {
	if declType == "" {
		return "any"
	}

	// Convert to lowercase for comparison
	lower := strings.ToLower(declType)

	// SQLite type affinity rules (https://www.sqlite.org/datatype3.html)
	switch {
	case strings.Contains(lower, "int"):
		return "integer"
	case strings.Contains(lower, "char"),
		strings.Contains(lower, "clob"),
		strings.Contains(lower, "text"):
		return "text"
	case strings.Contains(lower, "blob"):
		return "blob"
	case strings.Contains(lower, "real"),
		strings.Contains(lower, "floa"),
		strings.Contains(lower, "doub"):
		return "real"
	case strings.Contains(lower, "bool"):
		return "boolean"
	case strings.Contains(lower, "date"),
		strings.Contains(lower, "time"):
		return "datetime"
	default:
		// Return as-is for numeric or other types
		return lower
	}
}
