package compiler

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/analyzer"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/engine"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/plugin"
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

	// databaseOnlyMode indicates that the compiler should use database-only analysis
	// and skip building the internal catalog from schema files (analyzer.database: only)
	databaseOnlyMode bool
	// expander is used to expand SELECT * and RETURNING * in database-only mode
	expander *expander.Expander
}

func NewCompiler(conf config.SQL, combo config.CombinedSettings, parserOpts opts.Parser) (*Compiler, error) {
	c := &Compiler{conf: conf, combo: combo}

	if conf.Database != nil && conf.Database.Managed {
		client := dbmanager.NewClient(combo.Global.Servers)
		c.client = client
	}

	// Check for database-only mode (analyzer.database: only)
	// This feature requires the analyzerv2 experiment to be enabled
	databaseOnlyMode := conf.Analyzer.Database.IsOnly() && parserOpts.Experiment.AnalyzerV2

	switch conf.Engine {
	case config.EngineSQLite:
		parser := sqlite.NewParser()
		c.parser = parser
		c.catalog = sqlite.NewCatalog()
		c.selector = newSQLiteSelector()

		if databaseOnlyMode {
			// Database-only mode requires a database connection
			if conf.Database == nil {
				return nil, fmt.Errorf("analyzer.database: only requires database configuration")
			}
			if conf.Database.URI == "" && !conf.Database.Managed {
				return nil, fmt.Errorf("analyzer.database: only requires database.uri or database.managed")
			}
			c.databaseOnlyMode = true
			// Create the SQLite analyzer (implements Analyzer interface)
			sqliteAnalyzer := sqliteanalyze.New(*conf.Database)
			c.analyzer = analyzer.Cached(sqliteAnalyzer, combo.Global, *conf.Database)
			// Create the expander using the analyzer as the column getter
			c.expander = expander.New(c.analyzer, parser, parser)
		} else if conf.Database != nil {
			if conf.Analyzer.Database.IsEnabled() {
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

		if databaseOnlyMode {
			// Database-only mode requires a database connection
			if conf.Database == nil {
				return nil, fmt.Errorf("analyzer.database: only requires database configuration")
			}
			if conf.Database.URI == "" && !conf.Database.Managed {
				return nil, fmt.Errorf("analyzer.database: only requires database.uri or database.managed")
			}
			c.databaseOnlyMode = true
			// Create the PostgreSQL analyzer (implements Analyzer interface)
			pgAnalyzer := pganalyze.New(c.client, *conf.Database)
			c.analyzer = analyzer.Cached(pgAnalyzer, combo.Global, *conf.Database)
			// Create the expander using the analyzer as the column getter
			c.expander = expander.New(c.analyzer, parser, parser)
		} else if conf.Database != nil {
			if conf.Analyzer.Database.IsEnabled() {
				c.analyzer = analyzer.Cached(
					pganalyze.New(c.client, *conf.Database),
					combo.Global,
					*conf.Database,
				)
			}
		}
	default:
		// Check if this is a plugin engine
		if enginePlugin, found := config.FindEnginePlugin(&combo.Global, string(conf.Engine)); found {
			eng, err := createPluginEngine(enginePlugin, combo.Dir)
			if err != nil {
				return nil, err
			}
			c.parser = eng.Parser()
			c.catalog = eng.Catalog()
			sel := eng.Selector()
			if sel != nil {
				c.selector = &engineSelectorAdapter{sel}
			} else {
				c.selector = newDefaultSelector()
			}
		} else {
			return nil, fmt.Errorf("unknown engine: %s\n\nTo use a custom database engine, add it to the 'engines' section of sqlc.yaml:\n\n  engines:\n    - name: %s\n      process:\n        cmd: sqlc-engine-%s\n\nThen install the plugin: go install github.com/example/sqlc-engine-%s@latest",
				conf.Engine, conf.Engine, conf.Engine, conf.Engine)
		}
	}
	return c, nil
}

// createPluginEngine creates an engine from an engine plugin configuration.
func createPluginEngine(ep *config.EnginePlugin, dir string) (engine.Engine, error) {
	switch {
	case ep.Process != nil:
		return plugin.NewPluginEngine(ep.Name, ep.Process.Cmd, dir, ep.Env), nil
	case ep.WASM != nil:
		return plugin.NewWASMPluginEngine(ep.Name, ep.WASM.URL, ep.WASM.SHA256, ep.Env), nil
	default:
		return nil, fmt.Errorf("engine plugin %s has no process or wasm configuration", ep.Name)
	}
}

// engineSelectorAdapter adapts engine.Selector to the compiler's selector interface.
type engineSelectorAdapter struct {
	sel engine.Selector
}

func (a *engineSelectorAdapter) ColumnExpr(name string, column *Column) string {
	return a.sel.ColumnExpr(name, column.DataType)
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
