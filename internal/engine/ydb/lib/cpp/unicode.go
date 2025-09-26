package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func UnicodeFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, unicodeCheckFuncs()...)
	funcs = append(funcs, unicodeLengthFuncs()...)
	funcs = append(funcs, unicodeFindFuncs()...)
	funcs = append(funcs, unicodeSubstringFuncs()...)
	funcs = append(funcs, unicodeNormalizeFuncs()...)
	funcs = append(funcs, unicodeTranslitFuncs()...)
	funcs = append(funcs, unicodeLevensteinFuncs()...)
	funcs = append(funcs, unicodeFoldFuncs()...)
	funcs = append(funcs, unicodeReplaceFuncs()...)
	funcs = append(funcs, unicodeRemoveFuncs()...)
	funcs = append(funcs, unicodeCodePointFuncs()...)
	funcs = append(funcs, unicodeReverseFuncs()...)
	funcs = append(funcs, unicodeCaseFuncs()...)
	funcs = append(funcs, unicodeSplitJoinFuncs()...)
	funcs = append(funcs, unicodeToUint64Funcs()...)
	funcs = append(funcs, unicodeStripFuncs()...)
	funcs = append(funcs, unicodeIsFuncs()...)
	funcs = append(funcs, unicodeIsUnicodeSetFuncs()...)

	return funcs
}

func unicodeCheckFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::IsUtf",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func unicodeLengthFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::GetLength",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
	}
}

func unicodeFindFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::Find",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Unicode::Find",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Unicode::RFind",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Unicode::RFind",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
	}
}

func unicodeSubstringFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::Substring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeNormalizeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::Normalize",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::NormalizeNFD",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::NormalizeNFC",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::NormalizeNFKD",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::NormalizeNFKC",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeTranslitFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::Translit",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::Translit",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeLevensteinFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::LevensteinDistance",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
	}
}

func unicodeFoldFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::Fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::Fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::Fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::Fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::Fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::Fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeReplaceFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::ReplaceAll",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::ReplaceFirst",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::ReplaceLast",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeRemoveFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::RemoveAll",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::RemoveFirst",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::RemoveLast",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeCodePointFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::ToCodePointList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Unicode::FromCodePointList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeReverseFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::Reverse",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeCaseFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::ToLower",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::ToUpper",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Unicode::ToTitle",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeSplitJoinFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::SplitToList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Unicode::SplitToList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Unicode::SplitToList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Unicode::SplitToList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Unicode::JoinFromList",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeToUint64Funcs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::ToUint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Unicode::ToUint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Uint16"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Unicode::TryToUint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Unicode::TryToUint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Uint16"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
	}
}

func unicodeStripFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::Strip",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
	}
}

func unicodeIsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::IsAscii",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Unicode::IsSpace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Unicode::IsUpper",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Unicode::IsLower",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Unicode::IsAlpha",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Unicode::IsAlnum",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Unicode::IsHex",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func unicodeIsUnicodeSetFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Unicode::IsUnicodeSet",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}
