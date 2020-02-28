package dinosql

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/kyleconroy/sqlc/internal/catalog"
	"github.com/kyleconroy/sqlc/internal/config"
	core "github.com/kyleconroy/sqlc/internal/pg"

	"github.com/jinzhu/inflection"
)

var identPattern = regexp.MustCompile("[^a-zA-Z0-9_]+")

type GoConstant struct {
	Name  string
	Type  string
	Value string
}

type GoEnum struct {
	Name      string
	Comment   string
	Constants []GoConstant
}

type GoField struct {
	Name    string
	Type    string
	Tags    map[string]string
	Comment string
}

func (gf GoField) Tag() string {
	tags := make([]string, 0, len(gf.Tags))
	for key, val := range gf.Tags {
		tags = append(tags, fmt.Sprintf("%s\"%s\"", key, val))
	}
	if len(tags) == 0 {
		return ""
	}
	sort.Strings(tags)
	return strings.Join(tags, ",")
}

type GoStruct struct {
	Table   core.FQN
	Name    string
	Fields  []GoField
	Comment string
}

type GoQueryValue struct {
	Emit   bool
	Name   string
	Struct *GoStruct
	Typ    string
}

func (v GoQueryValue) EmitStruct() bool {
	return v.Emit
}

func (v GoQueryValue) IsStruct() bool {
	return v.Struct != nil
}

func (v GoQueryValue) isEmpty() bool {
	return v.Typ == "" && v.Name == "" && v.Struct == nil
}

func (v GoQueryValue) Pair() string {
	if v.isEmpty() {
		return ""
	}
	return v.Name + " " + v.Type()
}

func (v GoQueryValue) Type() string {
	if v.Typ != "" {
		return v.Typ
	}
	if v.Struct != nil {
		return v.Struct.Name
	}
	panic("no type for GoQueryValue: " + v.Name)
}

func (v GoQueryValue) Params() string {
	if v.isEmpty() {
		return ""
	}
	var out []string
	if v.Struct == nil {
		if strings.HasPrefix(v.Typ, "[]") && v.Typ != "[]byte" {
			out = append(out, "pq.Array("+v.Name+")")
		} else {
			out = append(out, v.Name)
		}
	} else {
		for _, f := range v.Struct.Fields {
			if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
				out = append(out, "pq.Array("+v.Name+"."+f.Name+")")
			} else {
				out = append(out, v.Name+"."+f.Name)
			}
		}
	}
	if len(out) <= 3 {
		return strings.Join(out, ",")
	}
	out = append(out, "")
	return "\n" + strings.Join(out, ",\n")
}

func (v GoQueryValue) Scan() string {
	var out []string
	if v.Struct == nil {
		if strings.HasPrefix(v.Typ, "[]") && v.Typ != "[]byte" {
			out = append(out, "pq.Array(&"+v.Name+")")
		} else {
			out = append(out, "&"+v.Name)
		}
	} else {
		for _, f := range v.Struct.Fields {
			if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
				out = append(out, "pq.Array(&"+v.Name+"."+f.Name+")")
			} else {
				out = append(out, "&"+v.Name+"."+f.Name)
			}
		}
	}
	if len(out) <= 3 {
		return strings.Join(out, ",")
	}
	out = append(out, "")
	return "\n" + strings.Join(out, ",\n")
}

// A struct used to generate methods and fields on the Queries struct
type GoQuery struct {
	Cmd          string
	Comments     []string
	MethodName   string
	FieldName    string
	ConstantName string
	SQL          string
	SourceName   string
	Ret          GoQueryValue
	Arg          GoQueryValue
}

type Generateable interface {
	Structs(settings config.CombinedSettings) []GoStruct
	GoQueries(settings config.CombinedSettings) []GoQuery
	Enums(settings config.CombinedSettings) []GoEnum
}

func UsesType(r Generateable, typ string, settings config.CombinedSettings) bool {
	for _, strct := range r.Structs(settings) {
		for _, f := range strct.Fields {
			fType := strings.TrimPrefix(f.Type, "[]")
			if strings.HasPrefix(fType, typ) {
				return true
			}
		}
	}
	return false
}

func UsesArrays(r Generateable, settings config.CombinedSettings) bool {
	for _, strct := range r.Structs(settings) {
		for _, f := range strct.Fields {
			if strings.HasPrefix(f.Type, "[]") {
				return true
			}
		}
	}
	return false
}

