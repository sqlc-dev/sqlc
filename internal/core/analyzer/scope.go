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

	cols, err := a.cat.ClassColumns(classOID)
	if err != nil {
		return scopeRel{}, err
	}
	rel.cols = cols
	return rel, nil
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
