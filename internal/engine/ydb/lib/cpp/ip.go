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
			Name: "Ip::FromString",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Ip::SubnetFromString",
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
			Name: "Ip::ToString",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "String"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Ip::ToString",
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
			Name: "Ip::IsIPv4",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Ip::IsIPv6",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Ip::IsEmbeddedIPv4",
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
			Name: "Ip::ConvertToIPv6",
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
			Name: "Ip::GetSubnet",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "Ip::GetSubnet",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint8"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "Ip::GetSubnetByMask",
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
			Name: "Ip::SubnetMatch",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}
