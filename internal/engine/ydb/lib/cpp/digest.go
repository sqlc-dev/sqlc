package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func DigestFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, digestCrcFuncs()...)
	funcs = append(funcs, digestFnvFuncs()...)
	funcs = append(funcs, digestMurmurFuncs()...)
	funcs = append(funcs, digestCityFuncs()...)

	return funcs
}

func digestCrcFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Digest::Crc32c",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "Digest::Crc64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Digest::Crc64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
	}
}

func digestFnvFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Digest::Fnv32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "Digest::Fnv32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "Digest::Fnv64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Digest::Fnv64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
	}
}

func digestMurmurFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Digest::MurMurHash",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Digest::MurMurHash",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Digest::MurMurHash32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "Digest::MurMurHash32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "Digest::MurMurHash2A",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Digest::MurMurHash2A",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Digest::MurMurHash2A32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "Digest::MurMurHash2A32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
	}
}

func digestCityFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Digest::CityHash",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Digest::CityHash",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Digest::CityHash128",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
