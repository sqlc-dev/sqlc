// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: batch.go

package querytest

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBatchAlreadyClosed = errors.New("batch already closed")
)

const getValues = `-- name: GetValues :batchmany
SELECT a, b
FROM myschema.foo
WHERE b = $1
`

type GetValuesBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

func (q *Queries) GetValues(ctx context.Context, b []pgtype.Int4) *GetValuesBatchResults {
	batch := &pgx.Batch{}
	for _, a := range b {
		vals := []interface{}{
			a,
		}
		batch.Queue(getValues, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetValuesBatchResults{br, len(b), false}
}

func (b *GetValuesBatchResults) Query(f func(int, []MyschemaFoo, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var items []MyschemaFoo
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
				var i MyschemaFoo
				if err := rows.Scan(&i.A, &i.B); err != nil {
					return err
				}
				items = append(items, i)
			}
			return rows.Err()
		}()
		if f != nil {
			f(t, items, err)
		}
	}
}

func (b *GetValuesBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}

const insertValues = `-- name: InsertValues :batchone
INSERT INTO myschema.foo (a, b)
VALUES ($1, $2)
RETURNING a
`

type InsertValuesBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type InsertValuesParams struct {
	A pgtype.Text
	B pgtype.Int4
}

func (q *Queries) InsertValues(ctx context.Context, arg []InsertValuesParams) *InsertValuesBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.A,
			a.B,
		}
		batch.Queue(insertValues, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &InsertValuesBatchResults{br, len(arg), false}
}

func (b *InsertValuesBatchResults) QueryRow(f func(int, pgtype.Text, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var a pgtype.Text
		if b.closed {
			if f != nil {
				f(t, a, ErrBatchAlreadyClosed)
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(&a)
		if f != nil {
			f(t, a, err)
		}
	}
}

func (b *InsertValuesBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}
