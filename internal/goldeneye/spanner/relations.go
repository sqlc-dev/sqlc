package spanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// relationQuery lists the relations of every schema but the database's
// own, unnamed one — INFORMATION_SCHEMA and SPANNER_SYS — with their
// columns in declared order, so that the output is stable. A column's
// type is spelled the way SPANNER_TYPE spells a declaration: STRING(MAX),
// ARRAY<INT64>.
const relationQuery = `
SELECT t.TABLE_SCHEMA, t.TABLE_NAME, t.TABLE_TYPE,
  c.COLUMN_NAME, c.SPANNER_TYPE, c.IS_NULLABLE
FROM INFORMATION_SCHEMA.TABLES AS t
JOIN INFORMATION_SCHEMA.COLUMNS AS c
  ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME
WHERE t.TABLE_SCHEMA != ''
ORDER BY t.TABLE_SCHEMA, t.TABLE_NAME, c.ORDINAL_POSITION`

// readRelations reads the views of the system schemas. Their names are
// kept as the catalog spells them, in upper case, which is how a query
// names them; a column's type is spelled in lower case, the way the
// dialect's types.jsonl spells its types, with an ARRAY<T> written as T
// with the array flag, since that is how a seed spells an array.
func readRelations(ctx context.Context, s *server, session string) ([]dialect.Relation, error) {
	rows, err := s.query(ctx, session, relationQuery)
	if err != nil {
		return nil, err
	}
	var relations []dialect.Relation
	var cur *dialect.Relation
	for _, row := range rows {
		if len(row) != 6 {
			return nil, fmt.Errorf("expected 6 columns, got %d", len(row))
		}
		schema, name, kind := row[0].GetStringValue(), row[1].GetStringValue(), row[2].GetStringValue()
		if cur == nil || cur.Schema != schema || cur.Name != name {
			rel := dialect.Relation{Schema: schema, Name: name}
			if kind == "VIEW" {
				rel.Kind = "v"
			}
			relations = append(relations, rel)
			cur = &relations[len(relations)-1]
		}
		col := dialect.Column{
			Name:    row[3].GetStringValue(),
			NotNull: row[5].GetStringValue() == "NO",
		}
		col.Type, col.Array = typeName(row[4].GetStringValue())
		cur.Columns = append(cur.Columns, col)
	}
	return relations, nil
}

// typeName spells a SPANNER_TYPE the way a seed spells a column's type:
// in lower case, with an ARRAY<T> as its element and the array flag, a
// STRUCT<a T, b U> as struct(a: t, b: u) and a PROTO<p.M> as proto('p.M'),
// since a seed spells a type's arguments in parentheses.
func typeName(spannerType string) (string, bool) {
	t := strings.TrimSpace(spannerType)
	if element, ok := strings.CutPrefix(t, "ARRAY<"); ok && strings.HasSuffix(element, ">") {
		name, _ := typeName(strings.TrimSuffix(element, ">"))
		return name, true
	}
	if fields, ok := strings.CutPrefix(t, "STRUCT<"); ok && strings.HasSuffix(fields, ">") {
		var args []string
		for _, f := range splitTop(strings.TrimSuffix(fields, ">"), ',') {
			f = strings.TrimSpace(f)
			label, typ := "", f
			if i := strings.IndexAny(f, " \t"); i > 0 && !strings.ContainsAny(f[:i], "<(") {
				label, typ = f[:i], strings.TrimSpace(f[i+1:])
			}
			name, array := typeName(typ)
			if array {
				name += "[]"
			}
			if label != "" {
				name = label + ": " + name
			}
			args = append(args, name)
		}
		return "struct(" + strings.Join(args, ", ") + ")", false
	}
	if message, ok := strings.CutPrefix(t, "PROTO<"); ok && strings.HasSuffix(message, ">") {
		return "proto('" + strings.TrimSuffix(message, ">") + "')", false
	}
	return strings.ToLower(t), false
}
