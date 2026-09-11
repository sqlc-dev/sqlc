package dolphin

import (
	"embed"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// The dialect directory describes MySQL's type system, and functions.jsonl is
// its standard library. Both the analysis core and the catalog the legacy
// compiler builds read them.
//
//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds MySQL's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}

func init() {
	core.RegisterCanonicalizer("mysql", canonicalize)
}

// canonicalize rewrites a type the way MySQL stores it: a decimal declared
// without a precision is decimal(10,0), a float declared with one is a
// float or a double depending on it, and a boolean is a tinyint(1).
func canonicalize(t *core.TypeExpr) *core.TypeExpr {
	switch {
	case (t.Name == "decimal" || t.Name == "decimal unsigned") && len(t.Args) == 0:
		p, s := int64(10), int64(0)
		return &core.TypeExpr{Name: t.Name, Nullable: t.Nullable, Args: []core.TypeArg{{Int: &p}, {Int: &s}}}
	case (t.Name == "decimal" || t.Name == "decimal unsigned") && len(t.Args) == 1 && t.Args[0].Int != nil:
		s := int64(0)
		return &core.TypeExpr{Name: t.Name, Nullable: t.Nullable, Args: []core.TypeArg{t.Args[0], {Int: &s}}}
	case t.Name == "float" && len(t.Args) == 1 && t.Args[0].Int != nil:
		name := "float"
		if *t.Args[0].Int > 24 {
			name = "double"
		}
		return &core.TypeExpr{Name: name, Nullable: t.Nullable}
	case t.Name == "bool" || t.Name == "boolean":
		one := int64(1)
		return &core.TypeExpr{Name: "tinyint", Nullable: t.Nullable, Args: []core.TypeArg{{Int: &one}}}
	}
	return t
}

// stdlib is MySQL's functions in the form the catalog uses. They are embedded
// in the binary and never change within a run, so they are read once.
var stdlib = sync.OnceValue(func() []*catalog.Function {
	funcs, err := seed.Functions(dialectFS, "dialect")
	if err != nil {
		// A failure here means sqlc was built from a broken tree, which no
		// caller can do anything about.
		panic(err)
	}
	return funcs
})
