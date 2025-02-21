// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const getColumns = `-- name: GetColumns :many
SELECT table_name::text, column_name::text from information_schema.columns
`

type GetColumnsRow struct {
	TableName  string
	ColumnName string
}

func (q *Queries) GetColumns(ctx context.Context) ([]GetColumnsRow, error) {
	rows, err := q.db.QueryContext(ctx, getColumns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetColumnsRow
	for rows.Next() {
		var i GetColumnsRow
		if err := rows.Scan(&i.TableName, &i.ColumnName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTables = `-- name: GetTables :many
SELECT table_name::text from information_schema.tables
`

func (q *Queries) GetTables(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getTables)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var table_name string
		if err := rows.Scan(&table_name); err != nil {
			return nil, err
		}
		items = append(items, table_name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTimezones = `-- name: GetTimezones :many
SELECT name, abbrev, utc_offset, is_dst from pg_catalog.pg_timezone_names
`

type GetTimezonesRow struct {
	Name      sql.NullString
	Abbrev    sql.NullString
	UtcOffset sql.NullInt64
	IsDst     sql.NullBool
}

func (q *Queries) GetTimezones(ctx context.Context) ([]GetTimezonesRow, error) {
	rows, err := q.db.QueryContext(ctx, getTimezones)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTimezonesRow
	for rows.Next() {
		var i GetTimezonesRow
		if err := rows.Scan(
			&i.Name,
			&i.Abbrev,
			&i.UtcOffset,
			&i.IsDst,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
