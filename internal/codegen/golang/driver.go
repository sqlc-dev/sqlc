package golang

import "github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"

func parseDriver(sqlPackage string) opts.SQLDriver {
	switch sqlPackage {
	case opts.SQLPackagePGXV4:
		return opts.SQLDriverPGXV4
	case opts.SQLPackagePGXV5:
		return opts.SQLDriverPGXV5
	case opts.SQLPackageYugaBytePGXV5:
		return opts.SQLDriverYugaBytePGXV5
	default:
		return opts.SQLDriverLibPQ
	}
}

// custom packages based on pgx/v5 (e.g., YugabyteDB smart drivers) should return
// only the original driver when determining columnType, refer to the postgresType func.
func parseDriverPGType(sqlPackage string) opts.SQLDriver {
	switch sqlPackage {
	case opts.SQLPackagePGXV4:
		return opts.SQLDriverPGXV4
	case opts.SQLPackagePGXV5, opts.SQLPackageYugaBytePGXV5:
		return opts.SQLDriverPGXV5
	default:
		return opts.SQLDriverLibPQ
	}
}
