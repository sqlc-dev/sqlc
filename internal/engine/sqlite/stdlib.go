package sqlite

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// defaultSchema is SQLite's standard library, read from the dialect
// directory's functions.jsonl, which internal/goldeneye generates from the
// sqlite3 shell's pragma_function_list.
func defaultSchema(name string) *catalog.Schema {
	return &catalog.Schema{Name: name, Funcs: stdlib()}
}
