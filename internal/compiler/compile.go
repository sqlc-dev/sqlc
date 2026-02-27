package compiler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
		contents := migrations.RemoveRollbackStatements(string(blob))
		contents = removePsqlMetaCommands(contents)
		c.schema = append(c.schema, contents)

		// In database-only mode, we parse the schema to validate syntax
		// but don't update the catalog - the database will be the source of truth
		stmts, err := c.parser.Parse(strings.NewReader(contents))
		if err != nil {
			merr.Add(filename, contents, 0, err)
			continue
		}

		// Skip catalog updates in database-only mode
		if c.databaseOnlyMode {
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

func removePsqlMetaCommands(contents string) string {
	if contents == "" {
		return contents
	}
	var out strings.Builder
	out.Grow(len(contents))

	lineStart := true
	inSingle := false
	inDollar := false
	var dollarTag string
	blockDepth := 0
	n := len(contents)
	for i := 0; ; {
		if lineStart && !inSingle && blockDepth == 0 && !inDollar {
			start := i
			for i < n {
				c := contents[i]
				if c == ' ' || c == '\t' || c == '\r' {
					i++
					continue
				}
				break
			}
			if i < n && contents[i] == '\\' {
				for i < n && contents[i] != '\n' {
					i++
				}
				if i < n && contents[i] == '\n' {
					out.WriteByte('\n')
					i++
				}
				lineStart = true
				continue
			}
			if start < i {
				out.WriteString(contents[start:i])
			}
			if i >= n {
				break
			}
		}
		if i >= n {
			break
		}
		c := contents[i]
		if inSingle {
			out.WriteByte(c)
			if c == '\'' {
				if i+1 < n && contents[i+1] == '\'' {
					out.WriteByte(contents[i+1])
					i += 2
					lineStart = false
					continue
				}
				inSingle = false
			}
			if c == '\n' {
				lineStart = true
			} else {
				lineStart = false
			}
			i++
			continue
		}
		if inDollar {
			if strings.HasPrefix(contents[i:], dollarTag) {
				out.WriteString(dollarTag)
				i += len(dollarTag)
				inDollar = false
				lineStart = false
				continue
			}
			out.WriteByte(c)
			if c == '\n' {
				lineStart = true
			} else {
				lineStart = false
			}
			i++
			continue
		}
		if blockDepth > 0 {
			if c == '/' && i+1 < n && contents[i+1] == '*' {
				blockDepth++
				out.WriteString("/*")
				i += 2
				lineStart = false
				continue
			}
			if c == '*' && i+1 < n && contents[i+1] == '/' {
				blockDepth--
				out.WriteString("*/")
				i += 2
				lineStart = false
				continue
			}
			out.WriteByte(c)
			if c == '\n' {
				lineStart = true
			} else {
				lineStart = false
			}
			i++
			continue
		}
		switch c {
		case '\'':
			inSingle = true
			out.WriteByte(c)
			lineStart = false
			i++
			continue
		case '$':
			tagEnd := i + 1
			for tagEnd < n && isDollarTagChar(contents[tagEnd]) {
				tagEnd++
			}
			if tagEnd < n && contents[tagEnd] == '$' {
				dollarTag = contents[i : tagEnd+1]
				inDollar = true
				out.WriteString(dollarTag)
				i = tagEnd + 1
				lineStart = false
				continue
			}
		case '/':
			if i+1 < n && contents[i+1] == '*' {
				blockDepth = 1
				out.WriteString("/*")
				i += 2
				lineStart = false
				continue
			}
		}
		out.WriteByte(c)
		if c == '\n' {
			lineStart = true
		} else {
			lineStart = false
		}
		i++
	}
	return out.String()
}

func isDollarTagChar(b byte) bool {
	return b == '_' || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
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
