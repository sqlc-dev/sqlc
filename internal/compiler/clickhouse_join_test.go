package compiler

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestClickHouseJoinColumnResolution tests that column names are properly resolved
// in JOIN queries now that JoinExpr is correctly converted
func TestClickHouseJoinColumnResolution(t *testing.T) {
	parser := clickhouse.NewParser()
	cat := clickhouse.NewCatalog()

	// Create database and tables
	schemaSQL := `CREATE DATABASE IF NOT EXISTS test_db;
CREATE TABLE test_db.users (
	id UInt32,
	name String,
	email String
);
CREATE TABLE test_db.posts (
	id UInt32,
	user_id UInt32,
	title String,
	content String
)`

	stmts, err := parser.Parse(strings.NewReader(schemaSQL))
	if err != nil {
		t.Fatalf("Parse schema failed: %v", err)
	}

	for _, stmt := range stmts {
		if err := cat.Update(stmt, nil); err != nil {
			t.Fatalf("Update catalog failed: %v", err)
		}
	}

	// Create compiler
	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}
	combo := config.CombinedSettings{
		Global: config.Config{},
	}

	c, err := NewCompiler(conf, combo)
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	// Replace catalog
	c.catalog = cat

	// Parse a JOIN query
	querySQL := "SELECT u.id, u.name, p.id as post_id, p.title FROM test_db.users u LEFT JOIN test_db.posts p ON u.id = p.user_id WHERE u.id = 1"
	queryStmts, err := parser.Parse(strings.NewReader(querySQL))
	if err != nil {
		t.Fatalf("Parse query failed: %v", err)
	}

	if len(queryStmts) == 0 {
		t.Fatal("No queries parsed")
	}

	selectStmt := queryStmts[0].Raw.Stmt
	if selectStmt == nil {
		t.Fatal("Select statement is nil")
	}

	selectAst, ok := selectStmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Expected SelectStmt, got %T", selectStmt)
	}

	// Build query catalog and get output columns
	qc, err := c.buildQueryCatalog(c.catalog, selectAst, nil)
	if err != nil {
		t.Fatalf("Failed to build query catalog: %v", err)
	}

	cols, err := c.outputColumns(qc, selectAst)
	if err != nil {
		t.Fatalf("Failed to get output columns: %v", err)
	}

	if len(cols) != 4 {
		t.Errorf("Expected 4 columns, got %d", len(cols))
	}

	expectedNames := []string{"id", "name", "post_id", "title"}
	for i, expected := range expectedNames {
		if i < len(cols) {
			if cols[i].Name != expected {
				t.Errorf("Column %d: expected name %q, got %q", i, expected, cols[i].Name)
			}
		}
	}
}

// TestClickHouseLeftJoinNullability tests that LEFT JOIN correctly marks right-side columns as nullable
// In ClickHouse, columns are non-nullable by default unless wrapped in Nullable(T)
func TestClickHouseLeftJoinNullability(t *testing.T) {
	parser := clickhouse.NewParser()
	cat := clickhouse.NewCatalog()

	schemaSQL := `CREATE TABLE orders (
		order_id UInt32,
		customer_name String,
		amount Float64,
		created_at DateTime
	);
	CREATE TABLE shipments (
		shipment_id UInt32,
		order_id UInt32,
		address String,
		shipped_at DateTime
	)`

	stmts, err := parser.Parse(strings.NewReader(schemaSQL))
	if err != nil {
		t.Fatalf("Parse schema failed: %v", err)
	}

	for _, stmt := range stmts {
		if err := cat.Update(stmt, nil); err != nil {
			t.Fatalf("Update catalog failed: %v", err)
		}
	}

	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}
	combo := config.CombinedSettings{
		Global: config.Config{},
	}

	c, err := NewCompiler(conf, combo)
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}
	c.catalog = cat

	querySQL := "SELECT o.order_id, o.customer_name, o.amount, o.created_at, s.shipment_id, s.address, s.shipped_at FROM orders o LEFT JOIN shipments s ON o.order_id = s.order_id ORDER BY o.created_at DESC"
	queryStmts, err := parser.Parse(strings.NewReader(querySQL))
	if err != nil {
		t.Fatalf("Parse query failed: %v", err)
	}

	selectAst := queryStmts[0].Raw.Stmt.(*ast.SelectStmt)
	qc, err := c.buildQueryCatalog(c.catalog, selectAst, nil)
	if err != nil {
		t.Fatalf("Failed to build query catalog: %v", err)
	}

	cols, err := c.outputColumns(qc, selectAst)
	if err != nil {
		t.Fatalf("Failed to get output columns: %v", err)
	}

	if len(cols) != 7 {
		t.Errorf("Expected 7 columns, got %d", len(cols))
	}

	// Left table columns should be non-nullable
	leftTableNonNull := map[string]bool{
		"order_id":      true,
		"customer_name": true,
		"amount":        true,
		"created_at":    true,
	}

	// Right table columns should be nullable (because of LEFT JOIN)
	rightTableNullable := map[string]bool{
		"shipment_id": true,
		"address":     true,
		"shipped_at":  true,
	}

	for _, col := range cols {
		if expected, ok := leftTableNonNull[col.Name]; ok {
			if col.NotNull != expected {
				t.Errorf("Column %q: expected NotNull=%v, got %v", col.Name, expected, col.NotNull)
			}
		}
		if expected, ok := rightTableNullable[col.Name]; ok {
			if col.NotNull == expected {
				t.Errorf("Column %q: expected NotNull=%v, got %v", col.Name, !expected, col.NotNull)
			}
		}
	}
}
