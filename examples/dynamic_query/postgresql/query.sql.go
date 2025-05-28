package postgresql

import (
	"context"
	"fmt" // Ensure fmt is imported
	"strings"
)

const getProducts = `-- name: GetProducts :many
SELECT id, name, category, price, is_available, created_at FROM products
WHERE 1=1
`

// GetProductsParams is a placeholder as the function takes optional params directly.
// It's not used by the generated GetProducts function itself but might be useful
// for users if they wanted to wrap the call.
type GetProductsParams struct {
	Category    interface{} `json:"category"`
	MinPrice    interface{} `json:"min_price"`
	IsAvailable interface{} `json:"is_available"`
}

func (q *Queries) GetProducts(ctx context.Context, category interface{}, minPrice interface{}, isAvailable interface{}) ([]Product, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(getProducts) // Base query

	var queryParams []interface{}

	// Optional 'Category'
	if category != nil {
		sqlBuilder.WriteString(" AND category = $")
		queryParams = append(queryParams, category)
		sqlBuilder.WriteString(fmt.Sprintf("%d", len(queryParams)))
	}

	// Optional 'MinPrice'
	if minPrice != nil {
		sqlBuilder.WriteString(" AND price >= $")
		queryParams = append(queryParams, minPrice)
		sqlBuilder.WriteString(fmt.Sprintf("%d", len(queryParams)))
	}

	// Optional 'IsAvailable'
	if isAvailable != nil {
		sqlBuilder.WriteString(" AND is_available = $")
		queryParams = append(queryParams, isAvailable)
		sqlBuilder.WriteString(fmt.Sprintf("%d", len(queryParams)))
	}

	rows, err := q.db.Query(ctx, sqlBuilder.String(), queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Category,
			&i.Price,
			&i.IsAvailable,
			&i.CreatedAt,
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

const addProduct = `-- name: AddProduct :one
INSERT INTO products (name, category, price, is_available)
VALUES ($1, $2, $3, $4)
RETURNING id, name, category, price, is_available, created_at
`

type AddProductParams struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Price       int32  `json:"price"`
	IsAvailable bool   `json:"is_available"`
}

func (q *Queries) AddProduct(ctx context.Context, arg AddProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, addProduct,
		arg.Name,
		arg.Category,
		arg.Price,
		arg.IsAvailable,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Category,
		&i.Price,
		&i.IsAvailable,
		&i.CreatedAt,
	)
	return i, err
}
