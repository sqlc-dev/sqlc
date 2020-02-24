package info

import (
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

// Provide a read-only view into the catalog
func Newo(c *catalog.Catalog) InformationSchema {
	return InformationSchema{c: c}
}

type InformationSchema struct {
	c *catalog.Catalog
}

type Table struct {
	Catalog string
	Schema  string
	Name    string
}

// SELECT * FROM information_schema.tables;
func (i *InformationSchema) Tables() []Table {
	var tables []Table
	for _, s := range i.c.Schemas {
		for _, t := range s.Tables {
			tables = append(tables, Table{
				Catalog: i.c.Name,
				Schema:  s.Name,
				Name:    t.Rel.Name,
			})
		}
	}
	return tables
}

type Column struct {
	Catalog    string
	Schema     string
	Table      string
	Name       string
	DataType   string
	IsNullable bool
}

// SELECT * FROM information_schema.columns;
func (i *InformationSchema) Columns() []Column {
	return []Column{}
}

type Schema struct {
	Catalog string
	Name    string
}

// SELECT * FROM information_schema.schemata;
func (i *InformationSchema) Schemata() []Schema {
	return []Schema{}
}
