package golang

import (
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/metadata"
)

type fileImports struct {
	Std []string
	Dep []string
}

func mergeImports(imps ...fileImports) [][]string {
	if len(imps) == 1 {
		return [][]string{imps[0].Std, imps[0].Dep}
	}

	var stds, pkgs []string
	seenStd := map[string]struct{}{}
	seenPkg := map[string]struct{}{}
	for i := range imps {
		for _, std := range imps[i].Std {
			if _, ok := seenStd[std]; ok {
				continue
			}
			stds = append(stds, std)
			seenStd[std] = struct{}{}
		}
		for _, pkg := range imps[i].Dep {
			if _, ok := seenPkg[pkg]; ok {
				continue
			}
			pkgs = append(pkgs, pkg)
			seenPkg[pkg] = struct{}{}
		}
	}
	return [][]string{stds, pkgs}
}

type importer struct {
	Settings config.CombinedSettings
	Queries  []Query
	Enums    []Enum
	Structs  []Struct
}

func (i *importer) usesType(typ string) bool {
	for _, strct := range i.Structs {
		for _, f := range strct.Fields {
			fType := strings.TrimPrefix(f.Type, "[]")
			if strings.HasPrefix(fType, typ) {
				return true
			}
		}
	}
	return false
}

func (i *importer) usesArrays() bool {
	for _, strct := range i.Structs {
		for _, f := range strct.Fields {
			if strings.HasPrefix(f.Type, "[]") {
				return true
			}
		}
	}
	return false
}

func (i *importer) Imports(filename string) [][]string {
	switch filename {
	case "db.go":
		return mergeImports(i.dbImports())
	case "models.go":
		return mergeImports(i.modelImports())
	case "querier.go":
		return mergeImports(i.interfaceImports())
	default:
		return mergeImports(i.queryImports(filename))
	}
}

func (i *importer) dbImports() fileImports {
	std := []string{"context", "database/sql"}
	if i.Settings.Go.EmitPreparedQueries {
		std = append(std, "fmt")
	}
	return fileImports{Std: std}
}

func (i *importer) interfaceImports() fileImports {
	uses := func(name string) bool {
		for _, q := range i.Queries {
			if q.hasRetType() {
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				if strings.HasPrefix(q.Arg.Type(), name) {
					return true
				}
			}
		}
		return false
	}

	std := map[string]struct{}{
		"context": struct{}{},
	}
	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}
	for _, q := range i.Queries {
		if q.Cmd == metadata.CmdExecResult {
			std["database/sql"] = struct{}{}
		}
	}
	if uses("json.RawMessage") {
		std["encoding/json"] = struct{}{}
	}
	if uses("time.Time") {
		std["time"] = struct{}{}
	}
	if uses("net.IP") {
		std["net"] = struct{}{}
	}
	if uses("net.HardwareAddr") {
		std["net"] = struct{}{}
	}

	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range i.Settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if uses("pq.NullTime") && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}
	_, overrideUUID := overrideTypes["uuid.UUID"]
	if uses("uuid.UUID") && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	// Custom imports
	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && uses(goType) {
			pkg[importPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(pkg))
	for p, _ := range pkg {
		pkgs = append(pkgs, p)
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	sort.Strings(pkgs)
	return fileImports{stds, pkgs}
}

func (i *importer) modelImports() fileImports {
	std := make(map[string]struct{})
	if i.usesType("sql.Null") {
		std["database/sql"] = struct{}{}
	}
	if i.usesType("json.RawMessage") {
		std["encoding/json"] = struct{}{}
	}
	if i.usesType("time.Time") {
		std["time"] = struct{}{}
	}
	if i.usesType("net.IP") {
		std["net"] = struct{}{}
	}
	if i.usesType("net.HardwareAddr") {
		std["net"] = struct{}{}
	}
	if len(i.Enums) > 0 {
		std["fmt"] = struct{}{}
	}

	// Custom imports
	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range i.Settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if i.usesType("pq.NullTime") && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}

	_, overrideUUID := overrideTypes["uuid.UUID"]
	if i.usesType("uuid.UUID") && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && i.usesType(goType) {
			pkg[importPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(pkg))
	for p, _ := range pkg {
		pkgs = append(pkgs, p)
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	sort.Strings(pkgs)
	return fileImports{stds, pkgs}
}

func (i *importer) queryImports(filename string) fileImports {
	var gq []Query
	for _, query := range i.Queries {
		if query.SourceName == filename {
			gq = append(gq, query)
		}
	}

	uses := func(name string) bool {
		for _, q := range gq {
			if q.hasRetType() {
				if q.Ret.EmitStruct() {
					for _, f := range q.Ret.Struct.Fields {
						fType := strings.TrimPrefix(f.Type, "[]")
						if strings.HasPrefix(fType, name) {
							return true
						}
					}
				}
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				if q.Arg.EmitStruct() {
					for _, f := range q.Arg.Struct.Fields {
						fType := strings.TrimPrefix(f.Type, "[]")
						if strings.HasPrefix(fType, name) {
							return true
						}
					}
				}
				if strings.HasPrefix(q.Arg.Type(), name) {
					return true
				}
			}
		}
		return false
	}

	sliceScan := func() bool {
		for _, q := range gq {
			if q.hasRetType() {
				if q.Ret.IsStruct() {
					for _, f := range q.Ret.Struct.Fields {
						if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Ret.Type(), "[]") && q.Ret.Type() != "[]byte" {
						return true
					}
				}
			}
			if !q.Arg.isEmpty() {
				if q.Arg.IsStruct() {
					for _, f := range q.Arg.Struct.Fields {
						if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Arg.Type(), "[]") && q.Arg.Type() != "[]byte" {
						return true
					}
				}
			}
		}
		return false
	}

	std := map[string]struct{}{
		"context": struct{}{},
	}
	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}
	for _, q := range gq {
		if q.Cmd == metadata.CmdExecResult {
			std["database/sql"] = struct{}{}
		}
	}
	if uses("json.RawMessage") {
		std["encoding/json"] = struct{}{}
	}
	if uses("time.Time") {
		std["time"] = struct{}{}
	}
	if uses("net.IP") {
		std["net"] = struct{}{}
	}

	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range i.Settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
	}

	if sliceScan() {
		pkg["github.com/lib/pq"] = struct{}{}
	}
	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if uses("pq.NullTime") && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}
	_, overrideUUID := overrideTypes["uuid.UUID"]
	if uses("uuid.UUID") && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	// Custom imports
	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && uses(goType) {
			pkg[importPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(pkg))
	for p, _ := range pkg {
		pkgs = append(pkgs, p)
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	sort.Strings(pkgs)
	return fileImports{stds, pkgs}
}
