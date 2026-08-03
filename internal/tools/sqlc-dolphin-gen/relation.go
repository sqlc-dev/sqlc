package main

import (
	"context"
	"database/sql"
	"strings"
)

const relationQuery = `
WITH relations AS (
    SELECT TABLE_SCHEMA, TABLE_NAME FROM information_schema.TABLES
    UNION ALL
    SELECT TABLE_SCHEMA, TABLE_NAME FROM information_schema.VIEWS
)
SELECT
    c.TABLE_SCHEMA,
    c.TABLE_NAME,
    c.ORDINAL_POSITION,
    c.COLUMN_NAME,
    c.IS_NULLABLE = 'YES',
    c.DATA_TYPE,
    c.CHARACTER_MAXIMUM_LENGTH
FROM relations r
LEFT JOIN information_schema.COLUMNS c ON (r.TABLE_SCHEMA = c.TABLE_SCHEMA AND r.TABLE_NAME = c.TABLE_NAME)

WHERE c.TABLE_SCHEMA = ?
ORDER BY 1,2,3,4,5,6,7;
`

type Relation struct {
	SchemaName string
	Name       string
	Columns    []RelationColumn
}

type RelationColumn struct {
	Name      string
	Type      string
	IsNotNull bool
	Length    *int
}

type tableRow struct {
	Schema     string
	Name       string
	Position   int
	ColumnName string
	IsNullable bool
	Type       string
	Length     *int
}

func scanRelations(rows *sql.Rows) ([]Relation, error) {
	defer rows.Close()
	// Iterate through the result set
	var relations []Relation
	var prevRel *Relation

	for rows.Next() {
		var rowData tableRow
		err := rows.Scan(
			&rowData.Schema,
			&rowData.Name,
			&rowData.Position,
			&rowData.ColumnName,
			&rowData.IsNullable,
			&rowData.Type,
			&rowData.Length,
		)
		if err != nil {
			return nil, err
		}

		if prevRel == nil || rowData.Name != prevRel.Name {
			// We are on the same table, just keep adding columns
			r := Relation{
				SchemaName: strings.ToLower(rowData.Schema),

				// Doing SELECT * FROM information_schema.TABLES looks for `tables` lowercase
				// when resolving queries
				Name: strings.ToLower(rowData.Name),
			}

			relations = append(relations, r)
			prevRel = &relations[len(relations)-1]
		}

		prevRel.Columns = append(prevRel.Columns, RelationColumn{
			Name:      rowData.ColumnName,
			Type:      rowData.Type,
			IsNotNull: !rowData.IsNullable,
			Length:    rowData.Length,
		})
	}

	return relations, rows.Err()
}

func readRelations(ctx context.Context, conn *sql.DB, schemaName string) ([]Relation, error) {
	rows, err := conn.Query(relationQuery, schemaName)
	if err != nil {
		return nil, err
	}

	return scanRelations(rows)
}
