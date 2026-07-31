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
	for _, item := range items {
		if err := a.appendFromItem(sc, item); err != nil {
			return nil, err
		}
	}
	return sc, nil
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
		if err := a.appendFromItem(sc, v.Larg); err != nil {
			return err
		}
		return a.appendFromItem(sc, v.Rarg)
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
	// gives a placeholder passed to the function its type.
	if _, err := a.typeFuncCall(call); err != nil {
		return scopeRel{}, err
	}

	p := overloads[0]
	rel := scopeRel{
		alias: name,
		cols:  []core.ClassColumn{{Name: name, TypeOID: p.ReturnTypeOID, NotNull: !p.ReturnNullable}},
	}
	if rf.Alias != nil && rf.Alias.Aliasname != nil && *rf.Alias.Aliasname != "" {
		rel.alias = *rf.Alias.Aliasname
	}
	return rel, nil
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
