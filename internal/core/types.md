# Types in the analysis core

How the core catalog, the analyzer and `sqlc analyze` represent a type today,
where that falls short of each engine's type system, how each engine's own
catalog represents one, and the design that closes the gap: every type the
schema or the dialect declares is a row holding its full expression,
canonicalized the way the engine would and denormalized so that a bare name
like `integer` and a structured expression like `numeric(10, 2)` both
resolve to one.

## What a type is today

The catalog (`catalogdef/schema.sql`) has one row per *name* in `sql_type`:
a lowercased string, a category letter, a `typtype` that is only ever `b` or
`e`, and an `element_oid` that is written for arrays but never read. A
dialect's `types.jsonl` seeds one row per type and one more per alias, then
joins every spelling of a type to every other with implicit casts. Arrays are
rows named after their element with `[]` appended, created on first use.
Anything else — a type argument, a struct field, an enum label, a domain's
base, a range's subtype — has nowhere to go: `sql_attribute` keeps the
column's verbatim spelling in `decl_type` (only SQLite and ClickHouse set it)
and has `type_length` and `type_scale` columns that no DDL path fills.

The analyzer (`analyzer/expr.go`) types an expression as `exprType`: a type
OID, or a bare name when the catalog has no row, plus nullability. A type is
therefore a row, and a type that is not a row degrades to whatever row its
spelling's first word finds.

The output (`TypeExpr` in `typeexpr.go`, printed by `sqlc analyze`) already
has the right shape: a name applied to labelled arguments that are types,
integers, booleans or strings, with `nullable` at any depth. It is complete
only where a source column's `decl_type` spelling exists to be parsed; every
other column and every parameter is rebuilt from the flat row name.

## What each engine loses

The table is what `sqlc analyze` reports today for a schema exercising each
engine's type system, against what the engine itself says the type is.

| Engine | Declared | Reported | The engine says |
|---|---|---|---|
| PostgreSQL | `numeric(10,2)`, `varchar(255)`, `timestamp(3)`, `bit(8)` | `numeric`, `varchar`, `timestamp`, `bit` | typmods are part of the type |
| PostgreSQL | `int[][]` | `array(int4)` | two dimensions |
| PostgreSQL | `CREATE DOMAIN posint AS integer` | `posint`, category U | an integer with a constraint |
| PostgreSQL | `CREATE TYPE point2 AS (x float8, y float8)` | `point2`, category U | two named fields |
| PostgreSQL | `CREATE TYPE mood AS ENUM (...)` | `mood` | the labels |
| PostgreSQL | `CREATE TYPE floatrange AS RANGE (subtype = float8)` | `floatrange`, category U | a range over float8 |
| PostgreSQL | `myschema.mood` | `mood` for a column, `myschema.mood` for a cast: two rows, both in `public` | one type in a namespace |
| PostgreSQL | `interval day to second` | `interval` | fields are a typmod |
| MySQL | `BIGINT UNSIGNED`, `INT UNSIGNED` | `bigint`, `int` | a different value range; codegen picks `int64` over `uint64` |
| MySQL | `TINYINT(1)` | `tinyint` | the display width is how drivers and codegen spot a boolean |
| MySQL | `DECIMAL(10,2) UNSIGNED`, `DATETIME(6)`, `VARCHAR(255)` | `decimal`, `datetime`, `varchar` | precision, fractional seconds, length |
| MySQL | `ENUM('a','b')`, `SET('x','y')` | `enum`, `set` | the members |
| MySQL | `CAST(? AS CHAR(10))` | `var_string` | `char`; the parser's internal name leaks |
| SQLite | `FOO BAR(3)`, `VARCHAR(255)` | `foo bar(3)`, `varchar(255)` | correct spelling, but each is a row of category U that compares with nothing, and the affinity SQLite gives it (NUMERIC, TEXT) is not modelled |
| ClickHouse | every column type | complete | complete, from the spelling |
| ClickHouse | `CAST(x AS Nullable(String))` | `nullable` | `Nullable(String)` |
| ClickHouse | `CAST(x AS Array(UInt8))` | `array` | `Array(UInt8)` |
| ClickHouse | `toDecimal64(x, 4)`, `toDateTime64(x, 3)` | `decimal64`, `datetime64` | `Decimal(18, 4)`, `DateTime64(3)`: the result depends on an argument's value |
| ClickHouse | `SimpleAggregateFunction(sum, UInt64)` | `simpleaggregatefunction(sum, uint64)` | `sum` is a function name, which the expression reads as a type |
| ClickHouse | `n Nested(a UInt8, b String)` | one column `n` | two columns `n.a Array(UInt8)`, `n.b Array(String)` |
| DuckDB | `STRUCT(a INTEGER, b VARCHAR)`, `MAP(VARCHAR, INTEGER)`, `UNION(num INTEGER, str VARCHAR)` | `struct`, `map`, `union` | the fields |
| DuckDB | `INTEGER[]`, `INTEGER[3]`, `INTEGER[][]` | `array(integer)` for all three | LIST, fixed-size ARRAY, nested LIST |
| DuckDB | `DECIMAL(18,3)`, `VARCHAR(10)`, `ENUM('a','b')` | `decimal`, `varchar`, `enum` | arguments and members |
| GoogleSQL | `ARRAY<INT64>` | a row *named* `array<int64>` | an array of int64 |
| GoogleSQL | `STRUCT<a INT64, b STRING>`, `ARRAY<STRUCT<x INT64>>` | `struct`, `array<struct>` | the fields |
| GoogleSQL | `STRING(10)`, `NUMERIC(10,2)`, `[1, 2]`, `STRUCT(1 AS x)` | `string`, `numeric`, untyped, untyped | parameters; a constructed array and struct |
| SQL Server | `NVARCHAR(MAX)`, `VARBINARY(MAX)` | `nvarchar`, `varbinary` | MAX decides the Go type |
| SQL Server | `DECIMAL(10,2)`, `DATETIME2(3)`, `FLOAT(24)`, `VECTOR(3)` | `decimal`, `datetime2`, `float`, `vector` | arguments; `FLOAT(24)` is `real` |
| SQL Server | `CREATE TYPE dbo.PhoneNumber FROM varchar(20) NOT NULL` | `phonenumber`, category U | a `varchar(20)` that is never null, in schema `dbo` |

