package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func YsonFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, ysonParseFuncs()...)
	funcs = append(funcs, ysonFromFuncs()...)
	funcs = append(funcs, ysonWithAttributesFuncs()...)
	funcs = append(funcs, ysonEqualsFuncs()...)
	funcs = append(funcs, ysonGetHashFuncs()...)
	funcs = append(funcs, ysonIsFuncs()...)
	funcs = append(funcs, ysonGetLengthFuncs()...)
	funcs = append(funcs, ysonConvertToFuncs()...)
	funcs = append(funcs, ysonConvertToListFuncs()...)
	funcs = append(funcs, ysonConvertToDictFuncs()...)
	funcs = append(funcs, ysonContainsFuncs()...)
	funcs = append(funcs, ysonLookupFuncs()...)
	funcs = append(funcs, ysonYPathFuncs()...)
	funcs = append(funcs, ysonAttributesFuncs()...)
	funcs = append(funcs, ysonSerializeFuncs()...)
	funcs = append(funcs, ysonOptionsFuncs()...)

	return funcs
}

func ysonParseFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_parse",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Yson"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_parsejson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Json"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_parsejsondecodeutf8",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Json"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_parse",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_parsejson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_parsejsondecodeutf8",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func ysonFromFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_from",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func ysonWithAttributesFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_withattributes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func ysonEqualsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_equals",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func ysonGetHashFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_gethash",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
	}
}

func ysonIsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_isentity",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "yson_isstring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "yson_isdouble",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "yson_isuint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "yson_isint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "yson_isbool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "yson_islist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "yson_isdict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func ysonGetLengthFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_getlength",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
	}
}

func ysonConvertToFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_convertto",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttobool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_converttoint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_converttouint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_converttodouble",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_converttostring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_converttolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func ysonConvertToListFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_converttoboollist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttoint64list",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttouint64list",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttodoublelist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttostringlist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func ysonConvertToDictFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_converttodict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttobooldict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttoint64dict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttouint64dict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttodoubledict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_converttostringdict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func ysonContainsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_contains",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
	}
}

func ysonLookupFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_lookup",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_lookupbool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_lookupint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_lookupuint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_lookupdouble",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_lookupstring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_lookupdict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_lookuplist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func ysonYPathFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_ypath",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_ypathbool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_ypathint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_ypathuint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_ypathdouble",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_ypathstring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_ypathdict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_ypathlist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func ysonAttributesFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_attributes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func ysonSerializeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_serialize",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Yson"},
		},
		{
			Name: "yson_serializetext",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Yson"},
		},
		{
			Name: "yson_serializepretty",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Yson"},
		},
		{
			Name: "yson_serializejson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Json"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_serializejson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Json"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_serializejson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Json"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_serializejson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Json"},
			ReturnTypeNullable: true,
		},
		{
			Name: "yson_serializejson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Json"},
			ReturnTypeNullable: true,
		},
	}
}

func ysonOptionsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "yson_options",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "yson_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
