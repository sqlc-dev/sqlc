package dinosql

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"text/template"

	core "github.com/kyleconroy/sqlc/internal/pg"

	"github.com/jinzhu/inflection"
)

var ktIdentPattern = regexp.MustCompile("[^a-zA-Z0-9_]+")

type KtConstant struct {
	Name  string
	Type  string
	Value string
}

type KtEnum struct {
	Name      string
	Comment   string
	Constants []KtConstant
}

type KtField struct {
	Name    string
	Type    string
	Comment string
}

type KtStruct struct {
	Table   core.FQN
	Name    string
	Fields  []KtField
	Comment string
}

type KtQueryValue struct {
	Emit   bool
	Name   string
	Struct *KtStruct
	Typ    string
}

func (v KtQueryValue) EmitStruct() bool {
	return v.Emit
}

func (v KtQueryValue) IsStruct() bool {
	return v.Struct != nil
}

func (v KtQueryValue) isEmpty() bool {
	return v.Typ == "" && v.Name == "" && v.Struct == nil
}

func (v KtQueryValue) Pair() string {
	if v.isEmpty() {
		return ""
	}
	return v.Name + ": " + v.Type()
}

func (v KtQueryValue) Type() string {
	if v.Typ != "" {
		return v.Typ
	}
	if v.Struct != nil {
		return v.Struct.Name
	}
	panic("no type for KtQueryValue: " + v.Name)
}

func (v KtQueryValue) Params() []KtQueryParam {
	if v.isEmpty() {
		return nil
	}
	var out []KtQueryParam
	if v.Struct == nil {
		if strings.HasPrefix(v.Typ, "[]") && v.Typ != "[]byte" {
			// TODO: this won't compile
			out = append(out, KtQueryParam{
				Name: "pq.Array(" + v.Name + ")",
				Typ:  v.Typ,
			})
		} else {
			out = append(out, KtQueryParam{
				Name: v.Name,
				Typ:  v.Typ,
			})
		}
	} else {
		for _, f := range v.Struct.Fields {
			if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
				out = append(out, KtQueryParam{
					Name: "pq.Array(" + v.Name + "." + f.Name + ")",
					Typ:  f.Type,
				})
			} else {
				out = append(out, KtQueryParam{
					Name: v.Name + "." + f.Name,
					Typ:  f.Type,
				})
			}
		}
	}
	return out
}

func (v KtQueryValue) Scan() string {
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

type KtQueryParam struct {
	Name string
	Typ  string
}

func (p KtQueryParam) Getter() string {
	return "get" + strings.TrimSuffix(p.Typ, "?")
}

func (p KtQueryParam) Setter() string {
	return "set" + strings.TrimSuffix(p.Typ, "?")
}

// A struct used to generate methods and fields on the Queries struct
type KtQuery struct {
	ClassName    string
	Cmd          string
	Comments     []string
	MethodName   string
	FieldName    string
	ConstantName string
	SQL          string
	SourceName   string
	Ret          KtQueryValue
	Arg          KtQueryValue
}

type KtGenerateable interface {
	KtDataClasses(settings CombinedSettings) []KtStruct
	KtQueries(settings CombinedSettings) []KtQuery
	KtEnums(settings CombinedSettings) []KtEnum
}

func KtUsesType(r KtGenerateable, typ string, settings CombinedSettings) bool {
	for _, strct := range r.KtDataClasses(settings) {
		for _, f := range strct.Fields {
			fType := strings.TrimPrefix(f.Type, "[]")
			if strings.HasPrefix(fType, typ) {
				return true
			}
		}
	}
	return false
}

func KtImports(r KtGenerateable, settings CombinedSettings) func(string) [][]string {
	return func(filename string) [][]string {
		if filename == "Models.kt" {
			return ModelKtImports(r, settings)
		}

		if filename == "Querier.kt" {
			return InterfaceKtImports(r, settings)
		}

		return QueryKtImports(r, settings, filename)
	}
}

func InterfaceKtImports(r KtGenerateable, settings CombinedSettings) [][]string {
	gq := r.KtQueries(settings)
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
		"java.sql.Connection":   {},
		"java.sql.SQLException": {},
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

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds}
}

