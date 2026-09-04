package sqlite

import (
	"github.com/sqlc-dev/sqlc/internal/core/seed"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// loadExtension returns the functions a compile option adds, read from the
// option's directory under the dialect's extensions/. The name is either the
// directory's — enable_fts5 — or a virtual table module the dialect maps to
// one, since a schema that says CREATE VIRTUAL TABLE ... USING fts5 has said
// which build of SQLite it runs on. An option the dialect has no data for
// adds nothing.
func loadExtension(name string) *catalog.Schema {
	dir, ok := seed.ExtensionDir(dialectFS, "dialect", name)
	if !ok {
		return nil
	}
	funcs, err := seed.Functions(dialectFS, dir)
	if err != nil {
		// The list is embedded in the binary: a failure here means sqlc was
		// built from a broken tree, which no caller can do anything about.
		panic(err)
	}
	return &catalog.Schema{Name: "main", Funcs: funcs}
}
