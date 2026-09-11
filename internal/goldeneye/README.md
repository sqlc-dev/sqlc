# goldeneye

`goldeneye` generates the dialect seeds under `internal/engine/<engine>/dialect`
from a live database, and verifies the committed ones against it. A dialect
is the JSONL that gives an engine its type system and standard library —
`types.jsonl`, `functions.jsonl`, `operators.jsonl`, `relations.jsonl` and
the `extensions/` bundles — read by `internal/core/seed`. Each engine package
here asks the database what it knows, writes the answer in that shape, and
the tests compare it with what is committed, byte for byte. A difference
means the committed dialect has drifted from the database.

It is a nested Go module, so its only dependencies beyond the standard
library are the database drivers — PostgreSQL's, MySQL's, SQL Server's and
the Spanner client — and it never shares code with the analysis that reads
the files: the files are the contract. Run it from this directory:

```bash
go run ./cmd/goldeneye install clickhouse   # download the pinned clickhouse binary once
go run ./cmd/goldeneye install duckdb       # download the current DuckDB 2.0 preview build once
go run ./cmd/goldeneye install sqlite       # build the pinned sqlite3 shells once; needs a C compiler
go run ./cmd/goldeneye check                # check every engine whose database is available
go run ./cmd/goldeneye check postgresql     # check one engine
go run ./cmd/goldeneye check spanner        # SPANNER_SERVER_URI=localhost:15000, a Spanner Omni container
go run ./cmd/goldeneye generate [engine]    # rewrite the generated files from the database
go test ./...                               # the same checks as tests; engines without a database skip
```

`generate` and `check` say which engines they skipped for lack of a
database; naming an engine makes its database required.

## What is generated, and what is not

A generator owns only the files it produces; `dialect.json` is always written
by hand, and so are the lists an engine cannot describe. Both commands leave
the hand-written files alone, and the checks do not look at them.

- **`postgresql`** reads a live server named by `POSTGRESQL_SERVER_URI`:
  `functions.jsonl` is `pg_catalog`'s functions, `relations.jsonl` the tables
  and views of `pg_catalog` and `information_schema`, and each contrib
  extension gets a directory under `extensions/` holding the types and
  functions `CREATE EXTENSION` adds, so the server needs contrib installed.
  A function that one of those extensions puts in `pg_catalog`, as
  `adminpack` does, belongs to the extension's directory rather than the
  catalog's list. The server has to be the major release pinned in
  `postgresql.Major`, since every release adds to the catalogs; the top-level
  `types.jsonl` and `operators.jsonl` are hand-written.
- **`mysql`** reads a live server named by `MYSQL_SERVER_URI`, a go-sql-driver
  DSN such as `root:mysecretpassword@tcp(127.0.0.1:3306)/mysql`, and writes
  into `internal/engine/dolphin/dialect`, since sqlc's MySQL engine is named
  after its parser. MySQL keeps no catalog of its types, functions or
  operators — the help tables describe a built-in function no further than
  its name — so `types.jsonl`, `functions.jsonl` and `operators.jsonl` are
  hand-written. What it does describe is its data dictionary:
  `relations.jsonl` is every view of `information_schema`, read from
  `information_schema` itself, with their names in lower case, since MySQL
  matches them in any case and sqlc's parser lowercases every identifier. A
  column's type is spelled the way a declaration does, `bigint unsigned`
  included, which `types.jsonl` lists as an alias of `bigint`. The other
  schemas `mysqld --initialize` creates — `mysql`, `performance_schema`, `sys`
  — are tables rather than views, and every table a dialect seeds is one the
  analysis core hands codegen as a model, so they are left out until codegen
  knows a system schema when it sees one. The server has to be the major
  release pinned in `mysql.Major`, since every release adds to
  `information_schema`.
