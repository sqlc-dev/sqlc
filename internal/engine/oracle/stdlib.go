package oracle

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// Oracle built-in functions and datatypes.
//
// References:
//   - Datatypes:  https://docs.oracle.com/en/database/oracle/oracle-database/21/sqlrf/Data-Types.html
//   - Functions:  https://docs.oracle.com/en/database/oracle/oracle-database/21/sqlrf/Single-Row-Functions.html
//                 https://docs.oracle.com/en/database/oracle/oracle-database/21/sqlrf/Aggregate-Functions.html
//
// The type names here are stored lower-cased to match the normalization applied
// by the converter (see identifier / normalizeDatatype in convert.go).

// OracleTypes is the set of Oracle built-in scalar type names the engine
// recognizes (lower-cased to match the converter's normalization). It is used by
// the code generator (Phase 5) to map Oracle column types to Go types.
var OracleTypes = []string{
	"number",
	"numeric",
	"dec",
	"decimal",
	"integer",
	"int",
	"smallint",
	"float",
	"real",
	"double precision",
	"binary_float",
	"binary_double",
	"varchar2",
	"nvarchar2",
	"varchar",
	"char",
	"nchar",
	"clob",
	"nclob",
	"long",
	"blob",
	"bfile",
	"raw",
	"long raw",
	"date",
	"timestamp",
	"rowid",
	"urowid",
	"xmltype",
	"boolean",
}

// IsBuiltinType reports whether name (case-insensitive) is a known Oracle
// built-in scalar type.
func IsBuiltinType(name string) bool {
	lower := toLowerASCII(name)
	for _, t := range OracleTypes {
		if t == lower {
			return true
		}
	}
	return false
}

