package compiler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/rpc"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlfile"
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
		c.schema = append(c.schema, contents)
		stmts, err := c.parser.Parse(strings.NewReader(contents))
		if err != nil {
			merr.Add(filename, contents, 0, err)
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
		stmts, err := c.parser.Parse(strings.NewReader(src))
		if err != nil {
			merr.Add(filename, src, 0, err)
			continue
		}
		for _, stmt := range stmts {
			query, err := c.parseQuery(stmt.Raw, src, o)
			if err != nil {
				var e *sqlerr.Error
				loc := stmt.Raw.Pos()
				if errors.As(err, &e) && e.Location != 0 {
					loc = e.Location
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
					merr.Add(filename, src, stmt.Raw.Pos(), fmt.Errorf("duplicate query name: %s", queryName))
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

// parseQueriesWithAnalyzer parses queries using only the database analyzer,
// skipping the parser and catalog entirely. Uses sqlfile.Split to extract
// individual queries from .sql files.
func (c *Compiler) parseQueriesWithAnalyzer(o opts.Parser) (*Result, error) {
	ctx := context.Background()
	var q []*Query
	merr := multierr.New()
	set := map[string]struct{}{}
	files, err := sqlpath.Glob(c.conf.Queries)
	if err != nil {
		return nil, err
	}

	if c.analyzer == nil {
		return nil, fmt.Errorf("database analyzer is required when skip_parser is enabled")
	}

	for _, filename := range files {
		blob, err := os.ReadFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		src := string(blob)

		// Use sqlfile.Split to extract individual queries
		queries, err := sqlfile.Split(ctx, strings.NewReader(src))
		if err != nil {
			merr.Add(filename, src, 0, err)
			continue
		}

		for _, queryText := range queries {
			// Extract metadata from comments
			name, cmd, err := metadata.ParseQueryNameAndType(queryText, metadata.CommentSyntax{Dash: true})
			if err != nil {
				merr.Add(filename, queryText, 0, err)
				continue
			}

			// Skip queries without names (not marked with sqlc comments)
			if name == "" {
				continue
			}

			// Check for duplicate query names
			if _, exists := set[name]; exists {
				merr.Add(filename, queryText, 0, fmt.Errorf("duplicate query name: %s", name))
				continue
			}
			set[name] = struct{}{}

			// Extract additional metadata from comments
			cleanedComments, err := source.CleanedComments(queryText, source.CommentSyntax{Dash: true})
			if err != nil {
				merr.Add(filename, queryText, 0, err)
				continue
			}

			md := metadata.Metadata{
				Name:     name,
				Cmd:      cmd,
				Filename: filepath.Base(filename),
			}

			md.Params, md.Flags, md.RuleSkiplist, err = metadata.ParseCommentFlags(cleanedComments)
			if err != nil {
				merr.Add(filename, queryText, 0, err)
				continue
			}

			// Use the database analyzer to analyze the query
			// We pass an empty AST node since we're not using the parser
			result, err := c.analyzer.Analyze(ctx, nil, queryText, c.schema, &named.ParamSet{})
			if err != nil {
				merr.Add(filename, queryText, 0, err)
				// If this rpc unauthenticated error bubbles up, then all future parsing/analysis will fail
				if errors.Is(err, rpc.ErrUnauthenticated) {
					return nil, merr
				}
				continue
			}

			// Convert analyzer results to Query format
			var cols []*Column
			for _, col := range result.Columns {
				cols = append(cols, convertColumn(col))
			}

			var params []Parameter
			for _, p := range result.Params {
				params = append(params, Parameter{
					Number: int(p.Number),
					Column: convertColumn(p.Column),
				})
			}

			// Strip comments from the final SQL
			trimmed, comments, err := source.StripComments(queryText)
			if err != nil {
				merr.Add(filename, queryText, 0, err)
				continue
			}
			md.Comments = comments

			query := &Query{
				SQL:      trimmed,
				Metadata: md,
				Columns:  cols,
				Params:   params,
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
		Catalog: nil, // No catalog when skip_parser is enabled
		Queries: q,
	}, nil
}
