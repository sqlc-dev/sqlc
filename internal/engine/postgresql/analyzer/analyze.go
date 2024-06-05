package analyzer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

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
	pool     *pgxpool.Pool
	dbg      opts.Debug
	replacer *shfmt.Replacer
	formats  sync.Map
	columns  sync.Map
	tables   sync.Map
}

func New(client dbmanager.Client, db config.Database) *Analyzer {
	return &Analyzer{
		db:       db,
		dbg:      opts.DebugFromEnv(),
		client:   client,
		replacer: shfmt.NewReplacer(nil),
	}
}

const columnQuery = `
SELECT
    pg_catalog.format_type(pg_attribute.atttypid, pg_attribute.atttypmod) AS data_type,
	pg_attribute.attnotnull as not_null,
	pg_attribute.attndims as array_dims
FROM
    pg_catalog.pg_attribute
WHERE
    attrelid = $1
	AND attnum = $2;
`

const tableQuery = `
SELECT
    pg_class.relname as table_name,
    pg_namespace.nspname as schema_name
FROM
    pg_catalog.pg_class 
JOIN
    pg_catalog.pg_namespace ON pg_namespace.oid = pg_class.relnamespace
WHERE
    pg_class.oid = $1;
`

type pgTable struct {
	TableName  string `db:"table_name"`
	SchemaName string `db:"schema_name"`
}

// Cache these types in memory
func (a *Analyzer) tableInfo(ctx context.Context, oid uint32) (*pgTable, error) {
	ctbl, ok := a.tables.Load(oid)
	if ok {
		return ctbl.(*pgTable), nil
	}
	rows, err := a.pool.Query(ctx, tableQuery, oid)
	if err != nil {
		return nil, err
	}
	tbl, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pgTable])
	if err != nil {
		return nil, err
	}
	a.tables.Store(oid, &tbl)
	return &tbl, nil
}

type pgColumn struct {
	DataType  string `db:"data_type"`
	NotNull   bool   `db:"not_null"`
	ArrayDims int    `db:"array_dims"`
}

type columnKey struct {
	OID  uint32
	Attr uint16
}

// Cache these types in memory
func (a *Analyzer) columnInfo(ctx context.Context, field pgconn.FieldDescription) (*pgColumn, error) {
	key := columnKey{field.TableOID, field.TableAttributeNumber}
	cinfo, ok := a.columns.Load(key)
	if ok {
		return cinfo.(*pgColumn), nil
	}
	rows, err := a.pool.Query(ctx, columnQuery, field.TableOID, int16(field.TableAttributeNumber))
	if err != nil {
		return nil, err
	}
	col, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pgColumn])
	if err != nil {
		return nil, err
	}
	a.columns.Store(key, &col)
	return &col, nil
}

type formatKey struct {
	OID      uint32
	Modified int32
}

// TODO: Use PGX to do the lookup for basic OID types
func (a *Analyzer) formatType(ctx context.Context, oid uint32, modifier int32) (string, error) {
	key := formatKey{oid, modifier}
	ftyp, ok := a.formats.Load(key)
	if ok {
		return ftyp.(string), nil
	}
	rows, err := a.pool.Query(ctx, `SELECT format_type($1, $2)`, oid, modifier)
	if err != nil {
		return "", err
	}
	dt, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	if err != nil {
		return "", err
	}
	a.formats.Store(key, dt)
	return dt, err
}

// TODO: This is bad
func rewriteType(dt string) string {
	switch {
	case strings.HasPrefix(dt, "character("):
		return "pg_catalog.bpchar"
	case strings.HasPrefix(dt, "character varying"):
		return "pg_catalog.varchar"
	case strings.HasPrefix(dt, "bit varying"):
		return "pg_catalog.varbit"
	case strings.HasPrefix(dt, "bit("):
		return "pg_catalog.bit"
	}
	switch dt {
	case "bpchar":
		return "pg_catalog.bpchar"
	case "timestamp without time zone":
		return "pg_catalog.timestamp"
	case "timestamp with time zone":
		return "pg_catalog.timestamptz"
	case "time without time zone":
		return "pg_catalog.time"
	case "time with time zone":
		return "pg_catalog.timetz"
	}
	return dt
}

