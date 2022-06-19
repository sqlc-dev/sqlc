// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package hstore

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db, observer: noopObserver}
}

type Queries struct {
	db DBTX

	observer func(ctx context.Context, methodName string) (context.Context, func(err error) error)
}

func noopObserver(ctx context.Context, methodName string) (context.Context, func(err error) error) {
	return ctx, func(err error) error { return err }
}

// WithObserver can be used to observe queries (metric, log, trace, ...)
// Example usage:
// 	queries.WithObserver(func (ctx context.Context, methodName string) (context.Context, func(err error) error) {
// 		spanCtx, span := tracer.Start(ctx, methodName)
// 		startTime := time.New()
// 		return spanCtx, func(err error) error {
// 			log.Println("Query %q executed in %s", methodName, time.Since(startTime))
// 			span.End()
// 			return err
// 		}
// 	})
func (q *Queries) WithObserver(observer func(ctx context.Context, methodName string) (context.Context, func(err error) error)) *Queries {
	return &Queries{

		db: q.db,

		observer: observer,
	}
}

func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{
		db:       tx,
		observer: q.observer,
	}
}
