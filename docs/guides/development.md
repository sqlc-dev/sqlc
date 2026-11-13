# Developing sqlc

## Building

For local development, install `sqlc` under an alias. We suggest `sqlc-dev`.

```
go build -o ~/go/bin/sqlc-dev ./cmd/sqlc
```

Install `sqlc-gen-json` to avoid test failure.

```
go build -o ~/go/bin/sqlc-gen-json ./cmd/sqlc-gen-json
```

## Dialect type-combination policies

The analysis core owns traversal, column names, catalog resolution and
nullability propagation. Database-specific common-type algorithms belong to
the engine: SQL Server's precedence, length and decimal rules live in
`internal/engine/mssql/typepolicy`, independent of its parser.

`set_operations.strategy` in a dialect's `dialect.json` explicitly selects an
implemented policy (currently only `"mssql"`). The remaining fields configure
that SQL Server policy; they are not a universal coercion language. Seeds
reject missing or unsupported strategies rather than accidentally applying
SQL Server's formula to another engine. The selector is serialized in the
catalog, so cache restores use the same policy as freshly seeded catalogs.
Omitting `set_operations` preserves the old inference fallback.

To support another dialect, add its engine-owned policy and explicit
dispatch/seed validation, then opt in its dialect settings and add end-to-end
cases. Do not copy SQL Server's settings into other seeds. Some engines need
the entire candidate list, cast context, domain information or recursive
structured types; extend the policy inputs/orchestration for those semantics
rather than forcing them through pairwise precedence. Set operations and
`CASE`-like expressions need not share one algorithm. The current cross-engine
gaps and database references are in
[`analyze`](../howto/analyze.md#common-type-inference).

SQL Server's converter also performs syntax/semantic lowering before seeded
function lookup. `convertMethodCall` handles receivers, static `type::method`
syntax and XML `value`'s type argument; `convertAlternatives` produces the
shared AST nodes that combine candidates and propagate nullability. Ordinary
function calls retain `OVER` rather than being rewritten as alternatives.
JSONL signatures already describe the lowered calls' argument and return
types, but cannot currently request these AST transformations. The method
recognition map must stay synchronized with the leading-dot signatures.
The datepart recognition map distinguishes bare keyword arguments from column
references; `ident_args` in the dialect settings applies to *type* arguments,
not function arguments. Moving these maps into seed metadata would require
new converter metadata and a way to consume it before analysis, not merely
additional function signatures.

## Running Tests

```
go test ./...
```

To run the tests in the examples folder, which is a module of its own, use the `examples` tag from inside it.

```
cd examples && go test --tags=examples ./...
```

These tests require locally-running database instances. Run these databases
using [Docker Compose](https://docs.docker.com/compose/).

```
docker compose up -d
```

The tests use the following environment variables to connect to the
database

### For PostgreSQL

```
Variable     Default Value
-------------------------
PG_HOST      127.0.0.1
PG_PORT      5432
PG_USER      postgres
PG_PASSWORD  mysecretpassword
PG_DATABASE  dinotest
```

### For MySQL

```
Variable     Default Value
-------------------------
MYSQL_HOST      127.0.0.1
MYSQL_PORT      3306
MYSQL_USER      root
MYSQL_ROOT_PASSWORD  mysecretpassword
MYSQL_DATABASE  dinotest
```
