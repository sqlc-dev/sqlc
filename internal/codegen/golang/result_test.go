package golang

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func boolPtr(b bool) *bool { return &b }

func catalogSchemaRequest() *plugin.GenerateRequest {
	col := func(name, typ string) *plugin.Column {
		return &plugin.Column{Name: name, Type: &plugin.Identifier{Name: typ}}
	}
	table := func(schema, name string, cols ...*plugin.Column) *plugin.Schema {
		return &plugin.Schema{
			Name: schema,
			Tables: []*plugin.Table{
				{Rel: &plugin.Identifier{Schema: schema, Name: name}, Columns: cols},
			},
		}
	}
	return &plugin.GenerateRequest{
		Settings: &plugin.Settings{Engine: "postgresql"},
		Catalog: &plugin.Catalog{
			DefaultSchema: "public",
			Schemas: []*plugin.Schema{
				table("pg_catalog", "pg_class", col("oid", "oid")),
				table("information_schema", "tables", col("table_name", "text")),
				table("public", "users", col("id", "int4")),
			},
		},
	}
}

func TestBuildStructs_OmitCatalogSchema(t *testing.T) {
	req := catalogSchemaRequest()

	t.Run("unset (default) skips pg_catalog and information_schema", func(t *testing.T) {
		structs := buildStructs(req, &opts.Options{})
		for _, s := range structs {
			if s.Table != nil && (s.Table.Schema == "pg_catalog" || s.Table.Schema == "information_schema") {
				t.Errorf("unexpected catalog struct: %s.%s", s.Table.Schema, s.Table.Name)
			}
		}
	})

	t.Run("omit=true skips pg_catalog and information_schema", func(t *testing.T) {
		structs := buildStructs(req, &opts.Options{OmitCatalogSchema: boolPtr(true)})
		for _, s := range structs {
			if s.Table != nil && (s.Table.Schema == "pg_catalog" || s.Table.Schema == "information_schema") {
				t.Errorf("unexpected catalog struct: %s.%s", s.Table.Schema, s.Table.Name)
			}
		}
	})

	t.Run("omit=false includes pg_catalog and information_schema", func(t *testing.T) {
		structs := buildStructs(req, &opts.Options{OmitCatalogSchema: boolPtr(false)})
		schemas := make(map[string]bool)
		for _, s := range structs {
			if s.Table != nil {
				schemas[s.Table.Schema] = true
			}
		}
		for _, want := range []string{"pg_catalog", "information_schema", "public"} {
			if !schemas[want] {
				t.Errorf("expected structs for schema %q, got none", want)
			}
		}
	})
}

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
