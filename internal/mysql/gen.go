package mysql

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jinzhu/inflection"
	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/kyleconroy/sqlc/internal/codegen"
	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/core"
)

type PackageGenerator struct {
	*Schema
	config.CombinedSettings
	packageName string
}

type Result struct {
	PackageGenerator
	Queries []*Query
}

// Enums generates parser-agnostic GoEnum types
func (r *Result) Enums(settings config.CombinedSettings) []golang.Enum {
	var enums []golang.Enum
	for _, table := range r.Schema.tables {
		for _, col := range table {
			if strings.ToLower(col.Type.Type) == "enum" {
				constants := []golang.Constant{}
				enumName := r.enumNameFromColDef(col)
				for _, c := range col.Type.EnumValues {
					stripped := stripInnerQuotes(c)
					constants = append(constants, golang.Constant{
						// TODO: maybe add the struct name call to capitalize the name here
						Name:  stripped,
						Value: stripped,
						Type:  enumName,
					})
				}

				goEnum := golang.Enum{
					Name:      enumName,
					Comment:   "",
					Constants: constants,
				}
				enums = append(enums, goEnum)
			}
		}
	}
	return enums
}

func stripInnerQuotes(identifier string) string {
	return strings.Replace(identifier, "'", "", 2)
}

func (pGen PackageGenerator) enumNameFromColDef(col *sqlparser.ColumnDefinition) string {
	return fmt.Sprintf("%sType",
		golang.StructName(col.Name.String(), pGen.CombinedSettings))
}

// Structs marshels each query into a go struct for generation
func (r *Result) Structs(settings config.CombinedSettings) []golang.Struct {
	var structs []golang.Struct
	for tableName, cols := range r.Schema.tables {
		structName := golang.StructName(tableName, settings)
		if !(settings.Go.EmitExactTableNames || settings.Kotlin.EmitExactTableNames) {
			structName = inflection.Singular(structName)
		}
		s := golang.Struct{
			Name:  structName,
			Table: core.FQN{tableName, "", ""}, // TODO: Complete hack. Only need for equality check to see if struct can be reused between queries
		}

		for _, col := range cols {
			tags := map[string]string{}
			if settings.Go.EmitDBTags {
				tags["db:"] = col.Name.String()
			}
			if settings.Go.EmitJSONTags {
				tags["json:"] = col.Name.String()
			}
			s.Fields = append(s.Fields, golang.Field{
				Name:    golang.StructName(col.Name.String(), settings),
				Type:    r.goTypeCol(Column{col, tableName}),
				Tags:    tags,
				Comment: "",
			})
		}
		structs = append(structs, s)
	}
	sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	return structs
}

// GoQueries generates parser-agnostic query information for code generation
func (r *Result) GoQueries(settings config.CombinedSettings) []golang.Query {
	structs := r.Structs(settings)

	qs := make([]golang.Query, 0, len(r.Queries))
	for ix, query := range r.Queries {
		if query == nil {
			panic(fmt.Sprintf("query is nil on index: %v, len: %v", ix, len(r.Queries)))
		}
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}

		gq := golang.Query{
			Cmd:          query.Cmd,
			ConstantName: codegen.LowerTitle(query.Name),
			FieldName:    codegen.LowerTitle(query.Name) + "Stmt",
			MethodName:   query.Name,
			SourceName:   query.Filename,
			SQL:          query.SQL,
			// Comments:     query.Comments,
		}

		if len(query.Params) == 1 {
			p := query.Params[0]
			gq.Arg = golang.QueryValue{
				Name: p.Name,
				Typ:  p.Typ,
			}
		} else if len(query.Params) > 1 {

			structInfo := make([]structParams, len(query.Params))
			for i := range query.Params {
				qp := query.Params[i]
				if qp.Typ == "" {
					// if the param doesn't have a type, check to see if there is
					// another param with the same name that does have a type.
					// Because of the way params are parsed and named this only works for sqlc.arg(x) named params, not :x or ?
					func(ps []*Param) {
						for j := range ps {
							if ps[j].OriginalName == qp.OriginalName &&
								ps[j].Typ != "" {
								query.Params[i].Typ = ps[j].Typ
							}
						}
					}(query.Params)
				}
				structInfo[i] = structParams{
					originalName: query.Params[i].Name,
					goType:       query.Params[i].Typ,
				}
			}

			gq.Arg = golang.QueryValue{
				Emit:   true,
				Name:   "arg",
				Struct: r.columnsToStruct(gq.MethodName+"Params", structInfo, settings),
			}
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = golang.QueryValue{
				Name: columnName(c.ColumnDefinition, 0),
				Typ:  r.goTypeCol(c),
			}
		} else if len(query.Columns) > 1 {
			var gs *golang.Struct
			var emit bool

			for _, s := range structs {
				if len(s.Fields) != len(query.Columns) {
					continue
				}
				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					sameName := f.Name == golang.StructName(columnName(c.ColumnDefinition, i), settings)
					sameType := f.Type == r.goTypeCol(c)

					hackedFQN := core.FQN{c.Table, "", ""} // TODO: only check needed here is equality to see if struct can be reused, this type should be removed or properly used
					sameTable := s.Table.Catalog == hackedFQN.Catalog && s.Table.Schema == hackedFQN.Schema && s.Table.Rel == hackedFQN.Rel

					if !sameName || !sameType || !sameTable {
						same = false
					}
				}
				if same {
					gs = &s
					break
				}
			}

			if gs == nil {
				structInfo := make([]structParams, len(query.Columns))
				for i := range query.Columns {
					structInfo[i] = structParams{
						originalName: query.Columns[i].Name.String(),
						goType:       r.goTypeCol(query.Columns[i]),
					}
				}
				gs = r.columnsToStruct(gq.MethodName+"Row", structInfo, settings)
				emit = true
			}
			gq.Ret = golang.QueryValue{
				Emit:   emit,
				Name:   "i",
				Struct: gs,
			}
		}

		qs = append(qs, gq)
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].MethodName < qs[j].MethodName })
	return qs
}

