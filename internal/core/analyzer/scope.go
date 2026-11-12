package analyzer

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type scope struct {
	rels []scopeRel
}

type scopeRel struct {
	alias    string
	classOID int64
	// cols is the catalog's column list, held as-is rather than copied into a
	// scope-local column type.
	cols []core.ClassColumn
}

func (a *analyzer) buildScope(from *ast.List) (*scope, error) {
	items := listItems(from)
	sc := &scope{rels: make([]scopeRel, 0, len(items))}
	defer a.binding(sc)()
	for _, item := range items {
		if err := a.appendFromItem(sc, item); err != nil {
			return nil, err
		}
	}
	return sc, nil
}

// binding makes the scope under construction the one the analyzer resolves
// against, and returns the func that puts back the scope it replaced. A FROM
// item can refer to the ones before it — a set-returning function takes its
// arguments from them — so binding an item has to see what is bound so far
// rather than no scope at all.
func (a *analyzer) binding(sc *scope) func() {
	prev := a.scope
	a.scope = sc
	return func() { a.scope = prev }
}

func (a *analyzer) appendFromItem(sc *scope, item ast.Node) error {
	switch v := item.(type) {
	case *ast.RangeVar:
		rel, err := a.bindRangeVar(v)
		if err != nil {
			return err
		}
		sc.rels = append(sc.rels, rel)
		return nil
	case *ast.JoinExpr:
		start := len(sc.rels)
		if err := a.appendFromItem(sc, v.Larg); err != nil {
			return err
		}
		mid := len(sc.rels)
		if err := a.appendFromItem(sc, v.Rarg); err != nil {
			return err
		}
		if a.cat.OuterJoinsNullable() {
			switch v.Jointype {
			case ast.JoinTypeLeft:
				nullableRels(sc.rels[mid:])
			case ast.JoinTypeRight:
				nullableRels(sc.rels[start:mid])
			case ast.JoinTypeFull:
				nullableRels(sc.rels[start:])
			}
		}
		return nil
	case *ast.RangeFunction:
		rel, err := a.bindRangeFunction(v)
		if err != nil {
			return err
		}
		sc.rels = append(sc.rels, rel)
		return nil
	case *ast.RangeSubselect:
		rel, err := a.bindRangeSubselect(v)
		if err != nil {
			return err
		}
		sc.rels = append(sc.rels, rel)
		return nil
	case *ast.RangePivot:
		rel, err := a.bindRangePivot(v)
		if err != nil {
			return err
		}
		sc.rels = append(sc.rels, rel)
		return nil
	case *ast.List:
		// Some engines report a comma-separated FROM as a nested list.
		for _, item := range listItems(v) {
			if err := a.appendFromItem(sc, item); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("scope: unsupported FROM item %T", item)
	}
}

// nullableRels makes every column of rels nullable, as the optional side of
// an outer join is. The columns are copied: a relation's column list is the
// catalog's own.
func nullableRels(rels []scopeRel) {
	for i := range rels {
		cols := make([]core.ClassColumn, len(rels[i].cols))
		for j, c := range rels[i].cols {
			c.NotNull = false
			cols[j] = c
		}
		rels[i].cols = cols
	}
}

// bindRangeFunction binds a function called in FROM. A set-returning function
// stands in for a relation of a single column named after it.
func (a *analyzer) bindRangeFunction(rf *ast.RangeFunction) (scopeRel, error) {
	call := findFuncCall(rf.Functions)
	if call == nil {
		return scopeRel{}, fmt.Errorf("range function: no function call")
	}
	name := funcCallName(call)
	overloads, err := a.cat.FindProcs(name, nil)
	if err != nil {
		return scopeRel{}, err
	}
	if len(overloads) == 0 {
		return scopeRel{}, fmt.Errorf("unknown function %q", name)
	}

	// The arguments are typed against the scope built so far, which is what
	// gives a placeholder passed to the function its type, and what
	// resolves a polymorphic result: unnest of an array of strings yields
	// strings.
	t, err := a.typeFuncCall(call)
	if err != nil {
		return scopeRel{}, err
	}

	// The relation's alias names its single column too, as "unnest(tags)
	// AS tag" does.
	colName := name
	rel := scopeRel{alias: name}
	if rf.Alias != nil && rf.Alias.Aliasname != nil && *rf.Alias.Aliasname != "" {
		rel.alias = *rf.Alias.Aliasname
		colName = rel.alias
	}
	col := core.ClassColumn{Name: colName, TypeOID: t.typeOID, Type: t.expr, NotNull: !t.nullable}
	if t.typeOID == 0 && t.expr == nil {
		p := overloads[0]
		col.TypeOID, col.NotNull = p.ReturnTypeOID, !p.ReturnNullable
	}
	rel.cols = []core.ClassColumn{col}
	if rf.Coldeflist != nil {
		rel.cols = a.coldefColumns(rf.Coldeflist, col)
	}
	return rel, nil
}

// coldefColumns are the columns a function's column definition list
// declares, for a function returning rows of more than one column. A
// definition with no type is the function's own result.
func (a *analyzer) coldefColumns(defs *ast.List, result core.ClassColumn) []core.ClassColumn {
	var cols []core.ClassColumn
	for _, item := range listItems(defs) {
		def, ok := item.(*ast.ColumnDef)
		if !ok {
			continue
		}
		if def.TypeName == nil {
			col := result
			col.Name = def.Colname
			cols = append(cols, col)
			continue
		}
		t := a.lookupType(core.TypeExprOfTypeName(def.TypeName))
		cols = append(cols, core.ClassColumn{Name: def.Colname, TypeOID: t.typeOID, Type: t.expr, NotNull: def.IsNotNull})
	}
	return cols
}

func findFuncCall(l *ast.List) *ast.FuncCall {
	for _, item := range listItems(l) {
		switch v := item.(type) {
		case *ast.FuncCall:
			return v
		case *ast.List:
			if call := findFuncCall(v); call != nil {
				return call
			}
		}
	}
	return nil
}

// bindRangePivot binds PIVOT and UNPIVOT. The columns of the source that
// are not rotated pass through, followed by those the rotation makes, typed
// against the source.
func (a *analyzer) bindRangePivot(p *ast.RangePivot) (scopeRel, error) {
	src := &scope{}
	if err := a.appendFromItem(src, p.Source); err != nil {
		return scopeRel{}, err
	}
	defer a.binding(src)()
	rotated := map[string]bool{}
	var made []core.ClassColumn
	if p.Unpivot {
		// The value column has the type the rotated columns share, and is
		// nullable when any of them is.
		value := core.ClassColumn{Name: p.ValueColumn, NotNull: true}
		for i, item := range listItems(p.InList) {
			t, err := a.typeExpr(&ast.ColumnRef{Fields: &ast.List{Items: []ast.Node{item}}})
			if err != nil {
				return scopeRel{}, err
			}
			if i == 0 {
				value.TypeOID, value.Type = t.typeOID, t.expr
			}
			value.NotNull = value.NotNull && !t.nullable
			rotated[item.(*ast.String).Str] = true
		}
		name := a.lookupType(core.TypeExprOfTypeName(p.ForType))
		made = []core.ClassColumn{value, {Name: p.ForColumn, TypeOID: name.typeOID, Type: name.expr}}
	} else {
		// Each value in the list names a column of the aggregate over the
		// rows with that value, null when there are none.
		t, err := a.typeExpr(p.Aggregate)
		if err != nil {
			return scopeRel{}, err
		}
		rotated[p.ForColumn] = true
		for _, arg := range listItems(p.Aggregate.Args) {
			if cr, ok := arg.(*ast.ColumnRef); ok {
				fields := flattenFields(cr.Fields)
				rotated[fields[len(fields)-1]] = true
			}
		}
		for _, item := range listItems(p.InList) {
			made = append(made, core.ClassColumn{Name: item.(*ast.String).Str, TypeOID: t.typeOID, Type: t.expr})
		}
	}
	var rel scopeRel
	if p.Alias != nil && p.Alias.Aliasname != nil {
		rel.alias = *p.Alias.Aliasname
	}
	for _, r := range src.rels {
		for _, c := range r.cols {
			if !rotated[c.Name] {
				rel.cols = append(rel.cols, c)
			}
		}
	}
	rel.cols = append(rel.cols, made...)
	return rel, nil
}

// bindRangeSubselect binds a subquery used as a relation in FROM.
func (a *analyzer) bindRangeSubselect(rs *ast.RangeSubselect) (scopeRel, error) {
	sel, ok := rs.Subquery.(*ast.SelectStmt)
	if !ok {
		return scopeRel{}, fmt.Errorf("subquery: unsupported %T", rs.Subquery)
	}
	cols, err := a.subqueryColumns(sel)
	if err != nil {
		return scopeRel{}, err
	}
	alias := ""
	if rs.Alias != nil && rs.Alias.Aliasname != nil {
		alias = *rs.Alias.Aliasname
	}
	rel := derivedRel(alias, cols)
	if rs.Alias != nil {
		renameColumns(&rel, rs.Alias.Colnames)
	}
	return rel, nil
}

func (a *analyzer) bindRangeVar(rv *ast.RangeVar) (scopeRel, error) {
	if rv.Relname == nil {
		return scopeRel{}, fmt.Errorf("range var: missing relation name")
	}
	relName := *rv.Relname

	// A WITH clause defines a relation that shadows nothing but is not in the
	// catalog, so it is resolved before one.
	if cte, ok := a.ctes[relName]; ok && (rv.Schemaname == nil || *rv.Schemaname == "") {
		if rv.Alias != nil && rv.Alias.Aliasname != nil && *rv.Alias.Aliasname != "" {
			cte.alias = *rv.Alias.Aliasname
		}
		return cte, nil
	}

	schema := ""
	if rv.Schemaname != nil {
		schema = *rv.Schemaname
	}
	if schema == "" {
		schema = "public"
	}
	nsOID, err := a.cat.NamespaceOID(schema)
	if err != nil {
		return scopeRel{}, fmt.Errorf("schema %q: %w", schema, err)
	}
	classOID, err := a.cat.ClassOID(nsOID, relName)
	if err != nil {
		return scopeRel{}, fmt.Errorf("relation %q.%q: %w", schema, relName, err)
	}
	rel := scopeRel{
		alias:    relName,
		classOID: classOID,
	}
	if rv.Alias != nil && rv.Alias.Aliasname != nil && *rv.Alias.Aliasname != "" {
		rel.alias = *rv.Alias.Aliasname
	}

	cols, err := a.cat.ClassColumns(classOID)
	if err != nil {
		return scopeRel{}, err
	}
	rel.cols = cols
	return rel, nil
}

// resolveColumn finds a column in this query's scope, falling back to the scope
// of the query it is nested in, which is what a correlated subquery refers to.
func (a *analyzer) resolveColumn(relation, column string) (scopeRel, core.ClassColumn, bool, error) {
	rel, col, ok, err := a.scope.resolveColumn(relation, column)
	if err != nil || ok {
		return rel, col, ok, err
	}
	if a.outer != nil {
		return a.outer.resolveColumn(relation, column)
	}
	return rel, col, false, nil
}

// resolveColumn finds the single column named column, optionally qualified by
// relation. It reports an error when more than one relation in scope offers
// that name.
func (s *scope) resolveColumn(relation, column string) (rel scopeRel, col core.ClassColumn, ok bool, err error) {
	// A statement whose scope is not built yet offers no columns. Report the
	// column as unresolved and let the caller say so, rather than crashing on
	// the way to the same answer.
	if s == nil {
		return rel, col, false, nil
	}
	found := 0
	for _, r := range s.rels {
		if relation != "" && r.alias != relation {
			continue
		}
		for _, c := range r.cols {
			if c.Name != column {
				continue
			}
			found++
			if found > 1 {
				return scopeRel{}, core.ClassColumn{}, false, fmt.Errorf("ambiguous column reference %q", column)
			}
			rel, col = r, c
		}
	}
	if found == 0 {
		return scopeRel{}, core.ClassColumn{}, false, nil
	}
	return rel, col, true, nil
}
