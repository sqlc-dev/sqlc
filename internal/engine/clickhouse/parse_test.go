package clickhouse

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		count int // expected statement count
	}{
		{
			name:  "simple select",
			sql:   "SELECT id, name FROM users",
			count: 1,
		},
		{
			name:  "select with where",
			sql:   "SELECT id, name FROM users WHERE id = 1",
			count: 1,
		},
		{
			name:  "create table",
			sql:   "CREATE TABLE users (id UInt64, name String) ENGINE = MergeTree() ORDER BY id",
			count: 1,
		},
		{
			name:  "insert",
			sql:   "INSERT INTO users (id, name) VALUES (1, 'test')",
			count: 1,
		},
		{
			name:  "multiple statements",
			sql:   "SELECT 1; SELECT 2",
			count: 2,
		},
	}

	p := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(tt.sql))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(stmts) != tt.count {
				t.Errorf("Parse() returned %d statements, want %d", len(stmts), tt.count)
			}
		})
	}
}

func TestCommentSyntax(t *testing.T) {
	p := NewParser()
	cs := p.CommentSyntax()
	if !cs.Dash {
		t.Error("expected Dash comment to be supported")
	}
	if !cs.SlashStar {
		t.Error("expected SlashStar comment to be supported")
	}
	if !cs.Hash {
		t.Error("expected Hash comment to be supported")
	}
}

func TestIsReservedKeyword(t *testing.T) {
	p := NewParser()
	tests := []struct {
		word     string
		reserved bool
	}{
		{"select", true},
		{"from", true},
		{"where", true},
		{"table", true},
		{"engine", true},
		{"foobar", false},
		{"mycolumn", false},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			if got := p.IsReservedKeyword(tt.word); got != tt.reserved {
				t.Errorf("IsReservedKeyword(%q) = %v, want %v", tt.word, got, tt.reserved)
			}
		})
	}
}
