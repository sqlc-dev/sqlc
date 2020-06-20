package compiler

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/engine/dolphin"
	"github.com/kyleconroy/sqlc/internal/engine/postgresql"
	"github.com/kyleconroy/sqlc/internal/engine/sqlite"
	"github.com/kyleconroy/sqlc/internal/opts"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

// The Engine type only exists as a compatibility shim between the old dinosql
// package and the new compiler package.
type Engine struct {
	conf    config.SQL
	combo   config.CombinedSettings
	catalog *catalog.Catalog
	parser  Parser
	result  *Result
}

func NewEngine(conf config.SQL, combo config.CombinedSettings) *Engine {
	e := &Engine{conf: conf, combo: combo}
	switch conf.Engine {
	case config.EngineXLemon:
		e.parser = sqlite.NewParser()
		e.catalog = catalog.New("main")
	case config.EngineMySQL, config.EngineXDolphin:
		e.parser = dolphin.NewParser()
		e.catalog = catalog.New("public") // TODO: What is the default database for MySQL?
	case config.EnginePostgreSQL:
		e.parser = postgresql.NewParser()
		e.catalog = postgresql.NewCatalog()
	default:
		panic(fmt.Sprintf("unknown engine: %s", conf.Engine))
	}
	return e
}

func (e *Engine) ParseCatalog(schema []string) error {
	return parseCatalog(e.parser, e.catalog, schema)
}

func (e *Engine) ParseQueries(queries []string, o opts.Parser) error {
	r, err := parseQueries(e.parser, e.catalog, e.conf.Queries, o)
	if err != nil {
		return err
	}
	e.result = r
	return nil
}

func (e *Engine) Result() *Result {
	return e.result
}
