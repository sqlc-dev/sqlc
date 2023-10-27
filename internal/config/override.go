package config

import (
	"fmt"
	"os"
	"strings"

	gopluginopts "github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/pattern"
	"gopkg.in/yaml.v3"
)

type Override struct {
	// name of the golang type to use, e.g. `github.com/segmentio/ksuid.KSUID`
	GoType gopluginopts.GoType `json:"go_type" yaml:"go_type"`

	// additional Go struct tags to add to this field, in raw Go struct tag form, e.g. `validate:"required" x:"y,z"`
	// see https://github.com/sqlc-dev/sqlc/issues/534
	GoStructTag gopluginopts.GoStructTag `json:"go_struct_tag" yaml:"go_struct_tag"`

	// fully qualified name of the Go type, e.g. `github.com/segmentio/ksuid.KSUID`
	DBType                  string `json:"db_type" yaml:"db_type"`
	Deprecated_PostgresType string `json:"postgres_type" yaml:"postgres_type"`

	// for global overrides only when two different engines are in use
	Engine Engine `json:"engine,omitempty" yaml:"engine"`

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

	// For passing plugin-specific configuration
	Plugin  string    `json:"plugin,omitempty"`
	Options yaml.Node `json:"options,omitempty"`
}

func (o Override) hasGoOptions() bool {
	hasGoTypePath := o.GoType.Path != ""
	hasGoTypePackage := o.GoType.Package != ""
	hasGoTypeName := o.GoType.Name != ""
	hasGoType := hasGoTypePath || hasGoTypePackage || hasGoTypeName
	hasGoStructTag := o.GoStructTag != ""
	return hasGoType || hasGoStructTag
}

func (o *Override) Parse() (err error) {
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

	// validate option combinations
	switch {
	case o.Column != "" && o.DBType != "":
		return fmt.Errorf("Override specifying both `column` (%q) and `db_type` (%q) is not valid.", o.Column, o.DBType)
	case o.Column == "" && o.DBType == "":
		return fmt.Errorf("Override must specify one of either `column` or `db_type`")
	case o.hasGoOptions() && !o.Options.IsZero():
		return fmt.Errorf("Override can specify go_type/go_struct_tag or options but not both")
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
			if o.TableSchema, err = pattern.MatchCompile("public"); err != nil {
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

	return nil
}
