package compiler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/sqlc-dev/sqlc/internal/core"
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

// schemaFile is a schema file's contents, after the rewrites that are applied
// before any of it is parsed.
type schemaFile struct {
	name     string
	contents string
}

// queryFile is a query file as it was read.
type queryFile struct {
	name string
	src  string
}

// querySet is the query files and the macros they define. A macro can be
// called from any file, so every file is read, and its macros collected,
// before the first one is rewritten.
type querySet struct {
	files  []queryFile
	macros *preprocess.Macros

	// globErr is why the files could not be listed, and errs is what went
	// wrong reading them or collecting their macros.
	globErr error
	errs    *multierr.Error
}

// loadQueries reads the query files, once. The schema is parsed before the
// queries are, and it needs their macros: a schema file is often a query file
// too.
func (c *Compiler) loadQueries() *querySet {
	if c.queries != nil {
		return c.queries
	}
	set := &querySet{
		macros: preprocess.NewMacros(c.conf.Engine),
		errs:   multierr.New(),
	}
	c.queries = set

	files, err := sqlpath.Glob(c.conf.Queries)
	if err != nil {
		set.globErr = err
		return set
	}
	for _, filename := range files {
		blob, err := os.ReadFile(filename)
		if err != nil {
			set.errs.Add(filename, "", 0, err)
			continue
		}
		src := string(blob)
		for _, err := range set.macros.Collect(filename, src) {
			set.errs.Add(filename, src, err.Location, err)
		}
		set.files = append(set.files, queryFile{name: filename, src: src})
	}
	return set
}

func (c *Compiler) parseCatalog(schemas []string) error {
	paths, err := sqlpath.Glob(schemas)
	if err != nil {
		return err
	}

	// A macro is not SQL, and a schema file that is a query file too may
	// define or call one. Whatever is wrong with the queries themselves is
	// reported when they are parsed.
	macros := c.loadQueries().macros

	merr := multierr.New()
	files := make([]schemaFile, 0, len(paths))
	for _, path := range paths {
		blob, err := os.ReadFile(path)
		if err != nil {
			merr.Add(path, "", 0, err)
			continue
		}
		contents := migrations.RemoveRollbackStatements(string(blob))
		contents = migrations.RemovePsqlMetaCommands(contents)
		contents = macros.Expand(contents)
		files = append(files, schemaFile{name: path, contents: contents})
		c.schema = append(c.schema, contents)
	}

	if c.coreAnalysis {
		// A schema file that could not be read is not part of the schema the
		// action was keyed on, so there is no catalog to restore and nothing
		// further to report.
		if len(merr.Errs()) > 0 {
			return merr
		}
		if err := c.parseCatalogCore(files, merr); err != nil {
			return err
		}
	} else {
		c.parseCatalogLegacy(files, merr)
	}

	if len(merr.Errs()) > 0 {
		return merr
	}
	return nil
}

func (c *Compiler) parseCatalogLegacy(files []schemaFile, merr *multierr.Error) {
	for _, file := range files {
		stmts, err := c.parser.Parse(strings.NewReader(file.contents))
		if err != nil {
			// A schema file and a query file are often the same file, so a
			// query's sqlc syntax can fail here. Look for an explanation
			// before reporting the syntax error it caused.
			if reported := addSyntaxErrors(merr, file.name, file.contents, c.loadQueries().macros.File(file.contents)); !reported {
				merr.Add(file.name, file.contents, 0, err)
			}
			continue
		}

		for i := range stmts {
			if err := c.catalog.Update(stmts[i], c); err != nil {
				merr.Add(file.name, file.contents, stmts[i].Pos(), err)
				continue
			}
		}
	}
}

