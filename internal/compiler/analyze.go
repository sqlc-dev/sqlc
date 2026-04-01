package compiler

import (
	"fmt"
	"sort"

	analyzer "github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
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
		if err := check(validate.InsertStmt(n)); err != nil {
			return nil, err
		}
		var err error
		table, err = ParseTableName(n.Relation)
		if err := check(err); err != nil {
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
	if n, ok := raw.Stmt.(*ast.InsertStmt); ok {
		if err := check(c.validateOnConflictClause(n, params)); err != nil {
			return nil, err
		}
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

// validateOnConflictClause validates an ON CONFLICT DO UPDATE clause against
// the target table. It checks:
//   - ON CONFLICT (col, ...) conflict target columns exist
//   - DO UPDATE SET col = ... assignment target columns exist
//   - EXCLUDED.col references exist
//   - For PostgreSQL: $N param and EXCLUDED.col type compatibility with SET target
func (c *Compiler) validateOnConflictClause(n *ast.InsertStmt, params []Parameter) error {
	if n.OnConflictClause == nil || n.OnConflictClause.Action != ast.OnConflictActionUpdate {
		return nil
	}
	fqn, err := ParseTableName(n.Relation)
	if err != nil {
		return err
	}
	table, err := c.catalog.GetTable(fqn)
	if err != nil {
		return err
	}

	// Build column name → DataType from catalog for existence and type checks.
	colDataTypes := make(map[string]string, len(table.Columns))
	for _, col := range table.Columns {
		colDataTypes[col.Name] = dataType(&col.Type)
	}

	// Validate ON CONFLICT (col, ...) conflict target columns.
	if n.OnConflictClause.Infer != nil {
		for _, item := range n.OnConflictClause.Infer.IndexElems.Items {
			elem, ok := item.(*ast.IndexElem)
			if !ok || elem.Name == nil {
				continue
			}
			if _, exists := colDataTypes[*elem.Name]; !exists {
				e := sqlerr.ColumnNotFound(table.Rel.Name, *elem.Name)
				e.Location = n.OnConflictClause.Infer.Location
				return e
			}
		}
	}

	// Validate DO UPDATE SET col = ... assignment target columns and EXCLUDED.col references.
	for _, item := range n.OnConflictClause.TargetList.Items {
		target, ok := item.(*ast.ResTarget)
		if !ok || target.Name == nil {
			continue
		}
		if _, exists := colDataTypes[*target.Name]; !exists {
			e := sqlerr.ColumnNotFound(table.Rel.Name, *target.Name)
			e.Location = target.Location
			return e
		}
		if ref, ok := target.Val.(*ast.ColumnRef); ok {
			if excludedCol, ok := excludedColumn(ref); ok {
				if _, exists := colDataTypes[excludedCol]; !exists {
					e := sqlerr.ColumnNotFound(table.Rel.Name, excludedCol)
					e.Location = ref.Location
					return e
				}
			}
		}
	}

	// Type compatibility checks (PostgreSQL only).
	// To remove type checking: delete this block and validateOnConflictSetTypes.
	if c.conf.Engine == config.EnginePostgreSQL {
		if err := c.validateOnConflictSetTypes(n, params, colDataTypes); err != nil {
			return err
		}
	}

	return nil
}

// validateOnConflictSetTypes checks that values in DO UPDATE SET assignments
// are type-compatible with their target columns (PostgreSQL only).
// It handles $N params (typed from INSERT VALUES) and EXCLUDED.col references.
func (c *Compiler) validateOnConflictSetTypes(n *ast.InsertStmt, params []Parameter, colDataTypes map[string]string) error {
	// Build param number → resolved DataType from already-resolved params.
	// Skips params with "any" type (unresolved).
	paramDataTypes := make(map[int]string, len(params))
	for i := range params {
		if params[i].Column != nil && params[i].Column.DataType != "any" {
			paramDataTypes[params[i].Number] = params[i].Column.DataType
		}
	}

	for _, item := range n.OnConflictClause.TargetList.Items {
		target, ok := item.(*ast.ResTarget)
		if !ok || target.Name == nil {
			continue
		}
		colDT, ok := colDataTypes[*target.Name]
		if !ok {
			continue
		}
		switch val := target.Val.(type) {
		case *ast.ParamRef:
			paramDT, ok := paramDataTypes[val.Number]
			if !ok {
				continue
			}
			if postgresql.TypeFamily(paramDT) != postgresql.TypeFamily(colDT) {
				return &sqlerr.Error{
					Message:  fmt.Sprintf("parameter $%d has type %q but column %q has type %q", val.Number, paramDT, *target.Name, colDT),
					Location: val.Location,
				}
			}
		case *ast.ColumnRef:
			excludedCol, ok := excludedColumn(val)
			if !ok {
				continue
			}
			excludedDT, ok := colDataTypes[excludedCol]
			if !ok {
				continue
			}
			if postgresql.TypeFamily(excludedDT) != postgresql.TypeFamily(colDT) {
				return &sqlerr.Error{
					Message:  fmt.Sprintf("EXCLUDED.%s has type %q but column %q has type %q", excludedCol, excludedDT, *target.Name, colDT),
					Location: val.Location,
				}
			}
		}
	}
	return nil
}

// excludedColumn returns the column name if the ColumnRef is an EXCLUDED.col
// reference, and ok=true. Returns "", false otherwise.
func excludedColumn(ref *ast.ColumnRef) (string, bool) {
	if ref.Fields == nil || len(ref.Fields.Items) != 2 {
		return "", false
	}
	first, ok := ref.Fields.Items[0].(*ast.String)
	if !ok || first.Str != "excluded" {
		return "", false
	}
	second, ok := ref.Fields.Items[1].(*ast.String)
	if !ok {
		return "", false
	}
	return second.Str, true
}
