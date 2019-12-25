package pg

// String Functions and Operators
//
// https://www.postgresql.org/docs/current/functions-string.html
//
// Table 9.9. SQL String Functions and Operators
func stringFunctions() []Function {
	return []Function{
		argN("position", 2),
		{
			Name:       "lower",
			ReturnType: "text",
			Arguments: []Argument{
				{
					DataType: "string",
				},
			},
		},
		{
			Name:       "upper",
			ReturnType: "text",
			Arguments: []Argument{
				{
					DataType: "string",
				},
			},
		},
	}
}
