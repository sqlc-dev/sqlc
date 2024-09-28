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

	name, cmd, err := metadata.ParseQueryNameAndType(rawSQL, metadata.CommentSyntax(c.parser.CommentSyntax()))
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, nil
	}

	if err := validate.Cmd(raw.Stmt, name, cmd); err != nil {
		return nil, err
	}

	md := metadata.Metadata{
		Name: name,
		Cmd:  cmd,
	}

	// TODO eventually can use this for name and type/cmd parsing too
	cleanedComments, err := source.CleanedComments(rawSQL, c.parser.CommentSyntax())
	if err != nil {
		return nil, err
	}

	md.Params, md.Flags, md.RuleSkiplist, err = metadata.ParseCommentFlags(cleanedComments)
	if err != nil {
		return nil, err
	}

	var anlys *analysis
	if c.analyzer != nil {
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

		// If the query uses star expansion, verify that it was edited. If not,
		// return an error.
		stars := astutils.Search(raw, func(node ast.Node) bool {
			_, ok := node.(*ast.A_Star)
			return ok
		})
		hasStars := len(stars.Items) > 0
		unchanged := inference.Query == rawSQL
		if unchanged && hasStars {
			return nil, fmt.Errorf("star expansion failed for query")
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

	md.Comments = comments

	return &Query{
		RawStmt:         raw,
		Metadata:        md,
		Params:          anlys.Parameters,
		Columns:         anlys.Columns,
		SQL:             trimmed,
		InsertIntoTable: anlys.Table,
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
