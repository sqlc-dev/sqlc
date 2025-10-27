package golang

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type QueryValue struct {
	Emit        bool
	EmitPointer bool
	Name        string
	DBName      string // The name of the field in the database. Only set if Struct==nil.
	Struct      *Struct
	Typ         string
	SQLDriver   opts.SQLDriver

	// Column is kept so late in the generation process around to differentiate
	// between mysql slices and pg arrays
	Column *plugin.Column
}

func (v QueryValue) EmitStruct() bool {
	return v.Emit
}

func (v QueryValue) IsStruct() bool {
	return v.Struct != nil
}

func (v QueryValue) IsPointer() bool {
	return v.EmitPointer && v.Struct != nil
}

func (v QueryValue) isEmpty() bool {
	return v.Typ == "" && v.Name == "" && v.Struct == nil
}

func (v QueryValue) IsEmpty() bool {
	return v.isEmpty()
}

type Argument struct {
	Name string
	Type string
}

func (v QueryValue) Pair() string {
	var out []string
	for _, arg := range v.Pairs() {
		out = append(out, arg.Name+" "+arg.Type)
	}
	return strings.Join(out, ",")
}

// Return the argument name and type for query methods. Should only be used in
// the context of method arguments.
func (v QueryValue) Pairs() []Argument {
	if v.isEmpty() {
		return nil
	}
	if !v.EmitStruct() && v.IsStruct() {
		var out []Argument
		for _, f := range v.Struct.Fields {
			out = append(out, Argument{
				Name: escape(toLowerCase(f.Name)),
				Type: f.Type,
			})
		}
		return out
	}
	return []Argument{
		{
			Name: escape(v.Name),
			Type: v.DefineType(),
		},
	}
}

func (v QueryValue) SlicePair() string {
	if v.isEmpty() {
		return ""
	}
	return v.Name + " []" + v.DefineType()
}

func (v QueryValue) Type() string {
	if v.Typ != "" {
		return v.Typ
	}
	if v.Struct != nil {
		return v.Struct.Name
	}
	panic("no type for QueryValue: " + v.Name)
}

func (v *QueryValue) DefineType() string {
	t := v.Type()
	if v.IsPointer() {
		return "*" + t
	}
	return t
}

func (v *QueryValue) ReturnName() string {
	if v.IsPointer() {
		return "&" + escape(v.Name)
	}
	return escape(v.Name)
}

func (v QueryValue) UniqueFields() []Field {
	seen := map[string]struct{}{}
	fields := make([]Field, 0, len(v.Struct.Fields))

	for _, field := range v.Struct.Fields {
		if _, found := seen[field.Name]; found {
			continue
		}
		seen[field.Name] = struct{}{}
		fields = append(fields, field)
	}

	return fields
}

// YDBUniqueFieldsWithComments returns unique fields for YDB struct generation with comments for interface{} fields
func (v QueryValue) YDBUniqueFieldsWithComments() []Field {
	seen := map[string]struct{}{}
	fields := make([]Field, 0, len(v.Struct.Fields))

	for _, field := range v.Struct.Fields {
		if _, found := seen[field.Name]; found {
			continue
		}
		seen[field.Name] = struct{}{}

		if strings.HasSuffix(field.Type, "interface{}") {
			field.Comment = "// sqlc couldn't resolve type, pass via params"
		}

		fields = append(fields, field)
	}

	return fields
}

func (v QueryValue) Params() string {
	if v.isEmpty() {
		return ""
	}
	var out []string
	if v.Struct == nil {
		if !v.Column.IsSqlcSlice && strings.HasPrefix(v.Typ, "[]") && v.Typ != "[]byte" && !v.SQLDriver.IsPGX() {
			out = append(out, "pq.Array("+escape(v.Name)+")")
		} else {
			out = append(out, escape(v.Name))
		}
	} else {
		for _, f := range v.Struct.Fields {
			if !f.HasSqlcSlice() && strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" && !v.SQLDriver.IsPGX() {
				out = append(out, "pq.Array("+escape(v.VariableForField(f))+")")
			} else {
				out = append(out, escape(v.VariableForField(f)))
			}
		}
	}
	if len(out) <= 3 {
		return strings.Join(out, ",")
	}
	out = append(out, "")
	return "\n" + strings.Join(out, ",\n")
}

func (v QueryValue) ColumnNames() []string {
	if v.Struct == nil {
		return []string{v.DBName}
	}
	names := make([]string, len(v.Struct.Fields))
	for i, f := range v.Struct.Fields {
		names[i] = f.DBName
	}
	return names
}

func (v QueryValue) ColumnNamesAsGoSlice() string {
	if v.Struct == nil {
		return fmt.Sprintf("[]string{%q}", v.DBName)
	}
	escapedNames := make([]string, len(v.Struct.Fields))
	for i, f := range v.Struct.Fields {
		if f.Column != nil && f.Column.OriginalName != "" {
			escapedNames[i] = fmt.Sprintf("%q", f.Column.OriginalName)
		} else {
			escapedNames[i] = fmt.Sprintf("%q", f.DBName)
		}
	}
	return "[]string{" + strings.Join(escapedNames, ", ") + "}"
}

