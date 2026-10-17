package golang

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func TestPutOutColumns_ForZeroColumns(t *testing.T) {
	tests := []struct {
		cmd  string
		want bool
	}{
		{
			cmd:  metadata.CmdExec,
			want: false,
		},
		{
			cmd:  metadata.CmdExecResult,
			want: false,
		},
		{
			cmd:  metadata.CmdExecRows,
			want: false,
		},
		{
			cmd:  metadata.CmdExecLastId,
			want: false,
		},
		{
			cmd:  metadata.CmdMany,
			want: true,
		},
		{
			cmd:  metadata.CmdOne,
			want: true,
		},
		{
			cmd:  metadata.CmdCopyFrom,
			want: false,
		},
		{
			cmd:  metadata.CmdBatchExec,
			want: false,
		},
		{
			cmd:  metadata.CmdBatchMany,
			want: true,
		},
		{
			cmd:  metadata.CmdBatchOne,
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.cmd, func(t *testing.T) {
			query := &plugin.Query{
				Cmd:     tc.cmd,
				Columns: []*plugin.Column{},
			}
			got := putOutColumns(query)
			if got != tc.want {
				t.Errorf("putOutColumns failed. want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestPutOutColumns_AlwaysTrueWhenQueryHasColumns(t *testing.T) {
	query := &plugin.Query{
		Cmd:     metadata.CmdMany,
		Columns: []*plugin.Column{{}},
	}
	if putOutColumns(query) != true {
		t.Error("should be true when we have columns")
	}
}

// TestBuildEnums_DeduplicatesConstantNames covers enum values that differ only
// by characters that get stripped or collapsed when building a Go identifier
// (e.g. "A+" and "A-"), which previously produced duplicate constant names and
// uncompilable output.
func TestBuildEnums_DeduplicatesConstantNames(t *testing.T) {
	req := &plugin.GenerateRequest{
		Catalog: &plugin.Catalog{
			DefaultSchema: "public",
			Schemas: []*plugin.Schema{
				{
					Name: "public",
					Enums: []*plugin.Enum{
						{
							Name: "blood_group_type",
							Vals: []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"},
						},
					},
				},
			},
		},
	}

	enums := buildEnums(req, &opts.Options{})
	if len(enums) != 1 {
		t.Fatalf("expected 1 enum, got %d", len(enums))
	}
	if got := len(enums[0].Constants); got != 8 {
		t.Fatalf("expected 8 constants, got %d", got)
	}

	seen := make(map[string]string, 8)
	for _, c := range enums[0].Constants {
		if prev, dup := seen[c.Name]; dup {
			t.Errorf("duplicate constant name %q for values %q and %q", c.Name, prev, c.Value)
		}
		seen[c.Name] = c.Value
	}
}
