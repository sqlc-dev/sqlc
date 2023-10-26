package opts

import (
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type GoOverride struct {
	*plugin.Override

	GoType *ParsedGoType
}

func (o *GoOverride) Convert() *plugin.Override {
	return &plugin.Override{
		DbType:     o.DbType,
		Nullable:   o.Nullable,
		Column:     o.Column,
		Table:      o.Table,
		ColumnName: o.ColumnName,
		Unsigned:   o.Unsigned,
	}
}

func (o *GoOverride) Matches(n *plugin.Identifier, defaultSchema string) bool {
	return sdk.Matches(o.Convert(), n, defaultSchema)
}
