package opts

import (
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type GlobalOverride struct {
	*plugin.Override

	GoType *ParsedGoType
}

func (o *GlobalOverride) Convert() *plugin.Override {
	return &plugin.Override{
		DbType:     o.DbType,
		Nullable:   o.Nullable,
		Column:     o.Column,
		Table:      o.Table,
		ColumnName: o.ColumnName,
		Unsigned:   o.Unsigned,
	}
}

func (o *GlobalOverride) Matches(n *plugin.Identifier, defaultSchema string) bool {
	return sdk.Matches(o.Convert(), n, defaultSchema)
}
