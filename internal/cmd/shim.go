package cmd

import (
	"strings"

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/plugin"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func pluginOverride(o config.Override) *plugin.Override {
	var column string
	var table plugin.Identifier

	if o.Column != "" {
		colParts := strings.Split(o.Column, ".")
		switch len(colParts) {
		case 2:
			table.Schema = "public"
			table.Name = colParts[0]
			column = colParts[1]
		case 3:
			table.Schema = colParts[0]
			table.Name = colParts[1]
			column = colParts[2]
		case 4:
			table.Catalog = colParts[0]
			table.Schema = colParts[1]
			table.Name = colParts[2]
			column = colParts[3]
		}
	}
	return &plugin.Override{
		CodeType:   "", // FIXME
		DbType:     o.DBType,
		Nullable:   o.Nullable,
		Column:     o.Column,
		ColumnName: column,
		Table:      &table,
		PythonType: pluginPythonType(o.PythonType),
	}
}

func pluginSettings(cs config.CombinedSettings) *plugin.Settings {
	var over []*plugin.Override
	for _, o := range cs.Overrides {
		over = append(over, pluginOverride(o))
	}
	return &plugin.Settings{
		Version:   cs.Global.Version,
		Engine:    string(cs.Package.Engine),
		Schema:    []string(cs.Package.Schema),
		Queries:   []string(cs.Package.Queries),
		Overrides: over,
		Rename:    cs.Rename,
		Python:    pluginPythonCode(cs.Python),
		Kotlin:    pluginKotlinCode(cs.Kotlin),
	}
}

func pluginPythonCode(s config.SQLPython) *plugin.PythonCode {
	return &plugin.PythonCode{
		Out:                 s.Out,
		Package:             s.Package,
		EmitExactTableNames: s.EmitExactTableNames,
		EmitSyncQuerier:     s.EmitSyncQuerier,
		EmitAsyncQuerier:    s.EmitAsyncQuerier,
	}
}

func pluginPythonType(pt config.PythonType) *plugin.PythonType {
	return &plugin.PythonType{
		Module: pt.Module,
		Name:   pt.Name,
	}
}

func pluginKotlinCode(s config.SQLKotlin) *plugin.KotlinCode {
	return &plugin.KotlinCode{
		Out:                 s.Out,
		Package:             s.Package,
		EmitExactTableNames: s.EmitExactTableNames,
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
					Table: &plugin.Identifier{
						Catalog: t.Rel.Catalog,
						Schema:  t.Rel.Schema,
						Name:    t.Rel.Name,
					},
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
			Filename: q.Filename,
		})
	}
	return out
}

func pluginQueryColumn(c *compiler.Column) *plugin.Column {
	l := -1
	if c.Length != nil {
		l = *c.Length
	}
	out := &plugin.Column{
		Name:    c.Name,
		Comment: c.Comment,
		NotNull: c.NotNull,
		IsArray: c.IsArray,
		Length:  int32(l),
	}

	if c.Type != nil {
		out.Type = &plugin.Identifier{
			Catalog: c.Type.Catalog,
			Schema:  c.Type.Schema,
			Name:    c.Type.Name,
		}
	} else {
		out.Type = &plugin.Identifier{
			Name: c.DataType,
		}
	}

	if c.Table != nil {
		out.Table = &plugin.Identifier{
			Catalog: c.Table.Catalog,
			Schema:  c.Table.Schema,
			Name:    c.Table.Name,
		}
	}

	return out
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
