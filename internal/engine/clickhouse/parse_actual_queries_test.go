package clickhouse

import (
	"strings"
	"testing"
)

// TestActualQueryBoundaries tests with the actual queries that are having issues
func TestActualQueryBoundaries(t *testing.T) {
	// These are the actual queries from examples/clickhouse/queries.sql that show the bug
	input := `-- name: UnfoldNestedData :many
SELECT
	record_id,
	nested_value
FROM sqlc_example.nested_table
ARRAY JOIN nested_array AS nested_value
WHERE record_id IN (sqlc.slice('record_ids'));

-- name: AnalyzeArrayElements :many
SELECT
	product_id,
	arrayJoin(categories) AS category,
	COUNT(*) OVER (PARTITION BY category) as category_count
FROM sqlc_example.products
WHERE product_id = ?
GROUP BY product_id, category;

-- name: ExtractMetadataFromJSON :many
SELECT
	MetadataPlatformId,
	arrayJoin(JSONExtract(JsonValue, 'Array(String)')) as self_help_id
FROM sqlc_example.events;`

	parser := NewParser()
	stmts, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if len(stmts) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(stmts))
	}

	// Define expected queries
	expectedQueries := []string{
		`-- name: UnfoldNestedData :many
SELECT
	record_id,
	nested_value
FROM sqlc_example.nested_table
ARRAY JOIN nested_array AS nested_value
WHERE record_id IN (sqlc.slice('record_ids'));`,
		`-- name: AnalyzeArrayElements :many
SELECT
	product_id,
	arrayJoin(categories) AS category,
	COUNT(*) OVER (PARTITION BY category) as category_count
FROM sqlc_example.products
WHERE product_id = ?
GROUP BY product_id, category;`,
		`-- name: ExtractMetadataFromJSON :many
SELECT
	MetadataPlatformId,
	arrayJoin(JSONExtract(JsonValue, 'Array(String)')) as self_help_id
FROM sqlc_example.events;`,
	}

	for i, stmt := range stmts {
		raw := stmt.Raw
		if raw == nil {
			t.Fatalf("statement %d has no RawStmt", i)
		}

		// Extract the SQL text
		location := raw.StmtLocation
		length := raw.StmtLen

		t.Logf("Statement %d: location=%d, length=%d", i, location, length)

		if location < 0 || location >= len(input) {
			t.Errorf("Statement %d: invalid location %d (input length: %d)", i, location, len(input))
			continue
		}

		if location+length > len(input) {
			t.Errorf("Statement %d: location+length (%d) exceeds input length (%d)",
				i, location+length, len(input))
			continue
		}

		extracted := input[location : location+length]

		// Normalize whitespace for comparison
		extracted = strings.TrimSpace(extracted)
		expected := strings.TrimSpace(expectedQueries[i])

		if extracted != expected {
			t.Errorf("Query %d boundary mismatch:\n\n=== EXPECTED ===\n%s\n\n=== GOT ===\n%s\n\n=== DIFF ===",
				i, expected, extracted)

			// Show first difference
			minLen := len(extracted)
			if len(expected) < minLen {
				minLen = len(expected)
			}
			for j := 0; j < minLen; j++ {
				if extracted[j] != expected[j] {
					start := j - 20
					if start < 0 {
						start = 0
					}
					end := j + 20
					if end > minLen {
						end = minLen
					}
					t.Errorf("First difference at position %d:\nExpected: %q\nGot: %q",
						j, expected[start:end], extracted[start:end])
					break
				}
			}
		}
	}
}
