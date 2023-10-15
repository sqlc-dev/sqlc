package compiler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/validate"
)

var ErrUnsupportedStatementType = errors.New("parseQuery: unsupported statement type")

func (c *Compiler) parseQuery(stmt ast.Node, src string, o opts.Parser) (*Query, error) {
	ctx := context.Background()

	if o.Debug.DumpAST {
		debug.Dump(stmt)
	}

	// validate sqlc-specific syntax
	if err := validate.SqlcFunctions(stmt); err != nil {
		return nil, err
	}

	// rewrite queries to remove sqlc.* functions

	raw, ok := stmt.(*ast.RawStmt)
	if !ok {
		return nil, errors.New("node is not a statement")
	}
	rawSQL, err := source.Pluck(src, raw.StmtLocation, raw.StmtLen)
	if err != nil {
		return nil, err
	}
	if rawSQL == "" {
		return nil, errors.New("missing semicolon at end of file")
	}

	name, cmd, err := metadata.ParseQueryNameAndType(strings.TrimSpace(rawSQL), c.parser.CommentSyntax())
	if err != nil {
		return nil, err
	}
	if err := validate.Cmd(raw.Stmt, name, cmd); err != nil {
		return nil, err
	}

	var anlys *analysis
	if c.analyzer != nil {
		// TODO: Handle panics
		inference, _ := c.inferQuery(raw, rawSQL)
		if inference == nil {
			inference = &analysis{}
		}
		if inference.Query == "" {
			inference.Query = rawSQL
		}

		result, err := c.analyzer.Analyze(ctx, raw, inference.Query, c.schema, inference.Named)
		if err != nil {
			return nil, err
		}

		// FOOTGUN: combineAnalysis mutates inference
		anlys = combineAnalysis(inference, result)
	} else {
		anlys, err = c.analyzeQuery(raw, rawSQL)
		if err != nil {
			return nil, err
		}
	}

	expanded := anlys.Query

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
		Params:          anlys.Parameters,
		Columns:         anlys.Columns,
		SQL:             trimmed,
		InsertIntoTable: anlys.Table,
	}, nil
}

func rawRangeTblRefs(root ast.Node) ([]*ast.RangeVar, []*ast.RangeSubselect, []*ast.RangeFunction) {
	var vars []*ast.RangeVar
	var subs []*ast.RangeSubselect
	var funs []*ast.RangeFunction
	find := astutils.SingleQueryVisitorFunc(func(node ast.Node) {
		switch n := node.(type) {
		case *ast.RangeVar:
			vars = append(vars, n)
		case *ast.RangeSubselect:
			subs = append(subs, n)
		case *ast.RangeFunction:
			funs = append(funs, n)
		}
	})
	astutils.Walk(find, root)
	return vars, subs, funs
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
