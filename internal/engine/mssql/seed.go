package mssql

import (
	"embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds SQL Server's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}

func init() {
	core.RegisterCanonicalizer("mssql", canonicalize)
}

// canonicalize rewrites a type the way sys.types stores it: a length or
// precision left out is filled in with SQL Server's default, float(p) is a
// real or a float by its mantissa, sysname is nvarchar(128), and the MAX
// of nvarchar(max) is the word it is rather than a type.
func canonicalize(t *core.TypeExpr) *core.TypeExpr {
	for i, a := range t.Args {
		if a.Type != nil && a.Type.Name == "max" && len(a.Type.Args) == 0 {
			t = t.Clone()
			max := "max"
			t.Args[i] = core.TypeArg{Label: a.Label, Ident: &max}
		}
	}
	switch t.Name {
	case "sysname":
		n := int64(128)
		return &core.TypeExpr{Name: "nvarchar", Nullable: t.Nullable, Args: []core.TypeArg{{Int: &n}}}
	case "float":
		if len(t.Args) == 1 && t.Args[0].Int != nil {
			if *t.Args[0].Int <= 24 {
				return &core.TypeExpr{Name: "real", Nullable: t.Nullable}
			}
			return &core.TypeExpr{Name: "float", Nullable: t.Nullable}
		}
	case "char", "varchar", "nchar", "nvarchar", "binary", "varbinary":
		if len(t.Args) == 0 {
			n := int64(1)
			return &core.TypeExpr{Name: t.Name, Nullable: t.Nullable, Args: []core.TypeArg{{Int: &n}}}
		}
	case "decimal", "numeric":
		switch {
		case len(t.Args) == 0:
			p, s := int64(18), int64(0)
			return &core.TypeExpr{Name: t.Name, Nullable: t.Nullable, Args: []core.TypeArg{{Int: &p}, {Int: &s}}}
		case len(t.Args) == 1 && t.Args[0].Int != nil:
			s := int64(0)
			return &core.TypeExpr{Name: t.Name, Nullable: t.Nullable, Args: []core.TypeArg{t.Args[0], {Int: &s}}}
		}
	case "datetime2", "time", "datetimeoffset":
		if len(t.Args) == 0 {
			p := int64(7)
			return &core.TypeExpr{Name: t.Name, Nullable: t.Nullable, Args: []core.TypeArg{{Int: &p}}}
		}
	}
	return t
}
