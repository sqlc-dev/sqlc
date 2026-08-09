package sqlite

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// defaultSchema is SQLite's standard library, read from the dialect
// directory's functions.jsonl.
//
// The functions are drawn from:
//
//	https://www.sqlite.org/lang_aggfunc.html
//	https://www.sqlite.org/lang_mathfunc.html
//	https://www.sqlite.org/lang_corefunc.html
func defaultSchema(name string) *catalog.Schema {
	return &catalog.Schema{Name: name, Funcs: stdlib()}
}
