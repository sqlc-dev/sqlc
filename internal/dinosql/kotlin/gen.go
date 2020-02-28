package kotlin

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/dinosql"
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
	Type    ktType
	Comment string
}

type KtStruct struct {
	Table             core.FQN
	Name              string
	Fields            []KtField
	JDBCParamBindings []KtField
	Comment           string
}

type KtQueryValue struct {
	Emit               bool
	Name               string
	Struct             *KtStruct
	Typ                ktType
	JDBCParamBindCount int
}

func (v KtQueryValue) EmitStruct() bool {
	return v.Emit
}

func (v KtQueryValue) IsStruct() bool {
	return v.Struct != nil
}

func (v KtQueryValue) isEmpty() bool {
	return v.Typ == (ktType{}) && v.Name == "" && v.Struct == nil
}

func (v KtQueryValue) Type() string {
	if v.Typ != (ktType{}) {
		return v.Typ.String()
	}
	if v.Struct != nil {
		return v.Struct.Name
	}
	panic("no type for KtQueryValue: " + v.Name)
}

func jdbcSet(t ktType, idx int, name string) string {
	if t.IsEnum && t.IsArray {
		return fmt.Sprintf(`stmt.setArray(%d, conn.createArrayOf("%s", %s.map { v -> v.value }.toTypedArray()))`, idx, t.DataType, name)
	}
	if t.IsEnum {
		return fmt.Sprintf("stmt.setObject(%d, %s.value, %s)", idx, name, "Types.OTHER")
	}
	if t.IsArray {
		return fmt.Sprintf(`stmt.setArray(%d, conn.createArrayOf("%s", %s.toTypedArray()))`, idx, t.DataType, name)
	}
	if t.IsTime() {
		return fmt.Sprintf("stmt.setObject(%d, %s)", idx, name)
	}
	return fmt.Sprintf("stmt.set%s(%d, %s)", t.Name, idx, name)
}

type KtParams struct {
	Struct *KtStruct
}

func (v KtParams) isEmpty() bool {
	return len(v.Struct.Fields) == 0
}

func (v KtParams) Args() string {
	if v.isEmpty() {
		return ""
	}
	var out []string
	for _, f := range v.Struct.Fields {
		out = append(out, f.Name+": "+f.Type.String())
	}
	if len(out) < 3 {
		return strings.Join(out, ", ")
	}
	return "\n" + indent(strings.Join(out, ",\n"), 6, -1)
}

func (v KtParams) Bindings() string {
	if v.isEmpty() {
		return ""
	}
	var out []string
	for i, f := range v.Struct.JDBCParamBindings {
		out = append(out, jdbcSet(f.Type, i+1, f.Name))
	}
	return indent(strings.Join(out, "\n"), 10, 0)
}

func jdbcGet(t ktType, idx int) string {
	if t.IsEnum && t.IsArray {
		return fmt.Sprintf(`(results.getArray(%d).array as Array<String>).map { v -> %s.lookup(v)!! }.toList()`, idx, t.Name)
	}
	if t.IsEnum {
		return fmt.Sprintf("%s.lookup(results.getString(%d))!!", t.Name, idx)
	}
	if t.IsArray {
		return fmt.Sprintf(`(results.getArray(%d).array as Array<%s>).toList()`, idx, t.Name)
	}
	if t.IsTime() {
		return fmt.Sprintf(`results.getObject(%d, %s::class.java)`, idx, t.Name)
	}
	return fmt.Sprintf(`results.get%s(%d)`, t.Name, idx)
}

func (v KtQueryValue) ResultSet() string {
	var out []string
	if v.Struct == nil {
		return jdbcGet(v.Typ, 1)
	}
	for i, f := range v.Struct.Fields {
		out = append(out, jdbcGet(f.Type, i+1))
	}
	ret := indent(strings.Join(out, ",\n"), 4, -1)
	ret = indent(v.Struct.Name+"(\n"+ret+"\n)", 12, 0)
	return ret
}

func indent(s string, n int, firstIndent int) string {
	lines := strings.Split(s, "\n")
	buf := bytes.NewBuffer(nil)
	for i, l := range lines {
		indent := n
		if i == 0 && firstIndent != -1 {
			indent = firstIndent
		}
		if i != 0 {
			buf.WriteRune('\n')
		}
		for i := 0; i < indent; i++ {
			buf.WriteRune(' ')
		}
		buf.WriteString(l)
	}
	return buf.String()
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
	Arg          KtParams
}

