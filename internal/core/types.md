# Types in the analysis core

How the core catalog, the analyzer and `sqlc analyze` represent a type today,
where that falls short of each engine's type system, and the design that
closes the gap: every type the schema or the dialect declares is a row holding
its full expression, denormalized so that a bare name like `integer` and a
structured expression like `numeric(10, 2)` both resolve to one.

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
    base_oid      INTEGER REFERENCES sql_type(oid),  -- what the type stands on: a domain's or alias type's base, a wrapper's inner type
    canonical_oid INTEGER REFERENCES sql_type(oid),  -- the row an alias spelling means: integer -> int4
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
is what the formatter prints back.

The four pointer columns are the denormalization. Each is derivable by
walking `sql_type_arg`, and each is what resolution asks for in one
statement:

- `family_oid` is what operator and function lookup fall back to when there
  is no overload on the instance: `numeric(10, 2) + numeric(5, 1)` finds no
  operator on either instance and resolves on `numeric`.
- `element_oid` is what a subscript, `ANY($1)`, `unnest` and a star over a
  map yield, and what `IsArray` was.
- `base_oid` is what a domain, an alias type, `LowCardinality(T)` or SQLite's
  affinity resolves through: `posint = 1` finds no operator on `posint` and
  resolves on `int4`.
- `canonical_oid` replaces the n² implicit casts between spellings of one
  type: `integer` and `int4` are two rows, so a column reports the spelling
  it was declared with, and one canonical row, so resolution treats them as
  one type.

