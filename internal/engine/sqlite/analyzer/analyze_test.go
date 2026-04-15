package analyzer

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestAnalyzer_Analyze(t *testing.T) {
	db := config.Database{
		Managed: true,
	}
	a := New(db)
	defer a.Close(context.Background())

	ctx := context.Background()

	migrations := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT
		)`,
	}

	query := `SELECT id, name, email FROM users WHERE id = ?`
	node := &ast.TODO{}

	result, err := a.Analyze(ctx, node, query, migrations, nil)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if len(result.Columns) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(result.Columns))
	}

	expectedCols := []struct {
		name     string
		dataType string
	}{
		{"id", "integer"},
		{"name", "text"},
		{"email", "text"},
	}

	for i, expected := range expectedCols {
		if i >= len(result.Columns) {
			break
		}
		col := result.Columns[i]
		if col.Name != expected.name {
			t.Errorf("Column %d: expected name %q, got %q", i, expected.name, col.Name)
		}
		if col.DataType != expected.dataType {
			t.Errorf("Column %d: expected dataType %q, got %q", i, expected.dataType, col.DataType)
		}
		if col.Table == nil || col.Table.Name != "users" {
			t.Errorf("Column %d: expected table 'users', got %v", i, col.Table)
		}
	}

	if len(result.Params) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(result.Params))
	}
}

func TestAnalyzer_InvalidQuery(t *testing.T) {
	db := config.Database{
		Managed: true,
	}
	a := New(db)
	defer a.Close(context.Background())

	ctx := context.Background()

	migrations := []string{
		`CREATE TABLE users (id INTEGER PRIMARY KEY)`,
	}

	query := `SELECT * FROM nonexistent`
	node := &ast.TODO{}

	_, err := a.Analyze(ctx, node, query, migrations, nil)
	if err == nil {
		t.Error("Expected error for invalid query, got nil")
	}
}

func TestNormalizeType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"INTEGER", "integer"},
		{"INT", "integer"},
		{"BIGINT", "integer"},
		{"TEXT", "text"},
		{"VARCHAR(255)", "text"},
		{"BLOB", "blob"},
		{"REAL", "real"},
		{"FLOAT", "real"},
		{"DOUBLE", "real"},
		{"BOOLEAN", "boolean"},
		{"DATETIME", "datetime"},
		{"", "any"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeType(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
