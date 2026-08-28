# sqlc fmt verification against the endtoend corpus

`sqlc fmt` was run over every PostgreSQL and MySQL endtoend test case, and the
result verified against what those tests already pin down. This report includes
**only the cases whose changed outcome needs to be fixed**; formatting changes
proven harmless are excluded (counts below).

## Method

For every config directory under `internal/endtoend/testdata` whose config
mentions the `postgresql` or `mysql` engine:

1. Run `sqlc generate` on a pristine copy (control).
2. Run `sqlc fmt` on a second copy, then `sqlc generate` there.
3. Compare the two generated trees with the embedded SQL text stripped
   (`tools/stripgo` empties string const/var initializers via the Go AST, so
   the expected reformatting of the query constant is ignored while any change
   to params structs, signatures, or scans is caught).
4. For MySQL, additionally prove the formatted SQL means the same thing:
   `tools/mysqlfp` parses original and formatted files with marino (the parser
   dolphin wraps) and compares canonical restore forms with identifier case
   preserved and redundant parentheses unwrapped. PostgreSQL already gets this
   proof inside fmt via the pg_query fingerprint.

A case is reported only when fmt crashes or self-reports a bug, generate breaks
after fmt, the generated code changes beyond the SQL text, or the MySQL
semantic check finds drift that live-MySQL testing doesn't clear.

## Coverage