type structParams struct {
	originalName string
	goType       string
}

func (r *Result) columnsToStruct(name string, items []structParams, settings config.CombinedSettings) *golang.Struct {
	gs := golang.Struct{
		Name: name,
	}
	seen := map[string]int{}
	for _, item := range items {
		name := item.originalName
		typ := item.goType
		tagName := name
		fieldName := golang.StructName(name, settings)
		if v := seen[name]; v > 0 {
			tagName = fmt.Sprintf("%s_%d", tagName, v+1)
			fieldName = fmt.Sprintf("%s_%d", fieldName, v+1)
		}
		tags := map[string]string{}
		if settings.Go.EmitDBTags {
			tags["db:"] = tagName
		}
		if settings.Go.EmitJSONTags {
			tags["json:"] = tagName
		}
		gs.Fields = append(gs.Fields, golang.Field{
			Name: fieldName,
			Type: typ,
			Tags: tags,
		})
		seen[name]++
	}
	return &gs
}

func (pGen PackageGenerator) goTypeCol(col Column) string {
	mySQLType := strings.ToLower(col.ColumnDefinition.Type.Type)
	notNull := bool(col.Type.NotNull)
	colName := col.Name.String()

	for _, oride := range pGen.Overrides {
		shouldOverride := (oride.DBType != "" && oride.DBType == mySQLType && oride.Nullable != notNull) ||
			(oride.ColumnName != "" && oride.ColumnName == colName && oride.Table.Rel == col.Table)
		if shouldOverride {
			return oride.GoTypeName
		}
	}
	switch t := mySQLType; {
	case "varchar" == t, "text" == t, "char" == t,
		"tinytext" == t, "mediumtext" == t, "longtext" == t:
		if col.Type.NotNull {
			return "string"
		}
		return "sql.NullString"
	case "int" == t, "integer" == t, t == "smallint",
		"mediumint" == t, "bigint" == t, "year" == t:
		if col.Type.NotNull {
			return "int"
		}
		return "sql.NullInt64"
	case "blob" == t, "binary" == t, "varbinary" == t, "tinyblob" == t,
		"mediumblob" == t, "longblob" == t:
		return "[]byte"
	case "float" == t, strings.HasPrefix(strings.ToLower(t), "decimal"):
		if col.Type.NotNull {
			return "float64"
		}
		return "sql.NullFloat64"
	case "enum" == t:
		return pGen.enumNameFromColDef(col.ColumnDefinition)
	case "date" == t, "timestamp" == t, "datetime" == t, "time" == t:
		if col.Type.NotNull {
			return "time.Time"
		}
		return "sql.NullTime"
	case "boolean" == t, "bool" == t, "tinyint" == t:
		if col.Type.NotNull {
			return "bool"
		}
		return "sql.NullBool"
	default:
		fmt.Printf("unknown MySQL type: %s\n", t)
		return "interface{}"
	}
}

func columnName(c *sqlparser.ColumnDefinition, pos int) string {
	if !c.Name.IsEmpty() {
		return c.Name.String()
	}
	return fmt.Sprintf("column_%d", pos+1)
}

func argName(name string) string {
	out := ""
	for i, p := range strings.Split(name, "_") {
		if i == 0 {
			out += strings.ToLower(p)
		} else if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}
