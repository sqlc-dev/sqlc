package postgresql

import (
	"context"
)

type Querier interface {
	GetProducts(ctx context.Context, category interface{}, minPrice interface{}, isAvailable interface{}) ([]Product, error)
	AddProduct(ctx context.Context, arg AddProductParams) (Product, error)
}

var _ Querier = (*Queries)(nil)
