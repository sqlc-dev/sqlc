package sqlite

import (
	"github.com/sqlc-dev/sqlc/internal/engine"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// sqliteEngine implements the engine.Engine interface for SQLite.
type sqliteEngine struct {
	parser *Parser
}

// NewEngine creates a new SQLite engine.
func NewEngine() engine.Engine {
	return &sqliteEngine{
		parser: NewParser(),
	}
}

// Name returns the engine name.
func (e *sqliteEngine) Name() string {
	return "sqlite"
}

// Parser returns the SQLite parser.
func (e *sqliteEngine) Parser() engine.Parser {
	return e.parser
}

// Catalog returns a new SQLite catalog.
func (e *sqliteEngine) Catalog() *catalog.Catalog {
	return NewCatalog()
}

// Selector returns a SQLite-specific selector for handling jsonb columns.
func (e *sqliteEngine) Selector() engine.Selector {
	return &sqliteSelector{}
}

// Dialect returns the parser which implements the Dialect interface.
func (e *sqliteEngine) Dialect() engine.Dialect {
	return e.parser
}

// sqliteSelector wraps jsonb columns with json() for proper output.
type sqliteSelector struct{}

// ColumnExpr wraps jsonb columns with json() function.
func (s *sqliteSelector) ColumnExpr(name string, dataType string) string {
	// Under SQLite, neither json nor jsonb are real data types, and rather just
	// of type blob, so database drivers just return whatever raw binary is
	// stored as values. This is a problem for jsonb, which is considered an
	// internal format to SQLite and no attempt should be made to parse it
	// outside of the database itself. For jsonb columns in SQLite, wrap values
	// in `json(col)` to coerce the internal binary format to JSON parsable by
	// the user-space application.
	if dataType == "jsonb" {
		return "json(" + name + ")"
	}
	return name
}