- **`duckdb`** reads the DuckDB CLI named by `DUCKDB`, or the one `install`
  put in the user cache directory, or `duckdb` on `PATH`: `types.jsonl`,
  `functions.jsonl` and `operators.jsonl` come from `duckdb_types()` and
  `duckdb_functions()`. The CLI has to be a DuckDB 2.0 build, the release
  darkwing is pinned against, which has no release to download yet: until
  2.0 is out, `install` downloads the current build of DuckDB's v2.0
  preview channel, `duckdb.DefaultVersion`, a rolling tarball per platform
  under `artifacts.duckdb.org` with no per-build download and no checksum
  to pin, so what a run logs is the version the CLI reports, and
  `duckdb.GeneratedFrom` records the build the committed dialect came
  from. A check against a later build reports what the later build added;
  regenerate, and update `GeneratedFrom`, to move the dialect along. Once
  2.0 is released, the installer should pin the release and its checksums
  the way the clickhouse one does.
- **`clickhouse`** needs no server: `types.jsonl` comes from
  `system.data_type_families` of an ephemeral `clickhouse local` process,
  every family that is not an alias becoming a type carrying the spellings
  that alias it, with a category decided by its name. The binary is
  downloaded once per pinned release by `install` into the user cache
  directory, or supplied through the `CLICKHOUSE` environment variable; the
  pinned release and the SHA-512 of each platform's download live in
  `clickhouse/install.go`, and a download that does not match is discarded.
  ClickHouse describes its functions no further than their names, so
  `functions.jsonl` is hand-written.
- **`mssql`** reads a live server named by `MSSQL_SERVER_URI`, in any form
  the go-mssqldb driver accepts, such as
  `sqlserver://sa:password@127.0.0.1:1433?encrypt=disable`. SQL Server keeps
  no catalog of its intrinsic functions or its operators — GETDATE and LEN
  are not objects — so `types.jsonl` and `functions.jsonl` are hand-written.
  What it does describe is its catalog: `relations.jsonl` is every view of
  the `sys` and `INFORMATION_SCHEMA` schemas, listed from a database of its
  own, since the views a query sees are the ones a user database has and
  `master` lists internal views no query can name; each view's columns are
  what `sys.dm_exec_describe_first_result_set` says a `SELECT *` from it
  returns, the type spelled the way a declaration spells it —
  `nvarchar(128)`, `decimal(10,2)`, `varbinary(max)` — and the nullability
  the server computes. Names are written in lower case, since SQL Server
  matches them in any case under its default collations and sqlc's parser
  lowercases every identifier. Both schemas hold nothing but views, so no
  table is seeded that codegen would take for a model. The server has to
  be the major release pinned in `mssql.Major`, since every release adds
  to the catalog views.
- **`spanner`** reads a live Spanner Omni server — the downloadable
  Spanner, run from its container image — named by `SPANNER_SERVER_URI`,
  the gRPC endpoint such as `localhost:15000`, reached without TLS or
  credentials as Omni is, and writes into `internal/engine/googlesql/dialect`,
  since sqlc's engine is named after the language Spanner speaks. Spanner
  keeps no catalog of its types, functions or operators, so `types.jsonl`,
  `functions.jsonl` and `operators.jsonl` are hand-written. What it does
  describe is its information schema: `relations.jsonl` is every view of
  `INFORMATION_SCHEMA` and `SPANNER_SYS`, read from `INFORMATION_SCHEMA`
  itself in a database created for the purpose in the instance Omni's
  single server provides, `projects/default/instances/default`. Names are
  kept as the catalog spells them, in upper case, which is how a query
  names them; a column's type is spelled in lower case the way the seed
  spells one, an `ARRAY<T>` as `T` with the array flag, a `STRUCT<a T>` as
  `struct(a: t)` and a `PROTO<p.M>` as `proto('p.M')`, since a seed writes
  a type's arguments in parentheses. The container image is pinned in the
  gen workflow and `docker-compose.yml`.
