package python

import (
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/core"
	"github.com/kyleconroy/sqlc/internal/inflection"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func modelName(name string, settings *plugin.Settings) string {
	if rename := settings.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	for _, p := range strings.Split(name, "_") {
		out += strings.Title(p)
	}
	return out
}

func makeEnums(req *plugin.CodeGenRequest) []Enum {
	var enums []Enum
	for _, schema := range req.Catalog.Schemas {
		if schema.Name == "pg_catalog" {
			continue
		}
		for _, enum := range schema.Enums {
			var enumName string
			if schema.Name == req.Catalog.DefaultSchema {
				enumName = enum.Name
			} else {
				enumName = schema.Name + "_" + enum.Name
			}
			e := Enum{
				Name:    modelName(enumName, req.Settings),
				Comment: enum.Comment,
			}
			for _, v := range enum.Vals {
				e.Constants = append(e.Constants, Constant{
					Name:  pyEnumValueName(v),
					Value: v,
					Type:  e.Name,
				})
			}
			enums = append(enums, e)
		}
	}
	if len(enums) > 0 {
		sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })
	}
	return enums
}

func makePyType2(req *plugin.CodeGenRequest, col *plugin.Column, settings *plugin.Settings) pyType {
	typ := pyInnerType2(req, col, settings)
	return pyType{
		InnerType: typ,
		IsArray:   col.IsArray,
		IsNull:    !col.NotNull,
	}
}

func pyInnerType2(req *plugin.CodeGenRequest, col *plugin.Column, settings *plugin.Settings) string {
	for _, oride := range settings.Overrides {
		// FIXME PythonType isn't sent to the plugin
		if oride.CodeType == "" {
			continue
		}
		// sameTable := oride.Matches(col.Table, req.Catalog.DefaultSchema)
		// if oride.Column != "" && oride.ColumnName.MatchString(col.Name) && sameTable {
		// 	return oride.PythonType.TypeString()
		// }
		// if oride.DBType != "" && oride.DBType == col.DataType && oride.Nullable != (col.NotNull || col.IsArray) {
		// 	return oride.PythonType.TypeString()
		// }
	}

	if settings.Engine == string(config.EnginePostgreSQL) {
		return postgresType(r, col, settings)
	} else {
		log.Println("unsupported engine type")
		return "Any"
	}
}

func makeModels(req *plugin.CodeGenRequest) []Struct {
	var structs []Struct
	for _, schema := range req.Catalog.Schemas {
		if schema.Name == "pg_catalog" {
			continue
		}
		for _, table := range schema.Tables {
			var tableName string
			if schema.Name == req.Catalog.DefaultSchema {
				tableName = table.Rel.Name
			} else {
				tableName = schema.Name + "_" + table.Rel.Name
			}
			structName := tableName
			// FIXME How do we deal with plugin specific settings?
			if false { // !req.Settings.Python.EmitExactTableNames {
				structName = inflection.Singular(structName)
			}
			s := Struct{
				Table:   core.FQN{Schema: schema.Name, Rel: table.Rel.Name},
				Name:    modelName(structName, req.Settings),
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				typ := makePyType2(req, column, req.Settings)
				typ.InnerType = strings.TrimPrefix(typ.InnerType, "models.")
				s.Fields = append(s.Fields, Field{
					Name:    column.Name,
					Type:    typ,
					Comment: column.Comment,
				})
			}
			structs = append(structs, s)
		}
	}
	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}
	return structs
}

func GenerateV2(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	makeEnums(req)
	makeModels(req)

	return &plugin.CodeGenResponse{}, nil
}
