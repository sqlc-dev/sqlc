package pg

func NewCatalog() Catalog {
	return Catalog{
		Schemas: map[string]Schema{
			"public":     NewSchema(),
			"pg_catalog": pgCatalog(),
			// Likewise, the current session's temporary-table schema, pg_temp_nnn, is
			// always searched if it exists. It can be explicitly listed in the path by
			// using the alias pg_temp. If it is not listed in the path then it is
			// searched first (even before pg_catalog). However, the temporary schema is
			// only searched for relation (table, view, sequence, etc) and data type
			// names. It is never searched for function or operator names.
			//
			// https://www.postgresql.org/docs/current/runtime-config-client.html
			"pg_temp": NewSchema(),
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

func (f FQN) String() string {
	s := f.Rel
	if f.Schema != "" {
		s = f.Schema + "." + s
	}
	if f.Catalog != "" {
		s = f.Catalog + "." + s
	}
	return s
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
		arity := fun.ArgN
		if fun.Arguments != nil {
			arity = len(fun.Arguments)
		}
		if arity == argn {
			return fun, nil
		}
	}
	return Function{}, ErrorRelationDoesNotExist(fqn.Rel)
}

type Schema struct {
	Name    string
	Tables  map[string]Table
	Enums   map[string]Enum
	Funcs   map[string][]Function
	Comment string
}

type Table struct {
	ID      FQN
	Name    string
	Columns []Column
	Comment string
}

type Column struct {
	Name     string
	DataType string
	NotNull  bool
	IsArray  bool
	Comment  string

	// XXX: Figure out what PostgreSQL calls `foo.id`
	Scope string
	Table FQN
}

type Enum struct {
	Name    string
	Vals    []string
	Comment string
}

type Function struct {
	Name       string
	ArgN       int
	Arguments  []Argument // not recorded for builtins
	ReturnType string
	Comment    string
	Desc       string
}

type Argument struct {
	Name       string
	DataType   string
	HasDefault bool
}