Three of these are regressions against the legacy compiler rather than gaps
shared with it: the plugin protocol's `Column` carries `unsigned`, `length`
and `array_dims`, codegen reads all three (`golang/mysql_type.go` turns
`tinyint` with length 1 into `bool` and `unsigned` into `uint64`;
`golang/go_type.go` nests one slice per dimension), and the core bridge in
`compiler/parse_core.go` sets none of them from what the core reports.

Two are broken outside column declarations only: ClickHouse's columns are
whole because the engine hands the core a spelling and the core keeps it on
the attribute. The same spelling in a cast, a function result or a typed
placeholder has no attribute to live on, so it degrades to its first word.
That is the tell: the expression belongs to the type, not to the column.

## What each engine's own catalog says

Each engine was asked, on a live system where one could be had, what it
records about a type, how it spells a column's type back, what it says a
query result's type is, and how it describes a function. PostgreSQL 16,
MySQL 8, SQLite 3.53, ClickHouse 25.8 and DuckDB 1.5 answered directly; SQL
Server and GoogleSQL are from their documented catalogs and the parsers
sqlc uses for them.

**PostgreSQL** has one `pg_type` row per family, and a second row per array
type (`_int4`, with `typelem` pointing at the element). Arguments are not
part of the type: they are an opaque `int32` typmod on the use site —
`pg_attribute.atttypmod`, a domain's `typtypmod`, a function argument's
type has none — whose encoding is the type's own business and which
`format_type(oid, typmod)` decodes back into a spelling. Dimensions are
likewise on the attribute (`attndims`) and informational: `int[3]` and
`int[][]` are both the type `integer[]`. Declared types are rows with a
`typtype`: `e` with labels and their order in `pg_enum`, `c` with fields as
the attributes of a hidden relation (`typrelid`), `d` with `typbasetype`,
`typtypmod` and `typnotnull`, `r` and `m` with the subtype in `pg_range`;
pseudo-types (`anyarray`, `record`) are rows of `typtype` `p` and category
`P`. What it reports is canonical and not what was written: `int` comes back
as `integer`, `varchar(255)` as `character varying(255)`, `timestamp(3)` as
`timestamp(3) without time zone`, `int[3]` as `integer[]`, and a domain
column as the domain's name. A view keeps its columns' typmods; a prepared
statement's result types (`pg_prepared_statements.result_types`, the
protocol's RowDescription) carry the typmod only for a column read straight
from a table and `-1` for an expression, so `numeric(10,2) + 1` is
`numeric`. `pg_proc` has full signatures over families and pseudo-types,
`pg_operator` likewise, and `pg_cast` the contexts.

**MySQL** has no catalog of types, functions or operators. What it has is
`information_schema.COLUMNS`, which describes a column twice: `DATA_TYPE` is
the family (`bigint`, `decimal`, `enum`) and `COLUMN_TYPE` is the whole
canonical spelling — `bigint unsigned`, `tinyint(1)`, `decimal(10,2)
unsigned`, `enum('a','b')`, `bigint(20) unsigned zerofill` — beside the
arguments decoded into `NUMERIC_PRECISION`, `NUMERIC_SCALE`,
`CHARACTER_MAXIMUM_LENGTH`, `DATETIME_PRECISION`, `CHARACTER_SET_NAME` and
`COLLATION_NAME`. The spelling is canonical: `BOOLEAN` becomes `tinyint(1)`,
`INTEGER` becomes `int`, `VARCHAR(10) CHARACTER SET binary` becomes
`varbinary(10)`. A view's columns show that the binder computes a full type
for every expression: `decimal(10,2) + 1` is `decimal(11,2)`, `decimal(10,2)
* decimal(10,2)` is `decimal(20,4)`, `int unsigned + 1` is `bigint
unsigned`, `CONCAT(varchar(255), 'x')` is `varchar(256)`, `SUM(int unsigned)`
is `decimal(32,0)`, `AVG(int)` is `decimal(14,4)`, `->>` is `longtext`. User
routines are described the same way in `ROUTINES` and `PARAMETERS`
(`DTD_IDENTIFIER` is `decimal(5,2)`, `bigint unsigned`); built-in functions
are not described at all. A result set on the wire is coarser than the
catalog: a storage type code (`TINY`, `LONG`, `NEWDECIMAL`, `VAR_STRING`),
a length in bytes, a decimals count and flags (`UNSIGNED`, `NOT_NULL`,
`ENUM`, `SET`, `BINARY`), so a `varchar(255)` column and `CAST(x AS
CHAR(10))` are both `VAR_STRING`, an enum is `STRING` with the `ENUM` flag,
and a `tinyint(1)` is `TINY` with length 1. That code is where the parser's
`var_string` comes from.

**SQLite** has no catalog of types either. `pragma_table_xinfo` returns a
column's declared type verbatim — `VARCHAR(255)`, `FOO BAR(3)`, `my type`,
or nothing — and the type system is affinity: a rule over the spelling's
substrings (INT anywhere is INTEGER; CHAR, CLOB or TEXT is TEXT; BLOB or no
type is BLOB; REAL, FLOA or DOUB is REAL; anything else is NUMERIC) that
decides how a stored value is coerced, so `FOO BAR(3)` holding `'12'`
stores the integer 12. Arguments are decoration: `CAST(x AS DECIMAL(5,2))`
applies NUMERIC affinity and nothing else. A `STRICT` table admits only
`INT`, `INTEGER`, `REAL`, `TEXT`, `BLOB` and `ANY`, and rejects
`VARCHAR(10)`. An expression has no type; a value has a storage class
(`typeof()` is `integer`, `real`, `text`, `blob` or `null`), and a result
column's `sqlite3_column_decltype` is the declared spelling for a column
read from a table and nothing for anything else. `pragma_function_list`
names functions and their arity and nothing more.