- **`sqlite`** needs no server either: `functions.jsonl` comes from
  `pragma_function_list` of a `sqlite3` shell run against an in-memory
  database. Which functions a SQLite has is decided when it is compiled, so
  `install` downloads the pinned release's amalgamation, checked against the
  SHA3-256 the download page lists, and compiles the shell from it with the
  compiler `CC` names, or `cc` — once with the options sqlite.org's own
  configure turns on by default, which gives `functions.jsonl`, and once
  more per option in `sqlite/install.go`'s extension list, each of which
  gets a directory under `extensions/` holding the functions its build adds
  over the default one, the way each PostgreSQL contrib extension holds what
  `CREATE EXTENSION` adds; a schema that says `CREATE VIRTUAL TABLE ... USING
  fts5` loads the option's directory, through the `modules` map in the
  hand-written `dialect.json`. SQLite describes its functions as far as their
  names, their kinds and the number of arguments each overload takes, and no
  further — it types values, not functions — so what each returns and what
  its arguments hold is read from the amalgamation: every function is
  registered with the C functions that implement it, and those set their
  result through `sqlite3_result_*` and read their arguments through
  `sqlite3_value_*`. A function a shell reports that the source does not
  register fails the run rather than being guessed at. Whether an aggregate
  returns NULL over no rows is found by running it over none; the scalar
  functions that return NULL for arguments that are not are a short list in
  `sqlite/signatures.go`, since a SQLite function returns NULL as often by
  setting no result as by saying so. The pinned release is the one the main module's
  driver embeds. SQLite has no catalog of types or operators, so
  `types.jsonl` and `operators.jsonl` are hand-written. A function that
  returns an integer for an integer and a real for a real — `abs`, `ceil`,
  `floor`, `trunc`, `sum` — is written once over `any`, returning the real
  that a text or blob argument gets, and once more per spelling
  `types.jsonl` gives integer and real, returning that type, the way
  PostgreSQL's catalog has a `sum` per numeric type; the overload over `any`
  comes first, since the legacy compiler resolves by arity alone and takes
  it. `install` builds one more shell, for the analysis check below;
  nothing is generated from it.

## Layout

- `dialect/` — the record types the files are made of, mirrored from
  `internal/core/seed`, and the helpers that write a generated set of files
  into an engine directory or diff it against what is committed.
- `endtoend/` — finds the analyze cases, splits their query files, and
  compares an engine's answer with a case's committed output.
- `analysis/` — the shape of that answer: the JSON `sqlc analyze` prints.
- `postgresql/`, `mysql/`, `mssql/`, `spanner/`, `duckdb/`, `clickhouse/`,
  `sqlite/` — one package per engine, each exposing `Locate`, `Version` and
  `Generate`, `Analyze` where the engine has an analysis check, and tests
  that run the checks.
- `cmd/goldeneye/` — the command.

## Analysis checks

`check` also verifies the `analyze_*` cases under `internal/endtoend/testdata`
against what the database itself reports. A case is an
`analyze_<name>/<engine>` directory whose `exec.json` runs the analyze
command; `endtoend/` finds them. The engine package loads the case's
`schema.sql` and optional `fixture.sql` into the database, runs `query.sql`
there, prints what the database reports in the JSON shape `sqlc analyze`
prints, and compares it with the committed `stdout.json` byte for byte. A
difference means sqlc's analysis disagrees with the database. A case that
asks for `--ast` is skipped, since only sqlc can print that.

- **`clickhouse`** runs each case in an ephemeral `clickhouse local` process.
  Column types come from the executed query's result header, provenance from
  `EXPLAIN QUERY TREE`, and parameters from sentinel constants substituted for
  `?`, `sqlc.arg()` and `sqlc.narg()`, since ClickHouse itself never sees a
  placeholder; `INSERT ... VALUES` parameters map onto `DESCRIBE TABLE`.
