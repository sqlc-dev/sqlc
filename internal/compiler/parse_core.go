package compiler

import (
	"errors"
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
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
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
	}
	if c.Source != nil && c.Source.Table != "" {
		col.Table = &ast.TableName{Schema: c.Source.Schema, Name: c.Source.Table}
		col.TableAlias = c.Source.TableAlias
		col.OriginalName = c.Source.Column
	}
	if c.TypeLength > 0 {
		l := c.TypeLength
		col.Length = &l
	}
	return col
}

func coreParamColumn(p core.Parameter, params *named.ParamSet) *Column {
	col := &Column{
		Name:     p.Name,
		DataType: p.DataType,
		NotNull:  p.NotNull,
		IsArray:  p.IsArray,
	}
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