**ClickHouse** models a type as a call expression, which is why the output
shape was chosen. `system.data_type_families` lists 66 families and 73
aliases (`INT` is an alias of `Int32`); a column's type in `system.columns`
is the whole expression, canonicalized: `Decimal32(4)` is stored as
`Decimal(9, 4)`, `Enum('a', 'b')` as `Enum8('a' = 1, 'b' = 2)`,
`Variant(String, Int64)` with its members sorted, and `Nested(a UInt8, b
String)` as two columns `n.a Array(UInt8)` and `n.b Array(String)`.
`Nullable` and `LowCardinality` are spelled as wrappers and reported as part
of the type, at any depth. Every expression has a full type computed by the
binder, from argument types with promotion — `Int8 + UInt8` is `Int16`,
`Int32 + UInt64` is `Int64`, `Int32 / Int32` is `Float64`, `Decimal(9, 4) *
Decimal(38, 10)` is `Decimal(38, 14)` — from argument values —
`toDecimal64(x, 4)` is `Decimal(18, 4)`, `toDateTime64(x, 3)` is
`DateTime64(3)` — and from literals by value: `1` is `UInt8`, `-1` is
`Int8`, `[1, NULL]` is `Array(Nullable(UInt8))`. Wrappers propagate:
`concat(lc, 'x')` is `LowCardinality(String)`, `ns = 'x'` is
`Nullable(UInt8)`, `NULL` is `Nullable(Nothing)`. A typed placeholder
`{p:Decimal(10, 2)}` is exactly its declared type. `system.functions` has a
name, an aggregate flag, an alias and a description, and no signature.

**DuckDB** has `duckdb_types()`, one row per spelling with a `logical_type`
naming the family (`int`, `int4` and `integer` are three rows over
`INTEGER`), a category (`NUMERIC`, `STRING`, `DATETIME`, `BOOLEAN`,
`COMPOSITE`, or none for `bit` and `enum`), and `labels` for an enum, which
a `CREATE TYPE ... AS ENUM` adds to as a non-internal row in the user's
schema. `duckdb_columns().data_type` is the whole canonical expression:
`INTEGER[]`, `INTEGER[3]`, `STRUCT(a INTEGER, b VARCHAR)`, `MAP(VARCHAR,
INTEGER)`, `UNION(num INTEGER, str VARCHAR)`, `DECIMAL(18,3)`, with the
precision and scale also decoded into their own columns. Canonicalization
goes further than anywhere else: `TEXT` is `VARCHAR`, `VARCHAR(10)` is
`VARCHAR` (the length is dropped), `NUMERIC` is `DECIMAL(18,3)`, `VARINT` is
`BIGNUM`, `FLOAT4` is `FLOAT`, `JSON` keeps its name over `VARCHAR`'s id,
and a column of the named enum `mood` is reported as `ENUM('sad', 'ok')`
with the name gone. Expressions have full binder-computed types:
`DECIMAL(18,3) * DECIMAL(18,3)` is `DECIMAL(18,6)`, `1.5` is
`DECIMAL(2,1)`, `SUM(INTEGER)` is `HUGEINT`, `SUM(DECIMAL(18,3))` is
`DECIMAL(38,3)`, `INTEGER / 2` is `DOUBLE`, `[1, NULL]` is `INTEGER[]` with
no inner nullability, and `NULL` is a type of its own. `duckdb_functions()`
has signatures over families and generics — `parameter_types` are
`DECIMAL`, `INTEGER[]`, `ANY`, `T`, `K`, `V` and `return_type` a family —
so the arguments of a result are the binder's, not the catalog's. A
prepared statement's parameter types are `UNKNOWN` until bound.

**SQL Server** describes a column in `sys.columns` as a family plus decoded
arguments: `system_type_id` and `user_type_id`, `max_length` in bytes with
`-1` for `MAX`, `precision`, `scale`, `collation_name`, `is_nullable`,
`is_identity`. `sys.types` has one row per system type — `decimal` and
`numeric` are two — and one per alias type from `CREATE TYPE ... FROM`, with
`is_user_defined`, its own `max_length`, `precision`, `scale` and
`is_nullable`, and its base's `system_type_id`; `sysname` is such a row over
`nvarchar(128)`, and CLR types like `geography` and `hierarchyid` are
assembly types sharing one `system_type_id`. Canonicalization fills in what
was left out and folds one family into another: `varchar` alone is
`varchar(1)` in a declaration and `varchar(30)` in a `CAST`, `decimal` is
`decimal(18,0)`, `datetime2` is `datetime2(7)`, `FLOAT(24)` is stored as
`real` and `FLOAT(25)` and up as `float`. `sys.all_parameters` describes a
function's parameters the same way as a column. The T-SQL parser sqlc uses
(`teesql`) reports a type as a name with a parameter list in which `MAX` is
a literal node of its own.

**GoogleSQL** has no catalog sqlc can read offline; its two hosts describe a
column as a whole string. BigQuery's `INFORMATION_SCHEMA.COLUMNS.DATA_TYPE`
is `ARRAY<STRUCT<a INT64, b STRING>>`, `NUMERIC(10, 2)` or `STRING(10)`,
and `COLUMN_FIELD_PATHS` lists every nested struct field with a path and a
`DATA_TYPE` of its own. Spanner's `SPANNER_TYPE` is `STRING(MAX)`,
`ARRAY<INT64>`, `PROTO<pkg.Message>` or `ENUM<pkg.Enum>`, and a column
cannot be a struct. The parser sqlc uses (`zetajones`) has a node per type
form — `SimpleType` with a `TypeParameterList`, `ArrayType`, `StructType`
with named fields, `RangeType`, `MapType`, `FunctionType` — and, as in T-SQL,
`MAX` is a literal node in the parameter list. Arrays do not nest.

What the review settles:

