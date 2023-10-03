// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const getTransaction = `-- name: GetTransaction :many
SELECT
	json_extract(transactions.data, '$.transaction.signatures[0]'),
	json_group_array(instructions.value)
FROM
  transactions, 
	json_each(json_extract(transactions.data, '$.transaction.message.instructions')) AS instructions
WHERE
	transactions.program_id = $1
	AND json_extract(transactions.data, '$.transaction.signatures[0]') > $2
	AND json_extract(json_extract(transactions.data, '$.transaction.message.accountKeys'), '$[' || json_extract(instructions.value, '$.programIdIndex') || ']') = transactions.program_id
GROUP BY transactions.id
LIMIT $3
`

type GetTransactionParams struct {
	ProgramID string
	Data      string
	Limit     int32
}

type GetTransactionRow struct {
	JsonExtract    interface{}
	JsonGroupArray interface{}
}

func (q *Queries) GetTransaction(ctx context.Context, arg GetTransactionParams, aq ...AdditionalQuery) ([]GetTransactionRow, error) {
	query := getTransaction
	queryParams := []interface{}{arg.ProgramID, arg.Data, arg.Limit}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.Query(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransactionRow
	for rows.Next() {
		var i GetTransactionRow
		if err := rows.Scan(&i.JsonExtract, &i.JsonGroupArray); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
