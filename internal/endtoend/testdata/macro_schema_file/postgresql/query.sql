-- This file is the schema and the queries both. A macro is not SQL, so it is
-- expanded before the schema is read, just as it is before the queries are.
CREATE TABLE authors (
  id         bigserial PRIMARY KEY,
  org_id     bigint NOT NULL,
  name       text NOT NULL,
  deleted_at timestamptz
);

-- name: NotDeleted :macro
-- @param deleted_at
sqlc.arg(deleted_at) IS NULL;

-- name: ActiveAuthors :macro
-- @param org bigint
SELECT * FROM authors a WHERE a.org_id = sqlc.arg(org) AND sqlc.macro.NotDeleted(a.deleted_at);

-- name: ListActiveAuthors :many
SELECT a.id, a.name FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a ORDER BY a.name;
