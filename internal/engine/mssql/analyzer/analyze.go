package analyzer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	_ "github.com/microsoft/go-mssqldb"

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
	conn     *sql.DB
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

// Analyze prepares the query against the database and extracts column and parameter information
func (a *Analyzer) Analyze(ctx context.Context, n ast.Node, query string, migrations []string, ps *named.ParamSet) (*core.Analysis, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.ensureConnLocked(ctx, migrations); err != nil {
		return nil, err
	}

	// For MSSQL, we use sp_describe_first_result_set to get column metadata
	// This stored procedure returns column information for a query without executing it
	result, err := a.analyzeQuery(ctx, n, query, ps)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Analyzer) analyzeQuery(ctx context.Context, n ast.Node, query string, ps *named.ParamSet) (*core.Analysis, error) {
	var result core.Analysis

	// Use sp_describe_first_result_set to get column metadata
	// This is MSSQL's equivalent of PostgreSQL's PREPARE for getting result set metadata
	rows, err := a.conn.QueryContext(ctx, "EXEC sp_describe_first_result_set @tsql = @p1", query)
	if err != nil {
		return nil, a.extractSqlErr(n, err)
	}
	defer rows.Close()

	for rows.Next() {
		var col columnInfo
		// sp_describe_first_result_set returns many columns, we only need a few
		// Columns: is_hidden, column_ordinal, name, is_nullable, system_type_id, system_type_name,
		// max_length, precision, scale, collation_name, user_type_id, user_type_database,
		// user_type_schema, user_type_name, assembly_qualified_type_name, xml_collection_id,
		// xml_collection_database, xml_collection_schema, xml_collection_name, is_xml_document,
		// is_case_sensitive, is_fixed_length_clr_type, source_server, source_database,
		// source_schema, source_table, source_column, is_identity_column, is_part_of_unique_key,
		// is_updateable, is_computed_column, is_sparse_column_set, ordinal_in_order_by_list,
		// order_by_is_descending, order_by_list_length, tds_type_id, tds_length,
		// tds_collation_id, tds_collation_sort_id

		var isHidden bool
		var colOrdinal int
		var name sql.NullString
		var isNullable bool
		var sysTypeId int
		var sysTypeName sql.NullString
		var maxLength int
		var precision int
		var scale int
		var collationName sql.NullString
		var userTypeId sql.NullInt64
		var userTypeDb sql.NullString
		var userTypeSchema sql.NullString
		var userTypeName sql.NullString
		var assemblyQualTypeName sql.NullString
		var xmlColId sql.NullInt64
		var xmlColDb sql.NullString
		var xmlColSchema sql.NullString
		var xmlColName sql.NullString
		var isXmlDoc bool
		var isCaseSensitive bool
		var isFixedLenClr bool
		var sourceServer sql.NullString
		var sourceDb sql.NullString
		var sourceSchema sql.NullString
		var sourceTable sql.NullString
		var sourceColumn sql.NullString
		var isIdentity bool
		var isPartOfUniqueKey sql.NullBool
		var isUpdateable bool
		var isComputed bool
		var isSparseColSet bool
		var ordinalInOrderBy sql.NullInt64
		var orderByDesc sql.NullBool
		var orderByLen sql.NullInt64
		var tdsTypeId sql.NullInt64
		var tdsLength sql.NullInt64
		var tdsCollationId sql.NullInt64
		var tdsCollationSortId sql.NullInt64

		err := rows.Scan(
			&isHidden, &colOrdinal, &name, &isNullable, &sysTypeId, &sysTypeName,
			&maxLength, &precision, &scale, &collationName, &userTypeId, &userTypeDb,
			&userTypeSchema, &userTypeName, &assemblyQualTypeName, &xmlColId,
			&xmlColDb, &xmlColSchema, &xmlColName, &isXmlDoc, &isCaseSensitive,
			&isFixedLenClr, &sourceServer, &sourceDb, &sourceSchema, &sourceTable,
			&sourceColumn, &isIdentity, &isPartOfUniqueKey, &isUpdateable,
			&isComputed, &isSparseColSet, &ordinalInOrderBy, &orderByDesc,
			&orderByLen, &tdsTypeId, &tdsLength, &tdsCollationId, &tdsCollationSortId,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning column info: %w", err)
		}

		if isHidden {
			continue
		}

		col.Name = name.String
		col.IsNullable = isNullable
		col.DataType = normalizeTypeName(sysTypeName.String, maxLength, precision, scale)
		col.Table = sourceTable.String
		col.Schema = sourceSchema.String

		coreCol := &core.Column{
			Name:         col.Name,
			OriginalName: col.Name,
			DataType:     col.DataType,
			NotNull:      !col.IsNullable,
		}

		if col.Table != "" {
			coreCol.Table = &core.Identifier{
				Schema: col.Schema,
				Name:   col.Table,
			}
		}

		result.Columns = append(result.Columns, coreCol)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating column info: %w", err)
	}

	// Get parameter information
	// MSSQL doesn't have a built-in way to get parameter metadata from a query string
	// We'll count the @pN placeholders in the query and create parameters
	paramCount := countParameters(query)
	for i := 1; i <= paramCount; i++ {
		paramName := ""
		if ps != nil {
			if n, ok := ps.NameFor(i); ok {
				paramName = n
			}
		}

		result.Params = append(result.Params, &core.Parameter{
			Number: int32(i),
			Column: &core.Column{
				Name:     paramName,
				DataType: "any", // MSSQL doesn't provide parameter type info
				NotNull:  false,
			},
		})
	}

	return &result, nil
}

