package postgresql

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func pgTemp() *catalog.Schema {
	return &catalog.Schema{Name: "pg_temp"}
}

func typeName(name string) *ast.TypeName {
	return &ast.TypeName{Name: name}
}

func argN(name string, n int) *catalog.Function {
	var args []*catalog.Argument
	for i := 0; i < n; i++ {
		args = append(args, &catalog.Argument{
			Type: &ast.TypeName{Name: "any"},
		})
	}
	return &catalog.Function{
		Name:       name,
		Args:       args,
		ReturnType: &ast.TypeName{Name: "any"},
	}
}

func pgCatalog() *catalog.Schema {
	s := &catalog.Schema{Name: "pg_catalog"}
	s.Funcs = []*catalog.Function{
		// Table 9.5. Mathematical Functions
		// https://www.postgresql.org/docs/current/functions-math.html#FUNCTIONS-MATH-FUNC-TABLE
		argN("abs", 1),
		argN("cbrt", 1),
		argN("ceil", 1),
		argN("ceiling", 1),
		argN("degrees", 1),
		argN("div", 2),
		argN("exp", 1),
		argN("floor", 1),
		argN("ln", 1),
		argN("log", 1),
		argN("log", 2),
		argN("mod", 2),
		argN("pi", 0),
		argN("power", 2),
		argN("radians", 1),
		argN("round", 1),
		argN("round", 2),
		argN("scale", 1),
		argN("sign", 1),
		argN("sqrt", 1),
		argN("trunc", 1),
		argN("trunc", 2),
		argN("width_bucket", 2),
		argN("width_bucket", 4),

		// Table 9.6. Random Functions
		// https://www.postgresql.org/docs/current/functions-math.html#FUNCTIONS-MATH-RANDOM-TABLE
		argN("random", 0),

		// Table 9.52. General-Purpose Aggregate Functions
		// https://www.postgresql.org/docs/current/functions-aggregate.html#FUNCTIONS-AGGREGATE-TABLE
		{
			Name: "bool_and",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: typeName("bool"),
		},
		{
			Name: "bool_or",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: typeName("bool"),
		},
		{
			Name:       "count",
			ReturnType: typeName("bigint"),
		},
		{
			Name: "count",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: typeName("bigint"),
		},
		{
			Name: "every",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: typeName("bool"),
		},

		// Table 9.9. SQL String Functions and Operators
		// https://www.postgresql.org/docs/current/functions-string.html
		argN("position", 2),
		{
			Name:       "lower",
			ReturnType: typeName("text"),
			Args: []*catalog.Argument{
				{Type: typeName("string")},
			},
		},
		{
			Name:       "upper",
			ReturnType: typeName("text"),
			Args: []*catalog.Argument{
				{Type: typeName("string")},
			},
		},

		// Advisory Lock Functions
		//
		// The functions shown in Table 9.95 manage advisory locks. For details about
		// proper use of these functions, see Section 13.3.5.
		//
		// https://www.postgresql.org/docs/current/functions-admin.html
		//
		// Table 9.95. Advisory Lock Functions
		{
			Name:       "pg_advisory_lock",
			Desc:       "Obtain exclusive session level advisory lock",
			ReturnType: typeName("void"),
			Args: []*catalog.Argument{
				{
					Name: "key",
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
		},
		{
			Name:       "pg_advisory_lock",
			Desc:       "Obtain exclusive session level advisory lock",
			ReturnType: typeName("void"),
			Args: []*catalog.Argument{
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
			},
		},
		{
			Name:       "pg_advisory_lock_shared",
			Desc:       "Obtain shared session level advisory lock",
			ReturnType: typeName("void"),
			Args: []*catalog.Argument{
				{
					Name: "key",
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
		},
		{
			Name:       "pg_advisory_lock_shared",
			Desc:       "Obtain shared session level advisory lock",
			ReturnType: typeName("void"),
			Args: []*catalog.Argument{
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
			},
		},
		{
			Name:       "pg_advisory_unlock",
			Desc:       "Release an exclusive session level advisory lock",
			ReturnType: typeName("bool"),
			Args: []*catalog.Argument{
				{
					Name: "key",
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
		},
		{
			Name:       "pg_advisory_unlock",
			Desc:       "Release an exclusive session level advisory lock",
			ReturnType: typeName("bool"),
			Args: []*catalog.Argument{
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
			},
		},
		{
			Name:       "pg_advisory_unlock_all",
			Desc:       "Release all session level advisory locks held by the current session",
			ReturnType: typeName("void"),
		},
		{
			Name:       "pg_advisory_unlock_shared",
			Desc:       "Unlock a shared session level advisory lock",
			ReturnType: typeName("bool"),
			Args: []*catalog.Argument{
				{
					Name: "key",
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
		},
		{
			Name:       "pg_advisory_unlock_shared",
			Desc:       "Unlock a shared session level advisory lock",
			ReturnType: typeName("bool"),
			Args: []*catalog.Argument{
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
			},
		},
		{
			Name:       "pg_advisory_xact_lock",
			Desc:       "Obtain exclusive transaction level advisory lock",
			ReturnType: typeName("void"),
			Args: []*catalog.Argument{
				{
					Name: "key",
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
		},
		{
			Name:       "pg_advisory_xact_lock",
			Desc:       "Obtain exclusive transaction level advisory lock",
			ReturnType: typeName("void"),
			Args: []*catalog.Argument{
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
			},
		},
		{
			Name:       "pg_advisory_xact_lock_shared",
			Desc:       "Obtain exclusive transaction level advisory lock",
			ReturnType: typeName("void"),
			Args: []*catalog.Argument{
				{
					Name: "key",
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
		},
		{
			Name:       "pg_advisory_xact_lock_shared",
			Desc:       "Obtain exclusive transaction level advisory lock",
			ReturnType: typeName("void"),
			Args: []*catalog.Argument{
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Name: "key1",
					Type: &ast.TypeName{Name: "int"},
				},
			},
		},
	}
	return s
}
