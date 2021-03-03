package python

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kyleconroy/sqlc/internal/codegen"
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/core"
	"github.com/kyleconroy/sqlc/internal/inflection"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"log"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type Constant struct {
	Name  string
	Type  string
	Value string
}

type Enum struct {
	Name      string
	Comment   string
	Constants []Constant
}

type pyType struct {
	InnerType string
	IsArray   bool
	IsNull    bool
}

func (t pyType) String() string {
	v := t.InnerType
	if t.IsArray {
		v = fmt.Sprintf("List[%s]", v)
	}
	if t.IsNull {
		v = fmt.Sprintf("Optional[%s]", v)
	}
	return v
}

type Field struct {
	Name    string
	Type    pyType
	Comment string
}

type Struct struct {
	Table   core.FQN
	Name    string
	Fields  []Field
	Comment string
}

func (s Struct) DedupFields() []Field {
	seen := map[string]struct{}{}
	dedupFields := make([]Field, 0)
	for _, f := range s.Fields {
		if _, ok := seen[f.Name]; ok {
			continue
		}
		seen[f.Name] = struct{}{}
		dedupFields = append(dedupFields, f)
	}
	return dedupFields
}

type QueryValue struct {
	Emit   bool
	Name   string
	Struct *Struct
	Typ    pyType
}

func (v QueryValue) EmitStruct() bool {
	return v.Emit
}

func (v QueryValue) IsStruct() bool {
	return v.Struct != nil
}

func (v QueryValue) isEmpty() bool {
	return v.Typ == (pyType{}) && v.Name == "" && v.Struct == nil
}

func (v QueryValue) Pair() string {
	if v.isEmpty() {
		return ""
	}
	return v.Name + ": " + v.Type()
}

func (v QueryValue) Type() string {
	if v.Typ != (pyType{}) {
		return v.Typ.String()
	}
	if v.Struct != nil {
		if v.Emit {
			return v.Struct.Name
		} else {
			return "models." + v.Struct.Name
		}
	}
	panic("no type for QueryValue: " + v.Name)
}

// A struct used to generate methods and fields on the Queries struct
type Query struct {
	Cmd          string
	Comments     []string
	MethodName   string
	FieldName    string
	ConstantName string
	SQL          string
	SourceName   string
	Ret          QueryValue
	Args         []QueryValue
}

func (q Query) ArgPairs() string {
	argPairs := make([]string, 0, len(q.Args))
	for _, a := range q.Args {
		argPairs = append(argPairs, a.Pair())
	}
	if len(argPairs) == 0 {
		return ""
	}
	return ", " + strings.Join(argPairs, ", ")
}

func (q Query) ArgParams() string {
	params := make([]string, 0, len(q.Args))
	for _, a := range q.Args {
		if a.isEmpty() {
			continue
		}
		if a.IsStruct() {
			for _, f := range a.Struct.Fields {
				params = append(params, a.Name+"."+f.Name)
			}
		} else {
			params = append(params, a.Name)
		}
	}
	if len(params) == 0 {
		return ""
	}
	return ", " + strings.Join(params, ", ")
}

func makePyType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) pyType {
	typ := pyInnerType(r, col, settings)
	return pyType{
		InnerType: typ,
		IsArray:   col.IsArray,
		IsNull:    !col.NotNull,
	}
}

func pyInnerType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	for _, oride := range settings.Overrides {
		if !oride.PythonType.IsSet() {
			continue
		}
		sameTable := sameTableName(col.Table, oride.Table, r.Catalog.DefaultSchema)
		if oride.Column != "" && oride.ColumnName == col.Name && sameTable {
			return oride.PythonType.TypeString()
		}
		if oride.DBType != "" && oride.DBType == col.DataType && oride.Nullable != (col.NotNull || col.IsArray) {
			return oride.PythonType.TypeString()
		}
	}

	switch settings.Package.Engine {
	case config.EnginePostgreSQL:
		return postgresType(r, col, settings)
	default:
		log.Println("unsupported engine type")
		return "Any"
	}
}

