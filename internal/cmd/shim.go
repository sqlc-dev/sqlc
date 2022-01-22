package cmd

import (
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/plugin"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func pluginSettings(cs config.CombinedSettings) *plugin.Settings {
	var over []*plugin.Override
	for _, o := range cs.Overrides {
		over = append(over, &plugin.Override{
			CodeType: "", // FIXME
			DbType:   o.DBType,
			Nullable: o.Nullable,
			Column:   o.Column,
		})
	}
	return &plugin.Settings{
		Version:   cs.Global.Version,
		Engine:    string(cs.Package.Engine),
		Schema:    []string(cs.Package.Schema),
		Queries:   []string(cs.Package.Queries),
		Overrides: over,
		Rename:    cs.Rename,
	}
}

func pluginCatalog(c *catalog.Catalog) *plugin.Catalog {
	var schemas []*plugin.Schema
	for _, s := range c.Schemas {
		var enums []*plugin.Enum
		for _, typ := range s.Types {
			enum, ok := typ.(*catalog.Enum)
			if !ok {
				continue
			}
			enums = append(enums, &plugin.Enum{
				Name:    enum.Name,
				Comment: enum.Comment,
				Vals:    enum.Vals,
			})
		}
		var tables []*plugin.Table
		for _, t := range s.Tables {
			var columns []*plugin.Column
			for _, c := range t.Columns {
				l := -1
				if c.Length != nil {
					l = *c.Length
				}
				columns = append(columns, &plugin.Column{
					Name: c.Name,
					Type: &plugin.Identifier{
						Catalog: c.Type.Catalog,
						Schema:  c.Type.Schema,
						Name:    c.Type.Name,
					},
					Comment: c.Comment,
					NotNull: c.IsNotNull,
					IsArray: c.IsArray,
					Length:  int32(l),
				})
			}
			tables = append(tables, &plugin.Table{
				Rel: &plugin.Identifier{
					Catalog: t.Rel.Catalog,
					Schema:  t.Rel.Schema,
					Name:    t.Rel.Name,
				},
				Columns: columns,
				Comment: t.Comment,
			})
		}
		schemas = append(schemas, &plugin.Schema{
			Comment: s.Comment,
			Name:    s.Name,
			Tables:  tables,
			Enums:   enums,
		})
	}
	return &plugin.Catalog{
		Name:          c.Name,
		DefaultSchema: c.DefaultSchema,
		Comment:       c.Comment,
		Schemas:       schemas,
	}
}

func pluginQueries(r *compiler.Result) []*plugin.Query {
	var out []*plugin.Query
	for _, q := range r.Queries {
		var params []*plugin.Parameter
		var columns []*plugin.Column
		for _, c := range q.Columns {
			columns = append(columns, pluginQueryColumn(c))
		}
		for _, p := range q.Params {
			params = append(params, pluginQueryParam(p))
		}
		out = append(out, &plugin.Query{
			Name:     q.Name,
			Cmd:      q.Cmd,
			Text:     q.SQL,
			Comments: q.Comments,
			Columns:  columns,
			Params:   params,
		})
	}
	return out
}

func pluginQueryColumn(c *compiler.Column) *plugin.Column {
	return &plugin.Column{
		Name: c.Name,
	}
}

func pluginQueryParam(p compiler.Parameter) *plugin.Parameter {
	return &plugin.Parameter{
		Number: int32(p.Number),
		Column: pluginQueryColumn(p.Column),
	}
}

func codeGenRequest(r *compiler.Result, settings config.CombinedSettings) *plugin.CodeGenRequest {
	return &plugin.CodeGenRequest{
		Settings: pluginSettings(settings),
		Catalog:  pluginCatalog(r.Catalog),
		Queries:  pluginQueries(r),
	}
}
