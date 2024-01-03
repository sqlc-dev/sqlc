# Changelog
All notable changes to this project will be documented in this file.

(v1-25-0)=
## [1.25.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.25.0)
Released 2024-01-03

### Release notes

#### Add tags to push and verify

You can add tags when [pushing](../howto/push.md) schema and queries to [sqlc Cloud](https://dashboard.sqlc.dev). Tags operate like git tags, meaning you can overwrite previously-pushed tag values. We suggest tagging pushes to associate them with something relevant from your environment, e.g. a git tag or branch name.

```
$ sqlc push --tag v1.0.0
```

Once you've created a tag, you can refer to it when [verifying](../howto/verify.md) changes, allowing you
to compare the existing schema against a known set of previous queries.

```
$ sqlc verify --against v1.0.0
```

#### C-ya, `cgo`

Over the last month, we've switched out a few different modules to remove our reliance on [cgo](https://go.dev/blog/cgo). Previously, we needed cgo for three separate functions:

- Parsing PostgreSQL queries with [pganalyze/pg_query_go](https://github.com/pganalyze/pg_query_go)
- Running SQLite databases with [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- Executing WASM / WASI code with [bytecodealliance/wasmtime-go](https://github.com/bytecodealliance/wasmtime-go)

With the help of the community, we found cgo-free alternatives for each module:

- Parsing PostgreSQL queries, now using [wasilibs/go-pgquery](https://github.com/wasilibs/go-pgquery)
- Running SQLite databases, now using [modernc.org/sqlite](https://gitlab.com/cznic/sqlite)
- Executing WASM / WASI code, now using [tetratelabs/wazero](https://github.com/tetratelabs/wazero)

For the first time, Windows users can enjoy full PostgreSQL support without using [WSL](https://learn.microsoft.com/en-us/windows/wsl/about). It's a Christmas miracle!

If you run into any issues with the updated dependencies, please [open an issue](https://github.com/sqlc-dev/sqlc/issues).

### Changes 

#### Bug Fixes

- (codegen) Wrong yaml annotation in go codegen options for output_querier_file_name (#3006)
- (codegen) Use derived ArrayDims instead of deprecated attndims (#3032)
- (codegen) Take the maximum array dimensions (#3034)
- (compiler) Skip analysis of queries without a `name` annotation (#3072)
- (codegen/golang) Don't import `"strings"` for `sqlc.slice()` with pgx (#3073)

### Documentation

- Add name to query set configuration (#3011)
- Add a sidebar link for `push`, add Go plugin link (#3023)
- Update banner for sqlc-gen-typescript (#3036)
- Add strict_order_by in doc (#3044)
- Re-order the migration tools list (#3064)

### Features

- (analyzer) Return zero values when encountering unexpected ast nodes (#3069)
- (codegen/go) add omit_sqlc_version to Go code generation (#3019)
- (codgen/go) Add `emit_sql_as_comment` option to Go code plugin (#2735)
- (plugins) Use wazero instead of wasmtime (#3042)
- (push) Add tag support (#3074)
- (sqlite) Support emit_pointers_for_null_types (#3026)

### Testing

- (endtoend) Enable for more build targets (#3041)
- (endtoend) Run MySQL and PostgreSQL locally on the runner (#3095)
- (typescript) Test against sqlc-gen-typescript (#3046)
- Add tests for omit_sqlc_version (#3020)
- Split schema and query for test (#3094)

### Build

- (deps) Bump idna from 3.4 to 3.6 in /docs (#3010)
- (deps) Bump sphinx-rtd-theme from 1.3.0 to 2.0.0 in /docs (#3016)
- (deps) Bump golang from 1.21.4 to 1.21.5 (#3043)
- (deps) Bump actions/setup-go from 4 to 5 (#3047)
- (deps) Bump github.com/jackc/pgx/v5 from 5.5.0 to 5.5.1 (#3050)
- (deps) Upgrade to latest version of github.com/wasilibs/go-pgquery (#3052)
- (deps) Bump google.golang.org/grpc from 1.59.0 to 1.60.0 (#3053)
- (deps) Bump babel from 2.13.1 to 2.14.0 in /docs (#3055)
- (deps) Bump actions/upload-artifact from 3 to 4 (#3061)
- (deps) Bump modernc.org/sqlite from 1.27.0 to 1.28.0 (#3062)
- (deps) Bump golang.org/x/crypto from 0.14.0 to 0.17.0 (#3068)
- (deps) Bump google.golang.org/grpc from 1.60.0 to 1.60.1 (#3070)
- (deps) Bump google.golang.org/protobuf from 1.31.0 to 1.32.0 (#3079)
- (deps) Bump github.com/tetratelabs/wazero from 1.5.0 to 1.6.0 (#3096)
- (sqlite) Update to antlr 4.13.1 (#3086)
- (sqlite) Disable modernc for WASM (#3048)
- (sqlite) Switch from mattn/go-sqlite3 to modernc.org/sqlite (#3040)

(v1-24-0)=
## [1.24.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.24.0)
Released 2023-11-22

### Release notes

#### Verifying database schema changes

Schema updates and poorly-written queries often bring down production databases. That’s bad.

Out of the box, `sqlc generate` catches some of these issues. Running `sqlc vet` with the `sqlc/db-prepare` rule catches more subtle problems. But there is a large class of issues that sqlc can’t prevent by looking at current schema and queries alone.

For instance, when a schema change is proposed, existing queries and code running in production might fail when the schema change is applied. Enter `sqlc verify`, which analyzes existing queries against new schema changes and errors if there are any issues.

Let's look at an example. Assume you have these two tables in production.

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY
);

CREATE TABLE user_actions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  action TEXT,
  created_at TIMESTAMP
);
```

Your application contains the following query to join user actions against the users table.

```sql
-- name: GetUserActions :many
SELECT * FROM users u
JOIN user_actions ua ON u.id = ua.user_id
ORDER BY created_at;
```

So far, so good. Then assume you propose this schema change:

```sql
ALTER TABLE users ADD COLUMN created_at TIMESTAMP;
```

Running `sqlc generate` fails with this change, returning a `column reference "created_at" is ambiguous` error. You update your query to fix the issue.

```sql
-- name: GetUserActions :many
SELECT * FROM users u
JOIN user_actions ua ON u.id = ua.user_id
ORDER BY u.created_at;
```

While that change fixes the issue, there's a production outage waiting to happen. When the schema change is applied, the existing `GetUserActions` query will begin to fail. The correct way to fix this is to deploy the updated query before applying the schema migration.

It ensures migrations are safe to deploy by sending your current schema and queries to sqlc cloud. There, we run the queries for your latest push against your new schema changes. This check catches backwards incompatible schema changes for existing queries.

Here `sqlc verify` alerts you to the fact that ORDER BY "created_at" is ambiguous.

```sh
$ sqlc verify
FAIL: app query.sql

=== Failed
=== FAIL: app query.sql GetUserActions
    ERROR: column reference "created_at" is ambiguous (SQLSTATE 42702)
```

By the way, this scenario isn't made up! It happened to us a few weeks ago. We've been happily testing early versions of `verify` for the last two weeks and haven't had any issues since.

This type of verification is only the start. If your application is deployed on-prem by your customers, `verify` could tell you if it's safe for your customers to rollback to an older version of your app, even after schema migrations have been run.

#### Rename `upload` command to `push`

We've renamed the `upload` sub-command to `push`. We changed the data sent along in a push request. Upload used to include the configuration file, migrations, queries, and all generated code. Push drops the generated code in favor of including the [plugin.GenerateRequest](https://buf.build/sqlc/sqlc/docs/main:plugin#plugin.GenerateRequest), which is the protocol buffer message we pass to codegen plugins.

We also add annotations to each push. By default, we include these environment variables if they are present:

```
GITHUB_REPOSITORY
GITHUB_REF
GITHUB_REF_NAME
GITHUB_REF_TYPE
GITHUB_SHA
```

Like upload, `push` should be run when you tag a release of your application. We run it on every push to main, as we continuously deploy those commits.

#### MySQL support in `createdb`

The `createdb` command, added in the last release, now supports MySQL. If you have a cloud project configured, you can use `sqlc createdb` to spin up a new ephemeral database with your schema and print its connection string to standard output. This is useful for integrating with other tools. Read more in the [managed databases](../howto/managed-databases.md#with-other-tools) documentation.

#### Plugin interface refactor

This release includes a refactored plugin interface to better support future functionality. Plugins now support different methods via a gRPC service interface, allowing plugins to support different functionality in a backwards-compatible way.

By using gRPC interfaces, we can even (theoretically) support [remote plugins](https://github.com/sqlc-dev/sqlc/pull/2938), but that's something for another day.

### Changes

#### Bug Fixes

- (engine/sqlite) Support CASE expr (#2926)
- (engine/sqlite) Support -> and ->> operators (#2927)
- (vet) Add a nil pointer check to prevent db/prepare panic (#2934)
- (compiler) Prevent panic when compiler is nil (#2942)
- (codegen/golang) Move more Go-specific config validation into the plugin (#2951)
- (compiler) No panic on full-qualified column names (#2956)
- (docs) Better discussion of type override nuances (#2972)
- (codegen) Never generate return structs for :exec (#2976)
- (generate) Update help text for generate to be more generic (#2981)
- (generate) Return an error instead of generating duplicate Go names (#2962)
- (codegen/golang) Pull opts into its own package (#2920)
- (config) Make some struct and field names less confusing (#2922)

#### Features

- (codegen) Remove Go specific overrides from codegen proto (#2929)
- (plugin) Use gRPC interface for codegen plugin communication (#2930)
- (plugin) Calculate SHA256 if it does not exist (#2935)
- (sqlc-gen-go) Add script to mirror code to sqlc-gen-go (#2952)
- (createdb) Add support for MySQL (#2980)
- (verify) Add new command to verify queries and migrations (#2986)

#### Testing

- (ci) New workflow for sqlc-gen-python (#2936)
- (ci) Rely on go.mod to determine which Go version to use (#2971)
- (tests) Add glob pattern tests to sqlpath.Glob (#2995)
- (examples) Use hosted MySQL databases for tests (#2982)
- (docs) Clean up a little, update LICENSE and README (#2941)

#### Build

- (deps) Bump babel from 2.13.0 to 2.13.1 in /docs (#2911)
- (deps) Bump github.com/spf13/cobra from 1.7.0 to 1.8.0 (#2944)
- (deps) Bump github.com/mattn/go-sqlite3 from 1.14.17 to 1.14.18 (#2945)
- (deps) Bump golang.org/x/sync from 0.4.0 to 0.5.0 (#2946)
- (deps) Bump github.com/jackc/pgx/v5 from 5.4.3 to 5.5.0 (#2947)
- (deps) Change github.com/pingcap/tidb/parser to github.com/pingcap/tidb/pkg/parser
- (deps) Bump github.com/google/cel-go from 0.18.1 to 0.18.2 (#2969)
- (deps) Bump urllib3 from 2.0.7 to 2.1.0 in /docs (#2975)
- (buf) Change root of Buf module (#2987)
- (deps) Bump certifi from 2023.7.22 to 2023.11.17 in /docs (#2993)
- (ci) Bump Go version from 1.21.3 to 1.21.4 in workflows and Dockerfile (#2961)

(v1-23-0)=
## [1.23.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.23.0)
Released 2023-10-24

### Release notes

#### Database-backed query analysis

With a [database connection](config.md#database) configured, `sqlc generate`
will gather metadata from that database to support its query analysis.
Turning this on resolves a [large number of
issues](https://github.com/sqlc-dev/sqlc/issues?q=is%3Aissue+label%3Aanalyzer)
in the backlog related to type inference and more complex queries. The easiest
way to try it out is with [managed databases](../howto/managed-databases.md).

The database-backed analyzer currently supports PostgreSQL, with [MySQL](https://github.com/sqlc-dev/sqlc/issues/2902) and [SQLite](https://github.com/sqlc-dev/sqlc/issues/2903)
support planned in the future.

#### New `createdb` command

When you have a cloud project configured, you can use the new `sqlc createdb`
command to spin up a new ephemeral database with your schema and print its
connection string to standard output. This is useful for integrating with other
tools. Read more in the [managed
databases](../howto/managed-databases.md#with-other-tools) documentation.

#### Support for pgvector

If you're using [pgvector](https://github.com/pgvector/pgvector), say goodbye to custom overrides! sqlc now generates code using [pgvector-go](https://github.com/pgvector/pgvector-go#pgx) as long as you're using `pgx`. The pgvector extension is also available in [managed databases](../howto/managed-databases.md).

#### Go build tags

With the new `emit_build_tags` configuration parameter you can set build tags
for sqlc to add at the top of generated source files.

### Changes

#### Bug Fixes

- (codegen) Correct column names in :copyfrom (#2838)
- (compiler) Search SELECT and UPDATE the same way (#2841)
- (dolphin) Support more UNIONs for MySQL (#2843)
- (compiler) Account for parameters without parents (#2844)
- (postgresql) Remove temporary pool config (#2851)
- (golang) Escape reserved keywords (#2849)
- (mysql) Handle simplified CASE statements (#2852)
- (engine/dolphin) Support enum in ALTER definition (#2680)
- (mysql) Add, drop, rename and change enum values (#2853)
- (config) Validate `database` config in all cases (#2856)
- (compiler) Use correct func signature for `CommentSyntax` on windows (#2867)
- (codegen/go) Prevent filtering of embedded struct fields (#2868)
- (compiler) Support functions with OUT params (#2865)
- (compiler) Pull in array information from analyzer (#2864)
- (analyzer) Error on unexpanded star expression (#2882)
- (vet) Remove rollback statements from DDL (#2895)

#### Documentation

- Add stable anchors to changelog (#2784)
- Fix typo in v1.22.0 changelog (#2796)
- Add sqlc upload to CI / CD guide (#2797)
- Fix broken link, add clarity to plugins doc (#2813)
- Add clarity and reference to JSON tags (#2819)
- Replace form with dashboard link (#2840)
- (examples) Update examples to use pgx/v5 (#2863)
- Use docker compose v2 and update MYSQL_DATABASE env var (#2870)
- Update getting started guides, use pgx for Postgres guide (#2891)
- Use managed databases in PostgreSQL getting started guide (#2892)
- Update managed databases doc to discuss codegen (#2897)
- Add managed dbs to CI/CD and vet guides (#2896)
- Document database-backed query analyzer (#2904)

#### Features

- (codegen) Support setting Go build tags (#2012) (#2807)
- (generate) Reorder codegen handlers to prefer plugins (#2814)
- (devenv) Add vscode settings.json with auto newline (#2834)
- (cmd) Support sqlc.yml configuration file (#2828)
- (analyzer) Analyze queries using a running PostgreSQL database (#2805)
- (sql/ast) Render AST to SQL (#2815)
- (codegen) Include plugin information (#2846)
- (postgresql) Add ALTER VIEW ... SET SCHEMA (#2855)
- (compiler) Parse query parameter metadata from comments (#2850)
- (postgresql) Support system columns on tables (#2871)
- (compiler) Support LEFT JOIN on aliased table (#2873)
- Improve messaging for common cloud config and rpc errors (#2885)
- Abort compiler when rpc fails as unauthenticated (#2887)
- (codegen) Add support for pgvector and pgvector-go (#2888)
- (analyzer) Cache query analysis (#2889)
- (createdb) Create ephemeral databases (#2894)
- (debug) Add databases=managed debug option (#2898)
- (config) Remove managed database validation (#2901)

#### Miscellaneous Tasks

- (endtoend) Fix test output for do tests (#2782)

#### Refactor

- (codegen) Remove golang and json settings from plugin proto (#2822)
- (codegen) Removed deprecated code and improved speed (#2899)

#### Testing

- (endtoend) Split shema and queries (#2803)
- Fix a few incorrect testcases (#2804)
- (analyzer) Add more database analyzer test cases (#2854)
- Add more analyzer test cases (#2866)
- Add more test cases for new analyzer (#2879)
- (endtoend) Enabled managed-db tests in CI (#2883)
- Enabled pgvector tests for managed dbs (#2893)

#### Build

- (deps) Bump packaging from 23.1 to 23.2 in /docs (#2791)
- (deps) Bump urllib3 from 2.0.5 to 2.0.6 in /docs (#2798)
- (deps) Bump babel from 2.12.1 to 2.13.0 in /docs (#2799)
- (deps) Bump golang.org/x/sync from 0.3.0 to 0.4.0 (#2810)
- (deps) Bump golang from 1.21.1 to 1.21.2 (#2811)
- (deps) Bump github.com/google/go-cmp from 0.5.9 to 0.6.0 (#2826)
- (deps) Bump golang from 1.21.2 to 1.21.3 (#2824)
- (deps) Bump google.golang.org/grpc from 1.58.2 to 1.58.3 (#2825)
- (deps) Bump golang.org/x/net from 0.12.0 to 0.17.0 (#2836)
- (deps) Bump urllib3 from 2.0.6 to 2.0.7 in /docs (#2872)
- (deps) Bump google.golang.org/grpc from 1.58.3 to 1.59.0 (#2876)
- (deps) Upgrade wasmtime-go from 13.0.0 to 14.0.0 (#2900)

#### Ci

- Bump go version in workflows (#2835)


(v1-22-0)=
## [1.22.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.22.0)
Released 2023-09-26

### Release notes

#### Managed databases for `sqlc vet`

If you're using [sqlc vet](../howto/vet.md) to write rules that require access to a running
database, `sqlc` can now start and manage that database for you. PostgreSQL
support is available today, with MySQL on the way.

When you turn on managed databases, `sqlc` will use your schema to create a
template database that it can copy to make future runs of `sqlc vet` very
performant.

This feature relies on configuration obtained via [sqlc
Cloud](https://dashboard.sqlc.dev).

Read more in the [managed databases](../howto/managed-databases.md) documentation.

### Changes

#### Bug Fixes

- (codegen/golang) Refactor imports code to match templates (#2709)
- (codegen/golang) Support name type (#2715)
- (wasm) Move Runner struct to shared file (#2725)
- (engine/sqlite) Fix grammer to avoid missing join_constraint (#2732)
- (convert) Support YAML anchors in plugin options (#2733)
- (mysql) Disallow time.Time in mysql :copyfrom queries, not all queries (#2768)
- (engine/sqlite) Fix convert process for VALUES (#2737)

#### Documentation

- Clarify nullable override behavior (#2753)
- Add managed databases to sidebar (#2764)
- Pull renaming and type overrides into separate sections (#2774)
- Update the docs banner for managed dbs (#2775)

#### Features

- (config) Enables the configuration of copyfrom.go similar to quierer and friends (#2727)
- (vet) Run rules against a managed database (#2751)
- (upload) Point upload command at new endpoint (#2772)
- (compiler) Support DO statements (#2777)

#### Miscellaneous Tasks

- (endtoend) Skip tests missing secrets (#2763)
- Skip certain tests on PRs (#2769)

#### Testing

- (endtoend) Verify all schemas in endtoend (#2744)
- (examples) Use a hosted database for example testing (#2749)
- (endtoend) Pull region from environment (#2750)

#### Build

- (deps) Bump golang from 1.21.0 to 1.21.1 (#2711)
- (deps) Bump google.golang.org/grpc from 1.57.0 to 1.58.1 (#2743)
- (deps) Bump wasmtime-go from v12 to v13 (#2756)
- (windows) Downgrade to mingw 11.2.0 (#2757)
- (deps) Bump urllib3 from 2.0.4 to 2.0.5 in /docs (#2747)
- (deps) Bump google.golang.org/grpc from 1.58.1 to 1.58.2 (#2758)
- (deps) Bump github.com/google/cel-go from 0.18.0 to 0.18.1 (#2778)

#### Ci

- Bump go version to latest in ci workflows (#2722)


(v1-21-0)=
## [1.21.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.21.0)
Released 2023-09-06

### Release notes

This is primarily a bugfix release, along with some documentation and testing improvements.

#### MySQL engine improvements

`sqlc` previously didn't know how to parse a `CALL` statement when using the MySQL engine,
which meant it was impossible to use sqlc with stored procedures in MySQL databases.

Additionally, `sqlc` now supports `IS [NOT] NULL` in queries. And `LIMIT` and `OFFSET` clauses
now work with `UNION`.

#### SQLite engine improvements

GitHub user [@orisano](https://github.com/orisano) continues to bring bugfixes and
improvements to `sqlc`'s SQLite engine. See the "Changes" section below for the
full list.

#### Plugin access to environment variables

If you're authoring a [sqlc plugin](../guides/plugins.html), you can now configure
sqlc to pass your plugin the values of specific environment variables.

For example, if your plugin
needs the `PATH` environment variable, add `PATH` to the `env` list in the
`plugins` collection.

```yaml
version: '2'
sql:
- schema: schema.sql
  queries: query.sql
  engine: postgresql
  codegen:
  - out: gen
    plugin: test
plugins:
- name: test
  env:
  - PATH
  wasm:
    url: https://github.com/sqlc-dev/sqlc-gen-test/releases/download/v0.1.0/sqlc-gen-test.wasm
    sha256: 138220eae508d4b65a5a8cea555edd155eb2290daf576b7a8b96949acfeb3790
```

A variable named `SQLC_VERSION` is always included in the plugin's
environment, set to the version of the `sqlc` executable invoking it.

### Changes

#### Bug Fixes

- Myriad string formatting changes (#2558)
- (engine/sqlite) Support quoted identifier (#2556)
- (engine/sqlite) Fix compile error (#2564)
- (engine/sqlite) Fixed detection of column alias without AS (#2560)
- (ci) Bump go version to 1.20.7 (#2568)
- Remove references to deprecated `--experimental` flag (#2567)
- (postgres) Fixed a problem with array dimensions disappearing when using "ALTER TABLE ADD COLUMN" (#2572)
- Remove GitHub sponsor integration (#2574)
- (docs) Improve discussion of prepared statements support (#2604)
- (docs) Remove multidimensional array qualification in datatypes.md (#2619)
- (config) Go struct tag parsing (#2606)
- (compiler) Fix to not scan children under ast.RangeSubselect when retrieving table listing (#2573)
- (engine/sqlite) Support NOT IN (#2587)
- (codegen/golang) Fixed detection of the used package (#2597)
- (engine/dolphin) Fixed problem that LIMIT OFFSET cannot be used with `UNION ALL` (#2613)
- (compiler) Support identifiers with schema (#2579)
- (compiler) Fix column expansion to work with quoted non-keyword identifiers (#2576)
- (codegen/go) Compare define type in codegen (#2263) (#2578)
- (engine/sqlite) Fix ast when using compound operator (#2673)
- (engine/sqlite) Fix to handle join clauses correctly (#2674)
- (codegen) Use correct Go types for bit strings and cid/oid/tid/xid with pgx/v4 (#2668)
- (endtoend) Ensure all SQL works against PostgreSQL (#2684)

#### Documentation

- Update Docker installation instructions (#2552)
- Missing emit_pointers_for_null_types configuration option in version 2 (#2682) (#2683)
- Fix typo (#2697)
- Document sqlc.* macros (#2698)
- (mysql) Document parseTimet=true requirement (#2699)
- Add atlas to the list of supported migration frameworks (#2700)
- Minor updates to insert howto (#2701)

#### Features

- (endtoend/testdata) Added two sqlite `CAST` tests and rearranged postgres tests for same (#2551)
- (docs) Add a reference to type overriding in datatypes.md (#2557)
- (engine/sqlite) Support COLLATE for sqlite WHERE clause (#2554)
- (mysql) Add parser support for IS [NOT] NULL (#2651)
- (engine/dolphin) Support CALL statement (#2614)
- (codegen) Allow plugins to access environment variables (#2669)
- (config) Add JSON schema files for configs (#2703)

#### Miscellaneous Tasks

- Ignore Vim swap files (#2616)
- Fix typo (#2696)

#### Refactor

- (astutils) Remove redundant nil check in `Walk` (#2660)

#### Build

- (deps) Bump wasmtime from v8.0.0 to v11.0.0 (#2553)
- (deps) Bump golang from 1.20.6 to 1.20.7 (#2563)
- (deps) Bump chardet from 5.1.0 to 5.2.0 in /docs (#2562)
- (deps) Bump github.com/pganalyze/pg_query_go/v4 (#2583)
- (deps) Bump golang from 1.20.7 to 1.21.0 (#2596)
- (deps) Bump github.com/jackc/pgx/v5 from 5.4.2 to 5.4.3 (#2582)
- (deps) Bump pygments from 2.15.1 to 2.16.1 in /docs (#2584)
- (deps) Bump sphinxcontrib-applehelp from 1.0.4 to 1.0.7 in /docs (#2620)
- (deps) Bump sphinxcontrib-qthelp from 1.0.3 to 1.0.6 in /docs (#2622)
- (deps) Bump github.com/google/cel-go from 0.17.1 to 0.17.6 (#2650)
- (deps) Bump sphinxcontrib-serializinghtml in /docs (#2641)
- Upgrade from Go 1.20 to Go 1.21 (#2665)
- (deps) Bump sphinxcontrib-devhelp from 1.0.2 to 1.0.5 in /docs (#2621)
- (deps) Bump github.com/bytecodealliance/wasmtime-go from v11.0.0 to v12.0.0 (#2666)
- (deps) Bump sphinx-rtd-theme from 1.2.2 to 1.3.0 in /docs (#2670)
- (deps) Bump sphinxcontrib-htmlhelp from 2.0.1 to 2.0.4 in /docs (#2671)
- (deps) Bump github.com/google/cel-go from 0.17.6 to 0.18.0 (#2691)
- (deps) Bump actions/checkout from 3 to 4 (#2694)
- (deps) Bump pytz from 2023.3 to 2023.3.post1 in /docs (#2695)
- (devenv) Bump go from 1.20.7 to 1.21.0 (#2702)

(v1-20-0)=
## [1.20.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.20.0)
Released 2023-07-31

### Release notes

#### `kyleconroy/sqlc` is now `sqlc-dev/sqlc`

We've completed our migration to the [sqlc-dev/sqlc](https://github.com/sqlc-dev/sqlc) repository. All existing links and installation instructions will continue to work. If you're using the `go` tool to install `sqlc`, you'll need to use the new import path to get v1.20.0 (and all future versions).

```sh
# INCORRECT: old import path
go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.20.0

# CORRECT: new import path
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.20.0
```

We designed the upgrade process to be as smooth as possible. If you run into any issues, please [file a bug report](https://github.com/sqlc-dev/sqlc/issues/new?assignees=&labels=bug%2Ctriage&projects=&template=BUG_REPORT.yml) via GitHub.

#### Use `EXPLAIN ...` output in lint rules

`sqlc vet` can now run `EXPLAIN` on your queries and include the results for use in your lint rules. For example, this rule checks that `SELECT` queries use an index.

```yaml
version: 2
sql:
  - schema: "query.sql"
    queries: "query.sql"
    engine: "postgresql"
    database:
      uri: "postgresql://postgres:postgres@localhost:5432/postgres"
    gen:
      go:
        package: "db"
        out: "db"
    rules:
      - has-index
rules:
- name: has-index
  rule: >
    query.sql.startsWith("SELECT") &&
    !(postgresql.explain.plan.plans.all(p, has(p.index_name) || p.plans.all(p, has(p.index_name))))
```

The expression environment has two variables containing `EXPLAIN ...` output, `postgresql.explain` and `mysql.explain`. `sqlc` only populates the variable associated with your configured database engine, and only when you have a [database connection configured](../reference/config.md#database).

For the `postgresql` engine, `sqlc` runs

```sql
EXPLAIN (ANALYZE false, VERBOSE, COSTS, SETTINGS, BUFFERS, FORMAT JSON) ...
```

where `"..."` is your query string, and parses the output into a [`PostgreSQLExplain`](https://buf.build/sqlc/sqlc/docs/v1.20.0:vet#vet.PostgreSQLExplain) proto message.

For the `mysql` engine, `sqlc` runs

```sql
EXPLAIN FORMAT=JSON ...
```

where `"..."` is your query string, and parses the output into a [`MySQLExplain`](https://buf.build/sqlc/sqlc/docs/v1.20.0:vet#vet.MySQLExplain) proto message.

These proto message definitions are too long to include here, but you can find them in the `protos` directory within the `sqlc` source tree.

The output from `EXPLAIN ...` depends on the structure of your query so it's a bit difficult to offer generic examples. Refer to the [PostgreSQL documentation](https://www.postgresql.org/docs/current/using-explain.html) and [MySQL documentation](https://dev.mysql.com/doc/refman/en/explain-output.html) for more information.

```yaml
...
rules:
- name: postgresql-query-too-costly
  message: "Query cost estimate is too high"
  rule: "postgresql.explain.plan.total_cost > 1.0"
- name: postgresql-no-seq-scan
  message: "Query plan results in a sequential scan"
  rule: "postgresql.explain.plan.node_type == 'Seq Scan'"
- name: mysql-query-too-costly
  message: "Query cost estimate is too high"
  rule: "has(mysql.explain.query_block.cost_info) && double(mysql.explain.query_block.cost_info.query_cost) > 2.0"
- name: mysql-must-use-primary-key
  message: "Query plan doesn't use primary key"
  rule: "has(mysql.explain.query_block.table.key) && mysql.explain.query_block.table.key != 'PRIMARY'"
```

When building rules that depend on `EXPLAIN ...` output, it may be helpful to see the actual JSON returned from the database. `sqlc` will print it When you set the environment variable `SQLCDEBUG=dumpexplain=1`. Use this environment variable together with a dummy rule to see `EXPLAIN ...` output for all of your queries.

#### Opting-out of lint rules

For any query, you can tell `sqlc vet` not to evaluate lint rules using the `@sqlc-vet-disable` query annotation.

```sql
/* name: GetAuthor :one */
/* @sqlc-vet-disable */
SELECT * FROM authors
WHERE id = ? LIMIT 1;
```

#### Bulk insert for MySQL

_Developed by [@Jille](https://github.com/Jille)_

MySQL now supports the `:copyfrom` query annotation. The generated code uses the [LOAD DATA](https://dev.mysql.com/doc/refman/8.0/en/load-data.html) command to insert data quickly and efficiently.

Use caution with this feature. Errors and duplicate keys are treated as warnings and insertion will continue, even without an error for some cases.  Use this in a transaction and use `SHOW WARNINGS` to check for any problems and roll back if necessary.

Check the [error handling](https://dev.mysql.com/doc/refman/8.0/en/load-data.html#load-data-error-handling) documentation for more information.

```sql
CREATE TABLE foo (a text, b integer, c DATETIME, d DATE);

-- name: InsertValues :copyfrom
INSERT INTO foo (a, b, c, d) VALUES (?, ?, ?, ?);
```

```go
func (q *Queries) InsertValues(ctx context.Context, arg []InsertValuesParams) (int64, error) {
	...
}
```

`LOAD DATA` support must be enabled in the MySQL server.

#### CAST support for MySQL

_Developed by [@ryanpbrewster](https://github.com/ryanpbrewster) and [@RadhiFadlillah](https://github.com/RadhiFadlillah)_

`sqlc` now understands `CAST` calls in MySQL queries, offering greater flexibility when generating code for complex queries.

```sql
CREATE TABLE foo (bar BOOLEAN NOT NULL);

-- name: SelectColumnCast :many
SELECT CAST(bar AS BIGINT) FROM foo;
```

```go
package querytest

import (
	"context"
)

const selectColumnCast = `-- name: SelectColumnCast :many
SELECT CAST(bar AS BIGINT) FROM foo
`

func (q *Queries) SelectColumnCast(ctx context.Context) ([]int64, error) {
  ...
}
```

#### SQLite improvements

A slew of fixes landed for our SQLite implementation, bringing it closer to parity with MySQL and PostgreSQL. We want to thank [@orisano](https://github.com/orisano) for their continued dedication to improving `sqlc`'s SQLite support.

### Changes

#### Features

- (debug) Add debug flag and docs for dumping vet rule variables (#2521)
- (mysql) :copyfrom support via LOAD DATA INFILE (#2545)
- (mysql) Implement cast function parser (#2473)
- (postgresql) Add support for PostgreSQL multi-dimensional arrays (#2338)
- (sql/catalog) Support ALTER TABLE IF EXISTS (#2542)
- (sqlite) Virtual tables and fts5 supported (#2531)
- (vet) Add default query parameters for explain queries (#2543)
- (vet) Add output from `EXPLAIN ...` for queries to the CEL program environment (#2489)
- (vet) Introduce a query annotation to opt out of sqlc vet rules (#2474)
- Parse comment lines starting with `@symbol` as boolean flags associated with a query (#2464)

#### Bug Fixes

- (codegen/golang) Fix sqlc.embed to work with pq.Array (#2544)
- (compiler) Correctly validate alias in order/group by clauses for joins (#2537)
- (engine/sqlite) Added function to convert cast node (#2470)
- (engine/sqlite) Fix join_operator rule (#2434)
- (engine/sqlite) Fix table_alias rules (#2465)
- (engine/sqlite) Fixed IN operator precedence (#2428)
- (engine/sqlite) Fixed to be able to find relation from WITH clause (#2444)
- (engine/sqlite) Lowercase ast.ResTarget.Name (#2433)
- (engine/sqlite) Put logging statement behind debug flag (#2488)
- (engine/sqlite) Support for repeated table_option (#2482)
- (mysql) Generate unsigned param (#2522)
- (sql/catalog) Support pg_dump output (#2508)
- (sqlite) Code generation for sqlc.slice (#2431)
- (vet) Clean up unnecessary `prepareable()` func and a var name (#2509)
- (vet) Query.cmd was always set to ":" (#2525)
- (vet) Report an error when a query is unpreparable, close prepared statement connection (#2486)
- (vet) Split vet messages out of codegen.proto (#2511)

#### Documentation

- Add a description to the document for cases when a query result has no rows (#2462)
- Update copyright and author (#2490)
- Add example sqlc.yaml for migration parsing (#2479)
- Small updates (#2506)
- Point GitHub links to new repository location (#2534)

#### Miscellaneous Tasks

- Rename kyleconroy/sqlc to sqlc-dev/sqlc (#2523)
- (proto) Reformat protos using `buf format -w` (#2536)
- Update FEATURE_REQUEST.yml to include SQLite engine option
- Finish migration to sqlc-dev/sqlc (#2548)
- (compiler) Remove some duplicate code (#2546)

#### Testing

- Add profiles to docker compose (#2503)

#### Build

- Run all supported versions of MySQL / PostgreSQL (#2463)
- (deps) Bump pygments from 2.7.4 to 2.15.0 in /docs (#2485)
- (deps) Bump github.com/jackc/pgconn from 1.14.0 to 1.14.1 (#2483)
- (deps) Bump github.com/google/cel-go from 0.16.0 to 0.17.1 (#2484)
- (docs) Check Python dependencies via dependabot (#2497)
- (deps) Bump idna from 2.10 to 3.4 in /docs (#2499)
- (deps) Bump packaging from 20.9 to 23.1 in /docs (#2498)
- (deps) Bump pygments from 2.15.0 to 2.15.1 in /docs (#2500)
- (deps) Bump certifi from 2022.12.7 to 2023.7.22 in /docs (#2504)
- (deps) Bump sphinx from 4.4.0 to 6.1.0 in /docs (#2505)
- Add psql and mysqlsh to devenv (#2507)
- (deps) Bump urllib3 from 1.26.5 to 2.0.4 in /docs (#2516)
- (deps) Bump chardet from 4.0.0 to 5.1.0 in /docs (#2517)
- (deps) Bump snowballstemmer from 2.1.0 to 2.2.0 in /docs (#2519)
- (deps) Bump pytz from 2021.1 to 2023.3 in /docs (#2520)
- (deps) Bump sphinxcontrib-htmlhelp from 2.0.0 to 2.0.1 in /docs (#2518)
- (deps) Bump pyparsing from 2.4.7 to 3.1.0 in /docs (#2530)
- (deps) Bump alabaster from 0.7.12 to 0.7.13 in /docs (#2526)
- (docs) Ignore updates for sphinx (#2532)
- (deps) Bump babel from 2.9.1 to 2.12.1 in /docs (#2527)
- (deps) Bump sphinxcontrib-applehelp from 1.0.2 to 1.0.4 in /docs (#2533)
- (deps) Bump google.golang.org/grpc from 1.56.2 to 1.57.0 (#2535)
- (deps) Bump pyparsing from 3.1.0 to 3.1.1 in /docs (#2547)


## [1.19.1](https://github.com/sqlc-dev/sqlc/releases/tag/v1.19.1)
Released 2023-07-13

### Bug Fixes

- Fix to traverse Sel in ast.In (#2414)
- (compiler) Validate UNION ... ORDER BY (#2446)
- (golang) Prevent duplicate enum output (#2447)

### Miscellaneous Tasks

- Replace codegen, test and docs references to github.com/tabbed repos (#2418)

### Build

- (deps) Bump google.golang.org/grpc from 1.56.1 to 1.56.2 (#2415)
- (deps) Bump golang from 1.20.5 to 1.20.6 (#2437)
- Pin Go to 1.20.6 (#2441)
- (deps) Bump github.com/jackc/pgx/v5 from 5.4.1 to 5.4.2 (#2436)

## [1.19.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.19.0)
Released 2023-07-06

### Release notes

#### sqlc vet

[`sqlc vet`](../howto/vet.md) runs queries through a set of lint rules.

Rules are defined in the `sqlc` [configuration](config.md) file. They consist
of a name, message, and a [Common Expression Language (CEL)](https://github.com/google/cel-spec)
expression. Expressions are evaluated using [cel-go](https://github.com/google/cel-go).
If an expression evaluates to `true`, an error is reported using the given message.

While these examples are simplistic, they give you a flavor of the types of
rules you can write.

```yaml
version: 2
sql:
  - schema: "query.sql"
    queries: "query.sql"
    engine: "postgresql"
    gen:
      go:
        package: "authors"
        out: "db"
    rules:
      - no-pg
      - no-delete
      - only-one-param
      - no-exec
rules:
  - name: no-pg
    message: "invalid engine: postgresql"
    rule: |
      config.engine == "postgresql"
  - name: no-delete
    message: "don't use delete statements"
    rule: |
      query.sql.contains("DELETE")
  - name: only-one-param
    message: "too many parameters"
    rule: |
      query.params.size() > 1
  - name: no-exec
    message: "don't use exec"
    rule: |
      query.cmd == "exec"
```

##### Database connectivity

`vet` also marks the first time that `sqlc` can connect to a live, running
database server. We'll expand this functionality over time, but for now it
powers the `sqlc/db-prepare` built-in rule.

When a [database](config.html#database) is configured, the
`sqlc/db-preapre` rule will attempt to prepare each of your
queries against the connected database and report any failures.

```yaml
version: 2
sql:
  - schema: "query.sql"
    queries: "query.sql"
    engine: "postgresql"
    gen:
      go:
        package: "authors"
        out: "db"
    database:
      uri: "postgresql://postgres:password@localhost:5432/postgres"
    rules:
      - sqlc/db-prepare
```

To see this in action, check out the [authors
example](https://github.com/sqlc-dev/sqlc/blob/main/examples/authors/sqlc.yaml).

Please note that `sqlc` does not manage or migrate your database. Use your
migration tool of choice to create the necessary database tables and objects
before running `sqlc vet`.

#### Omit unused structs

Added a new configuration parameter `omit_unused_structs` which, when set to
true, filters out table and enum structs that aren't used in queries for a given
package.

#### Suggested CI/CD setup

With the addition of `sqlc diff` and `sqlc vet`, we encourage users to run sqlc
in your CI/CD pipelines. See our [suggested CI/CD setup](../howto/ci-cd.md) for
more information.

#### Simplified plugin development

The [sqlc-gen-kotlin](https://github.com/sqlc-dev/sqlc-gen-kotlin) and
[sqlc-gen-python](https://github.com/sqlc-dev/sqlc-gen-python) plugins have been
updated use the upcoming [WASI](https://wasi.dev/) support in [Go
1.21](https://tip.golang.org/doc/go1.21#wasip1). Building these plugins no
longer requires [TinyGo](https://tinygo.org/).

### Changes

#### Bug Fixes

- Pointers overrides skip imports in generated query files (#2240)
- CASE-ELSE clause is not properly parsed when a value is constant (#2238)
- Fix toSnakeCase to handle input in CamelCase format (#2245)
- Add location info to sqlite ast (#2298)
- Add override tags to result struct (#1867) (#1887)
- Override types of aliased columns and named parameters (#1884)
- Resolve duplicate fields generated when inheriting multiple tables (#2089)
- Check column references in ORDER BY (#1411) (#1915)
- MySQL slice shadowing database/sql import (#2332)
- Don't defer rows.Close() if pgx.BatchResults.Query() failed  (#2362)
- Fix type overrides not working with sqlc.slice (#2351)
- Type overrides on columns for parameters inside an IN clause (#2352)
- Broken interaction between query_parameter_limit and pq.Array() (#2383)
- (codegen/golang) Bring :execlastid in line with the rest (#2378)

#### Documentation

- Update changelog.md with some minor edits (#2235)
- Add F# community plugin (#2295)
- Add a ReadTheDocs config file (#2327)
- Update query_parameter_limit documentation (#2374)
- Add launch announcement banner

#### Features
- PostgreSQL capture correct line and column numbers for parse error (#2289)
- Add supporting COMMENT ON VIEW (#2249)
- To allow spaces between function name and arguments of functions to be rewritten (#2250)
- Add support for pgx/v5 emit_pointers_for_null_types flag (#2269)
- (mysql) Support unsigned integers (#1746)
- Allow use of table and column aliases for table functions returning unknown types (#2156)
- Support "LIMIT ?" in UPDATE and DELETE for MySQL (#2365)
- (internal/codegen/golang) Omit unused structs from output (#2369)
- Improve default names for BETWEEN ? AND ? to have prefixes from_ and to_ (#2366)
- (cmd/sqlc) Add the vet subcommand (#2344)
- (sqlite) Add support for UPDATE/DELETE with a LIMIT clause (#2384)
- Add support for BETWEEN sqlc.arg(min) AND sqlc.arg(max) (#2373)
- (cmd/vet) Prepare queries against a database (#2387)
- (cmd/vet) Prepare queries for MySQL (#2388)
- (cmd/vet) Prepare SQLite queries (#2389)
- (cmd/vet) Simplify environment variable substiution (#2393)
- (cmd/vet) Add built-in db-prepare rule
- Add compiler support for NOTIFY and LISTEN (PostgreSQL) (#2363)

#### Miscellaneous Tasks

- A few small staticcheck fixes (#2361)
- Remove a bunch of dead code (#2360)
- (scripts/regenerate) Should also update stderr.txt (#2379)

#### Build

- (deps) Bump requests from 2.25.1 to 2.31.0 in /docs (#2283)
- (deps) Bump golang from 1.20.3 to 1.20.4 (#2256)
- (deps) Bump google.golang.org/grpc from 1.54.0 to 1.55.0 (#2265)
- (deps) Bump github.com/mattn/go-sqlite3 from 1.14.16 to 1.14.17 (#2293)
- (deps) Bump golang.org/x/sync from 0.1.0 to 0.2.0 (#2266)
- (deps) Bump golang from 1.20.4 to 1.20.5 (#2301)
- Configure dependencies via devenv.sh (#2319)
- Configure dependencies via devenv.sh (#2326)
- (deps) Bump golang.org/x/sync from 0.2.0 to 0.3.0 (#2328)
- (deps) Bump google.golang.org/grpc from 1.55.0 to 1.56.0 (#2333)
- (deps) Bump google.golang.org/protobuf from 1.30.0 to 1.31.0 (#2370)
- (deps) Bump actions/checkout from 2 to 3 (#2357)
- Run govulncheck on all builds (#2372)
- (deps) Bump google.golang.org/grpc from 1.56.0 to 1.56.1 (#2358)

#### Cmd/sqlc

- Show helpful output on missing subcommand (#2345)

#### Codegen

- Use catalog's default schema (#2310)
- (go) Add tests for tables with dashes (#2312)
- (go) Strip invalid characters from table and column names (#2314)
- (go) Support JSON tags for nullable enum structs (#2121)

#### Internal/config

- Support golang overrides for slices (#2339)

#### Kotlin

- Use latest version of sqlc-gen-kotlin (#2356)

#### Postgres

- Column merging for table inheritence (#2315)

#### Protos

- Add missing field name (#2354)

#### Python

- Use latest version of sqlc-gen-python (#2355)

#### Remote

- Use user-id/password auth (#2262)

#### Sqlite

- Fixed sqlite column type override (#1986)


## [1.18.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.18.0)
Released 2023-04-27

### Release notes

#### Remote code generation

_Developed by [@andrewmbenton](https://github.com/andrewmbenton)_

At its core, sqlc is powered by SQL engines, which include parsers, formatters,
analyzers and more. While our goal is to support each engine on each operating
system, it's not always possible. For example, the PostgreSQL engine does not
work on Windows.

To bridge that gap, we're announcing remote code generation, currently in
private alpha. To join the private alpha, [sign up for the waitlist](https://docs.google.com/forms/d/e/1FAIpQLScDWrGtTgZWKt3mdlF5R2XCX6tL1pMkB4yuZx5yq684tTNN1Q/viewform?usp=sf_link).

Remote code generation works like local code generation, except the heavy
lifting is performed in a consistent cloud environment. WASM-based plugins are
supported in the remote environment, but process-based plugins are not.

To configure remote generation, add a `cloud` block in `sqlc.json`.

```json
{
  "version": "2",
  "cloud": {
    "organization": "<org-id>",
    "project": "<project-id>",
  },
  ...
}
```

You'll also need to set the `SQLC_AUTH_TOKEN` environment variable.

```bash
export SQLC_AUTH_TOKEN=<token>
```

When the `cloud` configuration block exists, `sqlc generate` will default to remote
code generation. If you'd like to generate code locally without removing the `cloud`
block from your config, pass the `--no-remote` option.


```bash
sqlc generate --no-remote
```

Remote generation is off by default and requires an opt-in to use.

#### sqlc.embed

_Developed by [@nickjackson](https://github.com/nickjackson)_

Embedding allows you to reuse existing model structs in more queries, resulting
in less manual serialization work. First, imagine we have the following schema
with students and test scores.


```sql
CREATE TABLE students (
  id   bigserial PRIMARY KEY,
  name text,
  age  integer
)

CREATE TABLE test_scores (
  student_id bigint,
  score integer,
  grade text
)
```

We want to select the student record and the highest score they got on a test.
Here's how we'd usually do that:

```sql
-- name: HighScore :many
WITH high_scores AS (
  SELECT student_id, max(score) as high_score
  FROM test_scores
  GROUP BY 1
)
SELECT students.*, high_score::integer
FROM students
JOIN high_scores ON high_scores.student_id = students.id;
```

When using Go, sqlc will produce a struct like this:

```
type HighScoreRow struct {
	ID        int64
	Name      sql.NullString
	Age       sql.NullInt32
	HighScore int32
}
```

With embedding, the struct will contain a model for the table instead of a
flattened list of columns.

```sql
-- name: HighScoreEmbed :many
WITH high_scores AS (
  SELECT student_id, max(score) as high_score
  FROM test_scores
  GROUP BY 1
)
SELECT sqlc.embed(students), high_score::integer
FROM students
JOIN high_scores ON high_scores.student_id = students.id;
```

```
type HighScoreRow struct {
	Student   Student
	HighScore int32
}
```

#### sqlc.slice

_Developed by Paul Cameron and Jille Timmermans_

The MySQL Go driver does not support passing slices to the IN operator. The
`sqlc.slice` function generates a dynamic query at runtime with the correct
number of parameters.

```sql
/* name: SelectStudents :many */
SELECT * FROM students 
WHERE age IN (sqlc.slice("ages"))
```

```go
func (q *Queries) SelectStudents(ctx context.Context, ages []int32) ([]Student, error) {
```

This feature is only supported in MySQL and cannot be used with prepared
queries.

#### Batch operation improvements  

When using batches with pgx, the error returned when a batch is closed is
exported by the generated package. This change allows for cleaner error
handling using `errors.Is`.

```go
errors.Is(err, generated_package.ErrBatchAlreadyClosed)
```

Previously, you would have had to check match on the error message itself.

```
err.Error() == "batch already closed"
```

The generated code for batch operations always lived in `batch.go`. This file
name can now be configured via the `output_batch_file_name` configuration
option.

#### Configurable query parameter limits for Go

By default, sqlc will limit Go functions to a single parameter. If a query
includes more than one parameter, the generated method will use an argument
struct instead of positional arguments. This behavior can now be changed via
the `query_parameter_limit` configuration option.  If set to `0`, every
genreated method will use a argument struct. 

### Changes

#### Bug Fixes

- Prevent variable redeclaration in single param conflict for pgx (#2058)
- Retrieve Larg/Rarg join query after inner join (#2051)
- Rename argument when conflicted to imported package (#2048)
- Pgx closed batch return pointer if need #1959 (#1960)
- Correct singularization of "waves" (#2194)
- Honor Package level renames in v2 yaml config (#2001)
- (mysql) Prevent UPDATE ... JOIN panic #1590 (#2154)
- Mysql delete join panic (#2197)
- Missing import with pointer overrides, solves #2168 #2125 (#2217)

#### Documentation

- (config.md) Add `sqlite` as engine option (#2164)
- Add first pass at pgx documentation (#2174)
- Add missed configuration option (#2188)
- `specifies parameter ":one" without containing a RETURNING clause` (#2173)

#### Features

- Add `sqlc.embed` to allow model re-use (#1615)
- (Go) Add query_parameter_limit conf to codegen (#1558)
- Add remote execution for codegen (#2214)

#### Testing

- Skip tests if required plugins are missing (#2104)
- Add tests for reanme fix in v2 (#2196)
- Regenerate batch output for filename tests
- Remove remote test (#2232)
- Regenerate test output

#### Bin/sqlc

- Add SQLCTMPDIR environment variable (#2189)

#### Build

- (deps) Bump github.com/antlr/antlr4/runtime/Go/antlr (#2109)
- (deps) Bump github.com/jackc/pgx/v4 from 4.18.0 to 4.18.1 (#2119)
- (deps) Bump golang from 1.20.1 to 1.20.2 (#2135)
- (deps) Bump google.golang.org/protobuf from 1.28.1 to 1.29.0 (#2137)
- (deps) Bump google.golang.org/protobuf from 1.29.0 to 1.29.1 (#2143)
- (deps) Bump golang from 1.20.2 to 1.20.3 (#2192)
- (deps) Bump actions/setup-go from 3 to 4 (#2150)
- (deps) Bump google.golang.org/protobuf from 1.29.1 to 1.30.0 (#2151)
- (deps) Bump github.com/spf13/cobra from 1.6.1 to 1.7.0 (#2193)
- (deps) Bump github.com/lib/pq from 1.10.7 to 1.10.8 (#2211)
- (deps) Bump github.com/lib/pq from 1.10.8 to 1.10.9 (#2229)
- (deps) Bump github.com/go-sql-driver/mysql from 1.7.0 to 1.7.1 (#2228)

#### Cmd/sqlc

- Remove --experimental flag (#2170)
- Add option to disable process-based plugins (#2180)
- Bump version to v1.18.0

#### Codegen

- Correctly generate CopyFrom columns for single-column copyfroms (#2185)

#### Config

- Add top-level cloud configuration (#2204)

#### Engine/postgres

- Upgrade to pg_query_go/v4 (#2114)

#### Ext/wasm

- Check exit code on returned error (#2223)

#### Parser

- Generate correct types for `SELECT NOT EXISTS` (#1972)

#### Sqlite

- Add support for CREATE TABLE ... STRICT (#2175)

#### Wasm

- Upgrade to wasmtime v8.0.0 (#2222)

## [1.17.2](https://github.com/sqlc-dev/sqlc/releases/tag/v1.17.2)
Released 2023-02-22

### Bug Fixes

- Fix build on Windows (#2102)

## [1.17.1](https://github.com/sqlc-dev/sqlc/releases/tag/v1.17.1)
Released 2023-02-22

### Bug Fixes

- Prefer to use []T over pgype.Array[T] (#2090)
- Revert changes to Dockerfile (#2091)
- Do not throw error when IF NOT EXISTS is used on ADD COLUMN (#2092)

### MySQL

- Add `float` support to MySQL (#2097)

### Build

- (deps) Bump golang from 1.20.0 to 1.20.1 (#2082)

## [1.17.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.17.0)
Released 2023-02-13

### Bug Fixes

- Initialize generated code outside function (#1850)
- (engine/mysql) Take into account column's charset to distinguish text/blob, (var)char/(var)binary (#776) (#1895)
- The enum Value method returns correct type (#1996)
- Documentation for Inserting Rows (#2034)
- Add import statements even if only pointer types exist (#2046)
- Search from Rexpr if not found from Lexpr (#2056)

### Documentation

- Change ENTRYPOINT to CMD (#1943)
- Update samples for HOW-TO GUIDES (#1953)

### Features

- Add the diff command (#1963)

### Build

- (deps) Bump github.com/mattn/go-sqlite3 from 1.14.15 to 1.14.16 (#1913)
- (deps) Bump github.com/spf13/cobra from 1.6.0 to 1.6.1 (#1909)
- Fix devcontainer (#1942)
- Run sqlc-pg-gen via GitHub Actions (#1944)
- Move large arrays out of functions (#1947)
- Fix conflicts from pointer configuration (#1950)
- (deps) Bump github.com/go-sql-driver/mysql from 1.6.0 to 1.7.0 (#1988)
- (deps) Bump github.com/jackc/pgtype from 1.12.0 to 1.13.0 (#1978)
- (deps) Bump golang from 1.19.3 to 1.19.4 (#1992)
- (deps) Bump certifi from 2020.12.5 to 2022.12.7 in /docs (#1993)
- (deps) Bump golang from 1.19.4 to 1.19.5 (#2016)
- (deps) Bump golang from 1.19.5 to 1.20.0 (#2045)
- (deps) Bump github.com/jackc/pgtype from 1.13.0 to 1.14.0 (#2062)
- (deps) Bump github.com/jackc/pgx/v4 from 4.17.2 to 4.18.0 (#2063)

### Cmd

- Generate packages in parallel (#2026)

### Cmd/sqlc

- Bump version to v1.17.0

### Codegen

- Remove built-in Kotlin support (#1935)
- Remove built-in Python support (#1936)

### Internal/codegen

- Cache pattern matching compilations (#2028)

### Mysql

- Add datatype tests (#1948)
- Fix blob tests (#1949)

### Plugins

- Upgrade to wasmtime 3.0.1 (#2009)

### Sqlite

- Supported between expr (#1958) (#1967)

### Tools

- Regenerate scripts skips dirs that contains diff exec command (#1987)

### Wasm

- Upgrade to wasmtime 5.0.0 (#2065)

## [1.16.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.16.0)
Released 2022-11-09


### Bug Fixes

- (validate) Sqlc.arg & sqlc.narg are not "missing" (#1814)
- Emit correct comment for nullable enums (#1819)
- 🐛 Correctly switch `coalesce()` result `.NotNull` value (#1664)
- Prevent batch infinite loop with arg length (#1794)
- Support version 2 in error message (#1839)
- Handle empty column list in postgresql (#1843)
- Batch imports filter queries, update cmds having ret type (#1842)
- Named params contribute to batch parameter count (#1841)

### Documentation

- Add a getting started guide for SQLite (#1798)
- Various readability improvements (#1854)
- Add documentation for codegen plugins (#1904)
- Update migration guides with links (#1933)

### Features

- Add HAVING support to MySQL (#1806)

### Miscellaneous Tasks

- Upgrade wasmtime version (#1827)
- Bump wasmtime version to v1.0.0 (#1869)

### Build

- (deps) Bump github.com/jackc/pgconn from 1.12.1 to 1.13.0 (#1785)
- (deps) Bump github.com/mattn/go-sqlite3 from 1.14.13 to 1.14.15 (#1799)
- (deps) Bump github.com/jackc/pgx/v4 from 4.16.1 to 4.17.0 (#1786)
- (deps) Bump github.com/jackc/pgx/v4 from 4.17.0 to 4.17.1 (#1825)
- (deps) Bump github.com/bytecodealliance/wasmtime-go (#1826)
- (deps) Bump github.com/jackc/pgx/v4 from 4.17.1 to 4.17.2 (#1831)
- (deps) Bump golang from 1.19.0 to 1.19.1 (#1834)
- (deps) Bump github.com/google/go-cmp from 0.5.8 to 0.5.9 (#1838)
- (deps) Bump github.com/lib/pq from 1.10.6 to 1.10.7 (#1835)
- (deps) Bump github.com/bytecodealliance/wasmtime-go (#1857)
- (deps) Bump github.com/spf13/cobra from 1.5.0 to 1.6.0 (#1893)
- (deps) Bump golang from 1.19.1 to 1.19.3 (#1920)

### Cmd/sqlc

- Bump to v1.16.0

### Codgen

- Include serialized codegen options (#1890)

### Compiler

- Move Kotlin parameter logic into codegen (#1910)

### Examples

- Port Python examples to WASM plugin (#1903)

### Pg-gen

- Make sqlc-pg-gen the complete source of truth for pg_catalog.go (#1809)
- Implement information_schema shema (#1815)

### Python

- Port all Python tests to sqlc-gen-python (#1907)
- Upgrade to sqlc-gen-python v1.0.0 (#1932)

## [1.15.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.15.0)
Released 2022-08-07

### Bug Fixes

- (mysql) Typo (#1700)
- (postgresql) Add quotes for CamelCase columns (#1729)
- Cannot parse SQLite upsert statement (#1732)
- (sqlite) Regenerate test output for builtins (#1735)
- (wasm) Version modules by wasmtime version (#1734)
- Missing imports (#1637)
- Missing slice import for querier (#1773)

### Documentation

- Add process-based plugin docs (#1669)
- Add links to downloads.sqlc.dev (#1681)
- Update transactions how to example (#1775)

### Features

- More SQL Syntax Support for SQLite (#1687)
- (sqlite) Promote SQLite support to beta (#1699)
- Codegen plugins, powered by WASM (#1684)
- Set user-agent for plugin downloads (#1707)
- Null enums types (#1485)
- (sqlite) Support stdlib functions (#1712)
- (sqlite) Add support for returning (#1741)

### Miscellaneous Tasks

- Add tests for quoting columns (#1733)
- Remove catalog tests (#1762)

### Testing

- Add tests for fixing slice imports (#1736)
- Add test cases for returning (#1737)

### Build

- Upgrade to Go 1.19 (#1780)
- Upgrade to go-wasmtime 0.39.0 (#1781)

### Plugins

- (wasm) Change default cache location (#1709)
- (wasm) Change the SHA-256 config key (#1710)

## [1.14.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.14.0)
Released 2022-06-09

### Bug Fixes

- (postgresql) Remove extra newline with db argument (#1417)
- (sqlite) Fix DROP TABLE   (#1443)
- (compiler) Fix left join nullability with table aliases (#1491)
- Regenerate testdata for CREATE TABLE AS (#1516)
- (bundler) Only close multipart writer once (#1528)
- (endtoend) Regenerate testdata for exex_lastid
- (pgx) Copyfrom imports (#1626)
- Validate sqlc function arguments (#1633)
- Fixed typo `sql.narg` in doc (#1668)

### Features

- (golang) Add Enum.Valid and AllEnumValues (#1613)
- (sqlite) Start expanding support (#1410)
- (pgx) Add support for batch operations (#1437)
- (sqlite) Add support for delete statements (#1447)
- (codegen) Insert comments in interfaces (#1458)
- (sdk) Add the plugin SDK package (#1463)
- Upload projects (#1436)
- Add sqlc version to generated Kotlin code (#1512)
- Add sqlc version to generated Go code (#1513)
- Pass sqlc version in codegen request (#1514)
- (postgresql) Add materialized view support (#1509)
- (python) Graduate Python support to beta (#1520)
- Run sqlc with docker on windows cmd (#1557)
- Add JSON "codegen" output (#1565)
- Add sqlc.narg() for nullable named params (#1536)
- Process-based codegen plugins (#1578)

### Miscellaneous Tasks

- Fix extra newline in comments for copyfrom (#1438)
- Generate marshal/unmarshal with vtprotobuf (#1467)

### Refactor

- (codegen) Port Kotlin codegen package to use plugin types (#1416)
- (codegen) Port Go to plugin types (#1460)
- (cmd) Simplify codegen selection logic (#1466)
- (sql/catalog) Improve Readability (#1595)
- Add basic fuzzing for config / overrides (#1500)

## [1.13.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.13.0)
Released 2022-03-31

### Bug Fixes

- (compiler) Fix left join nullability with table aliases (#1491)
- (postgresql) Remove extra newline with db argument (#1417)
- (sqlite) Fix DROP TABLE (#1443)

### Features

- (cli) Upload projects (#1436)
- (codegen) Add sqlc version to generated Go code (#1513)
- (codegen) Add sqlc version to generated Kotlin code (#1512)
- (codegen) Insert comments in interfaces (#1458)
- (codegen) Pass sqlc version in codegen request (#1514)
- (pgx) Add support for batch operations (#1437)
- (postgresql) Add materialized view support (#1509)
- (python) Graduate Python support to beta (#1520)
- (sdk) Add the plugin SDK package (#1463)
- (sqlite) Add support for delete statements (#1447)
- (sqlite) Start expanding support (#1410)

### Miscellaneous Tasks

- Fix extra newline in comments for copyfrom (#1438)
- Generate marshal/unmarshal with vtprotobuf (#1467)

### Refactor

- (codegen) Port Kotlin codegen package to use plugin types (#1416)
- (codegen) Port Go to plugin types (#1460)
- (cmd) Simplify codegen selection logic (#1466)

### Config

- Add basic fuzzing for config / overrides (#1500)

## [1.12.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.12.0)
Released 2022-02-05

### Bug

- ALTER TABLE SET SCHEMA (#1409)

### Bug Fixes

- Update ANTLR v4 go.mod entry (#1336)
- Check delete statements for CTEs (#1329)
- Fix validation of GROUP BY on field aliases (#1348)
- Fix imports when non-copyfrom queries needed imports that copyfrom queries didn't (#1386)
- Remove extra comment newline (#1395)
- Enable strict function checking (#1405)

### Documentation

- Bump version to 1.11.0 (#1308)

### Features

- Inheritance (#1339)
- Generate query code using ASTs instead of templates (#1338)
- Add support for CREATE TABLE a ( LIKE b ) (#1355)
- Add support for sql.NullInt16 (#1376)

### Miscellaneous Tasks

- Add tests for :exec{result,rows} (#1344)
- Delete template-based codegen (#1345)

### Build

- Bump github.com/jackc/pgx/v4 from 4.14.0 to 4.14.1 (#1316)
- Bump golang from 1.17.3 to 1.17.4 (#1331)
- Bump golang from 1.17.4 to 1.17.5 (#1337)
- Bump github.com/spf13/cobra from 1.2.1 to 1.3.0 (#1343)
- Remove devel Docker build
- Bump golang from 1.17.5 to 1.17.6 (#1369)
- Bump github.com/google/go-cmp from 0.5.6 to 0.5.7 (#1382)
- Format all Go code (#1387)

## [1.11.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.11.0)
Released 2021-11-24


### Bug Fixes

- Update incorrect signatures (#1180)
- Correct aggregate func sig (#1182)
- Jsonb_build_object (#1211)
- Case-insensitive identifiers (#1216)
- Incorrect handling of meta (#1228)
- Detect invalid INSERT expression (#1231)
- Respect alias name for coalesce (#1232)
- Mark nullable when casting NULL (#1233)
- Support nullable fields in joins for MySQL engine (#1249)
- Fix between expression handling of table references (#1268)
- Support nullable fields in joins on same table (#1270)
- Fix missing binds in ORDER BY (#1273)
- Set RV for TargetList items on updates (#1252)
- Fix MySQL parser for query without trailing semicolon (#1282)
- Validate table alias references (#1283)
- Add support for MySQL ON DUPLICATE KEY UPDATE (#1286)
- Support references to columns in joined tables in UPDATE statements (#1289)
- Add validation for GROUP BY clause column references (#1285)
- Prevent variable redeclaration in single param conflict (#1298)
- Use common params struct field for same named params (#1296)

### Documentation

- Replace deprecated go get with go install (#1181)
- Fix package name referenced in tutorial (#1202)
- Add environment variables (#1264)
- Add go.17+ install instructions (#1280)
- Warn about golang-migrate file order (#1302)

### Features

- Instrument compiler via runtime/trace (#1258)
- Add MySQL support for BETWEEN arguments (#1265)

### Refactor

- Move from io/ioutil to io and os package (#1164)

### Styling

- Apply gofmt to sample code (#1261)

### Build

- Bump golang from 1.17.0 to 1.17.1 (#1173)
- Bump eskatos/gradle-command-action from 1 to 2 (#1220)
- Bump golang from 1.17.1 to 1.17.2 (#1227)
- Bump github.com/pganalyze/pg_query_go/v2 (#1234)
- Bump actions/checkout from 2.3.4 to 2.3.5 (#1238)
- Bump babel from 2.9.0 to 2.9.1 in /docs (#1245)
- Bump golang from 1.17.2 to 1.17.3 (#1272)
- Bump actions/checkout from 2.3.5 to 2.4.0 (#1267)
- Bump github.com/lib/pq from 1.10.3 to 1.10.4 (#1278)
- Bump github.com/jackc/pgx/v4 from 4.13.0 to 4.14.0 (#1303)

### Cmd/sqlc

- Bump version to v1.11.0

## [1.10.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.10.0)
Released 2021-09-07


### Documentation

- Fix invalid language support table (#1161)
- Add a getting started guide for MySQL (#1163)

### Build

- Bump golang from 1.16.7 to 1.17.0 (#1129)
- Bump github.com/lib/pq from 1.10.2 to 1.10.3 (#1160)

### Ci

- Upgrade Go to 1.17 (#1130)

### Cmd/sqlc

- Bump version to v1.10.0 (#1165)

### Codegen/golang

- Consolidate import logic (#1139)
- Add pgx support for range types (#1146)
- Use pgtype for hstore when using pgx (#1156)

### Codgen/golang

- Use p[gq]type for network address types (#1142)

### Endtoend

- Run `go test` in CI (#1134)

### Engine/mysql

- Add support for LIKE (#1162)

### Golang

- Output NullUUID when necessary (#1137)

## [1.9.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.9.0)
Released 2021-08-13


### Documentation

- Update documentation (a bit) for v1.9.0 (#1117)

### Build

- Bump golang from 1.16.6 to 1.16.7 (#1107)

### Cmd/sqlc

- Bump version to v1.9.0 (#1121)

### Compiler

- Add tests for COALESCE behavior (#1112)
- Handle subqueries in SELECT statements (#1113)

## [1.8.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.8.0)
Released 2021-05-03


### Documentation

- Add language support Matrix (#920)

### Features

- Add case style config option (#905)

### Python

- Eliminate runtime package and use sqlalchemy (#939)

### Build

- Bump github.com/google/go-cmp from 0.5.4 to 0.5.5 (#926)
- Bump github.com/lib/pq from 1.9.0 to 1.10.0 (#931)
- Bump golang from 1.16.0 to 1.16.1 (#935)
- Bump golang from 1.16.1 to 1.16.2 (#942)
- Bump github.com/jackc/pgx/v4 from 4.10.1 to 4.11.0 (#956)
- Bump github.com/go-sql-driver/mysql from 1.5.0 to 1.6.0 (#961)
- Bump github.com/pganalyze/pg_query_go/v2 (#965)
- Bump urllib3 from 1.26.3 to 1.26.4 in /docs (#968)
- Bump golang from 1.16.2 to 1.16.3 (#963)
- Bump github.com/lib/pq from 1.10.0 to 1.10.1 (#980)

### Cmd

- Add the --experimental flag (#929)
- Fix sqlc init (#959)

### Cmd/sqlc

- Bump version to v1.7.1-devel (#913)
- Bump version to v1.8.0

### Codegen

- Generate valid enum names for symbols (#972)

### Postgresql

- Support generated columns
- Add test for PRIMARY KEY INCLUDE
- Add tests for CREATE TABLE PARTITION OF
- CREATE TRIGGER EXECUTE FUNCTION
- Add support for renaming types (#971)

### Sql/ast

- Resolve return values from functions (#964)

### Workflows

- Only run tests once (#924)

## [1.7.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.7.0)
Released 2021-02-28


### Bug Fixes

- Struct tag formatting (#833)

### Documentation

- Include all the existing Markdown files (#877)
- Split docs into four sections (#882)
- Reorganize and consolidate documentation
- Add link to Windows download (#888)
- Shorten the README (#889)

### Features

- Adding support for pgx/v4
- Adding support for pgx/v4

### README

- Add Go Report Card badge (#891)

### Build

- Bump github.com/google/go-cmp from 0.5.3 to 0.5.4 (#813)
- Bump github.com/lib/pq from 1.8.0 to 1.9.0 (#820)
- Bump golang from 1.15.5 to 1.15.6 (#822)
- Bump github.com/jackc/pgx/v4 from 4.9.2 to 4.10.0 (#823)
- Bump github.com/jackc/pgx/v4 from 4.10.0 to 4.10.1 (#839)
- Bump golang from 1.15.6 to 1.15.7 (#855)
- Bump golang from 1.15.7 to 1.15.8 (#881)
- Bump github.com/spf13/cobra from 1.1.1 to 1.1.2 (#892)
- Bump golang from 1.15.8 to 1.16.0 (#897)
- Bump github.com/lfittl/pg_query_go from 1.0.1 to 1.0.2 (#901)
- Bump github.com/spf13/cobra from 1.1.2 to 1.1.3 (#893)

### Catalog

- Improve alter column type (#818)

### Ci

- Uprade to Go 1.15 (#887)

### Cmd

- Allow config file location to be specified (#863)

### Cmd/sqlc

- Bump to version v1.6.1-devel (#807)
- Bump version to v1.7.0 (#912)

### Codegen/golang

- Make sure to import net package (#858)

### Compiler

- Support UNION query

### Dolphin

- Generate bools for tinyint(1)
- Support joins in update statements (#883)
- Add support for union query

### Endtoend

- Add tests for INTERSECT and EXCEPT

### Go.mod

- Update to go 1.15 and run 'go mod tidy' (#808)

### Mysql

- Compile tinyint(1) to bool (#873)

### Sql/ast

- Add enum values for SetOperation

## [1.6.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.6.0)
Released 2020-11-23


### Dolphin

- Implement Rename (#651)
- Skip processing view drops (#653)

### README

- Update language / database support (#698)

### Astutils

- Fix Params rewrite call (#674)

### Build

- Bump golang from 1.14 to 1.15.3 (#765)
- Bump docker/build-push-action from v1 to v2.1.0 (#764)
- Bump github.com/google/go-cmp from 0.4.0 to 0.5.2 (#766)
- Bump github.com/spf13/cobra from 1.0.0 to 1.1.1 (#767)
- Bump github.com/jackc/pgx/v4 from 4.6.0 to 4.9.2 (#768)
- Bump github.com/lfittl/pg_query_go from 1.0.0 to 1.0.1 (#773)
- Bump github.com/google/go-cmp from 0.5.2 to 0.5.3 (#783)
- Bump golang from 1.15.3 to 1.15.5 (#782)
- Bump github.com/lib/pq from 1.4.0 to 1.8.0 (#769)

### Catalog

- Improve variadic argument support (#804)

### Cmd/sqlc

- Bump to version v1.6.0 (#806)

### Codegen

- Fix errant database/sql imports (#789)

### Compiler

- Use engine-specific reserved keywords (#677)

### Dolphi

- Add list of builtin functions (#795)

### Dolphin

- Update to the latest MySQL parser (#665)
- Add ENUM() support (#676)
- Add test for table aliasing (#684)
- Add MySQL ddl_create_table test (#685)
- Implete TRUNCATE table (#697)
- Represent tinyint as int32 (#797)
- Add support for coalesce (#802)
- Add function signatures (#796)

### Endtoend

- Add MySQL json test (#692)
- Add MySQL update set multiple test (#696)

### Examples

- Use generated enum constants in db_test (#678)
- Port ondeck to MySQL (#680)
- Add MySQL authors example (#682)

### Internal/cmd

- Print correct config file on parse failure (#749)

### Kotlin

- Remove runtime dependency (#774)

### Metadata

- Support multiple comment prefixes (#683)

### Postgresql

- Support string concat operator (#701)

### Sql/catalog

- Add support for variadic functions (#798)

## [1.5.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.5.0)
Released 2020-08-05


### Documentation

- Build sqlc using Go 1.14 (#549)

### Cmd

- Add debugging support (#573)

### Cmd/sqlc

- Bump version to v1.4.1-devel (#548)
- Bump version to v1.5.0

### Compiler

- Support calling functions with defaults (#635)
- Skip func args without a paramRef (#636)
- Return a single column from coalesce (#639)

### Config

- Add emit_empty_slices to version one (#552)

### Contrib

- Add generated code for contrib

### Dinosql

- Remove deprecated package (#554)

### Dolphin

- Add support for column aliasing (#566)
- Implement star expansion for subqueries (#619)
- Implement exapansion with reserved words (#620)
- Implement parameter refs (#621)
- Implement limit and offest (#622)
- Implement inserts (#623)
- Implement delete (#624)
- Implement simple update statements (#625)
- Implement INSERT ... SELECT (#626)
- Use test driver instead of TiDB driver (#629)
- Implement named parameters via sqlc.arg() (#632)

### Endtoend

- Add MySQL test for SELECT * JOIN (#565)
- Add MySQL test for inflection (#567)

### Engine

- Create engine package (#556)

### Equinox

- Use the new equinox-io/setup action (#586)

### Examples

- Run tests for MySQL booktest (#627)

### Golang

- Add support for the money type (#561)
- Generate correct types for int2 and int8 (#579)

### Internal

- Rm catalog, pg, postgres packages (#555)

### Mod

- Downgrade TiDB package to fix build (#603)

### Mysql

- Upgrade to the latest vitess commit (#562)
- Support to infer type of a duplicated arg (#615)
- Allow some builtin functions to be nullable (#616)

### Postgresql

- Generate all functions in pg_catalog (#550)
- Remove pg_catalog schema from tests (#638)
- Move contrib code to a package

### Sql/catalog

- Fix comparison of pg_catalog types (#637)

### Tools

- Generate functions for all of contrib

### Workflow

- Migrate to equinox-io/setup-release-tool (#614)

## [1.4.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.4.0)
Released 2020-06-17


### Dockerfile

- Add version build argument (#487)

### MySQL

- Prevent Panic when WHERE clause contains parenthesis.  (#531)

### README

- Document emit_exact_table_names (#486)

### All

- Remove the exp build tag (#507)

### Catalog

- Support functions with table parameters (#541)

### Cmd

- Bump to version 1.3.1-devel (#485)

### Cmd/sqlc

- Bump version to v1.4.0 (#547)

### Codegen

- Add the new codegen packages (#513)
- Add the :execresult query annotation (#542)

### Compiler

- Validate function calls (#505)
- Port bottom of parseQuery (#510)
- Don't mutate table name (#517)
- Enable experimental parser by default (#518)
- Apply rename rules to enum constants (#523)
- Temp fix for typecast function parameters (#530)

### Endtoend

- Standardize JSON formatting (#490)
- Add per-test configuration files (#521)
- Read expected stderr failures from disk (#527)

### Internal/dinosql

- Check parameter style before ref (#488)
- Remove unneeded column suffix (#492)
- Support named function arguments (#494)

### Internal/postgresql

- Fix NamedArgExpr rewrite (#491)

### Multierr

- Move dinosql.ParserErr to a new package (#496)

### Named

- Port parameter style validation to SQL (#504)

### Parser

- Support columns from subselect statements (#489)

### Rewrite

- Move parameter rewrite to package (#499)

### Sqlite

- Use convert functions instead of the listener (#519)

### Sqlpath

- Move ReadSQLFiles into a separate package (#495)

### Validation

- Move query validation to separate package (#498)

## [1.3.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.3.0)
Released 2020-05-12


### Makefile

- Update target (#449)

### README

- Add Myles as a sponsor (#469)

### Testing

- Make sure all Go examples build (#480)

### Cmd

- Bump version to v1.3.0 (#484)

### Cmd/sqlc

- Bump version to v1.2.1-devel (#442)

### Dinosql

- Inline addFile (#446)
- Add PostgreSQL support for TRUNCATE (#448)

### Gen

- Emit json.RawMessage for JSON columns (#461)

### Go.mod

- Use latest lib/pq (#471)

### Parser

- Use same function to load SQL files (#483)

### Postgresql

- Fix panic walking CreateTableAsStmt (#475)

## [1.2.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.2.0)
Released 2020-04-07


### Documentation

- Publish to Docker Hub (#422)

### README

- Docker installation docs (#424)

### Cmd/sqlc

- Bump version to v1.1.1-devel (#407)
- Bump version to v1.2.0 (#441)

### Gen

- Add special case for "campus" (#435)
- Properly quote reserved keywords on expansion (#436)

### Migrations

- Move migration parsing to new package (#427)

### Parser

- Generate correct types for SELECT EXISTS (#411)

## [1.1.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.1.0)
Released 2020-03-17


### README

- Add installation instructions (#350)
- Add section on running tests (#357)
- Fix typo (#371)

### Ast

- Add AST for ALTER TABLE ADD / DROP COLUMN (#376)
- Add support for CREATE TYPE as ENUM (#388)
- Add support for CREATE / DROP SCHEMA (#389)

### Astutils

- Apply changes to the ValuesList slice (#372)

### Cmd

- Return v1.0.0 (#348)
- Return next bug fix version (#349)

### Cmd/sqlc

- Bump version to v1.1.0 (#406)

### Compiler

- Wire up the experimental parsers

### Config

- Remove "emit_single_file" option (#367)

### Dolphin

- Add experimental parser for MySQL

### Gen

- Add option to emit single file for Go (#366)
- Add support for the ltree extension (#385)

### Go.mod

- Add packages for MySQL and SQLite parsers

### Internal/dinosql

- Support Postgres macaddr type in Go (#358)

### Internal/endtoend

- Remove %w (#354)

### Kotlin

- Add Query class to support timeout and cancellation (#368)

### Postgresql

- Add experimental parser for MySQL

### Sql

- Add generic SQL AST

### Sql/ast

- Port support for COMMENT ON (#391)
- Implement DROP TYPE (#397)
- Implement ALTER TABLE RENAME (#398)
- Implement ALTER TABLE RENAME column (#399)
- Implement ALTER TABLE SET SCHEMA (#400)

### Sql/catalog

- Port tests over from catalog pkg (#402)

### Sql/errors

- Add a new errors package (#390)

### Sqlite

- Add experimental parser for SQLite

## [1.0.0](https://github.com/sqlc-dev/sqlc/releases/tag/v1.0.0)
Released 2020-02-18


### Documentation

- Add documentation for query commands (#270)
- Add named parameter documentation (#332)

### README

- Add sponsors section (#333)

### Cmd

- Remove parse subcommand (#322)

### Config

- Parse V2 config format
- Add support for YAML (#336)

### Examples

- Add the jets and booktest examples (#237)
- Move sqlc.json into examples folder (#238)
- Add the authors example (#241)
- Add build tag to authors tests (#319)

### Internal

- Allow CTE to be used with UPDATE (#268)
- Remove the PackageMap from settings (#295)

### Internal/config

- Create new config package (#313)

### Internal/dinosql

- Emit Querier interface (#240)
- Strip leading "go-" or trailing "-go" from import (#262)
- Overrides can now be basic types (#271)
- Import needed types for Querier (#285)
- Handle schema-scoped enums (#310)
- Ignore golang-migrate rollbacks (#320)

### Internal/endtoend

- Move more tests to the record/replay framework
- Add update test for named params (#329)

### Internal/mysql

- Fix flaky test (#242)
- Port tests to endtoend package (#315)

### Internal/parser

- Resolve nested CTEs (#324)
- Error if last query is missing (#325)
- Support joins with aliases (#326)
- Remove print statement (#327)

### Internal/sqlc

- Add support for composite types (#311)

### Kotlin

- Support primitives
- Arrays, enums, and dates
- Generate examples
- README for examples
- Factor out db setup extension
- Fix enums, use List instead of Array
- Port Go tests for examples
- Rewrite numbered params to positional params
- Always use use, fix indents
- Unbox query params

### Parser

- Attach range vars to insert params
- Attach range vars to insert params (#342)
- Remove dead code (#343)

## [0.1.0](https://github.com/sqlc-dev/sqlc/releases/tag/v0.1.0)
Released 2020-01-07


### Documentation

- Replace remaining references to DinoSQL with sqlc (#149)

### README

- Fix download links (#66)
- Add LIMIT 1 to query that should return one (#99)

### Catalog

- Support "ALTER TABLE ... DROP CONSTRAINT ..." (#34)
- Differentiate functions with different argument types (#51)

### Ci

- Enable tests on pull requests

### Cmd

- Include filenames in error messages (#69)
- Do not output any changes on error (#72)

### Dinosql/internal

- Add lower and upper functions (#215)
- Ignore alter sequence commands (#219)

### Gen

- Add DO NOT EDIT comments to generated code (#50)
- Include all schemas when generating models (#90)
- Prefix structs with schema name (#91)
- Generate single import for uuid package (#98)
- Use same import logic for all Go files
- Pick correct struct to return for queries (#107)
- Create consistent JSON tags (#110)
- Add Close method to Queries struct (#127)
- Ignore empty override settings (#128)
- Turn SQL comments into Go comments (#136)

### Internal/catalog

- Parse unnamed function arguments (#166)

### Internal/dinosql

- Prepare() with no GoQueries still valid (#95)
- Fix multiline comment rendering (#142)
- Dereference alias nodes on walk (#158)
- Ignore sql-migrate rollbacks (#160)
- Sort imported packages (#165)
- Add support for timestamptz (#169)
- Error on missing queries (#180)
- Use more database/sql null types (#182)
- Support the pg_temp schema (#183)
- Override columns with array type (#184)
- Implement robust expansion
- Implement robust expansion (#186)
- Add COMMENT ON support (#191)
- Add DATE support
- Add DATE support (#196)
- Filter out invalid characters (#198)
- Quote reserved keywords (#205)
- Return parser errors first (#207)
- Implement advisory locks (#212)
- Error on duplicate query names (#221)
- Fix incorrect enum names (#223)
- Add support for numeric types
- Add support for numeric types (#228)

### Internal/dinosql/testdata/ondeck

- Add Makefile (#156)

### Ondeck

- Move all tests to GitHub CI (#58)

### ParseQuery

- Return either a query or an error (#178)

### Parser

- Use schema when resolving catalog refs (#82)
- Support function calls in expressions (#104)
- Correctly handle single files (#119)
- Return error if missing RETURNING (#131)
- Add support for mathmatical operators (#132)
- Add support for simple case expressions (#134)
- Error on mismatched INSERT input (#135)
- Set IsArray on joined columns (#139)

### Pg

- Store functions in the catalog (#41)
- Add location to errors (#73)

<!-- generated by git-cliff -->
