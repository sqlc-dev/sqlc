package golang

import "github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"

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

// usesPqArrays reports whether the queries pass array values through
// pq.Array. That is how lib/pq reads and writes PostgreSQL arrays, so it
// holds for PostgreSQL under database/sql; pgx and the other engines'
// drivers scan a slice directly.
func usesPqArrays(engine string, driver opts.SQLDriver) bool {
	switch engine {
	case "clickhouse", "duckdb", "googlesql", "mssql":
		return false
	}
	return !driver.IsPGX()
}
