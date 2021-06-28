package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/metadata"
)

type fileImports struct {
	Std []ImportSpec
	Dep []ImportSpec
}

type ImportSpec struct {
	ID   string
	Path string
}

func (s ImportSpec) String() string {
	if s.ID != "" {
		return fmt.Sprintf("%s \"%s\"", s.ID, s.Path)
	} else {
		return fmt.Sprintf("\"%s\"", s.Path)
	}
}

func mergeImports(imps ...fileImports) [][]ImportSpec {
	if len(imps) == 1 {
		return [][]ImportSpec{
			imps[0].Std,
			imps[0].Dep,
		}
	}

	var stds, pkgs []ImportSpec
	seenStd := map[string]struct{}{}
	seenPkg := map[string]struct{}{}
	for i := range imps {
		for _, spec := range imps[i].Std {
			if _, ok := seenStd[spec.Path]; ok {
				continue
			}
			stds = append(stds, spec)
			seenStd[spec.Path] = struct{}{}
		}
		for _, spec := range imps[i].Dep {
			if _, ok := seenPkg[spec.Path]; ok {
				continue
			}
			pkgs = append(pkgs, spec)
			seenPkg[spec.Path] = struct{}{}
		}
	}
	return [][]ImportSpec{stds, pkgs}
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

func (i *importer) Imports(filename string) [][]ImportSpec {
	dbFileName := "db.go"
	if i.Settings.Go.OutputDBFileName != "" {
		dbFileName = i.Settings.Go.OutputDBFileName
	}
	modelsFileName := "models.go"
	if i.Settings.Go.OutputModelsFileName != "" {
		modelsFileName = i.Settings.Go.OutputModelsFileName
	}
	querierFileName := "querier.go"
	if i.Settings.Go.OutputQuerierFileName != "" {
		querierFileName = i.Settings.Go.OutputQuerierFileName
	}

	switch filename {
	case dbFileName:
		return mergeImports(i.dbImports())
	case modelsFileName:
		return mergeImports(i.modelImports())
	case querierFileName:
		return mergeImports(i.interfaceImports())
	default:
		return mergeImports(i.queryImports(filename))
	}
}

func (i *importer) dbImports() fileImports {
	var pkg []ImportSpec
	std := []ImportSpec{
		{Path: "context"},
	}

	driver := DriverFromString(i.Settings.Go.Driver)
	switch driver {
	case PgxDriver:
		pkg = append(pkg, ImportSpec{Path: "github.com/jackc/pgconn"})
		pkg = append(pkg, ImportSpec{Path: "github.com/jackc/pgx/v4"})
	default:
		std = append(std, ImportSpec{Path: "database/sql"})
		if i.Settings.Go.EmitPreparedQueries {
			std = append(std, ImportSpec{Path: "fmt"})
		}
	}

	sort.Slice(std, func(i, j int) bool { return std[i].Path < std[j].Path })
	sort.Slice(pkg, func(i, j int) bool { return pkg[i].Path < pkg[j].Path })
	return fileImports{Std: std, Dep: pkg}
}

var stdlibTypes = map[string]string{
	"json.RawMessage":  "encoding/json",
	"time.Time":        "time",
	"net.IP":           "net",
	"net.HardwareAddr": "net",
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
		"context": {},
	}
	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}

	pkg := make(map[ImportSpec]struct{})

	driver := DriverFromString(i.Settings.Go.Driver)
	for _, q := range i.Queries {
		if q.Cmd == metadata.CmdExecResult {
			switch driver {
			case PgxDriver:
				pkg[ImportSpec{Path: "github.com/jackc/pgconn"}] = struct{}{}
			default:
				std["database/sql"] = struct{}{}
			}
		}
	}
	for typeName, pkg := range stdlibTypes {
		if uses(typeName) {
			std[pkg] = struct{}{}
		}
	}

	overrideTypes := map[string]string{}
	for _, o := range i.Settings.Overrides {
		if o.GoBasicType || o.GoTypeName == "" {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoImportPath
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if uses("pq.NullTime") && !overrideNullTime {
		pkg[ImportSpec{Path: "github.com/lib/pq"}] = struct{}{}
	}
	_, overrideUUID := overrideTypes["uuid.UUID"]
	if uses("uuid.UUID") && !overrideUUID {
		pkg[ImportSpec{Path: "github.com/google/uuid"}] = struct{}{}
	}

	// Custom imports
	for _, o := range i.Settings.Overrides {
		if o.GoBasicType || o.GoTypeName == "" {
			continue
		}
		_, alreadyImported := std[o.GoImportPath]
		hasPackageAlias := o.GoPackage != ""
		if (!alreadyImported || hasPackageAlias) && uses(o.GoTypeName) {
			pkg[ImportSpec{Path: o.GoImportPath, ID: o.GoPackage}] = struct{}{}
		}
	}

	pkgs := make([]ImportSpec, 0, len(pkg))
	for spec := range pkg {
		pkgs = append(pkgs, spec)
	}

	stds := make([]ImportSpec, 0, len(std))
	for path := range std {
		stds = append(stds, ImportSpec{Path: path})
	}

	sort.Slice(stds, func(i, j int) bool { return stds[i].Path < stds[j].Path })
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].Path < pkgs[j].Path })
	return fileImports{stds, pkgs}
}

