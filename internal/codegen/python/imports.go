package python

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
)

type importSpec struct {
	Module string
	Name   string
	Alias  string
}

func (i importSpec) String() string {
	if i.Alias != "" {
		if i.Name == "" {
			return fmt.Sprintf("import %s as %s", i.Module, i.Alias)
		}
		return fmt.Sprintf("from %s import %s as %s", i.Module, i.Name, i.Alias)
	}
	if i.Name == "" {
		return "import " + i.Module
	}
	return fmt.Sprintf("from %s import %s", i.Module, i.Name)
}

type importer struct {
	Settings config.CombinedSettings
	Models   []Struct
	Queries  []Query
	Enums    []Enum
}

func structUses(name string, s Struct) bool {
	for _, f := range s.Fields {
		if name == "typing.List" && f.Type.IsArray {
			return true
		}
		if name == "typing.Optional" && f.Type.IsNull {
			return true
		}
		if f.Type.InnerType == name {
			return true
		}
	}
	return false
}

func queryValueUses(name string, qv QueryValue) bool {
	if !qv.isEmpty() {
		if name == "typing.List" && qv.Typ.IsArray {
			return true
		}
		if name == "typing.Optional" && qv.Typ.IsNull {
			return true
		}
		if qv.IsStruct() && qv.EmitStruct() {
			if structUses(name, *qv.Struct) {
				return true
			}
		} else {
			if qv.Typ.InnerType == name {
				return true
			}
		}
	}
	return false
}

func (i *importer) Imports(fileName string) []string {
	if fileName == "models.py" {
		return i.modelImports()
	}
	return i.queryImports(fileName)
}

func (i *importer) modelImportSpecs() (map[string]importSpec, map[string]importSpec) {
	modelUses := func(name string) bool {
		for _, model := range i.Models {
			if structUses(name, model) {
				return true
			}
		}
		return false
	}

	std := stdImports(modelUses)
	std["dataclasses"] = importSpec{Module: "dataclasses"}
	if len(i.Enums) > 0 {
		std["enum"] = importSpec{Module: "enum"}
	}

	pkg := make(map[string]importSpec)

	for _, o := range i.Settings.Overrides {
		if o.PythonType.IsSet() && o.PythonType.Module != "" {
			if modelUses(o.PythonType.TypeString()) {
				pkg[o.PythonType.Module] = importSpec{Module: o.PythonType.Module}
			}
		}
	}
	return std, pkg
}

func (i *importer) modelImports() []string {
	std, pkg := i.modelImportSpecs()
	importLines := []string{
		buildImportBlock(std),
		"",
		buildImportBlock(pkg),
	}
	return importLines
}

func (i *importer) queryImportSpecs(fileName string) (map[string]importSpec, map[string]importSpec) {
	queryUses := func(name string) bool {
		for _, q := range i.Queries {
			if q.SourceName != fileName {
				continue
			}
			if queryValueUses(name, q.Ret) {
				return true
			}
			for _, arg := range q.Args {
				if queryValueUses(name, arg) {
					return true
				}
			}
		}
		return false
	}

	std := stdImports(queryUses)

	pkg := make(map[string]importSpec)
	pkg["sqlalchemy"] = importSpec{Module: "sqlalchemy"}
	if i.Settings.Python.EmitAsyncQuerier {
		pkg["sqlalchemy.ext.asyncio"] = importSpec{Module: "sqlalchemy.ext.asyncio"}
	}

	for _, o := range i.Settings.Overrides {
		if o.PythonType.IsSet() && o.PythonType.Module != "" {
			if queryUses(o.PythonType.TypeString()) {
				pkg[o.PythonType.Module] = importSpec{Module: o.PythonType.Module}
			}
		}
	}

	queryValueModelImports := func(qv QueryValue) {
		if qv.IsStruct() && qv.EmitStruct() {
			std["dataclasses"] = importSpec{Module: "dataclasses"}
		}
	}

	for _, q := range i.Queries {
		if q.SourceName != fileName {
			continue
		}
		if q.Cmd == ":one" {
			std["typing.Optional"] = importSpec{Module: "typing", Name: "Optional"}
		}
		if q.Cmd == ":many" {
			if i.Settings.Python.EmitSyncQuerier {
				std["typing.Iterator"] = importSpec{Module: "typing", Name: "Iterator"}
			}
			if i.Settings.Python.EmitAsyncQuerier {
				std["typing.AsyncIterator"] = importSpec{Module: "typing", Name: "AsyncIterator"}
			}
		}
		queryValueModelImports(q.Ret)
		for _, qv := range q.Args {
			queryValueModelImports(qv)
		}
	}

	return std, pkg
}

