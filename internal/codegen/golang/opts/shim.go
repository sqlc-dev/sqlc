package opts

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// The ShimOverride struct exists to bridge the gap between the Override struct
// and the previous Override struct defined in codegen.proto. Eventually these
// shim structs should be removed in favor of using the existing Override and
// GoType structs, but it's easier to provide these shim structs to not change
// the existing, working code.
type ShimOverride struct {
	DbType     string
	Nullable   bool
	Column     string
	Table      *plugin.Identifier
	ColumnName string
	Unsigned   bool
	GoType     *ShimGoType
}

func shimOverride(req *plugin.GenerateRequest, o *Override) *ShimOverride {
	var column string
	var table plugin.Identifier

	if o.Column != "" {
		colParts := strings.Split(o.Column, ".")
		switch len(colParts) {
		case 2:
			table.Schema = req.Catalog.DefaultSchema
			table.Name = colParts[0]
			column = colParts[1]
		case 3:
			table.Schema = colParts[0]
			table.Name = colParts[1]
			column = colParts[2]
		case 4:
			table.Catalog = colParts[0]
			table.Schema = colParts[1]
			table.Name = colParts[2]
			column = colParts[3]
		}
	}
	return &ShimOverride{
		DbType:     o.DBType,
		Nullable:   o.Nullable,
		Unsigned:   o.Unsigned,
		Column:     o.Column,
		ColumnName: column,
		Table:      &table,
		GoType:     shimGoType(o),
	}
}

type ShimGoType struct {
	ImportPath string
	Package    string
	TypeName   string
	BasicType  bool
	StructTags map[string]string
}

func shimGoType(o *Override) *ShimGoType {
	// Note that there is a slight mismatch between this and the
	// proto api. The GoType on the override is the unparsed type,
	// which could be a qualified path or an object, as per
	// https://docs.sqlc.dev/en/v1.18.0/reference/config.html#type-overriding
	return &ShimGoType{
		ImportPath: o.GoImportPath,
		Package:    o.GoPackage,
		TypeName:   o.GoTypeName,
		BasicType:  o.GoBasicType,
		StructTags: o.GoStructTags,
	}
}
