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
			Name: "Yson::Parse",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Yson"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ParseJson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Json"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ParseJsonDecodeUtf8",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Json"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::Parse",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::ParseJson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::ParseJsonDecodeUtf8",
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
			Name: "Yson::From",
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
			Name: "Yson::WithAttributes",
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
			Name: "Yson::Equals",
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
			Name: "Yson::GetHash",
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
			Name: "Yson::IsEntity",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Yson::IsString",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Yson::IsDouble",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Yson::IsUint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Yson::IsInt64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Yson::IsBool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Yson::IsList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Yson::IsDict",
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
			Name: "Yson::GetLength",
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
			Name: "Yson::ConvertTo",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToBool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::ConvertToInt64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::ConvertToUint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::ConvertToDouble",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::ConvertToString",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::ConvertToList",
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
			Name: "Yson::ConvertToBoolList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToInt64List",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToUint64List",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToDoubleList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToStringList",
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
			Name: "Yson::ConvertToDict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToBoolDict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToInt64Dict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToUint64Dict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToDoubleDict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::ConvertToStringDict",
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
			Name: "Yson::Contains",
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
			Name: "Yson::Lookup",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::LookupBool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::LookupInt64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::LookupUint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::LookupDouble",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::LookupString",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::LookupDict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::LookupList",
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
			Name: "Yson::YPath",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::YPathBool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::YPathInt64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::YPathUint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::YPathDouble",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::YPathString",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::YPathDict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::YPathList",
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
			Name: "Yson::Attributes",
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
			Name: "Yson::Serialize",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Yson"},
		},
		{
			Name: "Yson::SerializeText",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Yson"},
		},
		{
			Name: "Yson::SerializePretty",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Yson"},
		},
		{
			Name: "Yson::SerializeJson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Json"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::SerializeJson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Json"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::SerializeJson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Json"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Yson::SerializeJson",
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
			Name: "Yson::SerializeJson",
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
			Name:       "Yson::Options",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::Options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Yson::Options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
