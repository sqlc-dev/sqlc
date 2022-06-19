// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: pgcrypto.sql

package querytest

import (
	"context"
)

const encodeDigest = `-- name: EncodeDigest :one
SELECT encode(digest($1, 'sha1'), 'hex')
`

func (q *Queries) EncodeDigest(ctx context.Context, digest string) (string, error) {
	ctx, done := q.observer(ctx, "EncodeDigest")
	row := q.db.QueryRowContext(ctx, encodeDigest, digest)
	var encode string
	err := row.Scan(&encode)
	return encode, done(err)
}