type fileImports struct {
	Std []string
	Dep []string
}

func mergeImports(imps ...fileImports) [][]string {
	if len(imps) == 1 {
		return [][]string{imps[0].Std, imps[0].Dep}
	}

	var stds, pkgs []string
	seenStd := map[string]struct{}{}
	seenPkg := map[string]struct{}{}
	for i := range imps {
		for _, std := range imps[i].Std {
			if _, ok := seenStd[std]; ok {
				continue
			}
			stds = append(stds, std)
			seenStd[std] = struct{}{}
		}
		for _, pkg := range imps[i].Dep {
			if _, ok := seenPkg[pkg]; ok {
				continue
			}
			pkgs = append(pkgs, pkg)
			seenPkg[pkg] = struct{}{}
		}
	}
	return [][]string{stds, pkgs}
}

func Imports(r Generateable, settings config.CombinedSettings) func(string) [][]string {
	return func(filename string) [][]string {
		if filename == "db.go" {
			return mergeImports(dbImports(r, settings))
		}

		if filename == "models.go" {
			return mergeImports(modelImports(r, settings))
		}

		if filename == "querier.go" {
			return mergeImports(interfaceImports(r, settings))
		}

		return mergeImports(queryImports(r, settings, filename))
	}
}

func dbImports(r Generateable, settings config.CombinedSettings) fileImports {
	std := []string{"context", "database/sql"}
	if settings.Go.EmitPreparedQueries {
		std = append(std, "fmt")
	}
	return fileImports{Std: std}
}

