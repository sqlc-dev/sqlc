// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const insertCode = `-- name: InsertCode :one
WITH cc AS (
            INSERT INTO td3.codes(created_by, updated_by, code, hash, is_private)
            VALUES ($1, $1, $2, $3, false)
            RETURNING hash
)
INSERT INTO td3.test_codes(created_by, updated_by, test_id, code_hash)
VALUES(
            $1, $1, $4, (select hash from cc)
)
RETURNING id, ts_created, ts_updated, created_by, updated_by, test_id, code_hash
`

type InsertCodeParams struct {
	CreatedBy string
	Code      sql.NullString
	Hash      sql.NullString
	TestID    int32
}

func (q *Queries) InsertCode(ctx context.Context, arg InsertCodeParams) (Td3TestCode, error) {
	row := q.db.QueryRow(ctx, insertCode,
		arg.CreatedBy,
		arg.Code,
		arg.Hash,
		arg.TestID,
	)
	var i Td3TestCode
	err := row.Scan(
		&i.ID,
		&i.TsCreated,
		&i.TsUpdated,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.TestID,
		&i.CodeHash,
	)
	return i, err
}
