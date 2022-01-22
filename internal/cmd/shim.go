package cmd

import (
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/plugin"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func pluginSettings(settings config.CombinedSettings) *plugin.Settings {
	return nil
}

func pluginCatalog(c *catalog.Catalog) *plugin.Catalog {
	return nil
}

func pluginQueries(r *compiler.Result) []*plugin.Query {
	var out []*plugin.Query
	for _, q := range r.Queries {
		out = append(out, &plugin.Query{
			Text: q.SQL,
		})
	}
	return out
}

func codeGenRequest(r *compiler.Result, settings config.CombinedSettings) *plugin.CodeGenRequest {
	return &plugin.CodeGenRequest{
		Settings: pluginSettings(settings),
		Catalog:  pluginCatalog(r.Catalog),
		Queries:  pluginQueries(r),
	}
}
