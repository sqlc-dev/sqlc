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
			Name: "string_base32encode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "string_base32decode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "string_base32strictdecode",
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
			Name: "string_base64encode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "string_base64decode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "string_base64strictdecode",
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
			Name: "string_escapec",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "string_unescapec",
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
			Name: "string_hexencode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "string_hexdecode",
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
			Name: "string_encodehtml",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "string_decodehtml",
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
			Name: "string_cgiescape",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "string_cgiunescape",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}
