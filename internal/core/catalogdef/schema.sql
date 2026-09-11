-- sql_namespace: schemas / namespaces
CREATE TABLE sql_namespace (
    oid     INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT NOT NULL UNIQUE
);

-- sql_dialect: registered SQL dialects (postgresql, sqlite, mysql, ...).
CREATE TABLE sql_dialect (
    oid     INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT NOT NULL UNIQUE
);

-- sql_dialect_flag: per-dialect configuration knobs (case-folding,
-- identifier quoting, alias scoping rules, etc.). Values are opaque
-- strings; the analyzer interprets them per key.
CREATE TABLE sql_dialect_flag (
    dialect_oid INTEGER NOT NULL REFERENCES sql_dialect(oid),
    key         TEXT NOT NULL,
    value       TEXT NOT NULL,
    PRIMARY KEY (dialect_oid, key)
);

-- sql_type: data types, one row per type expression. Modeled on pg_type,
-- with the arguments PostgreSQL keeps as a typmod on the use site held in
-- the row instead.
--
-- A row is a family — a name the dialect or the schema declares: numeric,
-- array, mood — or an instance, a family applied to arguments: numeric(10, 2),
-- array(integer). expr is the canonical spelling of the whole expression and
-- the row's identity; name is the family's name, on an instance too, so a
-- lookup by name finds the family and an instance points at it.
--
--   typtype:       'b'ase | 'c'omposite | 'd'omain | 'e'num | 'p'seudo | 'r'ange
--   category:      'N'umeric | 'S'tring | 'B'oolean | 'D'atetime | 'A'rray |
--                  'C'omposite | 'E'num | 'U'serdef | 'X'unknown
--   preferred:     tie-breaker for implicit cast resolution within a category
--   family_oid:    NULL on a family; the family on an instance
--   element_oid:   what the type holds: an array's element, a map's value,
--                  a range's subtype
--   base_oid:      what the type stands on: a domain's or alias type's base,
--                  a wrapper's inner type, a SQLite spelling's affinity
--   canonical_oid: the row the engine reports this one as: an alias spelling
--                  points at the type it names
--   not_null:      a domain or alias type declared NOT NULL
--   dialect_oid:   NULL = standard / shared across dialects
CREATE TABLE sql_type (
    oid           INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace_oid INTEGER NOT NULL REFERENCES sql_namespace(oid),
    dialect_oid   INTEGER REFERENCES sql_dialect(oid),
    name          TEXT NOT NULL,
    expr          TEXT NOT NULL,
    typtype       TEXT NOT NULL DEFAULT 'b',
    category      TEXT,
    preferred     INTEGER NOT NULL DEFAULT 0,
    family_oid    INTEGER REFERENCES sql_type(oid),
    element_oid   INTEGER REFERENCES sql_type(oid),
    base_oid      INTEGER REFERENCES sql_type(oid),
    canonical_oid INTEGER REFERENCES sql_type(oid),
    not_null      INTEGER NOT NULL DEFAULT 0,
    UNIQUE (namespace_oid, expr)
);
CREATE INDEX idx_sql_type_name ON sql_type(name);

-- sql_type_arg: the arguments of an instance, or the fields, labels or
-- members of a declared composite, enum or set, in order. Exactly one of
-- arg_type_oid, int_value, bool_value, string_value and ident is set.
--   label:    a struct field, tuple element or enum label
--   nullable: the argument type is nullable at this position, as the inner
--             type of Array(Nullable(String)) is
--   ident:    a bare word that is not a type: max, sum, day to second
CREATE TABLE sql_type_arg (
    type_oid     INTEGER NOT NULL REFERENCES sql_type(oid),
    ord          INTEGER NOT NULL,
    label        TEXT NOT NULL DEFAULT '',
    arg_type_oid INTEGER REFERENCES sql_type(oid),
    nullable     INTEGER NOT NULL DEFAULT 0,
    int_value    INTEGER,
    bool_value   INTEGER,
    string_value TEXT,
    ident        TEXT,
    PRIMARY KEY (type_oid, ord)
);

