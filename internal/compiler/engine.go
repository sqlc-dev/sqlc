// +build exp

package compiler

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/dolphin"
	"github.com/kyleconroy/sqlc/internal/postgresql"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sqlite"
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
	case config.EngineXDolphin:
		e.parser = dolphin.NewParser()
		e.catalog = catalog.New("public") // TODO: What is the default database for MySQL?
	case config.EngineXElephant, config.EnginePostgreSQL:
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

func (e *Engine) ParseQueries(queries []string, opts dinosql.ParserOpts) error {
	r, err := parseQueries(e.parser, e.catalog, e.conf.Queries)
	if err != nil {
		return err
	}
	e.result = r
	return nil
}

func (e *Engine) Result() dinosql.Generateable {
	return e.result
}
