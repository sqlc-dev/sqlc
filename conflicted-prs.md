# Open PRs with Merge Conflicts

Generated: 2026-02-25

## PRs with confirmed merge conflicts (21 total)

| PR | Title | Author | Draft |
|----|-------|--------|-------|
| [#3668](https://github.com/sqlc-dev/sqlc/pull/3668) | feat: default schema configuration for postgresql | prog8 | No |
| [#3513](https://github.com/sqlc-dev/sqlc/pull/3513) | fix(normalized): table and column names should not be normalized | a-berahman | Yes |
| [#3370](https://github.com/sqlc-dev/sqlc/pull/3370) | hacky support for Postgres schemas w/ runtime interpolation | bgentry | Yes |
| [#3335](https://github.com/sqlc-dev/sqlc/pull/3335) | Added support to add open-tracing on go mysql | devlibx | No |
| [#3322](https://github.com/sqlc-dev/sqlc/pull/3322) | feat(compiler): Support subqueries in the FROM clause | Jille | No |
| [#3294](https://github.com/sqlc-dev/sqlc/pull/3294) | Draft: clickhouse engine | CNLHC | Yes |
| [#3287](https://github.com/sqlc-dev/sqlc/pull/3287) | enable go_struct_tag on column types | emanuelturis | No |
| [#3279](https://github.com/sqlc-dev/sqlc/pull/3279) | feat(sqlc): Support custom generic nullable types (e.g. `sql.Null`) | gregoryjjb | No |
| [#3130](https://github.com/sqlc-dev/sqlc/pull/3130) | Add support for ignoring DDL statements using `-- sqlc:ignore` | sgielen | No |
| [#3075](https://github.com/sqlc-dev/sqlc/pull/3075) | feat(analyzer): Implement parameter type annotations | andrewmbenton | Yes |
| [#3066](https://github.com/sqlc-dev/sqlc/pull/3066) | feat(ast): First-pass implementation of SQL MERGE | andrewmbenton | Yes |
| [#2966](https://github.com/sqlc-dev/sqlc/pull/2966) | Set the default schema name from configs | debugger84 | No |
| [#2859](https://github.com/sqlc-dev/sqlc/pull/2859) | Dynamic Queries | ovadbar | No |
| [#2858](https://github.com/sqlc-dev/sqlc/pull/2858) | Improve nullability evaluation | lorenz | No |
| [#2687](https://github.com/sqlc-dev/sqlc/pull/2687) | Improve sqlc.embed Go field naming | sapk | No |
| [#2636](https://github.com/sqlc-dev/sqlc/pull/2636) | Add `AdditionalQuery` as function parameter for `:one` and `:many` queries | RadhiFadlillah | No |
| [#2617](https://github.com/sqlc-dev/sqlc/pull/2617) | Add a config flag to add a context key with the current query name | Jille | No |
| [#2570](https://github.com/sqlc-dev/sqlc/pull/2570) | fix(engine/sqlite): added json_tree and json_each definitions | orisano | No |
| [#2510](https://github.com/sqlc-dev/sqlc/pull/2510) | feat(golang): support custom enum slices for pgx/v5 | toqueteos | No |
| [#2376](https://github.com/sqlc-dev/sqlc/pull/2376) | Support :many key=group_id to return a map instead of a slice | Jille | No |
| [#2375](https://github.com/sqlc-dev/sqlc/pull/2375) | Allow for extra parameters to the queryType | Jille | No |
| [#2275](https://github.com/sqlc-dev/sqlc/pull/2275) | Generate nullable value from subselect statements | ryu-ichiroh | No |

## Notes

- Out of ~133 open PRs total, 22 have merge conflicts (~17%)
- 5 of the conflicted PRs are drafts
- Most conflicted PRs are older PRs (from 2023-2024) that have fallen behind `main`
- Conflict status was determined via the GitHub API (`mergeable` field and `mergeable_state=dirty`)
