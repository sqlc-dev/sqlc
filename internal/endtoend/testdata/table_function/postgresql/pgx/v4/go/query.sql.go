// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgtype"
)

const getTransaction = `-- name: GetTransaction :many
SELECT
	jsonb_extract_path(transactions.data, '$.transaction.signatures[0]'),
	jsonb_agg(instructions.value)
FROM
  transactions, 
	jsonb_each(jsonb_extract_path(transactions.data, '$.transaction.message.instructions[0]')) AS instructions
WHERE
	transactions.program_id = $1
	AND jsonb_extract_path(transactions.data, '$.transaction.signatures[0]') @> to_jsonb($2::text)
	AND jsonb_extract_path(jsonb_extract_path(transactions.data, '$.transaction.message.accountKeys'), 'key') = to_jsonb(transactions.program_id)
GROUP BY transactions.id
`

type GetTransactionParams struct {
	ProgramID string
	Data      string
}

type GetTransactionRow struct {
	JsonbExtractPath pgtype.JSONB
	JsonbAgg         pgtype.JSONB
}

func (q *Queries) GetTransaction(ctx context.Context, arg GetTransactionParams) ([]GetTransactionRow, error) {
	rows, err := q.db.Query(ctx, getTransaction, arg.ProgramID, arg.Data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransactionRow
	for rows.Next() {
		var i GetTransactionRow
		if err := rows.Scan(&i.JsonbExtractPath, &i.JsonbAgg); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
