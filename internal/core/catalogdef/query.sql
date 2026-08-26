-- Queries against sqlc's own sql_* catalog tables, compiled by sqlc's
-- SQLite engine. Regenerate with `go generate ./internal/core/...`.

-- ============================ sql_namespace ============================

-- name: CreateNamespace :execlastid
INSERT INTO sql_namespace (name) VALUES (?);

-- name: NamespaceOID :one
SELECT oid FROM sql_namespace WHERE name = ?;

-- name: ListNamespaces :many
SELECT oid, name FROM sql_namespace ORDER BY oid;

-- ============================== sql_dialect ============================

-- name: CreateDialect :execlastid
INSERT INTO sql_dialect (name) VALUES (?);

-- name: DialectOID :one
SELECT oid FROM sql_dialect WHERE name = ?;

-- name: SeededDialect :one
-- The dialect a catalog was seeded with. A catalog is built for exactly one,
-- so a restored catalog finds it without being told its name.
SELECT oid, name FROM sql_dialect ORDER BY oid LIMIT 1;

-- name: SetDialectFlag :exec
INSERT INTO sql_dialect_flag (dialect_oid, key, value)
VALUES (?, ?, ?)
ON CONFLICT(dialect_oid, key) DO UPDATE SET value = excluded.value;

-- name: DialectFlag :one
SELECT value FROM sql_dialect_flag WHERE dialect_oid = ? AND key = ?;

-- =============================== sql_type ==============================

-- name: CreateType :execlastid
INSERT INTO sql_type
    (name, size, typtype, category, preferred, namespace_oid, dialect_oid, element_oid)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: TypeOIDByName :one
SELECT t.oid
FROM sql_type t
JOIN sql_namespace ns ON ns.oid = t.namespace_oid
WHERE t.name = sqlc.arg(name)
ORDER BY
    CASE ns.name
        WHEN 'pg_catalog' THEN 0
        WHEN 'public' THEN 1
        ELSE 2
    END,
    ns.name
LIMIT 1;

-- name: TypeNameByOID :one
SELECT name FROM sql_type WHERE oid = ?;

-- name: TypeOIDsInCategory :many
SELECT oid FROM sql_type
WHERE dialect_oid = ? AND category = ?
ORDER BY oid;

-- name: LookupType :one
SELECT oid, name, category, typtype, preferred
FROM sql_type
WHERE oid = ?;

-- =============================== sql_class =============================

-- name: CreateClass :execlastid
INSERT INTO sql_class (namespace_oid, name, kind) VALUES (?, ?, ?);

-- name: ClassOID :one
SELECT oid FROM sql_class WHERE namespace_oid = ? AND name = ?;

-- name: ClassOIDByName :one
SELECT oid FROM sql_class WHERE name = ? LIMIT 1;

-- name: ListTablesInNamespace :many
SELECT oid, name FROM sql_class
WHERE namespace_oid = ? AND kind = 'r'
ORDER BY oid;

-- name: DeleteClass :exec
DELETE FROM sql_class WHERE oid = ?;

-- name: RenameClass :exec
UPDATE sql_class SET name = sqlc.arg(new_name) WHERE oid = sqlc.arg(oid);

-- ============================= sql_attribute ===========================

