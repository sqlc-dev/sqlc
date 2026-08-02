# Merging rows

sqlc supports PostgreSQL's [`MERGE`](https://www.postgresql.org/docs/current/sql-merge.html)
statement (PostgreSQL 15+). Parameters are inferred in every clause: the source
subquery, the join condition, `WHEN` conditions, `UPDATE SET` assignments, and
`INSERT ... VALUES` lists.

```sql
CREATE TABLE inventory (
  sku        text PRIMARY KEY,
  quantity   integer NOT NULL,
  updated_by text
);

CREATE TABLE inventory_updates (
  sku   text PRIMARY KEY,
  delta integer NOT NULL
);
```

```sql
-- name: SyncInventory :exec
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED AND t.quantity + u.delta <= $1 THEN
    DELETE
WHEN MATCHED THEN
    UPDATE SET quantity = t.quantity + u.delta, updated_by = $2
WHEN NOT MATCHED THEN
    INSERT (sku, quantity, updated_by) VALUES (u.sku, u.delta, $2);
```

```go
type SyncInventoryParams struct {
	Quantity  int32
	UpdatedBy sql.NullString
}

func (q *Queries) SyncInventory(ctx context.Context, arg SyncInventoryParams) error {
	_, err := q.db.ExecContext(ctx, syncInventory, arg.Quantity, arg.UpdatedBy)
	return err
}
```

## Returning merged rows

On PostgreSQL 17 and later, `MERGE` supports a `RETURNING` clause, which can be
used with `:one` or `:many`. The clause can reference columns from both the
target and the source relations.

```sql
-- name: MergeInventory :many
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED THEN
    UPDATE SET quantity = u.delta
WHEN NOT MATCHED THEN
    INSERT (sku, quantity) VALUES (u.sku, u.delta)
RETURNING t.sku, t.quantity;
```

```go
type MergeInventoryRow struct {
	Sku      string
	Quantity int32
}

func (q *Queries) MergeInventory(ctx context.Context) ([]MergeInventoryRow, error) {
	// ...
}
```

Note that sqlc does not validate the server version: generating code for a
`MERGE ... RETURNING` query succeeds, but running it against a server older
than PostgreSQL 17 returns a syntax error at execution time.
