-- name: CreateAuthor :one
SELECT * FROM add_author (
  sqlc.arg(name), sqlc.arg(bio)
);
