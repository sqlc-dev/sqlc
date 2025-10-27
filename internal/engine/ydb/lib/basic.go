package lib

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func BasicFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, lengthFuncs()...)
	funcs = append(funcs, substringFuncs()...)
	funcs = append(funcs, findFuncs()...)
	funcs = append(funcs, rfindFuncs()...)
	funcs = append(funcs, startsWithFuncs()...)
	funcs = append(funcs, endsWithFuncs()...)
	funcs = append(funcs, ifFuncs()...)
	funcs = append(funcs, nanvlFuncs()...)
	funcs = append(funcs, randomFuncs()...)
	funcs = append(funcs, currentUtcFuncs()...)
	funcs = append(funcs, currentTzFuncs()...)
	funcs = append(funcs, addTimezoneFuncs()...)
	funcs = append(funcs, removeTimezoneFuncs()...)
	funcs = append(funcs, versionFuncs()...)
	funcs = append(funcs, ensureFuncs()...)
	funcs = append(funcs, assumeStrictFuncs()...)
	funcs = append(funcs, likelyFuncs()...)
	funcs = append(funcs, evaluateFuncs()...)
	funcs = append(funcs, simpleTypesLiteralsFuncs()...)
	funcs = append(funcs, toFromBytesFuncs()...)
	funcs = append(funcs, byteAtFuncs()...)
	funcs = append(funcs, testClearSetFlipBitFuncs()...)
	funcs = append(funcs, absFuncs()...)
	funcs = append(funcs, justUnwrapNothingFuncs()...)
	funcs = append(funcs, pickleUnpickleFuncs()...)
	funcs = append(funcs, asTableFuncs()...)

	// todo: implement functions:
	// Udf, AsTuple, AsStruct, AsList, AsDict, AsSet, AsListStrict, AsDictStrict, AsSetStrict,
	// Variant, AsVariant, Visit, VisitOrDefault, VariantItem, Way, DynamicVariant,
	// Enum, AsEnum, AsTagged, Untag, TableRow, Callable,
	// StaticMap, StaticZip, StaticFold, StaticFold1,
	// AggregationFactory, AggregateTransformInput, AggregateTransformOutput, AggregateFlatten

	return funcs
}

func lengthFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "LENGTH",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "LEN",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
	}
}

func substringFuncs() []*catalog.Function {
	funcs := []*catalog.Function{
		{
			Name: "Substring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "Substring",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
	return funcs
}

func findFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Find",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "Find",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
	}
}

func rfindFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "RFind",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "RFind",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
	}
}

func startsWithFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "StartsWith",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func endsWithFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "EndsWith",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func ifFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: false,
		},
	}
}

func nanvlFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "NANVL",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func randomFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Random",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "RandomNumber",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "RandomUuid",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "Uuid"},
		},
	}
}

func currentUtcFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "CurrentUtcDate",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "Date"},
		},
		{
			Name: "CurrentUtcDatetime",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "Datetime"},
		},
		{
			Name: "CurrentUtcTimestamp",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
	}
}

func currentTzFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "CurrentTzDate",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "TzDate"},
		},
		{
			Name: "CurrentTzDatetime",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "TzDatetime"},
		},
		{
			Name: "CurrentTzTimestamp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "TzTimestamp"},
		},
	}
}

func addTimezoneFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "AddTimezone",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func removeTimezoneFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "RemoveTimezone",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func versionFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name:       "Version",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}

func ensureFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Ensure",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "EnsureType",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "EnsureConvertibleTo",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func assumeStrictFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "AssumeStrict",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func likelyFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Likely",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func evaluateFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "EvaluateExpr",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "EvaluateAtom",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func simpleTypesLiteralsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Bool",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Uint8",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "Int32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Int32"},
		},
		{
			Name: "Uint32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "Int64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Int64"},
		},
		{
			Name: "Uint64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "Float",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Float"},
		},
		{
			Name: "Double",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Decimal",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "Uint8"}}, // precision
				{Type: &ast.TypeName{Name: "Uint8"}}, // scale
			},
			ReturnType: &ast.TypeName{Name: "Decimal"},
		},
		{
			Name: "String",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "Utf8",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Utf8"},
		},
		{
			Name: "Yson",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Yson"},
		},
		{
			Name: "Json",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Json"},
		},
		{
			Name: "Date",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Date"},
		},
		{
			Name: "Datetime",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Datetime"},
		},
		{
			Name: "Timestamp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "Interval",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "TzDate",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "TzDate"},
		},
		{
			Name: "TzDatetime",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "TzDatetime"},
		},
		{
			Name: "TzTimestamp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "TzTimestamp"},
		},
		{
			Name: "Uuid",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "Uuid"},
		},
	}
}

func toFromBytesFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "ToBytes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "FromBytes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func byteAtFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "ByteAt",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
	}
}

func testClearSetFlipBitFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "TestBit",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Uint8"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "ClearBit",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Uint8"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "SetBit",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Uint8"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "FlipBit",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Uint8"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func absFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Abs",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func justUnwrapNothingFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Just",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Unwrap",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Unwrap",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Nothing",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func pickleUnpickleFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Pickle",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "StablePickle",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "Unpickle",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func asTableFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "AS_TABLE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
		},
	}
}
