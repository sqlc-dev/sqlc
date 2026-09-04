package sqlite

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// defaultSchema is SQLite's standard library, read from the dialect
// directory's functions.jsonl, which internal/goldeneye generates from a
// default build of SQLite. The functions further compile options add live
// under the dialect's extensions/, one directory per option.
func defaultSchema(name string) *catalog.Schema {
	return &catalog.Schema{Name: name, Funcs: stdlib()}
}
