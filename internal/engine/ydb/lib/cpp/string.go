package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func StringFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, stringBase32Funcs()...)
	funcs = append(funcs, stringBase64Funcs()...)
	funcs = append(funcs, stringEscapeFuncs()...)
	funcs = append(funcs, stringHexFuncs()...)
	funcs = append(funcs, stringHtmlFuncs()...)
	funcs = append(funcs, stringCgiFuncs()...)

	return funcs
}

func stringBase32Funcs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "String::Base32Encode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "String::Base32Decode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "String::Base32StrictDecode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
	}
}

func stringBase64Funcs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "String::Base64Encode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "String::Base64Decode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "String::Base64StrictDecode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
	}
}

func stringEscapeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "String::EscapeC",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "String::UnescapeC",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}

func stringHexFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "String::HexEncode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "String::HexDecode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
	}
}

func stringHtmlFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "String::EncodeHtml",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "String::DecodeHtml",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}

func stringCgiFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "String::CgiEscape",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "String::CgiUnescape",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}