| | Unit of identity | Where arguments live | Column spelling reported | Result types | Function signatures |
|---|---|---|---|---|---|
| PostgreSQL | family row, plus one row per array type | opaque typmod and dims on the use site | canonical, decoded by `format_type` | family, typmod only for a table column | full, over families and pseudo-types |
| MySQL | none: a spelling | decoded columns beside the spelling | canonical spelling | full, computed by the binder | none for built-ins |
| SQLite | none: the declared text | none: decoration | verbatim | storage class of the value | arity only |
| ClickHouse | family, with aliases | in the spelling | canonical expression | full, computed by the binder | none |
| DuckDB | family, with aliases | in the spelling, plus decoded columns | canonical expression, some arguments dropped | full, computed by the binder | families and generics |
| SQL Server | family row and alias-type row | decoded columns on the use site | canonical, defaults filled in | family and decoded arguments | families and decoded arguments |
| GoogleSQL | none offline: a spelling | in the spelling | canonical spelling | full | none offline |

Two conclusions follow. First, no engine keeps arguments inside its type
table: PostgreSQL and SQL Server put them on the use site and everything
else puts them in the spelling. A row per expression is nonetheless the
right shape for sqlc, because sqlc's job is to *report* the expression, and
every engine that decodes its use-site arguments decodes them into exactly
the call expression `TypeExpr` already is; PostgreSQL's typmod is a
per-type encoding that only PostgreSQL can read, so sqlc would store the
decoded form either way. Second, every engine with a catalog canonicalizes
on the way in and reports the canonical form, never the declared spelling,
and the same is true of what `goldeneye` checks the analyze cases against.
The design below changes on that point.

## The design

### One row per type expression

`sql_type` keeps one row per distinct type expression. A row is either a
**family** — a name the dialect or the schema declares, such as `numeric`,
`array`, `struct`, `mood` — or an **instance**, a family applied to
arguments, such as `numeric(10, 2)`, `array(int4)` or
`struct(a: int4, b: text)`. The row's `expr` is the expression's canonical
string, which is its interning key; the row's `name` is the family's name,
so the index that turns `integer` into a row keeps working for instances,
and an instance points at its family.

```sql
CREATE TABLE sql_type (
    oid           INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace_oid INTEGER NOT NULL REFERENCES sql_namespace(oid),
    dialect_oid   INTEGER REFERENCES sql_dialect(oid),
    name          TEXT NOT NULL,   -- the family name: 'numeric', 'array', 'mood'
    expr          TEXT NOT NULL,   -- the whole expression, canonical: 'numeric(10, 2)'; equals name for a family
    typtype       TEXT NOT NULL DEFAULT 'b',  -- b base, c composite, d domain, e enum, r range, p pseudo
    category      TEXT,
    preferred     INTEGER NOT NULL DEFAULT 0,
    family_oid    INTEGER REFERENCES sql_type(oid),  -- NULL on a family; the family on an instance
    element_oid   INTEGER REFERENCES sql_type(oid),  -- what the type holds: an array's element, a map's value, a range's subtype
    base_oid      INTEGER REFERENCES sql_type(oid),  -- what the type stands on: a domain's or alias type's base, a wrapper's inner type, a SQLite spelling's affinity
    canonical_oid INTEGER REFERENCES sql_type(oid),  -- the row the engine reports this one as: integer -> int4, float(24) -> real, mood -> enum('sad', 'ok') in DuckDB
    not_null      INTEGER NOT NULL DEFAULT 0,        -- a domain or alias type declared NOT NULL
    UNIQUE (namespace_oid, expr)
);
CREATE INDEX idx_sql_type_name ON sql_type(name);

-- sql_type_arg: the arguments of an instance, or the fields, labels or
-- members of a declared composite, enum or set, in order. Exactly one of
-- arg_type_oid, int_value, bool_value, string_value and ident is set.
CREATE TABLE sql_type_arg (
    type_oid     INTEGER NOT NULL REFERENCES sql_type(oid),
    ord          INTEGER NOT NULL,
    label        TEXT NOT NULL DEFAULT '',  -- a struct field, tuple element or enum label
    arg_type_oid INTEGER REFERENCES sql_type(oid),
    nullable     INTEGER NOT NULL DEFAULT 0, -- the argument type is nullable here: Array(Nullable(String))
    int_value    INTEGER,
    bool_value   INTEGER,
    string_value TEXT,
    ident        TEXT,                        -- a bare word that is not a type: max, sum, day to second
    PRIMARY KEY (type_oid, ord)
);
```

`size` goes: nothing reads it. `type_length` and `type_scale` leave
`sql_attribute`: they were one engine's two arguments, and now every
engine's arguments are rows. `decl_type` stays, since the verbatim spelling
is what the formatter prints back and what SQLite reports.

The four pointer columns are the denormalization. Each is derivable by
walking `sql_type_arg`, and each is what resolution asks for in one
statement:

- `family_oid` is what operator and function lookup fall back to when there
  is no overload on the instance: `numeric(10, 2) + numeric(5, 1)` finds no
  operator on either instance and resolves on `numeric`, which is exactly
  what `pg_operator` and `duckdb_functions()` hold.
- `element_oid` is what a subscript, `ANY($1)`, `unnest` and a star over a
  map yield, and what `IsArray` was; it is PostgreSQL's `typelem` and
  `pg_range.rngsubtype`.
- `base_oid` is what a domain, an alias type, `LowCardinality(T)` or a
  SQLite spelling resolves through: `posint = 1` finds no operator on
  `posint` and resolves on `int4`; `FOO BAR(3) = 1` resolves on `numeric`.
  It is `typbasetype`, SQL Server's `system_type_id` behind a
  `user_type_id`, and SQLite's affinity rule.
