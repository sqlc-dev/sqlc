package golang

import (
	"github.com/kyleconroy/sqlc/internal/plugin"
)

type SQLDriver int

const (
	SQLDriverPGXV4 SQLDriver = iota
	SQLDriverLibPQ
)

func parseDriver(settings *plugin.Settings) SQLDriver {
	if settings.Go.SqlPackage == "pgx/v4" {
		return SQLDriverPGXV4
	} else {
		return SQLDriverLibPQ
	}
}
