package mysql

import (
	"fmt"

	"vitess.io/vitess/go/vt/sqlparser"
)

// NewSchema gives a newly instantiated MySQL schema map
func NewSchema() *Schema {
	return &Schema{
		tables: make(map[string]([]*sqlparser.ColumnDefinition)),
	}
}

// Schema proves that information for mapping columns in queries to their respective table definitions
// and validating that they are correct so as to map to the correct Go type
type Schema struct {
	tables map[string]([]*sqlparser.ColumnDefinition)
}

// returns a deep copy of the column definition for using as a query return type or param type
func (s *Schema) getColType(col *sqlparser.ColName, tableAliasMap FromTables, defaultTableName string) (*Column, error) {
	realTable, err := tableColReferences(col, defaultTableName, tableAliasMap)
	if err != nil {
		return nil, err
	}

	colDfn, err := s.schemaLookup(realTable.TrueName, col.Name.String())
	if err != nil {
		return nil, err
	}
	colDfnCopy := *colDfn.ColumnDefinition
	if realTable.IsLeftJoined {
		colDfnCopy.Type.NotNull = false
	}
	return &Column{&colDfnCopy, realTable.TrueName}, nil
}

func tableColReferences(col *sqlparser.ColName, defaultTable string, tableAliasMap FromTables) (FromTable, error) {
	var table FromTable
	if col.Qualifier.IsEmpty() {
		if defaultTable == "" {
			return FromTable{}, fmt.Errorf("column reference \"%s\" is ambiguous, add a qualifier", col.Name.String())
		}
		table = FromTable{defaultTable, false}
	} else {
		fromTable, ok := tableAliasMap[col.Qualifier.Name.String()]
		if !ok {
			return FromTable{}, fmt.Errorf("column qualifier \"%s\" is not in schema or is an invalid alias", col.Qualifier.Name.String())
		}
		return fromTable, nil
	}
	return table, nil
}

// Add add a MySQL table definition to the schema map
func (s *Schema) Add(ddl *sqlparser.DDL) {
	switch ddl.Action {
	case "create":
		name := ddl.Table.Name.String()
		if ddl.TableSpec == nil {
			panic(fmt.Sprintf("failed to parse table \"%s\" schema.", name))
		}
		s.tables[name] = ddl.TableSpec.Columns
	}
}

func (s *Schema) schemaLookup(table string, col string) (*Column, error) {
	cols, ok := s.tables[table]
	if !ok {
		return nil, fmt.Errorf("table \"%s\" not found in schema", table)
	}

	for _, colDef := range cols {
		if colDef.Name.EqualString(col) {
			return &Column{colDef, table}, nil
		}
	}

	return nil, fmt.Errorf("column \"%s\" not found in table \"%s\"", col, table)
}
