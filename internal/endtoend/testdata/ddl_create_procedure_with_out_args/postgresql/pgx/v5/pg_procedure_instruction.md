# PostgreSQL Procedure OUT Arguments Support

I have implemented support for capturing `OUT` and `INOUT` arguments from PostgreSQL `CALL` statements in `sqlc`.

## Implementation Details

The `sqlc` analyzer/compiler (`internal/compiler/output_columns.go`) was updated to handle `CALL` statements. Previously, `CALL` statements were treated as having no output columns. The updated logic now:
1. Resolves the procedure definition from the catalog.
2. Identifies arguments with `OUT`, `INOUT`, or `TABLE` modes.
3. Maps these arguments to output columns in the generated code.

## How to get OUT params

To retrieve `OUT` arguments in your Go code:

1. **Annotate your query with `:one`**: Use `:one` (or `:many` if the procedure returns multiple rows) instead of `:exec`. `:exec` tells `sqlc` to disregard any result rows, so it won't capture the `OUT` values.

   ```sql
   -- name: CallInsertData :one
   CALL insert_data($1, $2, null);
   ```

2. **Provide placeholders**: You still need to provide placeholders for the `OUT` arguments in the SQL call (e.g., `null`), as required by PostgreSQL's `CALL` syntax when not using named arguments for everything, or simply to satisfy `sqlc`'s signature matching.

3. **Use the return value**: The generated Go method will now return the `OUT` parameters.
   - If there is a single `OUT` parameter, it will be returned directly.
   - If there are multiple `OUT` parameters, a struct will be returned containing all of them.

### Example

**Schema:**
```sql
CREATE PROCEDURE insert_data(IN a integer, IN b integer, OUT c integer) ...
```

**Query:**
```sql
-- name: CallInsertData :one
CALL insert_data($1, $2, null);
```

**Generated Go:**
```go
func (q *Queries) CallInsertData(ctx context.Context, arg CallInsertDataParams) (pgtype.Int4, error) {
    // ...
    var c pgtype.Int4
    err := row.Scan(&c)
    return c, err
}
```
