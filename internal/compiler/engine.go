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

	schema []string
}

func NewCompiler(conf config.SQL, combo config.CombinedSettings) (*Compiler, error) {
	c := &Compiler{conf: conf, combo: combo}

	if conf.Database != nil && conf.Database.Managed {
		client := dbmanager.NewClient(combo.Global.Servers)
		c.client = client
	}

	// Check if skip_parser is enabled
	skipParser := conf.Analyzer.SkipParser != nil && *conf.Analyzer.SkipParser

	// If skip_parser is enabled, we must have database analyzer enabled
	if skipParser {
		if conf.Database == nil {
			return nil, fmt.Errorf("skip_parser requires database configuration")
		}
		if conf.Analyzer.Database != nil && !*conf.Analyzer.Database {
			return nil, fmt.Errorf("skip_parser requires database analyzer to be enabled")
		}
		// Only PostgreSQL is supported for now
		if conf.Engine != config.EnginePostgreSQL {
			return nil, fmt.Errorf("skip_parser is only supported for PostgreSQL")
		}
	}

	switch conf.Engine {
	case config.EngineSQLite:
		c.parser = sqlite.NewParser()
		c.catalog = sqlite.NewCatalog()
		c.selector = newSQLiteSelector()
	case config.EngineMySQL:
		c.parser = dolphin.NewParser()
		c.catalog = dolphin.NewCatalog()
		c.selector = newDefaultSelector()
	case config.EnginePostgreSQL:
		// Skip parser and catalog if skip_parser is enabled
		if !skipParser {
			c.parser = postgresql.NewParser()
			c.catalog = postgresql.NewCatalog()
		}
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
	return c.parseCatalog(schema)
}

func (c *Compiler) ParseQueries(queries []string, o opts.Parser) error {
	// Check if skip_parser is enabled
	skipParser := c.conf.Analyzer.SkipParser != nil && *c.conf.Analyzer.SkipParser

	var r *Result
	var err error

	if skipParser {
		// Use database analyzer only, skip parser and catalog
		r, err = c.parseQueriesWithAnalyzer(o)
	} else {
		// Use traditional parser-based approach
		r, err = c.parseQueries(o)
	}

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
