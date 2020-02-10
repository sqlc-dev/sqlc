package pg

func internalSchema() Schema {
	s := NewSchema()
	s.Name = "sqlc"
	fs := []Function{
		{
			Name:       "arg",
			Desc:       "Named argumented placeholder",
			ReturnType: "void",
			Arguments: []Argument{
				{
					Name:     "name",
					DataType: "id",
				},
			},
		},
	}
	s.Funcs = make(map[string][]Function, len(fs))
	for _, f := range fs {
		s.Funcs[f.Name] = append(s.Funcs[f.Name], f)
	}
	return s
}