// parseCatalogCore builds the catalog the analysis core works against: the
// dialect's seed with the schema's DDL applied. The two are one cached action,
// so a run that has built this catalog before restores it and never parses the
// schema at all.
func (c *Compiler) parseCatalogCore(files []schemaFile, merr *multierr.Error) error {
	contents := make([]string, 0, len(files))
	for _, file := range files {
		contents = append(contents, file.contents)
	}

	cat, err := core.NewCached(string(c.conf.Engine), contents, func(cat *core.Catalog) error {
		for _, file := range files {
			stmts, err := c.parser.Parse(strings.NewReader(file.contents))
			if err != nil {
				// A schema file and a query file are often the same file, so a
				// query's sqlc syntax can fail here. Look for an explanation
				// before reporting the syntax error it caused.
				if reported := addSyntaxErrors(merr, file.name, file.contents, c.loadQueries().macros.File(file.contents)); !reported {
					merr.Add(file.name, file.contents, 0, err)
				}
				continue
			}
			for i := range stmts {
				if err := coreschema.Apply(cat, stmts[i].Raw); err != nil {
					merr.Add(file.name, file.contents, stmts[i].Pos(), err)
					continue
				}
			}
		}
		// A catalog the whole schema did not make it into is not the catalog
		// this action describes, so it is not one to keep.
		if len(merr.Errs()) > 0 {
			return merr
		}
		return nil
	}, c.coreDialect)

	if cat == nil {
		return fmt.Errorf("%s: init catalog: %w", c.conf.Engine, err)
	}
	// Whatever apply reported is already in merr, which the caller returns.
	c.coreCatalog = cat

	// Codegen consumes the catalog through the Result, in the legacy shape.
	legacy, err := coreResultCatalog(cat)
	if err != nil {
		return fmt.Errorf("%s: dump catalog: %w", c.conf.Engine, err)
	}
	c.catalog = legacy
	return nil
}

// statement is a parsed query awaiting analysis, with what error reporting
// needs to place it back in the file it came from.
type statement struct {
	filename string
	src      string
	pp       *preprocess.Result
	raw      *ast.RawStmt
}

func (c *Compiler) parseQueries(o opts.Parser) (*Result, error) {
	loaded := c.loadQueries()
	if loaded.globErr != nil {
		return nil, loaded.globErr
	}
	merr := loaded.errs
	macros := loaded.macros

	var stmts []statement
	for _, file := range loaded.files {
		filename, src := file.name, file.src

		// sqlc syntax is not SQL. Expand the macros and rewrite the rest to
		// native placeholders before the engine parser ever sees the query,
		// so parsers only handle SQL.
		pp := macros.File(src)

		parsed, err := c.parser.Parse(strings.NewReader(pp.Text))
		if err != nil {
			if reported := addSyntaxErrors(merr, filename, src, pp); !reported {
				merr.Add(filename, src, 0, err)
			}
			continue
		}
		for _, stmt := range parsed {
			stmts = append(stmts, statement{filename: filename, src: src, pp: pp, raw: stmt.Raw})
		}
	}

	queries, errs := c.analyzeStatements(stmts, o)

	var q []*Query
	set := map[string]struct{}{}
	for i, stmt := range stmts {
		if err := errs[i]; err != nil {
			var e *sqlerr.Error
			loc := stmt.raw.Pos()
			if errors.As(err, &e) && e.Location != 0 {
				loc = e.Location
			}
			// Locations are reported against the rewritten query; map them
			// back so errors point at what the user wrote.
			loc = stmt.pp.Origin(loc)
			if e != nil && e.Location != 0 {
				e.Location = loc
			}
			merr.Add(stmt.filename, stmt.src, loc, err)
			// If this rpc unauthenticated error bubbles up, then all future parsing/analysis will fail
			if errors.Is(err, rpc.ErrUnauthenticated) {
				return nil, merr
			}
			continue
		}
		query := queries[i]
		if query == nil {
			continue
		}
		query.Metadata.Filename = filepath.Base(stmt.filename)
		queryName := query.Metadata.Name
		if queryName != "" {
			if _, exists := set[queryName]; exists {
				merr.Add(stmt.filename, stmt.src, stmt.pp.Origin(stmt.raw.Pos()), fmt.Errorf("duplicate query name: %s", queryName))
				continue
			}
			set[queryName] = struct{}{}
		}
		q = append(q, query)
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

// analyzeStatements analyzes each statement, returning the results in the
// order the statements were given so that queries and errors come out in
// source order however they were produced.
//
// The analysis core reads the catalog and nothing else, so its statements are
// analyzed concurrently. Every other path holds state that a second goroutine
// would race — a database connection, the legacy catalog — and stays serial.
func (c *Compiler) analyzeStatements(stmts []statement, o opts.Parser) ([]*Query, []error) {
	queries := make([]*Query, len(stmts))
	errs := make([]error, len(stmts))

	analyze := func(i int) {
		queries[i], errs[i] = c.parseQuery(stmts[i].raw, stmts[i].pp, o)
	}

	workers := runtime.GOMAXPROCS(0)
	if !c.coreAnalysis || workers < 2 || len(stmts) < 2 {
		for i := range stmts {
			analyze(i)
		}
		return queries, errs
	}

	var g errgroup.Group
	g.SetLimit(workers)
	for i := range stmts {
		g.Go(func() error {
			analyze(i)
			return nil
		})
	}
	g.Wait()
	return queries, errs
}

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
