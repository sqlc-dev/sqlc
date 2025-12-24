-- name: GetAuthor :one
SELECT id, name, bio FROM authors WHERE id = {id:UInt64};

-- name: ListAuthors :many
SELECT id, name, bio FROM authors ORDER BY name;

-- name: CreateAuthor :exec
INSERT INTO authors (id, name, bio) VALUES ({id:UInt64}, {name:String}, {bio:String});
