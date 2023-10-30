package opts

import (
	"fmt"
	"os"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/pattern"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type PluginGoType struct {
	ImportPath string
	Package    string
	TypeName   string
	BasicType  bool
	StructTags map[string]string
}

func pluginGoType(o *Override) *PluginGoType {
	// Note that there is a slight mismatch between this and the
	// proto api. The GoType on the override is the unparsed type,
	// which could be a qualified path or an object, as per
	// https://docs.sqlc.dev/en/v1.18.0/reference/config.html#type-overriding
	return &PluginGoType{
		ImportPath: o.GoImportPath,
		Package:    o.GoPackage,
		TypeName:   o.GoTypeName,
		BasicType:  o.GoBasicType,
		StructTags: o.GoStructTags,
	}
}

type PluginOverride struct {
	DbType     string
	Nullable   bool
	Column     string
	Table      *plugin.Identifier
	ColumnName string
	Unsigned   bool
	GoType     *PluginGoType
}

func pluginOverride(req *plugin.CodeGenRequest, o *Override) *PluginOverride {
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
	return &PluginOverride{
		DbType:     o.DBType,
		Nullable:   o.Nullable,
		Unsigned:   o.Unsigned,
		Column:     o.Column,
		ColumnName: column,
		Table:      &table,
		GoType:     pluginGoType(o),
	}
}

type Override struct {
	// name of the golang type to use, e.g. `github.com/segmentio/ksuid.KSUID`
	GoType GoType `json:"go_type" yaml:"go_type"`

	// additional Go struct tags to add to this field, in raw Go struct tag form, e.g. `validate:"required" x:"y,z"`
	// see https://github.com/sqlc-dev/sqlc/issues/534
	GoStructTag GoStructTag `json:"go_struct_tag" yaml:"go_struct_tag"`

	// fully qualified name of the Go type, e.g. `github.com/segmentio/ksuid.KSUID`
	DBType                  string `json:"db_type" yaml:"db_type"`
	Deprecated_PostgresType string `json:"postgres_type" yaml:"postgres_type"`

	// for global overrides only when two different engines are in use
	Engine string `json:"engine,omitempty" yaml:"engine"`

	// True if the GoType should override if the matching type is nullable
	Nullable bool `json:"nullable" yaml:"nullable"`

	// True if the GoType should override if the matching type is unsiged.
	Unsigned bool `json:"unsigned" yaml:"unsigned"`

	// Deprecated. Use the `nullable` property instead
	Deprecated_Null bool `json:"null" yaml:"null"`

	// fully qualified name of the column, e.g. `accounts.id`
	Column string `json:"column" yaml:"column"`

	ColumnName   *pattern.Match `json:"-"`
	TableCatalog *pattern.Match `json:"-"`
	TableSchema  *pattern.Match `json:"-"`
	TableRel     *pattern.Match `json:"-"`
	GoImportPath string         `json:"-"`
	GoPackage    string         `json:"-"`
	GoTypeName   string         `json:"-"`
	GoBasicType  bool           `json:"-"`

	// Parsed form of GoStructTag, e.g. {"validate:", "required"}
	GoStructTags map[string]string `json:"-"`
	Plugin       *PluginOverride   `json:"-"`
}

func (o *Override) Matches(n *plugin.Identifier, defaultSchema string) bool {
	if n == nil {
		return false
	}
	schema := n.Schema
	if n.Schema == "" {
		schema = defaultSchema
	}
	if o.TableCatalog != nil && !o.TableCatalog.MatchString(n.Catalog) {
		return false
	}
	if o.TableSchema == nil && schema != "" {
		return false
	}
	if o.TableSchema != nil && !o.TableSchema.MatchString(schema) {
		return false
	}
	if o.TableRel == nil && n.Name != "" {
		return false
	}
	if o.TableRel != nil && !o.TableRel.MatchString(n.Name) {
		return false
	}
	return true
}

func (o *Override) parse(req *plugin.CodeGenRequest) (err error) {
	// validate deprecated postgres_type field
	if o.Deprecated_PostgresType != "" {
		fmt.Fprintf(os.Stderr, "WARNING: \"postgres_type\" is deprecated. Instead, use \"db_type\" to specify a type override.\n")
		if o.DBType != "" {
			return fmt.Errorf(`Type override configurations cannot have "db_type" and "postres_type" together. Use "db_type" alone`)
		}
		o.DBType = o.Deprecated_PostgresType
	}

	// validate deprecated null field
	if o.Deprecated_Null {
		fmt.Fprintf(os.Stderr, "WARNING: \"null\" is deprecated. Instead, use the \"nullable\" field.\n")
		o.Nullable = true
	}

	schema := "public"
	if req != nil && req.Catalog != nil {
		schema = req.Catalog.DefaultSchema
	}

	// validate option combinations
	switch {
	case o.Column != "" && o.DBType != "":
		return fmt.Errorf("Override specifying both `column` (%q) and `db_type` (%q) is not valid.", o.Column, o.DBType)
	case o.Column == "" && o.DBType == "":
		return fmt.Errorf("Override must specify one of either `column` or `db_type`")
	}

	// validate Column
	if o.Column != "" {
		colParts := strings.Split(o.Column, ".")
		switch len(colParts) {
		case 2:
			if o.ColumnName, err = pattern.MatchCompile(colParts[1]); err != nil {
				return err
			}
			if o.TableRel, err = pattern.MatchCompile(colParts[0]); err != nil {
				return err
			}
			if o.TableSchema, err = pattern.MatchCompile(schema); err != nil {
				return err
			}
		case 3:
			if o.ColumnName, err = pattern.MatchCompile(colParts[2]); err != nil {
				return err
			}
			if o.TableRel, err = pattern.MatchCompile(colParts[1]); err != nil {
				return err
			}
			if o.TableSchema, err = pattern.MatchCompile(colParts[0]); err != nil {
				return err
			}
		case 4:
			if o.ColumnName, err = pattern.MatchCompile(colParts[3]); err != nil {
				return err
			}
			if o.TableRel, err = pattern.MatchCompile(colParts[2]); err != nil {
				return err
			}
			if o.TableSchema, err = pattern.MatchCompile(colParts[1]); err != nil {
				return err
			}
			if o.TableCatalog, err = pattern.MatchCompile(colParts[0]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Override `column` specifier %q is not the proper format, expected '[catalog.][schema.]tablename.colname'", o.Column)
		}
	}

	// validate GoType
	parsed, err := o.GoType.parse()
	if err != nil {
		return err
	}
	o.GoImportPath = parsed.ImportPath
	o.GoPackage = parsed.Package
	o.GoTypeName = parsed.TypeName
	o.GoBasicType = parsed.BasicType

	// validate GoStructTag
	tags, err := o.GoStructTag.parse()
	if err != nil {
		return err
	}
	o.GoStructTags = tags

	return nil
}
