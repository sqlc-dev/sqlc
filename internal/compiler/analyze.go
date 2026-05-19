package compiler

import (
	"sort"

	analyzer "github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
	"github.com/sqlc-dev/sqlc/internal/sql/validate"
)

type analysis struct {
	Table      *ast.TableName
	Columns    []*Column
	Parameters []Parameter
	Named      *named.ParamSet
	Query      string
}

func convertTableName(id *analyzer.Identifier) *ast.TableName {
	if id == nil {
		return nil
	}
	return &ast.TableName{
		Catalog: id.Catalog,
		Schema:  id.Schema,
		Name:    id.Name,
	}
}

func convertTypeName(id *analyzer.Identifier) *ast.TypeName {
	if id == nil {
		return nil
	}
	return &ast.TypeName{
		Catalog: id.Catalog,
		Schema:  id.Schema,
		Name:    id.Name,
	}
}

func convertColumn(c *analyzer.Column) *Column {
	length := int(c.Length)
	return &Column{
		Name:         c.Name,
		OriginalName: c.OriginalName,
		DataType:     c.DataType,
		NotNull:      c.NotNull,
		Unsigned:     c.Unsigned,
		IsArray:      c.IsArray,
		ArrayDims:    int(c.ArrayDims),
		Comment:      c.Comment,
		Length:       &length,
		IsNamedParam: c.IsNamedParam,
		IsFuncCall:   c.IsFuncCall,
		Scope:        c.Scope,
		Table:        convertTableName(c.Table),
		TableAlias:   c.TableAlias,
		Type:         convertTypeName(c.Type),
		EmbedTable:   convertTableName(c.EmbedTable),
		IsSqlcSlice:  c.IsSqlcSlice,
	}
}

func mergeColumnOrigin(dst, src *Column) {
	if dst == nil || src == nil {
		return
	}

	// Column overrides in the Go generator depend on the column's original
	// table identity. The analyzer can fill in missing origin metadata, but it
	// must not overwrite catalog-inferred origin metadata.
	//
	// In particular, CTE output columns are deliberately re-scoped to the CTE
	// relation by buildQueryCatalog. If we overwrite that with the analyzer's
	// underlying base table, queries like:
	//
	//   WITH expensive AS (SELECT * FROM products)
	//   SELECT * FROM expensive
	//
	// start looking like they return products directly, causing the Go generator
	// to reuse Product instead of emitting ListExpensiveProductsRow.
	if dst.OriginalName == "" && src.OriginalName != "" {
		dst.OriginalName = src.OriginalName
	}
	if dst.Table == nil && src.Table != nil {
		dst.Table = src.Table
	}
	if dst.TableAlias == "" && src.TableAlias != "" {
		dst.TableAlias = src.TableAlias
	}
	if dst.Scope == "" && src.Scope != "" {
		dst.Scope = src.Scope
	}
	if dst.EmbedTable == nil && src.EmbedTable != nil {
		dst.EmbedTable = src.EmbedTable
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
			Number: int(p.Number),
			Column: convertColumn(p.Column),
		})
	}
	if len(prev.Columns) == len(cols) {
		for i := range prev.Columns {
			mergeColumnOrigin(prev.Columns[i], cols[i])
			// Only override column types if the analyzer provides a specific type
			// (not "any"), since the catalog-based inference may have better info
			if cols[i].DataType != "any" {
				prev.Columns[i].DataType = cols[i].DataType
				prev.Columns[i].IsArray = cols[i].IsArray
				prev.Columns[i].ArrayDims = cols[i].ArrayDims
			}
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
			// Only override parameter types if the analyzer provides a specific type
			// (not "any"), since the catalog-based inference may have better info
			if params[i].Column.DataType != "any" {
				prev.Parameters[i].Column.DataType = params[i].Column.DataType
				prev.Parameters[i].Column.IsArray = params[i].Column.IsArray
				prev.Parameters[i].Column.ArrayDims = params[i].Column.ArrayDims
			}
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
	case *ast.InsertStmt:
		var err error
		table, err = ParseTableName(n.Relation)
		if err := check(err); err != nil {
			return nil, err
		}
		if err := check(validate.InsertStmt(c.catalog, table, n)); err != nil {
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

	params, err := c.resolveCatalogRefs(qc, rvs, refs, namedParams, embeds)
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
		Table:      table,
		Columns:    cols,
		Parameters: params,
		Query:      expanded,
		Named:      namedParams,
	}, rerr
}
