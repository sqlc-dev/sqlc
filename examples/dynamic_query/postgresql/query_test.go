package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx stdlib driver
)

var testQueries *Queries

// TestMain sets up the database connection and runs the tests.
// For a real test suite, you'd use a dedicated test database and potentially migrations.
func TestMain(m *testing.M) {
	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Provide a default for local testing if DATABASE_URL is not set.
		// Adjust if your local PostgreSQL instance is different.
		dbURL = "postgres://user:password@localhost:5432/testdb?sslmode=disable"
		log.Printf("DATABASE_URL not set, using default: %s\n", dbURL)
	}

	// For pgx/v5, it's common to use pgxpool directly.
	// However, sqlc can also generate code for database/sql, which might be simpler for some examples.
	// The sqlc.yaml specified pgx/v5, so we'll use pgxpool here.
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Or, if using database/sql with pgx (e.g., if sql_package was "database/sql" and driver was "pgx")
	// stdDB, err := sql.Open("pgx", dbURL)
	// if err != nil {
	// 	log.Fatalf("Unable to connect to database using database/sql: %v\n", err)
	// }
	// defer stdDB.Close()
	// testQueries = New(stdDB) // If using database/sql adapter

	testQueries = New(pool) // New takes DBTX, which pgxpool.Pool implements

	// Minimal schema setup - in a real test, use migrations.
	// This is a simplified version and might fail if the table already exists.
	// Consider dropping and recreating for idempotency in real tests.
	_, err = pool.Exec(ctx, `
		DROP TABLE IF EXISTS products;
		CREATE TABLE products (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			price INT NOT NULL,
			is_available BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create schema: %v\n", err)
	}

	// Insert some initial data
	initialProducts := []struct {
		Name        string
		Category    string
		Price       int32
		IsAvailable bool
	}{
		{"Laptop Pro", "electronics", 1200, true},
		{"Coffee Maker", "appliances", 80, true},
		{"Gaming Mouse", "electronics", 75, true},
		{"Desk Chair", "furniture", 150, false},
		{"Laptop Lite", "electronics", 800, true},
	}

	for _, p := range initialProducts {
		_, err := testQueries.AddProduct(ctx, AddProductParams{
			Name:        p.Name,
			Category:    p.Category,
			Price:       p.Price,
			IsAvailable: p.IsAvailable,
		})
		if err != nil {
			log.Fatalf("Failed to insert initial product %s: %v", p.Name, err)
		}
	}

	log.Println("Test database setup complete.")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetProducts(t *testing.T) {
	ctx := context.Background()

	if testQueries == nil {
		t.Fatal("testQueries not initialized. DB setup might have failed.")
	}

	// Example 1: Get all products (all optional params nil)
	t.Run("GetAllProducts", func(t *testing.T) {
		products, err := testQueries.GetProducts(ctx, nil, nil, nil)
		if err != nil {
			t.Fatalf("Failed to get all products: %v", err)
		}
		if len(products) < 5 { // Based on initial data
			t.Errorf("Expected at least 5 products, got %d", len(products))
		}
		t.Logf("All products count: %d", len(products))
		// for _, p := range products {
		// 	t.Logf("  Product: ID=%d, Name=%s, Category=%s, Price=%d, Available=%t", p.ID, p.Name, p.Category, p.Price, p.IsAvailable)
		// }
	})

	// Example 2: Get products in 'electronics' category
	t.Run("GetElectronicsProducts", func(t *testing.T) {
		electronicsCategory := "electronics"
		products, err := testQueries.GetProducts(ctx, &electronicsCategory, nil, nil)
		if err != nil {
			t.Fatalf("Failed to get electronics products: %v", err)
		}
		if len(products) < 3 { // Laptop Pro, Gaming Mouse, Laptop Lite
			t.Errorf("Expected at least 3 electronics products, got %d. Products: %+v", len(products), products)
		}
		for _, p := range products {
			if p.Category != electronicsCategory {
				t.Errorf("Expected category %s, got %s for product %s", electronicsCategory, p.Category, p.Name)
			}
		}
		t.Logf("Electronics products count: %d", len(products))
	})

	// Example 3: Get 'electronics' products with minPrice 1000
	t.Run("GetElectronicsMinPrice1000", func(t *testing.T) {
		electronicsCategory := "electronics"
		minPrice := int32(1000)
		products, err := testQueries.GetProducts(ctx, &electronicsCategory, &minPrice, nil)
		if err != nil {
			t.Fatalf("Failed to get electronics products >= 1000: %v", err)
		}
		if len(products) < 1 { // Laptop Pro
			t.Errorf("Expected at least 1 electronics product >= 1000, got %d. Products: %+v", len(products), products)
		}
		for _, p := range products {
			if p.Category != electronicsCategory || p.Price < minPrice {
				t.Errorf("Product %s (Cat: %s, Price: %d) does not match filters (Cat: %s, MinPrice: %d)",
					p.Name, p.Category, p.Price, electronicsCategory, minPrice)
			}
		}
		t.Logf("Electronics products >= 1000 count: %d", len(products))
	})

	// Example 4: Get available 'electronics' products with minPrice 1000
	t.Run("GetAvailableElectronicsMinPrice1000", func(t *testing.T) {
		electronicsCategory := "electronics"
		minPrice := int32(1000)
		isAvailable := true
		products, err := testQueries.GetProducts(ctx, &electronicsCategory, &minPrice, &isAvailable)
		if err != nil {
			t.Fatalf("Failed to get available electronics products >= 1000: %v", err)
		}
		if len(products) < 1 { // Laptop Pro
			t.Errorf("Expected at least 1 available electronics product >= 1000, got %d. Products: %+v", len(products), products)
		}
		for _, p := range products {
			if p.Category != electronicsCategory || p.Price < minPrice || !p.IsAvailable {
				t.Errorf("Product %s (Cat: %s, Price: %d, Avail: %t) does not match filters (Cat: %s, MinPrice: %d, Avail: %t)",
					p.Name, p.Category, p.Price, p.IsAvailable, electronicsCategory, minPrice, isAvailable)
			}
		}
		t.Logf("Available electronics products >= 1000 count: %d", len(products))
	})

	// Example 5: Get unavailable products (isAvailable = false)
	t.Run("GetUnavailableProducts", func(t *testing.T) {
		isAvailable := false
		products, err := testQueries.GetProducts(ctx, nil, nil, &isAvailable)
		if err != nil {
			t.Fatalf("Failed to get unavailable products: %v", err)
		}
		if len(products) < 1 { // Desk Chair
			t.Errorf("Expected at least 1 unavailable product, got %d. Products: %+v", len(products), products)
		}
		for _, p := range products {
			if p.IsAvailable != isAvailable {
				t.Errorf("Expected isAvailable %t, got %t for product %s", isAvailable, p.IsAvailable, p.Name)
			}
		}
		t.Logf("Unavailable products count: %d", len(products))
	})

	fmt.Println("TestGetProducts complete.")
}

// Example usage of AddProduct (not a test of GetProducts, but good for completeness)
func TestAddProduct(t *testing.T) {
	ctx := context.Background()
	if testQueries == nil {
		t.Fatal("testQueries not initialized.")
	}

	newProductParams := AddProductParams{
		Name:        "Test Book",
		Category:    "books",
		Price:       25,
		IsAvailable: true,
	}
	product, err := testQueries.AddProduct(ctx, newProductParams)
	if err != nil {
		t.Fatalf("AddProduct failed: %v", err)
	}
	if product.Name != newProductParams.Name {
		t.Errorf("Expected product name %s, got %s", newProductParams.Name, product.Name)
	}
	t.Logf("Added product: ID=%d, Name=%s", product.ID, product.Name)
}
