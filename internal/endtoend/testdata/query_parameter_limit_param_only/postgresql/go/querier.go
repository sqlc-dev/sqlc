// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"context"
	"database/sql"
	"time"
)

type Querier interface {
	CreateNotice(ctx context.Context, cnt int32, createdAt time.Time) error
	MarkNoticeDone(ctx context.Context, noticeAt sql.NullTime, iD int32) error
}

var _ Querier = (*Queries)(nil)
