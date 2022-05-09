// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: inventory.sql

package split

import (
	"context"

	"github.com/robinhooodmarkets/rh/storage-tools/playground/sharded-showoff/models"
)

const getStock = `-- name: GetStock :many
SELECT s_i_id, s_w_id, s_quantity, s_dist_01, s_dist_02, s_dist_03, s_dist_04, s_dist_05, s_dist_06, s_dist_07, s_dist_08, s_dist_09, s_dist_10, s_ytd, s_order_cnt, s_remote_cnt, s_data FROM stock
LIMIT 10
`

func (q *Queries) GetStock(ctx context.Context) ([]models.Stock, error) {
	rows, err := q.db.Query(ctx, getStock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Stock
	for rows.Next() {
		var i models.Stock
		if err := rows.Scan(
			&i.SIID,
			&i.SWID,
			&i.SQuantity,
			&i.SDist01,
			&i.SDist02,
			&i.SDist03,
			&i.SDist04,
			&i.SDist05,
			&i.SDist06,
			&i.SDist07,
			&i.SDist08,
			&i.SDist09,
			&i.SDist10,
			&i.SYtd,
			&i.SOrderCnt,
			&i.SRemoteCnt,
			&i.SData,
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
