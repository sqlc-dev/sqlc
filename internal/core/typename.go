package core

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TypeExprOfTypeName reads the type an AST node names into an expression.
// An engine that folds the whole type into a spelling — ClickHouse's
// Array(Nullable(String)), SQLite's VARYING CHARACTER(10) — hands it over
// in Spelling and the spelling is read as written. Otherwise the name comes
// from Name or the qualifying parts of Names, the type modifiers become
// integer or string arguments, and each array bound wraps the result in an
// array.
func TypeExprOfTypeName(tn *ast.TypeName) *TypeExpr {
	if tn == nil {
		return nil
	}
	if tn.Spelling != "" {
		return ParseTypeExpr(tn.Spelling)
	}
	name := strings.TrimSpace(tn.Name)
	if name == "" && tn.Names != nil {
		parts := make([]string, 0, len(tn.Names.Items))
		for _, item := range tn.Names.Items {
			s, ok := item.(*ast.String)
			if !ok || s.Str == "pg_catalog" {
				continue
			}
			parts = append(parts, s.Str)
		}
		name = strings.Join(parts, ".")
	}
	name = strings.ToLower(name)
	if name == "" {
		return nil
	}
	// A name an engine spelled with its own arguments or array suffix reads
	// the same way a spelling does.
	t := ParseTypeExpr(name)
	for _, item := range listItems(tn.Typmods) {
		if arg, ok := typmodArg(item); ok {
			t.Args = append(t.Args, arg)
		}
	}
	for range listItems(tn.ArrayBounds) {
		t = Array(t)
	}
	return t
}

// ColumnTypeExpr reads a column definition's type. Engines report an array
// column either on the type name or on the column itself.
func ColumnTypeExpr(col *ast.ColumnDef) *TypeExpr {
	if col == nil {
		return nil
	}
	t := TypeExprOfTypeName(col.TypeName)
	if t == nil {
		return nil
	}
	// MySQL reports an unsigned column on the definition, and the
	// members of an enum or set apart from the type's name.
	if col.IsUnsigned && !strings.Contains(t.Name, " unsigned") {
		t.Name += " unsigned"
	}
	if vals := listItems(col.Vals); len(vals) > 0 && len(t.Args) == 0 {
		for _, item := range vals {
			if s, ok := item.(*ast.String); ok {
				v := s.Str
				t.Args = append(t.Args, TypeArg{String: &v})
			}
		}
	}
	if col.TypeName.Spelling != "" || listItems(col.TypeName.ArrayBounds) != nil {
		return t
	}
	dims := col.ArrayDims
	if dims == 0 && col.IsArray {
		dims = 1
	}
	for i := 0; i < dims; i++ {
		t = Array(t)
	}
	return t
}

// typmodArg reads one type modifier as an argument: an integer, a quoted
// string, or a bare word, which is an identifier such as the max of
// nvarchar(max) or the day to second of an interval. A constant node holds
// a literal; a bare String node holds a word.
func typmodArg(n ast.Node) (TypeArg, bool) {
	switch v := n.(type) {
	case *ast.A_Const:
		switch val := v.Val.(type) {
		case *ast.String:
			s := val.Str
			return TypeArg{String: &s}, true
		default:
			return typmodArg(v.Val)
		}
	case *ast.Integer:
		i := v.Ival
		return TypeArg{Int: &i}, true
	case *ast.String:
		s := strings.ToLower(v.Str)
		return TypeArg{Ident: &s}, true
	case *ast.ColumnRef:
		parts := make([]string, 0, len(listItems(v.Fields)))
		for _, item := range listItems(v.Fields) {
			if s, ok := item.(*ast.String); ok {
				parts = append(parts, s.Str)
			}
		}
		if len(parts) == 0 {
			return TypeArg{}, false
		}
		ident := strings.ToLower(strings.Join(parts, "."))
		return TypeArg{Ident: &ident}, true
	}
	return TypeArg{}, false
}

func listItems(l *ast.List) []ast.Node {
	if l == nil {
		return nil
	}
	return l.Items
}

// ResolveType interns the type an AST node names, registering it when the
// dialect's seed did not: a schema is free to declare types of its own.
func (c *Catalog) ResolveType(tn *ast.TypeName) (int64, error) {
	t := TypeExprOfTypeName(tn)
	if t == nil {
		return 0, fmt.Errorf("missing type name")
	}
	return c.ResolveTypeExpr(t)
}
