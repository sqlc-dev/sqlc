package clickhouse

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// NewCatalog creates a new ClickHouse catalog with default settings
func NewCatalog() *catalog.Catalog {
	// ClickHouse uses "default" as the default database
	defaultSchemaName := "default"

	cat := &catalog.Catalog{
		DefaultSchema: defaultSchemaName,
		Schemas: []*catalog.Schema{
			newDefaultSchema(defaultSchemaName),
		},
		Extensions: map[string]struct{}{},
	}

	// Register ClickHouse built-in functions with fixed return types
	registerBuiltinFunctions(cat)

	return cat
}

// newDefaultSchema creates the default ClickHouse schema
func newDefaultSchema(name string) *catalog.Schema {
	return &catalog.Schema{
		Name:   name,
		Tables: make([]*catalog.Table, 0),
		Funcs:  make([]*catalog.Function, 0),
	}
}

// registerBuiltinFunctions registers ClickHouse built-in functions in the default schema
func registerBuiltinFunctions(cat *catalog.Catalog) {
	// Find the default schema
	var schema *catalog.Schema
	for _, s := range cat.Schemas {
		if s.Name == cat.DefaultSchema {
			schema = s
			break
		}
	}
	if schema == nil {
		return
	}

	if schema.Funcs == nil {
		schema.Funcs = make([]*catalog.Function, 0)
	}

	// Aggregate functions that always return uint64
	uint64Type := &ast.TypeName{Name: "uint64"}
	int64Type := &ast.TypeName{Name: "int64"}
	anyType := &ast.TypeName{Name: "any"}
	for _, name := range []string{"count", "uniqexact", "countif", "uniq", "uniqcombined", "uniquehll12"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: uint64Type,
			Args: []*catalog.Argument{
				{
					Name: "arg",
					Type: anyType,
					Mode: ast.FuncParamVariadic,
				},
			},
		})
	}

	// Statistical aggregate functions
	float64Type := &ast.TypeName{Name: "float64"}
	for _, name := range []string{"varsamp", "varpop", "stddevsamp", "stddevpop", "corr", "covariance", "avg", "avgif"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: float64Type,
			Args: []*catalog.Argument{
				{
					Name: "arg",
					Type: anyType,
					Mode: ast.FuncParamVariadic,
				},
			},
		})
	}

	// Date/Time functions
	dateType := &ast.TypeName{Name: "date"}
	timeType := &ast.TypeName{Name: "timestamp"}

	// Functions returning Date
	for _, name := range []string{"todate", "todate32", "today", "yesterday"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: dateType,
			Args: []*catalog.Argument{
				{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic},
			},
		})
	}

	// Functions returning Timestamp (DateTime/DateTime64)
	for _, name := range []string{"now", "todatetime", "todatetime64", "parseDateTime", "parseDateTimeBestEffort", "parseDateTime64BestEffort"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: timeType,
			Args: []*catalog.Argument{
				{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic},
			},
		})
	}

	// Functions returning integers (components)
	int32Type := &ast.TypeName{Name: "int32"}
	for _, name := range []string{"toyear", "tomonth", "todayofmonth", "tohour", "tominute", "tosecond", "tounixtimestamp"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: int32Type,
			Args: []*catalog.Argument{
				{Name: "arg", Type: anyType, Mode: ast.FuncParamIn},
			},
		})
	}

	// String functions
	stringType := &ast.TypeName{Name: "text"}
	for _, name := range []string{
		"concat", "substring", "lower", "upper", "trim", "ltrim", "rtrim", "reverse",
		"replace", "replaceall", "replaceregexpone", "replaceregexpall",
		"format", "tostring", "base64encode", "base64decode", "hex", "unhex",
		"extract", "extractall", "splitbychar", "splitbystring", "splitbyregexp",
		"domain", "domainwithoutwww", "topleveldomain", "protocol", "path",
		"cutquerystring", "cutfragment", "cutwww", "cutquerystringandfragment",
	} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: stringType,
			Args: []*catalog.Argument{
				{
					Name: "arg",
					Type: anyType,
					Mode: ast.FuncParamVariadic,
				},
			},
		})
	}

	// Hashing functions (returning strings usually, or numbers)
	for _, name := range []string{"md5", "sha1", "sha224", "sha256", "halfmd5"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: stringType,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}
	for _, name := range []string{"siphash64", "siphash128", "cityhash64", "inthash32", "inthash64", "farmhash64", "metrohash64"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: uint64Type, // hashes are often uint64
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}

	// JSON Functions
	// JSONExtract is special (handled by type resolver), but we register it with anyType fallback
	// JSONHas, JSONLength, JSONType, JSONExtractString, etc.

	// JSON functions returning generic/string
	for _, name := range []string{"jsonextract", "jsonextractstring", "jsonextractraw", "tojsonstring"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: stringType, // Fallback to string/text
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}

	// JSON functions returning bool
	boolType := &ast.TypeName{Name: "bool"}
	for _, name := range []string{"jsonhas", "jsonextractbool", "isvalidjson"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: boolType,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}

	// JSON functions returning int/float
	for _, name := range []string{"jsonlength", "jsonextractint", "jsonextractuint"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: int64Type,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}
	for _, name := range []string{"jsonextractfloat"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: float64Type,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}

	// Array Functions (generic return)
	// Many array functions return arrays, which sqlc handles as generic "any" mostly unless resolved
	for _, name := range []string{
		"array", "arrayconcat", "arrayslice", "arraypushback", "arraypushfront",
		"arraypopback", "arraypopfront", "arrayresize", "arrayfilter", "arraymap",
		"arrayreverse", "arraysort", "arraydistinct", "arrayuniq", "arrayjoin",
		"arrayenumerate", "arrayenumerateuniq", "arrayflatten", "arraycompact",
		"arrayzip", "arrayreduce", "arrayfold",
	} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: anyType,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}
	// Array functions returning bool
	for _, name := range []string{"has", "hasall", "hasany", "hassubstr", "arrayexists", "arrayall"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: boolType,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}
	// Array functions returning int/uint
	for _, name := range []string{"length", "empty", "notempty", "arraycount", "indexof"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: uint64Type,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamVariadic}},
		})
	}

	// Type Conversion
	// toInt*, toUInt*, toFloat*, toDecimal*
	for _, name := range []string{"toint8", "toint16", "toint32", "toint64", "toint128", "toint256"} {
		// Simplified to int64 for now for sqlc mapping purposes, though could be specific
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: int64Type,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamIn}},
		})
	}
	for _, name := range []string{"touint8", "touint16", "touint32", "touint64", "touint128", "touint256"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: uint64Type,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamIn}},
		})
	}
	for _, name := range []string{"tofloat32", "tofloat64"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: float64Type,
			Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamIn}},
		})
	}

	// UUID
	schema.Funcs = append(schema.Funcs, &catalog.Function{
		Name:       "generateuuidv4",
		ReturnType: &ast.TypeName{Name: "uuid"},
	})

	// IP
	schema.Funcs = append(schema.Funcs, &catalog.Function{
		Name:       "ipv4stringtonum",
		ReturnType: uint64Type,
		Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamIn}},
	})
	schema.Funcs = append(schema.Funcs, &catalog.Function{
		Name:       "ipv4numtostring",
		ReturnType: stringType,
		Args:       []*catalog.Argument{{Name: "arg", Type: anyType, Mode: ast.FuncParamIn}},
	})

	// Functions with context-dependent return types
	// These are registered with a placeholder return type and will be handled specially
	// by the compiler when analyzing query output columns

	// arrayJoin(Array(T)) returns T
	schema.Funcs = append(schema.Funcs, &catalog.Function{
		Name:       "arrayjoin",
		ReturnType: anyType,
		Args: []*catalog.Argument{
			{
				Name: "arr",
				Type: anyType,
				Mode: ast.FuncParamIn,
			},
		},
	})

	// argMin and argMax return the type of their first argument
	for _, name := range []string{"argmin", "argmax"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: anyType,
			Args: []*catalog.Argument{
				{Name: "val", Type: anyType, Mode: ast.FuncParamIn},
				{Name: "arg", Type: anyType, Mode: ast.FuncParamIn},
			},
		})
	}

	// argMinIf and argMaxIf return the type of their first argument
	for _, name := range []string{"argminif", "argmaxif"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: anyType,
			Args: []*catalog.Argument{
				{Name: "val", Type: anyType, Mode: ast.FuncParamIn},
				{Name: "arg", Type: anyType, Mode: ast.FuncParamIn},
				{Name: "cond", Type: anyType, Mode: ast.FuncParamIn},
			},
		})
	}

	// any, anyLast, anyHeavy return the type of their argument
	for _, name := range []string{"any", "anylast", "anyheavy", "min", "max", "sum"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: anyType,
			Args: []*catalog.Argument{
				{
					Name: "arg",
					Type: anyType,
					Mode: ast.FuncParamIn,
				},
			},
		})
	}

	// anyIf, anyLastIf, anyHeavyIf, minIf, maxIf, sumIf return the type of their argument
	for _, name := range []string{"anyif", "anylastif", "anyheavyif", "minif", "maxif", "sumif"} {
		schema.Funcs = append(schema.Funcs, &catalog.Function{
			Name:       name,
			ReturnType: anyType,
			Args: []*catalog.Argument{
				{
					Name: "arg",
					Type: anyType,
					Mode: ast.FuncParamIn,
				},
				{
					Name: "cond",
					Type: anyType,
					Mode: ast.FuncParamIn,
				},
			},
		})
	}
}
