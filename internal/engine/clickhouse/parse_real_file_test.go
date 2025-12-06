package clickhouse

import (
	"os"
	"strings"
	"testing"
)

// TestRealQueryFile tests parsing the actual queries.sql file
func TestRealQueryFile(t *testing.T) {
	// Read the actual queries file
	queriesPath := "../../../examples/clickhouse/queries.sql"
	content, err := os.ReadFile(queriesPath)
	if err != nil {
		t.Skipf("Could not read queries file: %v", err)
	}

	input := string(content)

	parser := NewParser()
	stmts, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	t.Logf("Parsed %d statements", len(stmts))

	// Find statements we know are having issues
	problemQueries := map[string]struct {
		expectedStart string
		expectedEnd   string
	}{
		"UnfoldNestedData": {
			expectedStart: "-- name: UnfoldNestedData :many",
			expectedEnd:   "WHERE record_id IN (sqlc.slice('record_ids'));",
		},
		"AnalyzeArrayElements": {
			expectedStart: "-- name: AnalyzeArrayElements :many",
			expectedEnd:   "GROUP BY product_id, category;",
		},
		"ExtractMetadataFromJSON": {
			expectedStart: "-- name: ExtractMetadataFromJSON :many",
			expectedEnd:   "FROM sqlc_example.events;",
		},
	}

	// Check each statement
	for i, stmt := range stmts {
		raw := stmt.Raw
		if raw == nil {
			continue
		}

		location := raw.StmtLocation
		length := raw.StmtLen

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
		extracted = strings.TrimSpace(extracted)

		// Look for the query name in the extracted text
		if strings.Contains(extracted, "-- name:") {
			// Extract the query name
			lines := strings.Split(extracted, "\n")
			if len(lines) > 0 {
				firstLine := strings.TrimSpace(lines[0])
				if strings.HasPrefix(firstLine, "-- name:") {
					parts := strings.Fields(firstLine)
					if len(parts) >= 3 {
						queryName := parts[2]

						if check, ok := problemQueries[queryName]; ok {
							// Check if it starts correctly
							if !strings.HasPrefix(extracted, check.expectedStart) {
								t.Errorf("Query %s: doesn't start correctly\nExpected start: %q\nGot: %q",
									queryName, check.expectedStart, extracted[:min(len(check.expectedStart)+20, len(extracted))])
							}

							// Check if it ends correctly
							if !strings.HasSuffix(extracted, check.expectedEnd) {
								t.Errorf("Query %s: doesn't end correctly\nExpected end: %q\nGot: %q",
									queryName, check.expectedEnd, extracted[max(0, len(extracted)-len(check.expectedEnd)-20):])
							}

							// Check for contamination from other queries
							nameCommentCount := 0
							for _, line := range lines {
								if strings.Contains(line, "-- name:") {
									nameCommentCount++
								}
							}
							if nameCommentCount > 1 {
								t.Errorf("Query %s contains %d '-- name:' comments (expected 1)",
									queryName, nameCommentCount)
							}
						}
					}
				}
			}
		}
	}
}
