package analyzer

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
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
	if s.OnDuplicateKeyUpdate != nil {
		if err := a.bindAssignments(rel, s.OnDuplicateKeyUpdate.TargetList); err != nil {
			return fmt.Errorf("on duplicate key update: %w", err)
		}
	}
	return a.projectReturning(s.ReturningList)
}

// bindAssignments types the values a SET-style list assigns to the
// relation's columns.
func (a *analyzer) bindAssignments(rel scopeRel, targets *ast.List) error {
	for _, item := range listItems(targets) {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name == nil {
			continue
		}
		col, ok := findColumn(rel, *rt.Name)
		if !ok {
			return fmt.Errorf("unknown column %q", *rt.Name)
		}
		if err := a.bindValue(rel, &col, rt.Val); err != nil {
			return fmt.Errorf("set %s: %w", *rt.Name, err)
		}
	}
	return nil
}

// analyzeCall types the arguments of a CALL from the procedure's declared
// parameters: a placeholder takes the parameter's type and name, whether
// it is passed by position or by name.
func (a *analyzer) analyzeCall(s *ast.CallStmt) error {
	if s.FuncCall == nil {
		return fmt.Errorf("call: missing procedure")
	}
	name := funcCallName(s.FuncCall)
	overloads, err := a.cat.FindProcs(name, nil)
	if err != nil {
		return err
	}
	args := listItems(s.FuncCall.Args)
	var procs []core.ProcOverload
	for _, p := range overloads {
		if p.Kind == "p" {
			procs = append(procs, p)
		}
	}
	if len(procs) == 0 {
		return fmt.Errorf("unknown procedure %q", name)
	}
	proc := procs[0]
	for _, p := range procs {
		if len(p.ArgTypes) == len(args) {
			proc = p
			break
		}
	}
	params, err := a.cat.ProcArgs(proc.OID)
	if err != nil {
		return err
	}
	for i, arg := range args {
		var param *core.ProcArg
		value := arg
		if named, ok := arg.(*ast.NamedArgExpr); ok {
			value = named.Arg
			if named.Name != nil {
				for j := range params {
					if params[j].Name == *named.Name {
						param = &params[j]
						break
					}
				}
			}
		} else if i < len(params) {
			param = &params[i]
		}
		pr, ok := value.(*ast.ParamRef)
		if !ok || param == nil {
			if _, err := a.typeExpr(value); err != nil {
				return err
			}
			continue
		}
		a.inferParam(pr.Number, exprType{typeOID: param.TypeOID})
		a.nameParam(pr.Number, param.Name)
	}
	return nil
}

func (a *analyzer) analyzeUpdate(s *ast.UpdateStmt) error {
	if err := a.bindCTEs(s.WithClause); err != nil {
		return err
	}
	sc, err := a.relationScope(s.Relations, s.FromClause, nil)
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
	if err := a.typeLimit(s.LimitCount); err != nil {
		return fmt.Errorf("limit: %w", err)
	}
	return a.projectReturning(s.ReturningList)
}

func (a *analyzer) analyzeDelete(s *ast.DeleteStmt) error {
	if err := a.bindCTEs(s.WithClause); err != nil {
		return err
	}
	sc, err := a.relationScope(s.Relations, s.UsingClause, s.FromClause)
	if err != nil {
		return err
	}
	a.scope = sc

	if s.WhereClause != nil {
		if _, err := a.typeExpr(s.WhereClause); err != nil {
			return fmt.Errorf("where: %w", err)
		}
	}
	if err := a.typeLimit(s.LimitCount); err != nil {
		return fmt.Errorf("limit: %w", err)
	}
	return a.projectReturning(s.ReturningList)
}

// relationScope builds the scope a DML statement operates on: the relations it
// targets, whatever a USING or FROM clause joins in, and — for the engines that
// report a multi-table DELETE that way — a single FROM node.
func (a *analyzer) relationScope(relations, extra *ast.List, from ast.Node) (*scope, error) {
	sc := &scope{}
	defer a.binding(sc)()
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
	if from != nil {
		if err := a.appendFromItem(sc, from); err != nil {
			return nil, err
		}
	}
	if len(sc.rels) == 0 {
		return nil, fmt.Errorf("missing target relation")
	}
	return sc, nil
}

func insertTargets(rel scopeRel, cols *ast.List) ([]core.ClassColumn, error) {
	items := listItems(cols)
	if len(items) == 0 {
		out := make([]core.ClassColumn, 0, len(rel.cols))
		for _, col := range rel.cols {
			if !col.Hidden {
				out = append(out, col)
			}
		}
		return out, nil
	}
	out := make([]core.ClassColumn, 0, len(items))
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

func (a *analyzer) bindInsertValues(n ast.Node, rel scopeRel, targets []core.ClassColumn) error {
	if n == nil {
		return nil
	}
	sel, ok := n.(*ast.SelectStmt)
	if !ok {
		return fmt.Errorf("insert: unsupported source %T", n)
	}
	// INSERT ... SELECT inserts whatever the query returns. The rows are not
	// the statement's result, but the query still holds placeholders, and a
	// placeholder selected directly stands in for the column it lands in.
	if len(listItems(sel.ValuesLists)) == 0 {
		if _, err := a.subqueryColumns(sel); err != nil {
			return err
		}
		for i, item := range listItems(sel.TargetList) {
			rt, ok := item.(*ast.ResTarget)
			if !ok || i >= len(targets) {
				continue
			}
			if pr, ok := rt.Val.(*ast.ParamRef); ok {
				a.inferParam(pr.Number, columnType(rel, targets[i]))
			}
		}
		return nil
	}
	for _, row := range listItems(sel.ValuesLists) {
		values, ok := row.(*ast.List)
		if !ok {
			continue
		}
		for i, v := range values.Items {
			var target *core.ClassColumn
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

func (a *analyzer) bindValue(rel scopeRel, target *core.ClassColumn, v ast.Node) error {
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

func findColumn(rel scopeRel, name string) (core.ClassColumn, bool) {
	for _, col := range rel.cols {
		if col.Name == name {
			return col, true
		}
	}
	return core.ClassColumn{}, false
}

func columnType(rel scopeRel, col core.ClassColumn) exprType {
	return exprType{
		typeOID:            col.TypeOID,
		nullable:           !col.NotNull,
		sourceClassOID:     rel.classOID,
		sourceAttributeOID: col.AttOID,
		sourceTableAlias:   rel.alias,
	}
}
