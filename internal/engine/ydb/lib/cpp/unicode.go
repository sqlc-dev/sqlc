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
			Name: "unicode_isutf",
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
			Name: "unicode_getlength",
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
			Name: "unicode_find",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "unicode_find",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "unicode_rfind",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "unicode_rfind",
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
			Name: "unicode_substring",
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
			Name: "unicode_normalize",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_normalizenfd",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_normalizenfc",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_normalizenfkd",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_normalizenfkc",
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
			Name: "unicode_translit",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_translit",
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
			Name: "unicode_levensteindistance",
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
			Name: "unicode_fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_fold",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_fold",
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
			Name: "unicode_fold",
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
			Name: "unicode_replaceall",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_replacefirst",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_replacelast",
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
			Name: "unicode_removeall",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_removefirst",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_removelast",
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
			Name: "unicode_tocodepointlist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "unicode_fromcodepointlist",
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
			Name: "unicode_reverse",
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
			Name: "unicode_tolower",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_toupper",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "unicode_totitle",
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
			Name: "unicode_splittolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "unicode_splittolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "unicode_splittolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "unicode_splittolist",
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
			Name: "unicode_joinfromlist",
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
			Name: "unicode_touint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "unicode_touint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Uint16"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "unicode_trytouint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "unicode_trytouint64",
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
			Name: "unicode_strip",
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
			Name: "unicode_isascii",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "unicode_isspace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "unicode_isupper",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "unicode_islower",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "unicode_isalpha",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "unicode_isalnum",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "unicode_ishex",
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
			Name: "unicode_isunicodeset",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Utf8"}},
				{Type: &ast.TypeName{Name: "Utf8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}
