package compiler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	coreschema "github.com/sqlc-dev/sqlc/internal/core/schema"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/rpc"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/preprocess"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

// TODO: Rename this interface Engine
type Parser interface {
	Parse(io.Reader) ([]ast.Statement, error)
	CommentSyntax() source.CommentSyntax
	IsReservedKeyword(string) bool
}

func (c *Compiler) parseCatalog(schemas []string) error {
	files, err := sqlpath.Glob(schemas)
	if err != nil {
		return err
	}
	merr := multierr.New()
	for _, filename := range files {
		blob, err := os.ReadFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		contents := migrations.RemoveRollbackStatements(string(blob))
		contents = migrations.RemovePsqlMetaCommands(contents)
		c.schema = append(c.schema, contents)

		// In database-only mode, we parse the schema to validate syntax
		// but don't update the catalog - the database will be the source of truth
		stmts, err := c.parser.Parse(strings.NewReader(contents))
		if err != nil {
			// A schema file and a query file are often the same file, so a
			// query's sqlc syntax can fail here. Look for an explanation
			// before reporting the syntax error it caused.
			if reported := addSyntaxErrors(merr, filename, contents, preprocess.File(c.conf.Engine, contents)); !reported {
				merr.Add(filename, contents, 0, err)
			}
			continue
		}

		// Skip catalog updates in database-only mode
		if c.databaseOnlyMode {
			continue
		}

		if c.coreCatalog != nil {
			for i := range stmts {
				if err := coreschema.Apply(c.coreCatalog, stmts[i].Raw); err != nil {
					merr.Add(filename, contents, stmts[i].Pos(), err)
					continue
				}
			}
			continue
		}

		for i := range stmts {
			if err := c.catalog.Update(stmts[i], c); err != nil {
				merr.Add(filename, contents, stmts[i].Pos(), err)
				continue
			}
		}
	}
	if len(merr.Errs()) > 0 {
		return merr
	}
	return nil
}

func (c *Compiler) parseQueries(o opts.Parser) (*Result, error) {
	ctx := context.Background()

	// In database-only mode, initialize the database connection before parsing queries
	if c.databaseOnlyMode && c.analyzer != nil {
		if err := c.analyzer.EnsureConn(ctx, c.schema); err != nil {
			return nil, fmt.Errorf("failed to initialize database connection: %w", err)
		}
	}

	var q []*Query
	merr := multierr.New()
	set := map[string]struct{}{}
	files, err := sqlpath.Glob(c.conf.Queries)
	if err != nil {
		return nil, err
	}
	for _, filename := range files {
		blob, err := os.ReadFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		src := string(blob)

		// sqlc syntax is not SQL. Rewrite it to native placeholders before the
		// engine parser ever sees the query, so parsers only handle SQL.
		pp := preprocess.File(c.conf.Engine, src)

		stmts, err := c.parser.Parse(strings.NewReader(pp.Text))
		if err != nil {
			if reported := addSyntaxErrors(merr, filename, src, pp); !reported {
				merr.Add(filename, src, 0, err)
			}
			continue
		}
		for _, stmt := range stmts {
			query, err := c.parseQuery(stmt.Raw, pp, o)
			if err != nil {
				var e *sqlerr.Error
				loc := stmt.Raw.Pos()
				if errors.As(err, &e) && e.Location != 0 {
					loc = e.Location
				}
				// Locations are reported against the rewritten query; map them
				// back so errors point at what the user wrote.
				loc = pp.Origin(loc)
				if e != nil && e.Location != 0 {
					e.Location = loc
				}
				merr.Add(filename, src, loc, err)
				// If this rpc unauthenticated error bubbles up, then all future parsing/analysis will fail
				if errors.Is(err, rpc.ErrUnauthenticated) {
					return nil, merr
				}
				continue
			}
			if query == nil {
				continue
			}
			query.Metadata.Filename = filepath.Base(filename)
			queryName := query.Metadata.Name
			if queryName != "" {
				if _, exists := set[queryName]; exists {
					merr.Add(filename, src, pp.Origin(stmt.Raw.Pos()), fmt.Errorf("duplicate query name: %s", queryName))
					continue
				}
				set[queryName] = struct{}{}
			}
			q = append(q, query)
		}
	}
	if len(merr.Errs()) > 0 {
		return nil, merr
	}
	if len(q) == 0 {
		return nil, fmt.Errorf("no queries contained in paths %s", strings.Join(c.conf.Queries, ","))
	}

	return &Result{
		Catalog: c.catalog,
		Queries: q,
	}, nil
}

// addSyntaxErrors reports every sqlc syntax error the preprocessor recorded
// for a file, in source order, and says whether there were any. Locations come
// back in the rewritten text's coordinates, so they are mapped through Origin
// to point at what the user wrote.
//
// A statement whose sqlc syntax did not validate is copied through for the
// engine to parse, which assumes the engine can parse it. SQLite cannot: it
// has no schema-qualified function call, so a bad sqlc.arg() is a syntax error
// there rather than a call the preprocessor's message can be attached to.
// These messages name the cause, so they are reported in place of the failure
// they produced; anything else wrong with the file surfaces on the next run.
func addSyntaxErrors(merr *multierr.Error, filename, src string, pp *preprocess.Result) bool {
	var found bool
	for _, stmt := range pp.Statements() {
		if stmt.Err == nil {
			continue
		}
		found = true
		loc := pp.Origin(stmt.Start)
		var e *sqlerr.Error
		if errors.As(stmt.Err, &e) && e.Location != 0 {
			loc = pp.Origin(e.Location)
			e.Location = loc
		}
		merr.Add(filename, src, loc, stmt.Err)
	}
	return found
}
