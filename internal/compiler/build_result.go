package compiler

import (
	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/config"
)

type BuildResult struct {
	enums   []golang.Enum
	structs []golang.Struct
	queries []golang.Query
}

func (r *BuildResult) Structs(settings config.CombinedSettings) []golang.Struct {
	return r.structs
}

func (r *BuildResult) GoQueries(settings config.CombinedSettings) []golang.Query {
	return r.queries
}

func (r *BuildResult) Enums(settings config.CombinedSettings) []golang.Enum {
	return r.enums
}
