package analyzer

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func Prepare(cat *core.Catalog, stmt ast.Node) (core.PrepareResult, error) {
	if rs, ok := stmt.(*ast.RawStmt); ok {
		stmt = rs.Stmt
	}
	a := &analyzer{
		cat:    cat,
		params: map[int]core.Parameter{},
	}
	switch s := stmt.(type) {
	case *ast.SelectStmt:
		if err := a.analyzeSelect(s); err != nil {
			return core.PrepareResult{}, err
		}
		a.command = core.CommandSelect
	case *ast.InsertStmt:
		if err := a.analyzeInsert(s); err != nil {
			return core.PrepareResult{}, err
		}
		a.command = core.CommandInsert
	case *ast.UpdateStmt:
		if err := a.analyzeUpdate(s); err != nil {
			return core.PrepareResult{}, err
		}
		a.command = core.CommandUpdate
	case *ast.DeleteStmt:
		if err := a.analyzeDelete(s); err != nil {
			return core.PrepareResult{}, err
		}
		a.command = core.CommandDelete
	case *ast.MergeStmt:
		if err := a.analyzeMerge(s); err != nil {
			return core.PrepareResult{}, err
		}
		a.command = core.CommandMerge
	default:
		return core.PrepareResult{}, fmt.Errorf("analyzer: unsupported statement %T", stmt)
	}
	return a.result(), nil
}

type analyzer struct {
	cat     *core.Catalog
	scope   *scope
	columns []core.Column
	params  map[int]core.Parameter
	command core.Command

	// outer is the scope of the query this one is nested in, which a
	// correlated subquery refers to.
	outer *scope

	// ctes are the relations a WITH clause defined, visible to this query and
	// to the ones nested in it.
	ctes map[string]scopeRel

	// aliases are the names the target list gives its expressions. GROUP BY,
	// HAVING and ORDER BY may refer to a result column by its output name.
	aliases map[string]ast.Node

	// resolving guards against an alias that refers to itself.
	resolving map[string]bool
}

// subquery analyzes a nested SELECT. It shares the parameter set, so a
// placeholder inside the subquery is reported with the rest, and it sees this
// query's scope, so a correlated reference resolves.
func (a *analyzer) subquery(s *ast.SelectStmt) (*analyzer, error) {
	sub := &analyzer{
		cat:    a.cat,
		params: a.params,
		outer:  a.scope,
		ctes:   a.ctes,
	}
	if err := sub.analyzeSelect(s); err != nil {
		return nil, err
	}
	return sub, nil
}

func (a *analyzer) subqueryColumns(s *ast.SelectStmt) ([]core.Column, error) {
	sub, err := a.subquery(s)
	if err != nil {
		return nil, err
	}
	return sub.columns, nil
}

// derivedRel turns a subquery's result columns into a relation, so that the
// query selecting from it resolves columns the same way as from a table.
func derivedRel(alias string, cols []core.Column) scopeRel {
	rel := scopeRel{alias: alias, cols: make([]core.ClassColumn, 0, len(cols))}
	for _, col := range cols {
		rel.cols = append(rel.cols, core.ClassColumn{
			AttOID:  col.SourceAttributeOID,
			Name:    col.Name,
			TypeOID: col.TypeOID,
			NotNull: col.NotNull,
		})
	}
	return rel
}

func (a *analyzer) result() core.PrepareResult {
	return core.PrepareResult{
		Command:    a.command,
		Columns:    a.columns,
		Parameters: orderedParams(a.params),
	}
}

func orderedParams(m map[int]core.Parameter) []core.Parameter {
	if len(m) == 0 {
		return nil
	}
	maxN := 0
	for n := range m {
		if n > maxN {
			maxN = n
		}
	}
	out := make([]core.Parameter, 0, len(m))
	for i := 1; i <= maxN; i++ {
		if p, ok := m[i]; ok {
			out = append(out, p)
		}
	}
	return out
}

