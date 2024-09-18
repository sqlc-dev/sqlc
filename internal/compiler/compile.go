package compiler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	// "github.com/ryboe/q"

	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/rpc"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
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
		src := string(blob)
		contents := migrations.RemoveRollbackStatements(src)
		c.schema = append(c.schema, contents)
		stmts, err := c.parser.Parse(strings.NewReader(contents))
		if err != nil {
			merr.Add(filename, contents, 0, err)
			continue
		}
		for i, stmt := range stmts {
			_, rawSQL, err := getRaw(stmt.Raw, src)
			if err != nil {
				return err
			}
			RawComments, err := source.RawComments(rawSQL, c.parser.CommentSyntax())
			if err != nil {
				return err
			}

			if err := c.catalog.Update(stmts[i], RawComments, c); err != nil {
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
					merr.Add(
						filename,
						src,
						stmt.Raw.Pos(),
						fmt.Errorf("duplicate query name: %s", queryName),
					)
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
		return nil, fmt.Errorf(
			"no queries contained in paths %s",
			strings.Join(c.conf.Queries, ","),
		)
	}
	return &Result{
		Catalog: c.catalog,
		Queries: q,
	}, nil
}

func getRaw(stmt ast.Node, src string) (*ast.RawStmt, string, error) {
	raw, ok := stmt.(*ast.RawStmt)
	if !ok {
		return nil, "", errors.New("node is not a statement")
	}
	rawSQL, err := source.Pluck(src, raw.StmtLocation, raw.StmtLen)
	if err != nil {
		return nil, "", err
	}
	if rawSQL == "" {
		return nil, "", errors.New("missing semicolon at end of file")
	}
	return raw, rawSQL, nil
}
