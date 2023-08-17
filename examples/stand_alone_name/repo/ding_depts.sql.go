// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: ding_depts.sql

package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type DingDeptDBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func NewDingDept(db DingDeptDBTX) *DingDeptQueries {
	return &DingDeptQueries{db: db}
}

type DingDeptQueries struct {
	db DingDeptDBTX
}

func (q *DingDeptQueries) WithTx(tx pgx.Tx) *DingDeptQueries {
	return &DingDeptQueries{
		db: tx,
	}
}

const dingDeptCountById = `-- name: DingDeptCountById :one
SELECT count(*)
FROM ding_depts
where id = $1
`

func (q *DingDeptQueries) DingDeptCountById(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRow(ctx, dingDeptCountById, id)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const dingDeptCreate = `-- name: DingDeptCreate :exec
insert into ding_depts (id, pid, title)
values ($1, $2, $3)
`

type DingDeptCreateParams struct {
	ID    int64       `json:"id"`
	Pid   pgtype.Int8 `json:"pid"`
	Title pgtype.Text `json:"title"`
}

func (q *DingDeptQueries) DingDeptCreate(ctx context.Context, arg DingDeptCreateParams) error {
	_, err := q.db.Exec(ctx, dingDeptCreate, arg.ID, arg.Pid, arg.Title)
	return err
}

const dingDeptListByPid = `-- name: DingDeptListByPid :many
SELECT id, pid, title
FROM ding_depts
where pid = $1
`

func (q *DingDeptQueries) DingDeptListByPid(ctx context.Context, pid pgtype.Int8) ([]DingDept, error) {
	rows, err := q.db.Query(ctx, dingDeptListByPid, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DingDept
	for rows.Next() {
		var i DingDept
		if err := rows.Scan(&i.ID, &i.Pid, &i.Title); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
