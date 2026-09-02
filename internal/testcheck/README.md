# testcheck

`testcheck` verifies the analyze cases under `internal/endtoend/testdata`
against a real database. It generates nothing. Each engine package reads a
case's `schema.sql`, `fixture.sql` and `query.sql`, asks the database what it
makes of the queries, prints the answer in the JSON shape `sqlc analyze`
prints, and compares it with the `output.json` the case committed, byte for
byte. A difference means sqlc's analysis disagrees with the database.

It is a nested Go module with no dependencies beyond the standard library,
so it never shares code with the analysis it checks. Run it from this
directory:

```bash
go run ./cmd/testcheck install clickhouse   # download the pinned clickhouse binary once
go run ./cmd/testcheck check                # check every engine whose database is available
go run ./cmd/testcheck check clickhouse     # check one engine
go test ./...                 # the same checks as tests; engines without a database skip
```

## Cases

A case is an `analyze_<name>/<engine>` directory whose `exec.json` runs the
analyze command. `fixture.sql` is optional and is loaded after the schema, so
the queries run against real rows. A case that asks for `--ast` is skipped,
since only sqlc can print that.

## Engines

Each engine is its own package.

- **`clickhouse`** needs no server. Each case runs in an ephemeral
  `clickhouse local` process, downloaded once per pinned version by
  `install` into the user cache directory, or supplied through the
  `CLICKHOUSE` environment variable. The pinned version and the SHA-512 of
  each platform's download live in `clickhouse/install.go`. Column types
  come from the executed query's result header, provenance from
  `EXPLAIN QUERY TREE`, and parameters from sentinel constants substituted
  for `?`, `sqlc.arg()` and `sqlc.narg()`, since ClickHouse itself never
  sees a placeholder.
