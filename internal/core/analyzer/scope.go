package analyzer

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// scope is the set of relations available to expression resolution at a
// given point in a query. For now it holds at most one parent (FROM
// list) — subquery / CTE / lateral parents are a follow-up.
type scope struct {
	rels []scopeRel
}

// scopeRel binds a FROM-list relation (table or aliased table) to the
// columns it contributes.
type scopeRel struct {
	alias    string // user-visible name for resolution; defaults to relation name
	classOID int64
	cols     []scopeCol
}

// scopeCol is a single column visible through a scopeRel.
type scopeCol struct {
	name   string
	attOID int64
	typeOID int64
	notNull bool
}

// buildScope walks a FROM clause (currently: RangeVar plus JoinExpr
// trees, which we flatten into the rels slice) and produces a flat
// scope. Subqueries and table-valued functions are still TODO.
func (a *analyzer) buildScope(from *ast.List) (*scope, error) {
	sc := &scope{}
	for _, item := range listItems(from) {
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
		// Flatten both sides into the same scope. The join condition
		// (Quals / USING) is type-checked against the assembled scope
		// in analyzeSelect after buildScope returns; we don't enforce
		// it here. Outer-join nullability is also a follow-up.
		if err := a.appendFromItem(sc, v.Larg); err != nil {
			return err
		}
		return a.appendFromItem(sc, v.Rarg)
	default:
		return fmt.Errorf("scope: unsupported FROM item %T", item)
	}
}

func (a *analyzer) bindRangeVar(rv *ast.RangeVar) (scopeRel, error) {
	if rv.Relname == nil {
		return scopeRel{}, fmt.Errorf("range var: missing relation name")
	}
	relName := *rv.Relname
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

	cols, err := a.classColumns(classOID)
	if err != nil {
		return scopeRel{}, err
	}
	rel.cols = cols
	return rel, nil
}

// classColumns reads the column layout of a relation directly from
// sql_attribute. We can't go through TableColumns because that resolves
// by name and doesn't return attribute OIDs.
func (a *analyzer) classColumns(classOID int64) ([]scopeCol, error) {
	rows, err := a.cat.DB().Query(
		`SELECT a.oid, a.name, a.type_oid, a.not_null
		 FROM sql_attribute a WHERE a.class_oid = ? ORDER BY a.num`,
		classOID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []scopeCol
	for rows.Next() {
		var c scopeCol
		var nn int
		if err := rows.Scan(&c.attOID, &c.name, &c.typeOID, &nn); err != nil {
			return nil, err
		}
		c.notNull = nn != 0
		out = append(out, c)
	}
	return out, rows.Err()
}

// resolveColumn locates a (relation, column) pair in the scope.
// relation may be empty for an unqualified reference; in that case we
// search every relation and error on ambiguity.
func (s *scope) resolveColumn(relation, column string) (rel scopeRel, col scopeCol, ok bool, err error) {
	var matches []struct {
		rel scopeRel
		col scopeCol
	}
	for _, r := range s.rels {
		if relation != "" && r.alias != relation {
			continue
		}
		for _, c := range r.cols {
			if c.name == column {
				matches = append(matches, struct {
					rel scopeRel
					col scopeCol
				}{r, c})
			}
		}
	}
	if len(matches) == 0 {
		return scopeRel{}, scopeCol{}, false, nil
	}
	if len(matches) > 1 {
		return scopeRel{}, scopeCol{}, false, fmt.Errorf("ambiguous column reference %q", column)
	}
	return matches[0].rel, matches[0].col, true, nil
}

// silence unused-import warning if core not referenced here directly
var _ = core.Column{}