- `canonical_oid` is the row the engine would report this one as. An alias
  spelling points at its family (`integer` at `int4`); a declared type
  points at what the engine reports instead of its name when it does so
  (DuckDB's `mood` at `enum('sad', 'ok')`); an instance points at the
  instance canonicalization rewrote it to when the rewrite changes the
  family (SQL Server's `float(24)` at `real`). It replaces the n² implicit
  casts between spellings of one type.

A declared type is a family row with arguments of its own. `CREATE TYPE
point2 AS (x float8, y float8)` is a row named `point2`, `typtype` `c`, with
two labelled type arguments, the way PostgreSQL keeps them as the attributes
of `typrelid`; `CREATE TYPE mood AS ENUM ('sad', 'ok')` is a row with
`typtype` `e` and two string arguments in order, which is `pg_enum` and
DuckDB's `labels`; ClickHouse's `Enum8('active' = 1, 'deleted' = 2)` is an
instance of `enum8` with two labelled integer arguments and `base_oid` at
`int8`; `CREATE DOMAIN posint AS integer` and SQL Server's `CREATE TYPE
PhoneNumber FROM varchar(20) NOT NULL` are rows with `typtype` `d`,
`base_oid` at the base instance and `not_null` set. The same table holds
what the schema names and what it constructs anonymously.

### Canonicalize on the way in; report the canonical form

Every engine with a catalog rewrites a declared type before storing it and
reports the rewritten form, and `goldeneye` checks the analyze cases against
what the engine reports. So interning canonicalizes, in three steps that
are dialect data or dialect code:

1. **Aliases**, which `types.jsonl` already lists: `int` and `integer`
   resolve to the canonical family. This also settles what the canonical
   *name* is: the one the engine reports, which for PostgreSQL is
   `format_type`'s (`integer`, `bigint`, `character varying`) rather than
   `pg_type`'s (`int4`, `int8`, `varchar`), so PostgreSQL's `types.jsonl`
   flips which spelling is the name and which the alias. Codegen accepts
   both spellings already.

   `format_type` is a string and the output is an expression, so the
   spelling is read into one, and its grammar is slightly wider than
   `name(args)`: the typmod may sit inside a multi-word name
   (`timestamp(3) without time zone`, `time(4) with time zone`), after
   trailing words (`interval day to second(3)`), or before an array suffix
   (`numeric(10,2)[]`, `character varying(255)[]`), and one name is quoted
   (`"char"`). The family is the words with the parenthesis lifted out, so
   `timestamp(3) without time zone` is `{name: "timestamp without time
   zone", args: [3]}`, `interval day to second(3)` is `interval` with an
   identifier argument and an integer one, and `numeric(10,2)[]` is
   `array(numeric(10, 2))`. Two spellings are the same family under
   different typmods, `bpchar` and `character(5)`, and canonicalize to
   `character`. The PostgreSQL goldens change with this: `bigserial`
   becomes `bigint` (a serial is a default, not a type PostgreSQL reports),
   `int4` becomes `integer`, `varchar` becomes `character varying`.
2. **Argument defaults and drops**, which are data: SQL Server's `varchar`
   is `varchar(1)` and `decimal` is `decimal(18, 0)`; DuckDB's `numeric` is
   `decimal(18, 3)` and `varchar(10)` is `varchar`; PostgreSQL's `int[3]` is
   `array(int4)`. A family in `types.jsonl` may say `"defaults": [18, 0]`
   or `"args": 0`.
3. **Rewrites that change the family by argument or member**, which are
   code, since they are ClickHouse's enum numbering and MySQL's charset
   folding: `Decimal32(s)` is `Decimal(9, s)`, `Enum('a', 'b')` is
   `Enum8('a' = 1, 'b' = 2)`, `Variant(...)` sorts its members, `varchar(n)
   character set binary` is `varbinary(n)`, `boolean` is `tinyint(1)`,
   `float(24)` is `real`. An engine package registers a `Canonicalize(*TypeExpr)`
   hook with its seed, and the catalog applies it before interning.

The reported type is the canonical row's expression. The declared spelling
is kept on the attribute in `decl_type`, for the formatter and for SQLite,
whose canonical form is the spelling itself. This reverses the earlier
choice of reporting a column as it was declared: the engines do not, and
the check that keeps the analyzer honest compares against the engines.

### Nullability is not a type

A row is never nullable. Outer nullability stays where it is: on the
attribute, on the analyzer's `exprType`, on the reported column. Inner
nullability is a flag on the argument position, so `Array(Nullable(String))`
is an instance of `array` whose one argument is `string` with `nullable`
set, and `Nullable(String)` in a cast is `string` with the expression's own
nullability set. This is what `TypeExpr` already says — `nullable` at
whatever depth it applies, never a wrapper — and it keeps `string` and
`string nullable` from being two types that need their own operators. The
review confirms only ClickHouse spells inner nullability at all: DuckDB's
`[1, NULL]` is `INTEGER[]` and PostgreSQL has no such thing.

### Array dimensions

PostgreSQL's type for `int[][]` is `integer[]`, with the dimensions on the
attribute and unenforced; DuckDB and ClickHouse nest, `INTEGER[][]` and
`Array(Array(UInt8))`, and DuckDB distinguishes the fixed-size
`INTEGER[3]`. The expression nests one `array` per declared dimension in
every dialect, with a fixed size as a second argument (`array(integer, 3)`),
because codegen renders one slice per dimension and a two-dimensional
PostgreSQL column has to come out `[][]int32` as it does on the legacy path.
PostgreSQL's canonicalization therefore keeps the dimensions its own
catalog drops, and a future PostgreSQL check in `goldeneye` reads `attndims`
to reproduce them. Subscripting follows the dialect: PostgreSQL yields the
innermost element however many subscripts are applied.

### The catalog interns at schema time; the analyzer looks up at query time

The cached catalog is opened read-only (`catalog.go` opens it with
`mode=ro&immutable=1`), and `exprType` carries a bare name precisely so that
analysis never has to write. That constraint stands, and it decides which
expressions become rows:

- **Interned**: what the dialect seeds and what the schema declares. Every
  column type, every declared type's fields and base, every function
  signature's argument and return type becomes a row on load, through one
  entry point, `ResolveTypeExpr(*TypeExpr) (oid, error)`, which
  canonicalizes, then walks the expression bottom-up, interning each
  argument type first. `ResolveType` and `ResolveTypeName` become callers
  of it.
- **Looked up**: what a query writes. A cast, a constructed array or struct,
  a typed placeholder and a function result are resolved by
  `LookupTypeExpr(*TypeExpr) (oid, familyOID, bool)`, which canonicalizes,
  finds the instance row when the schema happened to declare the same
  expression and otherwise the family row, and never writes.

`exprType` becomes the pair:

```go
type exprType struct {
	typeOID  int64          // the instance row when the catalog has one, else the family row, else 0
	expr     *core.TypeExpr // the whole expression, whenever anything is known about it
	nullable bool
	...
}
```

Resolution uses `typeOID` and its `family_oid`, `base_oid` and
`canonical_oid` chain; reporting uses `expr`. The `typeName` fallback and
the `[]` suffix convention go away: an array is `array` applied to its
element, in the catalog as in the output, and `TypeNameString` is replaced
by a function that reads an `ast.TypeName` into a `TypeExpr`, folding
`Typmods` into integer arguments, `ArrayBounds` into one `array` per
dimension, `Names` into a namespace and a name, and `Spelling` through
`ParseTypeExpr`.

### Result types: the family from the catalog, the arguments from the dialect

MySQL, ClickHouse and DuckDB compute a full type for every expression in the
binder — `decimal(10,2) + 1` is `decimal(11,2)`, `Int8 + UInt8` is `Int16`,
`SUM(INTEGER)` is `HUGEINT` — and none of them keeps those rules in a
catalog; PostgreSQL's catalog says `numeric + numeric` is `numeric` and its
results drop the typmod. The catalog can only ever answer at the family
level, and that is the baseline every dialect gets: an operator or function
result is the family the overload names, with the arguments of an `$n`
result carried over from the argument it stands for. A dialect that reports
more registers a `ResultType(op, args []*TypeExpr) *TypeExpr` hook beside
its `Canonicalize` hook, and the analyzer applies it to the expression it
reports. ClickHouse needs it for its arithmetic and for value-dependent
results like `toDecimal64(x, 4)`, since its check compares whole
expressions; MySQL's check reads the wire type, which is the family with
its flags, so the baseline passes it.

### What each engine hands the core

The contract with an engine is that its `ast.TypeName` reads into a
`TypeExpr` that says everything the engine's own catalog would. Where an
engine already folds its type into a spelling, `ParseTypeExpr` reads it;
where it does not, the converter has a small change to make.

| Engine | Form | Expression |
|---|---|---|
| PostgreSQL | `numeric(10,2)`, `varchar(255)`, `timestamp(3) with time zone` | `numeric(10, 2)`, `character varying(255)`, `timestamp with time zone(3)`: typmods become integer arguments on the canonical family, named as `format_type` names it |
| PostgreSQL | `int[]`, `int[][]`, `int[3]` | `array(integer)`, `array(array(integer))`, `array(integer)`: one per array bound, the bound itself dropped as PostgreSQL drops it |
| PostgreSQL | `interval day to second` | `interval('day to second')`: the fields are one identifier argument, as `format_type` prints them |
| PostgreSQL | domain, composite, enum, range | declared rows, as above; `CreateDomainStmt`, `CompositeTypeStmt` and `CreateRangeStmt` gain `schema.Apply` cases |
| PostgreSQL | `myschema.mood` | a row in namespace `myschema`; `Names` resolves to a namespace rather than a dotted name |
| MySQL | `BIGINT UNSIGNED`, `DECIMAL(10,2) UNSIGNED` | `bigint unsigned`, `decimal unsigned(10, 2)`: unsigned is a family of its own, as `COLUMN_TYPE` and the wire flag report it, rather than the alias of the signed type `types.jsonl` lists today; the converter puts it in the name instead of on `ColumnDef.IsUnsigned`, which the core ignores |
| MySQL | `TINYINT(1)`, `DATETIME(6)`, `VARCHAR(255)`, `BOOLEAN` | `tinyint(1)`, `datetime(6)`, `varchar(255)`, `tinyint(1)`: the converter's `Typmods` are read, and `boolean` canonicalizes as MySQL does |
| MySQL | `ENUM('a','b')`, `SET('x','y')` | `enum('a', 'b')`, `set('x', 'y')`: the converter renders `Vals` into the spelling |
| MySQL | `CAST(x AS CHAR(10))` | `char(10)`: the cast converter names the SQL type, not the wire code `var_string` |
| MySQL | `VARCHAR(10) CHARACTER SET binary` | `varbinary(10)`, by the canonicalization hook; collation is not part of the type |
| SQLite | any spelling | the spelling as an instance, `varchar(255)`, `foo bar(3)`, verbatim as `pragma_table_xinfo` reports it, with `base_oid` set by the affinity rule applied when the dialect resolves an unknown name; `types.jsonl`'s alias lists become the rule |
| SQLite | `STRICT` tables, `ANY` | the family rows; a strict table's column names one of them or fails, as SQLite does |
| SQLite | an expression | its storage class, `integer`, `real`, `text` or `blob`, which is what `typeof()` and the check report |
| ClickHouse | every parametric type | the spelling, read as today, now also for casts, `{p:T}` placeholders and results |
| ClickHouse | `Nullable(T)`, `LowCardinality(T)` | `T` with `nullable`; `lowcardinality(T)` with `base_oid` at `T` |
| ClickHouse | `Decimal32(4)`, `Enum('a', 'b')`, `Variant(String, Int64)`, `INT` | `decimal(9, 4)`, `enum8(a: 1, b: 2)`, `variant(int64, string)`, `int32`, by the canonicalization hook |
| ClickHouse | `SimpleAggregateFunction(sum, UInt64)` | `simpleaggregatefunction(sum, uint64)` with `sum` an identifier argument |
| ClickHouse | `toDecimal64(x, s)`, `Int8 + UInt8` | `decimal(18, s)`, `int16`, by the result-type hook |
| ClickHouse | `Nested(a UInt8, b String)` | a relation-shape rule, not a type: the column becomes `n.a array(uint8)` and `n.b array(string)` on load, as `system.columns` has them |
| DuckDB | `STRUCT(a INTEGER, b VARCHAR)`, `MAP(K, V)`, `UNION(...)` | `struct(a: integer, b: varchar)`, `map(varchar, integer)`, `union(num: integer, str: varchar)`: the converter renders the darkwing type expression it already has instead of keeping its name |
| DuckDB | `INTEGER[]`, `INTEGER[3]` | `array(integer)` and `array(integer, 3)`: a list is the cross-dialect array, a fixed size is its second argument |
| DuckDB | `TEXT`, `VARCHAR(10)`, `NUMERIC`, `mood` | `varchar`, `varchar`, `decimal(18, 3)`, `enum('sad', 'ok')`, by aliases, argument rules and `canonical_oid` on the declared enum |
| GoogleSQL | `ARRAY<INT64>`, `STRUCT<a INT64, b STRING>`, `RANGE<DATE>` | `array(int64)`, `struct(a: int64, b: string)`, `range(date)`: the converter renders the zetajones type node in call form, or `ParseTypeExpr` accepts `<...>` |
| GoogleSQL | `STRING(10)`, `STRING(MAX)`, `NUMERIC(10,2)`, `[1, 2]`, `STRUCT(1 AS x)` | the typmods are read, with `max` an identifier argument; the array and struct constructors are typed from their elements |
| SQL Server | `NVARCHAR(MAX)`, `VARCHAR`, `DECIMAL` | `nvarchar(max)` with `max` an identifier argument; `varchar(1)` and `decimal(18, 0)` by the defaults `sys.types` applies |
| SQL Server | `FLOAT(24)` | `float(24)` with `canonical_oid` at `real`, by the hook |
| SQL Server | `dbo.PhoneNumber`, `sysname` | a row in namespace `dbo`, `typtype` `d`, `base_oid` at `varchar(20)`, `not_null` set; `sysname` seeded the same way over `nvarchar(128)` |

The `ident` argument is the one addition to `TypeExpr` and to
`goldeneye/analysis.TypeExpr`, which mirrors it. Without it `max`, `sum` and
`day to second` read as types, and today they do: `SimpleAggregateFunction`
reports `sum` as a type. Both parsers that have `MAX` model it as a literal
node, and the converters hand it over as an identifier; `ParseTypeExpr`
writes a bare word that is not a known family as an identifier only when a
dialect says so, since a struct field's type is also a bare word.

### What the analyzer reports

`Column` and `Parameter` in `analysis.go` keep `Type` as the expression and
`TypeOID` as the row, instance or family. `DataType` and `IsArray` stay as
the flat view for the legacy compiler bridge, and that bridge derives what
codegen reads from the expression rather than dropping it: `ArrayDims` is
the depth of `array` nesting, `Length` is the first integer argument,
`Unsigned` is a family name ending in ` unsigned`. `sqlc analyze` prints the
expression as it does today, with `ident` as a fifth argument kind.

`TypeNameString`, `ArraySuffix`, `CreateArrayType`, `TypeLength` and
`TypeScale` are the API that goes; `ResolveTypeExpr`, `LookupTypeExpr` and
`TypeExprOf(oid)` — the expression a row stands for, read back from
`sql_type_arg` — are the API that replaces it.

### What the seed files gain

`types.jsonl` is unchanged for a family beyond the optional argument
defaults; an alias becomes a row with `canonical_oid` instead of a mesh of
casts. A function's argument and return types may be expressions, and the
return type may reference an argument's value as well as its type. The
category rules in `dialect.json` apply to families; an instance inherits
its family's category, which is how `numeric(10, 2)` joins the numeric
casts without being seeded. An engine package may register two hooks with
its seed, `Canonicalize` and `ResultType`, for what its catalog does in
code.

`goldeneye` checks the analyze cases against what each database reports, and
its answer shape is the same `TypeExpr`. ClickHouse and DuckDB report whole
expressions already. MySQL's driver reports the family, unsigned and
nullability but not precision or length, and `ColumnType.DecimalSize` and
`ColumnType.Length` can add them, with the length divided by the charset's
bytes per character since the wire reports bytes. SQLite reports a declared
spelling for a table column and a storage class for an expression, which
are an instance row and a family row respectively. A PostgreSQL check reads
`format_type(atttypid, atttypmod)` and `attndims`. Every check keeps passing
on the way, since a family with no arguments prints as it does now.

## A worked example: an array of arrays of integers

```sql
CREATE TABLE grids (
  id    bigint PRIMARY KEY,
  cells int[][] NOT NULL
);

-- name: GetGrid :one
SELECT id, cells, cells[1][2] AS cell FROM grids WHERE cells = $1;
```

The PostgreSQL parser hands over `pg_catalog.int4` with two array bounds,
which reads into `array(array(int4))`; canonicalization makes it
`array(array(integer))`. Interning walks it bottom-up, so the inner array
gets its row before the outer one. With illustrative OIDs, `sql_type`
holds the seeded families and the two rows the schema added:

| oid | name | expr | typtype | category | family_oid | element_oid | canonical_oid |
|---|---|---|---|---|---|---|---|
| 23 | integer | integer | b | N | | | |
| 24 | int4 | int4 | b | N | | | 23 |
| 20 | bigint | bigint | b | N | | | |
| 100 | array | array | b | A | | | |
| 1001 | array | array(integer) | b | A | 100 | 23 | |
| 1002 | array | array(array(integer)) | b | A | 100 | 1001 | |

Row 1001 is what PostgreSQL calls `_int4`; row 1002 is the one PostgreSQL
does not have, since its own type collapses the dimensions onto `attndims`.
An instance takes its family's namespace. `sql_type_arg` has one row per
argument position:

| type_oid | ord | label | arg_type_oid | nullable |
|---|---|---|---|---|
| 1001 | 1 | | 23 | 0 |
| 1002 | 1 | | 1001 | 0 |

and `sql_attribute` points `grids.cells` at row 1002, `not_null` set, with
`int[][]` in `decl_type`.

Analysis resolves `cells` to that attribute, so its `exprType` is
`{typeOID: 1002, expr: array(array(integer)), nullable: false}`, and the
parameter compared with it takes the same type. `cells[1][2]` follows the
dialect's subscript rule — PostgreSQL yields the innermost element however
many subscripts are applied — which is a walk down `element_oid` from 1002
to 1001 to 23, nullable because a subscript can miss. `sqlc analyze` prints:

```json
[
  {
    "name": "GetGrid",
    "cmd": ":one",
    "columns": [
      { "name": "id", "type": { "name": "bigint" }, "table": "grids" },
      {
        "name": "cells",
        "type": {
          "name": "array",
          "args": [
            { "type": { "name": "array", "args": [ { "type": { "name": "integer" } } ] } }
          ]
        },
        "table": "grids"
      },
      { "name": "cell", "type": { "name": "integer", "nullable": true } }
    ],
    "params": [
      {
        "number": 1,
        "column": {
          "name": "cells",
          "type": {
            "name": "array",
            "args": [
              { "type": { "name": "array", "args": [ { "type": { "name": "integer" } } ] } }
            ]
          },
          "table": "grids"
        }
      }
    ]
  }
]
```

Its string form is `array(array(integer))`; today the same column prints
`array(int4)`, one level, with `type_oid` pointing at a row named `int4[]`.
The legacy bridge derives `DataType` `integer` and `ArrayDims` 2 from the
nesting, so Go codegen renders `[][]int32` as the legacy path does.
DuckDB's `INTEGER[][]` and ClickHouse's `Array(Array(Int32))` produce the
same rows and output, with `int32` in ClickHouse's case.

## As implemented

The design above is implemented, engine by engine, with these departures
and details settled on the way:

- The three per-dialect hooks are registered by an engine package at init,
  by the name its `dialect.json` records, and looked up by that name, so a
  catalog restored from the cache — which runs no seed — has them:
  `core.RegisterCanonicalizer`, `core.RegisterUserTypeBase` (SQLite's
  affinity rule, applied when the schema declares a family the seed does
  not list) and `core.RegisterResultType` (ClickHouse's `toDecimal64(x, 4)`
  and `toDateTime64(x, 3)`).
- SQLite's `dialect.json` says `"alias": "base"`, which makes each alias in
  its `types.jsonl` a type of its own standing on the type it aliases,
  rather than another spelling of it.
- An engine hands the core either a spelling (`TypeName.Spelling`, read by
  `ParseTypeExpr`, which also reads words after a closing parenthesis as
  part of the name, as in `decimal(10,2) unsigned`) or a name with
  `Typmods` and `ArrayBounds`, where an integer constant is an integer
  argument, a bare `ast.String` is an identifier and a quoted constant a
  string. `ColumnDef.IsUnsigned` and `ColumnDef.Vals` add MySQL's unsigned
  and enum members. `ParamRef.Name` carries the name a `{name:Type}`
  placeholder gives itself.
- A cast is NULL when its operand is, or when its type says so, as
  `Nullable(String)` does; a cast of a placeholder types the placeholder
  and takes its name and source from what it is compared with.
- MySQL types `CAST(x AS CHAR(10))` as `varchar(10)` and `CAST(x AS
  BINARY(8))` as `varbinary(8)`, which is what its metadata and a view over
  the cast both report, rather than the `char(10)` the table above
  proposed. `goldeneye` reads a table column's type from `COLUMN_TYPE`, in
  the relations seed and in the analyze check, and compares an expression
  by family alone, since the wire carries no arguments.
- SQLite reports a cast to a spelling that is not a storage class, such as
  `DECIMAL(5,2)`, as that spelling, while the value's storage class is what
  a run would show; the check's cases keep to storage classes.
- DuckDB's `JSON` is grouped under `varchar` by `duckdb_types()`'s logical
  type, so a JSON column reports `varchar`; a named enum reports its name,
  not its labels, since the canonicalizer cannot see the catalog. A
  GoogleSQL array or struct constructor in a select list is still untyped.
- PostgreSQL's `relations.jsonl` still spells array columns as `pg_type`
  does (`_text`), which the canonicalizer reads as `array(text)`.

## Order of work

1. The tables and the interning entry point: `sql_type.expr`, `family_oid`,
   `element_oid`, `base_oid`, `canonical_oid`, `not_null`, `sql_type_arg`,
   `ResolveTypeExpr`, `LookupTypeExpr`, `TypeExprOf`, with the alias step
   of canonicalization. Arrays become instances of `array`; the `[]`
   convention goes. Every existing golden holds, since a bare name prints
   the same.
2. The analyzer: `exprType` carries the expression; casts, constructors,
   placeholders and function results report it; resolution falls back
   through the pointer chain. This is where ClickHouse's casts and
   parameters come right.
3. The engines, one at a time, each with an `analyze_types/<engine>` case
   alongside ClickHouse's and each with its `Canonicalize` hook: PostgreSQL
   typmods, dimensions, declared types and `format_type` names; MySQL
   unsigned, typmods, members and `boolean`, plus the `var_string` leak;
   DuckDB's nested types and dropped arguments; GoogleSQL's angle brackets;
   SQL Server's `max`, defaults and alias types; SQLite's affinity rule.
4. The legacy bridge: `parse_core.go` derives `Unsigned`, `Length` and
   `ArrayDims` from the expression, and the `experiment_coreanalyzer` cases
   grow MySQL unsigned and boolean columns and a PostgreSQL two-dimensional
   array, so the core path generates what the legacy path does.
5. Result-type hooks, ClickHouse first, and the value-dependent return
   types its seed needs; then MySQL's precision arithmetic once its check
   reads precision from the wire.

## Open questions

- How far a result-type hook goes. ClickHouse's arithmetic promotion and
  `toDecimal64(x, 4)` are finite rules; `arrayMap(f, arr)` returns an array
  of the lambda's result, which needs the lambda typed first.
- Whether DuckDB's dropped `VARCHAR(10)` length and reported anonymous enum
  should be canonicalized away as DuckDB does, or kept because a user
  declared them. The check decides for the former.
