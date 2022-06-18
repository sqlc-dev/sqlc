CREATE TABLE foo (name text not null, slug text not null);

-- name: UpdateSetMultiple :exec
UPDATE foo SET name = ?, slug = ?;
