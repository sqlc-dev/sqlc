package compiler

import (
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/dinosql"
)

type Result struct {
	structs []dinosql.GoStruct
	queries []dinosql.GoQuery
}

func (r *Result) Structs(settings config.CombinedSettings) []dinosql.GoStruct {
	return r.structs
}

func (r *Result) GoQueries(settings config.CombinedSettings) []dinosql.GoQuery {
	return r.queries
}

func (r *Result) Enums(settings config.CombinedSettings) []dinosql.GoEnum {
	return nil
}
