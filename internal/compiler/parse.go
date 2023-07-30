package compiler

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
	"github.com/sqlc-dev/sqlc/internal/sql/validate"
)

var ErrUnsupportedStatementType = errors.New("parseQuery: unsupported statement type")

func (c *Compiler) parseQuery(stmt ast.Node, src string, o opts.Parser) (*Query, error) {
	if o.Debug.DumpAST {
		debug.Dump(stmt)
	}
	if err := validate.ParamStyle(stmt); err != nil {
		return nil, err
	}
	numbers, dollar, err := validate.ParamRef(stmt)
	if err != nil {
		return nil, err
	}
	raw, ok := stmt.(*ast.RawStmt)
	if !ok {
		return nil, errors.New("node is not a statement")
	}
	var table *ast.TableName
	switch n := raw.Stmt.(type) {
	case *ast.CallStmt:
	case *ast.SelectStmt:
	case *ast.DeleteStmt:
	case *ast.InsertStmt:
		if err := validate.InsertStmt(n); err != nil {
			return nil, err
		}
		var err error
		table, err = ParseTableName(n.Relation)
		if err != nil {
			return nil, err
		}
	case *ast.ListenStmt:
	case *ast.NotifyStmt:
	case *ast.TruncateStmt:
	case *ast.UpdateStmt:
	case *ast.RefreshMatViewStmt:
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
	if err := validate.FuncCall(c.catalog, c.combo, raw); err != nil {
		return nil, err
	}
	if err := validate.In(c.catalog, raw); err != nil {
		return nil, err
	}
	name, cmd, err := metadata.ParseQueryNameAndType(strings.TrimSpace(rawSQL), c.parser.CommentSyntax())
	if err != nil {
		return nil, err
	}
	raw, namedParams, edits := rewrite.NamedParameters(c.conf.Engine, raw, numbers, dollar)
	if err := validate.Cmd(raw.Stmt, name, cmd); err != nil {
		return nil, err
	}
	rvs := rangeVars(raw.Stmt)
	refs, err := findParameters(raw.Stmt)
	if err != nil {
		return nil, err
	}
	refs = uniqueParamRefs(refs, dollar)
	if c.conf.Engine == config.EngineMySQL || !dollar {
		sort.Slice(refs, func(i, j int) bool { return refs[i].ref.Location < refs[j].ref.Location })
	} else {
		sort.Slice(refs, func(i, j int) bool { return refs[i].ref.Number < refs[j].ref.Number })
	}
	raw, embeds := rewrite.Embeds(raw)
	qc, err := c.buildQueryCatalog(c.catalog, raw.Stmt, embeds)
	if err != nil {
		return nil, err
	}

	params, err := c.resolveCatalogRefs(qc, rvs, refs, namedParams, embeds)
	if err != nil {
		return nil, err
	}
	cols, err := c.outputColumns(qc, raw.Stmt)
	if err != nil {
		return nil, err
	}

	expandEdits, err := c.expand(qc, raw)
	if err != nil {
		return nil, err
	}
	edits = append(edits, expandEdits...)
	expanded, err := source.Mutate(rawSQL, edits)
	if err != nil {
		return nil, err
	}

	// If the query string was edited, make sure the syntax is valid
	if expanded != rawSQL {
		if _, err := c.parser.Parse(strings.NewReader(expanded)); err != nil {
			return nil, fmt.Errorf("edited query syntax is invalid: %w", err)
		}
	}

	trimmed, comments, err := source.StripComments(expanded)
	if err != nil {
		return nil, err
	}

	flags, err := metadata.ParseQueryFlags(comments)
	if err != nil {
		return nil, err
	}

	return &Query{
		RawStmt:         raw,
		Cmd:             cmd,
		Comments:        comments,
		Name:            name,
		Flags:           flags,
		Params:          params,
		Columns:         cols,
		SQL:             trimmed,
		InsertIntoTable: table,
	}, nil
}

func rangeVars(root ast.Node) []*ast.RangeVar {
	var vars []*ast.RangeVar
	find := astutils.VisitorFunc(func(node ast.Node) {
		switch n := node.(type) {
		case *ast.RangeVar:
			vars = append(vars, n)
		}
	})
	astutils.Walk(find, root)
	return vars
}

func uniqueParamRefs(in []paramRef, dollar bool) []paramRef {
	m := make(map[int]bool, len(in))
	o := make([]paramRef, 0, len(in))
	for _, v := range in {
		if !m[v.ref.Number] {
			m[v.ref.Number] = true
			if v.ref.Number != 0 {
				o = append(o, v)
			}
		}
	}
	if !dollar {
		start := 1
		for _, v := range in {
			if v.ref.Number == 0 {
				for m[start] {
					start++
				}
				v.ref.Number = start
				o = append(o, v)
			}
		}
	}
	return o
}
