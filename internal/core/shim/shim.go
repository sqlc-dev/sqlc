package shim

import (
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func PluginCatalog(cc *core.Catalog) *plugin.Catalog {
	var schemas []*plugin.Schema

	namespaces, err := cc.Namespaces()
	if err != nil {
		return &plugin.Catalog{DefaultSchema: "public"}
	}
	for _, ns := range namespaces {
		tables, err := cc.TablesInNamespace(ns.OID)
		if err != nil {
			continue
		}
		var ptables []*plugin.Table
		for _, cl := range tables {
			rel := &plugin.Identifier{Schema: ns.Name, Name: cl.Name}
			cols, err := cc.ClassCodegenColumns(cl.OID)
			if err != nil {
				continue
			}
			var columns []*plugin.Column
			for _, col := range cols {
				columns = append(columns, &plugin.Column{
					Name:    col.Name,
					Type:    &plugin.Identifier{Name: col.TypeName},
					NotNull: col.NotNull,
					Table:   rel,
				})
			}
			ptables = append(ptables, &plugin.Table{Rel: rel, Columns: columns})
		}
		schemas = append(schemas, &plugin.Schema{Name: ns.Name, Tables: ptables})
	}

	return &plugin.Catalog{
		DefaultSchema: "public",
		Schemas:       schemas,
	}
}
