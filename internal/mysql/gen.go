package mysql

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"vitess.io/vitess/go/vt/sqlparser"
)

// Result holds the mysql validated queries schema
type Result struct {
	Queries     []*Query
	Schema      *Schema
	packageName string
}

func (r *Result) PkgName() string {
	return r.packageName
}

// Enums generates parser-agnostic GoEnum types
func (r *Result) Enums(settings dinosql.GenerateSettings) []dinosql.GoEnum {
	var enums []dinosql.GoEnum
	for _, table := range r.Schema.tables {
		for _, col := range table {
			if col.Type.Type == "enum" {
				constants := []dinosql.GoConstant{}
				enumName := enumNameFromColDef(col, settings)
				for _, c := range col.Type.EnumValues {
					stripped := stripInnerQuotes(c)
					constants = append(constants, dinosql.GoConstant{
						// TODO: maybe add the struct name call to capitalize the name here
						Name:  stripped,
						Value: stripped,
						Type:  enumName,
					})
				}

				goEnum := dinosql.GoEnum{
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

func enumNameFromColDef(col *sqlparser.ColumnDefinition, settings dinosql.GenerateSettings) string {
	return fmt.Sprintf("%sType",
		dinosql.StructName(col.Name.String(), settings))
}

// Structs marshels each query into a go struct for generation
func (r *Result) Structs(settings dinosql.GenerateSettings) []dinosql.GoStruct {
	var structs []dinosql.GoStruct
	for tableName, cols := range r.Schema.tables {
		s := dinosql.GoStruct{
			Name: inflection.Singular(dinosql.StructName(tableName, settings)),
		}

		for _, col := range cols {
			s.Fields = append(s.Fields, dinosql.GoField{
				Name:    dinosql.StructName(col.Name.String(), settings),
				Type:    goTypeCol(col, settings),
				Tags:    map[string]string{"json": col.Name.String()},
				Comment: "",
			})
		}
		structs = append(structs, s)
	}

	return structs
}

// GoQueries generates parser-agnostic query information for code generation
func (r *Result) GoQueries(settings dinosql.GenerateSettings) []dinosql.GoQuery {
	structs := r.Structs(settings)

	qs := make([]dinosql.GoQuery, 0, len(r.Queries))
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

		gq := dinosql.GoQuery{
			Cmd:          query.Cmd,
			ConstantName: dinosql.LowerTitle(query.Name),
			FieldName:    dinosql.LowerTitle(query.Name) + "Stmt",
			MethodName:   query.Name,
			SourceName:   "queries", // query.Filename,
			SQL:          query.SQL,
			// Comments:     query.Comments,
		}

		if len(query.Params) == 1 {
			p := query.Params[0]
			gq.Arg = dinosql.GoQueryValue{
				Name: p.name,
				Typ:  p.typ,
			}
		} else if len(query.Params) > 1 {

			structInfo := make([]structParams, len(query.Params))
			for i := range query.Params {
				structInfo[i] = structParams{
					originalName: query.Params[i].name,
					goType:       query.Params[i].typ,
				}
			}

			gq.Arg = dinosql.GoQueryValue{
				Emit:   true,
				Name:   "arg",
				Struct: r.columnsToStruct(gq.MethodName+"Params", structInfo, settings),
			}
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = dinosql.GoQueryValue{
				Name: columnName(c, 0),
				Typ:  goTypeCol(c, settings),
			}
		} else if len(query.Columns) > 1 {
			var gs *dinosql.GoStruct
			var emit bool

			for _, s := range structs {
				if len(s.Fields) != len(query.Columns) {
					continue
				}
				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					sameName := f.Name == dinosql.StructName(columnName(c, i), settings)
					sameType := f.Type == goTypeCol(c, settings)
					// TODO: consider making this deep equality from stdlib?
					// sameFQN := s.Table.EqualTo(&c.Table)
					if !sameName || !sameType || true { // !sameFQN
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
						goType:       goTypeCol(query.Columns[i], settings),
					}
				}
				gs = r.columnsToStruct(gq.MethodName+"Row", structInfo, settings)
				emit = true
			}
			gq.Ret = dinosql.GoQueryValue{
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

func (r *Result) columnsToStruct(name string, items []structParams, settings dinosql.GenerateSettings) *dinosql.GoStruct {
	gs := dinosql.GoStruct{
		Name: name,
	}
	seen := map[string]int{}
	for _, item := range items {
		name := item.originalName
		typ := item.goType
		tagName := name
		fieldName := dinosql.StructName(name, settings)
		if v := seen[name]; v > 0 {
			tagName = fmt.Sprintf("%s_%d", tagName, v+1)
			fieldName = fmt.Sprintf("%s_%d", fieldName, v+1)
		}
		gs.Fields = append(gs.Fields, dinosql.GoField{
			Name: fieldName,
			Type: typ,
			Tags: map[string]string{"json:": tagName},
		})
		seen[name]++
	}
	return &gs
}

func goTypeCol(col *sqlparser.ColumnDefinition, settings dinosql.GenerateSettings) string {
	switch t := col.Type.Type; {
	case "varchar" == t:
		if col.Type.NotNull {
			return "string"
		}
		return "sql.NullString"
	case "int" == t:
		if col.Type.NotNull {
			return "int"
		}
		return "sql.NullInt64"
	case "float" == t, strings.HasPrefix(strings.ToLower(t), "decimal"):
		if col.Type.NotNull {
			return "float64"
		}
		return "sql.NullFloat64"
	case "enum" == t:
		return enumNameFromColDef(col, settings)
	default:
		// TODO: remove panic here
		panic(fmt.Sprintf("Handle this col type directly: %v\n", col.Type))
		// return col.Type
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
