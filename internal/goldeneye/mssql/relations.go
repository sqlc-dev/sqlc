package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// systemSchemas are the schemas whose relations are generated: the catalog
// views, which are what a query of the server's catalog reads. Both hold
// nothing but views, so every relation a dialect seeds from them is one
// the analysis core never hands codegen as a model.
var systemSchemas = []string{
	"sys",
	"INFORMATION_SCHEMA",
}

// viewQuery lists a schema's views in name order, so that the output is
// stable.
const viewQuery = `
SELECT o.name
FROM sys.all_objects o
JOIN sys.schemas s ON s.schema_id = o.schema_id
WHERE o.type = 'V' AND s.name = @p1
ORDER BY o.name`

// describeQuery asks the server to describe a SELECT * from a view: each
// column's name, its type spelled the way a declaration spells it —
// nvarchar(128), decimal(10,2), varbinary(max) — and whether it can be
// NULL. A view the server cannot describe reports an error message on a
// row of its own rather than failing the query.
const describeQuery = `
SELECT column_ordinal, name, system_type_name, is_nullable, error_message
FROM sys.dm_exec_describe_first_result_set(@p1, NULL, 0)
ORDER BY column_ordinal`

// readRelations reads the views of the system schemas. Their names are
// written in lower case: SQL Server matches them in any case under its
// default collations, INFORMATION_SCHEMA spells its own in upper case, and
// sqlc's SQL Server parser lowercases every identifier, so lower case is
// how a query reaches them. The type of each column is spelled the way the
// server describes it to a driver, which is the way a declaration spells
// it, and a column of the type SQL Server calls timestamp is written as
// the rowversion the dialect names it.
func readRelations(ctx context.Context, conn *sql.Conn) ([]dialect.Relation, error) {
	var relations []dialect.Relation
	for _, schema := range systemSchemas {
		names, err := views(ctx, conn, schema)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", schema, err)
		}
		for _, name := range names {
			columns, err := describe(ctx, conn, schema, name)
			if err != nil {
				return nil, fmt.Errorf("%s.%s: %w", schema, name, err)
			}
			relations = append(relations, dialect.Relation{
				Schema:  strings.ToLower(schema),
				Name:    strings.ToLower(name),
				Kind:    "v",
				Columns: columns,
			})
		}
	}
	return relations, nil
}

func views(ctx context.Context, conn *sql.Conn, schema string) ([]string, error) {
	rows, err := conn.QueryContext(ctx, viewQuery, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

func describe(ctx context.Context, conn *sql.Conn, schema, name string) ([]dialect.Column, error) {
	rows, err := conn.QueryContext(ctx, describeQuery, "SELECT * FROM "+quote(schema)+"."+quote(name))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var columns []dialect.Column
	for rows.Next() {
		var (
			ordinal  int
			colName  sql.NullString
			typeName sql.NullString
			nullable sql.NullBool
			errMsg   sql.NullString
		)
		if err := rows.Scan(&ordinal, &colName, &typeName, &nullable, &errMsg); err != nil {
			return nil, err
		}
		if errMsg.Valid {
			return nil, fmt.Errorf("cannot be described: %s", errMsg.String)
		}
		columns = append(columns, dialect.Column{
			Name:    strings.ToLower(colName.String),
			Type:    typeName.String,
			NotNull: !nullable.Bool,
		})
	}
	return columns, rows.Err()
}

// quote brackets an identifier the way T-SQL does.
func quote(name string) string {
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
}
