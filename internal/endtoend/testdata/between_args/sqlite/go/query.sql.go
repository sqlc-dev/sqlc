// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
)

const getBetweenPrices = `-- name: GetBetweenPrices :many
SELECT  name, price
FROM    products
WHERE   price BETWEEN ? AND ?
`

type GetBetweenPricesParams struct {
	FromPrice int64
	ToPrice   int64
}

func (q *Queries) GetBetweenPrices(ctx context.Context, arg GetBetweenPricesParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getBetweenPrices, arg.FromPrice, arg.ToPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(&i.Name, &i.Price); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBetweenPricesTable = `-- name: GetBetweenPricesTable :many
SELECT  name, price
FROM    products
WHERE   products.price BETWEEN ? AND ?
`

type GetBetweenPricesTableParams struct {
	FromPrice int64
	ToPrice   int64
}

func (q *Queries) GetBetweenPricesTable(ctx context.Context, arg GetBetweenPricesTableParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getBetweenPricesTable, arg.FromPrice, arg.ToPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(&i.Name, &i.Price); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBetweenPricesTableAlias = `-- name: GetBetweenPricesTableAlias :many
SELECT  name, price
FROM    products as p
WHERE   p.price BETWEEN ? AND ?
`

type GetBetweenPricesTableAliasParams struct {
	FromPrice int64
	ToPrice   int64
}

func (q *Queries) GetBetweenPricesTableAlias(ctx context.Context, arg GetBetweenPricesTableAliasParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getBetweenPricesTableAlias, arg.FromPrice, arg.ToPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(&i.Name, &i.Price); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
