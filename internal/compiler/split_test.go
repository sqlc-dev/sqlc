package compiler

import (
	"context"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/sqlfile"
)

func TestSqlfileSplitForSkipParser(t *testing.T) {
	input := `
-- name: GetUser :one
SELECT id, name FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT id, name FROM users ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (name) VALUES ($1) RETURNING id, name;
`

	ctx := context.Background()
	queries, err := sqlfile.Split(ctx, strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to split queries: %s", err)
	}

	if len(queries) != 3 {
		t.Fatalf("expected 3 queries, got %d", len(queries))
	}

	// Check first query
	if !strings.Contains(queries[0], "GetUser") {
		t.Errorf("first query should contain GetUser, got: %s", queries[0])
	}
	if !strings.Contains(queries[0], "WHERE id = $1") {
		t.Errorf("first query should contain WHERE clause, got: %s", queries[0])
	}

	// Check second query
	if !strings.Contains(queries[1], "ListUsers") {
		t.Errorf("second query should contain ListUsers, got: %s", queries[1])
	}
	if !strings.Contains(queries[1], "ORDER BY id") {
		t.Errorf("second query should contain ORDER BY, got: %s", queries[1])
	}

	// Check third query
	if !strings.Contains(queries[2], "CreateUser") {
		t.Errorf("third query should contain CreateUser, got: %s", queries[2])
	}
	if !strings.Contains(queries[2], "RETURNING") {
		t.Errorf("third query should contain RETURNING, got: %s", queries[2])
	}
}

func TestSqlfileSplitWithComplexQueries(t *testing.T) {
	input := `
-- name: ComplexQuery :many
SELECT
    id,
    name,
    CASE
        WHEN price > 100 THEN 'expensive'
        ELSE 'affordable'
    END as category,
    tags
FROM products
WHERE
    name LIKE '%' || $1 || '%'
    AND $2 = ANY(tags)
ORDER BY created_at DESC
LIMIT $3;

-- name: UpdateWithJSON :one
UPDATE products
SET metadata = $2::jsonb
WHERE id = $1
RETURNING id, metadata;
`

	ctx := context.Background()
	queries, err := sqlfile.Split(ctx, strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to split complex queries: %s", err)
	}

	if len(queries) != 2 {
		t.Fatalf("expected 2 queries, got %d", len(queries))
	}

	// Check first complex query
	if !strings.Contains(queries[0], "CASE") {
		t.Errorf("first query should contain CASE statement, got: %s", queries[0])
	}
	if !strings.Contains(queries[0], "ANY(tags)") {
		t.Errorf("first query should contain ANY operator, got: %s", queries[0])
	}

	// Check second query with JSON
	if !strings.Contains(queries[1], "::jsonb") {
		t.Errorf("second query should contain jsonb cast, got: %s", queries[1])
	}
	if !strings.Contains(queries[1], "RETURNING") {
		t.Errorf("second query should contain RETURNING, got: %s", queries[1])
	}
}

func TestSqlfileSplitWithDollarQuotes(t *testing.T) {
	input := `
-- name: CreateFunction :exec
CREATE OR REPLACE FUNCTION calculate_total(p_price NUMERIC, p_quantity INTEGER)
RETURNS NUMERIC AS $$
BEGIN
    RETURN p_price * p_quantity;
END;
$$ LANGUAGE plpgsql;

-- name: GetValue :one
SELECT $$This is a dollar quoted string$$ as value;
`

	ctx := context.Background()
	queries, err := sqlfile.Split(ctx, strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to split queries with dollar quotes: %s", err)
	}

	if len(queries) != 2 {
		t.Fatalf("expected 2 queries, got %d", len(queries))
	}

	// Check function creation with dollar quotes
	if !strings.Contains(queries[0], "$$") {
		t.Errorf("first query should contain dollar quotes, got: %s", queries[0])
	}
	if !strings.Contains(queries[0], "plpgsql") {
		t.Errorf("first query should contain plpgsql, got: %s", queries[0])
	}

	// Check dollar quoted string
	if !strings.Contains(queries[1], "dollar quoted string") {
		t.Errorf("second query should contain dollar quoted string, got: %s", queries[1])
	}
}

func TestSqlfileSplitEmptyAndComments(t *testing.T) {
	input := `
-- name: OnlyQuery :one
SELECT 1;
`

	ctx := context.Background()
	queries, err := sqlfile.Split(ctx, strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to split queries: %s", err)
	}

	// Should get at least one query
	if len(queries) < 1 {
		t.Fatalf("expected at least 1 query, got %d", len(queries))
	}

	// Find the query with our name
	found := false
	for _, q := range queries {
		if strings.Contains(q, "OnlyQuery") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("should find query with OnlyQuery")
	}
}
