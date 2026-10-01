package golang

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
)

func parseDriver(sqlPackage string) opts.SQLDriver {
	switch sqlPackage {
	case opts.SQLPackagePGXV4:
		return opts.SQLDriverPGXV4
	case opts.SQLPackagePGXV5:
		return opts.SQLDriverPGXV5
	default:
		return opts.SQLDriverLibPQ
	}
}

// The engines, as the request names them.
const (
	engineClickHouse = "clickhouse"
	engineDuckDB     = "duckdb"
	engineGoogleSQL  = "googlesql"
	engineMSSQL      = "mssql"
	engineMySQL      = "mysql"
	enginePostgreSQL = "postgresql"
	engineSQLite     = "sqlite"
)

// usesPqArrays reports whether the queries pass array values through
// pq.Array. That is how lib/pq reads and writes PostgreSQL arrays, so it
// holds for PostgreSQL under database/sql; pgx and the other engines'
// drivers scan a slice directly.
func usesPqArrays(engine string, driver opts.SQLDriver) bool {
	if usesDriverTypes(engine) {
		return false
	}
	return !driver.IsPGX()
}

// usesDriverTypes reports whether the engine is one whose mapper reads the
// type expression and draws types from its driver's packages, the ones in
// driverTypes. These engines generate code for database/sql only.
func usesDriverTypes(engine string) bool {
	switch engine {
	case engineClickHouse, engineDuckDB, engineGoogleSQL, engineMSSQL:
		return true
	}
	return false
}

// validateSQLPackage rejects a sql_package the engine's generated code does
// not support. The engines whose mappers read the type expression target
// database/sql alone: their types are what those drivers scan into, and
// their named placeholders bind through sql.Named.
func validateSQLPackage(engine string, options *opts.Options) error {
	if !usesDriverTypes(engine) {
		return nil
	}
	switch options.SqlPackage {
	case "", opts.SQLPackageStandard:
		return nil
	}
	return fmt.Errorf("invalid options: engine %q supports sql_package %q only, not %q", engine, opts.SQLPackageStandard, options.SqlPackage)
}
