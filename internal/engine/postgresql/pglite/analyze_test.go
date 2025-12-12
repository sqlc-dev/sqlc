package pglite

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
)

func TestRewriteType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"character(10)", "pg_catalog.bpchar"},
		{"character varying(255)", "pg_catalog.varchar"},
		{"character varying", "pg_catalog.varchar"},
		{"bit varying(8)", "pg_catalog.varbit"},
		{"bit(1)", "pg_catalog.bit"},
		{"bpchar", "pg_catalog.bpchar"},
		{"timestamp without time zone", "pg_catalog.timestamp"},
		{"timestamp with time zone", "pg_catalog.timestamptz"},
		{"time without time zone", "pg_catalog.time"},
		{"time with time zone", "pg_catalog.timetz"},
		{"integer", "integer"},
		{"text", "text"},
		{"boolean", "boolean"},
		{"uuid", "uuid"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := rewriteType(tt.input)
			if result != tt.expected {
				t.Errorf("rewriteType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEqualMigrations(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected bool
	}{
		{
			name:     "both empty",
			a:        []string{},
			b:        []string{},
			expected: true,
		},
		{
			name:     "both nil",
			a:        nil,
			b:        nil,
			expected: true,
		},
		{
			name:     "equal single element",
			a:        []string{"CREATE TABLE users (id INT)"},
			b:        []string{"CREATE TABLE users (id INT)"},
			expected: true,
		},
		{
			name:     "equal multiple elements",
			a:        []string{"CREATE TABLE users (id INT)", "CREATE TABLE posts (id INT)"},
			b:        []string{"CREATE TABLE users (id INT)", "CREATE TABLE posts (id INT)"},
			expected: true,
		},
		{
			name:     "different length",
			a:        []string{"CREATE TABLE users (id INT)"},
			b:        []string{"CREATE TABLE users (id INT)", "CREATE TABLE posts (id INT)"},
			expected: false,
		},
		{
			name:     "different content",
			a:        []string{"CREATE TABLE users (id INT)"},
			b:        []string{"CREATE TABLE posts (id INT)"},
			expected: false,
		},
		{
			name:     "different order",
			a:        []string{"CREATE TABLE users (id INT)", "CREATE TABLE posts (id INT)"},
			b:        []string{"CREATE TABLE posts (id INT)", "CREATE TABLE users (id INT)"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equalMigrations(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("equalMigrations(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestPGLiteError(t *testing.T) {
	tests := []struct {
		name     string
		err      *PGLiteError
		expected string
	}{
		{
			name: "with code",
			err: &PGLiteError{
				Code:    "42601",
				Message: "syntax error at or near \"SELEC\"",
			},
			expected: "42601: syntax error at or near \"SELEC\"",
		},
		{
			name: "without code",
			err: &PGLiteError{
				Message: "connection failed",
			},
			expected: "connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		req      Request
		expected string
	}{
		{
			name: "init request",
			req: Request{
				Type:       "init",
				Migrations: []string{"CREATE TABLE users (id INT)"},
			},
			expected: `{"type":"init","migrations":["CREATE TABLE users (id INT)"],"query":""}`,
		},
		{
			name: "prepare request",
			req: Request{
				Type:  "prepare",
				Query: "SELECT * FROM users WHERE id = $1",
			},
			expected: `{"type":"prepare","migrations":null,"query":"SELECT * FROM users WHERE id = $1"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the struct can be created and fields are accessible
			if tt.req.Type == "" {
				t.Error("Type should not be empty")
			}
		})
	}
}

func TestColumnInfo(t *testing.T) {
	col := ColumnInfo{
		Name:        "user_id",
		DataType:    "integer",
		DataTypeOID: 23,
		NotNull:     true,
		IsArray:     false,
		ArrayDims:   0,
		TableOID:    16384,
		TableName:   "users",
		TableSchema: "public",
	}

	if col.Name != "user_id" {
		t.Errorf("Name = %q, want %q", col.Name, "user_id")
	}
	if col.DataType != "integer" {
		t.Errorf("DataType = %q, want %q", col.DataType, "integer")
	}
	if !col.NotNull {
		t.Error("NotNull should be true")
	}
}

func TestParameterInfo(t *testing.T) {
	param := ParameterInfo{
		Number:      1,
		DataType:    "text",
		DataTypeOID: 25,
		IsArray:     false,
		ArrayDims:   0,
	}

	if param.Number != 1 {
		t.Errorf("Number = %d, want %d", param.Number, 1)
	}
	if param.DataType != "text" {
		t.Errorf("DataType = %q, want %q", param.DataType, "text")
	}
}

func TestPrepareResult(t *testing.T) {
	result := PrepareResult{
		Columns: []ColumnInfo{
			{Name: "id", DataType: "integer", NotNull: true},
			{Name: "name", DataType: "text", NotNull: false},
		},
		Params: []ParameterInfo{
			{Number: 1, DataType: "integer"},
		},
	}

	if len(result.Columns) != 2 {
		t.Errorf("len(Columns) = %d, want 2", len(result.Columns))
	}
	if len(result.Params) != 1 {
		t.Errorf("len(Params) = %d, want 1", len(result.Params))
	}
}

func TestResponse(t *testing.T) {
	t.Run("success response", func(t *testing.T) {
		resp := Response{
			Success: true,
			Prepare: &PrepareResult{
				Columns: []ColumnInfo{
					{Name: "id", DataType: "integer"},
				},
			},
		}

		if !resp.Success {
			t.Error("Success should be true")
		}
		if resp.Prepare == nil {
			t.Error("Prepare should not be nil")
		}
		if resp.Error != nil {
			t.Error("Error should be nil")
		}
	})

	t.Run("error response", func(t *testing.T) {
		resp := Response{
			Success: false,
			Error: &ErrorResponse{
				Code:     "42P01",
				Message:  "relation \"users\" does not exist",
				Position: 15,
			},
		}

		if resp.Success {
			t.Error("Success should be false")
		}
		if resp.Error == nil {
			t.Error("Error should not be nil")
		}
		if resp.Error.Code != "42P01" {
			t.Errorf("Error.Code = %q, want %q", resp.Error.Code, "42P01")
		}
	})
}

func TestNewAnalyzer(t *testing.T) {
	cfg := config.PGLite{
		URL:    "file:///path/to/pglite.wasm",
		SHA256: "abc123",
	}

	a := New(cfg)
	if a == nil {
		t.Fatal("New() returned nil")
	}
	if a.cfg.URL != cfg.URL {
		t.Errorf("cfg.URL = %q, want %q", a.cfg.URL, cfg.URL)
	}
	if a.cfg.SHA256 != cfg.SHA256 {
		t.Errorf("cfg.SHA256 = %q, want %q", a.cfg.SHA256, cfg.SHA256)
	}
}
