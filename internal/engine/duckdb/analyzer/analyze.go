package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/marcboeker/go-duckdb"

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
	typeInfo sync.Map
}

func New(client dbmanager.Client, db config.Database) *Analyzer {
	return &Analyzer{
		db:       db,
		dbg:      opts.DebugFromEnv(),
		client:   client,
		replacer: shfmt.NewReplacer(nil),
	}
}

type duckdbColumn struct {
	ColumnName string
	DataType   string
	IsNullable string
	TableName  string
	SchemaName string
}

// Analyze uses DuckDB's PREPARE and DESCRIBE to analyze queries
func (a *Analyzer) Analyze(ctx context.Context, n ast.Node, query string, migrations []string, ps *named.ParamSet) (*core.Analysis, error) {
	extractSqlErr := func(e error) error {
		if e == nil {
			return nil
		}
		// DuckDB errors don't have the same structure as PostgreSQL
		// Return a basic error for now
		return &sqlerr.Error{
			Message:  e.Error(),
			Location: n.Pos(),
		}
	}

	if a.conn == nil {
		var uri string
		if a.db.Managed {
			if a.client == nil {
				return nil, fmt.Errorf("client is nil")
			}
			edb, err := a.client.CreateDatabase(ctx, &dbmanager.CreateDatabaseRequest{
				Engine:     "duckdb",
				Migrations: migrations,
			})
			if err != nil {
				return nil, err
			}
			uri = edb.Uri
		} else if a.dbg.OnlyManagedDatabases {
			return nil, fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
		} else {
			uri = a.replacer.Replace(a.db.URI)
		}

		// DuckDB connection string
		conn, err := sql.Open("duckdb", uri)
		if err != nil {
			return nil, err
		}
		a.conn = conn
	}

	// DuckDB supports PREPARE and DESCRIBE
	// First, prepare the statement
	stmt, err := a.conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, extractSqlErr(err)
	}
	defer stmt.Close()

	var result core.Analysis

	// For DuckDB, we need to use DESCRIBE to get column information
	// This is a workaround since database/sql doesn't expose column metadata
	// without executing the query
	descQuery := fmt.Sprintf("DESCRIBE %s", query)
	rows, err := a.conn.QueryContext(ctx, descQuery)
	if err != nil {
		// If DESCRIBE fails, fall back to executing with LIMIT 0
		limitQuery := fmt.Sprintf("SELECT * FROM (%s) LIMIT 0", query)
		rows, err = a.conn.QueryContext(ctx, limitQuery)
		if err != nil {
			return nil, extractSqlErr(err)
		}
	}
	defer rows.Close()

	// Get column types from the result set
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	for _, ct := range columnTypes {
		dataType := ct.DatabaseTypeName()
		notNull := false
		if nullable, ok := ct.Nullable(); ok {
			notNull = !nullable
		}

		// Parse array types
		isArray := strings.HasSuffix(dataType, "[]")
		if isArray {
			dataType = strings.TrimSuffix(dataType, "[]")
		}

		result.Columns = append(result.Columns, &core.Column{
			Name:         ct.Name(),
			OriginalName: ct.Name(),
			DataType:     normalizeDuckDBType(dataType),
			NotNull:      notNull,
			IsArray:      isArray,
			ArrayDims:    0,
		})
	}

	// For parameters, we don't have detailed type information from PREPARE
	// We'll need to infer from the query or use placeholders
	// DuckDB uses $1, $2, etc. for parameters
	paramCount := strings.Count(query, "$")
	for i := 0; i < paramCount; i++ {
		name := ""
		if ps != nil {
			name, _ = ps.NameFor(i + 1)
		}
		result.Params = append(result.Params, &core.Parameter{
			Number: int32(i + 1),
			Column: &core.Column{
				Name:     name,
				DataType: "any", // DuckDB doesn't provide parameter types without execution
				NotNull:  false,
			},
		})
	}

	return &result, nil
}

func (a *Analyzer) Close(_ context.Context) error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

// normalizeDuckDBType converts DuckDB types to sqlc-compatible types
func normalizeDuckDBType(duckdbType string) string {
	upper := strings.ToUpper(duckdbType)
	switch upper {
	case "INTEGER", "INT", "INT4":
		return "integer"
	case "BIGINT", "INT8", "LONG":
		return "bigint"
	case "SMALLINT", "INT2", "SHORT":
		return "smallint"
	case "TINYINT", "INT1":
		return "tinyint"
	case "DOUBLE", "FLOAT8":
		return "double"
	case "REAL", "FLOAT4", "FLOAT":
		return "real"
	case "VARCHAR", "TEXT", "STRING":
		return "varchar"
	case "BOOLEAN", "BOOL":
		return "boolean"
	case "DATE":
		return "date"
	case "TIME":
		return "time"
	case "TIMESTAMP":
		return "timestamp"
	case "TIMESTAMPTZ", "TIMESTAMP WITH TIME ZONE":
		return "timestamptz"
	case "BLOB", "BYTEA", "BINARY", "VARBINARY":
		return "bytea"
	case "UUID":
		return "uuid"
	case "JSON":
		return "json"
	case "DECIMAL", "NUMERIC":
		return "decimal"
	case "HUGEINT":
		return "hugeint"
	case "UINTEGER", "UINT4":
		return "uinteger"
	case "UBIGINT", "UINT8":
		return "ubigint"
	case "USMALLINT", "UINT2":
		return "usmallint"
	case "UTINYINT", "UINT1":
		return "utinyint"
	default:
		return strings.ToLower(duckdbType)
	}
}