func (i *importer) modelImports() fileImports {
	std := make(map[string]struct{})
	if i.usesType("sql.Null") {
		std["database/sql"] = struct{}{}
	}
	for typeName, pkg := range stdlibTypes {
		if i.usesType(typeName) {
			std[pkg] = struct{}{}
		}
	}
	if len(i.Enums) > 0 {
		std["fmt"] = struct{}{}
	}

	// Custom imports
	pkg := make(map[ImportSpec]struct{})
	overrideTypes := map[string]string{}
	for _, o := range i.Settings.Overrides {
		if o.GoBasicType || o.GoTypeName == "" {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoImportPath
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if i.usesType("pq.NullTime") && !overrideNullTime {
		pkg[ImportSpec{Path: "github.com/lib/pq"}] = struct{}{}
	}

	_, overrideUUID := overrideTypes["uuid.UUID"]
	if i.usesType("uuid.UUID") && !overrideUUID {
		pkg[ImportSpec{Path: "github.com/google/uuid"}] = struct{}{}
	}

	for _, o := range i.Settings.Overrides {
		if o.GoBasicType || o.GoTypeName == "" {
			continue
		}
		_, alreadyImported := std[o.GoImportPath]
		hasPackageAlias := o.GoPackage != ""
		if (!alreadyImported || hasPackageAlias) && i.usesType(o.GoTypeName) {
			pkg[ImportSpec{Path: o.GoImportPath, ID: o.GoPackage}] = struct{}{}
		}
	}

	pkgs := make([]ImportSpec, 0, len(pkg))
	for spec := range pkg {
		pkgs = append(pkgs, spec)
	}

	stds := make([]ImportSpec, 0, len(std))
	for path := range std {
		stds = append(stds, ImportSpec{Path: path})
	}

	sort.Slice(stds, func(i, j int) bool { return stds[i].Path < stds[j].Path })
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].Path < pkgs[j].Path })
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

	pkg := make(map[ImportSpec]struct{})
	std := map[string]struct{}{
		"context": {},
	}
	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}

	driver := DriverFromString(i.Settings.Go.Driver)

	for _, q := range gq {
		if q.Cmd == metadata.CmdExecResult {
			switch driver {
			case PgxDriver:
				pkg[ImportSpec{Path: "github.com/jackc/pgconn"}] = struct{}{}
			default:
				std["database/sql"] = struct{}{}
			}
		}
	}
	for typeName, pkg := range stdlibTypes {
		if uses(typeName) {
			std[pkg] = struct{}{}
		}
	}

	overrideTypes := map[string]string{}
	for _, o := range i.Settings.Overrides {
		if o.GoBasicType || o.GoTypeName == "" {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoImportPath
	}

	if sliceScan() && driver != PgxDriver {
		pkg[ImportSpec{Path: "github.com/lib/pq"}] = struct{}{}
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if uses("pq.NullTime") && !overrideNullTime {
		pkg[ImportSpec{Path: "github.com/lib/pq"}] = struct{}{}
	}
	_, overrideUUID := overrideTypes["uuid.UUID"]
	if uses("uuid.UUID") && !overrideUUID {
		pkg[ImportSpec{Path: "github.com/google/uuid"}] = struct{}{}
	}

	// Custom imports
	for _, o := range i.Settings.Overrides {
		if o.GoBasicType || o.GoTypeName == "" {
			continue
		}
		_, alreadyImported := std[o.GoImportPath]
		hasPackageAlias := o.GoPackage != ""
		if (!alreadyImported || hasPackageAlias) && uses(o.GoTypeName) {
			pkg[ImportSpec{Path: o.GoImportPath, ID: o.GoPackage}] = struct{}{}
		}
	}

	pkgs := make([]ImportSpec, 0, len(pkg))
	for spec := range pkg {
		pkgs = append(pkgs, spec)
	}

	stds := make([]ImportSpec, 0, len(std))
	for path := range std {
		stds = append(stds, ImportSpec{Path: path})
	}

	sort.Slice(stds, func(i, j int) bool { return stds[i].Path < stds[j].Path })
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].Path < pkgs[j].Path })
	return fileImports{stds, pkgs}
}
