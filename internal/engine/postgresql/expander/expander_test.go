package expander

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestExpand(t *testing.T) {
	// Skip if no database connection available
	uri := os.Getenv("POSTGRESQL_SERVER_URI")
	if uri == "" {
		uri = "postgres://postgres:mysecretpassword@localhost:5432/postgres"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		t.Skipf("could not connect to database: %v", err)
	}
	defer pool.Close()

	// Create a test table
	_, err = pool.Exec(ctx, `
		DROP TABLE IF EXISTS authors;
		CREATE TABLE authors (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			bio TEXT
		);
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	defer pool.Exec(ctx, "DROP TABLE IF EXISTS authors")

	exp := New(pool)

	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "simple select star",
			query:    "SELECT * FROM authors",
			expected: "SELECT id, name, bio FROM authors",
		},
		{
			name:     "select with no star",
			query:    "SELECT id, name FROM authors",
			expected: "SELECT id, name FROM authors",
		},
		{
			name:     "select star with where clause",
			query:    "SELECT * FROM authors WHERE id = 1",
			expected: "SELECT id, name, bio FROM authors WHERE id = 1",
		},
		{
			name:     "double star",
			query:    "SELECT *, * FROM authors",
			expected: "SELECT id, name, bio, id, name, bio FROM authors",
		},
		{
			name:     "table qualified star",
			query:    "SELECT authors.* FROM authors",
			expected: "SELECT authors.id, authors.name, authors.bio FROM authors",
		},
		{
			name:     "star in middle of columns",
			query:    "SELECT id, *, name FROM authors",
			expected: "SELECT id, id, name, bio, name FROM authors",
		},
		{
			name:     "insert returning star",
			query:    "INSERT INTO authors (name, bio) VALUES ('John', 'A writer') RETURNING *",
			expected: "INSERT INTO authors (name, bio) VALUES ('John', 'A writer') RETURNING id, name, bio",
		},
		{
			name:     "insert returning mixed",
			query:    "INSERT INTO authors (name, bio) VALUES ('John', 'A writer') RETURNING id, *",
			expected: "INSERT INTO authors (name, bio) VALUES ('John', 'A writer') RETURNING id, id, name, bio",
		},
		{
			name:     "update returning star",
			query:    "UPDATE authors SET name = 'Jane' WHERE id = 1 RETURNING *",
			expected: "UPDATE authors SET name = 'Jane' WHERE id = 1 RETURNING id, name, bio",
		},
		{
			name:     "delete returning star",
			query:    "DELETE FROM authors WHERE id = 1 RETURNING *",
			expected: "DELETE FROM authors WHERE id = 1 RETURNING id, name, bio",
		},
		{
			name:     "cte with select star",
			query:    "WITH a AS (SELECT * FROM authors) SELECT * FROM a",
			expected: "WITH a AS (SELECT id, name, bio FROM authors) SELECT id, name, bio FROM a",
		},
		{
			name:     "multiple ctes with dependency",
			query:    "WITH a AS (SELECT * FROM authors), b AS (SELECT * FROM a) SELECT * FROM b",
			expected: "WITH a AS (SELECT id, name, bio FROM authors), b AS (SELECT id, name, bio FROM a) SELECT id, name, bio FROM b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := exp.Expand(ctx, tc.query)
			if err != nil {
				t.Fatalf("Expand failed: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