func interfaceImports(r Generateable, settings config.CombinedSettings) fileImports {
	gq := r.GoQueries(settings)
	uses := func(name string) bool {
		for _, q := range gq {
			if !q.Ret.isEmpty() {
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				if strings.HasPrefix(q.Arg.Type(), name) {
					return true
				}
			}
		}
		return false
	}

	std := map[string]struct{}{
		"context": struct{}{},
	}
	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}
	if uses("json.RawMessage") {
		std["encoding/json"] = struct{}{}
	}
	if uses("time.Time") {
		std["time"] = struct{}{}
	}
	if uses("net.IP") {
		std["net"] = struct{}{}
	}
	if uses("net.HardwareAddr") {
		std["net"] = struct{}{}
	}

	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if uses("pq.NullTime") && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}
	_, overrideUUID := overrideTypes["uuid.UUID"]
	if uses("uuid.UUID") && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	// Custom imports
	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && uses(goType) {
			pkg[importPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(pkg))
	for p, _ := range pkg {
		pkgs = append(pkgs, p)
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	sort.Strings(pkgs)
	return fileImports{stds, pkgs}
}

func modelImports(r Generateable, settings config.CombinedSettings) fileImports {
	std := make(map[string]struct{})
	if UsesType(r, "sql.Null", settings) {
		std["database/sql"] = struct{}{}
	}
	if UsesType(r, "json.RawMessage", settings) {
		std["encoding/json"] = struct{}{}
	}
	if UsesType(r, "time.Time", settings) {
		std["time"] = struct{}{}
	}
	if UsesType(r, "net.IP", settings) {
		std["net"] = struct{}{}
	}
	if UsesType(r, "net.HardwareAddr", settings) {
		std["net"] = struct{}{}
	}

	// Custom imports
	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if UsesType(r, "pq.NullTime", settings) && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}

	_, overrideUUID := overrideTypes["uuid.UUID"]
	if UsesType(r, "uuid.UUID", settings) && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && UsesType(r, goType, settings) {
			pkg[importPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(pkg))
	for p, _ := range pkg {
		pkgs = append(pkgs, p)
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	sort.Strings(pkgs)
	return fileImports{stds, pkgs}
}

func queryImports(r Generateable, settings config.CombinedSettings, filename string) fileImports {
	// for _, strct := range r.Structs() {
	// 	for _, f := range strct.Fields {
	// 		if strings.HasPrefix(f.Type, "[]") {
	// 			return true
	// 		}
	// 	}
	// }
	var gq []GoQuery
	for _, query := range r.GoQueries(settings) {
		if query.SourceName == filename {
			gq = append(gq, query)
		}
	}

	uses := func(name string) bool {
		for _, q := range gq {
			if !q.Ret.isEmpty() {
				if q.Ret.EmitStruct() {
					for _, f := range q.Ret.Struct.Fields {
						fType := strings.TrimPrefix(f.Type, "[]")
						if strings.HasPrefix(fType, name) {
							return true
						}
					}
				}
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				if q.Arg.EmitStruct() {
					for _, f := range q.Arg.Struct.Fields {
						fType := strings.TrimPrefix(f.Type, "[]")
						if strings.HasPrefix(fType, name) {
							return true
						}
					}
				}
				if strings.HasPrefix(q.Arg.Type(), name) {
					return true
				}
			}
		}
		return false
	}

	sliceScan := func() bool {
		for _, q := range gq {
			if !q.Ret.isEmpty() {
				if q.Ret.IsStruct() {
					for _, f := range q.Ret.Struct.Fields {
						if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Ret.Type(), "[]") && q.Ret.Type() != "[]byte" {
						return true
					}
				}
			}
			if !q.Arg.isEmpty() {
				if q.Arg.IsStruct() {
					for _, f := range q.Arg.Struct.Fields {
						if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Arg.Type(), "[]") && q.Arg.Type() != "[]byte" {
						return true
					}
				}
			}
		}
		return false
	}

	std := map[string]struct{}{
		"context": struct{}{},
	}
	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}
	if uses("json.RawMessage") {
		std["encoding/json"] = struct{}{}
	}
	if uses("time.Time") {
		std["time"] = struct{}{}
	}
	if uses("net.IP") {
		std["net"] = struct{}{}
	}

	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
	}

	if sliceScan() {
		pkg["github.com/lib/pq"] = struct{}{}
	}
	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if uses("pq.NullTime") && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}
	_, overrideUUID := overrideTypes["uuid.UUID"]
	if uses("uuid.UUID") && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	// Custom imports
	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && uses(goType) {
			pkg[importPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(pkg))
	for p, _ := range pkg {
		pkgs = append(pkgs, p)
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	sort.Strings(pkgs)
	return fileImports{stds, pkgs}
}

func enumValueName(value string) string {
	name := ""
	id := strings.Replace(value, "-", "_", -1)
	id = strings.Replace(id, ":", "_", -1)
	id = strings.Replace(id, "/", "_", -1)
	id = identPattern.ReplaceAllString(id, "")
	for _, part := range strings.Split(id, "_") {
		name += strings.Title(part)
	}
	return name
}

func (r Result) Enums(settings config.CombinedSettings) []GoEnum {
	var enums []GoEnum
	for name, schema := range r.Catalog.Schemas {
		if name == "pg_catalog" {
			continue
		}
		for _, enum := range schema.Enums() {
			var enumName string
			if name == "public" {
				enumName = enum.Name
			} else {
				enumName = name + "_" + enum.Name
			}
			e := GoEnum{
				Name:    StructName(enumName, settings),
				Comment: enum.Comment,
			}
			for _, v := range enum.Vals {
				e.Constants = append(e.Constants, GoConstant{
					Name:  e.Name + enumValueName(v),
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

func StructName(name string, settings config.CombinedSettings) string {
	if rename := settings.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

func (r Result) Structs(settings config.CombinedSettings) []GoStruct {
	var structs []GoStruct
	for name, schema := range r.Catalog.Schemas {
		if name == "pg_catalog" {
			continue
		}
		for _, table := range schema.Tables {
			var tableName string
			if name == "public" {
				tableName = table.Name
			} else {
				tableName = name + "_" + table.Name
			}
			s := GoStruct{
				Table:   core.FQN{Schema: name, Rel: table.Name},
				Name:    inflection.Singular(StructName(tableName, settings)),
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, GoField{
					Name:    StructName(column.Name, settings),
					Type:    r.goType(column, settings),
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

func (r Result) goType(col core.Column, settings config.CombinedSettings) string {
	// package overrides have a higher precedence
	for _, oride := range settings.Overrides {
		if oride.Column != "" && oride.ColumnName == col.Name && oride.Table == col.Table {
			return oride.GoTypeName
		}
	}
	typ := r.goInnerType(col, settings)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func (r Result) goInnerType(col core.Column, settings config.CombinedSettings) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	// package overrides have a higher precedence
	for _, oride := range settings.Overrides {
		if oride.DBType != "" && oride.DBType == columnType && oride.Null != notNull {
			return oride.GoTypeName
		}
	}

	switch columnType {
	case "serial", "pg_catalog.serial4":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigserial", "pg_catalog.serial8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "smallserial", "pg_catalog.serial2":
		return "int16"

	case "integer", "int", "int4", "pg_catalog.int4":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigint", "pg_catalog.int8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "smallint", "pg_catalog.int2":
		return "int16"

	case "float", "double precision", "pg_catalog.float8":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "real", "pg_catalog.float4":
		if notNull {
			return "float32"
		}
		return "sql.NullFloat64" // TODO: Change to sql.NullFloat32 after updating the go.mod file

	case "pg_catalog.numeric":
		// Since the Go standard library does not have a decimal type, lib/pq
		// returns numerics as strings.
		//
		// https://github.com/lib/pq/issues/648
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "bool", "pg_catalog.bool":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	case "jsonb":
		return "json.RawMessage"

	case "bytea", "blob", "pg_catalog.bytea":
		return "[]byte"

	case "date":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.time", "pg_catalog.timetz":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.timestamp", "pg_catalog.timestamptz", "timestamptz":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "uuid":
		return "uuid.UUID"

	case "inet":
		return "net.IP"

	case "macaddr", "macaddr8":
		return "net.HardwareAddr"

	case "void":
		// A void value always returns NULL. Since there is no built-in NULL
		// value into the SQL package, we'll use sql.NullBool
		return "sql.NullBool"

	case "any":
		return "interface{}"

	default:
		fqn, err := catalog.ParseString(columnType)
		if err != nil {
			// TODO: Should this actually return an error here?
			return "interface{}"
		}

		for name, schema := range r.Catalog.Schemas {
			if name == "pg_catalog" {
				continue
			}
			for _, typ := range schema.Types {
				switch t := typ.(type) {
				case core.Enum:
					if fqn.Rel == t.Name && fqn.Schema == name {
						if name == "public" {
							return StructName(t.Name, settings)
						}
						return StructName(name+"_"+t.Name, settings)
					}
				case core.CompositeType:
					if notNull {
						return "string"
					}
					return "sql.NullString"
				}
			}
		}

		log.Printf("unknown PostgreSQL type: %s\n", columnType)
		return "interface{}"
	}
}

type goColumn struct {
	id int
	core.Column
}

// It's possible that this method will generate duplicate JSON tag values
//
//   Columns: count, count,   count_2
//    Fields: Count, Count_2, Count2
// JSON tags: count, count_2, count_2
//
// This is unlikely to happen, so don't fix it yet
func (r Result) columnsToStruct(name string, columns []goColumn, settings config.CombinedSettings) *GoStruct {
	gs := GoStruct{
		Name: name,
	}
	seen := map[string]int{}
	suffixes := map[int]int{}
	for i, c := range columns {
		tagName := c.Name
		fieldName := StructName(columnName(c.Column, i), settings)
		// Track suffixes by the ID of the column, so that columns referring to the same numbered parameter can be
		// reused.
		suffix := 0
		if o, ok := suffixes[c.id]; ok {
			suffix = o
		} else if v := seen[c.Name]; v > 0 {
			suffix = v + 1
		}
		suffixes[c.id] = suffix
		if suffix > 0 {
			tagName = fmt.Sprintf("%s_%d", tagName, suffix)
			fieldName = fmt.Sprintf("%s_%d", fieldName, suffix)
		}
		gs.Fields = append(gs.Fields, GoField{
			Name: fieldName,
			Type: r.goType(c.Column, settings),
			Tags: map[string]string{"json:": tagName},
		})
		seen[c.Name]++
	}
	return &gs
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

func columnName(c core.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}

func compareFQN(a *core.FQN, b *core.FQN) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Catalog == b.Catalog && a.Schema == b.Schema && a.Rel == b.Rel
}

func (r Result) GoQueries(settings config.CombinedSettings) []GoQuery {
	structs := r.Structs(settings)

	qs := make([]GoQuery, 0, len(r.Queries))
	for _, query := range r.Queries {
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}

		gq := GoQuery{
			Cmd:          query.Cmd,
			ConstantName: LowerTitle(query.Name),
			FieldName:    LowerTitle(query.Name) + "Stmt",
			MethodName:   query.Name,
			SourceName:   query.Filename,
			SQL:          query.SQL,
			Comments:     query.Comments,
		}

		if len(query.Params) == 1 {
			p := query.Params[0]
			gq.Arg = GoQueryValue{
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
			gq.Arg = GoQueryValue{
				Emit:   true,
				Name:   "arg",
				Struct: r.columnsToStruct(gq.MethodName+"Params", cols, settings),
			}
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = GoQueryValue{
				Name: columnName(c, 0),
				Typ:  r.goType(c, settings),
			}
		} else if len(query.Columns) > 1 {
			var gs *GoStruct
			var emit bool

			for _, s := range structs {
				if len(s.Fields) != len(query.Columns) {
					continue
				}
				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					sameName := f.Name == StructName(columnName(c, i), settings)
					sameType := f.Type == r.goType(c, settings)
					sameTable := s.Table.Catalog == c.Table.Catalog && s.Table.Schema == c.Table.Schema && s.Table.Rel == c.Table.Rel

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
			gq.Ret = GoQueryValue{
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

var templateSet = `
{{define "dbFile"}}// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{template "dbCode" . }}
{{end}}

{{define "dbCode"}}
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

{{if .EmitPreparedQueries}}
func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	{{- if eq (len .GoQueries) 0 }}
	_ = err
	{{- end }}
	{{- range .GoQueries }}
	if q.{{.FieldName}}, err = db.PrepareContext(ctx, {{.ConstantName}}); err != nil {
		return nil, fmt.Errorf("error preparing query {{.MethodName}}: %w", err)
	}
	{{- end}}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	{{- range .GoQueries }}
	if q.{{.FieldName}} != nil {
		if cerr := q.{{.FieldName}}.Close(); cerr != nil {
			err = fmt.Errorf("error closing {{.FieldName}}: %w", cerr)
		}
	}
	{{- end}}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Row) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}
{{end}}

type Queries struct {
	db DBTX

    {{- if .EmitPreparedQueries}}
	tx         *sql.Tx
	{{- range .GoQueries}}
	{{.FieldName}}  *sql.Stmt
	{{- end}}
	{{- end}}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
     	{{- if .EmitPreparedQueries}}
		tx: tx,
		{{- range .GoQueries}}
		{{.FieldName}}: q.{{.FieldName}},
		{{- end}}
		{{- end}}
	}
}
{{end}}

{{define "interfaceFile"}}// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{template "interfaceCode" . }}
{{end}}

{{define "interfaceCode"}}
type Querier interface {
	{{- range .GoQueries}}
	{{- if eq .Cmd ":one"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) ({{.Ret.Type}}, error)
	{{- end}}
	{{- if eq .Cmd ":many"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) ([]{{.Ret.Type}}, error)
	{{- end}}
	{{- if eq .Cmd ":exec"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) error
	{{- end}}
	{{- if eq .Cmd ":execrows"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) (int64, error)
	{{- end}}
	{{- end}}
}

var _ Querier = (*Queries)(nil)
{{end}}

{{define "modelsFile"}}// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{template "modelsCode" . }}
{{end}}

{{define "modelsCode"}}
{{range .Enums}}
{{if .Comment}}{{comment .Comment}}{{end}}
type {{.Name}} string

const (
	{{- range .Constants}}
	{{.Name}} {{.Type}} = "{{.Value}}"
	{{- end}}
)

func (e *{{.Name}}) Scan(src interface{}) error {
	*e = {{.Name}}(src.([]byte))
	return nil
}
{{end}}

{{range .Structs}}
{{if .Comment}}{{comment .Comment}}{{end}}
type {{.Name}} struct { {{- range .Fields}}
  {{- if .Comment}}
  // {{.Comment}}{{else}}
  {{- end}}
  {{.Name}} {{.Type}} {{if $.EmitJSONTags}}{{$.Q}}{{.Tag}}{{$.Q}}{{end}}
  {{- end}}
}
{{end}}
{{end}}

{{define "queryFile"}}// Code generated by sqlc. DO NOT EDIT.
// source: {{.SourceName}}

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{template "queryCode" . }}
{{end}}

{{define "queryCode"}}
{{range .GoQueries}}
{{if $.OutputQuery .SourceName}}
const {{.ConstantName}} = {{$.Q}}-- name: {{.MethodName}} {{.Cmd}}
{{.SQL}}
{{$.Q}}

{{if .Arg.EmitStruct}}
type {{.Arg.Type}} struct { {{- range .Arg.Struct.Fields}}
  {{.Name}} {{.Type}} {{if $.EmitJSONTags}}{{$.Q}}{{.Tag}}{{$.Q}}{{end}}
  {{- end}}
}
{{end}}

{{if .Ret.EmitStruct}}
type {{.Ret.Type}} struct { {{- range .Ret.Struct.Fields}}
  {{.Name}} {{.Type}} {{if $.EmitJSONTags}}{{$.Q}}{{.Tag}}{{$.Q}}{{end}}
  {{- end}}
}
{{end}}

{{if eq .Cmd ":one"}}
{{range .Comments}}//{{.}}
{{end -}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) ({{.Ret.Type}}, error) {
  	{{- if $.EmitPreparedQueries}}
	row := q.queryRow(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
	{{- else}}
	row := q.db.QueryRowContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
	{{- end}}
	var {{.Ret.Name}} {{.Ret.Type}}
	err := row.Scan({{.Ret.Scan}})
	return {{.Ret.Name}}, err
}
{{end}}

{{if eq .Cmd ":many"}}
{{range .Comments}}//{{.}}
{{end -}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) ([]{{.Ret.Type}}, error) {
  	{{- if $.EmitPreparedQueries}}
	rows, err := q.query(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
  	{{- else}}
	rows, err := q.db.QueryContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
  	{{- end}}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []{{.Ret.Type}}
	for rows.Next() {
		var {{.Ret.Name}} {{.Ret.Type}}
		if err := rows.Scan({{.Ret.Scan}}); err != nil {
			return nil, err
		}
		items = append(items, {{.Ret.Name}})
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
{{end}}

{{if eq .Cmd ":exec"}}
{{range .Comments}}//{{.}}
{{end -}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) error {
  	{{- if $.EmitPreparedQueries}}
	_, err := q.exec(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
  	{{- else}}
	_, err := q.db.ExecContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
  	{{- end}}
	return err
}
{{end}}

{{if eq .Cmd ":execrows"}}
{{range .Comments}}//{{.}}
{{end -}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) (int64, error) {
  	{{- if $.EmitPreparedQueries}}
	result, err := q.exec(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
  	{{- else}}
	result, err := q.db.ExecContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
  	{{- end}}
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
{{end}}
{{end}}
{{end}}
{{end}}
`

type tmplCtx struct {
	Q         string
	Package   string
	Enums     []GoEnum
	Structs   []GoStruct
	GoQueries []GoQuery
	Settings  config.Config

	// TODO: Race conditions
	SourceName string

	EmitJSONTags        bool
	EmitPreparedQueries bool
	EmitInterface       bool
}

func (t *tmplCtx) OutputQuery(sourceName string) bool {
	return t.SourceName == sourceName
}

func LowerTitle(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func DoubleSlashComment(s string) string {
	return "// " + strings.ReplaceAll(s, "\n", "\n// ")
}

func Generate(r Generateable, settings config.CombinedSettings) (map[string]string, error) {
	funcMap := template.FuncMap{
		"lowerTitle": LowerTitle,
		"comment":    DoubleSlashComment,
		"imports":    Imports(r, settings),
	}

	tmpl := template.Must(template.New("table").Funcs(funcMap).Parse(templateSet))

	golang := settings.Go
	tctx := tmplCtx{
		Settings:            settings.Global,
		EmitInterface:       golang.EmitInterface,
		EmitJSONTags:        golang.EmitJSONTags,
		EmitPreparedQueries: golang.EmitPreparedQueries,
		Q:                   "`",
		Package:             golang.Package,
		GoQueries:           r.GoQueries(settings),
		Enums:               r.Enums(settings),
		Structs:             r.Structs(settings),
	}

	output := map[string]string{}

	execute := func(name, templateName string) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := tmpl.ExecuteTemplate(w, templateName, &tctx)
		w.Flush()
		if err != nil {
			return err
		}
		code, err := format.Source(b.Bytes())
		if err != nil {
			fmt.Println(b.String())
			return fmt.Errorf("source error: %w", err)
		}
		if !strings.HasSuffix(name, ".go") {
			name += ".go"
		}
		output[name] = string(code)
		return nil
	}

	if err := execute("db.go", "dbFile"); err != nil {
		return nil, err
	}
	if err := execute("models.go", "modelsFile"); err != nil {
		return nil, err
	}
	if golang.EmitInterface {
		if err := execute("querier.go", "interfaceFile"); err != nil {
			return nil, err
		}
	}

	files := map[string]struct{}{}
	for _, gq := range r.GoQueries(settings) {
		files[gq.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, "queryFile"); err != nil {
			return nil, err
		}
	}
	return output, nil
}