- **`sqlite`** runs each case through one more shell `install` builds, under
  `analysis/`: every extension option at once, so that any case's schema
  loads, and `SQLITE_ENABLE_COLUMN_METADATA`, which lets the shell's `.stats
  stmt` say which table column each result column of a statement is read
  from. That column's declared type is the result column's, and its NOT
  NULL decides nullability, the rowid counting as NOT NULL. SQLite types
  values rather than expressions, so a column the library has no origin for
  — an aggregate, an arithmetic result — is typed by the storage class of
  the value it returns, which is why a case wants a `fixture.sql`: the query
  is run over the fixture, with each parameter bound to a value of the
  column it stands in for, and again over no rows, and a column is nullable
  when either run returns a NULL for it — an aggregate over nothing, the far
  side of an outer join. The library reports nothing about a parameter but
  its number, so parameters are found in the bytecode `EXPLAIN` prints, the
  way ClickHouse's are found in its query tree: each is followed from the
  register its `Variable` loads, through copies and the expressions it is an
  argument of, to the first opcode that uses it against something the
  catalog can name — the other operand of a comparison, the row a seek lands
  on, the position in the record an `Insert` writes, the column of an IN
  list's ephemeral table it comes back out of. One that reaches nothing
  nameable is described by what the program requires of it, when it
  requires anything: `MustBeInt` makes LIMIT's an integer. `sqlc.arg(name)`
  becomes `?N`, numbered as sqlc numbers them, so a repeated name is one
  parameter. Two things the check reports that sqlc does not: a bare column
  selected alongside an aggregate is NULL over no rows, and so nullable, and
  a comparison such as `x IS NULL` is an integer, since that is what SQLite
  returns.

- **`mysql`** runs each case in a database of its own on the server named by
  `MYSQL_SERVER_URI`, and asks the server three things about a query. What a
  driver sees: the query is run, with every parameter a user variable set to
  NULL, and each result column's name, type and nullability are read from
  the result set's metadata as go-sql-driver reports them. What the resolver
  made of it: the optimizer trace prints each query block back after name
  resolution and before optimisation, with every column qualified, every
  alias kept and every `SELECT *` expanded, which says which table a result
  column is read from and what a parameter is compared with or assigned to.
  And for a statement the trace does not expand — an `INSERT ... VALUES`, a
  single-table `UPDATE` or `DELETE` — the note `EXPLAIN` leaves, which prints
  the statement the same way; a `SELECT` cannot be read from the note, since
  it is printed after optimisation, and a lookup on a unique key against an
  empty table has folded to `NULL = (@x)` there. Views and derived tables are
  kept as the query wrote them rather than merged, so that a column read
  through one is reported as its column and an `information_schema` view is
  not resolved away into the dictionary tables behind it. MySQL itself
  reports nothing about a parameter but its position, so a parameter is
  described by its partner: a column's type and nullability come from
  `information_schema`, a column of a derived table or CTE from what its
  block projects, and an expression's from running it, over the tables it
  reads, as a query of its own; a `LIMIT` or `OFFSET` count is a `bigint
  unsigned`. Two things the driver keeps to itself: the table a result
  column comes from, which is why the trace is read for it, and the length
  beside the one wire type every size of `TEXT` and `BLOB` is sent as, which
  is why a column read from a table, directly or through a derived table,
  is spelled the way the table declares it.

- **`duckdb`** runs each case through the CLI, which loads the schema and
  fixture into an in-memory database of their own, one process per
  question, and is asked four things about each query. What its
  parameters are: the query is prepared and explained with a string
  sentinel bound to each parameter, `EXECUTE q('goldeneye_1', ...)`, and
  the unoptimized logical plan the CLI prints under
  `explain_output = 'all'` shows each as `CAST('goldeneye_k' AS T)`, `T`
  being the type the binder gave the parameter; a sentinel the binder
  converts on the spot, as an `INSERT`'s `VALUES` are, is bound to NULL
  instead. What its result columns are: `DESCRIBE`, with each parameter
  replaced by a NULL of its type, names and types them; DuckDB describes
  no DML, so a `RETURNING` column is the target table's column it names.
  Which table a result column is read from and which column a parameter
  stands in for: DuckDB prints a plan with every column by its bare name
  and every aliased expression by its alias, so these are read from the
  query text, the select list's items, a star expanded to its table's
  columns, and the operand beside each parameter, resolved against the
  `FROM` clause and the catalog, `duckdb_columns()`, from which a column
  read from a table takes its declared type and nullability; a parameter
  the query casts takes the cast's type as DuckDB spells it. And whether
  an expression can be NULL, which DuckDB does not track: the query is
  run, with each parameter bound to a value of its type, over the fixture
  and over no rows, and a column is nullable when either run returns a
  NULL for it. DuckDB spells an enum column by its labels whether the
  schema named the type or not, so labels that are those of an enum the
  schema created name that type, and a spelling `types.jsonl` lists as an
  alias — `json`, which DuckDB's own catalog lists as a spelling of
  `varchar` — is reported by the dialect's name for it.
- **`mssql`** describes each case in a database of its own on the server
  named by `MSSQL_SERVER_URI`, without running anything: the schema is
  loaded one statement at a time, since a `CREATE TYPE` has to be its own
  batch before a table can use the type, and the server is asked three
  things about each query. What a driver would see:
  `sys.dm_exec_describe_first_result_set` describes each result column —
  its name, its type spelled the way a declaration spells it, whether it
  can be NULL, and which table column it is read from. What each parameter
  would be: `sp_describe_undeclared_parameters` says what type the server
  would give each parameter the query leaves undeclared, which is the type
  of a `CAST(@p AS T)`; it describes a parameter only when it is used
  once, so each appearance of a repeated one becomes a variable of its
  own. And what each parameter stands in for: the estimated showplan,
  compiled with `SET SHOWPLAN_XML ON` and the parameters declared as those
  types, prints every column as a reference naming its table and every
  variable as `@variable`, and a parameter's partner is the column on the
  other side of the `Compare` it is an operand of, the column an `Assign`
  sets to it, or the column of a seek whose range expression it is, with a
  named expression such as `Expr1002` followed to its definition; a
  parameter under a `CONVERT` the query wrote has no partner, since the
  cast says what it is. A parameter with a partner is described as that
  column, from the catalog of the case's database. Three things the
  describing function keeps to itself: a `json` or `vector` column is
  described by the `nvarchar(max)` it is sent to a driver as, so a column
  read from a table is typed from `sys.columns` instead; a type the
  schema created is reported beside the system type it stands on, and the
  dialect reports it by its own name; and a spelling `types.jsonl` lists
  as an alias — `numeric`, `timestamp` — is reported by the dialect's name
  for it, `decimal`, `rowversion`. One thing it says that sqlc does not: a
  computed column — a cast, an arithmetic — is nullable whatever its
  arguments, since a conversion that fails under `ANSI_WARNINGS OFF`
  yields NULL, where sqlc follows the arguments, and a computed column
  names the column it is computed from, as the one an update through the
  result set would write, where sqlc gives an expression no table. The
  cases compute over nullable columns, and a computed column is reported
  without a table.
- **`spanner`** creates a database of each case's own from its schema in
  the instance named by `SPANNER_SERVER_URI`, writes its fixture there in a
  read-write transaction, and compiles each query in `PLAN` mode, which
  runs nothing: a DML statement is compiled in a read-write transaction
  that is rolled back. The server reports the name and type of each
  result column and the type of each parameter the query leaves
  undeclared, the way the wire spells them — `STRING`, `ARRAY<INT64>` —
  without the length a declaration gives a `STRING(10)` or whether a
  column can be NULL, and the query plan, which says where each comes
  from: the children of `Serialize Result` after the relation it
  serializes are the result columns, a `Scan` of a table defines a
  variable per column it reads, and a comparison is a `Function` whose
  description reads `($col = @param)`. A result column the plan reads
  from a table column, directly or through the variables a join or a
  batch passes it through, is described from the case's
  `INFORMATION_SCHEMA`, with its declared length and nullability; one the
  plan reads as a parameter is the column the parameter is compared with,
  which is why the optimizer substituted it. A parameter compared with a
  table column is described as that column. A DML plan lists the values
  it writes before the columns it returns — the table's key columns, then
  the columns an `UPDATE` sets or an `INSERT` inserts, read from the
  statement — and a parameter written to a column is described as it,
  while a `THEN RETURN` column is the table's column of that name. Two
  things Spanner does not say: whether an expression can be NULL, which
  is reported as it is not, and anything about a `STRUCT` returned as a
  column, which Spanner rejects.

The other engines have no analysis check yet.
