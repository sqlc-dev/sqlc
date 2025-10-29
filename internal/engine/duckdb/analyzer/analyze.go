package analyzer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	core "github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/shfmt"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

type Analyzer struct {
	db       config.Database
	client   dbmanager.Client
	conn     *sql.DB
	dbg      opts.Debug
	replacer *shfmt.Replacer
	mu       sync.Mutex
}

func New(client dbmanager.Client, db config.Database) *Analyzer {
	return &Analyzer{
		db:       db,
		dbg:      opts.DebugFromEnv(),
		client:   client,
		replacer: shfmt.NewReplacer(nil),
	}
}

// Analyze extracts column and parameter information by preparing the query
// against an in-memory DuckDB instance
func (a *Analyzer) Analyze(ctx context.Context, n ast.Node, query string, migrations []string, ps *named.ParamSet) (*core.Analysis, error) {
	extractSqlErr := func(e error) error {
		// DuckDB errors don't have the same structure as PostgreSQL errors
		// but we can still wrap them appropriately
		if e == nil {
			return nil
		}
		// Try to extract position information if available
		msg := e.Error()
		return &sqlerr.Error{
			Message:  msg,
			Location: n.Pos(),
		}
	}

	a.mu.Lock()
	if a.conn == nil {
		var uri string
		if a.db.Managed {
			if a.client == nil {
				a.mu.Unlock()
				return nil, fmt.Errorf("client is nil")
			}
			edb, err := a.client.CreateDatabase(ctx, &dbmanager.CreateDatabaseRequest{
				Engine:     "duckdb",
				Migrations: migrations,
			})
			if err != nil {
				a.mu.Unlock()
				return nil, err
			}
			uri = edb.Uri
		} else if a.dbg.OnlyManagedDatabases {
			a.mu.Unlock()
			return nil, fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
		} else {
			uri = a.replacer.Replace(a.db.URI)
		}

		// If no URI is provided, use an in-memory database
		if uri == "" {
			uri = ":memory:"
		}

		conn, err := sql.Open("duckdb", uri)
		if err != nil {
			a.mu.Unlock()
			return nil, err
		}

		// Execute migrations to set up the schema
		if len(migrations) > 0 {
			for _, migration := range migrations {
				if _, err := conn.ExecContext(ctx, migration); err != nil {
					conn.Close()
					a.mu.Unlock()
					return nil, fmt.Errorf("migration failed: %w", err)
				}
			}
		}

		a.conn = conn
	}
	a.mu.Unlock()

	// Prepare the query to extract metadata
	stmt, err := a.conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, extractSqlErr(err)
	}
	defer stmt.Close()

	// Get column types from the prepared statement
	// Note: DuckDB's database/sql driver should support this via ColumnTypes()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		// If the query can't be executed (e.g., it requires parameters),
		// we can still try to get column information from the prepared statement
		// For now, return the error
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, extractSqlErr(err)
		}
	}
	if rows != nil {
		defer rows.Close()
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		// Try getting columns from the prepared statement metadata
		columns, err := rows.Columns()
		if err != nil {
			return nil, extractSqlErr(err)
		}
		// Build basic column information without type details
		var result core.Analysis
		for _, name := range columns {
			result.Columns = append(result.Columns, &core.Column{
				Name:         name,
				OriginalName: name,
				DataType:     "text", // fallback type
				NotNull:      false,
				IsArray:      false,
			})
		}
		return &result, nil
	}

	var result core.Analysis
	for _, colType := range columnTypes {
		name := colType.Name()
		dataType := colType.DatabaseTypeName()

		// Parse array types (DuckDB arrays end with [])
		isArray := strings.HasSuffix(dataType, "[]")
		if isArray {
			dataType = strings.TrimSuffix(dataType, "[]")
		}

		notNull := false
		if nullable, ok := colType.Nullable(); ok {
			notNull = !nullable
		}

		result.Columns = append(result.Columns, &core.Column{
			Name:         name,
			OriginalName: name,
			DataType:     dataType,
			NotNull:      notNull,
			IsArray:      isArray,
		})
	}

	// Note: database/sql doesn't provide a standard way to get parameter types
	// Parameter type inference will be handled by the catalog-based compiler
	// We return an empty Params slice and let sqlc infer parameter types
	// from the query structure and catalog

	return &result, nil
}

func (a *Analyzer) Close(_ context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