func toLowerASCII(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

func defaultSchema(name string) *catalog.Schema {
	s := &catalog.Schema{Name: name}
	s.Funcs = stdFunctions()
	return s
}

// arg is a small helper for declaring a positional function argument.
func arg(typeName string) *catalog.Argument {
	return &catalog.Argument{Type: &ast.TypeName{Name: typeName}}
}

// ret is a small helper for declaring a function return type.
func ret(typeName string) *ast.TypeName {
	return &ast.TypeName{Name: typeName}
}

// stdFunctions returns Oracle's most common built-in functions. This is a
// starter set covering aggregates, string, numeric, date and conversion
// functions; more can be added incrementally.
func stdFunctions() []*catalog.Function {
	return []*catalog.Function{
		// --- Aggregate functions ---
		{Name: "COUNT", Args: []*catalog.Argument{}, ReturnType: ret("number")},
		{Name: "COUNT", Args: []*catalog.Argument{arg("any")}, ReturnType: ret("number")},
		{Name: "SUM", Args: []*catalog.Argument{arg("number")}, ReturnType: ret("number"), ReturnTypeNullable: true},
		{Name: "AVG", Args: []*catalog.Argument{arg("number")}, ReturnType: ret("number"), ReturnTypeNullable: true},
		{Name: "MIN", Args: []*catalog.Argument{arg("any")}, ReturnType: ret("any"), ReturnTypeNullable: true},
		{Name: "MAX", Args: []*catalog.Argument{arg("any")}, ReturnType: ret("any"), ReturnTypeNullable: true},

		// --- String functions ---
		{Name: "LENGTH", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("number")},
		{Name: "SUBSTR", Args: []*catalog.Argument{arg("varchar2"), arg("number")}, ReturnType: ret("varchar2")},
		{Name: "SUBSTR", Args: []*catalog.Argument{arg("varchar2"), arg("number"), arg("number")}, ReturnType: ret("varchar2")},
		{Name: "UPPER", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("varchar2")},
		{Name: "LOWER", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("varchar2")},
		{Name: "TRIM", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("varchar2")},
		{Name: "LTRIM", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("varchar2")},
		{Name: "RTRIM", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("varchar2")},
		{Name: "REPLACE", Args: []*catalog.Argument{arg("varchar2"), arg("varchar2"), arg("varchar2")}, ReturnType: ret("varchar2")},
		{Name: "CONCAT", Args: []*catalog.Argument{arg("varchar2"), arg("varchar2")}, ReturnType: ret("varchar2")},
		{Name: "INSTR", Args: []*catalog.Argument{arg("varchar2"), arg("varchar2")}, ReturnType: ret("number")},

		// --- Numeric functions ---
		{Name: "ABS", Args: []*catalog.Argument{arg("number")}, ReturnType: ret("number")},
		{Name: "CEIL", Args: []*catalog.Argument{arg("number")}, ReturnType: ret("number")},
		{Name: "FLOOR", Args: []*catalog.Argument{arg("number")}, ReturnType: ret("number")},
		{Name: "ROUND", Args: []*catalog.Argument{arg("number")}, ReturnType: ret("number")},
		{Name: "ROUND", Args: []*catalog.Argument{arg("number"), arg("number")}, ReturnType: ret("number")},
		{Name: "TRUNC", Args: []*catalog.Argument{arg("number")}, ReturnType: ret("number")},
		{Name: "MOD", Args: []*catalog.Argument{arg("number"), arg("number")}, ReturnType: ret("number")},
		{Name: "POWER", Args: []*catalog.Argument{arg("number"), arg("number")}, ReturnType: ret("number")},
		{Name: "SQRT", Args: []*catalog.Argument{arg("number")}, ReturnType: ret("number")},

		// --- Date/time functions ---
		{Name: "SYSDATE", Args: []*catalog.Argument{}, ReturnType: ret("date")},
		{Name: "SYSTIMESTAMP", Args: []*catalog.Argument{}, ReturnType: ret("timestamp")},
		{Name: "CURRENT_DATE", Args: []*catalog.Argument{}, ReturnType: ret("date")},
		{Name: "CURRENT_TIMESTAMP", Args: []*catalog.Argument{}, ReturnType: ret("timestamp")},
		{Name: "ADD_MONTHS", Args: []*catalog.Argument{arg("date"), arg("number")}, ReturnType: ret("date")},
		{Name: "MONTHS_BETWEEN", Args: []*catalog.Argument{arg("date"), arg("date")}, ReturnType: ret("number")},
		{Name: "EXTRACT", Args: []*catalog.Argument{arg("any")}, ReturnType: ret("number")},

		// --- Conversion functions ---
		{Name: "TO_CHAR", Args: []*catalog.Argument{arg("any")}, ReturnType: ret("varchar2")},
		{Name: "TO_CHAR", Args: []*catalog.Argument{arg("any"), arg("varchar2")}, ReturnType: ret("varchar2")},
		{Name: "TO_NUMBER", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("number")},
		{Name: "TO_DATE", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("date")},
		{Name: "TO_DATE", Args: []*catalog.Argument{arg("varchar2"), arg("varchar2")}, ReturnType: ret("date")},
		{Name: "TO_TIMESTAMP", Args: []*catalog.Argument{arg("varchar2")}, ReturnType: ret("timestamp")},
		{Name: "CAST", Args: []*catalog.Argument{arg("any")}, ReturnType: ret("any")},

		// --- NULL-handling / general functions ---
		{Name: "NVL", Args: []*catalog.Argument{arg("any"), arg("any")}, ReturnType: ret("any")},
		{Name: "NVL2", Args: []*catalog.Argument{arg("any"), arg("any"), arg("any")}, ReturnType: ret("any")},
		{Name: "COALESCE", Args: []*catalog.Argument{arg("any"), arg("any")}, ReturnType: ret("any"), ReturnTypeNullable: true},
		{Name: "DECODE", Args: []*catalog.Argument{arg("any"), arg("any"), arg("any")}, ReturnType: ret("any"), ReturnTypeNullable: true},
		{Name: "GREATEST", Args: []*catalog.Argument{arg("any")}, ReturnType: ret("any")},
		{Name: "LEAST", Args: []*catalog.Argument{arg("any")}, ReturnType: ret("any")},
	}
}
