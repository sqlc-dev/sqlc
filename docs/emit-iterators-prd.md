# PRD: Opt-in Iterator Generation for `:many` Queries

**Status:** Proposal  
**Target repo:** [sqlc-dev/sqlc](https://github.com/sqlc-dev/sqlc)  
**Related issues:** [#720](https://github.com/sqlc-dev/sqlc/issues/720), [#4464](https://github.com/sqlc-dev/sqlc/issues/4464), [#4108](https://github.com/sqlc-dev/sqlc/issues/4108)  
**Related PR (closed, reference only):** [#3631](https://github.com/sqlc-dev/sqlc/pull/3631)  
**Author:** Community proposal (reviving stalled discussion)  
**Last updated:** 2026-06-14

---

## 1. Problem Statement

sqlc generates excellent type-safe Go code, but its default API for `:many` queries **always materializes the full result set into a slice**:

```go
func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error)
```

For large result sets (exports, sync jobs, ETL pipelines, backfills), this forces **O(n) heap allocation** even when the caller only needs to process rows one at a time. Alternatives today:

| Workaround | Drawback |
|------------|----------|
| Manual paging with `LIMIT`/`OFFSET` | Extra query complexity; offset cost at scale; not always expressible |
| Fork sqlc or post-process generated code | Maintenance burden; loses upstream improvements |
| Skip sqlc for streaming paths | Loses type safety on the hot path |

Go 1.23 shipped **range-over-function iterators** (`iter.Seq`, `iter.Seq2`). sqlc maintainers [noted in #720](https://github.com/sqlc-dev/sqlc/issues/720) that this unblocks native iterator generation. As of June 2026, **no implementation has merged**; [PR #3631](https://github.com/sqlc-dev/sqlc/pull/3631) was closed without merge after API design remained unresolved.

---

## 2. Goals

1. **Zero breaking changes** — existing `:many` → `[]T` APIs remain the default.
2. **Opt-in streaming** — callers choose slice vs iterator via config or query annotation.
3. **Native Go 1.23 idioms** — generate `iter.Seq2[T, error]` (primary) with optional alternate styles.
4. **Lazy evaluation** — query execution begins on first iteration, not at method call (configurable).
5. **Correct resource lifecycle** — `rows.Close()` on normal completion, early break, panic (via `defer`), and error paths.
6. **Incremental rollout** — Go + `database/sql` first; pgx/stdlib variants and other languages follow.

## 3. Non-Goals (v1)

- Replacing or changing default `:many` behavior.
- Automatic streaming for `:one`, `:exec`, or `:copyfrom`.
- Memory pooling / object reuse (future enhancement; see #3631 discussion).
- Server-side PostgreSQL cursors (`DECLARE`/`FETCH`) — separate feature ([#1517](https://github.com/sqlc-dev/sqlc/issues/1517)).
- Python/Kotlin generators in the initial PR (coordinate separately; see #4464).

---

## 4. Proposed Configuration

### 4.1 `sqlc.yaml` options

```yaml
version: "2"
sql:
  - schema: schema.sql
    queries: queries.sql
    engine: postgresql
gen:
  go:
    package: db
    out: internal/db

    # --- Iterator options (all opt-in, defaults shown) ---

    emit_iterators: false
    # When true, generate a streaming companion method for each :many query.

    iterator_scope: global
    # global         — all :many queries get a streaming method
    # explicit_only  — only queries annotated with :many:stream (or :stream)

    iterator_method_prefix: "Iter"
    # ListAuthors → IterAuthors
    # Set to "Stream" for StreamAuthors if preferred.

    iterator_style: seq2
    # seq2      — iter.Seq2[T, error]  (recommended default)
    # callback  — EachAuthors(ctx, func(Author) error) error
    # rows      — *AuthorsRows with Next()/Scan()/Close()/Err() (legacy #720 style)

    iterator_start: lazy
    # lazy  — DB query runs on first iteration step (recommended)
    # eager — DB query runs at method call; returns (seq, error) or (*Rows, error)
```

### 4.2 Query-level override (optional, for `iterator_scope: explicit_only`)

```sql
-- name: ListAuthors :many:stream
SELECT id, name, bio FROM authors ORDER BY name;
```

Alternatively, a dedicated query kind (as proposed in #4464):

```sql
-- name: StreamAuthors :stream
SELECT id, name, bio FROM authors ORDER BY name;
```

**Recommendation:** support **both** `emit_iterators: global` and `explicit_only` + `:stream` annotation so teams can choose DX vs fine-grained control.

---

## 5. Generated API

### 5.1 Default: `seq2` + `lazy` (recommended)

**SQL (unchanged):**

```sql
-- name: ListAuthors :many
SELECT id, name, bio FROM authors ORDER BY name;
```

**Generated Go:**

```go
import "iter"

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, bio FROM authors ORDER BY name
`

// Existing — unchanged
func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
    // ... current implementation ...
}

// New — opt-in via emit_iterators
func (q *Queries) IterAuthors(ctx context.Context) iter.Seq2[Author, error] {
    return func(yield func(Author, error) bool) {
        rows, err := q.db.QueryContext(ctx, listAuthors)
        if err != nil {
            yield(Author{}, err)
            return
        }
        defer rows.Close()

        for rows.Next() {
            var i Author
            if err := rows.Scan(&i.ID, &i.Name, &i.Bio); err != nil {
                yield(Author{}, err)
                return
            }
            if !yield(i, nil) {
                return // early break; defer closes rows
            }
        }
        if err := rows.Err(); err != nil {
            yield(Author{}, err)
        }
    }
}
```

**Caller usage:**

```go
for author, err := range q.IterAuthors(ctx) {
    if err != nil {
        return fmt.Errorf("list authors: %w", err)
    }
    if err := process(author); err != nil {
        return err
    }
}
return nil
```

**Properties:**

- **Lazy:** no DB round-trip until `range` begins.
- **No wrapper type** for the common case — aligns with Kyle's [later preference](https://github.com/sqlc-dev/sqlc/pull/3631) for `for x, err := range q.Method(ctx)`.
- **Errors in-band** via `Seq2` — familiar Go 1.23 pattern.
- **`break` / `return` safe:** `defer rows.Close()` runs on all exit paths.

### 5.2 Alternate: `seq2` + `eager`

For callers who want connection errors before iteration:

```go
func (q *Queries) IterAuthors(ctx context.Context) (iter.Seq2[Author, error], error) {
    rows, err := q.db.QueryContext(ctx, listAuthors)
    if err != nil {
        return nil, err
    }
    return func(yield func(Author, error) bool) {
        defer rows.Close()
        // ... same loop ...
    }, nil
}
```

### 5.3 Alternate: `callback` style

Sugar for callers who prefer a single error return:

```go
func (q *Queries) EachAuthor(ctx context.Context, fn func(Author) error) error {
    for author, err := range q.IterAuthors(ctx) {
        if err != nil {
            return err
        }
        if err := fn(author); err != nil {
            return err
        }
    }
    return nil
}
```

**Note:** `Each*` can be generated optionally or left as a one-liner at call sites. Generating **both** `Iter*` and `Each*` for every query adds API surface without much benefit — recommend **`seq2` only** in v1, with `callback` as an opt-in `iterator_style`.

### 5.4 Alternate: `rows` style (compatibility with #720 / #3631)

For teams migrating from manual `sql.Rows` patterns:

```go
type IterAuthorsRows struct { /* rows, err */ }
func (q *Queries) IterAuthors(ctx context.Context) *IterAuthorsRows
func (r *IterAuthorsRows) All() iter.Seq2[Author, error]  // or Rows(), Items()
func (r *IterAuthorsRows) Err() error
func (r *IterAuthorsRows) Close() error
```

Useful when lazy start + separate error channel is required; more boilerplate than `seq2`.

---

## 6. Design Decisions & Rationale

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Break `:many` default? | **No** | Maintainer consensus (#720, #4464, Kyle) |
| Keyword vs yaml flag | **Both** | Global flag for DX; `:stream` for explicit control |
| Primary iterator type | `iter.Seq2[T, error]` | Go 1.23 stdlib idiom; Kyle referenced [Thibaut Rousseau's iterator post](https://blog.thibaut-rousseau.com/blog/writing-testing-a-paginated-api-iterator/) |
| Lazy vs eager default | **Lazy** | Avoids dangling queries if iterator is never consumed; matches pierrre/sgielen feedback in #3631 |
| Method naming | `Iter*` default, `Stream*` configurable | `Iter` matches PR #3631; `Stream` matches Kyle's early examples and #4464 |
| Generate 3 methods per query? | **No (v1)** | `List*` + `Iter*` sufficient; `Each*` is optional sugar |
| Min Go version | **1.23+** when `emit_iterators: true` | Required for `iter` package; document in release notes |
| pgx vs database/sql | **database/sql first** | Match existing codegen paths; pgx in follow-up |

---

## 7. Error Handling Semantics

### 7.1 `seq2` lazy mode

| Event | Behavior |
|-------|----------|
| Query fails | First `yield(zero, err)`; iteration ends |
| Scan fails | `yield(zero, err)`; iteration ends |
| `rows.Err()` after loop | Final `yield(zero, err)` |
| Caller `break` / `yield` returns false | Loop stops; `defer rows.Close()` |
| Panic in caller loop body | `defer rows.Close()` still runs |

### 7.2 Close-on-break concern (from #3631)

[gbarr noted](https://github.com/sqlc-dev/sqlc/pull/3631) that without `defer rows.Close()` inside the iterator closure, a recovered panic leaves the connection mid-fetch. **All generated iterators MUST use `defer rows.Close()`** inside the `Seq2` closure.

### 7.3 Context cancellation

Callers may cancel via `ctx`. Behavior depends on driver:

- `database/sql`: `rows.Next()` may block until cancel (driver-dependent).
- **v1:** pass `ctx` to `QueryContext`; document that full cancel propagation requires driver support.
- **Future:** optional `select { case <-ctx.Done(): ... }` in loop (as MatthiasKunnen uses with pgx).

---

## 8. Parameterized Queries

Iterators work identically for parameterized `:many` queries:

```sql
-- name: ListAuthorsByIDs :many
SELECT id, name, bio FROM authors WHERE id = ANY($1::int[]);
```

```go
func (q *Queries) IterAuthorsByIDs(ctx context.Context, ids []int32) iter.Seq2[Author, error]
func (q *Queries) ListAuthorsByIDs(ctx context.Context, ids []int32) ([]Author, error)
```

Same SQL constant, same prepared statement wiring — only the result consumption differs.

---

## 9. Implementation Plan

### Phase 1 — Design sign-off (this proposal)

- [ ] Post proposal to #720; cross-link #4464
- [ ] Maintainer confirmation on: naming, lazy default, global vs explicit scope
- [ ] Agree v1 scope: Go + `database/sql` + PostgreSQL example

### Phase 2 — PoC PR

- [ ] Add config parsing in `internal/codegen/golang/opts`
- [ ] Extend `:many` code generation in `internal/codegen/golang/query.go` (or templates)
- [ ] Generate `Iter*` method alongside existing `List*`
- [ ] End-to-end test in `examples/` (pattern from #3631)
- [ ] Document in sqlc.dev docs

### Phase 3 — Expand coverage

- [ ] MySQL, SQLite engines
- [ ] pgx/v5 driver variant
- [ ] `iterator_style: rows` and `callback` if requested
- [ ] Python generator (coordinate with borissmidt / sqlc-gen-python)

---

## 10. Testing Requirements

1. **Unit:** generated code compiles with `go test` under Go 1.23+.
2. **Integration:** iterator returns all rows; early `break` closes rows (verify via connection pool or mock).
3. **Error paths:** query error, scan error, `rows.Err()` — each yields exactly one error and stops.
4. **Parity:** for same fixture, `List*` and collecting `Iter*` produce identical slices.
5. **Opt-out:** `emit_iterators: false` produces byte-identical output to today (regression).

---

## 11. Open Questions for Maintainers

1. **Preferred method prefix:** `Iter` vs `Stream`?
2. **Lazy default:** agree lazy is correct for v1?
3. **Global flag vs `:stream` only:** ship both?
4. **Eager mode:** worth exposing in v1 or defer?
5. **pgx:** same PR or immediate follow-up?
6. **Min Go version bump:** gate behind `emit_iterators` or raise global minimum?

---

## 12. One-Line Pitch

> sqlc generates type-safe Go that materializes every `:many` query into a slice; with Go 1.23, an opt-in `emit_iterators` flag can generate lazy `iter.Seq2[T, error]` companions — same type safety, O(1) memory per row, zero breaking changes.

---

## Appendix A: Comparison with Closed PR #3631

| Aspect | PR #3631 | This proposal |
|--------|----------|---------------|
| Trigger | `:iter` query annotation | `emit_iterators` yaml + optional `:stream` |
| API | Wrapper type + `Iterate()` + `Err()` | `iter.Seq2` direct range (default) |
| Lazy start | Unclear / eager in PoC | Explicit `iterator_start: lazy` default |
| Config surface | None | Full yaml options |
| Status | Closed, not merged | — |

This proposal incorporates #3631's implementation lessons and resolves the API debates raised in its review thread.

## Appendix B: References

- [#720 — Ability to return an iterator on a "many" query](https://github.com/sqlc-dev/sqlc/issues/720)
- [#4464 — add :stream keyword](https://github.com/sqlc-dev/sqlc/issues/4464)
- [#4108 — low-level prepare/bind helpers for streaming](https://github.com/sqlc-dev/sqlc/issues/4108)
- [PR #3631 — Ability to return an iterator on rows (closed)](https://github.com/sqlc-dev/sqlc/pull/3631)
- [Go 1.23 — range-over-func](https://go.dev/doc/go1.23#range-over-function)
- [Eli Bendersky — Ranging over functions in Go 1.23](https://eli.thegreenplace.net/2024/ranging-over-functions-in-go-123/)
