package pg

func NewCatalog() Catalog {
	return Catalog{
		Schemas: map[string]Schema{
			"public":     NewSchema(),
			"pg_catalog": pgCatalog(),
		},
	}
}

func NewSchema() Schema {
	return Schema{
		Tables: map[string]Table{},
		Enums:  map[string]Enum{},
		Funcs:  map[string][]Function{},
	}
}

type FQN struct {
	Catalog string
	Schema  string
	Rel     string
}

type Catalog struct {
	Schemas map[string]Schema
}

func (c Catalog) LookupFunctions(fqn FQN) ([]Function, error) {
	schema, exists := c.Schemas[fqn.Schema]
	if !exists {
		return nil, ErrorSchemaDoesNotExist(fqn.Schema)
	}

	// pg_catalog is always effectively part of the search path. If it is not
	// named explicitly in the path then it is implicitly searched before
	// searching the path's schemas.
	//
	// https://www.postgresql.org/docs/current/ddl-schemas.html#DDL-SCHEMAS-PATH
	schemas := []Schema{c.Schemas["pg_catalog"], schema}

	var funs []Function
	for _, s := range schemas {
		// TODO: Efficient function search
		funs = append(funs, s.Funcs[fqn.Rel]...)
	}

	if len(funs) == 0 {
		return nil, ErrorRelationDoesNotExist(fqn.Rel)
	}

	return funs, nil
}

func (c Catalog) LookupFunctionN(fqn FQN, argn int) (Function, error) {
	funs, err := c.LookupFunctions(fqn)
	if err != nil {
		return Function{}, err
	}
	for _, fun := range funs {
		if fun.ArgN == argn {
			return fun, nil
		}
	}
	return Function{}, ErrorRelationDoesNotExist(fqn.Rel)
}

type Schema struct {
	Name   string
	Tables map[string]Table
	Enums  map[string]Enum
	Funcs  map[string][]Function
}

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name     string
	DataType string
	NotNull  bool
	IsArray  bool
}

type Enum struct {
	Name string
	Vals []string
}

type Function struct {
	Name       string
	ArgN       int
	ReturnType string
}
