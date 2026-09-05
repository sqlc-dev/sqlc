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
go run ./cmd/goldeneye install sqlite       # build the pinned sqlite3 shells once; needs a C compiler
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
- **`sqlite`** needs no server either: `functions.jsonl` comes from
  `pragma_function_list` of a `sqlite3` shell run against an in-memory
  database. Which functions a SQLite has is decided when it is compiled, so
  `install` downloads the pinned release's amalgamation, checked against the
  SHA3-256 the download page lists, and compiles the shell from it with the
  compiler `CC` names, or `cc` — once with the options sqlite.org's own
  configure turns on by default, which gives `functions.jsonl`, and once
  more per option in `sqlite/install.go`'s extension list, each of which
  gets a directory under `extensions/` holding the functions its build adds
  over the default one, the way each PostgreSQL contrib extension holds what
  `CREATE EXTENSION` adds; a schema that says `CREATE VIRTUAL TABLE ... USING
  fts5` loads the option's directory, through the `modules` map in the
  hand-written `dialect.json`. SQLite describes its functions as far as their
  names, their kinds and the number of arguments each overload takes, and no
  further — it types values, not functions — so what each returns and what
  its arguments hold is read from the amalgamation: every function is
  registered with the C functions that implement it, and those set their
  result through `sqlite3_result_*` and read their arguments through
  `sqlite3_value_*`. A function a shell reports that the source does not
  register fails the run rather than being guessed at. Whether an aggregate
  returns NULL over no rows is found by running it over none; the scalar
  functions that return NULL for arguments that are not are a short list in
  `sqlite/signatures.go`, since a SQLite function returns NULL as often by
  setting no result as by saying so. The pinned release is the one the main module's
  driver embeds. SQLite has no catalog of types or operators, so
  `types.jsonl` and `operators.jsonl` are hand-written.

## Layout

- `dialect/` — the record types the files are made of, mirrored from
  `internal/core/seed`, and the helpers that write a generated set of files
  into an engine directory or diff it against what is committed.
- `postgresql/`, `duckdb/`, `clickhouse/`, `sqlite/` — one package per
  engine, each exposing `Locate`, `Version` and `Generate`, and a test that
  runs the check.
- `cmd/goldeneye/` — the command.

The analysis checks — verifying the `analyze_*` cases under
`internal/endtoend/testdata` against what each database itself reports — are
meant to live here too, alongside the dialect checks.