-- sql_type_rewrite: the rewrites a dialect applies to a type before it is
-- interned, in order, which are how a dialect stores what only it spells:
-- SQL Server keeps float(24) as real and a bare decimal as decimal(18,0),
-- ClickHouse keeps Decimal32(4) as Decimal(9, 4). pattern is a type
-- expression whose arguments may be $1, $2... binding whatever stands there;
-- template is the expression the match becomes, with the bindings
-- substituted; cond bounds a binding, as "$1 <= 24" does.
CREATE TABLE sql_type_rewrite (
    dialect_oid INTEGER NOT NULL REFERENCES sql_dialect(oid),
    ord         INTEGER NOT NULL,
    pattern     TEXT NOT NULL,
    template    TEXT NOT NULL,
    cond        TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (dialect_oid, ord)
);

-- sql_type_affinity: the rule a dialect resolves an unseeded type family
-- by, in order: the first row one of whose words the family's name
-- contains names the type it stands on, and a row with no words is the
-- default. SQLite gives every declared spelling one of five affinities
-- this way.
CREATE TABLE sql_type_affinity (
    dialect_oid INTEGER NOT NULL REFERENCES sql_dialect(oid),
    ord         INTEGER NOT NULL,
    words       TEXT NOT NULL DEFAULT '', -- comma-separated, upper case
    type_oid    INTEGER NOT NULL REFERENCES sql_type(oid),
    PRIMARY KEY (dialect_oid, ord)
);

-- sql_class: relations (tables, views, indexes).
--   kind: 'r' = table, 'v' = view, 'i' = index, 'c' = composite type, 'f' = foreign
CREATE TABLE sql_class (
    oid           INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace_oid INTEGER NOT NULL REFERENCES sql_namespace(oid),
    name          TEXT NOT NULL,
    kind          TEXT NOT NULL DEFAULT 'r',
    UNIQUE(namespace_oid, name)
);

-- sql_attribute: columns of a relation.
--   decl_type:       the type as the schema spelled it, before
--                    canonicalization (VARCHAR(10), BIGINT UNSIGNED), which
--                    is what a formatter prints back and what SQLite reports.
--   auto_increment:  rowid alias (sqlite INTEGER PRIMARY KEY), AUTOINCREMENT,
--                    pg serial/bigserial/identity, mysql AUTO_INCREMENT.
--   is_primary_key:  this column participates in the relation's primary key.
--                    Set both for inline-column PK and for table-level PK.
--   is_unique:       column has a UNIQUE constraint or a single-column UNIQUE
--                    table constraint.
--   hidden:          resolvable by name but absent from a star expansion and
--                    from the relation's model, like the column an sqlite
--                    fts5 table names after itself.
CREATE TABLE sql_attribute (
    oid            INTEGER PRIMARY KEY AUTOINCREMENT,
    class_oid      INTEGER NOT NULL REFERENCES sql_class(oid),
    name           TEXT NOT NULL,
    type_oid       INTEGER NOT NULL REFERENCES sql_type(oid),
    not_null       INTEGER NOT NULL DEFAULT 0,
    has_default    INTEGER NOT NULL DEFAULT 0,
    num            INTEGER NOT NULL, -- ordinal position (1-based)
    decl_type      TEXT    NOT NULL DEFAULT '',
    auto_increment INTEGER NOT NULL DEFAULT 0,
    is_primary_key INTEGER NOT NULL DEFAULT 0,
    is_unique      INTEGER NOT NULL DEFAULT 0,
    hidden         INTEGER NOT NULL DEFAULT 0,
    UNIQUE(class_oid, name),
    UNIQUE(class_oid, num)
);

