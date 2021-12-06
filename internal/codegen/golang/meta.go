package golang

import (
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
)

type TableStruct struct {
	FieldName string
	TableName string
	Columns   []ColumnStruct
}

type ColumnStruct struct {
	FieldName  string
	ColumnName string
}

func BuildMetaStructs(r *compiler.Result, settings config.CombinedSettings) []TableStruct {
	var tableStructs []TableStruct
	for _, schema := range r.Catalog.Schemas {
		if schema.Name == "pg_catalog" {
			continue
		}
		for _, table := range schema.Tables {
			var tableStruct TableStruct
			var tableName string
			if schema.Name == r.Catalog.DefaultSchema {
				tableName = table.Rel.Name
			} else {
				tableName = schema.Name + "." + table.Rel.Name
			}

			tableStruct.TableName = tableName
			tableStruct.FieldName = StructName(tableName, settings)
			for _, column := range table.Columns {
				tableStruct.Columns = append(tableStruct.Columns, ColumnStruct{
					FieldName:  StructName(column.Name, settings),
					ColumnName: column.Name,
				})
			}
			tableStructs = append(tableStructs, tableStruct)
		}
	}
	return tableStructs
}
