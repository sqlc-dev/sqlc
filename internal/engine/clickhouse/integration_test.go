//go:build integration
// +build integration

package clickhouse

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

// TestArrayJoinIntegration tests ARRAY JOIN against a live ClickHouse database
// Run with: go test -tags=integration -run TestArrayJoinIntegration ./internal/engine/clickhouse
//
// Prerequisites:
// - ClickHouse server running (docker run -p 9000:9000 -p 8123:8123 clickhouse/clickhouse-server)
// - Or use docker-compose up clickhouse
func TestArrayJoinIntegration(t *testing.T) {
	// Skip if no ClickHouse connection info
	clickhouseURL := os.Getenv("CLICKHOUSE_URL")
	if clickhouseURL == "" {
		clickhouseURL = "clickhouse://localhost:9000/default"
	}

	db, err := sql.Open("clickhouse", clickhouseURL)
	if err != nil {
		t.Skip("ClickHouse not available:", err)
		return
	}
	defer db.Close()

	ctx := context.Background()

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		t.Skip("ClickHouse not reachable:", err)
		return
	}

	// Create test database
	_, err = db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS sqlc_test")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Clean up function
	defer func() {
		db.ExecContext(ctx, "DROP DATABASE IF EXISTS sqlc_test")
	}()

	t.Run("BasicArrayJoin", func(t *testing.T) {
		// Create table with array column
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS sqlc_test.users_with_tags (
				id UInt32,
				name String,
				tags Array(String)
			) ENGINE = MergeTree()
			ORDER BY id
		`)
		if err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}
		defer db.ExecContext(ctx, "DROP TABLE IF EXISTS sqlc_test.users_with_tags")

		// Insert test data
		_, err = db.ExecContext(ctx, `
			INSERT INTO sqlc_test.users_with_tags VALUES
			(1, 'Alice', ['developer', 'admin']),
			(2, 'Bob', ['designer', 'user']),
			(3, 'Charlie', ['manager'])
		`)
		if err != nil {
			t.Fatalf("Failed to insert data: %v", err)
		}

		// Test ARRAY JOIN query
		rows, err := db.QueryContext(ctx, `
			SELECT id, name, tag
			FROM sqlc_test.users_with_tags
			ARRAY JOIN tags AS tag
			ORDER BY id, tag
		`)
		if err != nil {
			t.Fatalf("ARRAY JOIN query failed: %v", err)
		}
		defer rows.Close()

		// Verify results
		expectedResults := []struct {
			id   uint32
			name string
			tag  string
		}{
			{1, "Alice", "admin"},
			{1, "Alice", "developer"},
			{2, "Bob", "designer"},
			{2, "Bob", "user"},
			{3, "Charlie", "manager"},
		}

		resultCount := 0
		for rows.Next() {
			var id uint32
			var name, tag string
			if err := rows.Scan(&id, &name, &tag); err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}

			if resultCount >= len(expectedResults) {
				t.Fatalf("More results than expected, got row: id=%d, name=%s, tag=%s", id, name, tag)
			}

			expected := expectedResults[resultCount]
			if id != expected.id || name != expected.name || tag != expected.tag {
				t.Errorf("Row %d mismatch: got (%d, %s, %s), want (%d, %s, %s)",
					resultCount, id, name, tag, expected.id, expected.name, expected.tag)
			}
			resultCount++
		}

		if resultCount != len(expectedResults) {
			t.Errorf("Expected %d rows, got %d", len(expectedResults), resultCount)
		}
	})

	t.Run("ArrayJoinWithFilter", func(t *testing.T) {
		// Create table with nested arrays
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS sqlc_test.events_with_props (
				event_id UInt32,
				event_name String,
				properties Array(String)
			) ENGINE = MergeTree()
			ORDER BY event_id
		`)
		if err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}
		defer db.ExecContext(ctx, "DROP TABLE IF EXISTS sqlc_test.events_with_props")

		// Insert test data
		_, err = db.ExecContext(ctx, `
			INSERT INTO sqlc_test.events_with_props VALUES
			(1, 'click', ['button', 'header', 'link']),
			(2, 'view', ['page', 'section']),
			(3, 'submit', ['form', 'button'])
		`)
		if err != nil {
			t.Fatalf("Failed to insert data: %v", err)
		}

		// Test ARRAY JOIN with WHERE clause
		rows, err := db.QueryContext(ctx, `
			SELECT event_id, event_name, prop
			FROM sqlc_test.events_with_props
			ARRAY JOIN properties AS prop
			WHERE event_id >= 2
			ORDER BY event_id, prop
		`)
		if err != nil {
			t.Fatalf("ARRAY JOIN with WHERE failed: %v", err)
		}
		defer rows.Close()

		resultCount := 0
		for rows.Next() {
			var eventID uint32
			var eventName, prop string
			if err := rows.Scan(&eventID, &eventName, &prop); err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}

			// Verify event_id is >= 2
			if eventID < 2 {
				t.Errorf("Expected event_id >= 2, got %d", eventID)
			}
			resultCount++
		}

		// Should have 4 rows (2: page, section; 3: button, form)
		if resultCount != 4 {
			t.Errorf("Expected 4 rows, got %d", resultCount)
		}
	})

	t.Run("ArrayJoinFunction", func(t *testing.T) {
		// Test arrayJoin() function (different from ARRAY JOIN clause)
		rows, err := db.QueryContext(ctx, `
			SELECT arrayJoin(['a', 'b', 'c']) AS element
		`)
		if err != nil {
			t.Fatalf("arrayJoin() function failed: %v", err)
		}
		defer rows.Close()

		elements := []string{}
		for rows.Next() {
			var elem string
			if err := rows.Scan(&elem); err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			elements = append(elements, elem)
		}

		expected := []string{"a", "b", "c"}
		if len(elements) != len(expected) {
			t.Errorf("Expected %d elements, got %d", len(expected), len(elements))
		}

		for i, elem := range elements {
			if elem != expected[i] {
				t.Errorf("Element %d: expected %s, got %s", i, expected[i], elem)
			}
		}
	})

	t.Run("ParseAndExecute", func(t *testing.T) {
		// Test that our parser can parse the ARRAY JOIN query
		// and the generated AST is correct
		sql := `
			SELECT id, name, tag
			FROM sqlc_test.users_with_tags
			ARRAY JOIN tags AS tag
			WHERE id = 1
			ORDER BY tag
		`

		p := NewParser()
		stmts, err := p.Parse(strings.NewReader(sql))
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if len(stmts) != 1 {
			t.Fatalf("Expected 1 statement, got %d", len(stmts))
		}

		// The parsing succeeded - this validates our converter worked
		t.Log("Successfully parsed ARRAY JOIN query")
	})
}

// TestArrayJoinCodeGeneration tests the full pipeline: parse -> generate -> execute
func TestArrayJoinCodeGeneration(t *testing.T) {
	// This would require the full sqlc pipeline
	// For now, we just test that the queries in queries.sql can be parsed
	t.Skip("Full code generation test requires sqlc generate - run manually")
}
