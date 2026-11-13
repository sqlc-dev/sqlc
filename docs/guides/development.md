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

`core.TypeCombiner` receives an operation context and the entire ordered list
of candidate types, retaining nil entries for untyped values. It returns a
type and an explicit `Handled` flag; declining a context preserves the old
inference fallback, while an error is propagated. Engines without a combiner
keep their existing behavior. `Prepare` accepts the combiner through an
optional analyzer option, inherited by nested queries.

When a converter lowers a function to a shared expression, it can retain the
source function name rather than encoding a result-type instruction in the
AST. `TypeCombinationRequest.FunctionName` passes that qualified identity to
the engine. For example, SQL Server's policy, not shared analysis, selects
ISNULL's first typed argument; ordinary COALESCE still combines its candidates.
The shared first-typed-value helper is only the compatibility fallback for
engines without a policy or for contexts a policy declines.

The compiler's existing engine assembly selects an optional combiner factory.
It calls that factory after `core.NewCached` returns, on both fresh and
restored catalogs, and passes the instance to analysis. No engine names,
imports or strategy registry are needed in shared core analysis or seeding.
Combiners must allow concurrent analysis and must not mutate candidate types.

`set_operations` in `dialect.json` is opaque JSON to shared seeding. Its schema
and validation belong to the engine implementation; SQL Server's precedence
and decimal settings are not a universal coercion language. Only the data is
cached; runtime behavior is constructed afresh. Malformed configuration fails
in the engine factory, and configuration without a provider is a setup error.

To support another dialect, implement its combiner and wire its factory at
the compiler's engine assembly, without editing shared policy dispatch or
validation. Add end-to-end cases rather than copying SQL Server's settings.
Some engines additionally need literal/domain metadata or recursive structured
types; extend the inputs for actual implementations rather than forcing those
semantics through pairwise precedence. Set operations and `CASE`-like
expressions need not share an algorithm. The current cross-engine gaps and
database references are in
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