type columnInfo struct {
	Name       string
	DataType   string
	IsNullable bool
	Table      string
	Schema     string
}

func countParameters(query string) int {
	count := 0
	for i := 1; i <= 100; i++ {
		param := fmt.Sprintf("@p%d", i)
		if strings.Contains(query, param) {
			count = i
		}
	}
	return count
}

func normalizeTypeName(typeName string, maxLen, precision, scale int) string {
	typeName = strings.ToLower(typeName)

	// Handle common MSSQL types
	switch typeName {
	case "int":
		return "int"
	case "bigint":
		return "bigint"
	case "smallint":
		return "smallint"
	case "tinyint":
		return "tinyint"
	case "bit":
		return "bit"
	case "decimal", "numeric":
		return fmt.Sprintf("decimal(%d,%d)", precision, scale)
	case "money":
		return "money"
	case "smallmoney":
		return "smallmoney"
	case "float":
		return "float"
	case "real":
		return "real"
	case "datetime":
		return "datetime"
	case "datetime2":
		return "datetime2"
	case "date":
		return "date"
	case "time":
		return "time"
	case "datetimeoffset":
		return "datetimeoffset"
	case "smalldatetime":
		return "smalldatetime"
	case "char":
		return fmt.Sprintf("char(%d)", maxLen)
	case "varchar":
		if maxLen == -1 {
			return "varchar(max)"
		}
		return fmt.Sprintf("varchar(%d)", maxLen)
	case "nchar":
		return fmt.Sprintf("nchar(%d)", maxLen/2) // nchar uses 2 bytes per char
	case "nvarchar":
		if maxLen == -1 {
			return "nvarchar(max)"
		}
		return fmt.Sprintf("nvarchar(%d)", maxLen/2)
	case "text":
		return "text"
	case "ntext":
		return "ntext"
	case "binary":
		return fmt.Sprintf("binary(%d)", maxLen)
	case "varbinary":
		if maxLen == -1 {
			return "varbinary(max)"
		}
		return fmt.Sprintf("varbinary(%d)", maxLen)
	case "image":
		return "image"
	case "uniqueidentifier":
		return "uniqueidentifier"
	case "xml":
		return "xml"
	case "sql_variant":
		return "sql_variant"
	default:
		return typeName
	}
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
	return a.ensureConnLocked(ctx, migrations)
}

func (a *Analyzer) ensureConnLocked(ctx context.Context, migrations []string) error {
	if a.conn != nil {
		return nil
	}

	if a.dbg.OnlyManagedDatabases {
		return fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
	}

	uri := a.replacer.Replace(a.db.URI)

	conn, err := sql.Open("sqlserver", uri)
	if err != nil {
		return fmt.Errorf("failed to open mssql database: %w", err)
	}

	// Test the connection
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return fmt.Errorf("failed to ping mssql database: %w", err)
	}

	a.conn = conn

	// Apply migrations
	for _, m := range migrations {
		if len(strings.TrimSpace(m)) == 0 {
			continue
		}
		// Split by GO statements for MSSQL batch separation
		batches := splitByGO(m)
		for _, batch := range batches {
			batch = strings.TrimSpace(batch)
			if len(batch) == 0 {
				continue
			}
			if _, err := a.conn.ExecContext(ctx, batch); err != nil {
				// Check if it's a "already exists" error and skip it
				if !isObjectExistsError(err) {
					a.conn.Close()
					a.conn = nil
					return fmt.Errorf("migration failed: %s: %w", batch, err)
				}
			}
		}
	}

	return nil
}

