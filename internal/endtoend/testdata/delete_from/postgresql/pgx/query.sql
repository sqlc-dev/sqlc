CREATE TABLE foo (id text not null);

-- name: DeleteFrom :exec
DELETE FROM foo WHERE id = $1;
