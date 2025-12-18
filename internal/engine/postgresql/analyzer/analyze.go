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
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
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

// SQL queries for schema introspection
const introspectTablesQuery = `
SELECT
    n.nspname AS schema_name,
    c.relname AS table_name,
    a.attname AS column_name,
    pg_catalog.format_type(a.atttypid, a.atttypmod) AS data_type,
    a.attnotnull AS not_null,
    a.attndims AS array_dims,
    COALESCE(
        (SELECT true FROM pg_index i
         WHERE i.indrelid = c.oid
         AND i.indisprimary
         AND a.attnum = ANY(i.indkey)),
        false
    ) AS is_primary_key
FROM
    pg_catalog.pg_class c
JOIN
    pg_catalog.pg_namespace n ON n.oid = c.relnamespace
JOIN
    pg_catalog.pg_attribute a ON a.attrelid = c.oid
WHERE
    c.relkind IN ('r', 'v', 'p')  -- tables, views, partitioned tables
    AND a.attnum > 0  -- skip system columns
    AND NOT a.attisdropped
    AND n.nspname = ANY($1)
ORDER BY
    n.nspname, c.relname, a.attnum
`

const introspectEnumsQuery = `
SELECT
    n.nspname AS schema_name,
    t.typname AS type_name,
    e.enumlabel AS enum_value
FROM
    pg_catalog.pg_type t
JOIN
    pg_catalog.pg_namespace n ON n.oid = t.typnamespace
JOIN
    pg_catalog.pg_enum e ON e.enumtypid = t.oid
WHERE
    t.typtype = 'e'
    AND n.nspname = ANY($1)
ORDER BY
    n.nspname, t.typname, e.enumsortorder
`

type introspectedColumn struct {
	SchemaName   string `db:"schema_name"`
	TableName    string `db:"table_name"`
	ColumnName   string `db:"column_name"`
	DataType     string `db:"data_type"`
	NotNull      bool   `db:"not_null"`
	ArrayDims    int    `db:"array_dims"`
	IsPrimaryKey bool   `db:"is_primary_key"`
}

type introspectedEnum struct {
	SchemaName string `db:"schema_name"`
	TypeName   string `db:"type_name"`
	EnumValue  string `db:"enum_value"`
}

// IntrospectSchema queries the database to build a catalog containing
// tables, columns, and enum types for the specified schemas.
func (a *Analyzer) IntrospectSchema(ctx context.Context, schemas []string) (*catalog.Catalog, error) {
	if a.pool == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	c, err := a.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Release()

	// Query tables and columns
	rows, err := c.Query(ctx, introspectTablesQuery, schemas)
	if err != nil {
		return nil, fmt.Errorf("introspect tables: %w", err)
	}
	columns, err := pgx.CollectRows(rows, pgx.RowToStructByName[introspectedColumn])
	if err != nil {
		return nil, fmt.Errorf("collect table rows: %w", err)
	}

	// Query enums
	enumRows, err := c.Query(ctx, introspectEnumsQuery, schemas)
	if err != nil {
		return nil, fmt.Errorf("introspect enums: %w", err)
	}
	enums, err := pgx.CollectRows(enumRows, pgx.RowToStructByName[introspectedEnum])
	if err != nil {
		return nil, fmt.Errorf("collect enum rows: %w", err)
	}

	// Build catalog
	cat := &catalog.Catalog{
		DefaultSchema: "public",
		SearchPath:    schemas,
	}

	// Create schema map for quick lookup
	schemaMap := make(map[string]*catalog.Schema)
	for _, schemaName := range schemas {
		schema := &catalog.Schema{Name: schemaName}
		cat.Schemas = append(cat.Schemas, schema)
		schemaMap[schemaName] = schema
	}

	// Group columns by table
	tableMap := make(map[string]*catalog.Table)
	for _, col := range columns {
		key := col.SchemaName + "." + col.TableName
		tbl, exists := tableMap[key]
		if !exists {
			tbl = &catalog.Table{
				Rel: &ast.TableName{
					Schema: col.SchemaName,
					Name:   col.TableName,
				},
			}
			tableMap[key] = tbl
			if schema, ok := schemaMap[col.SchemaName]; ok {
				schema.Tables = append(schema.Tables, tbl)
			}
		}

		dt, isArray, dims := parseType(col.DataType)
		tbl.Columns = append(tbl.Columns, &catalog.Column{
			Name:      col.ColumnName,
			Type:      ast.TypeName{Name: dt},
			IsNotNull: col.NotNull,
			IsArray:   isArray || col.ArrayDims > 0,
			ArrayDims: max(dims, col.ArrayDims),
		})
	}

	// Group enum values by type
	enumMap := make(map[string]*catalog.Enum)
	for _, e := range enums {
		key := e.SchemaName + "." + e.TypeName
		enum, exists := enumMap[key]
		if !exists {
			enum = &catalog.Enum{
				Name: e.TypeName,
			}
			enumMap[key] = enum
			if schema, ok := schemaMap[e.SchemaName]; ok {
				schema.Types = append(schema.Types, enum)
			}
		}
		enum.Vals = append(enum.Vals, e.EnumValue)
	}

	return cat, nil
}

// EnsureConn initializes the database connection pool if not already done.
// This is useful for database-only mode where we need to connect before analyzing queries.
func (a *Analyzer) EnsureConn(ctx context.Context, migrations []string) error {
	if a.pool != nil {
		return nil
	}

	var uri string
	if a.db.Managed {
		if a.client == nil {
			return fmt.Errorf("client is nil")
		}
		edb, err := a.client.CreateDatabase(ctx, &dbmanager.CreateDatabaseRequest{
			Engine:     "postgresql",
			Migrations: migrations,
		})
		if err != nil {
			return err
		}
		uri = edb.Uri
	} else if a.dbg.OnlyManagedDatabases {
		return fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
	} else {
		uri = a.replacer.Replace(a.db.URI)
	}

	conf, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return err
	}
	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return err
	}
	a.pool = pool
	return nil
}

// GetColumnNames implements the expander.ColumnGetter interface.
// It prepares a query and returns the column names from the result set description.
func (a *Analyzer) GetColumnNames(ctx context.Context, query string) ([]string, error) {
	if a.pool == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	conn, err := a.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	desc, err := conn.Conn().Prepare(ctx, "", query)
	if err != nil {
		return nil, err
	}

	columns := make([]string, len(desc.Fields))
	for i, field := range desc.Fields {
		columns[i] = field.Name
	}

	return columns, nil
}