// splitByGO splits a SQL script by GO batch separators
func splitByGO(script string) []string {
	lines := strings.Split(script, "\n")
	var batches []string
	var currentBatch strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.EqualFold(trimmed, "GO") {
			if currentBatch.Len() > 0 {
				batches = append(batches, currentBatch.String())
				currentBatch.Reset()
			}
		} else {
			currentBatch.WriteString(line)
			currentBatch.WriteString("\n")
		}
	}

	if currentBatch.Len() > 0 {
		batches = append(batches, currentBatch.String())
	}

	return batches
}

// isObjectExistsError checks if the error is about an object already existing
func isObjectExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "already exists") ||
		strings.Contains(errStr, "there is already an object")
}

// GetColumnNames implements the expander.ColumnGetter interface.
func (a *Analyzer) GetColumnNames(ctx context.Context, query string) ([]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		return nil, errors.New("database connection not initialized")
	}

	rows, err := a.conn.QueryContext(ctx, "EXEC sp_describe_first_result_set @tsql = @p1", query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var isHidden bool
		var colOrdinal int
		var name sql.NullString
		// We need to scan all columns but only care about name
		var dummy interface{}
		scanArgs := make([]interface{}, 39)
		scanArgs[0] = &isHidden
		scanArgs[1] = &colOrdinal
		scanArgs[2] = &name
		for i := 3; i < 39; i++ {
			scanArgs[i] = &dummy
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("scanning column name: %w", err)
		}

		if !isHidden && name.Valid {
			columns = append(columns, name.String)
		}
	}

	return columns, rows.Err()
}

// IntrospectSchema queries the database to build a catalog
func (a *Analyzer) IntrospectSchema(ctx context.Context, schemas []string) (*catalog.Catalog, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	// Build catalog
	cat := &catalog.Catalog{
		DefaultSchema: "dbo",
	}

	// Create schema map for quick lookup
	schemaMap := make(map[string]*catalog.Schema)
	for _, schemaName := range schemas {
		schema := &catalog.Schema{Name: schemaName}
		cat.Schemas = append(cat.Schemas, schema)
		schemaMap[schemaName] = schema
	}

	// Query tables and columns from INFORMATION_SCHEMA
	query := `
		SELECT
			c.TABLE_SCHEMA,
			c.TABLE_NAME,
			c.COLUMN_NAME,
			c.DATA_TYPE,
			CASE WHEN c.IS_NULLABLE = 'NO' THEN 1 ELSE 0 END AS NOT_NULL,
			COALESCE(
				(SELECT 1 FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
				 JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
				   ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
				  AND tc.TABLE_SCHEMA = kcu.TABLE_SCHEMA
				 WHERE tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
				   AND kcu.TABLE_SCHEMA = c.TABLE_SCHEMA
				   AND kcu.TABLE_NAME = c.TABLE_NAME
				   AND kcu.COLUMN_NAME = c.COLUMN_NAME),
				0
			) AS IS_PRIMARY_KEY
		FROM INFORMATION_SCHEMA.COLUMNS c
		WHERE c.TABLE_SCHEMA IN (SELECT value FROM STRING_SPLIT(@p1, ','))
		ORDER BY c.TABLE_SCHEMA, c.TABLE_NAME, c.ORDINAL_POSITION
	`

	rows, err := a.conn.QueryContext(ctx, query, strings.Join(schemas, ","))
	if err != nil {
		return nil, fmt.Errorf("introspect tables: %w", err)
	}
	defer rows.Close()

	// Group columns by table
	tableMap := make(map[string]*catalog.Table)
	for rows.Next() {
		var schemaName, tableName, columnName, dataType string
		var notNull, isPrimaryKey bool

		if err := rows.Scan(&schemaName, &tableName, &columnName, &dataType, &notNull, &isPrimaryKey); err != nil {
			return nil, fmt.Errorf("scanning column: %w", err)
		}

		key := schemaName + "." + tableName
		tbl, exists := tableMap[key]
		if !exists {
			tbl = &catalog.Table{
				Rel: &ast.TableName{
					Schema: schemaName,
					Name:   tableName,
				},
			}
			tableMap[key] = tbl
			if schema, ok := schemaMap[schemaName]; ok {
				schema.Tables = append(schema.Tables, tbl)
			}
		}

		tbl.Columns = append(tbl.Columns, &catalog.Column{
			Name:      columnName,
			Type:      ast.TypeName{Name: dataType},
			IsNotNull: notNull,
		})
	}

	return cat, rows.Err()
}
