package pg

// Advisory Lock Functions
//
// The functions shown in Table 9.95 manage advisory locks. For details about
// proper use of these functions, see Section 13.3.5.
//
// https://www.postgresql.org/docs/current/functions-admin.html
//
// Table 9.95. Advisory Lock Functions
func advisoryLockFunctions() []Function {
	return []Function{
		{
			Name:       "pg_advisory_lock",
			Desc:       "Obtain exclusive session level advisory lock",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "key",
					DataType: "bigint",
				},
			},
		},
		{
			Name:       "pg_advisory_lock",
			Desc:       "Obtain exclusive session level advisory lock",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "key1",
					DataType: "int",
				},
				{
					Name:     "key1",
					DataType: "int",
				},
			},
		},
		{
			Name:       "pg_advisory_lock_shared",
			Desc:       "Obtain shared session level advisory lock",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "key",
					DataType: "bigint",
				},
			},
		},
		{
			Name:       "pg_advisory_lock_shared",
			Desc:       "Obtain shared session level advisory lock",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "key1",
					DataType: "int",
				},
				{
					Name:     "key1",
					DataType: "int",
				},
			},
		},
		{
			Name:       "pg_advisory_unlock",
			Desc:       "Release an exclusive session level advisory lock",
			ReturnType: "bool",
			Arguments: []Argument{
				{
					Name:     "key",
					DataType: "bigint",
				},
			},
		},
		{
			Name:       "pg_advisory_unlock",
			Desc:       "Release an exclusive session level advisory lock",
			ReturnType: "bool",
			Arguments: []Argument{
				{
					Name:     "key1",
					DataType: "int",
				},
				{
					Name:     "key1",
					DataType: "int",
				},
			},
		},
		{
			Name:       "pg_advisory_unlock_all",
			Desc:       "Release all session level advisory locks held by the current session",
			ReturnType: "void",
		},
		{
			Name:       "pg_advisory_unlock_shared",
			Desc:       "Unlock a shared session level advisory lock",
			ReturnType: "bool",
			Arguments: []Argument{
				{
					Name:     "key",
					DataType: "bigint",
				},
			},
		},
		{
			Name:       "pg_advisory_unlock_shared",
			Desc:       "Unlock a shared session level advisory lock",
			ReturnType: "bool",
			Arguments: []Argument{
				{
					Name:     "key1",
					DataType: "int",
				},
				{
					Name:     "key1",
					DataType: "int",
				},
			},
		},
		{
			Name:       "pg_advisory_xact_lock",
			Desc:       "Obtain exclusive transaction level advisory lock",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "key",
					DataType: "bigint",
				},
			},
		},
		{
			Name:       "pg_advisory_xact_lock",
			Desc:       "Obtain exclusive transaction level advisory lock",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "key1",
					DataType: "int",
				},
				{
					Name:     "key1",
					DataType: "int",
				},
			},
		},
		{
			Name:       "pg_advisory_xact_lock",
			Desc:       "Obtain exclusive transaction level advisory lock",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "key",
					DataType: "bigint",
				},
			},
		},
		{
			Name:       "pg_advisory_xact_lock",
			Desc:       "Obtain exclusive transaction level advisory lock",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "key1",
					DataType: "int",
				},
				{
					Name:     "key1",
					DataType: "int",
				},
			},
		},
	}
}
