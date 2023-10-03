// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const selectTest = `-- name: SelectTest :many
SELECT v_daterange_null, v_datemultirange_null, v_tsrange_null, v_tsmultirange_null, v_tstzrange_null, v_tstzmultirange_null, v_numrange_null, v_nummultirange_null, v_int4range_null, v_int4multirange_null, v_int8range_null, v_int8multirange_null, v_daterange, v_datemultirange, v_tsrange, v_tsmultirange, v_tstzrange, v_tstzmultirange, v_numrange, v_nummultirange, v_int4range, v_int4multirange, v_int8range, v_int8multirange from test_table
`

func (q *Queries) SelectTest(ctx context.Context, aq ...AdditionalQuery) ([]TestTable, error) {
	query := selectTest
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.Query(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TestTable
	for rows.Next() {
		var i TestTable
		if err := rows.Scan(
			&i.VDaterangeNull,
			&i.VDatemultirangeNull,
			&i.VTsrangeNull,
			&i.VTsmultirangeNull,
			&i.VTstzrangeNull,
			&i.VTstzmultirangeNull,
			&i.VNumrangeNull,
			&i.VNummultirangeNull,
			&i.VInt4rangeNull,
			&i.VInt4multirangeNull,
			&i.VInt8rangeNull,
			&i.VInt8multirangeNull,
			&i.VDaterange,
			&i.VDatemultirange,
			&i.VTsrange,
			&i.VTsmultirange,
			&i.VTstzrange,
			&i.VTstzmultirange,
			&i.VNumrange,
			&i.VNummultirange,
			&i.VInt4range,
			&i.VInt4multirange,
			&i.VInt8range,
			&i.VInt8multirange,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