-- name: CreateAttribute :exec
INSERT INTO sql_attribute (
    class_oid, name, type_oid, not_null, has_default, num,
    decl_type, type_length, type_scale,
    auto_increment, is_primary_key, is_unique, hidden
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: SetAttributePrimaryKey :exec
UPDATE sql_attribute SET is_primary_key = 1, not_null = 1
WHERE class_oid = ? AND name = ?;

-- name: SetAttributeUnique :exec
UPDATE sql_attribute SET is_unique = 1
WHERE class_oid = ? AND name = ?;

-- name: DeleteAttributesByClass :exec
DELETE FROM sql_attribute WHERE class_oid = ?;

-- name: DeleteAttribute :exec
DELETE FROM sql_attribute WHERE class_oid = ? AND name = ?;

-- name: RenameAttribute :exec
UPDATE sql_attribute SET name = sqlc.arg(new_name)
WHERE class_oid = sqlc.arg(class_oid) AND name = sqlc.arg(name);

-- name: SetAttributeType :exec
UPDATE sql_attribute SET type_oid = sqlc.arg(type_oid), decl_type = sqlc.arg(decl_type)
WHERE class_oid = sqlc.arg(class_oid) AND name = sqlc.arg(name);

-- name: SetAttributeNotNull :exec
UPDATE sql_attribute SET not_null = sqlc.arg(not_null)
WHERE class_oid = sqlc.arg(class_oid) AND name = sqlc.arg(name);

-- name: MaxAttributeNum :one
SELECT CAST(COALESCE(MAX(num), 0) AS INTEGER) AS num FROM sql_attribute WHERE class_oid = ?;

-- name: ResolveColumn :one
SELECT a.oid, a.class_oid, a.name, t.name AS type_name, a.type_oid, a.not_null,
       a.decl_type, a.type_length, a.type_scale,
       a.auto_increment, a.is_primary_key, a.is_unique
FROM sql_attribute a
JOIN sql_class c ON c.oid = a.class_oid
JOIN sql_type t ON t.oid = a.type_oid
WHERE c.name = sqlc.arg(table_name) AND a.name = sqlc.arg(column_name);

-- name: TableColumns :many
SELECT a.oid, a.class_oid, a.name, t.name AS type_name, a.type_oid, a.not_null,
       a.decl_type, a.type_length, a.type_scale,
       a.auto_increment, a.is_primary_key, a.is_unique
FROM sql_attribute a
JOIN sql_class c ON c.oid = a.class_oid
JOIN sql_type t ON t.oid = a.type_oid
WHERE c.name = ?
ORDER BY a.num;

-- name: ClassAttributes :many
SELECT oid, name, type_oid, not_null, hidden
FROM sql_attribute
WHERE class_oid = ?
ORDER BY num;

-- name: ListClassColumns :many
SELECT a.name AS column_name, t.name AS type_name, a.not_null
FROM sql_attribute a
JOIN sql_type t ON t.oid = a.type_oid
WHERE a.class_oid = ? AND a.hidden = 0
ORDER BY a.num;

-- name: LookupAttribute :one
SELECT ns.name AS schema_name, cls.name AS table_name, a.name AS column_name, a.num,
       a.decl_type, a.type_length, a.type_scale,
       a.auto_increment, a.is_primary_key, a.is_unique, a.not_null
FROM sql_attribute a
JOIN sql_class cls ON cls.oid = a.class_oid
JOIN sql_namespace ns ON ns.oid = cls.namespace_oid
WHERE a.oid = ?;

-- ============================ sql_constraint ===========================

-- name: CreateConstraint :exec
INSERT INTO sql_constraint (class_oid, name, kind, columns) VALUES (?, ?, ?, ?);

-- =============================== sql_proc ==============================

-- name: CreateProc :execlastid
INSERT INTO sql_proc
    (namespace_oid, dialect_oid, name, kind,
     return_type_oid, return_set, return_nullable, strict, variadic_kind)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateProcArg :exec
INSERT INTO sql_proc_arg (proc_oid, ord, name, type_oid, mode, has_default)
VALUES (?, ?, ?, ?, ?, ?);

-- name: ProcArgTypes :many
SELECT type_oid FROM sql_proc_arg
WHERE proc_oid = ? AND mode IN ('i', 'b', 'v')
ORDER BY ord;

-- name: FindProcsAnyNamespace :many
SELECT oid, name, kind, return_type_oid, return_nullable
FROM sql_proc
WHERE name = ?;

-- name: FindProcsInNamespaces :many
SELECT oid, name, kind, return_type_oid, return_nullable
FROM sql_proc
WHERE name = sqlc.arg(name)
  AND namespace_oid IN (sqlc.slice(namespace_oids));

-- ============================= sql_operator ============================

-- name: CreateOperator :execlastid
INSERT INTO sql_operator
    (namespace_oid, dialect_oid, name,
     left_type_oid, right_type_oid, result_type_oid, proc_oid)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: FindOperators :many
SELECT oid, name, left_type_oid, right_type_oid, result_type_oid, proc_oid
FROM sql_operator
WHERE name = sqlc.arg(name)
  AND (sqlc.arg(left_type_oid) = 0 OR left_type_oid = sqlc.arg(left_type_oid))
  AND (sqlc.arg(right_type_oid) = 0 OR right_type_oid = sqlc.arg(right_type_oid));

-- =============================== sql_cast ==============================

-- name: CreateCast :exec
INSERT INTO sql_cast (source_type_oid, target_type_oid, proc_oid, context, dialect_oid)
VALUES (?, ?, ?, ?, ?);

-- name: FindCast :one
SELECT source_type_oid, target_type_oid, proc_oid, context, dialect_oid
FROM sql_cast
WHERE source_type_oid = ? AND target_type_oid = ?;
