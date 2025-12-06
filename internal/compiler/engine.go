package compiler

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/analyzer"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	pganalyze "github.com/sqlc-dev/sqlc/internal/engine/postgresql/analyzer"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	sqliteanalyze "github.com/sqlc-dev/sqlc/internal/engine/sqlite/analyzer"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

type ResolveTypeFunc func(call *ast.FuncCall, fun *catalog.Function, resolve func(n ast.Node) (*catalog.Column, error)) *ast.TypeName

type Compiler struct {
	conf     config.SQL
	combo    config.CombinedSettings
	catalog  *catalog.Catalog
	parser   Parser
	result   *Result
	analyzer analyzer.Analyzer
	client   dbmanager.Client
	selector selector

	schema       []string
	TypeResolver ResolveTypeFunc
}

func NewCompiler(conf config.SQL, combo config.CombinedSettings) (*Compiler, error) {
	c := &Compiler{conf: conf, combo: combo}

	if conf.Database != nil && conf.Database.Managed {
		client := dbmanager.NewClient(combo.Global.Servers)
		c.client = client
	}

	switch conf.Engine {
	case config.EngineClickHouse:
		c.parser = clickhouse.NewParser()
		c.catalog = clickhouse.NewCatalog()
		c.selector = newDefaultSelector()
		c.TypeResolver = clickhouse.TypeResolver
	case config.EngineSQLite:
		c.parser = sqlite.NewParser()
		c.catalog = sqlite.NewCatalog()
		c.selector = newSQLiteSelector()
		if conf.Database != nil {
			if conf.Analyzer.Database == nil || *conf.Analyzer.Database {
				c.analyzer = analyzer.Cached(
					sqliteanalyze.New(*conf.Database),
					combo.Global,
					*conf.Database,
				)
			}
		}
	case config.EngineMySQL:
		c.parser = dolphin.NewParser()
		c.catalog = dolphin.NewCatalog()
		c.selector = newDefaultSelector()
	case config.EnginePostgreSQL:
		c.parser = postgresql.NewParser()
		c.catalog = postgresql.NewCatalog()
		c.selector = newDefaultSelector()
		if conf.Database != nil {
			if conf.Analyzer.Database == nil || *conf.Analyzer.Database {
				c.analyzer = analyzer.Cached(
					pganalyze.New(c.client, *conf.Database),
					combo.Global,
					*conf.Database,
				)
			}
		}
	default:
		return nil, fmt.Errorf("unknown engine: %s", conf.Engine)
	}
	return c, nil
}

func (c *Compiler) Catalog() *catalog.Catalog {
	return c.catalog
}

func (c *Compiler) ParseCatalog(schema []string) error {
	err := c.parseCatalog(schema)
	if err == nil && c.conf.Engine == config.EngineClickHouse {
		// Set the catalog on the ClickHouse parser so it can register
		// context-dependent functions during query parsing
		if chParser, ok := c.parser.(*clickhouse.Parser); ok {
			chParser.Catalog = c.catalog
		}
	}
	return err
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

func (c *Compiler) Close(ctx context.Context) {
	if c.analyzer != nil {
		c.analyzer.Close(ctx)
	}
	if c.client != nil {
		c.client.Close(ctx)
	}
}
