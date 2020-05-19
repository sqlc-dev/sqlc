package compiler

import (
	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/config"
)

type Result struct {
	enums   []golang.Enum
	structs []golang.Struct
	queries []golang.Query
}

func (r *Result) Structs(settings config.CombinedSettings) []golang.Struct {
	return r.structs
}

func (r *Result) GoQueries(settings config.CombinedSettings) []golang.Query {
	return r.queries
}

func (r *Result) Enums(settings config.CombinedSettings) []golang.Enum {
	return r.enums
}
