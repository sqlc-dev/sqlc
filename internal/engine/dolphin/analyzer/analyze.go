package analyzer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"hash/fnv"
	"io"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"

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

	// Get metadata directly from prepared statement via driver connection
	result, err := a.getStatementMetadata(ctx, n, query, ps)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// getStatementMetadata uses the MySQL driver's prepared statement metadata API
// to get column and parameter type information without executing the query
func (a *Analyzer) getStatementMetadata(ctx context.Context, n ast.Node, query string, ps *named.ParamSet) (*core.Analysis, error) {
	var result core.Analysis

	// Get a raw connection to access driver-level prepared statement
	conn, err := a.conn.Conn(ctx)
	if err != nil {
		return nil, a.extractSqlErr(n, fmt.Errorf("failed to get connection: %w", err))
	}
	defer conn.Close()

	err = conn.Raw(func(driverConn any) error {
		// Get the driver connection that supports PrepareContext
		preparer, ok := driverConn.(driver.ConnPrepareContext)
		if !ok {
			return fmt.Errorf("driver connection does not support PrepareContext")
		}

		// Prepare the statement - this sends COM_STMT_PREPARE to MySQL
		// and receives column and parameter metadata
		stmt, err := preparer.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		// Access the metadata via the StmtMetadata interface from our forked driver
		meta, ok := stmt.(mysql.StmtMetadata)
		if !ok {
			// Fallback: just use param count from NumInput
			paramCount := stmt.NumInput()
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
			return nil
		}

		// Get column metadata
		for _, col := range meta.ColumnMetadata() {
			result.Columns = append(result.Columns, &core.Column{
				Name:     col.Name,
				DataType: strings.ToLower(col.DatabaseTypeName),
				NotNull:  !col.Nullable,
				Unsigned: col.Unsigned,
				Length:   int32(col.Length),
			})
		}

		// Get parameter metadata
		paramMeta := meta.ParamMetadata()
		for i, param := range paramMeta {
			name := ""
			if ps != nil {
				name, _ = ps.NameFor(i + 1)
			}
			result.Params = append(result.Params, &core.Parameter{
				Number: int32(i + 1),
				Column: &core.Column{
					Name:     name,
					DataType: strings.ToLower(param.DatabaseTypeName),
					NotNull:  !param.Nullable,
					Unsigned: param.Unsigned,
					Length:   int32(param.Length),
				},
			})
		}

		return nil
	})

	if err != nil {
		return nil, a.extractSqlErr(n, err)
	}

	return &result, nil
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

