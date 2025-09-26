package lib

import (
	"github.com/sqlc-dev/sqlc/internal/engine/ydb/lib/cpp"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func CppFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, cpp.DateTimeFunctions()...)
	funcs = append(funcs, cpp.DigestFunctions()...)
	funcs = append(funcs, cpp.HyperscanFunctions()...)
	funcs = append(funcs, cpp.IpFunctions()...)
	funcs = append(funcs, cpp.MathFunctions()...)
	funcs = append(funcs, cpp.PcreFunctions()...)
	funcs = append(funcs, cpp.PireFunctions()...)
	funcs = append(funcs, cpp.Re2Functions()...)
	funcs = append(funcs, cpp.StringFunctions()...)
	funcs = append(funcs, cpp.UnicodeFunctions()...)
	funcs = append(funcs, cpp.UrlFunctions()...)
	funcs = append(funcs, cpp.YsonFunctions()...)

	// TODO: Histogram library, KNN library, PostgeSQL library

	return funcs
}
