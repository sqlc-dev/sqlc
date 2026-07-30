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
	cols     []scopeCol
}

type scopeCol struct {
	name    string
	attOID  int64
	typeOID int64
	notNull bool
}

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

func (a *analyzer) classColumns(classOID int64) ([]scopeCol, error) {
	cols, err := a.cat.ClassColumns(classOID)
	if err != nil {
		return nil, err
	}
	out := make([]scopeCol, 0, len(cols))
	for _, c := range cols {
		out = append(out, scopeCol{
			name:    c.Name,
			attOID:  c.AttOID,
			typeOID: c.TypeOID,
			notNull: c.NotNull,
		})
	}
	return out, nil
}

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

var _ = core.Column{}
