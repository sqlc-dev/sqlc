// Package catalogdef holds the SQL that defines sqlc's internal catalog:
// schema.sql (the sql_* table definitions) and query.sql (the queries
// against them). Both are the input to sqlc's own code generation — sqlc
// analyzing sqlc's catalog — and schema.sql is embedded here so the core
// package builds its in-memory catalog from the same source of truth that
// codegen reads.
package catalogdef

import _ "embed"

// Schema is the DDL that creates the sql_* catalog tables.
//
//go:embed schema.sql
var Schema string