A declared type is a family row with arguments of its own. `CREATE TYPE
point2 AS (x float8, y float8)` is a row named `point2`, `typtype` `c`, with
two labelled type arguments; `CREATE TYPE mood AS ENUM ('sad', 'ok')` is a
row with `typtype` `e` and two string arguments; ClickHouse's
`Enum8('active' = 1, 'deleted' = 2)` is an instance of `enum8` with two
labelled integer arguments and `base_oid` pointing at `int8`; `CREATE DOMAIN
posint AS integer` and SQL Server's `CREATE TYPE PhoneNumber FROM varchar(20)
NOT NULL` are rows with `typtype` `d`, `base_oid` at the base instance and
`not_null` set. The same table holds what the schema names and what it
constructs anonymously.

### Nullability is not a type

A row is never nullable. Outer nullability stays where it is: on the
attribute, on the analyzer's `exprType`, on the reported column. Inner
nullability is a flag on the argument position, so `Array(Nullable(String))`
is an instance of `array` whose one argument is `string` with `nullable`
set, and `Nullable(String)` in a cast is `string` with the expression's own
nullability set. This is what `TypeExpr` already says — `nullable` at
whatever depth it applies, never a wrapper — and it keeps `string` and
`string nullable` from being two types that need their own operators.

### The catalog interns at schema time; the analyzer looks up at query time

The cached catalog is opened read-only (`catalog.go` opens it with
`mode=ro&immutable=1`), and `exprType` carries a bare name precisely so that
analysis never has to write. That constraint stands, and it decides which
expressions become rows:

- **Interned**: what the dialect seeds and what the schema declares. Every
  column type, every declared type's fields and base, every function
  signature's argument and return type becomes a row on load, through one
  entry point, `ResolveTypeExpr(*TypeExpr) (oid, error)`, which walks the
  expression bottom-up, interning each argument type first. `ResolveType`
  and `ResolveTypeName` become callers of it.
- **Looked up**: what a query writes. A cast, a constructed array or struct,
  a typed placeholder and a function result are resolved by
  `LookupTypeExpr(*TypeExpr) (oid, familyOID, bool)`, which finds the
  instance row when the schema happened to declare the same expression and
  otherwise the family row, and never writes.

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

### What each engine hands the core

The contract with an engine is that its `ast.TypeName` reads into a
`TypeExpr` that says everything the engine's own catalog would. Where an
engine already folds its type into a spelling, `ParseTypeExpr` reads it;
where it does not, the converter has a small change to make.

| Engine | Form | Expression |
|---|---|---|
| PostgreSQL | `numeric(10,2)`, `varchar(255)`, `timestamp(3) with time zone` | `numeric(10, 2)`, `varchar(255)`, `timestamptz(3)`: typmods become integer arguments on the canonical family |
| PostgreSQL | `int[]`, `int[][]` | `array(int4)`, `array(array(int4))`: one per array bound |
| PostgreSQL | `interval day to second` | `interval('day to second')`: the fields are one identifier argument, as `format_type` prints them |
| PostgreSQL | domain, composite, enum, range | declared rows, as above; `CreateDomainStmt`, `CompositeTypeStmt` and `CreateRangeStmt` gain `schema.Apply` cases |
| PostgreSQL | `myschema.mood` | a row in namespace `myschema`; `Names` resolves to a namespace rather than a dotted name |
| MySQL | `BIGINT UNSIGNED`, `DECIMAL(10,2) UNSIGNED` | `bigint unsigned`, `decimal unsigned(10, 2)`: unsigned is a family of its own, as the server reports it, rather than the alias of the signed type `types.jsonl` lists today, since the value range differs; the converter puts it in the name instead of on `ColumnDef.IsUnsigned`, which the core ignores |
| MySQL | `TINYINT(1)`, `DATETIME(6)`, `VARCHAR(255)` | `tinyint(1)`, `datetime(6)`, `varchar(255)`: the converter's `Typmods` are read |
| MySQL | `ENUM('a','b')`, `SET('x','y')` | `enum('a', 'b')`, `set('x', 'y')`: the converter renders `Vals` into the spelling |
| MySQL | `CAST(x AS CHAR(10))` | `char(10)`: the cast converter names the SQL type, not TiDB's `var_string` |
| MySQL | `CHARACTER SET binary`, `COLLATE` | a labelled string argument, `varchar(255, charset: 'binary')`, if codegen ever needs it; open |
| SQLite | any spelling | the spelling as an instance, `varchar(255)`, `foo bar(3)`, with `base_oid` set by the affinity rules — INT anywhere is INTEGER, CHAR, CLOB or TEXT is TEXT, BLOB or nothing is BLOB, REAL, FLOA or DOUB is REAL, else NUMERIC — applied when the dialect resolves an unknown name; `types.jsonl`'s alias lists become the rule |
| SQLite | `STRICT` tables, `ANY` | the family rows; a strict table's column names one of them or fails |
| ClickHouse | every parametric type | the spelling, read as today, now also for casts, `{p:T}` placeholders and results |
| ClickHouse | `Nullable(T)`, `LowCardinality(T)` | `T` with `nullable`; `lowcardinality(T)` with `base_oid` at `T` |
| ClickHouse | `SimpleAggregateFunction(sum, UInt64)` | `simpleaggregatefunction(sum, uint64)` with `sum` an identifier argument |
| ClickHouse | `toDecimal64(x, s)` | the seed's return type may name an argument's *value*: `"returns": "Decimal(18, $2)"`; the analyzer substitutes the literal when it is one and reports the family otherwise |
| ClickHouse | `Nested(a UInt8, b String)` | a relation-shape rule, not a type: the column becomes `n.a array(uint8)` and `n.b array(string)` on load |
| DuckDB | `STRUCT(a INTEGER, b VARCHAR)`, `MAP(K, V)`, `UNION(...)` | `struct(a: integer, b: varchar)`, `map(varchar, integer)`, `union(num: integer, str: varchar)`: the converter renders the darkwing type expression it already has instead of keeping its name |
| DuckDB | `INTEGER[]`, `INTEGER[3]` | `array(integer)` and `array(integer, 3)`: a list is the cross-dialect array, a fixed size is its second argument |
| GoogleSQL | `ARRAY<INT64>`, `STRUCT<a INT64, b STRING>`, `RANGE<DATE>` | `array(int64)`, `struct(a: int64, b: string)`, `range(date)`: `ParseTypeExpr` accepts `<...>` as well as `(...)`, or the converter renders the column schema in call form |
| GoogleSQL | `STRING(10)`, `NUMERIC(10,2)`, `[1, 2]`, `STRUCT(1 AS x)` | the typmods are read; the array and struct constructors are typed from their elements |
| SQL Server | `NVARCHAR(MAX)` | `nvarchar(max)` with `max` an identifier argument |
| SQL Server | `FLOAT(24)` | `float(24)`, with `canonical_oid` at `real`: a dialect may declare an instance canonical to a family, `"canonical": {"float(1..24)": "real"}`; open |
| SQL Server | `dbo.PhoneNumber` | a row in namespace `dbo`, `typtype` `d`, `base_oid` at `varchar(20)`, `not_null` set |

The `ident` argument is the one addition to `TypeExpr` and to
`goldeneye/analysis.TypeExpr`, which mirrors it. Without it `max`, `sum` and
`day to second` read as types, and today they do: `SimpleAggregateFunction`
reports `sum` as a type. `ParseTypeExpr` writes a bare word that is not a
known family as an identifier only when a dialect says so — ClickHouse for
the first argument of an aggregate-function type, SQL Server for `max` — and
otherwise as a type, since a struct field's type is also a bare word.

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

`types.jsonl` is unchanged for a family; an alias becomes a row with
`canonical_oid` instead of a mesh of casts. A function's argument and return
types may be expressions, and the return type may reference an argument's
value as well as its type. The category rules in `dialect.json` apply to
families; an instance inherits its family's category, which is how
`numeric(10, 2)` joins the numeric casts without being seeded.

`goldeneye` checks the analyze cases against what each database reports, and
its answer shape is the same `TypeExpr`. ClickHouse reports whole
expressions already. MySQL's driver reports the family, unsigned and
nullability but not precision or length, and `ColumnType.DecimalSize` and
`ColumnType.Length` can add them, with the length divided by the charset's
bytes per character. SQLite reports a declared spelling for a table column
and a storage class for an expression, which are a family row and an
instance row respectively. Every check keeps passing on the way, since a
family with no arguments prints as it does now.

## Order of work

1. The tables and the interning entry point: `sql_type.expr`, `family_oid`,
   `element_oid`, `base_oid`, `canonical_oid`, `not_null`, `sql_type_arg`,
   `ResolveTypeExpr`, `LookupTypeExpr`, `TypeExprOf`. Arrays become
   instances of `array`; the `[]` convention goes. Every existing golden
   holds, since a bare name prints the same.
2. The analyzer: `exprType` carries the expression; casts, constructors,
   placeholders and function results report it; resolution falls back
   through the pointer chain. This is where ClickHouse's casts and
   parameters come right.
3. The engines, one at a time, each with an `analyze_types/<engine>` case
   alongside ClickHouse's: PostgreSQL typmods, dimensions and declared types;
   MySQL unsigned, typmods and members, plus the `var_string` leak; DuckDB's
   nested types; GoogleSQL's angle brackets; SQL Server's `max` and alias
   types; SQLite's affinity rule.
4. The legacy bridge: `parse_core.go` derives `Unsigned`, `Length` and
   `ArrayDims` from the expression, and the `experiment_coreanalyzer` cases
   grow MySQL unsigned and boolean columns and a PostgreSQL two-dimensional
   array, so the core path generates what the legacy path does.
5. Aliases through `canonical_oid`, dropping the alias casts, and the
   value-dependent return types for ClickHouse.

## Open questions

- Whether MySQL's character set and collation are worth an argument. Codegen
  reads neither today, but a `binary` charset changes what a driver returns.
- Whether a dialect should be able to declare an instance canonical to
  another family, as SQL Server's `float(24)` is `real`, or whether reporting
  `float(24)` and leaving the mapping to codegen is enough.
- How far a value-dependent return type goes. ClickHouse's `toDecimal64(x,
  4)` is the common case and a literal covers it; `arrayMap(f, arr)` returns
  an array of the lambda's result, which no seed can spell.
