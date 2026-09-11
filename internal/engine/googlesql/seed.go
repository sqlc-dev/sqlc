package googlesql

import (
	"embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds GoogleSQL's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}

func init() {
	core.RegisterCanonicalizer("googlesql", canonicalize)
}

// canonicalize reads the MAX of STRING(MAX) as the word it is rather than a
// type, which is the only thing a spelling cannot say for itself.
func canonicalize(t *core.TypeExpr) *core.TypeExpr {
	return maxIdent(t)
}

// maxIdent rewrites an argument that is the bare word max into an
// identifier argument.
func maxIdent(t *core.TypeExpr) *core.TypeExpr {
	for i, a := range t.Args {
		if a.Type != nil && a.Type.Name == "max" && len(a.Type.Args) == 0 {
			out := t.Clone()
			max := "max"
			out.Args[i] = core.TypeArg{Label: a.Label, Ident: &max}
			return out
		}
	}
	return t
}
