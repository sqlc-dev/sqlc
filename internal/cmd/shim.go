package cmd

import (
	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/config/convert"
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func pluginSettings(r *compiler.Result, cs config.CombinedSettings) *plugin.Settings {
	return &plugin.Settings{
		Version: cs.Global.Version,
		Engine:  string(cs.Package.Engine),
		Schema:  []string(cs.Package.Schema),
		Queries: []string(cs.Package.Queries),
		Codegen: pluginCodegen(cs, cs.Codegen),
	}
}

func pluginCodegen(cs config.CombinedSettings, s config.Codegen) *plugin.Codegen {
	opts, err := convert.YAMLtoJSON(s.Options)
	if err != nil {
		panic(err)
	}
	cg := &plugin.Codegen{
		Out:     s.Out,
		Plugin:  s.Plugin,
		Options: opts,
	}
	for _, p := range cs.Global.Plugins {
		if p.Name == s.Plugin {
			cg.Env = p.Env
			cg.Process = pluginProcess(p)
			cg.Wasm = pluginWASM(p)
			return cg
		}
	}
	return cg
}

func pluginProcess(p config.Plugin) *plugin.Codegen_Process {
	if p.Process != nil {
		return &plugin.Codegen_Process{
			Cmd: p.Process.Cmd,
		}
	}
	return nil
}

func pluginWASM(p config.Plugin) *plugin.Codegen_WASM {
	if p.WASM != nil {
		return &plugin.Codegen_WASM{
			Url:    p.WASM.URL,
			Sha256: p.WASM.SHA256,
		}
	}
	return nil
}

func pluginCatalog(c *catalog.Catalog) *plugin.Catalog {
	var schemas []*plugin.Schema
	for _, s := range c.Schemas {
		var enums []*plugin.Enum
		var cts []*plugin.CompositeType
		for _, typ := range s.Types {
			switch typ := typ.(type) {
			case *catalog.Enum:
				enums = append(enums, &plugin.Enum{
					Name:    typ.Name,
					Comment: typ.Comment,
					Vals:    typ.Vals,
				})
			case *catalog.CompositeType:
				cts = append(cts, &plugin.CompositeType{
					Name:    typ.Name,
					Comment: typ.Comment,
				})
			}
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
					Comment:   c.Comment,
					NotNull:   c.IsNotNull,
					Unsigned:  c.IsUnsigned,
					IsArray:   c.IsArray,
					ArrayDims: int32(c.ArrayDims),
					Length:    int32(l),
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
			Comment:        s.Comment,
			Name:           s.Name,
			Tables:         tables,
			Enums:          enums,
			CompositeTypes: cts,
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
		var iit *plugin.Identifier
		if q.InsertIntoTable != nil {
			iit = &plugin.Identifier{
				Catalog: q.InsertIntoTable.Catalog,
				Schema:  q.InsertIntoTable.Schema,
				Name:    q.InsertIntoTable.Name,
			}
		}
		out = append(out, &plugin.Query{
			Name:            q.Metadata.Name,
			Cmd:             q.Metadata.Cmd,
			Text:            q.SQL,
			Comments:        q.Metadata.Comments,
			Columns:         columns,
			Params:          params,
			Filename:        q.Metadata.Filename,
			InsertIntoTable: iit,
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
		Name:         c.Name,
		OriginalName: c.OriginalName,
		Comment:      c.Comment,
		NotNull:      c.NotNull,
		Unsigned:     c.Unsigned,
		IsArray:      c.IsArray,
		ArrayDims:    int32(c.ArrayDims),
		Length:       int32(l),
		IsNamedParam: c.IsNamedParam,
		IsFuncCall:   c.IsFuncCall,
		IsSqlcSlice:  c.IsSqlcSlice,
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

	if c.EmbedTable != nil {
		out.EmbedTable = &plugin.Identifier{
			Catalog: c.EmbedTable.Catalog,
			Schema:  c.EmbedTable.Schema,
			Name:    c.EmbedTable.Name,
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

func codeGenRequest(r *compiler.Result, settings config.CombinedSettings) *plugin.GenerateRequest {
	// Engines on the xqlc core (ClickHouse) project the codegen catalog
	// from the core catalog rather than the in-memory sql/catalog, which is
	// nil on that path.
	var cat *plugin.Catalog
	if r.CoreCatalog != nil {
		cat = pluginCatalogFromCore(r.CoreCatalog)
	} else {
		cat = pluginCatalog(r.Catalog)
	}
	return &plugin.GenerateRequest{
		Settings:    pluginSettings(r, settings),
		Catalog:     cat,
		Queries:     pluginQueries(r),
		SqlcVersion: info.Version,
	}
}

// pluginCatalogFromCore projects a core.Catalog (the xqlc SQLite-backed
// catalog) into the plugin.Catalog that codegen consumes to emit models
// and enums. It reads the namespace / class / attribute / type tables
// directly. Projection is best-effort: an unexpected query error against
// the in-memory catalog yields a partial catalog rather than aborting
// generation.
func pluginCatalogFromCore(cc *core.Catalog) *plugin.Catalog {
	db := cc.DB()
	var schemas []*plugin.Schema

	type row struct {
		oid  int64
		name string
	}
	readRows := func(query string, args ...any) []row {
		rows, err := db.Query(query, args...)
		if err != nil {
			return nil
		}
		defer rows.Close()
		var out []row
		for rows.Next() {
			var r row
			if err := rows.Scan(&r.oid, &r.name); err != nil {
				return out
			}
			out = append(out, r)
		}
		return out
	}

	for _, ns := range readRows(`SELECT oid, name FROM sql_namespace ORDER BY oid`) {
		var tables []*plugin.Table
		for _, cl := range readRows(
			`SELECT oid, name FROM sql_class WHERE namespace_oid = ? AND kind = 'r' ORDER BY oid`, ns.oid,
		) {
			rel := &plugin.Identifier{Schema: ns.name, Name: cl.name}
			var columns []*plugin.Column
			crows, err := db.Query(
				`SELECT a.name, t.name, a.not_null
				 FROM sql_attribute a JOIN sql_type t ON t.oid = a.type_oid
				 WHERE a.class_oid = ? ORDER BY a.num`, cl.oid,
			)
			if err != nil {
				continue
			}
			for crows.Next() {
				var name, typeName string
				var notNull int
				if err := crows.Scan(&name, &typeName, &notNull); err != nil {
					break
				}
				columns = append(columns, &plugin.Column{
					Name:    name,
					Type:    &plugin.Identifier{Name: typeName},
					NotNull: notNull != 0,
					Table:   rel,
				})
			}
			crows.Close()
			tables = append(tables, &plugin.Table{Rel: rel, Columns: columns})
		}
		schemas = append(schemas, &plugin.Schema{Name: ns.name, Tables: tables})
	}

	return &plugin.Catalog{
		DefaultSchema: "public",
		Schemas:       schemas,
	}
}
