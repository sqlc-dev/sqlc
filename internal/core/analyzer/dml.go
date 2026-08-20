package analyzer

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
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
	return a.projectReturning(s.ReturningList)
}

func (a *analyzer) analyzeMerge(s *ast.MergeStmt) error {
	if err := a.bindCTEs(s.WithClause); err != nil {
		return err
	}
	if s.Relation == nil {
		return fmt.Errorf("merge: missing target relation")
	}
	if s.SourceRelation == nil {
		return fmt.Errorf("merge: missing source relation")
	}

	sourceScope := &scope{}
	a.scope = sourceScope
	if err := a.appendFromItem(sourceScope, s.SourceRelation); err != nil {
		return fmt.Errorf("merge source: %w", err)
	}
	if len(sourceScope.rels) == 0 {
		return fmt.Errorf("merge: empty source relation")
	}
	if err := a.typeJoinConditions(s.SourceRelation); err != nil {
		return fmt.Errorf("merge source join: %w", err)
	}

	target, err := a.bindRangeVar(s.Relation)
	if err != nil {
		return fmt.Errorf("merge target: %w", err)
	}
	fullScope := mergeScope(sourceScope, target)
	targetScope := &scope{rels: []scopeRel{target}}

	a.scope = fullScope
	if _, err := a.typeExpr(s.JoinCondition); err != nil {
		return fmt.Errorf("merge ON: %w", err)
	}

	sourceMayBeMissing := false
	for i, item := range listItems(s.MergeWhenClauses) {
		clause, ok := item.(*ast.MergeWhenClause)
		if !ok {
			return fmt.Errorf("merge WHEN %d: unsupported clause %T", i+1, item)
		}

		switch clause.MatchKind {
		case ast.MergeWhenMatched:
			a.scope = fullScope
		case ast.MergeWhenNotMatchedByTarget:
			a.scope = sourceScope
		case ast.MergeWhenNotMatchedBySource:
			a.scope = targetScope
			sourceMayBeMissing = true
		default:
			return fmt.Errorf("merge WHEN %d: unsupported match kind %d", i+1, clause.MatchKind)
		}

		if _, err := a.typeExpr(clause.Condition); err != nil {
			return fmt.Errorf("merge WHEN %d condition: %w", i+1, err)
		}
		if !mergeActionAllowed(clause.MatchKind, clause.CommandType) {
			return fmt.Errorf("merge WHEN %d: action %d is invalid for match kind %d", i+1, clause.CommandType, clause.MatchKind)
		}

		switch clause.CommandType {
		case ast.CmdTypeUpdate:
			if err := a.bindMergeUpdate(target, clause.TargetList); err != nil {
				return err
			}
		case ast.CmdTypeInsert:
			if err := a.bindMergeInsert(target, clause.TargetList, clause.Values); err != nil {
				return err
			}
		case ast.CmdTypeDelete, ast.CmdTypeNothing:
			// These actions have no values to analyze.
		}
	}

	a.scope = mergeReturningScope(sourceScope, target, sourceMayBeMissing)
	return a.projectReturning(s.ReturningList)
}

func mergeScope(source *scope, target scopeRel) *scope {
	sc := &scope{rels: make([]scopeRel, 0, len(source.rels)+1)}
	sc.rels = append(sc.rels, source.rels...)
	sc.rels = append(sc.rels, target)
	return sc
}

func mergeReturningScope(source *scope, target scopeRel, sourceMayBeMissing bool) *scope {
	if !sourceMayBeMissing {
		return mergeScope(source, target)
	}
	sc := &scope{rels: make([]scopeRel, 0, len(source.rels)+1)}
	for _, rel := range source.rels {
		rel.cols = append([]core.ClassColumn(nil), rel.cols...)
		for i := range rel.cols {
			rel.cols[i].NotNull = false
		}
		sc.rels = append(sc.rels, rel)
	}
	sc.rels = append(sc.rels, target)
	return sc
}

func mergeActionAllowed(kind ast.MergeMatchKind, command ast.CmdType) bool {
	switch kind {
	case ast.MergeWhenMatched, ast.MergeWhenNotMatchedBySource:
		return command == ast.CmdTypeUpdate || command == ast.CmdTypeDelete || command == ast.CmdTypeNothing
	case ast.MergeWhenNotMatchedByTarget:
		return command == ast.CmdTypeInsert || command == ast.CmdTypeNothing
	default:
		return false
	}
}

func (a *analyzer) bindMergeUpdate(target scopeRel, targetList *ast.List) error {
	for i, item := range listItems(targetList) {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name == nil {
			return fmt.Errorf("merge update target %d: unsupported %T", i+1, item)
		}
		col, ok := findColumn(target, *rt.Name)
		if !ok {
			return fmt.Errorf("unknown column %q", *rt.Name)
		}
		value, err := mergeAssignmentValue(rt.Val)
		if err != nil {
			return fmt.Errorf("set %s: %w", *rt.Name, err)
		}
		if err := a.bindValue(target, &col, value); err != nil {
			return fmt.Errorf("set %s: %w", *rt.Name, err)
		}
	}
	return nil
}

func mergeAssignmentValue(value ast.Node) (ast.Node, error) {
	multi, ok := value.(*ast.MultiAssignRef)
	if !ok {
		return value, nil
	}
	row, ok := multi.Source.(*ast.RowExpr)
	if !ok {
		return nil, fmt.Errorf("unsupported multi-assignment source %T", multi.Source)
	}
	values := listItems(row.Args)
	if len(values) != multi.Ncolumns {
		return nil, fmt.Errorf("MERGE UPDATE has %d expressions for %d target columns", len(values), multi.Ncolumns)
	}
	index := multi.Colno - 1
	if index < 0 || index >= len(values) {
		return nil, fmt.Errorf("multi-assignment column %d is outside row of %d values", multi.Colno, len(values))
	}
	return values[index], nil
}

func (a *analyzer) bindMergeInsert(target scopeRel, targetList, values *ast.List) error {
	targetItems := listItems(targetList)
	if len(targetItems) == 0 && hasParamRef(values) {
		return fmt.Errorf("MERGE INSERT with parameters requires an explicit column list")
	}
	targets, err := insertTargets(target, targetList)
	if err != nil {
		return err
	}
	valueItems := listItems(values)
	if len(valueItems) > len(targets) {
		return fmt.Errorf("MERGE INSERT has more expressions than target columns")
	}
	if len(targetItems) > 0 && len(valueItems) > 0 && len(valueItems) < len(targets) {
		return fmt.Errorf("MERGE INSERT has fewer expressions than target columns")
	}
	for i, value := range valueItems {
		if err := a.bindValue(target, &targets[i], value); err != nil {
			return fmt.Errorf("merge insert value %d for %s: %w", i+1, targets[i].Name, err)
		}
	}
	return nil
}

func hasParamRef(list *ast.List) bool {
	if list == nil {
		return false
	}
	found := astutils.Search(list, func(node ast.Node) bool {
		_, ok := node.(*ast.ParamRef)
		return ok
	})
	return len(found.Items) > 0
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
		return rel.cols, nil
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
	// the statement's result, but the query still holds placeholders.
	if sel.ValuesLists == nil {
		_, err := a.subqueryColumns(sel)
		return err
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
	if _, ok := v.(*ast.SetToDefault); ok {
		return nil
	}
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
