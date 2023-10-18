-- name: UpdateAuthor :one
update authors
set names[$1] = $2
where id=$3
RETURNING *;
