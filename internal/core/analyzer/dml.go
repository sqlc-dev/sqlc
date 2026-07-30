package analyzer

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func (a *analyzer) analyzeInsert(s *ast.InsertStmt) error {
	if s.Relation == nil {
		return fmt.Errorf("insert: missing relation")
	}
	rel, err := a.bindRangeVar(s.Relation)
	if err != nil {
		return err
	}
	a.scope = &scope{rels: []scopeRel{rel}}

	targets, err := insertTargets(rel, s.Cols)
	if err != nil {
		return err
	}
	if err := a.bindInsertValues(s.SelectStmt, rel, targets); err != nil {
		return err
	}
	return a.projectReturning(s.ReturningList)
}

func (a *analyzer) analyzeUpdate(s *ast.UpdateStmt) error {
	sc, err := a.relationScope(s.Relations, s.FromClause)
	if err != nil {
		return err
	}
	a.scope = sc
	target := sc.rels[0]

	for _, item := range listItems(s.TargetList) {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name == nil {
			continue
		}
		col, ok := findColumn(target, *rt.Name)
		if !ok {
			return fmt.Errorf("unknown column %q", *rt.Name)
		}
		if err := a.bindValue(target, &col, rt.Val); err != nil {
			return fmt.Errorf("set %s: %w", *rt.Name, err)
		}
	}

	if s.WhereClause != nil {
		if _, err := a.typeExpr(s.WhereClause); err != nil {
			return fmt.Errorf("where: %w", err)
		}
	}
	return a.projectReturning(s.ReturningList)
}

func (a *analyzer) analyzeDelete(s *ast.DeleteStmt) error {
	sc, err := a.relationScope(s.Relations, s.UsingClause)
	if err != nil {
		return err
	}
	a.scope = sc

	if s.WhereClause != nil {
		if _, err := a.typeExpr(s.WhereClause); err != nil {
			return fmt.Errorf("where: %w", err)
		}
	}
	return a.projectReturning(s.ReturningList)
}

func (a *analyzer) relationScope(relations, extra *ast.List) (*scope, error) {
	sc := &scope{}
	for _, item := range listItems(relations) {
		if err := a.appendFromItem(sc, item); err != nil {
			return nil, err
		}
	}
	for _, item := range listItems(extra) {
		if err := a.appendFromItem(sc, item); err != nil {
			return nil, err
		}
	}
	if len(sc.rels) == 0 {
		return nil, fmt.Errorf("missing target relation")
	}
	return sc, nil
}

func insertTargets(rel scopeRel, cols *ast.List) ([]scopeCol, error) {
	items := listItems(cols)
	if len(items) == 0 {
		return rel.cols, nil
	}
	out := make([]scopeCol, 0, len(items))
	for _, item := range items {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name == nil {
			return nil, fmt.Errorf("insert: unsupported column target %T", item)
		}
		col, ok := findColumn(rel, *rt.Name)
		if !ok {
			return nil, fmt.Errorf("unknown column %q", *rt.Name)
		}
		out = append(out, col)
	}
	return out, nil
}

func (a *analyzer) bindInsertValues(n ast.Node, rel scopeRel, targets []scopeCol) error {
	if n == nil {
		return nil
	}
	sel, ok := n.(*ast.SelectStmt)
	if !ok {
		return fmt.Errorf("insert: unsupported source %T", n)
	}
	if sel.ValuesLists == nil {
		return fmt.Errorf("insert: INSERT ... SELECT is not supported")
	}
	for _, row := range listItems(sel.ValuesLists) {
		values, ok := row.(*ast.List)
		if !ok {
			continue
		}
		for i, v := range values.Items {
			var target *scopeCol
			if i < len(targets) {
				target = &targets[i]
			}
			if err := a.bindValue(rel, target, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *analyzer) bindValue(rel scopeRel, target *scopeCol, v ast.Node) error {
	if target != nil {
		switch value := v.(type) {
		case *ast.ParamRef:
			a.inferParam(value.Number, columnType(rel, *target))
			return nil
		case *ast.A_Const:
			return nil
		}
	}
	_, err := a.typeExpr(v)
	return err
}

func (a *analyzer) projectReturning(l *ast.List) error {
	for _, item := range listItems(l) {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		if err := a.projectTarget(rt); err != nil {
			return err
		}
	}
	return nil
}

func findColumn(rel scopeRel, name string) (scopeCol, bool) {
	for _, col := range rel.cols {
		if col.name == name {
			return col, true
		}
	}
	return scopeCol{}, false
}

func columnType(rel scopeRel, col scopeCol) exprType {
	return exprType{
		typeOID:            col.typeOID,
		nullable:           !col.notNull,
		sourceClassOID:     rel.classOID,
		sourceAttributeOID: col.attOID,
		sourceTableAlias:   rel.alias,
	}
}