type KtGenerateable interface {
	KtDataClasses(settings config.CombinedSettings) []KtStruct
	KtQueries(settings config.CombinedSettings) []KtQuery
	KtEnums(settings config.CombinedSettings) []KtEnum
}

func KtUsesType(r KtGenerateable, typ string, settings config.CombinedSettings) bool {
	for _, strct := range r.KtDataClasses(settings) {
		for _, f := range strct.Fields {
			if f.Type.Name == typ {
				return true
			}
		}
	}
	return false
}

func KtImports(r KtGenerateable, settings config.CombinedSettings) func(string) [][]string {
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

func InterfaceKtImports(r KtGenerateable, settings config.CombinedSettings) [][]string {
	kq := r.KtQueries(settings)
	uses := func(name string) bool {
		for _, q := range kq {
			if !q.Ret.isEmpty() {
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if strings.HasPrefix(f.Type.Name, name) {
						return true
					}
				}
			}
		}
		return false
	}

	std := stdImports(uses)
	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds, runtimeImports(kq)}
}

func ModelKtImports(r KtGenerateable, settings config.CombinedSettings) [][]string {
	std := make(map[string]struct{})
	if KtUsesType(r, "LocalDate", settings) {
		std["java.time.LocalDate"] = struct{}{}
	}
	if KtUsesType(r, "LocalTime", settings) {
		std["java.time.LocalTime"] = struct{}{}
	}
	if KtUsesType(r, "LocalDateTime", settings) {
		std["java.time.LocalDateTime"] = struct{}{}
	}
	if KtUsesType(r, "OffsetDateTime", settings) {
		std["java.time.OffsetDateTime"] = struct{}{}
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds}
}

func stdImports(uses func(name string) bool) map[string]struct{} {
	std := map[string]struct{}{
		"java.sql.SQLException": {},
	}
	if uses("LocalDate") {
		std["java.time.LocalDate"] = struct{}{}
	}
	if uses("LocalTime") {
		std["java.time.LocalTime"] = struct{}{}
	}
	if uses("LocalDateTime") {
		std["java.time.LocalDateTime"] = struct{}{}
	}
	if uses("OffsetDateTime") {
		std["java.time.OffsetDateTime"] = struct{}{}
	}
	return std
}

func runtimeImports(kq []KtQuery) []string {
	rt := map[string]struct{}{}
	for _, q := range kq {
		switch q.Cmd {
		case ":one":
			rt["sqlc.runtime.RowQuery"] = struct{}{}
		case ":many":
			rt["sqlc.runtime.ListQuery"] = struct{}{}
		case ":exec":
			rt["sqlc.runtime.ExecuteQuery"] = struct{}{}
		case ":execUpdate":
			rt["sqlc.runtime.ExecuteUpdateQuery"] = struct{}{}
		default:
			panic(fmt.Sprintf("invalid command %q", q.Cmd))
		}
	}
	rts := make([]string, 0, len(rt))
	for s, _ := range rt {
		rts = append(rts, s)
	}
	sort.Strings(rts)
	return rts
}

func QueryKtImports(r KtGenerateable, settings config.CombinedSettings, filename string) [][]string {
	kq := r.KtQueries(settings)
	uses := func(name string) bool {
		for _, q := range kq {
			if !q.Ret.isEmpty() {
				if q.Ret.Struct != nil {
					for _, f := range q.Ret.Struct.Fields {
						if f.Type.Name == name {
							return true
						}
					}
				}
				if q.Ret.Type() == name {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if f.Type.Name == name {
						return true
					}
				}
			}
		}
		return false
	}

	hasEnum := func() bool {
		for _, q := range kq {
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if f.Type.IsEnum {
						return true
					}
				}
			}
		}
		return false
	}

	std := stdImports(uses)
	std["java.sql.Connection"] = struct{}{}
	if hasEnum() {
		std["java.sql.Types"] = struct{}{}
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds, runtimeImports(kq)}
}

