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

func (v QueryValue) StructRowParser(rowVar string, indentCount int) string {
	if !v.IsStruct() {
		panic("StructRowParse called on non-struct QueryValue")
	}
	indent := strings.Repeat(" ", indentCount+4)
	params := make([]string, 0, len(v.Struct.Fields))
	for i, f := range v.Struct.Fields {
		params = append(params, fmt.Sprintf("%s%s=%s[%v],", indent, f.Name, rowVar, i))
	}
	indent = strings.Repeat(" ", indentCount)
	return v.Type() + "(\n" + strings.Join(params, "\n") + "\n" + indent + ")"
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
	// A single struct arg does not need to be passed as a keyword argument
	if len(q.Args) == 1 && q.Args[0].IsStruct() {
		return ", " + q.Args[0].Pair()
	}
	argPairs := make([]string, 0, len(q.Args))
	for _, a := range q.Args {
		argPairs = append(argPairs, a.Pair())
	}
	if len(argPairs) == 0 {
		return ""
	}
	return ", *, " + strings.Join(argPairs, ", ")
}

func (q Query) ArgDict() string {
	params := make([]string, 0, len(q.Args))
	i := 1
	for _, a := range q.Args {
		if a.isEmpty() {
			continue
		}
		if a.IsStruct() {
			for _, f := range a.Struct.Fields {
				params = append(params, fmt.Sprintf("\"p%v\": %s", i, a.Name+"."+f.Name))
				i++
			}
		} else {
			params = append(params, fmt.Sprintf("\"p%v\": %s", i, a.Name))
			i++
		}
	}
	if len(params) == 0 {
		return ""
	}
	if len(params) < 4 {
		return ", {" + strings.Join(params, ", ") + "}"
	}
	return ", {\n            " + strings.Join(params, ",\n            ") + ",\n        }"
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

var postgresPlaceholderRegexp = regexp.MustCompile(`\B\$(\d+)\b`)

// Sqlalchemy uses ":name" for placeholders, so "$N" is converted to ":pN"
// This also means ":" has special meaning to sqlalchemy, so it must be escaped.
func sqlalchemySQL(s string, engine config.Engine) string {
	s = strings.ReplaceAll(s, ":", `\\:`)
	if engine == config.EnginePostgreSQL {
		return postgresPlaceholderRegexp.ReplaceAllString(s, ":p$1")
	}
	return s
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
			SQL:          sqlalchemySQL(query.SQL, settings.Package.Engine),
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
					// HACK: models do not have "models." on their types, so trim that so we can find matches
					trimmedPyType := makePyType(r, c, settings)
					trimmedPyType.InnerType = strings.TrimPrefix(trimmedPyType.InnerType, "models.")
					sameName := f.Name == columnName(c, i)
					sameType := f.Type == trimmedPyType
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


{{range .Enums}}
{{- if .Comment}}{{comment .Comment}}{{- end}}
class {{.Name}}(str, enum.Enum):
    {{- range .Constants}}
    {{.Name}} = "{{.Value}}"
    {{- end}}
{{end}}

{{- range .Models}}
{{if .Comment}}{{comment .Comment}}{{- end}}
@dataclasses.dataclass()
class {{.Name}}: {{- range .Fields}}
    {{- if .Comment}}
    {{comment .Comment}}{{else}}
    {{- end}}
    {{.Name}}: {{.Type}}
    {{- end}}

{{end}}
`

var queriesTmpl = `
{{- define "dataclassParse"}}

{{end}}
# Code generated by sqlc. DO NOT EDIT.
{{- range imports .SourceName}}
{{.}}
{{- end}}

{{range .Queries}}
{{- if $.OutputQuery .SourceName}}
{{.ConstantName}} = """-- name: {{.MethodName}} \\{{.Cmd}}
{{.SQL}}
"""
{{range .Args}}
{{- if .EmitStruct}}

@dataclasses.dataclass()
class {{.Type}}: {{- range .Struct.Fields}}
    {{.Name}}: {{.Type}}
    {{- end}}
{{end}}{{end}}
{{- if .Ret.EmitStruct}}

@dataclasses.dataclass()
class {{.Ret.Type}}: {{- range .Ret.Struct.Fields}}
    {{.Name}}: {{.Type}}
    {{- end}}
{{end}}
{{end}}
{{- end}}

{{- if .EmitSync}}
class Querier:
    def __init__(self, conn: sqlalchemy.engine.Connection):
        self._conn = conn
{{range .Queries}}
{{- if $.OutputQuery .SourceName}}
{{- if eq .Cmd ":one"}}
    def {{.MethodName}}(self{{.ArgPairs}}) -> Optional[{{.Ret.Type}}]:
        row = self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}}).first()
        if row is None:
            return None
        {{- if .Ret.IsStruct}}
        return {{.Ret.StructRowParser "row" 8}}
        {{- else}}
        return row[0]
        {{- end}}
{{end}}

