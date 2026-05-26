# Proposal: `sqlc.switch` — bounded dynamic ORDER BY / WHERE via compile-time branch expansion

Status: **draft** (evox-it fork)
Tracking upstream: discussions/364, issues #2061, #3414, #2060; PRs #4005, #4260, #2859 (all blocked on a canonical design)

## Problem

sqlc has no way to vary the *structure* of a query (sort order, filter shape) at
runtime. The community has asked for "dynamic queries" since 2020. The two
existing workarounds both fail:

1. **`CASE WHEN` in `ORDER BY`/`WHERE`** — defeats the query planner. Postgres
   will not use an index when the sort key is hidden inside a `CASE` expression.
2. **Runtime string interpolation** (`sqlc.raw`, `@filter::text`) — every such
   proposal has been closed because it reintroduces SQL injection and breaks
   sqlc's "type-safe, it's just SQL" guarantee.

## Design goals (derived from maintainer objections)

- **No SQL injection, ever.** User input must never reach the query string.
- **Schema-validated.** Every dynamic fragment must parse as real SQL and
  reference real columns at compile time. Bad column = compile error.
- **Planner/index friendly.** The emitted SQL must be a clean static
  `ORDER BY col DESC`, never a `CASE` wrapper.
- **Finite + enumerable.** The set of runtime choices is fixed at compile time.
- **Modeled on existing precedent** (`sqlc.slice`, the `sqlc.*` macro family).

## Syntax

```sql
-- name: ListAuthors :many
SELECT * FROM authors
WHERE deleted_at IS NULL
ORDER BY sqlc.switch(@sort,
    sqlc.when('name_asc',  'authors.name ASC'),
    sqlc.when('recent',    'authors.created_at DESC, authors.id DESC'),
    sqlc.else(             'authors.id ASC')
);
```

- `sqlc.switch(@selector, branches…)` sits where a value/expression is grammatically
  legal (`ORDER BY` position, `WHERE` position). `@selector` is the runtime chooser.
- `sqlc.when('key', 'sql-fragment')` — `'key'` is the enum value; `'sql-fragment'`
  is a **string literal** (must be — `ASC`/`DESC` are not valid inside a function
  arg list in any engine grammar). The fragment is an author-authored compile-time
  constant.
- `sqlc.else('sql-fragment')` — optional default branch.

## Why string-literal fragments are still safe

The fragment is a constant in the `.sql` file written by the developer, exactly
like the rest of the query. It is **not** runtime input. sqlc re-parses each
fragment in its grammatical context (e.g. `SELECT 1 FROM authors ORDER BY <frag>`)
and validates every column reference against the catalog. A typo or unknown
column fails `sqlc generate`. The only thing that varies at runtime is *which
already-validated branch* is chosen — a closed enum. Injection is structurally
impossible.

## Codegen strategy: compile-time expansion (one function per branch)

This follows Kyle Conroy's own `sqlc.switch()` suggestion in discussions/364
("generate multiple optimized queries at compile time rather than runtime CASE").

A query containing `sqlc.switch` is **expanded by the compiler into N concrete
queries**, one per branch. Each clone has the whole `sqlc.switch(...)` call
replaced in the SQL by that branch's fragment. The resulting query strings are
fully static constants — no runtime `strings.Replace`, no markers.

### Recognition is AST-based, like the other sqlc.* macros

`sqlc.switch`/`when`/`else` are recognized exactly the way `sqlc.arg` and
`sqlc.slice` are: by searching the parsed AST for a `FuncCall` whose schema is
`sqlc` (`astutils.Search`). There is no bespoke SQL lexer. A consequence is that
the feature works in precisely the clauses where the engine parser produces such
a node — i.e. **wherever `sqlc.arg` works**:

| Position | PostgreSQL | MySQL | SQLite |
|---|---|---|---|
| WHERE | ✅ | ✅ | ✅ |
| ORDER BY | ✅ | ✅ | ❌ parser drops the clause |

