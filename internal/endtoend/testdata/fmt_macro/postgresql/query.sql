-- A macro's own body is SQL, and is formatted like any other statement.
-- name: ActiveAuthors :macro
-- @param org bigint
select  *  from authors
where org_id = sqlc.arg(org) and deleted_at is null;

-- A statement that calls a macro is kept as it was written: the name of a
-- macro is matched exactly, and the call is not SQL until it is expanded.
-- name: ListActiveAuthors :many
select  a.id,a.name
from   sqlc.macro.ActiveAuthors( sqlc.arg(org_id) )  a
order by a.name;

-- name: GetAuthor :one
select  id,name
from   authors
where id =  $1;
