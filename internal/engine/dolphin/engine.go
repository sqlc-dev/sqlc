package dolphin

import (
	"github.com/sqlc-dev/sqlc/internal/engine"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// dolphinEngine implements the engine.Engine interface for MySQL.
type dolphinEngine struct {
	parser *Parser
}

// NewEngine creates a new MySQL engine.
func NewEngine() engine.Engine {
	return &dolphinEngine{
		parser: NewParser(),
	}
}

// Name returns the engine name.
func (e *dolphinEngine) Name() string {
	return "mysql"
}

// Parser returns the MySQL parser.
func (e *dolphinEngine) Parser() engine.Parser {
	return e.parser
}

// Catalog returns a new MySQL catalog.
func (e *dolphinEngine) Catalog() *catalog.Catalog {
	return NewCatalog()
}

// Selector returns nil because MySQL uses the default selector.
func (e *dolphinEngine) Selector() engine.Selector {
	return &engine.DefaultSelector{}
}

// Dialect returns the parser which implements the Dialect interface.
func (e *dolphinEngine) Dialect() engine.Dialect {
	return e.parser
}
