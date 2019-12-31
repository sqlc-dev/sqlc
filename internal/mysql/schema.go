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
func (s *Schema) getColType(col *sqlparser.ColName, defaultTableName string) (*sqlparser.ColumnDefinition, error) {
	if !col.Qualifier.IsEmpty() {
		colDfn, err := s.schemaLookup(col.Qualifier.Name.String(), col.Name.String())
		if err != nil {
			return nil, err
		}
		return &sqlparser.ColumnDefinition{
			Name: colDfn.Name, Type: colDfn.Type,
		}, nil
	}
	colDfn, err := s.schemaLookup(defaultTableName, col.Name.String())
	if err != nil {
		return nil, err
	}
	return &sqlparser.ColumnDefinition{
		Name: colDfn.Name, Type: colDfn.Type,
	}, nil
}

// Add add a MySQL table definition to the schema map
func (s *Schema) Add(table *sqlparser.DDL) {
	name := table.Table.Name.String()
	s.tables[name] = table.TableSpec.Columns
}

func (s *Schema) schemaLookup(table string, col string) (*sqlparser.ColumnDefinition, error) {
	cols, ok := s.tables[table]
	if !ok {
		return nil, fmt.Errorf("Table [%v] not found in Schema", table)
	}

	for _, colDef := range cols {
		if colDef.Name.EqualString(col) {
			return colDef, nil
		}
	}

	return nil, fmt.Errorf("Column [%v] not found in table [%v]", col, table)
}
