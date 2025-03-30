package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/metadata"
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
		return fmt.Sprintf("%s %q", s.ID, s.Path)
	} else {
		return fmt.Sprintf("%q", s.Path)
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
	Options *opts.Options
	Queries []Query
	Enums   []Enum
	Structs []Struct
}

func (i *importer) usesType(typ string) bool {
	for _, strct := range i.Structs {
		for _, f := range strct.Fields {
			if hasPrefixIgnoringSliceAndPointerPrefix(f.Type, typ) {
				return true
			}
		}
	}
	return false
}

func (i *importer) HasImports(filename string) bool {
	imports := i.Imports(filename)
	return len(imports[0]) != 0 || len(imports[1]) != 0
}

func (i *importer) Imports(filename string) [][]ImportSpec {
	dbFileName := "db.go"
	if i.Options.OutputDbFileName != "" {
		dbFileName = i.Options.OutputDbFileName
	}
	modelsFileName := "models.go"
	if i.Options.OutputModelsFileName != "" {
		modelsFileName = i.Options.OutputModelsFileName
	}
	querierFileName := "querier.go"
	if i.Options.OutputQuerierFileName != "" {
		querierFileName = i.Options.OutputQuerierFileName
	}
	copyfromFileName := "copyfrom.go"
	if i.Options.OutputCopyfromFileName != "" {
		copyfromFileName = i.Options.OutputCopyfromFileName
	}
	batchFileName := "batch.go"
	if i.Options.OutputBatchFileName != "" {
		batchFileName = i.Options.OutputBatchFileName
	}

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
		return mergeImports(i.batchImports())
	default:
		return mergeImports(i.queryImports(filename))
	}
}