func parseType(dt string) (string, bool, int) {
	size := 0
	for {
		trimmed := strings.TrimSuffix(dt, "[]")
		if trimmed == dt {
			return rewriteType(dt), size > 0, size
		}
		size += 1
		dt = trimmed
	}
}

// Don't create a database per query
func (a *Analyzer) Analyze(ctx context.Context, n ast.Node, query string, migrations []string, ps *named.ParamSet) (*core.Analysis, error) {
	extractSqlErr := func(e error) error {
		var pgErr *pgconn.PgError
		if errors.As(e, &pgErr) {
			return &sqlerr.Error{
				Code:     pgErr.Code,
				Message:  pgErr.Message,
				Location: max(n.Pos()+int(pgErr.Position)-1, 0),
			}
		}
		return e
	}

	if a.pool == nil {
		var uri string
		if a.db.Managed {
			if a.client == nil {
				return nil, fmt.Errorf("client is nil")
			}
			edb, err := a.client.CreateDatabase(ctx, &dbmanager.CreateDatabaseRequest{
				Engine:     "postgresql",
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
		conf, err := pgxpool.ParseConfig(uri)
		if err != nil {
			return nil, err
		}
		pool, err := pgxpool.NewWithConfig(ctx, conf)
		if err != nil {
			return nil, err
		}
		a.pool = pool
	}

	c, err := a.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Release()

	// TODO: Pick a random name
	desc, err := c.Conn().Prepare(ctx, "foo", query)
	if err != nil {
		return nil, extractSqlErr(err)
	}

	if err := c.Conn().Deallocate(ctx, "foo"); err != nil {
		return nil, err
	}

	var result core.Analysis
	for _, field := range desc.Fields {
		if field.TableOID > 0 {
			col, err := a.columnInfo(ctx, field)
			if err != nil {
				return nil, err
			}
			// debug.Dump(i, field, col)
			tbl, err := a.tableInfo(ctx, field.TableOID)
			if err != nil {
				return nil, err
			}
			dt, isArray, dims := parseType(col.DataType)
			notNull := col.NotNull
			name := field.Name
			result.Columns = append(result.Columns, &core.Column{
				Name:         name,
				OriginalName: field.Name,
				DataType:     dt,
				NotNull:      notNull,
				IsArray:      isArray,
				ArrayDims:    int32(max(col.ArrayDims, dims)),
				Table: &core.Identifier{
					Schema: tbl.SchemaName,
					Name:   tbl.TableName,
				},
			})
		} else {
			dataType, err := a.formatType(ctx, field.DataTypeOID, field.TypeModifier)
			if err != nil {
				return nil, err
			}
			// debug.Dump(i, field, dataType)
			notNull := false
			name := field.Name
			dt, isArray, dims := parseType(dataType)
			result.Columns = append(result.Columns, &core.Column{
				Name:         name,
				OriginalName: field.Name,
				DataType:     dt,
				NotNull:      notNull,
				IsArray:      isArray,
				ArrayDims:    int32(dims),
			})
		}
	}

	for i, oid := range desc.ParamOIDs {
		dataType, err := a.formatType(ctx, oid, -1)
		if err != nil {
			return nil, err
		}
		notNull := false
		dt, isArray, dims := parseType(dataType)
		name := ""
		if ps != nil {
			name, _ = ps.NameFor(i + 1)
		}
		result.Params = append(result.Params, &core.Parameter{
			Number: int32(i + 1),
			Column: &core.Column{
				Name:      name,
				DataType:  dt,
				IsArray:   isArray,
				ArrayDims: int32(dims),
				NotNull:   notNull,
			},
		})
	}

	return &result, nil
}

func (a *Analyzer) Close(_ context.Context) error {
	if a.pool != nil {
		a.pool.Close()
	}
	return nil
}