-- sql_constraint: constraints on a relation.
--   kind: 'p' = primary key, 'f' = foreign key, 'u' = unique, 'c' = check
CREATE TABLE sql_constraint (
    oid       INTEGER PRIMARY KEY AUTOINCREMENT,
    class_oid INTEGER NOT NULL REFERENCES sql_class(oid),
    name      TEXT NOT NULL DEFAULT '',
    kind      TEXT NOT NULL,
    columns   TEXT NOT NULL DEFAULT '' -- comma-separated attribute nums
);

-- sql_proc: functions, aggregates, window functions, procedures.
-- Modeled on pg_proc.
--   kind: 'f' = function, 'a' = aggregate, 'w' = window, 'p' = procedure
--   variadic_kind: 'n' = none, 'a' = array (VARIADIC any[]), 'v' = variadic-any
--   return_set: 1 if SETOF / table-returning
--   return_template: the result as an expression over the call's arguments,
--                    when it depends on their values: Decimal(18, $2) for
--                    toDecimal64(x, s). Empty for a result the type says.
CREATE TABLE sql_proc (
    oid             INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace_oid   INTEGER REFERENCES sql_namespace(oid),
    dialect_oid     INTEGER REFERENCES sql_dialect(oid),
    name            TEXT NOT NULL,
    kind            TEXT NOT NULL DEFAULT 'f',
    return_type_oid INTEGER NOT NULL REFERENCES sql_type(oid),
    return_set      INTEGER NOT NULL DEFAULT 0,
    return_nullable INTEGER NOT NULL DEFAULT 1,
    return_template TEXT NOT NULL DEFAULT '',
    strict          INTEGER NOT NULL DEFAULT 0,
    variadic_kind   TEXT NOT NULL DEFAULT 'n'
);

-- sql_proc_arg: ordered argument list for a proc.
--   mode: 'i' = in, 'o' = out, 'b' = both, 't' = table, 'v' = variadic
CREATE TABLE sql_proc_arg (
    proc_oid    INTEGER NOT NULL REFERENCES sql_proc(oid),
    ord         INTEGER NOT NULL,
    name        TEXT NOT NULL DEFAULT '',
    type_oid    INTEGER NOT NULL REFERENCES sql_type(oid),
    mode        TEXT NOT NULL DEFAULT 'i',
    has_default INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (proc_oid, ord)
);

-- sql_operator: operator overloads.
-- left_type_oid is NULL for prefix unary; right_type_oid is NULL for postfix.
CREATE TABLE sql_operator (
    oid             INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace_oid   INTEGER REFERENCES sql_namespace(oid),
    dialect_oid     INTEGER REFERENCES sql_dialect(oid),
    name            TEXT NOT NULL,
    left_type_oid   INTEGER REFERENCES sql_type(oid),
    right_type_oid  INTEGER REFERENCES sql_type(oid),
    result_type_oid INTEGER NOT NULL REFERENCES sql_type(oid),
    proc_oid        INTEGER REFERENCES sql_proc(oid),
    commutator_oid  INTEGER REFERENCES sql_operator(oid),
    negator_oid     INTEGER REFERENCES sql_operator(oid)
);

-- sql_cast: type coercion rules.
--   context: 'i' = implicit, 'a' = assignment-only, 'e' = explicit-only
--   proc_oid NULL = binary-coercible (no function needed)
CREATE TABLE sql_cast (
    source_type_oid INTEGER NOT NULL REFERENCES sql_type(oid),
    target_type_oid INTEGER NOT NULL REFERENCES sql_type(oid),
    proc_oid        INTEGER REFERENCES sql_proc(oid),
    context         TEXT NOT NULL DEFAULT 'e',
    dialect_oid     INTEGER REFERENCES sql_dialect(oid),
    PRIMARY KEY (source_type_oid, target_type_oid)
);

-- Resolution-speed indexes.
CREATE INDEX idx_sql_proc_name ON sql_proc(name, namespace_oid);
CREATE INDEX idx_sql_operator_name ON sql_operator(name, left_type_oid, right_type_oid);
CREATE INDEX idx_sql_attribute_name ON sql_attribute(name);
