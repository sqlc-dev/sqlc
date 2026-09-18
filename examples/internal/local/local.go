// Package local creates a database per test on the servers the examples
// for ClickHouse, SQL Server and Spanner run against, named by environment
// variable, and skips the test when none is named. The helpers for
// PostgreSQL, MySQL and SQLite live in the main module's
// internal/sqltest, which the rest of sqlc's tests share; these live with
// the examples so the drivers they need are dependencies of the examples
// module alone.
package local

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func id() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
