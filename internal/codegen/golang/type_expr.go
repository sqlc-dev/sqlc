package golang

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// columnType is the type a mapper that reads the type expression starts
// from, and whether a value of it may be null. The expression is only
// present when the analysis core typed the column; without it the flat
// description stands in, the type name wrapped in one array per dimension,
// with the elements taken as not null the way the legacy mappers do.
func columnType(col *plugin.Column) (*plugin.TypeExpr, bool) {
	if col.TypeExpr != nil {
		return col.TypeExpr, !col.NotNull
	}
	t := &plugin.TypeExpr{Name: col.Type.GetName()}
	dims := int(col.ArrayDims)
	if col.IsArray && dims == 0 {
		dims = 1
	}
	for range dims {
		t = &plugin.TypeExpr{Name: "array", Args: []*plugin.TypeArg{{Value: &plugin.TypeArg_Type{Type: t}}}}
	}
	return t, !col.NotNull
}

// typeExprName is the type's name as the mappers switch on it.
func typeExprName(t *plugin.TypeExpr) string {
	return strings.ToLower(t.GetName())
}

// typeExprArg is the i-th argument that is itself a type, or nil.
func typeExprArg(t *plugin.TypeExpr, i int) *plugin.TypeExpr {
	n := 0
	for _, arg := range t.GetArgs() {
		if arg.GetType() == nil {
			continue
		}
		if n == i {
			return arg.GetType()
		}
		n++
	}
	return nil
}

// typeExprLastArg is the last argument that is a type, or nil.
func typeExprLastArg(t *plugin.TypeExpr) *plugin.TypeExpr {
	var last *plugin.TypeExpr
	for _, arg := range t.GetArgs() {
		if arg.GetType() != nil {
			last = arg.GetType()
		}
	}
	return last
}

// nullableGoType is typ for a value that cannot be null, and for one that
// can, a pointer to it when the options say so or when nullTyp, the
// Null-prefixed wrapper the type has, is empty.
func nullableGoType(options *opts.Options, nullable bool, typ, nullTyp string) string {
	if !nullable {
		return typ
	}
	if options.EmitPointersForNullTypes || nullTyp == "" {
		return "*" + typ
	}
	return nullTyp
}
