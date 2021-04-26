package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/codegen"
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/core"
	"github.com/kyleconroy/sqlc/internal/inflection"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func buildEnums(r *compiler.Result, settings config.CombinedSettings) []Enum {
	var enums []Enum
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
			e := Enum{
				Name:    StructName(enumName, settings),
				Comment: enum.Comment,
			}
			seen := make(map[string]struct{}, len(enum.Vals))
			for i, v := range enum.Vals {
				value := EnumReplace(v)
				if _, found := seen[value]; found || value == "" {
					value = fmt.Sprintf("value_%d", i)
				}
				e.Constants = append(e.Constants, Constant{
					Name:  StructName(enumName+"_"+value, settings),
					Value: v,
					Type:  e.Name,
				})
				seen[value] = struct{}{}
			}
			enums = append(enums, e)
		}
	}
	if len(enums) > 0 {
		sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })
	}
	return enums
}

func buildStructs(r *compiler.Result, settings config.CombinedSettings) []Struct {
	var structs []Struct
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
			s := Struct{
				Table:   core.FQN{Schema: schema.Name, Rel: table.Rel.Name},
				Name:    StructName(structName, settings),
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				tags := map[string]string{}
				if settings.Go.EmitDBTags {
					tags["db:"] = column.Name
				}
				if settings.Go.EmitJSONTags {
					tags["json:"] = JSONTagName(column.Name, settings)
				}
				s.Fields = append(s.Fields, Field{
					Name:    StructName(column.Name, settings),
					Type:    goType(r, compiler.ConvertColumn(table.Rel, column), settings),
					Tags:    tags,
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

type goColumn struct {
	id    int
	Embed *Struct
	*compiler.Column
}

func columnName(c *compiler.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}

func paramName(p compiler.Parameter) string {
	if p.Column.Name != "" {
		return argName(p.Column.Name)
	}
	return fmt.Sprintf("dollar_%d", p.Number)
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

func buildQueries(r *compiler.Result, settings config.CombinedSettings, structs []Struct) []Query {
	qs := make([]Query, 0, len(r.Queries))
	for _, query := range r.Queries {
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}

		gq := Query{
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
			gq.Arg = QueryValue{
				Name: paramName(p),
				Typ:  goType(r, p.Column, settings),
			}
		} else if len(query.Params) > 1 {
			var cols []goColumn
			for _, p := range query.Params {
				cols = append(cols, goColumn{
					id:     p.Number,
					Column: p.Column,
				})
			}
			gq.Arg = QueryValue{
				Emit:   true,
				Name:   "arg",
				Struct: columnsToStruct(r, gq.MethodName+"Params", cols, settings),
			}
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = QueryValue{
				Name: columnName(c, 0),
				Typ:  goType(r, c, settings),
			}
		} else if len(query.Columns) > 1 {
			var columns []goColumn

			for ci := 0; ci < len(query.Columns); {
				c := query.Columns[ci]
				var embed *Struct

				// Checks for matching structs.
				for _, s := range structs {
					// Ensuring tables match the column selector.
					if c.Table == nil || c.Table.Name != s.Table.Rel {
						continue
					}
					// If the query doesn't have enough fields, it cannot
					// fufill the struct.
					if len(query.Columns) < len(s.Fields) {
						continue
					}
					same := true
					for fi, f := range s.Fields {
						fieldOffset := ci + fi
						// If the location of this field doesn't fit into our columns,
						// we know the struct can't fit either.
						if fieldOffset > len(query.Columns)-1 {
							break
						}
						c := query.Columns[fieldOffset]
						sameName := f.Name == StructName(columnName(c, fieldOffset), settings)
						sameType := f.Type == goType(r, c, settings)
						sameTable := sameTableName(c.Table, s.Table, r.Catalog.DefaultSchema)
						if !sameName || !sameType || !sameTable {
							same = false
						}
					}
					if same {
						embed = &s
						break
					}
				}

				// Used to track the amount of columns matched.
				// A struct could be embedded, and in that case
				// for performance we want to skip over those
				// matched columns.
				colsMatched := 1
				if embed != nil {
					colsMatched = len(embed.Fields)
				}
				for colID := ci; colID < ci+colsMatched; colID++ {
					columns = append(columns, goColumn{
						id:     colID,
						Embed:  embed,
						Column: query.Columns[colID],
					})
				}
				ci += colsMatched
			}

			var emit bool
			var gs *Struct
			// Check if all columns match a consistent embedded struct.
			// If they do, we don't need to generate a new struct for the row.
			for _, c := range columns {
				if gs == nil {
					gs = c.Embed
					continue
				}

				// Cheaper to compare the pointer instead of the name.
				if gs != c.Embed {
					gs = nil
					break
				}
			}

			if gs == nil {
				gs = columnsToStruct(r, gq.MethodName+"Row", columns, settings)
				emit = true
			}
			gq.Ret = QueryValue{
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
func columnsToStruct(r *compiler.Result, name string, columns []goColumn, settings config.CombinedSettings) *Struct {
	gs := Struct{
		Name: name,
	}
	embedded := map[string]interface{}{}
	seen := map[string]int{}
	suffixes := map[int]int{}
	for i, c := range columns {
		if c.Embed != nil {
			if _, ok := embedded[c.Embed.Name]; !ok {
				// We only want to include each embedded struct once.
				gs.Embedded = append(gs.Embedded, *c.Embed)
				embedded[c.Embed.Name] = nil
			}
		}

		colName := columnName(c.Column, i)
		tagName := colName
		fieldName := StructName(colName, settings)
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
		tags := map[string]string{}
		if settings.Go.EmitDBTags {
			tags["db:"] = tagName
		}
		if settings.Go.EmitJSONTags {
			tags["json:"] = JSONTagName(tagName, settings)
		}
		f := Field{
			Name: fieldName,
			Type: goType(r, c.Column, settings),
			Tags: tags,
		}
		if c.Embed != nil {
			f.Struct = c.Embed.Name
		}
		gs.Fields = append(gs.Fields, f)
		seen[colName]++
	}
	return &gs
}
