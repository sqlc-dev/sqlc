package compiler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/codegen"
	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/inflection"
	core "github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

type Result struct {
	Catalog *catalog.Catalog
	Queries []*Query
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

func paramName(p Parameter) string {
	if p.Column.Name != "" {
		return argName(p.Column.Name)
	}
	return fmt.Sprintf("dollar_%d", p.Number)
}

func (r *Result) Structs(settings config.CombinedSettings) []golang.Struct {
	var structs []golang.Struct
	for _, schema := range r.Catalog.Schemas {
		if schema.Name == "pg_catalog" {
			continue
		}
		for _, table := range schema.Tables {
			var tableName string
			if schema.Name == r.Catalog.DefaultSchema {
				tableName = table.Rel.Name
			} else {
				tableName = schema.Name + "_" + table.Rel.Name
			}
			structName := tableName
			if !settings.Go.EmitExactTableNames {
				structName = inflection.Singular(structName)
			}
			s := golang.Struct{
				Table:   core.FQN{Schema: schema.Name, Rel: table.Rel.Name},
				Name:    golang.StructName(structName, settings),
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, golang.Field{
					Name:    golang.StructName(column.Name, settings),
					Type:    r.goType(convertColumn(table.Rel, column), settings),
					Tags:    map[string]string{"json:": column.Name},
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

func (r *Result) Enums(settings config.CombinedSettings) []golang.Enum {
	var enums []golang.Enum
	for _, schema := range r.Catalog.Schemas {
		if schema.Name == "pg_catalog" {
			continue
		}
		for _, typ := range schema.Types {
			enum, ok := typ.(*catalog.Enum)
			if !ok {
				continue
			}
			var enumName string
			if schema.Name == r.Catalog.DefaultSchema {
				enumName = enum.Name
			} else {
				enumName = schema.Name + "_" + enum.Name
			}
			e := golang.Enum{
				Name:    golang.StructName(enumName, settings),
				Comment: enum.Comment,
			}
			for _, v := range enum.Vals {
				e.Constants = append(e.Constants, golang.Constant{
					Name:  e.Name + golang.EnumValueName(v),
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

type goColumn struct {
	id int
	*Column
}

func (r *Result) GoQueries(settings config.CombinedSettings) []golang.Query {
	structs := r.Structs(settings)

	qs := make([]golang.Query, 0, len(r.Queries))
	for _, query := range r.Queries {
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
			Comments:     query.Comments,
		}

		if len(query.Params) == 1 {
			p := query.Params[0]
			gq.Arg = golang.QueryValue{
				Name: paramName(p),
				Typ:  r.goType(p.Column, settings),
			}
		} else if len(query.Params) > 1 {
			var cols []goColumn
			for _, p := range query.Params {
				cols = append(cols, goColumn{
					id:     p.Number,
					Column: p.Column,
				})
			}
			gq.Arg = golang.QueryValue{
				Emit:   true,
				Name:   "arg",
				Struct: r.columnsToStruct(gq.MethodName+"Params", cols, settings),
			}
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = golang.QueryValue{
				Name: columnName(c, 0),
				Typ:  r.goType(c, settings),
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
					sameName := f.Name == golang.StructName(columnName(c, i), settings)
					sameType := f.Type == r.goType(c, settings)
					sameTable := sameTableName(c.Table, s.Table)

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
				var columns []goColumn
				for i, c := range query.Columns {
					columns = append(columns, goColumn{
						id:     i,
						Column: c,
					})
				}
				gs = r.columnsToStruct(gq.MethodName+"Row", columns, settings)
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

// It's possible that this method will generate duplicate JSON tag values
//
//   Columns: count, count,   count_2
//    Fields: Count, Count_2, Count2
// JSON tags: count, count_2, count_2
//
// This is unlikely to happen, so don't fix it yet
func (r Result) columnsToStruct(name string, columns []goColumn, settings config.CombinedSettings) *golang.Struct {
	gs := golang.Struct{
		Name: name,
	}
	seen := map[string]int{}
	suffixes := map[int]int{}
	for i, c := range columns {
		colName := columnName(c.Column, i)
		tagName := colName
		fieldName := golang.StructName(colName, settings)
		// Track suffixes by the ID of the column, so that columns referring to the same numbered parameter can be
		// reused.
		suffix := 0
		if o, ok := suffixes[c.id]; ok {
			suffix = o
		} else if v := seen[colName]; v > 0 {
			suffix = v + 1
		}
		suffixes[c.id] = suffix
		if suffix > 0 {
			tagName = fmt.Sprintf("%s_%d", tagName, suffix)
			fieldName = fmt.Sprintf("%s_%d", fieldName, suffix)
		}
		gs.Fields = append(gs.Fields, golang.Field{
			Name: fieldName,
			Type: r.goType(c.Column, settings),
			Tags: map[string]string{"json:": tagName},
		})
		seen[colName]++
	}
	return &gs
}

func columnName(c *Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}