func (i *importer) dbImports() fileImports {
	var pkg []ImportSpec
	std := []ImportSpec{
		{Path: "context"},
	}

	sqlpkg := parseDriver(i.Options.SqlPackage)
	switch sqlpkg {
	case opts.SQLDriverPGXV4:
		pkg = append(pkg, ImportSpec{Path: "github.com/jackc/pgconn"})
		pkg = append(pkg, ImportSpec{Path: "github.com/jackc/pgx/v4"})
	case opts.SQLDriverPGXV5:
		pkg = append(pkg, ImportSpec{Path: "github.com/jackc/pgx/v5/pgconn"})
		pkg = append(pkg, ImportSpec{Path: "github.com/jackc/pgx/v5"})
	default:
		std = append(std, ImportSpec{Path: "database/sql"})
		if i.Options.EmitPreparedQueries {
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
	"netip.Addr":       "net/netip",
	"netip.Prefix":     "net/netip",
}

var pqtypeTypes = map[string]struct{}{
	"pqtype.CIDR":           {},
	"pqtype.Inet":           {},
	"pqtype.Macaddr":        {},
	"pqtype.NullRawMessage": {},
}

func buildImports(options *opts.Options, queries []Query, uses func(string) bool) (map[string]struct{}, map[ImportSpec]struct{}) {
	pkg := make(map[ImportSpec]struct{})
	std := make(map[string]struct{})

	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}

	sqlpkg := parseDriver(options.SqlPackage)
	for _, q := range queries {
		if q.Cmd == metadata.CmdExecResult {
			switch sqlpkg {
			case opts.SQLDriverPGXV4:
				pkg[ImportSpec{Path: "github.com/jackc/pgconn"}] = struct{}{}
			case opts.SQLDriverPGXV5:
				pkg[ImportSpec{Path: "github.com/jackc/pgx/v5/pgconn"}] = struct{}{}
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

	if uses("pgtype.") {
		if sqlpkg == opts.SQLDriverPGXV5 {
			pkg[ImportSpec{Path: "github.com/jackc/pgx/v5/pgtype"}] = struct{}{}
		} else {
			pkg[ImportSpec{Path: "github.com/jackc/pgtype"}] = struct{}{}
		}
	}

	for typeName := range pqtypeTypes {
		if uses(typeName) {
			pkg[ImportSpec{Path: "github.com/sqlc-dev/pqtype"}] = struct{}{}
			break
		}
	}

	overrideTypes := map[string]string{}
	for _, override := range options.Overrides {
		o := override.ShimOverride
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
	_, overrideVector := overrideTypes["pgvector.Vector"]
	if uses("pgvector.Vector") && !overrideVector {
		pkg[ImportSpec{Path: "github.com/pgvector/pgvector-go"}] = struct{}{}
	}

	// Custom imports
	for _, override := range options.Overrides {
		o := override.ShimOverride

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
	std, pkg := buildImports(i.Options, i.Queries, func(name string) bool {
		for _, q := range i.Queries {
			if q.hasRetType() {
				if usesBatch([]Query{q}) {
					continue
				}
				if hasPrefixIgnoringSliceAndPointerPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			for _, f := range q.Arg.Pairs() {
				if hasPrefixIgnoringSliceAndPointerPrefix(f.Type, name) {
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
	std, pkg := buildImports(i.Options, nil, i.usesType)

	if len(i.Enums) > 0 {
		std["fmt"] = struct{}{}
		std["database/sql/driver"] = struct{}{}
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

	std, pkg := buildImports(i.Options, gq, func(name string) bool {
		for _, q := range gq {
			if q.hasRetType() {
				if q.Ret.EmitStruct() {
					for _, f := range q.Ret.Struct.Fields {
						if hasPrefixIgnoringSliceAndPointerPrefix(f.Type, name) {
							return true
						}
					}
				}
				if hasPrefixIgnoringSliceAndPointerPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			// Check the fields of the argument struct if it's emitted
			if q.Arg.EmitStruct() {
				for _, f := range q.Arg.Struct.Fields {
					if hasPrefixIgnoringSliceAndPointerPrefix(f.Type, name) {
						return true
					}
				}
			}
			// Check the argument pairs inside the method definition
			for _, f := range q.Arg.Pairs() {
				if hasPrefixIgnoringSliceAndPointerPrefix(f.Type, name) {
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
						for _, embed := range f.EmbedFields {
							if strings.HasPrefix(embed.Type, "[]") && embed.Type != "[]byte" {
								return true
							}
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
						if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" && !f.HasSqlcSlice() {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Arg.Type(), "[]") && q.Arg.Type() != "[]byte" && !q.Arg.HasSqlcSlices() {
						return true
					}
				}
			}
		}
		return false
	}

	// Search for sqlc.slice() calls
	sqlcSliceScan := func() bool {
		for _, q := range gq {
			if q.Arg.HasSqlcSlices() {
				return true
			}
		}
		return false
	}

	if anyNonCopyFrom {
		std["context"] = struct{}{}
	}

	sqlpkg := parseDriver(i.Options.SqlPackage)
	if sqlcSliceScan() && !sqlpkg.IsPGX() {
		std["strings"] = struct{}{}
	}
	if sliceScan() && !sqlpkg.IsPGX() {
		pkg[ImportSpec{Path: "github.com/lib/pq"}] = struct{}{}
	}

	if i.Options.WrapErrors {
		std["fmt"] = struct{}{}
	}

	return sortedImports(std, pkg)
}

func (i *importer) copyfromImports() fileImports {
	copyFromQueries := make([]Query, 0, len(i.Queries))
	for _, q := range i.Queries {
		if q.Cmd == metadata.CmdCopyFrom {
			copyFromQueries = append(copyFromQueries, q)
		}
	}
	std, pkg := buildImports(i.Options, copyFromQueries, func(name string) bool {
		for _, q := range copyFromQueries {
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
	if i.Options.SqlDriver == opts.SQLDriverGoSQLDriverMySQL {
		std["io"] = struct{}{}
		std["fmt"] = struct{}{}
		std["sync/atomic"] = struct{}{}
		pkg[ImportSpec{Path: "github.com/go-sql-driver/mysql"}] = struct{}{}
		pkg[ImportSpec{Path: "github.com/hexon/mysqltsv"}] = struct{}{}
	}

	return sortedImports(std, pkg)
}

func (i *importer) batchImports() fileImports {
	batchQueries := make([]Query, 0, len(i.Queries))
	for _, q := range i.Queries {
		if usesBatch([]Query{q}) {
			batchQueries = append(batchQueries, q)
		}
	}
	std, pkg := buildImports(i.Options, batchQueries, func(name string) bool {
		for _, q := range batchQueries {
			if q.hasRetType() {
				if q.Ret.EmitStruct() {
					for _, f := range q.Ret.Struct.Fields {
						if hasPrefixIgnoringSliceAndPointerPrefix(f.Type, name) {
							return true
						}
					}
				}
				if hasPrefixIgnoringSliceAndPointerPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if q.Arg.EmitStruct() {
				for _, f := range q.Arg.Struct.Fields {
					if hasPrefixIgnoringSliceAndPointerPrefix(f.Type, name) {
						return true
					}
				}
			}
			for _, f := range q.Arg.Pairs() {
				if hasPrefixIgnoringSliceAndPointerPrefix(f.Type, name) {
					return true
				}
			}
		}
		return false
	})

	std["context"] = struct{}{}
	std["errors"] = struct{}{}
	sqlpkg := parseDriver(i.Options.SqlPackage)
	switch sqlpkg {
	case opts.SQLDriverPGXV4:
		pkg[ImportSpec{Path: "github.com/jackc/pgx/v4"}] = struct{}{}
	case opts.SQLDriverPGXV5:
		pkg[ImportSpec{Path: "github.com/jackc/pgx/v5"}] = struct{}{}
	}

	return sortedImports(std, pkg)
}

func trimSliceAndPointerPrefix(v string) string {
	v = strings.TrimPrefix(v, "[]")
	v = strings.TrimPrefix(v, "*")
	return v
}

func hasPrefixIgnoringSliceAndPointerPrefix(s, prefix string) bool {
	trimmedS := trimSliceAndPointerPrefix(s)
	trimmedPrefix := trimSliceAndPointerPrefix(prefix)
	return strings.HasPrefix(trimmedS, trimmedPrefix)
}

func replaceConflictedArg(imports [][]ImportSpec, queries []Query) []Query {
	m := make(map[string]struct{})
	for _, is := range imports {
		for _, i := range is {
			paths := strings.Split(i.Path, "/")
			m[paths[len(paths)-1]] = struct{}{}
		}
	}

	replacedQueries := make([]Query, 0, len(queries))
	for _, query := range queries {
		if _, exist := m[query.Arg.Name]; exist {
			query.Arg.Name = toCamelCase(fmt.Sprintf("arg_%s", query.Arg.Name))
		}
		replacedQueries = append(replacedQueries, query)
	}
	return replacedQueries
}
