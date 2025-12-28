package postgresql

import (
	"github.com/sqlc-dev/sqlc/internal/engine"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// postgresqlEngine implements the engine.Engine interface for PostgreSQL.
type postgresqlEngine struct {
	parser *Parser
}

// NewEngine creates a new PostgreSQL engine.
func NewEngine() engine.Engine {
	return &postgresqlEngine{
		parser: NewParser(),
	}
}

// Name returns the engine name.
func (e *postgresqlEngine) Name() string {
	return "postgresql"
}

// Parser returns the PostgreSQL parser.
func (e *postgresqlEngine) Parser() engine.Parser {
	return e.parser
}

// Catalog returns a new PostgreSQL catalog.
func (e *postgresqlEngine) Catalog() *catalog.Catalog {
	return NewCatalog()
}

// Selector returns nil because PostgreSQL uses the default selector.
func (e *postgresqlEngine) Selector() engine.Selector {
	return &engine.DefaultSelector{}
}

// Dialect returns the parser which implements the Dialect interface.
func (e *postgresqlEngine) Dialect() engine.Dialect {
	return e.parser
}