func ktEnumValueName(value string) string {
	id := strings.Replace(value, "-", "_", -1)
	id = strings.Replace(id, ":", "_", -1)
	id = strings.Replace(id, "/", "_", -1)
	id = ktIdentPattern.ReplaceAllString(id, "")
	return strings.ToUpper(id)
}

// Result is a wrapper around *dinosql.Result that extends it with Kotlin support.
// It can be used to generate both Go and Kotlin code.
// TODO: This is a temporary hack to ensure minimal chance of merge conflicts while Kotlin support is forked.
// Once it is merged upstream, we can factor split out Go support from the core dinosql.Result.
type Result struct {
	*dinosql.Result
}

func (r Result) KtEnums(settings config.CombinedSettings) []KtEnum {
	var enums []KtEnum
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

func KtDataClassName(name string, settings config.CombinedSettings) string {
	if rename := settings.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	for _, p := range strings.Split(name, "_") {
		out += strings.Title(p)
	}
	return out
}

func KtMemberName(name string, settings config.CombinedSettings) string {
	return dinosql.LowerTitle(KtDataClassName(name, settings))
}

func (r Result) KtDataClasses(settings config.CombinedSettings) []KtStruct {
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

type ktType struct {
	Name     string
	IsEnum   bool
	IsArray  bool
	IsNull   bool
	DataType string
}

func (t ktType) String() string {
	v := t.Name
	if t.IsArray {
		v = fmt.Sprintf("List<%s>", v)
	} else if t.IsNull {
		v += "?"
	}
	return v
}

func (t ktType) jdbcSetter() string {
	return "set" + t.jdbcType()
}

func (t ktType) jdbcType() string {
	if t.IsArray {
		return "Array"
	}
	if t.IsEnum || t.IsTime() {
		return "Object"
	}
	return t.Name
}

func (t ktType) IsTime() bool {
	return t.Name == "LocalDate" || t.Name == "LocalDateTime" || t.Name == "LocalTime" || t.Name == "OffsetDateTime"
}

func (r Result) ktType(col core.Column, settings config.CombinedSettings) ktType {
	typ, isEnum := r.ktInnerType(col, settings)
	return ktType{
		Name:     typ,
		IsEnum:   isEnum,
		IsArray:  col.IsArray,
		IsNull:   !col.NotNull,
		DataType: col.DataType,
	}
}

func (r Result) ktInnerType(col core.Column, settings config.CombinedSettings) (string, bool) {
	columnType := col.DataType

	switch columnType {
	case "serial", "pg_catalog.serial4":
		return "Int", false

	case "bigserial", "pg_catalog.serial8":
		return "Long", false

	case "smallserial", "pg_catalog.serial2":
		return "Short", false

	case "integer", "int", "int4", "pg_catalog.int4":
		return "Int", false

	case "bigint", "pg_catalog.int8":
		return "Long", false

	case "smallint", "pg_catalog.int2":
		return "Short", false

	case "float", "double precision", "pg_catalog.float8":
		return "Double", false

	case "real", "pg_catalog.float4":
		return "Float", false

	case "pg_catalog.numeric":
		return "java.math.BigDecimal", false

	case "bool", "pg_catalog.bool":
		return "Boolean", false

	case "jsonb":
		// TODO: support json and byte types
		return "String", false

	case "bytea", "blob", "pg_catalog.bytea":
		return "String", false

	case "date":
		// Date and time mappings from https://jdbc.postgresql.org/documentation/head/java8-date-time.html
		return "LocalDate", false

	case "pg_catalog.time", "pg_catalog.timetz":
		return "LocalTime", false

	case "pg_catalog.timestamp":
		return "LocalDateTime", false

	case "pg_catalog.timestamptz", "timestamptz":
		// TODO
		return "OffsetDateTime", false

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string":
		return "String", false

	case "uuid":
		// TODO
		return "uuid.UUID", false

	case "inet":
		// TODO
		return "net.IP", false

	case "void":
		// TODO
		// A void value always returns NULL. Since there is no built-in NULL
		// value into the SQL package, we'll use sql.NullBool
		return "sql.NullBool", false

	case "any":
		// TODO
		return "Any", false

	default:
		for name, schema := range r.Catalog.Schemas {
			if name == "pg_catalog" {
				continue
			}
			for _, enum := range schema.Enums() {
				if columnType == enum.Name {
					if name == "public" {
						return KtDataClassName(enum.Name, settings), true
					}

					return KtDataClassName(name+"_"+enum.Name, settings), true
				}
			}
		}
		log.Printf("unknown PostgreSQL type: %s\n", columnType)
		return "interface{}", false
	}
}

type goColumn struct {
	id int
	core.Column
}

func (r Result) ktColumnsToStruct(name string, columns []goColumn, settings config.CombinedSettings, namer func(core.Column, int) string) *KtStruct {
	gs := KtStruct{
		Name: name,
	}
	idSeen := map[int]KtField{}
	nameSeen := map[string]int{}
	for _, c := range columns {
		if binding, ok := idSeen[c.id]; ok {
			gs.JDBCParamBindings = append(gs.JDBCParamBindings, binding)
			continue
		}
		fieldName := KtMemberName(namer(c.Column, c.id), settings)
		if v := nameSeen[c.Name]; v > 0 {
			fieldName = fmt.Sprintf("%s_%d", fieldName, v+1)
		}
		field := KtField{
			Name: fieldName,
			Type: r.ktType(c.Column, settings),
		}
		gs.Fields = append(gs.Fields, field)
		gs.JDBCParamBindings = append(gs.JDBCParamBindings, field)
		nameSeen[c.Name]++
		idSeen[c.id] = field
	}
	return &gs
}

func ktArgName(name string) string {
	out := ""
	for i, p := range strings.Split(name, "_") {
		if i == 0 {
			out += strings.ToLower(p)
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

func ktParamName(c core.Column, number int) string {
	if c.Name != "" {
		return ktArgName(c.Name)
	}
	return fmt.Sprintf("dollar_%d", number)
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

func (r Result) KtQueries(settings config.CombinedSettings) []KtQuery {
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
			ConstantName: dinosql.LowerTitle(query.Name),
			FieldName:    dinosql.LowerTitle(query.Name) + "Stmt",
			MethodName:   dinosql.LowerTitle(query.Name),
			SourceName:   query.Filename,
			SQL:          jdbcSQL(query.SQL),
			Comments:     query.Comments,
		}

		var cols []goColumn
		for _, p := range query.Params {
			cols = append(cols, goColumn{
				id:     p.Number,
				Column: p.Column,
			})
		}
		params := r.ktColumnsToStruct(gq.ClassName+"Bindings", cols, settings, ktParamName)
		gq.Arg = KtParams{
			Struct: params,
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = KtQueryValue{
				Name: "results",
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
				var columns []goColumn
				for i, c := range query.Columns {
					columns = append(columns, goColumn{
						id:     i,
						Column: c,
					})
				}
				gs = r.ktColumnsToStruct(gq.ClassName+"Row", columns, settings, ktColumnName)
				emit = true
			}
			gq.Ret = KtQueryValue{
				Emit:   emit,
				Name:   "results",
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

interface Queries {
  {{- range .KtQueries}}
  @Throws(SQLException::class)
  {{- if eq .Cmd ":one"}}
  fun {{.MethodName}}({{.Arg.Args}}): RowQuery<{{.Ret.Type}}>
  {{- end}}
  {{- if eq .Cmd ":many"}}
  fun {{.MethodName}}({{.Arg.Args}}): ListQuery<{{.Ret.Type}}>
  {{- end}}
  {{- if eq .Cmd ":exec"}}
  fun {{.MethodName}}({{.Arg.Args}}): ExecuteQuery
  {{- end}}
  {{- if eq .Cmd ":execrows"}}
  fun {{.MethodName}}({{.Arg.Args}}): ExecuteUpdateQuery
  {{- end}}
  {{end}}
}
`

var ktModelsTmpl = `// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

{{range imports .SourceName}}
{{range .}}import {{.}}
{{end}}
{{end}}

{{range .Enums}}
{{if .Comment}}{{comment .Comment}}{{end}}
enum class {{.Name}}(val value: String) {
  {{- range $i, $e := .Constants}}
  {{- if $i }},{{end}}
  {{.Name}}("{{.Value}}")
  {{- end}};

  companion object {
    private val map = {{.Name}}.values().associateBy({{.Name}}::value)
    fun lookup(value: String) = map[value]
  }
}
{{end}}

{{range .KtDataClasses}}
{{if .Comment}}{{comment .Comment}}{{end}}
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

{{if .Ret.EmitStruct}}
data class {{.Ret.Type}} ( {{- range $i, $e := .Ret.Struct.Fields}}
  {{- if $i }},{{end}}
  val {{.Name}}: {{.Type}}
  {{- end}}
)
{{end}}
{{end}}

class QueriesImpl(private val conn: Connection) : Queries {
{{range .KtQueries}}
{{if eq .Cmd ":one"}}
{{range .Comments}}//{{.}}
{{end}}
  @Throws(SQLException::class)
  override fun {{.MethodName}}({{.Arg.Args}}): RowQuery<{{.Ret.Type}}> {
    return object : RowQuery<{{.Ret.Type}}>() {
      override fun execute(): {{.Ret.Type}} {
        return conn.prepareStatement({{.ConstantName}}).use { stmt ->
          this.statement = stmt
          {{.Arg.Bindings}}

          val results = stmt.executeQuery()
          if (!results.next()) {
            throw SQLException("no rows in result set")
          }
          val ret = {{.Ret.ResultSet}}
          if (results.next()) {
              throw SQLException("expected one row in result set, but got many")
          }
          ret
        }
      }
    }
  }
{{end}}

{{if eq .Cmd ":many"}}
{{range .Comments}}//{{.}}
{{end}}
  @Throws(SQLException::class)
  override fun {{.MethodName}}({{.Arg.Args}}): ListQuery<{{.Ret.Type}}> {
    return object : ListQuery<{{.Ret.Type}}>() {
      override fun execute(): List<{{.Ret.Type}}> {
        return conn.prepareStatement({{.ConstantName}}).use { stmt ->
          this.statement = stmt
          {{.Arg.Bindings}}

          val results = stmt.executeQuery()
          val ret = mutableListOf<{{.Ret.Type}}>()
          while (results.next()) {
              ret.add({{.Ret.ResultSet}})
          }
          ret
        }
      }
    }
  }
{{end}}

{{if eq .Cmd ":exec"}}
{{range .Comments}}//{{.}}
{{end}}
  @Throws(SQLException::class)
  {{ if $.EmitInterface }}override {{ end -}}
  override fun {{.MethodName}}({{.Arg.Args}}): ExecuteQuery {
    return object : ExecuteQuery() {
      override fun execute() {
        conn.prepareStatement({{.ConstantName}}).use { stmt ->
          this.statement = stmt
          {{ .Arg.Bindings }}

          stmt.execute()
        }
      }
    }
  }
{{end}}

{{if eq .Cmd ":execrows"}}
{{range .Comments}}//{{.}}
{{end}}
  @Throws(SQLException::class)
  {{ if $.EmitInterface }}override {{ end -}}
  override fun {{.MethodName}}({{.Arg.Args}}): ExecuteUpdateQuery {
    return object : ExecUpdateQuery() {
      override fun execute(): Int {
        return conn.prepareStatement({{.ConstantName}}).use { stmt ->
          this.statement = stmt
          {{ .Arg.Bindings }}

          stmt.execute()
          stmt.updateCount
        }
      }
    }
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
	Settings      config.Config

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

func KtGenerate(r KtGenerateable, settings config.CombinedSettings) (map[string]string, error) {
	funcMap := template.FuncMap{
		"lowerTitle": dinosql.LowerTitle,
		"comment":    dinosql.DoubleSlashComment,
		"imports":    KtImports(r, settings),
		"offset":     Offset,
	}

	modelsFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktModelsTmpl))
	sqlFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktSqlTmpl))
	ifaceFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktIfaceTmpl))

	pkg := settings.Package
	tctx := ktTmplCtx{
		Settings:      settings.Global,
		Q:             `"""`,
		Package:       pkg.Gen.Kotlin.Package,
		KtQueries:     r.KtQueries(settings),
		Enums:         r.KtEnums(settings),
		KtDataClasses: r.KtDataClasses(settings),
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
	if err := execute("Queries.kt", ifaceFile); err != nil {
		return nil, err
	}
	if err := execute("QueriesImpl.kt", sqlFile); err != nil {
		return nil, err
	}

	return output, nil
}
