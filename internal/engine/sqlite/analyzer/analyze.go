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
		if a.db.Managed {
			// For managed databases, create an in-memory database
			uri = ":memory:"
		} else if a.dbg.OnlyManagedDatabases {
			return nil, fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
		} else {
			uri = a.replacer.Replace(a.db.URI)
		}

		conn, err := sqlite3.Open(uri)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite database: %w", err)
		}
		a.conn = conn

		// Apply migrations for managed databases
		if a.db.Managed {
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
