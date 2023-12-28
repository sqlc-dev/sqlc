// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: batch_gen.go

package querytest

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
)

var (
	ErrBatchAlreadyClosed = errors.New("batch already closed")
)

const usersB = `-- name: UsersB :batchmany
SELECT id FROM "user"
WHERE id = $1
`

type UsersBBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

func (q *Queries) UsersB(ctx context.Context, id []int64) *UsersBBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(usersB, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &UsersBBatchResults{br, len(id), false}
}

func (b *UsersBBatchResults) Query(f func(int, []int64, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var items []int64
		if b.closed {
			if f != nil {
				f(t, items, ErrBatchAlreadyClosed)
			}
			continue
		}
		err := func() error {
			rows, err := b.br.Query()
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var id int64
				if err := rows.Scan(&id); err != nil {
					return err
				}
				items = append(items, id)
			}
			return rows.Err()
		}()
		if f != nil {
			f(t, items, err)
		}
	}
}

func (b *UsersBBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}