func ModelName(name string, settings config.CombinedSettings) string {
	if rename := settings.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	for _, p := range strings.Split(name, "_") {
		out += strings.Title(p)
	}
	return out
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func MethodName(name string) string {
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var pyIdentPattern = regexp.MustCompile("[^a-zA-Z0-9_]+")

func pyEnumValueName(value string) string {
	id := strings.Replace(value, "-", "_", -1)
	id = strings.Replace(id, ":", "_", -1)
	id = strings.Replace(id, "/", "_", -1)
	id = pyIdentPattern.ReplaceAllString(id, "")
	return strings.ToUpper(id)
}

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
				Name:    ModelName(enumName, settings),
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

func buildModels(r *compiler.Result, settings config.CombinedSettings) []Struct {
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
			if !settings.Python.EmitExactTableNames {
				structName = inflection.Singular(structName)
			}
			s := Struct{
				Table:   core.FQN{Schema: schema.Name, Rel: table.Rel.Name},
				Name:    ModelName(structName, settings),
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				typ := makePyType(r, compiler.ConvertColumn(table.Rel, column), settings)
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

func columnName(c *compiler.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}

func paramName(p compiler.Parameter) string {
	if p.Column.Name != "" {
		return p.Column.Name
	}
	return fmt.Sprintf("dollar_%d", p.Number)
}

type pyColumn struct {
	id int
	*compiler.Column
}

func columnsToStruct(r *compiler.Result, name string, columns []pyColumn, settings config.CombinedSettings) *Struct {
	gs := Struct{
		Name: name,
	}
	seen := map[string]int{}
	suffixes := map[int]int{}
	for i, c := range columns {
		colName := columnName(c.Column, i)
		fieldName := colName
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
			fieldName = fmt.Sprintf("%s_%d", fieldName, suffix)
		}
		gs.Fields = append(gs.Fields, Field{
			Name: fieldName,
			Type: makePyType(r, c.Column, settings),
		})
		seen[colName]++
	}
	return &gs
}

func sameTableName(n *ast.TableName, f core.FQN, defaultSchema string) bool {
	if n == nil {
		return false
	}
	schema := n.Schema
	if n.Schema == "" {
		schema = defaultSchema
	}
	return n.Catalog == f.Catalog && schema == f.Schema && n.Name == f.Rel
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

		methodName := MethodName(query.Name)

		gq := Query{
			Cmd:          query.Cmd,
			Comments:     query.Comments,
			MethodName:   methodName,
			FieldName:    codegen.LowerTitle(query.Name) + "Stmt",
			ConstantName: strings.ToUpper(methodName),
			SQL:          query.SQL,
			SourceName:   query.Filename,
		}

		if len(query.Params) > 4 {
			var cols []pyColumn
			for _, p := range query.Params {
				cols = append(cols, pyColumn{
					id:     p.Number,
					Column: p.Column,
				})
			}
			gq.Args = []QueryValue{{
				Emit:   true,
				Name:   "arg",
				Struct: columnsToStruct(r, query.Name+"Params", cols, settings),
			}}
		} else {
			args := make([]QueryValue, 0, len(query.Params))
			for _, p := range query.Params {
				args = append(args, QueryValue{
					Name: paramName(p),
					Typ:  makePyType(r, p.Column, settings),
				})
			}
			gq.Args = args
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = QueryValue{
				Name: columnName(c, 0),
				Typ:  makePyType(r, c, settings),
			}
		} else if len(query.Columns) > 1 {
			var gs *Struct
			var emit bool

			for _, s := range structs {
				if len(s.Fields) != len(query.Columns) {
					continue
				}
				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					sameName := f.Name == columnName(c, i)
					sameType := f.Type == makePyType(r, c, settings)
					sameTable := sameTableName(c.Table, s.Table, r.Catalog.DefaultSchema)
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
				var columns []pyColumn
				for i, c := range query.Columns {
					columns = append(columns, pyColumn{
						id:     i,
						Column: c,
					})
				}
				gs = columnsToStruct(r, query.Name+"Row", columns, settings)
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

var modelsTmpl = `# Code generated by sqlc. DO NOT EDIT.
{{- range imports .SourceName}}
{{.}}
{{- end}}


# Enums
{{- range .Enums}}
{{- if .Comment}}{{comment .Comment}}{{- end}}
class {{.Name}}(str, enum.Enum):
    {{- range .Constants}}
    {{.Name}} = "{{.Value}}"
    {{- end}}
{{end}}

# Models
{{- range .Models}}
{{- if .Comment}}{{comment .Comment}}{{- end}}
class {{.Name}}(pydantic.BaseModel): {{- range .DedupFields}}
    {{- if .Comment}}
    {{comment .Comment}}{{else}}
    {{- end}}
    {{.Name}}: {{.Type}}
    {{- end}}

{{end}}
`

var queriesTmpl = `# Code generated by sqlc. DO NOT EDIT.
{{- range imports .SourceName}}
{{.}}
{{- end}}

{{range .Queries}}
{{- if $.OutputQuery .SourceName}}
{{.ConstantName}} = """-- name: {{.MethodName}} {{.Cmd}}
{{.SQL}}
"""
{{range .Args}}
{{- if .EmitStruct}}

class {{.Type}}(pydantic.BaseModel): {{- range .Struct.DedupFields}}
    {{.Name}}: {{.Type}}
    {{- end}}
{{end}}{{end}}
{{- if .Ret.EmitStruct}}

class {{.Ret.Type}}(pydantic.BaseModel): {{- range .Ret.Struct.DedupFields}}
    {{.Name}}: {{.Type}}
    {{- end}}
{{end}}
{{end}}
{{- end}}

{{- range .Queries}}
{{- if $.OutputQuery .SourceName}}
{{- if eq .Cmd ":one"}}
@overload
def {{.MethodName}}(conn: sqlc.Connection{{.ArgPairs}}) -> Optional[{{.Ret.Type}}]:
    pass


@overload
def {{.MethodName}}(conn: sqlc.AsyncConnection{{.ArgPairs}}) -> Awaitable[Optional[{{.Ret.Type}}]]:
    pass


def {{.MethodName}}(conn: sqlc.GenericConnection{{.ArgPairs}}) -> sqlc.ReturnType[Optional[{{.Ret.Type}}]]:
    {{- if .Ret.IsStruct}}
    return conn.execute_one_model({{.Ret.Type}}, {{.ConstantName}}{{.ArgParams}})
    {{- else}}
    return conn.execute_one({{.ConstantName}}{{.ArgParams}})
    {{- end}}
{{end}}

{{- if eq .Cmd ":many"}}
@overload
def {{.MethodName}}(conn: sqlc.Connection{{.ArgPairs}}) -> Iterator[{{.Ret.Type}}]:
    pass


@overload
def {{.MethodName}}(conn: sqlc.AsyncConnection{{.ArgPairs}}) -> AsyncIterator[{{.Ret.Type}}]:
    pass


def {{.MethodName}}(conn: sqlc.GenericConnection{{.ArgPairs}}) -> sqlc.IteratorReturn[{{.Ret.Type}}]:
    {{- if .Ret.IsStruct}}
    return conn.execute_many_model({{.Ret.Type}}, {{.ConstantName}}{{.ArgParams}})
    {{- else}}
    return conn.execute_many({{.ConstantName}}{{.ArgParams}})
    {{- end}}
{{end}}

{{- if eq .Cmd ":exec"}}
@overload
def {{.MethodName}}(conn: sqlc.Connection{{.ArgPairs}}) -> None:
    pass


@overload
def {{.MethodName}}(conn: sqlc.AsyncConnection{{.ArgPairs}}) -> Awaitable[None]:
    pass


def {{.MethodName}}(conn: sqlc.GenericConnection{{.ArgPairs}}) -> sqlc.ReturnType[None]:
    return conn.execute_none({{.ConstantName}}{{.ArgParams}})
{{end}}

{{- if eq .Cmd ":execrows"}}
@overload
def {{.MethodName}}(conn: sqlc.Connection{{.ArgPairs}}) -> int:
    pass


@overload
def {{.MethodName}}(conn: sqlc.AsyncConnection{{.ArgPairs}}) -> Awaitable[int]:
    pass


def {{.MethodName}}(conn: sqlc.GenericConnection{{.ArgPairs}}) -> sqlc.ReturnType[int]:
    return conn.execute_rowcount({{.ConstantName}}{{.ArgParams}})
{{end}}

{{- if eq .Cmd ":execresult"}}
@overload
def {{.MethodName}}(conn: sqlc.Connection{{.ArgPairs}}) -> sqlc.Cursor:
    pass


@overload
def {{.MethodName}}(conn: sqlc.AsyncConnection{{.ArgPairs}}) -> sqlc.AsyncCursor:
    pass


def {{.MethodName}}(conn: sqlc.GenericConnection{{.ArgPairs}}) -> sqlc.GenericCursor:
    return conn.execute({{.ConstantName}}{{.ArgParams}})
{{end}}
{{end}}
{{- end}}
`

type pyTmplCtx struct {
	Models     []Struct
	Queries    []Query
	Enums      []Enum
	SourceName string
}

func (t *pyTmplCtx) OutputQuery(sourceName string) bool {
	return t.SourceName == sourceName
}

func HashComment(s string) string {
	return "# " + strings.ReplaceAll(s, "\n", "\n# ")
}

func Generate(r *compiler.Result, settings config.CombinedSettings) (map[string]string, error) {
	enums := buildEnums(r, settings)
	models := buildModels(r, settings)
	queries := buildQueries(r, settings, models)

	i := &importer{
		Settings: settings,
		Models:   models,
		Queries:  queries,
		Enums:    enums,
	}

	funcMap := template.FuncMap{
		"lowerTitle": codegen.LowerTitle,
		"comment":    HashComment,
		"imports":    i.Imports,
	}

	modelsFile := template.Must(template.New("table").Funcs(funcMap).Parse(modelsTmpl))
	queriesFile := template.Must(template.New("table").Funcs(funcMap).Parse(queriesTmpl))

	tctx := pyTmplCtx{
		Models:  models,
		Queries: queries,
		Enums:   enums,
	}

	output := map[string]string{}

	execute := func(name string, t *template.Template) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := t.Execute(w, &tctx)
		w.Flush()
		if err != nil {
			return err
		}
		if !strings.HasSuffix(name, ".py") {
			name = strings.TrimSuffix(name, ".sql")
			name += ".py"
		}
		output[name] = b.String()
		return nil
	}

	if err := execute("models.py", modelsFile); err != nil {
		return nil, err
	}

	files := map[string]struct{}{}
	for _, q := range queries {
		files[q.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, queriesFile); err != nil {
			return nil, err
		}
	}

	return output, nil
}
