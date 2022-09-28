package main

import (
	"context"

	pgx "github.com/jackc/pgx/v4"
)

// Relations are the relations available in pg_tables and pg_views
// such as pg_catalog.pg_timezone_names
const relationQuery = `
with relations as (
	select schemaname, tablename as name from pg_catalog.pg_tables
	UNION ALL
	select schemaname, viewname as name from pg_catalog.pg_views
)
select
	relations.schemaname,
	relations.name as tablename,
	pg_attribute.attname as column_name,
	attnotnull as column_notnull,
	column_type.typname as column_type,
	nullif(column_type.typlen, -1) as column_length,
	column_type.typcategory = 'A' as column_isarray
from relations
inner join pg_catalog.pg_class on pg_class.relname = relations.name
left join pg_catalog.pg_attribute on pg_attribute.attrelid = pg_class.oid
inner join pg_catalog.pg_type column_type on pg_attribute.atttypid = column_type.oid
where relations.schemaname = $1
-- Make sure these columns are always generated in the same order
-- so that the output is stable
order by
 relations.schemaname ASC,
 relations.name ASC,
 pg_attribute.attnum ASC
`

type Relation struct {
	Catalog    string
	SchemaName string
	Name       string
	Columns    []RelationColumn
}

type RelationColumn struct {
	Name      string
	Type      string
	IsNotNull bool
	IsArray   bool
	Length    *int
}

func scanRelations(rows pgx.Rows) ([]Relation, error) {
	defer rows.Close()
	// Iterate through the result set
	var relations []Relation
	var prevRel *Relation

	for rows.Next() {
		var schemaName string
		var tableName string
		var columnName string
		var columnNotNull bool
		var columnType string
		var columnLength *int
		var columnIsArray bool
		err := rows.Scan(
			&schemaName,
			&tableName,
			&columnName,
			&columnNotNull,
			&columnType,
			&columnLength,
			&columnIsArray,
		)
		if err != nil {
			return nil, err
		}

		if prevRel == nil || tableName != prevRel.Name {
			// We are on the same table, just keep adding columns
			r := Relation{
				Catalog:    "pg_catalog",
				SchemaName: schemaName,
				Name:       tableName,
			}

			relations = append(relations, r)
			prevRel = &relations[len(relations)-1]
		}

		prevRel.Columns = append(prevRel.Columns, RelationColumn{
			Name:      columnName,
			Type:      columnType,
			IsNotNull: columnNotNull,
			IsArray:   columnIsArray,
			Length:    columnLength,
		})
	}

	return relations, rows.Err()
}

func readRelations(ctx context.Context, conn *pgx.Conn, schemaName string) ([]Relation, error) {
	rows, err := conn.Query(ctx, relationQuery, schemaName)
	if err != nil {
		return nil, err
	}

	return scanRelations(rows)
}
