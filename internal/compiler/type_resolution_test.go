package compiler

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/opts"
)

func TestNewCompilerClickHouse(t *testing.T) {
	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}

	combo := config.CombinedSettings{
		Global: config.Config{},
	}

	c, err := NewCompiler(conf, combo)
	if err != nil {
		t.Fatalf("unexpected error creating ClickHouse compiler: %v", err)
	}

	if c.parser == nil {
		t.Error("expected parser to be set")
	}

	if c.catalog == nil {
		t.Error("expected catalog to be set")
	}

	if c.parser.CommentSyntax().Dash == false {
		t.Error("expected ClickHouse parser to support dash comments")
	}

	if c.parser.CommentSyntax().SlashStar == false {
		t.Error("expected ClickHouse parser to support slash-star comments")
	}
}

func TestClickHouseTypeResolver(t *testing.T) {
	schema := `
CREATE TABLE events (
    id UInt64,
    tags Array(String),
    scores Array(UInt32)
) ENGINE = Memory;
`
	query := `
-- name: TestArrayJoin :many
SELECT arrayJoin(tags) as tag FROM events;

-- name: TestCount :one
SELECT count(*) FROM events;

-- name: TestArgMin :one
SELECT argMin(id, id) FROM events;

-- name: TestAny :one
SELECT any(tags) FROM events;
`

	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}

	c, err := NewCompiler(conf, config.CombinedSettings{})
	if err != nil {
		t.Fatal(err)
	}

	// Manually update catalog with schema
	c.schema = append(c.schema, schema)
	stmts, err := c.parser.Parse(strings.NewReader(schema))
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range stmts {
		if err := c.catalog.Update(stmt, c); err != nil {
			t.Fatal(err)
		}
	}

	// Parse queries
	queryStmts, err := c.parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	for _, stmt := range queryStmts {
		q, err := c.parseQuery(stmt.Raw, query, opts.Parser{})
		if err != nil {
			t.Fatal(err)
		}
		if q == nil {
			continue
		}

		// Verify types
		switch q.Metadata.Name {
		case "TestArrayJoin":
			if len(q.Columns) != 1 {
				t.Errorf("TestArrayJoin: expected 1 column, got %d", len(q.Columns))
			}
			col := q.Columns[0]
			if col.DataType != "text" { // String maps to text
				t.Errorf("TestArrayJoin: expected text, got %s", col.DataType)
			}
			if col.IsArray {
				t.Errorf("TestArrayJoin: expected not array")
			}

		case "TestCount":
			if len(q.Columns) != 1 {
				t.Errorf("TestCount: expected 1 column, got %d", len(q.Columns))
			}
			col := q.Columns[0]
			if col.DataType != "uint64" {
				t.Errorf("TestCount: expected uint64, got %s", col.DataType)
			}

		case "TestArgMin":
			if len(q.Columns) != 1 {
				t.Errorf("TestArgMin: expected 1 column, got %d", len(q.Columns))
			}
			col := q.Columns[0]
			if col.DataType != "uint64" {
				t.Errorf("TestArgMin: expected uint64, got %s", col.DataType)
			}

		case "TestAny":
			if len(q.Columns) != 1 {
				t.Errorf("TestAny: expected 1 column, got %d", len(q.Columns))
			}
			col := q.Columns[0]
			if col.DataType != "text" { // text[] in sqlc usually maps to text with IsArray=true
				t.Errorf("TestAny: expected text, got %s", col.DataType)
			}
			if !col.IsArray {
				t.Errorf("TestAny: expected array")
			}
		}
	}
}

func TestClickHouseLimitParameterType(t *testing.T) {
	schema := `
CREATE TABLE users (
    id UInt32,
    name String,
    email String,
    created_at DateTime
) ENGINE = Memory;
`
	query := `
-- name: ListUsers :many
SELECT id, name, email, created_at
FROM users
LIMIT ?;
`

	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}

	c, err := NewCompiler(conf, config.CombinedSettings{})
	if err != nil {
		t.Fatal(err)
	}

	// Manually update catalog with schema
	c.schema = append(c.schema, schema)
	stmts, err := c.parser.Parse(strings.NewReader(schema))
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range stmts {
		if err := c.catalog.Update(stmt, c); err != nil {
			t.Fatal(err)
		}
	}

	// Parse queries
	queryStmts, err := c.parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	for _, stmt := range queryStmts {
		q, err := c.parseQuery(stmt.Raw, query, opts.Parser{})
		if err != nil {
			t.Fatal(err)
		}
		if q == nil {
			continue
		}

		// Check parameters
		if len(q.Params) != 1 {
			t.Errorf("Expected 1 parameter, got %d", len(q.Params))
		}

		param := q.Params[0]
		t.Logf("Parameter: Name=%s, DataType=%s, Number=%d", param.Column.Name, param.Column.DataType, param.Number)

		if param.Column.DataType != "integer" {
			t.Errorf("Expected integer type for LIMIT parameter, got %s", param.Column.DataType)
		}
	}
}

