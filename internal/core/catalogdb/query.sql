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

-- name: LookupType :one
SELECT oid, name, category, typtype, preferred
FROM sql_type
WHERE oid = ?;

-- name: FindProcsAnyNamespace :many
SELECT oid, name, kind, return_type_oid, return_nullable
FROM sql_proc
WHERE name = ?;

-- name: FindProcsInNamespaces :many
SELECT oid, name, kind, return_type_oid, return_nullable
FROM sql_proc
WHERE name = sqlc.arg(name)
  AND namespace_oid IN (sqlc.slice(namespace_oids));
