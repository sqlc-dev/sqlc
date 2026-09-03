# goldeneye

`goldeneye` generates the dialect seeds under `internal/engine/<engine>/dialect`
from a live database, and verifies the committed ones against it. A dialect
is the JSONL that gives an engine its type system and standard library —
`types.jsonl`, `functions.jsonl`, `operators.jsonl`, `relations.jsonl` and
the `extensions/` bundles — read by `internal/core/seed`. Each engine package
here asks the database what it knows, writes the answer in that shape, and
the tests compare it with what is committed, byte for byte. A difference
means the committed dialect has drifted from the database.

It is a nested Go module, so its only dependency beyond the standard library
is the PostgreSQL driver, and it never shares code with the analysis that
reads the files: the files are the contract. Run it from this directory:

```bash
go run ./cmd/goldeneye install clickhouse   # download the pinned clickhouse binary once
go run ./cmd/goldeneye check                # check every engine whose database is available
go run ./cmd/goldeneye check postgresql     # check one engine
go run ./cmd/goldeneye generate [engine]    # rewrite the generated files from the database
go test ./...                               # the same checks as tests; engines without a database skip
```

`generate` and `check` say which engines they skipped for lack of a
database; naming an engine makes its database required.

## What is generated, and what is not

A generator owns only the files it produces; `dialect.json` is always written
by hand, and so are the lists an engine cannot describe. Both commands leave
the hand-written files alone, and the checks do not look at them.

- **`postgresql`** reads a live server named by `POSTGRESQL_SERVER_URI`:
  `functions.jsonl` is `pg_catalog`'s functions, `relations.jsonl` the tables
  and views of `pg_catalog` and `information_schema`, and each contrib
  extension gets a directory under `extensions/` holding the types and
  functions `CREATE EXTENSION` adds, so the server needs contrib installed.
  A function that one of those extensions puts in `pg_catalog`, as
  `adminpack` does, belongs to the extension's directory rather than the
  catalog's list. The server has to be the major release pinned in
  `postgresql.Major`, since every release adds to the catalogs; the top-level
  `types.jsonl` and `operators.jsonl` are hand-written.
- **`duckdb`** reads the DuckDB CLI named by `DUCKDB`, or `duckdb` on `PATH`:
  `types.jsonl`, `functions.jsonl` and `operators.jsonl` come from
  `duckdb_types()` and `duckdb_functions()`. The CLI has to be the DuckDB 2.0
  build darkwing is pinned against, which has no release to download yet.
- **`clickhouse`** needs no server: `types.jsonl` comes from
  `system.data_type_families` of an ephemeral `clickhouse local` process,
  every family that is not an alias becoming a type carrying the spellings
  that alias it, with a category decided by its name. The binary is
  downloaded once per pinned release by `install` into the user cache
  directory, or supplied through the `CLICKHOUSE` environment variable; the
  pinned release and the SHA-512 of each platform's download live in
  `clickhouse/install.go`, and a download that does not match is discarded.
  ClickHouse describes its functions no further than their names, so
  `functions.jsonl` is hand-written.

## Layout

- `dialect/` — the record types the files are made of, mirrored from
  `internal/core/seed`, and the helpers that write a generated set of files
  into an engine directory or diff it against what is committed.
- `postgresql/`, `duckdb/`, `clickhouse/` — one package per engine, each
  exposing `Locate`, `Version` and `Generate`, and a test that runs the check.
- `cmd/goldeneye/` — the command.

The analysis checks — verifying the `analyze_*` cases under
`internal/endtoend/testdata` against what each database itself reports — are
meant to live here too, alongside the dialect checks.