| bucket | configs |
|---|---|
| total configs in testdata | 960 |
| non-pg/mysql engines (out of scope) | 160 |
| processed | 800 |
| fmt made no change | 443 |
| reformatted, verified clean | 336 |
| **reformatted, needs fixing** | **21** |
| managed-db-only (codegen check not runnable locally; fmt's internal proofs + MySQL semantic check still ran) | 24 |
| expected-failure cases (no goldens; fmt ran cleanly) | 15 |
| files fmt skipped as unparseable | 7 (all deliberate `syntax_errors` fixtures + `select_empty_column_list/mysql`) |

No case crashed fmt and none tripped fmt's own "formatting would alter the
file" belt. All 21 failures reduce to six root causes: one in the PostgreSQL
printer, five in the MySQL (dolphin) printer. MySQL fails more because it has
no fingerprint proof — everything the reparse/fixpoint checks can't see gets
through.

---

## P1 — PostgreSQL: `@name` params get a space (`@ name`)

**Cases:** `named_param/pgx/v4`, `named_param/pgx/v5`, `unnest/postgresql/{stdlib,pgx/v4,pgx/v5}`, `unnest_star/postgresql/pgx`

The parser reads `@name` as the unary `@` operator applied to a column, and the
printer then puts a space after the operator. The fingerprint check cannot
catch it — `@name` and `@ name` parse to the same tree — but sqlc's named-param
scanner no longer matches, so **parameters silently vanish from the generated
code**:

```diff
 -- name: AtParams :many
-SELECT name FROM foo WHERE name = @slug AND @filter::bool;
+SELECT name FROM foo WHERE name = @slug AND @ filter::bool;
```

```diff
-func (q *Queries) AtParams(ctx context.Context, arg AtParamsParams) ([]string, error) {
-	rows, err := q.db.Query(ctx, atParams, arg.Slug, arg.Filter)
+func (q *Queries) AtParams(ctx context.Context, slug string) ([]string, error) {
+	rows, err := q.db.Query(ctx, atParams, slug)
```

(`@slug` survived only because `name = @slug` keeps `@` glued to its operand;
any `@param` in operator position breaks.)

## M1 — MySQL: identifier case folding

**Cases:** `identifier_case_sensitivity`, `coalesce_params/mysql`, `params_in_nested_func/mysql`, `overrides_go_types/mysql`, `vet_explain/mysql`

The dolphin printer prints identifiers lowercased (TiDB's `ident.L` instead of
`ident.O`), and drops the backticks with them. Table and database names are
case-sensitive on Linux MySQL, so this breaks queries at runtime:

```diff
--- name: GetAuthor :one
-SELECT * FROM `Authors` WHERE `ID` = ? LIMIT 1;
+SELECT * FROM authors WHERE id = ? LIMIT 1;
```

```diff
-SELECT * FROM foo WHERE retyped IN (sqlc.slice(paramName));
+SELECT * FROM foo WHERE retyped IN (sqlc.slice(paramname));
```

Column-only folding (`vet_explain`) is semantically safe in MySQL but is the
same root cause.

## M2 — MySQL: literal corruption

**Cases:** `case_named_params/mysql` (`NULL` → `''`), `selectstatic/mysql` (`1.0` → `0`, `true` → `1`)

```diff
-SELECT ... CASE WHEN sqlc.arg(email) = '' THEN NULL ELSE sqlc.arg(email) END ...
+SELECT ... CASE WHEN sqlc.arg(email) = '' THEN '' ELSE sqlc.arg(email) END ...
```

```diff
--- name: SelectStatic :one
-SELECT 'a', 'b' AS b, 1 AS num, true AS truefield, 1.0 AS floater
+SELECT 'a', 'b' AS b, 1 AS num, 1 AS truefield, 0 AS floater;
```

`NULL` becoming `''` and `1.0` becoming `0` change query results outright; the
`selectstatic` case also flips the generated field type (`Floater float64` →
`int32`). (`TRUE`/`FALSE` printing as `1`/`0` is technically equivalent in
MySQL and shows up in several otherwise-clean cases, but it is the same broken
literal path.)

## M3 — MySQL: clauses silently dropped

**Cases:** `join_left/mysql`, `order_by_union/mysql`, `mysql_reference_manual`, `mysql_optimizer_hints/mysql`, `update_inner_join`, `update_join/mysql`, `update_two_table/mysql`

- `SELECT DISTINCT` loses `DISTINCT` (`join_left`).
- The trailing `ORDER BY` of a `UNION` is dropped (`order_by_union`).
- `GROUP_CONCAT(DISTINCT x ORDER BY x DESC SEPARATOR ' ')` loses its
  `ORDER BY` (`mysql_reference_manual`).
- Optimizer hints vanish: `SELECT /*+ MAX_EXECUTION_TIME(1000) */ ...` →
  `SELECT ...` (`mysql_optimizer_hints`).
- Multi-table `UPDATE ... JOIN` loses its `ON` clause **and** the table
  qualifiers in `SET` — the formatted statement updates a cross join:

```diff
--- name: UpdateWithJoin :exec
-UPDATE join_table jt JOIN primary_table pt ON jt.primary_table_id = pt.id
-SET jt.is_active = ? WHERE ...
+UPDATE (join_table AS jt) JOIN primary_table AS pt
+SET is_active = ? WHERE ...
```

None of these change the generated Go (codegen doesn't depend on them), so
they are runtime-only breakage — exactly what the missing MySQL fingerprint
proof would have caught.

## M4 — MySQL: DDL type attributes dropped

**Case:** `datatype/mysql` (`sql/numeric.sql`)

```diff
-  h DECIMAL(10, 5), ...
+  h DECIMAL, ...
```

`DECIMAL(10,5)` → `DECIMAL`, and `INT UNSIGNED` → `INT` — precision, scale and
signedness are lost when fmt touches DDL that appears in query files.

## M5 — MySQL: `SHOW WARNINGS` rewritten

**Case:** `show_warnings/mysql`

```diff
--- name: GetWarnings :many
-SHOW WARNINGS;
+SELECT '' AS `level`, 0 AS code, '' AS message;
```

dolphin converts `SHOW` statements into a synthetic `SELECT` for analysis, and
the printer prints that internal form instead of the original statement.

---

## Verified-harmless changes (excluded from the report)

13 configs were flagged by the strict AST comparison but cleared on
inspection/live-MySQL testing, and are counted as clean above:
comma-joins printed as `JOIN` (`join_from`, `star_expansion_join`,
`star_expansion_core`, `cte_recursive`, `cte_count`), redundant parens added
or removed (`IN ((SELECT ...))` — verified equivalent against live MySQL 9 —
in `in_union`, `pattern_in_expr`, `sqlc_arg`, `cte_in_delete`,
`select_union`), `INSERT INTO t () VALUES ()` → `INSERT INTO t VALUES ()`
(`exec_lastid`), and `TRUE`/`FALSE` → `1`/`0` (`case_value_param`,
`cte_count`). One PostgreSQL case (`overrides_result_tag/stdlib`) only moves a
trailing `-- other fields` comment onto its own line, which then surfaces as a
doc comment in the generated Go.

## Reproducing

```bash
go build -o /tmp/sqlc ./cmd/sqlc
cd fmt-verify/tools/stripgo && go build -o /tmp/stripgo .
cd ../mysqlfp && go build -o /tmp/mysqlfp .
# then, per case: copy it, `sqlc fmt`, `sqlc generate`, and compare with
# stripgo / mysqlfp as in scripts/run_corpus.sh and scripts/postpass.sh
```

The scripts in `scripts/` are the exact harness used for this run (paths point
at the session scratchpad; adjust `SP`/`TD` at the top).
