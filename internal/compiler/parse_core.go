package compiler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	coreanalyzer "github.com/sqlc-dev/sqlc/internal/core/analyzer"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/preprocess"
	"github.com/sqlc-dev/sqlc/internal/sql/validate"
)

func (c *Compiler) parseQueryCore(raw *ast.RawStmt, src string, pre *preprocess.Statement) (*Query, error) {
	rawSQL, err := source.Pluck(src, raw.StmtLocation, raw.StmtLen)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(rawSQL) == "" {
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

	md := metadata.Metadata{Name: name, Cmd: cmd}
	cleanedComments, err := source.CleanedComments(rawSQL, c.parser.CommentSyntax())
	if err != nil {
		return nil, err
	}
	md.Params, md.Flags, md.RuleSkiplist, err = metadata.ParseCommentFlags(cleanedComments)
	if err != nil {
		return nil, err
	}

	if pre.ParamErr != nil {
		return nil, pre.ParamErr
	}
	namedParams := pre.Params
	expanded := rawSQL

	var cols []*Column
	var params []Parameter
	switch raw.Stmt.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
		res, err := coreanalyzer.Prepare(c.coreCatalog, raw)
		if err != nil {
			return nil, err
		}
		for _, col := range res.Columns {
			cols = append(cols, coreColumn(col))
		}
		for _, p := range res.Parameters {
			params = append(params, Parameter{Number: p.Number, Column: coreParamColumn(p, namedParams)})
		}
		expanded, err = source.Mutate(rawSQL, c.expandCore(raw, res.Stars))
		if err != nil {
			return nil, err
		}
	}

	// If the query string was edited, make sure the syntax is valid
	if expanded != rawSQL {
		if _, err := c.newParser().Parse(strings.NewReader(expanded)); err != nil {
			return nil, fmt.Errorf("edited query syntax is invalid: %w", err)
		}
	}

	trimmed, comments, err := source.StripComments(expanded)
	if err != nil {
		return nil, err
	}
	md.Comments = comments

	var insertTable *ast.TableName
	if ins, ok := raw.Stmt.(*ast.InsertStmt); ok {
		insertTable, _ = ParseTableName(ins.Relation)
	}

	return &Query{
		RawStmt:         raw,
		Metadata:        md,
		Params:          params,
		Columns:         cols,
		SQL:             trimmed,
		InsertIntoTable: insertTable,
	}, nil
}

func coreColumn(c core.Column) *Column {
	col := &Column{
		Name:     c.Name,
		DataType: c.DataType,
		NotNull:  c.NotNull,
		IsArray:  c.IsArray,
		TypeExpr: c.Type,
	}
	describeType(col, c.Type)
	if c.Source != nil && c.Source.Table != "" {
		col.Table = &ast.TableName{Schema: c.Source.Schema, Name: c.Source.Table}
		col.TableAlias = c.Source.TableAlias
		col.OriginalName = c.Source.Column
	}
	return col
}

// describeType fills in what codegen reads about a type from its
// expression: one array dimension per nesting, the length that is the
// innermost type's first integer argument (which is how a MySQL tinyint(1)
// is told from a tinyint), and whether the innermost type is unsigned.
func describeType(col *Column, t *core.TypeExpr) {
	if t == nil {
		if col.IsArray {
			col.ArrayDims = 1
		}
		return
	}
	col.ArrayDims = t.ArrayDims()
	inner := t.Innermost()
	if len(inner.Args) > 0 && inner.Args[0].Int != nil {
		l := int(*inner.Args[0].Int)
		col.Length = &l
	}
	col.Unsigned = strings.HasSuffix(inner.Name, " unsigned")
}

func coreParamColumn(p core.Parameter, params *named.ParamSet) *Column {
	col := &Column{
		Name:     p.Name,
		DataType: p.DataType,
		NotNull:  p.NotNull,
		IsArray:  p.IsArray,
		TypeExpr: p.Type,
	}
	describeType(col, p.Type)
	if p.Source != nil && p.Source.Table != "" {
		col.Table = &ast.TableName{Schema: p.Source.Schema, Name: p.Source.Table}
		col.OriginalName = p.Source.Column
	}
	if col.Name == "" && p.Source != nil {
		col.Name = p.Source.Column
	}
	// Merge in what the user asked for: sqlc.narg() makes the parameter
	// nullable and sqlc.slice() marks it as a slice, whichever way the
	// analyzer typed it.
	if param, isNamed := params.FetchMerge(p.Number, named.NewInferredParam(col.Name, p.NotNull)); isNamed {
		col.Name = param.Name()
		col.NotNull = param.NotNull()
		col.IsSqlcSlice = param.IsSqlcSlice()
		col.IsNamedParam = true
	}
	return col
}
