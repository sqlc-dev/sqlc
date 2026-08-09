package postgresql

import (
	"io/fs"
	"path"

	"github.com/sqlc-dev/sqlc/internal/core/seed"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// loadExtension reads the named extension's function list from the dialect
// directory, where sqlc-pg-gen writes one directory per extension. An
// extension sqlc does not know is nil, which CREATE EXTENSION treats as
// nothing to add.
func loadExtension(name string) *catalog.Schema {
	dir := path.Join("dialect", "extensions", name)
	if _, err := fs.Stat(dialectFS, path.Join(dir, seed.FunctionsFile)); err != nil {
		return nil
	}
	funcs, err := seed.Functions(dialectFS, dir)
	if err != nil {
		// The list is embedded in the binary: a failure here means sqlc was
		// built from a broken tree, which no caller can do anything about.
		panic(err)
	}
	return &catalog.Schema{Name: "pg_catalog", Funcs: funcs}
}
