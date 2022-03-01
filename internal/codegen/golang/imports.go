package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/plugin"
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
	Settings *plugin.Settings
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

func (i *importer) Imports(filename string) [][]ImportSpec {
	dbFileName := "db.go"
	if i.Settings.Go.OutputDbFileName != "" {
		dbFileName = i.Settings.Go.OutputDbFileName
	}
	modelsFileName := "models.go"
	if i.Settings.Go.OutputModelsFileName != "" {
		modelsFileName = i.Settings.Go.OutputModelsFileName
	}
	querierFileName := "querier.go"
	if i.Settings.Go.OutputQuerierFileName != "" {
		querierFileName = i.Settings.Go.OutputQuerierFileName
	}
	copyfromFileName := "copyfrom.go"
	batchFileName := "batch.go"

	switch filename {
	case dbFileName:
		return mergeImports(i.dbImports())
	case modelsFileName:
		return mergeImports(i.modelImports())
	case querierFileName:
		return mergeImports(i.interfaceImports())
	case copyfromFileName:
		return mergeImports(i.copyfromImports())
	case batchFileName:
		return mergeImports(i.batchImports(filename))
	default:
		return mergeImports(i.queryImports(filename))
	}
}

func (i *importer) dbImports() fileImports {
	var pkg []ImportSpec
	std := []ImportSpec{
		{Path: "context"},
	}

	sqlpkg := SQLPackageFromString(i.Settings.Go.SqlPackage)
	switch sqlpkg {
	case SQLPackagePGX:
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

var pgtypeTypes = map[string]struct{}{
	"pgtype.CIDR":      {},
	"pgtype.Daterange": {},
	"pgtype.Inet":      {},
	"pgtype.Int4range": {},
	"pgtype.Int8range": {},
	"pgtype.JSON":      {},
	"pgtype.JSONB":     {},
	"pgtype.Hstore":    {},
	"pgtype.Macaddr":   {},
	"pgtype.Numeric":   {},
	"pgtype.Numrange":  {},
	"pgtype.Tsrange":   {},
	"pgtype.Tstzrange": {},
}

var pqtypeTypes = map[string]struct{}{
	"pqtype.CIDR":           {},
	"pqtype.Inet":           {},
	"pqtype.Macaddr":        {},
	"pqtype.NullRawMessage": {},
}

func buildImports(settings *plugin.Settings, queries []Query, uses func(string) bool) (map[string]struct{}, map[ImportSpec]struct{}) {
	pkg := make(map[ImportSpec]struct{})
	std := make(map[string]struct{})

	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}

	sqlpkg := SQLPackageFromString(settings.Go.SqlPackage)
	for _, q := range queries {
		if q.Cmd == metadata.CmdExecResult {
			switch sqlpkg {
			case SQLPackagePGX:
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

	for typeName, _ := range pgtypeTypes {
		if uses(typeName) {
			pkg[ImportSpec{Path: "github.com/jackc/pgtype"}] = struct{}{}
		}
	}

	for typeName, _ := range pqtypeTypes {
		if uses(typeName) {
			pkg[ImportSpec{Path: "github.com/tabbed/pqtype"}] = struct{}{}
		}
	}

	overrideTypes := map[string]string{}
	for _, o := range settings.Overrides {
		if o.GoType.BasicType || o.GoType.TypeName == "" {
			continue
		}
		overrideTypes[o.GoType.TypeName] = o.GoType.ImportPath
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if uses("pq.NullTime") && !overrideNullTime {
		pkg[ImportSpec{Path: "github.com/lib/pq"}] = struct{}{}
	}
	_, overrideUUID := overrideTypes["uuid.UUID"]
	if uses("uuid.UUID") && !overrideUUID {
		pkg[ImportSpec{Path: "github.com/google/uuid"}] = struct{}{}
	}
	_, overrideNullUUID := overrideTypes["uuid.NullUUID"]
	if uses("uuid.NullUUID") && !overrideNullUUID {
		pkg[ImportSpec{Path: "github.com/google/uuid"}] = struct{}{}
	}

	// Custom imports
	for _, o := range settings.Overrides {
		if o.GoType.BasicType || o.GoType.TypeName == "" {
			continue
		}
		_, alreadyImported := std[o.GoType.ImportPath]
		hasPackageAlias := o.GoType.Package != ""
		if (!alreadyImported || hasPackageAlias) && uses(o.GoType.TypeName) {
			pkg[ImportSpec{Path: o.GoType.ImportPath, ID: o.GoType.Package}] = struct{}{}
		}
	}

	return std, pkg
}

func (i *importer) interfaceImports() fileImports {
	std, pkg := buildImports(i.Settings, i.Queries, func(name string) bool {
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
	})

	std["context"] = struct{}{}

	return sortedImports(std, pkg)
}

func (i *importer) modelImports() fileImports {
	std, pkg := buildImports(i.Settings, nil, func(prefix string) bool {
		return i.usesType(prefix)
	})

	if len(i.Enums) > 0 {
		std["fmt"] = struct{}{}
	}

	return sortedImports(std, pkg)
}

func sortedImports(std map[string]struct{}, pkg map[ImportSpec]struct{}) fileImports {
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
	anyNonCopyFrom := false
	for _, query := range i.Queries {
		if usesBatch([]Query{query}) {
			continue
		}
		if query.SourceName == filename {
			gq = append(gq, query)
			if query.Cmd != metadata.CmdCopyFrom {
				anyNonCopyFrom = true
			}
		}
	}

	std, pkg := buildImports(i.Settings, gq, func(name string) bool {
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
	})

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

	if anyNonCopyFrom {
		std["context"] = struct{}{}
	}

	sqlpkg := SQLPackageFromString(i.Settings.Go.SqlPackage)
	if sliceScan() && sqlpkg != SQLPackagePGX {
		pkg[ImportSpec{Path: "github.com/lib/pq"}] = struct{}{}
	}

	return sortedImports(std, pkg)
}

func (i *importer) copyfromImports() fileImports {
	std, pkg := buildImports(i.Settings, i.Queries, func(name string) bool {
		for _, q := range i.Queries {
			if q.Cmd != metadata.CmdCopyFrom {
				continue
			}
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
	})

	std["context"] = struct{}{}

	return sortedImports(std, pkg)
}

func (i *importer) batchImports(filename string) fileImports {
	std, pkg := buildImports(i.Settings, i.Queries, func(name string) bool {
		for _, q := range i.Queries {
			if !usesBatch([]Query{q}) {
				continue
			}
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
	})

	std["context"] = struct{}{}
	pkg[ImportSpec{Path: "github.com/jackc/pgx/v4"}] = struct{}{}

	return sortedImports(std, pkg)
}
