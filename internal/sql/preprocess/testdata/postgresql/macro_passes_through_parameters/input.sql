-- tenant is not a parameter of the macro, so it is one of the query's.
-- name: Visible :macro
-- @param owner
sqlc.arg(owner) = sqlc.arg(tenant) OR sqlc.narg(admin)::bool IS TRUE;

-- name: ListDocs :many
SELECT * FROM docs d
WHERE sqlc.macro.Visible(d.owner_id)
  AND d.kind = ANY(sqlc.slice(kinds))
  AND d.tenant_id = @tenant;
