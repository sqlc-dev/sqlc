package clickhouse

import (
	"strings"
	"testing"
)

// TestQueryBoundaryDetection tests that the parser correctly identifies
// statement boundaries including the -- name: annotation
func TestQueryBoundaryDetection(t *testing.T) {
	input := `-- name: QueryOne :one
SELECT id, name FROM table1
WHERE id = ?;

-- name: QueryTwo :many
SELECT id, value FROM table2
WHERE status = sqlc.arg('status')
ORDER BY id;

-- name: QueryThree :exec
INSERT INTO table3 (id, data)
VALUES (?, ?);`

	parser := NewParser()
	stmts, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if len(stmts) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(stmts))
	}

	// Extract the raw SQL for each statement
	type extractedQuery struct {
		name     string
		expected string
	}

	queries := []extractedQuery{
		{
			name: "QueryOne",
			expected: `-- name: QueryOne :one
SELECT id, name FROM table1
WHERE id = ?;`,
		},
		{
			name: "QueryTwo",
			expected: `-- name: QueryTwo :many
SELECT id, value FROM table2
WHERE status = sqlc.arg('status')
ORDER BY id;`,
		},
		{
			name: "QueryThree",
			expected: `-- name: QueryThree :exec
INSERT INTO table3 (id, data)
VALUES (?, ?);`,
		},
	}

	for i, stmt := range stmts {
		raw := stmt.Raw
		if raw == nil {
			t.Fatalf("statement %d has no RawStmt", i)
		}

		// Extract the SQL text using the same logic as the compiler
		location := raw.StmtLocation
		length := raw.StmtLen
		extracted := input[location : location+length]

		// Normalize whitespace for comparison
		extracted = strings.TrimSpace(extracted)
		expected := strings.TrimSpace(queries[i].expected)

		if extracted != expected {
			t.Errorf("Query %d (%s) boundary mismatch:\nExpected:\n%s\n\nGot:\n%s",
				i, queries[i].name, expected, extracted)
		}
	}
}

// TestComplexQueryBoundaries tests boundary detection with more complex queries
func TestComplexQueryBoundaries(t *testing.T) {
	input := `-- name: GetUserByID :one
SELECT id, name, email, created_at
FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT id, name, email, created_at
FROM users
ORDER BY created_at DESC
LIMIT ?;

-- name: InsertUser :exec
INSERT INTO users (id, name, email, created_at)
VALUES (?, ?, ?, ?);`

	parser := NewParser()
	stmts, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if len(stmts) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(stmts))
	}

	// Verify each extracted query
	for i, stmt := range stmts {
		raw := stmt.Raw
		location := raw.StmtLocation
		length := raw.StmtLen
		extracted := input[location : location+length]

		// Check that it starts with "-- name:"
		if !strings.HasPrefix(strings.TrimSpace(extracted), "-- name:") {
			t.Errorf("Query %d doesn't start with '-- name:': %q", i, extracted[:min(50, len(extracted))])
		}

		// Check that it ends with a semicolon
		trimmed := strings.TrimSpace(extracted)
		if !strings.HasSuffix(trimmed, ";") {
			t.Errorf("Query %d doesn't end with ';': %q", i, trimmed[max(0, len(trimmed)-50):])
		}

		// Check that it doesn't contain text from adjacent queries
		lines := strings.Split(extracted, "\n")
		nameCommentCount := 0
		for _, line := range lines {
			if strings.Contains(line, "-- name:") {
				nameCommentCount++
			}
		}
		if nameCommentCount != 1 {
			t.Errorf("Query %d contains %d '-- name:' comments, expected 1:\n%s",
				i, nameCommentCount, extracted)
		}
	}
}
