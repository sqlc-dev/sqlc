# ClickHouse testgen

`testgen` records what ClickHouse itself reports about a set of sqlc queries,
so the answer can be committed as a golden file and compared with what sqlc's
own analysis produces for the same schema and queries.

It is a nested Go module with no dependencies beyond the standard library.
Run it from this directory:

```bash
# Download the pinned clickhouse binary into the user cache directory.
go run . install

# Analyze every query in query.sql against schema.sql with fixture.sql loaded.
go run . analyze --schema schema.sql --fixture fixture.sql query.sql
```

The binary is looked up in the `CLICKHOUSE` environment variable first, then
in the cache populated by `install`. The version is pinned in `install.go`;
bumping it can change the query tree format and type inference, so regenerate
and review the goldens afterwards.

## How it works

Each query runs in its own `clickhouse local` process, which needs no server,
no network and no configuration: the schema and fixture are loaded fresh, the
query is explained and then executed, and the process exits.

- **Types and nullability** come from the executed query's result header,
  exactly as a driver would see them.
- **Provenance** comes from `EXPLAIN QUERY TREE`, whose resolved columns point
  at the table expression they read from. References are followed through
  subqueries, CTEs and unions to the base table.
- **Parameters** are invisible to ClickHouse, which substitutes `?` on the
  client. Each `?`, `sqlc.arg()` and `sqlc.narg()` is replaced by a constant
  expression that carries its ordinal, `(NULL + k)`, or `toUInt64(4294967295 + k)`
  after `LIMIT` and `OFFSET`. The query tree prints the expression each folded
  constant came from, so the placeholder is found again and described by the
  operand it is compared with. Parameters of `INSERT ... VALUES` map onto the
  target columns reported by `DESCRIBE TABLE`.

The output has the same shape as `sqlc analyze`, with ClickHouse types lowered
the way sqlc's ClickHouse engine lowers them: the lowercased base name with
parameters dropped, `Nullable` and `Array` folded into `not_null` and
`is_array`, and `LowCardinality` discarded.

## Tests

`testdata/<case>/` holds a `schema.sql`, `query.sql`, an optional
`fixture.sql` and the expected `stdout.txt`. The test skips unless a binary is
available.

```bash
go run . install
go test .
go test . -update   # rewrite every stdout.txt
```
