// Extracted from github.com/sqlc-dev/sqlc/internal/cmd/shim.go
package opts

import (
	"encoding/json"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func pluginOverride(defaultSchema string, o Override) *plugin.Override {
	var column string
	var table plugin.Identifier

	if o.Column != "" {
		colParts := strings.Split(o.Column, ".")
		switch len(colParts) {
		case 2:
			table.Schema = defaultSchema
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

	goTypeJSON, err := json.Marshal(pluginGoType(o))
	if err != nil {
		panic(err)
	}

	return &plugin.Override{
		CodeType:   goTypeJSON,
		DbType:     o.DBType,
		Nullable:   o.Nullable,
		Unsigned:   o.Unsigned,
		Column:     o.Column,
		ColumnName: column,
		Table:      &table,
	}
}

func pluginGoType(o Override) *ParsedGoType {
	// Note that there is a slight mismatch between this and the
	// proto api. The GoType on the override is the unparsed type,
	// which could be a qualified path or an object, as per
	// https://docs.sqlc.dev/en/v1.18.0/reference/config.html#type-overriding
	return &ParsedGoType{
		ImportPath: o.GoImportPath,
		Package:    o.GoPackage,
		TypeName:   o.GoTypeName,
		BasicType:  o.GoBasicType,
		StructTags: o.GoStructTags,
	}
}
