// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertOrders = `-- name: InsertOrders :exec
insert into Orders (id,name)
select id , CASE WHEN $1::BOOLEAN THEN $2 ELSE s.name END 
from Orders s
`

type InsertOrdersParams struct {
	NameDoUpdate pgtype.Bool
	Name         pgtype.Text
}

func (q *Queries) InsertOrders(ctx context.Context, arg InsertOrdersParams) error {
	_, err := q.db.Exec(ctx, insertOrders, arg.NameDoUpdate, arg.Name)
	return err
}
