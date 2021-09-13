package golang

import "github.com/kyleconroy/sqlc/internal/config"

type SQLDriver int

const (
	SQLDriverPGXV4 SQLDriver = iota
	SQLDriverLibPQ
)

func parseDriver(settings config.CombinedSettings) SQLDriver {
	if settings.Go.SQLPackage == "pgx/v4" {
		return SQLDriverPGXV4
	} else {
		return SQLDriverLibPQ
	}
}
