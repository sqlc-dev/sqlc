package compiler

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/analyzer"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	pganalyze "github.com/sqlc-dev/sqlc/internal/engine/postgresql/analyzer"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	sqliteanalyze "github.com/sqlc-dev/sqlc/internal/engine/sqlite/analyzer"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/x/expander"
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

	schema []string

	// accurateMode indicates that the compiler should use database-only analysis
	// and skip building the internal catalog from schema files
	accurateMode bool
	// pgAnalyzer is the PostgreSQL-specific analyzer used in accurate mode
	// for schema introspection
	pgAnalyzer *pganalyze.Analyzer
	// expander is used to expand SELECT * and RETURNING * in accurate mode
	expander *expander.Expander
}

func NewCompiler(conf config.SQL, combo config.CombinedSettings) (*Compiler, error) {
	c := &Compiler{conf: conf, combo: combo}

	if conf.Database != nil && conf.Database.Managed {
		client := dbmanager.NewClient(combo.Global.Servers)
		c.client = client
	}

	// Check for accurate mode
	accurateMode := conf.Analyzer.Accurate != nil && *conf.Analyzer.Accurate

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
		parser := postgresql.NewParser()
		c.parser = parser
		c.catalog = postgresql.NewCatalog()
		c.selector = newDefaultSelector()

		if accurateMode {
			// Accurate mode requires a database connection
			if conf.Database == nil {
				return nil, fmt.Errorf("accurate mode requires database configuration")
			}
			if conf.Database.URI == "" && !conf.Database.Managed {
				return nil, fmt.Errorf("accurate mode requires database.uri or database.managed")
			}
			c.accurateMode = true
			// Create the PostgreSQL analyzer for schema introspection
			c.pgAnalyzer = pganalyze.New(c.client, *conf.Database)
			// Use the analyzer wrapped with cache for query analysis
			c.analyzer = analyzer.Cached(
				c.pgAnalyzer,
				combo.Global,
				*conf.Database,
			)
			// Create the expander using the pgAnalyzer as the column getter
			// The parser implements both Parser and format.Dialect interfaces
			c.expander = expander.New(c.pgAnalyzer, parser, parser)
		} else if conf.Database != nil {
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
}