func ModelKtImports(r KtGenerateable, settings CombinedSettings) [][]string {
	std := make(map[string]struct{})
	if KtUsesType(r, "sql.Null", settings) {
		std["database/sql"] = struct{}{}
	}
	if KtUsesType(r, "json.RawMessage", settings) {
		std["encoding/json"] = struct{}{}
	}
	if KtUsesType(r, "time.Time", settings) {
		std["time"] = struct{}{}
	}
	if KtUsesType(r, "net.IP", settings) {
		std["net"] = struct{}{}
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds}
}

func QueryKtImports(r KtGenerateable, settings CombinedSettings, filename string) [][]string {
	// for _, strct := range r.KtDataClasses() {
	// 	for _, f := range strct.Fields {
	// 		if strings.HasPrefix(f.Type, "[]") {
	// 			return true
	// 		}
	// 	}
	// }
	var gq []KtQuery
	for _, query := range r.KtQueries(settings) {
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
		"java.sql.Connection":   {},
		"java.sql.SQLException": {},
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

	if sliceScan() {
		pkg["github.com/lib/pq"] = struct{}{}
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
	return [][]string{stds, pkgs}
}

func ktEnumValueName(value string) string {
	id := strings.Replace(value, "-", "_", -1)
	id = strings.Replace(id, ":", "_", -1)
	id = strings.Replace(id, "/", "_", -1)
	id = ktIdentPattern.ReplaceAllString(id, "")
	return strings.ToUpper(id)
}

func (r Result) KtEnums(settings CombinedSettings) []KtEnum {
	var enums []KtEnum
	for name, schema := range r.Catalog.Schemas {
		if name == "pg_catalog" {
			continue
		}
		for _, enum := range schema.Enums {
			var enumName string
			if name == "public" {
				enumName = enum.Name
			} else {
				enumName = name + "_" + enum.Name
			}
			e := KtEnum{
				Name:    KtDataClassName(enumName, settings),
				Comment: enum.Comment,
			}
			for _, v := range enum.Vals {
				e.Constants = append(e.Constants, KtConstant{
					Name:  ktEnumValueName(v),
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

func KtDataClassName(name string, settings CombinedSettings) string {
	if rename := settings.Global.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	for _, p := range strings.Split(name, "_") {
		out += strings.Title(p)
	}
	return out
}

func KtMemberName(name string, settings CombinedSettings) string {
	return LowerTitle(KtDataClassName(name, settings))
}

func (r Result) KtDataClasses(settings CombinedSettings) []KtStruct {
	var structs []KtStruct
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
			s := KtStruct{
				Table:   core.FQN{Schema: name, Rel: table.Name},
				Name:    inflection.Singular(KtDataClassName(tableName, settings)),
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, KtField{
					Name:    KtMemberName(column.Name, settings),
					Type:    r.ktType(column, settings),
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

func (r Result) ktType(col core.Column, settings CombinedSettings) string {
	typ := r.ktInnerType(col, settings)
	if col.IsArray {
		return fmt.Sprintf("Array<%s>", typ)
	}
	return typ
}

func (r Result) ktInnerType(col core.Column, settings CombinedSettings) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	switch columnType {
	case "serial", "pg_catalog.serial4":
		if notNull {
			return "Int"
		}
		return "Int?"

	case "bigserial", "pg_catalog.serial8":
		if notNull {
			return "Long"
		}
		return "Long?"

	case "smallserial", "pg_catalog.serial2":
		return "Short"

	case "integer", "int", "int4", "pg_catalog.int4":
		if notNull {
			return "Int"
		}
		return "Int?"

	case "bigint", "pg_catalog.int8":
		if notNull {
			return "Long"
		}
		return "Long?"

	case "smallint", "pg_catalog.int2":
		return "Short"

	case "float", "double precision", "pg_catalog.float8":
		if notNull {
			return "Double"
		}
		return "Double?"

	case "real", "pg_catalog.float4":
		if notNull {
			return "Float"
		}
		return "Float?"

	case "pg_catalog.numeric":
		if notNull {
			return "java.math.BigDecimal"
		}
		return "java.math.BigDecimal?"

	case "bool", "pg_catalog.bool":
		if notNull {
			return "Boolean"
		}
		return "Boolean?"

	case "jsonb":
		// TODO: support json and byte types
		return "String"

	case "bytea", "blob", "pg_catalog.bytea":
		return "String"

	case "date":
		// TODO
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.time", "pg_catalog.timetz":
		// TODO
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.timestamp", "pg_catalog.timestamptz", "timestamptz":
		// TODO
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string":
		if notNull {
			return "String"
		}
		return "String?"

	case "uuid":
		// TODO
		return "uuid.UUID"

	case "inet":
		// TODO
		return "net.IP"

	case "void":
		// TODO
		// A void value always returns NULL. Since there is no built-in NULL
		// value into the SQL package, we'll use sql.NullBool
		return "sql.NullBool"

	case "any":
		// TODO
		return "Any"

	default:
		for name, schema := range r.Catalog.Schemas {
			if name == "pg_catalog" {
				continue
			}
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					if name == "public" {
						return KtDataClassName(enum.Name, settings)
					}

					return KtDataClassName(name+"_"+enum.Name, settings)
				}
			}
		}
		log.Printf("unknown PostgreSQL type: %s\n", columnType)
		return "interface{}"
	}
}

// It's possible that this method will generate duplicate JSON tag values
//
//   Columns: count, count,   count_2
//    Fields: Count, Count_2, Count2
// JSON tags: count, count_2, count_2
//
// This is unlikely to happen, so don't fix it yet
func (r Result) ktColumnsToStruct(name string, columns []core.Column, settings CombinedSettings) *KtStruct {
	gs := KtStruct{
		Name: name,
	}
	seen := map[string]int{}
	for i, c := range columns {
		fieldName := KtMemberName(ktColumnName(c, i), settings)
		if v := seen[c.Name]; v > 0 {
			fieldName = fmt.Sprintf("%s_%d", fieldName, v+1)
		}
		gs.Fields = append(gs.Fields, KtField{
			Name: fieldName,
			Type: r.ktType(c, settings),
		})
		seen[c.Name]++
	}
	return &gs
}

func ktArgName(name string) string {
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

func ktParamName(p Parameter) string {
	if p.Column.Name != "" {
		return ktArgName(p.Column.Name)
	}
	return fmt.Sprintf("dollar_%d", p.Number)
}

func ktColumnName(c core.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}

var jdbcSQLRe = regexp.MustCompile(`\B\$\d+\b`)

// HACK: jdbc doesn't support numbered parameters, so we need to transform them to question marks...
// But there's no access to the SQL parser here, so we just do a dumb regexp replace instead. This won't work if
// the literal strings contain matching values, but good enough for a prototype.
func jdbcSQL(s string) string {
	return jdbcSQLRe.ReplaceAllString(s, "?")
}

func (r Result) KtQueries(settings CombinedSettings) []KtQuery {
	structs := r.KtDataClasses(settings)

	qs := make([]KtQuery, 0, len(r.Queries))
	for _, query := range r.Queries {
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}

		gq := KtQuery{
			Cmd:          query.Cmd,
			ClassName:    strings.Title(query.Name),
			ConstantName: LowerTitle(query.Name),
			FieldName:    LowerTitle(query.Name) + "Stmt",
			MethodName:   LowerTitle(query.Name),
			SourceName:   query.Filename,
			SQL:          jdbcSQL(query.SQL),
			Comments:     query.Comments,
		}

		if len(query.Params) == 1 {
			p := query.Params[0]
			gq.Arg = KtQueryValue{
				Name: ktParamName(p),
				Typ:  r.ktType(p.Column, settings),
			}
		} else if len(query.Params) > 1 {
			var cols []core.Column
			for _, p := range query.Params {
				cols = append(cols, p.Column)
			}
			gq.Arg = KtQueryValue{
				Emit:   true,
				Name:   "arg",
				Struct: r.ktColumnsToStruct(gq.ClassName+"Params", cols, settings),
			}
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = KtQueryValue{
				Name: ktColumnName(c, 0),
				Typ:  r.ktType(c, settings),
			}
		} else if len(query.Columns) > 1 {
			var gs *KtStruct
			var emit bool

			for _, s := range structs {
				if len(s.Fields) != len(query.Columns) {
					continue
				}
				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					sameName := f.Name == KtMemberName(ktColumnName(c, i), settings)
					sameType := f.Type == r.ktType(c, settings)
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
				gs = r.ktColumnsToStruct(gq.ClassName+"Row", query.Columns, settings)
				emit = true
			}
			gq.Ret = KtQueryValue{
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

var ktIfaceTmpl = `// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

{{range imports .SourceName}}
{{range .}}import {{.}}
{{end}}
{{end}}

interface Querier {
	{{- range .KtQueries}}
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
`

var ktModelsTmpl = `// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

{{range imports .SourceName}}
{{range .}}import {{.}}
{{end}}
{{end}}

{{range .Enums}}
{{if .Comment}}// {{.Comment}}{{end}}
enum class {{.Name}}(val value: String) {
  {{- range $i, $e := .Constants}}
  {{- if $i }},{{end}}
  {{.Name}}("{{.Value}}")
  {{- end}}
}
{{end}}

{{range .KtDataClasses}}
{{if .Comment}}// {{.Comment}}{{end}}
data class {{.Name}} ( {{- range $i, $e := .Fields}}
  {{- if $i }},{{end}}
  {{- if .Comment}}
  // {{.Comment}}{{else}}
  {{- end}}
  val {{.Name}}: {{.Type}}
  {{- end}}
)
{{end}}
`

var ktSqlTmpl = `// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

{{range imports .SourceName}}
{{range .}}import {{.}}
{{end}}
{{end}}

{{range .KtQueries}}
const val {{.ConstantName}} = {{$.Q}}-- name: {{.MethodName}} {{.Cmd}}
{{.SQL}}
{{$.Q}}

{{if .Arg.EmitStruct}}
data class {{.Arg.Type}} ( {{- range $i, $e := .Arg.Struct.Fields}}
  {{- if $i }},{{end}}
  val {{.Name}}: {{.Type}}
  {{- end}}
)
{{end}}

{{if .Ret.EmitStruct}}
data class {{.Ret.Type}} ( {{- range $i, $e := .Ret.Struct.Fields}}
  {{- if $i }},{{end}}
  val {{.Name}}: {{.Type}}
  {{- end}}
)
{{end}}
{{end}}

class Queries(private val conn: Connection) {
{{range .KtQueries}}
{{if eq .Cmd ":one"}}
{{range .Comments}}//{{.}}
{{end}}
  @Throws(SQLException::class)
  fun {{.MethodName}}({{.Arg.Pair}}): {{.Ret.Type}} {
    val stmt = conn.prepareStatement({{.ConstantName}}) {{- range $i, $e := .Arg.Params }}
    stmt.{{.Setter}}({{offset $i}}, {{.Name}})
    {{- end}}

    val results = stmt.executeQuery()
    if (!results.next()) {
      throw SQLException("no rows in result set")
    }
    {{ if .Ret.IsStruct }}
    val ret = {{.Ret.Type}}( {{- range $i, $e := .Ret.Params }}
      {{- if $i }},{{end}}
      results.{{.Getter}}({{offset $i}})
    {{- end -}}
    )
    {{ else }}
    val ret = results.{{(index .Ret.Params 0).Getter}}(1)
    {{ end }}
    if (results.next()) {
        throw SQLException("expected one row in result set, but got many")
    }
    return ret
  }
{{end}}

{{if eq .Cmd ":many"}}
{{range .Comments}}//{{.}}
{{end}}
  @Throws(SQLException::class)
  fun {{.MethodName}}({{.Arg.Pair}}): List<{{.Ret.Type}}> {
    val stmt = conn.prepareStatement({{.ConstantName}}) {{- range $i, $e := .Arg.Params }}
    stmt.{{.Setter}}({{offset $i}}, {{.Name}})
    {{- end}}

    val results = stmt.executeQuery()
    val ret = mutableListOf<{{.Ret.Type}}>()
    while (results.next()) {
    {{ if .Ret.IsStruct }}
        ret.add({{.Ret.Type}}( {{- range $i, $e := .Ret.Params }}
          {{- if $i }},{{end}}
          results.{{.Getter}}({{offset $i}})
          {{- end -}}
        ))
    {{ else }}
        ret.add(results.{{(index .Ret.Params 0).Getter}}(1)
    {{ end }}
    }
    return ret
  }
{{end}}

{{if eq .Cmd ":exec"}}
{{range .Comments}}//{{.}}
{{end}}
  @Throws(SQLException::class)
  fun {{.MethodName}}({{.Arg.Pair}}) {
    val stmt = conn.prepareStatement({{.ConstantName}}) {{- range $i, $e := .Arg.Params }}
    stmt.{{.Setter}}({{offset $i}}, {{.Name}})
    {{- end}}

    stmt.execute()
  }
{{end}}

{{if eq .Cmd ":execrows"}}
{{range .Comments}}//{{.}}
{{end}}
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
}
`

type ktTmplCtx struct {
	Q             string
	Package       string
	Enums         []KtEnum
	KtDataClasses []KtStruct
	KtQueries     []KtQuery
	Settings      GenerateSettings

	// TODO: Race conditions
	SourceName string

	EmitJSONTags        bool
	EmitPreparedQueries bool
	EmitInterface       bool
}

func Offset(v int) int {
	return v + 1
}

func ktFormat(s string) string {
	// TODO: do more than just skip multiple blank lines, like maybe run ktlint to format
	skipNextSpace := false
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		isSpace := len(strings.TrimSpace(l)) == 0
		if !isSpace || !skipNextSpace {
			lines = append(lines, l)
		}
		skipNextSpace = isSpace
	}
	o := strings.Join(lines, "\n")
	o += "\n"
	return o
}

func KtGenerate(r KtGenerateable, settings CombinedSettings) (map[string]string, error) {
	funcMap := template.FuncMap{
		"lowerTitle": LowerTitle,
		"imports":    KtImports(r, settings),
		"offset":     Offset,
	}

	modelsFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktModelsTmpl))
	sqlFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktSqlTmpl))
	ifaceFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktIfaceTmpl))

	pkg := settings.Package
	tctx := ktTmplCtx{
		Settings:            settings.Global,
		EmitInterface:       pkg.EmitInterface,
		EmitJSONTags:        pkg.EmitJSONTags,
		EmitPreparedQueries: pkg.EmitPreparedQueries,
		Q:                   `"""`,
		Package:             pkg.Name,
		KtQueries:           r.KtQueries(settings),
		Enums:               r.KtEnums(settings),
		KtDataClasses:       r.KtDataClasses(settings),
	}

	output := map[string]string{}

	execute := func(name string, t *template.Template) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := t.Execute(w, tctx)
		w.Flush()
		if err != nil {
			return err
		}
		if !strings.HasSuffix(name, ".kt") {
			name += ".kt"
		}
		output[name] = ktFormat(b.String())
		return nil
	}

	if err := execute("Models.kt", modelsFile); err != nil {
		return nil, err
	}
	if pkg.EmitInterface {
		if err := execute("Querier.kt", ifaceFile); err != nil {
			return nil, err
		}
	}
	if err := execute("Queries.kt", sqlFile); err != nil {
		return nil, err
	}

	return output, nil
}
