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

-- sql_type: data types. Modeled on pg_type.
--   typtype:  'b'ase | 'c'omposite | 'd'omain | 'e'num | 'p'seudo | 'r'ange
--   category: 'N'umeric | 'S'tring | 'B'oolean | 'D'atetime | 'A'rray |
--             'C'omposite | 'E'num | 'U'serdef | 'X'unknown
--   preferred: tie-breaker for implicit cast resolution within a category
--   element_oid: for arrays, points at the element type
--   dialect_oid: NULL = standard / shared across dialects
CREATE TABLE sql_type (
    oid           INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace_oid INTEGER NOT NULL REFERENCES sql_namespace(oid),
    dialect_oid   INTEGER REFERENCES sql_dialect(oid),
    name          TEXT NOT NULL,
    size          INTEGER NOT NULL DEFAULT 0,
    typtype       TEXT NOT NULL DEFAULT 'b',
    category      TEXT,
    preferred     INTEGER NOT NULL DEFAULT 0,
    element_oid   INTEGER REFERENCES sql_type(oid),
    UNIQUE (namespace_oid, name)
);
CREATE INDEX idx_sql_type_name ON sql_type(name);

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
--   decl_type:       original declared type string before normalization
--                    (e.g. VARCHAR(10), BIGINT UNSIGNED, INTEGER PRIMARY KEY).
--                    Useful for SQLite where multiple syntaxes collapse to
--                    one of five affinities, and as a debugging aid.
--   type_length:     length / precision (varchar(N), numeric(p,s).p,
--                    char(N), bit(N)). 0 = unspecified.
--   type_scale:      scale for numeric/decimal. 0 = unspecified.
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
    type_length    INTEGER NOT NULL DEFAULT 0,
    type_scale     INTEGER NOT NULL DEFAULT 0,
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
CREATE TABLE sql_proc (
    oid             INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace_oid   INTEGER REFERENCES sql_namespace(oid),
    dialect_oid     INTEGER REFERENCES sql_dialect(oid),
    name            TEXT NOT NULL,
    kind            TEXT NOT NULL DEFAULT 'f',
    return_type_oid INTEGER NOT NULL REFERENCES sql_type(oid),
    return_set      INTEGER NOT NULL DEFAULT 0,
    return_nullable INTEGER NOT NULL DEFAULT 1,
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
