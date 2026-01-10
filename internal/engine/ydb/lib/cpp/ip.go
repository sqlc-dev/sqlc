package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func IpFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, ipFromStringFuncs()...)
	funcs = append(funcs, ipToStringFuncs()...)
	funcs = append(funcs, ipCheckFuncs()...)
	funcs = append(funcs, ipConvertFuncs()...)
	funcs = append(funcs, ipSubnetFuncs()...)
	funcs = append(funcs, ipMatchFuncs()...)

	return funcs
}

func ipFromStringFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "ip_fromstring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "ip_subnetfromstring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
	}
}

func ipToStringFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "ip_tostring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "ip_tostring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
	}
}

func ipCheckFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "ip_isipv4",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "ip_isipv6",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "ip_isembeddedipv4",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func ipConvertFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "ip_converttoipv6",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}

func ipSubnetFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "ip_getsubnet",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "ip_getsubnet",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint8"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "ip_getsubnetbymask",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}

func ipMatchFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "ip_subnetmatch",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}
