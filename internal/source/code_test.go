package source

import (
	"strings"
	"testing"
)

func TestStripComments(t *testing.T) {
	type test struct {
		name        string
		input       string
		wantSQL     string
		wantComment []string
	}

	tests := []test{
		{
			name:        "plain block comment on its own line is stripped",
			input:       "SELECT 1\n/* a comment */\nFROM foo",
			wantSQL:     "SELECT 1\nFROM foo",
			wantComment: []string{" a comment "},
		},
		{
			name:    "inline optimizer hint is preserved",
			input:   "SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM t1",
			wantSQL: "SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM t1",
		},
		{
			name:    "multi-line optimizer hint is preserved",
			input:   "SELECT\n/*+ MAX_EXECUTION_TIME(1000) */\n*\nFROM t1",
			wantSQL: "SELECT\n/*+ MAX_EXECUTION_TIME(1000) */\n*\nFROM t1",
		},
		{
			name:    "query name comment is dropped",
			input:   "/* name: Foo :one */\nSELECT 1",
			wantSQL: "SELECT 1",
		},
		{
			name:        "dash comments are collected",
			input:       "-- helpful note\nSELECT 1",
			wantSQL:     "SELECT 1",
			wantComment: []string{" helpful note"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSQL, gotComments, err := StripComments(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotSQL != tc.wantSQL {
				t.Errorf("SQL mismatch\n got: %q\nwant: %q", gotSQL, tc.wantSQL)
			}
			if strings.Join(gotComments, "|") != strings.Join(tc.wantComment, "|") {
				t.Errorf("comments mismatch\n got: %q\nwant: %q", gotComments, tc.wantComment)
			}
		})
	}
}