func TestClickHouseLowCardinality(t *testing.T) {
	schema := `
CREATE TABLE products (
    id UInt32,
    category LowCardinality(String),
    status LowCardinality(String),
    priority LowCardinality(UInt8)
) ENGINE = Memory;
`
	query := `
-- name: GetProductsByCategory :many
SELECT id, category, status, priority FROM products WHERE category = ?;
`

	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}

	c, err := NewCompiler(conf, config.CombinedSettings{})
	if err != nil {
		t.Fatal(err)
	}

	// Manually update catalog with schema
	c.schema = append(c.schema, schema)
	stmts, err := c.parser.Parse(strings.NewReader(schema))
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range stmts {
		if err := c.catalog.Update(stmt, c); err != nil {
			t.Fatal(err)
		}
	}

	// Parse queries
	queryStmts, err := c.parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	for _, stmt := range queryStmts {
		q, err := c.parseQuery(stmt.Raw, query, opts.Parser{})
		if err != nil {
			t.Fatal(err)
		}
		if q == nil {
			continue
		}

		// Check that LowCardinality columns are correctly resolved
		if len(q.Columns) < 4 {
			t.Fatalf("Expected at least 4 columns, got %d", len(q.Columns))
		}

		tests := []struct {
			colIndex int
			name     string
			dataType string
		}{
			{0, "id", "uint32"},
			{1, "category", "text"},  // LowCardinality(String) -> String -> text
			{2, "status", "text"},    // LowCardinality(String) -> String -> text
			{3, "priority", "uint8"}, // LowCardinality(UInt8) -> UInt8 -> uint8
		}

		for _, test := range tests {
			if test.colIndex >= len(q.Columns) {
				t.Errorf("Column index %d out of bounds", test.colIndex)
				continue
			}

			col := q.Columns[test.colIndex]
			if col.Name != test.name {
				t.Errorf("Column %d: expected name %q, got %q", test.colIndex, test.name, col.Name)
			}
			if col.DataType != test.dataType {
				t.Errorf("Column %q: expected type %q, got %q", col.Name, test.dataType, col.DataType)
			}
		}

		// Check parameter (the WHERE clause)
		if len(q.Params) != 1 {
			t.Errorf("Expected 1 parameter, got %d", len(q.Params))
		} else {
			param := q.Params[0]
			if param.Column.DataType != "text" {
				t.Errorf("Expected text type for category parameter, got %s", param.Column.DataType)
			}
		}
	}
}

