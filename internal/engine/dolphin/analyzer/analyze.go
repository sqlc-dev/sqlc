package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"io"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	core "github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/shfmt"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

type Analyzer struct {
	db         config.Database
	servers    []config.Server
	conn       *sql.DB
	baseConn   *sql.DB // Connection to base database for creating/dropping temp DBs
	dbName     string  // Name of the created database (for cleanup)
	dbg        opts.Debug
	replacer   *shfmt.Replacer
	mu         sync.Mutex
}

func New(servers []config.Server, db config.Database) *Analyzer {
	return &Analyzer{
		db:       db,
		servers:  servers,
		dbg:      opts.DebugFromEnv(),
		replacer: shfmt.NewReplacer(nil),
	}
}

// dbid creates a unique hash from migration content
func dbid(migrations []string) string {
	h := fnv.New64()
	for _, query := range migrations {
		io.WriteString(h, query)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (a *Analyzer) Analyze(ctx context.Context, n ast.Node, query string, migrations []string, ps *named.ParamSet) (*core.Analysis, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		var uri string
		var applyMigrations bool

		if a.db.Managed {
			// Only require servers for managed databases
			// Non-managed use the database URI directly
			// Find MySQL server from configured servers
			var baseURI string
			for _, server := range a.servers {
				if server.Engine == config.EngineMySQL {
					baseURI = a.replacer.Replace(server.URI)
					break
				}
			}
			if baseURI == "" {
				return nil, fmt.Errorf("no MySQL database server configured")
			}

			// Create a unique database name based on migrations hash
			hash := dbid(migrations)
			a.dbName = fmt.Sprintf("sqlc_managed_%s", hash)

			// Connect to the base database to create our temp database
			baseConn, err := sql.Open("mysql", baseURI)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
			}
			if err := baseConn.PingContext(ctx); err != nil {
				baseConn.Close()
				return nil, fmt.Errorf("failed to ping MySQL server: %w", err)
			}
			a.baseConn = baseConn

			// Check if database already exists
			var dbExists int
			row := baseConn.QueryRowContext(ctx,
				"SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", a.dbName)
			if err := row.Scan(&dbExists); err != nil {
				return nil, fmt.Errorf("failed to check database existence: %w", err)
			}

			if dbExists == 0 {
				// Create the database
				if _, err := baseConn.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s`", a.dbName)); err != nil {
					return nil, fmt.Errorf("failed to create database %s: %w", a.dbName, err)
				}
				applyMigrations = true
			}

			// Build URI for the new database
			// Parse base URI to replace database name
			uri = replaceDatabase(baseURI, a.dbName)
		} else if a.dbg.OnlyManagedDatabases {
			return nil, fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
		} else {
			uri = a.replacer.Replace(a.db.URI)
			// If the URI is empty (e.g., environment variable not set), skip analysis
			if uri == "" {
				return nil, fmt.Errorf("database URI is empty (environment variable may not be set)")
			}
		}

		conn, err := sql.Open("mysql", uri)
		if err != nil {
			return nil, fmt.Errorf("failed to open mysql database: %w", err)
		}
		if err := conn.PingContext(ctx); err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to ping mysql database: %w", err)
		}
		a.conn = conn

		// Apply migrations for managed databases that were just created
		if applyMigrations {
			for _, m := range migrations {
				if len(strings.TrimSpace(m)) == 0 {
					continue
				}
				if _, err := a.conn.ExecContext(ctx, m); err != nil {
					return nil, fmt.Errorf("migration failed: %w", err)
				}
			}
		}
	}

	// Count parameters in the query
	paramCount := countParameters(query)

	// Try to prepare the statement first to validate syntax
	stmt, err := a.conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, a.extractSqlErr(n, err)
	}
	stmt.Close()

	var result core.Analysis

	// For SELECT queries, execute with default parameter values to get column metadata
	if isSelectQuery(query) {
		cols, err := a.getColumnMetadata(ctx, query, paramCount)
		if err == nil {
			result.Columns = cols
		}
		// If we fail to get column metadata, fall through to return empty columns
		// and let the catalog-based inference handle it
	}

	// Build parameter info
	for i := 1; i <= paramCount; i++ {
		name := ""
		if ps != nil {
			name, _ = ps.NameFor(i)
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

// isSelectQuery checks if a query is a SELECT statement
func isSelectQuery(query string) bool {
	trimmed := strings.TrimSpace(strings.ToUpper(query))
	return strings.HasPrefix(trimmed, "SELECT") ||
		strings.HasPrefix(trimmed, "WITH") // CTEs
}

// getColumnMetadata executes the query with default values to retrieve column information
func (a *Analyzer) getColumnMetadata(ctx context.Context, query string, paramCount int) ([]*core.Column, error) {
	// Generate default parameter values (use 1 for all - works for most types)
	args := make([]any, paramCount)
	for i := range args {
		args[i] = 1
	}

	// Wrap query to avoid fetching data: SELECT * FROM (query) AS _sqlc_wrapper LIMIT 0
	// This ensures we get column metadata without executing the actual query
	wrappedQuery := fmt.Sprintf("SELECT * FROM (%s) AS _sqlc_wrapper LIMIT 0", query)

	rows, err := a.conn.QueryContext(ctx, wrappedQuery, args...)
	if err != nil {
		// If wrapped query fails, try direct query with LIMIT 0
		// Some queries may not support being wrapped (e.g., queries with UNION at the end)
		return nil, err
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var columns []*core.Column
	for _, col := range colTypes {
		nullable, _ := col.Nullable()
		columns = append(columns, &core.Column{
			Name:     col.Name(),
			DataType: strings.ToLower(col.DatabaseTypeName()),
			NotNull:  !nullable,
		})
	}

	return columns, nil
}

// replaceDatabase replaces the database name in a MySQL DSN
func replaceDatabase(dsn string, newDB string) string {
	// MySQL DSN format: user:password@protocol(address)/dbname?params
	// We need to replace the dbname part

	// Find the slash before the database name
	slashIdx := strings.LastIndex(dsn, "/")
	if slashIdx == -1 {
		// No slash found, append /dbname
		if strings.Contains(dsn, "?") {
			// Has params, insert before ?
			paramIdx := strings.Index(dsn, "?")
			return dsn[:paramIdx] + "/" + newDB + dsn[paramIdx:]
		}
		return dsn + "/" + newDB
	}

	// Find the ? for parameters
	paramIdx := strings.Index(dsn[slashIdx:], "?")
	if paramIdx == -1 {
		// No params, replace everything after slash
		return dsn[:slashIdx+1] + newDB
	}

	// Replace database name between / and ?
	return dsn[:slashIdx+1] + newDB + dsn[slashIdx+paramIdx:]
}

// countParameters counts the number of ? placeholders in a query
func countParameters(query string) int {
	count := 0
	inString := false
	stringChar := byte(0)
	escaped := false

	for i := 0; i < len(query); i++ {
		c := query[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if inString {
			if c == stringChar {
				inString = false
			}
			continue
		}

		if c == '\'' || c == '"' || c == '`' {
			inString = true
			stringChar = c
			continue
		}

		if c == '?' {
			count++
		}
	}

	return count
}

func (a *Analyzer) extractSqlErr(n ast.Node, err error) error {
	if err == nil {
		return nil
	}
	return &sqlerr.Error{
		Message:  fmt.Sprintf("mysql: %s", err.Error()),
		Location: n.Pos(),
	}
}

func (a *Analyzer) Close(_ context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn != nil {
		a.conn.Close()
		a.conn = nil
	}

	// Note: We don't drop the database on close because:
	// 1. Other analyzers might be using the same database (based on migration hash)
	// 2. It can be reused for future runs with the same migrations
	// The databases are prefixed with sqlc_managed_ and can be cleaned up manually if needed

	if a.baseConn != nil {
		a.baseConn.Close()
		a.baseConn = nil
	}

	return nil
}