func (i *importer) queryImports(fileName string) []string {
	std, pkg := i.queryImportSpecs(fileName)

	modelImportStr := fmt.Sprintf("from %s import models", i.Settings.Python.Package)

	importLines := []string{
		buildImportBlock(std),
		"",
		buildImportBlock(pkg),
		"",
		modelImportStr,
	}
	return importLines
}

type importFromSpec struct {
	Module string
	Names  []string
	Alias  string
}

func buildImportBlock2(pkgs map[string]importSpec) []importFromSpec {
	pkgImports := make([]importFromSpec, 0)
	fromImports := make(map[string][]string)
	for _, is := range pkgs {
		if is.Name == "" || is.Alias != "" {
			pkgImports = append(pkgImports, importFromSpec{
				Module: is.Module,
				Names:  []string{is.Name},
				Alias:  is.Alias,
			})
		} else {
			names, ok := fromImports[is.Module]
			if !ok {
				names = make([]string, 0, 1)
			}
			names = append(names, is.Name)
			fromImports[is.Module] = names
		}
	}
	for modName, names := range fromImports {
		sort.Strings(names)
		pkgImports = append(pkgImports, importFromSpec{
			Module: modName,
			Names:  names,
		})
	}
	sort.Slice(pkgImports, func(i, j int) bool {
		return pkgImports[i].Module < pkgImports[j].Module || pkgImports[i].Names[0] < pkgImports[j].Names[0]
	})
	return pkgImports
}

func buildImportBlock(pkgs map[string]importSpec) string {
	pkgImports := make([]importSpec, 0)
	fromImports := make(map[string][]string)
	for _, is := range pkgs {
		if is.Name == "" || is.Alias != "" {
			pkgImports = append(pkgImports, is)
		} else {
			names, ok := fromImports[is.Module]
			if !ok {
				names = make([]string, 0, 1)
			}
			names = append(names, is.Name)
			fromImports[is.Module] = names
		}
	}

	importStrings := make([]string, 0, len(pkgImports)+len(fromImports))
	for _, is := range pkgImports {
		importStrings = append(importStrings, is.String())
	}
	for modName, names := range fromImports {
		sort.Strings(names)
		nameString := strings.Join(names, ", ")
		importStrings = append(importStrings, fmt.Sprintf("from %s import %s", modName, nameString))
	}
	sort.Strings(importStrings)
	return strings.Join(importStrings, "\n")
}

func stdImports(uses func(name string) bool) map[string]importSpec {
	std := make(map[string]importSpec)
	if uses("decimal.Decimal") {
		std["decimal"] = importSpec{Module: "decimal"}
	}
	if uses("datetime.date") || uses("datetime.time") || uses("datetime.datetime") || uses("datetime.timedelta") {
		std["datetime"] = importSpec{Module: "datetime"}
	}
	if uses("uuid.UUID") {
		std["uuid"] = importSpec{Module: "uuid"}
	}
	if uses("typing.List") {
		std["typing.List"] = importSpec{Module: "typing", Name: "List"}
	}
	if uses("typing.Optional") {
		std["typing.Optional"] = importSpec{Module: "typing", Name: "Optional"}
	}
	if uses("Any") {
		std["typing.Any"] = importSpec{Module: "typing", Name: "Any"}
	}
	return std
}