// When true, we have to build the arguments to q.db.QueryContext in addition to
// munging the SQL
func (v QueryValue) HasSqlcSlices() bool {
	if v.Struct == nil {
		return v.Column != nil && v.Column.IsSqlcSlice
	}
	for _, v := range v.Struct.Fields {
		if v.Column.IsSqlcSlice {
			return true
		}
	}
	return false
}

func (v QueryValue) Scan() string {
	var out []string
	if v.Struct == nil {
		if strings.HasPrefix(v.Typ, "[]") && v.Typ != "[]byte" && !v.SQLDriver.IsPGX() && !v.SQLDriver.IsYDBGoSDK() {
			out = append(out, "pq.Array(&"+v.Name+")")
		} else {
			out = append(out, "&"+v.Name)
		}
	} else {
		for _, f := range v.Struct.Fields {

			// append any embedded fields
			if len(f.EmbedFields) > 0 {
				for _, embed := range f.EmbedFields {
					if strings.HasPrefix(embed.Type, "[]") && embed.Type != "[]byte" && !v.SQLDriver.IsPGX() && !v.SQLDriver.IsYDBGoSDK() {
						out = append(out, "pq.Array(&"+v.Name+"."+f.Name+"."+embed.Name+")")
					} else {
						out = append(out, "&"+v.Name+"."+f.Name+"."+embed.Name)
					}
				}
				continue
			}

			if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" && !v.SQLDriver.IsPGX() && !v.SQLDriver.IsYDBGoSDK() {
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

// Deprecated: This method does not respect the Emit field set on the
// QueryValue. It's used by the go-sql-driver-mysql/copyfromCopy.tmpl and should
// not be used other places.
func (v QueryValue) CopyFromMySQLFields() []Field {
	// fmt.Printf("%#v\n", v)
	if v.Struct != nil {
		return v.Struct.Fields
	}
	return []Field{
		{
			Name:   v.Name,
			DBName: v.DBName,
			Type:   v.Typ,
		},
	}
}

func (v QueryValue) VariableForField(f Field) string {
	if !v.IsStruct() {
		return v.Name
	}
	if !v.EmitStruct() {
		return toLowerCase(f.Name)
	}
	return v.Name + "." + f.Name
}

func addDollarPrefix(name string) string {
	if name == "" {
		return name
	}
	if strings.HasPrefix(name, "$") {
		return name
	}
	return "$" + name
}

// ydbBuilderMethodForColumnType maps a YDB column data type to a ParamsBuilder method name.
func ydbBuilderMethodForColumnType(dbType string) string {
	baseType := extractBaseType(strings.ToLower(dbType))

	switch baseType {
	case "bool":
		return "Bool"
	case "uint64":
		return "Uint64"
	case "int64", "bigserial", "serial8":
		return "Int64"
	case "uint32":
		return "Uint32"
	case "int32", "serial", "serial4":
		return "Int32"
	case "uint16":
		return "Uint16"
	case "int16", "smallserial", "serial2":
		return "Int16"
	case "uint8":
		return "Uint8"
	case "int8":
		return "Int8"
	case "float":
		return "Float"
	case "double":
		return "Double"
	case "json":
		return "JSON"
	case "jsondocument":
		return "JSONDocument"
	case "utf8", "text":
		return "Text"
	case "string":
		return "Bytes"
	case "date":
		return "Date"
	case "date32":
		return "Date32"
	case "datetime":
		return "Datetime"
	case "timestamp":
		return "Timestamp"
	case "tzdate":
		return "TzDate"
	case "tzdatetime":
		return "TzDatetime"
	case "tztimestamp":
		return "TzTimestamp"
	case "uuid":
		return "Uuid"
	case "yson":
		return "YSON"

	case "integer": // LIMIT/OFFSET parameters support
		return "Uint64"

	//TODO: support other types
	default:
		// Check for decimal types
		if strings.HasPrefix(baseType, "decimal") {
			return "Decimal"
		}
		return ""
	}
}

// ydbIterateNamedParams iterates over named parameters and calls the provided function for each one.
// The function receives the field and method name, and should return true to continue iteration.
func (v QueryValue) ydbIterateNamedParams(fn func(field Field, method string) bool) bool {
	if v.isEmpty() {
		return false
	}

	for _, field := range v.getParameterFields() {
		if field.Column != nil {
			name := field.Column.GetName()
			if name == "" {
				continue
			}

			var method string
			if field.Column != nil && field.Column.Type != nil {
				method = ydbBuilderMethodForColumnType(sdk.DataType(field.Column.Type))
			}

			if !fn(field, method) {
				return false
			}
		}
	}
	return true
}

// YDBHasComplexContainers returns true if there are complex container types that sqlc cannot handle automatically.
func (v QueryValue) YDBHasComplexContainers() bool {
	hasComplex := false
	v.ydbIterateNamedParams(func(field Field, method string) bool {
		if method == "" {
			hasComplex = true
			return false
		}
		return true
	})
	return hasComplex
}

// YDBParamsBuilder emits Go code that constructs YDB params using ParamsBuilder.
func (v QueryValue) YDBParamsBuilder() string {
	var lines []string

	v.ydbIterateNamedParams(func(field Field, method string) bool {
		if method == "" {
			return true
		}

		name := field.Column.GetName()
		paramName := fmt.Sprintf("%q", addDollarPrefix(name))
		variable := escape(v.VariableForField(field))

		goType := field.Type
		isPtr := strings.HasPrefix(goType, "*")
		isArray := field.Column.IsArray || field.Column.IsSqlcSlice

		if method == "Decimal" {
			if isArray {
				lines = append(lines, fmt.Sprintf("\tvar list = parameters.Param(%s).BeginList()", paramName))
				lines = append(lines, fmt.Sprintf("\tfor _, param := range %s {", variable))
				lines = append(lines, "\t\tlist = list.Add().Decimal(param.Bytes, param.Precision, param.Scale)")
				lines = append(lines, "\t}")
				lines = append(lines, "\tparameters = list.EndList()")
			} else if isPtr {
				lines = append(lines, fmt.Sprintf("\tparameters = parameters.Param(%s).BeginOptional().Decimal(&%s.Bytes, %s.Precision, %s.Scale).EndOptional()", paramName, variable, variable, variable))
			} else {
				lines = append(lines, fmt.Sprintf("\tparameters = parameters.Param(%s).Decimal(%s.Bytes, %s.Precision, %s.Scale)", paramName, variable, variable, variable))
			}
			return true
		}

		if isArray {
			lines = append(lines, fmt.Sprintf("\tvar list = parameters.Param(%s).BeginList()", paramName))
			lines = append(lines, fmt.Sprintf("\tfor _, param := range %s {", variable))
			lines = append(lines, fmt.Sprintf("\t\tlist = list.Add().%s(param)", method))
			lines = append(lines, "\t}")
			lines = append(lines, "\tparameters = list.EndList()")
		} else if isPtr {
			lines = append(lines, fmt.Sprintf("\tparameters = parameters.Param(%s).BeginOptional().%s(%s).EndOptional()", paramName, method, variable))
		} else {
			lines = append(lines, fmt.Sprintf("\tparameters = parameters.Param(%s).%s(%s)", paramName, method, variable))
		}

		return true
	})

	params := strings.Join(lines, "\n")
	return fmt.Sprintf("\tparameters := ydb.ParamsBuilder()\n%s", params)
}

func (v QueryValue) YDBHasParams() bool {
	hasParams := false
	v.ydbIterateNamedParams(func(field Field, method string) bool {
		hasParams = true
		return false
	})
	return hasParams
}

func (v QueryValue) getParameterFields() []Field {
	if v.Struct == nil {
		return []Field{
			{
				Name:   v.Name,
				DBName: v.DBName,
				Type:   v.Typ,
				Column: v.Column,
			},
		}
	}
	return v.Struct.Fields
}

// YDBPair returns the argument name and type for YDB query methods, filtering out interface{} parameters
// that are handled by manualParams instead.
func (v QueryValue) YDBPair() string {
	if v.isEmpty() {
		return ""
	}

	var out []string
	for _, arg := range v.YDBPairs() {
		out = append(out, arg.Name+" "+arg.Type)
	}
	return strings.Join(out, ",")
}

// YDBPairs returns the argument name and type for YDB query methods, filtering out interface{} parameters
// that are handled by manualParams instead.
func (v QueryValue) YDBPairs() []Argument {
	if v.isEmpty() {
		return nil
	}

	if !v.EmitStruct() && v.IsStruct() {
		var out []Argument
		for _, f := range v.Struct.Fields {
			if strings.HasSuffix(f.Type, "interface{}") {
				continue
			}
			out = append(out, Argument{
				Name: escape(toLowerCase(f.Name)),
				Type: f.Type,
			})
		}
		return out
	}

	if strings.HasSuffix(v.Typ, "interface{}") {
		return nil
	}

	return []Argument{
		{
			Name: escape(v.Name),
			Type: v.DefineType(),
		},
	}
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
	Arg          QueryValue
	// Used for :copyfrom
	Table *plugin.Identifier
}

func (q Query) hasRetType() bool {
	scanned := q.Cmd == metadata.CmdOne || q.Cmd == metadata.CmdMany ||
		q.Cmd == metadata.CmdBatchMany || q.Cmd == metadata.CmdBatchOne
	return scanned && !q.Ret.isEmpty()
}

func (q Query) TableIdentifierAsGoSlice() string {
	escapedNames := make([]string, 0, 3)
	for _, p := range []string{q.Table.Catalog, q.Table.Schema, q.Table.Name} {
		if p != "" {
			escapedNames = append(escapedNames, fmt.Sprintf("%q", p))
		}
	}
	return "[]string{" + strings.Join(escapedNames, ", ") + "}"
}

func (q Query) TableIdentifierForMySQL() string {
	escapedNames := make([]string, 0, 3)
	for _, p := range []string{q.Table.Catalog, q.Table.Schema, q.Table.Name} {
		if p != "" {
			escapedNames = append(escapedNames, fmt.Sprintf("`%s`", p))
		}
	}
	return strings.Join(escapedNames, ".")
}
