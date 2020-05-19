package compiler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/source"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sql/validate"
)

type Query struct {
	SQL      string
	Name     string
	Cmd      string // TODO: Pick a better name. One of: one, many, exec, execrows
	Comments []string
}

var ErrUnsupportedStatementType = errors.New("parseQuery: unsupported statement type")

func parseQuery(p Parser, c *catalog.Catalog, stmt ast.Node, src string, rewriteParameters bool) (*Query, error) {
	if err := validate.ParamStyle(stmt); err != nil {
		return nil, err
	}
	if err := validate.ParamRef(stmt); err != nil {
		return nil, err
	}
	raw, ok := stmt.(*ast.RawStmt)
	if !ok {
		return nil, errors.New("node is not a statement")
	}
	switch n := raw.Stmt.(type) {
	case *pg.SelectStmt:
	case *pg.DeleteStmt:
	case *pg.InsertStmt:
		if err := validate.InsertStmt(n); err != nil {
			return nil, err
		}
	case *pg.TruncateStmt:
	case *pg.UpdateStmt:
	default:
		return nil, ErrUnsupportedStatementType
	}

	rawSQL, err := source.Pluck(src, raw.StmtLocation, raw.StmtLen)
	if err != nil {
		return nil, err
	}
	if rawSQL == "" {
		return nil, errors.New("missing semicolon at end of file")
	}
	if err := validate.FuncCall(c, raw); err != nil {
		return nil, err
	}
	name, cmd, err := metadata.Parse(strings.TrimSpace(rawSQL), metadata.CommentSyntaxDash)
	if err != nil {
		return nil, err
	}
	if err := validate.Cmd(raw.Stmt, name, cmd); err != nil {
		return nil, err
	}

	// TODO: Then a miracle occurs

	var edits []source.Edit
	expanded, err := source.Mutate(rawSQL, edits)
	if err != nil {
		return nil, err
	}

	// If the query string was edited, make sure the syntax is valid
	if expanded != rawSQL {
		if _, err := p.Parse(strings.NewReader(expanded)); err != nil {
			return nil, fmt.Errorf("edited query syntax is invalid: %w", err)
		}
	}

	trimmed, comments, err := source.StripComments(expanded)
	if err != nil {
		return nil, err
	}

	return &Query{
		Cmd:      cmd,
		Comments: comments,
		Name:     name,
		SQL:      trimmed,
	}, nil
}
