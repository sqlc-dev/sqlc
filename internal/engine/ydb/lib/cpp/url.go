package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func UrlFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, urlNormalizeFuncs()...)
	funcs = append(funcs, urlEncodeDecodeFuncs()...)
	funcs = append(funcs, urlParseFuncs()...)
	funcs = append(funcs, urlGetFuncs()...)
	funcs = append(funcs, urlDomainFuncs()...)
	funcs = append(funcs, urlCutFuncs()...)
	funcs = append(funcs, urlPunycodeFuncs()...)
	funcs = append(funcs, urlQueryStringFuncs()...)

	return funcs
}

func urlNormalizeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "url_normalize",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_normalizewithdefaulthttpscheme",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
	}
}

func urlEncodeDecodeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "url_encode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_decode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
	}
}

func urlParseFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "url_parse",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func urlGetFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "url_getscheme",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "url_gethost",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_gethostport",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_getschemehost",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_getschemehostport",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_getport",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_gettail",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_getpath",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_getfragment",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_getcgiparam",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_getdomain",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint8"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
	}
}

func urlDomainFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "url_gettld",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "url_isknowntld",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "url_iswellknowntld",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "url_getdomainlevel",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "url_getsignificantdomain",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "url_getsignificantdomain",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "url_getowner",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}

func urlCutFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "url_cutscheme",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_cutwww",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_cutwww2",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_cutquerystringandfragment",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}

func urlPunycodeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "url_hostnametopunycode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_forcehostnametopunycode",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "url_punycodetohostname",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "url_forcepunycodetohostname",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "url_canbepunycodehostname",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func urlQueryStringFuncs() []*catalog.Function {
	// fixme: rewrite with containers if possible
	return []*catalog.Function{
		{
			Name: "url_querystringtolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtolist",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtodict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtodict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtodict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtodict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_querystringtodict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "url_buildquerystring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "url_buildquerystring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}