{{- if eq .Cmd ":many"}}
    def {{.MethodName}}(self{{.ArgPairs}}) -> Iterator[{{.Ret.Type}}]:
        result = self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})
        for row in result:
            {{- if .Ret.IsStruct}}
            yield {{.Ret.StructRowParser "row" 12}}
            {{- else}}
            yield row[0]
            {{- end}}
{{end}}

{{- if eq .Cmd ":exec"}}
    def {{.MethodName}}(self{{.ArgPairs}}) -> None:
        self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})
{{end}}

{{- if eq .Cmd ":execrows"}}
    def {{.MethodName}}(self{{.ArgPairs}}) -> int:
        result = self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})
        return result.rowcount
{{end}}

{{- if eq .Cmd ":execresult"}}
    def {{.MethodName}}(self{{.ArgPairs}}) -> sqlalchemy.engine.Result:
        return self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})
{{end}}
{{- end}}
{{- end}}
{{- end}}

{{- if .EmitAsync}}

class AsyncQuerier:
    def __init__(self, conn: sqlalchemy.ext.asyncio.AsyncConnection):
        self._conn = conn
{{range .Queries}}
{{- if $.OutputQuery .SourceName}}
{{- if eq .Cmd ":one"}}
    async def {{.MethodName}}(self{{.ArgPairs}}) -> Optional[{{.Ret.Type}}]:
        row = (await self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})).first()
        if row is None:
            return None
        {{- if .Ret.IsStruct}}
        return {{.Ret.StructRowParser "row" 8}}
        {{- else}}
        return row[0]
        {{- end}}
{{end}}

{{- if eq .Cmd ":many"}}
    async def {{.MethodName}}(self{{.ArgPairs}}) -> AsyncIterator[{{.Ret.Type}}]:
        result = await self._conn.stream(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})
        async for row in result:
            {{- if .Ret.IsStruct}}
            yield {{.Ret.StructRowParser "row" 12}}
            {{- else}}
            yield row[0]
            {{- end}}
{{end}}

{{- if eq .Cmd ":exec"}}
    async def {{.MethodName}}(self{{.ArgPairs}}) -> None:
        await self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})
{{end}}

{{- if eq .Cmd ":execrows"}}
    async def {{.MethodName}}(self{{.ArgPairs}}) -> int:
        result = await self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})
        return result.rowcount
{{end}}

{{- if eq .Cmd ":execresult"}}
    async def {{.MethodName}}(self{{.ArgPairs}}) -> sqlalchemy.engine.Result:
        return await self._conn.execute(sqlalchemy.text({{.ConstantName}}){{.ArgDict}})
{{end}}
{{- end}}
{{- end}}
{{- end}}
`

type pyTmplCtx struct {
	Models     []Struct
	Queries    []Query
	Enums      []Enum
	EmitSync   bool
	EmitAsync  bool
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
		Models:    models,
		Queries:   queries,
		Enums:     enums,
		EmitSync:  settings.Python.EmitSyncQuerier,
		EmitAsync: settings.Python.EmitAsyncQuerier,
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