func (a *analyzer) analyzeSelect(s *ast.SelectStmt) error {
	if err := a.bindCTEs(s.WithClause); err != nil {
		return err
	}

	// A set operation reports the columns of its first branch. The other
	// branches still get analyzed, for the placeholders they hold.
	if s.Larg != nil && s.Rarg != nil {
		if err := a.analyzeSetOperation(s); err != nil {
			return err
		}
		return nil
	}

	sc, err := a.buildScope(s.FromClause)
	if err != nil {
		return err
	}
	a.scope = sc

	a.bindAliases(s.TargetList)

	for _, item := range listItems(s.FromClause) {
		if err := a.typeJoinConditions(item); err != nil {
			return fmt.Errorf("join: %w", err)
		}
	}

	if s.WhereClause != nil {
		if _, err := a.typeExpr(s.WhereClause); err != nil {
			return fmt.Errorf("where: %w", err)
		}
	}
	if items := listItems(s.GroupClause); items != nil {
		for _, g := range items {
			if _, err := a.typeExpr(g); err != nil {
				return fmt.Errorf("group by: %w", err)
			}
		}
	}
	if s.HavingClause != nil {
		if _, err := a.typeExpr(s.HavingClause); err != nil {
			return fmt.Errorf("having: %w", err)
		}
	}

	targets := listItems(s.TargetList)
	if targets == nil {
		// A VALUES list produces rows without naming columns.
		if s.ValuesLists != nil {
			return a.typeValuesLists(s.ValuesLists)
		}
		return fmt.Errorf("select: empty target list")
	}
	a.columns = make([]core.Column, 0, len(targets))
	for _, t := range targets {
		rt, ok := t.(*ast.ResTarget)
		if !ok {
			continue
		}
		if err := a.projectTarget(rt); err != nil {
			return err
		}
	}
	return nil
}

func (a *analyzer) typeValuesLists(l *ast.List) error {
	for _, row := range listItems(l) {
		if _, err := a.typeExpr(row); err != nil {
			return err
		}
	}
	return nil
}

// analyzeSetOperation types both sides of UNION, INTERSECT or EXCEPT and
// reports the first branch's columns, the way the databases name the result.
func (a *analyzer) analyzeSetOperation(s *ast.SelectStmt) error {
	left, err := a.subquery(s.Larg)
	if err != nil {
		return err
	}
	if _, err := a.subquery(s.Rarg); err != nil {
		return err
	}
	a.columns = left.columns
	a.scope = left.scope
	return nil
}

// bindCTEs analyzes a WITH clause, making each common table expression
// available to the query as a relation.
func (a *analyzer) bindCTEs(with *ast.WithClause) error {
	if with == nil {
		return nil
	}
	for _, item := range listItems(with.Ctes) {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok || cte.Ctename == nil {
			continue
		}
		sel, ok := cte.Ctequery.(*ast.SelectStmt)
		if !ok {
			continue
		}
		cols, err := a.subqueryColumns(sel)
		if err != nil {
			return fmt.Errorf("with %s: %w", *cte.Ctename, err)
		}
		rel := derivedRel(*cte.Ctename, cols)
		renameColumns(&rel, cte.Aliascolnames)
		if a.ctes == nil {
			a.ctes = map[string]scopeRel{}
		}
		a.ctes[*cte.Ctename] = rel
	}
	return nil
}

// bindAliases records the output names the target list assigns, so a later
// clause can refer to a result column by name.
func (a *analyzer) bindAliases(targets *ast.List) {
	for _, item := range listItems(targets) {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name == nil || *rt.Name == "" {
			continue
		}
		if a.aliases == nil {
			a.aliases = map[string]ast.Node{}
		}
		a.aliases[*rt.Name] = rt.Val
	}
}

// typeAlias types the expression an output name stands for, which is how a
// GROUP BY or HAVING clause refers to a result column.
func (a *analyzer) typeAlias(name string) (exprType, bool, error) {
	val, ok := a.aliases[name]
	if !ok || a.resolving[name] {
		return exprType{}, false, nil
	}
	if a.resolving == nil {
		a.resolving = map[string]bool{}
	}
	a.resolving[name] = true
	defer delete(a.resolving, name)

	t, err := a.typeExpr(val)
	if err != nil {
		return exprType{}, false, err
	}
	return t, true, nil
}

// renameColumns applies the column aliases a derived relation was declared
// with, as in "WITH t(a, b) AS (...)".
func renameColumns(rel *scopeRel, names *ast.List) {
	for i, item := range listItems(names) {
		if i >= len(rel.cols) {
			break
		}
		if s, ok := item.(*ast.String); ok && s.Str != "" {
			rel.cols[i].Name = s.Str
		}
	}
}

func listItems(l *ast.List) []ast.Node {
	if l == nil {
		return nil
	}
	return l.Items
}

func (a *analyzer) typeJoinConditions(item ast.Node) error {
	je, ok := item.(*ast.JoinExpr)
	if !ok {
		return nil
	}
	if err := a.typeJoinConditions(je.Larg); err != nil {
		return err
	}
	if err := a.typeJoinConditions(je.Rarg); err != nil {
		return err
	}
	if je.Quals != nil {
		if _, err := a.typeExpr(je.Quals); err != nil {
			return fmt.Errorf("ON: %w", err)
		}
	}
	return nil
}
