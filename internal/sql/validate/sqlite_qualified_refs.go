package validate

import (
	"fmt"
	"reflect"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// ValidateSQLiteQualifiedColumnRefs validates that qualified column references
// only use visible tables/aliases in the current or outer SELECT scopes.
func ValidateSQLiteQualifiedColumnRefs(root ast.Node) error {
	return validateNodeSQLite(root, nil)
}

type scope struct {
	parent *scope
	names  map[string]struct{}
}

func newScope(parent *scope) *scope {
	return &scope{parent: parent, names: map[string]struct{}{}}
}

func (s *scope) add(name string) {
	if name == "" {
		return
	}
	s.names[name] = struct{}{}
}

func (s *scope) has(name string) bool {
	for cur := s; cur != nil; cur = cur.parent {
		if _, ok := cur.names[name]; ok {
			return true
		}
	}
	return false
}

func stringSlice(list *ast.List) []string {
	if list == nil {
		return nil
	}
	out := make([]string, 0, len(list.Items))
	for _, it := range list.Items {
		if s, ok := it.(*ast.String); ok {
			out = append(out, s.Str)
		}
	}
	return out
}

func qualifierFromColumnRef(ref *ast.ColumnRef) (string, bool) {
	if ref == nil || ref.Fields == nil {
		return "", false
	}
	items := stringSlice(ref.Fields)
	switch len(items) {
	case 2:
		return items[0], true
	case 3:
		return items[1], true
	default:
		return "", false
	}
}

func addFromItemToScope(sc *scope, n ast.Node) {
	switch t := n.(type) {
	case *ast.RangeVar:
		if t.Relname != nil {
			sc.add(*t.Relname)
		}
		if t.Alias != nil && t.Alias.Aliasname != nil {
			sc.add(*t.Alias.Aliasname)
		}
	case *ast.JoinExpr:
		addFromItemToScope(sc, t.Larg)
		addFromItemToScope(sc, t.Rarg)
	case *ast.RangeSubselect:
		if t.Alias != nil && t.Alias.Aliasname != nil {
			sc.add(*t.Alias.Aliasname)
		}
	case *ast.RangeFunction:
		if t.Alias != nil && t.Alias.Aliasname != nil {
			sc.add(*t.Alias.Aliasname)
		}
	}
}

func validateNodeSQLite(node ast.Node, parent *scope) error {
	switch n := node.(type) {
	case *ast.SelectStmt:
		sc := newScope(parent)
		if n.FromClause != nil {
			for _, item := range n.FromClause.Items {
				addFromItemToScope(sc, item)
			}
		}
		return walkSQLite(n, sc)
	default:
		return nil
	}
}

func walkSQLite(node ast.Node, sc *scope) error {
	if node == nil {
		return nil
	}

	if ref, ok := node.(*ast.ColumnRef); ok {
		if qual, ok := qualifierFromColumnRef(ref); ok && !sc.has(qual) {
			return &sqlerr.Error{
				Code:     "42703",
				Message:  fmt.Sprintf("table alias %q does not exist", qual),
				Location: ref.Location,
			}
		}
	}

	switch n := node.(type) {
	case *ast.SubLink:
		if n.Subselect != nil {
			return validateNodeSQLite(n.Subselect, sc)
		}
		return nil
	case *ast.RangeSubselect:
		if n.Subquery != nil {
			return validateNodeSQLite(n.Subquery, sc)
		}
		return nil
	}

	return walkSQLiteReflect(node, sc)
}

func walkSQLiteReflect(node ast.Node, sc *scope) error {
	v := reflect.ValueOf(node)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).PkgPath != "" {
			continue
		}
		f := v.Field(i)
		if !f.IsValid() {
			continue
		}

		for f.Kind() == reflect.Pointer {
			if f.IsNil() {
				goto next
			}
			f = f.Elem()
		}

		if f.Type() == reflect.TypeOf(ast.List{}) {
			list := f.Addr().Interface().(*ast.List)
			for _, n := range list.Items {
				if err := walkSQLite(n, sc); err != nil {
					return err
				}
			}
			continue
		}

		if f.CanAddr() {
			if pl, ok := f.Addr().Interface().(**ast.List); ok && *pl != nil {
				for _, n := range (*pl).Items {
					if err := walkSQLite(n, sc); err != nil {
						return err
					}
				}
				continue
			}
		}

		if f.CanInterface() {
			if n, ok := f.Interface().(ast.Node); ok {
				if err := walkSQLite(n, sc); err != nil {
					return err
				}
				continue
			}
		}

		if f.Kind() == reflect.Slice {
			for j := 0; j < f.Len(); j++ {
				elem := f.Index(j)
				if elem.Kind() == reflect.Pointer && elem.IsNil() {
					continue
				}
				if elem.CanInterface() {
					if n, ok := elem.Interface().(ast.Node); ok {
						if err := walkSQLite(n, sc); err != nil {
							return err
						}
					}
				}
			}
		}

	next:
	}
	return nil
}
