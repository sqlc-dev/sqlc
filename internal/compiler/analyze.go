package compiler

import (
	"sort"

	"github.com/sqlc-dev/sqlc/internal/analyzer"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
	"github.com/sqlc-dev/sqlc/internal/sql/validate"
)

type analysis struct {
	Table        *ast.TableName
	Columns      []*Column
	QueryCatalog *QueryCatalog
	Parameters   []Parameter
	Named        *named.ParamSet
	Query        string
}

func convertColumn(c analyzer.Column) *Column {
	return &Column{
		Name:         c.Name,
		OriginalName: c.OriginalName,
		DataType:     c.DataType,
		NotNull:      c.NotNull,
		Unsigned:     c.Unsigned,
		IsArray:      c.IsArray,
		ArrayDims:    c.ArrayDims,
		Comment:      c.Comment,
		Length:       c.Length,
		IsNamedParam: c.IsNamedParam,
		IsFuncCall:   c.IsFuncCall,
		Scope:        c.Scope,
		Table:        c.Table,
		TableAlias:   c.TableAlias,
		Type:         c.Type,
		EmbedTable:   c.EmbedTable,
		IsSqlcSlice:  c.IsSqlcSlice,
	}
}

func combineAnalysis(prev *analysis, a *analyzer.Analysis) *analysis {
	var cols []*Column
	for _, c := range a.Columns {
		cols = append(cols, convertColumn(c))
	}
	var params []Parameter
	for _, p := range a.Params {
		params = append(params, Parameter{
			Number: p.Number,
			Column: convertColumn(*p.Column),
		})
	}
	if len(prev.Columns) == len(cols) {
		for i := range prev.Columns {
			prev.Columns[i].DataType = cols[i].DataType
		}
	} else {
		embedding := false
		for i := range prev.Columns {
			if prev.Columns[i].EmbedTable != nil {
				embedding = true
			}
		}
		if !embedding {
			prev.Columns = cols
		}
	}
	if len(prev.Parameters) == len(params) {
		for i := range prev.Parameters {
			prev.Parameters[i].Column.DataType = params[i].Column.DataType
		}
	} else {
		prev.Parameters = params
	}
	return prev
}

func (c *Compiler) analyzeQuery(raw *ast.RawStmt, query string) (*analysis, error) {
	return c._analyzeQuery(raw, query, true)
}

func (c *Compiler) inferQuery(raw *ast.RawStmt, query string) (*analysis, error) {
	return c._analyzeQuery(raw, query, false)
}

func (c *Compiler) _analyzeQuery(raw *ast.RawStmt, query string, failfast bool) (*analysis, error) {
	errors := make([]error, 0)
	check := func(err error) error {
		if failfast {
			return err
		}
		if err != nil {
			errors = append(errors, err)
		}
		return nil
	}

	numbers, dollar, err := validate.ParamRef(raw)
	if err := check(err); err != nil {
		return nil, err
	}

	raw, namedParams, edits := rewrite.NamedParameters(c.conf.Engine, raw, numbers, dollar)

	var table *ast.TableName
	switch n := raw.Stmt.(type) {
	case *ast.CallStmt:
	case *ast.SelectStmt:
	case *ast.DeleteStmt:
	case *ast.DoStmt:
	case *ast.InsertStmt:
		if err := check(validate.InsertStmt(n)); err != nil {
			return nil, err
		}
		var err error
		table, err = ParseTableName(n.Relation)
		if err := check(err); err != nil {
			return nil, err
		}
	case *ast.ListenStmt:
	case *ast.NotifyStmt:
	case *ast.TruncateStmt:
	case *ast.UpdateStmt:
	case *ast.RefreshMatViewStmt:
	default:
		if err := check(ErrUnsupportedStatementType); err != nil {
			return nil, err
		}
	}

	if err := check(validate.FuncCall(c.catalog, c.combo, raw)); err != nil {
		return nil, err
	}

	if err := check(validate.In(c.catalog, raw)); err != nil {
		return nil, err
	}
	rvs := rangeVars(raw.Stmt)
	refs, errs := findParameters(raw.Stmt)
	if len(errs) > 0 {
		if failfast {
			return nil, errs[0]
		}
		errors = append(errors, errs...)
	}
	refs = uniqueParamRefs(refs, dollar)
	if c.conf.Engine == config.EngineMySQL || !dollar {
		sort.Slice(refs, func(i, j int) bool { return refs[i].ref.Location < refs[j].ref.Location })
	} else {
		sort.Slice(refs, func(i, j int) bool { return refs[i].ref.Number < refs[j].ref.Number })
	}
	raw, embeds := rewrite.Embeds(raw)
	qc, err := c.buildQueryCatalog(c.catalog, raw.Stmt, embeds)
	if err := check(err); err != nil {
		return nil, err
	}

	params, err := c.resolveCatalogRefs(qc, rvs, refs, namedParams)
	if err := check(err); err != nil {
		return nil, err
	}
	err = c.resolveCatalogEmbeds(qc, rvs, embeds)
	if err := check(err); err != nil {
		return nil, err
	}
	cols, err := c.outputColumns(qc, raw.Stmt)
	if err := check(err); err != nil {
		return nil, err
	}

	expandEdits, err := c.expand(qc, raw)
	if check(err); err != nil {
		return nil, err
	}
	edits = append(edits, expandEdits...)
	expanded, err := source.Mutate(query, edits)
	if err != nil {
		return nil, err
	}

	var rerr error
	if len(errors) > 0 {
		rerr = errors[0]
	}

	return &analysis{
		Table:        table,
		Columns:      cols,
		Parameters:   params,
		QueryCatalog: qc,
		Query:        expanded,
		Named:        namedParams,
	}, rerr
}
