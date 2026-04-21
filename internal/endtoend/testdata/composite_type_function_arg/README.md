# Composite types as function arguments

sqlc generates a Go struct for every PostgreSQL composite type that is
reachable from a query (column, function argument, function return column,
etc.).  When a user-defined composite type is used as an IN parameter to a
function — scalar or array — the generated `Params` struct references that
struct directly.

## Runtime: registering composite types with pgx/v5

sqlc does *not* emit runtime registration code.  pgx needs the OID of each
user-defined type before it can encode/decode values.  Register the types at
pool init, e.g. via `pgxpool.Config.AfterConnect`:

```go
cfg, _ := pgxpool.ParseConfig(dsn)
cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
    for _, name := range []string{
        "masterdata.rolling_stock_number_input",
        "masterdata._rolling_stock_number_input", // array variant
    } {
        t, err := conn.LoadType(ctx, name)
        if err != nil {
            return err
        }
        conn.TypeMap().RegisterType(t)
    }
    return nil
}
```

Rules:

* Register the element type **before** the array type (pgx requires this).
* Array types in PostgreSQL are named `_<basetype>` (underscore prefix),
  possibly schema-qualified as `schema._<basetype>`.
* Field order on the generated Go struct mirrors the `CREATE TYPE` column
  order; pgx's `CompositeCodec` uses that order when encoding/decoding.
* Field types must be ones pgx already knows how to encode
  (`pgtype.Text`, `pgtype.Int8`, etc., or standard Go types).  This matches
  what sqlc already emits for the same column types when they appear in
  regular tables.

`Conn.LoadTypes` (pgx ≥ 5.5) is a more efficient alternative for
registering many types in one round-trip.
