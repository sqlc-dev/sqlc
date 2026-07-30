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
		params: map[int]*core.Parameter{},
	}
	switch s := stmt.(type) {
	case *ast.SelectStmt:
		if err := a.analyzeSelect(s); err != nil {
			return core.PrepareResult{}, err
		}
		a.command = core.CommandSelect
	default:
		return core.PrepareResult{}, fmt.Errorf("analyzer: unsupported statement %T", stmt)
	}
	return a.result(), nil
}

type analyzer struct {
	cat     *core.Catalog
	scope   *scope
	columns []core.Column
	params  map[int]*core.Parameter
	command core.Command
}

func (a *analyzer) result() core.PrepareResult {
	return core.PrepareResult{
		Command:    a.command,
		Columns:    a.columns,
		Parameters: orderedParams(a.params),
	}
}

func orderedParams(m map[int]*core.Parameter) []core.Parameter {
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
			out = append(out, *p)
		}
	}
	return out
}

func (a *analyzer) analyzeSelect(s *ast.SelectStmt) error {
	sc, err := a.buildScope(s.FromClause)
	if err != nil {
		return err
	}
	a.scope = sc

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
		return fmt.Errorf("select: empty target list")
	}
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
