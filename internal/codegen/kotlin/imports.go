package kotlin

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
)

type importer struct {
	Settings    config.CombinedSettings
	DataClasses []Struct
	Enums       []Enum
	Queries     []Query
}

func (i *importer) usesType(typ string) bool {
	for _, strct := range i.DataClasses {
		for _, f := range strct.Fields {
			if f.Type.Name == typ {
				return true
			}
		}
	}
	return false
}

func (i *importer) Imports(filename string) [][]string {
	switch filename {
	case "Models.kt":
		return i.modelImports()
	case "Querier.kt":
		return i.interfaceImports()
	default:
		return i.queryImports(filename)
	}
}

func (i *importer) interfaceImports() [][]string {
	uses := func(name string) bool {
		for _, q := range i.Queries {
			if !q.Ret.isEmpty() {
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if strings.HasPrefix(f.Type.Name, name) {
						return true
					}
				}
			}
		}
		return false
	}

	std := stdImports(uses)
	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds, runtimeImports(i.Queries)}
}

func (i *importer) modelImports() [][]string {
	std := make(map[string]struct{})
	if i.usesType("LocalDate") {
		std["java.time.LocalDate"] = struct{}{}
	}
	if i.usesType("LocalTime") {
		std["java.time.LocalTime"] = struct{}{}
	}
	if i.usesType("LocalDateTime") {
		std["java.time.LocalDateTime"] = struct{}{}
	}
	if i.usesType("OffsetDateTime") {
		std["java.time.OffsetDateTime"] = struct{}{}
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds}
}

func stdImports(uses func(name string) bool) map[string]struct{} {
	std := map[string]struct{}{
		"java.sql.SQLException": {},
	}
	if uses("LocalDate") {
		std["java.time.LocalDate"] = struct{}{}
	}
	if uses("LocalTime") {
		std["java.time.LocalTime"] = struct{}{}
	}
	if uses("LocalDateTime") {
		std["java.time.LocalDateTime"] = struct{}{}
	}
	if uses("OffsetDateTime") {
		std["java.time.OffsetDateTime"] = struct{}{}
	}
	return std
}

func runtimeImports(kq []Query) []string {
	rt := map[string]struct{}{}
	for _, q := range kq {
		switch q.Cmd {
		case ":one":
			rt["sqlc.runtime.RowQuery"] = struct{}{}
		case ":many":
			rt["sqlc.runtime.ListQuery"] = struct{}{}
		case ":exec":
			rt["sqlc.runtime.ExecuteQuery"] = struct{}{}
		case ":execUpdate":
			rt["sqlc.runtime.ExecuteUpdateQuery"] = struct{}{}
		default:
			panic(fmt.Sprintf("invalid command %q", q.Cmd))
		}
	}
	rts := make([]string, 0, len(rt))
	for s, _ := range rt {
		rts = append(rts, s)
	}
	sort.Strings(rts)
	return rts
}

func (i *importer) queryImports(filename string) [][]string {
	uses := func(name string) bool {
		for _, q := range i.Queries {
			if !q.Ret.isEmpty() {
				if q.Ret.Struct != nil {
					for _, f := range q.Ret.Struct.Fields {
						if f.Type.Name == name {
							return true
						}
					}
				}
				if q.Ret.Type() == name {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if f.Type.Name == name {
						return true
					}
				}
			}
		}
		return false
	}

	hasEnum := func() bool {
		for _, q := range i.Queries {
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if f.Type.IsEnum {
						return true
					}
				}
			}
		}
		return false
	}

	std := stdImports(uses)
	std["java.sql.Connection"] = struct{}{}
	if hasEnum() {
		std["java.sql.Types"] = struct{}{}
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds, runtimeImports(i.Queries)}
}
