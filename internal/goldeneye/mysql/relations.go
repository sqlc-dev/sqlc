package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// systemSchemas are the schemas whose relations are generated: the data
// dictionary's views, which are what a query of the server's catalog
// reads. The other schemas mysqld --initialize creates — mysql,
// performance_schema and sys — are tables rather than views, and every
// table a dialect seeds is one the analysis core hands codegen as a model,
// so they wait until codegen knows a system schema when it sees one.
var systemSchemas = []string{
	"information_schema",
}

// relationQuery lists a schema's relations and their columns, in the order
// the columns were declared, so that the output is stable.
const relationQuery = `
SELECT t.TABLE_NAME, t.TABLE_TYPE,
  c.COLUMN_NAME, c.DATA_TYPE, c.COLUMN_TYPE, c.IS_NULLABLE
FROM information_schema.TABLES t
JOIN information_schema.COLUMNS c
  ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME
WHERE t.TABLE_SCHEMA = ?
ORDER BY t.TABLE_NAME, c.ORDINAL_POSITION`

// readRelations reads a schema's tables and views. Their names are written
// in lower case: information_schema names its views in upper case and
// matches them in any case, MySQL matches every column name in any case,
// and sqlc's MySQL parser lowercases every identifier, so lower case is
// how a query reaches them. A column's type is its data type as MySQL
// names it, with UNSIGNED kept as part of the name the way a column
// declaration spells it; the length and precision a declaration adds are
// not part of the type.
func readRelations(ctx context.Context, conn *sql.Conn, schema string) ([]dialect.Relation, error) {
	rows, err := conn.QueryContext(ctx, relationQuery, schema)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", schema, err)
	}
	defer rows.Close()
	var relations []dialect.Relation
	var cur *dialect.Relation
	for rows.Next() {
		var name, kind, colName, dataType, columnType, nullable string
		if err := rows.Scan(&name, &kind, &colName, &dataType, &columnType, &nullable); err != nil {
			return nil, err
		}
		name = strings.ToLower(name)
		if cur == nil || cur.Name != name {
			rel := dialect.Relation{Catalog: "def", Schema: schema, Name: name}
			if kind != "BASE TABLE" {
				rel.Kind = "v"
			}
			relations = append(relations, rel)
			cur = &relations[len(relations)-1]
		}
		cur.Columns = append(cur.Columns, dialect.Column{
			Name:    strings.ToLower(colName),
			Type:    typeName(dataType, columnType),
			NotNull: nullable == "NO",
		})
	}
	return relations, rows.Err()
}

// typeName spells a column's type the way a declaration does, from the
// DATA_TYPE and COLUMN_TYPE information_schema reports for it: "bigint",
// or "bigint unsigned" when the column is unsigned.
func typeName(dataType, columnType string) string {
	name := strings.ToLower(dataType)
	if strings.Contains(strings.ToLower(columnType), " unsigned") {
		name += " unsigned"
	}
	return name
}