func TestClickHouseIPAddressTypes(t *testing.T) {
	schema := `
CREATE TABLE network_data (
    id UInt32,
    source_ip IPv4,
    dest_ip IPv4,
    ipv6_addr IPv6,
    nullable_ip Nullable(IPv4)
) ENGINE = Memory;
`
	query := `
-- name: GetNetworkData :one
SELECT id, source_ip, dest_ip, ipv6_addr, nullable_ip FROM network_data WHERE id = ?;

-- name: FilterByIPv4 :many
SELECT id, source_ip FROM network_data WHERE source_ip = ?;

-- name: FilterByIPv6 :many
SELECT id, ipv6_addr FROM network_data WHERE ipv6_addr = ?;
`

	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}

	c, err := NewCompiler(conf, config.CombinedSettings{})
	if err != nil {
		t.Fatal(err)
	}

	// Manually update catalog with schema
	c.schema = append(c.schema, schema)
	stmts, err := c.parser.Parse(strings.NewReader(schema))
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range stmts {
		if err := c.catalog.Update(stmt, c); err != nil {
			t.Fatal(err)
		}
	}

	// Parse queries
	queryStmts, err := c.parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	for _, stmt := range queryStmts {
		q, err := c.parseQuery(stmt.Raw, query, opts.Parser{})
		if err != nil {
			t.Fatal(err)
		}
		if q == nil {
			continue
		}

		switch q.Metadata.Name {
		case "GetNetworkData":
			if len(q.Columns) != 5 {
				t.Errorf("GetNetworkData: expected 5 columns, got %d", len(q.Columns))
			}

			tests := []struct {
				colIndex int
				name     string
				dataType string
			}{
				{0, "id", "uint32"},
				{1, "source_ip", "ipv4"},
				{2, "dest_ip", "ipv4"},
				{3, "ipv6_addr", "ipv6"},
				{4, "nullable_ip", "ipv4"},
			}

			for _, test := range tests {
				if test.colIndex >= len(q.Columns) {
					t.Errorf("Column index %d out of bounds", test.colIndex)
					continue
				}

				col := q.Columns[test.colIndex]
				if col.Name != test.name {
					t.Errorf("Column %d: expected name %q, got %q", test.colIndex, test.name, col.Name)
				}
				if col.DataType != test.dataType {
					t.Errorf("Column %q: expected type %q, got %q", col.Name, test.dataType, col.DataType)
				}
			}

			// Check parameter (id filter)
			if len(q.Params) != 1 {
				t.Errorf("Expected 1 parameter, got %d", len(q.Params))
			} else {
				param := q.Params[0]
				if param.Column.DataType != "uint32" {
					t.Errorf("Expected uint32 type for id parameter, got %s", param.Column.DataType)
				}
			}

		case "FilterByIPv4":
			if len(q.Columns) != 2 {
				t.Errorf("FilterByIPv4: expected 2 columns, got %d", len(q.Columns))
			}

			col := q.Columns[1]
			if col.Name != "source_ip" {
				t.Errorf("Expected column name source_ip, got %q", col.Name)
			}
			if col.DataType != "ipv4" {
				t.Errorf("Expected ipv4 type, got %s", col.DataType)
			}

			// Check parameter (IPv4 filter)
			if len(q.Params) != 1 {
				t.Errorf("Expected 1 parameter, got %d", len(q.Params))
			} else {
				param := q.Params[0]
				if param.Column.DataType != "ipv4" {
					t.Errorf("Expected ipv4 type for IP parameter, got %s", param.Column.DataType)
				}
			}

		case "FilterByIPv6":
			if len(q.Columns) != 2 {
				t.Errorf("FilterByIPv6: expected 2 columns, got %d", len(q.Columns))
			}

			col := q.Columns[1]
			if col.Name != "ipv6_addr" {
				t.Errorf("Expected column name ipv6_addr, got %q", col.Name)
			}
			if col.DataType != "ipv6" {
				t.Errorf("Expected ipv6 type, got %s", col.DataType)
			}

			// Check parameter (IPv6 filter)
			if len(q.Params) != 1 {
				t.Errorf("Expected 1 parameter, got %d", len(q.Params))
			} else {
				param := q.Params[0]
				if param.Column.DataType != "ipv6" {
					t.Errorf("Expected ipv6 type for IP parameter, got %s", param.Column.DataType)
				}
			}
		}
	}
}

func TestClickHouseMapType(t *testing.T) {
	schema := `
CREATE TABLE config (
    id UInt32,
    settings Map(String, String),
    metrics Map(String, UInt64),
    nested_data Map(String, Array(String)),
    invalid_key Map(Array(String), String)
) ENGINE = Memory;
`
	query := `
-- name: GetConfig :one
SELECT id, settings, metrics, nested_data, invalid_key FROM config WHERE id = ?;
`

	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}

	c, err := NewCompiler(conf, config.CombinedSettings{})
	if err != nil {
		t.Fatal(err)
	}

	// Manually update catalog with schema
	c.schema = append(c.schema, schema)
	stmts, err := c.parser.Parse(strings.NewReader(schema))
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range stmts {
		if err := c.catalog.Update(stmt, c); err != nil {
			t.Fatal(err)
		}
	}

	// Parse queries
	queryStmts, err := c.parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	for _, stmt := range queryStmts {
		q, err := c.parseQuery(stmt.Raw, query, opts.Parser{})
		if err != nil {
			t.Fatal(err)
		}
		if q == nil {
			continue
		}

		// Check that Map columns are correctly resolved
		if len(q.Columns) < 5 {
			t.Fatalf("Expected at least 5 columns, got %d", len(q.Columns))
		}

		tests := []struct {
			colIndex int
			name     string
			dataType string
		}{
			{0, "id", "uint32"},
			{1, "settings", "map[string]string"},         // Map(String, String) -> map[string]string
			{2, "metrics", "map[string]uint64"},          // Map(String, UInt64) -> map[string]uint64
			{3, "nested_data", "map[string][]string"},    // Map(String, Array(String)) -> map[string][]string
			{4, "invalid_key", "map[string]interface{}"}, // Map(Array(String), String) -> falls back due to invalid key
		}

		for _, test := range tests {
			if test.colIndex >= len(q.Columns) {
				t.Errorf("Column index %d out of bounds", test.colIndex)
				continue
			}

			col := q.Columns[test.colIndex]
			if col.Name != test.name {
				t.Errorf("Column %d: expected name %q, got %q", test.colIndex, test.name, col.Name)
			}
			if col.DataType != test.dataType {
				t.Errorf("Column %q: expected type %q, got %q", col.Name, test.dataType, col.DataType)
			}
		}
	}
}
