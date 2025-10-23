//go:build examples

package skipparser

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
)

func TestSkipParser(t *testing.T) {
	ctx := context.Background()
	uri := local.PostgreSQL(t, []string{"schema.sql"})
	db, err := pgx.Connect(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)

	q := New(db)

	// Test CountProducts on empty database
	count, err := q.CountProducts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0 products, got %d", count)
	}

	// Test CreateProduct
	product, err := q.CreateProduct(ctx, CreateProductParams{
		Name:  "Test Product",
		Price: "99.99",
		Tags:  []string{"electronics", "test"},
		Metadata: []byte(`{"color": "blue", "weight": 1.5}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if product.ID == 0 {
		t.Error("expected product ID to be non-zero")
	}
	if product.Name != "Test Product" {
		t.Errorf("expected name 'Test Product', got %s", product.Name)
	}
	t.Logf("Created product: %+v", product)

	// Test GetProduct
	fetchedProduct, err := q.GetProduct(ctx, product.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetchedProduct.ID != product.ID {
		t.Errorf("expected ID %d, got %d", product.ID, fetchedProduct.ID)
	}
	t.Logf("Fetched product: %+v", fetchedProduct)

	// Test ListProducts
	products, err := q.ListProducts(ctx, ListProductsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 1 {
		t.Errorf("expected 1 product, got %d", len(products))
	}
	t.Logf("Listed products: %+v", products)

	// Test UpdateProduct
	updatedProduct, err := q.UpdateProduct(ctx, UpdateProductParams{
		ID:    product.ID,
		Name:  "Updated Product",
		Price: "149.99",
		Tags:  []string{"electronics", "updated"},
		Metadata: []byte(`{"color": "red", "weight": 2.0}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedProduct.Name != "Updated Product" {
		t.Errorf("expected name 'Updated Product', got %s", updatedProduct.Name)
	}
	t.Logf("Updated product: %+v", updatedProduct)

	// Test SearchProductsByTag
	tagProducts, err := q.SearchProductsByTag(ctx, "electronics")
	if err != nil {
		t.Fatal(err)
	}
	if len(tagProducts) != 1 {
		t.Errorf("expected 1 product with tag 'electronics', got %d", len(tagProducts))
	}
	t.Logf("Products with tag 'electronics': %+v", tagProducts)

	// Test CountProducts after insert
	count, err = q.CountProducts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 product, got %d", count)
	}

	// Test DeleteProduct
	err = q.DeleteProduct(ctx, product.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Verify deletion
	count, err = q.CountProducts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0 products after deletion, got %d", count)
	}
}

func TestSkipParserComplexTypes(t *testing.T) {
	ctx := context.Background()
	uri := local.PostgreSQL(t, []string{"schema.sql"})
	db, err := pgx.Connect(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)

	q := New(db)

	// Test with empty arrays and JSON
	product, err := q.CreateProduct(ctx, CreateProductParams{
		Name:     "Minimal Product",
		Price:    "19.99",
		Tags:     []string{},
		Metadata: []byte(`{}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created minimal product: %+v", product)

	// Test with nil values where allowed (using pgtype for nullable fields)
	product2, err := q.CreateProduct(ctx, CreateProductParams{
		Name:     "Another Product",
		Price:    "29.99",
		Tags:     nil,
		Metadata: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created product with nil arrays: %+v", product2)

	// Cleanup
	_ = q.DeleteProduct(ctx, product.ID)
	_ = q.DeleteProduct(ctx, product2.ID)
}
