package pg

func argN(name string, n int) Function {
	return Function{
		Name:       name,
		ArgN:       n,
		ReturnType: "any",
	}
}

func pgCatalog() Schema {
	s := NewSchema()
	s.Name = "pg_catalog"
	fs := []Function{

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

		// Table 9.8. SQL String Functions and Operators
		// https://www.postgresql.org/docs/current/functions-string.html#FUNCTIONS-STRING-SQL
		argN("position", 2),

		// Table 9.52. General-Purpose Aggregate Functions
		// https://www.postgresql.org/docs/current/functions-aggregate.html#FUNCTIONS-AGGREGATE-TABLE
		{
			Name:       "bool_and",
			ArgN:       1,
			ReturnType: "bool",
		},
		{
			Name:       "bool_or",
			ArgN:       1,
			ReturnType: "bool",
		},
		{
			Name:       "count",
			ArgN:       0,
			ReturnType: "bigint",
		},
		{
			Name:       "count",
			ArgN:       1,
			ReturnType: "bigint",
		},
		{
			Name:       "every",
			ArgN:       1,
			ReturnType: "bool",
		},
	}

	fs = append(fs, advisoryLockFunctions()...)

	s.Funcs = make(map[string][]Function, len(fs))
	for _, f := range fs {
		s.Funcs[f.Name] = append(s.Funcs[f.Name], f)
	}
	return s
}
