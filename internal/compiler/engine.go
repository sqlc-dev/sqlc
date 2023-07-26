package compiler

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

type Compiler struct {
	conf    config.SQL
	combo   config.CombinedSettings
	catalog *catalog.Catalog
	parser  Parser
	result  *Result
}

func NewCompiler(conf config.SQL, combo config.CombinedSettings) *Compiler {
	c := &Compiler{conf: conf, combo: combo}
	switch conf.Engine {
	case config.EngineSQLite:
		c.parser = sqlite.NewParser()
		c.catalog = sqlite.NewCatalog()
	case config.EngineMySQL:
		c.parser = dolphin.NewParser()
		c.catalog = dolphin.NewCatalog()
	case config.EnginePostgreSQL:
		c.parser = postgresql.NewParser()
		c.catalog = postgresql.NewCatalog()
	default:
		panic(fmt.Sprintf("unknown engine: %s", conf.Engine))
	}
	return c
}

func (c *Compiler) Catalog() *catalog.Catalog {
	return c.catalog
}

func (c *Compiler) ParseCatalog(schema []string) error {
	return c.parseCatalog(schema)
}

func (c *Compiler) ParseQueries(queries []string, o opts.Parser) error {
	r, err := c.parseQueries(o)
	if err != nil {
		return err
	}
	c.result = r
	return nil
}

func (c *Compiler) Result() *Result {
	return c.result
}
