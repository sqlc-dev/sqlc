package compiler

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/analyzer"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/googlesql"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	pganalyze "github.com/sqlc-dev/sqlc/internal/engine/postgresql/analyzer"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	sqliteanalyze "github.com/sqlc-dev/sqlc/internal/engine/sqlite/analyzer"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

type Compiler struct {
	conf     config.SQL
	combo    config.CombinedSettings
	catalog  *catalog.Catalog
	parser   Parser
	result   *Result
	analyzer analyzer.Analyzer
	client   dbmanager.Client
	selector selector

	coreCatalog *core.Catalog

	// coreAnalysis routes every engine through the core catalog and analyzer,
	// and coreDialect seeds the catalog it does so with. The catalog itself is
	// not built until the schema is loaded, because the two are cached
	// together.
	coreAnalysis bool
	coreDialect  core.Option

	schema []string
}

// Option configures a Compiler.
type Option func(*Compiler)

// WithCoreAnalysis analyzes queries with the core catalog and analyzer rather
// than each engine's own analysis path. ClickHouse and GoogleSQL always do;
// this opts the remaining engines in.
func WithCoreAnalysis() Option {
	return func(c *Compiler) { c.coreAnalysis = true }
}

func NewCompiler(conf config.SQL, combo config.CombinedSettings, parserOpts opts.Parser, options ...Option) (*Compiler, error) {
	c := &Compiler{conf: conf, combo: combo}
	for _, o := range options {
		o(c)
	}

	// ClickHouse and GoogleSQL have no legacy analysis path to fall back to.
	switch conf.Engine {
	case config.EngineClickHouse, config.EngineGoogleSQL:
		c.coreAnalysis = true
	}
	if c.coreAnalysis {
		if err := c.initCore(); err != nil {
			return nil, err
		}
		return c, nil
	}

	if conf.Database != nil && conf.Database.Managed {
		client := dbmanager.NewClient(combo.Global.Servers)
		c.client = client
	}

	switch conf.Engine {
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

// initCore wires up the engine's parser and a core catalog seeded with its
// dialect. No legacy catalog and no analyzer connection are involved.
func (c *Compiler) initCore() error {
	var dialect core.Option
	switch c.conf.Engine {
	case config.EngineSQLite:
		c.parser = sqlite.NewParser()
		c.selector = newSQLiteSelector()
		dialect = sqlite.Dialect()
	case config.EngineMySQL:
		c.parser = dolphin.NewParser()
		c.selector = newDefaultSelector()
		dialect = dolphin.Dialect()
	case config.EnginePostgreSQL:
		c.parser = postgresql.NewParser()
		c.selector = newDefaultSelector()
		dialect = postgresql.Dialect()
	case config.EngineClickHouse:
		c.parser = clickhouse.NewParser()
		c.selector = newDefaultSelector()
		dialect = clickhouse.Dialect()
	case config.EngineGoogleSQL:
		c.parser = googlesql.NewParser()
		c.selector = newDefaultSelector()
		dialect = googlesql.Dialect()
	default:
		return fmt.Errorf("unknown engine: %s", c.conf.Engine)
	}
	c.coreDialect = dialect
	return nil
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

func (c *Compiler) Close(ctx context.Context) {
	if c.analyzer != nil {
		c.analyzer.Close(ctx)
	}
	if c.client != nil {
		c.client.Close(ctx)
	}
	if c.coreCatalog != nil {
		c.coreCatalog.Close()
	}
}