SQLite's parser discards *any* function call in `ORDER BY` (true of plain
`ORDER BY upper(name)` too — see upstream PR #4429), so `sqlc.switch` there is a
compile error rather than silently emitting the unexpanded call. This is the
same limitation `sqlc.arg` has.

Once recognized, the compiler replaces the `sqlc.switch(...)` text span with each
branch's fragment, renames the `-- name:` comment to `<QueryName><BranchKey>`,
and **re-parses each branch as an ordinary query**. Every branch therefore goes
through the normal parser + analyzer, so a bad column reference in a fragment is
a compile error, and the generated query strings are fully static constants — no
runtime markers, no `strings.Replace`.

### Generated Go — v1 (implemented)

One static function per branch:

```go
const listAuthorsNameAsc = `SELECT ... ORDER BY authors.name ASC`
func (q *Queries) ListAuthorsNameAsc(ctx context.Context) ([]Author, error) { ... }

const listAuthorsRecent = `SELECT ... ORDER BY authors.created_at DESC, authors.id DESC`
func (q *Queries) ListAuthorsRecent(ctx context.Context) ([]Author, error) { ... }

const listAuthorsElse = `SELECT ... ORDER BY authors.id ASC`
func (q *Queries) ListAuthorsElse(ctx context.Context) ([]Author, error) { ... }
```

This is the whole upstreamable primitive: pure compile-time expansion in the
compiler, **zero codegen changes** (the branches are ordinary queries, so every
language's codegen gets them for free). It is the conservative core to propose
first.

### Generated Go — v2 (proposed extension)

A generated enum for the selector plus one exported dispatcher that switches on
it. This is a codegen convenience layered on top of v1; it adds codegen surface
(enum synthesis + dispatcher emission) and is best proposed as a follow-up once
the primitive is accepted:

```go
type ListAuthorsSort string
const (
    ListAuthorsSortNameAsc ListAuthorsSort = "name_asc"
    ListAuthorsSortRecent  ListAuthorsSort = "recent"
)
func (q *Queries) ListAuthors(ctx context.Context, sort ListAuthorsSort) ([]Author, error) {
    switch sort {
    case ListAuthorsSortNameAsc: return q.ListAuthorsNameAsc(ctx)
    case ListAuthorsSortRecent:  return q.ListAuthorsRecent(ctx)
    default:                     return q.ListAuthorsElse(ctx)
    }
}
```

## Naming rules

| Element | Source | Example |
|---|---|---|
| Branch fn | `<QueryName>` + camelize(key) | `ListAuthorsNameAsc` |
| `sqlc.else` fn | `<QueryName>` + `Else` | `ListAuthorsElse` |
| Enum type (v2) | `<QueryName>` + selector name | `ListAuthorsSort` |
| Enum const (v2) | enum type + camelize(key) | `ListAuthorsSortNameAsc` |
| Dispatcher (v2) | `<QueryName>` | `ListAuthors` |

## v1 scope (implemented + tested)

- Recognized in WHERE (all engines) and ORDER BY (PostgreSQL, MySQL); SQLite
  ORDER BY is a clear compile error (parser limitation, parity with `sqlc.arg`).
- **Not allowed in the SELECT projection** — branches there could change the
  result columns; rejected at compile time.
- One static function per branch. Enum + dispatcher (v2) are the proposed
  follow-up.
- Engine-agnostic: the compiler emits ordinary queries, so all codegens benefit;
  no codegen changes in v1.
- Golden end-to-end tests for PostgreSQL (stdlib + pgx), MySQL, and SQLite.

## Open questions for upstream

1. Dispatcher on by default, or opt-in via a query annotation?
2. Should `sqlc.else` be mandatory (compile error if a non-exhaustive switch is
   possible) or optional (zero-value selector → else)?
3. Fragment validation depth: parse-only vs full type-check of the sort expr.
